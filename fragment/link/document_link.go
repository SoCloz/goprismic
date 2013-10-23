package link

import (
	"fmt"
	"reflect"
)

// A link to a prismic document
type DocumentLink struct {
	Document struct {
		Id   string
		Type string
		Slug string
	}
	Url      string
	IsBroken bool
}

func (l *DocumentLink) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to decode link content : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
	}
	if v, found := dec["document"]; found {
		doc, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a map", v)
		}
		if v, found := doc["id"]; found {
			l.Document.Id = v.(string)
		}
		if v, found := doc["type"]; found {
			l.Document.Type = v.(string)
		}
		if v, found := doc["slug"]; found {
			l.Document.Slug = v.(string)
		}
	}
	if v, found := dec["isBroken"]; found {
		l.IsBroken = v.(bool)
	}
	return nil
}

func (l *DocumentLink) GetUrl() string {
	if l.Url == "" {
		panic("url with no link - link might have not been resolved")
	}
	return l.Url
}

func (l *DocumentLink) GetText() string {
	if l.Url == "" {
		panic("url with no link - link might have not been resolved")
	}
	return l.Url
}

func (l *DocumentLink) Resolve(r Resolver) {
	l.Url = r(l)
}
