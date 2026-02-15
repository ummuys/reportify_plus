package service

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/cache"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/repository"
)

type reportService struct {
	db     repository.ReportDB
	cache  cache.ReportCache
	logger zerolog.Logger
}

func NewReportService(db repository.ReportDB, cache cache.ReportCache, baseLogger zerolog.Logger) ReportService {
	logger := baseLogger.With().Str("component", "report-service").Logger()
	return &reportService{db: db, cache: cache, logger: logger}
}

func (rs *reportService) CreateReport(ctx context.Context, in dto.CreateReportParams) (dto.CreateReportResult, error) {
	rs.logger.Debug().Str("evt", "call CreateReport").Msg("")

	out, err := rs.db.CreateReport(ctx, in)
	if err != nil {
		rs.logger.Error().
			Err(err).
			Str("db-method", "CreateReport").
			Str("author_id", in.AuthorID).
			Msg("create report failed")

		return out, errs.ParsePgError(err)
	}

	if err := rs.cache.Set(ctx, out.ReportID, out.Status); err != nil {
		rs.logger.Warn().
			Err(err).
			Str("cache-method", "Set").
			Str("report_id", out.ReportID).
			Msg("cache set failed")
	}

	return out, nil
}

func (rs *reportService) RecreateReport(ctx context.Context, in dto.RecreateReportParams) (dto.RecreateReportResult, error) {
	rs.logger.Debug().Str("evt", "call RecreateReport").Msg("")

	out, err := rs.db.RecreateReport(ctx, in)
	if err != nil {
		rs.logger.Error().
			Err(err).
			Str("db-method", "RecreateReport").
			Str("author_id", in.AuthorID).
			Str("report_id", in.ReportID).
			Msg("recreate report failed")
		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ListReports(ctx context.Context, in dto.ListReportsParams) (dto.ListReportsResult, error) {
	rs.logger.Debug().Str("evt", "call ListReports").Msg("")

	out, err := rs.db.ListReports(ctx, in)
	if err != nil {

		rs.logger.Error().
			Err(err).
			Str("db-method", "ListReports").
			Str("author_id", in.AuthorID).
			Msg("list reports failed")

		return out, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) ReportStatus(ctx context.Context, in dto.ReportStatusParams) (dto.ReportStatusResult, error) {
	rs.logger.Debug().Str("evt", "call ReportStatus").Msg("")

	status, err := rs.cache.Get(ctx, in.ReportID)
	if err != nil && !errors.Is(err, redis.Nil) {
		rs.logger.Warn().
			Err(err).
			Str("cache-method", "Get").
			Str("report_id", in.ReportID).
			Msg("cache get failed")
	}

	if status != nil {
		return dto.ReportStatusResult{ReportID: in.ReportID, Status: *status}, nil
	}

	out, err := rs.db.ReportStatus(ctx, in)
	if err != nil {

		rs.logger.Error().
			Err(err).
			Str("db-method", "ReportStatus").
			Str("report_id", in.ReportID).
			Msg("report status failed")

		return out, errs.ParsePgError(err)
	}

	if err := rs.cache.Set(ctx, in.ReportID, out.Status); err != nil {
		rs.logger.Warn().
			Err(err).
			Str("cahce-method", "Set").
			Str("report_id", in.ReportID).
			Msg("cache set failed")
	}

	return out, nil
}

func (rs *reportService) ReportInfo(ctx context.Context, in dto.ReportInfoParams) (dto.ReportInfoResult, error) {
	rs.logger.Debug().Str("evt", "call ReportInfo").Msg("")

	out, err := rs.db.ReportInfo(ctx, in)
	if err != nil {

		rs.logger.Error().
			Err(err).
			Str("db-method", "ReportInfo").
			Str("report_id", in.ReportID).
			Msg("report info failed")

		return dto.ReportInfoResult{}, errs.ParsePgError(err)
	}
	return out, nil
}

func (rs *reportService) DeleteReports(ctx context.Context, in dto.DeleteReportsParams) error {
	rs.logger.Debug().Str("evt", "call DeleteReports").Msg("")

	if err := rs.db.DeleteReports(ctx, in); err != nil {

		rs.logger.Error().
			Err(err).
			Str("db-method", "DeleteReports").
			Str("author_id", in.AuthorID).
			Msg("delete reports failed")

		return errs.ParsePgError(err)
	}
	return nil
}

func (rs *reportService) DeleteReport(ctx context.Context, in dto.DeleteReportParams) (dto.DeleteReportResult, error) {
	rs.logger.Debug().Str("evt", "call DeleteReport").Msg("")

	out, err := rs.db.DeleteReport(ctx, in)
	if err != nil {

		rs.logger.Error().
			Err(err).
			Str("db-method", "DeleteReport").
			Str("report_id", in.ReportID).
			Str("author_id", in.AuthorID).
			Msg("delete report failed")

		return dto.DeleteReportResult{}, errs.ParsePgError(err)
	}
	return out, nil
}
