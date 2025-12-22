package service

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/services/report/internal/cache"
	"github.com/ummuys/reportify/services/report/internal/dto"
)

type reportCacheService struct {
	cache  cache.ReportCache
	logger zerolog.Logger
}

func NewReportCacheService(cacheSvc cache.ReportCache, baseLogger zerolog.Logger) ReportCacheService {
	logger := baseLogger.With().Str("component", "report-cache-service").Logger()
	return &reportCacheService{cache: cacheSvc, logger: logger}
}

func (cs *reportCacheService) Get(ctx context.Context, key string) (dto.CacheGetResult, error) {
	cs.logger.Debug().Str("evt", "call Get").Str("key", key).Msg("")

	out, err := cs.cache.Get(ctx, key)
	if err != nil {
		return dto.CacheGetResult{}, err
	}

	return dto.CacheGetResult{Values: out}, nil
}

func (cs *reportCacheService) Delete(ctx context.Context, key string, value string) error {
	cs.logger.Debug().Str("evt", "call Delete").Str("key", key).Msg("")
	return cs.cache.Delete(ctx, key, []byte(value))
}

func (cs *reportCacheService) DeleteAll(ctx context.Context, key string) error {
	cs.logger.Debug().Str("evt", "call DeleteAll").Str("key", key).Msg("")
	return cs.cache.DeleteAll(ctx, key)
}
