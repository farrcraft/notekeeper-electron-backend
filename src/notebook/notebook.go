package notebook

import (
	"encoding/json"
	"time"

	"../crypto"
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
	Title        *title.Title   `json:"title"`          // Title is the title of the notebook
	Default      bool           `json:"default"`        // Default indicates whether this is the default notebook
	EncryptedKey []byte         `json:"encryption_key"` // EncryptedKey is the encrypted version of the notebook's encryption key
	Notes        []*note.Note   `json:"-"`              // Notes is the set of notes that belong to this notebook
	Tags         []*tag.Tag     `json:"tags"`           // Tags is the set of tags assigned to this notebook
	Created      time.Time      `json:"created"`        // Created is the time when the notebook was created
	Updated      time.Time      `json:"updated"`        // Updated is the time when the notebook was last updated
	Locked       bool           `json:"locked"`         // Locked indicates whether the notebook can be modified
	DB           *bolt.DB       `json:"-"`
	Logger       *logrus.Logger `json:"-"`
}

// NewNotebook creates a new notebook object
func NewNotebook(db *bolt.DB, logger *logrus.Logger) *Notebook {
	now := time.Now()
	notebook := &Notebook{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
		Default: false,
		Locked:  false,
		DB:      db,
		Logger:  logger,
	}
	return notebook
}

// Save saves a notebook to the database
// Account.ActiveUser.PassphraseKey
func (notebook *Notebook) Save(passphraseKey []byte) error {
	notebook.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		bucket, err := tx.CreateBucketIfNotExists([]byte("notebooks"))
		if err != nil {
			notebook.Logger.Error("Error creating notebook bucket - ", err)
			return err
		}

		// serialize notebook data
		data, err := json.Marshal(notebook)
		if err != nil {
			notebook.Logger.Error("Error marshaling notebook - ", err)
			return err
		}

		// retrieve the encryption key
		decryptedKey, err := crypto.Open(passphraseKey, notebook.EncryptedKey)
		if err != nil {
			notebook.Logger.Error("Error retrieving notebook key - ", err)
			return err
		}

		// encrypt the data
		encryptedData, err := crypto.Seal(decryptedKey, data)
		if err != nil {
			notebook.Logger.Error("Error encrypting notebook data - ", err)
			return err
		}

		// finally, save it
		err = bucket.Put(notebook.ID.Bytes(), encryptedData)
		if err != nil {
			notebook.Logger.Error("Error saving notebook - ", err)
			return err
		}
		return nil
	})

	return nil
}