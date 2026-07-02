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

func TestAuthService_DeleteUser_Success(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.DeleteUserParams{UserID: "u1"}
	dbOut := dto.DeleteUserResult{UserID: "u1"}

	db.EXPECT().DeleteUser(mock.Anything, in).Return(dbOut, nil).Once()

	res, err := svc.DeleteUser(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
}

func TestAuthService_DeleteUser_NotFound_ReturnsParsedPgError(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.DeleteUserParams{UserID: "u1"}

	dbErr := errs.ErrPgNotFound
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().DeleteUser(mock.Anything, in).Return(dto.DeleteUserResult{}, dbErr).Once()

	res, err := svc.DeleteUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.DeleteUserResult{}, res)
}

func TestAuthService_DeleteUser_InsufficientPrivilege_ReturnsParsedPgError(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.DeleteUserParams{UserID: "u1"}

	dbErr := errs.ErrPgInsufficientPrivilege
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().DeleteUser(mock.Anything, in).Return(dto.DeleteUserResult{}, dbErr).Once()

	res, err := svc.DeleteUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.DeleteUserResult{}, res)
}

func TestAuthService_DeleteUser_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.DeleteUserParams{UserID: "u1"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().DeleteUser(mock.Anything, in).Return(dto.DeleteUserResult{}, dbErr).Once()

	res, err := svc.DeleteUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.DeleteUserResult{}, res)
}

func TestAuthService_DeleteUser_DbErrorWrapped_ReturnsParsedPgError(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.DeleteUserParams{UserID: "u1"}

	base := errs.ErrPgDeadlock
	wrapped := errors.New("wrapped db err")
	expected := errs.ParsePgError(wrapped)

	db.EXPECT().DeleteUser(mock.Anything, in).Return(dto.DeleteUserResult{}, wrapped).Once()

	res, err := svc.DeleteUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.DeleteUserResult{}, res)
	require.NotErrorIs(t, err, base)
}
