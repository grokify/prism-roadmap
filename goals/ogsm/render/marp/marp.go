// Package marp renders an OGSM as Marp slides: a title slide, the plan on one
// page, the objective with its goals, and a slide per strategy with its
// measures and actions. See https://marp.app/.
package marp

import (
	"fmt"
	"strings"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-roadmap/goals/ogsm"
	"github.com/grokify/prism-roadmap/goals/ogsm/render"
	"github.com/grokify/prism-roadmap/goals/ogsm/render/markdown"
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

// Render converts the OGSM to Marp Markdown.
func (r *Renderer) Render(o *ogsm.OGSM, opts *render.Options) ([]byte, error) {
	if o == nil {
		return nil, fmt.Errorf("marp: nil OGSM")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	theme := sdmarp.GetTheme(opts.Theme)
	title := markdown.Title(o)
	var b strings.Builder

	b.WriteString("---\nmarp: true\n")
	fmt.Fprintf(&b, "theme: %s\npaginate: true\n", theme.Name)
	if m := o.Metadata; m != nil && m.Team != "" {
		header := m.Team
		if m.Period != "" {
			header += " | " + m.Period
		}
		fmt.Fprintf(&b, "header: %q\n", header)
	}
	fmt.Fprintf(&b, "footer: %q\n", "OGSM | "+title)
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
    font-size: 0.75em;
    width: 100%%;
  }
  blockquote {
    border-left: 4px solid %[2]s;
    padding-left: 1em;
    font-style: italic;
  }
---

`, theme.PrimaryBgColor, theme.AccentColor, theme.PrimaryTextColor)

	fmt.Fprintf(&b, "<!-- _class: title -->\n\n# %s\n\n**Objective · Goals · Strategies · Measures**\n", title)
	if m := o.Metadata; m != nil {
		for _, f := range []struct{ label, value string }{
			{"Team", m.Team}, {"Period", m.Period}, {"Author", m.Author}, {"Status", m.Status},
		} {
			if f.value != "" {
				fmt.Fprintf(&b, "\n**%s:** %s\n", f.label, f.value)
			}
		}
	}
	b.WriteString("\n---\n\n")

	b.WriteString("## The Plan on One Page\n\n")
	b.WriteString(markdown.OnePage(o))
	b.WriteString("\n---\n\n")

	b.WriteString("## Objective and Goals\n\n")
	fmt.Fprintf(&b, "> %s\n\n", strings.Join(strings.Fields(o.Objective.Statement), " "))
	if o.Objective.DueDate != "" {
		fmt.Fprintf(&b, "**By:** %s\n\n", o.Objective.DueDate)
	}
	b.WriteString("| Goal | Owner | Baseline | Target | Current | Status |\n|------|-------|----------|--------|---------|--------|\n")
	for _, g := range o.Goals {
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n", cell(g.Name), cell(g.Owner),
			cell(markdown.Value(g.Baseline, g.Unit)), cell(markdown.Value(g.Target, g.Unit)),
			cell(markdown.Value(g.Current, g.Unit)), cell(g.Status))
	}
	b.WriteString("\n---\n\n")

	for i, s := range o.Strategies {
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
			fmt.Fprintf(&b, "%s\n\n", strings.Join(strings.Fields(s.Description), " "))
		}
		if measures := o.MeasuresFor(s.ID); len(measures) > 0 {
			b.WriteString("| Measure | Owner | Baseline | Target | Current | Status |\n|---------|-------|----------|--------|---------|--------|\n")
			for _, m := range measures {
				fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s |\n", cell(m.Name), cell(m.Owner),
					cell(markdown.Value(m.Baseline, m.Unit)), cell(markdown.Value(m.Target, m.Unit)),
					cell(markdown.Value(m.Current, m.Unit)), cell(m.Status))
			}
			b.WriteString("\n")
		}
		if actions := o.ActionsFor(s.ID); opts.IncludeActions && len(actions) > 0 {
			b.WriteString("| Action | Owner | Timing | Status |\n|--------|-------|--------|--------|\n")
			for _, a := range actions {
				fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", cell(a.Name), cell(a.Owner),
					cell(markdown.Timing(a.StartDate, a.EndDate)), cell(a.Status))
			}
			b.WriteString("\n")
		}
		b.WriteString("---\n\n")
	}
	return []byte(strings.TrimSuffix(b.String(), "---\n\n")), nil
}

func cell(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "|", `\|`)
}
