package fragment

import (
	"fmt"

	"github.com/SoCloz/goprismic/fragment/link"
)

// A number
type Number float64

func (n *Number) Decode(_ string, enc interface{}) error {
	dec, ok := enc.(float64)
	if !ok {
		*n = Number(0)
		return nil
	}
	*n = Number(dec)
	return nil
}

func (n *Number) AsText() string {
	return fmt.Sprintf("%f", *n)
}

func (n *Number) AsHtml() string {
	return fmt.Sprintf("<span class=\"number\">%f</span>", *n)
}

func (n *Number) ResolveLinks(_ link.Resolver) {}
