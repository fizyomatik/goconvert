# Technology Stack

**Analysis Date:** 2026-06-06

## Languages

**Primary:**
- Go 1.23.0 - Full CLI application and format conversion engine

## Runtime

**Environment:**
- Go 1.23.0 runtime

**Build:**
- Go build system
- Platform: Builds to Windows (`goconvert.exe`), cross-compilable to other platforms

## Frameworks

**CLI:**
- Go standard library `flag` package - Command-line argument parsing
  - Location: `cli/parser.go`
  - Manages `-from`, `-to`, `-in`, `-out` flags

**Data Processing:**
- Go standard library `encoding/json` - JSON encoding/decoding
- Go standard library `encoding/xml` - XML encoding/decoding
- Standard library `io` package - Stream-based I/O operations

**Serialization:**
- `gopkg.in/yaml.v3 v3.0.1` - YAML parsing and encoding
  - Used throughout `formats/` package for YAML conversions
  - Provides `yaml.Decoder` and `yaml.Encoder` for streaming

## Key Dependencies

**Critical:**
- `gopkg.in/yaml.v3 v3.0.1` - YAML serialization library
  - Direct dependency
  - Used in: `formats/json.go`, `formats/yaml.go`, `formats/xml.go`
  - Provides RFC 1513 YAML parsing and encoding

**Infrastructure:**
- `github.com/clbanning/mxj/v2 v2.7.0` - Map-based XML library
  - Dependency for XML conversions
  - Used in: `formats/xml.go`
  - Provides `mxj.Map` type for XML ↔ JSON bridge conversions
  - Functions: `mxj.NewMapJson()`, `mxj.NewMapXmlReader()`, `mv.XmlIndent()`

**Testing:**
- `github.com/google/go-cmp v0.7.0` - Indirect dependency via go.mod
  - Available for test comparisons if needed

**Testing Support:**
- `gopkg.in/check.v1 v0.0.0-20161208181325-20d25e280405` - Indirect dependency
  - Available for advanced test frameworks if adopted

## Configuration

**Environment:**
- CLI flags only, no environment variables
- No `.env` file support
- Configuration passed entirely via command-line arguments

**Build:**
- `go.mod` - Module definition and dependency management
- `go.sum` - Lockfile for reproducible builds
- Standard Go build tools (`go build`, `go test`, `go run`)

## Platform Requirements

**Development:**
- Go 1.23.0 or later
- Standard development tools: Go compiler, `go` command-line tool
- No additional runtime dependencies required

**Production:**
- Standalone binary (`goconvert.exe` on Windows)
- Cross-platform compatible (Linux, macOS, Windows)
- No external runtime requirements beyond Go standard library

## Build & Distribution

**Build Output:**
- Standalone executable (`goconvert.exe` on Windows)
- No external dependencies shipped with binary
- Executable size: ~4.2 MB (Windows binary)

**Supported Conversions:**
- JSON ↔ YAML (handled via in-memory AST)
- JSON ↔ XML (handled via `mxj.Map` bridging)
- YAML ↔ XML (routed through JSON intermediate format)
- YAML ↔ JSON (supports multi-document YAML)

---

*Stack analysis: 2026-06-06*
