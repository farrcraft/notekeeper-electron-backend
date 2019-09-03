package user

import (
	"crypto/subtle"

	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/crypto"
	"notekeeper-electron-backend/db"

	"github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

// Index contains an index of users
type Index struct {
	AccountID  uuid.UUID      `json:"account_id"`
	Users      []*User        `json:"users"`
	DBRegistry *db.Registry   `json:"-"` // DBRegistry provides access to the database
	Logger     *logrus.Logger `json:"-"` // Logger is the logging facility
}

// NewIndex returns a new index object
func NewIndex(accountID uuid.UUID, dbRegistry *db.Registry, logger *logrus.Logger) *Index {
	index := &Index{
		AccountID:  accountID,
		DBRegistry: dbRegistry,
		Logger:     logger,
	}
	return index
}

// Save saves the user to the database
func (index *Index) Save(user *User, passphraseKey []byte) error {
	c := crypto.New(index.Logger)
	// Since users are keyed by only an unencrypted id in the db
	// we also need to store a mapping between a key derived from the email address and the id
	// otherwise there is no way to look up a user without taking a brute force decryption test approach
	accountKey := db.Key{
		Type: db.TypeAccount,
		ID:   user.AccountID,
	}
	accountDBHandle, err := index.DBRegistry.GetHandle(accountKey)
	if err != nil {
		return err
	}
	err = accountDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("user_index"))
		if err != nil {
			index.Logger.Debug("Error creating user index bucket - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorCreateBucket)
			return code
		}
		encryptedEmail, err := c.DeriveKey([]byte(user.Profile.Email), user.Salt)
		if err != nil {
			index.Logger.Debug("Error creating user index key - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
			return code
		}
		saltedKey := c.EmbedSalt(encryptedEmail, user.Salt)
		err = bucket.Put(saltedKey, user.ID.Bytes())
		if err != nil {
			index.Logger.Debug("Error saving user index - ", err)
			code := codes.New(codes.ScopeUser, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		index.Logger.Debug("Error mapping user - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorSave)
		return code
	}
	return nil
}

// Lookup a user id in the account database
func (index *Index) Lookup(user *User) error {
	originalID := user.ID
	user.ID = uuid.Nil
	accountKey := db.Key{
		Type: db.TypeAccount,
		ID:   user.AccountID,
	}
	accountDBHandle, err := index.DBRegistry.GetHandle(accountKey)
	if err != nil {
		return err
	}
	err = accountDBHandle.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("user_index"))
		if bucket == nil {
			code := codes.New(codes.ScopeUser, codes.ErrorBucketMissing)
			return code
		}
		cursor := bucket.Cursor()

		c := crypto.New(index.Logger)
		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			// extract the salt from the existing encrypted email
			salt, encryptedEmail := crypto.ExtractSalt(key)
			// create a new key using the extracted salt and the unencrypted email we're searching for
			checkEmail, err := c.DeriveKey([]byte(user.Profile.Email), salt[:])
			if err != nil {
				index.Logger.Warn("Error deriving user index key - ", err)
				code := codes.New(codes.ScopeUser, codes.ErrorDeriveKey)
				return code
			}
			// the new key should match the existing key if we have the right email and salt
			if subtle.ConstantTimeCompare(encryptedEmail[:], checkEmail[:]) == 1 {
				user.ID, err = uuid.FromBytes(value)
				if err != nil {
					index.Logger.Warn("Error converting email index uuid - ", err)
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
		index.Logger.Warn("Error looking up user index - ", err)
		code := codes.New(codes.ScopeUser, codes.ErrorLookup)
		return code
	}
	if user.ID == uuid.Nil {
		user.ID = originalID
		code := codes.New(codes.ScopeUser, codes.ErrorUserMissing)
		return code
	}
	return nil
}
