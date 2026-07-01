package adapter

import (
	"context"

	"github.com/rs/zerolog"
	dsv1 "github.com/ummuys/reportify/api/pb/datasource/service/v1"
	"github.com/ummuys/reportify/services/report/internal/dto"
	"github.com/ummuys/reportify/services/report/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type DatasourceAdapter struct {
	logger zerolog.Logger
	dsv1.UnimplementedDatasourceServiceServer
	datasourceSVC service.DatasourceService
}

func NewDatasourceAdapter(datasourceSVC service.DatasourceService, baseLogger zerolog.Logger) *DatasourceAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &DatasourceAdapter{datasourceSVC: datasourceSVC, logger: logger}
}

func (da *DatasourceAdapter) ListSchemas(ctx context.Context, in *dsv1.ListSchemasRequest) (*dsv1.ListSchemasResponse, error) {
	da.logger.Debug().Str("evt", "call ListSchemas").Msg("")

	out, err := da.datasourceSVC.ListSchemas(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &dsv1.ListSchemasResponse{
		Schemas: make([]*dsv1.Schema, 0, len(out.Schemas)),
	}

	for _, s := range out.Schemas {
		resp.Schemas = append(resp.Schemas, &dsv1.Schema{
			SchemaName: s.Name,
			SchemaComm: s.Comment,
		})
	}

	return resp, nil
}

func (da *DatasourceAdapter) ListTables(ctx context.Context, in *dsv1.ListTablesRequest) (*dsv1.ListTablesResponse, error) {
	da.logger.Debug().Str("evt", "call ListTables").Str("schema", in.SchemaName).Msg("")

	out, err := da.datasourceSVC.ListTables(ctx, dto.ListTablesParams{Schema: in.SchemaName})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &dsv1.ListTablesResponse{
		Tables: make([]*dsv1.Table, 0, len(out.Tables)),
	}

	for _, t := range out.Tables {
		resp.Tables = append(resp.Tables, &dsv1.Table{
			TableName: t.Name,
			TableComm: t.Comment,
		})
	}

	return resp, nil
}

func (da *DatasourceAdapter) ListColumns(ctx context.Context, in *dsv1.ListColumnsRequest) (*dsv1.ListColumnsResponse, error) {
	da.logger.Debug().Str("evt", "call ListColumns").Str("schema", in.SchemaName).Str("table", in.TableName).Msg("")

	out, err := da.datasourceSVC.ListColumns(ctx, dto.ListColumnsParams{
		Schema: in.SchemaName,
		Table:  in.TableName,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &dsv1.ListColumnsResponse{
		Columns: make([]*dsv1.Column, 0, len(out.Columns)),
	}

	for _, c := range out.Columns {
		resp.Columns = append(resp.Columns, &dsv1.Column{
			ColumnName: c.Name,
			ColumnComm: c.Comment,
		})
	}

	return resp, nil
}
