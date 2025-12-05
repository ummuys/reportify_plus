package dto

type RefreshTokenRequest struct {
	RefreshToken string
}

type RefreshTokenResponse struct {
	AccessToken string
}
