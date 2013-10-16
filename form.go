package goprismic

// A form field
type Field struct {
	Type    string `json:"type"`
	Default string `json:"default"`
}

// A form
type Form struct {
	Name    string           `json:"name"`
	Method  string           `json:"method"`
	Rel     string           `json:"rel"`
	EncType string           `json:"enctype"`
	Action  string           `json:"action"`
	Fields  map[string]Field `json:"fields"`
}
