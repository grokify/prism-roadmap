# Executive View

The Executive View provides a high-level summary optimized for executive stakeholders and decision-makers.

## Purpose

The Exec View distills the PRD into:

- Quick go/no-go signals
- Investment decision support
- Risk awareness
- Resource planning

## Structure

```go
type ExecView struct {
    Header      ExecHeader    `json:"header"`
    Decision    string        `json:"decision"`     // approve, revise, reject
    Confidence  string        `json:"confidence"`   // High, Medium, Low
    Summary     string        `json:"summary"`
    KeyActions  []ExecAction  `json:"key_actions"`
    Strengths   []string      `json:"strengths"`
    Concerns    []string      `json:"concerns"`
    Risks       []ExecRisk    `json:"risks"`
    ScoreSummary string       `json:"score_summary,omitempty"`
}
```

## Generate Executive View

```go
import "github.com/grokify/prism-roadmap/requirements/prd"

// Score the PRD first
scores := prd.Score(doc)

// Generate exec view with scores
execView := prd.GenerateExecView(doc, scores)

fmt.Printf("Decision: %s\n", execView.Decision)
fmt.Printf("Confidence: %s\n", execView.Confidence)
```

## Render as Markdown

```go
markdown := prd.RenderExecMarkdown(execView)
fmt.Println(markdown)
```

### Example Output

```markdown
# Executive Summary: Customer Portal Redesign

**Status:** Draft | **Decision:** REVISE | **Confidence:** Medium

---

## Recommendation

**REVISE** - The PRD shows promise but needs additional detail in key areas.

## Summary

Customer Portal Redesign aims to reduce bounce rate from 40% to under 20%
through a modern, responsive redesign.

## Strengths

- Clear problem statement with data
- Well-defined success metrics
- Strong persona research

## Concerns

- Limited technical feasibility analysis
- No competitive analysis included
- Risk mitigation plans need detail

## Key Actions Required

1. Add technical architecture section
2. Include competitive analysis
3. Detail risk mitigation strategies

## Risk Summary

| Risk | Impact | Status |
|------|--------|--------|
| Migration complexity | High | Identified |
| Resource constraints | Medium | Mitigating |

## Score Summary

Overall: 68% | Problem: 85% | Users: 70% | Solution: 55%
```

## Decision Logic

The decision is derived from scoring:

| Score Range | Decision | Meaning |
|-------------|----------|---------|
| ≥80% | **APPROVE** | Ready for development |
| 60-79% | **REVISE** | Needs improvements |
| <60% | **REJECT** | Significant rework needed |

!!! note "Blockers Override Score"
    If critical blockers exist, decision becomes REJECT regardless of score.

## Confidence Levels

```go
func determineConfidence(score float64, blockers int) string {
    if blockers > 0 {
        return "Low"
    }
    if score >= 0.8 {
        return "High"
    }
    if score >= 0.6 {
        return "Medium"
    }
    return "Low"
}
```

## ExecView Fields

### ExecHeader

```go
type ExecHeader struct {
    Title       string `json:"title"`
    Status      string `json:"status"`
    Author      string `json:"author"`
    LastUpdated string `json:"last_updated"`
}
```

### ExecAction

```go
type ExecAction struct {
    Action   string `json:"action"`
    Owner    string `json:"owner,omitempty"`
    Priority string `json:"priority,omitempty"`
}
```

### ExecRisk

```go
type ExecRisk struct {
    Description string `json:"description"`
    Impact      string `json:"impact"`
    Status      string `json:"status"`
}
```

## Without Scores

If you don't have scores, the view still works:

```go
// Generate without scores (nil)
execView := prd.GenerateExecView(doc, nil)

// Decision will be based on content analysis only
fmt.Printf("Decision: %s\n", execView.Decision)
```

## Use Cases

### Board Presentation

```go
execView := prd.GenerateExecView(doc, scores)
markdown := prd.RenderExecMarkdown(execView)

// Include in presentation
slides.AddSlide("PRD Review", markdown)
```

### Automated Review

```go
scores := prd.Score(doc)
execView := prd.GenerateExecView(doc, scores)

if execView.Decision == "REJECT" {
    notifyAuthors(doc.Metadata.Authors, execView.Concerns)
}
```

## Next Steps

- [PM View](pm-view.md)
- [6-Pager](six-pager.md)
- [Scoring](../features/scoring.md)
