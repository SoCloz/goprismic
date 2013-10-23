package span

import(
	"fmt"

	"github.com/SoCloz/goprismic/fragment/link"
)

type Hyperlink struct {
	Span
	Link link.Link
}

func (s *Hyperlink) Decode(enc interface{}) error {
	err := s.decodeSpan(enc)
	if err != nil {
		return err
	}

	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
	}
	v, found := dec["data"]
	if !found {
		return fmt.Errorf("No data found for link %+v", enc)
	}
	dec, ok = v.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", v)
	}
	t, found := dec["type"]
	if !found {
		return fmt.Errorf("No type found for link %+v", v)
	}
	l, foundlink := dec["value"]
	if !foundlink {
		err = fmt.Errorf("No link value found in %+v", enc)
	}
	s.Link, err = link.Decode(t.(string), l)
	fmt.Println("%s / %+v => %+v", t.(string), v, s.Link)
	return err
}

func (h *Hyperlink) HtmlBeginTag() string {
	return fmt.Sprintf("<a href=\"%s\">", h.Link.GetUrl())
}

func (h *Hyperlink) HtmlEndTag() string {
	return "</a>"
}

func (h *Hyperlink) ResolveLinks(r link.Resolver) {
	h.Link.Resolve(r)
}