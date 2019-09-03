package collection

import (
	"time"

	"notekeeper-electron-backend/db"
	"notekeeper-electron-backend/notebook"
	"notekeeper-electron-backend/tag"
	"notekeeper-electron-backend/title"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
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
	ID           uuid.UUID            `json:"id"`             // ID is the unique collection identifier
	Title        *title.Title         `json:"title"`          // Title is a title for the collection
	Notebooks    []*notebook.Notebook `json:"-"`              // Notebooks is the set of notebooks in the collection
	OwnerID      uuid.UUID            `json:"owner_id"`       // OwnerID is the account or user that the collection belongs to
	EncryptedKey []byte               `json:"encryption_key"` // EncryptedKey is the encrypted encryption key for the collection DB
	Scope        Scope                `json:"scope"`          // Scope is the ownership scope of the collection (user or account)
	ShelfID      uuid.UUID            `json:"shelf_id"`       // ShelfID is the shelf that contains the collection
	Tags         []*tag.Tag           `json:"tags"`           // Tags is the set of tags assigned to the collection
	Created      time.Time            `json:"created"`        // Created is the time when the collection was first created
	Updated      time.Time            `json:"updated"`        // Updated is the time when the collection was last updated
	Locked       bool                 `json:"locked"`         // Locked indicates whether the collection can be modified
	DBRegistry   *db.Registry         `json:"-"`              // DBRegistry provides database access
	Logger       *logrus.Logger       `json:"-"`              // Logger is the logging facility
}

// New creates a new collection object
func New(title *title.Title, scope Scope, dbRegistry *db.Registry, logger *logrus.Logger) (*Collection, error) {
	now := time.Now()

	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	collection := &Collection{
		ID:         id,
		Scope:      scope,
		Title:      title,
		Created:    now,
		Updated:    now,
		Locked:     false,
		DBRegistry: dbRegistry,
		Logger:     logger,
	}

	return collection, nil
}
