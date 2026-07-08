package errs

import "errors"

var (

	// SERVER ERR
	ErrServerInternal  = errors.New("something wrong with server, try again later")
	ErrBadRefreshToken = errors.New("bad refresh token")
	ErrBadAccessToken  = errors.New("bad access token")
	ErrInvalidJSON     = errors.New("invalid JSON format")

	// ANOTHER
	ErrDeleteAdmin        = errors.New("can't delete main admin")
	ErrInvalidPaylod      = errors.New("invalid payload")
	ErrInvalidUserID      = errors.New("invalid user_id")
	ErrInvalidReportID    = errors.New("invalid report_id")
	ErrInvalidData        = errors.New("invalid data")
	ErrUserUnauthorized   = errors.New("you need to auth")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUsernameExists     = errors.New("username already taken")
	ErrInvalidCredentials = errors.New("user or password are incorrect")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmptyUsername      = errors.New("empty username")
	ErrEmptySchemaName    = errors.New("schema name requeired")
	ErrEmptyTableName     = errors.New("table name required")
	ErrRoleNotFound       = errors.New("role not found")
)
