# Completeness Check

The completeness check evaluates how thoroughly a PRD is filled out across all sections.

## Purpose

Unlike scoring (which evaluates quality), completeness checks:

- How many sections are filled
- Required vs optional section coverage
- Field-level completeness
- Overall document maturity

## Quick Start

```go
import "github.com/grokify/prism-execution/requirements/prd"

report := prd.CheckCompleteness(doc)

fmt.Printf("Overall: %.0f%%\n", report.OverallScore*100)
fmt.Printf("Grade: %s\n", report.Grade)
fmt.Printf("Status: %s\n", report.Status)
```

## CompletenessReport Structure

```go
type CompletenessReport struct {
    OverallScore  float64            `json:"overall_score"`  // 0.0-1.0
    Grade         string             `json:"grade"`          // A, B, C, D, F
    Status        string             `json:"status"`         // Complete, Partial, Minimal
    Sections      []SectionReport    `json:"sections"`
    MissingSections []string         `json:"missing_sections,omitempty"`
    Recommendations []string         `json:"recommendations,omitempty"`
}

type SectionReport struct {
    Name          string  `json:"name"`
    Score         float64 `json:"score"`
    Required      bool    `json:"required"`
    FieldsFilled  int     `json:"fields_filled"`
    FieldsTotal   int     `json:"fields_total"`
    Issues        []string `json:"issues,omitempty"`
}
```

## Grading Scale

| Score | Grade | Status |
|-------|-------|--------|
| ≥90% | A | Complete |
| ≥80% | B | Near Complete |
| ≥70% | C | Partial |
| ≥60% | D | Partial |
| <60% | F | Minimal |

## Section Weights

| Section | Weight | Required |
|---------|--------|----------|
| Metadata | 10% | Yes |
| Executive Summary | 15% | Yes |
| Objectives | 15% | Yes |
| Personas | 15% | Yes |
| Requirements | 20% | Yes |
| Roadmap | 10% | No |
| Risks | 5% | No |
| Technical Architecture | 5% | No |
| Goals Alignment | 5% | No |

## Example Report

```go
report := prd.CheckCompleteness(doc)

// Overall status
fmt.Printf("Completeness: %.0f%% (%s)\n",
    report.OverallScore*100, report.Grade)

// Section breakdown
for _, section := range report.Sections {
    icon := "✓"
    if section.Score < 0.7 {
        icon = "!"
    }
    fmt.Printf("%s %s: %.0f%% (%d/%d fields)\n",
        icon,
        section.Name,
        section.Score*100,
        section.FieldsFilled,
        section.FieldsTotal,
    )
}

// Recommendations
fmt.Println("\nRecommendations:")
for _, rec := range report.Recommendations {
    fmt.Printf("  - %s\n", rec)
}
```

### Output

```
Completeness: 72% (C)

✓ Metadata: 100% (5/5 fields)
✓ Executive Summary: 80% (4/5 fields)
! Objectives: 60% (3/5 fields)
✓ Personas: 85% (6/7 fields)
! Requirements: 50% (5/10 fields)
! Roadmap: 0% (0/4 fields)
! Risks: 33% (1/3 fields)

Recommendations:
  - Add success metrics to objectives
  - Add acceptance criteria to requirements
  - Define roadmap phases
  - Add risk mitigation strategies
```

## Checking Individual Sections

```go
// Check metadata completeness
metadataScore := prd.CheckMetadataCompleteness(doc.Metadata)

// Check persona completeness
for _, persona := range doc.Personas {
    score := prd.CheckPersonaCompleteness(persona)
    if score < 0.7 {
        fmt.Printf("Incomplete persona: %s (%.0f%%)\n",
            persona.Name, score*100)
    }
}

// Check requirements completeness
for _, req := range doc.Requirements.Functional {
    if !prd.IsRequirementComplete(req) {
        fmt.Printf("Incomplete: %s\n", req.Title)
    }
}
```

## Format Report

Generate formatted output:

```go
// Plain text
text := prd.FormatCompletenessReport(report, "text")
fmt.Println(text)

// Markdown
markdown := prd.FormatCompletenessReport(report, "markdown")

// JSON
json := prd.FormatCompletenessReport(report, "json")
```

## Persona Completeness

Personas are evaluated for:

| Field | Weight |
|-------|--------|
| Name | Required |
| Role | Required |
| Goals | 20% |
| Pain Points | 20% |
| Description | 10% |
| Demographics | 10% |
| Technical Proficiency | 10% |

```go
func IsPersonaComplete(p Persona) bool {
    return len(p.Goals) > 0 &&
           len(p.PainPoints) > 0 &&
           p.Role != ""
}
```

## Requirements Completeness

Requirements are evaluated for:

| Field | Weight |
|-------|--------|
| Title | Required |
| Description | Required |
| Priority | 20% |
| Acceptance Criteria | 30% |
| Persona Links | 10% |

## NFR Category Coverage

Non-functional requirements should cover key categories:

```go
categories := prd.GetNFRCategoryCoverage(doc.Requirements.NonFunctional)

// Returns map of coverage
// {
//   "performance": true,
//   "security": true,
//   "scalability": false,
//   ...
// }
```

## Filter by Priority

Focus on high-priority sections:

```go
report := prd.CheckCompleteness(doc)

// Filter to only must-have sections
criticalSections := prd.FilterByPriority(report.Sections, "critical")

for _, s := range criticalSections {
    if s.Score < 1.0 {
        fmt.Printf("Critical section incomplete: %s\n", s.Name)
    }
}
```

## Completeness vs Scoring vs Validation

| Check | Purpose | Output |
|-------|---------|--------|
| **Validation** | Structure correctness | Valid/Invalid + Errors |
| **Completeness** | Document coverage | Percentage + Grade |
| **Scoring** | Content quality | Score + Decision |

```go
// 1. First validate structure
result := prd.Validate(doc)
if !result.Valid {
    // Fix structural issues first
}

// 2. Then check completeness
report := prd.CheckCompleteness(doc)
if report.Grade < "C" {
    // Fill in missing sections
}

// 3. Finally score quality
scores := prd.Score(doc)
if scores.Decision == prd.ReviewReject {
    // Improve content quality
}
```

## Next Steps

- [Scoring](scoring.md)
- [Persona Library](persona-library.md)
- [PRD Documentation](../documents/prd.md)
