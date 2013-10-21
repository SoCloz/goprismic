package fragment

import (
	"fmt"
)

type Number struct {
	Value int `json:"value"`
}

func (n *Number) Decode(enc interface{}) error {
	dec, ok := enc.(Number)
	if !ok {
		return fmt.Errorf("%+v is not a string")
	}
	*n = dec
	return nil
}
