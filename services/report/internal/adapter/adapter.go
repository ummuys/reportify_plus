package adapter

import (
	"github.com/rs/zerolog"
	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"github.com/ummuys/reportify/services/report/internal/service"
)

type AuthAdapter struct {
	logger zerolog.Logger
	reportv1.UnimplementedReportServiceServer
	svc service.ReportService
}

func NewReportAdapter(svc service.ReportService, baseLogger zerolog.Logger) *AuthAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &AuthAdapter{svc: svc, logger: logger}
}
