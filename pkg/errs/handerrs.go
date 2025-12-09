package errs

import "errors"

var (

	// SERVER ERR
	ErrInternalServer  = errors.New("something wrong with server, try again later")
	ErrBadRefreshToken = errors.New("bad refresh token")
	ErrBadAccessToken  = errors.New("bad access token")
	ErrInvalidJSON     = errors.New("invalid JSON format")

	// ANOTHER
	ErrInvalidData        = errors.New("invalid data")
	ErrUserUnauthorized   = errors.New("you need to auth")
	ErrUserExists         = errors.New("user already exists")
	ErrUsernameExists     = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("user or password are incorrect")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmptyUsername      = errors.New("empty username")
	ErrEmptySchemaName    = errors.New("schema name requeired")
	ErrEmptyTableName     = errors.New("table name required")
)
