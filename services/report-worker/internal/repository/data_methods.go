package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/db"
	"github.com/ummuys/reportify/services/report-worker/internal/dto"
)

type dDB struct {
	logger zerolog.Logger
	pool   *pgxpool.Pool
}

func NewDataDB(ctx context.Context, baseLogger *zerolog.Logger) (DataDB, error) {
	qctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	cfg, err := config.ParseDataDBEnv()
	if err != nil {
		return nil, err
	}

	pool, err := db.PoolFromConfig(qctx, cfg, "report")
	if err != nil {
		return nil, err
	}

	logger := baseLogger.With().Str("component", "data-db").Logger()

	return &dDB{
		logger: logger,
		pool:   pool,
	}, nil
}

func (d *dDB) GetData(ctx context.Context, in dto.GetDataParams) (dto.GetDataResult, error) {
	d.logger.Debug().Str("evt", "call CreateReport").Str("Query", in.Query).Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*180)
	defer cancel()

	rows, err := d.pool.Query(qctx, in.Query)
	if err != nil {
		return dto.GetDataResult{}, err
	}
	defer rows.Close()

	fds := rows.FieldDescriptions()
	columns := make([]string, len(fds))
	for i := range fds {
		columns[i] = string(fds[i].Name)
	}

	var data [][]any

	for rows.Next() {
		vals, err := rows.Values()
		if err != nil {
			return dto.GetDataResult{}, err
		}

		row := make([]any, len(vals))
		for i, v := range vals {
			row[i] = convertPGIntoGo(v)
		}
		data = append(data, row)
	}

	if err := rows.Err(); err != nil {
		return dto.GetDataResult{}, err
	}

	return dto.GetDataResult{
		Columns: columns,
		Rows:    data,
	}, nil
}

func convertPGIntoGo(v any) any {
	switch val := v.(type) {
	case nil:
		return nil

	case pgtype.Numeric:
		if f, err := val.Float64Value(); err == nil {
			return f.Float64
		}
		return nil

	case pgtype.Int2:
		if val.Valid {
			return val.Int16
		}
		return nil

	case pgtype.Int4:
		if val.Valid {
			return val.Int32
		}
		return nil

	case pgtype.Int8:
		if val.Valid {
			return val.Int64
		}
		return nil

	case pgtype.Bool:
		if val.Valid {
			return val.Bool
		}
		return nil

	case pgtype.Text:
		if val.Valid {
			return val.String
		}
		return nil

	case pgtype.Timestamp:
		if val.Valid {
			return val.Time
		}
		return nil

	case pgtype.Date:
		if val.Valid {
			return val.Time
		}
		return nil

	default:
		return v
	}
}
