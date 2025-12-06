package dto

type LoginParams struct {
	Username string
	Password string
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}
