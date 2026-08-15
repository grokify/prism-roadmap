# Roadmap Items (RMI)

The `rmi` package provides the `RoadmapItem` type for roadmap planning with integrated prioritization, effort estimation, and cross-module references.

## Overview

A Roadmap Item (RMI) combines:

- **Prioritization** - MoSCoW and RICE scoring
- **Market Signal** - Customer demand data
- **Effort & Complexity** - Implementation estimates
- **Cross-References** - Links to capabilities, goals, and ideas

## RoadmapItem Type

```go
type RoadmapItem struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`

    // Prioritization. MoSCoW is optional (v0.17.0+): empty means "not yet
    // prioritized" (e.g. items imported from an external PM tool before
    // triage). When set, it must be a valid MoSCoWPriority value.
    MoSCoW prioritization.MoSCoWPriority `json:"moscow,omitempty"`
    RICE   *prioritization.RICEScore     `json:"rice,omitempty"`

    // Customer Demand
    MarketSignal *signal.MarketSignal `json:"market_signal,omitempty"`

    // Implementation
    Effort     *effort.EffortEstimate   `json:"effort,omitempty"`
    Complexity *effort.ComplexityFactors `json:"complexity,omitempty"`

    // Cross-Module References
    CapabilityRefs []string `json:"capability_refs,omitempty"`
    GoalRefs       []string `json:"goal_refs,omitempty"`
    IdeaRefs       []string `json:"idea_refs,omitempty"`

    // Lifecycle
    Status    RMIStatus `json:"status"`
    CreatedAt time.Time `json:"created_at,omitempty"`
    UpdatedAt time.Time `json:"updated_at,omitempty"`

    // Metadata
    Tags     []string          `json:"tags,omitempty"`
    Metadata map[string]string `json:"metadata,omitempty"`
}
```

## RMI Status

| Status | Description |
|--------|-------------|
| `proposed` | Initial state, under consideration |
| `approved` | Approved for roadmap |
| `in_progress` | Currently being implemented |
| `completed` | Successfully delivered |
| `deferred` | Postponed to future release |
| `cancelled` | Will not be implemented |

## Quick Start

```go
import (
    "github.com/grokify/prism-roadmap/rmi"
    "github.com/grokify/prism-roadmap/prioritization"
    "github.com/grokify/prism-roadmap/signal"
    "github.com/grokify/prism-roadmap/effort"
)

item := &rmi.RoadmapItem{
    ID:          "rmi-bulk-export",
    Name:        "Bulk Export Feature",
    Description: "Allow users to export data in bulk",

    MoSCoW: prioritization.MoSCoWShouldHave,

    MarketSignal: &signal.MarketSignal{
        TotalVotes:    150,
        CustomerCount: 25,
        TotalARR:      500000000, // $5M
    },

    Effort: &effort.EffortEstimate{
        PersonDays: 15,
        TShirtSize: effort.TShirtSizeMedium,
        Confidence: effort.ConfidenceMedium,
    },

    CapabilityRefs: []string{"cap-data-export"},
    GoalRefs:       []string{"goal-q2-efficiency"},
    IdeaRefs:       []string{"idea-123", "idea-456"},

    Status: rmi.RMIStatusApproved,
}
```

## Priority Score

Calculate combined priority score:

```go
score := item.PriorityScore()
```

The score combines multiple factors:

```
PriorityScore = MoSCoW.Weight() × MarketSignal.Score / (Effort.PersonDays × Complexity.Score)
```

## RoadmapItemSet

Manage collections of roadmap items:

```go
set := rmi.NewRoadmapItemSet()

set.Add(item1)
set.Add(item2)
set.Add(item3)

// Get items by status
approved := set.ByStatus(rmi.RMIStatusApproved)

// Get items by MoSCoW priority
mustHaves := set.ByMoSCoW(prioritization.MoSCoWMustHave)

// Sort by priority score
set.SortByPriority()
topItems := set.TopN(5)

// Get items linked to a capability
capItems := set.ByCapability("cap-data-export")
```

## Cross-Module References

### Capability References

Link to prism-capability entities:

```go
item.CapabilityRefs = []string{
    "cap-data-export",
    "cap-reporting",
}
```

### Goal References

Link to prism-maturity goals:

```go
item.GoalRefs = []string{
    "goal-q2-efficiency",
    "goal-revenue-growth",
}
```

### Idea References

Link to ProductContext ideas:

```go
item.IdeaRefs = []string{
    "idea-aha-123",
    "idea-pb-456",
}
```

## JSON Example

```json
{
  "id": "rmi-bulk-export",
  "name": "Bulk Export Feature",
  "description": "Allow users to export data in bulk",
  "moscow": "should_have",
  "market_signal": {
    "total_votes": 150,
    "customer_count": 25,
    "total_arr": 500000000,
    "score": 40.5
  },
  "effort": {
    "person_days": 15,
    "tshirt_size": "M",
    "confidence": "medium"
  },
  "complexity": {
    "new_architecture": false,
    "new_design_ux": true,
    "new_billing_sku": false,
    "dependencies": []
  },
  "capability_refs": ["cap-data-export"],
  "goal_refs": ["goal-q2-efficiency"],
  "idea_refs": ["idea-123", "idea-456"],
  "status": "approved",
  "tags": ["q2", "data"]
}
```

## Integration Points

### With prism-capability

```go
// RMI references capability
item.CapabilityRefs = []string{"cap-reporting"}

// Capability defines what we're building toward
capability := &prismcap.Capability{
    ID:   "cap-reporting",
    Name: "Reporting Capability",
}
```

### With prism-maturity

```go
// RMI references maturity goal
item.GoalRefs = []string{"goal-ops-efficiency"}

// Goal defines why we're building it
goal := &maturity.Goal{
    ID:       "goal-ops-efficiency",
    Category: "efficiency",
}
```

### With ProductContext

```go
// RMI references customer ideas
item.IdeaRefs = []string{"idea-aha-123"}

// Idea captures customer demand
idea := &types.Idea{
    ID:     "idea-aha-123",
    Source: types.IdeaSourceAha,
    Votes:  150,
}
```

## Best Practices

!!! tip "RMI Design"
    - Set MoSCoW for release planning; items imported from external PM
      tools (e.g. via [omniroadmap](https://github.com/grokify/omniroadmap))
      may arrive unprioritized — empty MoSCoW is valid and means
      "not yet triaged"
    - Include market signal for customer-driven items
    - Link to capabilities to show strategic alignment
    - Track status changes for roadmap visibility

!!! warning "Reference Integrity"
    - Validate capability/goal/idea refs exist
    - Update refs when source entities change
    - Use consistent ID formats across systems

## See Also

- [Market Signal](market-signal.md) - Customer demand aggregation
- [Effort Estimation](effort-estimation.md) - Effort and complexity
- [Prioritization](../canvas/prioritization.md) - MoSCoW, RICE, Kano
