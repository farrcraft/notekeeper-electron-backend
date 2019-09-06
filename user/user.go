package user

import (
	"encoding/json"
	"time"

	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/crypto"
	"notekeeper-electron-backend/db"
	"notekeeper-electron-backend/shelf"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

// EncryptionKeyType indicates the type of encryption key
type EncryptionKeyType int

// Encryption Key Types
const (
	TypeAccount EncryptionKeyType = iota
	TypeUser
	TypePassphrase
)

// User is a single user in an account
type User struct {
	ID            uuid.UUID      `json:"id"`             // ID is the unique identifier of the user
	AccountID     uuid.UUID      `json:"account_id"`     // ID is the unique identifier of the account
	Profile       *Profile       `json:"profile"`        // Profile is the user information that is visible to all users in an account
	Active        bool           `json:"-"`              // Active indicates whether the user is active or not
	Created       time.Time      `json:"created"`        // Created is the time when the user was created
	Updated       time.Time      `json:"updated"`        // Updated is the time when the user was last created
	AccountKey    []byte         `json:"account_key"`    // AccountKey account-level encryption key encrypted with the passphrase key
	UserKey       []byte         `json:"encryption_key"` // UserKey is the user-level encryption key encrypted with the passphrase key
	PassphraseKey []byte         `json:"-"`              // PassphraseKey is the key derived from the passphrase
	Salt          []byte         `json:"-"`              // Salt is the unique salt for generating the passphrase key
	Shelves       []*shelf.Shelf `json:"-"`              // Shelves is the set of shelves that belong to the user
	Logger        *logrus.Logger `json:"-"`              // Logger is a log instance
	DBRegistry    *db.Registry   `json:"-"`              // DBRegistry provides access to dbs
}

// New creates a new user object
func New(dbRegistry *db.Registry, logger *logrus.Logger, accountID uuid.UUID, email string) (*User, error) {
	now := time.Now()
	id := uuid.NewV4()

	user := &User{
		ID:        id,
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

	return user, nil
}

// Load the user data for a user from the account database
func (user *User) Load(passphrase string) error {
	// salt will have already been provided from a previous Lookup() operation
	c := crypto.New(user.Logger)
	passphraseKey, err := c.DeriveKey([]byte(passphrase), user.Salt)
	if err != nil {
		user.Logger.Warn("Error deriving key from passphrase - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
		return code
	}

	userKey := db.Key{
		Type: db.TypeUser,
		ID:   user.ID,
	}
	user.Logger.Debug("Getting DB handle for user ID - ", userKey.ID)
	userDBHandle, err := user.DBRegistry.GetHandle(userKey)
	if err != nil {
		return err
	}
	err = userDBHandle.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("profile"))
		cursor := bucket.Cursor()
		key, value := cursor.Seek(user.ID.Bytes())
		if key == nil {
			user.Logger.Warn("Error loading user")
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
			code := codes.New(codes.ScopeUser, codes.ErrorUnauthorized)
			return code
		}

		err = json.Unmarshal(decryptedData, user)
		if err != nil {
			user.Logger.Warn("Error decoding user json - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorDecode)
			return code
		}
		user.PassphraseKey = passphraseKey[:]
		userDBHandle.EncryptedKey = user.UserKey

		// update the account DB handle with the account-level encryption key
		accountDBKey := db.Key{
			Type: db.TypeAccount,
			ID:   user.AccountID,
		}
		accountDBHandle, err := user.DBRegistry.GetHandle(accountDBKey)
		if err != nil {
			return err
		}
		accountDBHandle.EncryptedKey = user.AccountKey

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
	userDBHandle, err := user.DBRegistry.GetHandle(userKey)
	if err != nil {
		return err
	}
	c := crypto.New(user.Logger)
	err = userDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("profile"))
		if err != nil {
			user.Logger.Warn("Error creating users bucket - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorCreateBucket)
			return code
		}
		data, err := json.Marshal(user)
		if err != nil {
			user.Logger.Warn("Error marshaling user - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorMarshal)
			return code
		}
		encryptedData, err := c.Seal(user.PassphraseKey, data)
		if err != nil {
			user.Logger.Warn("Error encrypting user content - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorEncrypt)
			return code
		}

		err = bucket.Put(user.ID.Bytes(), encryptedData)
		if err != nil {
			user.Logger.Warn("Error writing user - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		user.Logger.Warn("Error saving user - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorSave)
		return code
	}
	return nil
}

// UnsealKey unseals a key sealed with one of either the user or passphrase key
// The EncryptionKeyType denotes which key sealedKey was sealed with.
func (user *User) UnsealKey(keyType EncryptionKeyType, sealedKey []byte) ([]byte, error) {
	var emptyKey []byte
	var unsealKey []byte
	c := crypto.New(user.Logger)
	if keyType == TypeUser {
		userKey, err := c.Open(user.PassphraseKey, user.UserKey)
		if err != nil {
			user.Logger.Warn("Error opening user key - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorOpenKey)
			return emptyKey, code
		}
		unsealKey = userKey
	} else if keyType == TypePassphrase {
		unsealKey = user.PassphraseKey
	} else {
		user.Logger.Warn("Cannot unseal key of unknown key type.")
		code := codes.New(codes.ScopeUser, codes.ErrorOpenKey)
		return emptyKey, code
	}

	unsealedKey, err := c.Open(unsealKey, sealedKey)
	if err != nil {
		user.Logger.Warn("Error opening key - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorOpenKey)
		return emptyKey, code
	}
	return unsealedKey, nil
}

// CreateEncryptedKey generates a new encryption key that is encrypted with the user's passphrase key
// The user must already have their own encryption key
// EncryptionKeyType denotes which key is used to seal the generated key
func (user *User) CreateEncryptedKey(keyType EncryptionKeyType) ([]byte, error) {
	var encryptedKey []byte
	var sealingKey []byte
	c := crypto.New(user.Logger)
	newKey, err := c.GenerateKey()
	if err != nil {
		return encryptedKey, err
	}

	if keyType == TypeUser {
		userKey, err := c.Open(user.PassphraseKey, user.UserKey)
		if err != nil {
			user.Logger.Warn("Error opening user key - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorOpenKey)
			return encryptedKey, code
		}
		sealingKey = userKey
	} else if keyType == TypePassphrase {
		sealingKey = user.PassphraseKey
	} else {
		user.Logger.Warn("Cannot unseal key of unknown key type.")
		code := codes.New(codes.ScopeUser, codes.ErrorOpenKey)
		return encryptedKey, code
	}

	encryptedKey, err = c.Seal(sealingKey, newKey[:])
	if err != nil {
		return encryptedKey, err
	}

	return encryptedKey, nil
}

// CreateUserKey generates a user-specific encryption key
func (user *User) CreateUserKey(passphrase []byte) error {
	c := crypto.New(user.Logger)

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

	user.UserKey, err = c.Seal(user.PassphraseKey, userKey[:])
	if err != nil {
		return err
	}

	return nil
}

// CreateKeys for both account and user.
// Both keys are encrypted with a key derived from the passphrase.
func (user *User) CreateKeys(passphrase []byte) error {
	err := user.CreateUserKey(passphrase)
	if err != nil {
		return err
	}

	// generate an account-wide encryption key
	user.AccountKey, err = user.CreateEncryptedKey(TypePassphrase)
	if err != nil {
		return err
	}

	return nil
}
