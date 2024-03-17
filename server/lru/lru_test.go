package lru

import (
	"testing"
	"time"
)

func TestLRUCache_PutGet(t *testing.T) {
	cache := NewLRU[int, string](3, time.Second)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Put(3, "three")

	val, ok := cache.Get(1)
	if !ok || val != "one" {
		t.Errorf("Expected value 'one' for key 1, got: %s", val)
	}

	cache.Put(4, "four")

	val, ok = cache.Get(2) // key 2 should be evicted
	if ok {
		t.Errorf("Expected key 2 to be evicted, but it was found in the cache")
	}

	// Checking if key 1 is still present after exceeding the size limit
	val, ok = cache.Get(1)
	if !ok || val != "one" {
		t.Errorf("Expected value 'one' for key 1, got: %s", val)
	}
}

func TestLRUCache_PurgeAll(t *testing.T) {
	cache := NewLRU[int, string](3, time.Second)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Put(3, "three")

	cache.PurgeAll()

	if _, ok := cache.Get(1); ok {
		t.Error("Expected cache to be empty after PurgeAll, but it still contains values")
	}

	if _, ok := cache.Get(2); ok {
		t.Error("Expected cache to be empty after PurgeAll, but it still contains values")
	}

	if _, ok := cache.Get(3); ok {
		t.Error("Expected cache to be empty after PurgeAll, but it still contains values")
	}
}

func TestLRUCache_TTLExpiration(t *testing.T) {
	cache := NewLRU[int, string](3, time.Second)

	cache.Put(1, "one")

	time.Sleep(2 * time.Second)

	if _, ok := cache.Get(1); ok {
		t.Error("Expected key 1 to be expired and evicted from the cache")
	}
}

func TestLRUCache_FullCapacity(t *testing.T) {
	cache := NewLRU[int, string](3, time.Second)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Put(3, "three")

	// Cache is now at full capacity, adding another entry should evict the least recently used item (key: 1)
	cache.Put(4, "four")

	// Accessing key 1 should return false, as it has been evicted from the cache
	if _, ok := cache.Get(1); ok {
		t.Errorf("Expected key 1 to be evicted, but it was found in the cache")
	}

	// Accessing key 2, 3, and 4 should return true, as they are still present in the cache
	for _, key := range []int{2, 3, 4} {
		if _, ok := cache.Get(key); !ok {
			t.Errorf("Expected key %d to be in the cache, but it was not found", key)
		}
	}
}

func TestLRUCache_TimerReset(t *testing.T) {
	cache := NewLRU[int, string](3, time.Second)

	cache.Put(1, "one")
	cache.Put(2, "two")
	cache.Put(3, "three")

	// Access key 1 to make it the most recently used
	cache.Get(1)

	// Add a new entry with the same key as before (key: 1)
	cache.Put(1, "newOne")

	// Wait for the timer to reset (in this case, we wait for 2 seconds to ensure the timer expires)
	time.Sleep(2 * time.Second)

	// Accessing key 1 should return false, as its timer should have expired and it should have been evicted
	if _, ok := cache.Get(1); ok {
		t.Errorf("Expected key 1 to be evicted after timer reset, but it was found in the cache")
	}
}
