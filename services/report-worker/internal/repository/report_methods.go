package repository

import (
	"context"
	"fmt"
	"strings"
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
	qctx, cancel := context.WithTimeout(ctx, time.Second*10)
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
	db.logger.Debug().Str("evt", "call GetReportInfo").Msg("")
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
	db.logger.Debug().Str("evt", "call SetReportStatus").Msg("")
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
	args := make([]any, 0, 5)
	set := make([]string, 0, 4)

	// updated_at
	set = append(set, "updated_at = NOW()")

	// status
	args = append(args, in.UpdateStatus)
	set = append(set, fmt.Sprintf("status = $%d", len(args)))

	if in.FilePath != nil {
		args = append(args, *in.FilePath) // или sql.NullString
		set = append(set, fmt.Sprintf("file_path = $%d", len(args)))
	}

	if in.ErrMsg != nil {
		args = append(args, *in.ErrMsg)
		set = append(set, fmt.Sprintf("error_message = $%d", len(args)))
	}

	// WHERE
	args = append(args, in.UUID)
	whereReportID := fmt.Sprintf("report_id = $%d", len(args))

	args = append(args, in.BeforeStatus)
	whereBeforeStatus := fmt.Sprintf("status = $%d", len(args))

	q := fmt.Sprintf(`
UPDATE report_metadata.report_requests
SET %s
WHERE %s
  AND %s`,
		strings.Join(set, ",\n    "),
		whereReportID,
		whereBeforeStatus,
	)

	return q, args
}
