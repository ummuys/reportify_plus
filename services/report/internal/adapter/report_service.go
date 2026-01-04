package adapter

import (
	"context"
	"errors"

	"github.com/rs/zerolog"
	rsv1 "github.com/ummuys/reportify/api/pb/report/service/v1"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ReportAdapter struct {
	logger zerolog.Logger
	rsv1.UnimplementedReportServiceServer
	reportSvc service.ReportService
	tunnel    chan<- dto.KafkaMessage
}

func NewReportAdapter(reportSvc service.ReportService, tunnelIn chan dto.KafkaMessage, baseLogger zerolog.Logger) *ReportAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &ReportAdapter{reportSvc: reportSvc, tunnel: tunnelIn, logger: logger}
}

func (ra *ReportAdapter) CreateReport(ctx context.Context, in *rsv1.CreateReportRequest) (*rsv1.CreateReportResponse, error) {
	ra.logger.Debug().Str("evt", "call CreateReport").Msg("")

	out, err := ra.reportSvc.CreateReport(ctx, dto.CreateReportParams{
		AuthorID: in.AuthorId,
		Name:     in.Name,
		Comm:     in.Comm,
		Query:    in.Query,
		Format:   in.Format,
		CSVSep:   in.CsvSep,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ra.tunnel <- dto.KafkaMessage{
		Key:   nil,
		Value: []byte(out.ReportID),
	}

	ra.logger.Info().
		Str("report_id", out.ReportID).
		Str("author_id", in.AuthorId).
		Str("format", in.Format).
		Msg("report task created")

	return &rsv1.CreateReportResponse{ReportId: out.ReportID, Status: out.Status}, nil
}

func (ra *ReportAdapter) ReportStatus(ctx context.Context, in *rsv1.ReportStatusRequest) (*rsv1.ReportStatusResponse, error) {
	ra.logger.Debug().Str("evt", "call ReportStatus").Msg("")

	out, err := ra.reportSvc.ReportStatus(ctx, dto.ReportStatusParams{AuthorID: in.AuthorId, ReportID: in.ReportId})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &rsv1.ReportStatusResponse{
		ReportId: out.ReportID,
		Status:   out.Status,
	}, nil
}

func (ra *ReportAdapter) ListReports(ctx context.Context, in *rsv1.ListReportsRequest) (*rsv1.ListReportsResponse, error) {
	ra.logger.Debug().Str("evt", "call ListReports").Msg("")

	out, err := ra.reportSvc.ListReports(ctx, dto.ListReportsParams{AuthorID: in.AuthorId})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := rsv1.ListReportsResponse{
		Reports: make([]*rsv1.ReportMetadata, 0, len(out.Reports)),
	}

	for _, r := range out.Reports {
		resp.Reports = append(resp.Reports, &rsv1.ReportMetadata{
			ReportId:  r.ReportID,
			Name:      r.Name,
			Comm:      r.Comm,
			Query:     r.Query,
			Format:    r.Format,
			CsvSep:    r.CSVSep,
			Status:    r.Status,
			CreatedAt: timestamppb.New(r.CreatedAt),
			FilePath:  r.FilePath,
			ErrMsg:    r.ErrMsg,
		})
	}

	return &resp, nil
}

func (ra *ReportAdapter) ReportInfo(ctx context.Context, in *rsv1.ReportInfoRequest) (*rsv1.ReportInfoResponse, error) {
	ra.logger.Debug().Str("evt", "call ReportInfo").Msg("")

	out, err := ra.reportSvc.ReportInfo(ctx, dto.ReportInfoParams{
		AuthorID: in.AuthorId,
		ReportID: in.ReportId,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &rsv1.ReportInfoResponse{
		Report: &rsv1.ReportMetadata{
			ReportId:  out.Report.ReportID,
			Name:      out.Report.Name,
			Comm:      out.Report.Comm,
			Query:     out.Report.Query,
			Format:    out.Report.Format,
			CsvSep:    out.Report.CSVSep,
			Status:    out.Report.Status,
			CreatedAt: timestamppb.New(out.Report.CreatedAt),
			FilePath:  out.Report.FilePath,
			ErrMsg:    out.Report.ErrMsg,
		},
	}, nil
}

func (ra *ReportAdapter) DeleteReports(ctx context.Context, in *rsv1.DeleteReportsRequest) (*emptypb.Empty, error) {
	ra.logger.Debug().Str("evt", "call DeleteReports").Msg("")

	err := ra.reportSvc.DeleteReports(ctx, dto.DeleteReportsParams{
		AuthorID: in.AuthorId,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &emptypb.Empty{}, nil
}

func (ra *ReportAdapter) DeleteReport(ctx context.Context, in *rsv1.DeleteReportRequest) (*rsv1.DeleteReportResponse, error) {
	ra.logger.Debug().Str("evt", "call DeleteReport").Msg("")

	out, err := ra.reportSvc.DeleteReport(ctx, dto.DeleteReportParams{
		AuthorID: in.AuthorId,
		ReportID: in.ReportId,
	})
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrNotFound):
			return nil, status.Error(codes.NotFound, err.Error())
		default:
			return nil, status.Error(codes.Internal, err.Error())
		}
	}

	return &rsv1.DeleteReportResponse{
		ReportId: out.ReportID,
	}, nil
}
