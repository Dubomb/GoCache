package cache

type TimeUnit int

type EvictionPolicy interface {
	OnSet(key, value string)

	OnSetWithTTL(key, value string, ttl uint, unit TimeUnit)

	Evict() string

	Remove(key string)

	Len() int
}
