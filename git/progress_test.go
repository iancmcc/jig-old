package git_test

import (
	"bytes"

	"github.com/iancmcc/jig/fixtures"
	. "github.com/iancmcc/jig/git"
	. "gopkg.in/check.v1"
)

type ProgressSuite struct{}

var _ = Suite(&ProgressSuite{})

func (s *ProgressSuite) TestParseProgress(c *C) {
	parser := &ProgressParser{}
	p, err := parser.ParseLine("Receiving objects:   0% (1/2026)   ")
	c.Assert(err, Equals, nil)
	c.Assert(p.Current, Equals, 1)
	c.Assert(p.Total, Equals, 2026)
	c.Assert(p.CurrentStage(), Equals, RECEIVING)
	c.Assert(p.IsBegin(), Equals, true)
	p2, err := parser.ParseLine("Receiving objects:   1% (10/2026)   ")
	c.Assert(err, Equals, nil)
	c.Assert(p2.Current, Equals, 10)
	c.Assert(p2.Total, Equals, 2026)
	c.Assert(p2.CurrentStage(), Equals, RECEIVING)
	c.Assert(p2.IsBegin(), Equals, false)
	p3, err := parser.ParseLine("Receiving objects:   100% (2026/2026), done. ")
	c.Assert(err, Equals, nil)
	c.Assert(p3.Current, Equals, 2026)
	c.Assert(p3.Total, Equals, 2026)
	c.Assert(p3.CurrentStage(), Equals, RECEIVING)
	c.Assert(p3.IsBegin(), Equals, false)
	c.Assert(p3.IsEnd(), Equals, true)
	p4, err := parser.ParseLine("Resolving deltas:  18% (34/187)   ")
	c.Assert(err, Equals, nil)
	c.Assert(p4.Current, Equals, 34)
	c.Assert(p4.Total, Equals, 187)
	c.Assert(p4.IsBegin(), Equals, true)
	c.Assert(p4.CurrentStage(), Equals, RESOLVING)
}

func (s *ProgressSuite) TestParser(c *C) {
	clone_output := fixtures.CloneOutput()
	parser := NewProgressParser(bytes.NewReader(clone_output))
	outchan := parser.Parse()
	for p := range outchan {
		if p.IsBegin() {
			c.Assert(p.CurrentStage(), Equals, COUNTING)
			c.Assert(p.IsEnd(), Equals, true)
			break
		}
	}
	var count int
	for p := range outchan {
		count += 1
		if p.IsBegin() {
			c.Assert(p.Current, Equals, 1)
			c.Assert(p.IsEnd(), Equals, false)
		}
		if p.IsEnd() {
			c.Assert(p.Current, Equals, 2036)
			break
		}
		c.Assert(p.CurrentStage(), Equals, RECEIVING)
		c.Assert(p.Total, Equals, 2036)
	}
	c.Assert(count, Equals, 102)
	count = 0
	for p := range outchan {
		count += 1
		if p.IsBegin() {
			c.Assert(p.Current, Equals, 0)
			c.Assert(p.IsEnd(), Equals, false)
		}
		if p.IsEnd() {
			c.Assert(p.Current, Equals, 193)
			break
		}
		c.Assert(p.CurrentStage(), Equals, RESOLVING)
		c.Assert(p.Total, Equals, 193)
	}
	c.Assert(count, Equals, 69)
	for p := range outchan {
		stage := p.CurrentStage()
		c.Assert(stage, Not(Equals), COUNTING)
		c.Assert(stage, Not(Equals), RESOLVING)
		c.Assert(stage, Not(Equals), RECEIVING)
	}
}
