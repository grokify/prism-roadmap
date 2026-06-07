# Opportunity Solution Tree

The Opportunity Solution Tree (OST) follows Teresa Torres's framework for continuous product discovery, connecting outcomes to opportunities to solutions.

## Overview

The OST is a tree structure that helps product teams:

1. Start with a clear **Outcome** (measurable goal)
2. Identify **Opportunities** (customer needs/problems)
3. Generate **Solutions** for each opportunity
4. Design **Experiments** to validate solutions

## Tree Structure

```
Outcome
├── Opportunity 1
│   ├── Solution A
│   │   └── Experiment 1
│   └── Solution B
└── Opportunity 2
    └── Solution C
        └── Experiment 2
```

## Structure

```go
type OpportunitySolutionTree struct {
    Metadata Metadata
    Outcome  OSTOutcome
}

type OSTOutcome struct {
    ID            string
    Description   string
    Metric        string           // How to measure
    Target        string           // Success target
    Opportunities []OSTOpportunity
}

type OSTOpportunity struct {
    ID          string
    Description string
    Source      string        // interview, analytics, support
    Priority    int
    Solutions   []OSTSolution
}

type OSTSolution struct {
    ID          string
    Description string
    Type        string         // feature, improvement, experiment
    Status      string         // proposed, validated, building, shipped
    Experiments []OSTExperiment
}

type OSTExperiment struct {
    ID         string
    Hypothesis string
    Method     string  // prototype, A/B test, interview
    Status     string  // planned, running, completed
    Result     string
    Learning   string
}
```

## Example

```go
ost := canvas.NewOpportunitySolutionTree("ost-activation", "User Activation")

ost.Outcome = canvas.OSTOutcome{
    ID:          "O1",
    Description: "Increase user activation rate",
    Metric:      "% users completing onboarding",
    Target:      "60%",
    Opportunities: []canvas.OSTOpportunity{
        {
            ID:          "OP1",
            Description: "Users don't understand core value proposition",
            Source:      "User interviews",
            Solutions: []canvas.OSTSolution{
                {
                    ID:          "S1",
                    Description: "Interactive product tour",
                    Status:      "building",
                    Experiments: []canvas.OSTExperiment{
                        {
                            ID:         "E1",
                            Hypothesis: "Guided tour increases completion by 20%",
                            Method:     "A/B test",
                            Status:     "planned",
                        },
                    },
                },
            },
        },
    },
}

// Render
c := canvas.NewOST(ost)
d2Output, _ := render.Render(c, render.FormatD2, render.OSTOptions())
```

## D2 Output

```d2
direction: down

O1: "Increase user activation rate\nTarget: 60%"

OP1: "Users don't understand value prop"
O1 -> OP1

S1: "Interactive product tour [building]"
OP1 -> S1

E1: "Guided tour increases completion by 20%"
S1 -> E1
```

## Color Scheme

| Element | Color | Hex |
|---------|-------|-----|
| Outcome | Green | `#E8F5E9` |
| Opportunity | Blue | `#E3F2FD` |
| Solution | Orange | `#FFF3E0` |
| Experiment | Purple | `#F3E5F5` |

## Examples

See `examples/canvas/ost/` for complete examples.
