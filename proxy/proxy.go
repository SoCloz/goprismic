package proxy

import (
	"fmt"
	"log"
	"sort"
	"time"

	"github.com/SoCloz/goprismic"
)

type Config struct {
	// Cache size
	CacheSize int
	// TTL
	// Cache is asynchronous during the TTL duration after an update.
	TTL time.Duration
	// API Master ref refresh frequency
	// The proxy will check at this interval if something has changed.
	MasterRefresh time.Duration
	// Base refresh chance after a master reference update, between 0 and 1
	BaselineRefreshChance float32
	// inherited from api
	debug bool
}

type Proxy struct {
	cache *Cache
	api   *goprismic.Api

	Config Config

	lastRefresh           time.Time
}

// Creates a proxy
//
// All documents are cached in a LRU cache of cacheSize elements.
func New(u, accessToken string, apiCfg goprismic.Config, cfg Config) (*Proxy, error) {
	a, err := goprismic.Get(u, accessToken, apiCfg)
	if err != nil {
		return nil, err
	}
	if cfg.MasterRefresh == 0 {
		cfg.MasterRefresh = time.Minute
	}
	cfg.debug = apiCfg.Debug
	if cfg.BaselineRefreshChance == 0 {
		cfg.BaselineRefreshChance = 1.0
	}
	c := NewCache(cfg.CacheSize, cfg.TTL, 1.0)
	c.revision = a.GetMasterRef()
	p := &Proxy{
		cache:                 c,
		api:                   a,
		Config:                cfg,
		lastRefresh:           time.Now(),
	}
	go p.loopRefresh()
	return p, nil
}

// Returns the cache stats
func (p *Proxy) GetStats() Stats {
	return p.cache.Stats
}

// Queries the API
//
//   proxy.Direct().Master().Form("everything").Submit()
func (p *Proxy) Direct() *goprismic.Api {
	return p.api
}

// Fetches a document by id
func (p *Proxy) GetDocument(id string) (*goprismic.Document, error) {
	sr, err := p.Search().Form("everything").Query("[[:d = at(document.id, \"" + id + "\")]]").Submit()
	if err != nil || sr.TotalResults == 0 {
		return nil, err
	}
	return &sr.Results[0], nil

}

// Fetches a document of a specific type by a field value
func (p *Proxy) GetDocumentBy(docType, field string, value interface{}) (*goprismic.Document, error) {
	query := fmt.Sprintf("[[:d = at(my.%s.%s, \"%v\")][:d = any(document.type, [\"%s\"])]]", docType, field, value, docType)
	sr, err := p.Search().Form("everything").Query(query).Submit()
	if err != nil || sr.TotalResults == 0 {
		return nil, err
	}
	return &sr.Results[0], nil
}

// Search documents
func (p *Proxy) Search() *SearchForm {
	f := &SearchForm{sf: p.api.Master(), p: p}
	f.sig = sort.StringSlice{}
	return f
}

// Refreshes the master ref, returns true if master ref has changed
func (p *Proxy) Refresh() bool {
	p.api.Refresh()
	if p.cache.updateRevision(p.api.GetMasterRef()) {
		p.lastRefresh = time.Now()
		// refresh : 100% cache miss expected => we switch to baseline
		p.cache.refreshChance = p.Config.BaselineRefreshChance
		if p.Config.debug {
			log.Printf("Prismic - refreshing master ref and lowering refresh chance to %.1f%%", p.cache.refreshChance*100)
		}
		return true
	}
	return false
}

// Refreshes the master ref, returns true if master ref has changed
func (p *Proxy) loopRefresh() {
	prevRefreshError, prevRefresh, lastNoError := 0, 0, 0
	tick := time.Tick(p.Config.MasterRefresh)
	for {
		select {
		case <-tick:
			refreshed := p.Refresh()
			if refreshed {
				lastNoError = 0
				continue
			}

			deltaRefreshError := p.cache.Stats.RefreshError - prevRefreshError
			deltaRefresh := p.cache.Stats.Refresh - prevRefresh
			if deltaRefresh+deltaRefreshError == 0 && p.cache.refreshChance == 1.0 {
				continue
			}

			prevRefreshError, prevRefresh = p.cache.Stats.RefreshError, p.cache.Stats.Refresh
			var refreshChance float32
			if deltaRefreshError > 0 {
				lastNoError = 0
				refreshChance = float32(deltaRefresh) / float32(deltaRefresh+deltaRefreshError) * p.cache.refreshChance
			} else {
				lastNoError++
				refreshChance = p.Config.BaselineRefreshChance*(1.0+float32(lastNoError*(lastNoError-1))/2)
				if refreshChance > 1.0 {
					refreshChance = 1.0
				}
			}
			if refreshChance < p.cache.refreshChance {
				p.cache.refreshChance = refreshChance
				if p.Config.debug {
					log.Printf("Prismic - lowering refresh chance to %.2f%%", p.cache.refreshChance*100.0)
				}
			}
			if refreshChance > p.cache.refreshChance {
				if refreshChance >= 1.0 && p.cache.refreshChance >= 0.99 {
					p.cache.refreshChance = 1.0
				} else {
					p.cache.refreshChance = (p.cache.refreshChance + refreshChance) / 2
				}
				if p.Config.debug {
					log.Printf("Prismic - raising refresh chance to %.2f%%", p.cache.refreshChance*100.0)
				}
			}
		}
	}
}

// Get something from the proxy
func (p *Proxy) Get(key string, refresh RefreshFn) (interface{}, error) {
	return p.cache.Get(key, refresh)
}

// Clears the cache
func (p *Proxy) Clear() {
	p.cache.Clear()
}
