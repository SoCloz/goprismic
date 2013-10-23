package goprismic

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SoCloz/goprismic/fragment"
	"github.com/SoCloz/goprismic/fragment/link"
)

type FragmentTree map[string]Fragments

type Fragments map[string]FragmentList
type FragmentList []FragmentInterface

type FragmentEnvelope struct {
	Type  string `json:"type"`
	Value interface{}
}

type FragmentInterface interface {
	Decode(string, interface{}) error
	AsText() string
	AsHtml() string
	ResolveLinks(link.Resolver)
}

func (fs *FragmentList) UnmarshalJSON(data []byte) error {
	*fs = make(FragmentList, 0, 128)
	raw := []FragmentEnvelope{}
	if data[0] == '{' {
		data = append([]byte{byte('[')}, data...)
		data = append(data, byte(']'))
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	for _, v := range raw {
		var n FragmentInterface

		switch v.Type {
		case "StructuredText":
			st := make(fragment.StructuredText, 0, 128)
			n = &st
		case "Image":
			n = new(fragment.Image)
		case "Color":
			n = new(fragment.Color)
		case "Number":
			n = new(fragment.Number)
		case "Date":
			n = new(fragment.Date)
		case "Text":
			n = new(fragment.Text)
		case "Link.web":
			n = new(fragment.Link)
		case "Link.document":
			n = new(fragment.Link)
		case "Link.media":
			n = new(fragment.Link)
		case "Embed":
			n = new(fragment.Embed)
		case "Select":
			n = new(fragment.Text)
		default:
			return fmt.Errorf("Unknown fragment type %s", v.Type)
		}
		err := n.Decode(v.Type, v.Value)
		if err != nil {
			log.Printf("goprismic: unable to decode fragment : %s\n", err)
			return err
		}
		(*fs) = append(*fs, n)
	}
	return nil
}
