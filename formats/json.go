package formats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)

type JSONToYAML struct{}

func (c *JSONToYAML) Convert(r io.Reader, w io.Writer) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// strip UTF-8 BOM if present (0xEF 0xBB 0xBF)
	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})

	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	encoder := yaml.NewEncoder(w)
	encoder.SetIndent(2)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}

	return encoder.Close()
}
