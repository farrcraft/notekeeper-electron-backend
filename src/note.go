package main

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	// NoteTypePlainText indicates that a note is just plain old text
	NoteTypePlainText = "plain"
	// NoteTypeRichText indicates that a note contains rich text
	NoteTypeRichText = "rich"
	// NoteTypeMarkdown indicates that a note contains markdown text
	NoteTypeMarkdown = "markdown"
	// NoteTypeHTML indicates that a note contains HTML
	NoteTypeHTML = "html"
	// NoteTypeImage indicates that a note is an image file
	NoteTypeImage = "image"
	// NoteTypeFile indicates that a note is an arbitrary file attachment
	NoteTypeFile = "file"
	// NoteTypePdf indicates that a note is a PDF file
	NoteTypePdf = "pdf"
	// NoteTypeAudio indicates that a note is an audio file
	NoteTypeAudio = "audio"
	// NoteTypeReminder indicates that a note is a reminder
	NoteTypeReminder = "reminder"
	// NoteTypeList indicates that a note is a list
	NoteTypeList = "list"
)

// Note is the primary content type
type Note struct {
	ID         uuid.UUID `json:"id"`          // ID is the unique identifier for this note
	Notebook   *Notebook `json:"-"`           // Notebook is the notebook this note belongs to
	Title      *Title    `json:"title"`       // Title is the title of the note
	Type       string    `json:"type"`        // Type is one of the NoteType* identifier values
	Content    string    `json:"content"`     // Content is the content of the note
	Tags       []*Tag    `json:"tags"`        // Tags is the set of tags assigned to the note
	Revisions  []*Note   `json:"-"`           // Revisions is the set of previously saved note revisions
	Created    time.Time `json:"created"`     // Created is the time when the note was created
	Updated    time.Time `json:"updated"`     // Updated is the time when note was last updated
	Locked     bool      `json:"locked"`      // Locked indicates whether the note can be modified
	TemplateID string    `json:"template_id"` // TemplateID indicates the ID of a template (if the note was created from a template)
}

// NewNote creates a new note object
func NewNote() *Note {
	now := time.Now()
	note := &Note{
		ID:      uuid.NewV4(),
		Created: now,
		Updated: now,
	}
	return note
}
