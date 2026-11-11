// Package marp renders a VSEM as Marp slides: a title slide, a one-page
// overview of the four sections, a slide per strategy with its resources,
// execution, and metrics, and a slide for the metrics that measure the vision
// as a whole. See https://marp.app/.
package marp

import (
	"fmt"
	"strings"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-roadmap/goals/vsem"
	"github.com/grokify/prism-roadmap/goals/vsem/render"
	"github.com/grokify/prism-roadmap/goals/vsem/render/markdown"
)

// Renderer implements render.Renderer for Marp slides.
type Renderer struct{}

// New returns a Marp renderer.
func New() *Renderer {
	return &Renderer{}
}

// Format returns "marp".
func (r *Renderer) Format() string {
	return "marp"
}

// FileExtension returns ".md".
func (r *Renderer) FileExtension() string {
	return ".md"
}

// Render converts the VSEM to Marp Markdown.
func (r *Renderer) Render(v *vsem.VSEM, opts *render.Options) ([]byte, error) {
	if v == nil {
		return nil, fmt.Errorf("marp: nil VSEM")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	theme := sdmarp.GetTheme(opts.Theme)
	var b strings.Builder

	title := "VSEM"
	if v.Metadata != nil && v.Metadata.Name != "" {
		title = v.Metadata.Name
	}

	// Front matter.
	b.WriteString("---\nmarp: true\n")
	fmt.Fprintf(&b, "theme: %s\npaginate: true\n", theme.Name)
	if v.Metadata != nil && v.Metadata.Team != "" {
		header := v.Metadata.Team
		if v.Metadata.Period != "" {
			header += " | " + v.Metadata.Period
		}
		fmt.Fprintf(&b, "header: %q\n", header)
	}
	fmt.Fprintf(&b, "footer: %q\n", "VSEM | "+title)
	fmt.Fprintf(&b, `style: |
  section {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
  section.title {
    text-align: center;
    background: linear-gradient(135deg, %[1]s 0%%, %[2]s 100%%);
    color: %[3]s;
  }
  section.title h1 {
    color: %[3]s;
  }
  table {
    font-size: 0.8em;
    width: 100%%;
  }
  blockquote {
    border-left: 4px solid %[2]s;
    padding-left: 1em;
    font-style: italic;
  }
---

`, theme.PrimaryBgColor, theme.AccentColor, theme.PrimaryTextColor)

	// Title slide.
	fmt.Fprintf(&b, "<!-- _class: title -->\n\n# %s\n\n**Vision · Strategy · Execution · Metrics**\n", title)
	if m := v.Metadata; m != nil {
		for _, p := range []struct{ label, value string }{
			{"Team", m.Team}, {"Period", m.Period}, {"Author", m.Author}, {"Status", m.Status},
		} {
			if p.value != "" {
				fmt.Fprintf(&b, "\n**%s:** %s\n", p.label, p.value)
			}
		}
	}
	b.WriteString("\n---\n\n")

	// One-page overview.
	b.WriteString("## The Plan on One Page\n\n")
	if s := v.Statement(); s != "" {
		fmt.Fprintf(&b, "> %s\n\n", s)
	}
	b.WriteString("| Strategy | Execution | Metrics |\n|----------|-----------|---------|\n")
	for _, s := range v.Strategies {
		fmt.Fprintf(&b, "| %s | %s | %s |\n", cell(s.Name), names(executionNames(v.ExecutionFor(s.ID))), names(metricNames(v.MetricsFor(s.ID))))
	}
	b.WriteString("\n---\n\n")

	// Vision and mission.
	if v.Vision != "" || v.Mission != "" {
		b.WriteString("## Vision\n\n")
		if v.Vision != "" {
			fmt.Fprintf(&b, "> %s\n\n", v.Vision)
		}
		if v.Mission != "" {
			fmt.Fprintf(&b, "**Mission:** %s\n\n", v.Mission)
		}
		b.WriteString("---\n\n")
	}

	// One slide per strategy.
	for i, s := range v.Strategies {
		fmt.Fprintf(&b, "## Strategy %d: %s\n\n", i+1, s.Name)
		var facts []string
		if s.Owner != "" {
			facts = append(facts, "**Owner:** "+s.Owner)
		}
		if s.Priority != "" {
			facts = append(facts, "**Priority:** "+s.Priority)
		}
		if len(facts) > 0 {
			fmt.Fprintf(&b, "%s\n\n", strings.Join(facts, " · "))
		}
		if opts.IncludeDescriptions && s.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", s.Description)
		}
		if opts.IncludeResources && len(s.Resources) > 0 {
			var res []string
			for _, r := range s.Resources {
				if r.Allocation != "" {
					res = append(res, r.Name+": "+r.Allocation)
				} else {
					res = append(res, r.Name)
				}
			}
			fmt.Fprintf(&b, "**Resources:** %s\n\n", strings.Join(res, "; "))
		}
		if items := v.ExecutionFor(s.ID); len(items) > 0 {
			b.WriteString("| Execution | Owner | Timing | Status |\n|-----------|-------|--------|--------|\n")
			for _, e := range items {
				fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", cell(e.Name), cell(e.Owner), cell(markdown.Timing(e.StartDate, e.EndDate)), cell(e.Status))
			}
			b.WriteString("\n")
		}
		if metrics := v.MetricsFor(s.ID); len(metrics) > 0 {
			writeMetrics(&b, metrics)
		}
		b.WriteString("---\n\n")
	}

	// Metrics for the vision as a whole.
	if metrics := v.VisionMetrics(); len(metrics) > 0 {
		b.WriteString("## Metrics for the Vision\n\n")
		writeMetrics(&b, metrics)
		b.WriteString("---\n\n")
	}

	return []byte(strings.TrimSuffix(b.String(), "---\n\n")), nil
}

func writeMetrics(b *strings.Builder, metrics []vsem.Metric) {
	b.WriteString("| Metric | Accountable | Baseline | Target | Current | Status |\n|--------|-------------|----------|--------|---------|--------|\n")
	for _, m := range metrics {
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s |\n", cell(m.Name), cell(m.Owner),
			cell(markdown.Value(m.Baseline, m.Unit)), cell(markdown.Value(m.Target, m.Unit)),
			cell(markdown.Value(m.Current, m.Unit)), cell(m.Status))
	}
	b.WriteString("\n")
}

func executionNames(items []vsem.Execution) []string {
	out := make([]string, 0, len(items))
	for _, e := range items {
		out = append(out, e.Name)
	}
	return out
}

func metricNames(metrics []vsem.Metric) []string {
	out := make([]string, 0, len(metrics))
	for _, m := range metrics {
		out = append(out, m.Name)
	}
	return out
}

func names(list []string) string {
	cells := make([]string, 0, len(list))
	for _, s := range list {
		cells = append(cells, cell(s))
	}
	return strings.Join(cells, "<br>")
}

func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.Join(strings.Fields(s), " ")
}
