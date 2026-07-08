package cache

import "context"

type ReportCache interface {
	Set(ctx context.Context, key string, value string) error
	Get(ctx context.Context, key string) (*string, error)
	Delete(ctx context.Context, keys ...string) error
}
