// Package markdown renders a MEL plan as a Markdown document: the purpose, the
// results framework, the activities, the indicators and their reference sheets, the
// assumptions to monitor, the evaluation plan, the learning agenda, and the
// review and reporting schedule. The output is deterministic: the same plan
// always renders to the same bytes.
package markdown

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/mel"
	"github.com/grokify/prism-roadmap/goals/mel/render"
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

// Render converts the plan to Markdown, with its name as the title.
func (r *Renderer) Render(p *mel.Plan, opts *render.Options) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("markdown: nil MEL plan")
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", Title(p))
	if line := metadataLine(p.Metadata); line != "" {
		fmt.Fprintf(&b, "%s\n\n", line)
	}
	b.Write(Body(p, opts))
	return []byte(strings.TrimRight(b.String(), "\n") + "\n"), nil
}

// Title returns the plan's name, or "MEL Plan".
func Title(p *mel.Plan) string {
	if p.Metadata != nil && p.Metadata.Name != "" {
		return p.Metadata.Name
	}
	return "MEL Plan"
}

// Body renders every section of the plan, starting at level-two headings and
// without the title, so other renderers can wrap it.
func Body(p *mel.Plan, opts *render.Options) []byte {
	if opts == nil {
		opts = render.DefaultOptions()
	}
	var b strings.Builder
	if p.Purpose != "" {
		fmt.Fprintf(&b, "## Purpose\n\n%s\n\n", p.Purpose)
	}
	writeResults(&b, p, opts)
	writeActivities(&b, p)
	writeIndicators(&b, p)
	if opts.IncludeReferenceSheets {
		writeReferenceSheets(&b, p)
	}
	writeAssumptions(&b, p)
	writeEvaluations(&b, p)
	writeLearning(&b, p, opts)
	writeSchedule(&b, p)
	if p.Budget != "" {
		fmt.Fprintf(&b, "## Budget\n\n%s\n\n", p.Budget)
	}
	return []byte(b.String())
}

func metadataLine(m *mel.Metadata) string {
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

// writeResults draws the results framework as a nested list from the top
// results down. A result that contributes to several is drawn in full the
// first time; later it is a one-line reference, so a shared branch is not
// repeated.
func writeResults(b *strings.Builder, p *mel.Plan, opts *render.Options) {
	if len(p.Results) == 0 {
		return
	}
	b.WriteString("## Results Framework\n\n")
	listed := map[string]bool{}
	var walk func(r mel.Result, depth int)
	walk = func(r mel.Result, depth int) {
		indent := strings.Repeat("  ", depth)
		line := fmt.Sprintf("%s- **%s:** %s (`%s`)", indent, mel.LevelLabel(r.Level), oneLine(r.Statement), r.ID)
		if listed[r.ID] {
			b.WriteString(line + ", shown above\n")
			return
		}
		listed[r.ID] = true // also stops a cycle, which validation reports
		if opts.IncludeDescriptions && r.Description != "" {
			line += " — " + oneLine(r.Description)
		}
		b.WriteString(line + "\n")
		for _, c := range p.Contributors(r.ID) {
			walk(c, depth+1)
		}
	}
	for _, r := range p.TopResults() {
		walk(r, 0)
	}
	b.WriteString("\n")
}

func writeActivities(b *strings.Builder, p *mel.Plan) {
	if len(p.Activities) == 0 {
		return
	}
	b.WriteString("## Activities\n\n")
	b.WriteString("| Activity | Produces | Lead | Timing | Status |\n")
	b.WriteString("|----------|----------|------|--------|--------|\n")
	for _, a := range p.Activities {
		var results []string
		for _, id := range a.ResultIDs {
			results = append(results, resultName(p, id))
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", cell(a.Name), cell(strings.Join(results, "; ")),
			cell(a.Owner), cell(Timing(a.StartDate, a.EndDate)), cell(a.Status))
	}
	b.WriteString("\n")
}

func writeIndicators(b *strings.Builder, p *mel.Plan) {
	if len(p.Indicators) == 0 {
		return
	}
	b.WriteString("## Indicators\n\n")
	b.WriteString("| Indicator | Measures | Baseline | Target | Frequency | Source | Responsible |\n")
	b.WriteString("|-----------|----------|----------|--------|-----------|--------|-------------|\n")
	for _, ind := range p.Indicators {
		baseline := ""
		if ind.Baseline != nil {
			baseline = Value(ind.Baseline.Value, ind.Unit)
		}
		target := ""
		if t, ok := FinalTarget(ind); ok {
			target = Value(t.Value, ind.Unit)
			if t.Date != "" {
				target += " by " + t.Date
			}
		}
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s | %s |\n",
			cell(ind.Name), cell(elementName(p, ind.Measures)), cell(baseline), cell(target),
			cell(ind.Frequency), cell(ind.DataSource), cell(ind.Owner))
	}
	b.WriteString("\n")
}

func writeReferenceSheets(b *strings.Builder, p *mel.Plan) {
	if len(p.Indicators) == 0 {
		return
	}
	b.WriteString("## Indicator Reference Sheets\n\n")
	for _, ind := range p.Indicators {
		fmt.Fprintf(b, "### %s\n\n", ind.Name)
		b.WriteString("| | |\n|---|---|\n")
		row := func(label, value string) {
			if value != "" {
				fmt.Fprintf(b, "| **%s** | %s |\n", label, cell(value))
			}
		}
		row("ID", "`"+ind.ID+"`")
		row("Measures", elementName(p, ind.Measures))
		row("Definition", ind.Definition)
		row("Type", ind.Type)
		row("Unit", ind.Unit)
		row("Direction", ind.Direction)
		row("Disaggregation", strings.Join(ind.Disaggregation, ", "))
		if bl := ind.Baseline; bl != nil {
			v := Value(bl.Value, ind.Unit)
			if bl.Date != "" {
				v += " (" + bl.Date + ")"
			}
			if bl.Source != "" {
				v += "; source: " + bl.Source
			}
			row("Baseline", v)
		}
		var targets []string
		for _, t := range ind.Targets {
			s := Value(t.Value, ind.Unit)
			if t.Date != "" {
				s += " by " + t.Date
			}
			if t.Label != "" {
				s += " (" + t.Label + ")"
			}
			targets = append(targets, s)
		}
		row("Targets", strings.Join(targets, "; "))
		row("Current", Value(ind.Current, ind.Unit))
		row("Data source", ind.DataSource)
		row("Collection method", ind.CollectionMethod)
		row("Frequency", ind.Frequency)
		row("Responsible", ind.Owner)
		if dq := ind.DataQuality; dq != nil {
			row("Data quality checks", strings.Join(dq.Checks, "; "))
			row("Known limitations", strings.Join(dq.Limitations, "; "))
			row("Quality last assessed", dq.LastAssessed)
		}
		b.WriteString("\n")
	}
}

func writeAssumptions(b *strings.Builder, p *mel.Plan) {
	if len(p.Assumptions) == 0 {
		return
	}
	b.WriteString("## Assumptions to Monitor\n\n")
	b.WriteString("| Assumption | Results that depend on it | How it is monitored |\n")
	b.WriteString("|------------|---------------------------|---------------------|\n")
	for _, a := range p.Assumptions {
		var results []string
		for _, id := range a.ResultIDs {
			results = append(results, resultName(p, id))
		}
		monitoring := a.Monitoring
		if len(a.IndicatorIDs) > 0 {
			var names []string
			for _, id := range a.IndicatorIDs {
				names = append(names, indicatorName(p, id))
			}
			if monitoring != "" {
				monitoring += " "
			}
			monitoring += "Indicators: " + strings.Join(names, ", ")
		}
		fmt.Fprintf(b, "| %s | %s | %s |\n", cell(a.Statement), cell(strings.Join(results, "; ")), cell(monitoring))
	}
	b.WriteString("\n")
}

func writeEvaluations(b *strings.Builder, p *mel.Plan) {
	if len(p.Evaluations) == 0 {
		return
	}
	b.WriteString("## Evaluation Plan\n\n")
	b.WriteString("| Evaluation | Type | Timing | Evaluator | Lead | Status |\n")
	b.WriteString("|------------|------|--------|-----------|------|--------|\n")
	for _, e := range p.Evaluations {
		fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s |\n", cell(e.Title), cell(e.Type),
			cell(Timing(e.StartDate, e.EndDate)), cell(e.Evaluator), cell(e.Owner), cell(e.Status))
	}
	b.WriteString("\n")
	for _, e := range p.Evaluations {
		if len(e.Questions) == 0 && e.Method == "" {
			continue
		}
		fmt.Fprintf(b, "### %s\n\n", e.Title)
		if e.Method != "" {
			fmt.Fprintf(b, "**Method:** %s\n\n", oneLine(e.Method))
		}
		if e.Budget != "" {
			fmt.Fprintf(b, "**Budget:** %s\n\n", oneLine(e.Budget))
		}
		for i, q := range e.Questions {
			fmt.Fprintf(b, "%d. %s\n", i+1, oneLine(q))
		}
		if len(e.Questions) > 0 {
			b.WriteString("\n")
		}
	}
}

func writeLearning(b *strings.Builder, p *mel.Plan, opts *render.Options) {
	if len(p.Learning) == 0 {
		return
	}
	b.WriteString("## Learning Agenda\n\n")
	if opts.IncludeDescriptions {
		b.WriteString("| Question | Why it matters | How it will be answered | Lead | Due | Status |\n")
		b.WriteString("|----------|----------------|-------------------------|------|-----|--------|\n")
	} else {
		b.WriteString("| Question | How it will be answered | Lead | Due | Status |\n")
		b.WriteString("|----------|-------------------------|------|-----|--------|\n")
	}
	for _, q := range p.Learning {
		if opts.IncludeDescriptions {
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s |\n", cell(q.Question), cell(q.Rationale),
				cell(strings.Join(q.Methods, "; ")), cell(q.Owner), cell(q.DueDate), cell(q.Status))
		} else {
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s |\n", cell(q.Question),
				cell(strings.Join(q.Methods, "; ")), cell(q.Owner), cell(q.DueDate), cell(q.Status))
		}
	}
	b.WriteString("\n")
	for _, q := range p.Learning {
		if q.Answer != "" {
			fmt.Fprintf(b, "**%s** %s\n\n", oneLine(q.Question), oneLine(q.Answer))
		}
	}
}

func writeSchedule(b *strings.Builder, p *mel.Plan) {
	if len(p.Reviews) == 0 && len(p.Reports) == 0 {
		return
	}
	b.WriteString("## Reviews and Reporting\n\n")
	if len(p.Reviews) > 0 {
		b.WriteString("| Review | Frequency | Participants | Lead | Next | Questions taken up |\n")
		b.WriteString("|--------|-----------|--------------|------|------|--------------------|\n")
		for _, r := range p.Reviews {
			var qs []string
			for _, id := range r.QuestionIDs {
				qs = append(qs, "`"+id+"`")
			}
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s |\n", cell(r.Name), cell(r.Frequency),
				cell(strings.Join(r.Participants, ", ")), cell(r.Owner), cell(r.NextDate), strings.Join(qs, ", "))
		}
		b.WriteString("\n")
	}
	if len(p.Reports) > 0 {
		b.WriteString("| Audience | Needs to know | Receives | Frequency | Prepared by | Indicators |\n")
		b.WriteString("|----------|---------------|----------|-----------|-------------|------------|\n")
		for _, r := range p.Reports {
			var names []string
			for _, id := range r.IndicatorIDs {
				names = append(names, indicatorName(p, id))
			}
			receives := r.Product
			if r.Channel != "" {
				if receives != "" {
					receives += ", in the " + r.Channel
				} else {
					receives = r.Channel
				}
			}
			fmt.Fprintf(b, "| %s | %s | %s | %s | %s | %s |\n", cell(r.Audience), cell(r.Needs), cell(receives),
				cell(r.Frequency), cell(r.Owner), cell(strings.Join(names, "; ")))
		}
		b.WriteString("\n")
	}
}

// FinalTarget returns the indicator's latest-dated target; an undated target
// counts as earliest.
func FinalTarget(ind mel.Indicator) (mel.Target, bool) {
	if len(ind.Targets) == 0 {
		return mel.Target{}, false
	}
	final := ind.Targets[0]
	for _, t := range ind.Targets[1:] {
		if t.Date > final.Date {
			final = t
		}
	}
	return final, true
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

func resultName(p *mel.Plan, id string) string {
	if r := p.Result(id); r != nil {
		return r.Statement
	}
	return id
}

// elementName names what an indicator measures: a result's statement, an
// activity's name, or an assumption's statement.
func elementName(p *mel.Plan, id string) string {
	if r := p.Result(id); r != nil {
		return r.Statement
	}
	if a := p.Activity(id); a != nil {
		return a.Name
	}
	for _, a := range p.Assumptions {
		if a.ID == id {
			return a.Statement
		}
	}
	return id
}

func indicatorName(p *mel.Plan, id string) string {
	if ind := p.Indicator(id); ind != nil {
		return ind.Name
	}
	return id
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// cell makes text safe for a Markdown table cell.
func cell(s string) string {
	return strings.ReplaceAll(oneLine(s), "|", `\|`)
}
