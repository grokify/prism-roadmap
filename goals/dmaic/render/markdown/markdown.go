// Package markdown renders a DMAIC metrics document as portable Markdown: an
// overview with the health of each category, then each category with its
// metrics and the detail of every metric that needs attention, then
// optionally the phase distribution, the root causes, the improvement
// initiatives, and the process capability appendix. It covers the content of
// the pandoc renderer without its front matter or LaTeX, so the output reads
// on any Markdown viewer and can be embedded in a larger document. The output
// is deterministic: the same document always renders to the same bytes.
package markdown

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/dmaic"
	"github.com/grokify/prism-roadmap/goals/dmaic/render"
	"github.com/grokify/prism-roadmap/goals/internal/mdshift"
)

// Renderer implements render.Renderer for Markdown.
type Renderer struct{}

// New returns a Markdown renderer.
func New() *Renderer {
	return &Renderer{}
}

// Format returns "markdown".
func (r *Renderer) Format() string {
	return "markdown"
}

// FileExtension returns ".md".
func (r *Renderer) FileExtension() string {
	return ".md"
}

// Render converts the DMAIC document to Markdown.
func (r *Renderer) Render(doc *dmaic.DMAICDocument, opts *render.Options) ([]byte, error) {
	if doc == nil {
		return nil, fmt.Errorf("markdown: nil DMAIC document")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	var b strings.Builder
	if !opts.OmitHeader {
		fmt.Fprintf(&b, "# %s\n\n", Title(doc))
		if line := metadataLine(doc.Metadata); line != "" {
			fmt.Fprintf(&b, "%s\n\n", line)
		}
	}

	writeOverview(&b, doc)
	for _, c := range doc.Categories {
		writeCategory(&b, c)
	}
	if opts.GroupByPhase {
		writePhases(&b, doc)
	}
	if opts.IncludeRootCauses {
		writeRootCauses(&b, doc)
	}
	if opts.IncludeInitiatives {
		writeInitiatives(&b, doc)
	}
	if opts.ShowCapabilityMetrics {
		writeCapability(&b, doc)
	}
	return []byte(mdshift.Shift(strings.TrimRight(b.String(), "\n")+"\n", opts.HeadingOffset)), nil
}

// Title returns the document's name, or "DMAIC Metrics Report".
func Title(doc *dmaic.DMAICDocument) string {
	if doc.Metadata != nil && doc.Metadata.Name != "" {
		return doc.Metadata.Name
	}
	return "DMAIC Metrics Report"
}

func metadataLine(m *dmaic.Metadata) string {
	if m == nil {
		return ""
	}
	var parts []string
	for _, f := range []struct{ label, value string }{
		{"Team", m.Team}, {"Period", m.Period}, {"Owner", m.Owner}, {"Status", m.Status},
		{"Review cadence", m.ReviewCadence}, {"Version", m.Version},
	} {
		if f.value != "" {
			parts = append(parts, fmt.Sprintf("**%s:** %s", f.label, f.value))
		}
	}
	return strings.Join(parts, " · ")
}

func writeOverview(b *strings.Builder, doc *dmaic.DMAICDocument) {
	b.WriteString("## Overview\n\n")
	if doc.Metadata != nil && doc.Metadata.Description != "" {
		fmt.Fprintf(b, "%s\n\n", oneLine(doc.Metadata.Description))
	}
	fmt.Fprintf(b, "**Overall health:** %s · **Categories:** %d · **Metrics:** %d\n\n",
		Health(doc.CalculateOverallHealth()), len(doc.Categories), len(doc.AllMetrics()))
	if len(doc.Categories) == 0 {
		return
	}
	b.WriteString("| Category | Health | Metrics | Green | Yellow | Red |\n")
	b.WriteString("|:---------|-------:|--------:|------:|-------:|----:|\n")
	for _, c := range doc.Categories {
		counts := map[string]int{}
		for _, m := range c.Metrics {
			counts[m.CalculateStatus()]++
		}
		fmt.Fprintf(b, "| %s | %s | %d | %d | %d | %d |\n", cell(c.Name), Health(c.CalculateCategoryHealth()),
			len(c.Metrics), counts[dmaic.StatusGreen], counts[dmaic.StatusYellow], counts[dmaic.StatusRed])
	}
	b.WriteString("\n")
}

func writeCategory(b *strings.Builder, c dmaic.Category) {
	fmt.Fprintf(b, "## %s\n\n", oneLine(c.Name))
	if c.Description != "" {
		fmt.Fprintf(b, "%s\n\n", oneLine(c.Description))
	}
	var facts []string
	if c.Owner != "" {
		facts = append(facts, "**Owner:** "+c.Owner)
	}
	facts = append(facts, "**Health:** "+Health(c.CalculateCategoryHealth()))
	fmt.Fprintf(b, "%s\n\n", strings.Join(facts, " · "))
	if len(c.Metrics) == 0 {
		return
	}
	b.WriteString("| Metric | Phase | Baseline | Current | Target | Unit | Status |\n")
	b.WriteString("|:-------|:-----:|---------:|--------:|-------:|:-----|:------:|\n")
	for _, m := range c.Metrics {
		unit := m.Unit
		if unit == "" {
			unit = "-"
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s | %s |\n", cell(m.Name), cell(m.Phase),
			Value(m.Baseline, m.Unit), Value(m.Current, m.Unit), Value(m.Target, m.Unit), cell(unit),
			StatusText(m.CalculateStatus()))
	}
	b.WriteString("\n")
	for _, m := range c.Metrics {
		if status := m.CalculateStatus(); status == dmaic.StatusYellow || status == dmaic.StatusRed {
			writeMetricDetail(b, m, status)
		}
	}
}

// writeMetricDetail details a metric that needs attention: its values,
// thresholds, control limits, and root causes.
func writeMetricDetail(b *strings.Builder, m dmaic.Metric, status string) {
	fmt.Fprintf(b, "### %s (%s)\n\n", oneLine(m.Name), status)
	if m.Description != "" {
		fmt.Fprintf(b, "%s\n\n", oneLine(m.Description))
	}
	b.WriteString("| Measure | Value |\n")
	b.WriteString("|:--------|------:|\n")
	fmt.Fprintf(b, "| Baseline | %s |\n", Value(m.Baseline, m.Unit))
	fmt.Fprintf(b, "| Current | %s |\n", Value(m.Current, m.Unit))
	fmt.Fprintf(b, "| Target | %s |\n", Value(m.Target, m.Unit))
	fmt.Fprintf(b, "| Trend direction | %s |\n", cell(TrendText(m.TrendDirection)))
	if t := m.Thresholds; t != nil {
		fmt.Fprintf(b, "| Warning threshold | %s |\n", Value(t.Warning, m.Unit))
		fmt.Fprintf(b, "| Critical threshold | %s |\n", Value(t.Critical, m.Unit))
	}
	b.WriteString("\n")
	if cl := m.ControlLimits; cl != nil {
		b.WriteString("**Control limits**\n\n")
		b.WriteString("| Limit | Value |\n")
		b.WriteString("|:------|------:|\n")
		fmt.Fprintf(b, "| Upper control limit (UCL) | %s |\n", Value(cl.UCL, m.Unit))
		fmt.Fprintf(b, "| Center line | %s |\n", Value(cl.CenterLine, m.Unit))
		fmt.Fprintf(b, "| Lower control limit (LCL) | %s |\n", Value(cl.LCL, m.Unit))
		if cl.Sigma != 0 {
			fmt.Fprintf(b, "| Sigma | %s |\n", Value(cl.Sigma, ""))
		}
		b.WriteString("\n")
	}
	if len(m.RootCauses) > 0 {
		b.WriteString("**Root causes**\n\n")
		b.WriteString("| Root cause | Category | Impact | Validated |\n")
		b.WriteString("|:-----------|:---------|:-------|:---------:|\n")
		for _, rc := range m.RootCauses {
			fmt.Fprintf(b, "| %s | %s | %s | %s |\n", cell(rc.Description), cell(rc.Category), cell(rc.Impact), yesNo(rc.Validated))
		}
		b.WriteString("\n")
	}
}

func writePhases(b *strings.Builder, doc *dmaic.DMAICDocument) {
	byPhase := doc.MetricsByPhase()
	b.WriteString("## DMAIC Phase Distribution\n\n")
	b.WriteString("| Phase | Description | Metrics |\n")
	b.WriteString("|:------|:------------|--------:|\n")
	for _, p := range dmaic.Phases() {
		fmt.Fprintf(b, "| **%s** - %s | %s | %d |\n", PhaseAbbrev(p), p, dmaic.PhaseDescription(p), len(byPhase[p]))
	}
	b.WriteString("\n")
}

func writeRootCauses(b *strings.Builder, doc *dmaic.DMAICDocument) {
	var rows []string
	for _, c := range doc.Categories {
		for _, m := range c.Metrics {
			for _, rc := range m.RootCauses {
				rows = append(rows, fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
					cell(m.Name), cell(rc.Description), cell(rc.Category), cell(rc.Impact), yesNo(rc.Validated)))
			}
		}
	}
	if len(rows) == 0 {
		return
	}
	b.WriteString("## Root Cause Analysis\n\n")
	b.WriteString("| Metric | Root cause | Category | Impact | Validated |\n")
	b.WriteString("|:-------|:-----------|:---------|:-------|:---------:|\n")
	for _, row := range rows {
		b.WriteString(row)
	}
	b.WriteString("\n")
}

func writeInitiatives(b *strings.Builder, doc *dmaic.DMAICDocument) {
	if len(doc.Initiatives) == 0 {
		return
	}
	names := map[string]string{}
	for _, m := range doc.AllMetrics() {
		if m.ID != "" {
			names[m.ID] = m.Name
		}
	}
	b.WriteString("## Improvement Initiatives\n\n")
	for _, in := range doc.Initiatives {
		fmt.Fprintf(b, "### %s\n\n", oneLine(in.Name))
		if in.Description != "" {
			fmt.Fprintf(b, "%s\n\n", oneLine(in.Description))
		}
		status := in.Status
		if status == "" {
			status = "N/A"
		}
		b.WriteString("| Field | Value |\n")
		b.WriteString("|:------|:------|\n")
		row := func(label, value string) {
			if value != "" {
				fmt.Fprintf(b, "| **%s** | %s |\n", label, cell(value))
			}
		}
		row("Owner", in.Owner)
		row("Status", status)
		row("Start date", in.StartDate)
		row("End date", in.EndDate)
		row("Expected impact", in.ExpectedImpact)
		row("Actual impact", in.ActualImpact)
		b.WriteString("\n")
		if len(in.MetricIDs) > 0 {
			var linked []string
			for _, id := range in.MetricIDs {
				if name, ok := names[id]; ok {
					linked = append(linked, oneLine(name))
				} else {
					linked = append(linked, "`"+id+"`")
				}
			}
			fmt.Fprintf(b, "**Linked metrics:** %s\n\n", strings.Join(linked, ", "))
		}
	}
}

func writeCapability(b *strings.Builder, doc *dmaic.DMAICDocument) {
	b.WriteString("## Appendix: Process Capability Metrics\n\n")
	b.WriteString("| Metric | Cp | Cpk | Sigma level | DPMO |\n")
	b.WriteString("|:-------|---:|----:|------------:|-----:|\n")
	for _, m := range doc.AllMetrics() {
		if pc := m.ProcessCapability; pc != nil {
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", cell(m.Name),
				Value(pc.Cp, ""), Value(pc.Cpk, ""), Value(pc.SigmaLevel, ""), Value(pc.DPMO, ""))
		}
	}
	b.WriteString("\n")
	b.WriteString("### Control Limits Summary\n\n")
	b.WriteString("| Metric | UCL | Center | LCL | In control |\n")
	b.WriteString("|:-------|----:|-------:|----:|:----------:|\n")
	for _, m := range doc.AllMetrics() {
		if cl := m.ControlLimits; cl != nil {
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", cell(m.Name),
				Value(cl.UCL, m.Unit), Value(cl.CenterLine, m.Unit), Value(cl.LCL, m.Unit), yesNo(m.IsInControl()))
		}
	}
	b.WriteString("\n")
}

// Health formats a health score between 0 and 1 as a whole percentage.
func Health(h float64) string {
	return fmt.Sprintf("%.0f%%", h*100)
}

// Value formats a metric value: one decimal place and a "%" sign for a
// percentage, no decimals for a whole number, and two decimals otherwise.
func Value(v float64, unit string) string {
	switch {
	case unit == "%":
		return fmt.Sprintf("%.1f%%", v)
	case v == float64(int64(v)):
		return fmt.Sprintf("%.0f", v)
	default:
		return fmt.Sprintf("%.2f", v)
	}
}

// StatusText names a status, or "Unknown".
func StatusText(status string) string {
	switch status {
	case dmaic.StatusGreen, dmaic.StatusYellow, dmaic.StatusRed:
		return status
	default:
		return "Unknown"
	}
}

// PhaseAbbrev returns a phase's initial letter, or "?".
func PhaseAbbrev(phase string) string {
	switch phase {
	case dmaic.PhaseDefine:
		return "D"
	case dmaic.PhaseMeasure:
		return "M"
	case dmaic.PhaseAnalyze:
		return "A"
	case dmaic.PhaseImprove:
		return "I"
	case dmaic.PhaseControl:
		return "C"
	default:
		return "?"
	}
}

// TrendText describes a trend direction, or "N/A" when there is none.
func TrendText(trend string) string {
	switch trend {
	case dmaic.TrendHigherBetter:
		return "Higher is better"
	case dmaic.TrendLowerBetter:
		return "Lower is better"
	case dmaic.TrendTargetValue:
		return "Closest to target"
	case "":
		return "N/A"
	default:
		return trend
	}
}

func yesNo(v bool) string {
	if v {
		return "Yes"
	}
	return "No"
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// cell makes text safe for a Markdown table cell.
func cell(s string) string {
	return strings.ReplaceAll(oneLine(s), "|", `\|`)
}
