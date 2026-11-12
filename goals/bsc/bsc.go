// Package bsc provides types and utilities for Balanced Scorecards.
//
// A Balanced Scorecard describes a strategy as objectives in a few
// perspectives, balanced so that the results an organization reports are not
// pursued at the expense of what produces them:
//
//   - Financial: how the organization succeeds for those who fund it.
//   - Customer: how it must appear to the people it serves.
//   - Internal process: which processes it must excel at.
//   - Learning and growth: which people, systems, and culture it must build.
//
// A strategy map links objectives from cause to effect, usually upward from
// learning and growth to the financial results. Each objective has measures,
// lead or lag, with targets, and initiatives close the gap between where a
// measure is and where it should be. Strategic themes group objectives that
// cross the perspectives.
//
// Perspectives are listed top first and may be renamed, reordered, or
// replaced: a public agency or nonprofit often puts its mission on top, with
// the financial perspective as a constraint below it.
package bsc

import (
	"encoding/json"
	"fmt"
	"os"
	"slices"
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

// Perspective kind constants. A kind says which classic perspective a
// perspective plays, whatever it is called.
const (
	KindFinancial       = "financial"
	KindCustomer        = "customer"
	KindInternalProcess = "internal_process"
	KindLearningGrowth  = "learning_growth"
	// KindMission is a mission or stakeholder perspective, which public and
	// nonprofit scorecards put on top.
	KindMission = "mission"
	KindCustom  = "custom"
)

// Measure type constants.
const (
	// TypeLag is an outcome measure: it reports a result after the fact.
	TypeLag = "lag"
	// TypeLead is a performance driver: it moves before the results do.
	TypeLead = "lead"
)

// Status constants for measures.
const (
	MeasureOnTrack  = "On Track"
	MeasureAtRisk   = "At Risk"
	MeasureBehind   = "Behind"
	MeasureAchieved = "Achieved"
	MeasureMissed   = "Missed"
)

// Initiative status constants.
const (
	InitiativePlanned    = "Planned"
	InitiativeInProgress = "In Progress"
	InitiativeAtRisk     = "At Risk"
	InitiativeCompleted  = "Completed"
	InitiativeCancelled  = "Cancelled"
)

// Direction constants say which way a measure should move.
const (
	DirectionIncrease = "increase"
	DirectionDecrease = "decrease"
	DirectionMaintain = "maintain"
)

// BSC is a complete Balanced Scorecard.
type BSC struct {
	Schema   string    `json:"$schema,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	// Mission is why the organization exists.
	Mission string `json:"mission,omitempty"`
	// Vision is what it wants to become.
	Vision string `json:"vision,omitempty"`
	// Strategy is how it will get there, in a sentence or two.
	Strategy string `json:"strategy,omitempty"`
	// Perspectives are listed top first.
	Perspectives []Perspective `json:"perspectives"`
	Themes       []Theme       `json:"themes,omitempty"`
	Objectives   []Objective   `json:"objectives"`
	Measures     []Measure     `json:"measures,omitempty"`
	Initiatives  []Initiative  `json:"initiatives,omitempty"`
}

// Metadata contains document metadata.
type Metadata struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Author  string `json:"author,omitempty"`
	Team    string `json:"team,omitempty"`
	Period  string `json:"period,omitempty"` // e.g. "FY2027"
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
	// ParentID is the scorecard this one cascades from.
	ParentID  string    `json:"parentId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}

// Perspective is one view of the strategy, such as the customer's.
type Perspective struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Kind is the classic perspective this one plays: financial, customer,
	// internal_process, learning_growth, mission, or custom.
	Kind string `json:"kind,omitempty"`
	// Question is the question the perspective answers, such as "To achieve
	// our vision, how should we appear to our customers?"
	Question    string `json:"question,omitempty"`
	Description string `json:"description,omitempty"`
}

// Theme is a strategic theme: a group of objectives, usually across the
// perspectives, that together carry out one part of the strategy.
type Theme struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

// Objective is a strategic objective in one perspective.
type Objective struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// PerspectiveID is the perspective the objective sits in.
	PerspectiveID string `json:"perspectiveId"`
	// ThemeIDs are the strategic themes the objective belongs to.
	ThemeIDs []string `json:"themeIds,omitempty"`
	Owner    string   `json:"owner,omitempty"`
	Priority string   `json:"priority,omitempty"` // P0, P1, P2, P3
	// ContributesTo are the objectives this one helps achieve: the arrows of
	// the strategy map, from cause to effect.
	ContributesTo []string `json:"contributesTo,omitempty"`
}

// Measure is how progress on an objective is measured.
type Measure struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// ObjectiveID is the objective the measure tracks.
	ObjectiveID string `json:"objectiveId"`
	// Type is lag, a result, or lead, a driver of results.
	Type      string `json:"type,omitempty"`
	Owner     string `json:"owner,omitempty"`
	Baseline  string `json:"baseline,omitempty"`
	Target    string `json:"target,omitempty"`
	Current   string `json:"current,omitempty"`
	Unit      string `json:"unit,omitempty"`
	Direction string `json:"direction,omitempty"`
	DueDate   string `json:"dueDate,omitempty"`   // ISO 8601 date
	Frequency string `json:"frequency,omitempty"` // how often it is measured, e.g. "Monthly"
	Source    string `json:"source,omitempty"`    // where the data comes from
	// Thresholds are the stoplight bands around the target.
	Thresholds *Thresholds `json:"thresholds,omitempty"`
	Progress   float64     `json:"progress,omitempty"` // 0.0-1.0
	Status     string      `json:"status,omitempty"`
}

// Thresholds are the values at which a measure turns green or red, in its
// unit. Between them it is yellow.
type Thresholds struct {
	Green string `json:"green,omitempty"`
	Red   string `json:"red,omitempty"`
}

// Initiative is work that closes the gap between where measures are and their
// targets.
type Initiative struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// ObjectiveIDs are the objectives the initiative serves.
	ObjectiveIDs []string `json:"objectiveIds"`
	Owner        string   `json:"owner,omitempty"`
	StartDate    string   `json:"startDate,omitempty"` // ISO 8601 date
	EndDate      string   `json:"endDate,omitempty"`   // ISO 8601 date
	Budget       string   `json:"budget,omitempty"`
	Status       string   `json:"status,omitempty"`
	// Dependencies are the IDs of initiatives this one waits on.
	Dependencies  []string          `json:"dependencies,omitempty"`
	ExternalLinks map[string]string `json:"externalLinks,omitempty"`
}

// DefaultPerspectives returns the four classic perspectives, top first.
func DefaultPerspectives() []Perspective {
	return []Perspective{
		{ID: "financial", Name: "Financial", Kind: KindFinancial,
			Question: "To succeed financially, how should we appear to those who fund us?"},
		{ID: "customer", Name: "Customer", Kind: KindCustomer,
			Question: "To achieve our vision, how should we appear to our customers?"},
		{ID: "internal-process", Name: "Internal Process", Kind: KindInternalProcess,
			Question: "To satisfy our customers, which processes must we excel at?"},
		{ID: "learning-growth", Name: "Learning and Growth", Kind: KindLearningGrowth,
			Question: "To excel at those processes, how will we sustain our ability to change and improve?"},
	}
}

// New returns a Balanced Scorecard with the given ID and name and the four
// classic perspectives.
func New(id, name string) *BSC {
	return &BSC{
		Metadata:     &Metadata{ID: id, Name: name, Status: StatusDraft},
		Perspectives: DefaultPerspectives(),
	}
}

// ReadFile reads a Balanced Scorecard from a JSON file.
func ReadFile(filepath string) (*BSC, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses Balanced Scorecard JSON data.
func Parse(data []byte) (*BSC, error) {
	var b BSC
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &b, nil
}

// JSON returns the scorecard as formatted JSON.
func (b *BSC) JSON() ([]byte, error) {
	return json.MarshalIndent(b, "", "  ")
}

// WriteFile writes the scorecard to a JSON file.
func (b *BSC) WriteFile(filepath string) error {
	data, err := b.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// Perspective returns the perspective with the given ID, or nil.
func (b *BSC) Perspective(id string) *Perspective {
	for i := range b.Perspectives {
		if b.Perspectives[i].ID == id {
			return &b.Perspectives[i]
		}
	}
	return nil
}

// Theme returns the theme with the given ID, or nil.
func (b *BSC) Theme(id string) *Theme {
	for i := range b.Themes {
		if b.Themes[i].ID == id {
			return &b.Themes[i]
		}
	}
	return nil
}

// Objective returns the objective with the given ID, or nil.
func (b *BSC) Objective(id string) *Objective {
	for i := range b.Objectives {
		if b.Objectives[i].ID == id {
			return &b.Objectives[i]
		}
	}
	return nil
}

// ObjectivesIn returns the objectives in the perspective, in document order.
func (b *BSC) ObjectivesIn(perspectiveID string) []Objective {
	var out []Objective
	for _, o := range b.Objectives {
		if o.PerspectiveID == perspectiveID {
			out = append(out, o)
		}
	}
	return out
}

// ObjectivesInTheme returns the objectives in the theme, in document order.
func (b *BSC) ObjectivesInTheme(themeID string) []Objective {
	var out []Objective
	for _, o := range b.Objectives {
		if slices.Contains(o.ThemeIDs, themeID) {
			out = append(out, o)
		}
	}
	return out
}

// ContributorsTo returns the objectives that contribute to the objective, in
// document order.
func (b *BSC) ContributorsTo(objectiveID string) []Objective {
	var out []Objective
	for _, o := range b.Objectives {
		if slices.Contains(o.ContributesTo, objectiveID) {
			out = append(out, o)
		}
	}
	return out
}

// MeasuresFor returns the measures that track the objective, in document
// order.
func (b *BSC) MeasuresFor(objectiveID string) []Measure {
	var out []Measure
	for _, m := range b.Measures {
		if m.ObjectiveID == objectiveID {
			out = append(out, m)
		}
	}
	return out
}

// InitiativesFor returns the initiatives that serve the objective, in
// document order.
func (b *BSC) InitiativesFor(objectiveID string) []Initiative {
	var out []Initiative
	for _, in := range b.Initiatives {
		if slices.Contains(in.ObjectiveIDs, objectiveID) {
			out = append(out, in)
		}
	}
	return out
}

// PerspectiveIndex returns the position of the perspective, top first, or -1.
func (b *BSC) PerspectiveIndex(id string) int {
	return slices.IndexFunc(b.Perspectives, func(p Perspective) bool { return p.ID == id })
}

// DocumentStatuses returns the document statuses.
func DocumentStatuses() []string {
	return []string{StatusDraft, StatusInReview, StatusApproved, StatusActive, StatusCompleted, StatusArchived}
}

// PerspectiveKinds returns the perspective kinds.
func PerspectiveKinds() []string {
	return []string{KindFinancial, KindCustomer, KindInternalProcess, KindLearningGrowth, KindMission, KindCustom}
}

// MeasureTypes returns the measure types.
func MeasureTypes() []string {
	return []string{TypeLag, TypeLead}
}

// MeasureStatuses returns the measure statuses.
func MeasureStatuses() []string {
	return []string{MeasureOnTrack, MeasureAtRisk, MeasureBehind, MeasureAchieved, MeasureMissed}
}

// InitiativeStatuses returns the initiative statuses.
func InitiativeStatuses() []string {
	return []string{InitiativePlanned, InitiativeInProgress, InitiativeAtRisk, InitiativeCompleted, InitiativeCancelled}
}

// Priorities returns the priorities.
func Priorities() []string {
	return []string{PriorityP0, PriorityP1, PriorityP2, PriorityP3}
}

// Directions returns the directions of measures.
func Directions() []string {
	return []string{DirectionIncrease, DirectionDecrease, DirectionMaintain}
}
