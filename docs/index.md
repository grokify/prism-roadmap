# PRISM Roadmap

A Go library for creating, validating, and transforming structured planning documents. Part of the [PRISM ecosystem](https://github.com/grokify/prism).

## Overview

PRISM Roadmap provides typed data structures and utilities for planning documents used in product development:

| Document | Purpose | Primary Audience |
|----------|---------|------------------|
| **PRD** | Product Requirements Document | Product Managers, Engineers |
| **MRD** | Market Requirements Document | Product Marketing, Sales |
| **TRD** | Technical Requirements Document | Engineers, Architects |

## Key Features

- **JSON format** - Machine-readable documents with defined schemas
- **Multiple output views** - PM View, Executive View, Amazon 6-Pager, PR/FAQ
- **Goals alignment** - Integrate with V2MOM and OKR frameworks
- **Scoring & validation** - Automated quality assessment
- **Persona library** - Reusable persona definitions across documents

## Architecture

```mermaid
graph TD
    A[PRD] --> B[Views]
    A --> C[Scoring]
    A --> D[Goals]

    B --> B1[PM View]
    B --> B2[Exec View]
    B --> B3[6-Pager]
    B --> B4[PR/FAQ]

    D --> D1[V2MOM]
    D --> D2[OKR]

    E[MRD] --> F[Markdown]
    G[TRD] --> H[Markdown]
```

## Quick Example

```go
package main

import (
    "fmt"
    "github.com/grokify/prism-roadmap/requirements/prd"
)

func main() {
    // Create a new PRD
    doc := prd.New("PRD-2025-001", "User Authentication System",
        prd.Person{Name: "Alice Smith", Role: "Product Manager"})

    // Set problem statement
    doc.ExecutiveSummary.ProblemStatement = "Users cannot securely access their accounts"
    doc.ExecutiveSummary.ProposedSolution = "Implement OAuth 2.0 authentication"

    // Score the PRD
    scores := prd.Score(doc)
    fmt.Printf("Overall Score: %.0f%%\n", scores.OverallScore*100)

    // Generate views
    pmView := prd.GeneratePMView(doc)
    markdown := prd.RenderPMMarkdown(pmView)
    fmt.Println(markdown)
}
```

## Installation

```bash
go get github.com/grokify/prism-roadmap
```

## Document Relationships

Planning documents can reference each other and align with strategic goals:

```mermaid
graph LR
    V[V2MOM] --> P[PRD]
    O[OKR] --> P
    P --> T[TRD]
    P --> M[MRD]
    P --> |views| PM[PM View]
    P --> |views| EX[Exec View]
    P --> |views| SP[6-Pager]
```

## Next Steps

- [Installation Guide](getting-started/installation.md)
- [Quick Start Tutorial](getting-started/quickstart.md)
- [PRD Documentation](documents/prd.md)
- [Goals Integration](goals/overview.md)
