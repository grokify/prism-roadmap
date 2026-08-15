# Feature Prioritization

prism-roadmap includes multiple industry-standard prioritization frameworks that integrate with OpportunitySpec: **RICE Scoring**, the **Kano Model**, and **MoSCoW Prioritization**.

## Overview

| Framework | Purpose | Output | Best For |
|-----------|---------|--------|----------|
| **RICE** | Quantitative prioritization | Numeric score | Ranking features by ROI |
| **Kano** | Qualitative classification | Category | Understanding customer impact |
| **MoSCoW** | Release planning | Priority category | Sprint/release scoping |

All frameworks are embedded in `OpportunitySpec` to enable prioritization at the feature discovery stage.

## MoSCoW Prioritization

MoSCoW is a prioritization technique used in release planning to categorize requirements by necessity.

**Reference:** [Wikipedia - MoSCoW Method](https://en.wikipedia.org/wiki/MoSCoW_method)

### Categories

| Category | Weight | Description |
|----------|--------|-------------|
| **Must Have** | 4 | Critical requirements. Release fails without them. |
| **Should Have** | 3 | Important but not critical. Can be worked around. |
| **Could Have** | 2 | Desirable. Include if time permits. |
| **Won't Have** | 0 | Explicitly excluded from current release scope. |

### MoSCoWPriority Type

```go
type MoSCoWPriority string

const (
    MoSCoWMustHave   MoSCoWPriority = "must_have"
    MoSCoWShouldHave MoSCoWPriority = "should_have"
    MoSCoWCouldHave  MoSCoWPriority = "could_have"
    MoSCoWWontHave   MoSCoWPriority = "wont_have"
)
```

### Usage Example

```go
import "github.com/grokify/prism-roadmap/prioritization"

// Get priority weight for sorting
priority := prioritization.MoSCoWMustHave
weight := priority.Weight() // Returns 4

// Check if valid
if priority.IsValid() {
    fmt.Printf("Priority: %s (weight: %d)\n", priority, weight)
}

// Get all valid priorities
allPriorities := prioritization.ValidMoSCoWPriorities()
```

### Weight Method

The `Weight()` method returns a numeric value for sorting:

| Priority | Weight |
|----------|--------|
| `must_have` | 4 |
| `should_have` | 3 |
| `could_have` | 2 |
| `wont_have` | 0 |

### When to Use MoSCoW

- Sprint or release planning sessions
- Negotiating scope with stakeholders
- Making explicit trade-offs under time pressure
- Clarifying what "out of scope" means

## RICE Scoring

RICE is a quantitative framework for prioritizing features based on four factors.

**Reference:** [Intercom - RICE Scoring](https://www.intercom.com/blog/rice-simple-prioritization-for-product-managers/)

### Formula

```
Score = (Reach × Impact × Confidence) / Effort
```

### Components

| Component | Description | Range |
|-----------|-------------|-------|
| **Reach** | Number of users affected per time period | Integer (users/quarter) |
| **Impact** | Effect on each user | massive (3x), high (2x), medium (1x), low (0.5x), minimal (0.25x) |
| **Confidence** | Certainty in estimates | high (100%), medium (80%), low (50%) |
| **Effort** | Resources required | Float (person-months) |

### Impact Levels

| Level | Multiplier | Description |
|-------|------------|-------------|
| `massive` | 3.0 | Game-changing impact for every user |
| `high` | 2.0 | Significant positive impact |
| `medium` | 1.0 | Standard improvement |
| `low` | 0.5 | Minor benefit |
| `minimal` | 0.25 | Marginal improvement |

### Confidence Levels

| Level | Multiplier | When to Use |
|-------|------------|-------------|
| `high` | 1.0 (100%) | Strong data, validated assumptions |
| `medium` | 0.8 (80%) | Some data, reasonable assumptions |
| `low` | 0.5 (50%) | Gut feel, limited data |

### Parsing and Validation

`ImpactLevel` and `ConfidenceLevel` each have `IsValid()` and a
`ParseXxxLevel()` constructor that rejects unrecognized input instead of
silently accepting it:

```go
impact, err := prioritization.ParseImpactLevel("High") // case-insensitive
if err != nil {
    // unrecognized level
}

prioritization.ImpactHigh.IsValid()      // true
prioritization.ImpactLevel("urgent").IsValid() // false
```

`RICEScore.Validate()` uses the same check: a non-empty but unrecognized
`Impact` or `Confidence` value is rejected rather than silently falling
back to a default multiplier in `Calculate()`. An empty value is still
allowed (not yet estimated).

### Example Calculation

```go
// Feature: Bulk export
score := prioritization.NewRICEScore(
    "bulk-export",
    1000,                        // Reach: 1000 users/quarter
    prioritization.ImpactHigh,   // Impact: 2x
    prioritization.ConfidenceHigh, // Confidence: 100%
    2.0,                         // Effort: 2 person-months
)
score.Calculate()
// Score = (1000 × 2.0 × 1.0) / 2.0 = 1000
```

### RICEScore Fields

```go
type RICEScore struct {
    FeatureID   string          `json:"featureId"`
    FeatureName string          `json:"featureName,omitempty"`

    // RICE components
    Reach      int             `json:"reach"`
    ReachUnit  string          `json:"reachUnit"`   // "quarter", "month"
    Impact     ImpactLevel     `json:"impact"`
    Confidence ConfidenceLevel `json:"confidence"`
    Effort     float64         `json:"effort"`
    EffortUnit string          `json:"effortUnit"`  // "person-months"

    // Calculated
    Score float64 `json:"score"`

    // Justifications
    ReachJustification      string `json:"reachJustification,omitempty"`
    ImpactJustification     string `json:"impactJustification,omitempty"`
    ConfidenceJustification string `json:"confidenceJustification,omitempty"`
    EffortJustification     string `json:"effortJustification,omitempty"`
}
```

### Ranking Features

```go
set := prioritization.NewRICEScoreSet()
set.Add(*prioritization.NewRICEScore("feat-1", 1000, ImpactHigh, ConfidenceHigh, 1))
set.Add(*prioritization.NewRICEScore("feat-2", 500, ImpactMassive, ConfidenceMedium, 3))
set.Add(*prioritization.NewRICEScore("feat-3", 2000, ImpactLow, ConfidenceLow, 0.5))

// Get top 2 features
top2 := set.TopN(2)

// Get rank of specific feature
rank := set.Rank("feat-2") // Returns 1-based rank
```

## Kano Model

The Kano Model classifies features by their impact on customer satisfaction.

**Reference:** [ProductPlan - Kano Model](https://www.productplan.com/glossary/kano-model/)

### Categories

| Category | Description | Priority |
|----------|-------------|----------|
| **Must-Be** | Basic expectation. Absence causes dissatisfaction, presence expected. | 5 (highest) |
| **Performance** | Linear satisfaction. More is better. | 4 |
| **Attractive** | Delighter. Unexpected positive surprise. | 3 |
| **Indifferent** | No significant impact on satisfaction. | 1 |
| **Reverse** | Unwanted. Presence causes dissatisfaction. | 0 |
| **Questionable** | Contradictory response, needs clarification. | 0 |

### Kano Questionnaire

The Kano Model uses paired questions:

1. **Functional question:** "If the product HAS this feature, how do you feel?"
2. **Dysfunctional question:** "If the product DOES NOT HAVE this feature, how do you feel?"

**Response options:**

- `like` - I like it
- `expect` - I expect it
- `neutral` - I am neutral
- `tolerate` - I can tolerate it
- `dislike` - I dislike it

### Evaluation Table

The response pair maps to a category:

```
                    | Dysfunctional Response
Functional Response | Like    | Expect  | Neutral | Tolerate | Dislike
--------------------|---------|---------|---------|----------|--------
Like                | Q       | A       | A       | A        | O
Expect              | R       | I       | I       | I        | M
Neutral             | R       | I       | I       | I        | M
Tolerate            | R       | I       | I       | I        | M
Dislike             | R       | R       | R       | R        | Q

Legend: M=Must-be, O=Performance, A=Attractive, I=Indifferent, R=Reverse, Q=Questionable
```

### Classification Example

```go
// If user LIKES having dark mode (functional)
// and is NEUTRAL about not having it (dysfunctional)
// → Category is ATTRACTIVE (delighter)

category := prioritization.ClassifyKano(
    prioritization.KanoLike,
    prioritization.KanoNeutral,
)
// category == KanoAttractive
```

### KanoFeature Fields

```go
type KanoFeature struct {
    FeatureID   string `json:"featureId"`
    FeatureName string `json:"featureName"`
    Description string `json:"description,omitempty"`

    // Questionnaire responses
    FunctionalResponse    KanoResponse `json:"functionalResponse"`
    DysfunctionalResponse KanoResponse `json:"dysfunctionalResponse"`

    // Classified category
    Category KanoCategory `json:"category"`

    // Coefficients (from aggregate analysis)
    SatisfactionCoefficient    float64 `json:"satisfactionCoefficient,omitempty"`
    DissatisfactionCoefficient float64 `json:"dissatisfactionCoefficient,omitempty"`

    // Response counts (for aggregate analysis)
    MustBeCount      int `json:"mustBeCount,omitempty"`
    PerformanceCount int `json:"performanceCount,omitempty"`
    AttractiveCount  int `json:"attractiveCount,omitempty"`
    IndifferentCount int `json:"indifferentCount,omitempty"`
}
```

### Satisfaction Coefficients

When analyzing multiple respondents, calculate coefficients:

```go
// Satisfaction: How much satisfaction increases WITH the feature
// Formula: (A + O) / (A + O + M + I)

// Dissatisfaction: How much satisfaction decreases WITHOUT the feature
// Formula: -1 * (O + M) / (A + O + M + I)
```

**Example:**

```go
feature := prioritization.KanoFeature{
    FeatureID:        "dark-mode",
    AttractiveCount:  20,
    PerformanceCount: 30,
    MustBeCount:      40,
    IndifferentCount: 10,
}
feature.CalculateCoefficients()

// SatisfactionCoefficient = (20+30)/(20+30+40+10) = 0.5
// DissatisfactionCoefficient = -1*(30+40)/(100) = -0.7
```

### Aggregate Analysis

```go
analysis := prioritization.NewKanoAnalysis()

analysis.Add(KanoFeature{
    FeatureID:             "feature-1",
    FunctionalResponse:    KanoExpect,
    DysfunctionalResponse: KanoDislike,
}) // → Must-Be

analysis.Add(KanoFeature{
    FeatureID:             "feature-2",
    FunctionalResponse:    KanoLike,
    DysfunctionalResponse: KanoNeutral,
}) // → Attractive

// Filter by category
mustHaves := analysis.MustHaves()
delighters := analysis.Delighters()

// Sort by priority
analysis.SortByPriority()

// Get summary
summary := analysis.Summary()
// map[KanoCategory]int{"must-be": 1, "attractive": 1}
```

## Integration with OpportunitySpec

OpportunitySpec includes both frameworks as optional fields:

```go
type OpportunitySpec struct {
    // ... 12 boxes ...

    // Prioritization Frameworks
    RICE *prioritization.RICEScore   `json:"rice,omitempty"`
    Kano *prioritization.KanoFeature `json:"kano,omitempty"`
}
```

### Setting RICE Score

```go
spec := canvas.NewOpportunitySpec("mobile-app", "Mobile App Feature")

// Set RICE using convenience method
spec.SetRICE(1000, ImpactHigh, ConfidenceHigh, 2.0)

// Check RICE
if spec.HasRICE() {
    score := spec.GetRICEScore()
    fmt.Printf("RICE Score: %.0f\n", score)
}
```

### Setting Kano Category

```go
// Set Kano from questionnaire responses
spec.SetKano(KanoLike, KanoDislike)

// Check Kano
if spec.HasKano() {
    category := spec.GetKanoCategory()
    fmt.Printf("Kano: %s\n", category) // "performance"
}

// Convenience checks
if spec.IsMustHave() {
    // Prioritize this feature
}
if spec.IsDelighter() {
    // Consider for differentiation
}
```

### Prioritization Summary

```go
summary := spec.GetPrioritizationSummary()
fmt.Printf("RICE: %.0f, Kano: %s, Priority: %d\n",
    summary.RICEScore,
    summary.KanoCategory,
    summary.KanoPriority,
)
```

## Decision Framework

### When to Use RICE

- Comparing multiple features for roadmap prioritization
- Need quantitative justification for investment
- Effort and reach estimates are reasonably reliable

### When to Use Kano

- Understanding customer expectations vs. delighters
- Deciding between "must-have" and "nice-to-have"
- Product-market fit validation

### Combined Approach

Use both for comprehensive prioritization:

1. **Kano first:** Classify features by customer impact
2. **RICE for tiebreakers:** Within each Kano category, rank by RICE score

```go
// Filter must-haves, then sort by RICE
mustHaves := analysis.MustHaves()
riceSet := prioritization.NewRICEScoreSet()
for _, f := range mustHaves {
    // Assume RICE scores are attached
    riceSet.Add(f.RICE)
}
topMustHave := riceSet.TopN(1)[0]
```

### Priority Matrix

| Kano Category | RICE Score | Action |
|---------------|------------|--------|
| Must-Be | Any | Implement first |
| Performance | High (>1000) | High priority |
| Performance | Low (<500) | Medium priority |
| Attractive | High | Differentiation opportunity |
| Attractive | Low | Nice to have |
| Indifferent | Any | Deprioritize |
| Reverse | Any | Do not implement |

## See Also

- [OpportunitySpec](opportunity-spec.md) - 12-box discovery framework
- [RICEScore API](../api/reference.md#ricescore) - API reference
- [KanoFeature API](../api/reference.md#kanofeature) - API reference
