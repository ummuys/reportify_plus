package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"github.com/ummuys/reportify/pkg/config"
)

type repCache struct {
	logger zerolog.Logger
	cli    *redis.Client
	ttl    time.Duration
}

func NewReportCache(ctx context.Context, baseLogger zerolog.Logger) (ReportCache, error) {
	qctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cfg, err := config.ParseReportCacheEnv()
	if err != nil {
		return nil, err
	}

	cli := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := cli.Ping(qctx).Err(); err != nil {
		return nil, fmt.Errorf("redis didn't pinged: %v", err)
	}

	logger := baseLogger.With().Str("component", "report-cache").Logger()

	return &repCache{
		cli:    cli,
		logger: logger,
		ttl:    cfg.TTL,
	}, nil
}

func (rc *repCache) Set(ctx context.Context, key string, value string) error {
	rc.logger.Debug().Str("evt", "call Set").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := rc.cli.Set(qctx, key, value, rc.ttl).Err(); err != nil {
		return err
	}
	return nil
}
