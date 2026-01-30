package memory

import (
	"GoAuth/internal/ports/outbounds"
	"context"
	"time"
)

type InMemoryCache struct {
	lru *LRU[string, any]
}

var _ outbounds.CacheService = (*InMemoryCache)(nil)

func NewInMemoryCache(capacity int, defaultTTL time.Duration) *InMemoryCache {
	return &InMemoryCache{
		lru: NewLRU[string, any](capacity, defaultTTL),
	}
}

func (c *InMemoryCache) Get(ctx context.Context, key string) (any, bool) {
	return c.lru.Get(key)
}

func (c *InMemoryCache) Set(ctx context.Context, key string, value any, ttl time.Duration) {
	c.lru.Put(key, value, ttl)
}

func (c *InMemoryCache) Delete(ctx context.Context, key string) {
	c.lru.Remove(key)
}

func (c *InMemoryCache) Clear(ctx context.Context) {
	c.lru.Clear()
}
