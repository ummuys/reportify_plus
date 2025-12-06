package dto

type CreateUserParams struct {
	Username string
	Password string
	Role     string
}

type CreateUserResult struct {
	UserID int64
}
