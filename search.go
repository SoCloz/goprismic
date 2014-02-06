package goprismic

import (
	"encoding/json"
	"fmt"
)

type SearchForm struct {
	err  error
	api  *Api
	form *Form
	data map[string]string
	ref  Ref
}

type SearchResult struct {
	Page           int        `json:"page"`
	TotalPages     int        `json:"total_pages"`
	ResultsPerPage int        `json:"results_per_page"`
	TotalResults   int        `json:"total_results_size"`
	Results        []Document `json:"results"`
}

const (
	OrderAsc = iota
	OrderDesc
)

// Returns the error
func (s *SearchForm) Error() error {
	return s.err
}

// Selects the form on which to use
func (s *SearchForm) Form(name string) *SearchForm {
	if s.err != nil {
		return s
	}
	form, found := s.api.Data.Forms[name]
	if !found {
		s.err = fmt.Errorf("form %s not found", name)
	}
	s.form = &form
	for name, f := range form.Fields {
		if f.Default != "" {
			s.data[name] = f.Default
		}
	}
	return s
}

// Query the form using a predicate query ([[:d = any(document.type, ["article"])]])
func (s *SearchForm) Query(query string) *SearchForm {
	if s.err != nil {
		return s
	}
	if s.form == nil {
		s.err = fmt.Errorf("no form set !")
		return s
	}
	query = stripQuery(query)
	if field, found := s.form.Fields["q"]; found && field.Default != "" {
		query += stripQuery(field.Default)
	}
	s.data["q"] = "[" + query + "]"
	return s
}

// Adds form data
func (s *SearchForm) Data(data map[string]string) *SearchForm {
	if s.err != nil {
		return s
	}
	if s.form == nil {
		s.err = fmt.Errorf("no form set !")
		return s
	}
	for k, v := range data {
		if _, found := s.form.Fields[k]; !found {
			s.err = fmt.Errorf("field %s not found in form %s !", k, s.form.Name)
			return s
		}
		s.data[k] = v
	}
	return s
}

// Sets the page number
func (s *SearchForm) Page(page int) *SearchForm {
	s.data["page"] = fmt.Sprintf("%d", page)
	return s
}

// Sets the page size
func (s *SearchForm) PageSize(pageSize int) *SearchForm {
	s.data["pageSize"] = fmt.Sprintf("%d", pageSize)
	return s
}

// Order result - can be chained multiple times
//
// example : form.Order("my.product.name", OrderAsc).Order("my.product.size", OrderDesc)
func (s *SearchForm) Order(field string, order int) *SearchForm {
	if _, found := s.data["orderings"]; !found {
		s.data["orderings"] = ""
	}
	if order == OrderDesc {
		field = fmt.Sprintf("%s desc", field)
	}
	s.data["orderings"] = fmt.Sprintf("%s[%s]", s.data["orderings"], field)
	return s
}

// Searches the repository
func (s *SearchForm) Submit() (*SearchResult, error) {
	sr := SearchResult{}
	sr.Results = make([]Document, 0, 1024)
	if s.err != nil {
		return &sr, s.err
	}
	s.data["ref"] = s.ref.Ref
	err := s.api.call(s.form.Action, s.data, &sr)
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		err = s.api.call(s.form.Action, s.data, &sr.Results)
	}
	return &sr, err
}
