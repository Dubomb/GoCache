package cache

type TimeUnit int

type EvictionPolicy interface {
	OnSet(key string)

	Evict() (string, bool)

	Remove(key string)

	Len() int
}
