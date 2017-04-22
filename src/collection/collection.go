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

// Collection holds a collection of notebooks
type Collection struct {
	ID        uuid.UUID            `json:"id"`         // ID is the unique collection identifier
	Title     *title.Title         `json:"title"`      // Title is a title for the collection
	Notebooks []*notebook.Notebook `json:"-"`          // Notebooks is the set of notebooks in the collection
	AccountID uuid.UUID            `json:"account_id"` // AccountID is the account that the collection belongs to
	UserID    uuid.UUID            `json:"user_id"`    // UserID is the individual account user that the collection belongs to
	ShelfID   uuid.UUID            `json:"shelf_id"`   // ShelfID is the shelf that contains the collection
	Tags      []*tag.Tag           `json:"tags"`       // Tags is the set of tags assigned to the collection
	Created   time.Time            `json:"created"`    // Created is the time when the collection was first created
	Updated   time.Time            `json:"updated"`    // Updated is the time when the collection was last updated
	Locked    bool                 `json:"locked"`     // Locked indicates whether the collection can be modified
	DBFactory *db.Factory          `json:"-"`          // DBFactory provides database access
	Logger    *logrus.Logger       `json:"-"`          // Logger is the logging facility
}

// New creates a new collection object
func New(title *title.Title, dbFactory *db.Factory, logger *logrus.Logger) *Collection {
	now := time.Now()
	collection := &Collection{
		ID:        uuid.NewV4(),
		Title:     title,
		Created:   now,
		Updated:   now,
		Locked:    false,
		DBFactory: dbFactory,
		Logger:    logger,
	}
	return collection
}

// Save a collection
func (collection *Collection) Save(passphraseKey []byte) error {
	shelfDB := collection.DBFactory.Find(db.TypeShelf, collection.ShelfID)
	err := shelfDB.DB.Update(func(tx *bolt.Tx) error {
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
		decryptedKey, err := crypto.Open(passphraseKey, shelfDB.EncryptedKey)
		if err != nil {
			collection.Logger.Debug("Error retrieving collection key - ", err)
			code := codes.New(codes.ScopeCollection, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := crypto.Seal(decryptedKey, data)
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
	return collections, nil
}

// Delete a collection
func (collection *Collection) Delete() error {
	return nil
}
