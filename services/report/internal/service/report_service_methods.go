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
		return out, errs.ParsePgError(err)
	}

	rs.cache.Set(ctx, out.ReportID, out.Status)

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

	status, err := rs.cache.Get(ctx, in.ReportID)
	if err != nil && !errors.Is(err, redis.Nil) {
		rs.logger.Err(err).Msg("can't load data from cache")
	}

	if status != nil {
		return dto.ReportStatusResult{ReportID: in.ReportID, Status: *status}, nil
	}

	out, err := rs.db.ReportStatus(ctx, in)
	if err != nil {
		return out, errs.ParsePgError(err)
	}

	if err = rs.cache.Set(ctx, in.ReportID, out.Status); err != nil {
		rs.logger.Err(err).Msg("can't set data to cache")
	}

	return out, nil
}

func (rs *reportService) ReportInfo(ctx context.Context, in dto.ReportInfoParams) (dto.ReportInfoResult, error) {
	rs.logger.Debug().Str("evt", "call ReportInfo").Msg("")

	out, err := rs.db.ReportInfo(ctx, in)
	if err != nil {
		return dto.ReportInfoResult{}, errs.ParsePgError(err)
	}

	return out, nil
}

func (rs *reportService) DeleteUserReports(ctx context.Context, in dto.DeleteUserReportsParams) error {
	rs.logger.Debug().Str("evt", "call ReportInfo").Msg("")

	if err := rs.db.DeleteUserReports(ctx, in); err != nil {
		return errs.ParsePgError(err)
	}

	return nil
}

func (rs *reportService) DeleteUserReport(ctx context.Context, in dto.DeleteUserReportParams) (dto.DeleteUserReportResult, error) {
	rs.logger.Debug().Str("evt", "call ReportInfo").Msg("")

	out, err := rs.db.DeleteUserReport(ctx, in)
	if err != nil {
		return dto.DeleteUserReportResult{}, errs.ParsePgError(err)
	}

	return out, nil
}
