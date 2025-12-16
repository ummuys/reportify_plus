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
	ListSchemas(pCtx context.Context) (dto.ListSchemasResult, error)
	ListTables(pCtx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error)
	ListColumns(pCtx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error)

	// QUERIES
	SetCacheQueries(pCtx context.Context, cache map[string][]byte) error
	GetCacheQueries(pCtx context.Context) (map[string][]byte, error)

	Close()
}
