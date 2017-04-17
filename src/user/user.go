package user

import (
	"crypto/subtle"
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../shelf"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Profile contains the minimal user information that is visible to all users of an account
type Profile struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// User is a single user in an account
type User struct {
	ID            uuid.UUID      `json:"id"`          // ID is the unique identifier of the user
	AccountID     uuid.UUID      `json:"account_id"`  // ID is the unique identifier of the account
	Profile       *Profile       `json:"profile"`     // Profile is the user information that is visible to all users in an account
	Active        bool           `json:"-"`           // Active indicates whether the user is active or not
	Created       time.Time      `json:"created"`     // Created is the time when the user was created
	Updated       time.Time      `json:"updated"`     // Updated is the time when the user was last created
	AccountKey    []byte         `json:"account_key"` // AccountKey is the encrypted version of the account-level encryption key
	UserKey       []byte         `json:"user_key"`    // UserKey is the encrypted version of the user-level encryption key
	PassphraseKey []byte         `json:"-"`           // PassphraseKey is the key derived from the passphrase
	Salt          []byte         `json:"-"`           // Salt is the unique salt for generating the passphrase key
	Shelves       []*shelf.Shelf `json:"-"`           // Shelves is the set of shelves that belong to the user
	Logger        *logrus.Logger `json:"-"`           // Logger is a log instance
	DBFactory     *db.Factory    `json:"-"`           // DBFactory provides access to dbs
}

// New creates a new user object
func New(dbFactory *db.Factory, logger *logrus.Logger, accountID uuid.UUID, email string) *User {
	now := time.Now()
	user := &User{
		ID:        uuid.NewV4(),
		AccountID: accountID,
		Profile: &Profile{
			Email: email,
		},
		Active:    true,
		Created:   now,
		Updated:   now,
		DBFactory: dbFactory,
		Logger:    logger,
	}
	return user
}

// Lookup a user id in the account database
func (user *User) Lookup() error {
	originalID := user.ID
	user.ID = uuid.Nil
	db := user.DBFactory.Find(db.TypeAccount, user.AccountID)
	err := db.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("user_map"))
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			// extract the salt from the existing encrypted email
			salt, encryptedEmail := crypto.ExtractSalt(key)
			// create a new key using the extracted salt and the unencrypted email we're searching for
			checkEmail, err := crypto.DeriveKey([]byte(user.Profile.Email), salt[:])
			if err != nil {
				user.Logger.Debug("Error deriving user map key - ", err)
				code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
				return code
			}
			// the new key should match the existing key if we have the right email and salt
			if subtle.ConstantTimeCompare(encryptedEmail[:], checkEmail[:]) == 1 {
				user.ID, err = uuid.FromBytes(value)
				if err != nil {
					user.Logger.Debug("Error converting email map uuid - ", err)
					code := codes.New(codes.ScopeUser, codes.ErrorConvertID)
					return code
				}
				// the salt stored in the email key is also the primary user passphrase key salt
				user.Salt = salt[:]
				return nil
			}
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		// [FIXME] - handle unknown error
	}
	if user.ID == uuid.Nil {
		user.ID = originalID
		code := codes.New(codes.ScopeUser, codes.ErrorLookup)
		return code
	}
	return nil
}

// Load the user data for a user from the account database
func (user *User) Load(passphrase string) error {
	db := user.DBFactory.Find(db.TypeUser, user.ID)
	err := db.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("user"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek(user.ID.Bytes())
		if key == nil {
			user.Logger.Debug("Error loading user")
			code := codes.New(codes.ScopeUser, codes.ErrorLoad)
			return code
		}

		// salt will have already been provided from a previous Lookup() operation
		passphraseKey, err := crypto.DeriveKey([]byte(passphrase), user.Salt)
		if err != nil {
			user.Logger.Debug("Error deriving key from passphrase - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
			return code
		}
		// need to decrypt value
		// the user data is a special case - it is sealed with the passphrase key instead
		// of the user db key - that key won't be available until the user data is unsealed.
		decryptedData, err := crypto.Open(passphraseKey[:], value)
		if err != nil {
			// this error condition is most likely caused by an incorrect passphrase
			user.Logger.Debug("Error decrypting user data - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorDecrypt)
			return code
		}

		err = json.Unmarshal(decryptedData, user)
		if err != nil {
			user.Logger.Debug("Error decoding user json - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorDecode)
			return code
		}
		user.PassphraseKey = passphraseKey[:]
		db.EncryptedKey = user.UserKey
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		// [FIXME] handle unknown error
	}
	return nil
}

// Save saves the user to the database
func (user *User) Save() error {
	userDB := user.DBFactory.DB(db.TypeUser, user.ID)
	err := userDB.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("user"))
		if err != nil {
			user.Logger.Debug("Error creating users bucket - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorCreateBucket)
			return code
		}
		data, err := json.Marshal(user)
		if err != nil {
			user.Logger.Debug("Error marshaling user - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorMarshal)
			return code
		}
		encryptedData, err := crypto.Seal(user.PassphraseKey, data)
		if err != nil {
			user.Logger.Debug("Error encrypting user content - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorEncrypt)
			return code
		}

		err = bucket.Put(user.ID.Bytes(), encryptedData)
		if err != nil {
			user.Logger.Debug("Error writing user - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		user.Logger.Debug("Error saving user - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorSave)
		return code
	}
	// Since users are keyed by only an unencrypted id in the db
	// we also need to store a mapping between a key derived from the email address and the id
	// otherwise there is no way to look up a user without taking a brute force decryption test approach
	accountDB := user.DBFactory.Find(db.TypeAccount, user.AccountID)
	err = accountDB.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("user_map"))
		if err != nil {
			user.Logger.Debug("Error creating user_map bucket - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorCreateBucket)
			return code
		}
		encryptedEmail, err := crypto.DeriveKey([]byte(user.Profile.Email), user.Salt)
		if err != nil {
			user.Logger.Debug("Error creating user map key - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
			return code
		}
		saltedKey := crypto.EmbedSalt(encryptedEmail, user.Salt)
		err = bucket.Put(saltedKey, user.ID.Bytes())
		if err != nil {
			user.Logger.Debug("Error saving user map - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		user.Logger.Debug("Error mapping user - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorSave)
		return code
	}
	return nil
}

// CreateKeys creates the account and user key from a passphrase
func (user *User) CreateKeys(passphrase []byte) error {
	// generate account-level encryption key
	accountKey, err := crypto.GenerateKey()
	if err != nil {
		// [FIXME] - internal error?
		return err
	}

	// we already have a user key in the form of the passphrase key
	// so having a separate user key is a bit redundant, but it does
	// make things consistent and easier overall in the long run.
	userKey, err := crypto.GenerateKey()
	if err != nil {
		// [FIXME] - internal error?
		return err
	}

	// derive key from passphrase
	var key = new([crypto.KeySize]byte)
	key, user.Salt, err = crypto.DeriveKeyAndSalt(passphrase)
	if err != nil {
		// [FIXME] - internal error?
		return err
	}

	// slicedKey := key[:]
	//user.PassphraseKey = append(user.Salt, slicedKey...)
	user.PassphraseKey = key[:]

	user.AccountKey, err = crypto.Seal(user.PassphraseKey, accountKey[:])
	if err != nil {
		// [FIXME] - internal error?
		return err
	}

	user.UserKey, err = crypto.Seal(user.PassphraseKey, userKey[:])
	if err != nil {
		// [FIXME] - internal error?
		return err
	}

	return nil
}
