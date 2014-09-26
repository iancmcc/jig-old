package vcs_test

import (
	"io/ioutil"
	"testing"

	. "github.com/iancmcc/jig/vcs"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func createTestRepo(t *testing.T) *Repository {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "git2go")
	checkFatal(t, err)
	repo, err := InitRepository(path, false)
	checkFatal(t, err)

	tmpfile := "README"
	err = ioutil.WriteFile(path+"/"+tmpfile, []byte("foo\n"), 0644)

	checkFatal(t, err)

	return repo
}

var _ = Describe("Git", func() {

	var repo GitRepository

	BeforeEach(func() {
		repo = GitRepository{}
	})

})
