package notebook

import (
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../note"
	"../tag"
	"../title"
	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"

	uuid "github.com/satori/go.uuid"
)

// Notebook contains a set of notes
type Notebook struct {
	ID           uuid.UUID      `json:"id"`             // ID is the unique identifier for this notebook
	UserID       uuid.UUID      `json:"user_id"`        // UserID is the user that owns the notebook (only if this is a user-scoped notebook, otherwise nil)
	ShelfID      uuid.UUID      `json:"shelf_id"`       // ShelfID is the shelf that contains the notebook
	CollectionID uuid.UUID      `json:"collection_id"`  // CollectionID is the collection that this notebook belongs to
	Title        *title.Title   `json:"title"`          // Title is the title of the notebook
	Default      bool           `json:"default"`        // Default indicates whether this is the default notebook
	EncryptedKey []byte         `json:"encryption_key"` // EncryptedKey is the encrypted version of the notebook's encryption key
	Notes        []*note.Note   `json:"-"`              // Notes is the set of notes that belong to this notebook
	NoteCount    int            `json:"note_count"`     // NoteCount keeps track of the number of notes in the notebook
	Tags         []*tag.Tag     `json:"tags"`           // Tags is the set of tags assigned to this notebook
	Created      time.Time      `json:"created"`        // Created is the time when the notebook was created
	Updated      time.Time      `json:"updated"`        // Updated is the time when the notebook was last updated
	Locked       bool           `json:"locked"`         // Locked indicates whether the notebook can be modified
	DBFactory    *db.Factory    `json:"-"`
	Logger       *logrus.Logger `json:"-"`
}

// NewNotebook creates a new notebook object
func NewNotebook(dbFactory *db.Factory, logger *logrus.Logger) *Notebook {
	now := time.Now()
	notebook := &Notebook{
		ID:        uuid.NewV4(),
		Created:   now,
		Updated:   now,
		NoteCount: 0,
		Default:   false,
		Locked:    false,
		DBFactory: dbFactory,
		Logger:    logger,
	}
	return notebook
}

// Save a notebook to the database
// Account.ActiveUser.PassphraseKey
func (notebook *Notebook) Save(passphraseKey []byte) error {
	// [FIXME] - notebook will be stored either in a shelf or collection db
	var notebookDB *db.DB
	if notebook.CollectionID != uuid.Nil {
		notebookDB = notebook.DBFactory.Find(db.TypeCollection, notebook.CollectionID)
	} else {
		notebookDB = notebook.DBFactory.Find(db.TypeShelf, notebook.ShelfID)
	}
	err := notebookDB.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("notebooks"))
		if err != nil {
			notebook.Logger.Debug("Error creating notebook bucket - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorCreateBucket)
			return code
		}

		// serialize notebook data
		data, err := json.Marshal(notebook)
		if err != nil {
			notebook.Logger.Debug("Error marshaling notebook - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		decryptedKey, err := crypto.Open(passphraseKey, notebookDB.EncryptedKey)
		if err != nil {
			notebook.Logger.Debug("Error retrieving notebook key - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := crypto.Seal(decryptedKey, data)
		if err != nil {
			notebook.Logger.Debug("Error encrypting notebook data - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(notebook.ID.Bytes(), encryptedData)
		if err != nil {
			notebook.Logger.Debug("Error writing notebook - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		notebook.Logger.Debug("Error saving notebook - err")
		code := codes.New(codes.ScopeNotebook, codes.ErrorSave)
		return code
	}

	return nil
}
