package rpc

import (
	"time"

	"../codes"
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

func messageToTitle(msg *messages.Title) *title.Title {
	t := title.New(msg.Text)
	t.Formatting.Bold = msg.Bold
	t.Formatting.Italics = msg.Italics
	t.Formatting.Underscore = msg.Underscore
	t.Formatting.Strike = msg.Strike
	t.Formatting.Color = msg.Color
	t.Formatting.Background = msg.Background
	return t
}

func setInternalError(header *messages.ResponseHeader, err error) {
	code := codes.ToInternalError(err)
	header.Code = int32(code.Code)
	header.Scope = int32(code.Scope)
	header.Status = code.Error()
}

func setRPCError(header *messages.ResponseHeader, c codes.Code) {
	header.Code = int32(c)
	header.Scope = int32(codes.ScopeRPC)
	header.Status = codes.StatusError
}

func newResponseHeader() *messages.ResponseHeader {
	header := &messages.ResponseHeader{
		Code:   int32(codes.ErrorOK),
		Status: codes.StatusOK,
	}
	return header
}
