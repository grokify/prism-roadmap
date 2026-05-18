// Package marp provides a Marp markdown renderer for V2MOM documents.
// Marp is a presentation ecosystem that converts Markdown to slides.
// See https://marp.app/ for more information.
package marp

import (
	"bytes"
	"fmt"
	"text/template"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-execution/goals/v2mom"
	"github.com/grokify/prism-execution/goals/v2mom/render"
)

// Renderer implements the render.Renderer interface for Marp output.
type Renderer struct{}

// New creates a new Marp renderer.
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

// Render converts a V2MOM to Marp markdown.
func (r *Renderer) Render(v *v2mom.V2MOM, opts *render.Options) ([]byte, error) {
	if opts == nil {
		opts = render.DefaultOptions()
	}

	term := v2mom.GetTerminologyLabels(opts.GetTerminology(v))
	structure := opts.GetStructure(v)

	data := &templateData{
		V2MOM:       v,
		Options:     opts,
		Term:        term,
		Structure:   structure,
		Theme:       sdmarp.GetTheme(opts.Theme),
		HasProjects: len(v.Projects) > 0 && opts.IncludeProjects,
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

	// Render vision slide
	if err := visionSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering vision slide: %w", err)
	}

	// Render values slide
	if err := valuesSlideTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering values slide: %w", err)
	}

	// Render methods overview slide
	if err := methodsOverviewTmpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("rendering methods overview: %w", err)
	}

	// Render method detail slides (for nested/hybrid structure)
	if structure == v2mom.StructureNested || structure == v2mom.StructureHybrid {
		for i, method := range v.Methods {
			if len(method.Measures) > 0 || len(method.Obstacles) > 0 {
				methodData := &methodDetailData{
					templateData: data,
					Method:       method,
					Index:        i + 1,
				}
				if err := methodDetailTmpl.Execute(&buf, methodData); err != nil {
					return nil, fmt.Errorf("rendering method detail %d: %w", i, err)
				}
			}
		}
	}

	// Render obstacles slide (global obstacles)
	if len(v.Obstacles) > 0 {
		if err := obstaclesSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering obstacles slide: %w", err)
		}
	}

	// Render measures slide (global measures for flat structure)
	if len(v.Measures) > 0 {
		if err := measuresSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering measures slide: %w", err)
		}
	}

	// Render measures dashboard (all measures summary)
	allMeasures := v.AllMeasures()
	if len(allMeasures) > 0 {
		dashData := &measuresDashboardData{
			templateData: data,
			AllMeasures:  allMeasures,
		}
		if err := measuresDashboardTmpl.Execute(&buf, dashData); err != nil {
			return nil, fmt.Errorf("rendering measures dashboard: %w", err)
		}
	}

	// Render projects/roadmap slide
	if data.HasProjects {
		if err := projectsSlideTmpl.Execute(&buf, data); err != nil {
			return nil, fmt.Errorf("rendering projects slide: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// templateData holds data for template rendering.
type templateData struct {
	V2MOM       *v2mom.V2MOM
	Options     *render.Options
	Term        v2mom.Terminology
	Structure   string
	Theme       sdmarp.ThemeConfig
	HasProjects bool
}

// methodDetailData holds data for method detail slide rendering.
type methodDetailData struct {
	*templateData
	Method v2mom.Method
	Index  int
}

// measuresDashboardData holds data for measures dashboard rendering.
type measuresDashboardData struct {
	*templateData
	AllMeasures []v2mom.Measure
}

// funcMap uses the shared CommonFuncMap from structureddocs.
var funcMap = sdmarp.CommonFuncMap

// Templates - updated to use structureddocs ThemeConfig field names
var frontMatterTmpl = template.Must(template.New("frontMatter").Parse(`---
marp: true
theme: {{.Theme.Name}}
paginate: true
{{- if and .V2MOM.Metadata .V2MOM.Metadata.Team}}
header: "{{.V2MOM.Metadata.Team}}{{if .V2MOM.Metadata.FiscalYear}} | {{.V2MOM.Metadata.FiscalYear}}{{end}}{{if .V2MOM.Metadata.Quarter}} {{.V2MOM.Metadata.Quarter}}{{end}}"
{{- end}}
{{- if and .V2MOM.Metadata .V2MOM.Metadata.Name}}
footer: "V2MOM | {{.V2MOM.Metadata.Name}}"
{{- end}}
style: |
  section {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  section.title {
    text-align: center;
  }
  section.title h1 {
    font-size: 2.5em;
  }
  section.vision {
    background: linear-gradient(135deg, {{.Theme.PrimaryBgColor}} 0%, {{.Theme.AccentColor}} 100%);
    color: {{.Theme.PrimaryTextColor}};
  }
  section.vision h2 {
    color: {{.Theme.PrimaryTextColor}};
  }
  table {
    font-size: 0.85em;
    width: 100%;
  }
  th {
    background: #f7fafc;
  }
  .status-done { color: {{.Theme.SuccessColor}}; }
  .status-progress { color: {{.Theme.AccentColor}}; }
  .status-risk { color: {{.Theme.DangerColor}}; }
  .status-planning { color: {{.Theme.WarningColor}}; }
  blockquote {
    font-size: 1.3em;
    font-style: italic;
    border-left: 4px solid {{.Theme.AccentColor}};
    padding-left: 1em;
  }
---

`))

var titleSlideTmpl = template.Must(template.New("titleSlide").Parse(`<!-- _class: title -->

# {{if and .V2MOM.Metadata .V2MOM.Metadata.Name}}{{.V2MOM.Metadata.Name}}{{else}}V2MOM{{end}}

{{- if and .V2MOM.Metadata .V2MOM.Metadata.Author}}
**Author:** {{.V2MOM.Metadata.Author}}
{{- end}}
{{- if and .V2MOM.Metadata .V2MOM.Metadata.Team}}
**Team:** {{.V2MOM.Metadata.Team}}
{{- end}}
{{- if and .V2MOM.Metadata .V2MOM.Metadata.FiscalYear}}
**Period:** {{.V2MOM.Metadata.FiscalYear}}{{if .V2MOM.Metadata.Quarter}} {{.V2MOM.Metadata.Quarter}}{{end}}
{{- end}}
{{- if and .V2MOM.Metadata .V2MOM.Metadata.Status}}
**Status:** {{.V2MOM.Metadata.Status}}
{{- end}}

---

`))

var visionSlideTmpl = template.Must(template.New("visionSlide").Parse(`<!-- _class: vision -->

## Vision

> {{.V2MOM.Vision}}

---

`))

var valuesSlideTmpl = template.Must(template.New("valuesSlide").Funcs(funcMap).Parse(`## Values

{{range $i, $v := .V2MOM.Values -}}
{{$num := (printf "%d" (add $i 1))}}**{{$num}}. {{$v.Name}}**{{if $v.Description}} - {{$v.Description}}{{end}}

{{end}}
---

`))

var methodsOverviewTmpl = template.Must(template.New("methodsOverview").Funcs(funcMap).Parse(`## {{.Term.Methods}}

| # | {{.Term.MethodSingular}} | Priority | Status |
|---|--------|----------|--------|
{{range $i, $m := .V2MOM.Methods -}}
| {{add $i 1}} | {{$m.Name}} | {{if $m.Priority}}{{priorityLabel $m.Priority}}{{else}}-{{end}} | {{if $m.Status}}{{$m.Status}}{{else}}-{{end}} |
{{end}}
---

`))

var methodDetailTmpl = template.Must(template.New("methodDetail").Funcs(funcMap).Parse(`## {{.Term.MethodSingular}} {{.Index}}: {{.Method.Name}}

{{if .Method.Description}}{{.Method.Description}}

{{end -}}
{{if .Method.Owner}}**Owner:** {{.Method.Owner}}{{end}}
{{if .Method.Priority}}**Priority:** {{priorityLabel .Method.Priority}}{{end}}
{{if .Method.Status}}**Status:** {{.Method.Status}}{{end}}

{{if .Method.Measures}}### {{.Term.Measures}}

| {{.Term.MeasureSingular}} | Target | Status | Progress |
|---------|--------|--------|----------|
{{range .Method.Measures -}}
| {{.Name}} | {{if .Target}}{{.Target}}{{else}}-{{end}} | {{if .Status}}{{.Status}}{{else}}-{{end}} | {{if .Progress}}{{progressPercent .Progress}}{{else}}-{{end}} |
{{end}}
{{end}}
{{if .Method.Obstacles}}### {{.Term.Obstacles}}

{{range .Method.Obstacles -}}
- **{{.Name}}**{{if .Severity}} ({{.Severity}}){{end}}{{if .Mitigation}}: {{.Mitigation}}{{end}}
{{end}}
{{end}}
---

`))

var obstaclesSlideTmpl = template.Must(template.New("obstaclesSlide").Funcs(funcMap).Parse(`## {{.Term.Obstacles}}

| {{.Term.ObstacleSingular}} | Severity | Likelihood | Status |
|----------|----------|------------|--------|
{{range .V2MOM.Obstacles -}}
| {{.Name}} | {{if .Severity}}{{.Severity}}{{else}}-{{end}} | {{if .Likelihood}}{{.Likelihood}}{{else}}-{{end}} | {{if .Status}}{{.Status}}{{else}}-{{end}} |
{{end}}
{{if gt (len .V2MOM.Obstacles) 0}}
### Mitigation Strategies

{{range .V2MOM.Obstacles -}}
{{if .Mitigation}}- **{{.Name}}:** {{.Mitigation}}
{{end -}}
{{end}}
{{end}}
---

`))

var measuresSlideTmpl = template.Must(template.New("measuresSlide").Funcs(funcMap).Parse(`## {{.Term.Measures}} (Global)

| {{.Term.MeasureSingular}} | Baseline | Target | Current | Status |
|---------|----------|--------|---------|--------|
{{range .V2MOM.Measures -}}
| {{.Name}} | {{if .Baseline}}{{.Baseline}}{{else}}-{{end}} | {{if .Target}}{{.Target}}{{else}}-{{end}} | {{if .Current}}{{.Current}}{{else}}-{{end}} | {{if .Status}}{{.Status}}{{else}}-{{end}} |
{{end}}
---

`))

var measuresDashboardTmpl = template.Must(template.New("measuresDashboard").Funcs(funcMap).Parse(`## {{.Term.Measures}} Dashboard

| {{.Term.MeasureSingular}} | Target | Progress | Status |
|---------|--------|----------|--------|
{{range .AllMeasures -}}
| {{.Name}} | {{if .Target}}{{.Target}}{{else}}-{{end}} | {{if .Progress}}[{{progressBar .Progress}}] {{progressPercent .Progress}}{{else}}-{{end}} | {{if .Status}}{{.Status}}{{else}}-{{end}} |
{{end}}
---

`))

var projectsSlideTmpl = template.Must(template.New("projectsSlide").Funcs(funcMap).Parse(`## Roadmap Projects

| Project | Priority | Quarter | Status |
|---------|----------|---------|--------|
{{range .V2MOM.Projects -}}
| {{.Name}} | {{if .Priority}}{{.Priority}}{{else}}-{{end}} | {{if .Quarter}}{{.Quarter}}{{else}}-{{end}} | {{if .Status}}{{.Status}}{{else}}-{{end}} |
{{end}}
---

`))
