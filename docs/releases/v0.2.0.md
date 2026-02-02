# Release Notes - v0.2.0

**Release Date:** January 25, 2026

## Overview

This release adds the **schema package** for programmatic JSON Schema access, implementing a **Go-first approach** where Go structs are the source of truth and JSON Schema is generated from them.

Key additions:

- Embedded JSON Schema via `//go:embed` for runtime access
- Schema generator using `invopop/jsonschema`
- CLI command for schema generation
- Extended PRD sections for multi-agent workflows

## Highlights

- **Programmatic Schema Access** - Import and use JSON Schema directly in Go code
- **Go-First Schema Generation** - Generate JSON Schema from Go types, ensuring consistency
- **Extended PRD Sections** - Support for problem, market, solution, decisions, reviews, and goals

## Installation

```bash
go get github.com/grokify/structured-plan@v0.2.0
```

## New Features

### Schema Embed Package

Access the PRD JSON Schema programmatically:

```go
import "github.com/grokify/structured-plan/schema"

// Get schema as string
schemaJSON := schema.PRDSchema()

// Get schema as bytes
schemaBytes := schema.PRDSchemaBytes()

// Get canonical schema URL
schemaID := schema.PRDSchemaID
// "https://github.com/grokify/structured-plan/schema/prd.schema.json"
```

### Schema Generator

Generate JSON Schema from Go types:

```go
import "github.com/grokify/structured-plan/schema"

gen := schema.NewGenerator()

// Generate schema
prdSchema, err := gen.GeneratePRDSchema()

// Write to file
err := gen.WritePRDSchema("prd.schema.json")

// Generate all schemas to directory
err := gen.GenerateAll("./schema/")
```

### CLI Schema Command

Generate schemas from the command line:

```bash
# Generate all schemas
srequirements schema generate

# Generate PRD schema only
srequirements schema generate --type prd

# Specify output location
srequirements schema generate --type prd -o ./schema/prd.schema.json
```

### Extended PRD Sections

New sections added to PRD schema for multi-agent PRD workflows:

| Section | Description |
|---------|-------------|
| `problem` | Evidence-backed problem definition with root causes |
| `market` | Competitive analysis, alternatives, differentiation |
| `solution` | Solution options with tradeoffs and selection rationale |
| `decisions` | Decision records with rationale and alternatives |
| `reviews` | Quality scores and review outcomes |
| `revision_history` | Audit trail of PRD changes |
| `goals` | V2MOM/OKR alignment references |

## Go-First Workflow

This release establishes the Go-first approach for schema management:

1. **Go structs are the source of truth** - Define types in Go with `json` tags
2. **Generate schema from Go types** - Use `schema.NewGenerator()` or CLI
3. **Validate generated schema** - Use `schemago lint` for Go-friendliness
4. **Embed for runtime access** - Use `//go:embed` for programmatic access

## Breaking Changes

None. This release is fully backward compatible.

## Dependencies

New dependency added:

- `github.com/invopop/jsonschema v0.13.0` - JSON Schema generation from Go types

## Commits

| Commit | Description |
|--------|-------------|
| [`c219bed`](https://github.com/grokify/structured-plan/commit/c219bed) | feat(schema): add schema embed package for programmatic access |
| [`699ef11`](https://github.com/grokify/structured-plan/commit/699ef11) | feat(schema): add JSON Schema generator from Go types |
| [`1970cb9`](https://github.com/grokify/structured-plan/commit/1970cb9) | feat(cli): add schema generate command |
| [`0911fd5`](https://github.com/grokify/structured-plan/commit/0911fd5) | feat(schema): add extended PRD sections to schema |

## Contributors

- John Wang (@grokify)
- Claude Opus 4.5 (Co-Author)

## Links

- [GitHub Repository](https://github.com/grokify/structured-plan)
- [Go Package Documentation](https://pkg.go.dev/github.com/grokify/structured-plan)
- [Changelog](CHANGELOG.md)
- [v0.1.0 Release Notes](RELEASE_NOTES_v0.1.0.md)
