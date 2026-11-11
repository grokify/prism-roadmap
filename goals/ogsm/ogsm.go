// Package ogsm provides types and utilities for OGSM (Objective, Goals,
// Strategies, Measures) plans.
//
// An OGSM puts a plan on one page in four parts, two in words and two in
// numbers:
//
//   - Objective, in words: the one thing to achieve over the plan's horizon.
//   - Goals, in numbers: how the organization will know the objective has
//     been reached, each with a baseline, a target, and a due date.
//   - Strategies, in words: the choices about how to reach the goals.
//   - Measures, in numbers: for each strategy, the progress indicators that
//     show whether it is working, before the goals move.
//
// Most OGSMs also list, for each strategy, the actions that carry it out, with
// owners and dates; this package holds them as Actions.
//
// OGSMs cascade: a strategy at one level becomes the objective of the level
// below. An objective names the strategy it carries out with
// ParentStrategyID, and the metadata names the parent OGSM with ParentID.
// Measures and actions refer to their strategy by ID, and a strategy names the
// goals it serves.
package ogsm

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

// Status constants for goals and measures.
const (
	MetricOnTrack  = "On Track"
	MetricAtRisk   = "At Risk"
	MetricBehind   = "Behind"
	MetricAchieved = "Achieved"
	MetricMissed   = "Missed"
)

// Action status constants.
const (
	ActionPlanned    = "Planned"
	ActionInProgress = "In Progress"
	ActionAtRisk     = "At Risk"
	ActionCompleted  = "Completed"
	ActionCancelled  = "Cancelled"
)

// Direction constants say which way a goal or measure should move.
const (
	DirectionIncrease = "increase"
	DirectionDecrease = "decrease"
	DirectionMaintain = "maintain"
)

// OGSM is a complete OGSM plan.
type OGSM struct {
	Schema     string     `json:"$schema,omitempty"`
	Metadata   *Metadata  `json:"metadata,omitempty"`
	Objective  Objective  `json:"objective"`
	Goals      []Goal     `json:"goals"`
	Strategies []Strategy `json:"strategies"`
	Measures   []Measure  `json:"measures,omitempty"`
	Actions    []Action   `json:"actions,omitempty"`
}

// Metadata contains document metadata.
type Metadata struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Author  string `json:"author,omitempty"`
	Team    string `json:"team,omitempty"`
	Period  string `json:"period,omitempty"` // e.g. "2027-2029"
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
	// ParentID is the OGSM this one cascades from.
	ParentID  string    `json:"parentId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}

// Objective is the one thing the plan sets out to achieve, in words.
type Objective struct {
	Statement   string `json:"statement"`
	Description string `json:"description,omitempty"`
	// DueDate is the end of the plan's horizon, an ISO 8601 date.
	DueDate string `json:"dueDate,omitempty"`
	// ParentStrategyID is the strategy of the parent OGSM that this objective
	// carries out, in a cascade.
	ParentStrategyID string `json:"parentStrategyId,omitempty"`
}

// Goal is a number that says the objective has been reached.
type Goal struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Owner       string  `json:"owner,omitempty"`
	Baseline    string  `json:"baseline,omitempty"`
	Target      string  `json:"target,omitempty"`
	Current     string  `json:"current,omitempty"`
	Unit        string  `json:"unit,omitempty"`
	Direction   string  `json:"direction,omitempty"`
	DueDate     string  `json:"dueDate,omitempty"`  // ISO 8601 date
	Progress    float64 `json:"progress,omitempty"` // 0.0-1.0
	Status      string  `json:"status,omitempty"`
}

// Strategy is a choice about how to reach the goals, in words.
type Strategy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
	Priority    string `json:"priority,omitempty"` // P0, P1, P2, P3
	// GoalIDs are the goals the strategy serves.
	GoalIDs []string `json:"goalIds,omitempty"`
}

// Measure is a number that shows whether a strategy is working.
type Measure struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// StrategyID is the strategy the measure tracks.
	StrategyID string  `json:"strategyId"`
	Owner      string  `json:"owner,omitempty"`
	Baseline   string  `json:"baseline,omitempty"`
	Target     string  `json:"target,omitempty"`
	Current    string  `json:"current,omitempty"`
	Unit       string  `json:"unit,omitempty"`
	Direction  string  `json:"direction,omitempty"`
	DueDate    string  `json:"dueDate,omitempty"`   // ISO 8601 date
	Frequency  string  `json:"frequency,omitempty"` // how often it is measured, e.g. "Monthly"
	Source     string  `json:"source,omitempty"`    // where the data comes from
	Progress   float64 `json:"progress,omitempty"`  // 0.0-1.0
	Status     string  `json:"status,omitempty"`
}

// Action is work that carries out a strategy.
type Action struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// StrategyID is the strategy the action carries out.
	StrategyID string `json:"strategyId"`
	Owner      string `json:"owner,omitempty"`
	StartDate  string `json:"startDate,omitempty"` // ISO 8601 date
	EndDate    string `json:"endDate,omitempty"`   // ISO 8601 date
	Status     string `json:"status,omitempty"`
	// Dependencies are the IDs of actions this one waits on.
	Dependencies  []string          `json:"dependencies,omitempty"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"`
}

// New returns an OGSM with the given ID, name, and objective.
func New(id, name, objective string) *OGSM {
	return &OGSM{
		Metadata:  &Metadata{ID: id, Name: name, Status: StatusDraft},
		Objective: Objective{Statement: objective},
	}
}

// ReadFile reads an OGSM from a JSON file.
func ReadFile(filepath string) (*OGSM, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses OGSM JSON data.
func Parse(data []byte) (*OGSM, error) {
	var o OGSM
	if err := json.Unmarshal(data, &o); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &o, nil
}

// JSON returns the OGSM as formatted JSON.
func (o *OGSM) JSON() ([]byte, error) {
	return json.MarshalIndent(o, "", "  ")
}

// WriteFile writes the OGSM to a JSON file.
func (o *OGSM) WriteFile(filepath string) error {
	data, err := o.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// Goal returns the goal with the given ID, or nil.
func (o *OGSM) Goal(id string) *Goal {
	for i := range o.Goals {
		if o.Goals[i].ID == id {
			return &o.Goals[i]
		}
	}
	return nil
}

// Strategy returns the strategy with the given ID, or nil.
func (o *OGSM) Strategy(id string) *Strategy {
	for i := range o.Strategies {
		if o.Strategies[i].ID == id {
			return &o.Strategies[i]
		}
	}
	return nil
}

// MeasuresFor returns the measures that track the strategy, in document order.
func (o *OGSM) MeasuresFor(strategyID string) []Measure {
	var out []Measure
	for _, m := range o.Measures {
		if m.StrategyID == strategyID {
			out = append(out, m)
		}
	}
	return out
}

// ActionsFor returns the actions that carry out the strategy, in document
// order.
func (o *OGSM) ActionsFor(strategyID string) []Action {
	var out []Action
	for _, a := range o.Actions {
		if a.StrategyID == strategyID {
			out = append(out, a)
		}
	}
	return out
}

// StrategiesFor returns the strategies that serve the goal, in document order.
func (o *OGSM) StrategiesFor(goalID string) []Strategy {
	var out []Strategy
	for _, s := range o.Strategies {
		for _, id := range s.GoalIDs {
			if id == goalID {
				out = append(out, s)
				break
			}
		}
	}
	return out
}

// DocumentStatuses returns the document statuses.
func DocumentStatuses() []string {
	return []string{StatusDraft, StatusInReview, StatusApproved, StatusActive, StatusCompleted, StatusArchived}
}

// MetricStatuses returns the statuses of goals and measures.
func MetricStatuses() []string {
	return []string{MetricOnTrack, MetricAtRisk, MetricBehind, MetricAchieved, MetricMissed}
}

// ActionStatuses returns the action statuses.
func ActionStatuses() []string {
	return []string{ActionPlanned, ActionInProgress, ActionAtRisk, ActionCompleted, ActionCancelled}
}

// Priorities returns the priorities.
func Priorities() []string {
	return []string{PriorityP0, PriorityP1, PriorityP2, PriorityP3}
}

// Directions returns the directions of goals and measures.
func Directions() []string {
	return []string{DirectionIncrease, DirectionDecrease, DirectionMaintain}
}
