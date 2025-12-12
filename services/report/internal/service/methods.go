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
