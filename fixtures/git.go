package fixtures

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
)

var pwd, home string

func init() {
	// caching `pwd` and $HOME and reset them after test repo is teared down
	// `pwd` is changed to the bin temp dir during test run
	pwd, _ = os.Getwd()
	home = os.Getenv("HOME")
}

type TestRepo struct {
	pwd    string
	dir    string
	home   string
	Remote string
}

func (r *TestRepo) Setup() {
	dir, err := ioutil.TempDir("", "test-repo")
	if err != nil {
		panic(err)
	}
	r.dir = dir

	os.Setenv("HOME", r.dir)

	targetPath := filepath.Join(r.dir, "test.git")
	err = r.clone(r.Remote, targetPath)
	if err != nil {
		panic(err)
	}

	err = os.Chdir(targetPath)
	if err != nil {
		panic(err)
	}
}

func (r *TestRepo) clone(repo, dir string) error {
	output, err := exec.Command("git", "clone", repo, dir).CombinedOutput()
	if err != nil {
		err = fmt.Errorf("error cloning %s to %s: %s", repo, dir, output)
	}
	return err
}

func (r *TestRepo) TearDown() {
	err := os.Chdir(r.pwd)
	if err != nil {
		panic(err)
	}

	os.Setenv("HOME", r.home)

	err = os.RemoveAll(r.dir)
	if err != nil {
		panic(err)
	}

}

func SetupTestRepo() *TestRepo {
	remotePath := filepath.Join(pwd, "..", "fixtures", "test.git")
	repo := &TestRepo{pwd: pwd, home: home, Remote: remotePath}
	repo.Setup()
	return repo
}

func CloneOutput() []byte {
	remotePath := filepath.Join(pwd, "..", "fixtures", "clone_output.txt")
	s, _ := ioutil.ReadFile(remotePath)
	return s
}
