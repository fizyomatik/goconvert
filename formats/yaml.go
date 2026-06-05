package formats

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type YAMLToJSON struct{}

func (c *YAMLToJSON) Convert(r io.Reader, w io.Writer) error {
	decoder := yaml.NewDecoder(r)

	var docs []interface{}
	for {
		var doc interface{}
		err := decoder.Decode(&doc)
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to parse YAML: %w", err)
		}
		docs = append(docs, doc)
	}

	// Single document: encode directly. Multiple documents: encode as array.
	var out interface{}
	switch len(docs) {
	case 0:
		out = nil
	case 1:
		out = docs[0]
	default:
		out = docs
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(out); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}
