package block

import (
	"fmt"
)

type Heading struct {
	BaseBlock
}

func (h *Heading) AsHtml() string {
	switch h.Type {
	case "heading1":
		return fmt.Sprintf("<h1>%s</h1>", h.FormatHtmlText())
	case "heading2":
		return fmt.Sprintf("<h2>%s</h2>", h.FormatHtmlText())
	case "heading3":
		return fmt.Sprintf("<h3>%s</h3>", h.FormatHtmlText())
	case "heading4":
		return fmt.Sprintf("<h4>%s</h4>", h.FormatHtmlText())
	}
	return ""
}

func (h *Heading) ParentHtmlTag() string {
	return ""
}
