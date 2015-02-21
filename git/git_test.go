package git_test

import (
	"fmt"
	"testing"

	. "github.com/iancmcc/jig/git"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GitSuite struct{}

var _ = Suite(&GitSuite{})

func (s *GitSuite) TestParseProgress(c *C) {
	p, err := ParseProgressLine("Receiving objects:   0% (1/2026)   ")
	fmt.Printf("%+v\n", err)
	c.Assert(err, Equals, nil)
	c.Assert(p.Current, Equals, 1)
	c.Assert(p.Total, Equals, 2026)
	c.Assert(p.Message, Equals, "Receiving objects")
	p2, err := ParseProgressLine("Resolving deltas:  18% (34/187)   ")
	c.Assert(err, Equals, nil)
	c.Assert(p2.Current, Equals, 34)
	c.Assert(p2.Total, Equals, 187)
	c.Assert(p2.Message, Equals, "Resolving deltas")
}
