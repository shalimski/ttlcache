package ttlcache

import (
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	t.Run("it returns a new cache", func(t *testing.T) {
		c, done := New[int, int]()
		defer close(done)

		if c == nil {
			t.Error("expected a cache, got nil")
		}
	})

	t.Run("it returns a done channel", func(t *testing.T) {
		_, done := New[int, int]()

		if done == nil {
			t.Error("expected a done channel, got nil")
		}

		close(done)

	})

	t.Run("it sets the now provider", func(t *testing.T) {
		now := func() time.Time {
			return time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		}

		c, done := New[int, int](WithNowProvider[int, int](now))
		defer close(done)

		if c.now() != now() {
			t.Errorf("expected now provider to be %v, got %v", now(), c.now())
		}
	})

	t.Run("it sets the clean up interval", func(t *testing.T) {
		c, done := New[int, int](WithCleanUpInterval[int, int](time.Second))
		defer close(done)

		if c.cleanUpInterval != time.Second {
			t.Errorf("expected clean up interval to be %v, got %v", time.Second, c.cleanUpInterval)
		}
	})

}

func TestCache_Set(t *testing.T) {
	t.Run("it sets a value in the cache", func(t *testing.T) {
		c, done := New[int, int]()
		defer close(done)

		c.Set(1, 1, time.Second)

		if c.Len() != 1 {
			t.Errorf("expected cache to have 1 item, got %v", c.Len())
		}
	})

	t.Run("it sets the value with the correct expiration time two times", func(t *testing.T) {
		c, done := New[int, int]()
		defer close(done)

		c.Set(1, 1, time.Second)
		c.Set(1, 2, time.Second)

		if c.Len() != 1 {
			t.Errorf("expected cache to have 1 item, got %v", c.Len())
		}
	})
}

func TestCache_Get(t *testing.T) {
	t.Run("it returns the value from the cache", func(t *testing.T) {
		c, done := New[string, int]()
		defer close(done)

		c.Set("test", 1, time.Second)

		v, ok := c.Get("test")

		if !ok {
			t.Error("expected ok to be true, got false")
		}

		if v != 1 {
			t.Errorf("expected value to be 1, got %v", v)
		}
	})

	t.Run("it returns false if the key is not in the cache", func(t *testing.T) {
		c, done := New[int, int]()
		defer close(done)

		_, ok := c.Get(1)

		if ok {
			t.Error("expected ok to be false, got true")
		}
	})

	t.Run("it returns false if the key is expired", func(t *testing.T) {
		c, done := New[int, int]()
		defer close(done)

		c.Set(1000, 1, -2*time.Hour)

		v, ok := c.Get(1000)

		if ok {
			t.Error("expected ok to be false, got true")
		}

		if v != 1 {
			t.Errorf("expected value to be 1, got %v", v)
		}
	})
}
