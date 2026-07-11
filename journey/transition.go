package journey

// CapabilityTransition is the core abstraction of the journey model.
// It represents a capability moving from one maturity state to another,
// enabled by initiatives, producing outcomes, with required evidence.
//
// The transition creates a deterministic chain:
// Current state → transition → initiative → capability gain → outcome → evidence
type CapabilityTransition struct {
	ID               string           `json:"id,omitempty"`
	Type             string           `json:"type,omitempty"` // "capability_transition"
	CapabilityID     string           `json:"capabilityId"`
	CapabilityName   string           `json:"capabilityName,omitempty"`
	From             TransitionState  `json:"from"`
	To               TransitionState  `json:"to"`
	EnabledBy        []string         `json:"enabledBy,omitempty"`        // Initiative IDs
	Produces         []string         `json:"produces,omitempty"`         // Outcome IDs
	EvidenceRequired []string         `json:"evidenceRequired,omitempty"` // Evidence to confirm transition
	Status           TransitionStatus `json:"status,omitempty"`
	Owner            string           `json:"owner,omitempty"` // Team or person
	Description      string           `json:"description,omitempty"`
}

// TransitionState represents a capability's state at one end of a transition.
type TransitionState struct {
	Period   string `json:"period"`   // Period ID: "now", "2026-q3"
	Maturity string `json:"maturity"` // Maturity level: "M1", "M2", etc.
	Summary  string `json:"summary,omitempty"`
}

// TransitionStatus tracks the progress of a transition.
type TransitionStatus string

const (
	TransitionStatusPlanned    TransitionStatus = "planned"
	TransitionStatusInProgress TransitionStatus = "in_progress"
	TransitionStatusCompleted  TransitionStatus = "completed"
	TransitionStatusBlocked    TransitionStatus = "blocked"
	TransitionStatusAtRisk     TransitionStatus = "at_risk"
	TransitionStatusCancelled  TransitionStatus = "cancelled"
)

// TransitionChain represents a sequence of transitions for a capability.
type TransitionChain struct {
	CapabilityID  string                 `json:"capabilityId"`
	Transitions   []CapabilityTransition `json:"transitions"`
	StartMaturity string                 `json:"startMaturity"` // Initial maturity
	EndMaturity   string                 `json:"endMaturity"`   // Final maturity
	TotalPeriods  int                    `json:"totalPeriods"`
}

// BuildTransitionChain creates a chain of transitions from a capability journey.
func BuildTransitionChain(journey CapabilityJourney) *TransitionChain {
	if journey.CurrentState == nil || len(journey.TargetStates) == 0 {
		return nil
	}

	chain := &TransitionChain{
		CapabilityID:  journey.CapabilityID,
		StartMaturity: journey.CurrentState.MaturityLevel,
		Transitions:   make([]CapabilityTransition, 0, len(journey.TargetStates)),
	}

	prevState := journey.CurrentState

	for _, target := range journey.TargetStates {
		// Only create transition if maturity changes
		if target.MaturityLevel != prevState.MaturityLevel {
			transition := CapabilityTransition{
				CapabilityID:   journey.CapabilityID,
				CapabilityName: journey.Name,
				From: TransitionState{
					Period:   prevState.PeriodID,
					Maturity: prevState.MaturityLevel,
					Summary:  prevState.Summary,
				},
				To: TransitionState{
					Period:   target.PeriodID,
					Maturity: target.MaturityLevel,
					Summary:  target.Summary,
				},
				EnabledBy: target.Initiatives,
				Owner:     journey.Owner,
			}
			chain.Transitions = append(chain.Transitions, transition)
		}

		// Update prev state for next iteration
		prevState = &MaturityState{
			PeriodID:      target.PeriodID,
			MaturityLevel: target.MaturityLevel,
			Summary:       target.Summary,
		}
	}

	if len(chain.Transitions) > 0 {
		chain.EndMaturity = chain.Transitions[len(chain.Transitions)-1].To.Maturity
		chain.TotalPeriods = len(journey.TargetStates)
	}

	return chain
}

// MaturityDelta calculates the maturity level change in a transition.
// Returns positive for improvement, negative for regression.
func (t *CapabilityTransition) MaturityDelta() int {
	fromLevel := parseMaturityLevel(t.From.Maturity)
	toLevel := parseMaturityLevel(t.To.Maturity)
	return toLevel - fromLevel
}

// parseMaturityLevel extracts the numeric level from "M1", "M2", etc.
func parseMaturityLevel(level string) int {
	if len(level) < 2 || level[0] != 'M' {
		return 0
	}
	n := int(level[1] - '0')
	if n < 0 || n > 9 {
		return 0
	}
	return n
}
