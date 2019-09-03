package title

import (
	"errors"
	"regexp"
)

// Formatting contains formatting options for a title
type Formatting struct {
	Bold       bool   `json:"bold"`       // Bold indicates that the title text is bolded
	Italics    bool   `json:"italics"`    // Italics indicates that the title text is in italics
	Underscore bool   `json:"underscore"` // Underscore indicates that the title text is underscored
	Strike     bool   `json:"strike"`     // Strike indicates that the title text contains a strikethrough line
	Background string `json:"background"` // Background indicates the background color of the text
	Color      string `json:"color"`      // Color indicates the color of the text
}

// Format applies formatting options
func (formatting *Formatting) Format(options int) {
	if options&FormatBold != 0 {
		formatting.Bold = true
	} else {
		formatting.Bold = false
	}

	if options&FormatItalics != 0 {
		formatting.Italics = true
	} else {
		formatting.Italics = false
	}

	if options&FormatStrike != 0 {
		formatting.Strike = true
	} else {
		formatting.Strike = false
	}

	if options&FormatUnderscore != 0 {
		formatting.Underscore = true
	} else {
		formatting.Underscore = false
	}

	if options&FormatDefault != 0 {
		formatting.Bold = false
		formatting.Italics = false
		formatting.Underscore = false
		formatting.Strike = false
		formatting.Background = ""
		formatting.Color = ""
	}
}

// SetColor sets the foreground or background color
func (formatting *Formatting) SetColor(color string, background bool) error {
	ok, err := formatting.ValidateColor(color)
	if err != nil {
		return err
	}
	if !ok {
		err = errors.New("Invalid color")
		return err
	}

	if background {
		formatting.Background = color
	} else {
		formatting.Color = color
	}

	return nil
}

// ValidateColor validates that a string is a valid hex color string
func (formatting *Formatting) ValidateColor(color string) (bool, error) {
	regex, err := regexp.Compile(`(?i)^#([0-9a-f]{3}){1,2}$`)
	if err != nil {
		return false, err
	}
	if regex.MatchString(color) {
		return true, nil
	}
	return false, nil
}
