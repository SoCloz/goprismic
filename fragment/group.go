package fragment

import (
	"fmt"

	"github.com/SoCloz/goprismic/fragment/link"
)

type GroupFragments map[string]Interface
type Group []GroupFragments

func NewGroup(enc interface{}) (*Group, error) {
	dec, ok := enc.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%#v is not a slice", enc)
	}
	g := make(Group, 0, len(dec))
	return &g, nil
}

func (g *Group) Decode(_ string, enc interface{}) error {
	list := enc.([]interface{})
	for _, v := range list {
		gf := make(GroupFragments)
		lenc, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%#v is not a map", v)
		}
		for name, enci := range lenc {
			enc, ok := enci.(map[string]interface{})
			if !ok {
				return fmt.Errorf("%#v is not a map", enci)
			}
			n, err := decode(enc["type"].(string), enc["value"])
			if err != nil {
				return err
			}
			gf[name] = n
		}
		*g = append(*g, gf)
	}
	return nil
}

func (g *Group) AsHtml() string {
	html := ""
	for _, v := range *g {
		for name, fragment := range v {
			html += fmt.Sprintf("<section data-field=\"%s\">%s</section>", name, fragment.AsHtml())
		}
	}
	return html
}

func (g *Group) AsText() string {
	text := ""
	for _, v := range *g {
		for _, fragment := range v {
			if text != "" {
				text += "\n"
			}
			text += fragment.AsText()
		}
	}
	return text
}

func (g *Group) ResolveLinks(r link.Resolver) {
	for _, v := range *g {
		for _, fragment := range v {
			fragment.ResolveLinks(r)
		}
	}
}
