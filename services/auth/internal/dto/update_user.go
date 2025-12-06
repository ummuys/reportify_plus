package dto

type UpdateUserParams struct {
	UserID   int64
	Username string
	Password string
	Role     string
	IsActive bool
}

type UpdateUserResult struct {
	UserID   int64
	Username string
	Password string
	Role     string
	IsActive bool
}
