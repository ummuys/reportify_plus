package dto

type RefreshTokenParams struct {
	RefreshToken string
}

type RefreshTokenResult struct {
	AccessToken string
}
