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

	// For kafka
	tunnel chan<- dto.KafkaMessage
}

func NewReportAdapter(svc service.ReportService, tunnelIn chan dto.KafkaMessage, baseLogger zerolog.Logger) *ReportAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &ReportAdapter{svc: svc, tunnel: tunnelIn, logger: logger}
}

func (ra *ReportAdapter) CreateReport(ctx context.Context, in *reportv1.CreateReportRequest) (*reportv1.CreateReportResponse, error) {
	ra.logger.Debug().Str("evt", "call CreateReport").Msg("")
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

	ra.tunnel <- dto.KafkaMessage{
		Key:   nil,
		Value: []byte(out.UUID),
	}

	ra.logger.Info().Str("uuid", out.UUID).Str("author_id", in.AuthorId).Str("format", in.Format).Msg("report task created")

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

func (ra *ReportAdapter) ListSchemas(ctx context.Context, in *reportv1.ListSchemasRequest) (*reportv1.ListSchemasResponse, error) {
	ra.logger.Debug().Str("evt", "call ListSchemas").Msg("")

	out, err := ra.svc.ListSchemas(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &reportv1.ListSchemasResponse{
		Schemas: make([]*reportv1.Schema, 0, len(out.Schemas)),
	}

	for _, s := range out.Schemas {
		resp.Schemas = append(resp.Schemas, &reportv1.Schema{
			SchemaName: s.Name,
			SchemaComm: s.Comment,
		})
	}

	return resp, nil
}

func (ra *ReportAdapter) ListTables(ctx context.Context, in *reportv1.ListTablesRequest) (*reportv1.ListTablesResponse, error) {
	ra.logger.Debug().
		Str("evt", "call ListTables").
		Str("schema", in.SchemaName).
		Msg("")

	out, err := ra.svc.ListTables(ctx, dto.ListTablesParams{Schema: in.SchemaName})
	if err != nil {
		// если хочешь отдельно обрабатывать "нет такой схемы" — добавь errs.ErrNotFound
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &reportv1.ListTablesResponse{
		Tables: make([]*reportv1.Table, 0, len(out.Tables)),
	}

	for _, t := range out.Tables {
		resp.Tables = append(resp.Tables, &reportv1.Table{
			TableName: t.Name,
			TableComm: t.Comment,
		})
	}

	return resp, nil
}

func (ra *ReportAdapter) ListColumns(ctx context.Context, in *reportv1.ListColumnsRequest) (*reportv1.ListColumnsResponse, error) {
	ra.logger.Debug().
		Str("evt", "call ListColumns").
		Str("schema", in.SchemaName).
		Str("table", in.TableName).
		Msg("")

	out, err := ra.svc.ListColumns(ctx, dto.ListColumnsParams{
		Schema: in.SchemaName,
		Table:  in.TableName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &reportv1.ListColumnsResponse{
		Columns: make([]*reportv1.Column, 0, len(out.Columns)),
	}

	for _, c := range out.Columns {
		resp.Columns = append(resp.Columns, &reportv1.Column{
			ColumnName: c.Name,
			ColumnComm: c.Comment,
		})
	}

	return resp, nil
}
