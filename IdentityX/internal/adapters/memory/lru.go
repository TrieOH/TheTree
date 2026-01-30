package memory

import (
	"container/list"
	"sync"
	"time"
)

type LRU[K comparable, V any] struct {
	capacity int
	ttl      time.Duration
	cache    map[K]*list.Element
	list     *list.List
	mu       sync.RWMutex
}

type entry[K comparable, V any] struct {
	key       K
	value     V
	expiresAt time.Time
}

// NewLRU creates a new LRU cache.
// If ttl > 0, items will expire after the duration.
// Accessing an item (Get or Put) refreshes its expiration time.
func NewLRU[K comparable, V any](capacity int, ttl time.Duration) *LRU[K, V] {
	return &LRU[K, V]{
		capacity: capacity,
		ttl:      ttl,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
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
	if c.ttl > 0 && time.Now().After(ent.expiresAt) {
		c.list.Remove(elem)
		delete(c.cache, key)
		var zero V
		return zero, false
	}

	// Move to front and refresh TTL
	c.list.MoveToFront(elem)
	if c.ttl > 0 {
		ent.expiresAt = time.Now().Add(c.ttl)
	}

	return ent.value, true
}

func (c *LRU[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	var expiresAt time.Time
	if c.ttl > 0 {
		expiresAt = time.Now().Add(c.ttl)
	}

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		ent := elem.Value.(*entry[K, V])
		ent.value = value
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