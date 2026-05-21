	package cli

	import (
		"flag"
		"os"
		"path/filepath"
		"strings"
	)


	type Config struct {
		// Here are 5 configuration options for the CLI application
		From       string // is the source of the data
		InputFile  string // is the file to read data from
		To         string // is the destination of the data
		OutputFile string // is the file to write data to
		IsPipe     bool   // indicates if the data is being piped in

	}

	func ParseCLI(args []string) (*Config, error) {
		config := &Config{}
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, err
		}

		// check if the standard input is being piped in
		isPiped := (stat.Mode() & os.ModeCharDevice) == 0

		// first we check if the user has provided flags 
		hasFlags := false
		for _, arg := range args {
			if strings.HasPrefix(arg, "-") {
				hasFlags = true
				break
			}
		}

		if hasFlags {
			fs := flag.NewFlagSet("goconvert", flag.ContinueOnError)

			fs.StringVar(&config.From, "from", "", "the source format (json or yaml)")
			fs.StringVar(&config.InputFile, "input", "", "the file to read data from")
			fs.StringVar(&config.To, "to", "", "the destination format (json or yaml)")
			fs.StringVar(&config.OutputFile, "output", "", "the file to write data to")
			
			err := fs.Parse(args[1:])
			if err != nil {
				return nil, err
			}

			return config, nil
		}

		if len(args) == 2 {
			config.InputFile = args[1]

			ext := filepath.Ext(args[1])
			config.From = ext[1:] // remove the dot from the extension

			if config.From == "json" {
				config.To = "yaml"	
				config.OutputFile = args[1][:len(args[1])-len(ext)] + ".yaml"
			} else if config.From == "yaml" {
				config.To = "json"
				config.OutputFile = args[1][:len(args[1])-len(ext)] + ".json"
			}
			
		} else if len(args) == 3 {
			config.InputFile = args[1]
			config.OutputFile = args[2]

			ext := filepath.Ext(args[1]) // is the extension of the input file
			config.From = ext[1:] // remove the dot from the extension

			ext = filepath.Ext(args[2]) // is the extension of the output file
			config.To = ext[1:] // remove the dot from the extension

		} else if len(args) == 1 && isPiped {
			config.IsPipe = true
		
		}

		return config, nil
	}