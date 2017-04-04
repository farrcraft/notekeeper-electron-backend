package codes

import (
	"strconv"
)

// ErrorCode is the error code type
type ErrorCode int32

// InternalError is a custom error type
type InternalError struct {
	Code    ErrorCode
	Message string
}

// These are the status types that can be passed to the front end
const (
	StatusOK    = "OK"
	StatusError = "ERROR"
)

// These are the error codes that can be passed to the front end
const (
	ErrorOK ErrorCode = iota + 1
	ErrorUnknown
	ErrorInternalEscape

	ErrorMasterDbOpen

	ErrorMasterDbOpenDecode
	ErrorUnlockAccountDecode
	ErrorCreateAccountDecode
	ErrorSigninAccountDecode
	ErrorSaveUIStateDecode

	ErrorUnlockActiveAccount
	ErrorUnlockActiveUser

	ErrorSignoutActiveAccount
	ErrorSignoutActiveUser

	ErrorLockActiveAccount
	ErrorLockActiveUser
	ErrorLockActiveDb

	ErrorUIStateMissingDb
	ErrorUIStateCreateBucket
	ErrorDefaultUIStateMarshal
	ErrorDefaultUIStateWrite
	ErrorDefaultUIStateSave
	ErrorCreateUIState
	ErrorUIStateDecode
	ErrorUIStateBucket
	ErrorLoadUIState
	ErrorUIStatemarshal
	ErrorUIStateWrite
	ErrorUIStateSave

	ErrorAccountDbOpen
	ErrorAccountBucket
	ErrorAccountMarshal
	ErrorAccountKey
	ErrorAccountEncrypt
	ErrorAccountWrite
	ErrorAccountSave

	ErrorAccountMapBucket
	ErrorAccountMapKey
	ErrorAccountMapWrite
	ErrorAccountMapSave

	ErrorAccountMapDerive
	ErrorAccountMapConvert
	ErrorAccountMapLookup

	ErrorAccountLookup
	ErrorAccountBucketMissing

	ErrorAccountLoad
	ErrorAccountKeyOpen
	ErrorAccountDecrypt
	ErrorAccountDecode
	ErrorAccountLoadView

	ErrorUserMapDerive
	ErrorUserMapConvert
	ErrorUserLookup
	ErrorUserLoad
	ErrorUserKeyDerive
	ErrorUserDecrypt
	ErrorUserDecode
	ErrorUserBucketCreate
	ErrorUserMarshal
	ErrorUserEncrypt
	ErrorUserWrite
	ErrorUserSave
	ErrorUserMapBucket
	ErrorUserMapKey
	ErrorUserMapWrite
	ErrorUserMapSave

	ErrorNotebookBucket
	ErrorNotebookMarshal
	ErrorNotebookKey
	ErrorNotebookDecrypt
	ErrorNotebookWrite
	ErrorNotebookSave
)

// String converts error code to a string
func (e ErrorCode) String() string {
	s := strconv.Itoa(int(e))
	return s
}

// New creates a new InternalError
func New(code ErrorCode) *InternalError {
	msg := messageFromCode(code)
	err := &InternalError{
		Code:    code,
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
	code := New(ErrorInternalEscape)
	return code
}

func messageFromCode(code ErrorCode) string {
	msg := "unknown internal error"
	switch code {
	case ErrorUnknown:
		// default value
		break
	case ErrorInternalEscape:
		msg = "internal error escape"
		break

	case ErrorMasterDbOpen:
		msg = "error opening master db"
		break

	case ErrorMasterDbOpenDecode:
		msg = "error decoding payload for open master db request"
		break
	case ErrorCreateAccountDecode:
		msg = "error decoding payload for create account request"
		break
	case ErrorUnlockAccountDecode:
		msg = "error decoding payload for unlock account request"
		break
	case ErrorSigninAccountDecode:
		msg = "error decoding payload for signin account request"
		break
	case ErrorSaveUIStateDecode:
		msg = "error decoding payload for signin account request"
		break

	case ErrorUnlockActiveAccount:
		msg = "error no active account to unlock"
		break
	case ErrorUnlockActiveUser:
		msg = "error no active user to unlock"

	case ErrorSignoutActiveAccount:
		msg = "error no active account to signout"
		break
	case ErrorSignoutActiveUser:
		msg = "error no active user to signout"
		break

	case ErrorLockActiveAccount:
		msg = "error no active account to lock"
		break
	case ErrorLockActiveUser:
		msg = "error no active user to lock"
		break
	case ErrorLockActiveDb:
		msg = "error no active db to lock"
		break

	case ErrorCreateUIState:
		msg = "error initializing ui state"
		break
	case ErrorLoadUIState:
		msg = "error loading ui state"
		break
	case ErrorUIStateMissingDb:
		msg = "error missing ui state db"
		break
	case ErrorUIStateCreateBucket:
		msg = "error creating default ui state bucket"
		break
	case ErrorDefaultUIStateMarshal:
		msg = "error marshaling default ui state"
		break
	case ErrorDefaultUIStateWrite:
		msg = "error writing default ui state"
		break
	case ErrorDefaultUIStateSave:
		msg = "error saving default ui state"
		break
	case ErrorUIStateDecode:
		msg = "error decoding ui state"
		break
	case ErrorUIStateBucket:
		msg = "error creating ui state bucket"
		break
	case ErrorUIStatemarshal:
		msg = "error marshaling ui state"
		break
	case ErrorUIStateWrite:
		msg = "error writing ui state"
		break
	case ErrorUIStateSave:
		msg = "error saving ui state"
		break

	case ErrorAccountDbOpen:
		msg = "error opening account db"
		break

	case ErrorAccountBucket:
		msg = "error creating account bucket"
		break
	case ErrorAccountMarshal:
		msg = "error marshaling account"
		break
	case ErrorAccountKey:
		msg = "error opening account key"
		break
	case ErrorAccountEncrypt:
		msg = "error encrypting account"
		break
	case ErrorAccountWrite:
		msg = "error writing account bucket"
		break
	case ErrorAccountSave:
		msg = "error saving account"
		break

	case ErrorAccountMapBucket:
		msg = "error creating account map bucket"
		break
	case ErrorAccountMapKey:
		msg = "error creating account map key"
		break
	case ErrorAccountMapWrite:
		msg = "error saving account map"
		break
	case ErrorAccountMapSave:
		msg = "error mapping account"
		break

	case ErrorAccountMapDerive:
		msg = "error deriving account map key"
		break
	case ErrorAccountMapConvert:
		msg = "error converting account map"
		break
	case ErrorAccountMapLookup:
		msg = "error looking up account map"
		break

	case ErrorAccountLookup:
		msg = "error looking up account"
		break
	case ErrorAccountBucketMissing:
		msg = "error account bucket missing"
		break
	case ErrorAccountLoad:
		msg = "error loading account"
		break
	case ErrorAccountKeyOpen:
		msg = "error opening account key"
		break
	case ErrorAccountDecrypt:
		msg = "error decrypting account data"
		break
	case ErrorAccountDecode:
		msg = "error decoding account data"
		break
	case ErrorAccountLoadView:
		msg = "internal account load error"
		break

	case ErrorUserMapDerive:
		msg = "error deriving user map key"
		break
	case ErrorUserMapConvert:
		msg = "error converting account map"
		break
	case ErrorUserLookup:
		msg = "error looking up user"
		break
	case ErrorUserLoad:
		msg = "error loading user"
		break
	case ErrorUserKeyDerive:
		msg = "error deriving user key"
		break
	case ErrorUserDecrypt:
		msg = "error decrypting user"
		break
	case ErrorUserDecode:
		msg = "error decoding user"
		break
	case ErrorUserBucketCreate:
		msg = "error creating user bucket"
		break
	case ErrorUserMarshal:
		msg = "error marshaling user"
		break
	case ErrorUserEncrypt:
		msg = "error encrypting user"
		break
	case ErrorUserWrite:
		msg = "error writing user bucket"
		break
	case ErrorUserSave:
		msg = "error saving user"
		break
	case ErrorUserMapBucket:
		msg = "error creating user map bucket"
		break
	case ErrorUserMapKey:
		msg = "error creating user map key"
		break
	case ErrorUserMapWrite:
		msg = "error writing user map"
		break
	case ErrorUserMapSave:
		msg = "error saving user map"
		break

	case ErrorNotebookBucket:
		msg = "error creating notebook bucket"
		break
	case ErrorNotebookMarshal:
		msg = "error marshaling notebook"
		break
	case ErrorNotebookKey:
		msg = "error retrieving notebook key"
		break
	case ErrorNotebookDecrypt:
		msg = "error decrypting notebook"
		break
	case ErrorNotebookWrite:
		msg = "error writing notebook"
		break
	case ErrorNotebookSave:
		msg = "error saving notebook"
		break
	}

	return msg
}
