package auth

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	pkg "github.com/ummuys/reportify/pkg/tm"
	"github.com/ummuys/reportify/services/auth/internal/dto"
	"github.com/ummuys/reportify/services/auth/internal/repository"
	"github.com/ummuys/reportify/services/auth/internal/secure"
)

type authService struct {
	ph     secure.PasswordHasher
	tm     pkg.TokenManager
	db     repository.AuthDB
	logger zerolog.Logger
}

func NewAuthService(ph secure.PasswordHasher, tm pkg.TokenManager, db repository.AuthDB, baseLogger zerolog.Logger) AuthService {
	logger := baseLogger.With().Str("component", "svc").Logger()
	return &authService{ph: ph, tm: tm, db: db, logger: logger}
}

func (as *authService) Login(ctx context.Context, in dto.LoginParams) (dto.LoginResult, error) {
	as.logger.Debug().Str("evt", "call Login").Msg("")
	out, err := as.db.Login(ctx, in.Username)
	if err != nil {
		return dto.LoginResult{}, errs.ParsePgError(err)
	}

	if !as.ph.CheckHash(in.Password, out.Password) {
		return dto.LoginResult{}, errs.ErrInvalidCredentials
	}

	at, err := as.tm.GenerateAccessToken(out.UserID, out.Role)
	if err != nil {
		return dto.LoginResult{}, err
	}

	rt, err := as.tm.GenerateRefreshToken(out.UserID, out.Role)
	if err != nil {
		return dto.LoginResult{}, err
	}
	return dto.LoginResult{AccessToken: at, RefreshToken: rt}, nil
}

func (as *authService) CreateBaseAdmin(ctx context.Context, in dto.CreateBaseAdminParams) error {
	as.logger.Debug().Str("evt", "call CreateUser").Msg("")

	hashPass, err := as.ph.Hash(in.Password)
	if err != nil {
		return err
	}
	in.Password = hashPass

	err = as.db.CreateBaseAdmin(ctx, in)
	if err != nil {
		return errs.ParsePgError(err)
	}
	return nil
}

func (as *authService) CreateUser(ctx context.Context, in dto.CreateUserParams) (dto.CreateUserResult, error) {
	as.logger.Debug().Str("evt", "call CreateUser").Msg("")

	hashPass, err := as.ph.Hash(in.Password)
	if err != nil {
		return dto.CreateUserResult{}, err
	}
	in.Password = hashPass

	out, err := as.db.CreateUser(ctx, in)
	if err != nil {
		return dto.CreateUserResult{}, errs.ParsePgError(err)
	}
	return out, nil
}

func (as *authService) UpdateUser(ctx context.Context, in dto.UpdateUserParams) (dto.UpdateUserResult, error) {
	as.logger.Debug().Str("evt", "call UpdateUser").Msg("")

	if in.Password != "" {
		hashPass, err := as.ph.Hash(in.Password)
		if err != nil {
			return dto.UpdateUserResult{}, err
		}
		in.Password = hashPass
	}

	out, err := as.db.UpdateUser(ctx, in)
	if err != nil {
		return dto.UpdateUserResult{}, errs.ParsePgError(err)
	}
	return out, nil
}

func (as *authService) DeleteUser(ctx context.Context, in dto.DeleteUserParams) (dto.DeleteUserResult, error) {
	as.logger.Debug().Str("evt", "call DeleteUser").Msg("")

	out, err := as.db.DeleteUser(ctx, in)
	if err != nil {
		return dto.DeleteUserResult{}, errs.ParsePgError(err)
	}
	return out, nil
}

func (as *authService) ListUsers(ctx context.Context) (dto.ListUsersResult, error) {
	as.logger.Debug().Str("evt", "call ListUsers").Msg("")

	out, err := as.db.ListUsers(ctx)
	if err != nil {
		return dto.ListUsersResult{}, errs.ParsePgError(err)
	}
	return out, nil
}

func (as *authService) RefreshToken(ctx context.Context, in dto.RefreshTokenParams) (dto.RefreshTokenResult, error) {
	as.logger.Debug().Str("evt", "call RefreshToken").Msg("")

	rc, err := as.tm.ValidateRefreshToken(in.RefreshToken)
	if err != nil {
		return dto.RefreshTokenResult{}, err
	}

	userID := rc.UserID
	role := rc.Role

	access, err := as.tm.GenerateAccessToken(userID, role)
	if err != nil {
		return dto.RefreshTokenResult{}, err
	}

	return dto.RefreshTokenResult{AccessToken: access}, nil
}
