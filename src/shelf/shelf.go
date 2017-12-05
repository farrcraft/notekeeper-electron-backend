package shelf

import (
	"time"

	"../collection"
	"../db"
	"../notebook"
	"../tag"
	"../title"
	"github.com/Sirupsen/logrus"
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
	OwnerID      uuid.UUID                `json:"-"`              // OwnerID is the ID of the account or user owning the shelf
	EncryptedKey []byte                   `json:"encryption_key"` // EncryptedKey is the encrypted encryption key for the shelf DB
	Notebooks    []*notebook.Notebook     `json:"-"`              // Notebooks is the set of notebooks in the shelf
	Collections  []*collection.Collection `json:"-"`              // Collections is the set of collections in the shelf
	Tags         []*tag.Tag               `json:"tags"`           // Tags is the set of tags assigned to the shelf
	Created      time.Time                `json:"created"`        // Created is the time when the shelf was created
	Updated      time.Time                `json:"updated"`        // Updated is the time when the shelf was last updated
	Locked       bool                     `json:"locked"`         // Locked indicates whether the shelf can be modified
	DBRegistry   *db.Registry             `json:"-"`              // DBRegistry provides access to the database
	Logger       *logrus.Logger           `json:"-"`              // Logger is the logging facility
}

// New creates a new shelf object
func New(title *title.Title, scope Scope, dbRegistry *db.Registry, logger *logrus.Logger) *Shelf {
	now := time.Now()
	shelf := &Shelf{
		ID:         uuid.NewV4(),
		Title:      title,
		Scope:      scope,
		Created:    now,
		Updated:    now,
		Default:    false,
		Trash:      false,
		Locked:     false,
		DBRegistry: dbRegistry,
		Logger:     logger,
	}
	return shelf
}
