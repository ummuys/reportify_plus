package pkg

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/ummuys/reportify/pkg/config"
)

type tokMan struct {
	cfg config.TMConfig
}

func NewTokenManager() (TokenManager, error) {
	c, err := config.ParseTMConfig()
	if err != nil {
		return nil, err
	}
	return &tokMan{cfg: c}, nil
}

func (tm *tokMan) GenerateRefreshToken(user_id string, role string) (string, error) {
	claims := RefreshClaims{
		UserID: user_id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.cfg.RefreshTokenLimit)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app-token-manager",
			Audience:  []string{"app-users"},
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return refresh.SignedString([]byte(tm.cfg.RefreshSecret))
}

func (tm *tokMan) GetRefreshLifetime() time.Duration {
	return tm.cfg.RefreshTokenLimit
}

func (tm *tokMan) GenerateAccessToken(user_id string, role string) (string, error) {
	claims := AccessClaims{
		UserID: user_id,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tm.cfg.AccessTokenLimit)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app-token-manager",
			Audience:  []string{"app-users"},
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return access.SignedString([]byte(tm.cfg.AccessSecret))
}

func (tm *tokMan) ValidateAccessToken(rawToken string) (AccessClaims, error) {
	token, claims, err := tm.unhashAccessToken(rawToken)
	if err != nil {
		return AccessClaims{}, err
	}
	if !token.Valid {
		return AccessClaims{}, fmt.Errorf("invalid token")
	}
	return claims, err
}

func (tm *tokMan) ValidateRefreshToken(rawToken string) (RefreshClaims, error) {
	token, claims, err := tm.unhashRefreshToken(rawToken)
	if err != nil {
		return RefreshClaims{}, err
	}
	if !token.Valid {
		return RefreshClaims{}, fmt.Errorf("invalid token")
	}
	return claims, err
}

func (tm *tokMan) unhashAccessToken(token string) (*jwt.Token, AccessClaims, error) {
	var claims AccessClaims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(tm.cfg.AccessSecret), nil
	})
	return jwtToken, claims, err
}

func (tm *tokMan) unhashRefreshToken(token string) (*jwt.Token, RefreshClaims, error) {
	var claims RefreshClaims
	jwtToken, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(tm.cfg.RefreshSecret), nil
	})
	return jwtToken, claims, err
}
