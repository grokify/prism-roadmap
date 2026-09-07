# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go build ./...                                  # build all packages
go test ./...                                   # run all tests
go test ./assessment/ -run TestPDIMToRGT        # run a single test
golangci-lint run                               # lint (config: .golangci.yaml)
```

JSON Schema regeneration (Go structs are the source of truth; never hand-edit `*.schema.json`):

```bash
cd schema && go run cmd/generate/main.go        # regenerates schema/*.schema.json via invopop/jsonschema
```

Changelog and release docs (CHANGELOG.json is the source of truth; CHANGELOG.md is generated):

```bash
schangelog parse-commits --since=<tag>          # analyze commits since last release
schangelog validate CHANGELOG.json
schangelog generate CHANGELOG.json -o CHANGELOG.md
```

A release also needs `docs/releases/vX.Y.Z.md` (follow the previous version's format) and a nav entry in `mkdocs.yml` under `Releases:`. MkDocs is not installed locally — the site build is verified by CI.

**Caution:** root-level `generate`, `splan`, `srequirements`, and `validate` are checked-in compiled binaries, not source directories — never edit or grep them as code. Source lives in `cmd/` and `cli/`.

## Architecture

A unified planning-document library: Go types + JSON serialization + markdown/Marp generation + JSON Schema, part of the PRISM ecosystem (`prism-core` primitives, `prism-capability`, `prism-maturity`). Dependencies flow toward leaf primitives with no cycles.

- **Leaf packages**: `prioritization` (RICE/Kano/MoSCoW types), `signal` (MarketSignal), `effort`, `common` (Person/Decision/Risk over prism-core), `roadmap`, `journey`, `goals` (`okr`, `v2mom`).
- **rmi** — `RoadmapItem`, the core roadmap unit combining MoSCoW+RICE, MarketSignal, and Effort, plus a file-backed `Service`/`Set` CRUD layer used by both the CLI and the MCP server.
- **canvas** — strategic canvases (BMC, OST, Opportunity, SVPG OpportunityAssessment, Lean UX, …) with `render`/`export`.
- **requirements** — MRD/PRD/TRD document types; PRD has `render`, `render/marp`, `render/terminal`.
- **assessment** — evidence-backed opportunity prioritization IR (see below). Depends on canvas, goals/okr, prioritization, compass-rice, structured-evaluation.
- **schema** — JSON Schema hub: reflects Go types from assessment, canvas, effort, goals, journey, prd, rmi, signal into embedded `*.schema.json`.
- **rubrics** / **templates** — `go:embed` bundles exposed via `FS()`.
- **cli** — exported cobra tree (`cli.RootCmd`); **cmd/splan** wraps it (version injected via GoReleaser ldflags), **cmd/validate** is a separate CLI validating cross-file idea/goal/capability references.
- **mcp** — MCP server (official `modelcontextprotocol/go-sdk`) exposing RMI CRUD tools (`list_rmis`, `get_rmi`, `create_rmi`, …) over file-backed `rmi.Service`.

Each document type keeps rendering in parallel subpackages: `render` (markdown/terminal) and `render/marp` (presentations via structureddocs).

### The assessment package's design principle

**The judge never invents a number.** Every score, tier, or category is resolved by deterministic Go code from bounded Y/N judge answers with evidence citations; a `true` answer without `EvidenceIDs` is ignored by every resolver. Related invariants to preserve when extending it:

- Portfolio dimensions (`DimensionDefinition`) are versioned and referenced by ID+version — a definition change never retroactively reinterprets a past assignment. Evolve a dimension by bumping its `Version`, keeping rollup helpers (`MIHRollup`, `PDIMRollup`, `SREWorkRollup`) compatible with legacy option IDs.
- Dimensions are descriptive only — they never enter `RankingPolicy.Rank` (MoSCoW tier first, then RICE descending). Exclusions carry reasons; overrides (`RankOverride`) are auditable records, never quiet input reweighting.
- Category dimensions surface multiple satisfied options as `Ambiguous` rather than silently picking one.
- Cross-framework projections are typed helpers with documented epistemic status (`PDIMToSREWork` is definitional, `PDIMToRGT` is a convention); an explicit native assignment takes precedence over a projection.
- Assessments are never edited in place — a correction is a new cycle (`NextCycle` + `MarkSuperseded`).

## Conventions

- JSON tags are being migrated to **camelCase** (see `MIGRATION_CAMELCASE.md`); new code uses camelCase, but some packages (e.g. `rmi`) still carry snake_case — don't "fix" them outside a deliberate migration change.
- New portfolio dimensions follow the pattern in `assessment/` (one file per framework: `XxxDimension() *DimensionDefinition`, option-ID constants, optional rollup helper, matching `_test.go`). External frameworks are modeled faithfully with source references in the doc comment; house frameworks are marked as prism-roadmap's own.
