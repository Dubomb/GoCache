package cache

import (
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestSetGet(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.Set("key", "value")
	value, _ := cache.Get("key")

	if value != "value" {
		t.Errorf("expected value to be \"value\", got '%s'", value)
	}
}

func TestSet_KeyExists(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.Set("key", "value1")
	cache.Set("key", "value2")
	value, _ := cache.Get("key")

	if value != "value2" {
		t.Errorf("expected value to be \"value\", got '%s'", value)
	}
}

func TestSet_Concurrent(t *testing.T) {
	cache := NewGoCache(500, NewLRUPolicy())

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)

		go func(i int) {
			cache.Set(strconv.Itoa(i), "value")
			wg.Done()
		}(i)
	}

	wg.Wait()

	if cache.policy.Len() != 500 {
		t.Errorf("not concurrent")
	}
}

func TestGet_NonExistentKey(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	value, _ := cache.Get("missing-key")

	if value != "" {
		t.Errorf("expected empty string for non-existent key, got '%s'", value)
	}
}

func TestSetWithTTLGet_KeyExpired(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.SetWithTTL("key", "value", 1, Second)

	time.Sleep(2 * time.Second)

	value, _ := cache.Get("key")

	if value != "" {
		t.Errorf("expected empty string for expired key, got '%s'", value)
	}
}

func TestSetWithTTL(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.SetWithTTL("key", "value1", 1, Second)

	value, _ := cache.Get("key")

	if value != "value1" {
		t.Errorf("expected value to be \"value\", got '%s'", value)
	}
}

func TestSetWithTTL_KeyExists(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.SetWithTTL("key", "value1", 1, Second)
	cache.SetWithTTL("key", "value2", 1, Second)
	value, _ := cache.Get("key")

	if value != "value2" {
		t.Errorf("expected value to be \"value\", got '%s'", value)
	}
}

func TestSetWithTTL_InvalidTime(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.SetWithTTL("key", "value1", 1, 5)

	time.Sleep(2 * time.Second)

	value, _ := cache.Get("key")

	if value != "" {
		t.Errorf("expected empty string for expired key, got '%s'", value)
	}
}

func TestDel(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.Set("key", "value")

	cache.Del("key")

	value, _ := cache.Get("key")

	if value != "" {
		t.Errorf("expected emptry string after delete, got '%s'", value)
	}
}

func TestExists(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.Set("key", "value")

	if exists := cache.Exists("key"); !exists {
		t.Errorf("expected Exists to return true, got '%t'", exists)
	}
}

func TestExists_KeyExpired(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.SetWithTTL("key", "value1", 0, Second)

	time.Sleep(1 * time.Second)

	if exists := cache.Exists("key"); exists {
		t.Errorf("expected Exists to return false, got '%t'", exists)
	}
}

func TestNotExists(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.Set("key", "value")
	cache.Del("key")

	if exists := cache.Exists("key"); exists {
		t.Errorf("expected Exists to return false, got '%t'", exists)
	}
}

func TestLRU(t *testing.T) {
	cache := NewGoCache(1, NewLRUPolicy())

	cache.Set("key1", "value1")
	cache.Set("key2", "value2")

	cache.Get("key1")

	if exists := cache.Exists("key2"); !exists {
		t.Errorf("expected key2 to be evicted, but Exists returned true")
	}

	if exists := cache.Exists("key1"); exists {
		t.Errorf("expected key1 to exist, but Exists returned false")
	}
}
