package formats

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type CSVToMD struct{}
type MDToCSV struct{}

func (c *CSVToMD) Convert(r io.Reader, w io.Writer) error {
	records, err := csv.NewReader(r).ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read CSV: %w", err)
	}
	if len(records) == 0 {
		return fmt.Errorf("CSV is empty")
	}

	// Calculate the widest value per column for aligned padding.
	widths := make([]int, len(records[0]))
	for _, row := range records {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	writeRow := func(cells []string) {
		fmt.Fprint(w, "|")
		for i, cell := range cells {
			fmt.Fprintf(w, " %-*s |", widths[i], cell)
		}
		fmt.Fprintln(w)
	}

	writeRow(records[0])

	sep := make([]string, len(records[0]))
	for i, width := range widths {
		sep[i] = strings.Repeat("-", width)
	}
	writeRow(sep)

	for _, row := range records[1:] {
		writeRow(row)
	}

	return nil
}

func (c *MDToCSV) Convert(r io.Reader, w io.Writer) error {
	cw := csv.NewWriter(w)
	scanner := bufio.NewScanner(r)
	rowCount := 0

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || !strings.HasPrefix(line, "|") {
			continue
		}
		if isMDSeparatorRow(line) {
			continue
		}
		cells := parseMDTableRow(line)
		if err := cw.Write(cells); err != nil {
			return fmt.Errorf("failed to write CSV row: %w", err)
		}
		rowCount++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read markdown: %w", err)
	}
	if rowCount == 0 {
		return fmt.Errorf("no markdown table found in input")
	}

	cw.Flush()
	return cw.Error()
}

func parseMDTableRow(line string) []string {
	line = strings.Trim(line, "|")
	parts := strings.Split(line, "|")
	cells := make([]string, len(parts))
	for i, p := range parts {
		cells[i] = strings.TrimSpace(p)
	}
	return cells
}

// isMDSeparatorRow detects lines like | --- | :--- | ---: | that divide header from data.
func isMDSeparatorRow(line string) bool {
	line = strings.Trim(line, "| ")
	for _, cell := range strings.Split(line, "|") {
		cell = strings.TrimSpace(cell)
		cell = strings.Trim(cell, ":-")
		if cell != "" {
			return false
		}
	}
	return true
}
