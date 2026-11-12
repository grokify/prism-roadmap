# OGSM Integration

OGSM (Objective, Goals, Strategies, Measures) puts a plan on one page in four parts, two in words and two in numbers. The `goals/ogsm` package models it as a document: one objective, the goals that say it has been reached, the strategies that will reach them, and the measures that show each strategy is working, with the actions that carry the strategies out.

## The Four Parts

| Part | Words or numbers | Question | In the document |
|------|------------------|----------|-----------------|
| **Objective** | Words | What one thing do we want to achieve? | `objective`, with a `dueDate` at the end of the horizon |
| **Goals** | Numbers | How will we know the objective has been reached? | `goals`, each with a baseline, target, and due date |
| **Strategies** | Words | What choices will get us there? | `strategies`, each naming the goals it serves |
| **Measures** | Numbers | How will we know each strategy is working, before the goals move? | `measures`, each naming its strategy |

Most OGSMs also list, for each strategy, the **actions** that carry it out, with owners and dates. The document holds them as `actions`, each naming its strategy.

Goals measure the objective; measures track the strategies. A goal such as 90-day retention moves slowly and late; a measure such as the share of new accounts that connect a bank in their first session moves within weeks, and says early whether the strategy is working.

### Cascading

OGSMs cascade: a strategy at one level becomes the objective of the level below. A team's OGSM names the parent document in `metadata.parentId` and the parent strategy its objective carries out in `objective.parentStrategyId`.

### How it differs from VSEM and V2MOM

| | OGSM | VSEM | V2MOM |
|---|------|------|-------|
| Top of the plan | One objective, in words | A vision or mission | A vision |
| Numbers for the top | Goals | Metrics with no strategy | Measures |
| Numbers per strategy | Measures | Metrics | Measures under methods |
| Work | Actions | Execution items | Projects |
| Also holds | Which goals each strategy serves | Resources each strategy commits | Values and obstacles |

## Document Structure

```json
{
  "$schema": "https://github.com/grokify/prism-roadmap/schema/ogsm.schema.json",
  "metadata": { "id": "bookkeeping-2027-2029", "name": "Bookkeeping App OGSM 2027-2029", "period": "2027-2029" },
  "objective": {
    "statement": "Become the bookkeeping app small businesses keep using",
    "dueDate": "2029-12-31"
  },
  "goals": [
    { "id": "goal-retention", "name": "Accounts still active at 90 days", "owner": "Head of Product",
      "baseline": "52", "target": "65", "unit": "%", "dueDate": "2029-12-31" }
  ],
  "strategies": [
    { "id": "onboarding", "name": "Get every new account to a reconciled month in its first week",
      "owner": "Head of Onboarding", "goalIds": ["goal-retention"] }
  ],
  "measures": [
    { "id": "first-session-connect", "name": "New accounts that connect a bank in their first session",
      "strategyId": "onboarding", "owner": "Head of Onboarding", "baseline": "45", "target": "75", "unit": "%" }
  ],
  "actions": [
    { "id": "guided-import", "name": "Ship guided bank import", "strategyId": "onboarding",
      "owner": "Onboarding squad", "startDate": "2027-01-04", "endDate": "2027-03-26" }
  ]
}
```

## Go Usage

```go
import "github.com/grokify/prism-roadmap/goals/ogsm"

o, err := ogsm.ReadFile("plan.ogsm.json")
if err != nil {
    return err
}
for _, g := range o.Goals {
    fmt.Printf("%s: %s -> %s\n", g.Name, g.Baseline, g.Target)
    for _, s := range o.StrategiesFor(g.ID) {
        fmt.Printf("  %s: %d measures, %d actions\n", s.Name, len(o.MeasuresFor(s.ID)), len(o.ActionsFor(s.ID)))
    }
}
```

## Validation

`Validate` returns errors and warnings in document order. Structural problems are always errors:

- no objective, no goals, or no strategies
- a missing ID or name, or an ID used twice anywhere in the document
- a strategy that serves an unknown goal, or a measure or action without a strategy or with an unknown one
- an action that depends on itself or on an unknown action, or that ends before it starts
- an unknown status, priority, or direction, a progress outside 0 to 1, or a date that is not `YYYY-MM-DD`

Gaps are warnings under `DefaultValidationOptions()` and errors under `StrictValidationOptions()`:

| Gap | Strict option |
|-----|---------------|
| A goal or measure without a target | `RequireTargets` |
| A goal, measure, or action without an owner | `RequireOwners` |
| An objective without a due date | `RequireObjectiveDate` |
| A strategy that no measure tracks | `RequireStrategyMeasure` |
| A strategy that no action carries out | `RequireStrategyAction` |
| A strategy that serves no goal, or a goal no strategy serves | `RequireGoalLinks` |

The strict options also warn at more than five goals or five strategies, more than fit on one page.

## Goals Alignment

A PRD or roadmap can carry an OGSM through the framework-agnostic `goals` package. Strategies become goal items; goals and then measures become result items, a goal with no `GoalID` because it measures the objective:

```go
g := goals.NewOGSM(o)
g.GoalLabel()   // "Strategies"
g.ResultLabel() // "Goals and Measures"
```

## Rendering

The `render/markdown` package writes a Markdown document that opens with the classic one-page OGSM, objective, goals, strategies, and measures side by side, followed by the goals and each strategy with its measures and actions. The `render/marp` package writes Marp slides. Both are deterministic.

```go
import (
    "github.com/grokify/prism-roadmap/goals/ogsm/render"
    "github.com/grokify/prism-roadmap/goals/ogsm/render/markdown"
    "github.com/grokify/prism-roadmap/goals/ogsm/render/marp"
)

doc, err := markdown.New().Render(o, render.DefaultOptions())
slides, err := marp.New().Render(o, render.OnePageOptions())
```

`DefaultOptions` includes descriptions and actions; `OnePageOptions` leaves them out.

To embed the Markdown document in a larger one, set `OmitHeader` to leave out the title and the metadata line under it, and `HeadingOffset` to push every heading down that many levels, capped at level 6; fenced code blocks are left alone. With `OmitHeader` and a `HeadingOffset` of 1, the sections that were level-two headings become level three. The Marp renderer ignores both options.

## Command Line

```bash
splan goals ogsm init --name "OGSM 2027-2029" -o plan.ogsm.json
splan goals ogsm validate plan.ogsm.json --strict
splan goals ogsm generate markdown plan.ogsm.json -o plan.md
splan goals ogsm generate marp plan.ogsm.json --theme corporate -o slides.md
```

The template meets the strict options, with its objective due three years from today. `splan goals validate` also accepts an OGSM, detected from its single `objective`.

## Examples

- [`examples/minimal.ogsm.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/minimal.ogsm.json)
- [`examples/product.ogsm.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/product.ogsm.json), a fictional bookkeeping app's three-year plan that meets the strict options

The schema is `schema/ogsm.schema.json`, embedded as `schema.OGSMSchema()`.
