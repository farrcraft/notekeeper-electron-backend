package codes

import (
	"strconv"
)

// ErrorCode is the error code type
type ErrorCode int32

// These are the status types that can be passed to the front end
const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

// These are the error codes that can be passed to the front end
const (
	ErrorUnknown ErrorCode = iota + 1
	ErrorMasterDbOpen
	ErrorCreateUIState
	ErrorLoadUIState
	ErrorSaveUIState
)

// String converts error code to a string
func (e ErrorCode) String() string {
	s := strconv.Itoa(int(e))
	return s
}
