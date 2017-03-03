package main

// TitleFormatting contains formatting options for a title
type TitleFormatting struct {
	Bold       bool   `json:"bold"`       // Bold indicates that the title text is bolded
	Italics    bool   `json:"italics"`    // Italics indicates that the title text is in italics
	Underscore bool   `json:"underscore"` // Underscore indicates that the title text is underscored
	Strike     bool   `json:"strike"`     // Strike indicates that the title text contains a strikethrough line
	Background string `json:"background"` // Background indicates the background color of the text
	Color      string `json:"color"`      // Color indicates the color of the text
}

// Title represents the title of a piece of content
type Title struct {
	Title      string           // Title is the text of the title
	Formatting *TitleFormatting // Formatting contains the formatting options for the title text
}
