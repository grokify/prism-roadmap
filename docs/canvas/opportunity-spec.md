# OpportunitySpec

The OpportunitySpec is a merged 12-box framework combining Jeff Patton's Opportunity Canvas (discovery-focused) with Marty Cagan's SVPG Opportunity Assessment (business case-focused).

## Overview

OpportunitySpec bridges the gap between user discovery and business case rigor, providing a comprehensive evaluation for feature-level opportunities.

**References:**

- [Jeff Patton - Opportunity Canvas](https://www.jpattonassociates.com/opportunity-canvas/)
- [Marty Cagan - Assessing Product Opportunities](https://www.svpg.com/assessing-product-opportunities/)

## Grid Layout (3×4)

```
┌─────────────────────┬─────────────────────┬─────────────────────┐
│ 1. Users & Problem  │ 2. Current Solutions│ 3. Solution Ideas   │
│                     │                     │                     │
├─────────────────────┼─────────────────────┼─────────────────────┤
│ 4. User Value       │ 5. Business Value   │ 6. Competitive Edge │
│                     │                     │                     │
├─────────────────────┼─────────────────────┼─────────────────────┤
│ 7. Market & Timing  │ 8. Go-to-Market     │ 9. Success Metrics  │
│                     │                     │                     │
├─────────────────────┼─────────────────────┼─────────────────────┤
│ 10. Critical Reqs   │ 11. Risks & Assump. │ 12. Recommendation  │
│                     │                     │                     │
└─────────────────────┴─────────────────────┴─────────────────────┘
```

### Example with Content

![OpportunitySpec Example](opportunity_spec_grid_example.svg)

## Row Organization

| Row | Focus | Boxes |
|-----|-------|-------|
| **Row 1: Discovery** | Understanding the problem space | Users & Problem, Current Solutions, Solution Ideas |
| **Row 2: Value** | Articulating value to users and business | User Value, Business Value, Competitive Edge |
| **Row 3: Market** | Market positioning and go-to-market | Market & Timing, Go-to-Market, Success Metrics |
| **Row 4: Validation** | Requirements, risks, and decision | Critical Reqs, Risks & Assumptions, Recommendation |

## Box Details

### Row 1: Discovery

#### Box 1: Users & Problem

*"Who has the problem and what problem are we solving?"*

Combines Patton's Users/Problems with Cagan's Value Proposition.

| Field | Description |
|-------|-------------|
| `primaryUsers` | Primary user types affected |
| `secondaryUsers` | Secondary user types |
| `problemStatement` | Clear statement of the problem |
| `customerPains` | Specific pain points |
| `currentState` | How things work today |
| `desiredState` | How things should work |
| `evidence` | Evidence the problem exists |

#### Box 2: Current Solutions

*"How do people solve this today?"*

Combines Patton's Solutions Today with Cagan's Competitive Landscape.

| Field | Description |
|-------|-------------|
| `directCompetitors` | Head-to-head competitors |
| `indirectCompetitors` | Alternative solutions |
| `workarounds` | Manual workarounds |
| `internalSolutions` | Internal tools being used |
| `doNothing` | Cost of status quo |
| `marketDynamics` | Competitive dynamics |

#### Box 3: Solution Ideas

*"What are our solution concepts?"*

From Patton's Solution Ideas with evaluation criteria.

| Field | Description |
|-------|-------------|
| `ideas` | Solution concepts with pros/cons |
| `recommendedIdea` | Which idea is recommended |
| `selectionReason` | Why this idea |
| `alternativesPros` | Why other ideas might work |
| `alternativesCons` | Why other ideas might not work |

### Row 2: Value

#### Box 4: User Value

*"What value does this provide to users?"*

Combines Patton's User Value with Cagan's Value Proposition.

| Field | Description |
|-------|-------------|
| `valueStatement` | Core value proposition for users |
| `keyBenefits` | Main benefits |
| `painsRelieved` | Which pains are addressed |
| `gainsEnabled` | What new gains are enabled |
| `jobsToBeDone` | Jobs this helps accomplish |
| `adoptionBarrier` | What might prevent adoption |

#### Box 5: Business Value

*"What value does this provide to the business?"*

Combines Patton's Business Problem/Metrics with Cagan's Market Size.

| Field | Description |
|-------|-------------|
| `businessProblem` | Why this matters to the business |
| `businessOutcomes` | Expected business outcomes |
| `revenueImpact` | Revenue potential |
| `costImpact` | Cost savings/implications |
| `strategicFit` | How this fits strategy |
| `marketSize` | TAM/SAM/SOM |

#### Box 6: Competitive Edge

*"Why are we best suited to pursue this?"*

From Cagan's Our Differentiator.

| Field | Description |
|-------|-------------|
| `differentiator` | What makes us different |
| `coreStrengths` | What we do well |
| `uniqueCapabilities` | What only we can do |
| `strategicAssets` | Assets we can leverage |
| `unfairAdvantage` | What's hard to replicate |
| `entryBarriers` | Barriers we create |

### Row 3: Market

#### Box 7: Market & Timing

*"Who's the market and why now?"*

Combines Cagan's Target Market with Market Window.

| Field | Description |
|-------|-------------|
| `primarySegment` | Main target segment |
| `secondarySegments` | Additional segments |
| `whyNow` | Primary reason for timing |
| `marketTriggers` | Events creating opportunity |
| `technologyShifts` | Tech changes enabling this |
| `windowDuration` | How long window is open |
| `urgencyLevel` | high, medium, low |

#### Box 8: Go-to-Market

*"How will we get this to users?"*

Combines Patton's Adoption Strategy with Cagan's GTM Strategy.

| Field | Description |
|-------|-------------|
| `strategy` | Overall GTM approach |
| `adoptionPath` | How users will find/adopt |
| `channels` | Distribution channels |
| `salesModel` | Self-serve, sales-led, PLG |
| `launchApproach` | Beta, GA, etc. |
| `estimatedCAC` | Customer acquisition cost |

#### Box 9: Success Metrics

*"How will we measure success?"*

Combines Patton's User/Business Metrics with Cagan's Metrics/Revenue.

| Field | Description |
|-------|-------------|
| `primaryMetric` | North star metric |
| `userMetrics` | User behavior to track |
| `businessMetrics` | Business outcomes to measure |
| `leadingIndicators` | Early signals of success |
| `revenueModel` | Subscription, usage, etc. |
| `timeToValue` | When value is realized |

### Row 4: Validation

#### Box 10: Critical Requirements

*"What must be true for this to succeed?"*

From Cagan's Solution Requirements.

| Field | Description |
|-------|-------------|
| `mustHaveCapabilities` | Non-negotiable features |
| `technicalRequirements` | Technical needs |
| `integrationNeeds` | Required integrations |
| `complianceNeeds` | Regulatory/compliance |
| `keyDependencies` | External dependencies |
| `successConditions` | Conditions that must be met |

#### Box 11: Risks & Assumptions

*"What are we betting on and what could go wrong?"*

Combines Patton's Assumptions with Cagan's Key Assumptions.

| Field | Description |
|-------|-------------|
| `keyAssumptions` | Key assumptions to validate |
| `riskiestAssumption` | Most critical assumption |
| `risks` | Identified risks |
| `highestRisk` | Most critical risk |
| `mitigationStrategy` | How to mitigate top risks |
| `validationApproach` | How to validate assumptions |
| `experimentsNeeded` | Experiments to run |

#### Box 12: Recommendation

*"Given all the above, what's the recommendation?"*

Combines Cagan's Recommendation with Patton's Budget.

| Field | Description |
|-------|-------------|
| `decision` | go, no-go, conditional, defer |
| `rationale` | Why this decision |
| `confidence` | high, medium, low |
| `conditions` | Conditions for conditional go |
| `timeEstimate` | Time to deliver |
| `teamSize` | Team needed |
| `investmentAsk` | Budget requested |
| `nextSteps` | What to do next |
| `successCriteria` | How we'll know it worked |

## Prioritization

OpportunitySpec integrates with two prioritization frameworks for feature evaluation:

### RICE Scoring

RICE provides quantitative scoring based on Reach, Impact, Confidence, and Effort.

```go
// Set RICE score
spec.SetRICE(1000, prioritization.ImpactHigh, prioritization.ConfidenceHigh, 2.0)

// Check score
if spec.HasRICE() {
    score := spec.GetRICEScore()
    fmt.Printf("RICE Score: %.0f\n", score) // 1000
}
```

| Component | Description |
|-----------|-------------|
| Reach | Users affected per quarter |
| Impact | massive (3x), high (2x), medium (1x), low (0.5x), minimal (0.25x) |
| Confidence | high (100%), medium (80%), low (50%) |
| Effort | Person-months |

### Kano Model

Kano classifies features by customer satisfaction impact.

```go
// Set Kano from questionnaire responses
spec.SetKano(prioritization.KanoLike, prioritization.KanoDislike)

// Check category
if spec.HasKano() {
    category := spec.GetKanoCategory() // "performance"
}

// Convenience checks
if spec.IsMustHave() { /* Basic expectation */ }
if spec.IsDelighter() { /* Attractive feature */ }
if spec.IsPerformance() { /* More is better */ }
```

| Category | Priority | Description |
|----------|----------|-------------|
| Must-Be | 5 | Basic expectation - absence causes dissatisfaction |
| Performance | 4 | Linear satisfaction - more is better |
| Attractive | 3 | Delighter - unexpected positive surprise |
| Indifferent | 1 | No significant impact |
| Reverse | 0 | Unwanted - presence causes dissatisfaction |

### Prioritization Summary

```go
summary := spec.GetPrioritizationSummary()
// summary.HasRICE, summary.RICEScore
// summary.HasKano, summary.KanoCategory, summary.KanoPriority
```

See [Prioritization](prioritization.md) for detailed documentation on RICE and Kano frameworks.

## Usage Example

```go
spec := canvas.NewOpportunitySpec("os-mobile", "Mobile App OpportunitySpec")

// Row 1: Discovery
spec.UsersAndProblem = canvas.OSUsersAndProblem{
    PrimaryUsers: []canvas.OSUser{
        {Name: "Field Service Reps", Goals: []string{"Complete jobs faster"}},
    },
    ProblemStatement: "Field teams cannot access critical data without WiFi",
    CustomerPains:    []string{"Delayed decisions", "Manual workarounds"},
}

spec.CurrentSolutions = canvas.OSCurrentSolutions{
    DirectCompetitors: []canvas.OSCompetitor{
        {Name: "FieldPro", Strengths: []string{"Mobile-first"}},
    },
    Workarounds: []string{"Paper forms", "Phone calls to office"},
}

spec.SolutionIdeas = canvas.OSSolutionIdeas{
    Ideas: []canvas.OSSolutionIdea{
        {ID: "S1", Name: "Native App", Feasibility: "high"},
        {ID: "S2", Name: "PWA", Feasibility: "medium"},
    },
    RecommendedIdea: "S1",
}

// Row 2: Value
spec.UserValue = canvas.OSUserValue{
    ValueStatement: "Complete jobs 50% faster with offline access",
    KeyBenefits:    []string{"Work anywhere", "Real-time sync"},
}

spec.BusinessValue = canvas.OSBusinessValue{
    BusinessProblem:  "Losing deals to mobile-first competitors",
    BusinessOutcomes: []string{"Increase win rate 20%"},
    MarketSize:       "TAM: $5B, SAM: $1.2B",
}

spec.CompetitiveEdge = canvas.OSCompetitiveEdge{
    Differentiator:  "Deep enterprise integration",
    UnfairAdvantage: "Existing customer relationships",
}

// Row 3: Market
spec.MarketAndTiming = canvas.OSMarketAndTiming{
    PrimarySegment: "Enterprise field service",
    WhyNow:         "5G rollout enabling new use cases",
    UrgencyLevel:   "high",
}

spec.GoToMarket = canvas.OSGoToMarket{
    Strategy:    "Land and expand with existing customers",
    SalesModel:  "sales-led",
}

spec.SuccessMetrics = canvas.OSSuccessMetrics{
    PrimaryMetric:   "Weekly active mobile users",
    BusinessMetrics: []string{"Mobile-originated revenue"},
}

// Row 4: Validation
spec.CriticalRequirements = canvas.OSCriticalRequirements{
    MustHaveCapabilities: []string{"Offline mode", "Sync engine"},
}

spec.RisksAndAssumptions = canvas.OSRisksAndAssumptions{
    RiskiestAssumption: "Users will adopt mobile within 30 days",
    Risks: []canvas.OSRisk{
        {ID: "R1", Description: "Development longer than expected", Impact: "high"},
    },
}

spec.Recommendation = canvas.OSRecommendation{
    Decision:      "go",
    Rationale:     "Strong market need, clear differentiation",
    Confidence:    "high",
    TimeEstimate:  "6 months",
    InvestmentAsk: "$1.2M",
    NextSteps:     []string{"Form team", "Customer validation"},
}

// Check decision
if spec.IsGo() && spec.HasHighConfidence() {
    fmt.Println("Ready to proceed!")
}
```

## Integration with multispec

OpportunitySpec is designed to integrate with [multispec](https://github.com/plexusone/multispec) for LLM-as-a-Judge evaluation:

```bash
# Initialize with OpportunitySpec profile
multispec init my-feature --profile aws-feature

# Draft the OpportunitySpec
multispec draft opportunity-spec -p my-feature

# Evaluate against rubric
multispec eval opportunity-spec -p my-feature

# If passing, proceed to Press Release
multispec synthesize press -p my-feature
```

### Canonical Assets in prism-roadmap

| Asset | Path | Description |
|-------|------|-------------|
| Go Types | `canvas/opportunity_spec.go` | OpportunitySpec struct and helpers |
| Template | `templates/opportunity-spec.md` | Markdown template with placeholders |
| Rubric | `rubrics/opportunity-spec.rubric.yaml` | LLM-as-a-Judge evaluation criteria |

## Comparison

| Aspect | Opportunity Canvas | Opportunity Assessment | OpportunitySpec |
|--------|-------------------|----------------------|-----------------|
| Boxes | 10 (9+budget) | 10 | 12 |
| Focus | Discovery | Business case | Both |
| Output | Problem understanding | Go/no-go decision | Comprehensive evaluation |
| Best for | Early exploration | Investment decision | Feature-level opportunities |

## When to Use

Use **OpportunitySpec** when:

- Evaluating a feature-level opportunity (not a new product line)
- Need both discovery rigor and business case
- Bridging PM discovery with executive approval
- Feeding into multispec Working Backwards flow

Use **Opportunity Canvas** (Patton) when:

- Pure discovery/exploration phase
- Don't need business case rigor yet

Use **Opportunity Assessment** (Cagan) when:

- Evaluating a major new initiative
- Need formal investment approval
- Less focus on discovery, more on business viability

## Template

See `templates/opportunity-spec.md` for a complete markdown template.

## Rubric

See `rubrics/opportunity-spec.rubric.yaml` for LLM evaluation criteria with 5 weighted categories:

| Category | Weight | Description |
|----------|--------|-------------|
| Discovery Clarity | 20% | Problem definition, user identification, competitive analysis |
| Value Proposition | 25% | User value, business value, differentiation |
| Market & Timing | 20% | Market definition, timing rationale, success metrics |
| Validation Readiness | 20% | Requirements, risks, recommendation quality |
| Document Quality | 15% | Coherence, completeness, actionability |

## Rendering

The OpportunitySpec renders natively to SVG in its twelve-box (four-row) layout,
colored by row — Discovery, Value, Market, Validation:

```go
c := canvas.NewOpportunitySpec(spec)
svgOut, _ := render.Render(c, render.FormatSVG, render.DefaultOptions())
```

## See Also

- [Opportunity Canvas](opportunity.md) - Jeff Patton's discovery-focused canvas
- [Opportunity Assessment](opportunity-assessment.md) - Marty Cagan's SVPG framework
- [multispec AWS Workflow](https://productbuildershq.com/visionspec/frameworks/aws/) - Working Backwards integration
