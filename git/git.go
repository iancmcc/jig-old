package git

import (
	"io"
	"os/exec"

	"github.com/iancmcc/jig/jig"
)

var git Git

type Git struct {
}

func (g *Git) RunWithProgress(args ...string) (progress <-chan *Progress, err error) {
	var stderr io.Reader
	// Ensure --progress is passed
	gcmd := args[0]
	args = append([]string{gcmd, "--progress"}, args[1:]...)
	// Run the command
	cmd := exec.Command("git", args...)
	if stderr, err = cmd.StderrPipe(); err != nil {
		return nil, err
	}
	parser := NewProgressParser(stderr)
	go cmd.Run()
	return parser.Parse(), nil
}

func (g *Git) Clone(uri, target string) (<-chan *Progress, error) {
	return g.RunWithProgress("clone", uri, target)
}

// Validate that GitRepository satisfies the Repository interface
var _ jig.Repository = &GitRepository{}

type GitRepository struct {
	jig.BaseRepository
}

func (r *GitRepository) Clone() error {
	uri := *r.GetURI()
	path, err := r.GetPath()
	if err != nil {
		return err
	}
	c, err := git.Clone(uri.String(), path)
	if err != nil {
		return err
	}
	for _ = range c {
	}
	return nil
}