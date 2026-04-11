package redis

import (
	"TrieForms/internal/platform/memory"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cache struct {
	client *redis.Client
}

var _ memory.CacheService = (*Cache)(nil)

func NewRedisCache(rdb *redis.Client) *Cache {
	return &Cache{
		client: rdb,
	}
}

func (r *Cache) Get(ctx context.Context, key string) (any, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false
	} else if err != nil {
		return nil, false
	}

	var result any
	err = json.Unmarshal([]byte(val), &result)
	if err != nil {
		return val, true // Return as string if not JSON
	}

	return result, true
}

func (r *Cache) Set(ctx context.Context, key string, value any, ttl time.Duration) {
	var data []byte
	var err error

	if s, ok := value.(string); ok {
		data = []byte(s)
	} else {
		data, err = json.Marshal(value)
		if err != nil {
			return
		}
	}

	r.client.Set(ctx, key, data, ttl)
}

func (r *Cache) Delete(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}

func (r *Cache) DeleteByPrefix(ctx context.Context, prefix string) {
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		r.client.Del(ctx, iter.Val())
	}
}

func (r *Cache) Clear(ctx context.Context) {
	r.client.FlushDB(ctx)
}
