package rpc

import (
	"time"

	"notekeeper-electron-backend/codes"
	messages "notekeeper-electron-backend/proto"
	"notekeeper-electron-backend/title"
)

// TitleToMessage converts a title domain instance into a protobuf instance
func TitleToMessage(t *title.Title) *messages.Title {
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

// TimeToMessage converts a native time to a consistent string representation
func TimeToMessage(t time.Time) string {
	s := t.Format(time.RFC3339)
	return s
}

// MessageToTitle converts a protobuf title message to a native title
func MessageToTitle(msg *messages.Title) *title.Title {
	t := title.New(msg.Text)
	t.Formatting.Bold = msg.Bold
	t.Formatting.Italics = msg.Italics
	t.Formatting.Underscore = msg.Underscore
	t.Formatting.Strike = msg.Strike
	t.Formatting.Color = msg.Color
	t.Formatting.Background = msg.Background
	return t
}

// SetInternalError sets an error in a response header
func SetInternalError(header *messages.ResponseHeader, err error) {
	code := codes.ToInternalError(err)
	header.Code = int32(code.Code)
	header.Scope = int32(code.Scope)
	header.Status = code.Error()
}

// SetRPCError sets an rpc-specific error in a response header
func SetRPCError(header *messages.ResponseHeader, c codes.Code) {
	header.Code = int32(c)
	header.Scope = int32(codes.ScopeRPC)
	header.Status = codes.StatusSystemError
}

// NewResponseHeader creates a new response header
func NewResponseHeader() *messages.ResponseHeader {
	header := &messages.ResponseHeader{
		Code:   int32(codes.ErrorOK),
		Status: codes.StatusOK,
	}
	return header
}
