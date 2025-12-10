package dto

type CreateBaseAdminParams struct {
	UserID   int64
	Username string
	Password string
	Role     string
}

type CreateUserParams struct {
	Username string
	Password string
	Role     string
}

type CreateUserResult struct {
	UserID int64
}
