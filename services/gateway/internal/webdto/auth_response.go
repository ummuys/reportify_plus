package webdto

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type CreateUserResponse struct {
	UserID string `json:"user_id"`
}

type UpdateUserResponse struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}

type User struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}
