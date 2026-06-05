package formats

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/clbanning/mxj/v2"
	"gopkg.in/yaml.v3"
)

type JSONToXML struct{}
type XMLToJSON struct{}
type YAMLToXML struct{}
type XMLToYAML struct{}

// jsonBytesToMap parses JSON bytes into an mxj.Map.
// Top-level arrays and multi-key objects are wrapped in <root> because XML requires a single root element.
func jsonBytesToMap(raw []byte) (mxj.Map, error) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) > 0 && trimmed[0] == '[' {
		raw = append([]byte(`{"root":`), append(trimmed, '}')...)
	}
	mv, err := mxj.NewMapJson(raw)
	if err != nil {
		return nil, err
	}
	if len(mv) > 1 {
		mv = mxj.Map{"root": map[string]interface{}(mv)}
	}
	return mv, nil
}

func writeXML(mv mxj.Map, w io.Writer) error {
	xmlBytes, err := mv.XmlIndent("", "  ")
	if err != nil {
		return err
	}
	xmlBytes = bytes.TrimRight(xmlBytes, "\n")
	_, err = fmt.Fprintf(w, "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n%s\n", xmlBytes)
	return err
}

func (c *JSONToXML) Convert(r io.Reader, w io.Writer) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})

	mv, err := jsonBytesToMap(raw)
	if err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	if err := writeXML(mv, w); err != nil {
		return fmt.Errorf("failed to write XML: %w", err)
	}
	return nil
}

func (c *XMLToJSON) Convert(r io.Reader, w io.Writer) error {
	mv, err := mxj.NewMapXmlReader(r)
	if err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(map[string]interface{}(mv)); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	return nil
}

func (c *YAMLToXML) Convert(r io.Reader, w io.Writer) error {
	var data interface{}
	if err := yaml.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	// Route through JSON bytes so mxj can build its internal map representation.
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal intermediate: %w", err)
	}
	mv, err := jsonBytesToMap(jsonBytes)
	if err != nil {
		return fmt.Errorf("failed to build XML structure: %w", err)
	}
	if err := writeXML(mv, w); err != nil {
		return fmt.Errorf("failed to write XML: %w", err)
	}
	return nil
}

func (c *XMLToYAML) Convert(r io.Reader, w io.Writer) error {
	mv, err := mxj.NewMapXmlReader(r)
	if err != nil {
		return fmt.Errorf("failed to parse XML: %w", err)
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(map[string]interface{}(mv)); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	return enc.Close()
}
