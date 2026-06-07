package canvas

// LeanUXCanvas follows Jeff Gothelf's v2 structure for Lean UX planning.
// Focuses on outcomes over outputs and hypothesis-driven development.
type LeanUXCanvas struct {
	Metadata Metadata `json:"metadata"`

	// Top row - Business context
	BusinessProblem  string    `json:"businessProblem"`  // What business problem are we solving?
	BusinessOutcomes []Outcome `json:"businessOutcomes"` // What business outcomes do we expect?

	// Middle row - User context
	Users        []LeanUXUser     `json:"users"`        // Who are we building for?
	UserOutcomes []Outcome        `json:"userOutcomes"` // What user outcomes do we expect?
	Solutions    []LeanUXSolution `json:"solutions"`    // What solutions might work?

	// Bottom row - Validation
	Hypotheses         []Hypothesis `json:"hypotheses"`
	RiskiestAssumption string       `json:"riskiestAssumption"` // Which assumption is riskiest?
	Experiment         *Experiment  `json:"experiment,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// LeanUXUser represents a user type in the Lean UX Canvas.
type LeanUXUser struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Description    string   `json:"description,omitempty"`
	Behaviors      []string `json:"behaviors,omitempty"`      // Observable behaviors
	NeedsGoals     []string `json:"needsGoals,omitempty"`     // Needs and goals
	CurrentJourney string   `json:"currentJourney,omitempty"` // Current experience
	PersonaRef     string   `json:"personaRef,omitempty"`     // Link to PRD persona
}

// Outcome represents a measurable outcome.
type Outcome struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Metric      string `json:"metric,omitempty"`    // How will we measure?
	Baseline    string `json:"baseline,omitempty"`  // Current value
	Target      string `json:"target,omitempty"`    // Target value
	Timeframe   string `json:"timeframe,omitempty"` // By when?
}

// LeanUXSolution represents a potential solution.
type LeanUXSolution struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Type        string   `json:"type,omitempty"`       // feature, change, experiment
	Effort      string   `json:"effort,omitempty"`     // small, medium, large
	Confidence  string   `json:"confidence,omitempty"` // high, medium, low
	Outcomes    []string `json:"outcomes,omitempty"`   // Outcome IDs addressed
}

// Hypothesis represents a testable hypothesis.
type Hypothesis struct {
	ID           string `json:"id"`
	WeBelieve    string `json:"weBelieve"`    // We believe that...
	WillResultIn string `json:"willResultIn"` // Will result in...
	WeWillKnow   string `json:"weWillKnow"`   // We will know we're right when...
	Validated    *bool  `json:"validated,omitempty"`
	Evidence     string `json:"evidence,omitempty"`
	SolutionRef  string `json:"solutionRef,omitempty"` // Solution being tested
}

// ExperimentStatus represents the status of an experiment.
type ExperimentStatus string

// Experiment status constants.
const (
	ExperimentStatusPlanned   ExperimentStatus = "planned"
	ExperimentStatusRunning   ExperimentStatus = "running"
	ExperimentStatusCompleted ExperimentStatus = "completed"
	ExperimentStatusCancelled ExperimentStatus = "cancelled"
)

// Experiment represents a validation experiment.
type Experiment struct {
	ID              string           `json:"id"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description"`
	Method          string           `json:"method"`             // prototype, A/B test, survey, interview, wizard-of-oz
	Duration        string           `json:"duration,omitempty"` // e.g., "2 weeks"
	SuccessCriteria string           `json:"successCriteria"`
	Results         string           `json:"results,omitempty"`
	Learnings       string           `json:"learnings,omitempty"`
	Status          ExperimentStatus `json:"status"`
	HypothesisRef   string           `json:"hypothesisRef,omitempty"` // Hypothesis being tested
}

// NewLeanUXCanvas creates a new LeanUXCanvas with defaults.
func NewLeanUXCanvas(id, title string) *LeanUXCanvas {
	return &LeanUXCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionLeanUX2,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *LeanUXCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// ValidatedHypotheses returns hypotheses that have been validated as true.
func (c *LeanUXCanvas) ValidatedHypotheses() []Hypothesis {
	var validated []Hypothesis
	for _, h := range c.Hypotheses {
		if h.Validated != nil && *h.Validated {
			validated = append(validated, h)
		}
	}
	return validated
}

// InvalidatedHypotheses returns hypotheses that have been invalidated.
func (c *LeanUXCanvas) InvalidatedHypotheses() []Hypothesis {
	var invalidated []Hypothesis
	for _, h := range c.Hypotheses {
		if h.Validated != nil && !*h.Validated {
			invalidated = append(invalidated, h)
		}
	}
	return invalidated
}

// UntestedHypotheses returns hypotheses that haven't been tested yet.
func (c *LeanUXCanvas) UntestedHypotheses() []Hypothesis {
	var untested []Hypothesis
	for _, h := range c.Hypotheses {
		if h.Validated == nil {
			untested = append(untested, h)
		}
	}
	return untested
}

// AllOutcomes returns both business and user outcomes.
func (c *LeanUXCanvas) AllOutcomes() []Outcome {
	all := make([]Outcome, 0, len(c.BusinessOutcomes)+len(c.UserOutcomes))
	all = append(all, c.BusinessOutcomes...)
	all = append(all, c.UserOutcomes...)
	return all
}
