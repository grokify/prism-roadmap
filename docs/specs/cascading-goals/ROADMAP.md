# Cascading Goals Framework Roadmap

**Implementation Roadmap**

**Version:** 1.0.0
**Date:** 2026-05-22
**Status:** Draft

## Overview

This roadmap outlines the implementation phases for the unified cascading goals framework supporting OKR and V2MOM methodologies.

## Timeline Summary

```
Phase 1: OKR Cascading          ████████░░░░░░░░░░░░░░░░░░░░░░
Phase 2: Unified Interface      ░░░░░░░░████████░░░░░░░░░░░░░░
Phase 3: Maturity Planning      ░░░░░░░░░░░░░░░░████████░░░░░░
Documentation & Polish          ░░░░░░░░░░░░░░░░░░░░░░░░██████

                                |-------|-------|-------|------|
                                Week 1-2 Week 3-4 Week 5-6 Week 7-8
```

## Phase 1: Extend OKRs with Cascading

**Goal:** Add parent-child relationships and cascading support to OKR structure

### Milestone 1.1: OKR Type Extensions

**Files to modify:** `goals/okr/okr.go`

| Task | Description | Priority |
|------|-------------|----------|
| Add `ParentObjectiveID` | Link child objectives to parent | P0 |
| Add `CompanyOKRIDs` | Align with company-level OKRs | P0 |
| Add `ChildObjectiveIDs` | Track child objectives | P1 |
| Add `SupportsObjectiveIDs` to KeyResult | Track which parent objectives KR supports | P1 |
| Add `SLIIDs` to KeyResult | Link to PRISM SLIs | P1 |
| Add `AlignmentScore` and `AlignmentStatus` | Track alignment metadata | P2 |

**Deliverables:**

- [ ] Extended `Objective` type with cascading fields
- [ ] Extended `KeyResult` type with alignment fields
- [ ] JSON schema updates
- [ ] Unit tests for serialization

### Milestone 1.2: OKR Cascade Logic

**Files to create:** `goals/okr/cascade.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `Cascade()` method | Generate child OKR from parent | P0 |
| Add name templating | Customize child goal names | P1 |
| Add result inheritance | Copy parent KRs as templates | P1 |
| Add filtering | Include/exclude specific objectives | P1 |

**Deliverables:**

- [ ] `CascadeOKR()` function
- [ ] `CascadeOptions` configuration type
- [ ] Name template support
- [ ] Unit tests with various scenarios

### Milestone 1.3: OKR Validation Updates

**Files to modify:** `goals/okr/validate.go`

| Task | Description | Priority |
|------|-------------|----------|
| Validate `ParentObjectiveID` references | Ensure parent exists | P0 |
| Validate KR `SupportsObjectiveIDs` | Ensure referenced objectives exist | P1 |
| Add cascade-specific validation | Check child properly links to parent | P1 |

**Deliverables:**

- [ ] Extended validation for cascading fields
- [ ] Warning for orphaned objectives
- [ ] Error for invalid parent references

---

## Phase 2: Unified Goals Interface

**Goal:** Create abstract interface working with both OKR and V2MOM

### Milestone 2.1: Framework Interface

**Files to create:** `goals/framework.go`, `goals/goal_item.go`

| Task | Description | Priority |
|------|-------------|----------|
| Define `Framework` interface | Abstract over OKR/V2MOM | P0 |
| Define `GoalItem` type | Framework-agnostic goal | P0 |
| Define `ResultItem` type | Framework-agnostic result | P0 |
| Define `CascadeOptions` | Unified cascade configuration | P0 |
| Define `ValidationOptions` | Unified validation configuration | P1 |

**Deliverables:**

- [ ] `Framework` interface with all methods
- [ ] `GoalItem` and `ResultItem` types
- [ ] Configuration types
- [ ] Interface documentation

### Milestone 2.2: OKR Framework Implementation

**Files to create:** `goals/okr/framework.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `Type()` | Return "okr" | P0 |
| Implement `Metadata()` | Extract document metadata | P0 |
| Implement `GoalItems()` | Convert objectives to GoalItems | P0 |
| Implement `ResultItems()` | Convert KRs to ResultItems | P0 |
| Implement `Validate()` | Wrap existing validation | P0 |
| Implement `Cascade()` | Wrap cascade logic | P0 |
| Implement `AlignmentScore()` | Calculate alignment | P1 |
| Implement `ToPRISMGoals()` | Convert to PRISM goals | P1 |

**Deliverables:**

- [ ] Complete `Framework` implementation for OKR
- [ ] Unit tests for each method
- [ ] Integration tests with real OKR documents

### Milestone 2.3: V2MOM Framework Implementation

**Files to create:** `goals/v2mom/framework.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `Type()` | Return "v2mom" | P0 |
| Implement `Metadata()` | Extract document metadata | P0 |
| Implement `GoalItems()` | Convert methods to GoalItems | P0 |
| Implement `ResultItems()` | Convert measures to ResultItems | P0 |
| Implement `Validate()` | Wrap existing validation | P0 |
| Implement `Cascade()` | Use existing ParentID logic | P0 |
| Implement `AlignmentScore()` | Calculate alignment | P1 |
| Implement `ToPRISMGoals()` | Convert to PRISM goals | P1 |

**Deliverables:**

- [ ] Complete `Framework` implementation for V2MOM
- [ ] Unit tests for each method
- [ ] Parity with OKR implementation

### Milestone 2.4: Alignment Service

**Files to create:** `goals/align.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `AlignmentService` | Core alignment logic | P0 |
| Implement `Align()` method | Compare child to parent | P0 |
| Calculate `AlignmentResult` | Score and recommendations | P0 |
| Identify `MissingSupport` | Parent goals without children | P1 |
| Identify `UnalignedGoals` | Child goals without parent | P1 |
| Generate `Recommendations` | Actionable suggestions | P1 |

**Deliverables:**

- [ ] `AlignmentService` with full analysis
- [ ] `AlignmentResult` type
- [ ] Recommendation generation
- [ ] Unit tests for alignment scenarios

### Milestone 2.5: CLI Commands (Phase 2)

**Files to create:** `cli/goals/*.go`

| Task | Description | Priority |
|------|-------------|----------|
| Create `goals` command group | Parent command | P0 |
| Implement `cascade` command | Generate child frameworks | P0 |
| Implement `align` command | Validate alignment | P0 |
| Implement `import` command | Import to PRISM goals | P1 |
| Add `--format` flag | JSON, markdown output | P1 |

**Deliverables:**

- [ ] `prism goals cascade` command
- [ ] `prism goals align` command
- [ ] `prism goals import` command
- [ ] Help text and examples

---

## Phase 3: Maturity-Driven Planning

**Goal:** Enable goal suggestions and planning based on SLI maturity gaps

### Milestone 3.1: Suggestion Service

**Files to create:** `goals/suggest.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `SuggestionService` | Core suggestion logic | P0 |
| Analyze SLI gaps | Compare current vs threshold | P0 |
| Generate `Suggestion` objects | Title, description, KRs | P0 |
| Calculate priority | Based on gap and impact | P1 |
| Estimate effort | Based on gap size | P1 |
| Filter by domain/layer | Focused suggestions | P1 |

**Deliverables:**

- [ ] `SuggestionService` implementation
- [ ] `Suggestion` and `SLIGap` types
- [ ] Priority and effort estimation
- [ ] Domain/layer filtering

### Milestone 3.2: Plan Service

**Files to create:** `goals/plan.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `PlanService` | Generate phased plans | P0 |
| Calculate maturity progression | M2 → M3 → M4 phases | P0 |
| Generate phase-specific goals | Goals per phase | P0 |
| Output as OKR or V2MOM | Framework-specific output | P1 |
| Include initiative suggestions | Enablers for each goal | P2 |

**Deliverables:**

- [ ] `PlanService` implementation
- [ ] Multi-phase goal generation
- [ ] OKR and V2MOM output formats

### Milestone 3.3: Auto-Scoring from SLI State

**Files to create:** `goals/score.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `ScoreService` | Update scores from state | P1 |
| Calculate KR scores | Based on SLI current values | P1 |
| Calculate Objective scores | Average of KR scores | P1 |
| Update goal progress | Propagate to GoalItem | P1 |

**Deliverables:**

- [ ] `ScoreService` implementation
- [ ] Automatic score calculation
- [ ] Integration with state documents

### Milestone 3.4: CLI Commands (Phase 3)

**Files to modify:** `cli/goals/*.go`

| Task | Description | Priority |
|------|-------------|----------|
| Implement `suggest` command | Generate suggestions | P0 |
| Implement `plan` command | Generate maturity plan | P1 |
| Add `--target-level` flag | Specify target maturity | P0 |
| Add `--domain` filter | Filter by domain | P1 |
| Add `--format okr|v2mom` | Output format | P1 |

**Deliverables:**

- [ ] `prism goals suggest` command
- [ ] `prism goals plan` command
- [ ] Filter and format options

---

## Documentation & Polish

### Documentation Tasks

| Task | Description | Priority |
|------|-------------|----------|
| Update README.md | Overview and quick start | P0 |
| Create user guide | Full documentation | P0 |
| Add CLI examples | Common use cases | P0 |
| Create API reference | GoDoc documentation | P1 |
| Add architecture diagrams | Visual documentation | P2 |

### Polish Tasks

| Task | Description | Priority |
|------|-------------|----------|
| Error message improvements | User-friendly errors | P1 |
| Progress indicators | CLI feedback | P2 |
| JSON schema generation | For validation | P1 |
| Example documents | OKR and V2MOM samples | P1 |

---

## Success Criteria

### Phase 1 Complete When

- [ ] OKR documents support `ParentObjectiveID` field
- [ ] `prism goals cascade` generates valid child OKRs
- [ ] Unit test coverage >80%

### Phase 2 Complete When

- [ ] Both OKR and V2MOM implement `Framework` interface
- [ ] `prism goals align` validates alignment for both frameworks
- [ ] `prism goals import` converts either format to PRISM goals

### Phase 3 Complete When

- [ ] `prism goals suggest` generates actionable suggestions
- [ ] Suggestions include priority, effort, and impact estimates
- [ ] `prism goals plan` generates multi-phase improvement plans

### Project Complete When

- [ ] All CLI commands documented with examples
- [ ] Integration tests pass for common workflows
- [ ] README updated with cascading goals documentation

---

## Risk Register

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| Interface too complex | High | Medium | Start minimal, add methods incrementally |
| OKR/V2MOM differences too large | Medium | Low | Accept some framework-specific features |
| Suggestion quality poor | Medium | Medium | Validate with real data, iterate |
| Performance issues | Low | Low | Profile early, optimize hot paths |

---

## Dependencies

| Dependency | Required By | Status |
|------------|-------------|--------|
| prism-maturity | Phase 3 | Available |
| prism-roadmap goals/okr | Phase 1 | Available |
| prism-roadmap goals/v2mom | Phase 2 | Available |

---

## Related Documents

- [PRD.md](PRD.md) - Product Requirements Document
- [TRD.md](TRD.md) - Technical Requirements Document
