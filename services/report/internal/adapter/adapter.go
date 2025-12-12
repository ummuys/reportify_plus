package adapter

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	ra.logger.Debug().Str("evt", "call ReportStatus").Msg("")
	out, err := ra.svc.ReportStatus(ctx, dto.ReportStatusParams{UUID: in.Uuid})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	return &reportv1.ReportStatusResponse{Uuid: out.UUID, Status: out.Status}, nil
}

func (ra *ReportAdapter) ListUserReports(ctx context.Context, in *reportv1.ListUserReportsRequest) (*reportv1.ListUserReportsResponse, error) {
	ra.logger.Debug().Str("evt", "call ListUserReports").Msg("")
	out, err := ra.svc.ListUserReports(ctx, dto.ListUserReportsParams{AuthorID: in.AuthorId})
	if err != nil {
		switch {
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	resp := reportv1.ListUserReportsResponse{}
	resp.Reports = make([]*reportv1.ReportMetadata, 0, len(out.Reports))
	for _, r := range out.Reports {
		cat := timestamppb.New(r.CreatedAt)
		uat := timestamppb.New(r.UpdatedAt)

		resp.Reports = append(resp.Reports, &reportv1.ReportMetadata{
			ReportId:  r.ReportID,
			AuthorId:  r.AuthorID,
			Name:      r.Name,
			Comm:      r.Comm,
			Query:     r.Query,
			Format:    r.Format,
			CsvSep:    r.CSVSep,
			Status:    r.Status,
			CreatedAt: cat,
			UpdatedAt: uat,
			FilePath:  r.FilePath,
			ErrMsg:    r.ErrMsg,
		})
	}

	return &resp, nil
}
