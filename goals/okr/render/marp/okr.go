// Package marp provides a Marp markdown renderer for OKR documents.
// Marp is a presentation ecosystem that converts Markdown to slides.
// See https://marp.app/ for more information.
package marp

import (
	"bytes"
	"fmt"
	"text/template"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/goals/okr/render"
)

// Renderer implements the render.Renderer interface for OKR Marp output.
type Renderer struct{}

// New creates a new OKR Marp renderer.
func New() *Renderer {
	return &Renderer{}
}

// Format returns the output format name.
func (r *Renderer) Format() string {
	return "marp"
}

// FileExtension returns the file extension for Marp output.
func (r *Renderer) FileExtension() string {
	return ".md"
}

// Render converts an OKR document to Marp markdown slides.
func (r *Renderer) Render(doc *okr.OKRDocument, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	data := &templateData{
		OKR:             doc,
		Options:         opts,
		Theme:           sdmarp.GetTheme(opts.Theme),
		OverallProgress: doc.CalculateOverallProgress(),
		HasRisks:        opts.IncludeRisks && len(doc.AllRisks()) > 0,
	}

	var buf bytes.Buffer

	// Render front matter
	if err := frontMatterTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering front matter: %w", err)
	}

	// Render title slide
	if err := titleSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering title slide: %w", err)
	}

	// Render overall progress slide
	if err := overviewSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering overview slide: %w", err)
	}

	// Render objective slides
	objectives := doc.Objectives
	if opts.MaxObjectives > 0 && len(objectives) > opts.MaxObjectives {
		objectives = objectives[:opts.MaxObjectives]
	}

	for i, obj := range objectives {
		objData := &objectiveData{
			templateData: data,
			Objective:    obj,
			Index:        i + 1,
		}
		if err := objectiveSlideTmpl.Execute(&buf, objData); err != nil {
			return nil, fmt.Errorf("rendering objective slide %d: %w", i+1, err)
		}
	}

	// Render all key results summary
	if len(doc.AllKeyResults()) > 0 {
		if err := keyResultsSummaryTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering key results summary: %w", err)
		}
	}

	// Render risks slide
	if data.HasRisks {
		if err := risksSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering risks slide: %w", err)
		}
	}

	// Render summary slide
	if err := summarySlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering summary slide: %w", err)
	}

	return buf.Bytes(), nil
}

// templateData holds data for OKR template rendering.
type templateData struct {
	OKR             *okr.OKRDocument
	Options         *render.Options
	Theme           sdmarp.ThemeConfig
	OverallProgress float64
	HasRisks        bool
}

// objectiveData holds data for objective slide rendering.
type objectiveData struct {
	*templateData
	Objective okr.Objective
	Index     int
}

// funcMap merges structureddocs CommonFuncMap with OKR-specific functions.
var funcMap = mergeFuncMaps(sdmarp.CommonFuncMap, template.FuncMap{
	"scoreGrade":       okr.ScoreGrade,
	"scoreDescription": okr.ScoreDescription,
	"confidenceIcon": func(confidence string) string {
		switch confidence {
		case "High":
			return "[HIGH]"
		case "Medium":
			return "[MED]"
		case "Low":
			return "[LOW]"
		default:
			return ""
		}
	},
	"impactColor": func(impact string) string {
		switch impact {
		case "Critical":
			return "#e53e3e"
		case "High":
			return "#dd6b20"
		case "Medium":
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

// Templates
var frontMatterTmpl = template.Must(template.New("frontMatter").Parse(`---
marp: true
theme: {{.Theme.Name}}
paginate: true
{{- if and .OKR.Metadata .OKR.Metadata.Team}}
header: "{{.OKR.Metadata.Team}}{{if .OKR.Metadata.Period}} | {{.OKR.Metadata.Period}}{{end}}"
{{- end}}
{{- if and .OKR.Metadata .OKR.Metadata.Name}}
footer: "OKR | {{.OKR.Metadata.Name}}"
{{- end}}
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
  section.objective {
    background: linear-gradient(135deg, {{.Theme.PrimaryBgColor}} 0%, {{.Theme.AccentColor}} 100%);
    color: {{.Theme.PrimaryTextColor}};
  }
  table {
    font-size: 0.85em;
    width: 100%;
  }
  th {
    background: #f7fafc;
  }
  .score-excellent { color: {{.Theme.SuccessColor}}; }
  .score-good { color: #3182ce; }
  .score-partial { color: {{.Theme.WarningColor}}; }
  .score-failed { color: {{.Theme.DangerColor}}; }
  blockquote {
    font-size: 1.2em;
    border-left: 4px solid {{.Theme.AccentColor}};
    padding-left: 1em;
    font-style: italic;
  }
  .progress-container {
    background: #e2e8f0;
    border-radius: 4px;
    padding: 2px;
  }
  .progress-bar {
    background: {{.Theme.AccentColor}};
    border-radius: 2px;
    height: 20px;
  }
---

`))

var titleSlideTmpl = template.Must(template.New("titleSlide").Funcs(funcMap).Parse(`<!-- _class: title -->

# {{if and .OKR.Metadata .OKR.Metadata.Name}}{{.OKR.Metadata.Name}}{{else}}OKR{{end}}

**Objectives and Key Results**

{{- if and .OKR.Metadata .OKR.Metadata.Owner}}
**Owner:** {{.OKR.Metadata.Owner}}
{{- end}}
{{- if and .OKR.Metadata .OKR.Metadata.Team}}
**Team:** {{.OKR.Metadata.Team}}
{{- end}}
{{- if and .OKR.Metadata .OKR.Metadata.Period}}
**Period:** {{.OKR.Metadata.Period}}
{{- end}}
{{- if and .OKR.Metadata .OKR.Metadata.Status}}
**Status:** {{.OKR.Metadata.Status}}
{{- end}}

---

`))

var overviewSlideTmpl = template.Must(template.New("overviewSlide").Funcs(funcMap).Parse(`## OKR Overview

{{- if .OKR.Theme}}

> {{.OKR.Theme}}
{{- end}}

**Overall Progress:** [{{progressBar .OverallProgress}}] {{scorePercent .OverallProgress}}
{{- if .Options.ShowScoreGrades}}
**Grade:** {{scoreGrade .OverallProgress}} - {{scoreDescription .OverallProgress}}
{{- end}}

### Objectives Summary

| # | Objective | Progress | Key Results |
|---|-----------|----------|-------------|
{{- range $i, $obj := .OKR.Objectives}}
| {{add $i 1}} | {{truncate $obj.Title 35}} | {{scorePercent $obj.CalculateProgress}} | {{len $obj.KeyResults}} |
{{- end}}

---

`))

var objectiveSlideTmpl = template.Must(template.New("objectiveSlide").Funcs(funcMap).Parse(`<!-- _class: objective -->

## O{{.Index}}: {{.Objective.Title}}

{{- if .Objective.Description}}

{{.Objective.Description}}
{{- end}}

{{- if .Objective.Owner}}
**Owner:** {{.Objective.Owner}}
{{- end}}

**Progress:** [{{progressBar .Objective.CalculateProgress}}] {{scorePercent .Objective.CalculateProgress}}

---

## O{{.Index}} Key Results

| KR | Key Result | Target | Score | Status |
|----|------------|--------|-------|--------|
{{- range $i, $kr := .Objective.KeyResults}}
| {{add $i 1}} | {{truncate $kr.Title 30}} | {{if $kr.Target}}{{$kr.Target}}{{else}}-{{end}} | {{scorePercent $kr.Score}} | {{if $kr.Status}}{{statusEmoji $kr.Status}}{{else}}-{{end}} |
{{- end}}

{{- if .Objective.Risks}}

### Risks

{{range .Objective.Risks -}}
- **{{.Title}}** ({{.Impact}}): {{if .Mitigation}}{{.Mitigation}}{{else}}No mitigation defined{{end}}
{{end}}
{{- end}}

---

`))

var keyResultsSummaryTmpl = template.Must(template.New("keyResultsSummary").Funcs(funcMap).Parse(`## All Key Results

| Objective | Key Result | Target | Score | Confidence |
|-----------|------------|--------|-------|------------|
{{- range $oi, $obj := .OKR.Objectives}}
{{- range $ki, $kr := $obj.KeyResults}}
| O{{add $oi 1}} | {{truncate $kr.Title 25}} | {{if $kr.Target}}{{$kr.Target}}{{else}}-{{end}} | {{scorePercent $kr.Score}} | {{if $kr.Confidence}}{{confidenceIcon $kr.Confidence}}{{else}}-{{end}} |
{{- end}}
{{- end}}

---

`))

var risksSlideTmpl = template.Must(template.New("risksSlide").Funcs(funcMap).Parse(`## Risks & Challenges

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
{{- range .OKR.AllRisks}}
| {{.Title}} | **{{.Impact}}** | {{if .Likelihood}}{{.Likelihood}}{{else}}-{{end}} | {{if .Mitigation}}{{truncate .Mitigation 40}}{{else}}-{{end}} |
{{- end}}

---

`))

var summarySlideTmpl = template.Must(template.New("summarySlide").Funcs(funcMap).Parse(`<!-- _class: title -->

## Summary

**Overall Progress:** {{scorePercent .OverallProgress}}
{{- if .Options.ShowScoreGrades}}
**Grade:** {{scoreGrade .OverallProgress}}
{{- end}}

{{- range $i, $obj := .OKR.Objectives}}
- **O{{add $i 1}}:** {{truncate $obj.Title 40}} ({{scorePercent $obj.CalculateProgress}})
{{- end}}

---

## Questions?

{{- if and .OKR.Metadata .OKR.Metadata.Name}}
**{{.OKR.Metadata.Name}}**
{{- end}}
{{- if and .OKR.Metadata .OKR.Metadata.Period}}
{{.OKR.Metadata.Period}}
{{- end}}

{{- if and .OKR.Metadata .OKR.Metadata.Owner}}
{{.OKR.Metadata.Owner}}
{{- end}}

`))
