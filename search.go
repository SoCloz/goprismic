package goprismic

import (
	"fmt"
)

type SearchForm struct {
	err  error
	api  *Api
	form *Form
	data map[string]string
	ref  Ref
}

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

// Searches the repository
func (s *SearchForm) Submit() ([]Document, error) {
	docs := make([]Document, 0, 1024)
	if s.err != nil {
		return docs, s.err
	}
	s.data["ref"] = s.ref.Ref
	err := s.api.call(s.form.Action, s.data, &docs)
	return docs, err
}
