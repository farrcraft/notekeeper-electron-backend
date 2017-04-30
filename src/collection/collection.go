package collection

import (
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../notebook"
	"../tag"
	"../title"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Scope indicates the scope of the collection
type Scope int

const (
	// ScopeUser indicates that a collection belongs to a single user
	ScopeUser Scope = iota
	// ScopeAccount indicates that a collection belongs to a whole account
	ScopeAccount
)

// Collection holds a collection of notebooks
type Collection struct {
	ID        uuid.UUID            `json:"id"`       // ID is the unique collection identifier
	Title     *title.Title         `json:"title"`    // Title is a title for the collection
	Notebooks []*notebook.Notebook `json:"-"`        // Notebooks is the set of notebooks in the collection
	OwnerID   uuid.UUID            `json:"owner_id"` // OwnerID is the account or user that the collection belongs to
	Scope     Scope                `json:"scope"`    // Scope is the ownership scope of the collection (user or account)
	ShelfID   uuid.UUID            `json:"shelf_id"` // ShelfID is the shelf that contains the collection
	Tags      []*tag.Tag           `json:"tags"`     // Tags is the set of tags assigned to the collection
	Created   time.Time            `json:"created"`  // Created is the time when the collection was first created
	Updated   time.Time            `json:"updated"`  // Updated is the time when the collection was last updated
	Locked    bool                 `json:"locked"`   // Locked indicates whether the collection can be modified
	DBFactory *db.Factory          `json:"-"`        // DBFactory provides database access
	Logger    *logrus.Logger       `json:"-"`        // Logger is the logging facility
}

// New creates a new collection object
func New(title *title.Title, scope Scope, dbFactory *db.Factory, logger *logrus.Logger) *Collection {
	now := time.Now()
	collection := &Collection{
		ID:        uuid.NewV4(),
		Scope:     scope,
		Title:     title,
		Created:   now,
		Updated:   now,
		Locked:    false,
		DBFactory: dbFactory,
		Logger:    logger,
	}
	return collection
}

func (collection *Collection) getDB(passphraseKey []byte) (*db.DB, error) {
	// even though the *content* of a collection gets its own db, the collection
	// itself is stored in the parent db
	var dbType db.Type
	if collection.Scope == ScopeUser {
		dbType = db.TypeUser
	} else {
		dbType = db.TypeAccount
	}

	collectionDB := collection.DBFactory.Find(db.TypeShelf, collection.ShelfID)
	if collectionDB == nil {
		key := db.Key{
			ID:   collection.ID,
			Type: db.TypeCollection,
		}
		parentKey := db.Key{
			ID:   collection.ShelfID,
			Type: dbType,
		}
		var err error
		collectionDB, err = collection.DBFactory.Open(key, parentKey, collection.OwnerID, passphraseKey)
		if err != nil {
			return nil, err
		}
	}
	return collectionDB, nil
}

// Save a collection
func (collection *Collection) Save(passphraseKey []byte) error {
	shelfDB, err := collection.getDB(passphraseKey)
	if err != nil {
		return err
	}
	err = shelfDB.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("collection_index"))
		if err != nil {
			collection.Logger.Debug("Error creating collection bucket - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorCreateBucket)
			return code
		}

		// serialize collection data
		data, err := json.Marshal(collection)
		if err != nil {
			collection.Logger.Debug("Error marshaling collection - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		c := crypto.New(collection.Logger)
		decryptedKey, err := c.Open(passphraseKey, shelfDB.EncryptedKey)
		if err != nil {
			collection.Logger.Debug("Error retrieving collection key - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := c.Seal(decryptedKey, data)
		if err != nil {
			collection.Logger.Debug("Error encrypting collection data - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(collection.ID.Bytes(), encryptedData)
		if err != nil {
			collection.Logger.Debug("Error writing collection - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		collection.Logger.Debug("Error saving collection - err")
		code := codes.New(codes.ScopeCollection, codes.ErrorSave)
		return code
	}

	return nil
}

// LoadAll collections
func (collection *Collection) LoadAll(passphraseKey []byte) ([]*Collection, error) {
	var collections []*Collection

	shelfDB, err := collection.getDB(passphraseKey)
	if err != nil {
		return collections, err
	}

	c := crypto.New(collection.Logger)
	shelfKey, err := c.Open(passphraseKey, shelfDB.EncryptedKey)
	if err != nil {
		collection.Logger.Debug("Error opening collection key - ", err)
		code := codes.New(codes.ScopeCollection, codes.ErrorOpenKey)
		return collections, code
	}

	err = shelfDB.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("collection_index"))
		if bucket == nil {
			collection.Logger.Debug("collection index bucket does not exist")
			code := codes.New(codes.ScopeCollection, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newCollection := &Collection{
				DBFactory: collection.DBFactory,
				Logger:    collection.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(shelfKey, value)
			if err != nil {
				collection.Logger.Debug("Error decrypting collection data - ", err)
				code := codes.New(codes.ScopeCollection, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newCollection)
			if err != nil {
				collection.Logger.Debug("Error decoding collection json - ", err)
				code := codes.New(codes.ScopeCollection, codes.ErrorDecode)
				return code
			}

			collections = append(collections, newCollection)
		}

		return nil
	})

	return collections, nil
}

// Delete a collection
func (collection *Collection) Delete(passphraseKey []byte) error {
	shelfDB, err := collection.getDB(passphraseKey)
	if err != nil {
		return err
	}

	err = shelfDB.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("collection_index"))
		if bucket == nil {
			collection.Logger.Debug("collection index bucket does not exist")
			code := codes.New(codes.ScopeCollection, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(collection.ID.Bytes())
		if err != nil {
			collection.Logger.Debug("Error deleting collection - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		collection.Logger.Debug("Error deleting collection - ", err)
		code := codes.New(codes.ScopeCollection, codes.ErrorDelete)
		return code
	}

	return nil
}
