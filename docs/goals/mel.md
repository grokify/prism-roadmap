# MEL Plans

A monitoring, evaluation, and learning (MEL) plan says how a team will know whether its work is producing the results it intends, and what it will do with what it learns. The `goals/mel` package models it as one document of linked parts.

- **Monitoring** tracks indicators continuously against their baselines and targets.
- **Evaluation** asks, at set times, whether results happened and why.
- **Learning** turns both into decisions, through questions the team commits to answering and reviews where it does.

MEL plans are common in program work and are just as useful for a product or business team: any team that sets goals and has to show results needs to say how it will measure them and when it will stop to learn.

## The Seven Parts

| Part | Question | In the document |
|------|----------|-----------------|
| **Results** | What should change, from outputs up to impact? | `results`, each with a `level` and the higher results it `contributesTo` |
| **Activities** | What work produces the results, who does it, and when? | `activities`, each with the results it produces, an owner, dates, and a status |
| **Indicators** | How exactly is each result, activity, or assumption measured? | `indicators`: one reference sheet each, naming what it `measures`, with definition, baseline, targets, data source, collection method, frequency, owner, and data quality |
| **Assumptions** | What must hold for the results to follow, and how is it watched? | `assumptions`, each with how it is monitored |
| **Evaluations** | When will we ask whether and why results happened? | `evaluations`, each with questions, type, method, evaluator, and dates |
| **Learning agenda** | What do we most need to find out, by when, and how? | `learningAgenda`: questions with the decision each informs, methods, owner, and due date |
| **Reviews and reports** | When do we stop to learn, and who hears what? | `reviews`, and `reports` with what each audience needs to know, what it receives, and how often |

Results use five levels, from the most immediate to the most distant: `output`, `short_term`, `medium_term`, `long_term`, and `impact`. Every element has an ID, IDs are unique across the plan, and elements refer to each other by ID.

## Document Structure

```json
{
  "$schema": "https://github.com/grokify/prism-roadmap/schema/mel.schema.json",
  "metadata": { "id": "onboarding-mel", "name": "Onboarding MEL Plan", "period": "FY2027" },
  "results": [
    { "id": "active-90", "statement": "New accounts are still active at 90 days", "level": "long_term" },
    { "id": "first-session-connect", "statement": "New accounts connect a bank in their first session",
      "level": "short_term", "contributesTo": ["active-90"] }
  ],
  "indicators": [
    {
      "id": "ind-active-90",
      "name": "Accounts still active at 90 days",
      "measures": "active-90",
      "unit": "%",
      "direction": "increase",
      "baseline": { "value": "52", "date": "2026-12-31" },
      "targets": [{ "value": "60", "date": "2027-12-31" }],
      "dataSource": "Product analytics",
      "collectionMethod": "Monthly cohort query",
      "frequency": "Monthly",
      "owner": "Head of Product",
      "dataQuality": { "checks": ["Test accounts excluded"] }
    }
  ],
  "evaluations": [
    {
      "id": "eval-first-year",
      "title": "First-year results evaluation",
      "type": "impact",
      "questions": ["Did guided import raise 90-day activity, compared with the holdout cohort?"],
      "evaluator": "external",
      "startDate": "2027-10-01",
      "endDate": "2028-01-31"
    }
  ],
  "learningAgenda": [
    {
      "id": "lq-connect-cause",
      "question": "Does connecting a bank early lead accounts to stay, or do committed owners simply do both?",
      "rationale": "Decides whether to keep investing in onboarding.",
      "methods": ["Holdout cohort comparison"],
      "dueDate": "2027-06-30",
      "relatedIds": ["active-90"]
    }
  ],
  "reviews": [
    { "id": "rev-quarterly", "name": "Quarterly learning review", "frequency": "Quarterly",
      "questionIds": ["lq-connect-cause"] }
  ],
  "reports": [
    { "id": "rep-board", "audience": "Board", "frequency": "Quarterly", "indicatorIds": ["ind-active-90"] }
  ]
}
```

Values are strings, so a qualitative indicator can have a baseline such as "Emerging" and a target such as "Established".

An indicator usually measures a result. One that measures an activity tracks delivery, such as the share of new accounts offered a feature; one that measures an assumption watches whether it still holds, such as a partner's connection success rate.

## Evaluation Types

| Type | What it is for |
|------|----------------|
| `baseline` | Establishes the starting point |
| `process` | Checks whether the work is delivered as planned |
| `formative` | Improves the work while it runs |
| `midterm` | Takes stock partway through |
| `final` | Judges results at the end (a summative evaluation) |
| `impact` | Attributes change to the work, against a comparison or counterfactual |
| `developmental` | Supports work that is still being designed |

An evaluation's `evaluator` is `internal`, `external`, or `mixed`.

## Go Usage

```go
import "github.com/grokify/prism-roadmap/goals/mel"

p, err := mel.ReadFile("onboarding.mel.json")
if err != nil {
    return err
}

for _, top := range p.TopResults() {
    fmt.Println(mel.LevelLabel(top.Level), top.Statement)
    for _, r := range p.Contributors(top.ID) {
        fmt.Println("  ", mel.LevelLabel(r.Level), r.Statement, len(p.IndicatorsFor(r.ID)), "indicators")
    }
}
```

## Validation

`Validate` returns errors and warnings in document order. Structural problems are always errors:

- no results, or a result without a statement or a known level
- a missing or duplicate ID
- a reference to an unknown element, or to the wrong kind of element, such as an indicator that measures a review
- a result that contributes to itself, or a cycle in the results framework
- an unknown type, direction, evaluation type, evaluator, or status
- a baseline or target without a value, a date that is not `YYYY-MM-DD`, or an activity or evaluation that ends before it starts

A result that contributes to a lower-level result is a warning. Gaps in the plan are warnings under `DefaultValidationOptions()` and errors under `StrictValidationOptions()`:

| Gap | Strict option |
|-----|---------------|
| An indicator without a data source, collection method, frequency, or owner | `RequireCollectionPlan` |
| An indicator without a baseline | `RequireBaseline` |
| An indicator without a target | `RequireTarget` |
| A result that no indicator measures | `RequireResultIndicator` |
| An evaluation without dates, or a learning question without a due date | `RequireTiming` |
| An evaluation without questions | `RequireEvaluationQuestions` |
| An assumption with no way to monitor it | `RequireAssumptionMonitoring` |
| No learning questions | `RequireLearningAgenda` |
| No evaluations | `RequireEvaluationPlan` |
| No reviews, or a review or report without a frequency | `RequireSchedule` |

A learning question without methods, or one marked answered with no answer recorded, is always a warning.

## Rendering

The `render/markdown` package writes the plan as a Markdown document, and `render/pandoc` writes the same body with YAML front matter for a PDF:

```go
import (
    "github.com/grokify/prism-roadmap/goals/mel/render"
    "github.com/grokify/prism-roadmap/goals/mel/render/markdown"
    "github.com/grokify/prism-roadmap/goals/mel/render/pandoc"
)

doc, err := markdown.New().Render(p, render.DefaultOptions())
pdfSource, err := pandoc.New().Render(p, render.SummaryOptions())
```

```bash
pandoc plan.md -o plan.pdf --pdf-engine=lualatex
```

The document draws the results framework as a nested list from impact down, then the activities, an indicator summary table, a reference sheet per indicator, the assumptions to monitor, the evaluation plan with each evaluation's questions, the learning agenda, and the review and reporting schedule. `SummaryOptions` leaves out the reference sheets and descriptions. Both renderers are deterministic; the Pandoc date is the plan's `metadata.updatedAt`, when set.

To embed the Markdown document in a larger one, set `OmitHeader` to leave out the title and the metadata line under it, and `HeadingOffset` to push every heading down that many levels, capped at level 6; fenced code blocks are left alone. With `OmitHeader` and a `HeadingOffset` of 1, the sections that were level-two headings become level three. The Pandoc renderer ignores both options.

## Command Line

```bash
splan goals mel init --name "Onboarding MEL Plan" -o plan.mel.json
splan goals mel validate plan.mel.json --strict
splan goals mel generate markdown plan.mel.json -o plan.md
splan goals mel generate pandoc plan.mel.json --summary -o plan-summary.md
```

The template meets the strict options, with its dates set from today. `splan goals validate` also accepts a MEL plan, detected from its `results`.

## Examples

- [`examples/minimal.mel.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/minimal.mel.json)
- [`examples/product.mel.json`](https://github.com/grokify/prism-roadmap/blob/main/examples/product.mel.json), a fictional bookkeeping app's onboarding plan that meets the strict options

The schema is `schema/mel.schema.json`, embedded as `schema.MELSchema()`.
