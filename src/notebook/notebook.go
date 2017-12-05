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

// ContainerType indicates the type of container that holds the notebook
type ContainerType int

const (
	// ContainerTypeShelf indicates that a notebook is contained in a shelf
	ContainerTypeShelf ContainerType = iota
	// ContainerTypeCollection indicates that a notebook is contained in a collection
	ContainerTypeCollection
)

// Scope indicates the scope of the notebook
type Scope int

const (
	// ScopeUser indicates that a notebook belongs to a single user
	ScopeUser Scope = iota
	// ScopeAccount indicates that a notebook belongs to a whole account
	ScopeAccount
)

// Notebook contains a set of notes
// notebooks are stored within either a collection or a shelf
// the parent container is owned either by a user or an account
type Notebook struct {
	ID            uuid.UUID      `json:"id"`             // ID is the unique identifier for this notebook
	OwnerID       uuid.UUID      `json:"owner_id"`       // OwnerID is the user or account that owns the notebook
	ContainerID   uuid.UUID      `json:"container_id"`   // ContainerID is the collection or shelf that this notebook belongs to
	ContainerType ContainerType  `json:"container_type"` // ContainerType is the type of container the notebook is stored within (shelf or collection)
	Scope         Scope          `json:"scope"`          // Scope is whether the notebook is owned by a user or an account
	Title         *title.Title   `json:"title"`          // Title is the title of the notebook
	Default       bool           `json:"default"`        // Default indicates whether this is the default notebook
	EncryptedKey  []byte         `json:"encryption_key"` // EncryptedKey is the encrypted version of the notebook's encryption key
	Notes         []*note.Note   `json:"-"`              // Notes is the set of notes that belong to this notebook
	NoteCount     int            `json:"note_count"`     // NoteCount keeps track of the number of notes in the notebook
	Tags          []*tag.Tag     `json:"tags"`           // Tags is the set of tags assigned to this notebook
	Created       time.Time      `json:"created"`        // Created is the time when the notebook was created
	Updated       time.Time      `json:"updated"`        // Updated is the time when the notebook was last updated
	Locked        bool           `json:"locked"`         // Locked indicates whether the notebook can be modified
	DBRegistry    *db.Registry   `json:"-"`
	Logger        *logrus.Logger `json:"-"`
}

// New creates a new notebook object
func New(title *title.Title, scope Scope, container ContainerType, dbRegistry *db.Registry, logger *logrus.Logger) *Notebook {
	now := time.Now()
	notebook := &Notebook{
		ID:            uuid.NewV4(),
		Scope:         scope,
		ContainerType: container,
		Title:         title,
		Created:       now,
		Updated:       now,
		NoteCount:     0,
		Default:       false,
		Locked:        false,
		DBRegistry:    dbRegistry,
		Logger:        logger,
	}
	return notebook
}

func (notebook *Notebook) getDBHandle(passphraseKey []byte) (*db.Handle, error) {
	var key db.Key
	key.ID = notebook.ContainerID
	if notebook.ContainerType == ContainerTypeCollection {
		key.Type = db.TypeCollection
	} else {
		key.Type = db.TypeShelf
	}
	notebookDBHandle, err := notebook.DBRegistry.GetHandle(key, passphraseKey)
	/*
		// [FIXME] - open if db nil
		if notebookDB == nil {
			// notebook db is either a shelf or collection db
			key := db.Key{
				ID:   notebook.ContainerID,
				Type: dbType,
			}
			var parentDBType db.Type
			if notebook.Scope == ScopeAccount {
				parentDBType = db.TypeAccount
			} else {
				parentDBType = db.TypeUser
			}
			parentKey := db.Key{
				ID:   notebook.OwnerID,
				Type: parentDBType,
			}
			var err error
			notebookDB, err = notebook.DBFactory.Open(key, parentKey, parentKey, passphraseKey)
			if err != nil {
				return nil
			}
			notebook.Logger.Debug("opened notebook parent db")
		} else {
			notebook.Logger.Debug("notebook db already open")
		}
	*/
	return notebookDBHandle, err
}

// Save a notebook to the database
// Account.ActiveUser.PassphraseKey
func (notebook *Notebook) Save(passphraseKey []byte) error {
	notebookDBHandle, err := notebook.getDBHandle(passphraseKey)
	if err != nil {
		return err
	}
	err = notebookDBHandle.DB.Update(func(tx *bolt.Tx) error {
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
		c := crypto.New(notebook.Logger)
		decryptedKey, err := c.Open(passphraseKey, notebookDBHandle.EncryptedKey)
		if err != nil {
			notebook.Logger.Debug("Error retrieving notebook key - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := c.Seal(decryptedKey, data)
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

// LoadAll notebooks
func (notebook *Notebook) LoadAll(passphraseKey []byte) ([]*Notebook, error) {
	var notebooks []*Notebook

	notebookDBHandle, err := notebook.getDBHandle(passphraseKey)
	if err != nil {
		return notebooks, err
	}
	c := crypto.New(notebook.Logger)
	notebookKey, err := c.Open(passphraseKey, notebookDBHandle.EncryptedKey)
	if err != nil {
		notebook.Logger.Debug("Error opening notebook key - ", err)
		code := codes.New(codes.ScopeNotebook, codes.ErrorOpenKey)
		return notebooks, code
	}

	err = notebookDBHandle.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("notebooks"))
		if bucket == nil {
			notebook.Logger.Debug("notebook bucket does not exist")
			code := codes.New(codes.ScopeNotebook, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newNotebook := &Notebook{
				DBRegistry: notebook.DBRegistry,
				Logger:     notebook.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(notebookKey, value)
			if err != nil {
				notebook.Logger.Debug("Error decrypting notebook data - ", err)
				code := codes.New(codes.ScopeNotebook, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newNotebook)
			if err != nil {
				notebook.Logger.Debug("Error decoding notebook json - ", err)
				code := codes.New(codes.ScopeNotebook, codes.ErrorDecode)
				return code
			}

			notebooks = append(notebooks, newNotebook)
		}

		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return notebooks, err
		}
		notebook.Logger.Debug("Error loading all notebooks - ", err)
		code := codes.New(codes.ScopeNotebook, codes.ErrorLoadAll)
		return nil, code
	}

	return notebooks, nil
}

// Delete a notebook
func (notebook *Notebook) Delete(passphraseKey []byte) error {
	notebookDBHandle, err := notebook.getDBHandle(passphraseKey)
	if err != nil {
		return err
	}
	err = notebookDBHandle.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notebooks"))
		if bucket == nil {
			notebook.Logger.Debug("notebook bucket does not exist")
			code := codes.New(codes.ScopeNotebook, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(notebook.ID.Bytes())
		if err != nil {
			notebook.Logger.Debug("Error deleting notebook - ", err)
			code := codes.New(codes.ScopeNotebook, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		notebook.Logger.Debug("Error deleting notebook - ", err)
		code := codes.New(codes.ScopeNotebook, codes.ErrorDelete)
		return code
	}

	return nil
}
