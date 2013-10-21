package fragment

import(
	"fmt"
)

type WebLink struct {
	Url string
}

func (l *WebLink) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
	}
	if v, found := dec["url"]; found {
		l.Url = v.(string)
	}
	return nil
}

type DocumentLink struct {
	Document struct {
		Id string
		Type string
		Slug string
	}
	IsBroken bool
}

func (l *DocumentLink) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
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
