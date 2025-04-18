package main

import "fmt"

type GoCache struct {
	cache map[string]string
}

func NewGoCache() *GoCache {
	return &GoCache{
		make(map[string]string),
	}
}

func (g *GoCache) Set(key string, value string) {
	g.cache[key] = value
}

func (g *GoCache) Get(key string) string {
	return g.cache[key]
}

func (g *GoCache) Del(key string) {
	delete(g.cache, key)
}

func (g *GoCache) Exists(key string) bool {
	_, exists := g.cache[key]
	return exists
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
}
