# Business Model Canvas

The Business Model Canvas (BMC) follows Alexander Osterwalder's 9-block structure for describing, designing, and analyzing business models.

## Overview

The BMC provides a visual chart with elements describing a firm's value proposition, infrastructure, customers, and finances.

## 9-Block Layout

| Key Partners | Key Activities | Value Propositions | Customer Relationships | Customer Segments |
|--------------|----------------|-------------------|----------------------|-------------------|
| Key Resources || Channels ||
| Cost Structure ||| Revenue Streams ||

## Structure

```go
type BusinessModelCanvas struct {
    Metadata             Metadata
    CustomerSegments     []CustomerSegment
    ValuePropositions    []ValueProposition
    Channels             []Channel
    CustomerRelationships []Relationship
    RevenueStreams       []RevenueStream
    KeyResources         []Resource
    KeyActivities        []Activity
    KeyPartnerships      []Partnership
    CostStructure        []Cost
}
```

## Example

```go
bmc := canvas.NewBusinessModelCanvas("bmc-1", "SaaS Platform")

bmc.CustomerSegments = []canvas.CustomerSegment{
    {ID: "cs1", Name: "Enterprise", Size: "$10B TAM"},
    {ID: "cs2", Name: "SMB", Size: "$5B TAM"},
}

bmc.ValuePropositions = []canvas.ValueProposition{
    {ID: "vp1", Description: "Unified platform for team collaboration"},
}

bmc.KeyResources = []canvas.Resource{
    {ID: "kr1", Name: "Engineering Team", Type: "human"},
    {ID: "kr2", Name: "Cloud Infrastructure", Type: "infrastructure"},
}

// Render
c := canvas.NewBMC(bmc)
d2Output, _ := render.Render(c, render.FormatD2, render.BMCOptions())
```

## Rendering

```bash
d2 bmc_example.d2 bmc_example.svg
```

## Examples

See `examples/canvas/bmc/` for complete examples.
