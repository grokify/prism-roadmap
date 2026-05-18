# Installation

## Requirements

- Go 1.21 or later

## Install the Library

```bash
go get github.com/grokify/prism-execution
```

## Install the CLI (Optional)

The `srequirements` CLI tool provides commands for creating and validating documents:

```bash
go install github.com/grokify/prism-execution/cmd/srequirements@latest
```

## Verify Installation

```go
package main

import (
    "fmt"
    "github.com/grokify/prism-execution/requirements/prd"
)

func main() {
    doc := prd.New("TEST-001", "Test Document")
    fmt.Printf("Created: %s\n", doc.Metadata.Title)
}
```

## Package Structure

```
github.com/grokify/prism-execution/
├── prd/          # Product Requirements Document
├── mrd/          # Market Requirements Document
├── trd/          # Technical Requirements Document
└── cmd/srequirements/   # CLI tool
```

## Import Paths

```go
import (
    "github.com/grokify/prism-execution/requirements/prd"
    "github.com/grokify/prism-execution/requirements/mrd"
    "github.com/grokify/prism-execution/requirements/trd"
)
```

## Goals Integration (Optional)

To use V2MOM and OKR integration, also install:

```bash
go get github.com/grokify/prism-execution/goals
```

The PRD package automatically includes goals types when prism-execution/goals is available.
