package note

import (
	"time"

	"../tag"
	"../title"
	uuid "github.com/satori/go.uuid"
)

// Type is the type of note
type Type int

const (
	// TypePlainText indicates that a note is just plain old text
	TypePlainText Type = iota
	// TypeRichText indicates that a note contains rich text
	TypeRichText
	// TypeMarkdown indicates that a note contains markdown text
	TypeMarkdown
	// TypeHTML indicates that a note contains HTML
	TypeHTML
	// TypeImage indicates that a note is an image file
	TypeImage
	// TypeFile indicates that a note is an arbitrary file attachment
	TypeFile
	// TypePdf indicates that a note is a PDF file
	TypePdf
	// TypeAudio indicates that a note is an audio file
	TypeAudio
	// TypeReminder indicates that a note is a reminder
	TypeReminder
	// TypeList indicates that a note is a list
	TypeList
)

// Note is the primary content type
type Note struct {
	ID uuid.UUID `json:"id"` // ID is the unique identifier for this note
	//Notebook   *Notebook `json:"-"`           // Notebook is the notebook this note belongs to
	Title         *title.Title `json:"title"`          // Title is the title of the note
	Type          Type         `json:"type"`           // Type is one of the NoteType* identifier values
	Content       string       `json:"content"`        // Content is the content of the note
	Tags          []*tag.Tag   `json:"tags"`           // Tags is the set of tags assigned to the note
	Revisions     []*Note      `json:"-"`              // Revisions is the set of previously saved note revisions
	RevisionCount int          `json:"revision_count"` // RevisionCount keeps track of the number of saved note revisions
	Created       time.Time    `json:"created"`        // Created is the time when the note was created
	Updated       time.Time    `json:"updated"`        // Updated is the time when note was last updated
	Locked        bool         `json:"locked"`         // Locked indicates whether the note can be modified
	TemplateID    uuid.UUID    `json:"template_id"`    // TemplateID indicates the ID of a template (if the note was created from a template)
}

// NewNote creates a new note object
func NewNote() *Note {
	now := time.Now()
	note := &Note{
		ID:            uuid.NewV4(),
		Created:       now,
		Updated:       now,
		Locked:        false,
		RevisionCount: 0,
	}
	return note
}
