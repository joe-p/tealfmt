package tealfmt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
)

var (
	voidOps = []string{
		"assert", "err", "return", "app_global_put", "b", "bnz", "bz", "store",
		"stores", "app_local_put", "app_global_del", "app_local_del", "callsub", "retsub",
		"log", "itxn_submit", "itxn_next",
	}

	versionMatch = regexp.MustCompile(`^#pragma `)
	opMatch      = regexp.MustCompile(`\S+`)
	labelRegex   = regexp.MustCompile(`\S+:($| //)`)
	commentRegex = regexp.MustCompile(`^//`)
	voidOpMatch  = regexp.MustCompile("^(" + strings.Join(voidOps, "|") + ")$")
)

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
	// If we're inside the body of a label
	// we should add a tab
	indent := ""
	if insideBody {
		indent = "\t"
	}

	// A new line is added after a label, void op, or any op with comments
	newline := ""
	if l.IsVoid || l.IsLabel || len(l.Comments) > 0 {
		newline = "\n"
	}

	// Set appropriate indent on each comment
	comments := ""
	for _, comment := range l.Comments {
		comments += fmt.Sprintf("%s%s\n", indent, comment)
	}

	// Add any indent, and an extra newline if necessary
	line := fmt.Sprintf("%s%s%s", indent, l.Text, newline)

	return fmt.Sprintf("%s%s\n", comments, line)
}

func Format(r io.Reader) string {
	var (
		newLines    []Line
		commentBuff []string

		scanner = bufio.NewScanner(r)
	)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// If its the pragma line, add it without modification
		if versionMatch.MatchString(line) {
			newLines = append(newLines, Line{Text: line, IsVoid: true, Comments: commentBuff})
			continue
		}

		// If its a comment, add it to the buffer to be
		// associated with the first non comment line
		if commentRegex.MatchString(line) {
			commentBuff = append(commentBuff, line)
			continue
		}

		// Construct a Line from the comment buf and flags
		op := opMatch.FindString(line)
		newLines = append(newLines, Line{
			Text:     line,
			Comments: commentBuff,
			Op:       op,
			IsVoid:   voidOpMatch.MatchString(op),
			IsLabel:  labelRegex.MatchString(line),
		})

		// Reset the comment buff
		commentBuff = nil
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to scan file: %+v", err)
	}

	output := ""
	insideBody := false
	for idx, line := range newLines {
		if idx > 0 && newLines[idx-1].IsLabel {
			insideBody = true
		}

		// Remove newline from a previous isVoid
		if idx > 0 && line.IsVoid && newLines[idx-1].IsVoid {
			output = output[:len(output)-1]
		}

		// We should indent if we're inside a body and
		// haven't hit a new label
		indented := insideBody && !line.IsLabel

		output += line.ToString(indented)
	}

	return output
}
