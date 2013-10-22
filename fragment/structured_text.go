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
			panic(fmt.Sprintf("unknown block type %s", dec["type"]))
		}
		err := b.Decode(v)
		if err != nil {
			return err
		}
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

// Formats the fragment content as text
func (st StructuredText) AsText() string {
	text := ""
	for _, v := range st {
		if text != "" {
			text += "\n"
		}
		text += v.AsText()
	}
	return text
}

// Returns the first paragraph fragment
func (st StructuredText) GetFirstParagraph() (*block.Paragraph, bool) {
	for k := range st {
		b, ok := st[k].(*block.Paragraph)
		if ok {
			return b, true
		}
	}
	return nil, false
}

// Returns the first preformatted fragment
func (st StructuredText) GetFirstPreformatted() (*block.Preformatted, bool) {
	for k := range st {
		b, ok := st[k].(*block.Preformatted)
		if ok {
			return b, true
		}
	}
	return nil, false
}

// Returns the first image fragment
func (st StructuredText) GetFirstImage() (*block.Image, bool) {
	for k := range st {
		b, ok := st[k].(*block.Image)
		if ok {
			return b, true
		}
	}
	return nil, false
}