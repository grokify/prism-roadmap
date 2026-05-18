# Marp Slides

Generate presentation slides from PRD documents using [Marp](https://marp.app/).

## Overview

The slides package converts PRD documents into Marp-compatible markdown that can be rendered as presentation slides. Three renderers are available:

| Renderer | Description | Use Case |
|----------|-------------|----------|
| `PRDRenderer` | Basic PRD slides | Standard PRD presentations |
| `PRDGoalsRenderer` | PRD with expanded goals | PRDs with V2MOM/OKR alignment |
| `OKRRenderer` | Standalone OKR slides | Goal-focused presentations |

## Quick Start

### PRD Slides

```go
import (
    "github.com/grokify/prism-execution/requirements/prd"
    "github.com/grokify/prism-execution/requirements/prd/render"
    "github.com/grokify/prism-execution/requirements/prd/render/marp"
)

// Load PRD
doc, _ := prd.Load("my-prd.json")

// Create renderer
renderer := marp.NewPRDRenderer()

// Generate slides
slides, err := renderer.Render(doc, nil)
if err != nil {
    log.Fatal(err)
}

// Write to file
os.WriteFile("presentation.md", slides, 0600)
```

### PRD with Goals Alignment

```go
// Create PRD+Goals renderer
renderer := marp.NewPRDGoalsRenderer()

// Generate slides with V2MOM and OKR sections
slides, err := renderer.Render(doc, nil)
```

### OKR Slides (from prism-execution/goals)

```go
import (
    "github.com/grokify/prism-execution/goals/okr"
    "github.com/grokify/prism-execution/goals/okr/render"
    okrmarp "github.com/grokify/prism-execution/goals/okr/render/marp"
)

// Load OKR
doc, _ := okr.ReadFile("okr.json")

// Create renderer
renderer := okrmarp.New()

// Generate slides
slides, _ := renderer.Render(doc, nil)
```

## Render Options

### PRD Options

```go
opts := &render.Options{
    Theme:               "corporate", // "default", "corporate", "minimal"
    IncludeGoals:        true,        // Include goals alignment section
    IncludeRoadmap:      true,        // Include roadmap slide
    IncludeRisks:        true,        // Include risks slide
    IncludeRequirements: true,        // Include requirements slide
    MaxPersonas:         5,           // Limit personas shown
    MaxRequirements:     10,          // Limit requirements shown
}

slides, _ := renderer.Render(doc, opts)
```

### OKR Options

```go
import okrrender "github.com/grokify/prism-execution/goals/okr/render"

opts := &okrrender.Options{
    Theme:            "default",  // "default", "corporate", "minimal"
    IncludeRisks:     true,       // Include risks slide
    IncludeStatus:    true,       // Show status indicators
    ShowScoreGrades:  true,       // Show letter grades (A, B, C)
    ShowProgressBars: true,       // Show visual progress bars
    MaxObjectives:    0,          // Limit objectives (0 = all)
    MaxKeyResults:    0,          // Limit KRs per objective
}
```

## Themes

Three built-in themes are available:

### Default

```go
opts := &render.Options{Theme: "default"}
```

- Primary: Indigo (#4c51bf / #5a67d8)
- Best for: General presentations

### Corporate

```go
opts := &render.Options{Theme: "corporate"}
```

- Primary: Navy blue (#1a365d)
- Best for: Executive presentations

### Minimal

```go
opts := &render.Options{Theme: "minimal"}
```

- Primary: Dark gray (#2d3748)
- Best for: Developer presentations

## Slide Structure

### PRD Slides

1. **Title Slide** - PRD title, author, version, status
2. **Problem Slide** - Problem statement and impact
3. **Solution Slide** - Proposed solution and outcomes
4. **Personas Slide** - Target users and their pain points
5. **Objectives Slide** - Business and product objectives
6. **Metrics Slide** - Success metrics with targets
7. **Requirements Slide** - Key functional/non-functional requirements
8. **Roadmap Slide** - Implementation phases
9. **Risks Slide** - Risk assessment and mitigations
10. **Goals Slide** - V2MOM/OKR alignment (if present)
11. **Summary Slide** - Key takeaways

### PRD+Goals Slides (Additional)

- **Agenda Slide** - Overview of presentation sections
- **V2MOM Vision Slide** - Vision and values
- **V2MOM Methods Slide** - Methods and obstacles
- **OKR Overview Slide** - Objectives summary
- **OKR Key Results Slide** - All key results
- **Alignment Summary Slide** - How PRD maps to goals

### OKR Slides

1. **Title Slide** - OKR name, owner, period
2. **Overview Slide** - Theme and overall progress
3. **Objective Slides** - One per objective with KRs
4. **Key Results Summary** - All KRs in table format
5. **Risks Slide** - Challenges and mitigations
6. **Summary Slide** - Final overview

## Converting to PDF/PPTX

The generated markdown can be converted using Marp CLI:

```bash
# Install Marp CLI
npm install -g @marp-team/marp-cli

# Convert to PDF
marp presentation.md -o presentation.pdf

# Convert to PowerPoint
marp presentation.md -o presentation.pptx

# Convert to HTML
marp presentation.md -o presentation.html
```

## Custom CSS

Add custom styling via the `CustomCSS` option:

```go
opts := &render.Options{
    CustomCSS: `
        section {
            background-color: #f0f0f0;
        }
        h1 {
            color: #2563eb;
        }
    `,
}
```

## Example: Complete Workflow

```go
package main

import (
    "log"
    "os"

    "github.com/grokify/prism-execution/requirements/prd"
    "github.com/grokify/prism-execution/requirements/prd/render"
    "github.com/grokify/prism-execution/requirements/prd/render/marp"
)

func main() {
    // Load PRD
    doc, err := prd.Load("customer-portal.prd.json")
    if err != nil {
        log.Fatal(err)
    }

    // Configure rendering
    opts := &render.Options{
        Theme:               "corporate",
        IncludeGoals:        true,
        IncludeRoadmap:      true,
        IncludeRisks:        true,
        IncludeRequirements: true,
    }

    // Choose renderer based on content
    var slides []byte
    if doc.Goals != nil && (doc.Goals.V2MOM != nil || doc.Goals.OKR != nil) {
        // Use PRD+Goals renderer for PRDs with goals
        renderer := marp.NewPRDGoalsRenderer()
        slides, err = renderer.Render(doc, opts)
    } else {
        // Use basic PRD renderer
        renderer := marp.NewPRDRenderer()
        slides, err = renderer.Render(doc, opts)
    }

    if err != nil {
        log.Fatal(err)
    }

    // Save markdown
    if err := os.WriteFile("presentation.md", slides, 0600); err != nil {
        log.Fatal(err)
    }

    log.Println("Slides generated: presentation.md")
    log.Println("Convert to PDF: marp presentation.md -o presentation.pdf")
}
```

## Next Steps

- [PRD Documentation](../documents/prd.md)
- [Goals Integration](../goals/overview.md)
- [V2MOM Integration](../goals/v2mom.md)
- [OKR Integration](../goals/okr.md)
