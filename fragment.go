package goprismic

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/SoCloz/goprismic/fragment"
)

type FragmentTree map[string]Fragments

type Fragments map[string]FragmentList
type FragmentList []FragmentInterface

type FragmentEnvelope struct {
	Type  string `json:"type"`
	Value interface{}
}

type FragmentInterface interface {
	Decode(interface{}) error
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
			n = new(fragment.WebLink)
		case "Link.document":
			n = new(fragment.DocumentLink)
		default:
			panic(fmt.Sprintf("Unknown fragment type %s", v.Type))
		}
		err := n.Decode(v.Value)
		if err != nil {
			log.Printf("goprismic: unable to decode fragment : %s\n", err)
			return err
		}
		//fmt.Printf("\n%s => %+v\n", v.Type, n)
		(*fs) = append(*fs, n)
	}
	return nil
}
