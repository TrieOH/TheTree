package imc

import (
	"context"
	"payssage/internal/platform/memory"
	"time"
)

type InMemoryCache struct {
	lru *LRU[string, any]
}

var _ memory.CacheService = (*InMemoryCache)(nil)

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

func (c *InMemoryCache) DeleteByPrefix(ctx context.Context, prefix string) {
	c.lru.mu.Lock()
	defer c.lru.mu.Unlock()

	for key, elem := range c.lru.cache {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			c.lru.list.Remove(elem)
			delete(c.lru.cache, key)
		}
	}
}

func (c *InMemoryCache) Clear(ctx context.Context) {
	c.lru.Clear()
}
