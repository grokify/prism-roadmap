package journey

// OutcomeJourney tracks a business outcome's evolution over time.
// Outcomes describe why capability improvements matter - the user or business impact.
type OutcomeJourney struct {
	ID           string               `json:"id"`
	Name         string               `json:"name"`
	Description  string               `json:"description,omitempty"`
	Persona      string               `json:"persona,omitempty"` // Who experiences this outcome
	Category     OutcomeCategory      `json:"category,omitempty"`
	CurrentState *OutcomeState        `json:"currentState"`
	TargetStates []OutcomeTargetState `json:"targetStates"`
	Tags         []string             `json:"tags,omitempty"`
	Metadata     map[string]any       `json:"metadata,omitempty"`
}

// OutcomeCategory classifies the type of outcome.
type OutcomeCategory string

const (
	OutcomeCategoryEfficiency OutcomeCategory = "efficiency" // Time/cost savings
	OutcomeCategoryQuality    OutcomeCategory = "quality"    // Better outcomes
	OutcomeCategoryCapability OutcomeCategory = "capability" // New abilities
	OutcomeCategoryExperience OutcomeCategory = "experience" // User satisfaction
	OutcomeCategoryRevenue    OutcomeCategory = "revenue"    // Revenue impact
	OutcomeCategoryRisk       OutcomeCategory = "risk"       // Risk reduction
	OutcomeCategoryCompliance OutcomeCategory = "compliance" // Regulatory
)

// OutcomeState represents the current state of an outcome metric.
type OutcomeState struct {
	Value      float64 `json:"value"`
	Unit       string  `json:"unit"` // e.g., "developer-days", "percent", "dollars"
	Summary    string  `json:"summary,omitempty"`
	MeasuredAt string  `json:"measuredAt,omitempty"` // ISO 8601 date
}

// OutcomeTargetState represents a target state for an outcome.
type OutcomeTargetState struct {
	PeriodID           string   `json:"periodId"`
	Value              float64  `json:"value"`
	Unit               string   `json:"unit,omitempty"`
	Summary            string   `json:"summary,omitempty"`
	DriverCapabilities []string `json:"driverCapabilities,omitempty"` // Capability IDs that drive this
	DriverInitiatives  []string `json:"driverInitiatives,omitempty"`  // Initiative IDs
	Confidence         float64  `json:"confidence,omitempty"`         // 0.0-1.0
}

// OutcomeImpact calculates the improvement from current to target state.
type OutcomeImpact struct {
	OutcomeID      string  `json:"outcomeId"`
	FromPeriod     string  `json:"fromPeriod"`
	ToPeriod       string  `json:"toPeriod"`
	FromValue      float64 `json:"fromValue"`
	ToValue        float64 `json:"toValue"`
	Unit           string  `json:"unit"`
	PercentChange  float64 `json:"percentChange"`
	AbsoluteChange float64 `json:"absoluteChange"`
	Direction      string  `json:"direction"` // "improvement", "regression", "unchanged"
}

// CalculateImpact computes the impact between current and a target state.
func (oj *OutcomeJourney) CalculateImpact(targetPeriodID string) *OutcomeImpact {
	if oj.CurrentState == nil {
		return nil
	}

	var target *OutcomeTargetState
	for i := range oj.TargetStates {
		if oj.TargetStates[i].PeriodID == targetPeriodID {
			target = &oj.TargetStates[i]
			break
		}
	}

	if target == nil {
		return nil
	}

	fromValue := oj.CurrentState.Value
	toValue := target.Value
	absoluteChange := toValue - fromValue

	var percentChange float64
	if fromValue != 0 {
		percentChange = (absoluteChange / fromValue) * 100
	}

	direction := "unchanged"
	if absoluteChange > 0 {
		direction = "improvement"
	} else if absoluteChange < 0 {
		direction = "regression"
	}

	// For metrics where lower is better (e.g., time, cost), flip direction
	// This is a simplification; real implementation might use outcome category

	return &OutcomeImpact{
		OutcomeID:      oj.ID,
		FromPeriod:     oj.CurrentState.MeasuredAt,
		ToPeriod:       targetPeriodID,
		FromValue:      fromValue,
		ToValue:        toValue,
		Unit:           oj.CurrentState.Unit,
		PercentChange:  percentChange,
		AbsoluteChange: absoluteChange,
		Direction:      direction,
	}
}
