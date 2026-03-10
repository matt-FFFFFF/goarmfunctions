package goarmfunctions

import (
	"container/list"
	"sync"
)

// DefaultParseCacheSize is the default maximum number of parsed expressions
// that can be held in the cache.
const DefaultParseCacheSize = 1000

// lruCache is a concurrency-safe LRU cache for parsed ARM expressions.
type lruCache struct {
	mu       sync.Mutex
	capacity int
	items    map[string]*list.Element
	order    *list.List // front = most recently used
	inflight map[string]*call
}

// call represents an in-flight or completed cache computation.
type call struct {
	wg  sync.WaitGroup
	val any
	err error
}

type lruEntry struct {
	key   string
	value any
}

func newLRUCache(capacity int) *lruCache {
	if capacity < 0 {
		capacity = 0
	}
	return &lruCache{
		capacity: capacity,
		items:    make(map[string]*list.Element, capacity),
		order:    list.New(),
		inflight: make(map[string]*call),
	}
}

// load retrieves a value from the cache, returning the value and true if found.
// Accessing an entry promotes it to the front (most recently used).
func (c *lruCache) load(key string) (any, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(*lruEntry).value, true
	}
	return nil, false
}

// store adds or updates a value in the cache.
// If the cache is at capacity, the least recently used entry is evicted.
func (c *lruCache) store(key string, value any) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		elem.Value.(*lruEntry).value = value
		return
	}

	if c.order.Len() >= c.capacity {
		c.evictOldest()
	}

	entry := &lruEntry{key: key, value: value}
	elem := c.order.PushFront(entry)
	c.items[key] = elem
}

// evictOldest removes the least recently used entry. Must be called with mu held.
func (c *lruCache) evictOldest() {
	oldest := c.order.Back()
	if oldest == nil {
		return
	}
	c.order.Remove(oldest)
	delete(c.items, oldest.Value.(*lruEntry).key)
}

// loadOrCompute retrieves a value from the cache, or computes it using the
// provided function if not present. Only one goroutine will compute the value
// for a given key at a time; concurrent callers for the same key will wait
// for the in-flight computation to complete (singleflight pattern).
// On success the result is stored in the cache; on error it is not cached.
func (c *lruCache) loadOrCompute(key string, compute func() (any, error)) (any, error) {
	c.mu.Lock()

	// Fast path: already cached.
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		val := elem.Value.(*lruEntry).value
		c.mu.Unlock()
		return val, nil
	}

	// Another goroutine is already computing this key — wait for it.
	if cl, ok := c.inflight[key]; ok {
		c.mu.Unlock()
		cl.wg.Wait()
		return cl.val, cl.err
	}

	// We are the first — register an inflight entry.
	cl := &call{}
	cl.wg.Add(1)
	c.inflight[key] = cl
	c.mu.Unlock()

	// Compute the value outside the lock.
	val, err := compute()
	cl.val = val
	cl.err = err
	cl.wg.Done()

	c.mu.Lock()
	delete(c.inflight, key)
	if err == nil {
		if _, exists := c.items[key]; !exists {
			if c.order.Len() >= c.capacity {
				c.evictOldest()
			}
			entry := &lruEntry{key: key, value: val}
			elem := c.order.PushFront(entry)
			c.items[key] = elem
		}
	}
	c.mu.Unlock()

	return val, err
}

// reset clears all entries from the cache.
func (c *lruCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element, c.capacity)
	c.order.Init()
	c.inflight = make(map[string]*call)
}

// len returns the number of entries in the cache.
func (c *lruCache) len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.order.Len()
}

// resize changes the capacity of the cache.
// If the new capacity is smaller than the current number of entries,
// the least recently used entries are evicted.
// Negative values are treated as 0.
func (c *lruCache) resize(capacity int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if capacity < 0 {
		capacity = 0
	}
	c.capacity = capacity
	for c.order.Len() > c.capacity {
		c.evictOldest()
	}
}
