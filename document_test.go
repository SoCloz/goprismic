package goprismic

import(
	"testing"
	"time"

	"launchpad.net/gocheck"

	"github.com/SoCloz/goprismic/test"
)


func Test(t *testing.T) { gocheck.TestingT(t) }

type TestSuite struct {
	doc *Document
}

var _ = gocheck.Suite(&TestSuite{})

func (s *TestSuite) SetUpSuite(c *gocheck.C) {
	docs := []Document{}
	test.Load("search", &docs)
	if len(docs) == 0 {
		panic("no doc found")
	}
	s.doc = &docs[0]
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
	content := "If you ever met coconut taste on its bad day, you surely know that coconut, coming from bad-tempered islands, can be rough sometimes. That is why we like to soften it with a touch of caramel taste in its ganache. The result is the perfect encounter between the finest palm fruit and the most tasty of sugarcane's offspring.";
	c.Assert(p.Text, gocheck.Equals, content, gocheck.Commentf("Has the right content"))
}

func (s *TestSuite) TestGetFirstPreformatted(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("description")
	p, _ := text.GetFirstPreformatted()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first preformatted"))
	content := "If you ever met coconut taste on its bad day, you surely know that coconut, coming from bad-tempered islands, can be rough sometimes. That is why we like to soften it with a touch of caramel taste in its ganache. The result is the perfect encounter between the finest palm fruit and the most tasty of sugarcane's offspring.";
	c.Assert(p.Text, gocheck.Equals, content, gocheck.Commentf("Has the right content"))
}

func (s *TestSuite) TestParagraphRendering(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("render")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "<p>This <em>is</em> <strong>a</strong> test.</p>";
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestPreformattedRendering(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("render")
	p, _ := text.GetFirstPreformatted()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first preformatted"))
	content := "<pre>This <em>is</em> <strong>a</strong> test.</pre>";
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestParagraphEscape(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("escape")
	p, _ := text.GetFirstParagraph()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first paragraph"))
	content := "<p>This is &lt;a&gt; test.</p>";
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}

func (s *TestSuite) TestPreformattedEscape(c *gocheck.C) {
	text, _ := s.doc.GetStructuredTextFragment("escape")
	p, _ := text.GetFirstPreformatted()
	c.Assert(p, gocheck.NotNil, gocheck.Commentf("Has a first preformatted"))
	content := "<pre>This is &lt;a&gt; test.</pre>";
	c.Assert(p.AsHtml(), gocheck.Equals, content, gocheck.Commentf("Has the right rendering"))
}
