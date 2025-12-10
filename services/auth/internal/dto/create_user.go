package dto

type CreateUserParams struct {
	UserID   string
	Username string
	Password string
	Role     string
}

type CreateUserResult struct {
	UserID string
}
