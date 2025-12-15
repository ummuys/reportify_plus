package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type ReportDB interface {
	GetReportInfo(ctx context.Context, in dto.GetReportInfoParams) (dto.GetReportInfoResult, error)
	SetReportStatus(ctx context.Context, in dto.SetReportStatusParams) error
	SetReportFailedStatus(ctx context.Context, in dto.SetReportFailedStatusParams) error
	FinalizeReport(ctx context.Context, in dto.FinalizeReportParams) error
	Close()
}
