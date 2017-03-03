package main

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

// Tag is used for assigning labels to various object types
type Tag struct {
	ID      uuid.UUID `json:"id"`    // ID is the unique identifier of the tab
	Title   *Title    `json:"title"` // Title is the title of the tag
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

// NewTag creates a new tag object
func NewTag() *Tag {
	now := time.Now()
	tag := &Tag{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return tag
}
