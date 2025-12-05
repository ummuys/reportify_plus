package repository

import (
	"context"

	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type AuthDB interface {
	Login(ctx context.Context, in dto.LoginRequest) (dto.UserInfo, error)
	CreateUser(ctx context.Context, in dto.CreateUserRequest) (dto.CreateUserResponse, error)
	UpdateUser(ctx context.Context, in dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, in dto.DeleteUserRequest) (dto.DeleteUserResponse, error)
	ListUsers(ctx context.Context) (dto.ListUsersResponse, error)
	RefreshToken(ctx context.Context, in dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error)
}
