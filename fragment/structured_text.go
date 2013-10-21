package fragment

import (
	"fmt"

	"github.com/SoCloz/goprismic/fragment/block"
)

type StructuredText []block.Block

func (st *StructuredText) Decode(enc interface{}) error {
	dec, ok := enc.([]interface{})
	if !ok {
		return fmt.Errorf("%+v is not a slice", enc)
	}
	*st = make(StructuredText, 0, len(dec))
	for _, v := range dec {
		dec, ok := v.(map[string]interface{})
		if !ok {
			return fmt.Errorf("%+v is not a map", v)
		}
		var b block.Block
		switch dec["type"] {
		case "heading1", "heading2", "heading3", "heading4":
			b = new(block.Heading)
		case "paragraph":
			b = new(block.Paragraph)
		case "preformatted":
			b = new(block.Preformatted)
		case "list-item":
			b = new(block.ListItem)
		case "image":
			b = new(block.OrderedListItem)
		case "o-list-item":
			b = new(block.Image)
		case "embed":
			b = new(block.Embed)
		default:
			panic(fmt.Sprintf("Unknown block type %s", dec["type"]))
		}
		err := b.Decode(v)
		if err != nil {
			//fmt.Printf("\n!!! %s\n", err)
			return err
		}
		//fmt.Printf("\n%+v => %+v\n", dec, b)
		*st = append(*st, b)
	}
	return nil
}

// Formats the fragment content as html
func (st StructuredText) AsHtml() string {
	parentTag := ""
	html := ""
	for _, v := range st {
		if parentTag != v.ParentHtmlTag() {
			if parentTag != "" {
				html += fmt.Sprintf("</%s>", parentTag)
			}
			parentTag = v.ParentHtmlTag()
			if parentTag != "" {
				html += fmt.Sprintf("<%s>", parentTag)
			}
		}
		html += v.AsHtml()
	}
	if parentTag != "" {
		html += fmt.Sprintf("</%s>", parentTag)
	}
	return html
}
