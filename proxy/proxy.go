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
	// Documents are cached for a maximum time of ttl.
	TTL time.Duration
	// API Master ref refresh frequency
	MasterRefresh time.Duration
	// Base refresh chance after master reference chance
	BaselineRefreshChance float32
	// inherited from api
	debug bool
}

type Proxy struct {
	cache *Cache
	api   *goprismic.Api

	Config Config

	lastRefresh           time.Time
	baselineRefreshChance float32
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
		baselineRefreshChance: cfg.BaselineRefreshChance,
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
	oldRev := p.cache.revision
	p.cache.revision = p.api.GetMasterRef()
	if oldRev != p.cache.revision {
		p.lastRefresh = time.Now()
		// refresh : 100% cache miss expected => we switch to baseline
		p.cache.refreshChance = p.baselineRefreshChance
		if p.Config.debug {
			log.Printf("Prismic - refreshing master ref and lowering refresh chance to %.1f%%", p.cache.refreshChance*100)
		}
		return true
	}
	return false
}

// Refreshes the master ref, returns true if master ref has changed
func (p *Proxy) loopRefresh() {
	prevRefreshError, prevRefresh := 0, 0
	tick := time.Tick(p.Config.MasterRefresh)
	for {
		select {
		case <-tick:
			refreshed := p.Refresh()
			deltaRefreshError := p.cache.Stats.RefreshError - prevRefreshError
			deltaRefresh := p.cache.Stats.Refresh - prevRefresh
			// ensure that chance is always > 0 and valid
			refreshChance := float32(deltaRefresh+1) / float32(deltaRefresh+deltaRefreshError+1)
			if refreshChance < p.baselineRefreshChance {
				p.baselineRefreshChance = refreshChance
				if p.Config.debug {
					log.Printf("Prismic - lowering baseline refresh chance to %.2f%%", p.baselineRefreshChance*100.0)
				}
			}
			if !refreshed {
				if refreshChance < p.cache.refreshChance {
					p.cache.refreshChance = refreshChance
					if p.Config.debug {
						log.Printf("Prismic - lowering refresh chance to %.2f%%", p.cache.refreshChance*100.0)
					}
				}
				if refreshChance > p.cache.refreshChance {
					if refreshChance < p.cache.refreshChance*1.1 {
						p.cache.refreshChance = refreshChance
					} else {
						p.cache.refreshChance = p.cache.refreshChance * 1.1
					}
					if p.cache.refreshChance >= 1.0 {
						p.cache.refreshChance = 1.0
					}
					if p.Config.debug {
						log.Printf("Prismic - raising refresh chance to %.2f%%", p.cache.refreshChance*100.0)
					}
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
