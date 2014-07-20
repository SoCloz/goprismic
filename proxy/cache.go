package proxy

import (
	"container/list"
	"sync"
	"time"
)

type Stats struct {
	Get      int
	Hit      int
	Miss     int
	Refresh  int
	RefreshError  int
	Eviction int
	Error   int
}

// A cache Entry
type CacheEntry struct {
	key          string
	entry        interface{}
	revision     string
	refreshing   bool
	validUntil   time.Time
}

// A cache
type Cache struct {
	sync.Mutex
	entries map[string]*list.Element
	lru     *list.List
	size    int
	ttl     time.Duration
	Stats   Stats
	revision string
}

type RefreshFn func() (interface{}, error)

// Creates a cache
func NewCache(size int, ttl time.Duration) *Cache {
	c := Cache{size: size, ttl: ttl}
	c.entries = make(map[string]*list.Element)
	c.lru = list.New()
	return &c
}

// Fetches something from the cache. If not found, refresh() will be called
func (c *Cache) Get(key string, refresh RefreshFn) (interface{}, error) {
	c.Lock()
	defer c.Unlock()
	c.Stats.Get++

	var err error
	e, found := c.get(key)

	if found && e.validUntil.After(time.Now()) {
		c.Stats.Hit++
		// bad revision - old content is always returned
		// and content is refreshed
		if c.revision != e.revision && !e.refreshing {
			e.refreshing = true
			go c.asyncRefresh(e, refresh)
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
	e.validUntil = time.Now().Add(c.ttl)
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
			if time.Now().Before(e.validUntil) {
				c.Stats.Eviction++
			}
		}
	}
}

// Clears the cache
func (c *Cache) Clear() {
	c.lru.Init()
	c.entries = make(map[string]*list.Element)
	c.Stats = Stats{}
}
