package account

import (
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../shelf"
	"../user"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Account is the database holding one or more users and their collection of notes
type Account struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Users        []*user.Profile `json:"users"`
	ActiveUser   *user.User      `json:"-"`              // ActiveUser is the currently active user of the account
	DBRegistry   *db.Registry    `json:"-"`              // DBRegistry provides access to databases
	LicenseKey   string          `json:"license_key"`    // LicenseKey is the token which determines the available application features
	EncryptedKey []byte          `json:"encryption_key"` // EncryptedKey is the encrypted encryption key for the account DB
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
func New(dbRegistry *db.Registry, logger *logrus.Logger, name string) *Account {
	now := time.Now()
	account := &Account{
		ID:         uuid.NewV4(),
		Name:       name,
		DBRegistry: dbRegistry,
		Logger:     logger,
		Created:    now,
		Updated:    now,
	}
	return account
}

// Save saves the account to the database
func (account *Account) Save() error {
	// account data is stored in the master database
	accountDBKey := db.Key{
		ID:   account.ID,
		Type: db.TypeAccount,
	}
	accountDBHandle, err := account.DBRegistry.GetHandle(accountDBKey)
	if err != nil {
		return err
	}
	err = accountDBHandle.DB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("profile"))
		if err != nil {
			account.Logger.Debug("Error creating accounts bucket - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorCreateBucket)
			return code
		}
		data, err := json.Marshal(account)
		if err != nil {
			account.Logger.Debug("Error marshaling account - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorMarshal)
			return code
		}

		// account data must be encrypted with the account key and not the user key
		c := crypto.New(account.Logger)
		accountKey, err := c.Open(account.ActiveUser.PassphraseKey, accountDBHandle.EncryptedKey)
		if err != nil {
			account.Logger.Debug("Error opening account key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
			return code
		}
		encryptedData, err := c.Seal(accountKey, data)
		crypto.Zero(accountKey)
		if err != nil {
			account.Logger.Debug("Error encrypting account content - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorEncrypt)
			return code
		}

		err = bucket.Put(account.ID.Bytes(), encryptedData)
		if err != nil {
			account.Logger.Debug("Error writing to account bucket - ", err)
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
		account.Logger.Debug("Error saving account - ", err)
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
	err = handle.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("profile"))
		if bucket == nil {
			account.Logger.Debug("account bucket does not exist")
			code := codes.New(codes.ScopeAccount, codes.ErrorBucketMissing)
			return code
		}
		cursor := bucket.Cursor()
		key, value := cursor.Seek(account.ID.Bytes())
		if key == nil {
			account.Logger.Debug("Error loading account")
			code := codes.New(codes.ScopeAccount, codes.ErrorLoad)
			return code
		}

		// account data is encrypted with the account key and not the user key
		c := crypto.New(account.Logger)
		// EncryptedKey is the same as account.ActiveUser.AccountKey
		accountKey, err := c.Open(account.ActiveUser.PassphraseKey, handle.EncryptedKey)
		if err != nil {
			account.Logger.Debug("Error opening account key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorOpenKey)
			return code
		}

		// decrypt value
		decryptedData, err := c.Open(accountKey, value)
		if err != nil {
			account.Logger.Debug("Error decrypting account data - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorDecrypt)
			return code
		}

		err = json.Unmarshal(decryptedData, account)
		if err != nil {
			account.Logger.Debug("Error decoding account json - ", err)
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
		account.Logger.Debug("Error loading account - ", err)
		code := codes.New(codes.ScopeAccount, codes.ErrorLoad)
		return code
	}

	return nil
}
