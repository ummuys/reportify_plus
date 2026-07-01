package dto

type LoginParams struct {
	Username string
	Password string
}

type RefreshTokenParams struct {
	RefreshToken string
}
