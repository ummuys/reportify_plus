package service

import (
	"context"

	"github.com/ummuys/reportify/services/report/internal/dto"
)

type ReportCacheService interface {
	Get(ctx context.Context, key string) (dto.CacheGetResult, error)
	Delete(ctx context.Context, key string, value string) error
	DeleteAll(ctx context.Context, key string) error
}
