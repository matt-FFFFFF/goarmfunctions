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

// reset clears all entries from the cache.
func (c *lruCache) reset() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = make(map[string]*list.Element, c.capacity)
	c.order.Init()
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
