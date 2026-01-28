package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/auth/internal/dto"
	"github.com/ummuys/reportify/services/auth/internal/mocks"
)

func newSvc(t *testing.T) (AuthService, *mocks.MockAuthDB, *mocks.MockPasswordHasher, *mocks.MockTokenManager) {
	t.Helper()

	db := mocks.NewMockAuthDB(t)
	ph := mocks.NewMockPasswordHasher(t)
	tm := mocks.NewMockTokenManager(t)

	svc := NewAuthService(ph, tm, db, zerolog.Nop())

	return svc, db, ph, tm
}

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
	require.True(t, errors.Is(err, errs.ErrInvalidCredentials))
	require.Empty(t, out.AccessToken)

	ph.AssertNotCalled(t, "CheckHash", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateAccessToken", mock.Anything, mock.Anything)
	tm.AssertNotCalled(t, "GenerateRefreshToken", mock.Anything, mock.Anything)
}
