package repository

import (
	"context"

	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type AuthDB interface {
	Login(ctx context.Context, username string) (dto.AuthUser, error)
	CreateUser(ctx context.Context, in dto.CreateUserParams) (dto.CreateUserResult, error)
	UpdateUser(ctx context.Context, in dto.UpdateUserParams) (out dto.UpdateUserResult, err error)
	DeleteUser(ctx context.Context, in dto.DeleteUserParams) (dto.DeleteUserResult, error)
	ListUsers(ctx context.Context) (dto.ListUsersResult, error)
}
