# PRD Examples

This page provides complete PRD examples for common scenarios.

## Minimal PRD

The smallest valid PRD:

```go
package main

import "github.com/grokify/prism-execution/requirements/prd"

func main() {
    doc := prd.New("PRD-001", "My Product",
        prd.Person{Name: "Author"})

    doc.ExecutiveSummary.ProblemStatement = "Users need X"
    doc.ExecutiveSummary.ProposedSolution = "Build Y"

    doc.Personas = []prd.Persona{{
        ID:        "PER-1",
        Name:      "User",
        Role:      "End User",
        IsPrimary: true,
    }}

    prd.Save(doc, "minimal.prd.json")
}
```

## Feature PRD

A typical feature PRD:

```go
doc := prd.New(prd.GenerateID(), "Dark Mode Feature",
    prd.Person{Name: "Jane Doe", Role: "PM"})

doc.ExecutiveSummary = prd.ExecutiveSummary{
    ProblemStatement: "Users report eye strain during extended use, " +
        "especially in low-light environments. 35% of support tickets " +
        "mention this issue.",
    ProposedSolution: "Implement a dark mode theme that users can " +
        "enable manually or auto-switch based on system preferences.",
    ExpectedOutcomes: []string{
        "Reduce eye-strain complaints by 80%",
        "Increase evening session duration by 25%",
    },
}

doc.Personas = []prd.Persona{
    {
        ID:        "PER-1",
        Name:      "Night Owl Developer",
        Role:      "Software Developer",
        IsPrimary: true,
        Goals:     []string{"Code comfortably at night", "Reduce eye strain"},
        PainPoints: []string{
            "Bright screens cause headaches",
            "Have to use browser extensions for dark mode",
        },
    },
}

doc.Objectives = prd.Objectives{
    ProductGoals: []prd.Objective{
        {ID: "PG-1", Description: "Reduce eye-strain support tickets by 80%"},
        {ID: "PG-2", Description: "Increase user retention for night sessions"},
    },
    SuccessMetrics: []prd.SuccessMetric{
        {
            ID:              "SM-1",
            Name:            "Eye Strain Tickets",
            Target:          "<10/month",
            CurrentBaseline: "50/month",
        },
        {
            ID:     "SM-2",
            Name:   "Dark Mode Adoption",
            Target: ">40% of users",
        },
    },
}

doc.Requirements = prd.Requirements{
    Functional: []prd.FunctionalRequirement{
        {
            ID:          "FR-1",
            Title:       "Manual Dark Mode Toggle",
            Description: "Users can toggle dark mode from settings",
            Priority:    prd.PriorityHigh,
            MoSCoW:      prd.MoSCoWMust,
        },
        {
            ID:          "FR-2",
            Title:       "System Preference Sync",
            Description: "Auto-switch based on OS dark mode setting",
            Priority:    prd.PriorityMedium,
            MoSCoW:      prd.MoSCoWShould,
        },
    },
    NonFunctional: []prd.NonFunctionalRequirement{
        {
            ID:          "NFR-1",
            Category:    prd.NFRUsability,
            Description: "Theme switch must complete in <100ms",
        },
    },
}

doc.OutOfScope = []string{
    "Custom color themes",
    "Per-page theme settings",
}
```

## Platform PRD

A larger platform initiative:

```go
doc := prd.New(prd.GenerateID(), "Customer Portal 2.0",
    prd.Person{Name: "Jane Doe", Role: "Senior PM", Email: "jane@example.com"})

doc.Metadata.Tags = []string{"platform", "customer-experience", "Q1-2025"}

doc.ExecutiveSummary = prd.ExecutiveSummary{
    ProblemStatement: "Current portal has 40% bounce rate. Users spend " +
        "average 15 minutes on tasks that should take 2 minutes. " +
        "This costs us $2M annually in support overhead.",
    ProposedSolution: "Complete redesign with modern tech stack, " +
        "optimized UX, and self-service capabilities.",
    ExpectedOutcomes: []string{
        "Reduce bounce rate to <20%",
        "Reduce average task time to <3 minutes",
        "Cut support ticket volume by 50%",
    },
    TargetAudience:   "Enterprise customers (ARR >$50K)",
    ValueProposition: "Fastest, easiest way to manage your account",
}

// Extended problem definition with evidence
doc.Problem = &prd.ProblemDefinition{
    Statement: "Portal performance and usability issues...",
    Impact:    "$2M annual support cost, 15% churn attributed to poor UX",
    Evidence: []prd.Evidence{
        {
            Type:     prd.EvidenceAnalytics,
            Source:   "Google Analytics",
            Summary:  "40% bounce rate, 3.2s avg load time",
            Strength: prd.StrengthHigh,
        },
        {
            Type:     prd.EvidenceSurvey,
            Source:   "Q4 NPS Survey",
            Summary:  "Portal UX cited by 45% of detractors",
            Strength: prd.StrengthHigh,
        },
    },
}

doc.Personas = []prd.Persona{
    {
        ID:          "PER-1",
        Name:        "Power User Pat",
        Role:        "Account Manager",
        IsPrimary:   true,
        Description: "Uses portal daily to manage client accounts",
        Goals: []string{
            "Quickly access client data",
            "Generate reports for meetings",
            "Update account settings efficiently",
        },
        PainPoints: []string{
            "3+ second page loads",
            "5+ clicks to complete common tasks",
            "Can't find search results",
        },
        TechnicalProficiency: prd.ProficiencyMedium,
    },
    {
        ID:          "PER-2",
        Name:        "Occasional User Oliver",
        Role:        "Customer",
        IsPrimary:   false,
        Description: "Logs in monthly to check statements",
        Goals:       []string{"Download invoices", "Update payment method"},
        PainPoints:  []string{"Forgets password", "Can't find billing section"},
    },
}

doc.Objectives = prd.Objectives{
    BusinessObjectives: []prd.Objective{
        {
            ID:          "BO-1",
            Description: "Reduce customer churn by 15%",
            AlignedWith: "O1",  // OKR reference
        },
        {
            ID:          "BO-2",
            Description: "Cut support costs by $1M annually",
        },
    },
    ProductGoals: []prd.Objective{
        {ID: "PG-1", Description: "Achieve <1s page load time"},
        {ID: "PG-2", Description: "Reach 90% task success rate"},
    },
    SuccessMetrics: []prd.SuccessMetric{
        {
            ID:              "SM-1",
            Name:            "Bounce Rate",
            Metric:          "bounce_rate",
            Target:          "<20%",
            CurrentBaseline: "40%",
            MeasurementMethod: "Google Analytics",
        },
        {
            ID:     "SM-2",
            Name:   "Task Completion Time",
            Target: "<3 minutes",
            CurrentBaseline: "15 minutes",
        },
    },
}

doc.Requirements = prd.Requirements{
    Functional: []prd.FunctionalRequirement{
        {
            ID:          "FR-1",
            Title:       "Instant Search",
            Description: "Full-text search across all user data",
            Priority:    prd.PriorityCritical,
            MoSCoW:      prd.MoSCoWMust,
            AcceptanceCriteria: []string{
                "Results appear within 500ms",
                "Supports fuzzy matching",
                "Highlights matching terms",
            },
            PersonaIDs: []string{"PER-1"},
        },
        {
            ID:          "FR-2",
            Title:       "Dashboard Redesign",
            Description: "New dashboard with most-used actions",
            Priority:    prd.PriorityCritical,
            MoSCoW:      prd.MoSCoWMust,
        },
    },
    NonFunctional: []prd.NonFunctionalRequirement{
        {
            ID:          "NFR-1",
            Category:    prd.NFRPerformance,
            Description: "Page load time <1 second (P95)",
        },
        {
            ID:          "NFR-2",
            Category:    prd.NFRAvailability,
            Description: "99.9% uptime SLA",
        },
        {
            ID:          "NFR-3",
            Category:    prd.NFRSecurity,
            Description: "SOC2 Type 2 compliant",
        },
    },
}

doc.Roadmap = prd.Roadmap{
    Phases: []prd.Phase{
        {
            ID:        "P1",
            Name:      "Foundation",
            Type:      prd.PhaseTypeQuarter,
            StartDate: "2025-01-01",
            EndDate:   "2025-03-31",
            Goals:     []string{"Infrastructure", "Auth system"},
            Deliverables: []prd.Deliverable{
                {ID: "D1", Title: "New API Gateway"},
                {ID: "D2", Title: "SSO Integration"},
            },
        },
        {
            ID:        "P2",
            Name:      "Core Features",
            Type:      prd.PhaseTypeQuarter,
            StartDate: "2025-04-01",
            EndDate:   "2025-06-30",
            Goals:     []string{"Dashboard", "Search"},
        },
    },
}

doc.Risks = []prd.Risk{
    {
        ID:          "R-1",
        Description: "Legacy data migration complexity",
        Impact:      prd.RiskImpactHigh,
        Probability: prd.RiskProbabilityMedium,
        Mitigation:  "Phased migration with rollback capability",
    },
}

// Goals alignment
doc.Goals = &prd.GoalsAlignment{
    OKR: &okr.OKRDocument{
        Objectives: []okr.Objective{
            {
                ID:    "O1",
                Title: "Delight customers with self-service",
                KeyResults: []okr.KeyResult{
                    {Title: "Reduce bounce to <20%", Target: "20%"},
                    {Title: "Achieve 90% task success", Target: "90%"},
                },
            },
        },
    },
    AlignedObjectives: map[string]string{
        "BO-1": "O1",
        "SM-1": "O1-KR1",
    },
}
```

## JSON Output

All examples can be saved as JSON:

```go
prd.Save(doc, "my-prd.prd.json")
```

Example output structure:

```json
{
  "metadata": {
    "id": "PRD-2025-022",
    "title": "Customer Portal 2.0",
    "version": "1.0.0",
    "status": "draft",
    "created_at": "2025-01-22T10:00:00Z",
    "authors": [{"name": "Jane Doe", "role": "Senior PM"}],
    "tags": ["platform", "customer-experience"]
  },
  "executive_summary": {
    "problem_statement": "Current portal has 40% bounce rate...",
    "proposed_solution": "Complete redesign with modern tech stack...",
    "expected_outcomes": ["Reduce bounce rate to <20%", "..."]
  },
  "personas": [...],
  "objectives": {...},
  "requirements": {...},
  "roadmap": {...},
  "risks": [...],
  "goals": {
    "okr": {...},
    "aligned_objectives": {"BO-1": "O1"}
  }
}
```

## Next Steps

- [Integration Examples](integration.md)
- [PRD Documentation](../documents/prd.md)
- [Quick Start](../getting-started/quickstart.md)
