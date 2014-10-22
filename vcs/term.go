package vcs

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	tm "github.com/buger/goterm"
	"github.com/cheggaaa/pb"
)

type lastWriteBuffer struct {
	buffer *bytes.Buffer
}

func newLastWriteBuffer() *lastWriteBuffer {
	return &lastWriteBuffer{&bytes.Buffer{}}
}

type ProgressBarBank struct {
	bars       map[*pb.ProgressBar]*lastWriteBuffer
	barorder   []*pb.ProgressBar
	names      map[*pb.ProgressBar]string
	started    bool
	isFinished bool
}

func NewProgressBarBank() *ProgressBarBank {
	return &ProgressBarBank{
		bars:     make(map[*pb.ProgressBar]*lastWriteBuffer),
		barorder: []*pb.ProgressBar{},
		names:    map[*pb.ProgressBar]string{},
	}
}

func (b *lastWriteBuffer) Write(p []byte) (n int, err error) {
	b.buffer.Reset()
	return b.buffer.Write(p)
}

func (b *lastWriteBuffer) String() string {
	return b.buffer.String()
}

func (b *ProgressBarBank) updateWidth() {
	var max int
	for _, name := range b.names {
		l := len(name)
		if l > max {
			max = l
		}
	}
	w := tm.Width() - max - 2
	for bar, _ := range b.bars {
		bar.SetWidth(w)
	}
}

func (b *ProgressBarBank) StartNew(count int, prefix string) *pb.ProgressBar {
	defer b.Start()
	bar := pb.New(count)
	defer bar.Start()
	defer b.updateWidth()
	buffer := newLastWriteBuffer()
	b.bars[bar] = buffer
	b.barorder = append(b.barorder, bar)
	b.names[bar] = prefix
	bar.ShowCounters = false
	bar.Output = buffer
	bar.NotPrint = true
	return bar
}

func clearScreen() {
	tm.Clear()
	tm.MoveCursor(1, 1)
}

func (b *ProgressBarBank) Render() {
	clearScreen()
	table := tm.NewTable(0, 5, 1, ' ', 0)
	for _, k := range b.barorder {
		fmt.Fprintf(table, "%s\t%s\n", tm.Color(b.names[k], tm.WHITE), strings.Trim(b.bars[k].String(), "\r"))
	}
	tm.Println(table)
	tm.Flush()
}

func (b *ProgressBarBank) writer() {
	for {
		if b.isFinished {
			break
		}
		b.Render()
		time.Sleep(200 * time.Millisecond)
	}
}

func (b *ProgressBarBank) Start() {
	if !b.started {
		go b.writer()
	}
	b.started = true
}