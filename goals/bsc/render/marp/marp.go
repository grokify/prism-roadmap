// Package marp renders a Balanced Scorecard as Marp slides: a title slide,
// the mission, vision, and strategy, the strategy map as a table, a slide per
// perspective with its scorecard, and the initiatives. See https://marp.app/.
package marp

import (
	"fmt"
	"strings"

	sdmarp "github.com/grokify/structureddocs/marp"

	"github.com/grokify/prism-roadmap/goals/bsc"
	"github.com/grokify/prism-roadmap/goals/bsc/render"
	"github.com/grokify/prism-roadmap/goals/bsc/render/markdown"
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

// Render converts the scorecard to Marp Markdown.
func (r *Renderer) Render(b *bsc.BSC, opts *render.Options) ([]byte, error) {
	if b == nil {
		return nil, fmt.Errorf("marp: nil Balanced Scorecard")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	theme := sdmarp.GetTheme(opts.Theme)
	title := markdown.Title(b)
	var s strings.Builder

	s.WriteString("---\nmarp: true\n")
	fmt.Fprintf(&s, "theme: %s\npaginate: true\n", theme.Name)
	if m := b.Metadata; m != nil && m.Team != "" {
		header := m.Team
		if m.Period != "" {
			header += " | " + m.Period
		}
		fmt.Fprintf(&s, "header: %q\n", header)
	}
	fmt.Fprintf(&s, "footer: %q\n", "Balanced Scorecard | "+title)
	fmt.Fprintf(&s, `style: |
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
    font-size: 0.7em;
    width: 100%%;
  }
  blockquote {
    border-left: 4px solid %[2]s;
    padding-left: 1em;
    font-style: italic;
  }
---

`, theme.PrimaryBgColor, theme.AccentColor, theme.PrimaryTextColor)

	fmt.Fprintf(&s, "<!-- _class: title -->\n\n# %s\n\n**Balanced Scorecard**\n", title)
	if m := b.Metadata; m != nil {
		for _, f := range []struct{ label, value string }{
			{"Team", m.Team}, {"Period", m.Period}, {"Author", m.Author}, {"Status", m.Status},
		} {
			if f.value != "" {
				fmt.Fprintf(&s, "\n**%s:** %s\n", f.label, f.value)
			}
		}
	}
	s.WriteString("\n---\n\n")

	if purpose := markdown.Purpose(b); len(purpose) > 0 {
		s.WriteString("## Mission, Vision, and Strategy\n\n")
		for _, p := range purpose {
			fmt.Fprintf(&s, "**%s**\n\n> %s\n\n", p.Label, p.Text)
		}
		s.WriteString("---\n\n")
	}

	if opts.IncludeStrategyMap && len(b.Objectives) > 0 {
		s.WriteString("## Strategy Map\n\n")
		s.WriteString("| Perspective | Objectives | Contribute to |\n|-------------|------------|---------------|\n")
		for _, p := range b.Perspectives {
			objectives := b.ObjectivesIn(p.ID)
			if len(objectives) == 0 {
				continue
			}
			var names, targets []string
			for _, o := range objectives {
				names = append(names, cell(o.Name))
				targets = append(targets, cell(markdown.ObjectiveNamesByID(b, o.ContributesTo)))
			}
			fmt.Fprintf(&s, "| **%s** | %s | %s |\n", cell(p.Name), strings.Join(names, "<br>"), strings.Join(targets, "<br>"))
		}
		s.WriteString("\n---\n\n")
	}

	for _, p := range b.Perspectives {
		objectives := b.ObjectivesIn(p.ID)
		if len(objectives) == 0 {
			continue
		}
		fmt.Fprintf(&s, "## %s\n\n", p.Name)
		if p.Question != "" {
			fmt.Fprintf(&s, "*%s*\n\n", strings.Join(strings.Fields(p.Question), " "))
		}
		s.WriteString("| Objective | Measure | Type | Target | Current | Status |\n|-----------|---------|------|--------|---------|--------|\n")
		for _, o := range objectives {
			name := cell(o.Name)
			measures := b.MeasuresFor(o.ID)
			if len(measures) == 0 {
				fmt.Fprintf(&s, "| %s |  |  |  |  |  |\n", name)
				continue
			}
			for i, m := range measures {
				if i > 0 {
					name = ""
				}
				fmt.Fprintf(&s, "| %s | %s | %s | %s | %s | %s |\n", name, cell(m.Name), cell(m.Type),
					cell(markdown.Value(m.Target, m.Unit)), cell(markdown.Value(m.Current, m.Unit)), cell(m.Status))
			}
		}
		s.WriteString("\n---\n\n")
	}

	if opts.IncludeInitiatives && len(b.Initiatives) > 0 {
		s.WriteString("## Initiatives\n\n")
		s.WriteString("| Initiative | Objectives | Owner | Timing | Status |\n|------------|------------|-------|--------|--------|\n")
		for _, in := range b.Initiatives {
			fmt.Fprintf(&s, "| %s | %s | %s | %s | %s |\n", cell(in.Name), cell(markdown.ObjectiveNamesByID(b, in.ObjectiveIDs)),
				cell(in.Owner), cell(markdown.Timing(in.StartDate, in.EndDate)), cell(in.Status))
		}
		s.WriteString("\n---\n\n")
	}
	return []byte(strings.TrimSuffix(s.String(), "---\n\n")), nil
}

func cell(s string) string {
	return strings.ReplaceAll(strings.Join(strings.Fields(s), " "), "|", `\|`)
}
