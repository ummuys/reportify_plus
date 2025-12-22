package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/repository"
)

type datasourceService struct {
	db     repository.DatasourceDB
	logger zerolog.Logger
}

func NewDatasourceService(db repository.DatasourceDB, baseLogger zerolog.Logger) DatasourceService {
	logger := baseLogger.With().Str("component", "datasource-service").Logger()
	return &datasourceService{db: db, logger: logger}
}

func (ds *datasourceService) ListSchemas(ctx context.Context) (dto.ListSchemasResult, error) {
	ds.logger.Debug().Str("evt", "call ListSchemas").Msg("")

	out, err := ds.db.ListSchemas(ctx)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (ds *datasourceService) ListTables(ctx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error) {
	ds.logger.Debug().Str("evt", "call ListTables").Str("schema", in.Schema).Msg("")

	out, err := ds.db.ListTables(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (ds *datasourceService) ListColumns(ctx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error) {
	ds.logger.Debug().Str("evt", "call ListColumns").Str("schema", in.Schema).Str("table", in.Table).Msg("")

	out, err := ds.db.ListColumns(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}
