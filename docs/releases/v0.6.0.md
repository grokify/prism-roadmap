# Release Notes - v0.6.0

**Release Date:** 2026-01-30

## Overview

This release standardizes all JSON field names to **camelCase** for consistency with JSON/JavaScript conventions and modern API standards (OpenAPI, GraphQL, Google APIs). It also adds Goals CLI commands (V2MOM and OKR) for validation, Marp slide generation, and template creation, along with JSON Schema generation support for both frameworks.

## Breaking Changes

### JSON Field Names Migrated to camelCase

All 252 JSON field names have been migrated from snake_case to camelCase. This is a **breaking change** that requires updating existing JSON files.

**Examples of changes:**

| Before (snake_case) | After (camelCase) |
|---------------------|-------------------|
| `created_at` | `createdAt` |
| `updated_at` | `updatedAt` |
| `user_stories` | `userStories` |
| `executive_summary` | `executiveSummary` |
| `problem_statement` | `problemStatement` |
| `proposed_solution` | `proposedSolution` |
| `expected_outcomes` | `expectedOutcomes` |
| `success_metrics` | `successMetrics` |
| `key_results` | `keyResults` |
| `non_functional` | `nonFunctional` |
| `out_of_scope` | `outOfScope` |
| `target_audience` | `targetAudience` |

See [MIGRATION_CAMELCASE.md](MIGRATION_CAMELCASE.md) for the complete field mapping.

## Migration Guide

### Automatic Migration Script

A jq-based migration script is provided to convert existing JSON files:

```bash
# Convert a single file (in place)
./scripts/migrate_json_to_camelcase.sh myproduct.prd.json

# Convert to a new file
./scripts/migrate_json_to_camelcase.sh myproduct.prd.json myproduct-new.prd.json

# Batch convert all PRD files
for f in *.prd.json; do
  ./scripts/migrate_json_to_camelcase.sh "$f"
done
```

**Requirements:** `jq` (install with `brew install jq`)

### Manual Migration

If you prefer manual migration, update JSON keys according to the pattern:

```json
// Before
{
  "executive_summary": {
    "problem_statement": "...",
    "proposed_solution": "..."
  },
  "user_stories": [...]
}

// After
{
  "executiveSummary": {
    "problemStatement": "...",
    "proposedSolution": "..."
  },
  "userStories": [...]
}
```

## Rationale

The migration to camelCase provides:

1. **JSON Convention Alignment**: camelCase is the de facto standard for JSON (derived from JavaScript)
2. **API Compatibility**: Matches OpenAPI, GraphQL, and Google API style guides
3. **Smaller Payloads**: No underscores means slightly smaller JSON files
4. **IDE Support**: Better autocomplete in JavaScript/TypeScript editors

## New Features

### V2MOM CLI Commands

The `splan goals v2mom` commands are now available for working with V2MOM (Vision, Values, Methods, Obstacles, Measures) documents:

```bash
# Validate a V2MOM JSON file
splan goals v2mom validate my-v2mom.json
splan goals v2mom validate my-v2mom.json --structure=nested

# Generate Marp presentation slides
splan goals v2mom generate marp my-v2mom.json
splan goals v2mom generate marp my-v2mom.json -o slides.md --theme=corporate

# Create a new V2MOM template
splan goals v2mom init
splan goals v2mom init --name "FY2026 Product Strategy" -o product-v2mom.json
```

**Structure modes:**
- `flat` - Traditional V2MOM with measures/obstacles at top level
- `nested` - OKR-aligned with measures under methods
- `hybrid` - Both levels allowed (default)

**Terminology modes:**

- `v2mom` - Use V2MOM terms: Methods, Measures, Obstacles
- `okr` - Use OKR terms: Objectives, Key Results, Risks
- `hybrid` - Show both terminologies

### OKR CLI Commands

The `splan goals okr` commands are now available for working with OKR (Objectives and Key Results) documents:

```bash
# Validate an OKR JSON file
splan goals okr validate my-okrs.json

# Generate Marp presentation slides
splan goals okr generate marp my-okrs.json
splan goals okr generate marp my-okrs.json -o slides.md --theme=corporate

# Create a new OKR template
splan goals okr init
splan goals okr init --name "Q1 2026 Goals" -o q1-okrs.json
```

### JSON Schema Generation

Generate JSON Schemas from Go types for validation:

```bash
# Generate all schemas (PRD, OKR, V2MOM)
splan schema generate -o ./schema/

# Generate specific schema
splan schema generate --type okr -o okr.schema.json
splan schema generate --type v2mom -o v2mom.schema.json
splan schema generate --type prd -o prd.schema.json
```

Embedded schemas are also available programmatically:

```go
import "github.com/grokify/structured-plan/schema"

// Get schemas as strings
prdSchema := schema.PRDSchema()
okrSchema := schema.OKRSchema()
v2momSchema := schema.V2MOMSchema()

// Get schemas as byte slices
prdBytes := schema.PRDSchemaBytes()
okrBytes := schema.OKRSchemaBytes()
v2momBytes := schema.V2MOMSchemaBytes()
```

## New Files

- `MIGRATION_CAMELCASE.md` - Complete migration documentation with field mapping
- `scripts/migrate_json_to_camelcase.sh` - jq-based migration script for user JSON files
- `scripts/migrate_to_camelcase.sh` - Migration script for Go source files (for reference)
- `schema/okr.schema.json` - JSON Schema for OKR documents
- `schema/v2mom.schema.json` - JSON Schema for V2MOM documents
- `examples/minimal.v2mom.json` - Minimal V2MOM example
- `examples/product.v2mom.json` - Product strategy V2MOM example
- `examples/agentplexus.v2mom.json` - Comprehensive V2MOM example

## Installation

### Homebrew (macOS/Linux)

```bash
brew upgrade grokify/tap/splan
```

### Go Install

```bash
go install github.com/grokify/structured-plan/cmd/splan@v0.6.0
```

### Go Module

```bash
go get github.com/grokify/structured-plan@v0.6.0
```

## Full Changelog

See [CHANGELOG.md](CHANGELOG.md) for complete details.
