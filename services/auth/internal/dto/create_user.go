package dto

type CreateUserRequest struct {
	Username string
	Password string
	Role     string
}

type CreateUserResponse struct {
	UserID int64
}
