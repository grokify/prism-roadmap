// Package rmi provides the RoadmapItem (RMI) type with full prioritization support.
//
// A RoadmapItem represents a single item on a product roadmap, combining:
//   - MoSCoW prioritization (strategic priority)
//   - RICE scoring (quantitative prioritization)
//   - Market signals (customer demand)
//   - Effort estimation (implementation cost)
//   - Complexity factors (risk/dependency tracking)
//   - Cross-module references (capabilities, goals, ideas)
package rmi

import (
	"fmt"
	"time"

	"github.com/grokify/prism-roadmap/effort"
	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/grokify/prism-roadmap/signal"
)

// RoadmapItem represents a single item on the roadmap with full prioritization.
type RoadmapItem struct {
	// Identity
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Prioritization. MoSCoW is optional — empty means not yet prioritized
	// (e.g. items imported from an external PM tool before triage). When
	// set, it must be a valid prioritization.MoSCoWPriority value.
	MoSCoW       prioritization.MoSCoWPriority `json:"moscow,omitempty"`
	RICE         *prioritization.RICEScore     `json:"rice,omitempty"`
	MarketSignal *signal.MarketSignal          `json:"market_signal,omitempty"`

	// Effort & Complexity
	Effort     *effort.EffortEstimate    `json:"effort,omitempty"`
	Complexity *effort.ComplexityFactors `json:"complexity,omitempty"`

	// Links to other entities
	CapabilityRefs []string `json:"capability_refs,omitempty"` // prism-capability IDs
	GoalRefs       []string `json:"goal_refs,omitempty"`       // prism-maturity goal IDs
	IdeaRefs       []string `json:"idea_refs,omitempty"`       // ProductContext idea IDs

	// Timeline
	Quarter   string `json:"quarter,omitempty"`    // e.g., "Q3 2026"
	StartDate string `json:"start_date,omitempty"` // ISO 8601
	DueDate   string `json:"due_date,omitempty"`   // ISO 8601

	// Status
	Status   RMIStatus `json:"status"`
	Progress *int      `json:"progress,omitempty"` // 0-100 percentage

	// Metadata
	Owner     string    `json:"owner,omitempty"`
	Tags      []string  `json:"tags,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Notes     string    `json:"notes,omitempty"`
}

// RMIStatus represents the status of a roadmap item.
type RMIStatus string

const (
	RMIStatusPlanned    RMIStatus = "planned"
	RMIStatusInProgress RMIStatus = "in_progress"
	RMIStatusCompleted  RMIStatus = "completed"
	RMIStatusBlocked    RMIStatus = "blocked"
	RMIStatusCancelled  RMIStatus = "cancelled"
	RMIStatusDeferred   RMIStatus = "deferred"
)

// String returns the string representation.
func (s RMIStatus) String() string {
	return string(s)
}

// IsActive returns true if the item is in an active state.
func (s RMIStatus) IsActive() bool {
	return s == RMIStatusPlanned || s == RMIStatusInProgress || s == RMIStatusBlocked
}

// IsClosed returns true if the item is in a closed state.
func (s RMIStatus) IsClosed() bool {
	return s == RMIStatusCompleted || s == RMIStatusCancelled
}

// PriorityScore calculates a composite priority score.
// Higher score = higher priority.
// Formula: MoSCoW weight + (RICE score / 100) + (market signal score × 0.5) - complexity penalty
func (r *RoadmapItem) PriorityScore() float64 {
	score := float64(r.MoSCoW.Weight())

	if r.RICE != nil && r.RICE.Score > 0 {
		score += r.RICE.Score / 100
	}

	if r.MarketSignal != nil && r.MarketSignal.Score > 0 {
		score += r.MarketSignal.Score * 0.5
	}

	if r.Complexity != nil {
		score -= r.Complexity.ComplexityScore() * 0.2
	}

	return score
}

// EffectiveEffortDays returns the effort in person-days.
func (r *RoadmapItem) EffectiveEffortDays() int {
	if r.Effort == nil {
		return 0
	}
	return r.Effort.EffectivePersonDays()
}

// HasBlockingDependencies returns true if there are unresolved blocking dependencies.
func (r *RoadmapItem) HasBlockingDependencies() bool {
	if r.Complexity == nil {
		return false
	}
	for _, d := range r.Complexity.Dependencies {
		if d.IsBlocking() {
			return true
		}
	}
	return false
}

// IsActionable returns true if the item should be worked on.
func (r *RoadmapItem) IsActionable() bool {
	return r.MoSCoW.IsActionable() && r.Status.IsActive() && !r.HasBlockingDependencies()
}

// Validate returns an error if the item is invalid.
func (r *RoadmapItem) Validate() error {
	if r.ID == "" {
		return fmt.Errorf("id is required")
	}
	if r.Name == "" {
		return fmt.Errorf("name is required")
	}
	if r.MoSCoW != "" && !prioritization.IsValidMoSCoWPriority(r.MoSCoW) {
		return fmt.Errorf("invalid moscow priority: %s", r.MoSCoW)
	}
	if r.RICE != nil {
		if err := r.RICE.Validate(); err != nil {
			return fmt.Errorf("rice: %w", err)
		}
	}
	if r.Effort != nil {
		if err := r.Effort.Validate(); err != nil {
			return fmt.Errorf("effort: %w", err)
		}
	}
	return nil
}

// NewRoadmapItem creates a new roadmap item with required fields.
func NewRoadmapItem(id, name string, moscow prioritization.MoSCoWPriority) *RoadmapItem {
	now := time.Now()
	return &RoadmapItem{
		ID:        id,
		Name:      name,
		MoSCoW:    moscow,
		Status:    RMIStatusPlanned,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// WithRICE adds RICE scoring to the item.
func (r *RoadmapItem) WithRICE(rice *prioritization.RICEScore) *RoadmapItem {
	r.RICE = rice
	r.UpdatedAt = time.Now()
	return r
}

// WithMarketSignal adds market signal to the item.
func (r *RoadmapItem) WithMarketSignal(ms *signal.MarketSignal) *RoadmapItem {
	r.MarketSignal = ms
	r.UpdatedAt = time.Now()
	return r
}

// WithEffort adds effort estimate to the item.
func (r *RoadmapItem) WithEffort(e *effort.EffortEstimate) *RoadmapItem {
	r.Effort = e
	r.UpdatedAt = time.Now()
	return r
}

// WithComplexity adds complexity factors to the item.
func (r *RoadmapItem) WithComplexity(c *effort.ComplexityFactors) *RoadmapItem {
	r.Complexity = c
	r.UpdatedAt = time.Now()
	return r
}

// AddCapabilityRef adds a capability reference.
func (r *RoadmapItem) AddCapabilityRef(capID string) {
	r.CapabilityRefs = append(r.CapabilityRefs, capID)
	r.UpdatedAt = time.Now()
}

// AddGoalRef adds a goal reference.
func (r *RoadmapItem) AddGoalRef(goalID string) {
	r.GoalRefs = append(r.GoalRefs, goalID)
	r.UpdatedAt = time.Now()
}

// AddIdeaRef adds an idea reference.
func (r *RoadmapItem) AddIdeaRef(ideaID string) {
	r.IdeaRefs = append(r.IdeaRefs, ideaID)
	r.UpdatedAt = time.Now()
}

// ValidRMIStatuses returns all valid RMI statuses.
func ValidRMIStatuses() []RMIStatus {
	return []RMIStatus{
		RMIStatusPlanned,
		RMIStatusInProgress,
		RMIStatusCompleted,
		RMIStatusBlocked,
		RMIStatusCancelled,
		RMIStatusDeferred,
	}
}
