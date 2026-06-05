package cli

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var SupportedFormats = []string{"json", "yaml"}

type Config struct {
	From       string
	InputFile  string
	To         string
	OutputFile string
	IsPipe     bool
}

func isSupportedFormat(format string) bool {
	for _, f := range SupportedFormats {
		if f == format {
			return true
		}
	}
	return false
}

func ParseCLI(args []string) (*Config, error) {
	config := &Config{}
	stat, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}

	isPiped := (stat.Mode() & os.ModeCharDevice) == 0

	// Only scan args[1:] — the binary name (args[0]) is never a flag.
	hasFlags := false
	for _, arg := range args[1:] {
		if strings.HasPrefix(arg, "-") {
			hasFlags = true
			break
		}
	}

	if hasFlags {
		fs := flag.NewFlagSet("goconvert", flag.ContinueOnError)

		fs.StringVar(&config.From, "from", "", "the source format (e.g. json, yaml)")
		fs.StringVar(&config.InputFile, "in", "", "the file to read data from")
		fs.StringVar(&config.To, "to", "", "the destination format (e.g. json, yaml)")
		fs.StringVar(&config.OutputFile, "out", "", "the file to write data to")

		if err := fs.Parse(args[1:]); err != nil {
			return nil, fmt.Errorf("error parsing CLI arguments: %w", err)
		}

		if config.From != "" && !isSupportedFormat(config.From) {
			return nil, fmt.Errorf("unsupported format: %s. Supported: %v", config.From, SupportedFormats)
		}
		if config.To != "" && !isSupportedFormat(config.To) {
			return nil, fmt.Errorf("unsupported format: %s. Supported: %v", config.To, SupportedFormats)
		}

		// Carry pipe detection into flag mode so main.go reads from stdin.
		config.IsPipe = isPiped
		return config, nil
	}

	if len(args) == 2 {
		config.InputFile = args[1]

		ext := filepath.Ext(args[1])
		if len(ext) < 2 {
			return nil, fmt.Errorf("no valid file extension found for: %s", args[1])
		}
		config.From = ext[1:]

		if !isSupportedFormat(config.From) {
			return nil, fmt.Errorf("unsupported file format: '%s'. Supported: %v", config.From, SupportedFormats)
		}

		if config.From == "json" {
			config.To = "yaml"
			config.OutputFile = strings.TrimSuffix(args[1], ext) + ".yaml"
		} else if config.From == "yaml" {
			config.To = "json"
			config.OutputFile = strings.TrimSuffix(args[1], ext) + ".json"
		}

	} else if len(args) == 3 {
		config.InputFile = args[1]
		config.OutputFile = args[2]

		extIn := filepath.Ext(args[1])
		if len(extIn) < 2 {
			return nil, fmt.Errorf("no valid file extension found for: %s", args[1])
		}
		config.From = extIn[1:]

		extOut := filepath.Ext(args[2])
		if len(extOut) < 2 {
			return nil, fmt.Errorf("no valid file extension found for: %s", args[2])
		}
		config.To = extOut[1:]

		if !isSupportedFormat(config.From) {
			return nil, fmt.Errorf("unsupported file format: '%s'. Supported: %v", config.From, SupportedFormats)
		}
		if !isSupportedFormat(config.To) {
			return nil, fmt.Errorf("unsupported file format: '%s'. Supported: %v", config.To, SupportedFormats)
		}
	}

	return config, nil
}
