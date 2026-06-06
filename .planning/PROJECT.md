# goconverter

## What This Is

goconverter is a smart file format conversion CLI tool built on the goconvert codebase. It only offers conversions that are lossless — users specify a source file and see exactly which target formats are safe, or run a conversion directly. A VSCode extension brings the same experience into the editor via right-click and command palette.

## Core Value

Only offer conversions that work cleanly — no data loss, no guesswork, no surprises.

## Requirements

### Validated

- ✓ Convert JSON ↔ YAML, JSON → XML/CSV/TOML/HCL — existing
- ✓ Convert YAML → XML/CSV — existing
- ✓ Convert XML ↔ JSON/YAML — existing
- ✓ Convert CSV ↔ JSON/YAML, CSV → MD — existing
- ✓ Convert MD → CSV — existing
- ✓ Convert TOML → JSON, HCL → JSON — existing
- ✓ File and stdin/stdout I/O — existing
- ✓ Flag-based CLI (`-from`, `-to`, `-in`, `-out`) — existing
- ✓ Strategy pattern via `Converter` interface — existing

### Active

- [ ] Lossless-only compatibility registry — only conversions with no data loss are allowed or displayed
- [ ] Discovery mode — `goconverter file.json` lists compatible target formats without converting
- [ ] Auto-detect source format from file extension
- [ ] Short command syntax: `goconverter file.json yaml` (no flags required)
- [ ] Helpful error messages — human-readable, specific, no raw Go errors shown
- [ ] Colored terminal output — success (green), errors (red), format lists (cyan)
- [ ] `goconverter --help` with usage examples, not just flag descriptions
- [ ] New format: TSV (lossless ↔ CSV)
- [ ] New format: XLSX (Excel — lossless ↔ CSV/TSV)
- [ ] VSCode extension: right-click file → Convert + command palette support
- [ ] Test suite covering all conversion paths and CLI parser

### Out of Scope

- Lossy conversions (e.g., JSON → HTML) — by design; violates the core value
- JetBrains plugin in v1 — deferred to v2 after VSCode ships
- Web UI / REST API wrapper — different product category
- Batch / directory conversion — future phase
- Streaming large files — performance optimization, not v1 correctness

## Context

**Existing codebase (goconvert):**
- Clean pipeline architecture: `main.go` → `cli/parser.go` → `engine/engine.go` → `formats/*.go`
- Each format pair is a separate struct implementing the `Converter` interface
- I/O is stream-based (`io.Reader` / `io.Writer`) — already works with files and stdin/stdout
- Currently compiled as `goconvert.exe` (Windows); no runtime dependencies

**Known issues to address during this work:**
- Zero test coverage — every conversion path is untested
- Silent `defer file.Close()` errors in `main.go` — write failures go undetected
- Positional-arg mode ignores formats other than JSON/YAML, giving misleading errors
- HCL v1 library (`github.com/hashicorp/hcl v1.0.0`) is archived/unmaintained — needs migration to HCL v2
- `goconvert.exe` binary committed to git — should be gitignored and removed from history
- `go.mod` marks direct dependencies as `// indirect` — needs `go mod tidy`
- Mixed-type JSON arrays produce broken HCL output

**New formats rationale:**
- TSV is structurally identical to CSV (delimiter swap) — perfect lossless candidate
- XLSX is the dominant spreadsheet format; CSV ↔ XLSX is a high-demand conversion

## Constraints

- **Tech stack**: Go — existing codebase, no language change
- **Compatibility**: Must keep existing `-from`/`-to`/`-in`/`-out` flags working (breaking change would hurt existing users)
- **Lossless rule**: The compatibility registry must be conservative — when in doubt, a conversion is NOT listed
- **VSCode extension**: TypeScript/JavaScript (VS Code extension API requirement); shells out to the goconverter binary

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| Build on goconvert | Existing architecture is clean; Converter interface is the right abstraction | — Pending |
| Lossless-only (no tiered approach) | Simpler mental model; users trust the tool completely | — Pending |
| Auto-detect from file extension | Reduces typing; matches user expectation from filename | — Pending |
| VSCode extension in v1 | Primary user environment mentioned; JetBrains deferred | — Pending |
| Test suite as v1 requirement | Zero tests today; unsafe to refactor or add formats without them | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state

---
*Last updated: 2026-06-06 after initialization*
