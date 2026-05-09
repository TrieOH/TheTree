package imc

import (
	"container/list"
	"sync"
	"time"
)

type LRU[K comparable, V any] struct {
	capacity   int
	defaultTTL time.Duration
	cache      map[K]*list.Element
	list       *list.List
	mu         sync.RWMutex
}

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
	duration  time.Duration
}

// NewLRU creates a new LRU cache.
// defaultTTL is used if Put is called with 0.
func NewLRU[K comparable, V any](capacity int, defaultTTL time.Duration) *LRU[K, V] {
	return &LRU[K, V]{
		capacity:   capacity,
		defaultTTL: defaultTTL,
		cache:      make(map[K]*list.Element),
		list:       list.New(),
	}
}

func (c *LRU[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.cache[key]
	if !ok {
		var zero V
		return zero, false
	}

	ent := elem.Value.(*entry[K, V])

	// Check expiration
	if ent.duration > 0 && time.Now().After(ent.expiresAt) {
		c.list.Remove(elem)
		delete(c.cache, key)
		var zero V
		return zero, false
	}

	// Move to front and refresh TTL
	c.list.MoveToFront(elem)
	if ent.duration > 0 {
		ent.expiresAt = time.Now().Add(ent.duration)
	}

	return ent.value, true
}

func (c *LRU[K, V]) Put(key K, value V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	duration := ttl
	if duration == 0 {
		duration = c.defaultTTL
	}

	var expiresAt time.Time
	if duration > 0 {
		expiresAt = time.Now().Add(duration)
	}

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		ent := elem.Value.(*entry[K, V])
		ent.value = value
		ent.duration = duration
		ent.expiresAt = expiresAt
		return
	}

	if c.list.Len() >= c.capacity {
		back := c.list.Back()
		if back != nil {
			c.list.Remove(back)
			kv := back.Value.(*entry[K, V])
			delete(c.cache, kv.key)
		}
	}

	elem := c.list.PushFront(&entry[K, V]{
		key:       key,
		value:     value,
		expiresAt: expiresAt,
		duration:  duration,
	})
	c.cache[key] = elem
}

func (c *LRU[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.Remove(elem)
		delete(c.cache, key)
	}
}

func (c *LRU[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.list.Init()
	c.cache = make(map[K]*list.Element)
}
