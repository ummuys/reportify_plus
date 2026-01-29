package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/auth/internal/dto"
)

func TestAuthService_ListUsers_Success(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	dbOut := dto.ListUsersResult{}

	db.EXPECT().ListUsers(mock.Anything).Return(dbOut, nil).Once()

	res, err := svc.ListUsers(ctx)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)

}

func TestAuthService_ListUsers_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _, _ := newSvc(t)
	ctx := context.Background()

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().ListUsers(mock.Anything).Return(dto.ListUsersResult{}, dbErr).Once()

	res, err := svc.ListUsers(ctx)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.ListUsersResult{}, res)

}
