# Effort Estimation

The `effort` package provides types for estimating implementation effort and tracking complexity factors.

## Overview

Effort estimation combines:

- **Person-days estimate** - Numeric effort in person-days
- **T-shirt sizing** - Quick categorical sizing (XS, S, M, L, XL)
- **Confidence levels** - How certain the estimate is
- **Complexity factors** - What makes implementation complex

## EffortEstimate Type

```go
type EffortEstimate struct {
    PersonDays int        `json:"person_days"`
    TShirtSize TShirtSize `json:"tshirt_size,omitempty"`
    Confidence Confidence `json:"confidence"`
    Note       string     `json:"note,omitempty"`
}
```

## T-Shirt Sizing

Quick categorical sizing for rough estimates:

| Size | Description | Typical Days |
|------|-------------|--------------|
| `XS` | Extra Small | 1-2 days |
| `S` | Small | 3-5 days |
| `M` | Medium | 1-2 weeks |
| `L` | Large | 2-4 weeks |
| `XL` | Extra Large | 4+ weeks |

```go
import "github.com/grokify/prism-roadmap/effort"

size := effort.TShirtSizeMedium
days := size.EstimatedDays() // Returns 7-10
```

## Confidence Levels

Indicate certainty of estimates:

| Level | Multiplier | When to Use |
|-------|------------|-------------|
| `high` | 1.0 | Well-understood work, similar past projects |
| `medium` | 1.3 | Some unknowns, reasonable assumptions |
| `low` | 1.5 | Many unknowns, exploratory work |

```go
conf := effort.ConfidenceMedium
multiplier := conf.Multiplier() // Returns 1.3
```

## Quick Start

```go
import "github.com/grokify/prism-roadmap/effort"

estimate := &effort.EffortEstimate{
    PersonDays: 15,
    TShirtSize: effort.TShirtSizeMedium,
    Confidence: effort.ConfidenceMedium,
    Note:       "Includes API and UI work",
}

// Get adjusted estimate with confidence buffer
adjusted := estimate.AdjustedDays()
// 15 × 1.3 = 19.5 days
```

## ComplexityFactors Type

Track what makes implementation complex:

```go
type ComplexityFactors struct {
    NewArchitecture bool         `json:"new_architecture"`
    NewDesignUX     bool         `json:"new_design_ux"`
    NewBillingSKU   bool         `json:"new_billing_sku"`
    Dependencies    []Dependency `json:"dependencies,omitempty"`
}

type Dependency struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Type        string `json:"type"`  // internal, external, vendor
    Team        string `json:"team,omitempty"`
    Status      string `json:"status"` // unknown, in_progress, blocked, complete
    RiskLevel   string `json:"risk_level,omitempty"`
    Description string `json:"description,omitempty"`
}
```

## Complexity Score

Calculate overall complexity:

```go
complexity := &effort.ComplexityFactors{
    NewArchitecture: true,  // +0.3
    NewDesignUX:     true,  // +0.2
    NewBillingSKU:   false, // +0.0
    Dependencies: []effort.Dependency{
        {Name: "Auth Service", Type: "internal"},
        {Name: "Payment Provider", Type: "vendor"},
    },
}

score := complexity.ComplexityScore()
// 1.0 (base) + 0.3 + 0.2 + (2 × 0.1) = 1.7
```

### Complexity Weights

| Factor | Weight | Description |
|--------|--------|-------------|
| Base | 1.0 | Starting score |
| New Architecture | +0.3 | Major architectural changes |
| New Design/UX | +0.2 | Design collaboration required |
| New Billing SKU | +0.2 | Revenue/billing changes |
| Per Dependency | +0.1 | Each external dependency |

## Dependency Types

| Type | Description | Risk Level |
|------|-------------|------------|
| `internal` | Another internal team | Lower |
| `external` | External API/service | Medium |
| `vendor` | Third-party vendor | Higher |

## Integration with RoadmapItem

```go
import (
    "github.com/grokify/prism-roadmap/rmi"
    "github.com/grokify/prism-roadmap/effort"
)

item := &rmi.RoadmapItem{
    ID:   "feature-1",
    Name: "Payment Integration",
    Effort: &effort.EffortEstimate{
        PersonDays: 20,
        TShirtSize: effort.TShirtSizeLarge,
        Confidence: effort.ConfidenceLow,
    },
    Complexity: &effort.ComplexityFactors{
        NewArchitecture: false,
        NewBillingSKU:   true,
        Dependencies: []effort.Dependency{
            {Name: "Stripe API", Type: "vendor", Status: "unknown"},
        },
    },
}

// Priority considers effort in denominator
score := item.PriorityScore()
```

## Best Practices

!!! tip "Estimation Tips"
    - Use T-shirt sizes for initial roadmap discussions
    - Add person-day estimates when planning sprints
    - Lower confidence = higher buffer needed
    - Track dependencies early to identify risks

!!! warning "Common Pitfalls"
    - Don't ignore complexity factors when estimating
    - Vendor dependencies often take longer than expected
    - New architecture work compounds effort
    - Update confidence as you learn more

## See Also

- [Market Signal](market-signal.md) - Customer demand signals
- [Roadmap Items](roadmap-items.md) - Complete RMI type
- [Prioritization](../canvas/prioritization.md) - RICE scoring with effort
