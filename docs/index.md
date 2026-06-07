# PRISM Roadmap

A Go library for creating, validating, and transforming structured planning documents. Part of the [PRISM ecosystem](https://github.com/grokify/prism).

## Overview

PRISM Roadmap provides typed data structures and utilities for planning documents used in product development:

### Requirements Documents

| Document | Purpose | Primary Audience |
|----------|---------|------------------|
| **PRD** | Product Requirements Document | Product Managers, Engineers |
| **MRD** | Market Requirements Document | Product Marketing, Sales |
| **TRD** | Technical Requirements Document | Engineers, Architects |

### Strategic Canvases

| Canvas | Framework | Use Case |
|--------|-----------|----------|
| **BMC** | Business Model Canvas | Business model visualization |
| **OST** | Opportunity Solution Tree | Outcome-driven discovery |
| **Opportunity** | Opportunity Canvas | Opportunity assessment |
| **Feature** | Feature Canvas | Feature definition |
| **Lean UX** | Lean UX Canvas | Hypothesis-driven design |

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

    subgraph Canvases
        CV1[BMC]
        CV2[OST]
        CV3[Opportunity]
        CV4[Feature]
        CV5[Lean UX]
    end

    CV1 --> R[Renderers]
    CV2 --> R
    CV3 --> R
    CV4 --> R
    CV5 --> R

    R --> R1[D2/SVG]
    R --> R2[Mermaid]
    R --> R3[Lit/JSON]
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
- [Strategic Canvases](canvas/overview.md)
