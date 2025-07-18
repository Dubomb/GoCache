package cache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	value     string
	expiresAt time.Time
}

type GoCache struct {
	Capacity uint
	Cache    map[string]*cacheEntry
	policy   EvictionPolicy
	mutex    sync.Mutex
}

func NewGoCache(capacity uint, policy EvictionPolicy) *GoCache {
	return &GoCache{
		Capacity: capacity,
		Cache:    make(map[string]*cacheEntry),
		policy:   policy,
	}
}

func (cache *GoCache) newEntry(key, value string, expiresAt time.Time) {
	entry := &cacheEntry{
		value:     value,
		expiresAt: expiresAt,
	}

	cache.Cache[key] = entry
	cache.policy.OnSet(key)

	if cache.policy.Len() > int(cache.Capacity) {
		if evictedKey, found := cache.policy.Evict(); found {
			delete(cache.Cache, evictedKey)
		}
	}
}

func (cache *GoCache) Set(key, value string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.newEntry(key, value, time.Time{})
}

func (cache *GoCache) SetWithTTL(key, value string, ttl uint) {
	duration := time.Duration(ttl) * time.Millisecond

	expiresAt := time.Now().Add(duration)

	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	cache.newEntry(key, value, expiresAt)
}

func (cache *GoCache) Get(key string) (string, bool) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	elem, exists := cache.Cache[key]

	if !exists {
		return "", false
	}

	if !elem.expiresAt.IsZero() && time.Now().After(elem.expiresAt) {
		cache.policy.Remove(key)
		delete(cache.Cache, key)
		return "", false
	}

	cache.policy.OnGet(key)

	return cache.Cache[key].value, true
}

func (cache *GoCache) Del(key string) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	if _, exists := cache.Cache[key]; exists {
		cache.policy.Remove(key)
		delete(cache.Cache, key)
	}
}

func (cache *GoCache) Exists(key string) bool {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()

	elem, exists := cache.Cache[key]

	if !exists {
		return false
	}

	if !elem.expiresAt.IsZero() && time.Now().After(elem.expiresAt) {
		cache.policy.Remove(key)
		delete(cache.Cache, key)
		return false
	}

	return true
}
