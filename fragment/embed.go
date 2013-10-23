package fragment

import (
	"github.com/SoCloz/goprismic/fragment/embed"
	"github.com/SoCloz/goprismic/fragment/link"
)

// A embed fragment (see http://oembed.com/)
type Embed struct {
	Embed *embed.Embed
}

func (e *Embed) Decode(_ string, enc interface{}) error {
	e.Embed = new(embed.Embed)
	return e.Embed.Decode(enc)
}

func (e *Embed) AsHtml() string {
	return e.Embed.AsHtml()
}

func (e *Embed) AsText() string {
	return ""
}

func (e *Embed) ResolveLinks(_ link.Resolver) {}