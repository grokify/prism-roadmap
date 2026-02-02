# Release Notes - v0.1.0

**Release Date:** January 25, 2026

## Overview

This is the initial release of **structured-requirements**, a Go library for managing product requirements documents. It supports PRD (Product Requirements Document), MRD (Market Requirements Document), and TRD (Technical Requirements Document) formats with validation, rendering, and goals alignment.

structured-requirements integrates with:
- [structureddocs](https://github.com/grokify/structureddocs) for consistent Marp presentation rendering
- [structured-goals](https://github.com/grokify/structured-goals) for OKR and V2MOM alignment

## Highlights

- **PRD Support** - Comprehensive Product Requirements Documents with personas, objectives, and roadmaps
- **MRD Support** - Market Requirements Documents with competitive analysis
- **TRD Support** - Technical Requirements Documents with architecture specifications
- **Goals Alignment** - Native support for OKR and V2MOM goal frameworks
- **Marp Integration** - Generate presentation slides with consistent theming

## Installation

```bash
go get github.com/grokify/structured-plan
```

## Features

### PRD (Product Requirements Document)

The `prd` package provides:

- Document structure: metadata, executive summary, personas, objectives, requirements
- Functional and non-functional requirements with MoSCoW prioritization
- Success metrics with baselines and targets
- Roadmap with phases and milestones
- Risk management with mitigation strategies

### Goals Alignment

Align PRDs with strategic goals:

- **OKR Integration** - Link requirements to Objectives and Key Results
- **V2MOM Integration** - Align with Vision, Values, Methods, Obstacles, Measures

### Document Views

Generate multiple formats from a single PRD:

| View | Description |
|------|-------------|
| Executive Summary | High-level overview for stakeholders |
| PM View | Detailed view for product managers |
| PR/FAQ | Amazon-style press release format |
| Six-Pager | Detailed narrative document |
| Marp Slides | Presentation-ready slides |

### MRD (Market Requirements Document)

The `mrd` package provides:

- Market analysis and sizing
- Customer segments and personas
- Competitive landscape
- Business objectives and success criteria

### TRD (Technical Requirements Document)

The `trd` package provides:

- Architecture decisions and rationale
- Technical specifications
- Performance requirements
- Security and compliance requirements
- Scalability considerations

### Marp Slide Generation

Generate professional presentation slides:

- Problem and solution slides
- Persona overview
- Objectives and metrics
- Requirements summary
- Roadmap visualization
- Goals alignment slides

## Quick Start

```go
package main

import (
    "fmt"
    "github.com/grokify/structured-plan/prd"
    "github.com/grokify/structured-plan/prd/render/marp"
)

func main() {
    // Load PRD from file
    doc, err := prd.ReadFile("product.prd.json")
    if err != nil {
        panic(err)
    }

    // Check completeness
    score := prd.CalculateCompleteness(doc)
    fmt.Printf("PRD Completeness: %.0f%%\n", score*100)

    // Generate Marp slides
    renderer := marp.NewPRDRenderer()
    slides, err := renderer.Render(doc, nil)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(slides))
}
```

## Examples

The `examples/` directory contains sample documents:

- `agent-control-plane.prd.json` - Full-featured PRD example
- `agent-compute-plane.prd.json` - Another PRD example
- `agent-platform.mrd.json` - MRD example
- `agent-control-plane.trd.json` - TRD example

## Dependencies

- Go 1.24+
- github.com/grokify/structureddocs v0.1.0
- github.com/grokify/structured-goals v0.1.0
- github.com/spf13/cobra v1.10.2

## Contributors

- John Wang (@grokify)
- Claude Opus 4.5 (Co-Author)

## Links

- [GitHub Repository](https://github.com/grokify/structured-plan)
- [Go Package Documentation](https://pkg.go.dev/github.com/grokify/structured-plan)
- [Changelog](CHANGELOG.md)
- [structureddocs](https://github.com/grokify/structureddocs) - Shared rendering utilities
- [structured-goals](https://github.com/grokify/structured-goals) - Goals frameworks
