package cache

type TimeUnit int

type EvictionPolicy interface {
	OnSet(key string)

	Evict() string

	Remove(key string)

	Len() int
}
