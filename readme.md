# gofpdi

[![Go Report Card](https://goreportcard.com/badge/github.com/ArExi0n/gofpdi)](https://goreportcard.com/report/github.com/ArExi0n/gofpdi)
[![GoDoc](https://godoc.org/github.com/ArExi0n/gofpdi?status.svg)](https://godoc.org/github.com/ArExi0n/gofpdi)
[![License](https://img.shields.io/github/license/ArExi0n/gofpdi.svg)](LICENSE)

A Go port of [FPDI](https://www.setasign.com/products/fpdi/about/) — a PDF importer library that lets your PDF generation code use pages from existing PDF files as templates.

## Features

- Import pages from an existing PDF file or an `io.ReadSeeker` stream
- Retrieve page sizes and boxes (MediaBox, CropBox, etc.) of the source PDF
- Use imported pages as templates at any position, width, and height
- Emit imported objects as form XObjects for embedding into a generated PDF
- Ordered (integer object IDs) and unordered (SHA-1 hash) output modes
- Handles compressed object streams and classic xref tables

## Installation

```bash
go get github.com/ArExi0n/gofpdi
```

## Usage

```go
package main

import (
    "github.com/ArExi0n/gofpdi/gofpdi"
)

func main() {
    imp := gofpdi.NewImporter()

    // Load the source PDF
    imp.SetSourceFile("template.pdf")

    // Get page sizes (page -> box -> "w"/"h")
    sizes := imp.GetPageSizes()

    // Import page 1 with its media box as a template
    tplID := imp.ImportPage(1, "/MediaBox")

    // Grab the objects to embed in your generated PDF
    formXObjects := imp.PutFormXobjects()

    // Get the template name and draw values for rendering
    tplName, scaleX, scaleY, tx, ty := imp.UseTemplate(tplID, 0, 0, 210, 297)

    _ = sizes
    _ = formXObjects
    _, _, _, _, _ = tplName, scaleX, scaleY, tx, ty
}
```

See [gofpdf](https://github.com/ArExi0n/gofpdf) (or any FPDI-compatible generator) for a full end-to-end example of embedding imported pages into a new PDF.

## API Overview

### Importer

| Method | Description |
| --- | --- |
| `NewImporter()` | Create a new importer instance |
| `SetSourceFile(path)` | Load the source PDF from a file |
| `SetSourceStream(rs)` | Load the source PDF from an `io.ReadSeeker` |
| `GetPageSizes()` | Return all page boxes as `map[int]map[string]map[string]float64` |
| `ImportPage(pageno, box)` | Import a page and return a template ID |
| `SetNextObjectID(id)` | Offset the next template object ID |
| `PutFormXobjects()` | Return form XObjects as `map[string]int` (name → object ID) |
| `PutFormXobjectsUnordered()` | Return form XObjects as `map[string]string` (name → SHA-1 hash) |
| `GetImportedObjects()` | Return imported objects as `map[int]string` (ID → content) |
| `GetImportedObjectsUnordered()` | Return imported objects as `map[string][]byte` (hash → content) |
| `GetImportedObjHashPos()` | Return positions of hashes within objects for ID replacement |
| `UseTemplate(tplID, x, y, w, h)` | Return the template name and the 4 values needed to draw it |

### Reader / Writer

| Type | Description |
| --- | --- |
| `PDFreader` | Low-level PDF parser — resolves objects, compressed streams, and xref tables |
| `PdfWriter` | Serializes imported pages into templates and form XObjects |

## Project Structure

```
gofpdi/
├── gofpdi/       # Library package
│   ├── importer.go   # High-level Importer API
│   ├── reader.go     # PDF parsing and object resolution
│   ├── writer.go     # Template and form XObject writing
│   ├── helper.go     # Utility helpers (PNG filter, arrays, etc.)
│   └── const.go      # Constants
└── pkg/
    └── errors/       # Error helpers (github.com/pkg/errors)
```

## License

[MIT](LICENSE)
