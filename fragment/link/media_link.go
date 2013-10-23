package link

import (
	"fmt"
	"reflect"
)

// A link to a file
type MediaLink struct {
	File struct {
		Url   string
		Kind string
		Filename string
	}
}

func (l *MediaLink) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to decode link content : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
	}
	if v, found := dec["file"]; found {
		doc, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a map", v)
		}
		if v, found := doc["url"]; found {
			l.File.Url = v.(string)
		}
		if v, found := doc["kind"]; found {
			l.File.Kind = v.(string)
		}
		if v, found := doc["filename"]; found {
			l.File.Filename = v.(string)
		}
	}
	return nil
}

func (l *MediaLink) GetUrl() string {
	return l.File.Url
}

func (l *MediaLink) GetText() string {
	return l.File.Filename
}

func (l *MediaLink) Resolve(_ Resolver) {}