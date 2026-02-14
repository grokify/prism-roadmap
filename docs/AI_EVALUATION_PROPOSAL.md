# AI-Assisted PRD Evaluation - Design Proposal

## Overview

Extend the existing rule-based scoring with an optional AI evaluation layer that assesses semantic quality of PRD content.

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    PRD Evaluation                        │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  ┌─────────────────┐      ┌─────────────────────────┐  │
│  │ Rule-Based      │      │ AI-Based                │  │
│  │ (Deterministic) │      │ (Semantic)              │  │
│  ├─────────────────┤      ├─────────────────────────┤  │
│  │ ✓ Fast          │      │ ✓ Content quality       │  │
│  │ ✓ Consistent    │      │ ✓ Clarity assessment    │  │
│  │ ✓ No cost       │      │ ✓ Gap identification    │  │
│  │ ✓ Offline       │      │ ✓ Suggestions           │  │
│  │                 │      │ ✓ Cross-reference check │  │
│  │ ✗ No semantics  │      │                         │  │
│  │ ✗ No context    │      │ ✗ Slower                │  │
│  │                 │      │ ✗ API cost              │  │
│  │                 │      │ ✗ Non-deterministic     │  │
│  └────────┬────────┘      └────────────┬────────────┘  │
│           │                            │                │
│           └──────────┬─────────────────┘                │
│                      ▼                                  │
│           ┌─────────────────────┐                       │
│           │ Combined Score      │                       │
│           │ + Recommendations   │                       │
│           └─────────────────────┘                       │
└─────────────────────────────────────────────────────────┘
```

## AI Evaluation Categories

### 1. Problem Clarity (Weight: 20%)
**Prompt Focus:**
- Is the problem statement specific and measurable?
- Is there clear cause-and-effect reasoning?
- Would a new team member understand the problem?

**Sample Prompt:**
```
Evaluate this problem statement for clarity and specificity:
"{problem_statement}"

Score 1-10 on:
1. Specificity (vague vs concrete)
2. Measurability (can we tell when it's solved?)
3. Root cause clarity (symptom vs cause)
4. User impact clarity

Provide specific suggestions for improvement.
```

### 2. Persona Validity (Weight: 10%)
**Prompt Focus:**
- Do personas represent real user archetypes?
- Are pain points specific or generic?
- Is there evidence these users exist?

### 3. Requirements Testability (Weight: 15%)
**Prompt Focus:**
- Can each requirement be objectively verified?
- Are acceptance criteria specific enough?
- Are there ambiguous terms ("fast", "easy", "intuitive")?

### 4. Solution-Problem Fit (Weight: 20%)
**Prompt Focus:**
- Does the solution actually address the stated problems?
- Are there gaps where problems aren't addressed?
- Are there solution features without corresponding problems?

### 5. Risk Completeness (Weight: 10%)
**Prompt Focus:**
- Are obvious risks missing?
- Are mitigations realistic?
- Are there unstated assumptions?

### 6. Internal Consistency (Weight: 15%)
**Prompt Focus:**
- Do personas align with user stories?
- Do requirements trace to goals?
- Are there contradictions?

### 7. Competitive Differentiation (Weight: 10%)
**Prompt Focus:**
- Is the differentiation actually meaningful?
- Are competitor weaknesses realistic?
- Is the "why us" compelling?

## Implementation

### Package Split

Types are split between two packages:

| Package | Types | Purpose |
|---------|-------|---------|
| `agentplexus/multi-agent-spec/sdk/go` | Generic report types | Reusable by any LLM evaluation |
| `grokify/structured-plan/requirements/prd` | PRD-specific types | Domain-specific evaluation |

### multi-agent-spec/sdk/go (generic, already exists + extensions)

**Existing types (no changes needed):**

```go
// Already in multi-agent-spec/sdk/go/report.go
type Status string              // GO, WARN, NO-GO, SKIP
type TaskResult struct { ... }  // id, status, detail, metadata
type TeamSection struct { ... } // id, name, tasks, status, contentBlocks, narrative
type TeamReport struct { ... }  // project, version, teams, status

// Already in multi-agent-spec/sdk/go/narrative.go
type NarrativeSection struct {  // problem, analysis, recommendation
    Problem        string `json:"problem,omitempty"`
    Analysis       string `json:"analysis,omitempty"`
    Recommendation string `json:"recommendation,omitempty"`
}
```

**New types to add (in multi-agent-spec/sdk/go/llm.go):**

```go
package multiagentspec

// EvaluationType discriminates between rule-based and LLM evaluation.
type EvaluationType string

const (
    EvaluationTypeRule     EvaluationType = "rule"
    EvaluationTypeLLM      EvaluationType = "llm"
    EvaluationTypeCombined EvaluationType = "combined"
)

// LLMEvaluation contains LLM-based evaluation results.
// This is the nested "llm" object in evaluation results.
type LLMEvaluation struct {
    Score         float64  `json:"score"`
    MaxScore      float64  `json:"maxScore,omitempty"`
    Confidence    float64  `json:"confidence,omitempty"`
    Reasoning     string   `json:"reasoning,omitempty"`
    Strengths     []string `json:"strengths,omitempty"`
    Concerns      []string `json:"concerns,omitempty"`
    Suggestions   []string `json:"suggestions,omitempty"`
    Model         string   `json:"model"`
    Provider      string   `json:"provider,omitempty"`
    TokensUsed    int      `json:"tokensUsed,omitempty"`
    LatencyMs     int      `json:"latencyMs,omitempty"`
    PromptVersion string   `json:"promptVersion,omitempty"`
}

// CombinedWeights specifies weighting for combined evaluation.
type CombinedWeights struct {
    Rule float64 `json:"rule"` // e.g., 0.6
    LLM  float64 `json:"llm"`  // e.g., 0.4
}

// Issue represents a specific problem identified in evaluation.
// Used in NarrativeReport.Issues for detailed fix guidance.
type Issue struct {
    ID             string   `json:"id"`
    Category       string   `json:"category"`
    Severity       string   `json:"severity"` // critical, major, minor, suggestion
    Problem        string   `json:"problem"`
    Location       string   `json:"location,omitempty"`
    Analysis       string   `json:"analysis,omitempty"`
    Recommendation string   `json:"recommendation,omitempty"`
    Example        string   `json:"example,omitempty"`
    Effort         string   `json:"effort,omitempty"` // trivial, low, medium, high
    RelatedIssues  []string `json:"relatedIssues,omitempty"`
}

// Severity constants for Issue.Severity
const (
    SeverityCritical   = "critical"
    SeverityMajor      = "major"
    SeverityMinor      = "minor"
    SeveritySuggestion = "suggestion"
)

// Effort constants for Issue.Effort
const (
    EffortTrivial = "trivial"
    EffortLow     = "low"
    EffortMedium  = "medium"
    EffortHigh    = "high"
)
```

**Extensions to existing types (in multi-agent-spec):**

```go
// Extended TeamSection (add to existing struct)
type TeamSection struct {
    // ... existing fields ...

    // New LLM fields
    EvaluationType EvaluationType  `json:"evaluationType,omitempty"`
    LLM            *LLMEvaluation  `json:"llm,omitempty"`
    CombinedScore  float64         `json:"combinedScore,omitempty"`
}

// Extended TeamReport (add to existing struct)
type TeamReport struct {
    // ... existing fields ...

    // New LLM fields
    EvaluationType EvaluationType   `json:"evaluationType,omitempty"`
    Weights        *CombinedWeights `json:"weights,omitempty"`
    Issues         []Issue          `json:"issues,omitempty"` // Detailed issues for narrative
}
```

### structured-plan/requirements/prd (PRD-specific)

```go
package prd

import (
    mas "github.com/agentplexus/multi-agent-spec/sdk/go"
)

// AIEvaluationConfig configures AI-based PRD evaluation.
type AIEvaluationConfig struct {
    Provider    string  `json:"provider"` // bedrock, openai, anthropic
    Model       string  `json:"model"`
    MaxTokens   int     `json:"maxTokens,omitempty"`
    Temperature float64 `json:"temperature,omitempty"` // 0.0 for consistency
    Enabled     bool    `json:"enabled"`
    Weights     mas.CombinedWeights `json:"weights,omitempty"`
}

// PRD evaluation categories
const (
    CategoryProblemClarity      = "problem-clarity"
    CategoryPersonaValidity     = "persona-validity"
    CategoryRequirementsQuality = "requirements-quality"
    CategorySolutionFit         = "solution-problem-fit"
    CategoryRiskCompleteness    = "risk-completeness"
    CategoryInternalConsistency = "internal-consistency"
    CategoryDifferentiation     = "competitive-differentiation"
)

// PRDEvaluationCategory defines PRD-specific evaluation dimensions.
var PRDEvaluationCategories = []struct {
    ID          string
    Name        string
    Weight      float64
    Description string
}{
    {CategoryProblemClarity, "Problem Clarity", 0.20, "Is the problem specific and measurable?"},
    {CategoryPersonaValidity, "Persona Validity", 0.10, "Do personas represent real user archetypes?"},
    {CategoryRequirementsQuality, "Requirements Quality", 0.15, "Can requirements be objectively verified?"},
    {CategorySolutionFit, "Solution-Problem Fit", 0.20, "Does solution address stated problems?"},
    {CategoryRiskCompleteness, "Risk Completeness", 0.10, "Are obvious risks identified?"},
    {CategoryInternalConsistency, "Internal Consistency", 0.15, "Do sections align and trace?"},
    {CategoryDifferentiation, "Differentiation", 0.10, "Is differentiation meaningful?"},
}

// PRDEvaluator evaluates PRD documents using rule-based and/or LLM analysis.
type PRDEvaluator struct {
    Config AIEvaluationConfig
}

// Evaluate runs evaluation and returns a multi-agent-spec TeamReport.
func (e *PRDEvaluator) Evaluate(doc *Document) (*mas.TeamReport, error) {
    // Implementation:
    // 1. Run rule-based scoring (existing Score() function)
    // 2. If config.Enabled, run LLM evaluation per category
    // 3. Combine scores using config.Weights
    // 4. Build TeamReport with TeamSections per category
    // 5. Generate Issues for narrative report
    return nil, nil // TODO: implement
}
```

### Dependency Direction

```
structured-plan
    └── imports → multi-agent-spec/sdk/go
                      ├── Status, TaskResult, TeamSection, TeamReport
                      ├── NarrativeSection, ContentBlock
                      └── EvaluationType, LLMEvaluation, Issue (new)
```

### CLI Extension

```bash
# Rule-based only (default, fast, free)
splan requirements prd score myproduct.prd.json

# With AI evaluation
splan requirements prd score myproduct.prd.json --ai

# AI with specific provider
splan requirements prd score myproduct.prd.json --ai --provider=bedrock --model=claude-3-sonnet

# AI evaluation only (skip rule-based)
splan requirements prd score myproduct.prd.json --ai-only

# Adjust weighting
splan requirements prd score myproduct.prd.json --ai --ai-weight=0.4
```

### Environment Variables

```bash
# AWS Bedrock (default for AWS-first)
export SPLAN_AI_PROVIDER=bedrock
export AWS_REGION=us-east-1

# OpenAI
export SPLAN_AI_PROVIDER=openai
export OPENAI_API_KEY=sk-...

# Anthropic
export SPLAN_AI_PROVIDER=anthropic
export ANTHROPIC_API_KEY=sk-ant-...
```

## Prompt Engineering Strategy

### Structured Output
Use JSON mode to get consistent, parseable responses:

```json
{
  "score": 7.5,
  "confidence": 0.85,
  "reasoning": "The problem statement clearly identifies...",
  "strengths": [
    "Specific user impact quantified",
    "Root cause identified"
  ],
  "weaknesses": [
    "No baseline metrics provided",
    "Affected user count is estimated, not measured"
  ],
  "suggestions": [
    "Add current state metrics to establish baseline",
    "Include user research citations"
  ]
}
```

### Calibration
- Use few-shot examples of good/bad PRDs
- Include scoring rubric in system prompt
- Request confidence scores to flag uncertain evaluations

### Cost Control
- Evaluate categories independently (parallel, can stop early)
- Cache evaluations by content hash
- Skip AI for categories that score high in rule-based

## Sample Combined Output

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                           PRD EVALUATION (Combined)                          ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ Document: PRD.json                                                           ║
║ Title:    Enterprise Platform                                                ║
║                                                                              ║
║ Rule-Based Score:  7.2 / 10.0  (structural completeness)                     ║
║ AI Score:          6.8 / 10.0  (semantic quality)                            ║
║ Combined Score:    7.0 / 10.0  (60% rule + 40% AI)                           ║
║                                                                              ║
║ Decision: REVISE                                                             ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ CATEGORY BREAKDOWN                                                           ║
╠══════════════════════════════════════════════════════════════════════════════╣
║                           Rule-Based    AI        Combined                   ║
║   Problem Definition         8.0       6.5         7.4                       ║
║   │ AI: "Problem statement is present but lacks quantified impact"           ║
║   │                                                                          ║
║   Solution Fit               7.0       5.0         6.2                       ║
║   │ AI: "Solution addresses 3/5 stated problems; gaps in cost control"       ║
║   │                                                                          ║
║   User Understanding         9.0       8.5         8.8                       ║
║   │ AI: "Personas are specific and pain points are actionable"               ║
║   │                                                                          ║
║   Requirements Quality       7.5       6.0         6.9                       ║
║   │ AI: "3 requirements use ambiguous terms: 'fast', 'seamless'"             ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ AI-IDENTIFIED CONCERNS                                                       ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ • FR-003 "fast response" - define specific latency target                    ║
║ • Risk section missing: "vendor lock-in" given AWS-first strategy            ║
║ • Persona "Feature Engineer" pain points don't map to any requirements       ║
║ • No success metrics for 2 of 5 stated goals                                 ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ AI-GENERATED SUGGESTIONS                                                     ║
╠══════════════════════════════════════════════════════════════════════════════╣
║ 1. Add baseline metrics: "Current time-to-deploy is X days"                  ║
║ 2. Quantify problem impact: "Affecting N users, costing $X/month"            ║
║ 3. Add requirement for cost visibility (maps to Feature Engineer pain)       ║
║ 4. Consider adding "vendor portability" as explicit non-goal                 ║
╚══════════════════════════════════════════════════════════════════════════════╝
```

## Implementation Phases

### Phase 1: Foundation
- [ ] Add `AIEvaluationConfig` to config
- [ ] Implement Bedrock client (AWS-first)
- [ ] Single category evaluation (Problem Definition)
- [ ] `--ai` flag for CLI

### Phase 2: Full Coverage
- [ ] All 7 AI evaluation categories
- [ ] Parallel evaluation
- [ ] Combined scoring logic
- [ ] OpenAI/Anthropic providers

### Phase 3: Advanced
- [ ] Caching by content hash
- [ ] Confidence-based fallback
- [ ] Custom evaluation prompts
- [ ] Batch evaluation for multiple PRDs

## Cost Estimation

| Provider | Model | Est. Tokens/PRD | Est. Cost/PRD |
|----------|-------|-----------------|---------------|
| Bedrock | Claude 3 Sonnet | ~15,000 | ~$0.05 |
| Bedrock | Claude 3 Haiku | ~15,000 | ~$0.01 |
| OpenAI | GPT-4 Turbo | ~15,000 | ~$0.20 |
| OpenAI | GPT-3.5 Turbo | ~15,000 | ~$0.02 |

## Report Schema Integration (agentplexus/multi-agent-spec)

Use the [multi-agent-spec](https://github.com/agentplexus/multi-agent-spec) report schemas for evaluation output. This provides:

- Consistent report format across deterministic and LLM evaluations
- Two report types: pass/fail (GO/NO-GO) and narrative (what to fix)
- Extensible schema with LLM-specific fields

### Evaluation Type Discriminator

```json
{
  "evaluationType": "rule | llm | combined",
  "status": "GO | WARN | NO-GO | SKIP",
  "detail": "Rule-based evaluation detail",

  "llm": {
    "score": 6.5,
    "confidence": 0.85,
    "reasoning": "...",
    "suggestions": ["..."],
    "concerns": ["..."],
    "model": "claude-sonnet-4"
  },

  "weights": {
    "rule": 0.6,
    "llm": 0.4
  }
}
```

### Report Types

| Report | Purpose | Key Fields |
|--------|---------|------------|
| **Pass/Fail Report** | Go/No-Go decision | `status`, `tasks[]`, `content_blocks[]` |
| **Narrative Report** | What to fix | `narrative.issues[]` with problem/analysis/recommendation/example |

### Schema Files

| File | Description |
|------|-------------|
| `llm-evaluation.schema.json` | LLM extension definitions (`LLMEvaluation`, `Issue`, `NarrativeReport`) |
| `agent-result.schema.json` | Individual agent results (extended with `evaluationType`, `llm`) |
| `team-report.schema.json` | Aggregate team report (extended with `narrative`, `weights`) |

### Issue Structure (Narrative Report)

```json
{
  "id": "ISS-001",
  "category": "problem-clarity",
  "severity": "major | minor | critical | suggestion",
  "problem": "Problem impact is not quantified",
  "location": "executiveSummary.problemStatement",
  "analysis": "Why this is a problem...",
  "recommendation": "How to fix it...",
  "example": "Example improved text...",
  "effort": "trivial | low | medium | high"
}
```

## Open Questions

1. ~~Should AI evaluation be opt-in or opt-out?~~ **Decision: Opt-in via `--ai` flag**
2. How to handle AI evaluation failures gracefully?
3. Should we store/cache AI evaluations?
4. How to version prompts for reproducibility? **Partial: `promptVersion` field added**
5. ~~Should AI suggest specific text improvements?~~ **Decision: Yes, via `issue.example` field**
6. Should we support streaming output for long evaluations?
