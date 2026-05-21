# GoConvert 🚀

GoConvert is a high-performance, lightweight, and architecture-first CLI tool written in Go. It seamlessly converts data and configuration files between various formats (**JSON, YAML, CSV, and Markdown Tables**).

Unlike static conversion tools, GoConvert is built on a **Strategy-Based Hybrid Architecture**. It dynamically evaluates the input/output formats and selects the most optimal data pipeline (Zero-Memory Streaming vs. Structural Tree Parsing) at runtime.

---

## 📐 Hybrid System Architecture

GoConvert implements the **Strategy and Factory Design Patterns**. The core engine does not enforce a single global data structure. Instead, it decouples the execution into specialized processing pipelines based on the nature of the data:



1. **Stream Pipeline (O(1) Memory):** Used for flat, sequential data (e.g., `CSV` ⇄ `Markdown Tables`). Data is read line-by-line via `io.Reader` and flushed immediately via `io.Writer` without allocating massive heap memory.
2. **Tree-Parser Pipeline:** Used for hierarchical, nested data (e.g., `JSON` ⇄ `YAML`). It constructs an in-memory Abstract Syntax Tree (AST) to securely map complex structures, arrays, and datatypes.

---

## 🛠 Features

- **Multi-Format Universe:** Initial support for `JSON`, `YAML`, `CSV`, and Github-flavored `Markdown Tables`.
- **Smart Strategy Selector:** Automatically detects the most efficient processing pipeline based on your `-from` and `-to` flags.
- **Flexible Stream Routing:** Supports reading from direct physical files or piping data on-the-fly via standard input (`stdin`). Output can be saved to a file or streamed to the console (`stdout`).
- **Production-Ready Error Handling:** Zero panic-driven design. Validates schema and structure format boundaries, returning graceful exit codes and clear developer-friendly error logs.

---

## 🚀 Target CLI UX

Compile the binary and use simple CLI flags to control the smart engine:

```bash
# Executing a Tree-Parser Pipeline (Complex structures)
goconvert -from=json -to=yaml -in=config.json -out=config.yaml

# Executing a Zero-Memory Stream Pipeline (Sequential data)
cat metrics.csv | goconvert -from=csv -to=md