package block

type Embed struct {
	BaseBlock
}

func (e *Embed) AsHtml() string {
	return e.Text
}

func (e *Embed) ParentHtmlTag() string {
	return ""
}
