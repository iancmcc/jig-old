package jig_test

import (
	"github.com/iancmcc/jig/jig"
	. "gopkg.in/check.v1"
)

func (s *JigSuite) TestMatcher(c *C) {
	paths := []string{
		"github.com/powerline/fonts",
		"github.com/fonts/powerline",
		"github.com/abc/fonts",
		"github.com/fonts/fonts",
		"github.com/powerline/powerline",
		"powerline-fonts",
	}

	c.Assert(jig.SortedMatches("powerline", paths)[0], Equals, "github.com/powerline/powerline")
	c.Assert(jig.SortedMatches("powerline/fonts", paths)[0], Equals, "github.com/powerline/fonts")
	c.Assert(jig.SortedMatches("fonts/powerline", paths)[0], Equals, "github.com/fonts/powerline")
	c.Assert(jig.SortedMatches("p/fonts", paths)[0], Equals, "github.com/powerline/fonts")
	c.Assert(jig.SortedMatches("fonts/p", paths)[0], Equals, "github.com/fonts/powerline")
	c.Assert(jig.SortedMatches("fonts", paths)[0], Equals, "github.com/fonts/fonts")
	c.Assert(jig.SortedMatches("abc", paths)[0], Equals, "github.com/abc/fonts")
}
