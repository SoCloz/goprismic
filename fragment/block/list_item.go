package block

import(
	"fmt"
)

type ListItem struct {
	BaseBlock
}

func (l *ListItem) AsHtml() string {
	return fmt.Sprintf("<li>%s</li>", l.FormatHtmlText());
}

func (l *ListItem) ParentHtmlTag() string {
	return "ul"
}

type OrderedListItem struct {
	BaseBlock
}

func (l *OrderedListItem) AsHtml() string {
	return fmt.Sprintf("<li>%s</li>", l.FormatHtmlText());
}

func (l *OrderedListItem) ParentHtmlTag() string {
	return "ol"
}
