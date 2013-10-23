package fragment

import (
	"fmt"
	"reflect"

	"github.com/SoCloz/goprismic/fragment/link"
)

type ImageView struct {
	Url        string `json:"url"`
	Alt        string `json:"alt"`
	Copyright  string `json:"copyright"`
	Dimensions struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
}

func (i *ImageView) Decode(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("unable to decode image view : %+v is a %s, not a map", enc, reflect.TypeOf(enc))
	}
	if v, found := dec["url"]; found {
		i.Url = v.(string)
	}
	if v, found := dec["alt"]; found {
		i.Alt = v.(string)
	}
	if v, found := dec["copyright"]; found {
		i.Copyright = v.(string)
	}
	if d, found := dec["dimensions"]; found {
		dim, ok := d.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a map", d)
		}
		if v, found := dim["width"]; found {
			i.Dimensions.Width = int(v.(float64))
		}
		if v, found := dim["height"]; found {
			i.Dimensions.Height = int(v.(float64))
		}
	}
	return nil
}

func (i *ImageView) AsText() string {
	return i.Url
}

func (i *ImageView) AsHtml() string {
	return fmt.Sprintf("<img src=\"%s\" width=\"%d\" height=\"%d\"/>", i.Url, i.Dimensions.Width, i.Dimensions.Height)
}

func (i *ImageView) Ratio() float64 {
	return float64(i.Dimensions.Width)/float64(i.Dimensions.Height)
}

type Image struct {
	Main  ImageView            `json:"main"`
	Views map[string]ImageView `json:"views"`
}

// Returns a view of this image
func (i *Image) GetView(view string) (*ImageView, bool) {
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
		i.Views = make(map[string]ImageView)
		for k, view := range views {
			iv := &ImageView{}
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
