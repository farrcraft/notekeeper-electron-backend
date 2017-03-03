package main

import ()

type TitleFormatting struct {
	Bold       bool
	Italics    bool
	Underscore bool
	Strike     bool
	Background string
	Color      string
}

type Title struct {
	Title      string
	Formatting *TitleFormatting
}
