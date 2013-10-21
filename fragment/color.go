package fragment

import (
	"fmt"
)

type Color string

func (c *Color) Decode(enc interface{}) error {
	res, ok := enc.(Color)
	if !ok {
		return fmt.Errorf("%+v is not a string")
	}
	*c = res
	return nil
}
