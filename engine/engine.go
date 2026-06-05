package engine

import (
	"fmt"
	"goconvert/formats"
	"io"
)

// Converter is the interface that every format pair must implement.
// Convert reads from r and writes the converted output to w.
type Converter interface {
	Convert(r io.Reader, w io.Writer) error
}

// Select returns the right Converter for the given from/to format pair.
// Returns an error if the pair is not supported.
func Select(from, to string) (Converter, error) {
	switch from + "->" + to {
	case "json->yaml":
		return &formats.JSONToYAML{}, nil
	case "yaml->json":
		return &formats.YAMLToJSON{}, nil
	case "json->xml":
		return &formats.JSONToXML{}, nil
	case "xml->json":
		return &formats.XMLToJSON{}, nil
	case "yaml->xml":
		return &formats.YAMLToXML{}, nil
	case "xml->yaml":
		return &formats.XMLToYAML{}, nil
	case "json->csv":
		return &formats.JSONToCSV{}, nil
	case "csv->json":
		return &formats.CSVToJSON{}, nil
	case "yaml->csv":
		return &formats.YAMLToCSV{}, nil
	case "csv->yaml":
		return &formats.CSVToYAML{}, nil
	case "csv->md":
		return &formats.CSVToMD{}, nil
	case "md->csv":
		return &formats.MDToCSV{}, nil
	default:
		return nil, fmt.Errorf("unsupported conversion: %s -> %s", from, to)
	}
}
