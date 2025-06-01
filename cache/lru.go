package cache

import (
	"container/list"
	"sync"
)

type lruItem struct {
	key string
}

type lruPolicy struct {
	cache map[string]*list.Element
	order *list.List
	mutex sync.Mutex
}

func NewLRUPolicy() EvictionPolicy {
	return &lruPolicy{
		cache: make(map[string]*list.Element),
		order: list.New(),
	}
}

func (lru *lruPolicy) OnGet(key string) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	elem := lru.cache[key]
	lru.order.MoveToFront(elem)
}

func (lru *lruPolicy) OnSet(key string) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if elem, exists := lru.cache[key]; exists {
		lru.order.MoveToFront(elem)
	} else {
		newItem := &lruItem{
			key: key,
		}

		elem := lru.order.PushFront(newItem)

		lru.cache[key] = elem
	}
}

func (lru *lruPolicy) Evict() (string, bool) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	back := lru.order.Back()

	if back != nil {
		lru.order.Remove(back)
		entry := back.Value.(*lruItem)
		delete(lru.cache, entry.key)
		return entry.key, true
	}

	return "", false
}

func (lru *lruPolicy) Remove(key string) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	if elem, exists := lru.cache[key]; exists {
		lru.order.Remove(elem)
		delete(lru.cache, key)
	}
}

func (lru *lruPolicy) Len() int {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	return lru.order.Len()
}
