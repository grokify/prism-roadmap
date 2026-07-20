# Market Signal

The `signal` package provides types for aggregating customer demand signals from product ideas.

## Overview

Market signals aggregate customer demand data to help prioritize features based on real customer input:

- Vote counts from idea portals
- Customer count (how many unique customers requested)
- ARR impact (revenue at risk or opportunity)
- Source tracking (which systems contributed data)

## MarketSignal Type

```go
type MarketSignal struct {
    TotalVotes    int      `json:"total_votes"`
    CustomerCount int      `json:"customer_count"`
    TotalARR      int64    `json:"total_arr"`       // in cents
    Score         float64  `json:"score"`
    IdeaCount     int      `json:"idea_count"`
    Sources       []string `json:"sources,omitempty"`
}
```

## Signal Score Formula

The score combines multiple signals with configurable weights:

```
Score = (TotalVotes × 0.1) + (CustomerCount × 1.0) + (TotalARR / 10,000,000)
```

This formula:

- Weights customer count highly (1.0 per customer)
- Normalizes votes (0.1 per vote to prevent gaming)
- Normalizes ARR to a reasonable scale ($10M = 1.0 points)

## Quick Start

```go
import "github.com/grokify/prism-roadmap/signal"

// Create a new market signal
sig := &signal.MarketSignal{
    TotalVotes:    150,
    CustomerCount: 25,
    TotalARR:      500000000, // $5M in cents
    IdeaCount:     3,
    Sources:       []string{"aha", "productboard"},
}

// Calculate the score
score := sig.CalculateScore()
// Score = (150 × 0.1) + (25 × 1.0) + (5M / 10M)
// Score = 15 + 25 + 0.5 = 40.5
```

## Merging Signals

Combine signals from multiple sources:

```go
sig1 := &signal.MarketSignal{
    TotalVotes:    100,
    CustomerCount: 20,
    TotalARR:      300000000,
    Sources:       []string{"aha"},
}

sig2 := &signal.MarketSignal{
    TotalVotes:    50,
    CustomerCount: 10,
    TotalARR:      200000000,
    Sources:       []string{"productboard"},
}

// Merge sig2 into sig1
sig1.Merge(sig2)
// sig1 now has:
// - TotalVotes: 150
// - CustomerCount: 30
// - TotalARR: 500000000
// - Sources: ["aha", "productboard"]
```

## Integration with RoadmapItem

Market signals attach to roadmap items for prioritization:

```go
import (
    "github.com/grokify/prism-roadmap/rmi"
    "github.com/grokify/prism-roadmap/signal"
)

item := &rmi.RoadmapItem{
    ID:   "feature-1",
    Name: "Bulk Export",
    MarketSignal: &signal.MarketSignal{
        TotalVotes:    150,
        CustomerCount: 25,
        TotalARR:      500000000,
    },
}

// Priority score includes market signal
score := item.PriorityScore()
```

## Data Sources

Common sources for market signals:

| Source | Description |
|--------|-------------|
| `aha` | Aha! Ideas Portal |
| `productboard` | Productboard feedback |
| `intercom` | Intercom conversations |
| `zendesk` | Zendesk tickets |
| `salesforce` | Salesforce opportunities |
| `manual` | Manually entered data |

## Best Practices

!!! tip "Signal Quality"
    - Track unique customers, not just votes
    - Include ARR for revenue-weighted prioritization
    - Merge duplicates before calculating scores
    - Document source systems for traceability

!!! warning "Avoid Gaming"
    - The 0.1 vote weight reduces impact of vote manipulation
    - Customer count is harder to game than votes
    - ARR ties to real revenue data

## See Also

- [Roadmap Items](roadmap-items.md) - Using signals in RMIs
- [Effort Estimation](effort-estimation.md) - Estimating implementation cost
- [Prioritization Frameworks](../canvas/prioritization.md) - RICE, Kano, MoSCoW
