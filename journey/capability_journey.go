package journey

// CapabilityJourney tracks a single capability's evolution through maturity levels
// over time. Each target state represents a planned maturity improvement.
type CapabilityJourney struct {
	ID              string         `json:"id"`
	CapabilityID    string         `json:"capabilityId"` // Reference to capability definition
	Name            string         `json:"name"`
	Description     string         `json:"description,omitempty"`
	Owner           string         `json:"owner,omitempty"` // Team or person ID
	CurrentState    *MaturityState `json:"currentState"`
	TargetStates    []TargetState  `json:"targetStates"`
	DesiredEndState *MaturityState `json:"desiredEndState,omitempty"`
	Tags            []string       `json:"tags,omitempty"`
	Metadata        map[string]any `json:"metadata,omitempty"`
}

// MaturityState represents a capability's state at a point in time.
type MaturityState struct {
	PeriodID      string   `json:"periodId"`
	MaturityLevel string   `json:"maturityLevel"` // e.g., "M1", "M2", "M3", "M4", "M5"
	Summary       string   `json:"summary,omitempty"`
	Evidence      []string `json:"evidence,omitempty"` // Evidence that this state was achieved
}

// TargetState represents a planned future state for a capability.
type TargetState struct {
	PeriodID        string           `json:"periodId"`
	MaturityLevel   string           `json:"maturityLevel"`
	Summary         string           `json:"summary,omitempty"`
	Changes         []string         `json:"changes,omitempty"`     // What changes in this transition
	Initiatives     []string         `json:"initiatives,omitempty"` // Initiative IDs that enable this
	SuccessMeasures []SuccessMeasure `json:"successMeasures,omitempty"`
	Confidence      float64          `json:"confidence,omitempty"` // 0.0-1.0
	Commitment      CommitmentLevel  `json:"commitment,omitempty"`
	Assumptions     []string         `json:"assumptions,omitempty"`
	ScenarioID      string           `json:"scenarioId,omitempty"` // For scenario-specific states
}

// SuccessMeasure defines how to measure success for a target state.
type SuccessMeasure struct {
	Metric      string  `json:"metric"`
	Target      float64 `json:"target"`
	Unit        string  `json:"unit,omitempty"`
	Description string  `json:"description,omitempty"`
}

// CommitmentLevel indicates how committed the organization is to a target state.
type CommitmentLevel string

const (
	CommitmentCommitted CommitmentLevel = "committed" // Locked in, resourced
	CommitmentPlanned   CommitmentLevel = "planned"   // Agreed, pending resources
	CommitmentTargeted  CommitmentLevel = "targeted"  // Goal, not yet planned
	CommitmentAspirant  CommitmentLevel = "aspirant"  // Stretch goal
)

// MaturityLevel constants for common maturity models.
const (
	MaturityM0 = "M0" // Initial/Ad-hoc
	MaturityM1 = "M1" // Developing/Repeatable
	MaturityM2 = "M2" // Defined
	MaturityM3 = "M3" // Managed/Measured
	MaturityM4 = "M4" // Optimizing
	MaturityM5 = "M5" // Innovating/Leading
)
