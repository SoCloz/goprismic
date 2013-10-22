package fragment

import (
	"fmt"
	"reflect"
)

type Number float64

func (n *Number) Decode(enc interface{}) error {
	dec, ok := enc.(Number)
	if !ok {
		return fmt.Errorf("unable to decode number fragment : %+v is a %s, not a number", enc, reflect.TypeOf(enc))
	}
	*n = dec
	return nil
}
