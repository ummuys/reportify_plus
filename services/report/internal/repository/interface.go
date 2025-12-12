package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportDB interface {
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
	ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error)
}
