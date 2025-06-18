package cache

import (
	"container/list"
	"sync"
)

type lfuItem struct {
	key  string
	freq int
}

type lfuPolicy struct {
	cache   map[string]*list.Element
	order   map[int]*list.List
	minFreq int
	size    int
	mutex   sync.Mutex
}

func NewLFUPolicy() EvictionPolicy {
	return &lfuPolicy{
		cache:   make(map[string]*list.Element),
		order:   make(map[int]*list.List),
		minFreq: 0,
		size:    0,
	}
}

func (lfu *lfuPolicy) OnGet(key string) {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	elem := lfu.cache[key]
	item := elem.Value.(*lfuItem)

	oldFreq := item.freq
	item.freq++
	newFreq := item.freq

	lfu.order[oldFreq].Remove(elem)

	if _, exists := lfu.order[newFreq]; !exists {
		lfu.order[newFreq] = list.New()
	}

	elem = lfu.order[newFreq].PushFront(item)
	lfu.cache[key] = elem

	lfu.updateMinFreq(oldFreq)
}

func (lfu *lfuPolicy) OnSet(key string) {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	if elem, exists := lfu.cache[key]; exists {
		item := elem.Value.(*lfuItem)

		oldFreq := item.freq
		item.freq++
		newFreq := item.freq

		item.key = key

		lfu.order[oldFreq].Remove(elem)

		if _, exists := lfu.order[newFreq]; !exists {
			lfu.order[newFreq] = list.New()
		}

		elem := lfu.order[newFreq].PushFront(item)
		lfu.cache[key] = elem

		lfu.updateMinFreq(oldFreq)
	} else {
		item := &lfuItem{
			key:  key,
			freq: 1,
		}

		if _, exists := lfu.order[1]; !exists {
			lfu.order[1] = list.New()
		}

		elem := lfu.order[1].PushFront(item)
		lfu.cache[key] = elem

		lfu.minFreq = 1
		lfu.size++
	}
}

func (lfu *lfuPolicy) Evict() (string, bool) {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	leastFreqOrder := lfu.order[lfu.minFreq]

	back := leastFreqOrder.Back()

	if back != nil {
		leastFreqOrder.Remove(back)
		item := back.Value.(*lfuItem)
		delete(lfu.cache, item.key)

		lfu.size--

		lfu.updateMinFreq(item.freq)

		return item.key, true
	}

	return "", false
}

func (lfu *lfuPolicy) Remove(key string) {
	lfu.mutex.Lock()
	defer lfu.mutex.Unlock()

	if elem, exists := lfu.cache[key]; exists {
		freq := elem.Value.(*lfuItem).freq
		lfu.order[freq].Remove(elem)
		delete(lfu.cache, key)

		lfu.size--

		lfu.updateMinFreq(freq)
	}
}

func (lfu *lfuPolicy) Len() int {
	return lfu.size
}

func (lfu *lfuPolicy) updateMinFreq(prevFreq int) {
	if lfu.order[prevFreq].Len() == 0 {
		delete(lfu.order, prevFreq)

		if prevFreq == lfu.minFreq && lfu.size > 0 {
			for i := lfu.minFreq + 1; i <= lfu.size+1; i++ {
				if _, exists := lfu.order[i]; exists && lfu.order[i].Len() > 0 {
					lfu.minFreq = i
					break
				}
			}
		}
	} else if lfu.size == 0 {
		lfu.minFreq = 0
	}
}
