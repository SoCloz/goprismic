package fragment

import(
	"fmt"
	"time"
)

type Date struct {
	Value time.Time `json:"value"`
}

func (d *Date) Decode(enc interface{}) error {
	_, ok := enc.(string)
	if !ok {
		return fmt.Errorf("%+v is not a string")
	}
	//*d = res
	return nil
}