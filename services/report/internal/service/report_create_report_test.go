package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

func TestReportService_CreateReport_Success_CacheSetOk(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.CreateReportParams{AuthorID: "a1"}
	out := dto.CreateReportResult{ReportID: "r1", Status: "created"}

	db.EXPECT().CreateReport(mock.Anything, in).Return(out, nil).Once()
	cache.EXPECT().Set(mock.Anything, "r1", "created").Return(nil).Once()

	res, err := svc.CreateReport(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}

func TestReportService_CreateReport_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.CreateReportParams{AuthorID: "a1"}
	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().CreateReport(mock.Anything, in).Return(dto.CreateReportResult{}, dbErr).Once()

	res, err := svc.CreateReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.CreateReportResult{}, res)

	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestReportService_CreateReport_Success_CacheSetError_Ignored(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.CreateReportParams{AuthorID: "a1"}
	out := dto.CreateReportResult{ReportID: "r1", Status: "created"}
	cacheErr := errs.ErrPgDeadlock

	db.EXPECT().CreateReport(mock.Anything, in).Return(out, nil).Once()
	cache.EXPECT().Set(mock.Anything, "r1", "created").Return(cacheErr).Once()

	res, err := svc.CreateReport(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}
