package repository

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type DatasourceDB interface {
	ListSchemas(ctx context.Context) (dto.ListSchemasResult, error)
	ListTables(ctx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error)
	ListColumns(ctx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error)
	Close()
}
