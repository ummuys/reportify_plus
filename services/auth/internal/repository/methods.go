package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/db"
	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type authDB struct {
	logger zerolog.Logger
	pool   *pgxpool.Pool
}

func NewAuthDB(ctx context.Context, baseLogger zerolog.Logger) (AuthDB, error) {
	dctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	cfg, err := config.ParseAuthDBEnv()
	if err != nil {
		return nil, err
	}

	pool, err := db.PoolFromConfig(dctx, cfg, "AUTH_DB")
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "db").Logger()

	return &authDB{
		logger: logger,
		pool:   pool,
	}, nil
}

func (db *authDB) Login(ctx context.Context, username string) (dto.AuthUser, error) {
	db.logger.Debug().Str("evt", "call Login").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var out dto.AuthUser
	if err := db.pool.QueryRow(qctx, loginQuery, username).Scan(&out.UserID, &out.Password, &out.Role); err != nil {
		return dto.AuthUser{}, err
	}

	return out, nil
}

func (db *authDB) CreateUser(ctx context.Context, in dto.CreateUserParams) (dto.CreateUserResult, error) {
	db.logger.Debug().Str("evt", "call CreateUser").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var out dto.CreateUserResult
	if err := db.pool.QueryRow(qctx, createUserQuery, in.Username, in.Password, in.Role).Scan(&out.UserID); err != nil {
		return dto.CreateUserResult{}, err
	}

	return out, nil
}

func (db *authDB) UpdateUser(ctx context.Context, in dto.UpdateUserParams) (out dto.UpdateUserResult, err error) {
	db.logger.Debug().Str("evt", "call UpdateUser").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var tx pgx.Tx

	// stack: rollback(3) - commit(2) - close batch(1)

	// tx failed -> rollback
	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				db.logger.Error().Err(rbErr).Msg("rollback failed")
			}
		}
	}()

	//commit tx
	defer func() {
		err = tx.Commit(qctx)
	}()

	tx, err = db.pool.Begin(qctx)
	if err != nil {
		return
	}

	b := &pgx.Batch{}

	if in.Username != "" {
		b.Queue(updateUsernameQuery, in.Username)
	}

	if in.Password != "" {
		b.Queue(updatePasswordQuery, in.Password)
	}

	if in.Role != "" {
		b.Queue(updateRoleQuery, in.Role)
	}

	br := tx.SendBatch(qctx, b)

	// close batch
	defer func() {
		err = br.Close()
	}()

	for range b.Len() {
		if _, err = br.Exec(); err != nil {
			return
		}
	}

	return dto.UpdateUserResult{}, nil
}

func (db *authDB) DeleteUser(ctx context.Context, in dto.DeleteUserParams) (dto.DeleteUserResult, error) {
	return dto.DeleteUserResult{}, nil
}

func (db *authDB) ListUsers(ctx context.Context) (dto.ListUsersResult, error) {
	return dto.ListUsersResult{}, nil
}

func (db *authDB) RefreshToken(ctx context.Context, in dto.RefreshTokenParams) (dto.RefreshTokenResult, error) {
	return dto.RefreshTokenResult{}, nil
}
