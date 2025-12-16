package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/db"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type reportDB struct {
	logger zerolog.Logger
	pool   *pgxpool.Pool
}

func NewReportDB(ctx context.Context, baseLogger zerolog.Logger) (ReportDB, error) {
	qctx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	cfg, err := config.ParseReportDBEnv()
	if err != nil {
		return nil, err
	}

	pool, err := db.PoolFromConfig(qctx, cfg, "REPORT_DB")
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "report-db").Logger()

	return &reportDB{
		logger: logger,
		pool:   pool,
	}, nil
}

func (db *reportDB) GetReportInfo(ctx context.Context, in dto.GetReportInfoParams) (dto.GetReportInfoResult, error) {
	db.logger.Debug().Str("evt", "call GetReportInfo")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	var sep string
	out := dto.GetReportInfoResult{}
	if err := db.pool.QueryRow(qctx, ReportInfoQuery, in.UUID).Scan(&out.AuthorID,
		&out.Name, &out.Comm, &out.Query, &out.Format, &sep); err != nil {
		db.logger.Error().Err(err).Str("evt", "call GetReportInfo").Msg("")
		return dto.GetReportInfoResult{}, err
	}
	if sep != "" {
		out.CSVSep = sep[0]
	}
	return out, nil
}

func (db *reportDB) SetReportStatus(ctx context.Context, in dto.SetReportStatusParams) error {
	db.logger.Debug().Str("evt", "call SetReportStatus")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	query, vars := buildStatusQuery(in)

	if _, err := db.pool.Exec(qctx, query, vars...); err != nil {
		db.logger.Error().Err(err).Str("evt", "call SetReportStatus").Msg("")
		return err
	}
	return nil
}

func (db *reportDB) Close() {
	db.pool.Close()
}

func buildStatusQuery(in dto.SetReportStatusParams) (string, []any) {
	c := 1
	vars := []any{}

	base := `UPDATE report_metadata.report_requests
	SET
		updated_at = NOW()`

	base += fmt.Sprintf("\n status = %d", c)
	c++
	vars = append(vars, in.UpdateStatus)

	if in.FilePath != nil {
		base += fmt.Sprintf(",\n file_path = %d", c)
		c++
		vars = append(vars, in.FilePath)
	}

	if in.ErrMsg != nil {
		base += fmt.Sprintf(",\n error_message = %d", c)
		vars = append(vars, in.ErrMsg)
		c++
	}

	base += fmt.Sprintf("\n where report_id = %d", c)
	vars = append(vars, in.UUID)
	c++

	base += fmt.Sprintf("\n and status = %d", c)
	vars = append(vars, in.BeforeStatus)
	c++

	return base, vars
}
