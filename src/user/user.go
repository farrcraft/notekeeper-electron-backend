package user

import (
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
	ID            uuid.UUID      `json:"id"`             // ID is the unique identifier of the user
	AccountID     uuid.UUID      `json:"account_id"`     // ID is the unique identifier of the account
	Profile       *Profile       `json:"profile"`        // Profile is the user information that is visible to all users in an account
	Active        bool           `json:"-"`              // Active indicates whether the user is active or not
	Created       time.Time      `json:"created"`        // Created is the time when the user was created
	Updated       time.Time      `json:"updated"`        // Updated is the time when the user was last created
	AccountKey    []byte         `json:"account_key"`    // AccountKey is the encrypted version of the account-level encryption key
	UserKey       []byte         `json:"encryption_key"` // UserKey is the encrypted version of the user-level encryption key
	PassphraseKey []byte         `json:"-"`              // PassphraseKey is the key derived from the passphrase
	Salt          []byte         `json:"-"`              // Salt is the unique salt for generating the passphrase key
	Shelves       []*shelf.Shelf `json:"-"`              // Shelves is the set of shelves that belong to the user
	Logger        *logrus.Logger `json:"-"`              // Logger is a log instance
	DBRegistry    *db.Registry   `json:"-"`              // DBRegistry provides access to dbs
}

// New creates a new user object
func New(dbRegistry *db.Registry, logger *logrus.Logger, accountID uuid.UUID, email string) *User {
	now := time.Now()
	user := &User{
		ID:        uuid.NewV4(),
		AccountID: accountID,
		Profile: &Profile{
			Email: email,
		},
		Active:     true,
		Created:    now,
		Updated:    now,
		DBRegistry: dbRegistry,
		Logger:     logger,
	}
	return user
}

// Load the user data for a user from the account database
func (user *User) Load(passphrase string) error {
	// salt will have already been provided from a previous Lookup() operation
	c := crypto.New(user.Logger)
	passphraseKey, err := c.DeriveKey([]byte(passphrase), user.Salt)
	if err != nil {
		user.Logger.Debug("Error deriving key from passphrase - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
		return code
	}

	userKey := db.Key{Type: db.TypeUser, ID: user.ID}
	userDBHandle, err := user.DBRegistry.GetHandle(userKey, passphraseKey[:])
	if err != nil {
		return err
	}
	err = userDBHandle.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("user"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek(user.ID.Bytes())
		if key == nil {
			user.Logger.Debug("Error loading user")
			code := codes.New(codes.ScopeUser, codes.ErrorLoad)
			return code
		}

		// need to decrypt value
		// the user data is a special case - it is sealed with the passphrase key instead
		// of the user db key - that key won't be available until the user data is unsealed.
		decryptedData, err := c.Open(passphraseKey[:], value)
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
		userDBHandle.EncryptedKey = user.UserKey
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
	userKey := db.Key{Type: db.TypeUser, ID: user.ID}
	userDBHandle, err := user.DBRegistry.GetHandle(userKey, user.PassphraseKey)
	if err != nil {
		return err
	}
	c := crypto.New(user.Logger)
	err = userDBHandle.DB.Update(func(tx *bolt.Tx) error {
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
		encryptedData, err := c.Seal(user.PassphraseKey, data)
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
	return nil
}

// CreateKeys creates the account and user key from a passphrase
func (user *User) CreateKeys(passphrase []byte) error {
	// generate account-level encryption key
	c := crypto.New(user.Logger)
	accountKey, err := c.GenerateKey()
	if err != nil {
		return err
	}

	// we already have a user key in the form of the passphrase key
	// so having a separate user key is a bit redundant, but it does
	// make things consistent and easier overall in the long run.
	userKey, err := c.GenerateKey()
	if err != nil {
		return err
	}

	// derive key from passphrase
	var key = new([crypto.KeySize]byte)
	key, user.Salt, err = c.DeriveKeyAndSalt(passphrase)
	if err != nil {
		return err
	}

	// slicedKey := key[:]
	//user.PassphraseKey = append(user.Salt, slicedKey...)
	user.PassphraseKey = key[:]

	user.AccountKey, err = c.Seal(user.PassphraseKey, accountKey[:])
	if err != nil {
		return err
	}

	user.UserKey, err = c.Seal(user.PassphraseKey, userKey[:])
	if err != nil {
		return err
	}

	return nil
}
