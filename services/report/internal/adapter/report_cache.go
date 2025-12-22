package adapter

import (
	"context"

	"github.com/rs/zerolog"
	rcv1 "github.com/ummuys/reportify/api/pb/report/cache/v1"
	"github.com/ummuys/reportify/services/report/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ReportCacheAdapter struct {
	logger zerolog.Logger
	rcv1.UnimplementedReportCacheServiceServer
	reportCacheSvc service.ReportCacheService
}

func NewReportCacheAdapter(reportCacheSvc service.ReportCacheService, baseLogger zerolog.Logger) *ReportCacheAdapter {
	logger := baseLogger.With().Str("component", "adpt").Logger()
	return &ReportCacheAdapter{reportCacheSvc: reportCacheSvc, logger: logger}
}

func (ra *ReportCacheAdapter) Get(ctx context.Context, in *rcv1.GetRequest) (*rcv1.GetResponse, error) {
	ra.logger.Debug().Str("evt", "call Get").Str("key", in.Key).Msg("")

	out, err := ra.reportCacheSvc.Get(ctx, in.Key)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &rcv1.GetResponse{Values: out.Values}, nil
}

func (ra *ReportCacheAdapter) Delete(ctx context.Context, in *rcv1.DeleteRequest) (*rcv1.DeleteResponse, error) {
	ra.logger.Debug().Str("evt", "call Delete").Str("key", in.Key).Msg("")

	if err := ra.reportCacheSvc.Delete(ctx, in.Key, in.Value); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &rcv1.DeleteResponse{}, nil
}

func (ra *ReportCacheAdapter) DeleteAll(ctx context.Context, in *rcv1.DeleteAllRequest) (*rcv1.DeleteAllResponse, error) {
	ra.logger.Debug().Str("evt", "call DeleteAll").Str("key", in.Key).Msg("")

	if err := ra.reportCacheSvc.DeleteAll(ctx, in.Key); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &rcv1.DeleteAllResponse{}, nil
}
