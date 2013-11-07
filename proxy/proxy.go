package proxy

import (
	"fmt"
	"time"

	"github.com/SoCloz/goprismic"
)

type Proxy struct {
	cache *Cache
	api   *goprismic.Api
}

// Creates a proxy
//
// All documents are cached in a LRU cache of cacheSize elements.
// Documents are cached for a maximum time of ttl, and will be asynchronously refreshed between ttl-grace and ttl.
func New(u, accessToken string, cacheSize int, ttl, grace time.Duration) (*Proxy, error) {
	a, err := goprismic.Get(u, accessToken)
	if err != nil {
		return nil, err
	}
	c := NewCache(cacheSize, ttl, grace)
	return &Proxy{cache: c, api: a}, nil
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
	key := fmt.Sprintf("byid++%s", id)
	d, err := p.getDoc(key, "[[:d = at(document.id, \""+id+"\")]]")
	return d, err
}

// Fetched a document by a field value
func (p *Proxy) GetDocumentBy(form, field, value interface{}) (*goprismic.Document, error) {
	key := fmt.Sprintf("by++%s++%s++%s", form, field, value)
	query := fmt.Sprintf("[[:d = at(my.%s.%s, \"%v\")]]", form, field, value)
	d, err := p.getDoc(key, query)
	return d, err
}

func (p *Proxy) getDoc(key, query string) (*goprismic.Document, error) {
	d, err := p.cache.Get(key, func() (interface{}, error) {
		docs, err := p.api.Master().Form("everything").Query(query).Submit()
		if err != nil {
			p.api.Refresh()
			docs, err = p.api.Master().Form("everything").Query(query).Submit()
			if err != nil {
				return nil, err
			}
		}
		if len(docs) == 0 {
			return nil, nil
		}
		return &docs[0], nil
	})
	if d != nil {
		return d.(*goprismic.Document), nil
	}
	return nil, err
}

// Clears the cache
func (p *Proxy) Clear() {
	p.cache.Clear()
}
