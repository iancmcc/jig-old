package vcs_test

import (
	"io/ioutil"
	"runtime"
	"testing"

	. "github.com/iancmcc/jig/vcs"
	"github.com/libgit2/git2go"

	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

func checkFatal(t *testing.T, err error) {
	if err == nil {
		return
	}

	// The failure happens at wherever we were called, not here
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		t.Fatal()
	}

	t.Fatalf("Fail at %v:%v; %v", file, line, err)
}

func createTestRepo(t *testing.T) *git.Repository {
	// figure out where we can create the test repo
	path, err := ioutil.TempDir("", "git2go")
	checkFatal(t, err)
	repo, err := git.InitRepository(path, false)
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
