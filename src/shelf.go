package main

import (
	"time"

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
	ID        uuid.UUID   `json:"id"`      // ID is the unique identifier for the shelf
	Title     *Title      `json:"title"`   // Title is the title of the shelf
	Type      string      `json:"type"`    // Type is one of the ShelfType* identifier values
	Default   bool        `json:"default"` // Default indicates whether this is the default shelf
	Trash     bool        `json:"trash"`   // Trash indicates whether this is the trash shelf
	Account   *Account    `json:"-"`       // Account is the account that the shelf belongs to
	User      *User       `json:"-"`       // User is the user that the shelf belongs to
	Notebooks []*Notebook `json:"-"`       // Notebooks is the set of notebooks in the shelf
	Tags      []*Tag      `json:"tags"`    // Tags is the set of tags assigned to the shelf
	Created   time.Time   `json:"created"` // Created is the time when the shelf was created
	Updated   time.Time   `json:"updated"` // Updated is the time when the shelf was last updated
	Locked    bool        `json:"locked"`  // Locked indicates whether the shelf can be modified
}

// NewShelf creates a new shelf object
func NewShelf() *Shelf {
	now := time.Now()
	shelf := &Shelf{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return shelf
}
