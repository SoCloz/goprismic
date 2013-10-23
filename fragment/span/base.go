package span

import (
	"fmt"
)

type Span struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

func (s *Span) decodeSpan(enc interface{}) error {
	dec, ok := enc.(map[string]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a map", enc)
	}
	if v, found := dec["start"]; found {
		s.Start = int(v.(float64))
	}
	if v, found := dec["end"]; found {
		s.End = int(v.(float64))
	}
	return nil
}

func (s *Span) GetStart() int {
	return s.Start
}

func (s *Span) GetEnd() int {
	return s.End
}
