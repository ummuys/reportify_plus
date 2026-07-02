package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/db"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

type datasourceDB struct {
	logger zerolog.Logger
	pool   *pgxpool.Pool
}

func NewDatasourceDB(ctx context.Context, baseLogger zerolog.Logger) (DatasourceDB, error) {
	dctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cfg, err := config.ParseDatasourceDBEnv()
	if err != nil {
		return nil, err
	}

	pool, err := db.PoolFromConfig(dctx, cfg, "DATASOURCE_DB")
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "datasource-db").Logger()

	return &datasourceDB{
		logger: logger,
		pool:   pool,
	}, nil
}

func (db *datasourceDB) ListSchemas(pCtx context.Context) (dto.ListSchemasResult, error) {
	db.logger.Debug().Str("evt", "call ListSchemas").Msg("")

	ctx, cancel := context.WithTimeout(pCtx, 2*time.Second)
	defer cancel()

	rows, err := db.pool.Query(ctx, schemaWithCommentQuery)
	if err != nil {
		return dto.ListSchemasResult{}, err
	}
	defer rows.Close()

	out := dto.ListSchemasResult{Schemas: make([]dto.Schema, 0, 16)}
	for rows.Next() {
		var s dto.Schema
		if err := rows.Scan(&s.Name, &s.Comment); err != nil {
			return dto.ListSchemasResult{}, err
		}
		out.Schemas = append(out.Schemas, s)
	}

	if err := rows.Err(); err != nil {
		return dto.ListSchemasResult{}, err
	}

	return out, nil
}

func (db *datasourceDB) ListTables(pCtx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error) {
	db.logger.Debug().Str("evt", "call ListTables").Str("schema", in.Schema).Msg("")

	ctx, cancel := context.WithTimeout(pCtx, 2*time.Second)
	defer cancel()

	rows, err := db.pool.Query(ctx, tablesWithCommentQuery, in.Schema)
	if err != nil {
		return dto.ListTablesResult{}, err
	}
	defer rows.Close()

	out := dto.ListTablesResult{Tables: make([]dto.Table, 0, 128)}
	for rows.Next() {
		var t dto.Table
		if err := rows.Scan(&t.Name, &t.Comment); err != nil {
			return dto.ListTablesResult{}, err
		}
		out.Tables = append(out.Tables, t)
	}

	if err := rows.Err(); err != nil {
		return dto.ListTablesResult{}, err
	}

	return out, nil
}

func (db *datasourceDB) ListColumns(pCtx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error) {
	db.logger.Debug().Str("evt", "call ListColumns").Str("schema", in.Schema).Str("table", in.Table).Msg("")

	ctx, cancel := context.WithTimeout(pCtx, 5*time.Second)
	defer cancel()

	rows, err := db.pool.Query(ctx, columnsWithCommentQuery, in.Schema, in.Table)
	if err != nil {
		return dto.ListColumnsResult{}, err
	}
	defer rows.Close()

	out := dto.ListColumnsResult{Columns: make([]dto.Column, 0, 256)}
	for rows.Next() {
		var c dto.Column
		if err := rows.Scan(&c.Name, &c.Comment); err != nil {
			return dto.ListColumnsResult{}, err
		}
		out.Columns = append(out.Columns, c)
	}

	if err := rows.Err(); err != nil {
		return dto.ListColumnsResult{}, err
	}

	return out, nil
}

func (db *datasourceDB) Close() {
	db.pool.Close()
}
