# DMAIC Integration

DMAIC (Define, Measure, Analyze, Improve, Control) is a data-driven quality improvement methodology used in Six Sigma. The `goals/dmaic` package models it as a metrics document: categories of metrics, each in a DMAIC phase, with baselines, targets, control limits, root causes, and the initiatives that improve them.

## Overview

DMAIC provides a structured approach to process improvement:

| Phase | Description | Key Activities |
|-------|-------------|----------------|
| **Define** | Identify the problem | Project charter, scope, goals |
| **Measure** | Quantify the problem | Data collection, baseline metrics |
| **Analyze** | Find root causes | Statistical analysis, root cause identification |
| **Improve** | Implement solutions | Solution design, pilot testing |
| **Control** | Sustain improvements | Control charts, monitoring, documentation |

## Document Structure

```json
{
  "$schema": "https://github.com/grokify/prism-roadmap/schema/dmaic.schema.json",
  "metadata": {
    "id": "DMAIC-2025-001",
    "name": "Manufacturing Quality Metrics",
    "owner": "Quality Manager",
    "team": "Quality Engineering",
    "period": "FY2025-Q1",
    "reviewCadence": "Weekly"
  },
  "categories": [
    {
      "id": "production",
      "name": "Production Quality",
      "metrics": [
        {
          "id": "prod-001",
          "name": "First Pass Yield",
          "baseline": 85.0,
          "current": 94.2,
          "target": 98.0,
          "unit": "%",
          "phase": "Control",
          "trendDirection": "higher_better",
          "controlLimits": {
            "ucl": 99.0,
            "lcl": 90.0,
            "centerLine": 95.0,
            "sigma": 1.5
          }
        }
      ]
    }
  ]
}
```

## Go Usage

```go
import "github.com/grokify/prism-roadmap/goals/dmaic"

// Create a new DMAIC document
doc := dmaic.New("DMAIC-2025-001", "Quality Metrics", "Quality Manager")

// Add a category with metrics
doc.Categories = append(doc.Categories, dmaic.Category{
    ID:   "production",
    Name: "Production Quality",
    Metrics: []dmaic.Metric{
        {
            ID:             "prod-001",
            Name:           "First Pass Yield",
            Baseline:       85.0,
            Current:        94.2,
            Target:         98.0,
            Unit:           "%",
            Phase:          dmaic.PhaseControl,
            TrendDirection: dmaic.TrendHigherBetter,
            ControlLimits: &dmaic.ControlLimits{
                UCL:        99.0,
                LCL:        90.0,
                CenterLine: 95.0,
                Sigma:      1.5,
            },
        },
    },
})

// Calculate overall health
health := doc.CalculateOverallHealth()
fmt.Printf("Overall Health: %.0f%%\n", health*100)
```

## Trend Directions

Trend direction determines how thresholds are interpreted:

| Direction | Description | Warning | Critical |
|-----------|-------------|---------|----------|
| `higher_better` | Higher values are better | Below warning threshold | Below critical threshold |
| `lower_better` | Lower values are better | Above warning threshold | Above critical threshold |
| `target_value` | Specific target is best | Deviation from target | Large deviation from target |

## Statistical Process Control

For metrics in the Control phase, SPC limits define acceptable variation:

```go
type ControlLimits struct {
    UCL        float64 // Upper Control Limit
    LCL        float64 // Lower Control Limit
    CenterLine float64 // Target/mean value
    Sigma      float64 // Standard deviation
}
```

A metric is "in control" when:

```
LCL ≤ Current Value ≤ UCL
```

## Process Capability

Six Sigma process capability metrics:

```go
type ProcessCapability struct {
    Cp         float64 // Process Capability
    Cpk        float64 // Process Capability Index
    SigmaLevel float64 // Sigma level (1-6)
    DPMO       float64 // Defects Per Million Opportunities
}
```

| Sigma Level | DPMO | Yield |
|-------------|------|-------|
| 1σ | 691,462 | 30.9% |
| 2σ | 308,538 | 69.1% |
| 3σ | 66,807 | 93.3% |
| 4σ | 6,210 | 99.4% |
| 5σ | 233 | 99.98% |
| 6σ | 3.4 | 99.9997% |

## Root Cause Analysis

Track root causes during the Analyze phase:

```go
type RootCause struct {
    ID          string
    Description string
    Category    string  // People, Process, Technology, etc.
    Impact      string  // High, Medium, Low
    Validated   bool
}
```

## Initiatives

Link improvement initiatives to metrics:

```go
type Initiative struct {
    ID             string
    Name           string
    Description    string
    Owner          string
    Status         string
    StartDate      string
    EndDate        string
    MetricIDs      []string // Linked metrics
    ExpectedImpact string
    ActualImpact   string
}
```

## Validation Options

Three validation presets are available:

| Preset | Use Case |
|--------|----------|
| `DefaultValidationOptions()` | Permissive, good for drafts |
| `StrictValidationOptions()` | Production-ready documents |
| `SixSigmaValidationOptions()` | Full Six Sigma compliance |

## Rendering

The `render/marp` package writes Marp slides and `render/pandoc` writes Pandoc
Markdown for PDF or Word:

```go
import (
    "github.com/grokify/prism-roadmap/goals/dmaic/render"
    "github.com/grokify/prism-roadmap/goals/dmaic/render/marp"
)

slides, err := marp.New().Render(doc, render.DefaultOptions())
```

## Examples

- [`examples/manufacturing.dmaic.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/manufacturing.dmaic.json)
- [`examples/software-development.dmaic.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/software-development.dmaic.json)
- [`examples/customer-service.dmaic.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/customer-service.dmaic.json)

The schema is `schema/dmaic.schema.json`, embedded as `schema.DMAICSchema()`.
