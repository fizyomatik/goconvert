# Testing Patterns

**Analysis Date:** 2026-06-06

## Test Framework

**Runner:**
- Go testing framework (built-in `testing` package)
- Go version: 1.23.0 (supports latest testing features)
- No third-party test framework configured

**Test File Location:**
- Co-located with source files using `_test.go` suffix
- Example: `main_test.go` in root directory alongside `main.go`

**Run Commands:**
```bash
go test ./...              # Run all tests in all packages
go test -v ./...           # Run with verbose output
go test -cover ./...       # Show test coverage percentage
go test -coverprofile=coverage.out ./...  # Generate coverage file
go tool cover -html=coverage.out          # View coverage as HTML
go test -run TestName      # Run specific test by name
```

## Test File Organization

**Location:**
- Co-located pattern: tests live in same directory as source code
- No separate `tests/` directory

**Naming:**
- Test functions: `Test` prefix followed by function name (Go convention)
- Example patterns: `TestParseCLI`, `TestJSONToYAML`, `TestSelect`
- Subtests: Use `t.Run("subtest name", func(t *testing.T) {...})` pattern
- Benchmark functions: `Benchmark` prefix (when performance testing needed)

**Structure:**
```
[package]/
├── module.go           # Source file
├── module_test.go      # Test file for this module
└── testdata/           # Test data files (optional, currently uses ../data/)
```

**Current Test Data:**
- Location: `data/` directory in project root
- Files: `test_data.json`, `dummy_data.yaml`, `book_data.xml`
- Usage: Manual testing with CLI or future test implementations
- Path: `data/test_data.json`, `data/dummy_data.yaml`, `data/book_data.xml`

## Test Structure

**Go Testing Pattern:**
```go
package cli

import (
	"testing"
)

func TestParseCLI(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    *Config
		wantErr bool
	}{
		{
			name:    "flag mode",
			args:    []string{"prog", "-from", "json", "-to", "yaml"},
			want:    &Config{From: "json", To: "yaml", IsPipe: false},
			wantErr: false,
		},
		{
			name:    "invalid format",
			args:    []string{"prog", "-from", "invalid", "-to", "json"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCLI(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCLI() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCLI() got = %v, want %v", got, tt.want)
			}
		})
	}
}
```

**Setup/Teardown:**
- Use `testing.T` helper methods: no explicit setup/teardown needed for simple tests
- For fixtures: create helper functions to build test data
- For file operations: use `t.TempDir()` for temporary test directories (Go 1.15+)

**Patterns:**
- Table-driven tests: Define test cases as slice of structs with inputs and expected outputs
- Subtest organization: Use `t.Run()` to group related test cases
- Error assertions: Check `(err != nil) != wantErr` pattern to verify error presence
- Deep equality: Use `reflect.DeepEqual()` or similar for complex type comparison

## Mocking

**Framework:**
- No mocking library configured (testify/mock, gomock, etc. not present)
- Manual interface mocking recommended for tests

**Patterns:**
```go
// Example: Mock io.Reader for testing
type mockReader struct {
	data []byte
	err  error
}

func (m *mockReader) Read(p []byte) (n int, err error) {
	if m.err != nil {
		return 0, m.err
	}
	n = copy(p, m.data)
	return n, io.EOF
}

// Usage in test:
func TestJSONToYAMLWithError(t *testing.T) {
	conv := &JSONToYAML{}
	reader := &mockReader{err: fmt.Errorf("read error")}
	writer := &bytes.Buffer{}
	
	err := conv.Convert(reader, writer)
	if err == nil {
		t.Error("expected error, got nil")
	}
}
```

**What to Mock:**
- External I/O operations: `io.Reader`, `io.Writer`
- File system operations: use `t.TempDir()` instead of mocking
- Error scenarios: Create mock implementations that return specific errors
- Rate-limited operations: Create controlled mock implementations

**What NOT to Mock:**
- Format conversion logic: Test actual conversion implementations
- Standard library functions: Trust their behavior
- Internal helper functions: Test through public API instead

## Fixtures and Factories

**Test Data:**
```go
// Helper function to create test Config
func createTestConfig(from, to string) *Config {
	return &Config{
		From: from,
		To: to,
		InputFile: "test.json",
		OutputFile: "test.yaml",
		IsPipe: false,
	}
}

// Helper to load test JSON data
func loadTestJSON(t *testing.T, filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Fatalf("failed to load test file: %v", err)
	}
	return data
}
```

**Location:**
- Helper functions in same `*_test.go` file
- Shared test utilities: Create `testutil.go` or similar if used across packages
- Test data files: Store in `data/` directory (currently used for manual testing)

**Naming:**
- Helpers prefixed with `helper` or `create`: `helperLoadJSON`, `createTestConfig`
- Data builders: `NewTestConfig()`, `WithJSONInput()` (fluent style if needed)

## Coverage

**Requirements:**
- No coverage target configured
- Coverage not enforced in this codebase

**View Coverage:**
```bash
# Terminal output
go test -cover ./...

# HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Console coverage details
go tool cover -func=coverage.out
```

**Coverage Gaps (Once Tests Implemented):**
- `main.go`: Entry point (should have integration tests, not unit tests)
- `cli/parser.go`: All code paths (CLI parsing logic, all flag combinations)
- `engine/engine.go`: All conversion pair selections
- `formats/*.go`: All format conversion implementations
- Error paths: Input validation, malformed data, I/O errors

## Test Types

**Unit Tests:**
- Scope: Individual functions and methods (`ParseCLI`, `isSupportedFormat`, converter implementations)
- Approach: Table-driven tests with various inputs
- Files: Located next to source (`cli/parser_test.go`, `formats/json_test.go`)
- Focus: Input validation, error handling, edge cases
- Example: Test `ParseCLI` with all CLI flag combinations and invalid inputs

**Integration Tests:**
- Scope: End-to-end conversion workflows (JSON → YAML, YAML → XML, etc.)
- Approach: Use actual file I/O or `bytes.Buffer` for in-memory testing
- Files: Can be in `main_test.go` or separate package tests
- Focus: Verify complete data conversion pipeline
- Example: Convert sample JSON to YAML and verify output structure

**E2E Tests:**
- Framework: Not formally configured
- Approach: Manual testing via CLI or shell scripts in `usage/` directory
- Example commands in `usage/commands.md`
- When formalized: Could use `os/exec` to test binary directly

## Common Patterns

**Async Testing:**
- Converters are synchronous; no goroutines used in current implementation
- If async patterns are added: Use `sync.WaitGroup` or channels with timeout detection
- Example pattern:
```go
func TestAsyncConversion(t *testing.T) {
	done := make(chan error)
	go func() {
		done <- converter.Convert(reader, writer)
	}()
	
	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("conversion failed: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("conversion timeout")
	}
}
```

**Error Testing:**
```go
func TestConvertWithInvalidJSON(t *testing.T) {
	conv := &JSONToYAML{}
	reader := strings.NewReader("{invalid json}")
	writer := &bytes.Buffer{}
	
	err := conv.Convert(reader, writer)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
	if !strings.Contains(err.Error(), "failed to parse JSON") {
		t.Errorf("unexpected error message: %v", err)
	}
}
```

**Comparing I/O Output:**
```go
func TestYAMLOutput(t *testing.T) {
	input := bytes.NewBufferString(`{"name": "test"}`)
	output := &bytes.Buffer{}
	
	conv := &JSONToYAML{}
	if err := conv.Convert(input, output); err != nil {
		t.Fatalf("conversion failed: %v", err)
	}
	
	result := output.String()
	if !strings.Contains(result, "name: test") {
		t.Errorf("output missing expected YAML: %q", result)
	}
}
```

## Test Data

**Test Files:**
- `data/test_data.json`: Small JSON sample for quick tests
- `data/dummy_data.yaml`: Large YAML file for performance testing
- `data/book_data.xml`: XML test sample

**Usage Pattern:**
```go
// In tests, load data relative to test file location
testJSONPath := filepath.Join("..","data", "test_data.json")
data, err := ioutil.ReadFile(testJSONPath)
if err != nil {
	t.Fatalf("failed to load test data: %v", err)
}
```

**Creating Test Data:**
- Keep test files small for unit tests (< 1 KB)
- Use real-world samples but minimal: actual structure, few records
- Performance test data in separate large files (dummy_data.yaml)

---

*Testing analysis: 2026-06-06*
