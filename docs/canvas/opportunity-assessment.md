# Opportunity Assessment

The Opportunity Assessment follows Marty Cagan's SVPG 10-question framework for evaluating product opportunities before committing resources.

## Overview

This canvas implements the framework from [Assessing Product Opportunities](https://www.svpg.com/assessing-product-opportunities/) by Marty Cagan of Silicon Valley Product Group (SVPG).

Unlike Jeff Patton's Opportunity Canvas which focuses on discovery and user value, Cagan's framework emphasizes business case rigor and go/no-go decision making.

## The 10 Questions

| # | Question | Box Name |
|---|----------|----------|
| 1 | Exactly what problem will this solve? | Value Proposition |
| 2 | For whom do we solve that problem? | Target Market |
| 3 | How big is the opportunity? | Market Size |
| 4 | What alternatives are out there? | Competitive Landscape |
| 5 | Why are we best suited to pursue this? | Our Differentiator |
| 6 | Why now? | Market Window |
| 7 | How will we get this product to market? | Go-to-Market Strategy |
| 8 | How will we measure success/make money? | Metrics/Revenue Strategy |
| 9 | What factors are critical to success? | Solution Requirements |
| 10 | Given the above, what's the recommendation? | Go/No-Go Decision |

## Grid Layout

```
| 1. Value Proposition     | 2. Target Market         |
| 3. Market Size           | 4. Competitive Landscape |
| 5. Our Differentiator    | 6. Market Window         |
| 7. Go-to-Market Strategy | 8. Metrics/Revenue       |
| 9. Solution Requirements | 10. Recommendation       |
```

## Go Types

```go
type OpportunityAssessment struct {
    Metadata             Metadata
    ValueProposition     OAValueProposition     // Box 1
    TargetMarket         OATargetMarket         // Box 2
    MarketSize           OAMarketSize           // Box 3
    CompetitiveLandscape OACompetitiveLandscape // Box 4
    Differentiator       OADifferentiator       // Box 5
    MarketWindow         OAMarketWindow         // Box 6
    GoToMarket           OAGoToMarketStrategy   // Box 7
    MetricsRevenue       OAMetricsRevenue       // Box 8
    SolutionRequirements OASolutionRequirements // Box 9
    Recommendation       OARecommendation       // Box 10
    PRDRef               *PRDReference
}
```

## Box Details

### Box 1: Value Proposition

*"Exactly what problem will this solve?"*

| Field | Description |
|-------|-------------|
| `problemStatement` | Clear statement of the problem |
| `customerPains` | Specific pain points addressed |
| `desiredOutcomes` | What customers want to achieve |
| `currentState` | How things work today |
| `futureState` | How things will work with solution |

### Box 2: Target Market

*"For whom do we solve that problem?"*

| Field | Description |
|-------|-------------|
| `primarySegment` | Main target segment |
| `secondarySegments` | Additional segments |
| `personas` | Specific user personas |
| `industries` | Target industries |
| `geography` | Geographic focus |
| `companySize` | SMB, Mid-market, Enterprise |

### Box 3: Market Size

*"How big is the opportunity?"*

| Field | Description |
|-------|-------------|
| `tam` | Total Addressable Market |
| `sam` | Serviceable Addressable Market |
| `som` | Serviceable Obtainable Market |
| `growthRate` | Market growth rate |
| `marketTrends` | Key trends affecting size |
| `dataSources` | Where estimates come from |

### Box 4: Competitive Landscape

*"What alternatives are out there?"*

| Field | Description |
|-------|-------------|
| `directCompetitors` | Head-to-head competitors |
| `indirectCompetitors` | Alternative solutions |
| `substitutes` | Non-product alternatives |
| `marketDynamics` | Competitive dynamics |
| `entryBarriers` | Barriers to entry |

### Box 5: Our Differentiator

*"Why are we best suited to pursue this?"*

| Field | Description |
|-------|-------------|
| `coreStrengths` | What we do well |
| `uniqueCapabilities` | What only we can do |
| `strategicAssets` | Assets we can leverage |
| `teamExpertise` | Relevant team experience |
| `technologyEdge` | Technical advantages |
| `unfairAdvantage` | What's hard to replicate |

### Box 6: Market Window

*"Why now?"*

| Field | Description |
|-------|-------------|
| `whyNow` | Primary reason for timing |
| `marketTriggers` | Events creating opportunity |
| `technologyShifts` | Tech changes enabling this |
| `regulatoryChanges` | Regulatory drivers |
| `competitorMoves` | Competitor timing |
| `windowDuration` | How long window is open |
| `urgencyLevel` | high, medium, low |

### Box 7: Go-to-Market Strategy

*"How will we get this product to market?"*

| Field | Description |
|-------|-------------|
| `strategy` | Overall GTM approach |
| `channels` | Distribution channels |
| `salesModel` | Self-serve, sales-led, PLG |
| `partnerships` | Key partnerships needed |
| `launchApproach` | Beta, GA, etc. |
| `customerAcquisition` | How to acquire customers |
| `estimatedCAC` | Customer acquisition cost |

### Box 8: Metrics/Revenue Strategy

*"How will we measure success/make money?"*

| Field | Description |
|-------|-------------|
| `primaryMetric` | North star metric |
| `secondaryMetrics` | Supporting metrics |
| `leadingIndicators` | Early signals of success |
| `revenueModel` | Subscription, usage, etc. |
| `pricingStrategy` | How to price |
| `estimatedRevenue` | Revenue potential |
| `timeToRevenue` | When revenue starts |

### Box 9: Solution Requirements

*"What factors are critical to success?"*

| Field | Description |
|-------|-------------|
| `mustHaveCapabilities` | Non-negotiable features |
| `technicalRequirements` | Technical needs |
| `integrationNeeds` | Required integrations |
| `complianceNeeds` | Regulatory/compliance |
| `scalabilityNeeds` | Scale requirements |
| `keyDependencies` | External dependencies |
| `criticalRisks` | Risks to mitigate |

### Box 10: Recommendation

*"Given the above, what's the recommendation?"*

| Field | Description |
|-------|-------------|
| `decision` | go, no-go, conditional, defer |
| `rationale` | Why this decision |
| `confidence` | high, medium, low |
| `keyAssumptions` | Assumptions behind decision |
| `nextSteps` | What to do next |
| `conditions` | Conditions for conditional go |
| `reviewDate` | When to revisit if deferred |
| `timelineEstimate` | Time to deliver |
| `resourcesRequired` | Team/resources needed |
| `investmentAsk` | Budget requested |

## Usage Example

```go
oa := canvas.NewOpportunityAssessment("oa-mobile", "Mobile App Assessment")

oa.ValueProposition = canvas.OAValueProposition{
    ProblemStatement: "Field teams cannot access critical data without WiFi",
    CustomerPains:    []string{"Delayed decisions", "Manual workarounds"},
    DesiredOutcomes:  []string{"Real-time access", "Offline capability"},
}

oa.TargetMarket = canvas.OATargetMarket{
    PrimarySegment: "Enterprise field service teams",
    CompanySize:    "Enterprise (1000+ employees)",
}

oa.MarketSize = canvas.OAMarketSize{
    TAM: "$5B",
    SAM: "$1.2B",
    SOM: "$150M",
}

oa.Recommendation = canvas.OARecommendation{
    Decision:   "go",
    Rationale:  "Strong market need, clear differentiation, aligned with strategy",
    Confidence: "high",
    NextSteps:  []string{"Form team", "Validate with 5 customers"},
}

// Check decision
if oa.IsGo() && oa.HasHighConfidence() {
    fmt.Println("Ready to proceed!")
}
```

## Comparison with Opportunity Canvas

| Aspect | Opportunity Canvas (Patton) | Opportunity Assessment (Cagan) |
|--------|----------------------------|-------------------------------|
| Focus | Discovery & validation | Business case rigor |
| Output | Understanding opportunity | Go/no-go decision |
| Users | PM/UX discovery | Executive investment decision |
| Depth | User problems, solutions | Market size, competitive moats |
| Budget | Learning investment | Full investment ask |

## When to Use

Use the **Opportunity Assessment** when:

- Seeking investment approval for a new initiative
- Evaluating whether to enter a new market
- Making a go/no-go decision on a significant project
- Need to present a business case to executives

Use the **Opportunity Canvas** (Patton) when:

- Exploring and validating a problem space
- Running discovery with users
- Iterating on solution concepts
- Team alignment on opportunity understanding

## See Also

- [Opportunity Canvas](opportunity.md) - Jeff Patton's discovery-focused canvas
- [OpportunitySpec](opportunity-spec.md) - Merged Patton + Cagan framework
- [SVPG Article](https://www.svpg.com/assessing-product-opportunities/) - Original framework
