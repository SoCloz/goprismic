package link

type Resolver func(Link) string

type Link interface {
	Decode(interface{}) error
	GetUrl() string
	GetText() string
	Resolve(Resolver)
}