package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

func TestReportService_ListReports_Success(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.ListReportsParams{AuthorID: "a1"}
	dbOut := dto.ListReportsResult{}

	db.EXPECT().ListReports(mock.Anything, in).Return(dbOut, nil).Once()

	res, err := svc.ListReports(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
}

func TestReportService_ListReports_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.ListReportsParams{AuthorID: "a1"}
	dbOut := dto.ListReportsResult{}

	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().ListReports(mock.Anything, in).Return(dbOut, dbErr).Once()

	res, err := svc.ListReports(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dbOut, res)
}
