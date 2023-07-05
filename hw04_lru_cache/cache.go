package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type CacheItem struct {
	CacheItemKey Key
	CacheItemVal interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	listItem, ok := c.items[key]
	if ok {
		// update
		cacheItem, ok := listItem.Value.(*CacheItem)
		if !ok {
			// todo: panic ?
		}

		cacheItem.CacheItemVal = value
		c.queue.MoveToFront(listItem)

		return true
	} else {
		newCacheItem := &CacheItem{CacheItemKey: key, CacheItemVal: value}
		newListItem := c.queue.PushFront(newCacheItem)
		c.items[key] = newListItem

		if c.queue.Len() > c.capacity {
			lastListItem := c.queue.Back()
			cacheItem, ok := lastListItem.Value.(*CacheItem)
			if !ok {
				// todo: panic ?
			}

			c.queue.Remove(lastListItem)
			delete(c.items, cacheItem.CacheItemKey)
		}
	}

	return ok
}

// todo: make GetItem func

func (c *lruCache) Get(key Key) (interface{}, bool) {
	listItem, ok := c.items[key]
	if ok {
		c.queue.MoveToFront(listItem)

		cacheItem, ok := listItem.Value.(*CacheItem)
		if !ok {
			// todo: panic
		}

		return cacheItem.CacheItemVal, true
	}

	return nil, false
}

func (c *lruCache) Clear() {
	c.queue = nil
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
