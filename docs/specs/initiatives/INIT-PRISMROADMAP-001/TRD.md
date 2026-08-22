# TRD — Opportunity Assessment IR: Evidence-Backed Prioritization and Portfolio Reporting

**Initiative:** `INIT-PRISMROADMAP-001`

## Architecture

Four repos, one evaluation system, strict layering:

```
structured-evaluation (plexusone)      generic eval machinery: rubric/, claims/, summary/
        ▲ imported by
prism-roadmap (grokify)                domain semantics: OpportunityAssessment IR, rubric
                                       definitions, deterministic mappings, ranking policy,
                                       portfolio dimensions, OKR links, report contracts
        ▲ references                   ▲ imported by
prism-capability (grokify)             omniroadmap (grokify)
capability IDs in IR                   persistence, evidence store, compile, PM review gate,
                                       rank materialization, report rendering
```

Cross-org import (grokify → plexusone/structured-evaluation) is approved; SE is standard stack for all evaluation work (same tier as Cobra/Huma).

## Core design decisions

### D1 — Assessment is a sibling of OpportunitySpec, not an extension

`canvas.OpportunitySpec` (Patton+Cagan 12-box) is a human-authored discovery **document**: revised occasionally, one per opportunity. `OpportunityAssessment` is a **judgment record**: regenerated per review cycle, versioned against rubric versions, many per opportunity over time. The assessment references the spec by ID; the spec is one of the judge's evidence sources (citable per box, e.g. `#businessValue`). The spec's existing inline `RICE`/`Kano` pointers remain the lightweight manual mode; when an assessment exists it is authoritative.

New package: `assessment/` (or `oppassessment/`) in prism-roadmap, composing `prioritization`, `goals`, `rmi`, `canvas` and SE types.

### D2 — LLM classifies, code computes

Judges answer bounded rubric questions (SE `rubric.Criterion` with checklist thresholds) and must cite evidence per YES; deterministic Go maps classifications to values:

| Dimension | Judge output | Mapping |
|---|---|---|
| MoSCoW | Must-criteria booleans (KTLO, compliance, contractual, critical-risk, EOL); Should/Could criteria | any Must-criterion YES → Must; else highest satisfied class |
| RICE Impact | nested thresholds Massive→None, evaluated top-down, evidence required per YES | Massive 3.0 / High 2.0 / Medium 1.0 / Low 0.5 / Minimal 0.25 / None 0 |
| RICE Confidence | evidence-quality class | High 1.0 / Medium 0.8 / Low 0.5 / INSUFFICIENT_EVIDENCE → review state (not 0) |
| RICE Reach | % of customer accounts, evidence-verified | used directly (0..1) |
| Effort | decomposed Person-Days per work item + estimability gate | sum; gate failure → INSUFFICIENT_DETAIL, no estimate emitted |

Consistency check: `Massive=YES, High=NO` is a contradiction → flag for review. Multi-judge (SE `MultiJudgeResult`): median with disagreement flagging. Judge confidence (did I apply the rubric correctly) is a separate field from RICE Confidence (does evidence support the claims).

### D3 — Ranking policy is data, execution is deterministic

```
tier   = MoSCoW class (Must=1, Should=2, Could=3, Won't=excluded)
within = RICE = (Reach × Impact × Confidence) / EffortPD, descending
ties   = |RICE_a − RICE_b| / max ≤ 0.05 → Kano-informed tie-break (recorded, not automatic)
final  = calculatedRank unless explicit override {finalRank, rationale, approvedBy}
```

Policy struct is versioned in prism-roadmap; omniroadmap executes it. Kano/MIH/themes never enter the formula.

### D4 — Evidence entity (SE claims-based)

```jsonc
{
  "id": "EV-042",
  "claim": "Customer X contract requires FedRAMP by Q3 FY27",
  "uri": "https://docs.google.com/...#heading=h.abc",   // wiki/GDocs/GitLab/GitHub; prefer SHA permalinks for git
  "system": "google-docs",
  "sourceType": "contract",                              // maps to SE ReliabilityTier ladder
  "excerpt": "Section 4.2: Provider shall achieve...",   // what the judge actually evaluated; survives link rot
  "capturedAt": "2026-08-19", "capturedBy": "...",
  "verification": "direct",                              // SE Verdict
  "sensitivity": "restricted"                            // gates excerpt rendering in reports
}
```

Evidence is stored once, referenced by ID from any number of rubric answers across assessments; reverse lookup ("12 assessments cite EV-042") is required. Staleness: per-sourceType validity window; expired evidence flags citing assessments and degrades their Confidence display (deployment inventory ~1 quarter; signed contract effectively unlimited).

### D5 — Portfolio dimensions: two primitives

`CategoryDimension` (0..1 selection per opportunity — pie by % Person-Days) and `TagDimension` (0..N — bar/coverage, sums may exceed 100%). Definitions are versioned entities referenced from assessments (`dimensionId` + `version` + selection + question results); Kano and MIH ship as built-in Category definitions with judge rubrics; custom dimensions (e.g. `strategic-priority-2026: AI/Growth/Excellence`) need no schema change. Kano is not an ordered ladder — classification derives from the answer pattern (expected/absence-dissatisfaction/more-is-better/...), with an insufficient-evidence outcome since code cannot reveal customer satisfaction.

### D6 — Reports are pure functions of the IR

prism-roadmap defines: report data contract (per-section computed facts: ranked table, % PD distributions per dimension, capability overlay, stakeholder impact, deltas vs. previous dataset, override log) + deterministic TOC (conditional sections) + typed narrative slots. omniroadmap: compiles the dataset, holds it for PM review, renders markdown (6-pager + appendices) and a presentation projection from the same dataset. PM edits are IR deltas (override/effort/reclassification/defer) that trigger recompute — narrative-slot text is the only render-side editable. Every material narrative claim carries `derivedFrom` (computation ref) or `evidenceRefs`.

## Persistence (omniroadmap)

- Canonical IR: JSON column (normative record).
- Projections: `moscow_class`, `rice_score`, `rice_reach/impact/confidence/effort_pd`, `opportunity_rank_calculated/final`, `kano_category`, `mih_category` — indexed, rebuilt from IR, never hand-written.
- Custom dimensions: `portfolio_dimensions` / `dimension_options` / `item_dimension_assignments` tables.
- Evidence: own table + `assessment_evidence_refs` join for reverse lookup.
- Ent over Dolt per stack convention; follows existing store/upsert patterns.

## Schema and conventions

- Go structs are source of truth; `invopop/jsonschema` generation; `//go:embed`; `tools.go` guard for the generator dep.
- `schemakit lint --property-case camelCase` on prism-roadmap surfaces; SE `$defs` with divergent case are upstream convention — lint findings isolated there are accepted.
- No `Date.now`-style nondeterminism in compile: dataset carries `generatedAt` supplied by caller; recompute with same IR + policy version is byte-stable.

## Risks

| Risk | Mitigation |
|---|---|
| Rubric gaming moves arguments to evidence-qualification | calibration exemplars per rubric level; periodic rubric review; INSUFFICIENT_EVIDENCE path |
| Stale evidence lends false authority | capturedAt + validity windows + visible Confidence degradation (FR11) |
| Effort PD estimates contested | estimability gate refuses to score under-specified plans; range (low/expected/high) retained |
| SE API evolution breaks IR contract | pin SE version; IR schema version gates; assessment records carry rubric+policy versions |
| Composite-score pressure ("just give one number") | Opportunity Rank is the number; its derivation is the answer — policy doc states only MoSCoW+RICE rank |
