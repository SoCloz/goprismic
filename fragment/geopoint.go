package fragment

import (
	"fmt"
	"github.com/SoCloz/goprismic/fragment/link"
)

// A Geopoint
type GeoPoint struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (g *GeoPoint) Decode(_ string, enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
	}
	if v, found := dec["latitude"]; found {
		g.Latitude = v.(float64)
	}
	if v, found := dec["longitude"]; found {
		g.Longitude = v.(float64)
	}
	return nil
}

func (g *GeoPoint) AsText() string {
	return fmt.Sprintf("%f,%f", g.Latitude, g.Longitude)
}

func (g *GeoPoint) AsHtml() string {
	return fmt.Sprintf(`<div class="geopoint"><span class="latitude">%f</span><span class="longitude">%f</span></div>`, g.Latitude, g.Longitude)
}

func (g *GeoPoint) ResolveLinks(_ link.Resolver) {}
