package ogsm

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ValidationError is one validation finding.
type ValidationError struct {
	Path    string // JSON path to the field, e.g. "goals[0].target"
	Message string
	IsError bool // true for errors, false for warnings
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationOptions configures which gaps are errors rather than warnings.
// Structural problems, such as a missing objective, name, or ID, a duplicate
// ID, an unknown reference, an unknown value, or a bad date, are always
// errors.
type ValidationOptions struct {
	// RequireTargets makes a goal or measure without a target an error: in an
	// OGSM, goals and measures are numbers.
	RequireTargets bool
	// RequireOwners makes a goal, measure, or action without an owner an
	// error.
	RequireOwners bool
	// RequireObjectiveDate makes an objective without a due date an error.
	RequireObjectiveDate bool
	// RequireStrategyMeasure makes a strategy that no measure tracks an error.
	RequireStrategyMeasure bool
	// RequireStrategyAction makes a strategy that no action carries out an
	// error.
	RequireStrategyAction bool
	// RequireGoalLinks makes a strategy that serves no goal, or a goal that no
	// strategy serves, an error.
	RequireGoalLinks bool
	// MaxGoals and MaxStrategies warn when an OGSM has more than fit on one
	// page (0 = no limit).
	MaxGoals      int
	MaxStrategies int
}

// DefaultValidationOptions reports gaps as warnings.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{}
}

// StrictValidationOptions makes every gap an error, and warns at more than
// five goals or five strategies.
func StrictValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireTargets:         true,
		RequireOwners:          true,
		RequireObjectiveDate:   true,
		RequireStrategyMeasure: true,
		RequireStrategyAction:  true,
		RequireGoalLinks:       true,
		MaxGoals:               5,
		MaxStrategies:          5,
	}
}

type validator struct {
	o    *OGSM
	opts *ValidationOptions
	errs []ValidationError
	ids  map[string]string // ID to the path that first used it
}

func (v *validator) add(path, msg string, isError bool) {
	v.errs = append(v.errs, ValidationError{Path: path, Message: msg, IsError: isError})
}

func (v *validator) id(path, id string) {
	if strings.TrimSpace(id) == "" {
		v.add(path+".id", "id is required", true)
		return
	}
	if first, ok := v.ids[id]; ok {
		v.add(path+".id", fmt.Sprintf("duplicate id %q (first used at %s)", id, first), true)
		return
	}
	v.ids[id] = path
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

// number checks the fields goals and measures share.
func (v *validator) number(path, name, owner, target, direction, status, dueDate string, progress float64) {
	if strings.TrimSpace(name) == "" {
		v.add(path+".name", "name is required", true)
	}
	if strings.TrimSpace(target) == "" {
		v.add(path+".target", "no target", v.opts.RequireTargets)
	}
	if strings.TrimSpace(owner) == "" {
		v.add(path+".owner", "no owner", v.opts.RequireOwners)
	}
	v.oneOf(path+".direction", direction, Directions())
	v.oneOf(path+".status", status, MetricStatuses())
	v.date(path+".dueDate", dueDate)
	if progress < 0 || progress > 1 {
		v.add(path+".progress", fmt.Sprintf("progress %g is outside 0 to 1", progress), true)
	}
}

// Validate checks the OGSM and returns its errors and warnings in document
// order.
func (o *OGSM) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}
	v := &validator{o: o, opts: opts, ids: map[string]string{}}

	if m := o.Metadata; m != nil {
		if strings.TrimSpace(m.Name) == "" {
			v.add("metadata.name", "OGSM name is recommended", false)
		}
		v.oneOf("metadata.status", m.Status, DocumentStatuses())
	}
	if strings.TrimSpace(o.Objective.Statement) == "" {
		v.add("objective.statement", "an objective is required", true)
	}
	v.date("objective.dueDate", o.Objective.DueDate)
	if o.Objective.DueDate == "" {
		v.add("objective.dueDate", "objective has no due date", opts.RequireObjectiveDate)
	}

	if len(o.Goals) == 0 {
		v.add("goals", "at least one goal is required", true)
	}
	if opts.MaxGoals > 0 && len(o.Goals) > opts.MaxGoals {
		v.add("goals", fmt.Sprintf("%d goals is more than fit on one page (max: %d)", len(o.Goals), opts.MaxGoals), false)
	}
	for i, g := range o.Goals {
		path := fmt.Sprintf("goals[%d]", i)
		v.id(path, g.ID)
		v.number(path, g.Name, g.Owner, g.Target, g.Direction, g.Status, g.DueDate, g.Progress)
		if g.ID != "" && len(o.StrategiesFor(g.ID)) == 0 {
			v.add(path, "no strategy serves this goal", opts.RequireGoalLinks)
		}
	}

	if len(o.Strategies) == 0 {
		v.add("strategies", "at least one strategy is required", true)
	}
	if opts.MaxStrategies > 0 && len(o.Strategies) > opts.MaxStrategies {
		v.add("strategies", fmt.Sprintf("%d strategies is more than fit on one page (max: %d)", len(o.Strategies), opts.MaxStrategies), false)
	}
	strategies := map[string]bool{}
	for i, s := range o.Strategies {
		path := fmt.Sprintf("strategies[%d]", i)
		v.id(path, s.ID)
		strategies[s.ID] = true
		if strings.TrimSpace(s.Name) == "" {
			v.add(path+".name", "strategy name is required", true)
		}
		v.oneOf(path+".priority", s.Priority, Priorities())
		if len(s.GoalIDs) == 0 {
			v.add(path+".goalIds", "strategy serves no goal", opts.RequireGoalLinks)
		}
		for j, id := range s.GoalIDs {
			if o.Goal(id) == nil {
				v.add(fmt.Sprintf("%s.goalIds[%d]", path, j), fmt.Sprintf("unknown goal %q", id), true)
			}
		}
		if s.ID != "" && len(o.MeasuresFor(s.ID)) == 0 {
			v.add(path, "no measure tracks this strategy", opts.RequireStrategyMeasure)
		}
		if s.ID != "" && len(o.ActionsFor(s.ID)) == 0 {
			v.add(path, "no action carries out this strategy", opts.RequireStrategyAction)
		}
	}

	strategyRef := func(path, id, what string) {
		switch {
		case id == "":
			v.add(path+".strategyId", what+" must name its strategy", true)
		case !strategies[id]:
			v.add(path+".strategyId", fmt.Sprintf("unknown strategy %q", id), true)
		}
	}
	for i, m := range o.Measures {
		path := fmt.Sprintf("measures[%d]", i)
		v.id(path, m.ID)
		strategyRef(path, m.StrategyID, "a measure")
		v.number(path, m.Name, m.Owner, m.Target, m.Direction, m.Status, m.DueDate, m.Progress)
	}

	actions := map[string]bool{}
	for _, a := range o.Actions {
		actions[a.ID] = true
	}
	for i, a := range o.Actions {
		path := fmt.Sprintf("actions[%d]", i)
		v.id(path, a.ID)
		if strings.TrimSpace(a.Name) == "" {
			v.add(path+".name", "action name is required", true)
		}
		strategyRef(path, a.StrategyID, "an action")
		if strings.TrimSpace(a.Owner) == "" {
			v.add(path+".owner", "no owner", opts.RequireOwners)
		}
		v.oneOf(path+".status", a.Status, ActionStatuses())
		okStart := v.date(path+".startDate", a.StartDate)
		okEnd := v.date(path+".endDate", a.EndDate)
		if okStart && okEnd && a.EndDate < a.StartDate {
			v.add(path+".endDate", fmt.Sprintf("end date %s is before start date %s", a.EndDate, a.StartDate), true)
		}
		for j, dep := range a.Dependencies {
			dp := fmt.Sprintf("%s.dependencies[%d]", path, j)
			switch {
			case dep == a.ID:
				v.add(dp, "action depends on itself", true)
			case !actions[dep]:
				v.add(dp, fmt.Sprintf("unknown action %q", dep), true)
			}
		}
	}
	return v.errs
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
