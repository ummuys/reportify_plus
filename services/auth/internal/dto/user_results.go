package dto

type CreateUserResult struct {
	UserID string
}

type DeleteUserResult struct {
	UserID string
}

type UpdateUserResult struct {
	UserID   string
	Username string
	Role     string
	IsActive bool
}

type User struct {
	UserID   string
	Username string
	Role     string
	IsActive bool
}

type ListUsersResult struct {
	Users []User
}

// Struct for validate password
type AuthUser struct {
	UserID   string
	Password string
	Role     string
}
