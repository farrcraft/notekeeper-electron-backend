package tag

import (
	"time"

	"../title"
	uuid "github.com/satori/go.uuid"
)

// Scope indicates the scope of the tag
type Scope int

const (
	// ScopeUser indicates that a tag belongs to a single user
	ScopeUser Scope = iota
	// ScopeAccount indicates that a tag belongs to a whole account
	ScopeAccount
)

// Tag is used for assigning labels to various object types
type Tag struct {
	ID      uuid.UUID    `json:"id"`    // ID is the unique identifier of the tag
	Title   *title.Title `json:"title"` // Title is the title of the tag
	Created time.Time    `json:"created"`
	Updated time.Time    `json:"updated"`
	Scope   Scope        `json:"scope"`
}

// New creates a new tag object
func New() *Tag {
	now := time.Now()
	tag := &Tag{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return tag
}
