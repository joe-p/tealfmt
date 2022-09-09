package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/docopt/docopt-go"
)

const version = "tealfmt v0.1.0"
const usage = `tealfmt

Usage:
  tealfmt <file> [ -e | --edit ]
  tealfmt -h | --help
  tealfmt --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  -e --edit     Edit the file in place.`

var config struct {
	File string
	Edit bool
}

var voidOps = [...]string{"assert", "app_global_put", "b", "bnz", "bz", "store"}

func contains(strs []string, str string) bool {
	for _, ss := range strs {
		if ss == str {
			return true
		}
	}
	return false
}

func handleLabel(newLinesPtr *[]string, commentLines int, trimmedLastLine string, line string) {
	newLines := *newLinesPtr

	// Undo indentation of comments if they are headers comments of a label
	if commentLines > 0 {
		for i := 1; i <= commentLines; i++ {
			idx := len(newLines) - i
			newLines[idx] = strings.TrimSpace(newLines[idx])
		}

		if commentLines > 1 {
			newLines[len(newLines)-commentLines] = "\n" + newLines[len(newLines)-commentLines]
		}

	} else if trimmedLastLine != "" {
		newLines = append(newLines, "")
	}

	newLines = append(newLines, line)

	*newLinesPtr = newLines
}

func handleLine(newLinesPtr *[]string, line string, commentLinesPtr *int, voidOpLinesPtr *bool) {
	newLines := *newLinesPtr
	commentLines := *commentLinesPtr
	voidOpLines := *voidOpLinesPtr

	lastLine := ""
	trimmedLastLine := ""

	if len(newLines) > 0 {
		lastLine = newLines[(len(newLines))-1]
		trimmedLastLine = strings.TrimSpace(lastLine)
	}

	opRegex, _ := regexp.Compile(`\S+`)
	opcode := opRegex.FindString(line)

	labelRegex, _ := regexp.Compile(`\S+:($| //)`)
	commentRegex, _ := regexp.Compile(`^//`)

	// Add a space after any sequence of voidOps
	if contains(voidOps[:], opcode) {
		voidOpLines = true
	} else if voidOpLines == true {
		voidOpLines = false
		newLines = append(newLines, "")
	}

	if labelRegex.MatchString(line) {
		handleLabel(&newLines, commentLines, trimmedLastLine, line)
	} else {
		newLines = append(newLines, "    "+line)
	}

	if commentRegex.MatchString(line) {
		commentLines++
	} else {
		commentLines = 0
	}

	*newLinesPtr = newLines
	*commentLinesPtr = commentLines
	*voidOpLinesPtr = voidOpLines
}

func main() {
	opts, err := docopt.ParseArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatal(err)
	}

	opts.Bind(&config)

	file, err := os.Open(config.File)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	newLines := []string{}
	commentLines := 0
	voidOpLines := false

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if regexp.MustCompile(`^#pragma `).MatchString(line) {
			newLines = append(newLines, line)
			newLines = append(newLines, "")
			continue
		}

		handleLine(&newLines, line, &commentLines, &voidOpLines)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	newContent := strings.Join(newLines[:], "\n")
	if config.Edit {
		os.WriteFile(config.File, []byte(newContent), 0)
	} else {
		fmt.Println(newContent)
	}
}
