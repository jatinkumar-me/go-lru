package lru

import (
	"container/list"
	"sync"
	"time"
)

// There are two approaches of evicting the cache value after
// expiration time, either by periodic cleanups, or by assigning a
// timer to each cache entry, there are pros and cos for both
// of these approach, I've chosen the later one.

type ListVal[K comparable, V any] struct {
	key   K
	value V
	kill  *time.Timer
}

type LRUCache[K comparable, V any] struct {
	size  int
	ll    *list.List
	cache map[K]*list.Element
	ttl   time.Duration
	mu    sync.Mutex
}

// NO_EVICTION_TTL - very long ttl to prevent eviction
const NO_EVICTION_TTL = time.Hour * 24 * 365 * 10

func NewLRU[K comparable, V any](
	size int,
	ttl time.Duration,
) *LRUCache[K, V] {
	if size < 0 {
		size = 0
	}

	if ttl <= 0 {
		ttl = NO_EVICTION_TTL
	}

	return &LRUCache[K, V]{
		size:  size,
		ll:    list.New(),
		cache: make(map[K]*list.Element),
		ttl:   ttl,
		mu:    sync.Mutex{},
	}
}

// Purge clears the cache completely.
func (c *LRUCache[K, V]) PurgeAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k := range c.cache {
		delete(c.cache, k)
	}
}

// Retrieve value from cache
func (c *LRUCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	if !ok {
		return
	}
	c.ll.MoveToFront(entry)
	return entry.Value.(ListVal[K, V]).value, ok
}

func (c *LRUCache[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, ok := c.cache[key]
	// handle collision
	if ok {
        listVal := entry.Value.(ListVal[K, V])
        listVal.value = value
        listVal.kill.Reset(c.ttl)
		c.ll.MoveToFront(entry)
		return
	}

	// cache overflow
	if len(c.cache) == c.size {
		last := c.ll.Back()
		last.Value.(ListVal[K, V]).kill.Stop()
		delete(c.cache, last.Value.(ListVal[K, V]).key)
	}

	c.addEntry(key, value)
}

// Make sure to lock the mutex before calling this function
func (c *LRUCache[K, V]) addEntry(key K, value V) {
	killTimer := time.NewTimer(c.ttl)
	c.cache[key] = c.ll.PushFront(ListVal[K, V]{
		key:   key,
		value: value,
		kill:  killTimer,
	})

	go func() {
		<-killTimer.C
		c.mu.Lock()
		defer c.mu.Unlock()

		if entry, ok := c.cache[key]; ok {
			c.ll.Remove(entry)
			delete(c.cache, key)
		}
	}()
}
