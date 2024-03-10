package utils

import (
	"time"
)

type Cache struct {
	cache   map[string]cacheEntry
	entries int
	cfg     CacheConfig

	qchan  chan struct{}
	ticker *time.Ticker
}

type cacheEntry struct {
	value   any
	cfg     EntryConfig
	hits    int
	lastHit time.Time
	dirty   bool
}

type CacheConfig struct {
	DefaultTtl       time.Duration
	MaxEntries       int
	RenewTtlOnHit    bool
	ExpiresAfterHits int           // the number of times it can renew until expires forever. Defaults to 100 hits
	ExpiresAfterTtl  time.Duration // the amount of time it can renew until expires forever. Defaults to 5 hours.
	CheckEach        time.Duration // the amount of time between cache entry expiration checks
}

type EntryConfig struct {
	Ttl              time.Duration // the time to live of the entry. If zero, it won't expire
	RenewTtlOnHit    bool          // A flag that indicates if the ttl must be renew when the entry is hit
	ExpiresAfterHits int           // the number of times it can renew until expires forever. Defaults to 100 hits
	ExpiresAfterTtl  time.Duration // the amount of time it can renew until expires forever. Defaults to 5 hours.

	expiresAt      time.Time
	realExpiration time.Time
}

func NewCache(config ...CacheConfig) *Cache {
	var cfg CacheConfig
	if len(config) == 1 {
		cfg = config[0]
	} else {
		cfg = CacheConfig{
			DefaultTtl:       1 * time.Minute,
			MaxEntries:       50,
			RenewTtlOnHit:    false,
			CheckEach:        1 * time.Minute,
			ExpiresAfterHits: 100,
			ExpiresAfterTtl:  5 * time.Hour,
		}
	}

	cache := &Cache{
		cache: make(map[string]cacheEntry, cfg.MaxEntries),
		cfg:   cfg,
	}

	return cache
}

// Put adds a new entry to the cache
func (c *Cache) Put(key string, value any, entryConfig ...EntryConfig) {
	var entryCfg EntryConfig
	if len(entryConfig) == 1 {
		entryCfg = entryConfig[0]
	} else {
		entryCfg = EntryConfig{
			Ttl:              c.cfg.DefaultTtl,
			RenewTtlOnHit:    c.cfg.RenewTtlOnHit,
			ExpiresAfterHits: c.cfg.ExpiresAfterHits,
			ExpiresAfterTtl:  c.cfg.ExpiresAfterTtl,
		}
	}

	// delete older
	if c.cfg.MaxEntries != 0 && c.entries == c.cfg.MaxEntries {
		c.deleteOlder()
	}

	now := time.Now()
	entryCfg.expiresAt = now.Add(entryCfg.Ttl)

	if entryCfg.ExpiresAfterTtl != 0 {
		entryCfg.realExpiration = now.Add(entryCfg.ExpiresAfterTtl)
	}

	c.cache[key] = cacheEntry{
		value: value,
		cfg:   entryCfg,
	}
	c.entries++
}

// Delete removes a cached entry
func (c *Cache) Delete(key string) {
	c.entries--
	delete(c.cache, key)
}

// Get retrieves a value. Returns nil if not found
func (c *Cache) Get(key string) any {
	now := time.Now()
	if entry, ok := c.cache[key]; ok {
		if entry.dirty {
			return nil
		}

		if !entry.cfg.expiresAt.IsZero() && now.After(entry.cfg.expiresAt) {
			entry.dirty = true
			return nil
		}

		if entry.cfg.ExpiresAfterHits != 0 && entry.hits == entry.cfg.ExpiresAfterHits {
			entry.dirty = true
			return nil
		}

		if entry.cfg.RenewTtlOnHit {
			entry.cfg.expiresAt = now.Add(entry.cfg.Ttl)
		}

		entry.hits++
		entry.lastHit = now
		return entry.value
	}

	return nil
}

// Close disposes the cache and clears it
func (c *Cache) Close() {
	c.cache = nil
	c.entries = 0
	close(c.qchan)
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

func (c *Cache) reviewEntries() {
	deleteKeys := make([]string, 0, len(c.cache)/2)
	for k, v := range c.cache {
		if v.dirty {
			deleteKeys = append(deleteKeys, k)
		}
	}

	for _, k := range deleteKeys {
		delete(c.cache, k)
	}
	c.entries -= len(deleteKeys)
}

func (c *Cache) timer() {
	c.ticker = time.NewTicker(c.cfg.CheckEach)
	c.qchan = make(chan struct{})
	go func() {
		for {
			select {
			case <-c.ticker.C:
				c.reviewEntries()
			case <-c.qchan:
				c.ticker.Stop()
				return
			}
		}
	}()
}
