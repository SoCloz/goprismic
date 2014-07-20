package goprismic

import (
	"fmt"
)

type PrismicError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
	Until   int64  `json:"until"`
}

func (e *PrismicError) Error() string {
	return fmt.Sprintf("%s error at line %d/column %d : \"%s\"", e.Type, e.Line, e.Column, e.Message)
}

func (e *PrismicError) IsOverCapacity() bool {
	return e.Message == "Too many requests"
}
