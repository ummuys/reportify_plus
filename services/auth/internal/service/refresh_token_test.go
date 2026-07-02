package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	pkg "github.com/ummuys/reportify/pkg/tm"
	"github.com/ummuys/reportify/services/auth/internal/dto"
)

func TestAuthService_RefreshToken_Success(t *testing.T) {
	svc, _, _, tm := newSvc(t)
	ctx := context.Background()

	in := dto.RefreshTokenParams{RefreshToken: "RT"}
	rc := pkg.RefreshClaims{UserID: "u1", Role: "admin"}

	tm.EXPECT().ValidateRefreshToken("RT").Return(rc, nil).Once()
	tm.EXPECT().GenerateAccessToken("u1", "admin").Return("AT", nil).Once()

	res, err := svc.RefreshToken(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dto.RefreshTokenResult{AccessToken: "AT"}, res)
}

func TestAuthService_RefreshToken_InvalidRefreshToken_ReturnsUnauthorized(t *testing.T) {
	svc, _, _, tm := newSvc(t)

	ctx := context.Background()

	in := dto.RefreshTokenParams{RefreshToken: "RT"}
	valErr := errors.New("bad token")

	tm.EXPECT().ValidateRefreshToken("RT").Return(pkg.RefreshClaims{}, valErr).Once()

	res, err := svc.RefreshToken(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, errs.ErrPgUnauthorized)
	require.Equal(t, dto.RefreshTokenResult{}, res)

	tm.AssertNotCalled(t, "GenerateAccessToken", mock.Anything, mock.Anything)
}

func TestAuthService_RefreshToken_AccessTokenError_ReturnsError(t *testing.T) {
	svc, _, _, tm := newSvc(t)
	ctx := context.Background()

	in := dto.RefreshTokenParams{RefreshToken: "RT"}
	rc := pkg.RefreshClaims{UserID: "u1", Role: "admin"}
	jwtErr := errors.New("jwt err")

	tm.EXPECT().ValidateRefreshToken("RT").Return(rc, nil).Once()
	tm.EXPECT().GenerateAccessToken("u1", "admin").Return("", jwtErr).Once()

	res, err := svc.RefreshToken(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, jwtErr)
	require.Equal(t, dto.RefreshTokenResult{}, res)
}
