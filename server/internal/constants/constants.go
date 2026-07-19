package constants

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrDuplicatedEmail    = errors.New("email is already used")
	UserIDContextKey      = "user_id"
)
