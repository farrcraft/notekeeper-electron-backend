package shelf

import (
	"time"

	"../collection"
	"../notebook"
	"../tag"
	"../title"
	uuid "github.com/satori/go.uuid"
)

const (
	// ShelfTypeUser indicates that a shelf belongs to a single user
	ShelfTypeUser = "user"
	// ShelfTypeAccount indicates that a shelf belongs to a whole account
	ShelfTypeAccount = "account"
)

// Shelf contains a collection of notebooks
type Shelf struct {
	ID          uuid.UUID                `json:"id"`      // ID is the unique identifier for the shelf
	Title       *title.Title             `json:"title"`   // Title is the title of the shelf
	Type        string                   `json:"type"`    // Type is one of the ShelfType* identifier values
	Default     bool                     `json:"default"` // Default indicates whether this is the default shelf
	Trash       bool                     `json:"trash"`   // Trash indicates whether this is the trash shelf
	UserID      uuid.UUID                `json:"user_id"` // UserID is the ID of the user owning the shelf
	Notebooks   []*notebook.Notebook     `json:"-"`       // Notebooks is the set of notebooks in the shelf
	Collections []*collection.Collection `json:"-"`       // Collections is the set of collections in the shelf
	Tags        []*tag.Tag               `json:"tags"`    // Tags is the set of tags assigned to the shelf
	Created     time.Time                `json:"created"` // Created is the time when the shelf was created
	Updated     time.Time                `json:"updated"` // Updated is the time when the shelf was last updated
	Locked      bool                     `json:"locked"`  // Locked indicates whether the shelf can be modified
}

// New creates a new shelf object
func New() *Shelf {
	now := time.Now()
	shelf := &Shelf{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return shelf
}
