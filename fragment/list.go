package fragment

import (
	"encoding/json"
	"fmt"
	"log"
)

type Tree map[string]Fragments

type Fragments map[string]List
type List []Interface

type Envelope struct {
	Type  string `json:"type"`
	Value interface{}
}

func (fs *List) UnmarshalJSON(data []byte) error {
	*fs = make(List, 0, 128)
	raw := []Envelope{}
	if data[0] == '{' {
		data = append([]byte{byte('[')}, data...)
		data = append(data, byte(']'))
	}
	err := json.Unmarshal(data, &raw)
	if err != nil {
		return err
	}
	for _, v := range raw {
		n, err := decode(v.Type, v.Value)
		if err != nil {
			return err
		}
		(*fs) = append(*fs, n)
	}
	return nil
}

func decode(t string, value interface{}) (Interface, error) {
	var n Interface
	var err error

	switch t {
	case "StructuredText":
		n, err = NewStructuredText(value)
	case "Image":
		n = new(Image)
	case "Color":
		n = new(Color)
	case "Number":
		n = new(Number)
	case "Date":
		n = new(Date)
	case "Text":
		n = new(Text)
	case "Link.web":
		n = new(Link)
	case "Link.document":
		n = new(Link)
	case "Link.media":
		n = new(Link)
	case "Embed":
		n = new(Embed)
	case "Select":
		n = new(Text)
	case "GeoPoint":
		n = new(GeoPoint)	
	case "Group":
		n, err = NewGroup(value)
	}
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, fmt.Errorf("goprismic: unable to decode fragment type %s", t)
	}
	err = n.Decode(t, value)
	if err != nil {
		log.Printf("goprismic: unable to decode fragment : %s\n", err)
		return nil, err
	}
	return n, nil
}
