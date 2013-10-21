package span

type Hyperlink struct {
	Span
}

func (h *Hyperlink) HtmlBeginTag() string {
	return "<a href=\"#\">"
}

func (h *Hyperlink) HtmlEndTag() string {
	return "</a>"
}
