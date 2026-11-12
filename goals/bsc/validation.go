package bsc

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ValidationError is one validation finding.
type ValidationError struct {
	Path    string // JSON path to the field, e.g. "measures[0].target"
	Message string
	IsError bool // true for errors, false for warnings
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationOptions configures which gaps are errors rather than warnings.
// Structural problems, such as a missing name or ID, a duplicate ID, an
// unknown reference, a cycle in the strategy map, an unknown value, or a bad
// date, are always errors.
type ValidationOptions struct {
	// RequireVision makes a scorecard with neither a mission nor a vision an
	// error.
	RequireVision bool
	// RequireTargets makes a measure without a target an error.
	RequireTargets bool
	// RequireOwners makes an objective, measure, or initiative without an
	// owner an error.
	RequireOwners bool
	// RequirePerspectiveObjective makes a perspective with no objective an
	// error: a scorecard with an empty perspective is not balanced.
	RequirePerspectiveObjective bool
	// RequireObjectiveMeasure makes an objective that no measure tracks an
	// error.
	RequireObjectiveMeasure bool
	// RequireStrategyMap makes an objective outside the top perspective that
	// contributes to nothing, or one outside the bottom perspective that
	// nothing contributes to, an error.
	RequireStrategyMap bool
	// RequireMeasureType makes a measure that is neither lead nor lag an
	// error.
	RequireMeasureType bool
	// RequireLeadAndLag makes a scorecard without both lead and lag measures
	// an error: results alone do not say how they will be reached.
	RequireLeadAndLag bool
	// MaxMeasures warns when a scorecard has more measures than a team can
	// steer by (0 = no limit).
	MaxMeasures int
}

// DefaultValidationOptions reports gaps as warnings.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{}
}

// StrictValidationOptions makes every gap an error, and warns at more than 25
// measures. An objective that no initiative serves stays a warning, since an
// objective already on track needs none.
func StrictValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireVision:               true,
		RequireTargets:              true,
		RequireOwners:               true,
		RequirePerspectiveObjective: true,
		RequireObjectiveMeasure:     true,
		RequireStrategyMap:          true,
		RequireMeasureType:          true,
		RequireLeadAndLag:           true,
		MaxMeasures:                 25,
	}
}

type validator struct {
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

func (v *validator) name(path, name, what string) {
	if strings.TrimSpace(name) == "" {
		v.add(path+".name", what+" name is required", true)
	}
}

func (v *validator) owner(path, owner string) {
	if strings.TrimSpace(owner) == "" {
		v.add(path+".owner", "no owner", v.opts.RequireOwners)
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

// Validate checks the scorecard and returns its errors and warnings in
// document order.
func (b *BSC) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}
	v := &validator{opts: opts, ids: map[string]string{}}

	if m := b.Metadata; m != nil {
		if strings.TrimSpace(m.Name) == "" {
			v.add("metadata.name", "scorecard name is recommended", false)
		}
		v.oneOf("metadata.status", m.Status, DocumentStatuses())
	}
	if strings.TrimSpace(b.Mission) == "" && strings.TrimSpace(b.Vision) == "" {
		v.add("vision", "the scorecard has neither a mission nor a vision", opts.RequireVision)
	}

	b.validatePerspectives(v)
	for i, t := range b.Themes {
		path := fmt.Sprintf("themes[%d]", i)
		v.id(path, t.ID)
		v.name(path, t.Name, "theme")
		if t.ID != "" && len(b.ObjectivesInTheme(t.ID)) == 0 {
			v.add(path, "no objective belongs to this theme", false)
		}
	}
	b.validateObjectives(v)
	b.validateMeasures(v)
	b.validateInitiatives(v)
	return v.errs
}

func (b *BSC) validatePerspectives(v *validator) {
	if len(b.Perspectives) == 0 {
		v.add("perspectives", "at least one perspective is required", true)
	}
	for i, p := range b.Perspectives {
		path := fmt.Sprintf("perspectives[%d]", i)
		v.id(path, p.ID)
		v.name(path, p.Name, "perspective")
		v.oneOf(path+".kind", p.Kind, PerspectiveKinds())
		if p.ID != "" && len(b.ObjectivesIn(p.ID)) == 0 {
			v.add(path, "no objective in this perspective", v.opts.RequirePerspectiveObjective)
		}
	}
}

func (b *BSC) validateObjectives(v *validator) {
	if len(b.Objectives) == 0 {
		v.add("objectives", "at least one objective is required", true)
	}
	last := len(b.Perspectives) - 1
	for i, o := range b.Objectives {
		path := fmt.Sprintf("objectives[%d]", i)
		v.id(path, o.ID)
		v.name(path, o.Name, "objective")
		v.owner(path, o.Owner)
		v.oneOf(path+".priority", o.Priority, Priorities())
		level := b.PerspectiveIndex(o.PerspectiveID)
		switch {
		case o.PerspectiveID == "":
			v.add(path+".perspectiveId", "an objective must name its perspective", true)
		case level < 0:
			v.add(path+".perspectiveId", fmt.Sprintf("unknown perspective %q", o.PerspectiveID), true)
		}
		for j, id := range o.ThemeIDs {
			if b.Theme(id) == nil {
				v.add(fmt.Sprintf("%s.themeIds[%d]", path, j), fmt.Sprintf("unknown theme %q", id), true)
			}
		}
		for j, id := range o.ContributesTo {
			lp := fmt.Sprintf("%s.contributesTo[%d]", path, j)
			target := b.Objective(id)
			switch {
			case id == o.ID:
				v.add(lp, "an objective cannot contribute to itself", true)
			case target == nil:
				v.add(lp, fmt.Sprintf("unknown objective %q", id), true)
			case level >= 0 && b.PerspectiveIndex(target.PerspectiveID) > level:
				v.add(lp, fmt.Sprintf("contributes to %q in a lower perspective; a strategy map runs from cause at the bottom to effect at the top", id), false)
			}
		}
		if level > 0 && len(o.ContributesTo) == 0 {
			v.add(path+".contributesTo", "contributes to no objective in the strategy map", v.opts.RequireStrategyMap)
		}
		if level >= 0 && level < last && o.ID != "" && len(b.ContributorsTo(o.ID)) == 0 {
			v.add(path, "no objective contributes to this one in the strategy map", v.opts.RequireStrategyMap)
		}
		if o.ID != "" && len(b.MeasuresFor(o.ID)) == 0 {
			v.add(path, "no measure tracks this objective", v.opts.RequireObjectiveMeasure)
		}
		if o.ID != "" && len(b.InitiativesFor(o.ID)) == 0 {
			v.add(path, "no initiative serves this objective", false)
		}
	}
	if cycle := b.strategyMapCycle(); cycle != nil {
		v.add("objectives", "the strategy map has a cycle: "+strings.Join(cycle, " -> "), true)
	}
}

// strategyMapCycle returns the first cycle in the strategy map, as the IDs
// along it with the first repeated at the end, or nil.
func (b *BSC) strategyMapCycle() []string {
	const (
		unvisited = iota
		onPath
		done
	)
	state := map[string]int{}
	var path []string
	var visit func(id string) []string
	visit = func(id string) []string {
		state[id] = onPath
		path = append(path, id)
		if o := b.Objective(id); o != nil {
			for _, next := range o.ContributesTo {
				if next == id || b.Objective(next) == nil {
					continue // reported on the link itself
				}
				switch state[next] {
				case onPath:
					start := slices.Index(path, next)
					return append(slices.Clone(path[start:]), next)
				case unvisited:
					if c := visit(next); c != nil {
						return c
					}
				}
			}
		}
		path = path[:len(path)-1]
		state[id] = done
		return nil
	}
	for _, o := range b.Objectives {
		if o.ID != "" && state[o.ID] == unvisited {
			if c := visit(o.ID); c != nil {
				return c
			}
		}
	}
	return nil
}

func (b *BSC) validateMeasures(v *validator) {
	if v.opts.MaxMeasures > 0 && len(b.Measures) > v.opts.MaxMeasures {
		v.add("measures", fmt.Sprintf("%d measures is more than a team can steer by (max: %d)", len(b.Measures), v.opts.MaxMeasures), false)
	}
	var lead, lag bool
	for i, m := range b.Measures {
		path := fmt.Sprintf("measures[%d]", i)
		v.id(path, m.ID)
		v.name(path, m.Name, "measure")
		switch {
		case m.ObjectiveID == "":
			v.add(path+".objectiveId", "a measure must name its objective", true)
		case b.Objective(m.ObjectiveID) == nil:
			v.add(path+".objectiveId", fmt.Sprintf("unknown objective %q", m.ObjectiveID), true)
		}
		if m.Type == "" {
			v.add(path+".type", "neither lead nor lag", v.opts.RequireMeasureType)
		}
		v.oneOf(path+".type", m.Type, MeasureTypes())
		lead = lead || m.Type == TypeLead
		lag = lag || m.Type == TypeLag
		if strings.TrimSpace(m.Target) == "" {
			v.add(path+".target", "no target", v.opts.RequireTargets)
		}
		v.owner(path, m.Owner)
		v.oneOf(path+".direction", m.Direction, Directions())
		v.oneOf(path+".status", m.Status, MeasureStatuses())
		v.date(path+".dueDate", m.DueDate)
		if m.Progress < 0 || m.Progress > 1 {
			v.add(path+".progress", fmt.Sprintf("progress %g is outside 0 to 1", m.Progress), true)
		}
	}
	if len(b.Measures) > 0 && (!lead || !lag) {
		v.add("measures", "a balanced scorecard has both lead and lag measures", v.opts.RequireLeadAndLag)
	}
}

func (b *BSC) validateInitiatives(v *validator) {
	initiatives := map[string]bool{}
	for _, in := range b.Initiatives {
		initiatives[in.ID] = true
	}
	for i, in := range b.Initiatives {
		path := fmt.Sprintf("initiatives[%d]", i)
		v.id(path, in.ID)
		v.name(path, in.Name, "initiative")
		if len(in.ObjectiveIDs) == 0 {
			v.add(path+".objectiveIds", "an initiative must name the objectives it serves", true)
		}
		for j, id := range in.ObjectiveIDs {
			if b.Objective(id) == nil {
				v.add(fmt.Sprintf("%s.objectiveIds[%d]", path, j), fmt.Sprintf("unknown objective %q", id), true)
			}
		}
		v.owner(path, in.Owner)
		v.oneOf(path+".status", in.Status, InitiativeStatuses())
		okStart := v.date(path+".startDate", in.StartDate)
		okEnd := v.date(path+".endDate", in.EndDate)
		if okStart && okEnd && in.EndDate < in.StartDate {
			v.add(path+".endDate", fmt.Sprintf("end date %s is before start date %s", in.EndDate, in.StartDate), true)
		}
		for j, dep := range in.Dependencies {
			dp := fmt.Sprintf("%s.dependencies[%d]", path, j)
			switch {
			case dep == in.ID:
				v.add(dp, "initiative depends on itself", true)
			case !initiatives[dep]:
				v.add(dp, fmt.Sprintf("unknown initiative %q", dep), true)
			}
		}
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
