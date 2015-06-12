package repository_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/iancmcc/jig/fixtures"
	. "github.com/iancmcc/jig/git"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type GitSuite struct {
	dir  string
	repo *fixtures.TestRepo
}

var _ = Suite(&GitSuite{})

func (s *GitSuite) SetUpSuite(c *C) {
	s.repo = fixtures.SetupTestRepo()
}

func (s *GitSuite) TearDownSuite(c *C) {
	s.repo.TearDown()
}

func (s *GitSuite) SetUpTest(c *C) {
	s.dir = c.MkDir()
}

type fileExists struct {
	info *CheckerInfo
}

type isGitRepo struct {
	info *CheckerInfo
}

func (f *fileExists) Info() *CheckerInfo {
	return f.info
}
func (f *isGitRepo) Info() *CheckerInfo {
	return f.info
}

func (f *fileExists) Check(params []interface{}, names []string) (result bool, error string) {
	if _, err := os.Stat(params[0].(string)); os.IsNotExist(err) {
		return false, "File does not exist"
	}
	return true, ""
}

func (f *isGitRepo) Check(params []interface{}, names []string) (result bool, error string) {
	gitpath := filepath.Join(params[0].(string), ".git")
	return FileExists.Check([]interface{}{gitpath}, names)
}

var FileExists Checker = &fileExists{
	&CheckerInfo{Name: "FileExists", Params: []string{"path"}},
}

var IsGitRepo Checker = &isGitRepo{
	&CheckerInfo{Name: "IsGitRepo", Params: []string{"path"}},
}

func (s *GitSuite) TestCloneProgress(c *C) {
	git := &Git{}
	progress, err := git.Clone("file://"+s.repo.Remote, s.dir)
	progresses := []*Progress{}
	for p := range progress {
		progresses = append(progresses, p)
	}
	// TODO: A non-sucky test
	c.Assert(err, IsNil)
	c.Assert(s.dir, IsGitRepo)
}
