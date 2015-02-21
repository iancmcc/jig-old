package git

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

type Progress struct {
	Message string
	Current int
	Total   int
}

var progre = regexp.MustCompile(`(?P<Message>.*?):.*\((?P<Current>\d+)/(?P<Total>\d+)\)`)

type Git struct {
}

func ParseProgressLine(line string) (*Progress, error) {
	p := &Progress{}
	stuff := progre.FindStringSubmatch(line)
	if len(stuff) != 4 {
		return nil, fmt.Errorf("Unable to parse progress line")
	}
	p.Message = stuff[1]
	if current, err := strconv.ParseInt(stuff[2], 0, 32); err != nil {
		return nil, err
	} else {
		p.Current = int(current)
	}
	if total, err := strconv.ParseInt(stuff[3], 0, 32); err != nil {
		return nil, err
	} else {
		p.Total = int(total)
	}
	return p, nil
}

func (g *Git) Run(args []string) error {
	cmd := exec.Command("git", args)
}

func (g *Git) RunWithProgress(args []string) (chan<- Progress, error) {
	cmd := exec.Command("git", args)
	cmd.Stderr = nil

}

func (g *Git) Clone(uri string) chan<- Progress {
	return nil
}