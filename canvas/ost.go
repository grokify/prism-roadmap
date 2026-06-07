package canvas

// OpportunitySolutionTree follows Teresa Torres's structure for
// continuous discovery and outcome-driven product development.
// The tree flows: Outcome -> Opportunities -> Solutions -> Experiments
type OpportunitySolutionTree struct {
	Metadata Metadata   `json:"metadata"`
	Outcome  OSTOutcome `json:"outcome"` // Root node - the outcome we're targeting

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// OSTOutcome represents the desired outcome at the root of the tree.
type OSTOutcome struct {
	ID            string           `json:"id"`
	Description   string           `json:"description"`
	Metric        string           `json:"metric,omitempty"`    // How to measure
	Baseline      string           `json:"baseline,omitempty"`  // Current value
	Target        string           `json:"target,omitempty"`    // Target value
	Timeframe     string           `json:"timeframe,omitempty"` // By when
	Opportunities []OSTOpportunity `json:"opportunities"`

	// Link to PRD OKR
	OKRRef string `json:"okrRef,omitempty"`
}

// OSTOpportunity represents an opportunity discovered through research.
// Opportunities are customer problems, needs, or desires.
type OSTOpportunity struct {
	ID          string        `json:"id"`
	Description string        `json:"description"`
	Source      string        `json:"source,omitempty"`    // Interview, analytics, support, survey
	Frequency   string        `json:"frequency,omitempty"` // How often mentioned/observed
	Priority    int           `json:"priority,omitempty"`  // Ranking (lower = higher priority)
	Notes       string        `json:"notes,omitempty"`
	Solutions   []OSTSolution `json:"solutions"`
}

// OSTSolution represents a potential solution to an opportunity.
type OSTSolution struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Type        string          `json:"type,omitempty"`   // feature, improvement, experiment, quick-win
	Effort      string          `json:"effort,omitempty"` // small, medium, large
	Impact      string          `json:"impact,omitempty"` // high, medium, low
	Status      string          `json:"status,omitempty"` // proposed, testing, validated, building, shipped
	Experiments []OSTExperiment `json:"experiments"`

	// Link to PRD requirement
	RequirementRef string `json:"requirementRef,omitempty"`
}

// OSTExperiment represents an experiment to validate a solution.
type OSTExperiment struct {
	ID           string `json:"id"`
	Hypothesis   string `json:"hypothesis"`
	Method       string `json:"method"`                 // prototype, A/B test, survey, interview, fake-door
	Duration     string `json:"duration,omitempty"`     // e.g., "1 week"
	Participants string `json:"participants,omitempty"` // e.g., "10 users"
	Status       string `json:"status"`                 // planned, running, completed
	Result       string `json:"result,omitempty"`       // success, failure, inconclusive
	Learning     string `json:"learning,omitempty"`     // What did we learn?
	NextStep     string `json:"nextStep,omitempty"`     // What should we do next?
}

// NewOpportunitySolutionTree creates a new OST with defaults.
func NewOpportunitySolutionTree(id, title string) *OpportunitySolutionTree {
	return &OpportunitySolutionTree{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionOST1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *OpportunitySolutionTree) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// AllOpportunities returns all opportunities in the tree.
func (c *OpportunitySolutionTree) AllOpportunities() []OSTOpportunity {
	return c.Outcome.Opportunities
}

// AllSolutions returns all solutions across all opportunities.
func (c *OpportunitySolutionTree) AllSolutions() []OSTSolution {
	var solutions []OSTSolution
	for _, opp := range c.Outcome.Opportunities {
		solutions = append(solutions, opp.Solutions...)
	}
	return solutions
}

// AllExperiments returns all experiments across all solutions.
func (c *OpportunitySolutionTree) AllExperiments() []OSTExperiment {
	var experiments []OSTExperiment
	for _, opp := range c.Outcome.Opportunities {
		for _, sol := range opp.Solutions {
			experiments = append(experiments, sol.Experiments...)
		}
	}
	return experiments
}

// PrioritizedOpportunities returns opportunities sorted by priority (lower first).
func (c *OpportunitySolutionTree) PrioritizedOpportunities() []OSTOpportunity {
	// Create a copy to avoid modifying the original
	opps := make([]OSTOpportunity, len(c.Outcome.Opportunities))
	copy(opps, c.Outcome.Opportunities)

	// Simple bubble sort for small lists
	for i := 0; i < len(opps)-1; i++ {
		for j := 0; j < len(opps)-i-1; j++ {
			// Treat 0 as lowest priority (unprioritized)
			pi := opps[j].Priority
			pj := opps[j+1].Priority
			if pi == 0 {
				pi = 999999
			}
			if pj == 0 {
				pj = 999999
			}
			if pi > pj {
				opps[j], opps[j+1] = opps[j+1], opps[j]
			}
		}
	}
	return opps
}

// ValidatedSolutions returns solutions with status "validated".
func (c *OpportunitySolutionTree) ValidatedSolutions() []OSTSolution {
	var validated []OSTSolution
	for _, opp := range c.Outcome.Opportunities {
		for _, sol := range opp.Solutions {
			if sol.Status == "validated" {
				validated = append(validated, sol)
			}
		}
	}
	return validated
}

// ExperimentsByStatus returns experiments with the given status.
func (c *OpportunitySolutionTree) ExperimentsByStatus(status string) []OSTExperiment {
	var filtered []OSTExperiment
	for _, exp := range c.AllExperiments() {
		if exp.Status == status {
			filtered = append(filtered, exp)
		}
	}
	return filtered
}

// RunningExperiments returns experiments with status "running".
func (c *OpportunitySolutionTree) RunningExperiments() []OSTExperiment {
	return c.ExperimentsByStatus("running")
}

// CompletedExperiments returns experiments with status "completed".
func (c *OpportunitySolutionTree) CompletedExperiments() []OSTExperiment {
	return c.ExperimentsByStatus("completed")
}
