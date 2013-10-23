package fragment

import (
	"fmt"
	"reflect"

	"github.com/SoCloz/goprismic/fragment/image"
	"github.com/SoCloz/goprismic/fragment/link"
)

// An image, with multiple views
type Image struct {
	Main  image.View
	Views map[string]image.View
}

// Returns a view of this image
func (i *Image) GetView(view string) (*image.View, bool) {
	v, found := i.Views[view]
	return &v, found
}

func (i *Image) Decode(_ string, enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
	}
	if v, found := dec["main"]; found {
		(&i.Main).Decode(v)
	}
	if v, found := dec["views"]; found {
		views, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("unable to decode image fragment : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
		}
		i.Views = make(map[string]image.View)
		for k, view := range views {
			iv := &image.View{}
			iv.Decode(view)
			i.Views[k] = *iv
		}
	}
	return nil
}

func (i *Image) AsText() string {
	return i.Main.AsText()
}

func (i *Image) AsHtml() string {
	return i.Main.AsHtml()
}

func (i *Image) ResolveLinks(_ link.Resolver) {}
