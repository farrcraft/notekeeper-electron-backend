package title

// Formatting options
const (
	FormatDefault = 1 << iota
	FormatBold
	FormatItalics
	FormatUnderscore
	FormatStrike
)

// Title represents the title of a piece of content
type Title struct {
	Title      string      `json:"title"`      // Title is the text of the title
	Formatting *Formatting `json:"formatting"` // Formatting contains the formatting options for the title text
}

// New creates a new Title object
func New(text string) *Title {
	title := &Title{
		Title:      text,
		Formatting: &Formatting{},
	}
	return title
}

// Format applies formatting options to the title
func (title *Title) Format(options int) {
	title.Formatting.Format(options)
}
