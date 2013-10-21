package span

type Em struct {
	Span
}

func (e *Em) HtmlBeginTag() string {
	return "<em>"
}

func (e *Em) HtmlEndTag() string {
	return "</em>"
}
