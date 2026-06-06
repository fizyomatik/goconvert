# External Integrations

**Analysis Date:** 2026-06-06

## APIs & External Services

**None detected**

This project has no external API dependencies or cloud service integrations. All functionality is self-contained within the CLI application.

## Data Storage

**Databases:**
- Not applicable - No database connectivity
- All data is processed in-memory as streams or ASTs

**File Storage:**
- Local filesystem only
- Input: Read from files via `os.Open()` or stdin
- Output: Write to files via `os.Create()` or stdout
- Location of I/O logic: `main.go` (lines 29-61)

**Caching:**
- None - No caching layer

## Authentication & Identity

**Auth Provider:**
- Not applicable - No external authentication
- No API keys, tokens, or credentials required

## Monitoring & Observability

**Error Tracking:**
- Not applicable - No error tracking service
- All errors logged to stderr via `fmt.Fprintf(os.Stderr, ...)`

**Logs:**
- Stderr only for error messages
- Exit codes: 0 (success), 1 (error)
- No structured logging framework

## CI/CD & Deployment

**Hosting:**
- Not applicable - Standalone CLI executable
- Distributed as binary (`goconvert.exe`)

**CI Pipeline:**
- Not detected - No CI configuration files present
- Standard Go testing available via `go test`

## Environment Configuration

**Required env vars:**
- None - Application uses only CLI flags

**Secrets location:**
- Not applicable - No secrets or credentials required

## Webhooks & Callbacks

**Incoming:**
- Not applicable - Standalone CLI tool

**Outgoing:**
- Not applicable - No remote communication

## Stream-Based Processing

**Input Streams:**
- `os.Stdin` - Supports piped input for real-time data processing
  - Detection: `os.Stdin.Stat()` (line 30 in `cli/parser.go`)
  - Used when no input file specified and data is piped

**Output Streams:**
- `os.Stdout` - Default output target when no output file specified
- `os.Stderr` - Error messages and completion status
- File writing: Standard `os.File` for direct file output

## Data Flow

**Primary Processing Path:**
1. User invokes CLI with format flags (`-from`, `-to`)
2. Input source determined: file (`-in`) or stdin
3. Output destination determined: file (`-out`) or stdout
4. Engine selects appropriate converter via `engine.Select(from, to)`
5. Converter streams input to output with format transformation
6. No external services involved at any stage

---

*Integration audit: 2026-06-06*
