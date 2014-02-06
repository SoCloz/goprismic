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
		sr, _ := s.proxy.Direct().Master().Form("everything").Submit()
		s.docs = sr.Results
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

func (s *ProxyTestSuite) TestLRU(c *gocheck.C) {
	c.Assert(s.proxy, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	sr, err := s.proxy.Search("product", "")
	c.Assert(sr, gocheck.Not(gocheck.IsNil), gocheck.Commentf("Submit did not return an error - %s", err))
	s.proxy.Clear()
	stats := s.proxy.GetStats()
	// accessing ds[0]
	s.proxy.GetDocument(sr.Results[0].Id)
	stats1 := s.proxy.GetStats()
	c.Assert(stats1.Hit, gocheck.Equals, stats.Hit, gocheck.Commentf("on an empty cache, no hit"))
	c.Assert(stats1.Miss, gocheck.Equals, stats.Miss+1, gocheck.Commentf("on an empty cache, miss"))
	// accessing ds[0] again
	s.proxy.GetDocument(sr.Results[0].Id)
	stats1 = s.proxy.GetStats()
	c.Assert(stats1.Hit, gocheck.Equals, stats.Hit+1, gocheck.Commentf("doc in cache, hit"))
	c.Assert(stats1.Miss, gocheck.Equals, stats.Miss+1, gocheck.Commentf("doc in cache, no miss"))
	// accessing another doc. LRU size is 1 => ds[0] should be evicted
	s.proxy.GetDocument(sr.Results[1].Id)
	stats1 = s.proxy.GetStats()
	c.Assert(stats1.Hit, gocheck.Equals, stats.Hit+1, gocheck.Commentf("doc not in cache, no hit"))
	c.Assert(stats1.Miss, gocheck.Equals, stats.Miss+2, gocheck.Commentf("doc not in cache, miss"))
	// accessing ds[0] again, should not be found in cache
	s.proxy.GetDocument(sr.Results[0].Id)
	stats1 = s.proxy.GetStats()
	c.Assert(stats1.Hit, gocheck.Equals, stats.Hit+1, gocheck.Commentf("doc evicted from cache, no hit"))
	c.Assert(stats1.Miss, gocheck.Equals, stats.Miss+3, gocheck.Commentf("doc evicted from cache, miss"))
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
