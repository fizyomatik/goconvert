# Codebase Structure

**Analysis Date:** 2026-06-06

## Directory Layout

```
goconvert/
├── cli/                    # CLI argument parsing and configuration
│   └── parser.go          # ParseCLI function, Config struct, format validation
├── engine/                # Format selection and converter factory
│   └── engine.go          # Converter interface, Select() routing function
├── formats/               # Format conversion implementations
│   ├── json.go           # JSONToYAML converter
│   ├── yaml.go           # YAMLToJSON converter
│   └── xml.go            # JSONToXML, XMLToJSON, YAMLToXML converters
├── data/                  # Test data files (not executed)
│   ├── book_data.xml
│   ├── dummy_data.yaml
│   └── test_data.json
├── usage/                 # Documentation
│   └── commands.md        # CLI usage modes and examples
├── main.go               # Entry point, I/O orchestration
├── main_test.go          # Test file (mostly empty)
├── go.mod                # Go module definition
├── go.sum                # Dependency checksums
├── README.md             # Project documentation
├── LICENSE               # License file
└── .gitignore            # Git ignore rules
```

## Directory Purposes

**cli/:**
- Purpose: Parse and validate command-line input
- Contains: Config struct definition, CLI parsing logic, format validation
- Key files: `cli/parser.go`

**engine/:**
- Purpose: Route format pairs to appropriate converters using factory pattern
- Contains: Converter interface definition, Select() function for converter selection
- Key files: `engine/engine.go`

**formats/:**
- Purpose: Implement bidirectional format conversions
- Contains: Converter implementations for JSON ↔ YAML, JSON ↔ XML, YAML ↔ XML
- Key files: `formats/json.go`, `formats/yaml.go`, `formats/xml.go`

**data/:**
- Purpose: Test data for manual verification (not executed)
- Contains: Sample JSON, YAML, and XML files for testing conversions
- Key files: Not executable; for development reference

**usage/:**
- Purpose: Documentation and CLI usage examples
- Contains: CLI modes, examples, and flag descriptions
- Key files: `usage/commands.md`

## Key File Locations

**Entry Points:**
- `main.go`: Binary entry point; orchestrates CLI parsing, converter selection, I/O, and error handling

**Configuration & Parsing:**
- `cli/parser.go`: Defines Config struct, ParseCLI() function, SupportedFormats list, format validation

**Core Logic:**
- `engine/engine.go`: Converter interface, Select() factory function for format pair routing
- `formats/json.go`: JSON ↔ YAML conversion implementation
- `formats/yaml.go`: YAML ↔ JSON conversion implementation (handles single/multi-document YAML)
- `formats/xml.go`: JSON ↔ XML, YAML ↔ XML conversions (includes XML root wrapping logic)

**Testing:**
- `main_test.go`: Currently empty package declaration only

**Documentation:**
- `README.md`: Project overview, features, architecture explanation
- `usage/commands.md`: CLI usage modes and examples

## Naming Conventions

**Files:**
- Go source files: lowercase with underscores (e.g., `parser.go`, `main_test.go`)
- Directory names: lowercase (e.g., `cli/`, `formats/`, `data/`)
- No file name prefixes; organization by directory

**Functions:**
- Exported (public): PascalCase (e.g., `ParseCLI`, `Select`, `Convert`)
- Unexported (private): camelCase (e.g., `isSupportedFormat`, `jsonBytesToMap`, `writeXML`)

**Types:**
- Exported: PascalCase (e.g., `Config`, `Converter`, `JSONToYAML`)
- Unexported: none currently used (all types are exported for use across packages)

**Variables & Constants:**
- Package-level exported: PascalCase (e.g., `SupportedFormats`)
- Function-scoped: camelCase (e.g., `config`, `converter`, `reader`)

**Packages:**
- Package names match directory names and are lowercase (e.g., `package cli`, `package formats`)

## Where to Add New Code

**New Format Pair (e.g., JSON ↔ CSV):**

1. **Implementation:** Create converter types in appropriate file under `formats/`:
   - If new format family: Create `formats/csv.go` with `JSONToCSV` and `CSVToJSON` structs
   - If extending existing format: Add to existing file (e.g., add `YAMLToCSV` to `formats/csv.go`)
   - Implement `Convert(r io.Reader, w io.Writer) error` method on each type

2. **Registration:** Add cases to `engine.Select()` in `engine/engine.go`:
   ```go
   case "json->csv":
       return &formats.JSONToCSV{}, nil
   case "csv->json":
       return &formats.CSVToJSON{}, nil
   ```

3. **Validation:** Add new format to `SupportedFormats` slice in `cli/parser.go:11`:
   ```go
   var SupportedFormats = []string{"json", "yaml", "xml", "csv"}
   ```

4. **Tests:** Create test data file in `data/` directory (e.g., `data/test_data.csv`)

**New CLI Feature (e.g., Custom output encoding):**

1. Add field to `Config` struct in `cli/parser.go`
2. Add corresponding flag in `flag.NewFlagSet()` call at `cli/parser.go:47`
3. Handle logic in `main.go` after converter selection

**Error Handling Enhancement:**

- Centralize error types/wrapping in new file `error.go` at package level
- Use consistently across `cli/`, `engine/`, `formats/`

## Special Directories

**data/:**
- Purpose: Test/sample data files
- Generated: No
- Committed: Yes
- Note: Files are static test data, not auto-generated

**.planning/:**
- Purpose: Architecture documentation (created by GSD tools)
- Generated: Yes (by codebase mapper)
- Committed: Yes

**go.mod & go.sum:**
- Purpose: Dependency management (Go modules)
- Generated: Yes (by go mod commands)
- Committed: Yes

## Module Structure

**Module Name:** `goconvert` (defined in `go.mod`)

**Packages:**
- `goconvert` (main) - entry point, orchestration
- `goconvert/cli` - argument parsing
- `goconvert/engine` - converter selection
- `goconvert/formats` - conversion implementations

**Import Paths:**
- Internal: `"goconvert/cli"`, `"goconvert/engine"`, `"goconvert/formats"`
- External: `"gopkg.in/yaml.v3"`, `"github.com/clbanning/mxj/v2"`

---

*Structure analysis: 2026-06-06*
