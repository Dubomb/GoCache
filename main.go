package main

import (
	"fmt"
	"time"
)

type TimeUnit int

const (
	Second TimeUnit = iota
	Minute
	Hour
)

type CacheItem struct {
	data      string
	expiresAt time.Time
}

type GoCache struct {
	Cache map[string]*CacheItem
}

func NewGoCache() *GoCache {
	return &GoCache{
		make(map[string]*CacheItem),
	}
}

func (g *GoCache) Set(key, value string) {
	if _, exists := g.Cache[key]; !exists {
		g.Cache[key] = &CacheItem{}
	}

	g.Cache[key].data = value
}

func (g *GoCache) SetWithTTL(key, value string, ttl uint, unit TimeUnit) {
	if _, exists := g.Cache[key]; !exists {
		g.Cache[key] = &CacheItem{}
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

	g.Cache[key].data = value
	g.Cache[key].expiresAt = time.Now().Add(duration)
}

func (g *GoCache) Get(key string) string {
	item, exists := g.Cache[key]

	if !exists {
		return ""
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(g.Cache, key)
		return ""
	}

	return g.Cache[key].data
}

func (g *GoCache) Del(key string) {
	delete(g.Cache, key)
}

func (g *GoCache) Exists(key string) bool {
	item, exists := g.Cache[key]

	if !exists {
		return false
	}

	if !item.expiresAt.IsZero() && time.Now().After(item.expiresAt) {
		delete(g.Cache, key)
		return false
	}

	return true
}

func main() {
	fmt.Println("Hello world!")

	g := NewGoCache()
	g.Set("hello", "world")
	fmt.Println(g)
	fmt.Println(g.Get("hello"))
	fmt.Println(g.Exists("hello"))
	g.Del("hello")
	fmt.Println(g.Exists("hello"))
	fmt.Println(g)
	g.SetWithTTL("world", "my", 1, Second)
	fmt.Println(g)

	fmt.Println(g.Get("world"))
	time.Sleep(2 * time.Second)
	fmt.Println(g.Get("world"))
}
