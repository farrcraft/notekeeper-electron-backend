package user

// Profile contains the minimal user information that is visible to all users of an account
type Profile struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
