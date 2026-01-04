package dto

type LoginResult struct {
	AccessToken  string
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken string
}
