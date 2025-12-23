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

func (rc *repCache) Init(ctx context.Context, queries map[string][][]byte) error {
	rc.logger.Debug().Str("evt", "call Init").Msg("")
	qctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	for key, values := range queries {
		for _, val := range values {
			if err := rc.cli.LPush(qctx, key, val).Err(); err != nil {
				return fmt.Errorf("can't add a query (Init): %v", err)
			}
		}
	}

	return nil
}

func (rc *repCache) Set(ctx context.Context, key string, value []byte) error {
	rc.logger.Debug().Str("evt", "call Set").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if err := rc.cli.LPush(qctx, key, value).Err(); err != nil {
		return fmt.Errorf("can't add a query: %v", err)
	}
	return nil
}

func (rc *repCache) Get(ctx context.Context, key string) ([][]byte, error) {
	rc.logger.Debug().Str("evt", "call Get").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	value, err := rc.cli.LRange(qctx, key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("can't get a value: %v", err)
	}

	out := make([][]byte, len(value))
	for i, v := range value {
		out[i] = []byte(v)
	}

	return out, nil
}

func (rc *repCache) GetAll(ctx context.Context) (map[string][]byte, error) {
	rc.logger.Debug().Str("evt", "call GetAll").Msg("")
	qctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	keys, err := rc.cli.Keys(qctx, "*").Result()
	if err != nil {
		return nil, err
	}

	m := make(map[string][]byte, len(keys))

	for _, k := range keys {
		val, err := rc.cli.LRange(qctx, k, 0, -1).Result()
		if err != nil {
			return nil, fmt.Errorf("can't get a value: %v", err)
		}

		var queries []byte
		queries = append(queries, '[')
		for i, q := range val {
			queries = append(queries, []byte(q)...)
			if i+1 < len(val) {
				queries = append(queries, ',')
			}
		}
		queries = append(queries, ']')

		m[k] = queries
	}

	return m, nil
}

func (rc *repCache) Delete(ctx context.Context, key string, value []byte) error {
	rc.logger.Debug().Str("evt", "call Delete").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if _, err := rc.cli.LRem(qctx, key, 0, value).Result(); err != nil {
		return fmt.Errorf("can't delete a value: %v", err)
	}

	return nil
}

func (rc *repCache) DeleteAll(ctx context.Context, key string) error {
	rc.logger.Debug().Str("evt", "call DeleteAll").Msg("")
	qctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	if _, err := rc.cli.Del(qctx, key).Result(); err != nil {
		return fmt.Errorf("can't delete all values: %v", err)
	}

	return nil
}
