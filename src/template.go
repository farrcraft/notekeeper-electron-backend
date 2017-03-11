package main

import (
	"time"

	"./tag"
	"./title"
	uuid "github.com/satori/go.uuid"
)

// Template is a template note structure that can be used for creating actual notes
type Template struct {
	ID        uuid.UUID    `json:"id"`      // ID is the unique identifier for the template
	Title     *title.Title `json:"title"`   // Title is the title of the note template
	Type      string       `json:"type"`    // Type is the type of the note
	Content   string       `json:"content"` // Content is the default note content
	Tags      []*tag.Tag   `json:"tags"`    // Tags is the default set of tags for the note
	Revisions []*Template  `json:"-"`       // Revisions is the set of previously saved template revisions
	Created   time.Time    `json:"created"` // Created is the time when the template was created
	Updated   time.Time    `json:"updated"` // Updated is the time when the template was last updated
	Locked    bool         `json:"locked"`  // Locked indicates whether the template can be modified
}

// NewTemplate creates a new note object
func NewTemplate() *Template {
	now := time.Now()
	template := &Template{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return template
}
