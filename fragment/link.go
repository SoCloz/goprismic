package fragment

import (
	"fmt"

	"github.com/SoCloz/goprismic/fragment/link"
)

// A link fragment
type Link struct {
	Link link.Link
}

func (l *Link) Decode(t string, enc interface{}) error {
	var err error
	l.Link, err = link.Decode(t, enc)
	return err
}

func (l *Link) AsText() string {
	return l.Link.GetUrl()
}

func (l *Link) AsHtml() string {
	return fmt.Sprintf("<a href=\"%s\">%s</a>", l.Link.GetUrl(), l.Link.GetText())
}

func (l *Link) ResolveLinks(r link.Resolver) {
	l.Link.Resolve(r)
}