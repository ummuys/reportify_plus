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

func TestAuthService_UpdateUser_Success_NoPasswordChange(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.UpdateUserParams{UserID: "u1", Username: "Bob", Password: "", Role: "Admin"}
	dbOut := dto.UpdateUserResult{UserID: "u1"}

	db.EXPECT().UpdateUser(mock.Anything, in).Return(dbOut, nil).Once()

	res, err := svc.UpdateUser(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
	require.Equal(t, "", in.Password)
}

func TestAuthService_UpdateUser_Success_WithPasswordChange(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.UpdateUserParams{UserID: "u1", Username: "Bob", Password: "123", Role: "Admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbOut := dto.UpdateUserResult{UserID: "u1"}

	ph.EXPECT().Hash("123").Return("HASH", nil).Once()
	db.EXPECT().UpdateUser(mock.Anything, afterHash).Return(dbOut, nil).Once()

	res, err := svc.UpdateUser(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
	require.Equal(t, "123", in.Password)
}

func TestAuthService_UpdateUser_HashError_ReturnsError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.UpdateUserParams{UserID: "u1", Username: "Bob", Password: "123", Role: "Admin"}
	hashErr := errors.New("hash err")

	ph.EXPECT().Hash("123").Return("", hashErr).Once()

	res, err := svc.UpdateUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, hashErr)
	require.Equal(t, dto.UpdateUserResult{}, res)
	require.Equal(t, "123", in.Password)

	db.AssertNotCalled(t, "UpdateUser", mock.Anything, mock.Anything)
}

func TestAuthService_UpdateUser_DbError_NoPasswordChange_ReturnsParsedPgError(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	in := dto.UpdateUserParams{UserID: "u1", Username: "Bob", Password: "", Role: "Admin"}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().UpdateUser(mock.Anything, in).Return(dto.UpdateUserResult{}, dbErr).Once()

	res, err := svc.UpdateUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.UpdateUserResult{}, res)
}

func TestAuthService_UpdateUser_DbError_WithPasswordChange_ReturnsParsedPgError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.UpdateUserParams{UserID: "u1", Username: "Bob", Password: "123", Role: "Admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	ph.EXPECT().Hash("123").Return("HASH", nil).Once()
	db.EXPECT().UpdateUser(mock.Anything, afterHash).Return(dto.UpdateUserResult{}, dbErr).Once()

	res, err := svc.UpdateUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.UpdateUserResult{}, res)
	require.Equal(t, "123", in.Password)
}
