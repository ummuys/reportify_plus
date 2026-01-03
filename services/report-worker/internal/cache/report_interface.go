package cache

import "context"

type ReportCache interface {
	Set(ctx context.Context, key string, value string) error
}
