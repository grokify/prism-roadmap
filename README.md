# PRISM Roadmap

[![Go CI][go-ci-svg]][go-ci-url]
[![Go Lint][go-lint-svg]][go-lint-url]
[![Go SAST][go-sast-svg]][go-sast-url]
[![Docs][docs-godoc-svg]][docs-godoc-url]
[![Docs][docs-mkdoc-svg]][docs-mkdoc-url]
[![Visualization][viz-svg]][viz-url]
[![License][license-svg]][license-url]

 [go-ci-svg]: https://github.com/grokify/prism-roadmap/actions/workflows/go-ci.yaml/badge.svg?branch=main
 [go-ci-url]: https://github.com/grokify/prism-roadmap/actions/workflows/go-ci.yaml
 [go-lint-svg]: https://github.com/grokify/prism-roadmap/actions/workflows/go-lint.yaml/badge.svg?branch=main
 [go-lint-url]: https://github.com/grokify/prism-roadmap/actions/workflows/go-lint.yaml
 [go-sast-svg]: https://github.com/grokify/prism-roadmap/actions/workflows/go-sast-codeql.yaml/badge.svg?branch=main
 [go-sast-url]: https://github.com/grokify/prism-roadmap/actions/workflows/go-sast-codeql.yaml
 [docs-godoc-svg]: https://pkg.go.dev/badge/github.com/grokify/prism-roadmap
 [docs-godoc-url]: https://pkg.go.dev/github.com/grokify/prism-roadmap
 [docs-mkdoc-svg]: https://img.shields.io/badge/Go-dev%20guide-blue.svg
 [docs-mkdoc-url]: https://grokify.github.io/prism-roadmap/
 [viz-svg]: https://img.shields.io/badge/visualization-Go-blue.svg
 [viz-url]: https://mango-dune-07a8b7110.1.azurestaticapps.net/?repo=grokify%2Fprism-roadmap
 [loc-svg]: https://tokei.rs/b1/github/grokify/prism-roadmap
 [repo-url]: https://github.com/grokify/prism-roadmap
 [license-svg]: https://img.shields.io/badge/license-MIT-blue.svg
 [license-url]: https://github.com/grokify/prism-roadmap/blob/main/LICENSE

A unified planning system with Go data types, JSON serialization, and markdown generation. Supports requirements documents, goal frameworks, and roadmaps.

## PRISM Ecosystem

This project is part of the [PRISM ecosystem](https://github.com/grokify/prism):

| Project | Purpose |
|---------|---------|
| [prism-core](https://github.com/grokify/prism-core) | Shared primitives (Person, Status, Priority, Risk, Domain) |
| [prism-capability](https://github.com/grokify/prism-capability) | Define what capabilities exist |
| [prism-intelligence](https://github.com/grokify/prism-intelligence) | Measure maturity with SLIs/SLOs |
| **prism-roadmap** | Roadmaps, requirements, and goal frameworks (this project) |

```
prism-core (shared types)
         │
         ├── prism-capability ── "What we need"
         ├── prism-intelligence ── "How we measure"
         └── prism-roadmap ── "What we build"
```

## Overview

This library provides comprehensive, machine-readable formats for planning documents:

### Requirements Documents

- **MRD** - Market Requirements Document: Market analysis, competitive landscape, buyer personas, positioning
- **PRD** - Product Requirements Document: Personas, user stories, functional/non-functional requirements, roadmap
- **TRD** - Technical Requirements Document: Architecture, technology stack, APIs, security design, deployment

### Goal Frameworks

- **OKR** - Objectives and Key Results: Objectives with measurable key results and phase targets
- **V2MOM** - Vision, Values, Methods, Obstacles, Measures: Salesforce-style goal alignment

### Roadmap

- **Roadmap** - Standalone roadmaps with phases, deliverables, and swimlane visualization
- **Journey Roadmap** - Capability maturity evolution planning with periods, targets, and narratives

### Strategic Planning Canvases

Visual canvas frameworks for strategic planning and opportunity assessment:

| Canvas | Framework | Description |
|--------|-----------|-------------|
| **BMC** | Business Model Canvas (Osterwalder) | 9-block business model visualization |
| **OST** | Opportunity Solution Tree (Torres) | Tree structure for outcome-driven discovery |
| **Opportunity** | Opportunity Canvas (Patton) | 9-block opportunity assessment |
| **OpportunityAssessment** | SVPG Opportunity Assessment (Cagan) | 10-question go/no-go evaluation |
| **OpportunitySpec** | Merged Patton + Cagan | 12-box discovery + business case |
| **Feature** | Feature Canvas (Efimov) | Feature definition and validation |
| **Lean UX** | Lean UX Canvas (Gothelf) | Hypothesis-driven product design |

The **OpportunitySpec** is designed for feature-level opportunities, combining Patton's discovery rigor with Cagan's business case validation. It includes a canonical template and LLM evaluation rubric in `templates/` and `rubrics/`.

Each canvas supports multiple output formats:

- **D2** - D2 diagram language for high-quality diagrams
- **SVG** - Rendered vector graphics (native for BMC, OpportunitySpec, Opportunity Canvas, and Lean UX; others via D2 CLI)
- **Mermaid** - Mermaid diagram syntax for documentation
- **Lit/JSON** - Structured data for web components

### Journey Roadmaps

Plan capability maturity evolution over time with the `journey` package:

```go
import "github.com/grokify/prism-roadmap/journey"

roadmap := journey.JourneyRoadmap{
    ID:      "security-2026",
    Name:    "Security Capability Roadmap 2026",
    Periods: []journey.Period{
        {ID: "q1-2026", Name: "Q1 2026", StartDate: "2026-01-01", EndDate: "2026-03-31"},
        {ID: "q2-2026", Name: "Q2 2026", StartDate: "2026-04-01", EndDate: "2026-06-30"},
    },
    Capabilities: []journey.CapabilityJourney{
        {
            ID:           "sast",
            Name:         "Static Analysis",
            CurrentLevel: 2,
            TargetLevel:  4,
            Targets: []journey.CapabilityMaturityTarget{
                {PeriodID: "q1-2026", Level: 3, Confidence: "high"},
                {PeriodID: "q2-2026", Level: 4, Confidence: "medium"},
            },
        },
    },
}
```

Key types:

| Type | Description |
|------|-------------|
| `JourneyRoadmap` | Root type with periods, capabilities, dependencies, teams |
| `CapabilityJourney` | Tracks capability maturity progression |
| `CapabilityMaturityTarget` | Target level for a period with confidence |
| `Dependency` | Cross-capability or external dependencies |
| `Team` | Team hierarchy with capacity assignments |
| `Initiative` | Work items driving capability improvements |
| `JourneyNarrative` | Executive storytelling with chapters |

### Feature Prioritization

Prioritization frameworks integrated with OpportunitySpec:

| Framework | Type | Description |
|-----------|------|-------------|
| **RICE** | Quantitative | Score = (Reach × Impact × Confidence) / Effort |
| **Kano** | Qualitative | Categories: Must-Be, Performance, Attractive, Indifferent, Reverse |
| **MoSCoW** | Release Planning | Categories: Must Have, Should Have, Could Have, Won't Have |

```go
import "github.com/grokify/prism-roadmap/prioritization"

// RICE scoring
rice := prioritization.NewRICEScore("feature-1", 1000, ImpactHigh, ConfidenceHigh, 2.0)
rice.Calculate() // Score = 1000

// Kano classification
category := prioritization.ClassifyKano(KanoLike, KanoDislike) // Performance

// MoSCoW prioritization
priority := prioritization.MoSCoWMustHave
weight := priority.Weight() // Returns 4
```

### Market Signals & Effort Estimation

Aggregate customer demand and estimate implementation effort:

```go
import (
    "github.com/grokify/prism-roadmap/signal"
    "github.com/grokify/prism-roadmap/effort"
    "github.com/grokify/prism-roadmap/rmi"
)

// Market signal from customer ideas
sig := &signal.MarketSignal{
    TotalVotes:    150,
    CustomerCount: 25,
    TotalARR:      500000000, // $5M in cents
}
sig.CalculateScore() // Weighted score

// Effort estimation with complexity
est := &effort.EffortEstimate{
    PersonDays: 15,
    TShirtSize: effort.TShirtSizeMedium,
    Confidence: effort.ConfidenceMedium,
}

// Complexity factors
complexity := &effort.ComplexityFactors{
    NewArchitecture: false,
    NewDesignUX:     true,
    Dependencies:    []effort.Dependency{{Name: "Auth API", Type: "internal"}},
}

// Roadmap Item combining all components
item := &rmi.RoadmapItem{
    ID:           "rmi-feature-1",
    Name:         "Bulk Export",
    MoSCoW:       prioritization.MoSCoWShouldHave,
    MarketSignal: sig,
    Effort:       est,
    Complexity:   complexity,
    Status:       rmi.RMIStatusApproved,
}
```

The natural workflow from market to implementation:

**MRD** (Market) → **PRD** (Product) → **TRD** (Technical)

Each document type supports:

- Mandatory and optional sections for flexibility
- JSON serialization with Go types (camelCase field names)
- Markdown generation with Pandoc-compatible YAML frontmatter
- Validation of required fields
- Framework-agnostic goals (OKR or V2MOM)

## Installation

### Homebrew (macOS/Linux)

```bash
brew install grokify/tap/splan
```

### Go Install

```bash
go install github.com/grokify/prism-roadmap/cmd/splan@latest
```

### Download Binary

Pre-built binaries for Linux, macOS, and Windows are available on the [releases page](https://github.com/grokify/prism-roadmap/releases).

## CLI Usage

The `splan` CLI provides commands for working with planning documents:

```bash
# PRD commands
splan requirements prd generate <file.json>   # Generate markdown from PRD
splan requirements prd validate <file.json>   # Validate PRD structure
splan requirements prd check <file.json>      # Check PRD completeness
splan requirements prd score <file.json>      # Score PRD quality
splan requirements prd filter <file.json>     # Filter PRD by tags

# MRD commands
splan requirements mrd generate <file.json>   # Generate markdown from MRD
splan requirements mrd validate <file.json>   # Validate MRD structure

# TRD commands
splan requirements trd generate <file.json>   # Generate markdown from TRD
splan requirements trd validate <file.json>   # Validate TRD structure

# Utility commands
splan merge file1.json file2.json -o out.json # Merge JSON files
splan schema generate                          # Generate JSON schemas
```

**Shorthand:** Use `req` instead of `requirements` (e.g., `splan req prd generate`).

### Generate Options

```bash
splan req prd generate input.json -o output.md    # Custom output path
splan req prd generate input.json --no-frontmatter # Without YAML frontmatter
splan req prd generate input.json --margin 1in    # Custom page margin
splan req prd generate input.json --mainfont Arial # Custom font
splan req prd generate input.json --text-icons    # ASCII icons for Pandoc PDF
```

### Section Ordering (PRD only)

Control the order of sections in generated markdown:

```bash
# Use a PRD type template for section ordering
splan req prd generate input.json --type=strategy   # Context-first (CurrentState, Problem, Market early)
splan req prd generate input.json --type=feature    # User-needs-first (Problem, Personas, UserStories early)
splan req prd generate input.json --type=technical  # Architecture-focused (TechArchitecture early)

# Custom section order (comma-separated IDs)
splan req prd generate input.json --order=executiveSummary,problem,solution,objectives

# List available section IDs
splan req prd list-sections
```

### Check Options (PRD only)

```bash
splan req prd check input.json          # Human-readable completeness report
splan req prd check input.json --json   # JSON output for programmatic use
```

### Examples

```bash
# Validate and generate markdown
splan req mrd validate examples/agent-platform.mrd.json
splan req mrd generate examples/agent-platform.mrd.json

splan req prd validate examples/agent-control-plane.prd.json
splan req prd generate examples/agent-control-plane.prd.json
splan req prd check examples/agent-control-plane.prd.json

splan req trd validate examples/agent-control-plane.trd.json
splan req trd generate examples/agent-control-plane.trd.json
```

## Library Usage

### Requirements Documents

```go
package main

import (
    "encoding/json"
    "os"

    "github.com/grokify/prism-roadmap/requirements/prd"
    "github.com/grokify/prism-roadmap/requirements/mrd"
    "github.com/grokify/prism-roadmap/requirements/trd"
)

func main() {
    // Create a PRD programmatically
    doc := prd.Document{
        Metadata: prd.Metadata{
            ID:      "prd-001",
            Title:   "User Authentication System",
            Version: "1.0.0",
            Status:  prd.StatusDraft,
            Authors: []prd.Person{{Name: "Jane Doe"}},
        },
        ExecutiveSummary: prd.ExecutiveSummary{
            ProblemStatement: "Users need secure authentication",
            ProposedSolution: "Implement OAuth 2.0 with MFA",
        },
        // ... additional fields
    }

    // Generate markdown
    opts := prd.MarkdownOptions{
        IncludeFrontmatter: true,
        Margin:             "2cm",
    }
    markdown := doc.ToMarkdown(opts)

    // Or marshal to JSON
    data, _ := json.MarshalIndent(doc, "", "  ")
    os.WriteFile("output.prd.json", data, 0600)
}
```

### Goals (OKR and V2MOM)

The `goals` package provides a framework-agnostic interface for both OKR and V2MOM:

```go
import (
    "github.com/grokify/prism-roadmap/goals"
    "github.com/grokify/prism-roadmap/goals/okr"
    "github.com/grokify/prism-roadmap/goals/v2mom"
)

// Create OKR-based goals
okrSet := okr.OKRSet{
    Objectives: []okr.Objective{
        {
            ID:          "obj-1",
            Description: "Increase customer satisfaction",
            KeyResults: []okr.KeyResult{
                {ID: "kr-1", Description: "NPS score", Target: "> 50"},
            },
        },
    },
}
g := goals.NewOKR(okrSet)

// Or create V2MOM-based goals
v := v2mom.V2MOM{
    Vision: "Be the market leader",
    Methods: []v2mom.Method{
        {ID: "m-1", Description: "Launch enterprise features"},
    },
}
g := goals.NewV2MOM(v)

// Framework-agnostic access
for _, item := range g.GoalItems() {
    fmt.Println(item.Description())  // Works with both OKR and V2MOM
}

// Dynamic labels based on framework
fmt.Println(g.GoalLabel())   // "Objectives" (OKR) or "Methods" (V2MOM)
fmt.Println(g.ResultLabel()) // "Key Results" (OKR) or "Measures" (V2MOM)
```

### PRD with Goals

PRDs support framework-agnostic goals via the `ProductGoals` field:

```go
import (
    "github.com/grokify/prism-roadmap/requirements/prd"
    "github.com/grokify/prism-roadmap/goals"
)

doc := prd.Document{
    // ... metadata, executive summary, etc.
    ProductGoals: goals.NewOKR(okrSet),  // or goals.NewV2MOM(v2mom)
}

// Roadmap tables use correct terminology automatically
table := doc.ToSwimlaneTableWithGoals(opts)  // Uses "Objectives" or "Methods"
```

### Strategic Planning Canvases

```go
import (
    "github.com/grokify/prism-roadmap/canvas"
    "github.com/grokify/prism-roadmap/canvas/render"
    "github.com/grokify/prism-roadmap/canvas/render/d2"
)

// Create an Opportunity Solution Tree
ost := canvas.NewOpportunitySolutionTree("ost-1", "User Onboarding")
ost.Outcome = canvas.OSTOutcome{
    ID:          "O1",
    Description: "Increase activation to 60%",
    Opportunities: []canvas.OSTOpportunity{
        {
            ID:          "OP1",
            Description: "Users don't understand value prop",
            Solutions: []canvas.OSTSolution{
                {ID: "S1", Description: "Interactive tutorial"},
            },
        },
    },
}

// Wrap and render to D2
c := canvas.NewOST(ost)
d2Output, _ := render.Render(c, render.FormatD2, render.OSTOptions())

// Create an Opportunity Canvas with grid layout (BMC-style)
opp := canvas.NewOpportunityCanvas("opp-1", "Mobile App")
opp.Users = []canvas.User{{ID: "u1", Name: "Mobile Users"}}
opp.Problems = []canvas.Problem{{ID: "p1", Description: "Desktop-only access"}}
opp.SolutionIdeas = []string{"Native app", "PWA"}

// Render as grid (no arrows) or flow (with arrows)
gridD2, _ := render.Render(canvas.NewOpportunity(opp), render.FormatD2, render.OpportunityGridOptions())
flowD2, _ := render.Render(canvas.NewOpportunity(opp), render.FormatD2, render.OpportunityOptions())
```

See `examples/canvas/` for complete examples with rendered D2, SVG, Mermaid, and HTML outputs.

## Evaluation Integration

The library integrates with `structured-evaluation` for standardized quality reports:

```go
import "github.com/grokify/prism-roadmap/requirements/prd"

// Load and score a PRD
doc, _ := prd.Load("my-product.prd.json")

// Convert deterministic scoring to Rubric format
report := prd.ScoreToRubric(doc, "my-product.prd.json")

// Or generate a template for LLM judge evaluation
template := prd.GenerateEvaluationTemplate(doc, "my-product.prd.json")

// Validate the generated report (v0.7.0+)
result := rubric.ValidateReport(report)
if !result.Valid {
    // Handle validation errors
}
```

**Standard Evaluation Categories:**

| Category | Weight | Description |
|----------|--------|-------------|
| problem_definition | 20% | Problem statement clarity and evidence |
| solution_fit | 15% | Solution alignment with problem |
| user_understanding | 10% | Persona depth and user insights |
| market_awareness | 10% | Competitive analysis |
| scope_discipline | 10% | Clear objectives and boundaries |
| requirements_quality | 10% | Functional and non-functional specs |
| metrics_quality | 10% | Success metrics with targets |
| ux_coverage | 5% | Design and accessibility |
| technical_feasibility | 5% | Architecture and integrations |
| risk_management | 5% | Risk identification and mitigation |

## Document Types

### MRD - Market Requirements Document

Defines the market opportunity and business justification.

| Section | Required | Description |
|---------|----------|-------------|
| `metadata` | Yes | Document ID, title, version, authors |
| `executiveSummary` | Yes | Market opportunity, proposed offering, key findings |
| `marketOverview` | Yes | TAM/SAM/SOM, growth rate, trends |
| `targetMarket` | Yes | Primary/secondary segments, buyer personas |
| `competitiveLandscape` | Yes | Competitors, strengths/weaknesses, differentiators |
| `marketRequirements` | Yes | Market-level requirements with priorities |
| `positioning` | Yes | Positioning statement, key benefits |
| `goToMarket` | No | Launch strategy, pricing, distribution |
| `successMetrics` | Yes | Revenue targets, market share goals |
| `risks` | No | Market and competitive risks |
| `glossary` | No | Term definitions |

### PRD - Product Requirements Document

Defines what the product should do and for whom.

| Section | Required | Description |
|---------|----------|-------------|
| `metadata` | Yes | Document ID, title, version, authors |
| `executiveSummary` | Yes | Problem statement, proposed solution, outcomes |
| `objectives` | Yes | Business objectives, product goals, success metrics |
| `personas` | Yes | User personas with goals and pain points |
| `userStories` | Yes | User stories with acceptance criteria |
| `requirements.functional` | Yes | Functional requirements (MoSCoW priority) |
| `requirements.nonFunctional` | Yes | NFRs (performance, security, etc.) |
| `requirements.compliance` | No | Compliance requirements (GDPR, SOC2, HIPAA, etc.) |
| `roadmap` | Yes | Phases with deliverables and success criteria |
| `assumptions` | No | Assumptions, constraints, dependencies |
| `inScope` | No | Explicitly included items |
| `outOfScope` | No | Explicitly excluded items |
| `technicalArchitecture` | No | System overview, integrations, services, APIs |
| `relatedDocuments` | No | Links to related PRDs, TRDs, design docs |
| `risks` | No | Product and technical risks |
| `glossary` | No | Term definitions |
| `problem` | No | Extended problem definition with evidence and root causes |
| `market` | No | Market analysis with alternatives and differentiation |
| `solution` | No | Solution options with selection rationale |
| `decisions` | No | Decision records with alternatives considered |
| `reviews` | No | Review outcomes with quality scores |
| `revisionHistory` | No | Document revision history |
| `nonGoals` | No | Structured non-goals with rationale |
| `successMetrics` | No | Success metrics organized by type (north star, supporting, guardrail) |

### TRD - Technical Requirements Document

Defines how the product will be built.

| Section | Required | Description |
|---------|----------|-------------|
| `metadata` | Yes | Document ID, title, version, authors |
| `executiveSummary` | Yes | Purpose, scope, technical approach |
| `architecture` | Yes | Overview, principles, components, data flows |
| `technologyStack` | Yes | Languages, frameworks, databases, infrastructure |
| `apiSpecifications` | No | API definitions with endpoints |
| `dataModel` | No | Entities, attributes, data stores |
| `securityDesign` | Yes | AuthN, AuthZ, encryption, compliance |
| `performance` | Yes | Performance requirements and benchmarks |
| `scalability` | No | Horizontal/vertical scaling, limits |
| `deployment` | Yes | Environments, strategy, regions |
| `integrations` | No | External system integrations |
| `development` | No | Coding standards, branch strategy |
| `testing` | No | Testing strategy and coverage |
| `risks` | No | Technical risks |
| `glossary` | No | Term definitions |

## File Naming Convention

Use these extensions for automatic type detection:

- `*.prd.json` - Product Requirements Document
- `*.mrd.json` - Market Requirements Document
- `*.trd.json` - Technical Requirements Document

## PRD Details

### Personas

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique persona identifier |
| `name` | Yes | Persona name (e.g., "Developer Dan") |
| `role` | Yes | Job title or role |
| `description` | Yes | Background and context |
| `goals` | Yes | What they want to achieve |
| `painPoints` | Yes | Current frustrations |
| `behaviors` | No | Typical behaviors and patterns |
| `technicalProficiency` | No | Low, Medium, High, Expert |

### User Stories

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique story identifier |
| `personaId` | Yes | Reference to persona |
| `title` | Yes | Short descriptive title |
| `story` | Yes | "As a [persona], I want [goal] so that [reason]" |
| `acceptanceCriteria` | Yes | Testable conditions (Given/When/Then) |
| `priority` | Yes | Critical, High, Medium, Low |
| `phaseId` | Yes | Reference to roadmap phase |

### Roadmap and Swimlane Table

The PRD roadmap is rendered as a swimlane table with phases as columns and deliverable types as rows.

#### Roadmap Structure

```json
{
  "roadmap": {
    "phases": [
      {
        "id": "phase-1",
        "name": "MVP",
        "deliverables": [
          {
            "id": "d1",
            "title": "User Authentication",
            "type": "feature",
            "status": "completed"
          }
        ]
      }
    ]
  }
}
```

#### Ensuring Items Appear in the Roadmap Table

For a deliverable to appear in the swimlane table:

1. **Add to a Phase**: The deliverable must be in a phase's `deliverables` array
2. **Set the Type**: The `type` field determines which swimlane row the item appears in
3. **Set Status (optional)**: The `status` field adds a status icon

#### Deliverable Types (Swimlanes)

| Type Value | Swimlane Row | Description |
|------------|--------------|-------------|
| `feature` | Features | Product features and capabilities |
| `integration` | Integrations | Third-party integrations |
| `infrastructure` | Infrastructure | Platform, CI/CD, monitoring |
| `documentation` | Documentation | User guides, API docs |
| `milestone` | Milestones | Release milestones, checkpoints |
| `rollout` | Rollout | Customer/segment deployment phases |

#### Deliverable Status Icons

| Status Value | Icon | Description |
|--------------|------|-------------|
| `completed` | ✅ | Work is done |
| `in_progress` | 🔄 | Currently being worked on |
| `not_started` | ⏳ | Planned but not started |
| `blocked` | 🚫 | Blocked by dependency |

#### Example: Complete Deliverable

```json
{
  "id": "auth-feature",
  "title": "OAuth 2.0 Authentication",
  "description": "Implement OAuth 2.0 with support for Google and GitHub providers",
  "type": "feature",
  "status": "in_progress"
}
```

This appears in the **Features** row under the phase it belongs to, with a 🔄 icon.

#### Common Issues

| Problem | Cause | Solution |
|---------|-------|----------|
| Item not appearing | Missing or invalid `type` | Set `type` to a valid value |
| Item in wrong row | Wrong `type` value | Check spelling (e.g., `feature` not `Feature`) |
| Item in wrong column | Wrong phase | Move deliverable to correct phase's array |
| No status icon | Missing `status` field | Add `status` field with valid value |

#### Operational Rollout Swimlane

The `rollout` type enables tracking customer/segment deployments across phases. This is useful for phased go-to-market strategies where features are deployed to different customer segments over time.

**Recommended approach for calendar-tied phases:**

When phases represent calendar periods (quarters, months), place rollouts in the phase when deployment actually occurs, not when development completes:

| Swimlane | **Phase 1**<br>Q1 2026 | **Phase 2**<br>Q2 2026 | **Phase 3**<br>Q3 2026 |
|----------|------------------------|------------------------|------------------------|
| **Features** | • Auth<br>• Dashboard | • Reporting<br>• API v2 | • Analytics |
| **Rollout** | | • ✅ Auth → Enterprise<br>• 🔄 Dashboard → Pilot | • Reporting → All |

**Rationale:**

- **Reflects reality** - Development and rollout rarely happen in the same calendar window
- **Shows dependencies** - Clearly communicates "build first, then deploy"
- **Planning accuracy** - Resource allocation aligns with actual work timing

**Naming convention for rollout deliverables:**

Use the `→` notation to distinguish rollout targets:

```json
{
  "id": "rollout-auth-enterprise",
  "title": "Auth → Enterprise customers",
  "description": "Roll out Phase 1 Auth feature to enterprise segment",
  "type": "rollout",
  "status": "completed"
}
```

**Example: Multi-phase customer rollout**

```json
{
  "phases": [
    {
      "id": "phase-1",
      "name": "Q1 2026 - Build",
      "deliverables": [
        { "id": "f1", "title": "User Authentication", "type": "feature", "status": "completed" },
        { "id": "f2", "title": "Dashboard", "type": "feature", "status": "completed" }
      ]
    },
    {
      "id": "phase-2",
      "name": "Q2 2026 - Pilot",
      "deliverables": [
        { "id": "f3", "title": "Reporting", "type": "feature", "status": "in_progress" },
        { "id": "r1", "title": "Auth → Enterprise (Acme, TechCo)", "type": "rollout", "status": "completed" },
        { "id": "r2", "title": "Dashboard → Pilot customers", "type": "rollout", "status": "in_progress" }
      ]
    },
    {
      "id": "phase-3",
      "name": "Q3 2026 - GA",
      "deliverables": [
        { "id": "r3", "title": "Auth → All customers", "type": "rollout", "status": "not_started" },
        { "id": "r4", "title": "Dashboard → All customers", "type": "rollout", "status": "not_started" },
        { "id": "r5", "title": "Reporting → Enterprise", "type": "rollout", "status": "not_started" }
      ]
    }
  ]
}
```

### Non-Functional Requirements

| Category | Description | Example Metrics |
|----------|-------------|-----------------|
| `performance` | Response time, throughput | P95 < 200ms |
| `scalability` | Scaling capability | 10K concurrent users |
| `reliability` | Uptime, MTBF, MTTR | 99.9% uptime |
| `security` | AuthN, AuthZ, encryption | SOC 2 compliance |
| `multiTenancy` | Tenant isolation | Schema-per-tenant |
| `observability` | Logging, metrics, tracing | 100% trace coverage |
| `compliance` | Regulatory requirements | GDPR, HIPAA |

### Compliance Requirements

For tracking regulatory and standards compliance:

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier (e.g., "CR-001") |
| `title` | Yes | Requirement title |
| `description` | Yes | Detailed description |
| `category` | Yes | `data_privacy`, `security`, `healthcare`, `financial`, `accessibility`, `government`, `industry` |
| `standard` | Yes | Standard name (GDPR, SOC2, HIPAA, PCI-DSS, WCAG, FedRAMP) |
| `controlReference` | No | Specific control reference (e.g., "GDPR Article 17") |
| `geographicScope` | No | Applicable regions (EU, US, California, Global) |
| `effectiveDate` | No | When compliance is required |
| `priority` | Yes | MoSCoW priority |
| `phaseId` | Yes | Target roadmap phase |
| `status` | No | `not_started`, `in_progress`, `compliant`, `non_compliant` |
| `auditFrequency` | No | `annual`, `quarterly`, `continuous` |
| `evidenceRequirements` | No | Documentation needed for compliance |
| `certificationRequired` | No | Whether third-party certification is required |
| `thirdPartyAssessment` | No | Assessor type or name |
| `penalties` | No | Business risk of non-compliance |

**Compliance Categories:**

| Category | Description | Example Standards |
|----------|-------------|-------------------|
| `data_privacy` | Data protection regulations | GDPR, CCPA |
| `security` | Security certifications | SOC2, ISO 27001 |
| `healthcare` | Healthcare regulations | HIPAA, HITRUST |
| `financial` | Financial regulations | PCI-DSS, SOX |
| `accessibility` | Accessibility standards | WCAG, ADA |
| `government` | Government certifications | FedRAMP, StateRAMP |
| `industry` | Industry-specific standards | Varies by sector |

### Technical Architecture (Platform PRDs)

For platform and infrastructure PRDs, the `technicalArchitecture` section supports microservices documentation:

#### Services Inventory

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique service identifier |
| `name` | Yes | Service name |
| `description` | Yes | What the service does |
| `layer` | No | `control-plane`, `execution-plane`, `data-plane`, `gateway` |
| `protocol` | No | `REST`, `gRPC`, `GraphQL`, `WebSocket` |
| `language` | No | Primary programming language |
| `languageRationale` | No | Why this language was chosen |
| `responsibilities` | No | List of service responsibilities |
| `dependencies` | No | IDs of dependent services |

#### API Specifications

| Field | Required | Description |
|-------|----------|-------------|
| `name` | Yes | API name |
| `protocol` | Yes | `REST`, `gRPC`, `GraphQL`, `WebSocket` |
| `basePath` | No | Base URL path |
| `version` | No | API version |
| `endpoints` | No | List of endpoints (method, path, description, auth) |
| `openApiSpec` | No | URL to OpenAPI/Swagger spec |
| `protobufSpec` | No | URL to protobuf definitions |

#### Storage Architecture

| Field | Required | Description |
|-------|----------|-------------|
| `category` | Yes | `metadata`, `artifacts`, `state`, `cache`, `observability`, `audit`, `secrets` |
| `purpose` | Yes | What this storage is for |
| `technology` | Yes | Storage technology (DynamoDB, S3, etc.) |
| `encryption` | No | Encryption approach |
| `retention` | No | Data retention policy |
| `perTenant` | No | Whether storage is isolated per tenant |

#### GitOps Configuration

| Field | Required | Description |
|-------|----------|-------------|
| `enabled` | Yes | Whether GitOps is used |
| `provider` | No | GitOps provider (ArgoCD, Flux, etc.) |
| `workflow` | No | GitOps workflow description |
| `sourcesOfTruth` | No | List of artifacts and their locations (`git`, `s3`, `database`, `secrets-manager`, `registry`) |

#### Workflow Orchestration

| Field | Required | Description |
|-------|----------|-------------|
| `shortLived` | No | Engine for short-lived workflows (Step Functions, etc.) |
| `longRunning` | No | Engine for long-running workflows (Temporal, etc.) |
| `description` | No | Orchestration approach description |

### Related Documents

Link to related PRDs, TRDs, and design documents:

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Document identifier |
| `title` | Yes | Document title |
| `type` | Yes | `prd`, `trd`, `mrd`, `design-doc`, `rfc` |
| `relationship` | Yes | `child`, `parent`, `sibling`, `implements`, `supersedes`, `related` |
| `path` | No | File path to document |
| `url` | No | URL to document |
| `description` | No | Relationship context |

## MRD Details

### Market Size (TAM/SAM/SOM)

| Field | Required | Description |
|-------|----------|-------------|
| `value` | Yes | Market size (e.g., "$10B") |
| `year` | No | Reference year |
| `source` | No | Data source citation |
| `notes` | No | Additional context |

### Buyer Personas

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier |
| `name` | Yes | Persona name |
| `title` | Yes | Job title |
| `buyingRole` | Yes | Decision Maker, Influencer, User, Gatekeeper |
| `budgetAuthority` | Yes | Has budget authority (boolean) |
| `painPoints` | Yes | Business pain points |
| `goals` | Yes | Business goals |
| `buyingCriteria` | No | Purchase decision criteria |

### Competitors

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Unique identifier |
| `name` | Yes | Competitor name |
| `category` | No | Direct, Indirect, Substitute |
| `strengths` | Yes | Competitive strengths |
| `weaknesses` | Yes | Competitive weaknesses |
| `marketShare` | No | Market share percentage |
| `threatLevel` | No | High, Medium, Low |

## TRD Details

### Architecture Components

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | Component identifier |
| `name` | Yes | Component name |
| `description` | Yes | What it does |
| `type` | No | Service, Library, Database, Queue, etc. |
| `responsibilities` | No | List of responsibilities |
| `dependencies` | No | IDs of dependent components |
| `technology` | No | Implementation technology |

### API Specifications

| Field | Required | Description |
|-------|----------|-------------|
| `id` | Yes | API identifier |
| `name` | Yes | API name |
| `type` | Yes | REST, gRPC, GraphQL, WebSocket |
| `version` | No | API version |
| `baseUrl` | No | Base URL |
| `auth` | No | Authentication method |
| `endpoints` | No | List of endpoints |

### Security Design

| Field | Required | Description |
|-------|----------|-------------|
| `overview` | Yes | Security approach summary |
| `authentication` | No | AuthN method, provider, MFA |
| `authorization` | No | AuthZ model (RBAC, ABAC) |
| `encryption` | No | At-rest and in-transit encryption |
| `compliance` | No | Compliance standards (SOC2, GDPR) |

## PRD Completeness Check

The `splan req prd check` command analyzes a PRD for completeness and quality, providing:

- **Overall score** (0-100%) and letter grade (A-F)
- **Section-by-section breakdown** for both required and optional sections
- **Specific recommendations** prioritized by severity

### Scoring

The completeness check evaluates:

| Section | Weight | What's Checked |
|---------|--------|----------------|
| Metadata | 10% | ID, title, version, status, authors |
| Executive Summary | 10% | Problem statement depth, proposed solution, outcomes |
| Objectives | 10% | Business objectives, product goals, success metrics with targets |
| Personas | 10% | Number of personas, completeness of goals/pain points |
| User Stories | 10% | Acceptance criteria coverage, persona/phase linkage |
| Requirements | 10% | Functional/non-functional count, essential NFR categories |
| Roadmap | 10% | Phases with deliverables, success criteria, goals |
| Optional sections | 30% | 18 optional sections including assumptions, scope, architecture, risks, problem, market, solution, decisions, reviews, non-goals, success metrics, compliance requirements, requirements by phase |

### Example Output

```
=============================================================
PRD COMPLETENESS REPORT
=============================================================

Overall Score: 85.5% (Grade: B)
Required Sections: 7/7 complete
Optional Sections: 8/18 complete

-------------------------------------------------------------
SECTION BREAKDOWN
-------------------------------------------------------------

Required Sections:
  [+] Metadata                  100.0% (complete)
  [+] Executive Summary         100.0% (complete)
  [+] Objectives                100.0% (complete)
  [+] Personas                  100.0% (complete)
  [+] User Stories              100.0% (complete)
  [+] Requirements               83.3% (complete)
  [+] Roadmap                   100.0% (complete)

Optional Sections:
  [+] Assumptions & Constraints 100.0% (complete)
  [+] In Scope                  100.0% (complete)
  [+] Out of Scope              100.0% (complete)
  [~] Technical Architecture     50.0% (partial)
  [ ] UX Requirements             0.0% (missing)
  [+] Risks                     100.0% (complete)
  [+] Glossary                  100.0% (complete)
  [ ] Related Documents           0.0% (missing)
  [+] Problem Definition        100.0% (complete)
  [ ] Market Analysis             0.0% (missing)
  [ ] Solution                    0.0% (missing)
  [ ] Decisions                   0.0% (missing)
  [ ] Reviews                     0.0% (missing)
  [ ] Revision History            0.0% (missing)
  [ ] Non-Goals                   0.0% (missing)
  [ ] Success Metrics             0.0% (missing)

-------------------------------------------------------------
RECOMMENDATIONS
-------------------------------------------------------------

HIGH (should fix):
  [*] Requirements: Missing NFR categories: reliability

=============================================================
```

## PDF Generation

The generated markdown includes YAML frontmatter compatible with Pandoc.

```bash
# Generate markdown
splan req prd generate myproduct.prd.json -o myproduct.md

# Convert to PDF
pandoc myproduct.md -o myproduct.pdf
```

**Requirements:**

- [Pandoc](https://pandoc.org/installing.html)
- A LaTeX distribution (TeX Live, MacTeX, or MiKTeX)

### Handling Status Icons (Emoji)

The default output uses emoji status icons (✅, 🔄, ⏳, etc.) which display correctly in HTML but may not render in PDF output. You have two options:

**Option 1: Use text icons for PDF compatibility**

```go
opts := prd.DefaultMarkdownOptions()
opts.UseTextIcons = true  // Uses [DONE], [WIP], [TODO] instead of emoji
markdown := doc.ToMarkdown(opts)
```

Text icon mappings:

| Emoji | Text Icon |
|-------|-----------|
| ✅ | [DONE] |
| 🔄 | [WIP] |
| ⏳ | [TODO] |
| 🚫 | [BLOCKED] |
| ❌ | [MISSED] |

**Option 2: Use XeLaTeX with emoji-capable fonts**

```bash
pandoc myproduct.md -o myproduct.pdf --pdf-engine=xelatex \
  -V mainfont="Noto Sans" \
  -V monofont="Noto Sans Mono"
```

Note: This requires fonts with emoji support (e.g., Noto Color Emoji, Apple Color Emoji).

## Templates and Rubrics

The `templates/` and `rubrics/` directories contain canonical assets for LLM-assisted document authoring and evaluation:

### Templates

| Template | Description |
|----------|-------------|
| `opportunity-spec.md` | OpportunitySpec 12-box canvas template |

### Rubrics

| Rubric | Description |
|--------|-------------|
| `opportunity-spec.rubric.yaml` | LLM-as-a-Judge evaluation criteria for OpportunitySpec |

These assets are designed for integration with [multispec](https://github.com/plexusone/multispec) for structured document workflows.

## Examples

See the `examples/` directory for complete examples:

- `examples/agent-platform.mrd.json` - Market requirements for an AI governance platform
- `examples/agent-control-plane.prd.json` - Product requirements for the control plane
- `examples/agent-control-plane.trd.json` - Technical requirements for implementation

## TypeScript / JavaScript

For TypeScript/JavaScript projects, use the unified `@grokify/prism` npm package:

```bash
npm install @grokify/prism
```

```typescript
import { JourneyRoadmapSchema } from '@grokify/prism/schema/roadmap';
import { renderTimelineView } from '@grokify/prism/html/roadmap';

// Validate data
const roadmap = JourneyRoadmapSchema.parse(jsonData);

// Render to HTML
const html = renderTimelineView(roadmap);
```

See [@grokify/prism on npm](https://www.npmjs.com/package/@grokify/prism) for full documentation.

## References

### Requirements Documents

- [Modern Analyst - 9 Types of Requirements Documents](https://modernanalyst.com/Resources/Articles/tabid/115/ID/5464/9-Types-Of-Requirements-Documents-What-They-Mean-And-Who-Writes-Them.aspx)
- [Product School - PRD Template](https://productschool.com/blog/product-strategy/product-template-requirements-document-prd)
- [Atlassian - Product Requirements](https://www.atlassian.com/agile/product-management/requirements)

### Technical Documentation

- [AWS - Architecture Documentation](https://docs.aws.amazon.com/wellarchitected/latest/framework/welcome.html)
- [C4 Model - Software Architecture](https://c4model.com/)
- [ADR - Architecture Decision Records](https://adr.github.io/)

## License

MIT License
