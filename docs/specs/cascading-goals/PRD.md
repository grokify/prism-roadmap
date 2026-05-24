# Cascading Goals Framework PRD

**Product Requirements Document**

**Version:** 1.0.0
**Date:** 2026-05-22
**Author:** PRISM Team
**Status:** Draft

## Executive Summary

This PRD defines requirements for a unified cascading goals framework that supports both OKR and V2MOM methodologies. The framework enables organizations to cascade strategic goals from company level through teams to individuals, with automatic alignment validation and maturity-driven goal suggestion.

## Problem Statement

### Current Challenges

1. **Framework Fragmentation** - Organizations using OKRs lack built-in cascading; V2MOM has cascading but lower adoption
2. **Manual Alignment** - Child goals must be manually aligned with parent goals
3. **Disconnected Maturity** - PRISM maturity goals don't automatically cascade through organizational hierarchy
4. **No Goal Suggestions** - Teams manually identify SLI gaps without automated recommendations

### User Pain Points

| Persona | Pain Point |
|---------|------------|
| Engineering Director | Cannot cascade reliability goals to team OKRs automatically |
| Team Lead | Manual work to ensure team OKRs align with org objectives |
| SRE Manager | No automated suggestions for which SLIs need improvement |
| Product Manager | Cannot visualize goal alignment across organization |

## Goals and Success Metrics

### Primary Goals

| Goal | Success Metric |
|------|----------------|
| Enable goal cascading | 100% of goal frameworks support parent-child relationships |
| Automate alignment validation | Alignment check completes in <1s for 1000 goals |
| Suggest maturity improvements | 80% of suggestions rated useful by users |
| Support both OKR and V2MOM | Unified CLI commands work with either framework |

### Non-Goals

- Real-time collaborative editing (out of scope for v1)
- AI-generated goal text (may be added later)
- Integration with external OKR tools (Lattice, 15Five, etc.)

## User Stories

### Epic 1: Cascading Goals

```
As an Engineering Director
I want to cascade my OKRs to team-level OKRs
So that team goals automatically align with organizational objectives
```

**Acceptance Criteria:**

- [ ] Can create child OKR linked to parent via `ParentObjectiveID`
- [ ] Child Key Results map to parent Objectives
- [ ] Validation warns when child goals don't support parent
- [ ] `prism goals cascade` generates child goal templates

### Epic 2: Unified Goals Interface

```
As a Platform Team Lead
I want to use the same commands for OKR and V2MOM
So that I don't need to learn two different workflows
```

**Acceptance Criteria:**

- [ ] `prism goals list` works for both OKR and V2MOM documents
- [ ] `prism goals align` validates alignment in either framework
- [ ] `prism goals import` imports either format to PRISM
- [ ] Framework auto-detected from document structure

### Epic 3: Maturity-Driven Planning

```
As an SRE Manager
I want PRISM to suggest goals based on SLI gaps
So that I can prioritize maturity improvements effectively
```

**Acceptance Criteria:**

- [ ] `prism goals suggest` analyzes current SLI state vs targets
- [ ] Suggestions include estimated effort and impact
- [ ] Suggestions map to specific maturity level improvements
- [ ] Can filter suggestions by domain, layer, or priority

## Functional Requirements

### FR1: OKR Cascading Support

| ID | Requirement | Priority |
|----|-------------|----------|
| FR1.1 | Add `ParentObjectiveID` field to Objective type | P0 |
| FR1.2 | Add `CompanyOKRIDs` field for company-level alignment | P0 |
| FR1.3 | Validate child KRs support parent Objectives | P1 |
| FR1.4 | Generate child OKR templates from parent | P1 |

### FR2: Unified Goals Interface

| ID | Requirement | Priority |
|----|-------------|----------|
| FR2.1 | Create `goals.Framework` interface abstracting OKR/V2MOM | P0 |
| FR2.2 | Implement `GoalItems()` returning goals in unified format | P0 |
| FR2.3 | Implement `ResultItems()` returning results in unified format | P0 |
| FR2.4 | Auto-detect framework from document structure | P1 |

### FR3: Maturity-Driven Planning

| ID | Requirement | Priority |
|----|-------------|----------|
| FR3.1 | Analyze SLI current values vs maturity level thresholds | P0 |
| FR3.2 | Identify SLIs below target maturity level | P0 |
| FR3.3 | Generate goal suggestions with improvement targets | P1 |
| FR3.4 | Estimate effort based on gap size and complexity | P2 |

### FR4: CLI Commands

| ID | Requirement | Priority |
|----|-------------|----------|
| FR4.1 | `prism goals cascade <parent> --output <child>` | P0 |
| FR4.2 | `prism goals align <child> <parent>` | P0 |
| FR4.3 | `prism goals import <file> --to-prism <output>` | P1 |
| FR4.4 | `prism goals suggest --state <state.json>` | P1 |
| FR4.5 | `prism goals plan --target-level <M3>` | P2 |

## Non-Functional Requirements

### Performance

| Requirement | Target |
|-------------|--------|
| Cascade generation | <100ms for 50 goals |
| Alignment validation | <1s for 1000 goals |
| Suggestion generation | <500ms for 100 SLIs |

### Compatibility

- Go 1.21+
- JSON Schema validation
- Existing OKR/V2MOM document formats preserved

## Data Model Changes

### OKR Extensions

```go
type Objective struct {
    // ... existing fields ...

    // Cascading support
    ParentObjectiveID string   `json:"parentObjectiveId,omitempty"`
    CompanyOKRIDs     []string `json:"companyOkrIds,omitempty"`
    ChildObjectiveIDs []string `json:"childObjectiveIds,omitempty"`
}

type KeyResult struct {
    // ... existing fields ...

    // Alignment tracking
    SupportsObjectiveIDs []string `json:"supportsObjectiveIds,omitempty"`
}
```

### Unified Interface

```go
type Framework interface {
    Type() string                    // "okr" or "v2mom"
    GoalItems() []GoalItem           // Objectives or Methods
    ResultItems() []ResultItem       // Key Results or Measures
    Validate(opts ValidationOptions) []ValidationError
    Cascade(parent Framework) (Framework, error)
    AlignmentScore(parent Framework) float64
}

type GoalItem struct {
    ID          string
    Name        string
    Description string
    Owner       string
    Status      string
    Progress    float64
    ParentID    string
    Children    []string
    Results     []ResultItem
}

type ResultItem struct {
    ID       string
    Name     string
    Baseline string
    Target   string
    Current  string
    Score    float64
    GoalID   string
}
```

## User Interface

### CLI Examples

```bash
# Cascade company OKRs to team OKRs
prism goals cascade company-okrs.json \
  --team "Platform Team" \
  --output team-okrs.json

# Validate team alignment with company
prism goals align team-okrs.json company-okrs.json

# Suggest goals from maturity gaps
prism goals suggest \
  --model model.json \
  --state state.json \
  --target-level M4

# Import OKRs to PRISM goals
prism goals import team-okrs.json \
  --format okr \
  --to-prism goals-output.json

# Generate maturity improvement plan
prism goals plan \
  --model model.json \
  --state state.json \
  --from-level M2 \
  --to-level M4 \
  --phases 4
```

### Output Formats

- JSON (default)
- Markdown tables
- Marp slides (for presentations)

## Dependencies

| Dependency | Purpose |
|------------|---------|
| prism-maturity | SLI/maturity model access |
| prism-roadmap | OKR/V2MOM types |
| structureddocs | Rendering output |

## Risks and Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
| Framework differences too large | High | Design minimal common interface |
| Cascading creates invalid goals | Medium | Validation at cascade time |
| Suggestions not actionable | Medium | Include specific SLI targets |

## Timeline

See [ROADMAP.md](ROADMAP.md) for detailed timeline.

## Appendix

### A. Framework Comparison

| Aspect | OKR | V2MOM |
|--------|-----|-------|
| Cascading | To be added | Existing (ParentID) |
| Values/Context | None | Native Values |
| Structure | Fixed hierarchy | Configurable |
| Adoption | Higher | Salesforce ecosystem |

### B. Related Documents

- [TRD.md](TRD.md) - Technical Requirements Document
- [ROADMAP.md](ROADMAP.md) - Implementation Roadmap
