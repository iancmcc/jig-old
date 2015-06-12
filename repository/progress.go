package repository

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"regexp"
	"strconv"
	"strings"
)

type ProgressStage int

const (
	BEGIN ProgressStage = 1 << iota
	END
	COUNTING
	COMPRESSING
	WRITING
	RECEIVING
	RESOLVING
	MASK = ^(BEGIN | END)
)

const (
	COUNTING_MSG    = "Counting objects"
	COMPRESSING_MSG = "Compressing objects"
	WRITING_MSG     = "Writing objects"
	RECEIVING_MSG   = "Receiving objects"
	RESOLVING_MSG   = "Resolving deltas"
	DONE_MARKER     = ", done."
)

var (
	stage_pattern        = regexp.MustCompile(`(remote: )?([\w\s]+):\s+()(\d+)()(.*)`)
	stage_pattern_totals = regexp.MustCompile(`(remote: )?([\w\s]+):\s+(\d+)% \((\d+)/(\d+)\)(.*)`)
)

type Progress struct {
	stage   ProgressStage // The stage of progress
	Message string        // The message at the end of the line
	Current int           // The current number of objects
	Total   int           // The total number of objects
	Line    string        // The actual output line, for reference
}

func (p *Progress) CurrentStage() ProgressStage {
	return p.stage & MASK
}

func (p *Progress) IsBegin() bool {
	return p.stage&BEGIN != 0
}

func (p *Progress) IsEnd() bool {
	return p.stage&END != 0
}

type ProgressParser struct {
	r    io.Reader
	seen ProgressStage
}

func NewProgressParser(r io.Reader) *ProgressParser {
	return &ProgressParser{r: r}
}

func progressSplitter(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexAny(data, "\n"); i >= 0 {
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

func (parser *ProgressParser) Parse() <-chan *Progress {
	output := make(chan *Progress)
	scanner := bufio.NewScanner(parser.r)
	scanner.Split(progressSplitter)
	go func() {
		defer close(output)
		for scanner.Scan() {
			prog, err := parser.ParseLine(scanner.Text())
			if err == nil {
				output <- prog
			}
		}
	}()
	return output
}

func (parser *ProgressParser) ParseLine(line string) (*Progress, error) {
	var (
		p     Progress
		match []string
	)
	p.Line = line
	match = stage_pattern_totals.FindStringSubmatch(line)
	if match == nil {
		match = stage_pattern.FindStringSubmatch(line)
	}
	if match == nil {
		return &p, errors.New("Unable to parse output")
	}
	switch match[2] {
	case COUNTING_MSG:
		p.stage |= COUNTING
	case COMPRESSING_MSG:
		p.stage |= COMPRESSING
	case WRITING_MSG:
		p.stage |= WRITING
	case RECEIVING_MSG:
		p.stage |= RECEIVING
	case RESOLVING_MSG:
		p.stage |= RESOLVING
	}
	if (parser.seen & p.stage) == 0 {
		// Haven't seen this one before
		parser.seen |= p.stage
		p.stage |= BEGIN
	}
	msg := strings.TrimSpace(match[6])
	p.Message = strings.TrimSuffix(msg, DONE_MARKER)
	if msg != p.Message {
		p.stage |= END
	}
	current, total := match[4], match[5]
	if c, err := strconv.ParseInt(current, 0, 32); err == nil {
		p.Current = int(c)
	}
	if t, err := strconv.ParseInt(total, 0, 32); err == nil {
		p.Total = int(t)
	}
	return &p, nil
}
