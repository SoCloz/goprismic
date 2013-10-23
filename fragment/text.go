package fragment

import (
	"fmt"
	"reflect"

	"github.com/SoCloz/goprismic/fragment/link"
)

type Text string

func (t *Text) Decode(_ string, enc interface{}) error {
	dec, ok := enc.(string)
	if !ok {
		return fmt.Errorf("unable to decode text fragment : %+v is a %s, not a string", enc, reflect.TypeOf(enc))
	}
	*t = Text(dec)
	return nil
}

func (t *Text) AsText() string {
	return string(*t)
}

func (t *Text) AsHtml() string {
	return fmt.Sprintf("<span class=\"text\">%s</span>", *t)
}

func (t *Text) ResolveLinks(_ link.Resolver) {}