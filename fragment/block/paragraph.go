package block

import (
	"fmt"
)

// A paragraph block
type Paragraph struct {
	BaseBlock
}

func (p *Paragraph) Decode(enc interface{}) error {
	return p.decodeBlock(enc)
}

func (p *Paragraph) AsHtml() string {
	return fmt.Sprintf("<p>%s</p>", p.FormatHtmlText())
}

func (p *Paragraph) ParentHtmlTag() string {
	return ""
}
