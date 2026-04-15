package redis

import (
	"IdentityX/internal/shared/ports"
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
}

var _ ports.RedisCacheService = (*RedisCache)(nil)

func NewRedisCache(rdb *redis.Client) *RedisCache {
	return &RedisCache{
		client: rdb,
	}
}

func (r *RedisCache) Get(ctx context.Context, key string) (any, bool, error) {
	val, err := r.client.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, err
	}

	var result any
	if err := json.Unmarshal([]byte(val), &result); err != nil {
		return val, true, nil
	}

	return result, true, nil
}

// New GetAny: returns raw value as []byte (no unmarshal) for compatibility with middleware
func (r *RedisCache) GetAny(ctx context.Context, key string) (any, bool, error) {
	val, err := r.client.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	return []byte(val), true, nil
}

func (r *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	var data []byte
	var err error

	if s, ok := value.(string); ok {
		data = []byte(s)
	} else {
		data, err = json.Marshal(value)
		if err != nil {
			return err
		}
	}

	cmd := r.client.Set(ctx, key, data, ttl)
	return cmd.Err()
}

func (r *RedisCache) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisCache) DeleteByPrefix(ctx context.Context, prefix string) error {

	iter := r.client.Scan(ctx, 0, prefix+"*", 0).Iterator()

	for iter.Next(ctx) {
		if err := r.client.Del(ctx, iter.Val()).Err(); err != nil {
			return err
		}
	}

	if err := iter.Err(); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Clear(ctx context.Context) error {
	return r.client.FlushDB(ctx).Err()
}
