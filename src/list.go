package main

import (
	"time"
)

const (
	// ListArrangeUp indicates that checked list entries are arranged to the top
	ListArrangeUp = "up"
	// ListArrangeDown indicates that checked list entries are arranged to the  bottom
	ListArrangeDown = "down"
)

// ListEntry is a single item in a list
type ListEntry struct {
	Content string    `json:"content"` // Content is the text of the list
	Created time.Time `json:"created"` // Created is the time when the entry was created
	Updated time.Time `json:"updated"` // Updated is the time when the entry was last updated
	Checked bool      `json:"checked"` // Checked is the current state of the entry
}

// List is a checklist which may contain other nested lists
type List struct {
	Entries    []*ListEntry `json:"entries"`      // Entries is the set of entries in the list
	Lists      []*List      `json:"lists"`        // Lists is the set of nested lists in this list
	AutoArrage bool         `json:"auto_arrange"` // AutoArrange indicates that checked entries should be automatically sorted in the list
	Direction  string       `json:"direction"`    // Direction indicates when auto arranging wether checked entries are moved to the top or the bottom
	Created    time.Time    `json:"created"`      // Created is the time when the list was created
	Updated    time.Time    `json:"updated"`      // Updated is the time when the list was last updated
}

// NewListEntry creates a new list entry object
func NewListEntry() *ListEntry {
	now := time.Now()
	entry := &ListEntry{
		Created: now,
		Updated: now,
	}
	return entry
}

// NewList creates a new list object
func NewList() *List {
	now := time.Now()
	list := &List{
		Created: now,
		Updated: now,
	}
	return list
}
