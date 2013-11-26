package proxy

import (
	"testing"
	"time"

	"launchpad.net/gocheck"

	"github.com/SoCloz/goprismic"
)

func TestProxy(t *testing.T) { gocheck.TestingT(t) }

type ProxyTestSuite struct {
	proxy *Proxy
	docs  []goprismic.Document
}

var _ = gocheck.Suite(&ProxyTestSuite{})

func (s *ProxyTestSuite) SetUpSuite(c *gocheck.C) {
	p, err := New("https://lesbonneschoses.prismic.io/api", "", 1, 10*time.Second, 5*time.Second)
	if err == nil {
		s.proxy = p
		s.docs, _ = s.proxy.Direct().Master().Form("everything").Submit()
	}
}

func (s *ProxyTestSuite) SetUpTest(c *gocheck.C) {
	s.proxy.Clear()
}

func (s *ProxyTestSuite) TestProxy(c *gocheck.C) {
	c.Assert(s.proxy, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
}

func (s *ProxyTestSuite) TestLoadReload(c *gocheck.C) {
	c.Assert(s.proxy, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	d := s.docs[0]
	d1, err := s.proxy.GetDocument(d.Id)
	c.Assert(err, gocheck.IsNil, gocheck.Commentf("Submit did not return an error - %s", err))
	c.Assert(d, gocheck.DeepEquals, *d1, gocheck.Commentf("Proxy returns the same doc"))
	d1, err = s.proxy.GetDocument(d.Id)
	c.Assert(err, gocheck.IsNil, gocheck.Commentf("Submit did not return an error - %s", err))
	c.Assert(d, gocheck.DeepEquals, *d1, gocheck.Commentf("Proxy returns the same doc again"))
}

func (s *ProxyTestSuite) TestGetBy(c *gocheck.C) {
	c.Assert(s.proxy, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	d, err := s.proxy.GetDocumentBy("product", "flavour", "Pistachio")
	c.Assert(d, gocheck.Not(gocheck.IsNil), gocheck.Commentf("Submit did not return an error - %s", err))
	f, _ := d.GetTextFragment("flavour")
	c.Assert(f.AsText(), gocheck.Equals, "Pistachio", gocheck.Commentf("Proxy returns the same doc"))
}

func (s *ProxyTestSuite) TestTtlAndGrace(c *gocheck.C) {
	c.Assert(s.proxy, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	d := s.docs[0]
	s.proxy.GetDocument(d.Id)
	c.Assert(s.proxy.GetStats(), gocheck.DeepEquals, Stats{Miss: 1}, gocheck.Commentf("miss"))
	s.proxy.GetDocument(d.Id)
	c.Assert(s.proxy.GetStats(), gocheck.DeepEquals, Stats{Hit: 1, Miss: 1}, gocheck.Commentf("miss+hit"))

	time.Sleep(5 * time.Second)
	s.proxy.GetDocument(d.Id)
	c.Assert(s.proxy.GetStats(), gocheck.DeepEquals, Stats{Hit: 2, Miss: 1, Refresh: 1}, gocheck.Commentf("miss+hit+refresh"))

	time.Sleep(12 * time.Second) // Add some delay to allow for async refresh to finish
	s.proxy.GetDocument(d.Id)
	c.Assert(s.proxy.GetStats(), gocheck.DeepEquals, Stats{Hit: 2, Miss: 2, Refresh: 1}, gocheck.Commentf("miss+hit+refresh+miss"))
}
