package ports

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

type RedisCacheService interface {
	Get(ctx context.Context, key string) (value any, found bool, err error)
	GetAny(ctx context.Context, key string) (any, bool, error)
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	DeleteByPrefix(ctx context.Context, prefix string) error
	Clear(ctx context.Context) error
}
