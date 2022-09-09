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

type Config struct {
	File    string `docopt:"<file>"`
	InPlace bool   `docopt:"-i,--inplace"`
}
type Line struct {
	Text     string
	Comments []string
	IsVoid   bool
	IsLabel  bool
	Op       string
}

// ToString returns a string representing the TEAL line
// With any comments above it matching its indent settings
func (l *Line) ToString(insideBody bool) string {
	indent := ""
	if insideBody {
		indent = "\t"
	}

	newline := ""
	if l.IsVoid || l.IsLabel {
		newline = "\n"
	}

	comments := ""
	for _, comment := range l.Comments {
		comments += fmt.Sprintf("%s%s\n", indent, comment)
	}

	line := fmt.Sprintf("%s%s%s", indent, l.Text, newline)

	return fmt.Sprintf("%s%s\n", comments, line)
}

const (
	version = "tealfmt v0.1.0"
	usage   = `tealfmt

Usage:
  tealfmt <file> [ -i | --inplace ]
  tealfmt -h | --help
  tealfmt --version

Options:
  -h --help     Show this screen.
  --version     Show version.
  -i --inplace  Edit the file in place.`
)

var (
	voidOps = []string{
		"assert", "err", "return", "app_global_put", "b", "bnz", "bz", "store",
		"stores", "app_local_put", "app_global_del", "app_local_del", "callsub",
		"log", "itxn_submit", "itxn_next",
	}

	versionMatch = regexp.MustCompile(`^#pragma `)
	opMatch      = regexp.MustCompile(`\S+`)
	labelRegex   = regexp.MustCompile(`\S+:($| //)`)
	commentRegex = regexp.MustCompile(`^//`)
	voidOpMatch  = regexp.MustCompile("^(" + strings.Join(voidOps, "|") + ")$")
)

func main() {
	opts, err := docopt.ParseArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatalf("failed to parse args: %s", err)
	}

	config := Config{}
	err = opts.Bind(&config)
	if err != nil {
		log.Fatalf("failed to bind args: %s:", err)
	}

	file, err := os.Open(config.File)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// First pass, add spacing
	newLines := []Line{}
	commentBuff := []string{}
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// If its the pragma line, add it without modification
		if versionMatch.MatchString(line) {
			newLines = append(newLines, Line{Text: line, IsVoid: true, Comments: commentBuff})
			continue
		}

		if commentRegex.MatchString(line) {
			commentBuff = append(commentBuff, line)
			continue
		}
		op := opMatch.FindString(line)
		newLines = append(newLines, Line{
			Text:     line,
			Comments: commentBuff,
			Op:       op,
			IsVoid:   voidOpMatch.MatchString(op),
			IsLabel:  labelRegex.MatchString(line),
		})

		commentBuff = nil
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	s := strings.Builder{}
	insideBody := false
	for idx, line := range newLines {
		if idx > 0 && newLines[idx-1].IsLabel {
			insideBody = true
		}

		s.Write([]byte(line.ToString(insideBody && !line.IsLabel)))
	}

	if config.InPlace {
		os.WriteFile(config.File, []byte(s.String()), 0)
	} else {
		fmt.Println(s.String())
	}
}
