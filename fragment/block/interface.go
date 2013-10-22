package block

type Block interface {
	Decode(interface{}) error
	AsHtml() string
	AsText() string
	ParentHtmlTag() string
}
