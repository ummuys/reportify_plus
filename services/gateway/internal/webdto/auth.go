package webdto

// LOGIN

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

// CREATE USER

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type CreateUserResponse struct {
	UserID int64 `json:"user_id"`
}

// UPDATE USER

type UpdateUserRequest struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type UpdateUserResponse struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

// DELETE

type DeleteUserResponse struct {
	UserID int64 `json:"user_id"`
}

// REFRESH TOKEN
type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// USERS

type User struct {
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	IsActive bool   `json:"is_active"`
}

type ListUsersResponse struct {
	Users []User `json:"users"`
}
