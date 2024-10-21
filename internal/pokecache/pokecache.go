package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type cache struct {
	mu      sync.RWMutex
	entries map[string]cacheEntry
}

func NewCache(interval time.Duration) *cache {
	cache := cache{}
	cache.reapLoop(interval)
	return &cache
}

func (c *cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
	c.entries[key] = entry
}

func (c *cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.Unlock()
	data, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return data.val, true
}

func (c *cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	go func() {
		for {
			select {
			case t := <-ticker.C:
				c.mu.Lock()
				for entry := range c.entries {
					if t.Sub(c.entries[entry].createdAt) > interval {
						delete(c.entries, entry)
					}
				}
				c.mu.Unlock()
			default:
			}
		}
	}()
}
