package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

func TestReportService_DeleteReports_Success(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	in := dto.DeleteReportsParams{AuthorID: "a1"}

	db.EXPECT().
		ListReports(mock.Anything, dto.ListReportsParams{AuthorID: "a1"}).
		Return(dto.ListReportsResult{}, nil).
		Once()

	db.EXPECT().DeleteReports(mock.Anything, in).Return(nil).Once()

	cache.EXPECT().
		Delete(mock.Anything).
		Return(nil).
		Once()

	err := svc.DeleteReports(ctx, in)
	require.NoError(t, err)
}

func TestReportService_DeleteReports_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.DeleteReportsParams{AuthorID: "a1"}
	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().
		ListReports(mock.Anything, dto.ListReportsParams{AuthorID: "a1"}).
		Return(dto.ListReportsResult{}, nil).
		Once()

	db.EXPECT().DeleteReports(mock.Anything, in).Return(dbErr).Once()

	err := svc.DeleteReports(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}

func TestReportService_DeleteReport_Success(t *testing.T) {
	svc, db, cache := newReportSvc(t)
	ctx := context.Background()

	reportID := "550e8400-e29b-41d4-a716-446655440000"
	in := dto.DeleteReportParams{ReportID: reportID, AuthorID: "a1"}
	out := dto.DeleteReportResult{ReportID: reportID}

	db.EXPECT().DeleteReport(mock.Anything, in).Return(out, nil).Once()

	cache.EXPECT().
		Delete(mock.Anything, reportID).
		Return(nil).
		Once()

	res, err := svc.DeleteReport(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}

func TestReportService_DeleteReport_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	reportID := "550e8400-e29b-41d4-a716-446655440000"
	in := dto.DeleteReportParams{ReportID: reportID, AuthorID: "a1"}
	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().DeleteReport(mock.Anything, in).Return(dto.DeleteReportResult{}, dbErr).Once()

	res, err := svc.DeleteReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.DeleteReportResult{}, res)
}
