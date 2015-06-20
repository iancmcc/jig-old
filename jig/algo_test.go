package jig_test

import (
	. "github.com/iancmcc/jig/jig"
	. "gopkg.in/check.v1"
)

func (s *JigSuite) TestFuzzyMatch(c *C) {
	str := []rune("this is a test string")
	sidx, _ := FuzzyMatch(true, &str, []rune("hsatststr"))
	c.Assert(sidx, Equals, 1)
}

func (s *JigSuite) TestMatchEmptyPattern(c *C) {
	str := []rune("this is a test string")
	sidx, _ := FuzzyMatch(true, &str, []rune{})
	c.Assert(sidx, Equals, 0)
}

func (s *JigSuite) TestCaseInsensitive(c *C) {
	str := []rune("THIS IS A TEST STRING")
	sidx, _ := FuzzyMatch(false, &str, []rune("test"))
	c.Assert(sidx, Equals, 10)
}

func (s *JigSuite) TestUnicode(c *C) {
	str := []rune("\u008F\u00F2")
	sidx, _ := FuzzyMatch(false, &str, []rune("\u00F2"))
	c.Assert(sidx, Equals, 1)
}

func (s *JigSuite) TestNoMatch(c *C) {
	str := []rune("abcdefg")
	sidx, _ := FuzzyMatch(false, &str, []rune("hijklmnop"))
	c.Assert(sidx, Equals, -1)
}
