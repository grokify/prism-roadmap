# Release Notes: v0.3.0

**Release Date:** 2026-01-26

## Overview

This release adds integration with `structured-evaluation` for standardized PRD quality reports, a new merge command for combining multiple JSON files, and enhanced user story fields.

## Highlights

- **Structured Evaluation Integration**: Convert deterministic scoring results to `EvaluationReport` format for consistent output across LLM-based and rule-based evaluations
- **Merge Command**: New CLI command to combine multiple PRD JSON files with deep merge support
- **Enhanced User Stories**: Additional fields for better Agile workflow integration

## New Features

### Evaluation Integration

The new evaluation integration enables PRD scoring results to be output in the standardized `EvaluationReport` format from `structured-evaluation`:

```go
import "github.com/grokify/structured-plan/prd"

// Convert deterministic scoring to EvaluationReport
doc, _ := prd.Load("my-product.prd.json")
report := prd.ScoreToEvaluationReport(doc, "my-product.prd.json")

// Or generate a template for LLM judge evaluation
template := prd.GenerateEvaluationTemplate(doc, "my-product.prd.json")
```

**Functions Added:**

| Function | Description |
|----------|-------------|
| `ScoreToEvaluationReport()` | Converts deterministic scoring results to EvaluationReport |
| `GenerateEvaluationTemplate()` | Creates empty template for LLM judge to fill in |
| `GenerateEvaluationTemplateWithWeights()` | Template with custom category weights |
| `StandardCategories()` | Returns 10 standard PRD evaluation categories |
| `CategoryDescriptions()` | Category ID to description map for LLM prompts |
| `CategoryOwners()` | Category ID to suggested owner map |
| `GetCategoriesFromDocument()` | Extracts categories including custom sections |

**Standard Categories (10 total):**

| Category | Weight | Owner |
|----------|--------|-------|
| problem_definition | 20% | problem-discovery |
| solution_fit | 15% | solution-ideation |
| user_understanding | 10% | user-research |
| market_awareness | 10% | market-intel |
| scope_discipline | 10% | prd-lead |
| requirements_quality | 10% | requirements |
| metrics_quality | 10% | metrics-success |
| ux_coverage | 5% | ux-journey |
| technical_feasibility | 5% | tech-feasibility |
| risk_management | 5% | risk-compliance |

### Merge Command

Combine multiple PRD JSON files into one with deep merge support:

```bash
# Merge multiple files
srequirements merge base.json overlay.json -o combined.json

# Merge with default output (merged.json)
srequirements merge part1.json part2.json part3.json
```

Deep merge behavior:
- Nested objects are recursively merged
- Later files override earlier values for scalar fields
- Arrays are replaced (not concatenated)

### Enhanced User Stories

User stories now support additional fields for Agile workflows:

```json
{
  "id": "US-001",
  "title": "User Authentication",
  "epic_id": "EPIC-AUTH",
  "story_points": 5,
  "labels": ["security", "mvp"],
  "dependencies": ["US-002", "US-003"]
}
```

## Dependencies

- Updated `structured-evaluation` to v0.2.0 (adds rubrics, judge metadata, pairwise comparison)

## Migration

No breaking changes. Existing code continues to work without modification.

## Installation

```bash
go get github.com/grokify/structured-plan@v0.3.0
```

Or install the CLI:

```bash
go install github.com/grokify/structured-plan/cmd/srequirements@v0.3.0
```
