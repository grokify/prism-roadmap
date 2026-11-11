// Package vsem provides types and utilities for VSEM (Vision, Strategy,
// Execution, Metrics) strategic plans.
//
// A VSEM fits a team's plan on one page in four sections:
//
//   - Vision: where the team or organization wants to end up. Some teams write
//     a mission, why they exist, in its place; a VSEM holds either or both.
//   - Strategy: how the vision will be reached. Each strategy names a choice
//     and the resources it commits, such as people, budget, or tooling.
//   - Execution: the initiatives, programs, and actions that carry out each
//     strategy, each with an owner and timing.
//   - Metrics: how progress is measured, each with the person accountable for
//     it, a baseline, and a target.
//
// VSEM has no single canonical source, and published versions differ: some
// put a mission in place of the vision, and some describe the fourth section
// as accountability rather than measurement. This package follows the form
// above, which keeps both readings: every metric names an accountable owner.
//
// Compared with V2MOM, a VSEM has no values or obstacles, states the
// resources each strategy commits, and keeps the work (execution) apart from
// the choices it serves (strategy). Execution items and metrics refer to their
// strategy by ID; a metric without a strategy measures the vision as a whole.
package vsem

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Document status constants.
const (
	StatusDraft     = "Draft"
	StatusInReview  = "In Review"
	StatusApproved  = "Approved"
	StatusActive    = "Active"
	StatusCompleted = "Completed"
	StatusArchived  = "Archived"
)

// Priority constants.
const (
	PriorityP0 = "P0"
	PriorityP1 = "P1"
	PriorityP2 = "P2"
	PriorityP3 = "P3"
)

// Execution kind constants: the forms of work a VSEM's execution section lists.
const (
	KindInitiative = "initiative" // a bounded effort with an outcome, such as a launch
	KindProgram    = "program"    // ongoing work that runs across periods
	KindAction     = "action"     // a single task or step
)

// Execution status constants.
const (
	ExecutionPlanned    = "Planned"
	ExecutionInProgress = "In Progress"
	ExecutionAtRisk     = "At Risk"
	ExecutionCompleted  = "Completed"
	ExecutionCancelled  = "Cancelled"
)

// Metric status constants.
const (
	MetricOnTrack  = "On Track"
	MetricAtRisk   = "At Risk"
	MetricBehind   = "Behind"
	MetricAchieved = "Achieved"
	MetricMissed   = "Missed"
)

// Direction constants say which way a metric should move.
const (
	DirectionIncrease = "increase"
	DirectionDecrease = "decrease"
	DirectionMaintain = "maintain"
)

// VSEM is a complete VSEM strategic plan.
type VSEM struct {
	Schema   string    `json:"$schema,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	// Vision is where the team or organization wants to end up.
	Vision string `json:"vision,omitempty"`
	// Mission is why the team or organization exists. A VSEM states a vision,
	// a mission, or both.
	Mission    string      `json:"mission,omitempty"`
	Strategies []Strategy  `json:"strategies"`
	Execution  []Execution `json:"execution,omitempty"`
	Metrics    []Metric    `json:"metrics,omitempty"`
}

// Metadata contains document metadata.
type Metadata struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Author  string `json:"author,omitempty"`
	Team    string `json:"team,omitempty"`
	Period  string `json:"period,omitempty"` // e.g. "FY2027" or "2027-H1"
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
	// ParentID is the VSEM this one cascades from.
	ParentID  string    `json:"parentId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}

// Strategy is one choice about how the vision will be reached.
type Strategy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Priority    string `json:"priority,omitempty"` // P0, P1, P2, P3
	// Resources are what the strategy commits: where people, money, and
	// tooling will go.
	Resources []Resource `json:"resources,omitempty"`
}

// Resource is one resource a strategy commits.
type Resource struct {
	// Name is the kind of resource, e.g. "Engineering" or "Marketing budget".
	Name string `json:"name"`
	// Allocation is how much is committed, as written, e.g. "3 engineers" or
	// "40% of the quarter's budget".
	Allocation  string `json:"allocation,omitempty"`
	Description string `json:"description,omitempty"`
}

// Execution is an initiative, program, or action that carries out a strategy.
type Execution struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// StrategyID is the strategy this work carries out.
	StrategyID string `json:"strategyId"`
	Kind       string `json:"kind,omitempty"` // initiative, program, action
	Owner      string `json:"owner,omitempty"`
	Status     string `json:"status,omitempty"`
	StartDate  string `json:"startDate,omitempty"` // ISO 8601 date
	EndDate    string `json:"endDate,omitempty"`   // ISO 8601 date
	// Dependencies are the IDs of execution items this one waits on.
	Dependencies  []string          `json:"dependencies,omitempty"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"`
}

// Metric is how progress is measured, and who answers for it.
type Metric struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// StrategyID is the strategy this metric measures; empty means the metric
	// measures the vision as a whole.
	StrategyID string `json:"strategyId,omitempty"`
	// Owner is the person accountable for the metric.
	Owner     string  `json:"owner,omitempty"`
	Baseline  string  `json:"baseline,omitempty"`
	Target    string  `json:"target,omitempty"`
	Current   string  `json:"current,omitempty"`
	Unit      string  `json:"unit,omitempty"`
	Direction string  `json:"direction,omitempty"` // increase, decrease, maintain
	DueDate   string  `json:"dueDate,omitempty"`   // ISO 8601 date the target is due
	Frequency string  `json:"frequency,omitempty"` // how often it is measured, e.g. "Monthly"
	Source    string  `json:"source,omitempty"`    // where the data comes from
	Progress  float64 `json:"progress,omitempty"`  // 0.0-1.0
	Status    string  `json:"status,omitempty"`
}

// New returns a VSEM with the given ID, name, and vision.
func New(id, name, vision string) *VSEM {
	return &VSEM{
		Metadata: &Metadata{ID: id, Name: name, Status: StatusDraft},
		Vision:   vision,
	}
}

// ReadFile reads a VSEM from a JSON file.
func ReadFile(filepath string) (*VSEM, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses VSEM JSON data.
func Parse(data []byte) (*VSEM, error) {
	var v VSEM
	if err := json.Unmarshal(data, &v); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &v, nil
}

// JSON returns the VSEM as formatted JSON.
func (v *VSEM) JSON() ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// WriteFile writes the VSEM to a JSON file.
func (v *VSEM) WriteFile(filepath string) error {
	data, err := v.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// Statement returns the vision, or the mission when there is no vision.
func (v *VSEM) Statement() string {
	if v.Vision != "" {
		return v.Vision
	}
	return v.Mission
}

// Strategy returns the strategy with the given ID, or nil.
func (v *VSEM) Strategy(id string) *Strategy {
	for i := range v.Strategies {
		if v.Strategies[i].ID == id {
			return &v.Strategies[i]
		}
	}
	return nil
}

// ExecutionFor returns the execution items that carry out the strategy, in
// document order.
func (v *VSEM) ExecutionFor(strategyID string) []Execution {
	var items []Execution
	for _, e := range v.Execution {
		if e.StrategyID == strategyID {
			items = append(items, e)
		}
	}
	return items
}

// MetricsFor returns the metrics that measure the strategy, in document order.
// An empty strategyID returns the metrics that measure the vision as a whole.
func (v *VSEM) MetricsFor(strategyID string) []Metric {
	var items []Metric
	for _, m := range v.Metrics {
		if m.StrategyID == strategyID {
			items = append(items, m)
		}
	}
	return items
}

// VisionMetrics returns the metrics that measure the vision as a whole.
func (v *VSEM) VisionMetrics() []Metric {
	return v.MetricsFor("")
}

// Kinds returns the execution kinds.
func Kinds() []string {
	return []string{KindInitiative, KindProgram, KindAction}
}

// DocumentStatuses returns the document statuses.
func DocumentStatuses() []string {
	return []string{StatusDraft, StatusInReview, StatusApproved, StatusActive, StatusCompleted, StatusArchived}
}

// ExecutionStatuses returns the execution statuses.
func ExecutionStatuses() []string {
	return []string{ExecutionPlanned, ExecutionInProgress, ExecutionAtRisk, ExecutionCompleted, ExecutionCancelled}
}

// MetricStatuses returns the metric statuses.
func MetricStatuses() []string {
	return []string{MetricOnTrack, MetricAtRisk, MetricBehind, MetricAchieved, MetricMissed}
}

// Priorities returns the priorities.
func Priorities() []string {
	return []string{PriorityP0, PriorityP1, PriorityP2, PriorityP3}
}

// Directions returns the metric directions.
func Directions() []string {
	return []string{DirectionIncrease, DirectionDecrease, DirectionMaintain}
}
