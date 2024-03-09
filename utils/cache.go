package utils

import (
	"time"
)

type Cache struct {
	cache   map[string]cacheEntry
	entries int
	cfg     CacheConfig
}

type cacheEntry struct {
	value   any
	cfg     EntryConfig
	hits    int
	lastHit time.Time
}

type CacheConfig struct {
	DefaultTtl    time.Duration
	MaxEntries    int
	RenewTtlOnHit bool
}

type EntryConfig struct {
	Ttl           time.Duration // the time to live of the entry. If zero, it won't expire
	RenewTtlOnHit bool          // A flag that indicates if the ttl must be renew when the entry is hit

	expiresAt time.Time
}

func NewCache(config ...CacheConfig) *Cache {
	var cfg CacheConfig
	if len(config) == 1 {
		cfg = config[0]
	} else {
		cfg = CacheConfig{
			DefaultTtl:    1 * time.Minute,
			MaxEntries:    50,
			RenewTtlOnHit: false,
		}
	}

	return &Cache{
		cache: make(map[string]cacheEntry, cfg.MaxEntries),
		cfg:   cfg,
	}
}

// Put adds a new entry to the cache
func (c *Cache) Put(key string, value any, entryConfig ...EntryConfig) {
	var entryCfg EntryConfig
	if len(entryConfig) == 1 {
		entryCfg = entryConfig[0]
	} else {
		entryCfg = EntryConfig{
			Ttl:           c.cfg.DefaultTtl,
			RenewTtlOnHit: c.cfg.RenewTtlOnHit,
		}
	}

	// delete older
	if c.cfg.MaxEntries != 0 && c.entries == c.cfg.MaxEntries {
		c.deleteOlder()
	}

	entryCfg.expiresAt = time.Now().Add(entryCfg.Ttl)
	c.cache[key] = cacheEntry{
		value: value,
		cfg:   entryCfg,
	}
	c.entries++
}

// Delete removes a cached entry
func (c *Cache) Delete(key string) {
	delete(c.cache, key)
}

// Get retrieves a value. Returns nil if not found
func (c *Cache) Get(key string) any {
	now := time.Now()
	if entry, ok := c.cache[key]; ok {
		if !entry.cfg.expiresAt.IsZero() && now.After(entry.cfg.expiresAt) {
			c.Delete(key)
			return nil
		}

		if entry.cfg.RenewTtlOnHit {
			entry.cfg.expiresAt = now.Add(entry.cfg.Ttl)
		}

		entry.hits++
		return entry.value
	}

	return nil
}

func (c *Cache) deleteOlder() {
	now := time.Now()
	var diff time.Duration = time.Duration(0)
	var toRemove string
	for k, v := range c.cache {
		sub := now.Sub(v.lastHit)
		if diff == 0 || sub > diff {
			toRemove = k
		}
	}

	c.Delete(toRemove)
}
