package service

import (
	"context"
	"encoding/json"

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

	value := dto.ReportCacheValue{
		Name:   in.Name,
		Comm:   in.Comm,
		Query:  in.Query,
		CSVSep: in.CSVSep,
	}

	bytes, err := json.Marshal(value)
	if err != nil {
		return dto.CreateReportResult{}, err
	}

	err = rs.cache.Set(ctx, in.AuthorID, bytes)
	if err != nil {
		return dto.CreateReportResult{}, err
	}

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
