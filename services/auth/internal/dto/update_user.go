package dto

type UpdateUserParams struct {
	UserID   string
	Username string
	Password string
	Role     string
	IsActive bool
}

type UpdateUserResult struct {
	UserID   string
	Username string
	Role     string
	IsActive bool
}
