package errlist

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)
