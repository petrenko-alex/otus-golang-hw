package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	mu sync.Mutex

	capacity int
	queue    List
	items    map[Key]*ListItem
}

type CacheItem struct {
	CacheItemKey Key
	CacheItemVal interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity <= 0 {
		return false
	}

	listItem, ok := c.hitListItem(key)
	if ok {
		c.updateListItem(listItem, value)
	} else {
		c.addListItem(key, value)

		if c.needToPurgeCache() {
			c.purgeCache()
		}
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.capacity <= 0 {
		return nil, false
	}

	listItem, ok := c.hitListItem(key)
	if ok {
		cacheItem := c.getCacheItem(listItem)

		return cacheItem.CacheItemVal, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.mu.Lock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)

	c.mu.Unlock()
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (c *lruCache) getCacheItem(listItem *ListItem) *CacheItem {
	return listItem.Value.(*CacheItem)
}

func (c *lruCache) hitListItem(key Key) (*ListItem, bool) {
	listItem, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(listItem)
	}

	return listItem, ok
}

func (c *lruCache) updateListItem(listItem *ListItem, value interface{}) {
	cacheItem := c.getCacheItem(listItem)
	cacheItem.CacheItemVal = value
}

func (c *lruCache) addListItem(key Key, value interface{}) {
	newCacheItem := &CacheItem{CacheItemKey: key, CacheItemVal: value}
	newListItem := c.queue.PushFront(newCacheItem)

	c.items[key] = newListItem
}

func (c *lruCache) needToPurgeCache() bool {
	return c.queue.Len() > c.capacity
}

func (c *lruCache) purgeCache() {
	lastListItem := c.queue.Back()
	cacheItem := c.getCacheItem(lastListItem)

	c.queue.Remove(lastListItem)
	delete(c.items, cacheItem.CacheItemKey)
}
