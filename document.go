package goprismic

import (
	"github.com/SoCloz/goprismic/fragment"
)

// A document is made of fragments of various types
type Document struct {
	Id        string       `json:"id"`
	Type      string       `json:"type"`
	Href      string       `json:"href"`
	Tags      []string     `json:"tags"`
	Slugs     []string     `json:"slugs"`
	Fragments FragmentTree `json:"data"`
}

// Returns the document slug
func (d *Document) GetSlug() string {
	return d.Slugs[0]
}

// Returns the list of fragments of a certain name
func (d *Document) GetFragments(field string) (FragmentList, bool) {
	frags, found := d.Fragments[d.Type]
	if !found {
		return nil, false
	}
	f, found := frags[field]
	return f, found
}

// Returns the nth fragment of a certain name
func (d *Document) GetFragmentAt(field string, index int) (FragmentInterface, bool) {
	frags, found := d.GetFragments(field)
	if !found || len(frags) < index {
		return nil, false
	}
	return frags[index], true
}

// Returns the first fragment of a certain name
func (d *Document) GetFragment(field string) (FragmentInterface, bool) {
	return d.GetFragmentAt(field, 0)
}

// Returns an image fragment (the first found)
func (d *Document) GetImage(field string) (*fragment.Image, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	i, ok := f.(*fragment.Image)
	if !ok {
		return nil, false
	}
	return i, true
}

// Returns s structured text fragment (the first found)
func (d *Document) GetStructuredText(field string) (*fragment.StructuredText, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.StructuredText)
	if !ok {
		return nil, false
	}
	return st, true
}

// Returns a color fragment (the first found)
func (d *Document) GetColor(field string) (*fragment.Color, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.Color)
	if !ok {
		return nil, false
	}
	return st, true
}

// Returns a number fragment (the first found)
func (d *Document) GetNumber(field string) (*fragment.Number, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.Number)
	if !ok {
		return nil, false
	}
	return st, true
}

// Returns a text fragment (the first found)
func (d *Document) GetText(field string) (*fragment.Text, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.Text)
	if !ok {
		return nil, false
	}
	return st, true
}

// Returns a document link fragment (the first found)
func (d *Document) GetDocumentLink(field string) (*fragment.DocumentLink, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.DocumentLink)
	if !ok {
		return nil, false
	}
	return st, true
}

// Returns a web link fragment (the first found)
func (d *Document) GetWebLink(field string) (*fragment.WebLink, bool) {
	f, found := d.GetFragment(field)
	if !found {
		return nil, false
	}
	st, ok := f.(*fragment.WebLink)
	if !ok {
		return nil, false
	}
	return st, true
}
