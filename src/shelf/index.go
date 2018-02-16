package shelf

import (
	"encoding/json"

	"../codes"
	"../crypto"
	"../db"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Index of shelf metadata
// While metadata about objects stored within a shelf is stored in the shelf's database,
// the metadata for the shelf itself is kept in an index within its parent database
type Index struct {
	Scope      Scope          `json:"scope"` // Type is one of the Type* identifier values
	OwnerID    uuid.UUID      `json:"-"`     // OwnerID is the ID of the account or user owning the shelf
	Shelves    []*Shelf       `json:"shelves"`
	DBRegistry *db.Registry   `json:"-"` // DBRegistry provides access to the database
	Logger     *logrus.Logger `json:"-"` // Logger is the logging facility
}

// NewIndex returns a new index object
func NewIndex(scope Scope, ownerID uuid.UUID, dbRegistry *db.Registry, logger *logrus.Logger) *Index {
	index := &Index{
		Scope:      scope,
		OwnerID:    ownerID,
		DBRegistry: dbRegistry,
		Logger:     logger,
	}
	return index
}

func (index *Index) getDBHandle() (*db.Handle, error) {
	var key db.Key
	key.ID = index.OwnerID
	if index.Scope == ScopeUser {
		key.Type = db.TypeUser
	} else {
		key.Type = db.TypeAccount
	}
	handle, err := index.DBRegistry.GetHandle(key)
	return handle, err
}

// Save a shelf in the index
func (index *Index) Save(shelf *Shelf, encryptionKey []byte) error {
	handle, err := index.getDBHandle()
	if err != nil {
		return err
	}
	err = handle.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("shelf_index"))
		if err != nil {
			index.Logger.Debug("Error creating shelf index bucket - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorCreateBucket)
			return code
		}

		// serialize shelf data
		data, err := json.Marshal(shelf)
		if err != nil {
			index.Logger.Debug("Error marshaling shelf - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorMarshal)
			return code
		}

		// encrypt the data
		c := crypto.New(index.Logger)
		encryptedData, err := c.Seal(encryptionKey, data)
		if err != nil {
			index.Logger.Debug("Error encrypting shelf data - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(shelf.ID.Bytes(), encryptedData)
		if err != nil {
			index.Logger.Debug("Error writing shelf - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		index.Logger.Debug("Error saving shelf - err")
		code := codes.New(codes.ScopeShelf, codes.ErrorSave)
		return code
	}

	return nil
}

// LoadAll of the shelves from an account or user DB
func (index *Index) LoadAll(passphraseKey []byte) error {
	handle, err := index.getDBHandle()
	if err != nil {
		return err
	}

	// [FIXME] - method should recieve unsealed encryption key directly
	c := crypto.New(index.Logger)
	shelfKey, err := c.Open(passphraseKey, handle.EncryptedKey)
	if err != nil {
		index.Logger.Debug("Error opening shelf key - ", err)
		code := codes.New(codes.ScopeShelf, codes.ErrorOpenKey)
		return code
	}

	err = handle.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("shelf_index"))
		if bucket == nil {
			index.Logger.Debug("shelf bucket does not exist")
			code := codes.New(codes.ScopeShelf, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newShelf := &Shelf{
				DBRegistry: index.DBRegistry,
				Logger:     index.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(shelfKey, value)
			if err != nil {
				index.Logger.Debug("Error decrypting shelf data - ", err)
				code := codes.New(codes.ScopeShelf, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newShelf)
			if err != nil {
				index.Logger.Debug("Error decoding shelf json - ", err)
				code := codes.New(codes.ScopeShelf, codes.ErrorDecode)
				return code
			}

			index.Shelves = append(index.Shelves, newShelf)
		}

		return nil
	})

	return err
}

// Delete a shelf from the index
func (index *Index) Delete(shelf *Shelf) error {
	handle, err := index.getDBHandle()
	if err != nil {
		return err
	}
	err = handle.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("shelf_index"))
		if bucket == nil {
			index.Logger.Debug("shelf bucket does not exist")
			code := codes.New(codes.ScopeShelf, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(shelf.ID.Bytes())
		if err != nil {
			index.Logger.Debug("Error deleting shelf - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		index.Logger.Debug("Error deleting shelf - ", err)
		code := codes.New(codes.ScopeShelf, codes.ErrorDelete)
		return code
	}

	return nil
}
