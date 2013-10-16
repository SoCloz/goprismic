package goprismic

type Span struct {
	Start int                    `json:"start"`
	End   int                    `json:"end"`
	Type  string                 `json:"type"`
	Data  map[string]interface{} `json:"data"`
}

type TextItem struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	Spans []Span `json:"spans"`
}

type StructuredTextContent struct {
	Value []TextItem `json:"value"`
}
