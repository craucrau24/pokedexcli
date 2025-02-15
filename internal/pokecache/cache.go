package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries map[string]cacheEntry
	mutex   sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	return Cache{entries: make(map[string]cacheEntry)}
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	entry, ok := c.entries[key]
	if !ok {
		return nil, ok
	}
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for (t := <- ticker.C) {
		c.mutex.Lock()
		defer c.mutex.Unlock()

		var toRemove []string{}
		for key, val := range c.entries {
			if val.createdAt < t - interval {
				toRemove = append(toRemove, key)
			}
		}
		for _, key := range(toRemove) {
			delete(c.entries, key)
		}
	}
}
