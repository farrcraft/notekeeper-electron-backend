package note

import (
	"encoding/json"
	"time"

	"../codes"
	"../crypto"
	"../db"
	"../tag"
	"../title"

	"github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	uuid "github.com/satori/go.uuid"
)

// Type is the type of note
type Type int

const (
	// TypePlainText indicates that a note is just plain old text
	TypePlainText Type = iota
	// TypeRichText indicates that a note contains rich text
	TypeRichText
	// TypeMarkdown indicates that a note contains markdown text
	TypeMarkdown
	// TypeHTML indicates that a note contains HTML
	TypeHTML
	// TypeImage indicates that a note is an image file
	TypeImage
	// TypeFile indicates that a note is an arbitrary file attachment
	TypeFile
	// TypePdf indicates that a note is a PDF file
	TypePdf
	// TypeAudio indicates that a note is an audio file
	TypeAudio
	// TypeReminder indicates that a note is a reminder
	TypeReminder
	// TypeList indicates that a note is a list
	TypeList
)

// Scope is the scope of the note (account or user)
type Scope int

const (
	// ScopeUser indicates that a note belongs to a single user
	ScopeUser Scope = iota
	// ScopeAccount indicates that a note belongs to a whole account
	ScopeAccount
)

// StoreType indicates the type of store that holds the note
type StoreType int

const (
	// StoreTypeShelf indicates that a note is stored in a shelf
	StoreTypeShelf StoreType = iota
	// StoreTypeCollection indicates that a note is stored in a collection
	StoreTypeCollection
)

// Note is the primary content type
type Note struct {
	ID            uuid.UUID      `json:"id"`             // ID is the unique identifier for this note
	OwnerID       uuid.UUID      `json:"owner_id"`       // OwnerID is the owner (an account or user id)
	Scope         Scope          `json:"scope"`          // Scope indicates whether the owner is an account or user
	NotebookID    uuid.UUID      `json:"notebook_id"`    // NotebookID is the id of the owning notebook
	StoreID       uuid.UUID      `json:"store_id"`       // StoreID is the id of the data store
	StoreType     StoreType      `json:"store_type"`     // StoreType indicates whether the data store is a shelf or collection
	Title         *title.Title   `json:"title"`          // Title is the title of the note
	Type          Type           `json:"type"`           // Type is one of the NoteType* identifier values
	Content       string         `json:"content"`        // Content is the content of the note
	Tags          []*tag.Tag     `json:"tags"`           // Tags is the set of tags assigned to the note
	Revisions     []*Note        `json:"-"`              // Revisions is the set of previously saved note revisions
	RevisionCount int            `json:"revision_count"` // RevisionCount keeps track of the number of saved note revisions
	Created       time.Time      `json:"created"`        // Created is the time when the note was created
	Updated       time.Time      `json:"updated"`        // Updated is the time when note was last updated
	Locked        bool           `json:"locked"`         // Locked indicates whether the note can be modified
	TemplateID    uuid.UUID      `json:"template_id"`    // TemplateID indicates the ID of a template (if the note was created from a template)
	DBFactory     *db.Factory    `json:"-"`
	Logger        *logrus.Logger `json:"-"`
}

// New creates a new note object
func New(title *title.Title, scope Scope, store StoreType, dbFactory *db.Factory, logger *logrus.Logger) *Note {
	now := time.Now()
	note := &Note{
		ID:            uuid.NewV4(),
		Scope:         scope,
		StoreType:     store,
		Title:         title,
		Created:       now,
		Updated:       now,
		Locked:        false,
		RevisionCount: 0,
		DBFactory:     dbFactory,
		Logger:        logger,
	}
	return note
}

func (note *Note) getDB() *db.DB {
	var noteDB *db.DB
	if note.StoreType == StoreTypeCollection {
		noteDB = note.DBFactory.Find(db.TypeCollection, note.StoreID)
	} else {
		noteDB = note.DBFactory.Find(db.TypeShelf, note.StoreID)
	}
	// [FIXME] - open if db nil
	return noteDB
}

// Save a note
func (note *Note) Save(passphraseKey []byte) error {
	noteDB := note.getDB()
	err := noteDB.DB.Update(func(tx *bolt.Tx) error {
		// get bucket, creating it if needed
		// [FIXME] - notes are grouped into unique buckets by notebook id
		bucket, err := tx.CreateBucketIfNotExists([]byte("notes"))
		if err != nil {
			note.Logger.Debug("Error creating notes bucket - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorCreateBucket)
			return code
		}

		// serialize note data
		data, err := json.Marshal(note)
		if err != nil {
			note.Logger.Debug("Error marshaling note - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		c := crypto.New(note.Logger)
		decryptedKey, err := c.Open(passphraseKey, noteDB.EncryptedKey)
		if err != nil {
			note.Logger.Debug("Error retrieving note key - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := c.Seal(decryptedKey, data)
		if err != nil {
			note.Logger.Debug("Error encrypting note data - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(note.ID.Bytes(), encryptedData)
		if err != nil {
			note.Logger.Debug("Error writing note - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		note.Logger.Debug("Error saving note - err")
		code := codes.New(codes.ScopeNote, codes.ErrorSave)
		return code
	}

	return nil
}

// LoadAll notes
func (note *Note) LoadAll(passphraseKey []byte) ([]*Note, error) {
	var notes []*Note

	noteDB := note.getDB()
	c := crypto.New(note.Logger)
	noteKey, err := c.Open(passphraseKey, noteDB.EncryptedKey)
	if err != nil {
		note.Logger.Debug("Error opening note key - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorOpenKey)
		return notes, code
	}

	err = noteDB.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			note.Logger.Debug("note bucket does not exist")
			code := codes.New(codes.ScopeNote, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newNote := &Note{
				DBFactory: note.DBFactory,
				Logger:    note.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(noteKey, value)
			if err != nil {
				note.Logger.Debug("Error decrypting note data - ", err)
				code := codes.New(codes.ScopeNote, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newNote)
			if err != nil {
				note.Logger.Debug("Error decoding note json - ", err)
				code := codes.New(codes.ScopeNote, codes.ErrorDecode)
				return code
			}

			notes = append(notes, newNote)
		}

		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return nil, err
		}
		note.Logger.Debug("Error loading all notes - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorLoadAll)
		return nil, code
	}

	return notes, nil
}

// Load a note
func (note *Note) Load(passphraseKey []byte) error {
	noteDB := note.getDB()
	c := crypto.New(note.Logger)
	noteKey, err := c.Open(passphraseKey, noteDB.EncryptedKey)
	if err != nil {
		note.Logger.Debug("Error opening note key - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorOpenKey)
		return code
	}

	err = noteDB.DB.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			note.Logger.Debug("note bucket does not exist")
			code := codes.New(codes.ScopeNote, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		key, value := cursor.Seek(note.ID.Bytes())
		if key == nil {
			note.Logger.Debug("Error loading note")
			code := codes.New(codes.ScopeNote, codes.ErrorLoad)
			return code
		}

		// decrypt value
		decryptedData, err := c.Open(noteKey, value)
		if err != nil {
			note.Logger.Debug("Error decrypting note data - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorDecrypt)
			return code
		}

		err = json.Unmarshal(decryptedData, note)
		if err != nil {
			note.Logger.Debug("Error decoding note json - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorDecode)
			return code
		}

		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		note.Logger.Debug("Error loading note - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorLoad)
		return code
	}

	return nil
}

// Delete a note
func (note *Note) Delete(passphraseKey []byte) error {
	noteDB := note.getDB()
	err := noteDB.DB.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			note.Logger.Debug("note bucket does not exist")
			code := codes.New(codes.ScopeNote, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(note.ID.Bytes())
		if err != nil {
			note.Logger.Debug("Error deleting note - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		note.Logger.Debug("Error deleting note - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorDelete)
		return code
	}

	return nil
}
