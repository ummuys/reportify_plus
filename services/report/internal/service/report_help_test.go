package service

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/services/report/internal/mocks"
)

func newReportSvc(t *testing.T) (ReportService, *mocks.MockReportDB, *mocks.MockReportCache) {
	t.Helper()

	db := mocks.NewMockReportDB(t)
	cache := mocks.NewMockReportCache(t)

	svc := NewReportService(
		db,
		cache,
		zerolog.Nop(),
	)
	return svc, db, cache
}
