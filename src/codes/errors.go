package codes

import (
	"fmt"
	"strconv"
)

// Code is the error code type
type Code int32

// Scope of an error code
type Scope int32

// InternalError is a custom error type
type InternalError struct {
	Scope   Scope
	Code    Code
	Message string
}

// These are the status types that can be passed to the front end
const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

// Valid error scopes
const (
	ScopeGeneral Scope = iota
	ScopeAccount
	ScopeAPI
	ScopeCollection
	ScopeCrypto
	ScopeDB
	ScopeNote
	ScopeNotebook
	ScopeRPC
	ScopeShelf
	ScopeTag
	ScopeTemplate
	ScopeTitle
	ScopeUIState
	ScopeUser
)

// These are the error codes that can be passed to the front end
const (
	ErrorOK Code = iota
	ErrorUnknown
	ErrorInternalEscape // This means a non-internal error tried to escape the RPC boundary

	ErrorUnauthorized
	ErrorInvalidType
	ErrorMissingDB
	ErrorDbOpen
	ErrorCreateBucket
	ErrorMarshal
	ErrorOpenKey
	ErrorEncrypt
	ErrorDecrypt
	ErrorCrypto
	ErrorWriteBucket
	ErrorSave
	ErrorBucketMissing
	ErrorDecode
	ErrorDeriveKey
	ErrorConvertID
	ErrorLookup
	ErrorLoad
	ErrorLoadAll
	ErrorDelete
	ErrorCreate
)

// String converts error code to a string
func (e Code) String() string {
	s := strconv.Itoa(int(e))
	return s
}

func (scope Scope) String() string {
	s := strconv.Itoa(int(scope))
	return s
}

// New creates a new InternalError
func New(scope Scope, code Code) *InternalError {
	msg := messageFromCode(scope, code)
	err := &InternalError{
		Code:    code,
		Scope:   scope,
		Message: msg,
	}
	return err
}

// Error satisfies the error type interface
func (error *InternalError) Error() string {
	return error.Message
}

// IsInternalError tests to see if this is an internal error or a native error type
func IsInternalError(err error) bool {
	if _, ok := err.(*InternalError); ok {
		return true
	}
	return false
}

// ToInternalError converts an error to an InternalError
func ToInternalError(err error) *InternalError {
	if internal, ok := err.(*InternalError); ok {
		return internal
	}
	code := New(ScopeGeneral, ErrorInternalEscape)
	return code
}

func messageFromScope(scope Scope) string {
	msgScope := "unknown scope"
	switch scope {
	case ScopeAccount:
		msgScope = "account"
	case ScopeAPI:
		msgScope = "api"
	case ScopeCollection:
		msgScope = "collection"
	case ScopeCrypto:
		msgScope = "crypto"
	case ScopeDB:
		msgScope = "db"
	case ScopeGeneral:
		msgScope = "general"
	case ScopeNote:
		msgScope = "note"
	case ScopeNotebook:
		msgScope = "notebook"
	case ScopeRPC:
		msgScope = "rpc"
	case ScopeShelf:
		msgScope = "shelf"
	case ScopeTag:
		msgScope = "tag"
	case ScopeTemplate:
		msgScope = "template"
	case ScopeTitle:
		msgScope = "title"
	case ScopeUIState:
		msgScope = "ui"
	case ScopeUser:
		msgScope = "user"
	default:
		msgScope = "default"
	}
	return msgScope
}

func messageFromCode(scope Scope, code Code) string {
	msgScope := messageFromScope(scope)
	msg := defaultMessageFromCode(code)
	errorString := fmt.Sprint(msgScope, " - ", msg)
	return errorString
}

func defaultMessageFromCode(code Code) string {
	msg := "unknown internal error"
	switch code {
	case ErrorUnknown:
		// default value
	case ErrorInternalEscape:
		msg = "internal error escape"
	case ErrorUnauthorized:
		msg = "error unauthorized"
	case ErrorInvalidType:
		msg = "error invalid type"
	case ErrorMissingDB:
		msg = "error missing db"
	case ErrorDbOpen:
		msg = "error opening db"
	case ErrorCreateBucket:
		msg = "error creating bucket"
	case ErrorMarshal:
		msg = "error marshaling"
	case ErrorOpenKey:
		msg = "error retrieving key"
	case ErrorEncrypt:
		msg = "error encrypting"
	case ErrorDecrypt:
		msg = "error decrypting"
	case ErrorCrypto:
		msg = "cryptography error"
	case ErrorWriteBucket:
		msg = "error writing to bucket"
	case ErrorSave:
		msg = "error saving"
	case ErrorBucketMissing:
		msg = "error bucket missing"
	case ErrorDecode:
		msg = "error decoding"
	case ErrorDeriveKey:
		msg = "error deriving key"
	case ErrorConvertID:
		msg = "error converting id"
	case ErrorLookup:
		msg = "error looking up"
	case ErrorLoad:
		msg = "error loading"
	case ErrorLoadAll:
		msg = "error loading all"
	case ErrorDelete:
		msg = "error deleting"
	case ErrorCreate:
		msg = "error creating"
	}

	return msg
}
