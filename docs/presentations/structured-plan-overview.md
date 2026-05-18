---
marp: true
theme: default
paginate: true
title: PRISM Execution - Unified Planning Artifacts
---

# PRISM Execution

## Unified Planning Artifacts for the Entire Organization

A cascading system for V2MOM, OKR, Roadmap, and PRD

---

# The Problem

Planning artifacts are often **siloed**:

- V2MOM in one tool
- OKRs in another spreadsheet
- Roadmaps in yet another tool
- PRDs in documents with no structure

**Result:** Drift, duplication, misalignment

---

# The Solution

**One unified system** where:

- Same types at every level
- Artifacts link and cascade
- Single source of truth
- Evolve together

---

# The Cascading Model

```
Organization
    ├── V2MOM (company vision)
    ├── OKRs (company objectives)
    └── Roadmap (strategic initiatives)
            │
            ▼
        Department / Team
            ├── V2MOM (team methods)
            ├── OKRs (team objectives)
            └── Roadmap (team deliverables)
                    │
                    ▼
                PRD (project level)
                    ├── OKRs (project KRs)
                    └── Roadmap (features)
```

---

# V2MOM: Cascading Goals

Created by Marc Benioff at Salesforce

| Level | Vision | Methods | Measures |
|-------|--------|---------|----------|
| **Company** | "Market leader in X" | Strategic initiatives | Revenue, market share |
| **Department** | Supports company vision | Dept programs | Dept KPIs |
| **Team** | Supports dept vision | Team projects | Team metrics |
| **Individual** | Supports team vision | Personal goals | Individual KPIs |

Each level's V2MOM **aligns with** the level above

---

# OKR: Aligned Objectives

| Level | Objective | Key Results |
|-------|-----------|-------------|
| **Company** | Become market leader | 20% market share, $100M ARR |
| **Team** | Launch enterprise product | 50 enterprise customers, 99.9% uptime |
| **PRD** | Build auth system | SSO for 10 providers, <100ms latency |

Team OKRs roll up to Company OKRs
PRD KRs contribute to Team KRs

---

# Roadmap: Nested Deliverables

**Portfolio Roadmap** (org level)

| Q1 2026 | Q2 2026 | Q3 2026 |
|---------|---------|---------|
| Auth PRD | Dashboard PRD | Analytics PRD |
| Platform PRD | API v2 PRD | Mobile PRD |

**Product Roadmap** (PRD level)

| Phase 1 | Phase 2 | Phase 3 |
|---------|---------|---------|
| SSO Integration | MFA Support | Audit Logging |
| User Management | Role-Based Access | Compliance Reports |

Same structure, different granularity

---

# PRD: The Project Level

A PRD contains:

- **OKRs** - Project objectives aligned to team/company OKRs
- **Roadmap** - Phased delivery of features
- **Personas** - Who we're building for
- **Requirements** - What we're building
- **Risks** - What could go wrong

PRD is one deliverable in the Portfolio Roadmap

---

# Unified Type System

```
prism-execution/
├── goals/
│   ├── okr/        # Objective, KeyResult, PhaseTarget
│   └── v2mom/      # V2MOM, Method, Measure, Obstacle
├── roadmap/        # Roadmap, Phase, Deliverable
└── requirements/
    ├── prd/        # Uses goals/ and roadmap/ types
    ├── mrd/
    └── trd/
```

**One set of types used at all levels**

---

# Linking Across Levels

```go
// Portfolio Roadmap deliverable references a PRD
Deliverable{
    Title: "Authentication Initiative",
    Type:  "prd",
    PRDRef: "prd-auth-2026",  // Link to PRD
}

// PRD OKR aligns to company OKR
Objective{
    Description: "Secure authentication for enterprise",
    AlignedWith: ["company-okr-security"],  // Link up
}
```

---

# Benefits

| Before | After |
|--------|-------|
| Siloed tools | Unified system |
| Manual alignment | Automatic linking |
| Drift between docs | Single source of truth |
| Different formats | Consistent structure |
| Hard to track | Clear traceability |

---

# Output Formats

From the same structured data:

- **Markdown** - Documentation
- **Marp Slides** - Presentations
- **JSON** - API/tooling
- **HTML** - Web views
- **Tables** - Swimlane roadmaps

<!-- TODO: Add screenshots when implemented -->

---

# CLI: `splan`

```bash
# Generate PRD markdown
splan requirements prd markdown input.json -o prd.md

# Generate roadmap slides
splan roadmap marp input.json -o roadmap.md

# Filter by tags
splan requirements prd filter input.json --include mvp

# Validate alignment
splan validate --check-alignment portfolio.json prd-*.json
```

<!-- TODO: Update with actual CLI when implemented -->

---

# Roadmap for PRISM Execution

| Phase | Deliverables |
|-------|--------------|
| **Phase 1** | Restructure packages (requirements/, roadmap/, goals/) |
| **Phase 2** | Rename repo to prism-execution |
| **Phase 3** | Consolidate OKR/V2MOM goals frameworks |
| **Phase 4** | Cross-level alignment validation |

---

# Summary

- **Cascading system**: Org → Dept → Team → PRD
- **Same types at all levels**: Just different content
- **Linked artifacts**: PRD refs, OKR alignment
- **Single source of truth**: No duplication
- **Multiple outputs**: Markdown, Marp, JSON, HTML

---

# Questions?

GitHub: `github.com/grokify/prism-execution`

<!-- TODO: Update URL when repo is renamed -->
