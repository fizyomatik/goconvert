## 🚀 Dynamic CLI UX & Edge Cases (The 5 Rules of GoConvert)

GoConvert doesn't force a single way to interact with the CLI. It dynamically parses flags, positional arguments, and system pipes (`stdin`) to evaluate your intent.

Here are the 5 execution modes our engine supports:

### 1. Explicit Mode (Fully Documented)
```bash
goconvert -from=json -to=yaml config.json config.yaml

### 2. Standard Mode
goconvert -from config.json -to config.md

### 3. Smart Extension Auto-Detection (Zero-Flag)
goconvert config.json convertedConfig.xml

### 4. Smart Default Output Mode
goconvert config.json

### 5. Pipe Stream Mode (POSIX Standard)
cat README.md | goconvert
```

