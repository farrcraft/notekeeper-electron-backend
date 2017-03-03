package main

import ()

const (
	LIST_ARRANGE_UP  = "up"
	LIST_ARRAGE_DOWN = "down"
)

type ListEntry struct {
	Content string
	Checked bool
}

type List struct {
	Entries    []*ListEntry
	Lists      []*List
	AutoArrage bool
	Direction  string // when auto arranging, checked entries are moved to the top or the bottom?
}
