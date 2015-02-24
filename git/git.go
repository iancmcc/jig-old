package git

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

func scanOutput(output io.Reader) <-chan string {
	outchan := make(chan string)
	scanner := bufio.NewScanner(output)
	go func() {
		defer close(outchan)
		for scanner.Scan() {
			outchan <- strings.TrimSpace(scanner.Text())
		}
	}()
	return outchan
}

func scanProgress(output io.Reader) <-chan Progress {
	outchan := make(chan Progress)
	scanner := bufio.NewScanner(output)
	split := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, ''); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, bytes.TrimSpace(data[0:i]), nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), bytes.TrimSpace(data), nil
		}
		// Request more data.
		return 0, nil, nil
	}
	scanner.Split(split)
	go func() {
		defer close(outchan)
		for scanner.Scan() {
			prog, err := ParseProgressLine(scanner.Text())
			if err != nil {
				continue
			}
			outchan <- *prog
		}
	}()
	return outchan
}

func (g *Git) Run(args ...string) (output <-chan string, err error) {
	var stderr, stdout io.Reader
	cmd := exec.Command("git", args...)
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, err
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return nil, err
	}
	outchan := scanOutput(stdout)
	errchan := scanOutput(stderr)
	combined := make(chan string)
	var wait sync.WaitGroup
	wait.Add(2)
	go func() {
		defer wait.Done()
		for line := range outchan {
			combined <- line
		}
	}()
	go func() {
		defer wait.Done()
		for line := range errchan {
			combined <- line
		}
	}()
	go func() {
		wait.Wait()
		close(combined)
	}()
	go cmd.Run()
	return combined, nil
}

func (g *Git) RunWithProgress(args ...string) (output <-chan string, progress <-chan Progress, err error) {
	var stdout, stderr io.Reader
	cmd := exec.Command("git", args...)
	if stdout, err = cmd.StdoutPipe(); err != nil {
		return nil, nil, err
	}
	if stderr, err = cmd.StderrPipe(); err != nil {
		return nil, nil, err
	}
	output = scanOutput(stdout)
	progress = scanProgress(stderr)
	return output, progress, nil
}

func (g *Git) Clone(uri string) (<-chan string, <-chan Progress, error) {
	return g.RunWithProgress("clone", uri)
}