package auth

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/auth/internal/dto"
	"github.com/ummuys/reportify/services/auth/internal/repository"
	"github.com/ummuys/reportify/services/auth/internal/secure"
)

type authService struct {
	ph     secure.PasswordHasher
	tm     secure.TokenManager
	db     repository.AuthDB
	logger zerolog.Logger
}

func NewAuthService(ph secure.PasswordHasher, tm secure.TokenManager, db repository.AuthDB, logger zerolog.Logger) AuthService {
	return &authService{ph: ph, tm: tm, db: db, logger: logger}
}

func (as *authService) Login(ctx context.Context, in dto.LoginRequest) (dto.LoginResponse, error) {
	as.logger.Debug().Str("evt", "call Login").Msg("")
	out, err := as.db.Login(ctx, in)
	if err != nil {
		return dto.LoginResponse{}, errs.ParsePgError(err)
	}

	if !as.ph.CheckHash(in.Password, out.Password) {
		return dto.LoginResponse{}, errs.ErrUserUnauthorized
	}

	at, err := as.tm.GenerateAccessToken(out.UserID, out.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}

	rt, err := as.tm.GenerateRefreshToken(out.UserID, out.Role)
	if err != nil {
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{AccessToken: at, RefreshToken: rt}, nil
}

func (as *authService) CreateUser(ctx context.Context, in dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	return dto.CreateUserResponse{}, nil
}

func (as *authService) UpdateUser(ctx context.Context, in dto.UpdateUserRequest) (dto.UpdateUserResponse, error) {
	return dto.UpdateUserResponse{}, nil
}

func (as *authService) DeleteUser(ctx context.Context, in dto.DeleteUserRequest) (dto.DeleteUserResponse, error) {
	return dto.DeleteUserResponse{}, nil
}

func (as *authService) ListUsers(ctx context.Context) (dto.ListUsersResponse, error) {
	return dto.ListUsersResponse{}, nil
}

func (as *authService) RefreshToken(ctx context.Context, in dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error) {
	return dto.RefreshTokenResponse{}, nil
}
