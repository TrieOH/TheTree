package memory

import (
	"context"
	"time"
)

type CacheService interface {
	Get(ctx context.Context, key string) (any, bool)
	Set(ctx context.Context, key string, value any, ttl time.Duration)
	Delete(ctx context.Context, key string)
	DeleteByPrefix(ctx context.Context, prefix string)
	Clear(ctx context.Context)
}
