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

var config struct {
	File string `docopt:"<file>"`
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

func main() {
	usage := `LinTEAL

Usage:
  linteal <file>
  linteal -h | --help
  linteal --version

Options:
  -h --help     Show this screen.
  --version     Show version.`

	opts, err := docopt.ParseArgs(usage, os.Args[1:], "")
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
	voidOpLines := true

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		pragmaRegex := regexp.MustCompile(`^#pragma `)
		if pragmaRegex.MatchString(line) {
			newLines = append(newLines, line)
			continue
		}

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
		} else {
			newLines = append(newLines, "    "+line)
		}

		if commentRegex.MatchString(line) {
			commentLines++
		} else {
			commentLines = 0
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(strings.Join(newLines[:], "\n"))
}
