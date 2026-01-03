package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportDB interface {
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
	ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error)
	ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error)
	ReportInfo(ctx context.Context, in dto.ReportInfoParams) (dto.ReportInfoResult, error)
	DeleteUserReports(ctx context.Context, in dto.DeleteUserReportsParams) error
	DeleteUserReport(ctx context.Context, in dto.DeleteUserReportParams) (dto.DeleteUserReportResult, error)
	Close()
}
