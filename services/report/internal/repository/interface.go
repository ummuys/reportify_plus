package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportDB interface {
	CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error)
	ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error)
	ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error)

	// REPORT
	GetSchemas(pCtx context.Context) (map[string]string, error)
	GetTables(pCtx context.Context, schemaName string) (map[string]string, error)
	GetColumns(pCtx context.Context, schemaName, tableName string) (map[string]string, error)

	// QUERIES
	SetCacheQueries(pCtx context.Context, cache map[string][]byte) error
	GetCacheQueries(pCtx context.Context) (map[string][]byte, error)

	Close()
}
