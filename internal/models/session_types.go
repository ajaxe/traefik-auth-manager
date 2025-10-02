package models

// Session holds the user's session information.
type Session struct {
	User SessionUser
}

// SessionUser holds profile information for the logged-in user.
type SessionUser struct {
	Name    string
	Sub     string
	Email   string
	Picture string
	IdpName string
}
