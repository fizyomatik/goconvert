# Codebase Concerns

**Analysis Date:** 2026-06-06

## Critical Issues

### Project Does Not Compile

**Issue:** The `Config` struct in `cli/parser.go` is missing the `From` and `To` fields, but the code references these fields throughout the parser.

**Files:** 
- `cli/parser.go` (lines 49, 51, 58-62, 77-112)
- `main.go` (lines 18, 23)

**Impact:** The project cannot be built. Running `go build` produces 10+ compilation errors:
```
cli/parser.go:49:24: config.From undefined (type *Config has no field or method From)
cli/parser.go:51:24: config.To undefined (type *Config has no field or method To)
```

**Fix approach:** Add `From` and `To` fields to the `Config` struct definition at `cli/parser.go:13`:
```go
type Config struct {
    InputFile  string
    OutputFile string
    From       string  // Add this field
    To         string  // Add this field
    IsPipe     bool
}
```

**Priority:** CRITICAL - Blocks all development and testing.

---

## Tech Debt

### README Documents Unimplemented Features

**Issue:** The `README.md` prominently advertises support for CSV and Markdown formats ("Multi-Format Universe: Initial support for `JSON`, `YAML`, `CSV`, and Github-flavored `Markdown Tables`"), but these formats are not implemented.

**Files:**
- `README.md` (lines 3, 22, 38)
- `cli/parser.go` (line 49, 51 - help text mentions csv and md)

**Why it's wrong:** Users will attempt to convert CSV or Markdown files and encounter "unsupported format" errors, creating a broken user experience. The code in `SupportedFormats` at `cli/parser.go:11` only includes `["json", "yaml", "xml"]`.

**Context:** Git history shows CSV and Markdown converters (`formats/csv.go`, `formats/md.go`) were attempted but then reverted in commit `c9f830b`, indicating they were not working correctly.

**Fix approach:**
1. Update README to accurately document only JSON, YAML, and XML support
2. Remove CSV and Markdown from CLI help text in `cli/parser.go` lines 49, 51
3. Either implement CSV/Markdown converters properly or close the feature permanently

**Priority:** HIGH - Affects user expectations and credibility.

### Missing Conversion Route: XML-to-YAML

**Issue:** The `YAMLToXML` type is defined in `formats/xml.go` (lines 15, 76-94), but its inverse conversion `XMLToYAML` is also defined (lines 16, 96-107) but NOT registered in the `Select()` function dispatcher.

**Files:**
- `formats/xml.go` (lines 96-107 - implementation exists)
- `engine/engine.go` (lines 17-32 - Select function missing case)

**Impact:** Users cannot convert XML to YAML (`xml->yaml`), even though the code to do so exists and is implemented.

**Fix approach:** Add this case to `engine/engine.go` Select function:
```go
case "xml->yaml":
    return &formats.XMLToYAML{}, nil
```

**Priority:** MEDIUM - Feature is partially implemented but unavailable.

### Incomplete Format Auto-Detection in Single-File Mode

**Issue:** When a single input file is provided (2 arguments), the parser only handles auto-conversion for JSON and YAML formats. XML files are supported but have no auto-conversion logic.

**Files:** `cli/parser.go` (lines 83-89)

**What happens:** 
```go
if config.From == "json" {
    config.To = "yaml"
    config.OutputFile = strings.TrimSuffix(args[1], ext) + ".yaml"
} else if config.From == "yaml" {
    config.To = "json"
    config.OutputFile = strings.TrimSuffix(args[1], ext) + ".json"
}
// No handling for xml, so config.To remains empty
```

When an XML file is passed as the only argument, `config.To` remains an empty string, and the main function validates this with an error at `main.go:18-20`.

**Fix approach:** Either:
1. Extend the parser to handle XML auto-conversion (defaulting to JSON), or
2. Document that single-file mode only works with JSON/YAML, or
3. Require explicit `-to` flag when using XML

**Priority:** MEDIUM - Limits usability of single-file mode.

---

## Error Handling Issues

### Deferred File Close Operations Silently Ignore Errors

**Issue:** In `main.go`, the deferred close operations for both reader (lines 38-43) and writer (lines 55-60) files have empty error handlers that silently discard close errors.

**Files:** `main.go` (lines 38-43, 55-60)

**Current code:**
```go
defer func(reader *os.File) {
    err := reader.Close()
    if err != nil {
        // Empty block - error is ignored
    }
}(reader)
```

**Why it's wrong:** If a file fails to close properly (e.g., write to disk fails, file is locked), the error is silently ignored. This could mask data loss or corruption issues.

**Fix approach:** Log or handle the error:
```go
defer func(reader *os.File) {
    if err := reader.Close(); err != nil {
        _, _ = fmt.Fprintf(os.Stderr, "Warning: failed to close input file: %s\n", err)
    }
}(reader)
```

**Priority:** MEDIUM - Data integrity concern.

---

## Testing Coverage Gaps

### No Test Coverage

**Issue:** The `main_test.go` file exists but is completely empty (only contains `package main` declaration).

**Files:** `main_test.go` (1 line)

**What's not tested:**
- CLI argument parsing (all code paths in `cli/parser.go`)
- Format conversions (all conversion functions in `formats/*.go`)
- Error handling and edge cases
- File I/O operations
- Main entry point logic

**Risk:** Any refactoring, dependency upgrades, or format changes could introduce silent bugs that are only discovered by end users.

**Priority:** HIGH - No safety net for changes.

---

## Known Limitations

### Missing Error Handling for Invalid Arguments

**Issue:** The argument parser has conditional logic for 2, 3, or more arguments, but doesn't explicitly handle invalid argument counts (e.g., 4+ arguments).

**Files:** `cli/parser.go` (lines 70-115)

**What happens:** If a user provides 4+ arguments without using flags, the parser silently returns `*Config` with zero values, and the caller must detect the missing required fields.

**Fix approach:** Add validation at the end of ParseCLI:
```go
if config.From == "" || config.To == "" {
    return nil, fmt.Errorf("unable to determine conversion format from arguments")
}
```

**Priority:** LOW - Caught later by main.go validation.

### XMLToYAML Implementation Not Registered

**Issue:** While `XMLToYAML` is implemented in `formats/xml.go`, it's never returned by the `Select()` function, making it impossible for users to invoke.

**Files:**
- `formats/xml.go` (lines 96-107)
- `engine/engine.go` (lines 17-32)

**Impact:** The inverse of the working `YAMLToXML` conversion is unavailable.

**Priority:** MEDIUM.

---

## Maintainability Concerns

### Hardcoded Format List in Multiple Places

**Issue:** The list of supported formats is duplicated:
- `cli/parser.go:11` - `SupportedFormats` array
- `main.go:14` - `fmt.Fprintf` mentions `-from` and `-to`
- Help text in `cli/parser.go:49, 51` - Mentions unsupported formats (csv, md)
- `engine/engine.go:18` - Switch statement with hardcoded conversions

**Why it's fragile:** Adding a new format requires changes in multiple locations. Easy to forget updating one or miss a conversion pair.

**Fix approach:** Define format support and conversion routes in a single configuration source, or use a table-driven approach in `Select()`.

**Priority:** LOW - Works but creates maintenance burden.

---

## Summary of Blockers

| Issue | Severity | Blocker | Details |
|-------|----------|---------|---------|
| Config struct missing fields | CRITICAL | Yes | Cannot compile |
| README promises unimplemented features | HIGH | No | User confusion |
| XML->YAML conversion missing | MEDIUM | No | Partial feature |
| Deferred close error handling | MEDIUM | No | Data integrity risk |
| No tests | HIGH | No | No regression detection |
| Single-file XML mode unsupported | MEDIUM | No | Usability limitation |

The **critical blocker must be fixed before any other work** — the project cannot currently be compiled.

---

*Concerns audit: 2026-06-06*
