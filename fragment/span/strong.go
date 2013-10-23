package span

import (
	"github.com/SoCloz/goprismic/fragment/link"
)

type Strong struct {
	Span
}

func (s *Strong) Decode(enc interface{}) error {
	return s.decodeSpan(enc)
}

func (s *Strong) HtmlBeginTag() string {
	return "<strong>"
}

func (s *Strong) HtmlEndTag() string {
	return "</strong>"
}

func (s *Strong) ResolveLinks(_ link.Resolver) {}
