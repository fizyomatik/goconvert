# goconvert

A simple command-line tool written in Go that converts data files between formats.

## Supported Formats

| From | To |
|------|----|
| JSON | YAML |
| JSON | XML |
| YAML | JSON |
| YAML | XML |
| XML  | JSON |

This is just the beginning...

## Quickstart
```bash
go install github.com/fizyomatik/goconvert@latest
```

## Usage

There are four ways to use goconvert:

### 1. Just pass a file — it picks the output format automatically

```bash
goconvert config.json        # outputs config.yaml
goconvert  config.yaml        # outputs config.json
```

### 2. Pass input and output files

```bash
goconvert config.json config.xml
```

### 3. Use flags for full control

```bash
goconvert -from=json -to=yaml -in="config.json" -out="config.yaml"
```

### 4. Pipe from stdin

```powershell
Get-Content config.json | goconvert.exe -from=json -to=xml -out="config.xml"
```

If no `-out` flag is given, the result is printed to the terminal.

## License

See [LICENSE](LICENSE).
