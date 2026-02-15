package redis

import (
	"context"
	"encoding/json"
	"time"
	"univents/internal/plataform/memory"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var _ memory.CacheService = (*RedisCache)(nil)

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		client: rdb,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, bool) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
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

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) {
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

func (r *RedisCache) Delete(ctx context.Context, key string) {
	r.client.Del(ctx, key)
}

func (r *RedisCache) DeleteByPrefix(ctx context.Context, prefix string) {
	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()
	for iter.Next(ctx) {
		r.client.Del(ctx, iter.Val())
	}
}

func (r *RedisCache) Clear(ctx context.Context) {
	r.client.FlushDB(ctx)
}
