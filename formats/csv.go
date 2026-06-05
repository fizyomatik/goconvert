package formats

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"sort"

	"gopkg.in/yaml.v3"
)

type JSONToCSV struct{}
type CSVToJSON struct{}
type YAMLToCSV struct{}
type CSVToYAML struct{}

// valueToCSVString converts a scalar value to its string representation.
// Returns an error for nested objects or arrays — CSV cells must be flat.
func valueToCSVString(v interface{}) (string, error) {
	if v == nil {
		return "", nil
	}
	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return "", fmt.Errorf("nested objects and arrays cannot be represented as a CSV cell — flatten your data first")
	default:
		return fmt.Sprintf("%v", v), nil
	}
}

// sliceToCSV writes a []interface{} of map objects to w as CSV.
// Headers are collected from all rows and sorted alphabetically for deterministic output.
func sliceToCSV(rows []interface{}, w io.Writer) error {
	if len(rows) == 0 {
		return nil
	}

	// Collect all unique keys across all rows.
	seen := make(map[string]bool)
	var headers []string
	for _, row := range rows {
		obj, ok := row.(map[string]interface{})
		if !ok {
			return fmt.Errorf("each element must be an object, got %T", row)
		}
		for k := range obj {
			if !seen[k] {
				seen[k] = true
				headers = append(headers, k)
			}
		}
	}
	sort.Strings(headers)

	cw := csv.NewWriter(w)
	if err := cw.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for i, row := range rows {
		obj := row.(map[string]interface{})
		record := make([]string, len(headers))
		for j, h := range headers {
			s, err := valueToCSVString(obj[h])
			if err != nil {
				return fmt.Errorf("row %d, column %q: %w", i+1, h, err)
			}
			record[j] = s
		}
		if err := cw.Write(record); err != nil {
			return fmt.Errorf("failed to write row %d: %w", i+1, err)
		}
	}

	cw.Flush()
	return cw.Error()
}

// csvToSlice reads CSV from r and returns a slice of string maps.
// Note: all values are strings — CSV carries no type information.
func csvToSlice(r io.Reader) ([]map[string]string, error) {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV: %w", err)
	}
	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}
	headers := records[0]
	result := make([]map[string]string, 0, len(records)-1)
	for _, record := range records[1:] {
		row := make(map[string]string, len(headers))
		for i, h := range headers {
			if i < len(record) {
				row[h] = record[i]
			}
		}
		result = append(result, row)
	}
	return result, nil
}

func (c *JSONToCSV) Convert(r io.Reader, w io.Writer) error {
	raw, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}
	raw = bytes.TrimPrefix(raw, []byte{0xEF, 0xBB, 0xBF})

	var data interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}
	rows, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("JSON input must be an array of objects")
	}
	return sliceToCSV(rows, w)
}

func (c *CSVToJSON) Convert(r io.Reader, w io.Writer) error {
	rows, err := csvToSlice(r)
	if err != nil {
		return err
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(rows); err != nil {
		return fmt.Errorf("failed to write JSON: %w", err)
	}
	return nil
}

func (c *YAMLToCSV) Convert(r io.Reader, w io.Writer) error {
	var data interface{}
	if err := yaml.NewDecoder(r).Decode(&data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}
	rows, ok := data.([]interface{})
	if !ok {
		return fmt.Errorf("YAML input must be a sequence of mappings")
	}
	return sliceToCSV(rows, w)
}

func (c *CSVToYAML) Convert(r io.Reader, w io.Writer) error {
	rows, err := csvToSlice(r)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(rows); err != nil {
		return fmt.Errorf("failed to write YAML: %w", err)
	}
	return enc.Close()
}
