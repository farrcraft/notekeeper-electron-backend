package main

import ()

const (
	NOTE_TYPE_PLAIN_TEXT = "plain"
	NOTE_TYPE_RICH_TEXT  = "rich"
	NOTE_TYPE_MARKDOWN   = "markdown"
	NOTE_TYPE_HTML       = "html"
	NOTE_TYPE_IMAGE      = "image"
	NOTE_TYPE_FILE       = "file"
	NOTE_TYPE_PDF        = "pdf"
	NOTE_TYPE_AUDIO      = "audio"
	NOTE_TYPE_REMINDER   = "reminder"
	NOTE_TYPE_LIST       = "list"
)

type Template struct {
	Id        string
	Title     *Title
	Type      string
	Content   string
	Tags      []*Tags
	Revisions []*Template
	Created   string
	Updated   string
	Locked    bool
}

type Note struct {
	Id         string
	Notebook   *Notebook
	Title      *Title
	Type       string
	Content    string
	Tags       []*Tag
	Revisions  []*Note
	Created    string
	Updated    string
	Locked     bool
	TemplateId string // only if note was created from a template
}

func NewNote() *Note {

}
