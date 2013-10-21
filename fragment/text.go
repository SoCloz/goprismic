package fragment

import(
	"fmt"
)

type Text string

func (t *Text) Decode(enc interface{}) error {
	dec, ok := enc.(Text)
	if !ok {
		return fmt.Errorf("%+v is not a string")
	}
	*t = dec
	return nil
}