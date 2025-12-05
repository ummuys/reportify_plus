package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type authDB struct {
	logger zerolog.Logger
	pool   *pgxpool.Conn
}

func NewAuthDB(logger zerolog.Logger) AuthDB {
	return &authDB{
		logger: logger,
		pool:   nil,
	}
}

func (db *authDB) Login(ctx context.Context, in dto.LoginRequest) (dto.UserInfo, error) {
	return dto.UserInfo{}, nil
}

func (db *authDB) CreateUser(ctx context.Context, in dto.CreateUserRequest) (dto.CreateUserResponse, error) {
	return dto.CreateUserResponse{}, nil
}

func (db *authDB) UpdateUser(ctx context.Context, in dto.UpdateUserRequest) (dto.UpdateUserResponse, error) {
	return dto.UpdateUserResponse{}, nil
}

func (db *authDB) DeleteUser(ctx context.Context, in dto.DeleteUserRequest) (dto.DeleteUserResponse, error) {
	return dto.DeleteUserResponse{}, nil
}

func (db *authDB) ListUsers(ctx context.Context) (dto.ListUsersResponse, error) {
	return dto.ListUsersResponse{}, nil
}

func (db *authDB) RefreshToken(ctx context.Context, in dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error) {
	return dto.RefreshTokenResponse{}, nil
}
