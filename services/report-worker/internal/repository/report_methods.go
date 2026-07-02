package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
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
	if err := db.pool.QueryRow(qctx, ReportInfoQuery, in.ReportID).Scan(&out.AuthorID,
		&out.Name, &out.Comm, &out.Query, &out.Format, &sep); err != nil {
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
		return err
	}
	return nil
}

func (db *reportDB) PickAndMarkDeletingFile(ctx context.Context, in dto.PickAndMarkDeletingFileParams) (dto.PickAndMarkDeletingFileResult, error) {
	db.logger.Debug().Str("evt", "call PickAndMarkDeletingFile").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	rows, err := db.pool.Query(qctx, MarkAsDeletingFileAndGetReportIdQuery, time.Now().Add(-in.TimeLife), in.CountBatch)
	if err != nil {
		return dto.PickAndMarkDeletingFileResult{}, err
	}
	defer rows.Close()

	var out dto.PickAndMarkDeletingFileResult
	out.ReportsId = make([]string, 0, in.CountBatch)

	for rows.Next() {
		var rid uuid.UUID
		if err := rows.Scan(&rid); err != nil {
			return dto.PickAndMarkDeletingFileResult{}, err
		}
		out.ReportsId = append(out.ReportsId, rid.String())
	}

	if err := rows.Err(); err != nil {
		return dto.PickAndMarkDeletingFileResult{}, err
	}

	return out, nil
}

func (db *reportDB) MarkArchived(ctx context.Context, in dto.MarkArchivedParams) error {
	db.logger.Debug().Str("evt", "call MarkArchived").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	if in.Error == nil {
		if _, err := db.pool.Exec(qctx, MarkAsFileDeletedQuery, in.ReportsId); err != nil {
			return err
		}
		return nil
	}

	if _, err := db.pool.Exec(qctx, MarkAsErrorFileDeletedQuery, in.Error.Error(), in.ReportsId); err != nil {
		return err
	}
	return nil
}

func (db *reportDB) Close() {
	db.pool.Close()
}
