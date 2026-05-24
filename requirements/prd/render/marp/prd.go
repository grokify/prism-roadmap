// Package marp provides a Marp markdown renderer for PRD documents.
// Marp is a presentation ecosystem that converts Markdown to slides.
// See https://marp.app/ for more information.
package marp

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-roadmap/requirements/prd"
	"github.com/grokify/prism-roadmap/requirements/prd/render"
)

// PRDRenderer implements the render.Renderer interface for PRD Marp output.
type PRDRenderer struct{}

// NewPRDRenderer creates a new PRD Marp renderer.
func NewPRDRenderer() *PRDRenderer {
	return &PRDRenderer{}
}

// Format returns the output format name.
func (r *PRDRenderer) Format() string {
	return "marp"
}

// FileExtension returns the file extension for Marp output.
func (r *PRDRenderer) FileExtension() string {
	return ".md"
}

// Render converts a PRD to Marp markdown slides.
func (r *PRDRenderer) Render(doc *prd.Document, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	data := &prdTemplateData{
		PRD:      doc,
		Options:  opts,
		Theme:    sdmarp.GetTheme(opts.Theme),
		Date:     time.Now().Format("January 2, 2006"),
		HasGoals: opts.IncludeGoals && doc.Goals != nil,
		HasRisks: opts.IncludeRisks && len(doc.Risks) > 0,
	}

	var buf bytes.Buffer

	// Render front matter
	if err := prdFrontMatterTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering front matter: %w", err)
	}

	// Render title slide
	if err := prdTitleSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering title slide: %w", err)
	}

	// Render problem slide
	if err := prdProblemSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering problem slide: %w", err)
	}

	// Render solution slide
	if err := prdSolutionSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering solution slide: %w", err)
	}

	// Render personas slide
	if len(doc.Personas) > 0 {
		if err := prdPersonasSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering personas slide: %w", err)
		}
	}

	// Render objectives slide
	if err := prdObjectivesSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering objectives slide: %w", err)
	}

	// Render success metrics slide (from OKR Key Results)
	hasKeyResults := false
	for _, okr := range doc.Objectives.OKRs {
		if len(okr.KeyResults) > 0 {
			hasKeyResults = true
			break
		}
	}
	if hasKeyResults {
		if err := prdMetricsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering metrics slide: %w", err)
		}
	}

	// Render requirements slides
	if opts.IncludeRequirements && (len(doc.Requirements.Functional) > 0 || len(doc.Requirements.NonFunctional) > 0) {
		if err := prdRequirementsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering requirements slide: %w", err)
		}
	}

	// Render roadmap slide
	if opts.IncludeRoadmap && len(doc.Roadmap.Phases) > 0 {
		if err := prdRoadmapSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering roadmap slide: %w", err)
		}
	}

	// Render risks slide
	if data.HasRisks {
		if err := prdRisksSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering risks slide: %w", err)
		}
	}

	// Render goals alignment slide
	if data.HasGoals {
		if err := prdGoalsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering goals slide: %w", err)
		}
	}

	// Render summary slide
	if err := prdSummarySlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering summary slide: %w", err)
	}

	return buf.Bytes(), nil
}

// prdTemplateData holds data for PRD template rendering.
type prdTemplateData struct {
	PRD      *prd.Document
	Options  *render.Options
	Theme    sdmarp.ThemeConfig
	Date     string
	HasGoals bool
	HasRisks bool
}

// prdFuncMap merges structureddocs CommonFuncMap with PRD-specific functions.
var prdFuncMap = mergeFuncMaps(sdmarp.CommonFuncMap, template.FuncMap{
	"moscowLabel": func(moscow interface{}) string {
		s := fmt.Sprintf("%v", moscow)
		switch s {
		case "must":
			return "Must Have"
		case "should":
			return "Should Have"
		case "could":
			return "Could Have"
		case "wont":
			return "Won't Have"
		default:
			return s
		}
	},
	"riskImpactColor": func(impact string) string {
		switch impact {
		case "critical":
			return "#e53e3e"
		case "high":
			return "#dd6b20"
		case "medium":
			return "#d69e2e"
		default:
			return "#38a169"
		}
	},
})

// mergeFuncMaps merges multiple template.FuncMaps into one.
// Later maps override earlier ones for duplicate keys.
func mergeFuncMaps(maps ...template.FuncMap) template.FuncMap {
	result := make(template.FuncMap)
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// Templates - updated to use structureddocs ThemeConfig field names
var prdFrontMatterTmpl = template.Must(template.New("prdFrontMatter").Parse(`---
marp: true
theme: {{.Theme.Name}}
paginate: true
{{- if .PRD.Metadata.Title}}
header: "PRD | {{.PRD.Metadata.Title}}"
{{- end}}
footer: "{{.PRD.Metadata.ID}} | v{{.PRD.Metadata.Version}}"
style: |
  section {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  section.title {
    text-align: center;
    background: linear-gradient(135deg, {{.Theme.PrimaryBgColor}} 0%, {{.Theme.AccentColor}} 100%);
    color: {{.Theme.PrimaryTextColor}};
  }
  section.title h1 {
    font-size: 2.5em;
    color: {{.Theme.PrimaryTextColor}};
  }
  section.problem {
    background: linear-gradient(135deg, {{.Theme.DangerColor}} 0%, #9b2c2c 100%);
    color: #ffffff;
  }
  section.solution {
    background: linear-gradient(135deg, {{.Theme.SuccessColor}} 0%, #276749 100%);
    color: #ffffff;
  }
  table {
    font-size: 0.85em;
    width: 100%;
  }
  th {
    background: #f7fafc;
  }
  .metric-card {
    background: #f7fafc;
    padding: 1em;
    border-radius: 8px;
    margin: 0.5em 0;
  }
  blockquote {
    font-size: 1.2em;
    border-left: 4px solid {{.Theme.AccentColor}};
    padding-left: 1em;
    font-style: italic;
  }
---

`))

var prdTitleSlideTmpl = template.Must(template.New("prdTitleSlide").Parse(`<!-- _class: title -->

# {{.PRD.Metadata.Title}}

**Product Requirements Document**

{{- if .PRD.Metadata.Authors}}
{{- range $i, $a := .PRD.Metadata.Authors}}
{{if eq $i 0}}**Author:** {{$a.Name}}{{if $a.Role}} ({{$a.Role}}){{end}}{{end}}
{{- end}}
{{- end}}
**Version:** {{.PRD.Metadata.Version}} | **Status:** {{.PRD.Metadata.Status}}
**Date:** {{.Date}}

---

`))

var prdProblemSlideTmpl = template.Must(template.New("prdProblemSlide").Parse(`<!-- _class: problem -->

## The Problem

> {{.PRD.ExecutiveSummary.ProblemStatement}}

{{- if .PRD.ExecutiveSummary.TargetAudience}}

**Who's Affected:** {{.PRD.ExecutiveSummary.TargetAudience}}
{{- end}}

{{- if .PRD.Problem}}
{{- if .PRD.Problem.Impact}}

**Impact:** {{.PRD.Problem.Impact}}
{{- end}}
{{- end}}

---

`))

var prdSolutionSlideTmpl = template.Must(template.New("prdSolutionSlide").Parse(`<!-- _class: solution -->

## The Solution

> {{.PRD.ExecutiveSummary.ProposedSolution}}

{{- if .PRD.ExecutiveSummary.ExpectedOutcomes}}

### Expected Outcomes

{{range .PRD.ExecutiveSummary.ExpectedOutcomes -}}
- {{.}}
{{end}}
{{- end}}

{{- if .PRD.ExecutiveSummary.ValueProposition}}

**Value Proposition:** {{.PRD.ExecutiveSummary.ValueProposition}}
{{- end}}

---

`))

var prdPersonasSlideTmpl = template.Must(template.New("prdPersonasSlide").Funcs(prdFuncMap).Parse(`## Target Personas

| Persona | Role | Primary | Key Goals |
|---------|------|---------|-----------|
{{- range .PRD.Personas}}
| **{{.Name}}** | {{.Role}} | {{if .IsPrimary}}Yes{{else}}No{{end}} | {{if .Goals}}{{index .Goals 0}}{{end}} |
{{- end}}

{{- if gt (len .PRD.Personas) 0}}
{{- with index .PRD.Personas 0}}

### Primary Persona: {{.Name}}

{{- if .PainPoints}}
**Pain Points:**
{{range .PainPoints -}}
- {{.}}
{{end}}
{{- end}}
{{- end}}
{{- end}}

---

`))

var prdObjectivesSlideTmpl = template.Must(template.New("prdObjectivesSlide").Funcs(prdFuncMap).Parse(`## OKRs (Objectives & Key Results)

{{- range $i, $okr := .PRD.Objectives.OKRs}}
### {{$okr.Objective.Description}}
{{- range $okr.KeyResults}}
- **{{.ID}}:** {{.Description}} → {{.Target}}
{{- end}}
{{end}}

---

`))

var prdMetricsSlideTmpl = template.Must(template.New("prdMetricsSlide").Parse(`## Key Results

| Key Result | Target | Baseline |
|------------|--------|----------|
{{- range .PRD.Objectives.OKRs}}
{{- range .KeyResults}}
| {{.Description}} | {{.Target}} | {{if .Baseline}}{{.Baseline}}{{else}}-{{end}} |
{{- end}}
{{- end}}

---

`))

var prdRequirementsSlideTmpl = template.Must(template.New("prdRequirementsSlide").Funcs(prdFuncMap).Parse(`## Key Requirements

{{- if .PRD.Requirements.Functional}}
### Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
{{- range $i, $r := .PRD.Requirements.Functional}}
{{- if lt $i 5}}
| {{$r.ID}} | {{truncate $r.Title 40}} | {{moscowLabel $r.Priority}} |
{{- end}}
{{- end}}
{{- end}}

{{- if .PRD.Requirements.NonFunctional}}
### Non-Functional Requirements

{{range $i, $r := .PRD.Requirements.NonFunctional -}}
{{- if lt $i 4}}
- **{{$r.Category}}:** {{truncate $r.Description 60}}
{{end}}
{{- end}}
{{- end}}

---

`))

var prdRoadmapSlideTmpl = template.Must(template.New("prdRoadmapSlide").Parse(`## Roadmap

| Phase | Timeline | Key Deliverables |
|-------|----------|------------------|
{{- range .PRD.Roadmap.Phases}}
| **{{.Name}}** | {{if .StartDate}}{{.StartDate}}{{end}}{{if .EndDate}} - {{.EndDate}}{{end}} | {{if .Goals}}{{index .Goals 0}}{{end}} |
{{- end}}

{{- if gt (len .PRD.Roadmap.Phases) 0}}

### Current Focus: {{(index .PRD.Roadmap.Phases 0).Name}}

{{- with index .PRD.Roadmap.Phases 0}}
{{- range .Goals}}
- {{.}}
{{- end}}
{{- end}}
{{- end}}

---

`))

var prdRisksSlideTmpl = template.Must(template.New("prdRisksSlide").Funcs(prdFuncMap).Parse(`## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
{{- range .PRD.Risks}}
| {{.Description}} | **{{.Impact}}** | {{if .Mitigation}}{{truncate .Mitigation 40}}{{else}}-{{end}} |
{{- end}}

---

`))

var prdGoalsSlideTmpl = template.Must(template.New("prdGoalsSlide").Parse(`## Goals Alignment

{{- if and .PRD.Goals .PRD.Goals.OKR}}
### OKR Alignment
{{range .PRD.Goals.OKR.Objectives -}}
**Objective:** {{.Title}}
{{range .KeyResults -}}
- KR: {{.Title}} (Target: {{.Target}})
{{end}}
{{end}}
{{- end}}

{{- if and .PRD.Goals .PRD.Goals.V2MOM}}
### V2MOM Alignment
**Vision:** {{.PRD.Goals.V2MOM.Vision}}

{{- if .PRD.Goals.V2MOM.Methods}}
**Methods:**
{{range .PRD.Goals.V2MOM.Methods -}}
- {{.Name}}
{{end}}
{{- end}}
{{- end}}

---

`))

var prdSummarySlideTmpl = template.Must(template.New("prdSummarySlide").Parse(`<!-- _class: title -->

## Summary

**Problem:** {{.PRD.ExecutiveSummary.ProblemStatement}}

**Solution:** {{.PRD.ExecutiveSummary.ProposedSolution}}

{{- if .PRD.Objectives.OKRs}}
{{- with index .PRD.Objectives.OKRs 0}}
{{- if .KeyResults}}
{{- with index .KeyResults 0}}
**Key Metric:** {{.Description}} → {{.Target}}
{{- end}}
{{- end}}
{{- end}}
{{- end}}

---

## Questions?

**{{.PRD.Metadata.Title}}**
{{.PRD.Metadata.ID}} | v{{.PRD.Metadata.Version}}

{{- if .PRD.Metadata.Authors}}
{{range .PRD.Metadata.Authors -}}
{{.Name}}{{if .Email}} ({{.Email}}){{end}}
{{end}}
{{- end}}

`))
