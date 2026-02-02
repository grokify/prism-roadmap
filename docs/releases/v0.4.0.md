# Release Notes - v0.4.0

**Release Date:** 2026-01-29

## Overview

This release introduces a major overhaul of the goals system, replacing legacy fields with a proper OKR (Objectives and Key Results) structure. It also adds comprehensive tag-based filtering, roadmap swimlane tables, and new PRD sections for decision tracking and security documentation.

## Highlights

- **OKR Structure**: Replace BusinessObjectives, ProductGoals, and SuccessMetrics with industry-standard OKR format
- **Tag-Based Filtering**: New CLI `filter` command with OR/AND logic for filtering documents by tags
- **Roadmap Swimlane Tables**: Visual roadmap tables with phase columns and deliverable type rows
- **Extended PRD Sections**: OpenItems for decision tracking, SecurityModel for threat modeling

## Breaking Changes

The following fields have been replaced with the new OKR structure:

| Removed Field | Replaced With |
|---------------|---------------|
| `BusinessObjectives` | `Objectives.OKRs[].Objective` |
| `ProductGoals` | `Objectives.OKRs[].Objective` |
| `SuccessMetrics` | `Objectives.OKRs[].KeyResults` |

### Migration Guide

**Before (v0.3.x):**
```json
{
  "objectives": {
    "business_objectives": ["Increase revenue"],
    "product_goals": ["Launch MVP"],
    "success_metrics": ["10K users"]
  }
}
```

**After (v0.4.0):**
```json
{
  "objectives": {
    "okrs": [
      {
        "objective": {
          "id": "obj-1",
          "description": "Become market leader",
          "category": "business"
        },
        "key_results": [
          {
            "id": "kr-1",
            "description": "Increase market share",
            "metric": "Market share",
            "target": "20%"
          }
        ]
      }
    ]
  }
}
```

## New Features

### OKR Structure

The new OKR structure follows industry-standard methodology:

- **Objectives**: Qualitative, inspirational goals with owner and timeframe
- **Key Results**: Quantitative, measurable outcomes with baseline, target, and confidence
- **Phase Targets**: Link Key Results to roadmap phases for progress tracking

### CLI Filter Command

Filter PRD documents by tags:

```bash
# OR-logic: include items with any of the tags
srequirements prd filter input.json --include mvp --include phase-1

# AND-logic: include items with all tags
srequirements prd filter input.json --include mvp --include auth --all

# Output to file
srequirements prd filter input.json -i mvp -o filtered.json
```

### Tag Validation

Tags are validated to kebab-case format:
- Lowercase alphanumeric with hyphens
- No leading/trailing hyphens
- No consecutive hyphens
- Leading digits allowed (e.g., `2024-q1`)

### Roadmap Swimlane Tables

Generate visual roadmap tables:

```go
opts := prd.DefaultRoadmapTableOptions()
table := doc.Roadmap.ToSwimlaneTable(opts)

// With OKR swimlanes
tableWithOKRs := doc.ToSwimlaneTableWithOKRs(opts)
```

Output example:
```
| Swimlane       | **Phase 1**<br>MVP | **Phase 2**<br>Beta |
|----------------|---------------------|---------------------|
| **Features**   | • Auth<br>• Search  | • Dashboard         |
| **Objectives** | • O1: Market leader |                     |
| **Key Results**| • KR1.1: Share→15%  | • KR1.1: Share→20%  |
```

### New PRD Sections

- **OpenItems**: Track pending decisions with options and tradeoffs
- **SecurityModel**: Document security architecture and threat model
- **CurrentState**: Document existing state before changes
- **Appendices**: Generic appendix system with helper functions

## Installation

```bash
go install github.com/grokify/structured-plan/cmd/srequirements@v0.4.0
```

## Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete details.
