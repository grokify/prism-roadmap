# Quick Start

This guide walks you through creating your first PRD with prism-roadmap.

## Create a PRD

```go
package main

import (
    "github.com/grokify/prism-roadmap/requirements/prd"
)

func main() {
    // Create new PRD with ID, title, and author
    doc := prd.New(
        prd.GenerateID(),  // Auto-generate ID like "PRD-2025-022"
        "Customer Portal Redesign",
        prd.Person{Name: "Jane Doe", Role: "Senior PM", Email: "jane@example.com"},
    )

    // Set executive summary
    doc.ExecutiveSummary = prd.ExecutiveSummary{
        ProblemStatement: "Current portal has 40% bounce rate due to poor UX",
        ProposedSolution: "Modern, responsive redesign with improved navigation",
        ExpectedOutcomes: []string{
            "Reduce bounce rate to under 20%",
            "Increase task completion by 50%",
        },
    }

    // Add personas
    doc.Personas = []prd.Persona{
        {
            ID:        "PER-1",
            Name:      "Power User",
            Role:      "Account Manager",
            IsPrimary: true,
            Goals:     []string{"Quickly access client data", "Generate reports"},
            PainPoints: []string{"Slow page loads", "Too many clicks to complete tasks"},
        },
    }

    // Add objectives
    doc.Objectives = prd.Objectives{
        BusinessObjectives: []prd.Objective{
            {ID: "BO-1", Description: "Increase customer retention by 15%"},
        },
        ProductGoals: []prd.Objective{
            {ID: "PG-1", Description: "Achieve 90% task success rate"},
        },
        SuccessMetrics: []prd.SuccessMetric{
            {
                ID:          "SM-1",
                Name:        "Bounce Rate",
                Metric:      "bounce_rate",
                Target:      "<20%",
                CurrentBaseline: "40%",
            },
        },
    }

    // Save to file
    prd.Save(doc, "portal-redesign.prd.json")
}
```

## Score Your PRD

```go
// Score the document
scores := prd.Score(doc)

fmt.Printf("Overall Score: %.0f%%\n", scores.OverallScore*100)
fmt.Printf("Decision: %s\n", scores.Decision)

// Check category scores
for _, cat := range scores.CategoryScores {
    fmt.Printf("  %s: %.0f%%\n", cat.Category, cat.Score*100)
}
```

## Generate Views

### PM View

```go
pmView := prd.GeneratePMView(doc)
markdown := prd.RenderPMMarkdown(pmView)
fmt.Println(markdown)
```

### Executive View

```go
execView := prd.GenerateExecView(doc, scores)
markdown := prd.RenderExecMarkdown(execView)
fmt.Println(markdown)
```

### Amazon 6-Pager

```go
sixPager := prd.GenerateSixPagerView(doc)
markdown := prd.RenderSixPagerMarkdown(sixPager)
fmt.Println(markdown)
```

### PR/FAQ

```go
prfaq := prd.GeneratePRFAQView(doc)
markdown := prd.RenderPRFAQMarkdown(prfaq)
fmt.Println(markdown)
```

## Load Existing PRD

```go
doc, err := prd.Load("my-product.prd.json")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Loaded: %s (v%s)\n", doc.Metadata.Title, doc.Metadata.Version)
```

## Add Goals Alignment

```go
import (
    "github.com/grokify/prism-roadmap/goals/okr"
)

// Create OKR alignment
doc.Goals = &prd.GoalsAlignment{
    OKR: &okr.OKRDocument{
        Objectives: []okr.Objective{
            {
                Title: "Improve customer experience",
                KeyResults: []okr.KeyResult{
                    {Title: "Reduce bounce rate to <20%", Target: "20%"},
                    {Title: "Increase NPS by 20 points", Target: "+20"},
                },
            },
        },
    },
    AlignedObjectives: map[string]string{
        "BO-1": "O1",  // Map business objective to OKR objective
    },
}
```

## Next Steps

- [PRD Documentation](../documents/prd.md) - Full PRD structure reference
- [Scoring](../features/scoring.md) - Understand the scoring system
- [Views](../views/pm-view.md) - All available output formats
- [Goals Integration](../goals/overview.md) - V2MOM and OKR alignment
