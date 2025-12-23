package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportDB interface {
	GetAllReports(ctx context.Context) (map[string][][]byte, error)
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
	ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error)
	ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error)
	Close()
}
