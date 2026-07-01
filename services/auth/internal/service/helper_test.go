package service

import (
	"testing"

	"github.com/rs/zerolog"
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
