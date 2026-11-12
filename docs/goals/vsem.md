# VSEM Integration

VSEM (Vision, Strategy, Execution, Metrics) puts a plan on one page in four sections. The `goals/vsem` package models it as a document: a vision, the strategies that will reach it, the work that carries out each strategy, and the metrics that say whether it is working, each with the person accountable for it.

## The Four Sections

| Section | Question | In the document |
|---------|----------|-----------------|
| **Vision** | Where do we want to end up? | `vision`, or `mission` when a team states why it exists instead |
| **Strategy** | How will we get there? | `strategies`, each with the `resources` it commits |
| **Execution** | What will we do, who does it, and when? | `execution`: initiatives, programs, and actions, each naming its strategy |
| **Metrics** | How will we know, and who answers for it? | `metrics`, each with an accountable `owner`, baseline, and target |

### The form this package follows

VSEM has no single canonical source, and published versions differ. Some write a mission in place of the vision, and some describe the fourth section as accountability rather than measurement. This package follows the form above and keeps both readings:

- A VSEM states a vision, a mission, or both.
- Every metric names an accountable owner. The strict options make that an error when it is missing.

### How it differs from V2MOM

| | VSEM | V2MOM |
|---|------|-------|
| Values | Not part of the plan | A section of their own |
| Obstacles | Not part of the plan | A section of their own |
| Resources | Each strategy states what it commits | Not recorded |
| Work | Execution items, apart from the strategy they serve | Projects linked to methods |
| Accountability | An owner on every metric | Owners on methods |

Use VSEM when the plan's job is to say where effort and money go and who answers for each number. Use [V2MOM](v2mom.md) when values and obstacles need to be part of the conversation.

## Document Structure

```json
{
  "$schema": "https://github.com/grokify/prism-roadmap/schema/vsem.schema.json",
  "metadata": {
    "id": "bookkeeping-fy2027",
    "name": "Bookkeeping App FY2027",
    "team": "Product",
    "period": "FY2027",
    "status": "Active"
  },
  "vision": "Every small business closes its books in an hour a month.",
  "strategies": [
    {
      "id": "onboarding",
      "name": "Get new accounts to a reconciled month in their first week",
      "owner": "Head of Onboarding",
      "priority": "P0",
      "resources": [
        { "name": "Engineering", "allocation": "2 squads" }
      ]
    }
  ],
  "execution": [
    {
      "id": "guided-import",
      "name": "Guided bank import",
      "strategyId": "onboarding",
      "kind": "initiative",
      "owner": "Onboarding squad",
      "startDate": "2027-01-04",
      "endDate": "2027-03-26"
    }
  ],
  "metrics": [
    {
      "id": "first-session-connect",
      "name": "New accounts that connect a bank in their first session",
      "strategyId": "onboarding",
      "owner": "Head of Onboarding",
      "baseline": "45",
      "target": "70",
      "unit": "%",
      "dueDate": "2027-06-30"
    },
    {
      "id": "active-90",
      "name": "Accounts still active at 90 days",
      "owner": "Head of Product",
      "baseline": "52",
      "target": "60",
      "unit": "%"
    }
  ]
}
```

Execution items and metrics refer to their strategy by `strategyId`. A metric without one measures the vision as a whole, like `active-90` above.

## Go Usage

```go
import "github.com/grokify/prism-roadmap/goals/vsem"

v := vsem.New("bookkeeping-fy2027", "Bookkeeping App FY2027",
    "Every small business closes its books in an hour a month.")

v.Strategies = append(v.Strategies, vsem.Strategy{
    ID:        "onboarding",
    Name:      "Get new accounts to a reconciled month in their first week",
    Owner:     "Head of Onboarding",
    Resources: []vsem.Resource{{Name: "Engineering", Allocation: "2 squads"}},
})
v.Execution = append(v.Execution, vsem.Execution{
    ID: "guided-import", Name: "Guided bank import",
    StrategyID: "onboarding", Kind: vsem.KindInitiative, Owner: "Onboarding squad",
})
v.Metrics = append(v.Metrics, vsem.Metric{
    ID: "first-session-connect", Name: "New accounts that connect a bank in their first session",
    StrategyID: "onboarding", Owner: "Head of Onboarding",
    Baseline: "45", Target: "70", Unit: "%",
})

for _, s := range v.Strategies {
    work := v.ExecutionFor(s.ID)
    metrics := v.MetricsFor(s.ID)
    fmt.Printf("%s: %d items of work, %d metrics\n", s.Name, len(work), len(metrics))
}
```

`VisionMetrics` returns the metrics that measure the vision as a whole, and `Statement` returns the vision, or the mission when there is no vision.

## Validation

`Validate` returns errors and warnings in document order. Structural problems are always errors:

- no vision or mission, or no strategies
- a missing ID or name, or an ID used twice anywhere in the document
- an execution item without a strategy, or one that names an unknown strategy or depends on an unknown item
- a metric that names an unknown strategy
- an unknown status, kind, priority, or direction, a progress outside 0 to 1, or a date that is not `YYYY-MM-DD`
- an execution item that ends before it starts

Gaps in ownership and coverage are warnings under `DefaultValidationOptions()` and errors under `StrictValidationOptions()`:

| Gap | Strict option |
|-----|---------------|
| A metric without an accountable owner | `RequireMetricOwner` |
| A metric without a target | `RequireMetricTarget` |
| An execution item without an owner | `RequireExecutionOwner` |
| A strategy with no execution item | `RequireStrategyExecution` |
| A strategy with no metric | `RequireStrategyMetric` |
| A strategy that commits no resources | `RequireResources` |

The strict options also warn when there are more than five strategies, more than fit on one page.

```go
errs := v.Validate(vsem.StrictValidationOptions())
if !vsem.IsValid(errs) {
    for _, e := range vsem.Errors(errs) {
        fmt.Println(e)
    }
}
```

## Goals Alignment

A PRD or roadmap can carry a VSEM through the framework-agnostic `goals` package. Strategies become goal items and metrics become result items:

```go
g := goals.NewVSEM(v)
g.GoalLabel()   // "Strategies"
g.ResultLabel() // "Metrics"
for _, r := range g.ResultItems() {
    fmt.Println(r.GoalID, r.Title, r.Target) // GoalID is empty for a vision metric
}
```

## Rendering

The `render/markdown` package writes a Markdown document, and `render/marp` writes Marp slides: a title slide, the plan on one page, the vision, a slide per strategy, and the vision's metrics. Both are deterministic.

```go
import (
    "github.com/grokify/prism-roadmap/goals/vsem/render"
    "github.com/grokify/prism-roadmap/goals/vsem/render/markdown"
    "github.com/grokify/prism-roadmap/goals/vsem/render/marp"
)

doc, err := markdown.New().Render(v, render.DefaultOptions())
slides, err := marp.New().Render(v, render.OnePageOptions())
```

`DefaultOptions` includes descriptions and resources; `OnePageOptions` leaves descriptions out.

To embed the Markdown document in a larger one, set `OmitHeader` to leave out the title and the metadata line under it, and `HeadingOffset` to push every heading down that many levels, capped at level 6; fenced code blocks are left alone. With `OmitHeader` and a `HeadingOffset` of 1, the sections that were level-two headings become level three. The Marp renderer ignores both options.

## Command Line

```bash
splan goals vsem init --name "FY2027 plan" -o plan.vsem.json
splan goals vsem validate plan.vsem.json --strict
splan goals vsem generate markdown plan.vsem.json -o plan.md
splan goals vsem generate marp plan.vsem.json --theme corporate -o slides.md
```

`splan goals validate` also accepts a VSEM, detected from its `strategies`.

## Examples

- [`examples/minimal.vsem.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/minimal.vsem.json)
- [`examples/product.vsem.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/product.vsem.json), a fictional bookkeeping app's plan that meets the strict options

The schema is `schema/vsem.schema.json`, embedded as `schema.VSEMSchema()`.
