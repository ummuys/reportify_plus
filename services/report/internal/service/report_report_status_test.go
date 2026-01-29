package service

import (
	"context"
	"errors"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

func TestReportService_ReportStatus_CacheHit_ReturnsFromCache(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.ReportStatusParams{ReportID: "r1"}
	st := "done"

	cache.EXPECT().Get(mock.Anything, "r1").Return(&st, nil).Once()

	res, err := svc.ReportStatus(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dto.ReportStatusResult{ReportID: "r1", Status: "done"}, res)

	db.AssertNotCalled(t, "ReportStatus", mock.Anything, mock.Anything)
	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestReportService_ReportStatus_CacheMiss_DbSuccess_CacheSetOk(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.ReportStatusParams{ReportID: "r1"}
	out := dto.ReportStatusResult{ReportID: "r1", Status: "running"}

	cache.EXPECT().Get(mock.Anything, "r1").Return(nil, redis.Nil).Once()
	db.EXPECT().ReportStatus(mock.Anything, in).Return(out, nil).Once()
	cache.EXPECT().Set(mock.Anything, "r1", "running").Return(nil).Once()

	res, err := svc.ReportStatus(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}

func TestReportService_ReportStatus_CacheError_DbSuccess_CacheSetOk(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.ReportStatusParams{ReportID: "r1"}
	out := dto.ReportStatusResult{ReportID: "r1", Status: "running"}
	cacheErr := errors.New("cache down")

	cache.EXPECT().Get(mock.Anything, "r1").Return(nil, cacheErr).Once()
	db.EXPECT().ReportStatus(mock.Anything, in).Return(out, nil).Once()
	cache.EXPECT().Set(mock.Anything, "r1", "running").Return(nil).Once()

	res, err := svc.ReportStatus(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}

func TestReportService_ReportStatus_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.ReportStatusParams{ReportID: "r1"}
	dbOut := dto.ReportStatusResult{}
	dbErr := errs.PgErrDeadlock
	expected := errs.ParsePgError(dbErr)

	cache.EXPECT().Get(mock.Anything, "r1").Return(nil, redis.Nil).Once()
	db.EXPECT().ReportStatus(mock.Anything, in).Return(dbOut, dbErr).Once()

	res, err := svc.ReportStatus(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dbOut, res)

	cache.AssertNotCalled(t, "Set", mock.Anything, mock.Anything, mock.Anything)
}

func TestReportService_ReportStatus_CacheSetError_Ignored(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.ReportStatusParams{ReportID: "r1"}
	out := dto.ReportStatusResult{ReportID: "r1", Status: "running"}
	cacheSetErr := errors.New("cache set err")

	cache.EXPECT().Get(mock.Anything, "r1").Return(nil, redis.Nil).Once()
	db.EXPECT().ReportStatus(mock.Anything, in).Return(out, nil).Once()
	cache.EXPECT().Set(mock.Anything, "r1", "running").Return(cacheSetErr).Once()

	res, err := svc.ReportStatus(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}
