package main

import (
	"log"
	"os"
	"strings"

	"github.com/gongo/9t"
)

func main() {
	labeledFiles := parseArgs(os.Args[1:])

	runner, err := ninetail.Runner(labeledFiles, ninetail.Config{Colorize: true}) // TODO use flags!!
	if err != nil {
		log.Fatal(err)
	}
	runner.Run()
}

func parseArgs(args []string) []*ninetail.LabeledFile {
	labeledFiles := make([]*ninetail.LabeledFile, 0, len(os.Args)-1)
	labeledFlag := false
	pipeFlag := false
	for _, arg := range os.Args[1:] {
		if arg == "-l" || arg == "--with-label" {
			labeledFlag = true
			continue
		}
		if arg == "-p" || arg == "--pipe" {
			pipeFlag = true
			continue
		}

		var labeledFile *ninetail.LabeledFile

		if labeledFlag {
			idx := strings.IndexByte(arg, ':')
			labeledFile = &ninetail.LabeledFile{
				FileName: arg[idx+1:],
				Label:    arg[:idx],
				Pipe:     pipeFlag,
			}
		} else {
			labeledFile = &ninetail.LabeledFile{
				FileName: arg,
				Pipe:     pipeFlag,
			}
		}

		labeledFiles = append(labeledFiles, labeledFile)
		labeledFlag = false
		pipeFlag = true
	}
	return labeledFiles
}
