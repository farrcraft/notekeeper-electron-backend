package account

import (
	"crypto/subtle"

	"../codes"
	"../crypto"
	"../db"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Index contains an index of accounts
type Index struct {
	Accounts   []*Account     `json:"accounts"`
	DBRegistry *db.Registry   `json:"-"` // DBRegistry provides access to the database
	Logger     *logrus.Logger `json:"-"` // Logger is the logging facility
}

// NewIndex returns a new index object
func NewIndex(dbRegistry *db.Registry, logger *logrus.Logger) *Index {
	index := &Index{
		DBRegistry: dbRegistry,
		Logger:     logger,
	}
	return index
}

// Count returns the number of records in the account_index table
func (index *Index) Count() int {
	count := 0
	key := db.Key{
		ID:   uuid.Nil,
		Type: db.TypeMaster,
	}
	handle, err := index.DBRegistry.GetHandle(key)
	if err != nil {
		return count
	}
	handle.DB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("account_index"))
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

// Save an account in the index
// Since accounts are keyed by only an unencrypted id in the db
// we also need to store a mapping between a key derived from the name and the id
// otherwise there is no way to look up an account without taking a brute force decryption test approach
func (index *Index) Save(account *Account, passphraseKey []byte) error {
	// account index is stored in the master db
	masterKey := db.Key{
		ID:   uuid.Nil,
		Type: db.TypeMaster,
	}
	masterDBHandle, err := index.DBRegistry.GetHandle(masterKey)
	if err != nil {
		return err
	}

	err = masterDBHandle.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("account_index"))
		if err != nil {
			index.Logger.Debug("Error creating account index bucket - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorCreateBucket)
			return code
		}

		c := crypto.New(index.Logger)
		encryptedName, err := c.DeriveSaltedKey([]byte(account.Name))
		if err != nil {
			index.Logger.Debug("Error creating account index key - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorDeriveKey)
			return code
		}
		err = bucket.Put(encryptedName, account.ID.Bytes())
		if err != nil {
			index.Logger.Debug("Error saving account index - ", err)
			code := codes.New(codes.ScopeAccount, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		index.Logger.Debug("Error saving account - err")
		code := codes.New(codes.ScopeAccount, codes.ErrorSave)
		return code
	}

	return nil
}

// Lookup searches the index for a matching account name and sets the account id if it exists
func (index *Index) Lookup(account *Account) error {
	originalID := account.ID
	account.ID = uuid.Nil
	key := db.Key{
		ID:   uuid.Nil,
		Type: db.TypeMaster,
	}
	handle, err := index.DBRegistry.GetHandle(key)
	if err != nil {
		return err
	}
	err = handle.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("account_index"))
		// If bucket doesn't exist, no accounts have been created yet
		if bucket == nil {
			index.Logger.Debug("account index bucket does not exist")
			return nil
		}
		// we don't expect many accounts to exist in the db (typically just one), so we iterate through them all
		cursor := bucket.Cursor()
		c := crypto.New(index.Logger)
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			// extract the salt from the existing encrypted name
			salt, encryptedName := crypto.ExtractSalt(key)
			// create a new key using the extracted salt and the unencrypted name we're searching for
			checkName, err := c.DeriveKey([]byte(account.Name), salt[:])
			if err != nil {
				index.Logger.Debug("Error deriving account index key - ", err)
				code := codes.New(codes.ScopeAccount, codes.ErrorDeriveKey)
				return code
			}
			// the new key should match the existing key if we have the right name and salt
			if subtle.ConstantTimeCompare(encryptedName[:], checkName[:]) == 1 {
				account.ID, err = uuid.FromBytes(value)
				if err != nil {
					index.Logger.Debug("Error converting account index uuid - ", err)
					code := codes.New(codes.ScopeAccount, codes.ErrorConvertID)
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
		index.Logger.Debug("Error looking up account index - ", err)
		code := codes.New(codes.ScopeAccount, codes.ErrorLookup)
		return code
	}
	if account.ID == uuid.Nil {
		account.ID = originalID
		index.Logger.Debug("account lookup - no account found for [", account.Name, "]")
		code := codes.New(codes.ScopeAccount, codes.ErrorLookup)
		return code
	}
	return nil
}
