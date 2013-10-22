package fragment

import (
	"fmt"
	"reflect"
	"time"
)

type Date time.Time

func (d *Date) Decode(enc interface{}) error {
	dec, ok := enc.(string)
	if !ok {
		return fmt.Errorf("unable to decode date fragment : %+v is a %s, not a string", enc, reflect.TypeOf(enc))
	}
	date, err := time.Parse("2006-01-02", dec)
	if err != nil {
		return err
	}
	*d = Date(date)
	return nil
}

func (d *Date) AsText() string {
	t := time.Time(*d)
	return t.Format("2006-01-02")
}

func (d *Date) AsHtml() string {
	t := time.Time(*d)
	return "<time>"+t.Format("2006-01-02")+"</time>"
}
