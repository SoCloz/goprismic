package goprismic

import (
	"testing"

	"launchpad.net/gocheck"
)

func TestApi(t *testing.T) { gocheck.TestingT(t) }

type ApiTestSuite struct {
	api *Api
}

var _ = gocheck.Suite(&ApiTestSuite{})

func (s *ApiTestSuite) SetUpSuite(c *gocheck.C) {
	s.api, _ = Get("https://lesbonneschoses.prismic.io/api", "", DefaultConfig)
}

func (s *ApiTestSuite) TestApi(c *gocheck.C) {
	c.Assert(s.api, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	c.Assert(len(s.api.Data.Bookmarks), gocheck.Equals, 3, gocheck.Commentf("3 bookmarks found"))
	c.Assert(len(s.api.Data.Forms), gocheck.Equals, 10, gocheck.Commentf("10 forms found"))
}

func (s *ApiTestSuite) TestSearch(c *gocheck.C) {
	c.Assert(s.api, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	sr, err := s.api.Master().Form("everything").PageSize(15).Submit()
	c.Assert(err, gocheck.IsNil, gocheck.Commentf("Submit did not return an error - %s", err))
	c.Assert(len(sr.Results), gocheck.Equals, 15, gocheck.Commentf("Submit did return 20 documents"))
	c.Assert(sr.ResultsPerPage, gocheck.Equals, 15, gocheck.Commentf("Submit did return the right page size"))

	sr, err = s.api.Master().Form("everything").PageSize(15).Page(2).Submit()
	c.Assert(err, gocheck.IsNil, gocheck.Commentf("Submit did not return an error - %s", err))
	c.Assert(sr.Page, gocheck.Equals, 2, gocheck.Commentf("Submit did return the second page"))
}

func (s *ApiTestSuite) TestQuery(c *gocheck.C) {
	c.Assert(s.api, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	sr, err := s.api.Master().Form("everything").Query("[[:d = at(document.tags, [\"Macaron\"])]]").Submit()
	c.Assert(err, gocheck.IsNil, gocheck.Commentf("Submit did not return an error - %s", err))
	c.Assert(len(sr.Results), gocheck.Equals, 7, gocheck.Commentf("Submit did return 7 documents"))

}

func (s *ApiTestSuite) TestError(c *gocheck.C) {
	c.Assert(s.api, gocheck.NotNil, gocheck.Commentf("Connection with api is OK"))
	_, err := s.api.Master().Form("everything").Query("[:d = at(document.tags, [\"Macaron\"])]").Submit()
	c.Assert(err, gocheck.NotNil, gocheck.Commentf("Submit should return an error"))
	c.Assert(err.(*PrismicError).Type, gocheck.Equals, "parsing-error", gocheck.Commentf("Submit should return a parse error"))
}
