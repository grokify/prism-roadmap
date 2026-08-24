# Report Contracts

`ReportDataset`, `OpportunityReport`, and `PortfolioReview` (all in `github.com/grokify/prism-roadmap/assessment`) are the renderer-ready contracts that turn a corpus of [`OpportunityAssessment`s](overview.md) into a report. The design principle carries over unchanged: **a report is a pure function of its dataset.** No fact is computed by a renderer, and no number is invented by an LLM narrative pass — every generated sentence is anchored to either a computed fact or a cited `Evidence` record.

These are the same contracts [omniroadmap](https://github.com/grokify/omniroadmap) compiles, reviews, and renders against — this package defines the shape; omniroadmap owns persistence, compilation, and the actual markdown/slide render.

## ReportDataset: the compile output

`ReportDataset` is what omniroadmap's compile step produces from a corpus of assessments — the single input every report render reads from:

```go
type ReportDataset struct {
    GeneratedAt          time.Time
    RankingPolicyID      string
    RankingPolicyVersion string

    Ranking             []OpportunityRank        // full calculated-and-governed ranking
    Distributions       []DimensionDistribution  // % Person-Days per portfolio dimension option
    CapabilityOverlay   []CapabilityInvestment    // % Person-Days per referenced capability
    ObjectiveInvestment []ObjectiveInvestment     // % Person-Days per OKR objective
    Deltas              *ReportDeltas             // vs. the previous review cycle, nil for a first review
    OverrideLog         []RankOverride
}
```

```go
dataset := assessment.NewReportDataset(time.Now(), assessment.DefaultRankingPolicy(), finalRanking)
dataset.Distributions = []assessment.DimensionDistribution{
    assessment.ComputeDimensionDistribution("kano", assessments),
    assessment.ComputeDimensionDistribution("market-investment-horizon", assessments),
}
dataset.CapabilityOverlay = assessment.ComputeCapabilityOverlay(assessments)
dataset.ObjectiveInvestment = assessment.ComputeObjectiveInvestment(assessments)
```

Every distribution is weighted by **Person-Days, not opportunity count** — ten 2-PD items shouldn't look more significant than one 100-PD item. `UnclassifiedPersonDays` on `DimensionDistribution` surfaces Person-Days from assessments with no assignment for that dimension explicitly, rather than silently excluding them from the percentages.

`ComputeDeltas(previous, current)` compares two datasets for the same portfolio and returns `ReportDeltas` — rank moves, additions, removals, and distribution shifts since the last review. This is what makes a recurring portfolio review start from "what changed and why" instead of from zero every time.

## OpportunityReport: the per-opportunity six-pager

`OpportunityReport` is the "audit this one ranking end-to-end" artifact — scoped to a single opportunity, assembled from one `OpportunityAssessment` and its `OpportunityRank`:

```go
report := assessment.NewOpportunityReport(time.Now(), *myAssessment, myRank)
```

The section list is fixed — every report has the same table of contents; only which sections actually render varies by data availability:

| Sections (fixed order) | Present when |
|---|---|
| Recommendation Summary | always |
| Opportunity Definition | always |
| Prioritization Rationale | RICE or MoSCoW answers recorded |
| Portfolio Context | any `Dimensions` recorded |
| Strategic & Capability Alignment | any `Contributions` or `Capabilities` recorded |
| Risks, Assumptions & Decision | always |

| Appendices (fixed order) | Present when |
|---|---|
| Full Evidence Log | `assessment.HasEvidence()` |
| Rubric Answer Trace | `assessment.HasRubricAnswers()` |
| Assessment Provenance | `assessment.Judge != nil` |
| Related Opportunities | set by a renderer with corpus access (needs cross-referencing beyond one assessment) |

`report.PresentSections()` / `report.PresentAppendices()` return only the sections a renderer should actually emit, in order.

Each section carries an optional `NarrativeSlot`:

```go
type NarrativeSlot struct {
    ID          string
    Text        string   // empty until an LLM narrative pass fills it in
    DerivedFrom []string // computed facts this narrative interprets, e.g. "rank.finalRank"
    EvidenceIDs []string // evidence cited beyond what DerivedFrom implies
}
```

The slot and its anchors (`DerivedFrom`/`EvidenceIDs`) are produced deterministically regardless of whether a narrative pass has run — a reviewer or a later regeneration can always verify a filled-in `Text` wasn't invented, by checking it against the anchors that traveled with it.

## PortfolioReview: the whole-roadmap counterpart

`PortfolioReview` is `OpportunityReport`'s portfolio-wide sibling — one `ReportDataset` plus a deterministic review agenda, reusing the same `ReportSection`/`NarrativeSlot` contract so document and slide renderers share one model:

```go
review := assessment.NewPortfolioReview(dataset)
```

The agenda is generic by design — it never names "Kano" or "Market Investment Horizon" specifically, since `ReportDataset.Distributions` already holds every portfolio dimension (built-in or custom) uniformly:

| Agenda (fixed order) | Present when |
|---|---|
| Executive Summary | always |
| Decision Requested | always |
| Prioritized Roadmap | ranking present |
| Changes Since Previous Review | `Deltas` set |
| Prioritization Methodology | always |
| Portfolio Composition | distributions present |
| Capability Stack | capability overlay present |
| Strategic Alignment | objective investment present |
| Key Opportunities | ranking present |
| Governance & Overrides | override log present |
| Recommendations | always |
| Decisions / Open Questions | always |

`review.PresentAgenda()` / `review.PresentAppendices()` mirror `OpportunityReport`'s present-only accessors.

### Presentation projection

The same `PortfolioReview` renders into slide form via `PresentationProjection()` — one slide per present agenda section by default, each traceable back to its source section:

```go
slides := review.PresentationProjection()
// []PresentationSlide{ID, SourceSectionID, Headline, Narrative}
```

This is "the presentation is a projection of the same document model," not a second authoring surface — a renderer may split a data-heavy section (e.g. Portfolio Composition, if it has many dimensions) across multiple slides, but the default one-slide-per-section layout is what the contract guarantees.

## Next Steps

- [Opportunity Prioritization Overview](overview.md) — Evidence, Ladder, MoSCoW/RICE, dimensions, ranking
- [Feature Prioritization](../canvas/prioritization.md) — the underlying RICE/Kano/MoSCoW types
