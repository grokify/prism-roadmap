package mel

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ValidationError is one validation finding.
type ValidationError struct {
	Path    string // JSON path to the field, e.g. "indicators[0].dataSource"
	Message string
	IsError bool // true for errors, false for warnings
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationOptions configures which gaps are errors rather than warnings.
// Structural problems, such as a missing ID or statement, a duplicate ID, an
// unknown reference, a cycle in the results framework, an unknown value, or a
// bad date, are always errors.
type ValidationOptions struct {
	// RequireCollectionPlan makes an indicator without a data source,
	// collection method, frequency, or owner an error.
	RequireCollectionPlan bool
	// RequireBaseline makes an indicator without a baseline an error.
	RequireBaseline bool
	// RequireTarget makes an indicator without a target an error.
	RequireTarget bool
	// RequireResultIndicator makes a result that no indicator measures an
	// error.
	RequireResultIndicator bool
	// RequireTiming makes an evaluation without dates, or a learning question
	// without a due date, an error.
	RequireTiming bool
	// RequireEvaluationQuestions makes an evaluation without questions an
	// error.
	RequireEvaluationQuestions bool
	// RequireAssumptionMonitoring makes an assumption without a way to monitor
	// it an error.
	RequireAssumptionMonitoring bool
	// RequireLearningAgenda makes a plan without learning questions an error.
	RequireLearningAgenda bool
	// RequireEvaluationPlan makes a plan without evaluations an error.
	RequireEvaluationPlan bool
	// RequireSchedule makes a plan without reviews, or a review or report
	// without a frequency, an error.
	RequireSchedule bool
}

// DefaultValidationOptions reports gaps in the collection plan, evaluation
// plan, learning agenda, and schedule as warnings.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{}
}

// StrictValidationOptions makes every gap an error: a plan ready to run.
func StrictValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireCollectionPlan:       true,
		RequireBaseline:             true,
		RequireTarget:               true,
		RequireResultIndicator:      true,
		RequireTiming:               true,
		RequireEvaluationQuestions:  true,
		RequireAssumptionMonitoring: true,
		RequireLearningAgenda:       true,
		RequireEvaluationPlan:       true,
		RequireSchedule:             true,
	}
}

// element kinds, for reference checks
const (
	kindResult     = "result"
	kindActivity   = "activity"
	kindIndicator  = "indicator"
	kindAssumption = "assumption"
	kindEvaluation = "evaluation"
	kindQuestion   = "learning question"
	kindReview     = "review"
	kindReport     = "report"
)

type validator struct {
	p    *Plan
	opts *ValidationOptions
	errs []ValidationError
	kind map[string]string // ID to element kind
}

func (v *validator) add(path, msg string, isError bool) {
	v.errs = append(v.errs, ValidationError{Path: path, Message: msg, IsError: isError})
}

func (v *validator) required(path, field, value, what string) {
	if strings.TrimSpace(value) == "" {
		v.add(path+"."+field, what+" is required", true)
	}
}

func (v *validator) oneOf(path, value string, allowed []string) {
	if value != "" && !slices.Contains(allowed, value) {
		v.add(path, fmt.Sprintf("invalid value %q (expected: %s)", value, strings.Join(allowed, ", ")), true)
	}
}

func (v *validator) date(path, value string) bool {
	if value == "" {
		return false
	}
	if _, err := time.Parse(time.DateOnly, value); err != nil {
		v.add(path, fmt.Sprintf("invalid date %q (expected YYYY-MM-DD)", value), true)
		return false
	}
	return true
}

// refs checks that each ID names an element of one of the kinds.
func (v *validator) refs(path string, ids []string, kinds ...string) {
	for i, id := range ids {
		v.ref(fmt.Sprintf("%s[%d]", path, i), id, kinds...)
	}
}

func (v *validator) ref(path, id string, kinds ...string) {
	k, ok := v.kind[id]
	switch {
	case !ok:
		v.add(path, fmt.Sprintf("unknown %s %q", strings.Join(kinds, " or "), id), true)
	case !slices.Contains(kinds, k):
		v.add(path, fmt.Sprintf("%q is a %s, not a %s", id, k, strings.Join(kinds, " or ")), true)
	}
}

// Validate checks the plan and returns its errors and warnings in document
// order.
func (p *Plan) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}
	v := &validator{p: p, opts: opts, kind: map[string]string{}}
	v.identifiers()
	if p.Metadata != nil {
		if strings.TrimSpace(p.Metadata.Name) == "" {
			v.add("metadata.name", "MEL plan name is recommended", false)
		}
		v.oneOf("metadata.status", p.Metadata.Status, DocumentStatuses())
	}
	v.results()
	v.activities()
	v.indicators()
	v.assumptions()
	v.evaluations()
	v.learning()
	v.schedule()
	return v.errs
}

// identifiers records each element's kind and reports missing and duplicate
// IDs.
func (v *validator) identifiers() {
	first := map[string]string{} // ID to the path that first used it
	check := func(path, id, kind string) {
		if strings.TrimSpace(id) == "" {
			v.add(path+".id", "id is required", true)
			return
		}
		if at, ok := first[id]; ok {
			v.add(path+".id", fmt.Sprintf("duplicate id %q (first used at %s)", id, at), true)
			return
		}
		first[id] = path
		v.kind[id] = kind
	}
	p := v.p
	for i, x := range p.Results {
		check(fmt.Sprintf("results[%d]", i), x.ID, kindResult)
	}
	for i, x := range p.Activities {
		check(fmt.Sprintf("activities[%d]", i), x.ID, kindActivity)
	}
	for i, x := range p.Indicators {
		check(fmt.Sprintf("indicators[%d]", i), x.ID, kindIndicator)
	}
	for i, x := range p.Assumptions {
		check(fmt.Sprintf("assumptions[%d]", i), x.ID, kindAssumption)
	}
	for i, x := range p.Evaluations {
		check(fmt.Sprintf("evaluations[%d]", i), x.ID, kindEvaluation)
	}
	for i, x := range p.Learning {
		check(fmt.Sprintf("learningAgenda[%d]", i), x.ID, kindQuestion)
	}
	for i, x := range p.Reviews {
		check(fmt.Sprintf("reviews[%d]", i), x.ID, kindReview)
	}
	for i, x := range p.Reports {
		check(fmt.Sprintf("reports[%d]", i), x.ID, kindReport)
	}
}

func (v *validator) results() {
	p := v.p
	if len(p.Results) == 0 {
		v.add("results", "at least one result is required", true)
	}
	for i, r := range p.Results {
		path := fmt.Sprintf("results[%d]", i)
		v.required(path, "statement", r.Statement, "result statement")
		if r.Level == "" {
			v.add(path+".level", "result level is required", true)
		}
		v.oneOf(path+".level", r.Level, Levels())
		for j, id := range r.ContributesTo {
			jp := fmt.Sprintf("%s.contributesTo[%d]", path, j)
			if id == r.ID {
				v.add(jp, "a result cannot contribute to itself", true)
				continue
			}
			v.ref(jp, id, kindResult)
			if up := p.Result(id); up != nil && LevelRank(up.Level) >= 0 && LevelRank(r.Level) > LevelRank(up.Level) {
				v.add(jp, fmt.Sprintf("a %s contributes to %q, a lower-level %s", LevelLabel(r.Level), id, LevelLabel(up.Level)), false)
			}
		}
		if r.ID != "" && len(p.IndicatorsFor(r.ID)) == 0 {
			v.add(path, "no indicator measures this result", v.opts.RequireResultIndicator)
		}
	}
	if id := v.cycle(); id != "" {
		v.add("results", fmt.Sprintf("the results framework has a cycle through %q", id), true)
	}
}

// cycle returns a result on a cycle of contributesTo links, or "".
func (v *validator) cycle() string {
	const (
		active = iota + 1 // the zero value is unseen
		done
	)
	state := map[string]int{}
	var visit func(id string) string
	visit = func(id string) string {
		switch state[id] {
		case active:
			return id
		case done:
			return ""
		}
		state[id] = active
		if r := v.p.Result(id); r != nil {
			for _, up := range r.ContributesTo {
				if up == id {
					continue // reported as contributing to itself
				}
				if found := visit(up); found != "" {
					return found
				}
			}
		}
		state[id] = done
		return ""
	}
	for _, r := range v.p.Results {
		if found := visit(r.ID); found != "" {
			return found
		}
	}
	return ""
}

func (v *validator) activities() {
	for i, a := range v.p.Activities {
		path := fmt.Sprintf("activities[%d]", i)
		v.required(path, "name", a.Name, "activity name")
		v.refs(path+".resultIds", a.ResultIDs, kindResult)
		v.oneOf(path+".status", a.Status, ActivityStatuses())
		okStart := v.date(path+".startDate", a.StartDate)
		okEnd := v.date(path+".endDate", a.EndDate)
		if okStart && okEnd && a.EndDate < a.StartDate {
			v.add(path+".endDate", fmt.Sprintf("end date %s is before start date %s", a.EndDate, a.StartDate), true)
		}
	}
}

func (v *validator) indicators() {
	for i, ind := range v.p.Indicators {
		path := fmt.Sprintf("indicators[%d]", i)
		v.required(path, "name", ind.Name, "indicator name")
		if ind.Measures == "" {
			v.add(path+".measures", "an indicator must name the result, activity, or assumption it measures", true)
		} else {
			v.ref(path+".measures", ind.Measures, kindResult, kindActivity, kindAssumption)
		}
		v.oneOf(path+".type", ind.Type, IndicatorTypes())
		v.oneOf(path+".direction", ind.Direction, Directions())
		if b := ind.Baseline; b != nil {
			v.required(path+".baseline", "value", b.Value, "baseline value")
			v.date(path+".baseline.date", b.Date)
		} else {
			v.add(path+".baseline", "indicator has no baseline", v.opts.RequireBaseline)
		}
		if len(ind.Targets) == 0 {
			v.add(path+".targets", "indicator has no target", v.opts.RequireTarget)
		}
		for j, t := range ind.Targets {
			tp := fmt.Sprintf("%s.targets[%d]", path, j)
			v.required(tp, "value", t.Value, "target value")
			v.date(tp+".date", t.Date)
		}
		for _, f := range []struct{ field, value, what string }{
			{"dataSource", ind.DataSource, "data source"},
			{"collectionMethod", ind.CollectionMethod, "collection method"},
			{"frequency", ind.Frequency, "frequency"},
			{"owner", ind.Owner, "owner"},
		} {
			if strings.TrimSpace(f.value) == "" {
				v.add(path+"."+f.field, "indicator has no "+f.what, v.opts.RequireCollectionPlan)
			}
		}
		if dq := ind.DataQuality; dq != nil {
			v.date(path+".dataQuality.lastAssessed", dq.LastAssessed)
		}
	}
}

func (v *validator) assumptions() {
	for i, a := range v.p.Assumptions {
		path := fmt.Sprintf("assumptions[%d]", i)
		v.required(path, "statement", a.Statement, "assumption statement")
		v.refs(path+".resultIds", a.ResultIDs, kindResult)
		v.refs(path+".indicatorIds", a.IndicatorIDs, kindIndicator)
		if strings.TrimSpace(a.Monitoring) == "" && len(a.IndicatorIDs) == 0 {
			v.add(path+".monitoring", "nothing says how this assumption is monitored", v.opts.RequireAssumptionMonitoring)
		}
	}
}

func (v *validator) evaluations() {
	if len(v.p.Evaluations) == 0 {
		v.add("evaluations", "the plan has no evaluations", v.opts.RequireEvaluationPlan)
	}
	for i, e := range v.p.Evaluations {
		path := fmt.Sprintf("evaluations[%d]", i)
		v.required(path, "title", e.Title, "evaluation title")
		v.oneOf(path+".type", e.Type, EvaluationTypes())
		v.oneOf(path+".evaluator", e.Evaluator, Evaluators())
		v.oneOf(path+".status", e.Status, EvaluationStatuses())
		if len(e.Questions) == 0 {
			v.add(path+".questions", "evaluation has no questions", v.opts.RequireEvaluationQuestions)
		}
		okStart := v.date(path+".startDate", e.StartDate)
		okEnd := v.date(path+".endDate", e.EndDate)
		if okStart && okEnd && e.EndDate < e.StartDate {
			v.add(path+".endDate", fmt.Sprintf("end date %s is before start date %s", e.EndDate, e.StartDate), true)
		}
		if e.StartDate == "" && e.EndDate == "" {
			v.add(path, "evaluation has no dates", v.opts.RequireTiming)
		}
		v.refs(path+".relatedIds", e.RelatedIDs, kindResult, kindActivity, kindIndicator, kindAssumption, kindQuestion)
	}
}

func (v *validator) learning() {
	if len(v.p.Learning) == 0 {
		v.add("learningAgenda", "the plan has no learning questions", v.opts.RequireLearningAgenda)
	}
	for i, q := range v.p.Learning {
		path := fmt.Sprintf("learningAgenda[%d]", i)
		v.required(path, "question", q.Question, "learning question")
		v.oneOf(path+".status", q.Status, QuestionStatuses())
		v.date(path+".dueDate", q.DueDate)
		if q.DueDate == "" {
			v.add(path+".dueDate", "learning question has no due date", v.opts.RequireTiming)
		}
		if len(q.Methods) == 0 {
			v.add(path+".methods", "nothing says how this question will be answered", false)
		}
		if q.Status == QuestionAnswered && strings.TrimSpace(q.Answer) == "" {
			v.add(path+".answer", "question is marked answered but records no answer", false)
		}
		v.refs(path+".relatedIds", q.RelatedIDs, kindResult, kindActivity, kindIndicator, kindAssumption)
	}
}

func (v *validator) schedule() {
	if len(v.p.Reviews) == 0 {
		v.add("reviews", "the plan schedules no reviews", v.opts.RequireSchedule)
	}
	for i, r := range v.p.Reviews {
		path := fmt.Sprintf("reviews[%d]", i)
		v.required(path, "name", r.Name, "review name")
		if strings.TrimSpace(r.Frequency) == "" {
			v.add(path+".frequency", "review has no frequency", v.opts.RequireSchedule)
		}
		v.date(path+".nextDate", r.NextDate)
		v.refs(path+".questionIds", r.QuestionIDs, kindQuestion)
	}
	for i, r := range v.p.Reports {
		path := fmt.Sprintf("reports[%d]", i)
		v.required(path, "audience", r.Audience, "report audience")
		if strings.TrimSpace(r.Frequency) == "" {
			v.add(path+".frequency", "report has no frequency", v.opts.RequireSchedule)
		}
		v.refs(path+".indicatorIds", r.IndicatorIDs, kindIndicator)
	}
}

// Errors returns only error-level validation results.
func Errors(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if e.IsError {
			result = append(result, e)
		}
	}
	return result
}

// Warnings returns only warning-level validation results.
func Warnings(errs []ValidationError) []ValidationError {
	var result []ValidationError
	for _, e := range errs {
		if !e.IsError {
			result = append(result, e)
		}
	}
	return result
}

// IsValid reports whether there are no error-level validation results.
func IsValid(errs []ValidationError) bool {
	return len(Errors(errs)) == 0
}
