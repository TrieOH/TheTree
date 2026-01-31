package outbounds

import (
	"context"
	"time"
)

type CacheService interface {
	Get(ctx context.Context, key string) (any, bool)
	Set(ctx context.Context, key string, value any, ttl time.Duration)
	Delete(ctx context.Context, key string)
	Clear(ctx context.Context)
}
