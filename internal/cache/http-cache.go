package cache

import (
	"sync"
	"time"
)

type CacheEntry struct {
	Value     []byte
	CreatedAt time.Time
}

type HttpCache struct {
	cache map[string]CacheEntry
	mu    sync.Mutex
	reapC <-chan time.Time
}

var once sync.Once
var instance *HttpCache

func GetHttpCacheInstance() *HttpCache {
	if instance == nil {
		once.Do(func() {
			instance = NewHttpCache()
		})
	}
	return instance
}

func NewHttpCache() *HttpCache {
	timer := time.NewTicker(5 * time.Minute)

	cache := &HttpCache{
		cache: make(map[string]CacheEntry),
		reapC: timer.C,
	}

	go cache.reapLoop()
	return cache
}

func (c *HttpCache) Get(url string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, exists := c.cache[url]
	return value.Value, exists
}

func (c *HttpCache) Add(url string, value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[url] = CacheEntry{
		Value:     value,
		CreatedAt: time.Now(),
	}
}

func (c *HttpCache) Delete(url string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, url)
}

func (c *HttpCache) reapLoop() {
	for range c.reapC {
		c.mu.Lock()
		for url, entry := range c.cache {
			if time.Since(entry.CreatedAt) > 1*time.Hour {
				delete(c.cache, url)
			}
		}
		c.mu.Unlock()
	}
}
