package service

import (
	"context"
	"errors"

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
	out, err := ds.db.ListSchemas(ctx)
	if err != nil {
		perr := errs.ParsePgError(err)

		if !isExpectedDatasourceErr(perr) {
			ds.logger.Error().
				Err(err).
				Str("db-method", "ListSchemas").
				Msg("list schemas failed")
		}

		return out, perr
	}
	return out, nil
}

func (ds *datasourceService) ListTables(ctx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error) {
	out, err := ds.db.ListTables(ctx, in)
	if err != nil {
		perr := errs.ParsePgError(err)

		if !isExpectedDatasourceErr(perr) {
			ds.logger.Error().
				Err(err).
				Str("db-method", "ListTables").
				Str("schema", in.Schema).
				Msg("list tables failed")
		}

		return out, perr
	}
	return out, nil
}

func (ds *datasourceService) ListColumns(ctx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error) {
	out, err := ds.db.ListColumns(ctx, in)
	if err != nil {
		perr := errs.ParsePgError(err)

		if !isExpectedDatasourceErr(perr) {
			ds.logger.Error().
				Err(err).
				Str("db-method", "ListColumns").
				Str("schema", in.Schema).
				Str("table", in.Table).
				Msg("list columns failed")
		}

		return out, perr
	}
	return out, nil
}

func isExpectedDatasourceErr(err error) bool {
	return errors.Is(err, errs.ErrNotFound) ||
		errors.Is(err, errs.ErrUndefinedTable) ||
		errors.Is(err, errs.ErrUndefinedColumn) ||
		errors.Is(err, errs.ErrInsufficientPrivilege)
}
