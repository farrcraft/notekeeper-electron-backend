package note

import (
	"encoding/json"
	"time"

	"notekeeper-electron-backend/codes"
	"notekeeper-electron-backend/crypto"
	"notekeeper-electron-backend/db"
	"notekeeper-electron-backend/tag"
	"notekeeper-electron-backend/title"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"go.etcd.io/bbolt"
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
	DBRegistry    *db.Registry   `json:"-"`
	Logger        *logrus.Logger `json:"-"`
}

// New creates a new note object
func New(title *title.Title, scope Scope, store StoreType, dbRegistry *db.Registry, logger *logrus.Logger) (*Note, error) {
	now := time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	note := &Note{
		ID:            id,
		Scope:         scope,
		StoreType:     store,
		Title:         title,
		Created:       now,
		Updated:       now,
		Locked:        false,
		RevisionCount: 0,
		DBRegistry:    dbRegistry,
		Logger:        logger,
	}

	return note, nil
}

func (note *Note) getDBHandle() (*db.Handle, error) {
	var key db.Key
	key.ID = note.StoreID
	if note.StoreType == StoreTypeCollection {
		key.Type = db.TypeCollection
	} else {
		key.Type = db.TypeShelf
	}
	handle, err := note.DBRegistry.GetHandle(key)
	return handle, err
}

// Save a note
func (note *Note) Save(passphraseKey []byte) error {
	noteDBHandle, err := note.getDBHandle()
	if err != nil {
		return err
	}
	err = noteDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		// get bucket, creating it if needed
		// [FIXME] - notes are grouped into unique buckets by notebook id
		bucket, err := tx.CreateBucketIfNotExists([]byte("notes"))
		if err != nil {
			note.Logger.Warn("Error creating notes bucket - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorCreateBucket)
			return code
		}

		// serialize note data
		data, err := json.Marshal(note)
		if err != nil {
			note.Logger.Warn("Error marshaling note - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorMarshal)
			return code
		}

		// retrieve the encryption key
		c := crypto.New(note.Logger)
		decryptedKey, err := c.Open(passphraseKey, noteDBHandle.EncryptedKey)
		if err != nil {
			note.Logger.Warn("Error retrieving note key - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorOpenKey)
			return code
		}

		// encrypt the data
		encryptedData, err := c.Seal(decryptedKey, data)
		if err != nil {
			note.Logger.Warn("Error encrypting note data - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorEncrypt)
			return code
		}

		// finally, save it
		err = bucket.Put(note.ID.Bytes(), encryptedData)
		if err != nil {
			note.Logger.Warn("Error writing note - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorWriteBucket)
			return code
		}
		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		note.Logger.Warn("Error saving note - err")
		code := codes.New(codes.ScopeNote, codes.ErrorSave)
		return code
	}

	return nil
}

// LoadAll notes
func (note *Note) LoadAll(passphraseKey []byte) ([]*Note, error) {
	var notes []*Note

	noteDBHandle, err := note.getDBHandle()
	if err != nil {
		return nil, err
	}
	c := crypto.New(note.Logger)
	noteKey, err := c.Open(passphraseKey, noteDBHandle.EncryptedKey)
	if err != nil {
		note.Logger.Warn("Error opening note key - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorOpenKey)
		return notes, code
	}

	err = noteDBHandle.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			note.Logger.Warn("note bucket does not exist")
			code := codes.New(codes.ScopeNote, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		for key, value := cursor.First(); key != nil; key, value = cursor.Next() {
			newNote := &Note{
				DBRegistry: note.DBRegistry,
				Logger:     note.Logger,
			}

			// decrypt value
			decryptedData, err := c.Open(noteKey, value)
			if err != nil {
				note.Logger.Warn("Error decrypting note data - ", err)
				code := codes.New(codes.ScopeNote, codes.ErrorDecrypt)
				return code
			}

			err = json.Unmarshal(decryptedData, newNote)
			if err != nil {
				note.Logger.Warn("Error decoding note json - ", err)
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
		note.Logger.Warn("Error loading all notes - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorLoadAll)
		return nil, code
	}

	return notes, nil
}

// Load a note
func (note *Note) Load(passphraseKey []byte) error {
	noteDBHandle, err := note.getDBHandle()
	if err != nil {
		return err
	}
	c := crypto.New(note.Logger)
	noteKey, err := c.Open(passphraseKey, noteDBHandle.EncryptedKey)
	if err != nil {
		note.Logger.Warn("Error opening note key - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorOpenKey)
		return code
	}

	err = noteDBHandle.DB.View(func(tx *bbolt.Tx) error {
		// Assume bucket exists and has keys
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			note.Logger.Warn("note bucket does not exist")
			code := codes.New(codes.ScopeNote, codes.ErrorBucketMissing)
			return code
		}

		cursor := bucket.Cursor()

		key, value := cursor.Seek(note.ID.Bytes())
		if key == nil {
			note.Logger.Warn("Error loading note")
			code := codes.New(codes.ScopeNote, codes.ErrorLoad)
			return code
		}

		// decrypt value
		decryptedData, err := c.Open(noteKey, value)
		if err != nil {
			note.Logger.Warn("Error decrypting note data - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorDecrypt)
			return code
		}

		err = json.Unmarshal(decryptedData, note)
		if err != nil {
			note.Logger.Warn("Error decoding note json - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorDecode)
			return code
		}

		return nil
	})
	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		note.Logger.Warn("Error loading note - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorLoad)
		return code
	}

	return nil
}

// Delete a note
func (note *Note) Delete(passphraseKey []byte) error {
	noteDBHandle, err := note.getDBHandle()
	if err != nil {
		return err
	}
	err = noteDBHandle.DB.Update(func(tx *bbolt.Tx) error {
		bucket := tx.Bucket([]byte("notes"))
		if bucket == nil {
			note.Logger.Warn("note bucket does not exist")
			code := codes.New(codes.ScopeNote, codes.ErrorBucketMissing)
			return code
		}

		err := bucket.Delete(note.ID.Bytes())
		if err != nil {
			note.Logger.Warn("Error deleting note - ", err)
			code := codes.New(codes.ScopeNote, codes.ErrorDelete)
			return code
		}

		return nil
	})

	if err != nil {
		if codes.IsInternalError(err) {
			return err
		}
		note.Logger.Warn("Error deleting note - ", err)
		code := codes.New(codes.ScopeNote, codes.ErrorDelete)
		return code
	}

	return nil
}
