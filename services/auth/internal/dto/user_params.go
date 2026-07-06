package dto

type CreateUserParams struct {
	Username    string
	Password    string
	Role        string
	IsProtected bool
}

type DeleteUserParams struct {
	UserID string
}

type UpdateUserParams struct {
	UserID   string
	Username string
	Password string
	Role     string
	IsActive bool
}
