// Package marp provides a Marp markdown renderer for PRD documents.
// This file contains the PRD+Goals combined renderer that includes
// V2MOM and OKR alignment slides.
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

// PRDGoalsRenderer renders PRD documents with expanded goals alignment slides.
type PRDGoalsRenderer struct{}

// NewPRDGoalsRenderer creates a new PRD+Goals Marp renderer.
func NewPRDGoalsRenderer() *PRDGoalsRenderer {
	return &PRDGoalsRenderer{}
}

// Format returns the output format name.
func (r *PRDGoalsRenderer) Format() string {
	return "marp-prd-goals"
}

// FileExtension returns the file extension for Marp output.
func (r *PRDGoalsRenderer) FileExtension() string {
	return ".md"
}

// Render converts a PRD with embedded Goals to Marp markdown slides.
func (r *PRDGoalsRenderer) Render(doc *prd.Document, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	hasV2MOM := opts.IncludeGoals && doc.Goals != nil && doc.Goals.V2MOM != nil
	hasOKR := opts.IncludeGoals && doc.Goals != nil && doc.Goals.OKR != nil

	data := &prdGoalsTemplateData{
		PRD:      doc,
		Options:  opts,
		Theme:    sdmarp.GetTheme(opts.Theme),
		Date:     time.Now().Format("January 2, 2006"),
		HasGoals: hasV2MOM || hasOKR,
		HasV2MOM: hasV2MOM,
		HasOKR:   hasOKR,
		HasRisks: opts.IncludeRisks && len(doc.Risks) > 0,
	}

	var buf bytes.Buffer

	// Render front matter
	if err := prdGoalsFrontMatterTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering front matter: %w", err)
	}

	// Render title slide
	if err := prdGoalsTitleSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering title slide: %w", err)
	}

	// Render agenda slide (shows structure with goals sections)
	if err := prdGoalsAgendaSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering agenda slide: %w", err)
	}

	// Render problem slide
	if err := prdGoalsProblemSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering problem slide: %w", err)
	}

	// Render solution slide
	if err := prdGoalsSolutionSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering solution slide: %w", err)
	}

	// Render personas slide
	if len(doc.Personas) > 0 {
		if err := prdGoalsPersonasSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering personas slide: %w", err)
		}
	}

	// Render objectives slide
	if err := prdGoalsObjectivesSlideTmpl.Execute(&buf, data); err != nil {
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
		if err := prdGoalsMetricsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering metrics slide: %w", err)
		}
	}

	// Render V2MOM alignment slides (if present)
	if hasV2MOM {
		if err := v2momVisionSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering V2MOM vision slide: %w", err)
		}
		if err := v2momMethodsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering V2MOM methods slide: %w", err)
		}
	}

	// Render OKR alignment slides (if present)
	if hasOKR {
		if err := okrOverviewSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering OKR overview slide: %w", err)
		}
		if err := okrKeyResultsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering OKR key results slide: %w", err)
		}
	}

	// Render requirements slides
	if opts.IncludeRequirements && (len(doc.Requirements.Functional) > 0 || len(doc.Requirements.NonFunctional) > 0) {
		if err := prdGoalsRequirementsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering requirements slide: %w", err)
		}
	}

	// Render roadmap slide
	if opts.IncludeRoadmap && len(doc.Roadmap.Phases) > 0 {
		if err := prdGoalsRoadmapSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering roadmap slide: %w", err)
		}
	}

	// Render risks slide
	if data.HasRisks {
		if err := prdGoalsRisksSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering risks slide: %w", err)
		}
	}

	// Render alignment summary (connecting PRD to goals)
	if data.HasGoals && doc.Goals.AlignedObjectives != nil && len(doc.Goals.AlignedObjectives) > 0 {
		if err := alignmentSummarySlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering alignment summary slide: %w", err)
		}
	}

	// Render summary slide
	if err := prdGoalsSummarySlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering summary slide: %w", err)
	}

	return buf.Bytes(), nil
}

// prdGoalsTemplateData holds data for PRD+Goals template rendering.
type prdGoalsTemplateData struct {
	PRD      *prd.Document
	Options  *render.Options
	Theme    sdmarp.ThemeConfig
	Date     string
	HasGoals bool
	HasV2MOM bool
	HasOKR   bool
	HasRisks bool
}

// prdGoalsFuncMap merges structureddocs CommonFuncMap with PRD+Goals-specific functions.
var prdGoalsFuncMap = mergeFuncMaps(sdmarp.CommonFuncMap, template.FuncMap{
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
})

// Templates
var prdGoalsFrontMatterTmpl = template.Must(template.New("prdGoalsFrontMatter").Parse(`---
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
    background: linear-gradient(135deg, #c53030 0%, #9b2c2c 100%);
    color: #ffffff;
  }
  section.solution {
    background: linear-gradient(135deg, #2f855a 0%, #276749 100%);
    color: #ffffff;
  }
  section.goals {
    background: linear-gradient(135deg, #5a67d8 0%, #7f9cf5 100%);
    color: #ffffff;
  }
  section.v2mom {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: #ffffff;
  }
  section.okr {
    background: linear-gradient(135deg, #38b2ac 0%, #319795 100%);
    color: #ffffff;
  }
  table {
    font-size: 0.85em;
    width: 100%;
  }
  th {
    background: #f7fafc;
  }
  blockquote {
    font-size: 1.2em;
    border-left: 4px solid {{.Theme.AccentColor}};
    padding-left: 1em;
    font-style: italic;
  }
---

`))

var prdGoalsTitleSlideTmpl = template.Must(template.New("prdGoalsTitleSlide").Parse(`<!-- _class: title -->

# {{.PRD.Metadata.Title}}

**Product Requirements Document**
{{- if .HasGoals}}
_with Goals Alignment_
{{- end}}

{{- if .PRD.Metadata.Authors}}
{{- range $i, $a := .PRD.Metadata.Authors}}
{{if eq $i 0}}**Author:** {{$a.Name}}{{if $a.Role}} ({{$a.Role}}){{end}}{{end}}
{{- end}}
{{- end}}
**Version:** {{.PRD.Metadata.Version}} | **Status:** {{.PRD.Metadata.Status}}
**Date:** {{.Date}}

---

`))

var prdGoalsAgendaSlideTmpl = template.Must(template.New("prdGoalsAgendaSlide").Parse(`## Agenda

1. **Problem & Solution**
2. **Target Users**
3. **Objectives & Metrics**
{{- if .HasV2MOM}}
4. **V2MOM Alignment** (Vision, Values, Methods)
{{- end}}
{{- if .HasOKR}}
{{if .HasV2MOM}}5{{else}}4{{end}}. **OKR Alignment** (Objectives & Key Results)
{{- end}}
{{- if .Options.IncludeRequirements}}
{{if .HasV2MOM}}{{if .HasOKR}}6{{else}}5{{end}}{{else}}{{if .HasOKR}}5{{else}}4{{end}}{{end}}. **Requirements**
{{- end}}
{{- if .Options.IncludeRoadmap}}
- **Roadmap**
{{- end}}
{{- if .HasRisks}}
- **Risks**
{{- end}}
- **Summary**

---

`))

var prdGoalsProblemSlideTmpl = template.Must(template.New("prdGoalsProblemSlide").Parse(`<!-- _class: problem -->

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

var prdGoalsSolutionSlideTmpl = template.Must(template.New("prdGoalsSolutionSlide").Parse(`<!-- _class: solution -->

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

var prdGoalsPersonasSlideTmpl = template.Must(template.New("prdGoalsPersonasSlide").Funcs(prdGoalsFuncMap).Parse(`## Target Personas

| Persona | Role | Primary | Key Goals |
|---------|------|---------|-----------
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

var prdGoalsObjectivesSlideTmpl = template.Must(template.New("prdGoalsObjectivesSlide").Funcs(prdGoalsFuncMap).Parse(`## OKRs (Objectives & Key Results)

{{- range $i, $okr := .PRD.Objectives.OKRs}}
### {{$okr.Objective.Description}}
{{- range $okr.KeyResults}}
- **{{.ID}}:** {{.Description}} → {{.Target}}
{{- end}}
{{end}}

---

`))

var prdGoalsMetricsSlideTmpl = template.Must(template.New("prdGoalsMetricsSlide").Parse(`## Key Results

| Key Result | Target | Baseline |
|------------|--------|----------|
{{- range .PRD.Objectives.OKRs}}
{{- range .KeyResults}}
| {{.Description}} | {{.Target}} | {{if .Baseline}}{{.Baseline}}{{else}}-{{end}} |
{{- end}}
{{- end}}

---

`))

// V2MOM alignment slides
var v2momVisionSlideTmpl = template.Must(template.New("v2momVisionSlide").Funcs(prdGoalsFuncMap).Parse(`<!-- _class: v2mom -->

## V2MOM: Vision & Values

### Vision
> {{.PRD.Goals.V2MOM.Vision}}

{{- if .PRD.Goals.V2MOM.Values}}

### Values
{{range $i, $v := .PRD.Goals.V2MOM.Values -}}
{{$num := (add $i 1)}}**{{$num}}. {{$v.Name}}**{{if $v.Description}} - {{$v.Description}}{{end}}
{{end}}
{{- end}}

---

`))

var v2momMethodsSlideTmpl = template.Must(template.New("v2momMethodsSlide").Funcs(prdGoalsFuncMap).Parse(`<!-- _class: v2mom -->

## V2MOM: Methods

{{- if .PRD.Goals.V2MOM.Methods}}

| # | Method | Priority | Status |
|---|--------|----------|--------|
{{- range $i, $m := .PRD.Goals.V2MOM.Methods}}
| {{add $i 1}} | {{$m.Name}} | {{if $m.Priority}}{{$m.Priority}}{{else}}-{{end}} | {{if $m.Status}}{{$m.Status}}{{else}}-{{end}} |
{{- end}}
{{- end}}

{{- if .PRD.Goals.V2MOM.Obstacles}}

### Key Obstacles
{{range .PRD.Goals.V2MOM.Obstacles -}}
- **{{.Name}}**{{if .Severity}} ({{.Severity}}){{end}}
{{end}}
{{- end}}

---

`))

// OKR alignment slides
var okrOverviewSlideTmpl = template.Must(template.New("okrOverviewSlide").Funcs(prdGoalsFuncMap).Parse(`<!-- _class: okr -->

## OKR Alignment

{{- if .PRD.Goals.OKR.Theme}}
> {{.PRD.Goals.OKR.Theme}}
{{- end}}

### Objectives Overview

| # | Objective | Progress | Key Results |
|---|-----------|----------|-------------|
{{- range $i, $obj := .PRD.Goals.OKR.Objectives}}
| {{add $i 1}} | {{truncate $obj.Title 35}} | {{scorePercent $obj.CalculateProgress}} | {{len $obj.KeyResults}} |
{{- end}}

---

`))

var okrKeyResultsSlideTmpl = template.Must(template.New("okrKeyResultsSlide").Funcs(prdGoalsFuncMap).Parse(`<!-- _class: okr -->

## OKR: Key Results

| Objective | Key Result | Target | Score |
|-----------|------------|--------|-------|
{{- range $oi, $obj := .PRD.Goals.OKR.Objectives}}
{{- range $ki, $kr := $obj.KeyResults}}
| O{{add $oi 1}} | {{truncate $kr.Title 25}} | {{if $kr.Target}}{{$kr.Target}}{{else}}-{{end}} | {{scorePercent $kr.Score}} |
{{- end}}
{{- end}}

---

`))

var prdGoalsRequirementsSlideTmpl = template.Must(template.New("prdGoalsRequirementsSlide").Funcs(prdGoalsFuncMap).Parse(`## Key Requirements

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

var prdGoalsRoadmapSlideTmpl = template.Must(template.New("prdGoalsRoadmapSlide").Parse(`## Roadmap

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

var prdGoalsRisksSlideTmpl = template.Must(template.New("prdGoalsRisksSlide").Funcs(prdGoalsFuncMap).Parse(`## Risks & Mitigations

| Risk | Impact | Mitigation |
|------|--------|------------|
{{- range .PRD.Risks}}
| {{.Description}} | **{{.Impact}}** | {{if .Mitigation}}{{truncate .Mitigation 40}}{{else}}-{{end}} |
{{- end}}

---

`))

var alignmentSummarySlideTmpl = template.Must(template.New("alignmentSummarySlide").Parse(`<!-- _class: goals -->

## Goals Alignment Summary

### How PRD Objectives Map to Goals

| PRD Objective | Aligned Goal |
|---------------|--------------|
{{- range $key, $value := .PRD.Goals.AlignedObjectives}}
| {{$key}} | {{$value}} |
{{- end}}

{{- if .HasV2MOM}}

**V2MOM Vision:** {{.PRD.Goals.V2MOM.Vision}}
{{- end}}

{{- if .HasOKR}}
{{- if .PRD.Goals.OKR.Theme}}

**OKR Theme:** {{.PRD.Goals.OKR.Theme}}
{{- end}}
{{- end}}

---

`))

var prdGoalsSummarySlideTmpl = template.Must(template.New("prdGoalsSummarySlide").Parse(`<!-- _class: title -->

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

{{- if .HasGoals}}
**Aligned to:** {{if .HasV2MOM}}V2MOM{{end}}{{if and .HasV2MOM .HasOKR}} + {{end}}{{if .HasOKR}}OKR{{end}}
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
