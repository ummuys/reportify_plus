package service

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportService interface {
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
}
