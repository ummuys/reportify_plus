package adapter

import (
	"context"

	"github.com/rs/zerolog"
	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReportAdapter struct {
	logger zerolog.Logger
	reportv1.UnimplementedReportServiceServer
	svc service.ReportService
}

func NewReportAdapter(svc service.ReportService, baseLogger zerolog.Logger) *ReportAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &ReportAdapter{svc: svc, logger: logger}
}

func (ra *ReportAdapter) CreateReport(ctx context.Context, in *reportv1.CreateReportRequest) (*reportv1.CreateReportResponse, error) {
	ra.logger.Debug().Str("evt", "call CreateLogin").Msg("")
	out, err := ra.svc.CreateReport(ctx, dto.CreateReportParams{
		AuthorID: in.AuthorId,
		Name:     in.Name,
		Comm:     in.Comm,
		Query:    in.Query,
		Format:   in.Format,
		CSVSep:   in.CsvSep,
	})
	if err != nil {
		switch {
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &reportv1.CreateReportResponse{Uuid: out.UUID, Status: out.Status}, nil
}

func (ra *ReportAdapter) ReportStatus(ctx context.Context, in *reportv1.ReportStatusRequest) (*reportv1.ReportStatusResponse, error) {
	return &reportv1.ReportStatusResponse{}, nil
}
