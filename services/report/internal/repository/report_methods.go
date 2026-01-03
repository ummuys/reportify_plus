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
	db.logger.Debug().Str("evt", "call CreateReport").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	out := dto.CreateReportResult{}

	reportID := uuid.New().String()
	if err := db.pool.QueryRow(qctx, createReportQuery, reportID, in.AuthorID, in.Name, in.Comm, in.Query, in.Format, in.CSVSep).Scan(&out.Status); err != nil {
		db.logger.Error().Err(err).Str("evt", "call CreateReport").Msg("")
		return dto.CreateReportResult{}, err
	}

	out.ReportID = reportID

	return out, nil
}

func (db *reportDB) ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error) {
	db.logger.Debug().Str("evt", "call ListUserReports").Msg("")
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
		if err := rows.Scan(&rmd.ReportID, &rmd.Name,
			&rmd.Comm, &rmd.Query, &rmd.Format, &rmd.CSVSep, &rmd.Status,
			&rmd.CreatedAt, &rmd.FilePath, &rmd.ErrMsg); err != nil {
			db.logger.Error().Err(err).Str("evt", "call ListUserReports").Msg("")
			return dto.ListReportsResult{}, nil
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

	if err := db.pool.QueryRow(qctx, reportStatusQuery, in.AuthorID, in.ReportID).Scan(&status); err != nil {
		db.logger.Error().Err(err).Str("evt", "call ReportStatus").Msg("")
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

	err := db.pool.QueryRow(qctx, reportInfoQuery, in.AuthorID, in.ReportID).Scan(
		&out.Report.Name,
		&out.Report.Comm,
		&out.Report.Query,
		&out.Report.Format,
		&out.Report.CSVSep,
		&out.Report.Status,
		&out.Report.CreatedAt,
		&out.Report.FilePath,
		&out.Report.ErrMsg,
	)
	if err != nil {
		return dto.ReportInfoResult{}, err
	}

	return out, nil
}

func (db *reportDB) DeleteUserReports(ctx context.Context, in dto.DeleteUserReportsParams) error {
	db.logger.Debug().Str("evt", "call DeleteUserReports").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if _, err := db.pool.Exec(qctx, deleteUserReportsQuery, in.AuthorID); err != nil {
		db.logger.Error().Err(err).Str("evt", "call DeleteUserReports").Msg("")
		return err
	}
	return nil
}

func (db *reportDB) DeleteUserReport(ctx context.Context, in dto.DeleteUserReportParams) (dto.DeleteUserReportResult, error) {
	db.logger.Debug().Str("evt", "call DeleteUserReport").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if _, err := db.pool.Exec(qctx, deleteUserReportQuery, in.AuthorID, in.ReportID); err != nil {
		db.logger.Error().Err(err).Str("evt", "call DeleteUserReport").Msg("")
		return dto.DeleteUserReportResult{}, err
	}
	return dto.DeleteUserReportResult{ReportID: in.ReportID}, nil
}

func (db *reportDB) Close() {
	db.pool.Close()
}
