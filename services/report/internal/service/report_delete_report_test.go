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
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.DeleteReportsParams{AuthorID: "a1"}

	db.EXPECT().DeleteReports(mock.Anything, in).Return(nil).Once()

	err := svc.DeleteReports(ctx, in)
	require.NoError(t, err)
}

func TestReportService_DeleteReports_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.DeleteReportsParams{AuthorID: "a1"}
	dbErr := errs.PgErrDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().DeleteReports(mock.Anything, in).Return(dbErr).Once()

	err := svc.DeleteReports(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
}

func TestReportService_DeleteReport_Success(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.DeleteReportParams{ReportID: "r1", AuthorID: "a1"}
	out := dto.DeleteReportResult{ReportID: "r1"}

	db.EXPECT().DeleteReport(mock.Anything, in).Return(out, nil).Once()

	res, err := svc.DeleteReport(ctx, in)
	require.NoError(t, err)
	require.Equal(t, out, res)
}

func TestReportService_DeleteReport_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.DeleteReportParams{ReportID: "r1", AuthorID: "a1"}
	dbErr := errs.PgErrDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().DeleteReport(mock.Anything, in).Return(dto.DeleteReportResult{}, dbErr).Once()

	res, err := svc.DeleteReport(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.DeleteReportResult{}, res)
}
