package collection

import (
	"time"

	"../db"
	"../notebook"
	"../tag"
	"../title"
	"github.com/Sirupsen/logrus"
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
