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

// Fetches a document of a specific type by a field value
func (p *Proxy) GetDocumentBy(docType, field string, value interface{}) (*goprismic.Document, error) {
	key := fmt.Sprintf("by++%s++%s++%s", docType, field, value)
	query := fmt.Sprintf("[[:d = at(my.%s.%s, \"%v\")][:d = any(document.type, [\"%s\"])]]", docType, field, value, docType)
	d, err := p.getDoc(key, query)
	return d, err
}

// Search documents
func (p *Proxy) Search(docType, q string) ([]goprismic.Document, error) {
	key := fmt.Sprintf("search++%s++%s", docType, q)
	query := fmt.Sprintf("[%s[:d = any(document.type, [\"%s\"])]]", q, docType)
	d, err := p.getDocs(key, query)
	return d, err
}

func (p *Proxy) getDoc(key, query string) (*goprismic.Document, error) {
	d, err := p.cache.Get(key, func() (interface{}, error) {
		docs, err := p.query(query)
		if err != nil || len(docs) == 0 {
			return nil, err
		}
		return &docs[0], nil
	})
	if d != nil {
		return d.(*goprismic.Document), nil
	}
	return nil, err
}

func (p *Proxy) getDocs(key, query string) ([]goprismic.Document, error) {
	d, err := p.cache.Get(key, func() (interface{}, error) {
		return p.query(query)
	})
	if d != nil {
		return d.([]goprismic.Document), nil
	}
	return []goprismic.Document{}, err
}

func (p *Proxy) query(query string) ([]goprismic.Document, error) {
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
		for _, d := range docs {
			p.addToCache(&d)
		}
		return docs, nil
}

func (p *Proxy) addToCache(d *goprismic.Document) {
	p.cache.Set(fmt.Sprintf("byid++%s", d.Id), d)
}

// Clears the cache
func (p *Proxy) Clear() {
	p.cache.Clear()
}
