package dto

type CreateUserParams struct {
	Username string
	Password string
	Role     string
}

type DeleteUserParams struct {
	UserID string
}

type UpdateUserParams struct {
	UserID   string
	Username string
	Password string
	Role     string
	IsActive bool
}
