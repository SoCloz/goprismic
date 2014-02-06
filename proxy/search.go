package proxy

import (
	"fmt"
	"sort"
	"strings"

	"github.com/SoCloz/goprismic"
)

type SearchForm struct {
	sf  *goprismic.SearchForm
	p   *Proxy
	sig sort.StringSlice
}

// Returns the error
func (s *SearchForm) Error() error {
	return s.sf.Error()
}

// Selects the form on which to use
func (s *SearchForm) Form(name string) *SearchForm {
	s.sig = append(s.sig, fmt.Sprintf("form-%s", name))
	s.sf.Form(name)
	return s
}

// Query the form using a predicate query ([[:d = any(document.type, ["article"])]])
func (s *SearchForm) Query(query string) *SearchForm {
	s.sig = append(s.sig, fmt.Sprintf("q-%s", query))
	s.sf.Query(query)
	return s
}

// Adds form data
func (s *SearchForm) Data(data map[string]string) *SearchForm {
	for k, v := range data {
		s.sig = append(s.sig, fmt.Sprintf("data-%s-%s", k, v))
	}
	s.sf.Data(data)
	return s
}

// Sets the page number
func (s *SearchForm) Page(page int) *SearchForm {
	s.sig = append(s.sig, fmt.Sprintf("page-%d", page))
	s.sf.Page(page)
	return s
}

// Sets the page size
func (s *SearchForm) PageSize(pageSize int) *SearchForm {
	s.sig = append(s.sig, fmt.Sprintf("pageSize-%d", pageSize))
	s.sf.PageSize(pageSize)
	return s
}

// Order result - can be chained multiple times
//
// example : form.Order("my.product.name", OrderAsc).Order("my.product.size", OrderDesc)
func (s *SearchForm) Order(field string, order int) *SearchForm {
	s.sig = append(s.sig, fmt.Sprintf("order-%s-%d", field, order))
	s.sf.Order(field, order)
	return s
}

// Searches the repository
func (s *SearchForm) Submit() (*goprismic.SearchResult, error) {
	sort.Sort(s.sig)
	key := strings.Join(s.sig, ",")
	sr, err := s.p.Get(key, func() (interface{}, error) {
		sr, err := s.sf.Submit()
		if err != nil {
			s.p.Refresh()
			sr, err = s.sf.Submit()
		}
		return sr, err
	})
	if sr != nil {
		return sr.(*goprismic.SearchResult), nil
	}
	return &goprismic.SearchResult{}, err
}
