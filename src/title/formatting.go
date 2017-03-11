package title

// Formatting contains formatting options for a title
type Formatting struct {
	Bold       bool   `json:"bold"`       // Bold indicates that the title text is bolded
	Italics    bool   `json:"italics"`    // Italics indicates that the title text is in italics
	Underscore bool   `json:"underscore"` // Underscore indicates that the title text is underscored
	Strike     bool   `json:"strike"`     // Strike indicates that the title text contains a strikethrough line
	Background string `json:"background"` // Background indicates the background color of the text
	Color      string `json:"color"`      // Color indicates the color of the text
}
