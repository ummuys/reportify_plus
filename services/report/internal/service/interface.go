package service

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportService interface {
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
	ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error)
	ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error)
}
