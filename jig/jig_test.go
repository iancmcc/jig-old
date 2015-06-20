package jig_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/iancmcc/jig/jig"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type JigSuite struct {
}

var _ = Suite(&JigSuite{})

func (s *JigSuite) TestCreation(c *C) {
	path := c.MkDir()
	j, err := jig.NewJig(path)
	c.Assert(err, IsNil)
	c.Assert(j.Path(), Equals, path)
	c.Assert(j.IsRoot(), Equals, false)

	j.Initialize()

	c.Assert(j.IsRoot(), Equals, true)
	_, e := os.Stat(filepath.Join(path, ".jig"))
	c.Assert(e, IsNil)
}

func (s *JigSuite) TestFindRoot(c *C) {
	path := c.MkDir()
	subpath := filepath.Join(path, "a/b/c")
	os.MkdirAll(filepath.Join(subpath, "d"), os.ModeDir|0777)

	var (
		j   *jig.Jig
		err error
	)

	// Initialize jig at root and a subpath
	j, _ = jig.NewJig(path)
	j.Initialize()
	c.Assert(j.IsRoot(), Equals, true)
	j, _ = jig.NewJig(subpath)
	j.Initialize()
	c.Assert(j.IsRoot(), Equals, true)

	j, err = jig.FindClosestJig(path)
	c.Assert(err, IsNil)
	c.Assert(j.Path(), Equals, path)

	j, err = jig.FindClosestJig(filepath.Join(path, "a"))
	c.Assert(err, IsNil)
	c.Assert(j.Path(), Equals, path)

	j, err = jig.FindClosestJig(filepath.Join(path, "a/b"))
	c.Assert(err, IsNil)
	c.Assert(j.Path(), Equals, path)

	j, err = jig.FindClosestJig(filepath.Join(path, "a/b/c"))
	c.Assert(err, IsNil)
	c.Assert(j.Path(), Equals, subpath)

	j, err = jig.FindClosestJig(filepath.Join(path, "a/b/c/d"))
	c.Assert(err, IsNil)
	c.Assert(j.Path(), Equals, subpath)

}
