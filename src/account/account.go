package account

import (
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../shelf"
	"../user"
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

// Account is the database holding one or more users and their collection of notes
type Account struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Users        []*user.Profile `json:"users"`
	ActiveUser   *user.User      `json:"-"`           // ActiveUser is the currently active user of the account
	DBRegistry   *db.Registry    `json:"-"`           // DBRegistry provides access to databases
	LicenseKey   string          `json:"license_key"` // LicenseKey is the token which determines the available application features
	EncryptedKey []byte          `json:"-"`           // EncryptedKey is the encrypted encryption key for the account DB
	Shelves      []*shelf.Shelf  `json:"-"`
	Created      time.Time       `json:"created"`
	Updated      time.Time       `json:"updated"`
	Logger       *logrus.Logger  `json:"-"`
}

// IsLocked returns whether the account is in a locked state or not
func (account *Account) IsLocked() bool {
	if account.ActiveUser == nil {
		return true
	}
	if len(account.ActiveUser.PassphraseKey) == 0 {
		return true
	}
	return false
}

// New creates a new Account object
func New(dbRegistry *db.Registry, logger *logrus.Logger, name string) (*Account, error) {
	now := time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	account := &Account{
		ID:         id,
		Name:       name,
		DBRegistry: dbRegistry,
		Logger:     logger,
		Created:    now,
		Updated:    now,
	}

	return account, nil
}

// CreateEncryptedKey generates a new encryption key that is encrypted with the account key
// The account key must already exist
func (account *Account) CreateEncryptedKey() ([]byte, error) {
	// create the new key
	var encryptedKey []byte
	c := crypto.New(account.Logger)
	newKey, err := c.GenerateKey()
	if err != nil {
		return encryptedKey, err
	}

	// we need to open up the account key
	// our copy of the account key has been sealed with the active user's passphrase key
	accountKey, err := c.Open(account.ActiveUser.PassphraseKey, account.EncryptedKey)
	if err != nil {
		account.Logger.Warn("Error opening account key while creating key - ", err)
		code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
		return encryptedKey, code
	}

	// now we can seal our newly created key
	encryptedKey, err = c.Seal(accountKey, newKey[:])
	if err != nil {
		return encryptedKey, err
	}

	return encryptedKey, nil
}

// UnsealKey unseals a key sealed with one of either the user, account or passphrase key
// The EncryptionKeyType denotes which key sealedKey was sealed with.
func (account *Account) UnsealKey(keyType EncryptionKeyType, sealedKey []byte) ([]byte, error) {
	var emptyKey []byte
	var unsealKey []byte
	c := crypto.New(account.Logger)
	if keyType == TypeAccount {
		accountKey, err := c.Open(account.ActiveUser.PassphraseKey, account.EncryptedKey)
		if err != nil {
			account.Logger.Warn("Error unsealing account key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
			return emptyKey, code
		}
		unsealKey = accountKey
	} else if keyType == TypeUser {
		userKey, err := c.Open(account.ActiveUser.PassphraseKey, account.ActiveUser.UserKey)
		if err != nil {
			account.Logger.Warn("Error opening user key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
			return emptyKey, code
		}
		unsealKey = userKey
	} else if keyType == TypePassphrase {
		unsealKey = account.ActiveUser.PassphraseKey
	}

	unsealedKey, err := c.Open(unsealKey, sealedKey)
	if err != nil {
		account.Logger.Warn("Error opening key - ", err)
		code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
		return emptyKey, code
	}
	return unsealedKey, nil
}

// Save saves the account to the database
func (account *Account) Save() error {
	accountDBKey := db.Key{
		ID:   account.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := account.DBRegistry.GetHandle(accountDBKey)
	if err != nil {
		return err
	}
	err = accountDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("profile"))
		if err != nil {
			account.Logger.Warn("Error creating accounts bucket - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorCreateBucket)
			return code
		}
		data, err := json.Marshal(account)
		if err != nil {
			account.Logger.Warn("Error marshaling account - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorMarshal)
			return code
		}

		// account data must be encrypted with the account key and not the user key
		c := crypto.New(account.Logger)
		accountKey, err := c.Open(account.ActiveUser.PassphraseKey, accountDBHandle.EncryptedKey)
		if err != nil {
			account.Logger.Warn("Error opening account key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
			return code
		}
		encryptedData, err := c.Seal(accountKey, data)
		crypto.Zero(accountKey)
		if err != nil {
			account.Logger.Warn("Error encrypting account content - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorEncrypt)
			return code
		}

		err = bucket.Put(account.ID.Bytes(), encryptedData)
		if err != nil {
			account.Logger.Warn("Error writing to account bucket - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorWriteBucket)
			return code
		}
		account.Logger.Debug("Saved account")
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		account.Logger.Warn("Error saving account - ", err)
		code := codes.New(codes.ScopeAccount, codes.ErrorSave)
		return code
	}
	return nil
}

// Load account profile data from the database
func (account *Account) Load() error {
	key := db.Key{
		ID:   account.ID,
		Type: db.TypeAccount,
	}
	handle, err := account.DBRegistry.GetHandle(key)
	if err != nil {
		return err
	}
	err = handle.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("profile"))
		if bucket == nil {
			account.Logger.Warn("account bucket does not exist")
			code := codes.New(codes.ScopeAccount, codes.ErrorBucketMissing)
			return code
		}
		cursor := bucket.Cursor()
		key, value := cursor.Seek(account.ID.Bytes())
		if key == nil {
			account.Logger.Warn("Error loading account")
			code := codes.New(codes.ScopeAccount, codes.ErrorLoad)
			return code
		}

		// account data is encrypted with the account key and not the user key
		c := crypto.New(account.Logger)
		if account.ActiveUser == nil {
			account.Logger.Warn("Error missing active user")
			code := codes.New(codes.ScopeAccount, codes.ErrorMissingUser)
			return code
		}
		// EncryptedKey is the same as account.ActiveUser.AccountKey
		accountKey, err := c.Open(account.ActiveUser.PassphraseKey, handle.EncryptedKey)
		if err != nil {
			account.Logger.Warn("Error opening account key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
			return code
		}

		// decrypt value
		decryptedData, err := c.Open(accountKey, value)
		if err != nil {
			account.Logger.Warn("Error decrypting account data - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorDecrypt)
			return code
		}

		err = json.Unmarshal(decryptedData, account)
		if err != nil {
			account.Logger.Warn("Error decoding account json - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorDecode)
			return code
		}
		handle.EncryptedKey = account.ActiveUser.AccountKey
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		account.Logger.Warn("Error loading account - ", err)
		code := codes.New(codes.ScopeAccount, codes.ErrorLoad)
		return code
	}

	return nil
}
