package vsem

import (
	"fmt"
	"slices"
	"strings"
	"time"
)

// ValidationError is one validation finding.
type ValidationError struct {
	Path    string // JSON path to the field, e.g. "strategies[0].name"
	Message string
	IsError bool // true for errors, false for warnings
}

// Error implements the error interface.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Path, e.Message)
}

// ValidationOptions configures which gaps are errors rather than warnings.
// Structural problems, such as a missing name, an unknown strategy ID, a
// duplicate ID, or an unknown status, are always errors.
type ValidationOptions struct {
	// RequireMetricOwner makes a metric without an accountable owner an error.
	RequireMetricOwner bool
	// RequireMetricTarget makes a metric without a target an error.
	RequireMetricTarget bool
	// RequireExecutionOwner makes an execution item without an owner an error.
	RequireExecutionOwner bool
	// RequireStrategyExecution makes a strategy that no execution item carries
	// out an error.
	RequireStrategyExecution bool
	// RequireStrategyMetric makes a strategy that no metric measures an error.
	RequireStrategyMetric bool
	// RequireResources makes a strategy that commits no resources an error.
	RequireResources bool
	// MaxStrategies warns when a VSEM has more strategies than fit on one
	// page (0 = no limit).
	MaxStrategies int
}

// DefaultValidationOptions reports gaps in accountability and coverage as
// warnings.
func DefaultValidationOptions() *ValidationOptions {
	return &ValidationOptions{}
}

// StrictValidationOptions makes every metric and execution item name an
// owner, every metric a target, and every strategy commit resources and have
// work and a metric.
func StrictValidationOptions() *ValidationOptions {
	return &ValidationOptions{
		RequireMetricOwner:       true,
		RequireMetricTarget:      true,
		RequireExecutionOwner:    true,
		RequireStrategyExecution: true,
		RequireStrategyMetric:    true,
		RequireResources:         true,
		MaxStrategies:            5,
	}
}

// Validate checks the VSEM and returns its errors and warnings in document
// order.
func (v *VSEM) Validate(opts *ValidationOptions) []ValidationError {
	if opts == nil {
		opts = DefaultValidationOptions()
	}
	var errs []ValidationError
	add := func(path, msg string, isError bool) {
		errs = append(errs, ValidationError{Path: path, Message: msg, IsError: isError})
	}

	if v.Metadata != nil {
		errs = append(errs, validateMetadata(v.Metadata)...)
	}
	if strings.TrimSpace(v.Vision) == "" && strings.TrimSpace(v.Mission) == "" {
		add("vision", "a vision or a mission is required", true)
	}
	if len(v.Strategies) == 0 {
		add("strategies", "at least one strategy is required", true)
	}
	if opts.MaxStrategies > 0 && len(v.Strategies) > opts.MaxStrategies {
		add("strategies", fmt.Sprintf("%d strategies is more than fit on one page (max: %d)", len(v.Strategies), opts.MaxStrategies), false)
	}

	seen := make(map[string]string) // ID to the path that first used it
	checkID := func(path, id string) {
		if strings.TrimSpace(id) == "" {
			add(path+".id", "id is required", true)
			return
		}
		if first, ok := seen[id]; ok {
			add(path+".id", fmt.Sprintf("duplicate id %q (first used at %s)", id, first), true)
			return
		}
		seen[id] = path
	}

	strategies := make(map[string]bool)
	for i, s := range v.Strategies {
		path := fmt.Sprintf("strategies[%d]", i)
		checkID(path, s.ID)
		strategies[s.ID] = true
		if strings.TrimSpace(s.Name) == "" {
			add(path+".name", "strategy name is required", true)
		}
		if s.Priority != "" && !slices.Contains(Priorities(), s.Priority) {
			add(path+".priority", fmt.Sprintf("invalid priority %q (expected: %s)", s.Priority, strings.Join(Priorities(), ", ")), true)
		}
		if len(s.Resources) == 0 {
			add(path+".resources", "strategy commits no resources", opts.RequireResources)
		}
		for j, r := range s.Resources {
			if strings.TrimSpace(r.Name) == "" {
				add(fmt.Sprintf("%s.resources[%d].name", path, j), "resource name is required", true)
			}
		}
		if s.ID != "" && len(v.ExecutionFor(s.ID)) == 0 {
			add(path, "no execution item carries out this strategy", opts.RequireStrategyExecution)
		}
		if s.ID != "" && len(v.MetricsFor(s.ID)) == 0 {
			add(path, "no metric measures this strategy", opts.RequireStrategyMetric)
		}
	}

	executions := make(map[string]bool)
	for _, e := range v.Execution {
		executions[e.ID] = true
	}
	for i, e := range v.Execution {
		path := fmt.Sprintf("execution[%d]", i)
		checkID(path, e.ID)
		if strings.TrimSpace(e.Name) == "" {
			add(path+".name", "execution item name is required", true)
		}
		switch {
		case e.StrategyID == "":
			add(path+".strategyId", "execution item must name the strategy it carries out", true)
		case !strategies[e.StrategyID]:
			add(path+".strategyId", fmt.Sprintf("unknown strategy %q", e.StrategyID), true)
		}
		if e.Kind != "" && !slices.Contains(Kinds(), e.Kind) {
			add(path+".kind", fmt.Sprintf("invalid kind %q (expected: %s)", e.Kind, strings.Join(Kinds(), ", ")), true)
		}
		if e.Status != "" && !slices.Contains(ExecutionStatuses(), e.Status) {
			add(path+".status", fmt.Sprintf("invalid status %q (expected: %s)", e.Status, strings.Join(ExecutionStatuses(), ", ")), true)
		}
		if strings.TrimSpace(e.Owner) == "" {
			add(path+".owner", "execution item has no owner", opts.RequireExecutionOwner)
		}
		errs = append(errs, validateDateRange(path, e.StartDate, e.EndDate)...)
		for j, dep := range e.Dependencies {
			depPath := fmt.Sprintf("%s.dependencies[%d]", path, j)
			switch {
			case dep == e.ID:
				add(depPath, "execution item depends on itself", true)
			case !executions[dep]:
				add(depPath, fmt.Sprintf("unknown execution item %q", dep), true)
			}
		}
	}

	for i, m := range v.Metrics {
		path := fmt.Sprintf("metrics[%d]", i)
		checkID(path, m.ID)
		if strings.TrimSpace(m.Name) == "" {
			add(path+".name", "metric name is required", true)
		}
		if m.StrategyID != "" && !strategies[m.StrategyID] {
			add(path+".strategyId", fmt.Sprintf("unknown strategy %q", m.StrategyID), true)
		}
		if strings.TrimSpace(m.Owner) == "" {
			add(path+".owner", "metric has no accountable owner", opts.RequireMetricOwner)
		}
		if strings.TrimSpace(m.Target) == "" {
			add(path+".target", "metric has no target", opts.RequireMetricTarget)
		}
		if m.Direction != "" && !slices.Contains(Directions(), m.Direction) {
			add(path+".direction", fmt.Sprintf("invalid direction %q (expected: %s)", m.Direction, strings.Join(Directions(), ", ")), true)
		}
		if m.Status != "" && !slices.Contains(MetricStatuses(), m.Status) {
			add(path+".status", fmt.Sprintf("invalid status %q (expected: %s)", m.Status, strings.Join(MetricStatuses(), ", ")), true)
		}
		if m.Progress < 0 || m.Progress > 1 {
			add(path+".progress", fmt.Sprintf("progress %g is outside 0 to 1", m.Progress), true)
		}
		if m.DueDate != "" {
			if _, err := time.Parse(time.DateOnly, m.DueDate); err != nil {
				add(path+".dueDate", fmt.Sprintf("invalid date %q (expected YYYY-MM-DD)", m.DueDate), true)
			}
		}
	}

	return errs
}

func validateMetadata(meta *Metadata) []ValidationError {
	var errs []ValidationError
	if strings.TrimSpace(meta.Name) == "" {
		errs = append(errs, ValidationError{Path: "metadata.name", Message: "VSEM name is recommended"})
	}
	if meta.Status != "" && !slices.Contains(DocumentStatuses(), meta.Status) {
		errs = append(errs, ValidationError{
			Path:    "metadata.status",
			Message: fmt.Sprintf("invalid status %q (expected: %s)", meta.Status, strings.Join(DocumentStatuses(), ", ")),
			IsError: true,
		})
	}
	return errs
}

// validateDateRange checks that the dates parse and that start is not after
// end.
func validateDateRange(path, start, end string) []ValidationError {
	var errs []ValidationError
	parse := func(field, value string) (time.Time, bool) {
		if value == "" {
			return time.Time{}, false
		}
		t, err := time.Parse(time.DateOnly, value)
		if err != nil {
			errs = append(errs, ValidationError{
				Path:    path + "." + field,
				Message: fmt.Sprintf("invalid date %q (expected YYYY-MM-DD)", value),
				IsError: true,
			})
			return time.Time{}, false
		}
		return t, true
	}
	s, okStart := parse("startDate", start)
	e, okEnd := parse("endDate", end)
	if okStart && okEnd && s.After(e) {
		errs = append(errs, ValidationError{
			Path:    path + ".endDate",
			Message: fmt.Sprintf("end date %s is before start date %s", end, start),
			IsError: true,
		})
	}
	return errs
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
