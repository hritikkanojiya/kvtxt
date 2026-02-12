// Package cache provides in-memory caching layer.
// It improves read performance and reduces database load.

package cache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key         string
	value       string
	contentType string
	expiresAt   *int64
}

type Cache struct {
	mu      sync.Mutex
	maxSize int
	ll      *list.List
	items   map[string]*list.Element
}

func New(maxSize int) *Cache {
	if maxSize <= 0 {
		maxSize = 1000
	}

	return &Cache{
		maxSize: maxSize,
		ll:      list.New(),
		items:   make(map[string]*list.Element),
	}
}

func (c *Cache) Get(key string) (value string, contentType string, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	el, ok := c.items[key]
	if !ok {
		return "", "", false
	}

	ent := el.Value.(*entry)

	if ent.expiresAt != nil && *ent.expiresAt <= time.Now().Unix() {
		c.removeElement(el)
		return "", "", false
	}

	c.ll.MoveToFront(el)
	return ent.value, ent.contentType, true
}

func (c *Cache) Set(key, value string, contentType string, expiresAt *int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if el, ok := c.items[key]; ok {
		c.ll.MoveToFront(el)
		el.Value.(*entry).value = value
		el.Value.(*entry).expiresAt = expiresAt
		return
	}

	el := c.ll.PushFront(&entry{
		key:         key,
		value:       value,
		contentType: contentType,
		expiresAt:   expiresAt,
	})
	c.items[key] = el

	if c.ll.Len() > c.maxSize {
		c.removeOldest()
	}
}

func (c *Cache) removeOldest() {
	el := c.ll.Back()
	if el != nil {
		c.removeElement(el)
	}
}

func (c *Cache) removeElement(el *list.Element) {
	c.ll.Remove(el)
	ent := el.Value.(*entry)
	delete(c.items, ent.key)
}
