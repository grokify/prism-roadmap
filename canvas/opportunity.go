package canvas

// OpportunityCanvas follows Jeff Patton's 9-block structure for
// evaluating product opportunities before committing resources.
//
// Grid layout (BMC-style):
//
//	| Users & Customers | Problems        | Solution Ideas    |
//	| Solutions Today   | User Value      | Adoption Strategy |
//	| User Metrics      | Business Problem| Business Metrics  |
//	| Budget (colspan 3)                                      |
type OpportunityCanvas struct {
	Metadata Metadata `json:"metadata"`

	// Row 1: Problem Space
	Users            []User     `json:"users"`            // Who has the problem
	Problems         []Problem  `json:"problems"`         // Pains to address
	SolutionIdeas    []string   `json:"solutionIdeas"`    // Ways to solve it
	CurrentSolutions []Solution `json:"currentSolutions"` // Current workarounds (Row 2, Col 1)

	// Row 2: Solution Space
	UserValue        []string `json:"userValue"`        // Benefits to users
	AdoptionStrategy []string `json:"adoptionStrategy"` // How they'll find it

	// Row 3: Metrics & Business
	UserMetrics     []string `json:"userMetrics"`     // Behaviour to track
	BusinessProblem string   `json:"businessProblem"` // Why it matters to us
	BusinessMetrics []string `json:"businessMetrics"` // Outcome to measure

	// Row 4: Budget (full width)
	Budget *Budget `json:"budget,omitempty"` // What you're willing to invest to learn

	// Legacy fields (for backward compatibility with arrow-based view)
	ValueProposition ValueProp    `json:"valueProposition,omitempty"`
	BusinessValue    []string     `json:"businessValue,omitempty"` // Benefits to business
	Assumptions      []Assumption `json:"assumptions,omitempty"`
	Risks            []Risk       `json:"risks,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// Problem describes a user problem.
type Problem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Frequency   string `json:"frequency,omitempty"` // How often does it occur?
	Severity    string `json:"severity,omitempty"`  // How painful is it?
	Evidence    string `json:"evidence,omitempty"`  // What evidence supports this?
}

// User describes a user type affected by the problem.
type User struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"` // User type/persona name
	Description string   `json:"description,omitempty"`
	Goals       []string `json:"goals,omitempty"`      // What are they trying to achieve?
	Context     string   `json:"context,omitempty"`    // In what context?
	PersonaRef  string   `json:"personaRef,omitempty"` // Link to PRD persona
}

// Solution describes an existing solution to the problem.
type Solution struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Strengths   []string `json:"strengths,omitempty"`
	Weaknesses  []string `json:"weaknesses,omitempty"`
	Type        string   `json:"type,omitempty"` // competitor, workaround, internal, manual
}

// ValueProp describes the proposed value proposition.
type ValueProp struct {
	Statement      string   `json:"statement"`
	Differentiator string   `json:"differentiator,omitempty"` // What makes us different?
	KeyBenefits    []string `json:"keyBenefits,omitempty"`
}

// Assumption represents something assumed to be true that needs validation.
type Assumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Validated   bool   `json:"validated"`
	Evidence    string `json:"evidence,omitempty"`
	RiskLevel   string `json:"riskLevel,omitempty"` // high, medium, low
}

// Risk represents a risk to the opportunity.
type Risk struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`    // market, technical, business, regulatory
	Probability string `json:"probability,omitempty"` // high, medium, low
	Impact      string `json:"impact,omitempty"`      // high, medium, low
	Mitigation  string `json:"mitigation,omitempty"`
}

// Budget describes resource requirements.
type Budget struct {
	TimeEstimate    string `json:"timeEstimate,omitempty"`    // e.g., "3 months"
	TeamSize        string `json:"teamSize,omitempty"`        // e.g., "5 engineers"
	CostEstimate    string `json:"costEstimate,omitempty"`    // e.g., "$500K"
	ResourcesNeeded string `json:"resourcesNeeded,omitempty"` // Additional resources
	Constraints     string `json:"constraints,omitempty"`     // Known constraints
}

// NewOpportunityCanvas creates a new OpportunityCanvas with defaults.
func NewOpportunityCanvas(id, title string) *OpportunityCanvas {
	return &OpportunityCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionOpportunity1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *OpportunityCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// UnvalidatedAssumptions returns assumptions that haven't been validated.
func (c *OpportunityCanvas) UnvalidatedAssumptions() []Assumption {
	var unvalidated []Assumption
	for _, a := range c.Assumptions {
		if !a.Validated {
			unvalidated = append(unvalidated, a)
		}
	}
	return unvalidated
}

// HighRisks returns risks with high probability or high impact.
func (c *OpportunityCanvas) HighRisks() []Risk {
	var high []Risk
	for _, r := range c.Risks {
		if r.Probability == "high" || r.Impact == "high" {
			high = append(high, r)
		}
	}
	return high
}
