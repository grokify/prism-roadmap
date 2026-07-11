// Package journey provides types for capability evolution roadmaps.
// A journey roadmap models how capabilities, outcomes, and business value
// evolve over time through transitions enabled by initiatives.
package journey

// JourneyRoadmap is the top-level container for a capability evolution roadmap.
// It combines capability journeys, outcome journeys, initiatives, and narrative
// into a cohesive strategic planning document.
type JourneyRoadmap struct {
	ID                 string              `json:"id"`
	Type               string              `json:"type,omitempty"` // "journey_roadmap"
	Name               string              `json:"name"`
	Vision             string              `json:"vision,omitempty"`
	Description        string              `json:"description,omitempty"`
	Scope              *Scope              `json:"scope,omitempty"`
	TimeModel          *TimeModel          `json:"timeModel,omitempty"`
	CapabilityJourneys []CapabilityJourney `json:"capabilityJourneys,omitempty"`
	OutcomeJourneys    []OutcomeJourney    `json:"outcomeJourneys,omitempty"`
	Initiatives        []Initiative        `json:"initiatives,omitempty"`
	Dependencies       []Dependency        `json:"dependencies,omitempty"`
	Teams              []Team              `json:"teams,omitempty"`
	Risks              []Risk              `json:"risks,omitempty"`
	Narrative          *RoadmapNarrative   `json:"narrative,omitempty"`
	Scenarios          []Scenario          `json:"scenarios,omitempty"`
	Metadata           map[string]any      `json:"metadata,omitempty"`
}

// Scope defines what the roadmap covers.
type Scope struct {
	Type        string   `json:"type,omitempty"` // "platform", "product", "team", "organization"
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// TimeModel defines the time structure for the roadmap.
type TimeModel struct {
	Type       TimeModelType `json:"type"`                // quarterly, monthly, sprint, milestone
	StartDate  string        `json:"startDate,omitempty"` // ISO 8601 date
	EndDate    string        `json:"endDate,omitempty"`   // ISO 8601 date
	Periods    []Period      `json:"periods"`
	FiscalYear string        `json:"fiscalYear,omitempty"` // e.g., "FY2026"
}

// TimeModelType represents the type of time periods used.
type TimeModelType string

const (
	TimeModelQuarterly TimeModelType = "quarterly"
	TimeModelMonthly   TimeModelType = "monthly"
	TimeModelSprint    TimeModelType = "sprint"
	TimeModelMilestone TimeModelType = "milestone"
	TimeModelHalf      TimeModelType = "half" // H1, H2
	TimeModelCustom    TimeModelType = "custom"
)

// Period represents a time period in the roadmap.
type Period struct {
	ID          string `json:"id"`                  // e.g., "now", "2026-q3"
	Label       string `json:"label"`               // e.g., "Q3 2026", "Current State"
	StartDate   string `json:"startDate,omitempty"` // ISO 8601 date
	EndDate     string `json:"endDate,omitempty"`   // ISO 8601 date
	IsCurrent   bool   `json:"isCurrent,omitempty"` // True for "now" period
	Description string `json:"description,omitempty"`
}

// Scenario represents an alternative planning scenario.
type Scenario struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsBase      bool   `json:"isBase,omitempty"` // True for the default scenario
}

// Risk represents a risk associated with the roadmap.
type Risk struct {
	ID              string   `json:"id"`
	Description     string   `json:"description"`
	Probability     string   `json:"probability,omitempty"` // low, medium, high
	Impact          string   `json:"impact,omitempty"`      // low, medium, high, critical
	Mitigation      string   `json:"mitigation,omitempty"`
	Status          string   `json:"status,omitempty"`          // identified, mitigating, resolved, accepted
	AffectedPeriods []string `json:"affectedPeriods,omitempty"` // Period IDs
	Tags            []string `json:"tags,omitempty"`
}
