# TTL Cache

A simple TTL cache implementation with generics in Go.

```go
cache, done := ttlcache.New[string, int]()
defer close(done)

cache.Set("key", 100, 10 * time.Second)

v, ok := cache.Get("test")
fmt.Println(v, ok) // 100, true
```