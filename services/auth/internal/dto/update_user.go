package dto

type UpdateUserRequest struct {
	UserID   int64
	Username string
	Password string
	Role     string
	IsActive bool
}

type UpdateUserResponse struct {
	UserID   int64
	Username string
	Password string
	Role     string
	IsActive bool
}
