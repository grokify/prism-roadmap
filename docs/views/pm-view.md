# PM View

The PM View provides a comprehensive summary of the PRD tailored for product managers.

## Purpose

The PM View extracts and organizes PRD information for:

- Sprint planning
- Stakeholder updates
- Feature prioritization
- Team alignment

## Structure

```go
type PMView struct {
    Header       PMHeader         `json:"header"`
    Problem      string           `json:"problem"`
    Solution     string           `json:"solution"`
    Personas     []PersonaSummary `json:"personas"`
    Goals        []string         `json:"goals"`
    Requirements RequirementsList `json:"requirements"`
    Metrics      MetricsSummary   `json:"metrics"`
    Risks        []RiskSummary    `json:"risks"`
    Timeline     string           `json:"timeline,omitempty"`
}
```

## Generate PM View

```go
import "github.com/grokify/prism-roadmap/requirements/prd"

// Generate view from PRD
pmView := prd.GeneratePMView(doc)

// Access view data
fmt.Printf("Title: %s\n", pmView.Header.Title)
fmt.Printf("Problem: %s\n", pmView.Problem)

for _, persona := range pmView.Personas {
    fmt.Printf("Persona: %s (%s)\n", persona.Name, persona.Role)
}
```

## Render as Markdown

```go
markdown := prd.RenderPMMarkdown(pmView)
fmt.Println(markdown)
```

### Example Output

```markdown
# Customer Portal Redesign

**Version:** 1.0.0 | **Status:** Draft | **Author:** Jane Doe

---

## Problem Statement

Current portal has 40% bounce rate due to poor UX and slow performance.

## Proposed Solution

Modern, responsive redesign with improved navigation and faster load times.

## Target Personas

| Persona | Role | Primary |
|---------|------|---------|
| Power User | Account Manager | Yes |
| Casual User | Customer | No |

## Goals

- Reduce bounce rate to under 20%
- Increase task completion by 50%

## Requirements Summary

### Must Have (3)
- User authentication
- Dashboard redesign
- Search functionality

### Should Have (2)
- Export to PDF
- Dark mode

## Success Metrics

| Metric | Target | Baseline |
|--------|--------|----------|
| Bounce Rate | <20% | 40% |
| Task Completion | >90% | 60% |

## Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Migration complexity | High | Phased rollout |
```

## PMView Fields

### Header

```go
type PMHeader struct {
    Title   string `json:"title"`
    Version string `json:"version"`
    Status  string `json:"status"`
    Author  string `json:"author"`
    PRDID   string `json:"prd_id"`
}
```

### PersonaSummary

```go
type PersonaSummary struct {
    Name      string `json:"name"`
    Role      string `json:"role"`
    IsPrimary bool   `json:"is_primary"`
}
```

### RequirementsList

```go
type RequirementsList struct {
    MustHave   []string `json:"must_have"`
    ShouldHave []string `json:"should_have"`
    CouldHave  []string `json:"could_have"`
    WontHave   []string `json:"wont_have"`
}
```

### MetricsSummary

```go
type MetricsSummary struct {
    Metrics []MetricItem `json:"metrics"`
}

type MetricItem struct {
    Name     string `json:"name"`
    Target   string `json:"target"`
    Baseline string `json:"baseline,omitempty"`
}
```

## Customization

### JSON Output

```go
json, err := pmView.ToJSON()
if err != nil {
    log.Fatal(err)
}
fmt.Println(string(json))
```

### Filtering

```go
// Get only must-have requirements
mustHaves := pmView.Requirements.MustHave

// Get primary persona
for _, p := range pmView.Personas {
    if p.IsPrimary {
        fmt.Printf("Primary: %s\n", p.Name)
    }
}
```

## Next Steps

- [Executive View](exec-view.md)
- [6-Pager](six-pager.md)
- [PR/FAQ](prfaq.md)
