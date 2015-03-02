package git

import (
	"io"
	"os/exec"
	"regexp"
)

var progre = regexp.MustCompile(`(?P<Message>.*?):.*\((?P<Current>\d+)/(?P<Total>\d+)\)`)

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