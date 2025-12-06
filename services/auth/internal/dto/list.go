package dto

type User struct {
	UserID   string
	Username string
	Role     string
	IsActive bool
}

type ListUsersResult struct {
	Users []User
}
