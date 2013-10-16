package goprismic

type ImageView struct {
	Url        string `json:"url"`
	Dimensions struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}
}

type ImageContent struct {
	Value struct {
		Main  ImageView            `json:"main"`
		Views map[string]ImageView `json:"views"`
	} `json:"value"`
}
