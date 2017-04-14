package account

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"../codes"
	"../crypto"
	"../shelf"
	"../user"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Account is the database holding one or more users and their collection of notes
type Account struct {
	ID         uuid.UUID       `json:"id"`
	Name       string          `json:"name"`
	Users      []*user.Profile `json:"users"`
	ActiveUser *user.User      `json:"-"`           // ActiveUser is the currently active user of the account
	MasterDB   *bolt.DB        `json:"-"`           // MasterDB is the application-wide master database
	DB         *bolt.DB        `json:"-"`           // DB is the account-local database
	LicenseKey string          `json:"license_key"` // LicenseKey is the token which determines the available application features
	Shelves    []*shelf.Shelf  `json:"-"`
	Created    time.Time       `json:"created"`
	Updated    time.Time       `json:"updated"`
	Logger     *logrus.Logger  `json:"-"`
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
		account.Logger.Debug("Unable to open account DB [", fileName, "] - ", err)
		code := codes.New(codes.ErrorAccountDbOpen)
		return code
	}
	return nil
}

// New creates a new Account object
func New(db *bolt.DB, logger *logrus.Logger, name string) *Account {
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
			account.Logger.Debug("Error creating accounts bucket - ", err)
			code := codes.New(codes.ErrorAccountBucket)
			return code
		}
		data, err := json.Marshal(account)
		if err != nil {
			account.Logger.Debug("Error marshaling account - ", err)
			code := codes.New(codes.ErrorAccountMarshal)
			return code
		}

		// account data must be encrypted with the account key and not the user key
		accountKey, err := crypto.Open(account.ActiveUser.PassphraseKey, account.ActiveUser.AccountKey)
		if err != nil {
			account.Logger.Debug("Error opening account key - ", err)
			code := codes.New(codes.ErrorAccountKey)
			return code
		}
		encryptedData, err := crypto.Seal(accountKey, data)
		crypto.Zero(accountKey)
		if err != nil {
			account.Logger.Debug("Error encrypting account content - ", err)
			code := codes.New(codes.ErrorAccountEncrypt)
			return code
		}

		err = bucket.Put(account.ID.Bytes(), encryptedData)
		if err != nil {
			account.Logger.Debug("Error writing to account bucket - ", err)
			code := codes.New(codes.ErrorAccountWrite)
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
		code := codes.New(codes.ErrorAccountSave)
		return code
	}

	// Since accounts are keyed by only an unencrypted id in the db
	// we also need to store a mapping between a key derived from the name and the id
	// otherwise there is no way to look up an account without taking a brute force decryption test approach
	err = account.MasterDB.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("account_map"))
		if err != nil {
			account.Logger.Debug("Error creating account_map bucket - ", err)
			code := codes.New(codes.ErrorAccountMapBucket)
			return code
		}
		encryptedName, err := crypto.DeriveSaltedKey([]byte(account.Name))
		if err != nil {
			account.Logger.Debug("Error creating account map key - ", err)
			code := codes.New(codes.ErrorAccountMapKey)
			return code
		}
		err = bucket.Put(encryptedName, account.ID.Bytes())
		if err != nil {
			account.Logger.Debug("Error saving account map - ", err)
			code := codes.New(codes.ErrorAccountMapWrite)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		account.Logger.Debug("Error mapping account - ", err)
		code := codes.New(codes.ErrorAccountMapSave)
		return code
	}
	return nil
}

// Lookup searches the database for a matching account name and loads it from the db
func (account *Account) Lookup() error {
	originalID := account.ID
	account.ID = uuid.Nil
	// we don't expect many accounts to exist in the db (typically just one), so we iterate through them all
	err := account.MasterDB.View(func(tx *bolt.Tx) error {
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
				account.Logger.Debug("Error deriving account map key - ", err)
				code := codes.New(codes.ErrorAccountMapDerive)
				return code
			}
			// the new key should match the existing key if we have the right name and salt
			if subtle.ConstantTimeCompare(encryptedName[:], checkName[:]) == 1 {
				account.ID, err = uuid.FromBytes(value)
				if err != nil {
					account.Logger.Debug("Error converting account map uuid - ", err)
					code := codes.New(codes.ErrorAccountMapConvert)
					return code
				}
				return nil
			}
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		account.Logger.Debug("Error looking up account map - ", err)
		code := codes.New(codes.ErrorAccountMapLookup)
		return code
	}
	if account.ID == uuid.Nil {
		account.ID = originalID
		account.Logger.Debug("account lookup - no account found for [", account.Name, "]")
		code := codes.New(codes.ErrorAccountLookup)
		return code
	}
	return nil
}

// Load loads an account from the database
func (account *Account) Load() error {
	err := account.MasterDB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("accounts"))
		if bucket == nil {
			account.Logger.Debug("account bucket does not exist")
			code := codes.New(codes.ErrorAccountBucketMissing)
			return code
		}
		cursor := bucket.Cursor()
		key, value := cursor.Seek(account.ID.Bytes())
		if key == nil {
			account.Logger.Debug("Error loading account")
			code := codes.New(codes.ErrorAccountLoad)
			return code
		}

		// account data is encrypted with the account key and not the user key
		accountKey, err := crypto.Open(account.ActiveUser.PassphraseKey, account.ActiveUser.AccountKey)
		if err != nil {
			account.Logger.Debug("Error opening account key - ", err)
			code := codes.New(codes.ErrorAccountKeyOpen)
			return code
		}

		// decrypt value
		decryptedData, err := crypto.Open(accountKey, value)
		if err != nil {
			account.Logger.Debug("Error decrypting account data - ", err)
			code := codes.New(codes.ErrorAccountDecrypt)
			return code
		}

		err = json.Unmarshal(decryptedData, account)
		if err != nil {
			account.Logger.Debug("Error decoding account json - ", err)
			code := codes.New(codes.ErrorAccountDecode)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		account.Logger.Debug("Error loading account - ", err)
		code := codes.New(codes.ErrorAccountLoadView)
		return code
	}
	return nil
}
