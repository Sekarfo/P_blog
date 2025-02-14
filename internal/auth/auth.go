package auth

// UserModel represents the subset of the users table needed for authentication
type UserModel struct {
	ID       uint
	Password string
}
