package adapter

import (
	"context"

	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/services/auth/internal/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthAdapter struct {
	authv1.UnimplementedAuthServiceServer
	svc auth.AuthService
}

func NewAuthAdapter(svc auth.AuthService) *AuthAdapter {
	return &AuthAdapter{svc: svc}
}

// Login: заглушка
func (a *AuthAdapter) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	return &authv1.LoginResponse{}, nil
}

// CreateUser: заглушка
func (a *AuthAdapter) CreateUser(ctx context.Context, in *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	return &authv1.CreateUserResponse{}, nil
}

// UpdateUser: заглушка
func (a *AuthAdapter) UpdateUser(ctx context.Context, in *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	return &authv1.UpdateUserResponse{}, nil
}

// DeleteUser: заглушка
func (a *AuthAdapter) DeleteUser(ctx context.Context, in *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	return &authv1.DeleteUserResponse{}, nil
}

// ListUsers: заглушка
func (a *AuthAdapter) ListUsers(ctx context.Context, in *emptypb.Empty) (*authv1.ListUsersResponse, error) {
	return &authv1.ListUsersResponse{}, nil
}

// RefreshToken: заглушка
func (a *AuthAdapter) RefreshToken(ctx context.Context, in *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	return &authv1.RefreshTokenResponse{}, nil
}
