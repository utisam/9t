package ninetail

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"strings"

	"github.com/hpcloud/tail"
	"github.com/mattn/go-runewidth"
)

var (
	// red, green, yellow, magenta, cyan
	ansiColorCodes  = [...]int{31, 32, 33, 35, 36}
	seekInfoOnStart = &tail.SeekInfo{Offset: 0, Whence: os.SEEK_END}
)

//Tailer contains watches tailed files and contains per-file output parameters
type Tailer struct {
	*tail.Tail
	label     string
	colorCode int
	padding   string
}

//NewTailers creates slice of Tailers from file names.
//Colors of file names are cycled through the list.
//maxWidth is a maximum widht of passed file names, for nice alignment
func NewTailers(labeledFiles []*LabeledFile) ([]*Tailer, error) {
	maxLength := maximumNameLength(labeledFiles)
	ts := make([]*Tailer, len(labeledFiles))

	for i, labeledFile := range labeledFiles {
		t, err := newTailer(labeledFile, getColorCode(i), maxLength)
		if err != nil {
			return nil, err
		}

		ts[i] = t
	}

	return ts, nil
}

func newTailer(labeledFile *LabeledFile, colorCode int, maxWidth int) (*Tailer, error) {
	var location *tail.SeekInfo
	if !labeledFile.Pipe {
		location = seekInfoOnStart
	}

	t, err := tail.TailFile(labeledFile.FileName, tail.Config{
		Follow:   true,
		Location: location,
		Logger:   tail.DefaultLogger,
		Pipe:     labeledFile.Pipe,
	})

	if err != nil {
		return nil, err
	}

	dispNameLength := displayFilenameLength(labeledFile.name())

	return &Tailer{
		Tail:      t,
		label:     labeledFile.Label,
		colorCode: colorCode,
		padding:   strings.Repeat(" ", maxWidth-dispNameLength),
	}, nil
}

//Do formats, colors and writes to stdout appended lines when they happen, exiting on write error
func (t Tailer) Do(output io.Writer) {
	for line := range t.Lines {
		_, err := fmt.Fprintf(
			output,
			"\x1b[%dm%s%s\x1b[0m: %s\n",
			t.colorCode,
			t.padding,
			t.name(),
			line.Text,
		)
		if err != nil {
			return
		}
	}
}

func (t Tailer) name() string {
	if t.label != "" {
		return t.label
	}
	return filepath.Base(t.Filename)
}

func getColorCode(index int) int {
	return ansiColorCodes[index%len(ansiColorCodes)]
}

func maximumNameLength(labeledFiles []*LabeledFile) int {
	max := 0
	for _, labeledFile := range labeledFiles {
		if current := displayFilenameLength(labeledFile.name()); current > max {
			max = current
		}
	}
	return max
}

func displayFilename(filename string) string {
	return filepath.Base(filename)
}

func displayFilenameLength(filename string) int {
	return runewidth.StringWidth(displayFilename(filename))
}
