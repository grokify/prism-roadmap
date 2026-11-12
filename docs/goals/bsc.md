# Balanced Scorecard Integration

The Balanced Scorecard describes a strategy as objectives in a few perspectives, balanced so that the results an organization reports are not pursued at the expense of what produces them. Robert Kaplan and David Norton introduced it in the *Harvard Business Review* in 1992, and later added the strategy map, which links the objectives from cause to effect. The `goals/bsc` package models a scorecard as a document: its perspectives, strategic themes, objectives and their strategy map, measures, and initiatives.

## The Parts

| Part | Question | In the document |
|------|----------|-----------------|
| **Mission, vision, strategy** | Why do we exist, what do we want to become, and how will we get there? | `mission`, `vision`, `strategy` |
| **Perspectives** | From which points of view must the strategy succeed? | `perspectives`, top first |
| **Strategic objectives** | What must we achieve in each perspective? | `objectives`, each naming its perspective |
| **Strategy map** | How does achieving one objective help achieve another? | `contributesTo` on each objective |
| **Strategic themes** | Which objectives together carry out one part of the strategy? | `themes`, and `themeIds` on objectives |
| **Measures** | How will we know each objective is being achieved? | `measures`, each naming its objective, lead or lag, with a target and thresholds |
| **Initiatives** | What work closes the gap between where a measure is and its target? | `initiatives`, each naming the objectives it serves |

### Perspectives

`bsc.New` and `bsc.DefaultPerspectives` give the four classic perspectives, top first, each with the question it answers:

| Perspective | Kind | Question |
|-------------|------|----------|
| Financial | `financial` | To succeed financially, how should we appear to those who fund us? |
| Customer | `customer` | To achieve our vision, how should we appear to our customers? |
| Internal Process | `internal_process` | To satisfy our customers, which processes must we excel at? |
| Learning and Growth | `learning_growth` | To excel at those processes, how will we sustain our ability to change and improve? |

Perspectives may be renamed, reordered, added, or replaced. A public agency or nonprofit often puts a `mission` perspective on top, with the financial perspective below it as a constraint rather than the goal. The `kind` says which classic role a perspective plays, whatever it is called; `custom` is for anything else.

### Strategy map

A strategy map runs from cause at the bottom to effect at the top: building skills (learning and growth) improves a process (internal process), which customers notice (customer), which shows in the results (financial). Each objective lists the objectives it helps achieve in `contributesTo`. A link may stay within a perspective, such as saving owners time keeping them as customers; a link that runs down the perspectives is allowed but warned about.

### Lead and lag measures

A **lag** measure reports a result after the fact, such as retention or revenue. A **lead** measure moves first and drives the result, such as the share of new accounts reconciled in their first week. A balanced scorecard has both: lag measures alone say whether the strategy worked, but not in time to change course.

Each measure may carry `thresholds`: the value at which it turns `green` and the value at which it turns `red`, in its unit. Between them it is yellow.

### Cascading

A team's scorecard names the scorecard it cascades from in `metadata.parentId`.

### How it differs from OGSM and VSEM

| | Balanced Scorecard | OGSM | VSEM |
|---|--------------------|------|------|
| Top of the plan | Mission, vision, and strategy | One objective | A vision or mission |
| Organized by | Perspectives | Goals and strategies | Strategies |
| Cause and effect | A strategy map between objectives | Which goals each strategy serves | Not recorded |
| Numbers | Lead and lag measures per objective | Goals for the objective, measures per strategy | Metrics per strategy and for the vision |
| Work | Initiatives, each serving several objectives | Actions per strategy | Execution items per strategy |

## Document Structure

```json
{
  "$schema": "https://github.com/grokify/prism-roadmap/schema/bsc.schema.json",
  "metadata": { "id": "bookkeeping-fy2027", "name": "Bookkeeping App Balanced Scorecard FY2027", "period": "FY2027" },
  "vision": "The bookkeeping app small businesses keep using",
  "perspectives": [
    { "id": "customer", "name": "Customer", "kind": "customer" },
    { "id": "internal-process", "name": "Internal Process", "kind": "internal_process" }
  ],
  "objectives": [
    { "id": "cus-retention", "name": "Keep small businesses using the app",
      "perspectiveId": "customer", "owner": "Head of Product" },
    { "id": "ip-onboarding", "name": "Get every new account to a reconciled month in its first week",
      "perspectiveId": "internal-process", "owner": "Head of Onboarding",
      "contributesTo": ["cus-retention"] }
  ],
  "measures": [
    { "id": "m-retention-90", "name": "Accounts still active at 90 days", "objectiveId": "cus-retention",
      "type": "lag", "baseline": "52", "target": "65", "unit": "%",
      "thresholds": { "green": "63", "red": "55" } },
    { "id": "m-first-week-reconciled", "name": "New accounts with a reconciled month in their first week",
      "objectiveId": "ip-onboarding", "type": "lead", "baseline": "18", "target": "40", "unit": "%" }
  ],
  "initiatives": [
    { "id": "init-guided-import", "name": "Ship guided bank import",
      "objectiveIds": ["ip-onboarding", "cus-retention"], "owner": "Onboarding squad",
      "startDate": "2027-01-04", "endDate": "2027-03-26" }
  ]
}
```

## Go Usage

```go
import "github.com/grokify/prism-roadmap/goals/bsc"

b, err := bsc.ReadFile("scorecard.bsc.json")
if err != nil {
    return err
}
for _, p := range b.Perspectives {
    fmt.Println(p.Name)
    for _, o := range b.ObjectivesIn(p.ID) {
        fmt.Printf("  %s: %d measures, %d initiatives, %d contributors\n", o.Name,
            len(b.MeasuresFor(o.ID)), len(b.InitiativesFor(o.ID)), len(b.ContributorsTo(o.ID)))
    }
}
```

## Validation

`Validate` returns errors and warnings in document order. Structural problems are always errors:

- no perspectives or no objectives
- a missing ID or name, or an ID used twice anywhere in the document
- an objective without a perspective or with an unknown one, or in an unknown theme
- a strategy-map link to an unknown objective or to the objective itself, or a cycle in the strategy map
- a measure without an objective or with an unknown one, or an initiative that serves no objective or an unknown one
- an initiative that depends on itself or on an unknown initiative, or that ends before it starts
- an unknown kind, type, status, priority, or direction, a progress outside 0 to 1, or a date that is not `YYYY-MM-DD`

Gaps are warnings under `DefaultValidationOptions()` and errors under `StrictValidationOptions()`:

| Gap | Strict option |
|-----|---------------|
| Neither a mission nor a vision | `RequireVision` |
| A measure without a target | `RequireTargets` |
| An objective, measure, or initiative without an owner | `RequireOwners` |
| A perspective with no objective | `RequirePerspectiveObjective` |
| An objective that no measure tracks | `RequireObjectiveMeasure` |
| An objective outside the top perspective that contributes to nothing, or outside the bottom one that nothing contributes to | `RequireStrategyMap` |
| A measure that is neither lead nor lag | `RequireMeasureType` |
| A scorecard without both lead and lag measures | `RequireLeadAndLag` |

The strict options also warn at more than 25 measures, more than a team can steer by. Some findings are always warnings: a strategy-map link that runs down the perspectives, a theme no objective belongs to, and an objective no initiative serves, since an objective already on track needs none.

## Goals Alignment

A PRD or roadmap can carry a Balanced Scorecard through the framework-agnostic `goals` package. Objectives become goal items, tagged with their perspective's name, and measures become result items under their objective:

```go
g := goals.NewBSC(b)
g.GoalLabel()   // "Objectives"
g.ResultLabel() // "Measures"
```

## Rendering

The `render/markdown` package writes a Markdown document: the mission, vision, and strategy, the strategy map as a Mermaid flowchart with a box per perspective, the strategic themes, the scorecard as a table per perspective, and the initiatives. The `render/marp` package writes Marp slides, with the strategy map as a table, since Marp does not draw Mermaid diagrams. Both are deterministic.

```go
import (
    "github.com/grokify/prism-roadmap/goals/bsc/render"
    "github.com/grokify/prism-roadmap/goals/bsc/render/markdown"
    "github.com/grokify/prism-roadmap/goals/bsc/render/marp"
)

doc, err := markdown.New().Render(b, render.DefaultOptions())
slides, err := marp.New().Render(b, render.CompactOptions())
```

`DefaultOptions` includes descriptions, the strategy map, and initiatives; `CompactOptions` keeps the strategy map and the scorecard only.

To embed the Markdown document in a larger one, set `OmitHeader` to leave out the title and the metadata line under it, and `HeadingOffset` to push every heading down that many levels, capped at level 6; fenced code blocks, such as the strategy map, are left alone. With `OmitHeader` and a `HeadingOffset` of 1, the sections that were level-two headings become level three. The Marp renderer ignores both options.

## Command Line

```bash
splan goals bsc init --name "Scorecard FY2027" -o scorecard.bsc.json
splan goals bsc validate scorecard.bsc.json --strict
splan goals bsc generate markdown scorecard.bsc.json -o scorecard.md
splan goals bsc generate marp scorecard.bsc.json --theme corporate -o slides.md
```

The template has the four classic perspectives with an objective in each, linked from learning and growth up to the financial result, and meets the strict options. `splan goals validate` also accepts a scorecard, detected from its `perspectives`.

## Examples

- [`examples/minimal.bsc.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/minimal.bsc.json)
- [`examples/product.bsc.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/product.bsc.json), a fictional bookkeeping app's scorecard with two strategic themes that meets the strict options

The schema is `schema/bsc.schema.json`, embedded as `schema.BSCSchema()`.
