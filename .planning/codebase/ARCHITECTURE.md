<!-- refreshed: 2026-06-06 -->
# Architecture

**Analysis Date:** 2026-06-06

## System Overview

```text
┌─────────────────────────────────────────────────────────────┐
│                      CLI Input Layer                         │
│  Argument Parsing & Configuration (os.Args -> Config)        │
│                  `cli/parser.go`                             │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                  Format Selection Engine                     │
│  Strategy/Factory: Routes from/to format pairs               │
│                  `engine/engine.go`                          │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│                   Format Converters                          │
│  Implements Converter interface for each format pair         │
│                  `formats/*.go`                              │
└──────────────────────────┬──────────────────────────────────┘
                           │
                           ▼
┌─────────────────────────────────────────────────────────────┐
│          I/O Layer (Files & Streams)                         │
│  os.File handles, os.Stdin, os.Stdout                        │
└─────────────────────────────────────────────────────────────┘
```

## Component Responsibilities

| Component | Responsibility | File |
|-----------|----------------|------|
| CLI Parser | Parse command-line flags, positional args, detect stdin pipe, validate format support | `cli/parser.go` |
| Converter Selector | Select appropriate converter based on from/to format pair | `engine/engine.go` |
| Format Converters | Implement bidirectional conversion logic for JSON/YAML/XML | `formats/json.go`, `formats/yaml.go`, `formats/xml.go` |
| Main Entry | Orchestrate CLI parsing, converter selection, I/O handling, error handling | `main.go` |

## Pattern Overview

**Overall:** Strategy Pattern + Factory Pattern (Hybrid conversion engine)

**Key Characteristics:**
- Converter interface abstracts all conversion implementations
- engine.Select() acts as factory to instantiate correct converter
- Supports multiple execution modes: explicit flags, auto-detection, piping
- Defers I/O to caller (io.Reader/io.Writer) for flexibility

## Layers

**CLI Layer:**
- Purpose: Parse and validate user input (flags, positional args, stdin detection)
- Location: `cli/parser.go`
- Contains: Config struct, ParseCLI function, format validation
- Depends on: Standard library (flag, os, filepath)
- Used by: main.go

**Selection Engine Layer:**
- Purpose: Route from/to format pairs to correct converter implementation
- Location: `engine/engine.go`
- Contains: Converter interface, Select() factory function
- Depends on: formats package
- Used by: main.go

**Format Implementation Layer:**
- Purpose: Perform actual data transformation between format pairs
- Location: `formats/*.go`
- Contains: Type definitions implementing Converter interface
- Depends on: gopkg.in/yaml.v3, github.com/clbanning/mxj/v2 (external libraries)
- Used by: engine.Select()

**I/O & Main Layer:**
- Purpose: Handle file operations, orchestrate pipeline, error handling
- Location: `main.go`
- Contains: File open/close, stdin/stdout detection, converter invocation
- Depends on: cli package, engine package
- Used by: CLI (executed as binary)

## Data Flow

### Primary Request Path (File-based Conversion)

1. User executes binary with CLI args (`main()` at `main.go:11`)
2. CLI parser processes arguments and detects input source (`cli.ParseCLI()` at `cli/parser.go:28`)
3. Config returned with From, To, InputFile, OutputFile, IsPipe flags
4. Main opens input file (file or stdin) based on config (`main.go:29-44`)
5. Engine selects appropriate converter (`engine.Select()` at `engine/engine.go:17`)
6. Main opens output file (file or stdout) based on config (`main.go:46-61`)
7. Converter.Convert(reader, writer) executes transformation (`main.go:63`)
8. Files closed with deferred cleanup (`main.go:38-43, 55-60`)
9. Success message written to stdout (if not piping output)

### Pipe Stream Mode

1. User pipes input: `cat file.json | goconvert -from=json -to=yaml`
2. CLI detection recognizes stdin (os.Stdin.Stat(), `cli/parser.go:30-35`)
3. Config.IsPipe = true, reader set to os.Stdin
4. Same converter selection and execution as above
5. Output streams to os.Stdout or file

### Auto-Detection Mode

1. User provides only filename: `goconvert config.json`
2. CLI infers From from file extension (`.json` -> `json`)
3. CLI infers default To format (json -> yaml, yaml -> json)
4. Auto-generates output filename with new extension
5. Proceeds with normal conversion

**State Management:**
- Config holds all user intent (formats, file paths)
- No global state; each conversion is isolated
- File handles are request-scoped (opened/closed in main)

## Key Abstractions

**Converter Interface:**
- Purpose: Abstract the conversion contract - every format pair must implement Convert(r io.Reader, w io.Writer)
- Examples: `formats.JSONToYAML`, `formats.YAMLToJSON`, `formats.JSONToXML`, `formats.XMLToJSON`, `formats.YAMLToXML`
- Pattern: Strategy Pattern - behavior encapsulation for each format pair

**Config Structure:**
- Purpose: Encapsulate all user intent from CLI parsing
- Fields: InputFile, OutputFile, IsPipe (and missing: From, To - see CONCERNS.md)
- Used by: main.go to orchestrate workflow

## Entry Points

**main():**
- Location: `main.go:11`
- Triggers: Binary execution with CLI arguments
- Responsibilities: 
  - Parse CLI arguments
  - Validate format support
  - Handle file I/O (open, close, cleanup)
  - Select converter
  - Execute conversion
  - Handle errors and exit codes

## Architectural Constraints

- **Threading:** Single-threaded sequential execution (no goroutines used)
- **Global state:** None - all state is request-scoped in main()
- **Circular imports:** None detected
- **I/O Strategy:** Uses io.Reader/io.Writer interfaces for maximum flexibility (works with files, pipes, buffers)
- **XML Special Case:** XML requires single root element; arrays and multi-key objects auto-wrapped in `<root>` tag (`formats/xml.go:20-32`)

## Anti-Patterns

### Missing Config Fields

**What happens:** Code references `config.From` and `config.To` in `cli/parser.go` (lines 49, 51, 58, etc.) but Config struct only defines InputFile, OutputFile, IsPipe. This causes compilation failure.

**Why it's wrong:** The struct definition is incomplete; the code cannot compile.

**Do this instead:** Add From and To fields to Config struct in `cli/parser.go:13`:
```go
type Config struct {
    From       string  // Source format (json, yaml, xml, etc.)
    To         string  // Destination format
    InputFile  string
    OutputFile string
    IsPipe     bool
}
```

### Unused Converter Type

**What happens:** `formats.XMLToYAML` type is defined at `formats/xml.go:16` but not registered in engine.Select() function (`engine/engine.go:17-32`).

**Why it's wrong:** Complete bidirectional conversion tree not implemented despite type existing.

**Do this instead:** Add case in `engine.Select()` for xml->yaml:
```go
case "xml->yaml":
    return &formats.XMLToYAML{}, nil
```

## Error Handling

**Strategy:** Explicit error propagation with formatted context

**Patterns:**
- CLI validation: Format existence checked against SupportedFormats list before processing
- Converter errors: Wrapped with context (e.g., `fmt.Errorf("failed to parse JSON: %w", err)`)
- File I/O errors: Checked immediately with early exit (os.Exit(1))
- No panic-driven design; all errors result in graceful stderr output and exit codes

## Cross-Cutting Concerns

**Logging:** Printf/Fprintf to stderr for errors, stdout for success messages

**Validation:** Format validation in CLI layer before engine selection; data schema validation deferred to format parsers (json.Unmarshal, yaml.Decode, etc.)

**Error Messages:** Consistent format: `Error: <detail>` to stderr with specific context

---

*Architecture analysis: 2026-06-06*
