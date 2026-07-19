package constants

import (
	"errors"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrDuplicatedEmail    = errors.New("email is already used")
	UserIDContextKey      = "user_id"
)

const (
	ONE_DAY_TIMESTAMP_SECS   = time.Second * 60 * 60 * 24
	ONE_MONTH_TIMESTAMP_SECS = ONE_DAY_TIMESTAMP_SECS * 30
)
