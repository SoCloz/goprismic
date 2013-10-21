package block

type Block interface {
	Decode(interface{}) error
	AsHtml() string
	ParentHtmlTag() string
}
