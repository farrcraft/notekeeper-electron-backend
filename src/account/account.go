package account

import (
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"../crypto"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Account is the database holding one or more users and their collection of notes
type Account struct {
	ID         uuid.UUID      `json:"id"`
	Name       string         `json:"name"`
	Users      []*UserProfile `json:"users"`
	ActiveUser *User          `json:"-"`           // ActiveUser is the currently active user of the account
	MasterDB   *bolt.DB       `json:"-"`           // MasterDB is the application-wide master database
	DB         *bolt.DB       `json:"-"`           // DB is the account-local database
	LicenseKey string         `json:"license_key"` // LicenseKey is the token which determines the available application features
	Shelves    []*Shelf       `json:"-"`
	Created    time.Time      `json:"created"`
	Updated    time.Time      `json:"updated"`
	Logger     *logrus.Logger `json:"-"`
}

// MapCount returns the number of records in the account_map table
func MapCount(db *bolt.DB) int {
	count := 0
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("account_map"))
		if bucket == nil {
			return nil
		}
		cursor := bucket.Cursor()
		for key, _ := cursor.First(); key != nil; key, _ = cursor.Next() {
			count++
		}
		return nil
	})
	return count
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

// OpenAccountDb opens the database file for a account
// The file is created if it doesn't already exist
func (account *Account) OpenAccountDb(dataPath string) error {
	dbFile := fmt.Sprint(account.ID.String(), ".db")
	fileName := filepath.Join(dataPath, dbFile)
	account.Logger.Info("Opening account db file [", fileName, "]")
	var err error
	account.DB, err = bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		account.Logger.Error("Unable to open account DB [", fileName, "] - ", err)
		return errors.New("unable to open account DB")
	}
	return nil
}

// NewAccount creates a new Account object
func NewAccount(db *bolt.DB, logger *logrus.Logger, name string) *Account {
	now := time.Now()
	account := &Account{
		ID:       uuid.NewV4(),
		Name:     name,
		MasterDB: db,
		Logger:   logger,
		Created:  now,
		Updated:  now,
	}
	return account
}

// Save saves the account to the database
func (account *Account) Save() error {
	// account data is stored in the master database
	err := account.MasterDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("accounts"))
		if err != nil {
			account.Logger.Error("Error creating accounts bucket - ", err)
			return err
		}
		data, err := json.Marshal(account)
		if err != nil {
			account.Logger.Error("Error marshaling account - ", err)
			return err
		}

		// account data must be encrypted with the account key and not the user key
		accountKey, err := crypto.Open(account.ActiveUser.PassphraseKey, account.ActiveUser.AccountKey)
		if err != nil {
			account.Logger.Error("Error opening account key - ", err)
			return err
		}
		encryptedData, err := crypto.Seal(accountKey, data)
		crypto.Zero(accountKey)
		if err != nil {
			account.Logger.Error("Error encrypting account content - ", err)
			return err
		}

		err = bucket.Put(account.ID.Bytes(), encryptedData)
		if err != nil {
			account.Logger.Error("Error saving account - ", err)
			return err
		}
		account.Logger.Debug("Saved account")
		return nil
	})
	if err != nil {
		account.Logger.Error("Error saving account - ", err)
		return err
	}
	// Since accounts are keyed by only an unencrypted id in the db
	// we also need to store a mapping between a key derived from the name and the id
	// otherwise there is no way to look up an account without taking a brute force decryption test approach
	err = account.MasterDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("account_map"))
		if err != nil {
			account.Logger.Error("Error creating account_map bucket - ", err)
			return err
		}
		encryptedName, err := crypto.DeriveSaltedKey([]byte(account.Name))
		if err != nil {
			account.Logger.Error("Error creating account map key - ", err)
			return err
		}
		err = bucket.Put(encryptedName, account.ID.Bytes())
		if err != nil {
			account.Logger.Error("Error saving account map - ", err)
			return err
		}
		return nil
	})
	if err != nil {
		account.Logger.Error("Error mapping account - ", err)
		return err
	}
	return err
}

// Lookup searches the database for a matching account name and loads it from the db
func (account *Account) Lookup() error {
	originalID := account.ID
	account.ID = uuid.Nil
	// we don't expect many accounts to exist in the db (typically just one), so we iterate through them all
	account.MasterDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("account_map"))
		// If bucket doesn't exist, no accounts have been created yet
		if bucket == nil {
			account.Logger.Debug("account map bucket does not exist")
			return nil
		}
		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			// extract the salt from the existing encrypted name
			salt, encryptedName := crypto.ExtractSalt(key)
			// create a new key using the extracted salt and the unencrypted name we're searching for
			checkName, err := crypto.DeriveKey([]byte(account.Name), salt[:])
			if err != nil {
				account.Logger.Error("Error deriving account map key - ", err)
				return err
			}
			// the new key should match the existing key if we have the right name and salt
			if subtle.ConstantTimeCompare(encryptedName[:], checkName[:]) == 1 {
				account.ID, err = uuid.FromBytes(value)
				if err != nil {
					account.Logger.Error("Error converting account map uuid - ", err)
					return err
				}
				return nil
			}
		}
		return nil
	})
	if account.ID == uuid.Nil {
		account.ID = originalID
		account.Logger.Debug("account lookup - no account found for [", account.Name, "]")
		return errors.New("no account found")
	}
	return nil
}

// Load loads an account from the database
func (account *Account) Load() error {
	err := account.MasterDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("accounts"))
		if bucket == nil {
			err := errors.New("account bucket does not exist")
			account.Logger.Error(err)
			return err
		}
		cursor := bucket.Cursor()
		key, value := cursor.Seek(account.ID.Bytes())
		if key == nil {
			err := errors.New("Error loading account")
			account.Logger.Error(err)
			return err
		}

		// account data is encrypted with the account key and not the user key
		account.Logger.Debug("account key [", account.ActiveUser.AccountKey, "] passphrase key [", account.ActiveUser.PassphraseKey, "]")
		accountKey, err := crypto.Open(account.ActiveUser.PassphraseKey, account.ActiveUser.AccountKey)
		if err != nil {
			account.Logger.Error("Error opening account key - ", err)
			return err
		}

		// decrypt value
		decryptedData, err := crypto.Open(accountKey, value)
		if err != nil {
			account.Logger.Error("Error decrypting account data - ", err)
			return err
		}

		err = json.Unmarshal(decryptedData, account)
		if err != nil {
			account.Logger.Error("Error decoding account json - ", err)
			return err
		}
		return nil
	})
	return err
}
