# PRD — Opportunity Assessment IR: Evidence-Backed Prioritization and Portfolio Reporting

**Initiative:** `INIT-PRISMROADMAP-001`

## Problem

Platform roadmap prioritization (Application Platform + Cloud Platform for enterprise SaaS) is high-stakes and contested. Today, prioritization arguments are "my Must vs. your Must": scores are entered by hand into spreadsheets, the reasoning behind them is not captured, and stakeholders — especially Engineering, who is simultaneously implementer, internal customer, and stakeholder — cannot audit why one investment outranks another. Leadership cannot see whether the portfolio composition (table stakes vs. differentiation, KTLO vs. market expansion) matches stated strategy.

The fight is never "what's the formula" — it's "why should I believe your inputs." We need a system where every score decomposes into rubric answers with cited evidence, so disagreements become "is this evidence sufficient" instead of competing assertions.

## Goals

1. **A canonical JSON IR** (`OpportunityAssessment`) in prism-roadmap that captures, per opportunity: definition references, OKR contribution links, ranking-framework assessments (MoSCoW, RICE), portfolio classifications (Kano, Market Investment Horizon, custom categories/tags), capability references, stakeholder/value mappings, and evidence — with full provenance (rubric version, judge model, timestamps, calculated vs. final rank with override rationale).
2. **Deterministic ranking from structured judgments.** LLM judges classify semantics against versioned rubrics (Is this Massive? Y/N + evidence); deterministic code maps classifications to values and computes Opportunity Rank. The LLM never emits a number.
3. **Evidence as a first-class entity** — referenced by ID across assessments, carrying claim, source link, excerpt captured at assessment time, capture date, reliability tier, and sensitivity flag.
4. **Execution in omniroadmap**: persist assessments, compile portfolio datasets across all opportunities, gate on Platform PM review (edits flow back as IR deltas — never edits to rendered output), materialize ranks, and render reports.
5. **Two report products** from one dataset: a per-opportunity narrative report (6-pager body + evidence appendices) and a portfolio roadmap review (deterministic agenda: ranking → Kano → MIH → strategy → capability stack → stakeholder impact → deltas → decisions), each with LLM narrative slots that interpret computed facts only.

## Users and Stakeholders

- **Platform PM** (primary): authors/curates assessments, reviews compiled datasets, applies governed overrides, presents the roadmap.
- **Engineering**: consumes the capability-stack and stakeholder-value views; must be able to see how the roadmap serves engineering needs and audit any ranking.
- **Feature product teams**: consume Kano/customer-value views; contribute evidence.
- **Leadership**: consumes portfolio composition (% Person-Days by Kano, MIH, strategic theme) and calculated-vs-final override signal; adjusts levers by challenging evidence or issuing explicit overrides.

## Functional Requirements

### prism-roadmap (specification layer)

- FR1: `OpportunityAssessment` type referencing (not extending) the existing `canvas.OpportunitySpec`; the spec is a definition document and an evidence source, the assessment is a regenerable judgment record with history.
- FR2: Rubric definitions as versioned `structured-evaluation` RubricSets: MoSCoW tier criteria (Must = any of KTLO / compliance / contractual / critical-risk / EOL booleans), RICE Impact nested thresholds (Massive→None, highest YES wins → 3.0/2.0/1.0/0.5/0.25), RICE Confidence evidence classification (High/Medium/Low → 1.0/0.8/0.5, plus INSUFFICIENT_EVIDENCE as a review state, never a zero), Effort in Person-Days rolled up from decomposed work with an estimability gate (insufficient detail → INSUFFICIENT_DETAIL, not a guess).
- FR3: Deterministic ranking policy: MoSCoW tier → RICE descending within tier → Kano-informed tie-break within a ±5% RICE equivalence band → explicit override (`calculatedRank` vs `finalRank` + rationale + approver).
- FR4: Portfolio-dimension primitives: `Category` (0..1 per opportunity, pie-chartable) and `Tags` (0..N, bar/coverage charts); Kano and Market Investment Horizon (KTLO / SAM+SOM / TAM Expansion) ship as built-in versioned Category definitions; organizations add custom dimensions (e.g. AI/Growth/Excellence) without schema changes. Definitions are referenced by version, never copied into assessments.
- FR5: Evidence entity (EV-NNN ids) built on `structured-evaluation/claims`: claim text, source URI + system type, excerpt, capturedAt/capturedBy, reliability tier, verification status, sensitivity/renderable flag. Queryable in reverse ("which assessments cite this source").
- FR6: OKR contribution links (`contributesTo` objective/KR with strength) kept out of the scoring formula; enables investment-by-objective portfolio rollups and post-delivery outcome measurement.
- FR7: prism-capability references (`enables`/`improves`/`dependsOn` capability IDs) for the Engineering capability-stack view.
- FR8: Report contracts: the deterministic TOC/agenda for both report types, the computed-facts data contract each section requires, and typed narrative slots. Reports are a pure function of the IR.
- FR9: JSON Schemas generated from Go types (invopop/jsonschema), linted with `schemakit lint --property-case camelCase`; SE `$defs` case deviations are accepted as upstream convention.

### omniroadmap (execution layer)

- FR10: Persist canonical assessment IR as the normative record plus indexed projection columns (`moscow_class`, `rice_score`, `opportunity_rank`, `kano_category`, `mih_category`); custom dimensions in normalized tables (no migration per new taxonomy).
- FR11: Evidence store with staleness sweep: evidence older than its source-type validity window flags the assessments citing it and degrades their Confidence classification visibly.
- FR12: Compile step producing a draft report dataset (rankings, distributions by % Person-Days, deltas since previous review, flagged items) across all opportunities.
- FR13: PM review gate: structured edits that flow back into the IR (override with rationale, corrected effort, reclassification with new evidence, defer/suppress) — the rendered report is never directly edited; narrative-slot text is the sole PM-editable rendering element.
- FR14: Renderers for the opportunity report (markdown 6-pager + appendices, every material claim traceable to computation or cited evidence) and the portfolio review (document + presentation projection of the same dataset).

## Non-Goals

- A single weighted composite score combining MoSCoW/RICE/Kano/MIH — explicitly rejected; only MoSCoW + RICE determine rank.
- Kano/MIH/strategic themes automatically altering rank (they describe the portfolio and inform tie-breaks/overrides).
- Maturity modeling (prism-maturity integration is future work; capability references suffice now).
- Automated evidence fetchers per source system (GitLab/Google Drive APIs) — manual excerpt + link is a valid v1; fetchers are a later enhancement.
- UI. This initiative is library + CLI surface; visualization lands later.

## Success Criteria

- An agent or PM can produce a complete, schema-valid OpportunityAssessment for a real platform initiative with every rubric answer citing evidence.
- Rank recomputation is reproducible: same IR + same policy version → same ranking, byte-stable report dataset.
- A contested ranking can be defended end-to-end: rank → framework classifications → rubric answers → evidence records → source links.
- Portfolio review renders with zero hand-authored numbers: every figure computed from the IR.
