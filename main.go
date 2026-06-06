package main

import (
	"fmt"
	"os"

	"github.com/fizyomatik/goconvert/cli"
	"github.com/fizyomatik/goconvert/engine"
)

func main() {
	config, err := cli.ParseCLI(os.Args)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if config.From == "" || config.To == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Error: -from and -to flags are required\n")
		os.Exit(1)
	}

	converter, err := engine.Select(config.From, config.To)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	var reader *os.File
	if config.IsPipe {
		reader = os.Stdin
	} else {
		reader, err = os.Open(config.InputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error opening input file: %s\n", err)
			os.Exit(1)
		}
		defer func(reader *os.File) {
			err := reader.Close()
			if err != nil {

			}
		}(reader)
	}

	var writer *os.File
	if config.OutputFile == "" {
		writer = os.Stdout
	} else {
		writer, err = os.Create(config.OutputFile)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Error creating output file: %s\n", err)
			os.Exit(1)
		}
		defer func(writer *os.File) {
			err := writer.Close()
			if err != nil {

			}
		}(writer)
	}

	if err := converter.Convert(reader, writer); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error converting: %s\n", err)
		os.Exit(1)
	}

	if config.OutputFile != "" {
		_, _ = fmt.Fprintf(os.Stdout, "Done! Output written to %s\n", config.OutputFile)
	}
}
