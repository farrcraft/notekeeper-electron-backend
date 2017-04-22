package rpc

// IsLocked returns true if the user state appears locked
func (rpc *Server) IsLocked() bool {
	if rpc.UserState == UserStateLocked {
		return true
	}
	return false
}

// IsSignedIn returns true if the user state appears signed in
func (rpc *Server) IsSignedIn() bool {
	if rpc.UserState == UserStateSignedIn {
		return true
	}
	return false
}
