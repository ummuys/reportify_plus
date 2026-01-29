package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/auth/internal/dto"
)

// ALL OK
func TestAuthService_Login_Success(t *testing.T) {
	svc, db, ph, tm := newSvc(t)
	ctx := context.Background()

	in := dto.LoginParams{Username: "bob", Password: "pass"}
	dbOut := dto.AuthUser{UserID: "u1", Password: "HASH", Role: "admin"}

	db.EXPECT().Login(mock.Anything, "bob").Return(dbOut, nil).Once()
	ph.EXPECT().CheckHash("pass", "HASH").Return(true).Once()
	tm.EXPECT().GenerateAccessToken("u1", "admin").Return("AT", nil).Once()
	tm.EXPECT().GenerateRefreshToken("u1", "admin").Return("RT", nil).Once()

	out, err := svc.Login(ctx, in)
	require.NoError(t, err)
	require.Equal(t, "AT", out.AccessToken)
	require.Equal(t, "RT", out.RefreshToken)
}

func TestAuthService_Login_NotFound_ReturnsInvalidCredentials(t *testing.T) {
	svc, db, ph, tm := newSvc(t)
	ctx := context.Background()

	db.EXPECT().Login(mock.Anything, "bob").Return(dto.AuthUser{}, errs.PgErrNotFound).Once()

	out, err := svc.Login(ctx, dto.LoginParams{Username: "bob", Password: "x"})
	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidCredentials)
	require.Equal(t, dto.LoginResult{}, out)

	ph.AssertNotCalled(t, "CheckHash", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateAccessToken", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateRefreshToken", mock.Anything, mock.Anything)
}

func TestAuthService_Login_Database_Err_ReturnsDbErr(t *testing.T) {
	svc, db, ph, tm := newSvc(t)
	ctx := context.Background()

	// err != errs.PgErrNotFound
	db.EXPECT().Login(mock.Anything, "bob").Return(dto.AuthUser{}, errs.PgErrDeadlock).Once()

	out, err := svc.Login(ctx, dto.LoginParams{Username: "bob", Password: "x"})
	require.Error(t, err)
	require.ErrorIs(t, err, errs.PgErrDeadlock)
	require.Equal(t, dto.LoginResult{}, out)

	ph.AssertNotCalled(t, "CheckHash", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateAccessToken", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateRefreshToken", mock.Anything, mock.Anything)
}

func TestAuthService_Incorrect_Hash_ReturnsInvalidCredentials(t *testing.T) {
	svc, db, ph, tm := newSvc(t)
	ctx := context.Background()

	in := dto.LoginParams{Username: "bob", Password: "pass"}
	dbOut := dto.AuthUser{UserID: "u1", Password: "ANOTHER_HASH", Role: "admin"}

	db.EXPECT().Login(mock.Anything, in.Username).Return(dbOut, nil).Once()
	ph.EXPECT().CheckHash(in.Password, dbOut.Password).Return(false).Once()

	out, err := svc.Login(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrInvalidCredentials)
	require.Equal(t, dto.LoginResult{}, out)

	tm.AssertNotCalled(t, "GenerateAccessToken", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateRefreshToken", mock.Anything, mock.Anything)
}

func TestAuthService_Login_AccessTokenError_ReturnsError(t *testing.T) {
	svc, db, ph, tm := newSvc(t)
	ctx := context.Background()

	in := dto.LoginParams{Username: "bob", Password: "pass"}
	dbOut := dto.AuthUser{UserID: "u1", Password: "HASH", Role: "admin"}
	jwtErr := errors.New("some JWT Err")

	db.EXPECT().Login(mock.Anything, in.Username).Return(dbOut, nil).Once()
	ph.EXPECT().CheckHash(in.Password, dbOut.Password).Return(true).Once()
	tm.EXPECT().GenerateAccessToken(dbOut.UserID, dbOut.Role).Return("", jwtErr).Once()

	out, err := svc.Login(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, jwtErr)
	require.Equal(t, dto.LoginResult{}, out)

	tm.AssertNotCalled(t, "GenerateRefreshToken", mock.Anything, mock.Anything)
}

func TestAuthService_Login_RefreshTokenError_ReturnsError(t *testing.T) {
	svc, db, ph, tm := newSvc(t)
	ctx := context.Background()

	in := dto.LoginParams{Username: "bob", Password: "pass"}
	dbOut := dto.AuthUser{UserID: "u1", Password: "HASH", Role: "admin"}
	jwtErr := errors.New("some JWT Err")

	db.EXPECT().Login(mock.Anything, in.Username).Return(dbOut, nil).Once()
	ph.EXPECT().CheckHash(in.Password, dbOut.Password).Return(true).Once()
	tm.EXPECT().GenerateAccessToken(dbOut.UserID, dbOut.Role).Return("access_token", nil).Once()
	tm.EXPECT().GenerateRefreshToken(dbOut.UserID, dbOut.Role).Return("", jwtErr).Once()

	out, err := svc.Login(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, jwtErr)
	require.Empty(t, out.AccessToken)
}
