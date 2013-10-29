package proxy

import (
	"container/list"
	"sync"
	"time"
)

type Stats struct {
	Hit     int
	Miss    int
	Refresh int
}

// A cache Entry
type CacheEntry struct {
	key          string
	entry        interface{}
	validUntil   time.Time
	refreshAfter time.Time
}

// A cache
type Cache struct {
	sync.Mutex
	entries map[string]*list.Element
	lru     *list.List
	size    int
	ttl     time.Duration
	grace   time.Duration
	Stats   Stats
}

type RefreshFn func() (interface{}, error)

// Creates a cache
func NewCache(size int, ttl, grace time.Duration) *Cache {
	c := Cache{size: size, ttl: ttl, grace: grace}
	c.entries = make(map[string]*list.Element)
	c.lru = list.New()
	return &c
}

// Fetches something from the cache. If not found, refresh() will be called
func (c *Cache) Get(key string, refresh RefreshFn) (interface{}, error) {
	c.Lock()
	defer c.Unlock()

	needRefresh := false
	valid := false

	e, found := c.get(key)

	now := time.Now()
	if found {
		if now.After(e.validUntil) {
			needRefresh = true
			c.Stats.Miss++
		} else if now.After(e.refreshAfter) {
			needRefresh = true
			valid = true
			c.Stats.Hit++
			c.Stats.Refresh++
		} else {
			valid = true
			c.Stats.Hit++
		}
	} else {
		c.Stats.Miss++
	}

	var err error
	if !found || needRefresh {
		if valid {
			e.refreshAfter = now.Add(time.Duration(1) * time.Second)
			go c.refresh(e, refresh)
		} else if !found {
			e = &CacheEntry{}
			err = c.refresh(e, refresh)
			c.add(key, e)
		} else {
			err = c.refresh(e, refresh)
		}
	}
	return e.entry, err
}

func (c *Cache) refresh(e *CacheEntry, refresh RefreshFn) error {
	v, err := refresh()
	if err != nil {
		return err
	}
	e.entry = v
	n := time.Now()
	e.validUntil = n.Add(c.ttl)
	e.refreshAfter = e.validUntil.Add(-c.grace)
	return nil
}

func (c *Cache) get(key string) (*CacheEntry, bool) {
	el, found := c.entries[key]
	if found {
		c.lru.MoveToFront(el)
		return el.Value.(*CacheEntry), true
	}
	return nil, false
}

func (c *Cache) add(key string, e *CacheEntry) {
	el := &list.Element{Value: e}
	c.entries[key] = el
	c.lru.PushFront(el)
	for c.lru.Len() > c.size {
		del := c.lru.Back()
		delete(c.entries, del.Value.(*CacheEntry).key)
		c.lru.Remove(del)
	}
}

// Clears the cache
func (c *Cache) Clear() {
	c.lru.Init()
	c.entries = make(map[string]*list.Element)
	c.Stats = Stats{}
}
