package prd

import pf "github.com/grokify/priority-frameworks"

// MoSCoW represents the MoSCoW prioritization method.
// Now backed by priority-frameworks MoSCoW framework.
type MoSCoW string

// MoSCoW constants using priority-frameworks.
const (
	MoSCoWMust   MoSCoW = "must"
	MoSCoWShould MoSCoW = "should"
	MoSCoWCould  MoSCoW = "could"
	MoSCoWWont   MoSCoW = "wont"
)

// ValidMoSCoW checks if a MoSCoW value is valid.
func ValidMoSCoW(moscow MoSCoW) bool {
	if moscow == "" {
		return true
	}
	f := pf.MoSCoW()
	return f.IndexOf(string(moscow)) >= 0
}

// MoSCoWWeight returns a numeric weight for sorting MoSCoW priorities.
// Higher weight = higher priority.
func MoSCoWWeight(moscow MoSCoW) int {
	f := pf.MoSCoW()
	idx := f.IndexOf(string(moscow))
	if idx < 0 {
		return 0
	}
	return len(f.Levels) - idx
}

// MoSCoWFramework returns the MoSCoW priority framework.
func MoSCoWFramework() *pf.Framework {
	return pf.MoSCoW()
}

// UserStory represents a user story with acceptance criteria.
type UserStory struct {
	ID                 string                `json:"id"`
	PersonaID          string                `json:"personaId"` // Reference to persona
	Title              string                `json:"title"`
	AsA                string                `json:"asA"`    // Persona role (e.g., "developer", "admin")
	IWant              string                `json:"iWant"`  // Desired action/feature
	SoThat             string                `json:"soThat"` // Benefit/reason
	AcceptanceCriteria []AcceptanceCriterion `json:"acceptanceCriteria"`
	Priority           Priority              `json:"priority"`
	PhaseID            string                `json:"phaseId"` // Reference to roadmap phase
	StoryPoints        *int                  `json:"storyPoints,omitempty"`
	Dependencies       []string              `json:"dependencies,omitempty"` // Dependent story IDs
	Epic               string                `json:"epic,omitempty"`         // Parent epic
	Tags               []string              `json:"tags,omitempty"`         // For filtering by topic/domain
	Notes              string                `json:"notes,omitempty"`
}

// Story returns the full user story string in standard format.
func (us UserStory) Story() string {
	return "As a " + us.AsA + ", I want " + us.IWant + " so that " + us.SoThat
}

// AcceptanceCriterion defines a testable condition for a user story.
type AcceptanceCriterion struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Given       string `json:"given,omitempty"` // Precondition
	When        string `json:"when,omitempty"`  // Action
	Then        string `json:"then,omitempty"`  // Expected result
}
