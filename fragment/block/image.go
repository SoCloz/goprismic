package block

import(
	"fmt"
)

type Image struct {
	BaseBlock
}

func (i *Image) AsHtml() string {
	return fmt.Sprintf("<img>%s</ig>", i.Text);
}

func (i *Image) ParentHtmlTag() string {
	return ""
}