package git_test

import (
	. "github.com/iancmcc/jig/git"
	"github.com/iancmcc/jig/jig"
	. "gopkg.in/check.v1"
)

type UriSuite struct{}

var _ = Suite(&UriSuite{})

func assertURIParts(uri jig.RepositoryURI, c *C) {
	c.Assert(uri, Not(IsNil))
	c.Assert(uri.Domain(), Equals, "github.com")
	c.Assert(uri.Owner(), Equals, "iancmcc")
	c.Assert(uri.Repository(), Equals, "jig")
}

func (s *UriSuite) TestParseSSHUri(c *C) {
	var (
		uri jig.RepositoryURI
		err error
	)
	uri, err = ParseGitURI("git@github.com:iancmcc/jig")
	c.Assert(err, IsNil)
	assertURIParts(uri, c)
	c.Assert(string(uri.Scheme()), Equals, jig.SCHEME_SSH)
	gituri := uri.(*GitURI)
	c.Assert(gituri.ToSSH(), Equals, "git@github.com:iancmcc/jig")
	c.Assert(gituri.ToHTTPS(), Equals, "https://github.com/iancmcc/jig")
}

func (s *UriSuite) TestParseNoRepoUri(c *C) {
	var (
		uri jig.RepositoryURI
		err error
	)
	uri, err = ParseGitURI("iancmcc/jig")
	c.Assert(err, IsNil)
	assertURIParts(uri, c)
	c.Assert(uri.Scheme(), Equals, jig.SCHEME_HTTPS)
}

func (s *UriSuite) TestParseHTTPSUri(c *C) {
	var (
		uri jig.RepositoryURI
		err error
	)
	uri, err = ParseGitURI("https://github.com/iancmcc/jig")
	c.Assert(err, IsNil)
	assertURIParts(uri, c)
	c.Assert(uri.Scheme(), Equals, jig.SCHEME_HTTPS)
}

func (s *UriSuite) TestParseGitUri(c *C) {
	var (
		uri jig.RepositoryURI
		err error
	)
	uri, err = ParseGitURI("git://github.com/iancmcc/jig")
	c.Assert(err, IsNil)
	assertURIParts(uri, c)
	c.Assert(string(uri.Scheme()), Equals, jig.SCHEME_GIT)
}

func (s *UriSuite) TestParseGithubStyle(c *C) {
	var (
		uri jig.RepositoryURI
		err error
	)
	uri, err = ParseGitURI("github.com/iancmcc/jig")
	c.Assert(err, IsNil)
	assertURIParts(uri, c)
	// Default to HTTPS
	c.Assert(uri.Scheme(), Equals, jig.SCHEME_HTTPS)
}
