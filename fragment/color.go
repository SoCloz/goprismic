package fragment

import (
	"fmt"
	"reflect"
)

type Color string

func (c *Color) Decode(enc interface{}) error {
	res, ok := enc.(string)
	if !ok {
		return fmt.Errorf("unable to decode color fragment : %+v is a %s, not a string", enc, reflect.TypeOf(enc))
	}
	*c = Color(res)
	return nil
}
