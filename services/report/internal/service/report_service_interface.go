package service

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportService interface {
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
	RecreateReport(ctx context.Context, in dto.RecreateReportParams) (dto.RecreateReportResult, error)
	ListReports(ctx context.Context, in dto.ListReportsParams) (dto.ListReportsResult, error)
	ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error)
	ReportInfo(ctx context.Context, in dto.ReportInfoParams) (dto.ReportInfoResult, error)
	DeleteReports(ctx context.Context, in dto.DeleteReportsParams) error
	DeleteReport(ctx context.Context, in dto.DeleteReportParams) (dto.DeleteReportResult, error)
}
