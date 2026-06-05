package main

import (
	"fmt"
	"os"

	"goconvert/cli"
	"goconvert/engine"
)

func main() {
	config, err := cli.ParseCLI(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	if config.From == "" || config.To == "" {
		fmt.Fprintf(os.Stderr, "Error: -from and -to flags are required\n")
		os.Exit(1)
	}

	converter, err := engine.Select(config.From, config.To)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	var reader *os.File
	if config.IsPipe {
		reader = os.Stdin
	} else {
		reader, err = os.Open(config.InputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening input file: %s\n", err)
			os.Exit(1)
		}
		defer reader.Close()
	}

	var writer *os.File
	if config.OutputFile == "" {
		writer = os.Stdout
	} else {
		writer, err = os.Create(config.OutputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output file: %s\n", err)
			os.Exit(1)
		}
		defer writer.Close()
	}

	if err := converter.Convert(reader, writer); err != nil {
		fmt.Fprintf(os.Stderr, "Error converting: %s\n", err)
		os.Exit(1)
	}

	if config.OutputFile != "" {
		fmt.Fprintf(os.Stdout, "Done! Output written to %s\n", config.OutputFile)
	}
}