package rpc

// UserState type
type UserState int

// Valid user states
const (
	UserStateSignedOut UserState = iota
	UserStateSignedIn
	UserStateLocked
)
