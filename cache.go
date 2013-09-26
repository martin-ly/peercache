package peercache

import (
	"sync"
	"time"
)

// Cache is a simple in-memory cache of key (string) to data ([]byte).
// Every key has a TTL specified on write, after which the data will be
// erased. There are no other constraints.
type Cache struct {
	sync.RWMutex
	data map[string][]byte
	ttls map[string]time.Time
}

// NewCache returns a new, empty Cache.
func NewCache(flushInterval time.Duration) *Cache {
	c := &Cache{
		data: map[string][]byte{},
		ttls: map[string]time.Time{},
	}
	go c.manage(flushInterval)
	return c
}

func (c *Cache) manage(flushInterval time.Duration) {
	for now := range time.Tick(flushInterval) {
		expired := []string{}

		func() {
			c.RLock()
			defer c.RUnlock()
			for key, expire := range c.ttls {
				if expire.Before(now) {
					expired = append(expired, key)
				}
			}
		}()

		func() {
			c.Lock()
			defer c.Unlock()
			for _, key := range expired {
				delete(c.data, key)
				delete(c.ttls, key)
			}
		}()
	}
}

// Write writes the given value to the given key, expiring after ttl.
func (c *Cache) Write(key string, val []byte, ttl time.Duration) {
	c.Lock()
	defer c.Unlock()
	c.data[key] = val
	c.ttls[key] = time.Now().Add(ttl)
}

// Read returns the data cached under the given key, if it exists.
func (c *Cache) Read(key string) ([]byte, bool) {
	c.RLock()
	defer c.RUnlock()
	val, ok := c.data[key]
	return val, ok
}
