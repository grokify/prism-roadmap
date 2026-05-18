package prd

import core "github.com/grokify/prism-core"

// Note: Priority type and constants are now defined in common/
// and aliased in document.go for backward compatibility.

// MoSCoW represents the MoSCoW prioritization method.
type MoSCoW string

// MoSCoW constants imported from prism-core.
const (
	MoSCoWMust   MoSCoW = MoSCoW(core.MoSCoWMust)
	MoSCoWShould MoSCoW = MoSCoW(core.MoSCoWShould)
	MoSCoWCould  MoSCoW = MoSCoW(core.MoSCoWCould)
	MoSCoWWont   MoSCoW = MoSCoW(core.MoSCoWWont)
)

// ValidMoSCoW checks if a MoSCoW value is valid.
func ValidMoSCoW(moscow MoSCoW) bool {
	return core.ValidMoSCoW(string(moscow))
}

// MoSCoWWeight returns a numeric weight for sorting MoSCoW priorities.
func MoSCoWWeight(moscow MoSCoW) int {
	return core.MoSCoWWeight(string(moscow))
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
