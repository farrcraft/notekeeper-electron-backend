package rpc

import (
	"time"

	messages "../proto"
	"../title"
)

// titleToMessage converts a title domain instance into a protobuf instance
func titleToMessage(t *title.Title) *messages.Title {
	m := &messages.Title{
		Text:       t.Title,
		Bold:       t.Formatting.Bold,
		Italics:    t.Formatting.Italics,
		Underscore: t.Formatting.Underscore,
		Strike:     t.Formatting.Strike,
		Color:      t.Formatting.Color,
		Background: t.Formatting.Background,
	}
	return m
}

// timeToMessage converts a native time to a consistent string representation
func timeToMessage(t time.Time) string {
	s := t.Format(time.RFC3339)
	return s
}
