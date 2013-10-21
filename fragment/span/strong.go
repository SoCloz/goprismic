package span

type Strong struct {
	Span
}

func (s *Strong) HtmlBeginTag() string {
	return "<strong>"
}

func (s *Strong) HtmlEndTag() string {
	return "</strong>"
}
