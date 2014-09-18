package proxy

import (
	"container/list"
	"math/rand"
	"sync"
	"time"
)

type Stats struct {
	Get          int
	Hit          int
	Miss         int
	Refresh      int
	RefreshError int
	Eviction     int
	Error        int
}

// A cache Entry
type CacheEntry struct {
	key        string
	entry      interface{}
	revision   string
	refreshing bool
}

// A cache
type Cache struct {
	sync.Mutex
	entries       map[string]*list.Element
	lru           *list.List
	size          int
	ttl           time.Duration
	Stats         Stats
	refreshChance float32
	revision      string
	asyncDeadline time.Time
}

type RefreshFn func() (interface{}, error)

// Creates a cache
func NewCache(size int, ttl time.Duration, refreshChance float32) *Cache {
	return &Cache{
		size:          size,
		ttl:           ttl,
		entries:       make(map[string]*list.Element),
		lru:           list.New(),
		refreshChance: refreshChance,
	}
}

// Fetches something from the cache. If not found, refresh() will be called
func (c *Cache) Get(key string, refresh RefreshFn) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	c.Stats.Get++

	var err error
	e, found := c.get(key)

	if found && (e.refreshing || c.asyncDeadline.IsZero() || c.asyncDeadline.After(time.Now())) {
		c.Stats.Hit++
		// bad revision - old content is always returned
		// and content is refreshed
		if c.revision != e.revision && !e.refreshing {
			if c.refreshChance == 1.0 || rand.Float32() <= c.refreshChance {
				e.refreshing = true
				go c.asyncRefresh(e, refresh)
			}
		}
	} else {
		c.Stats.Miss++
		if !found {
			e = &CacheEntry{}
			err = c.refresh(e, refresh)
			if err == nil {
				c.add(key, e)
			}
		} else {
			err = c.refresh(e, refresh)
		}
	}
	if err != nil {
		c.Stats.Error++
	}
	return e.entry, err
}

// Adds something to the cache
func (c *Cache) Set(key string, v interface{}) {
	e := &CacheEntry{}
	c.refresh(e, func() (interface{}, error) { return v, nil })
	c.add(key, e)
}

func (c *Cache) asyncRefresh(e *CacheEntry, refresh RefreshFn) {
	err := c.refresh(e, refresh)
	if err == nil {
		c.Stats.Refresh++
	} else {
		c.Stats.RefreshError++
	}
}
func (c *Cache) refresh(e *CacheEntry, refresh RefreshFn) error {
	v, err := refresh()
	e.refreshing = false
	if err != nil {
		return err
	}
	e.entry = v
	e.revision = c.revision
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
	if el, found := c.entries[key]; found {
		c.lru.MoveToFront(el)
	} else {
		e.key = key
		el := c.lru.PushFront(e)
		c.entries[key] = el
		for c.lru.Len() > c.size {
			del := c.lru.Back()
			delete(c.entries, del.Value.(*CacheEntry).key)
			c.lru.Remove(del)
			if c.revision == del.Value.(*CacheEntry).revision {
				c.Stats.Eviction++
			}
		}
	}
}

func (c *Cache) updateRevision(revision string) bool {
	if revision == c.revision {
		return false
	}
	c.revision = revision
	if c.ttl > 0 {
		c.asyncDeadline = time.Now().Add(c.ttl)
	}
	return true
}

// Clears the cache
func (c *Cache) Clear() {
	c.lru.Init()
	c.entries = make(map[string]*list.Element)
	c.Stats = Stats{}
}
