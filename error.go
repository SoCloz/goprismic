package goprismic

import (
	"fmt"
)

type PrismicError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}

func (e *PrismicError) Error() string {
	return fmt.Sprintf("%s error at line %d/column %d : \"%s\"", e.Type, e.Line, e.Column, e.Message)
}
