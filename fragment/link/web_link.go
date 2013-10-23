package link

import (
	"fmt"
	"reflect"
)

// A link to a website
type WebLink struct {
	Url string
}

func (l *WebLink) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to decode link content : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
	}
	if v, found := dec["url"]; found {
		l.Url = v.(string)
	}
	return nil
}

func (l *WebLink) GetUrl() string {
	return l.Url
}

func (l *WebLink) GetText() string {
	return l.Url
}

func (l *WebLink) Resolve(_ Resolver) {}
