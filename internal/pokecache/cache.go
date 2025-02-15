package pokecache

import (
	"fmt"
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

func NewCache(interval time.Duration) *Cache {
	c := Cache{entries: make(map[string]cacheEntry)}
	go c.reapLoop(interval)
	return &c
}

func (c *Cache) Add(key string, val []byte) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	fmt.Printf("Add -> %s: %d\n", key, len(c.entries))
}

func (c *Cache) Get(key string) ([]byte, bool) {
	fmt.Printf("Get -> %s: %d\n", key, len(c.entries))
	entry, ok := c.entries[key]
	if !ok {
		fmt.Println("cache missed")
		return nil, ok
	}
	fmt.Println("cache hit")
	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for t := range ticker.C {
		c.mutex.Lock()

		var toRemove []string
		for key, val := range c.entries {
			if val.createdAt.Add(interval).Before(t) {
				toRemove = append(toRemove, key)
			}
		}
		for _, key := range toRemove {
			delete(c.entries, key)
		}
		c.mutex.Unlock()
	}
}
