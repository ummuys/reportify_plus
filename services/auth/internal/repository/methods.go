package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/db"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/auth/internal/dto"
)

type authDB struct {
	logger    zerolog.Logger
	pool      *pgxpool.Pool
	adminUUID string
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

func (db *authDB) Close() {
	db.pool.Close()
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
	db.logger.Debug().Str("evt", "call CreateBaseAdmin").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	// P0-07 4 июля
	var roleID int

	err := db.pool.QueryRow(qctx, getRoleQuery, in.Role).Scan(&roleID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return dto.CreateUserResult{}, err
		}
	}
	// ------------------
	uuid := uuid.New().String()

	if _, err := db.pool.Exec(qctx, createUserQuery, uuid, in.Username, in.Password, in.Role); err != nil {
		return dto.CreateUserResult{}, err
	}

	return dto.CreateUserResult{UserID: uuid}, nil
}

func (db *authDB) SetAdminUUID(UUID string) {
	db.adminUUID = UUID
}

func (db *authDB) UpdateUser(ctx context.Context, in dto.UpdateUserParams) (out dto.UpdateUserResult, err error) {
	db.logger.Debug().Str("evt", "call UpdateUser").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	// P0-09 2 июля
	var tx pgx.Tx

	tx, err = db.pool.Begin(qctx)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				db.logger.Error().Err(rbErr).Str("evt", "call UpdateUser").Msg("")
			}
		}
	}()

	b := &pgx.Batch{}

	if in.Username != "" {
		b.Queue(updateUsernameQuery, in.UserID, in.Username)
	}

	if in.Password != "" {
		b.Queue(updatePasswordQuery, in.UserID, in.Password)
	}

	if in.Role != "" {
		b.Queue(updateRoleQuery, in.UserID, in.Role)
	}

	br := tx.SendBatch(qctx, b)

	var res pgconn.CommandTag
	for range b.Len() {

		if res, err = br.Exec(); err != nil {
			_ = br.Close()
			return
		}

		if res.RowsAffected() == 0 {
			_ = br.Close()
			err = pgx.ErrNoRows
			return
		}
	}

	if err = br.Close(); err != nil {
		return
	}

	err = tx.QueryRow(qctx, getUserQuery, in.UserID).Scan(
		&out.UserID,
		&out.Username,
		&out.Role,
	)
	if err != nil {
		err = pgx.ErrNoRows
		return
	}

	if err = tx.Commit(qctx); err != nil {
		return
	}

	// --------------------

	// out = dto.UpdateUserResult{
	// 	UserID:   in.UserID,
	// 	Username: in.Username,
	// 	Role:     in.Role,
	// IsActive: in.IsActive,
	// }
	return
}

func (db *authDB) DeleteUser(ctx context.Context, in dto.DeleteUserParams) (dto.DeleteUserResult, error) {
	db.logger.Debug().Str("evt", "call DeleteUser").Msg("")
	if in.UserID == db.adminUUID {
		return dto.DeleteUserResult{}, errs.ErrPgInsufficientPrivilege
	}
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	res, err := db.pool.Exec(qctx, deleteUserQuery, in.UserID)
	if err != nil {
		return dto.DeleteUserResult{}, err
	}

	if res.RowsAffected() == 0 {
		return dto.DeleteUserResult{}, pgx.ErrNoRows
	}

	return dto.DeleteUserResult(in), nil
}

func (db *authDB) ListUsers(ctx context.Context) (dto.ListUsersResult, error) {
	db.logger.Debug().Str("evt", "call GetUsers").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	rows, err := db.pool.Query(qctx, ListUsersQuery)
	if err != nil {
		return dto.ListUsersResult{}, err
	}
	defer rows.Close()

	var users []dto.User
	for rows.Next() {
		var u dto.User
		err = rows.Scan(&u.UserID, &u.Username, &u.Role)
		if err != nil {
			return dto.ListUsersResult{}, err
		}
		users = append(users, u)
	}

	fmt.Println(users)

	if err := rows.Err(); err != nil {
		return dto.ListUsersResult{}, err
	}

	return dto.ListUsersResult{Users: users}, nil
}
