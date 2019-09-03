package collection

import (
	"encoding/json"

	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/crypto"
	"notekeeper-electron-backend/db"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
)

// Index of collection metadata
// While metadata about objects stored within a collection is stored in the collection's database,
// the metadata for the collection itself is kept in an index within its parent database
type Index struct {
	Scope       Scope          `json:"scope"` // Type is one of the Type* identifier values
	ShelfID     uuid.UUID      `json:"-"`
	OwnerID     uuid.UUID      `json:"-"` // OwnerID is the ID of the account or user owning the collection
	Collections []*Collection  `json:"collections"`
	DBRegistry  *db.Registry   `json:"-"` // DBRegistry provides access to the database
	Logger      *logrus.Logger `json:"-"` // Logger is the logging facility
}

// NewIndex returns a new index object
func NewIndex(scope Scope, dbRegistry *db.Registry, logger *logrus.Logger) *Index {
	index := &Index{
		Scope:      scope,
		DBRegistry: dbRegistry,
		Logger:     logger,
	}
	return index
}

func (index *Index) getDBHandle() (*db.Handle, error) {
	// even though the *content* of a collection gets its own db, the collection
	// itself is stored in the parent db
	shelfKey := db.Key{
		Type: db.TypeShelf,
		ID:   index.ShelfID,
	}
	collectionDBHandle, err := index.DBRegistry.GetHandle(shelfKey)
	/*
		[FIXME]
		if err != nil {
			return collectionDBHandle, err
		}
		if collectionDBHandle != nil {
			return collectionDBHandle, nil
		}
		key := db.Key{
			ID:   collection.ID,
			Type: db.TypeCollection,
		}
		parentKey := db.Key{
			ID:   collection.ShelfID,
			Type: db.TypeShelf,
		}
		var ownerType db.Type
		if collection.Scope == ScopeUser {
			ownerType = db.TypeUser
		} else {
			ownerType = db.TypeAccount
		}
		ownerKey := db.Key{
			ID:   collection.OwnerID,
			Type: ownerType,
		}
		var err error
		collectionDBHandle, err = collection.DBRegistry.Open(key, parentKey, ownerKey, passphraseKey)
		if err != nil {
			return nil, err
		}

		return collectionDBHandle, nil
	*/
	return collectionDBHandle, err
}

// Save a collection in the collection index
func (index *Index) Save(collection *Collection, passphraseKey []byte) error {
	shelfDBHandle, err := index.getDBHandle()
	if err != nil {
		return err
	}
	err = shelfDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("collection_index"))
		if err != nil {
			index.Logger.Warn("Error creating collection bucket - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorCreateBucket)
			return code
		}

		// serialize collection data
		data, err := json.Marshal(collection)
		if err != nil {
			index.Logger.Warn("Error marshaling collection - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		c := crypto.New(index.Logger)
		decryptedKey, err := c.Open(passphraseKey, shelfDBHandle.EncryptedKey)
		if err != nil {
			index.Logger.Warn("Error retrieving collection key - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := c.Seal(decryptedKey, data)
		if err != nil {
			index.Logger.Warn("Error encrypting collection data - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(collection.ID.Bytes(), encryptedData)
		if err != nil {
			index.Logger.Warn("Error writing collection - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		index.Logger.Warn("Error saving collection - err")
		code := codes.New(codes.ScopeCollection, codes.ErrorSave)
		return code
	}

	return nil
}

// LoadAll collections
func (index *Index) LoadAll(passphraseKey []byte) error {
	shelfDBHandle, err := index.getDBHandle()
	if err != nil {
		return err
	}

	c := crypto.New(index.Logger)
	shelfKey, err := c.Open(passphraseKey, shelfDBHandle.EncryptedKey)
	if err != nil {
		index.Logger.Warn("Error opening collection key - ", err)
		code := codes.New(codes.ScopeCollection, codes.ErrorOpenKey)
		return code
	}

	err = shelfDBHandle.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("collection_index"))
		if bucket == nil {
			index.Logger.Warn("collection index bucket does not exist")
			code := codes.New(codes.ScopeCollection, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newCollection := &Collection{
				DBRegistry: index.DBRegistry,
				Logger:     index.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(shelfKey, value)
			if err != nil {
				index.Logger.Warn("Error decrypting collection data - ", err)
				code := codes.New(codes.ScopeCollection, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newCollection)
			if err != nil {
				index.Logger.Warn("Error decoding collection json - ", err)
				code := codes.New(codes.ScopeCollection, codes.ErrorDecode)
				return code
			}

			index.Collections = append(index.Collections, newCollection)
		}

		return nil
	})

	return nil
}

// Delete a collection from the collection index
func (index *Index) Delete(collection *Collection, passphraseKey []byte) error {
	shelfDBHandle, err := index.getDBHandle()
	if err != nil {
		return err
	}

	err = shelfDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("collection_index"))
		if bucket == nil {
			index.Logger.Warn("collection index bucket does not exist")
			code := codes.New(codes.ScopeCollection, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(collection.ID.Bytes())
		if err != nil {
			index.Logger.Warn("Error deleting collection - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		index.Logger.Warn("Error deleting collection - ", err)
		code := codes.New(codes.ScopeCollection, codes.ErrorDelete)
		return code
	}

	return nil
}
