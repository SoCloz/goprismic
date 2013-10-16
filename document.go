package goprismic

// A document is made of fragments of various types
type Document struct {
	Id        string               `json:"id"`
	Type      string               `json:"type"`
	Href      string               `json:"href"`
	Tags      []string             `json:"tags"`
	Slugs     []string             `json:"slugs"`
	Fragments map[string]Fragments `json:"data"`
}

func (d *Document) GetSlug() string {
	return d.Slugs[0]
}
