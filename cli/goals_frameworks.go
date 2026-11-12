package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/grokify/prism-roadmap/goals/bsc"
	bscrender "github.com/grokify/prism-roadmap/goals/bsc/render"
	bscmarkdown "github.com/grokify/prism-roadmap/goals/bsc/render/markdown"
	bscmarp "github.com/grokify/prism-roadmap/goals/bsc/render/marp"
	"github.com/grokify/prism-roadmap/goals/dmaic"
	dmaicrender "github.com/grokify/prism-roadmap/goals/dmaic/render"
	dmaicmarp "github.com/grokify/prism-roadmap/goals/dmaic/render/marp"
	dmaicpandoc "github.com/grokify/prism-roadmap/goals/dmaic/render/pandoc"
	"github.com/grokify/prism-roadmap/goals/mel"
	melrender "github.com/grokify/prism-roadmap/goals/mel/render"
	melmarkdown "github.com/grokify/prism-roadmap/goals/mel/render/markdown"
	melpandoc "github.com/grokify/prism-roadmap/goals/mel/render/pandoc"
	"github.com/grokify/prism-roadmap/goals/ogsm"
	ogsmrender "github.com/grokify/prism-roadmap/goals/ogsm/render"
	ogsmmarkdown "github.com/grokify/prism-roadmap/goals/ogsm/render/markdown"
	ogsmmarp "github.com/grokify/prism-roadmap/goals/ogsm/render/marp"
	"github.com/grokify/prism-roadmap/goals/okr"
	"github.com/grokify/prism-roadmap/goals/v2mom"
	"github.com/grokify/prism-roadmap/goals/vsem"
	vsemrender "github.com/grokify/prism-roadmap/goals/vsem/render"
	vsemmarkdown "github.com/grokify/prism-roadmap/goals/vsem/render/markdown"
	vsemmarp "github.com/grokify/prism-roadmap/goals/vsem/render/marp"
)

// GoalFrameworks are the goal document formats the goals commands read, in
// the order they are listed.
var GoalFrameworks = []string{"okr", "v2mom", "vsem", "ogsm", "bsc", "dmaic", "mel"}

// GoalFinding is one validation result, whatever the framework.
type GoalFinding struct {
	Path    string
	Message string
	IsError bool
}

func (f GoalFinding) String() string { return f.Path + ": " + f.Message }

// DetectGoalFramework tells a goal document's framework from its top-level
// fields: perspectives (Balanced Scorecard, checked first because it also
// has objectives), objectives (OKR), methods or values (V2MOM), a single
// objective (OGSM), strategies (VSEM), categories (DMAIC), or results (MEL).
func DetectGoalFramework(data []byte) (string, error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return "", fmt.Errorf("parsing JSON: %w", err)
	}
	for _, c := range []struct{ field, framework string }{
		{"perspectives", "bsc"}, {"objectives", "okr"}, {"methods", "v2mom"}, {"values", "v2mom"},
		{"objective", "ogsm"}, {"strategies", "vsem"}, {"categories", "dmaic"}, {"results", "mel"},
	} {
		if _, ok := fields[c.field]; ok {
			return c.framework, nil
		}
	}
	return "", fmt.Errorf("cannot tell the framework from the file's fields; pass --framework (%s)", strings.Join(GoalFrameworks, ", "))
}

// ValidateGoals validates a goal document with its framework's validation.
// strict uses the framework's strict options where it has them (V2MOM does
// not). An empty framework is detected from the document.
func ValidateGoals(data []byte, framework string, strict bool) (string, []GoalFinding, error) {
	if framework == "" {
		var err error
		if framework, err = DetectGoalFramework(data); err != nil {
			return "", nil, err
		}
	}
	var out []GoalFinding
	switch framework {
	case "okr":
		var doc okr.OKRDocument
		if err := json.Unmarshal(data, &doc); err != nil {
			return framework, nil, fmt.Errorf("parsing OKR document: %w", err)
		}
		opts := okr.DefaultValidationOptions()
		if strict {
			opts = okr.StrictValidationOptions()
		}
		for _, e := range doc.Validate(opts) {
			out = append(out, GoalFinding{e.Path, e.Message, e.IsError})
		}
	case "v2mom":
		v, err := v2mom.Parse(data)
		if err != nil {
			return framework, nil, err
		}
		for _, e := range v.Validate(v2mom.DefaultValidationOptions()) {
			out = append(out, GoalFinding{e.Path, e.Message, e.Severity == "error"})
		}
	case "vsem":
		v, err := vsem.Parse(data)
		if err != nil {
			return framework, nil, err
		}
		opts := vsem.DefaultValidationOptions()
		if strict {
			opts = vsem.StrictValidationOptions()
		}
		for _, e := range v.Validate(opts) {
			out = append(out, GoalFinding{e.Path, e.Message, e.IsError})
		}
	case "ogsm":
		o, err := ogsm.Parse(data)
		if err != nil {
			return framework, nil, err
		}
		opts := ogsm.DefaultValidationOptions()
		if strict {
			opts = ogsm.StrictValidationOptions()
		}
		for _, e := range o.Validate(opts) {
			out = append(out, GoalFinding{e.Path, e.Message, e.IsError})
		}
	case "bsc":
		b, err := bsc.Parse(data)
		if err != nil {
			return framework, nil, err
		}
		opts := bsc.DefaultValidationOptions()
		if strict {
			opts = bsc.StrictValidationOptions()
		}
		for _, e := range b.Validate(opts) {
			out = append(out, GoalFinding{e.Path, e.Message, e.IsError})
		}
	case "dmaic":
		doc, err := dmaic.Parse(data)
		if err != nil {
			return framework, nil, err
		}
		opts := dmaic.DefaultValidationOptions()
		if strict {
			opts = dmaic.StrictValidationOptions()
		}
		for _, e := range doc.Validate(opts) {
			out = append(out, GoalFinding{e.Path, e.Message, e.IsError})
		}
	case "mel":
		p, err := mel.Parse(data)
		if err != nil {
			return framework, nil, err
		}
		opts := mel.DefaultValidationOptions()
		if strict {
			opts = mel.StrictValidationOptions()
		}
		for _, e := range p.Validate(opts) {
			out = append(out, GoalFinding{e.Path, e.Message, e.IsError})
		}
	default:
		return framework, nil, fmt.Errorf("unknown framework %q (use %s)", framework, strings.Join(GoalFrameworks, ", "))
	}
	return framework, out, nil
}

// reportGoalFindings prints warnings, then errors, and returns an error when
// there is at least one error.
func reportGoalFindings(w io.Writer, file, framework string, findings []GoalFinding) error {
	var errs, warns []GoalFinding
	for _, f := range findings {
		if f.IsError {
			errs = append(errs, f)
		} else {
			warns = append(warns, f)
		}
	}
	for _, group := range []struct {
		title string
		items []GoalFinding
	}{{"Warnings", warns}, {"Errors", errs}} {
		if len(group.items) == 0 {
			continue
		}
		fmt.Fprintf(w, "%s:\n", group.title)
		for _, f := range group.items {
			fmt.Fprintf(w, "  - %s\n", f)
		}
		fmt.Fprintln(w)
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s: validation failed with %d error(s)", file, len(errs))
	}
	fmt.Fprintf(w, "✓ %s is a valid %s document", file, strings.ToUpper(framework))
	if len(warns) > 0 {
		fmt.Fprintf(w, " (%d warning(s))", len(warns))
	}
	fmt.Fprintln(w)
	return nil
}

func validateGoalsFile(cmd *cobra.Command, file, framework string, strict bool) error {
	data, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}
	framework, findings, err := ValidateGoals(data, framework, strict)
	if err != nil {
		return fmt.Errorf("%s: %w", file, err)
	}
	return reportGoalFindings(cmd.OutOrStdout(), file, framework, findings)
}

// writeGoalsOutput writes to the file, or to standard output when it is empty.
func writeGoalsOutput(cmd *cobra.Command, file string, data []byte) error {
	if file == "" {
		_, err := cmd.OutOrStdout().Write(data)
		return err
	}
	if err := os.WriteFile(file, data, 0600); err != nil {
		return fmt.Errorf("writing %s: %w", file, err)
	}
	fmt.Fprintf(cmd.ErrOrStderr(), "Generated: %s\n", file)
	return nil
}

// writeNewGoalsFile writes a template, refusing to overwrite an existing file.
func writeNewGoalsFile(cmd *cobra.Command, file string, v any, next string) error {
	if _, err := os.Stat(file); err == nil {
		return fmt.Errorf("file already exists: %s (use -o to choose another path)", file)
	}
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding template: %w", err)
	}
	if err := os.WriteFile(file, append(data, '\n'), 0600); err != nil {
		return fmt.Errorf("writing %s: %w", file, err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "Created: %s\n%s", file, next)
	return nil
}

// rendered reads a document with read and renders it with render, which
// receives the --theme value.
func rendered[T any](read func(string) (*T, error), render func(*T, string) ([]byte, error)) func(string, bool, string) ([]byte, error) {
	return func(file string, _ bool, theme string) ([]byte, error) {
		doc, err := read(file)
		if err != nil {
			return nil, err
		}
		return render(doc, theme)
	}
}

// frameworkCommands builds the validate, generate, and init subcommands for one
// framework.
type frameworkCommands struct {
	name, title, long string
	generators        []generator
	template          func(name string) any
}

type generator struct {
	format, short, long string
	// run renders the file with the flags given.
	run     func(file string, summary bool, theme string) ([]byte, error)
	summary bool // offers --summary
	theme   bool // offers --theme
}

func (f frameworkCommands) command() *cobra.Command {
	root := &cobra.Command{Use: f.name, Short: "Work with " + f.title + " documents", Long: f.long}

	var strict bool
	validate := &cobra.Command{
		Use:   "validate FILE",
		Short: "Validate a " + f.title + " JSON file",
		Long: fmt.Sprintf(`Validate a %s JSON file. Structural problems are errors; gaps are warnings,
or errors with --strict.

Examples:
  splan goals %s validate plan.json
  splan goals %s validate plan.json --strict`, f.title, f.name, f.name),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return validateGoalsFile(cmd, args[0], f.name, strict)
		},
	}
	validate.Flags().BoolVar(&strict, "strict", false, "treat gaps as errors")
	root.AddCommand(validate)

	generate := &cobra.Command{Use: "generate", Short: "Generate output from a " + f.title + " document"}
	for _, g := range f.generators {
		var out, theme string
		var summary bool
		c := &cobra.Command{
			Use:     g.format + " FILE",
			Short:   g.short,
			Long:    g.long,
			Example: fmt.Sprintf("  splan goals %s generate %s plan.json -o plan.md", f.name, g.format),
			Args:    cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				data, err := g.run(args[0], summary, theme)
				if err != nil {
					return err
				}
				return writeGoalsOutput(cmd, out, data)
			},
		}
		c.Flags().StringVarP(&out, "output", "o", "", "output file (default: stdout)")
		if g.summary {
			c.Flags().BoolVar(&summary, "summary", false, "leave out reference sheets and descriptions")
		}
		if g.theme {
			c.Flags().StringVar(&theme, "theme", "default", "slide theme (default, corporate, minimal)")
		}
		generate.AddCommand(c)
	}
	root.AddCommand(generate)

	var name, out string
	initCmd := &cobra.Command{
		Use:   "init",
		Short: "Create a " + f.title + " template",
		Long: fmt.Sprintf(`Create a %s JSON template with example content to edit.

Examples:
  splan goals %s init --name "FY2027 plan" -o plan.json`, f.title, f.name),
		RunE: func(cmd *cobra.Command, args []string) error {
			return writeNewGoalsFile(cmd, out, f.template(name), fmt.Sprintf(
				"Next: edit it, then run 'splan goals %s validate %s'\n", f.name, out))
		},
	}
	initCmd.Flags().StringVar(&name, "name", "My "+f.title, "document name")
	initCmd.Flags().StringVarP(&out, "output", "o", f.name+".json", "output file")
	root.AddCommand(initCmd)
	return root
}

var vsemCommands = frameworkCommands{
	name: "vsem", title: "VSEM",
	long: `Commands for VSEM (Vision, Strategy, Execution, Metrics) plans: validate them,
render them as a Markdown document or Marp slides, and create a template.`,
	generators: []generator{
		{format: "markdown", short: "Generate a Markdown document", long: "Render a VSEM as a Markdown document: the vision, each strategy with its resources, execution, and metrics, then the vision's metrics.",
			run: rendered(vsem.ReadFile, func(v *vsem.VSEM, _ string) ([]byte, error) {
				return vsemmarkdown.New().Render(v, vsemrender.DefaultOptions())
			})},
		{format: "marp", short: "Generate Marp slides", long: "Render a VSEM as Marp slides: the plan on one page, the vision, a slide per strategy, and the vision's metrics.", theme: true,
			run: rendered(vsem.ReadFile, func(v *vsem.VSEM, theme string) ([]byte, error) {
				opts := vsemrender.DefaultOptions()
				opts.Theme = theme
				return vsemmarp.New().Render(v, opts)
			})},
	},
	template: func(name string) any { return VSEMTemplate(name) },
}

var ogsmCommands = frameworkCommands{
	name: "ogsm", title: "OGSM",
	long: `Commands for OGSM (Objective, Goals, Strategies, Measures) plans: validate
them, render them as a Markdown document or Marp slides, and create a template.`,
	generators: []generator{
		{format: "markdown", short: "Generate a Markdown document", long: "Render an OGSM as a Markdown document: the plan on one page, the objective, the goals, and each strategy with its measures and actions.",
			run: rendered(ogsm.ReadFile, func(o *ogsm.OGSM, _ string) ([]byte, error) {
				return ogsmmarkdown.New().Render(o, ogsmrender.DefaultOptions())
			})},
		{format: "marp", short: "Generate Marp slides", long: "Render an OGSM as Marp slides: the plan on one page, the objective and goals, and a slide per strategy.", theme: true,
			run: rendered(ogsm.ReadFile, func(o *ogsm.OGSM, theme string) ([]byte, error) {
				opts := ogsmrender.DefaultOptions()
				opts.Theme = theme
				return ogsmmarp.New().Render(o, opts)
			})},
	},
	template: func(name string) any { return OGSMTemplate(name) },
}

var bscCommands = frameworkCommands{
	name: "bsc", title: "Balanced Scorecard",
	long: `Commands for Balanced Scorecards (perspectives, strategic objectives and
their strategy map, measures, and initiatives): validate them, render them as a
Markdown document or Marp slides, and create a template.`,
	generators: []generator{
		{format: "markdown", short: "Generate a Markdown document", long: "Render a Balanced Scorecard as a Markdown document: the mission, vision, and strategy, the strategy map as a Mermaid diagram, the strategic themes, the scorecard by perspective, and the initiatives.",
			run: rendered(bsc.ReadFile, func(b *bsc.BSC, _ string) ([]byte, error) {
				return bscmarkdown.New().Render(b, bscrender.DefaultOptions())
			})},
		{format: "marp", short: "Generate Marp slides", long: "Render a Balanced Scorecard as Marp slides: the mission, vision, and strategy, the strategy map, a slide per perspective, and the initiatives.", theme: true,
			run: rendered(bsc.ReadFile, func(b *bsc.BSC, theme string) ([]byte, error) {
				opts := bscrender.DefaultOptions()
				opts.Theme = theme
				return bscmarp.New().Render(b, opts)
			})},
	},
	template: func(name string) any { return BSCTemplate(name) },
}

var dmaicCommands = frameworkCommands{
	name: "dmaic", title: "DMAIC",
	long: `Commands for DMAIC (Define, Measure, Analyze, Improve, Control) metrics
documents: validate them, render them as Marp slides or Pandoc Markdown, and
create a template.`,
	generators: []generator{
		{format: "marp", short: "Generate Marp slides", long: "Render a DMAIC document as Marp slides.", theme: true,
			run: rendered(dmaic.ReadFile, func(doc *dmaic.DMAICDocument, theme string) ([]byte, error) {
				opts := dmaicrender.DefaultOptions()
				opts.Theme = theme
				return dmaicmarp.New().Render(doc, opts)
			})},
		{format: "pandoc", short: "Generate Pandoc Markdown for a PDF", long: "Render a DMAIC document as Pandoc Markdown with front matter for a PDF.",
			run: rendered(dmaic.ReadFile, func(doc *dmaic.DMAICDocument, _ string) ([]byte, error) {
				return dmaicpandoc.New().Render(doc, dmaicrender.DefaultOptions())
			})},
	},
	template: func(name string) any { return DMAICTemplate(name) },
}

var melCommands = frameworkCommands{
	name: "mel", title: "MEL",
	long: `Commands for monitoring, evaluation, and learning (MEL) plans: validate them,
render them as a Markdown document or Pandoc Markdown for a PDF, and create a
template.`,
	generators: []generator{
		{format: "markdown", short: "Generate a Markdown document", long: "Render a MEL plan as a Markdown document. --summary leaves out the indicator reference sheets and descriptions.", summary: true,
			run: func(file string, summary bool, _ string) ([]byte, error) {
				p, err := mel.ReadFile(file)
				if err != nil {
					return nil, err
				}
				return melmarkdown.New().Render(p, melOptions(summary))
			}},
		{format: "pandoc", short: "Generate Pandoc Markdown for a PDF", long: "Render a MEL plan as Pandoc Markdown with front matter for a PDF. --summary leaves out the indicator reference sheets and descriptions.", summary: true,
			run: func(file string, summary bool, _ string) ([]byte, error) {
				p, err := mel.ReadFile(file)
				if err != nil {
					return nil, err
				}
				return melpandoc.New().Render(p, melOptions(summary))
			}},
	},
	template: func(name string) any { return MELTemplate(name) },
}

func melOptions(summary bool) *melrender.Options {
	if summary {
		return melrender.SummaryOptions()
	}
	return melrender.DefaultOptions()
}

// VSEMTemplate returns a small VSEM to start from.
func VSEMTemplate(name string) *vsem.VSEM {
	v := vsem.New(slug(name), name, "Where we want to be at the end of the period")
	v.Strategies = []vsem.Strategy{{
		ID: "strategy-1", Name: "The first choice about how to get there", Owner: "Strategy owner", Priority: vsem.PriorityP0,
		Resources: []vsem.Resource{{Name: "Team", Allocation: "Who and how much"}},
	}}
	v.Execution = []vsem.Execution{{
		ID: "work-1", Name: "The first initiative that carries it out", StrategyID: "strategy-1",
		Kind: vsem.KindInitiative, Owner: "Initiative owner", Status: vsem.ExecutionPlanned,
	}}
	v.Metrics = []vsem.Metric{{
		ID: "metric-1", Name: "How you will know the strategy is working", StrategyID: "strategy-1",
		Owner: "Accountable person", Baseline: "0", Target: "1", Direction: vsem.DirectionIncrease,
	}}
	return v
}

// OGSMTemplate returns a small OGSM to start from. It meets the strict
// options, so its objective is due three years from today.
func OGSMTemplate(name string) *ogsm.OGSM {
	o := ogsm.New(slug(name), name, "The one thing to achieve over the plan's horizon")
	o.Objective.DueDate = time.Now().AddDate(3, 0, 0).Format(time.DateOnly)
	o.Goals = []ogsm.Goal{{
		ID: "goal-1", Name: "A number that says the objective is reached", Owner: "Goal owner",
		Baseline: "0", Target: "1", Direction: ogsm.DirectionIncrease, DueDate: o.Objective.DueDate,
	}}
	o.Strategies = []ogsm.Strategy{{
		ID: "strategy-1", Name: "A choice about how to reach the goals", Owner: "Strategy owner",
		Priority: ogsm.PriorityP0, GoalIDs: []string{"goal-1"},
	}}
	o.Measures = []ogsm.Measure{{
		ID: "measure-1", Name: "A number that shows the strategy is working", StrategyID: "strategy-1",
		Owner: "Measure owner", Baseline: "0", Target: "1", Direction: ogsm.DirectionIncrease, Frequency: "Monthly",
	}}
	o.Actions = []ogsm.Action{{
		ID: "action-1", Name: "The first action that carries out the strategy", StrategyID: "strategy-1",
		Owner: "Action owner", Status: ogsm.ActionPlanned,
	}}
	return o
}

// BSCTemplate returns a small Balanced Scorecard to start from: the four
// classic perspectives with an objective each, linked from learning and
// growth up to the financial result, a measure on each, and an initiative.
// It meets the strict options.
func BSCTemplate(name string) *bsc.BSC {
	b := bsc.New(slug(name), name)
	b.Vision = "What the organization wants to become"
	b.Strategy = "How it will get there, in a sentence or two"
	due := time.Now().AddDate(1, 0, 0).Format(time.DateOnly)
	steps := []struct{ id, perspective, name, measureType string }{
		{"financial-1", "financial", "The financial result the strategy delivers", bsc.TypeLag},
		{"customer-1", "customer", "What customers must experience", bsc.TypeLag},
		{"process-1", "internal-process", "The process the organization must excel at", bsc.TypeLead},
		{"learning-1", "learning-growth", "The capability the organization must build", bsc.TypeLead},
	}
	ids := make([]string, 0, len(steps))
	above := ""
	for i, st := range steps {
		o := bsc.Objective{ID: st.id, Name: st.name, PerspectiveID: st.perspective, Owner: "Objective owner"}
		if above != "" {
			o.ContributesTo = []string{above}
		}
		above = st.id
		ids = append(ids, st.id)
		b.Objectives = append(b.Objectives, o)
		b.Measures = append(b.Measures, bsc.Measure{
			ID: fmt.Sprintf("measure-%d", i+1), Name: "How you will know: " + strings.ToLower(st.name[:1]) + st.name[1:],
			ObjectiveID: st.id, Type: st.measureType, Owner: "Measure owner", Baseline: "0", Target: "1",
			Direction: bsc.DirectionIncrease, DueDate: due, Frequency: "Quarterly",
		})
	}
	b.Initiatives = []bsc.Initiative{{
		ID: "initiative-1", Name: "The first initiative that closes the gap", ObjectiveIDs: ids,
		Owner: "Initiative owner", Status: bsc.InitiativePlanned,
	}}
	return b
}

// DMAICTemplate returns a small DMAIC metrics document to start from.
func DMAICTemplate(name string) *dmaic.DMAICDocument {
	doc := dmaic.New(slug(name), name, "Process owner")
	doc.Categories = []dmaic.Category{{
		ID: "category-1", Name: "The process being improved",
		Metrics: []dmaic.Metric{{
			ID: "metric-1", Name: "The metric being improved", Baseline: 10, Current: 10, Target: 5,
			Phase: dmaic.PhaseDefine, TrendDirection: dmaic.TrendLowerBetter, Frequency: "Weekly",
		}},
	}}
	return doc
}

// MELTemplate returns a small MEL plan to start from. It meets the strict
// options, so its dates are set from today.
func MELTemplate(name string) *mel.Plan {
	date := func(months int) string { return time.Now().AddDate(0, months, 0).Format(time.DateOnly) }
	p := mel.New(slug(name), name)
	p.Results = []mel.Result{
		{ID: "impact-1", Statement: "The lasting change the work contributes to", Level: mel.LevelImpact},
		{ID: "outcome-1", Statement: "The change you expect first", Level: mel.LevelShortTerm, ContributesTo: []string{"impact-1"}},
	}
	p.Activities = []mel.Activity{{ID: "activity-1", Name: "The work that produces it", ResultIDs: []string{"outcome-1"}, Status: mel.ActivityPlanned}}
	indicator := func(id, name, measures string) mel.Indicator {
		return mel.Indicator{
			ID: id, Name: name, Measures: measures, Type: mel.TypeQuantitative,
			Baseline: &mel.Baseline{Value: "0"}, Targets: []mel.Target{{Value: "1", Date: date(12)}},
			DataSource: "Where the data comes from", CollectionMethod: "How it is gathered",
			Frequency: "Monthly", Owner: "Responsible person",
		}
	}
	p.Indicators = []mel.Indicator{
		indicator("indicator-1", "How you will know the lasting change is happening", "impact-1"),
		indicator("indicator-2", "How you will know the first change happened", "outcome-1"),
	}
	p.Assumptions = []mel.Assumption{{
		ID: "assumption-1", Statement: "What must hold for the first change to lead to the lasting one",
		ResultIDs: []string{"impact-1"}, Monitoring: "How you will know whether it still holds",
	}}
	p.Evaluations = []mel.Evaluation{{
		ID: "evaluation-1", Title: "Final evaluation", Type: mel.EvalFinal,
		Questions: []string{"Did the lasting change happen, and how much did the work contribute?"},
		StartDate: date(10), EndDate: date(12), Status: mel.EvaluationPlanned, RelatedIDs: []string{"impact-1", "assumption-1"},
	}}
	p.Learning = []mel.LearningQuestion{{
		ID: "question-1", Question: "What do you most need to find out?", Methods: []string{"How you will answer it"},
		DueDate: date(6), Status: mel.QuestionOpen, RelatedIDs: []string{"outcome-1", "assumption-1"},
	}}
	p.Reviews = []mel.Review{{ID: "review-1", Name: "Quarterly learning review", Frequency: "Quarterly", QuestionIDs: []string{"question-1"}}}
	return p
}

// slug makes a document ID from a name.
func slug(name string) string {
	var b strings.Builder
	dash := false
	for _, r := range strings.ToLower(name) {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			dash = false
		} else if !dash && b.Len() > 0 {
			b.WriteByte('-')
			dash = true
		}
	}
	return strings.TrimSuffix(b.String(), "-")
}

// isGoalFramework reports whether name is one of GoalFrameworks.
func isGoalFramework(name string) bool { return slices.Contains(GoalFrameworks, name) }
