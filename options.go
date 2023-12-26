package ttlcache

import "time"

type Option[K comparable, V any] func(*Cache[K, V])

// WithNowProvider sets the now provider for the cache.
// The default now provider is time.Now.
func WithNowProvider[K comparable, V any](
	now nowProvider) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.nowProvider = now
	}
}

// WithCleanUpInterval sets the interval between clean up runs.
// The default interval is 5 second.
func WithCleanUpInterval[K comparable, V any](
	interval time.Duration) Option[K, V] {
	return func(c *Cache[K, V]) {
		c.cleanUpInterval = interval
	}
}
