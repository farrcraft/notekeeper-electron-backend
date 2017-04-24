package shelf

import (
	"encoding/json"
	"time"

	"../codes"
	"../collection"
	"../crypto"
	"../db"
	"../notebook"
	"../tag"
	"../title"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Scope indicates the scope of the shelf
type Scope int

const (
	// ScopeUser indicates that a shelf belongs to a single user
	ScopeUser Scope = iota
	// ScopeAccount indicates that a shelf belongs to a whole account
	ScopeAccount
)

// Shelf contains notebooks & collections of notebooks
type Shelf struct {
	ID           uuid.UUID                `json:"id"`             // ID is the unique identifier for the shelf
	Title        *title.Title             `json:"title"`          // Title is the title of the shelf
	Scope        Scope                    `json:"scope"`          // Type is one of the Type* identifier values
	Default      bool                     `json:"default"`        // Default indicates whether this is the default shelf
	Trash        bool                     `json:"trash"`          // Trash indicates whether this is the trash shelf
	AccountID    uuid.UUID                `json:"-"`              // AccountID is the ID of the account owning the shelf
	UserID       uuid.UUID                `json:"-"`              // UserID is the ID of the user owning the shelf
	EncryptedKey []byte                   `json:"encryption_key"` // EncryptedKey is the encrypted encryption key for the shelf DB
	Notebooks    []*notebook.Notebook     `json:"-"`              // Notebooks is the set of notebooks in the shelf
	Collections  []*collection.Collection `json:"-"`              // Collections is the set of collections in the shelf
	Tags         []*tag.Tag               `json:"tags"`           // Tags is the set of tags assigned to the shelf
	Created      time.Time                `json:"created"`        // Created is the time when the shelf was created
	Updated      time.Time                `json:"updated"`        // Updated is the time when the shelf was last updated
	Locked       bool                     `json:"locked"`         // Locked indicates whether the shelf can be modified
	DBFactory    *db.Factory              `json:"-"`              // DBFactory provides access to the database
	Logger       *logrus.Logger           `json:"-"`              // Logger is the logging facility
}

// New creates a new shelf object
func New(title *title.Title, scope Scope, dbFactory *db.Factory, logger *logrus.Logger) *Shelf {
	now := time.Now()
	shelf := &Shelf{
		ID:        uuid.NewV4(),
		Title:     title,
		Scope:     scope,
		Created:   now,
		Updated:   now,
		Default:   false,
		Trash:     false,
		Locked:    false,
		DBFactory: dbFactory,
		Logger:    logger,
	}
	return shelf
}

func (shelf *Shelf) getDB(passphraseKey []byte) (*db.DB, error) {
	// even though the *content* of a shelf gets its own db, the shelf itself
	// is stored in the parent db
	var dbType db.Type
	var id uuid.UUID
	if shelf.Scope == ScopeUser {
		dbType = db.TypeUser
		id = shelf.UserID
	} else {
		dbType = db.TypeAccount
		id = shelf.AccountID
	}
	shelfDB := shelf.DBFactory.Find(dbType, id)
	if shelfDB == nil {
		key := db.Key{
			ID:   shelf.ID,
			Type: db.TypeShelf,
		}
		parentKey := db.Key{
			ID:   id,
			Type: dbType,
		}
		var err error
		shelfDB, err = shelf.DBFactory.Open(key, parentKey, id, passphraseKey)
		if err != nil {
			return shelfDB, err
		}
	}
	return shelfDB, nil
}

// Save a shelf to the DB
func (shelf *Shelf) Save(passphraseKey []byte) error {
	db, err := shelf.getDB(passphraseKey)
	if err != nil {
		return err
	}
	err = db.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("shelf_index"))
		if err != nil {
			shelf.Logger.Debug("Error creating shelf bucket - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorCreateBucket)
			return code
		}

		// serialize shelf data
		data, err := json.Marshal(shelf)
		if err != nil {
			shelf.Logger.Debug("Error marshaling shelf - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		decryptedKey, err := crypto.Open(passphraseKey, db.EncryptedKey)
		if err != nil {
			shelf.Logger.Debug("Error retrieving shelf key - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := crypto.Seal(decryptedKey, data)
		if err != nil {
			shelf.Logger.Debug("Error encrypting shelf data - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(shelf.ID.Bytes(), encryptedData)
		if err != nil {
			shelf.Logger.Debug("Error writing shelf - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		shelf.Logger.Debug("Error saving shelf - err")
		code := codes.New(codes.ScopeShelf, codes.ErrorSave)
		return code
	}

	return nil
}

/*
// Load a shelf from the DB
func (shelf *Shelf) Load() error {
	return nil
}
*/

// LoadAll of the shelves from an account or user DB
func (shelf *Shelf) LoadAll(passphraseKey []byte) ([]*Shelf, error) {
	var shelves []*Shelf

	shelfDB, err := shelf.getDB(passphraseKey)
	if err != nil {
		return shelves, err
	}

	shelfKey, err := crypto.Open(passphraseKey, shelfDB.EncryptedKey)
	if err != nil {
		shelf.Logger.Debug("Error opening shelf key - ", err)
		code := codes.New(codes.ScopeShelf, codes.ErrorOpenKey)
		return shelves, code
	}

	err = shelfDB.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("shelf_index"))
		if bucket == nil {
			shelf.Logger.Debug("shelf bucket does not exist")
			code := codes.New(codes.ScopeShelf, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newShelf := &Shelf{
				DBFactory: shelf.DBFactory,
				Logger:    shelf.Logger,
			}

			// decrypt value
			decryptedData, err := crypto.Open(shelfKey, value)
			if err != nil {
				shelf.Logger.Debug("Error decrypting shelf data - ", err)
				code := codes.New(codes.ScopeShelf, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newShelf)
			if err != nil {
				shelf.Logger.Debug("Error decoding shelf json - ", err)
				code := codes.New(codes.ScopeShelf, codes.ErrorDecode)
				return code
			}

			shelves = append(shelves, newShelf)
		}

		return nil
	})

	return shelves, err
}

// Delete a shelf
func (shelf *Shelf) Delete(passphraseKey []byte) error {
	shelfDB, err := shelf.getDB(passphraseKey)
	if err != nil {
		return err
	}
	err = shelfDB.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("shelf_index"))
		if bucket == nil {
			shelf.Logger.Debug("shelf bucket does not exist")
			code := codes.New(codes.ScopeShelf, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(shelf.ID.Bytes())
		if err != nil {
			shelf.Logger.Debug("Error deleting shelf - ", err)
			code := codes.New(codes.ScopeShelf, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		shelf.Logger.Debug("Error deleting shelf - ", err)
		code := codes.New(codes.ScopeShelf, codes.ErrorDelete)
		return code
	}

	return nil
}
