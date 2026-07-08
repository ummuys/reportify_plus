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
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if err := rc.cli.Set(qctx, key, value, rc.ttl).Err(); err != nil {
		return err
	}
	rc.logger.Debug().Str("key", key).Dur("ttl", rc.ttl).Str("value", value).Msg("redis set")

	return nil
}

func (rc *repCache) Get(ctx context.Context, key string) (*string, error) {
	rc.logger.Debug().Str("evt", "call Get").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	value, err := rc.cli.Get(qctx, key).Result()
	if err != nil {
		rc.logger.Debug().Err(err).Str("key", key).Msg("redis get")
		return nil, err
	}
	rc.logger.Debug().Str("key", key).Str("value", value).Msg("redis hit")

	return &value, nil
}

func (rc *repCache) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	rc.logger.Debug().
		Str("evt", "call Delete").
		Int("keys_count", len(keys)).
		Msg("")

	qctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	_, err := rc.cli.Unlink(qctx, keys...).Result()
	if err != nil {
		rc.logger.Debug().
			Err(err).
			Interface("keys", keys).
			Msg("redis delete failed")
		return err
	}

	return nil
}
