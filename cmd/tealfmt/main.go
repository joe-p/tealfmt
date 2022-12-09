package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/joe-p/tealfmt"
)

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

type Config struct {
	File    string `docopt:"<file>"`
	InPlace bool   `docopt:"-i,--inplace"`
}

func main() {
	opts, err := docopt.ParseArgs(usage, os.Args[1:], version)
	if err != nil {
		log.Fatalf("failed to parse args: %s", err)
	}

	config := Config{}
	if err = opts.Bind(&config); err != nil {
		log.Fatalf("failed to bind args: %s:", err)
	}

	file, err := os.Open(config.File)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	output := tealfmt.Format(file)

	if config.InPlace {
		os.WriteFile(config.File, []byte(output), 0)
	} else {
		fmt.Println(output)
	}
}
