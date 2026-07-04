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

func TestAuthService_CreateUser_Success(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Bob", Password: "123", Role: "admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbOut := dto.CreateUserResult{UserID: "u1"}

	ph.EXPECT().Hash(in.Password).Return("HASH", nil).Once()
	db.EXPECT().CreateUser(mock.Anything, afterHash).Return(dbOut, nil).Once()

	res, err := svc.CreateUser(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
	require.Equal(t, "123", in.Password)
}

func TestAuthService_CreateUser_HashError_ReturnsError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Bob", Password: "123", Role: "admin"}
	hashErr := errors.New("hash err")

	ph.EXPECT().Hash(in.Password).Return("", hashErr).Once()

	res, err := svc.CreateUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, hashErr)
	require.Equal(t, dto.CreateUserResult{}, res)
	require.Equal(t, "123", in.Password)

	db.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
}

func TestAuthService_CreateUser_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Bob", Password: "123", Role: "admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	ph.EXPECT().Hash(in.Password).Return("HASH", nil).Once()
	db.EXPECT().CreateUser(mock.Anything, afterHash).Return(dto.CreateUserResult{}, dbErr).Once()

	res, err := svc.CreateUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.CreateUserResult{}, res)
	require.Equal(t, "123", in.Password)
}

func TestAuthService_CreateUser_UserAlreadyExists_ReturnsParsedPgError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Bob", Password: "123", Role: "admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbErr := errs.ErrPgDuplicate
	expected := errs.ParsePgError(dbErr)

	ph.EXPECT().Hash(in.Password).Return("HASH", nil).Once()
	db.EXPECT().CreateUser(mock.Anything, afterHash).Return(dto.CreateUserResult{}, dbErr).Once()

	res, err := svc.CreateUser(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.CreateUserResult{}, res)
}

func TestAuthService_CreateBaseAdmin_Success(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Admin", Password: "123", Role: "admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbOut := dto.CreateUserResult{UserID: "u1"}

	ph.EXPECT().Hash(in.Password).Return("HASH", nil).Once()
	db.EXPECT().CreateUser(mock.Anything, afterHash).Return(dbOut, nil).Once()
	db.EXPECT().SetAdminUUID("u1").Once()

	err := svc.CreateBaseAdmin(ctx, in)
	require.NoError(t, err)
	require.Equal(t, "123", in.Password)
}

func TestAuthService_CreateBaseAdmin_HashError_ReturnsError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Admin", Password: "123", Role: "admin"}
	hashErr := errors.New("hash err")

	ph.EXPECT().Hash(in.Password).Return("", hashErr).Once()

	err := svc.CreateBaseAdmin(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, hashErr)
	require.Equal(t, "123", in.Password)

	db.AssertNotCalled(t, "CreateUser", mock.Anything, mock.Anything)
	db.AssertNotCalled(t, "SetAdminUUID", mock.Anything)
}

func TestAuthService_CreateBaseAdmin_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Admin", Password: "123", Role: "admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	ph.EXPECT().Hash(in.Password).Return("HASH", nil).Once()
	db.EXPECT().CreateUser(mock.Anything, afterHash).Return(dto.CreateUserResult{}, dbErr).Once()

	err := svc.CreateBaseAdmin(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, "123", in.Password)

	db.AssertNotCalled(t, "SetAdminUUID", mock.Anything)
}

func TestAuthService_CreateBaseAdmin_UserAlreadyExists_ReturnsParsedPgError(t *testing.T) {
	svc, db, ph, _ := newSvc(t)
	ctx := context.Background()

	in := dto.CreateUserParams{Username: "Admin", Password: "123", Role: "admin"}
	afterHash := in
	afterHash.Password = "HASH"

	dbErr := errs.ErrPgDuplicate
	expected := errs.ParsePgError(dbErr)

	ph.EXPECT().Hash(in.Password).Return("HASH", nil).Once()
	db.EXPECT().CreateUser(mock.Anything, afterHash).Return(dto.CreateUserResult{}, dbErr).Once()

	err := svc.CreateBaseAdmin(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)

	db.AssertNotCalled(t, "SetAdminUUID", mock.Anything)
}
