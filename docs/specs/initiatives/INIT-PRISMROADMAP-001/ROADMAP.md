# ROADMAP — Opportunity Assessment IR: Evidence-Backed Prioritization and Portfolio Reporting

**Initiative:** `INIT-PRISMROADMAP-001`
**Repository:** `github.com/grokify/prism-roadmap`

## Phase 1 — IR Foundations
**Theme:** Evidence entity, assessment aggregate, scoring rubrics, ranking policy (prism-roadmap)

- [ ] `RMI-PRISMROADMAP-001` Evidence entity on structured-evaluation claims (EV ids, claim, uri, excerpt, capturedAt, reliability tier, sensitivity, reverse refs)
- [ ] `RMI-PRISMROADMAP-002` OpportunityAssessment aggregate type referencing OpportunitySpec, with judge metadata and per-cycle history semantics
  - Depends on: `RMI-PRISMROADMAP-001`
- [ ] `RMI-PRISMROADMAP-003` MoSCoW and RICE rubric definitions as versioned SE RubricSets with deterministic mappings and Effort estimability gate
  - Depends on: `RMI-PRISMROADMAP-002`
- [ ] `RMI-PRISMROADMAP-004` Ranking policy type and reference implementation (MoSCoW tier, RICE within tier, ±5% tie band, override record)
  - Depends on: `RMI-PRISMROADMAP-003`

## Phase 2 — Portfolio and Linkage
**Theme:** Portfolio dimensions, OKR links, capability refs, schema generation (prism-roadmap)

- [ ] `RMI-PRISMROADMAP-005` Category and Tags portfolio-dimension primitives with versioned referenced definitions
  - Depends on: `RMI-PRISMROADMAP-002`
- [ ] `RMI-PRISMROADMAP-006` Built-in Kano and Market Investment Horizon Category definitions with judge rubrics and decision rules
  - Depends on: `RMI-PRISMROADMAP-005`
- [ ] `RMI-PRISMROADMAP-007` OKR contribution links (contributesTo objective/KR with strength) against the goals package
  - Depends on: `RMI-PRISMROADMAP-002`
- [ ] `RMI-PRISMROADMAP-008` prism-capability references (enables, improves, dependsOn) in the assessment IR
  - Depends on: `RMI-PRISMROADMAP-002`
- [ ] `RMI-PRISMROADMAP-009` JSON Schema generation and schemakit lint wiring for all new IR types
  - Depends on: `RMI-PRISMROADMAP-004`
  - Depends on: `RMI-PRISMROADMAP-006`

## Phase 3 — Report Contracts
**Theme:** Deterministic report dataset contract and templates (prism-roadmap)

- [ ] `RMI-PRISMROADMAP-010` Report dataset contract: computed facts per section, % Person-Day distributions, deltas, override log
  - Depends on: `RMI-PRISMROADMAP-004`
- [ ] `RMI-PRISMROADMAP-011` Opportunity report template: 6-pager TOC, appendix structure, typed narrative slots with provenance refs
  - Depends on: `RMI-PRISMROADMAP-010`
- [ ] `RMI-PRISMROADMAP-012` Portfolio review template: deterministic agenda and presentation projection contract
  - Depends on: `RMI-PRISMROADMAP-010`

## Phase 4 — Persistence
**Theme:** Canonical IR storage, projections, evidence store (omniroadmap)

- [ ] `RMI-OMNIROADMAP-001` Canonical assessment IR storage plus indexed Ent projection columns (moscow_class, rice_score, opportunity_rank, kano, mih)
  - Depends on: `RMI-PRISMROADMAP-009`
- [ ] `RMI-OMNIROADMAP-002` Normalized custom portfolio-dimension tables (no migration per new taxonomy)
  - Depends on: `RMI-OMNIROADMAP-001`
- [ ] `RMI-OMNIROADMAP-003` Evidence store with reverse assessment refs and staleness sweep degrading Confidence display
  - Depends on: `RMI-OMNIROADMAP-001`

## Phase 5 — Compile, Review, Rank
**Theme:** Portfolio compilation, PM review gate, rank materialization (omniroadmap)

- [ ] `RMI-OMNIROADMAP-004` Assessment compiler producing the draft report dataset across all opportunities
  - Depends on: `RMI-OMNIROADMAP-002`
  - Depends on: `RMI-PRISMROADMAP-010`
- [ ] `RMI-OMNIROADMAP-005` PM review gate: IR-delta edit commands (override, effort correction, reclassification, defer) with recompute
  - Depends on: `RMI-OMNIROADMAP-004`
- [ ] `RMI-OMNIROADMAP-006` Rank materialization executing the versioned prism-roadmap ranking policy
  - Depends on: `RMI-OMNIROADMAP-004`

## Phase 6 — Rendering
**Theme:** Report and presentation outputs (omniroadmap)

- [ ] `RMI-OMNIROADMAP-007` Opportunity report renderer: markdown 6-pager plus appendices, provenance footnotes, sensitivity-gated excerpts
  - Depends on: `RMI-OMNIROADMAP-006`
  - Depends on: `RMI-PRISMROADMAP-011`
- [ ] `RMI-OMNIROADMAP-008` Portfolio review renderer: distributions by % Person-Days, calculated-vs-final ranks, deltas since previous review
  - Depends on: `RMI-OMNIROADMAP-006`
  - Depends on: `RMI-PRISMROADMAP-012`
- [ ] `RMI-OMNIROADMAP-009` Presentation projection (Marp) rendered from the same report dataset
  - Depends on: `RMI-OMNIROADMAP-008`
