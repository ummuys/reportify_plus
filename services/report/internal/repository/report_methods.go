package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
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

	logger := baseLogger.With().Str("component", "report-db").Logger()

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
