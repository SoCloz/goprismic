package fragment

import (
	"github.com/SoCloz/goprismic/fragment/link"
)

type Interface interface {
	Decode(string, interface{}) error
	AsText() string
	AsHtml() string
	ResolveLinks(link.Resolver)
}
