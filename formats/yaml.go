package formats

import (
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type YAMLToJSON struct{}

func (c *YAMLToJSON) Convert(r io.Reader, w io.Writer) error {
	var data interface{}

	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}

	return nil
}
