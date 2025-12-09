package dto

type User struct {
	UserID   int64
	Username string
	Role     string
	IsActive bool
}

type ListUsersResult struct {
	Users []User
}
