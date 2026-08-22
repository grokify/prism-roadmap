# PLAN — Opportunity Assessment IR: Evidence-Backed Prioritization and Portfolio Reporting

**Initiative:** `INIT-PRISMROADMAP-001`

## Sequencing rationale

The IR is a one-way door (omniroadmap columns, report contracts, and downstream consumers all inherit it), so prism-roadmap phases land first and harden the contract before omniroadmap builds on it. Within prism-roadmap: evidence and the assessment aggregate first (everything cites evidence), then rubrics + ranking (the scoring core), then linkage (portfolio dimensions, OKR, capability), then report contracts. omniroadmap follows: persistence, then compile/review/rank, then rendering.

## Phase mapping

1. **IR foundations (prism-roadmap)** — Evidence entity on SE `claims`; `OpportunityAssessment` aggregate (definition refs incl. `canvas.OpportunitySpec` by ID, judge metadata, assessment history semantics); MoSCoW + RICE rubric definitions as versioned SE RubricSets with deterministic mapping functions and the Effort estimability gate; ranking policy type + reference implementation with the ±5% tie band and override record.
2. **Portfolio & linkage (prism-roadmap)** — Category/Tags dimension primitives with versioned definitions; built-in Kano + MIH definitions with judge rubrics and decision rules; OKR `contributesTo` links against `goals`; prism-capability ID references; JSON Schema generation + schemakit lint wiring.
3. **Report contracts (prism-roadmap)** — report dataset contract (computed facts per section, deltas, override log); opportunity 6-pager TOC + narrative slots; portfolio review agenda + presentation projection contract.
4. **Persistence (omniroadmap)** — canonical IR storage + Ent projection columns; normalized custom-dimension tables; evidence store with reverse refs and the staleness sweep.
5. **Compile, review, rank (omniroadmap)** — assessment compiler → draft dataset; PM review gate (IR-delta edit commands: override, effort correction, reclassification, defer); rank materialization executing the prism-roadmap policy.
6. **Rendering (omniroadmap)** — opportunity report renderer (markdown, provenance footnotes, sensitivity-gated excerpts); portfolio review renderer (% PD distributions, calculated-vs-final, deltas); presentation projection (Marp) from the same dataset.

## Verification approach

- Unit tests per package (library-first; no integration-test dependence). Rubric mapping tables are table-driven tests; the Massive-YES/High-NO contradiction and INSUFFICIENT_* paths are explicit cases.
- Determinism test: fixture IR → compile twice → byte-identical dataset; policy-version bump changes output, same version never does.
- Round-trip test: assessment → persist (omniroadmap) → load → recompute projections → identical to source IR values.
- Schema gate: `go generate` + `schemakit lint --property-case camelCase` in CI for prism-roadmap surfaces (SE `$defs` exemption documented in the lint invocation).
- Dogfood gate before closing: produce a real assessment set for 3+ actual platform opportunities (this initiative itself is a candidate), compile, review-edit, render both reports, and walk the rank→evidence chain end-to-end.

## Working agreements

- Commits carry `Refs: RMI-<REPOSLUG>-NNN` trailers per RMI (never the INIT ID).
- Scope discovered mid-build gets a new RMI with `--origin implementation`; human-proposed additions `--origin discussion`.
- prism-roadmap README/mkdocs pages for the new packages land with their phases, not deferred (the repo's ROADMAP.md/TODO.md staleness is a known issue — this initiative must not add to it).
- SE changes, if any prove necessary, are logged as `RMI-STRUCTUREDEVALUATION-NNN` under this initiative rather than worked untracked.
