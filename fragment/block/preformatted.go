package block

import (
	"fmt"
)

// A preformatted block
type Preformatted struct {
	BaseBlock
}

func (p *Preformatted) AsHtml() string {
	return fmt.Sprintf("<pre>%s</pre>", p.FormatHtmlText())
}

func (p *Preformatted) ParentHtmlTag() string {
	return ""
}
