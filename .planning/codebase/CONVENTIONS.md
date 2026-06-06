# Coding Conventions

**Analysis Date:** 2026-06-06

## Naming Patterns

**Files:**
- Package files use lowercase with no underscores: `parser.go`, `engine.go`, `json.go`
- Test files follow Go convention: `*_test.go` (e.g., `main_test.go`)
- Converter implementations use source-to-target format: `json.go`, `yaml.go`, `xml.go`

**Packages:**
- Lowercase, single word: `main`, `cli`, `engine`, `formats`
- Package names match directory names

**Types:**
- Exported types use PascalCase: `Config`, `Converter`, `JSONToYAML`, `XMLToJSON`, `YAMLToXML`
- Empty struct types are used for method receivers (e.g., `type JSONToYAML struct{}`)
- Interface names use agent noun pattern: `Converter` (not `ConvertInterface`)

**Functions & Methods:**
- Exported functions: PascalCase (e.g., `ParseCLI`, `Select`, `Convert`)
- Unexported functions: camelCase (e.g., `isSupportedFormat`, `jsonBytesToMap`, `writeXML`)
- Method receivers use single lowercase letter (typically `c` for converter): `func (c *JSONToYAML) Convert(...)`
- All methods use pointer receivers for consistency

**Variables:**
- Local variables: camelCase (e.g., `config`, `converter`, `reader`, `writer`, `raw`, `encoder`, `decoder`)
- Struct fields: PascalCase/exported when part of public struct (e.g., `InputFile`, `OutputFile`, `IsPipe`)
- Interface values: lowercase (e.g., `r io.Reader`, `w io.Writer`)

**Constants & Module-Level Variables:**
- Exported slices: PascalCase (e.g., `SupportedFormats`)
- Alphabetically organized when multiple variables

## Code Style

**Formatting:**
- Standard Go formatting via `gofmt` (implied, no `.editorconfig` or explicit formatter config)
- Indentation: tabs (Go standard)
- Line length: lines are generally kept under 100 characters but no enforced limit detected

**Linting:**
- No `.golangci.yml` or linting configuration file present
- No explicit linting tool configuration detected
- Rely on standard Go conventions and manual code review

**Comments:**
- Proper English sentences starting with the name being documented
- Example from `engine.go`: `// Converter is the interface that every format pair must implement.`
- Inline comments explain non-obvious logic
- Example from `formats/xml.go`: `// strip UTF-8 BOM if present (0xEF 0xBB 0xBF)`
- Comments above functions/types rather than inline where possible

## Import Organization

**Order:**
1. Standard library imports (alphabetically sorted)
2. Blank line separator
3. Third-party imports (alphabetically sorted)

**Example from `formats/json.go`:**
```go
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"gopkg.in/yaml.v3"
)
```

**Path Aliases:**
- No import aliases used in the codebase
- Full paths used: `"goconvert/cli"`, `"goconvert/engine"`, `"goconvert/formats"`
- Standard package imports use default names

## Error Handling

**Patterns:**
- All error-returning functions use `if err != nil { return err }` pattern
- Errors wrapped with context using `fmt.Errorf()` and `%w` verb
- Example from `cli/parser.go`: `return fmt.Errorf("error parsing CLI arguments: %w", err)`
- Descriptive error messages that identify the operation and reason for failure
- Error prefix matches operation context: "failed to read input", "failed to parse JSON", "failed to write YAML"

**Stderr Output:**
- Write errors to `os.Stderr` using `fmt.Fprintf()` for user-facing messages
- Discard write-to-stderr errors using blank identifier: `_, _ = fmt.Fprintf(os.Stderr, ...)`
- Example from `main.go`: `_, _ = fmt.Fprintf(os.Stderr, "Error: %s\n", err)`

**Exit Codes:**
- Program exits with `os.Exit(1)` on any error
- Success path completes without explicit exit (implicit exit code 0)

## Function Design

**Size:**
- Functions kept relatively compact (largest is `ParseCLI` at ~90 lines with complex logic)
- Helper functions extracted when logic is reused or adds clarity
- Example: `jsonBytesToMap` in `formats/xml.go` handles JSON parsing complexity

**Parameters:**
- Use simple types or interfaces: `string`, `[]string`, `io.Reader`, `io.Writer`
- Interface parameters (`io.Reader`, `io.Writer`) preferred over concrete types for flexibility
- Receiver methods use pointer receivers for consistency: `func (c *JSONToYAML) Convert(...)`

**Return Values:**
- Error returned as last return value (Go idiom)
- Single return type plus error: `(*Config, error)`, `(Converter, error)`, `error`
- All error-returning functions follow this pattern

## Module Design

**Package Structure:**
- Clear separation of concerns:
  - `main.go`: CLI entry point and I/O orchestration
  - `cli/`: Command-line argument parsing
  - `engine/`: Converter selection (factory pattern)
  - `formats/`: Conversion implementations

**Exports:**
- Each format module exports only the struct types needed by `engine.Select()`
- No factory functions per format (factory lives in `engine.Select()`)
- Helper functions kept unexported: `jsonBytesToMap()`, `writeXML()`

**Interfaces:**
- Single interface `Converter` defined in `engine/engine.go`
- All format structs implicitly satisfy `Converter` interface
- Interface consists of single method: `Convert(r io.Reader, w io.Writer) error`

## Data Handling

**Standard Library I/O:**
- All data flows through `io.Reader` and `io.Writer` interfaces
- `io.ReadAll()` used for complete data consumption before processing
- Encoders/decoders created from readers/writers for format-specific processing

**Struct Field Tags:**
- No struct tags used (JSON/YAML marshaling handled by encoding libraries directly)
- All data processing uses generic `interface{}` rather than typed structs

**Nil Handling:**
- Empty `interface{}` slices are appended to (not pre-allocated)
- No null pointer checks needed due to interface{} usage patterns

## Testing Notes

- No `_test.go` files with actual test implementations yet
- Test data files stored in `data/` directory: `test_data.json`, `dummy_data.yaml`, `book_data.xml`
- When tests are implemented, follow Go testing conventions:
  - Test functions named `Test*` (e.g., `TestParseJSON`)
  - Table-driven tests preferred for multiple input scenarios
  - Use `t.Helper()` for test utility functions

---

*Convention analysis: 2026-06-06*
