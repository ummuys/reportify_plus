package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/db"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

type reportDB struct {
	logger zerolog.Logger
	pool   *pgxpool.Pool
}

func NewReportDB(ctx context.Context, baseLogger zerolog.Logger) (ReportDB, error) {
	dctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	cfg, err := config.ParseReportDBEnv()
	if err != nil {
		return nil, err
	}

	pool, err := db.PoolFromConfig(dctx, cfg, "REPORT_DB")
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "db").Logger()

	return &reportDB{
		logger: logger,
		pool:   pool,
	}, nil
}

func (db *reportDB) CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error) {
	db.logger.Debug().Str("evt", "call CreateReport")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	out := dto.CreateReportResult{}

	uuid := uuid.New().String()
	if err := db.pool.QueryRow(qctx, createReportQuery, uuid, in.AuthorID, in.Name, in.Comm, in.Query, in.Format, in.CSVSep).Scan(&out.Status); err != nil {
		db.logger.Error().Err(err).Str("evt", "call CreateReport").Msg("")
		return dto.CreateReportResult{}, err
	}

	out.UUID = uuid

	return out, nil
}

func (db *reportDB) ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error) {
	db.logger.Debug().Str("evt", "call ListUserReports")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	rows, err := db.pool.Query(qctx, listUserReportsQuery, in.AuthorID)
	if err != nil {
		db.logger.Error().Err(err).Str("evt", "call ListUserReports").Msg("")
		return dto.ListReportsResult{}, err
	}
	defer rows.Close()

	out := dto.ListReportsResult{}
	for rows.Next() {
		rmd := dto.ReportMetadata{}
		if err := rows.Scan(&rmd.ReportID, &rmd.AuthorID, &rmd.Name,
			&rmd.Comm, &rmd.Query, &rmd.Format, &rmd.CSVSep, &rmd.Status,
			&rmd.CreatedAt, &rmd.UpdatedAt, &rmd.FilePath, &rmd.ErrMsg); err != nil {
			db.logger.Error().Err(err).Str("evt", "call ListUserReports").Msg("")
			return dto.ListReportsResult{}, nil
		}
		out.Reports = append(out.Reports, rmd)
	}

	return out, nil
}

func (db *reportDB) ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error) {
	db.logger.Debug().Str("evt", "call ReportStatus")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	out := dto.ReportStatusResult{}

	if err := db.pool.QueryRow(qctx, getReportStatusQuery, in.UUID).Scan(&out.Status, &out.ErrMsg, &out.FilePath); err != nil {
		db.logger.Error().Err(err).Str("evt", "call CreateReport").Msg("")
		return dto.ReportStatusResult{}, err
	}

	out.UUID = in.UUID
	return out, nil
}

func (db *reportDB) Close() {
	db.pool.Close()
}

func (db *reportDB) ListSchemas(pCtx context.Context) (dto.ListSchemasResult, error) {
	db.logger.Debug().Str("evt", "call ListSchemas").Msg("")

	ctx, cancel := context.WithTimeout(pCtx, time.Second*2)
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
			db.logger.Error().Err(err).Str("evt", "call ListSchemas").Msg("")
			return dto.ListSchemasResult{}, err
		}
		out.Schemas = append(out.Schemas, s)
	}

	if err := rows.Err(); err != nil {
		return dto.ListSchemasResult{}, err
	}

	return out, nil
}

func (db *reportDB) ListTables(pCtx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error) {
	db.logger.Debug().
		Str("evt", "call ListTables").
		Str("schema", in.Schema).
		Msg("")

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
			db.logger.Error().Err(err).Str("evt", "call ListTables").Msg("")
			return dto.ListTablesResult{}, err
		}
		out.Tables = append(out.Tables, t)
	}
	if err := rows.Err(); err != nil {
		return dto.ListTablesResult{}, err
	}

	return out, nil
}

func (db *reportDB) ListColumns(pCtx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error) {
	db.logger.Debug().
		Str("evt", "call ListColumns").
		Str("schema", in.Schema).
		Str("table", in.Table).
		Msg("")

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
			db.logger.Error().Err(err).Str("evt", "call ListColumns").Msg("")
			return dto.ListColumnsResult{}, err
		}
		out.Columns = append(out.Columns, c)
	}
	if err := rows.Err(); err != nil {
		return dto.ListColumnsResult{}, err
	}

	return out, nil
}

// TO CHECK AND FIX
func (d *reportDB) SetCacheQueries(pCtx context.Context, cache map[string][]byte) (err error) {
	d.logger.Debug().Str("evt", "call SetCacheQuerys").Msg("")
	ctx, cancel := context.WithTimeout(pCtx, 30*time.Second)
	defer cancel()

	var rows pgx.Rows
	allID := `select user_id from identity.users`
	rows, err = d.pool.Query(ctx, allID)
	if err != nil {
		return err
	}
	defer rows.Close()

	usersID := make([]string, 0, 64)
	for rows.Next() {
		var uid string
		if err = rows.Scan(&uid); err != nil {
			return err
		}
		usersID = append(usersID, uid)
	}

	if err := rows.Err(); err != nil {
		return err
	}

	var tx pgx.Tx
	tx, err = d.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rbErr := tx.Rollback(ctx); rbErr != nil && !errors.Is(rbErr, pgx.ErrTxClosed) {
				d.logger.Error().Err(rbErr).Msg("rollback failed")
			}
		}
	}()

	b := &pgx.Batch{}

	for _, uid := range usersID {
		val, ok := cache[uid]
		if !ok {
			b.Queue(setCacheQuery, uid, []byte("[]"))
		} else {
			b.Queue(setCacheQuery, uid, val)
		}
	}

	br := tx.SendBatch(ctx, b)

	for range usersID {
		if _, err = br.Exec(); err != nil {
			_ = br.Close()
			return err
		}
	}

	if err = br.Close(); err != nil {
		return
	}

	if err = tx.Commit(ctx); err != nil {
		return
	}

	return nil
}

// TO CHECK AND FIX
func (d *reportDB) GetCacheQueries(pCtx context.Context) (map[string][]byte, error) {
	d.logger.Debug().Str("evt", "call SaveCacheQuerys").Msg("")
	ctx, cancel := context.WithTimeout(pCtx, 5*time.Second)
	defer cancel()

	rows, err := d.pool.Query(ctx, getCacheQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	quer := make(map[string][]byte)
	for rows.Next() {
		var (
			key   string
			value []byte
		)
		err := rows.Scan(&key, &value)
		if err != nil {
			return nil, err
		}
		quer[key] = value
	}
	return quer, nil
}
