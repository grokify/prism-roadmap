// Package markdown renders a Balanced Scorecard as a Markdown document: the
// mission, vision, and strategy; the strategy map as a Mermaid diagram; the
// strategic themes; the scorecard, one table per perspective; and the
// initiatives. The output is deterministic: the same scorecard always renders
// to the same bytes.
package markdown

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/bsc"
	"github.com/grokify/prism-roadmap/goals/bsc/render"
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

// Render converts the scorecard to Markdown.
func (r *Renderer) Render(b *bsc.BSC, opts *render.Options) ([]byte, error) {
	if b == nil {
		return nil, fmt.Errorf("markdown: nil Balanced Scorecard")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	var s strings.Builder
	if !opts.OmitHeader {
		fmt.Fprintf(&s, "# %s\n\n", Title(b))
		if line := metadataLine(b.Metadata); line != "" {
			fmt.Fprintf(&s, "%s\n\n", line)
		}
	}

	if purpose := Purpose(b); len(purpose) > 0 {
		s.WriteString("## Mission, Vision, and Strategy\n\n")
		for _, p := range purpose {
			fmt.Fprintf(&s, "**%s:** %s\n\n", p.Label, p.Text)
		}
	}

	if opts.IncludeStrategyMap && len(b.Objectives) > 0 {
		s.WriteString("## Strategy Map\n\n")
		s.WriteString("Each arrow runs from an objective to one it helps achieve.\n\n")
		s.WriteString(StrategyMap(b))
		s.WriteString("\n")
	}

	if len(b.Themes) > 0 {
		s.WriteString("## Strategic Themes\n\n")
		s.WriteString("| Theme | Owner | Objectives |\n|-------|-------|------------|\n")
		for _, t := range b.Themes {
			name := cell(t.Name)
			if opts.IncludeDescriptions && t.Description != "" {
				name += "<br>" + cell(t.Description)
			}
			fmt.Fprintf(&s, "| %s | %s | %s |\n", name, cell(t.Owner), cell(objectiveNames(b.ObjectivesInTheme(t.ID))))
		}
		s.WriteString("\n")
	}

	s.WriteString("## Scorecard\n\n")
	for _, p := range b.Perspectives {
		fmt.Fprintf(&s, "### %s\n\n", p.Name)
		if p.Question != "" {
			fmt.Fprintf(&s, "*%s*\n\n", oneLine(p.Question))
		}
		if opts.IncludeDescriptions && p.Description != "" {
			fmt.Fprintf(&s, "%s\n\n", oneLine(p.Description))
		}
		objectives := b.ObjectivesIn(p.ID)
		if len(objectives) == 0 {
			s.WriteString("No objectives.\n\n")
			continue
		}
		s.WriteString("| Objective | Measure | Type | Owner | Baseline | Target | Current | Status |\n")
		s.WriteString("|-----------|---------|------|-------|----------|--------|---------|--------|\n")
		for _, o := range objectives {
			name := cell(o.Name)
			if o.Owner != "" {
				name += "<br>Owner: " + cell(o.Owner)
			}
			if opts.IncludeDescriptions && o.Description != "" {
				name += "<br>" + cell(o.Description)
			}
			measures := b.MeasuresFor(o.ID)
			if len(measures) == 0 {
				fmt.Fprintf(&s, "| %s |  |  |  |  |  |  |  |\n", name)
				continue
			}
			for i, m := range measures {
				if i > 0 {
					name = ""
				}
				fmt.Fprintf(&s, "| %s | %s | %s | %s | %s | %s | %s | %s |\n", name, cell(m.Name), cell(m.Type), cell(m.Owner),
					cell(Value(m.Baseline, m.Unit)), cell(Target(m)), cell(Value(m.Current, m.Unit)), cell(m.Status))
			}
		}
		s.WriteString("\n")
	}

	if opts.IncludeInitiatives && len(b.Initiatives) > 0 {
		s.WriteString("## Initiatives\n\n")
		s.WriteString("| Initiative | Objectives | Owner | Timing | Budget | Status |\n")
		s.WriteString("|------------|------------|-------|--------|--------|--------|\n")
		for _, in := range b.Initiatives {
			name := cell(in.Name)
			if opts.IncludeDescriptions && in.Description != "" {
				name += "<br>" + cell(in.Description)
			}
			fmt.Fprintf(&s, "| %s | %s | %s | %s | %s | %s |\n", name, cell(ObjectiveNamesByID(b, in.ObjectiveIDs)),
				cell(in.Owner), cell(Timing(in.StartDate, in.EndDate)), cell(in.Budget), cell(in.Status))
		}
		s.WriteString("\n")
	}
	return []byte(mdshift.Shift(strings.TrimRight(s.String(), "\n")+"\n", opts.HeadingOffset)), nil
}

// StrategyMap renders the strategy map as a Mermaid flowchart in a fenced
// block: one box per perspective, top first, and an arrow from each objective
// to each one it contributes to. Node IDs are positional, so any objective ID
// is safe.
func StrategyMap(b *bsc.BSC) string {
	node := map[string]string{}
	for i, o := range b.Objectives {
		node[o.ID] = fmt.Sprintf("o%d", i+1)
	}
	var s strings.Builder
	s.WriteString("```mermaid\nflowchart BT\n")
	for i, p := range b.Perspectives {
		objectives := b.ObjectivesIn(p.ID)
		if len(objectives) == 0 {
			continue
		}
		fmt.Fprintf(&s, "  subgraph p%d[\"%s\"]\n", i+1, label(p.Name))
		for _, o := range objectives {
			fmt.Fprintf(&s, "    %s[\"%s\"]\n", node[o.ID], label(o.Name))
		}
		s.WriteString("  end\n")
	}
	for _, o := range b.Objectives {
		for _, to := range o.ContributesTo {
			if target, ok := node[to]; ok && to != o.ID {
				fmt.Fprintf(&s, "  %s --> %s\n", node[o.ID], target)
			}
		}
	}
	s.WriteString("```\n")
	return s.String()
}

// Statement is one labelled statement of purpose.
type Statement struct {
	Label, Text string
}

// Purpose returns the mission, vision, and strategy that are set, in that
// order.
func Purpose(b *bsc.BSC) []Statement {
	var out []Statement
	for _, p := range []Statement{{"Mission", b.Mission}, {"Vision", b.Vision}, {"Strategy", b.Strategy}} {
		if t := oneLine(p.Text); t != "" {
			out = append(out, Statement{p.Label, t})
		}
	}
	return out
}

// Title returns the scorecard's name, or "Balanced Scorecard".
func Title(b *bsc.BSC) string {
	if b.Metadata != nil && b.Metadata.Name != "" {
		return b.Metadata.Name
	}
	return "Balanced Scorecard"
}

// Target formats a measure's target with its unit, and its thresholds when
// set.
func Target(m bsc.Measure) string {
	t := Value(m.Target, m.Unit)
	if th := m.Thresholds; th != nil {
		var bands []string
		if th.Green != "" {
			bands = append(bands, "green at "+Value(th.Green, m.Unit))
		}
		if th.Red != "" {
			bands = append(bands, "red at "+Value(th.Red, m.Unit))
		}
		if len(bands) > 0 {
			t += " (" + strings.Join(bands, ", ") + ")"
		}
	}
	return t
}

// ObjectiveNamesByID returns the names of the objectives, joined by "; ". An
// unknown ID is written as it is.
func ObjectiveNamesByID(b *bsc.BSC, ids []string) string {
	var names []string
	for _, id := range ids {
		if o := b.Objective(id); o != nil {
			names = append(names, o.Name)
		} else {
			names = append(names, id)
		}
	}
	return strings.Join(names, "; ")
}

func objectiveNames(objectives []bsc.Objective) string {
	names := make([]string, 0, len(objectives))
	for _, o := range objectives {
		names = append(names, o.Name)
	}
	return strings.Join(names, "; ")
}

func metadataLine(m *bsc.Metadata) string {
	if m == nil {
		return ""
	}
	var parts []string
	for _, f := range []struct{ label, value string }{
		{"Team", m.Team}, {"Period", m.Period}, {"Author", m.Author}, {"Status", m.Status}, {"Version", m.Version},
	} {
		if f.value != "" {
			parts = append(parts, fmt.Sprintf("**%s:** %s", f.label, f.value))
		}
	}
	return strings.Join(parts, " · ")
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

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// cell makes text safe for a Markdown table cell.
func cell(s string) string {
	return strings.ReplaceAll(oneLine(s), "|", `\|`)
}

// label makes text safe for a quoted Mermaid label.
func label(s string) string {
	return strings.ReplaceAll(oneLine(s), `"`, "#quot;")
}
