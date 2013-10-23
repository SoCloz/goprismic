package link

import(
	"fmt"
)

func Decode(t string, enc interface{}) (Link, error) {
	var n Link
	switch t {
	case "Link.web":
		n = new(WebLink)
	case "Link.document":
		n = new(DocumentLink)
	case "Link.file":
		n = new(MediaLink)
	default:
		return nil, fmt.Errorf("Unknown link type %s", t)
	}
	err := n.Decode(enc)
	return n, err
}