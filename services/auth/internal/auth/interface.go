package auth

import (
	"context"

	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type AuthService interface {
	Login(ctx context.Context, in dto.LoginRequest) (dto.LoginResponse, error)
	CreateUser(ctx context.Context, in dto.CreateUserRequest) (dto.CreateUserResponse, error)
	UpdateUser(ctx context.Context, in dto.UpdateUserRequest) (dto.UpdateUserResponse, error)
	DeleteUser(ctx context.Context, in dto.DeleteUserRequest) (dto.DeleteUserResponse, error)
	ListUsers(ctx context.Context) (dto.ListUsersResponse, error)
	RefreshToken(ctx context.Context, in dto.RefreshTokenRequest) (dto.RefreshTokenResponse, error)
}
