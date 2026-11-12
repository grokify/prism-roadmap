// Package markdown renders an OGSM as a Markdown document: the plan on one
// page as the classic four-column table, then the objective, the goals, and
// each strategy with its measures and actions. The output is deterministic:
// the same OGSM always renders to the same bytes.
package markdown

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/internal/mdshift"
	"github.com/grokify/prism-roadmap/goals/ogsm"
	"github.com/grokify/prism-roadmap/goals/ogsm/render"
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

// Render converts the OGSM to Markdown.
func (r *Renderer) Render(o *ogsm.OGSM, opts *render.Options) ([]byte, error) {
	if o == nil {
		return nil, fmt.Errorf("markdown: nil OGSM")
	}
	if opts == nil {
		opts = render.DefaultOptions()
	}
	var b strings.Builder
	if !opts.OmitHeader {
		fmt.Fprintf(&b, "# %s\n\n", Title(o))
		if line := metadataLine(o.Metadata); line != "" {
			fmt.Fprintf(&b, "%s\n\n", line)
		}
	}

	b.WriteString("## On One Page\n\n")
	b.WriteString(OnePage(o))
	b.WriteString("\n")

	b.WriteString("## Objective\n\n")
	fmt.Fprintf(&b, "%s\n\n", oneLine(o.Objective.Statement))
	var facts []string
	if o.Objective.DueDate != "" {
		facts = append(facts, "**By:** "+o.Objective.DueDate)
	}
	if o.Objective.ParentStrategyID != "" {
		facts = append(facts, "**Carries out:** the parent strategy `"+o.Objective.ParentStrategyID+"`")
	}
	if len(facts) > 0 {
		fmt.Fprintf(&b, "%s\n\n", strings.Join(facts, " · "))
	}
	if opts.IncludeDescriptions && o.Objective.Description != "" {
		fmt.Fprintf(&b, "%s\n\n", oneLine(o.Objective.Description))
	}

	b.WriteString("## Goals\n\n")
	b.WriteString("| Goal | Owner | Baseline | Target | Current | Due | Status |\n")
	b.WriteString("|------|-------|----------|--------|---------|-----|--------|\n")
	for _, g := range o.Goals {
		name := cell(g.Name)
		if opts.IncludeDescriptions && g.Description != "" {
			name += "<br>" + cell(g.Description)
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n", name, cell(g.Owner),
			cell(Value(g.Baseline, g.Unit)), cell(Value(g.Target, g.Unit)), cell(Value(g.Current, g.Unit)),
			cell(g.DueDate), cell(g.Status))
	}
	b.WriteString("\n")

	b.WriteString("## Strategies\n\n")
	for i, s := range o.Strategies {
		fmt.Fprintf(&b, "### %d. %s\n\n", i+1, s.Name)
		var facts []string
		if s.Owner != "" {
			facts = append(facts, "**Owner:** "+s.Owner)
		}
		if s.Priority != "" {
			facts = append(facts, "**Priority:** "+s.Priority)
		}
		if serves := goalNames(o, s.GoalIDs); serves != "" {
			facts = append(facts, "**Serves:** "+serves)
		}
		if len(facts) > 0 {
			fmt.Fprintf(&b, "%s\n\n", strings.Join(facts, " · "))
		}
		if opts.IncludeDescriptions && s.Description != "" {
			fmt.Fprintf(&b, "%s\n\n", oneLine(s.Description))
		}
		if measures := o.MeasuresFor(s.ID); len(measures) > 0 {
			b.WriteString("**Measures**\n\n")
			b.WriteString("| Measure | Owner | Baseline | Target | Current | Frequency | Status |\n")
			b.WriteString("|---------|-------|----------|--------|---------|-----------|--------|\n")
			for _, m := range measures {
				fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s |\n", cell(m.Name), cell(m.Owner),
					cell(Value(m.Baseline, m.Unit)), cell(Value(m.Target, m.Unit)), cell(Value(m.Current, m.Unit)),
					cell(m.Frequency), cell(m.Status))
			}
			b.WriteString("\n")
		}
		if actions := o.ActionsFor(s.ID); opts.IncludeActions && len(actions) > 0 {
			b.WriteString("**Actions**\n\n")
			b.WriteString("| Action | Owner | Timing | Status |\n")
			b.WriteString("|--------|-------|--------|--------|\n")
			for _, a := range actions {
				name := cell(a.Name)
				if opts.IncludeDescriptions && a.Description != "" {
					name += "<br>" + cell(a.Description)
				}
				fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", name, cell(a.Owner), cell(Timing(a.StartDate, a.EndDate)), cell(a.Status))
			}
			b.WriteString("\n")
		}
	}
	return []byte(mdshift.Shift(strings.TrimRight(b.String(), "\n")+"\n", opts.HeadingOffset)), nil
}

// OnePage renders the classic one-page OGSM: the objective, the goals, the
// strategies, and each strategy's measures, side by side.
func OnePage(o *ogsm.OGSM) string {
	var b strings.Builder
	b.WriteString("| Objective | Goals | Strategies | Measures |\n")
	b.WriteString("|-----------|-------|------------|----------|\n")
	rows := max(len(o.Goals), len(o.Strategies), 1)
	for i := range rows {
		objective, goal, strategy, measures := "", "", "", ""
		if i == 0 {
			objective = cell(o.Objective.Statement)
		}
		if i < len(o.Goals) {
			g := o.Goals[i]
			goal = cell(g.Name)
			if g.Target != "" {
				goal += ": " + cell(Value(g.Target, g.Unit))
				if g.DueDate != "" {
					goal += " by " + g.DueDate
				}
			}
		}
		if i < len(o.Strategies) {
			s := o.Strategies[i]
			strategy = cell(s.Name)
			var names []string
			for _, m := range o.MeasuresFor(s.ID) {
				name := cell(m.Name)
				if m.Target != "" {
					name += ": " + cell(Value(m.Target, m.Unit))
				}
				names = append(names, name)
			}
			measures = strings.Join(names, "<br>")
		}
		fmt.Fprintf(&b, "| %s | %s | %s | %s |\n", objective, goal, strategy, measures)
	}
	return b.String()
}

// Title returns the OGSM's name, or "OGSM".
func Title(o *ogsm.OGSM) string {
	if o.Metadata != nil && o.Metadata.Name != "" {
		return o.Metadata.Name
	}
	return "OGSM"
}

func metadataLine(m *ogsm.Metadata) string {
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

func goalNames(o *ogsm.OGSM, ids []string) string {
	var names []string
	for _, id := range ids {
		if g := o.Goal(id); g != nil {
			names = append(names, g.Name)
		} else {
			names = append(names, id)
		}
	}
	return strings.Join(names, "; ")
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
