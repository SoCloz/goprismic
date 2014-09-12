package goprismic

import (
	"fmt"
	"testing"
	"time"

	"launchpad.net/gocheck"

	"github.com/SoCloz/goprismic/fragment/block"
	"github.com/SoCloz/goprismic/fragment/link"
	"github.com/SoCloz/goprismic/test"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type TestSuite struct {
	doc *Document
}

var _ = gocheck.Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *gocheck.C) {
	sr := SearchResult{}
	sr.Results = []Document{}
	test.Load("search", &sr)
	if len(sr.Results) == 0 {
		panic("no doc found")
	}
	s.doc = &sr.Results[0]
}

func (s *TestSuite) TestSlug(c *gocheck.C) {
	c.Assert(s.doc.GetSlug(), gocheck.Equals, "cool-coconut-macaron", gocheck.Commentf("Has the right slug"))
	c.Assert(s.doc.HasSlug("coconut-macaron"), gocheck.Equals, true, gocheck.Commentf("Has the other slug"))
}

func (s *TestSuite) TestGetText(c *gocheck.C) {
	t, _ := s.doc.GetText("name")
	c.Assert(t, gocheck.Equals, "Cool Coconut Macaron", gocheck.Commentf("Found the right text"))
}

func (s *TestSuite) TestGetColor(c *gocheck.C) {
	t, _ := s.doc.GetColor("color")
	c.Assert(t, gocheck.Equals, "#d4bf43", gocheck.Commentf("Found the right text"))
}

func (s *TestSuite) TestGetNumber(c *gocheck.C) {
	t, _ := s.doc.GetNumber("price")
	c.Assert(t, gocheck.Equals, 2.5, gocheck.Commentf("Found the right number"))
}

func (s *TestSuite) TestGetBoolean(c *gocheck.C) {
	a, _ := s.doc.GetBool("adult")
	c.Assert(a, gocheck.Equals, true, gocheck.Commentf("Found the right bool"))
	a, _ = s.doc.GetBool("teenager")
	c.Assert(a, gocheck.Equals, true, gocheck.Commentf("Found the right bool"))
	a, _ = s.doc.GetBool("french")
	c.Assert(a, gocheck.Equals, false, gocheck.Commentf("Found the right bool"))
}

func (s *TestSuite) TestGetDate(c *gocheck.C) {
	df, _ := s.doc.GetDateFragment("birthdate")
	c.Assert(df.AsText(), gocheck.Equals, "2013-10-23", gocheck.Commentf("Found the right date"))
	t, _ := s.doc.GetDate("birthdate")
	c.Assert(t.Month(), gocheck.Equals, time.October, gocheck.Commentf("Date has the right month"))
}

func (s *TestSuite) TestGetImageView(c *gocheck.C) {
	i, _ := s.doc.GetImageFragment("image")
	c.Assert(i, gocheck.NotNil, gocheck.Commentf("Found an image"))
	c.Assert(i.Main.Url, gocheck.Equals, "https://prismicio.s3.amazonaws.com/lesbonneschoses/30214ac0c3a51e7516d13c929086c49f49af7988.png", gocheck.Commentf("Image has the right url"))
	c.Assert(i.AsHtml(), gocheck.Equals, i.Main.AsHtml(), gocheck.Commentf("The image html is the main view html"))
	c.Assert(i.AsHtml(), gocheck.Equals, "<img src=\"https://prismicio.s3.amazonaws.com/lesbonneschoses/30214ac0c3a51e7516d13c929086c49f49af7988.png\" width=\"500\" height=\"500\"/>", gocheck.Commentf("The image html is ok"))

	v, _ := i.GetView("icon")
	c.Assert(v, gocheck.NotNil, gocheck.Commentf("Found a view"))
	c.Assert(v.Url, gocheck.Equals, "https://prismicio.s3.amazonaws.com/lesbonneschoses/f428de9ed832705617063c9c69eb752f4ab92ac5.png", gocheck.Commentf("View has the right url"))
	c.Assert(v.AsHtml(), gocheck.Equals, "<img src=\"https://prismicio.s3.amazonaws.com/lesbonneschoses/f428de9ed832705617063c9c69eb752f4ab92ac5.png\" width=\"250\" height=\"250\"/>", gocheck.Commentf("The view html is ok"))
	c.Assert(v.Ratio(), gocheck.Equals, 1.0, gocheck.Commentf("The view has the right ratio"))
}

func (s *TestSuite) TestHeader(c *gocheck.C) {
	t, _ := s.doc.GetFragment("name")
	c.Assert(t, gocheck.NotNil, gocheck.Commentf("Has a first header"))
	c.Assert(t.AsHtml(), gocheck.Equals, "<h1>Cool Coconut Macaron</h1>", gocheck.Commentf("Has the right content"))
}

func (s *TestSuite) TestGetFirstParagraph(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("description")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "If you ever met coconut taste on its bad day, you surely know that coconut, coming from bad-tempered islands, can be rough sometimes. That is why we like to soften it with a touch of caramel taste in its ganache. The result is the perfect encounter between the finest palm fruit and the most tasty of sugarcane's offspring."
	c.Assert(p.Text, gocheck.Equals, content, gocheck.Commentf("Has the right content"))
}

func (s *TestSuite) TestGetFirstPreformatted(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("description")
	p, _ := text.GetFirstPreformatted()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first preformatted"))
	content := "If you ever met coconut taste on its bad day, you surely know that coconut, coming from bad-tempered islands, can be rough sometimes. That is why we like to soften it with a touch of caramel taste in its ganache. The result is the perfect encounter between the finest palm fruit and the most tasty of sugarcane's offspring."
	c.Assert(p.Text, gocheck.Equals, content, gocheck.Commentf("Has the right content"))
}

func (s *TestSuite) TestParagraphRendering(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("render")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "<p>This <em>is</em> <strong>a</strong> test.</p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestPreformattedRendering(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("render")
	p, _ := text.GetFirstPreformatted()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first preformatted"))
	content := "<pre>This <em>is</em> <strong>a</strong> test.</pre>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestParagraphEscape(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("escape")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "<p>This is &lt;a&gt; test.</p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestPreformattedEscape(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("escape")
	p, _ := text.GetFirstPreformatted()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first preformatted"))
	content := "<pre>This is &lt;a&gt; test.</pre>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestParagraphEscapeAndSpan(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("escapespan")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "<p>Thi&amp; i&amp; &lt;a&gt; <em>test</em>.</p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestParagraphUtf8Span(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("utf8span")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "<p>Thìs ìs &lt;à&gt; <em>tést</em>.</p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestLinkFragment(c *gocheck.C) {
	l, _ := s.doc.GetLinkFragment("related")
	r := func(l link.Link) string {
		return l.(*link.DocumentLink).Document.Slug
	}
	l.ResolveLinks(r)
	content := "<a href=\"vanilla-macaron\">vanilla-macaron</a>"
	c.Assert(l.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestLinkSpans(c *gocheck.C) {
	desc, _ := s.doc.GetStructuredTextFragment("description")
	r := func(l link.Link) string {
		return l.(*link.DocumentLink).Document.Slug
	}
	desc.ResolveLinks(r)

	blocks := *desc

	p := blocks[2].(*block.Paragraph)
	content := "<p><a href=\"http://apple.com\">link</a></p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Link.web has the right rendering"))

	p = blocks[3].(*block.Paragraph)
	content = "<p><a href=\"apricot-pie\">link1</a></p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Link.document has the right rendering"))

	p = blocks[4].(*block.Paragraph)
	content = "<p><a href=\"http://data.prismic.io/lesbonneschoses%2F1378998378075_medium_1374778510922_coconut.png\">link2</a></p>"
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Link.file has the right rendering"))
}

func (s *TestSuite) TestEmbedAndImage(c *gocheck.C) {
	blocks, _ := s.doc.GetStructuredTextBlocks("embedimage")

	i := blocks[0]
	content := "<img src=\"https://wroomdev.s3.amazonaws.com/lesbonneschoses/899162db70c73f11b227932b95ce862c63b9df2A.jpg\" width=\"800\" height=\"400\"/>"
	c.Assert(i.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Image block has the right rendering"))

	e := blocks[1]
	content = "<div data-oembed=\"http://www.youtube.com/watch?v=MmIKLlRE7n0\" data-oembed-type=\"video\" data-oembed-provider=\"YouTube\"><iframe width='459' height='344' src='http://www.youtube.com/embed/MmIKLlRE7n0?feature=oembed' frameborder='0' allowfullscreen></iframe></div>"
	c.Assert(e.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Embed block has the right rendering"))
}

func (s *TestSuite) TestList(c *gocheck.C) {
	desc, _ := s.doc.GetStructuredTextFragment("list")

	content := "<ol><li>ol1</li><li>ol2</li><li>ol3</li></ol><ul><li>l1</li><li>l2</li><li>l3</li></ul>"
	c.Assert(desc.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Lists have the right rendering"))
}

func (s *TestSuite) TestGeoPoint(c *gocheck.C) {
	geo, _ := s.doc.GetGeoPointFragment("location")

	c.Assert(geo.Latitude, gocheck.Equals, 48.87687670000001, gocheck.Commentf("Geopoint has the correct latitude"))
	c.Assert(geo.Longitude, gocheck.Equals, 2.3338801999999825, gocheck.Commentf("Geopoint has the correct longitude"))

	content := fmt.Sprintf(`<div class="geopoint"><span class="latitude">%f</span><span class="longitude">%f</span></div>`, geo.Latitude, geo.Longitude)
	c.Assert(geo.AsHtml(), gocheck.Equals, content, gocheck.Commentf("GeoPoints have the right rendering"))
}
