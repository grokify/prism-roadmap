// Package markdown renders a VSEM as a Markdown document: the vision, then
// each strategy with its resources, execution, and metrics, then the metrics
// that measure the vision as a whole. The output is deterministic: the same
// VSEM always renders to the same bytes.
package markdown

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/vsem"
	"github.com/grokify/prism-roadmap/goals/vsem/render"
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

// Render converts the VSEM to Markdown.
func (r *Renderer) Render(v *vsem.VSEM, opts *render.Options) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("markdown: nil VSEM")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	var b strings.Builder

	title := "VSEM"
	if v.Metadata != nil && v.Metadata.Name != "" {
		title = v.Metadata.Name
	}
	fmt.Fprintf(&b, "# %s\n\n", title)
	if line := metadataLine(v.Metadata); line != "" {
		fmt.Fprintf(&b, "%s\n\n", line)
	}

	b.WriteString("## Vision\n\n")
	if v.Vision != "" {
		fmt.Fprintf(&b, "%s\n\n", v.Vision)
	}
	if v.Mission != "" {
		fmt.Fprintf(&b, "**Mission:** %s\n\n", v.Mission)
	}

	b.WriteString("## Strategy\n\n")
	for i, s := range v.Strategies {
		fmt.Fprintf(&b, "### %d. %s\n\n", i+1, s.Name)
		if facts := strategyFacts(s); facts != "" {
			fmt.Fprintf(&b, "%s\n\n", facts)
		}
		if opts.IncludeDescriptions && s.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", s.Description)
		}
		if opts.IncludeResources && len(s.Resources) > 0 {
			b.WriteString("**Resources**\n\n")
			for _, res := range s.Resources {
				line := res.Name
				if res.Allocation != "" {
					line += ": " + res.Allocation
				}
				if opts.IncludeDescriptions && res.Description != "" {
					line += " (" + res.Description + ")"
				}
				fmt.Fprintf(&b, "- %s\n", line)
			}
			b.WriteString("\n")
		}
		if items := v.ExecutionFor(s.ID); len(items) > 0 {
			b.WriteString("**Execution**\n\n")
			writeExecution(&b, items, opts)
		}
		if metrics := v.MetricsFor(s.ID); len(metrics) > 0 {
			b.WriteString("**Metrics**\n\n")
			writeMetrics(&b, metrics, opts)
		}
	}

	if metrics := v.VisionMetrics(); len(metrics) > 0 {
		b.WriteString("## Metrics for the Vision\n\n")
		writeMetrics(&b, metrics, opts)
	}

	return []byte(strings.TrimRight(b.String(), "\n") + "\n"), nil
}

func metadataLine(m *vsem.Metadata) string {
	if m == nil {
		return ""
	}
	var parts []string
	for _, p := range []struct{ label, value string }{
		{"Team", m.Team},
		{"Period", m.Period},
		{"Author", m.Author},
		{"Status", m.Status},
		{"Version", m.Version},
	} {
		if p.value != "" {
			parts = append(parts, fmt.Sprintf("**%s:** %s", p.label, p.value))
		}
	}
	return strings.Join(parts, " · ")
}

func strategyFacts(s vsem.Strategy) string {
	var parts []string
	if s.Owner != "" {
		parts = append(parts, "**Owner:** "+s.Owner)
	}
	if s.Priority != "" {
		parts = append(parts, "**Priority:** "+s.Priority)
	}
	return strings.Join(parts, " · ")
}

func writeExecution(b *strings.Builder, items []vsem.Execution, opts *render.Options) {
	b.WriteString("| Work | Kind | Owner | Timing | Status |\n")
	b.WriteString("|------|------|-------|--------|--------|\n")
	for _, e := range items {
		name := cell(e.Name)
		if opts.IncludeDescriptions && e.Description != "" {
			name += "<br>" + cell(e.Description)
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n",
			name, cell(e.Kind), cell(e.Owner), cell(Timing(e.StartDate, e.EndDate)), cell(e.Status))
	}
	b.WriteString("\n")
}

func writeMetrics(b *strings.Builder, metrics []vsem.Metric, opts *render.Options) {
	b.WriteString("| Metric | Accountable | Baseline | Target | Current | Due | Status |\n")
	b.WriteString("|--------|-------------|----------|--------|---------|-----|--------|\n")
	for _, m := range metrics {
		name := cell(m.Name)
		if opts.IncludeDescriptions && m.Description != "" {
			name += "<br>" + cell(m.Description)
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s | %s |\n",
			name, cell(m.Owner), cell(Value(m.Baseline, m.Unit)), cell(Value(m.Target, m.Unit)),
			cell(Value(m.Current, m.Unit)), cell(m.DueDate), cell(m.Status))
	}
	b.WriteString("\n")
}

// Timing formats a start and end date as a range.
func Timing(start, end string) string {
	switch {
	case start != "" && end != "":
		return start + " to " + end
	case start != "":
		return "from " + start
	case end != "":
		return "by " + end
	default:
		return ""
	}
}

// Value appends the unit to a value, without a space for "%".
func Value(v, unit string) string {
	switch {
	case v == "" || unit == "":
		return v
	case unit == "%":
		return v + unit
	default:
		return v + " " + unit
	}
}

// cell makes text safe for a Markdown table cell.
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.Join(strings.Fields(s), " ")
}
