package service

import (
	"context"
	"errors"

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
		perr := errs.ParsePgError(err)
		if !errors.Is(perr, errs.ErrPgNotFound) {
			as.logger.Error().
				Err(err).
				Str("db-method", "Login").
				Msg("db login failed")
			return dto.LoginResult{}, perr
		}

		return dto.LoginResult{}, errs.ErrInvalidCredentials
	}

	if !as.ph.CheckHash(in.Password, out.Password) {
		as.logger.Warn().
			Str("user-id", out.UserID).
			Msg("invalid credentials")
		return dto.LoginResult{}, errs.ErrInvalidCredentials
	}

	at, err := as.tm.GenerateAccessToken(out.UserID, out.Role)
	if err != nil {
		as.logger.Error().
			Err(err).
			Str("op", "GenerateAccessToken").
			Str("user-id", out.UserID).
			Str("role", out.Role).
			Msg("token generation failed")
		return dto.LoginResult{}, err
	}

	rt, err := as.tm.GenerateRefreshToken(out.UserID, out.Role)
	if err != nil {
		as.logger.Error().
			Err(err).
			Str("op", "GenerateRefreshToken").
			Str("user-id", out.UserID).
			Str("role", out.Role).
			Msg("token generation failed")
		return dto.LoginResult{}, err
	}

	return dto.LoginResult{AccessToken: at, RefreshToken: rt}, nil
}

func (as *authService) CreateUser(ctx context.Context, in dto.CreateUserParams) (dto.CreateUserResult, error) {
	as.logger.Debug().Str("evt", "call CreateUser").Msg("")

	hashPass, err := as.ph.Hash(in.Password)
	if err != nil {
		as.logger.Error().
			Err(err).
			Str("op", "HashPassword").
			Msg("create user failed")
		return dto.CreateUserResult{}, err
	}
	in.Password = hashPass

	out, err := as.db.CreateUser(ctx, in)
	if err != nil {

		as.logger.Error().
			Err(err).
			Str("db-method", "CreateUser").
			Str("role", in.Role).
			Msg("create user failed")

		return dto.CreateUserResult{}, errs.ParsePgError(err)
	}

	as.logger.Info().
		Str("evt", "user.created").
		Str("user-id", out.UserID).
		Str("role", in.Role).
		Msg("user created")

	return out, nil
}

func (as *authService) CreateBaseAdmin(ctx context.Context, in dto.CreateUserParams) error {
	as.logger.Debug().Str("evt", "call CreateBaseAdmin").Msg("")

	hashPass, err := as.ph.Hash(in.Password)
	if err != nil {
		as.logger.Error().
			Err(err).
			Str("op", "HashPassword").
			Msg("create base admin failed")
		return err
	}
	in.Password = hashPass

	out, err := as.db.CreateUser(ctx, in)
	if err != nil {

		as.logger.Error().
			Err(err).
			Str("db-method", "CreateUser").
			Msg("create base admin failed")

		return errs.ParsePgError(err)
	}

	as.db.SetAdminUUID(out.UserID)
	as.logger.Warn().
		Str("user-id", out.UserID).
		Msg("Admin user created")

	return nil
}

func (as *authService) UpdateUser(ctx context.Context, in dto.UpdateUserParams) (dto.UpdateUserResult, error) {
	as.logger.Debug().Str("evt", "call UpdateUser").Msg("")

	changedPassword := false
	if in.Password != "" {
		changedPassword = true
		hashPass, err := as.ph.Hash(in.Password)
		if err != nil {
			as.logger.Error().
				Err(err).
				Str("op", "HashPassword").
				Str("user-id", in.UserID).
				Msg("update user failed")
			return dto.UpdateUserResult{}, err
		}
		in.Password = hashPass
	}

	out, err := as.db.UpdateUser(ctx, in)
	if err != nil {
		as.logger.Error().
			Err(err).
			Str("db-method", "UpdateUser").
			Str("user-id", in.UserID).
			Msg("update user failed")

		return dto.UpdateUserResult{}, errs.ParsePgError(err)
	}

	as.logger.Info().
		Str("evt", "user.updated").
		Str("user-id", in.UserID).
		Bool("changed_password", changedPassword).
		Bool("changed_username", in.Username != "").
		Bool("changed_role", in.Role != "").
		Msg("user updated")

	return out, nil
}

func (as *authService) DeleteUser(ctx context.Context, in dto.DeleteUserParams) (dto.DeleteUserResult, error) {
	as.logger.Debug().Str("evt", "call DeleteUser").Msg("")

	out, err := as.db.DeleteUser(ctx, in)
	if err != nil {
		perr := errs.ParsePgError(err)

		if !errors.Is(perr, errs.ErrPgNotFound) && !errors.Is(perr, errs.ErrPgInsufficientPrivilege) {
			as.logger.Error().
				Err(err).
				Str("db-method", "DeleteUser").
				Str("user-id", in.UserID).
				Msg("delete user failed")
		}

		if errors.Is(perr, errs.ErrPgInsufficientPrivilege) {
			as.logger.Warn().
				Str("evt", "user.delete_forbidden").
				Str("user-id", in.UserID).
				Msg("attempt to delete protected user")
		}

		return dto.DeleteUserResult{}, perr
	}

	as.logger.Info().
		Str("evt", "user.deleted").
		Str("user-id", in.UserID).
		Msg("user deleted")

	return out, nil
}

func (as *authService) ListUsers(ctx context.Context) (dto.ListUsersResult, error) {
	as.logger.Debug().Str("evt", "call ListUsers").Msg("")

	out, err := as.db.ListUsers(ctx)
	if err != nil {
		perr := errs.ParsePgError(err)

		as.logger.Error().
			Err(err).
			Str("db-method", "ListUsers").
			Msg("list users failed")

		return dto.ListUsersResult{}, perr
	}

	return out, nil
}

func (as *authService) RefreshToken(ctx context.Context, in dto.RefreshTokenParams) (dto.RefreshTokenResult, error) {
	as.logger.Debug().Str("evt", "call RefreshToken").Msg("")

	rc, err := as.tm.ValidateRefreshToken(in.RefreshToken)
	if err != nil {
		as.logger.Info().
			Err(err).
			Str("evt", "token.refresh_rejected").
			Msg("refresh token rejected")
		return dto.RefreshTokenResult{}, errs.ErrPgUnauthorized
	}

	access, err := as.tm.GenerateAccessToken(rc.UserID, rc.Role)
	if err != nil {
		as.logger.Error().
			Err(err).
			Str("op", "GenerateAccessToken").
			Str("user-id", rc.UserID).
			Str("role", rc.Role).
			Msg("refresh token failed")
		return dto.RefreshTokenResult{}, err
	}

	return dto.RefreshTokenResult{AccessToken: access}, nil
}
