package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/repository"
)

type reportService struct {
	db     repository.ReportDB
	logger zerolog.Logger
}

func NewReportService(db repository.ReportDB, baseLogger zerolog.Logger) ReportService {
	logger := baseLogger.With().Str("component", "svc").Logger()
	return &reportService{db: db, logger: logger}
}

func (rs *reportService) CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error) {
	rs.logger.Debug().Str("evt", "call CreateReport").Msg("")
	out, err := rs.db.CreateReport(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ListUserReports(ctx context.Context, in dto.ListUserReportsParams) (dto.ListReportsResult, error) {
	rs.logger.Debug().Str("evt", "call ListUserReports").Msg("")
	out, err := rs.db.ListUserReports(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error) {
	rs.logger.Debug().Str("evt", "call ReportStatus").Msg("")
	out, err := rs.db.ReportStatus(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ListSchemas(ctx context.Context) (dto.ListSchemasResult, error) {
	rs.logger.Debug().Str("evt", "call ListSchemas").Msg("")
	out, err := rs.db.ListSchemas(ctx)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ListTables(ctx context.Context, in dto.ListTablesParams) (dto.ListTablesResult, error) {
	rs.logger.Debug().
		Str("evt", "call ListTables").
		Str("schema", in.Schema).
		Msg("")

	out, err := rs.db.ListTables(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ListColumns(ctx context.Context, in dto.ListColumnsParams) (dto.ListColumnsResult, error) {
	rs.logger.Debug().
		Str("evt", "call ListColumns").
		Str("schema", in.Schema).
		Str("table", in.Table).
		Msg("")

	out, err := rs.db.ListColumns(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}
	return out, nil
}
