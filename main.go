package main

import (
	"container/list"
	"time"
)

type TimeUnit int

const (
	Second TimeUnit = iota
	Minute
	Hour
)

type cacheItem struct {
	data      string
	expiresAt time.Time
}

type cacheEntry struct {
	key   string
	value *cacheItem
}

type GoCache struct {
	Capacity uint
	Cache    map[string]*list.Element
	order    *list.List
}

func NewGoCache(capacity uint) *GoCache {
	return &GoCache{
		Capacity: capacity,
		Cache:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (g *GoCache) newEntry(key, value string, expiresAt time.Time) {
	item := &cacheItem{
		data:      value,
		expiresAt: expiresAt,
	}

	entry := &cacheEntry{
		key:   key,
		value: item,
	}

	elem := g.order.PushFront(entry)
	g.Cache[key] = elem

	if g.order.Len() > int(g.Capacity) {
		back := g.order.Back()

		if back != nil {
			g.order.Remove(back)
			entry := back.Value.(*cacheEntry)
			delete(g.Cache, entry.key)
		}
	}
}

func (g *GoCache) Set(key, value string) {
	if elem, exists := g.Cache[key]; exists {
		newItem := elem.Value.(*cacheEntry)
		newItem.value.data = value
		g.order.MoveToFront(elem)
		return
	}

	g.newEntry(key, value, time.Time{})
}

func (g *GoCache) SetWithTTL(key, value string, ttl uint, unit TimeUnit) {
	if _, exists := g.Cache[key]; exists {
		g.Del(key)
	}

	var duration time.Duration

	switch unit {
	case Second:
		duration = time.Duration(ttl) * time.Second
	case Minute:
		duration = time.Duration(ttl) * time.Minute
	case Hour:
		duration = time.Duration(ttl) * time.Hour
	default:
		duration = time.Duration(ttl) * time.Second
	}

	expiresAt := time.Now().Add(duration)

	g.newEntry(key, value, expiresAt)
}

func (g *GoCache) Get(key string) string {
	elem, exists := g.Cache[key]

	if !exists {
		return ""
	}

	item := elem.Value.(*cacheEntry).value

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		g.order.Remove(elem)
		delete(g.Cache, key)
		return ""
	}

	g.order.MoveToFront(elem)

	return item.data
}

func (g *GoCache) Del(key string) {
	if elem, exists := g.Cache[key]; exists {
		g.order.Remove(elem)
		delete(g.Cache, key)
	}
}

func (g *GoCache) Exists(key string) bool {
	elem, exists := g.Cache[key]

	if !exists {
		return false
	}

	item := elem.Value.(*cacheEntry).value

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(g.Cache, key)
		return false
	}

	return true
}

func main() {

}
