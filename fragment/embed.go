package fragment

import (
	"fmt"
	"reflect"

	"github.com/SoCloz/goprismic/fragment/link"
)

// A embed fragment (see http://oembed.com/)
type Embed struct {
	Type  string
	ProviderName string
	EmbedUrl string
	Width int
	Height int
	Html string
}

func (e *Embed) Decode(_ string, enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to decode embed fragment : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
	}
	if v, found := dec["oembed"]; found {
		doc, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to decode embed fragment : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
		}
		if v, found := doc["type"]; found {
			e.Type = v.(string)
		}
		if v, found := doc["provider_name"]; found {
			e.ProviderName = v.(string)
		}
		if v, found := doc["embed_url"]; found {
			e.EmbedUrl = v.(string)
		}
		if v, found := doc["width"]; found {
			e.Width = int(v.(float64))
		}
		if v, found := doc["height"]; found {
			e.Height = int(v.(float64))
		}
		if v, found := doc["html"]; found {
			e.Html = v.(string)
		}
	}
	return nil
}

func (e *Embed) AsText() string {
	return ""
}

func (e *Embed) AsHtml() string {
	if e.Html != "" {
		return "<div data-oembed=\""+e.EmbedUrl+"\" data-oembed-type=\""+e.Type+"\" data-oembed-provider=\""+e.ProviderName+"\">"+e.Html+"</div>";
	}
	return ""
}

func (e *Embed) ResolveLinks(_ link.Resolver) {}