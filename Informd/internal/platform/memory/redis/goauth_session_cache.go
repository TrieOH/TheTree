package redis

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionCache struct {
	client *redis.Client
}

func NewSessionCache(rdb *redis.Client) *SessionCache {
	return &SessionCache{
		client: rdb,
	}
}

func (r *SessionCache) GetSession(ctx context.Context, id string) ([]byte, error) {
	val, err := r.client.Get(ctx, id).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return []byte(val), nil
}

func (r *SessionCache) SetSession(ctx context.Context, id string, data []byte, ttl time.Duration) error {
	return r.client.Set(ctx, id, data, ttl).Err()
}

func (r *SessionCache) DeleteSession(ctx context.Context, id string) error {
	return r.client.Del(ctx, id).Err()
}
