package repository

import (
	"context"
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

	logger := baseLogger.With().Str("component", "report-db").Logger()

	return &reportDB{
		logger: logger,
		pool:   pool,
	}, nil
}

func (db *reportDB) CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error) {
	db.logger.Debug().Str("evt", "call CreateReport").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	out := dto.CreateReportResult{}
	reportID := uuid.New().String()

	if err := db.pool.QueryRow(
		qctx,
		createReportQuery,
		reportID,
		in.AuthorID,
		in.Name,
		in.Comm,
		in.Query,
		in.Format,
		in.CSVSep,
	).Scan(&out.Status); err != nil {
		return dto.CreateReportResult{}, err
	}

	out.ReportID = reportID
	return out, nil
}

func (db *reportDB) RecreateReport(ctx context.Context, in dto.RecreateReportParams) (dto.RecreateReportResult, error) {
	db.logger.Debug().Str("evt", "call RecreateReport ").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	ct, err := db.pool.Exec(qctx, recreateReportQuery, in.AuthorID, in.ReportID)
	if err != nil {
		return dto.RecreateReportResult{}, err
	}

	if ct.RowsAffected() == 0 {
		return dto.RecreateReportResult{}, pgx.ErrNoRows
	}

	return dto.RecreateReportResult{ReportID: in.ReportID, Status: "CREATED"}, nil
}

func (db *reportDB) ListReports(ctx context.Context, in dto.ListReportsParams) (dto.ListReportsResult, error) {
	db.logger.Debug().Str("evt", "call ListReports").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	rows, err := db.pool.Query(qctx, listUserReportsQuery, in.AuthorID)
	if err != nil {
		return dto.ListReportsResult{}, err
	}
	defer rows.Close()

	out := dto.ListReportsResult{}
	for rows.Next() {
		rmd := dto.ReportMetadata{}
		if err := rows.Scan(
			&rmd.ReportID,
			&rmd.Name,
			&rmd.Comm,
			&rmd.Query,
			&rmd.Format,
			&rmd.CSVSep,
			&rmd.Status,
			&rmd.CreatedAt,
			&rmd.UpdatedAt,
			&rmd.FilePath,
			&rmd.ErrMsg,
		); err != nil {
			return dto.ListReportsResult{}, err
		}
		out.Reports = append(out.Reports, rmd)
	}

	return out, nil
}

func (db *reportDB) ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error) {
	db.logger.Debug().Str("evt", "call ReportStatus").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var status string
	if err := db.pool.QueryRow(
		qctx,
		reportStatusQuery,
		in.AuthorID,
		in.ReportID,
	).Scan(&status); err != nil {
		return dto.ReportStatusResult{}, err
	}

	return dto.ReportStatusResult{
		ReportID: in.ReportID,
		Status:   status,
	}, nil
}

func (db *reportDB) ReportInfo(ctx context.Context, in dto.ReportInfoParams) (dto.ReportInfoResult, error) {
	db.logger.Debug().Str("evt", "call ReportInfo").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var out dto.ReportInfoResult
	out.Report.ReportID = in.ReportID

	if err := db.pool.QueryRow(
		qctx,
		reportInfoQuery,
		in.AuthorID,
		in.ReportID,
	).Scan(
		&out.Report.Name,
		&out.Report.Comm,
		&out.Report.Query,
		&out.Report.Format,
		&out.Report.CSVSep,
		&out.Report.Status,
		&out.Report.CreatedAt,
		&out.Report.FilePath,
		&out.Report.ErrMsg,
	); err != nil {
		return dto.ReportInfoResult{}, err
	}

	return out, nil
}

func (db *reportDB) DeleteReports(ctx context.Context, in dto.DeleteReportsParams) error {
	db.logger.Debug().Str("evt", "call DeleteReports").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if _, err := db.pool.Exec(qctx, deleteUserReportsQuery, in.AuthorID); err != nil {
		return err
	}
	return nil
}

func (db *reportDB) DeleteReport(ctx context.Context, in dto.DeleteReportParams) (dto.DeleteReportResult, error) {
	db.logger.Debug().Str("evt", "call DeleteReport").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if _, err := db.pool.Exec(
		qctx,
		deleteUserReportQuery,
		in.AuthorID,
		in.ReportID,
	); err != nil {
		return dto.DeleteReportResult{}, err
	}

	return dto.DeleteReportResult{ReportID: in.ReportID}, nil
}

func (db *reportDB) Close() {
	db.pool.Close()
}
