package service

import (
	"context"

	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type AuthService interface {
	Login(ctx context.Context, in dto.LoginParams) (dto.LoginResult, error)
	CreateUser(ctx context.Context, in dto.CreateUserParams) (dto.CreateUserResult, error)
	CreateBaseAdmin(ctx context.Context, in dto.CreateUserParams) error
	UpdateUser(ctx context.Context, in dto.UpdateUserParams) (dto.UpdateUserResult, error)
	DeleteUser(ctx context.Context, in dto.DeleteUserParams) (dto.DeleteUserResult, error)
	ListUsers(ctx context.Context) (dto.ListUsersResult, error)
	RefreshToken(ctx context.Context, in dto.RefreshTokenParams) (dto.RefreshTokenResult, error)
}
