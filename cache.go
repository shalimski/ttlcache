package ttlcache

import (
	"sync"
	"time"
)

type nowProvider func() time.Time

// Cache is a thread-safe in-memory cache with expiration.
type Cache[K comparable, V any] struct {
	data            map[K]value[V]
	mu              sync.RWMutex
	nowProvider     nowProvider
	cleanUpInterval time.Duration
	done            chan struct{}
}

// New returns a new cache.
func New[K comparable, V any](
	options ...Option[K, V]) (*Cache[K, V], chan struct{}) {

	done := make(chan struct{})

	c := Cache[K, V]{
		data:            make(map[K]value[V]),
		nowProvider:     time.Now,
		cleanUpInterval: 5 * time.Second,
		done:            done,
	}

	for _, option := range options {
		option(&c)
	}

	go c.cleanUp()

	return &c, done
}

type value[V any] struct {
	data      V
	expiredAt time.Time
}

func (v value[V]) isExpired(now time.Time) bool {
	return now.After(v.expiredAt)
}

// Set sets a value in the cache with duration.
func (c *Cache[K, V]) Set(key K, v V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value[V]{
		data:      v,
		expiredAt: c.now().Add(ttl),
	}
}

// Get gets a value from the cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	v, ok := c.data[key]
	if !ok {
		return v.data, false
	}

	if v.isExpired(c.now()) {
		return v.data, false
	}

	return v.data, true
}

// Delete deletes a value from the cache.
func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Len returns the number of items in the cache.
func (c *Cache[K, V]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.data)
}

func (c *Cache[K, V]) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data = make(map[K]value[V])
}

func (c *Cache[K, V]) cleanUp() {

	timer := time.NewTimer(c.cleanUpInterval)
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.mu.Lock()
			for k, v := range c.data {
				if v.isExpired(c.now()) {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		case <-c.done:
			return
		}
	}

}

// Keys returns all the keys in the cache.
func (c *Cache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()

	keys := make([]K, 0, len(c.data))
	for k := range c.data {
		keys = append(keys, k)
	}

	return keys
}

func (c *Cache[K, V]) now() time.Time {
	return c.nowProvider()
}
