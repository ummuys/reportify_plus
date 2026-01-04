package adapter

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	authv1 "github.com/ummuys/reportify/api/pb/auth/service/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/auth/internal/dto"
	"github.com/ummuys/reportify/services/auth/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthAdapter struct {
	logger zerolog.Logger
	authv1.UnimplementedAuthServiceServer
	svc service.AuthService
}

func NewAuthAdapter(svc service.AuthService, baseLogger zerolog.Logger) *AuthAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &AuthAdapter{svc: svc, logger: logger}
}

func (a *AuthAdapter) Login(ctx context.Context, in *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	a.logger.Debug().Str("evt", "call Login").Msg("")
	out, err := a.svc.Login(ctx, dto.LoginParams{
		Username: in.Username,
		Password: in.Password,
	})

	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.Unauthenticated, errs.ErrInvalidCredentials.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &authv1.LoginResponse{AccessToken: out.AccessToken, RefreshToken: out.RefreshToken}, nil
}

func (a *AuthAdapter) CreateUser(ctx context.Context, in *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	a.logger.Debug().Str("evt", "call CreateUser").Msg("")
	out, err := a.svc.CreateUser(ctx, dto.CreateUserParams{
		Username: in.Username,
		Password: in.Password,
		Role:     in.Role,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrDuplicate):
			return nil, status.Error(codes.AlreadyExists, errs.ErrUserAlreadyExists.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &authv1.CreateUserResponse{UserId: out.UserID}, nil
}

func (a *AuthAdapter) UpdateUser(ctx context.Context, in *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	a.logger.Debug().Str("evt", "call UpdateUser").Msg("")
	out, err := a.svc.UpdateUser(ctx, dto.UpdateUserParams{
		UserID:   in.UserId,
		Username: in.Username,
		Password: in.Password,
		Role:     in.Role,
		IsActive: in.IsActive,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, errs.ErrInvalidData.Error())
		case errors.Is(err, errs.ErrDuplicate):
			return nil, status.Error(codes.AlreadyExists, errs.ErrUsernameExists.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &authv1.UpdateUserResponse{UserId: out.UserID, Username: out.Username, Role: out.Role, IsActive: out.IsActive}, nil
}

func (a *AuthAdapter) DeleteUser(ctx context.Context, in *authv1.DeleteUserRequest) (*authv1.DeleteUserResponse, error) {
	a.logger.Debug().Str("evt", "call DeleteUser").Msg("")
	out, err := a.svc.DeleteUser(ctx, dto.DeleteUserParams{
		UserID: in.UserId,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, errs.ErrUserNotFound.Error())
		case errors.Is(err, errs.ErrInsufficientPrivilege):
			return nil, status.Error(codes.PermissionDenied, errs.ErrDeleteAdmin.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &authv1.DeleteUserResponse{UserId: out.UserID}, nil
}

func (a *AuthAdapter) ListUsers(ctx context.Context, in *emptypb.Empty) (*authv1.ListUsersResponse, error) {
	a.logger.Debug().Str("evt", "call ListUsers").Msg("")
	out, err := a.svc.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var resp authv1.ListUsersResponse
	resp.Users = make([]*authv1.User, 0, len(out.Users))
	for _, user := range out.Users {
		resp.Users = append(resp.Users, &authv1.User{
			UserId:   user.UserID,
			Username: user.Username,
			Role:     user.Role,
		})
	}

	return &resp, nil
}

// RefreshToken: заглушка
func (a *AuthAdapter) RefreshToken(ctx context.Context, in *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	a.logger.Debug().Str("evt", "call RefreshToken").Msg("")
	out, err := a.svc.RefreshToken(ctx, dto.RefreshTokenParams{RefreshToken: in.RefreshToken})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrBadRefreshToken):
			return nil, status.Error(codes.Unauthenticated, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &authv1.RefreshTokenResponse{AccessToken: out.AccessToken}, nil
}
