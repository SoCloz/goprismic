package goprismic

import (
	"encoding/json"
)

type Fragments map[string]Fragment

type Fragment struct {
	Type    string          `json:"type"`
	Content FragmentContent `json:"-"`
}

type FragmentContent interface{}

func (fs Fragment) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	fs.Type = raw.Type
	switch raw.Type {
	case "StructuredText":
		fs.Content = new(StructuredTextContent)
	case "Image":
		fs.Content = new(ImageContent)
	}
	if fs.Content != nil {
		err = json.Unmarshal(data, &fs.Content)
	}
	return err
}
