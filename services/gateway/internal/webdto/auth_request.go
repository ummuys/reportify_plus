package webdto

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type UpdateUserRequest struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type DeleteUserResponse struct {
	UserID string `json:"user_id"`
}
