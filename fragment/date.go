package fragment

import (
	"fmt"
	"reflect"
	"time"
)

type Date struct {
	Value time.Time `json:"value"`
}

func (d *Date) Decode(enc interface{}) error {
	_, ok := enc.(string)
	if !ok {
		return fmt.Errorf("unable to decode date fragment : %+v is a %s, not a string", enc, reflect.TypeOf(enc))
	}
	//*d = res
	return nil
}
