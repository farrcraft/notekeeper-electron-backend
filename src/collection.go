package main

import (
	"time"

	"./account"
	"./notebook"
	"./tag"
	"./title"
	uuid "github.com/satori/go.uuid"
)

// Collection holds a collection of notebooks
type Collection struct {
	ID        uuid.UUID            `json:"id"`      // ID is the unique collection identifier
	Title     *title.Title         `json:"title"`   // Title is a title for the collection
	Notebooks []*notebook.Notebook `json:"-"`       // Notebooks is the set of notebooks in the collection
	Account   *account.Account     `json:"-"`       // Account is the account that the collection belongs to
	User      *account.User        `json:"-"`       // User is the individual account user that the collection belongs to
	Shelf     *account.Shelf       `json:"-"`       // Shelf is the shelf that contains the collection
	Tags      []*tag.Tag           `json:"tags"`    // Tags is the set of tags assigned to the collection
	Created   time.Time            `json:"created"` // Created is the time when the collection was first created
	Updated   time.Time            `json:"updated"` // Updated is the time when the collection was last updated
	Locked    bool                 `json:"locked"`  // Locked indicates whether the collection can be modified
}

// NewCollection creates a new collection object
func NewCollection() *Collection {
	now := time.Now()
	collection := &Collection{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return collection
}
