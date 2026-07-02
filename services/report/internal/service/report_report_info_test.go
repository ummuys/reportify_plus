package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

func TestReportService_ReportInfo_Success(t *testing.T) {
	svc, db, _ := newReportSvc(t)

	ctx := context.Background()

	in := dto.ReportInfoParams{ReportID: "r1"}
	dbOut := dto.ReportInfoResult{}

	db.EXPECT().ReportInfo(mock.Anything, in).Return(dbOut, nil).Once()

	res, err := svc.ReportInfo(ctx, in)
	require.NoError(t, err)
	require.Equal(t, dbOut, res)
}

func TestReportService_ReportInfo_DbError_ReturnsParsedPgError(t *testing.T) {
	svc, db, _ := newReportSvc(t)
	ctx := context.Background()

	in := dto.ReportInfoParams{ReportID: "r1"}
	dbErr := errs.ErrPgDeadlock
	expected := errs.ParsePgError(dbErr)

	db.EXPECT().ReportInfo(mock.Anything, in).Return(dto.ReportInfoResult{}, dbErr).Once()

	res, err := svc.ReportInfo(ctx, in)
	require.Error(t, err)
	require.ErrorIs(t, err, expected)
	require.Equal(t, dto.ReportInfoResult{}, res)
}
