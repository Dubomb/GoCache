package cache

type TimeUnit int

type EvictionPolicy interface {
	OnGet(key string)

	OnSet(key string)

	Evict() (string, bool)

	Remove(key string)

	Len() int
}
