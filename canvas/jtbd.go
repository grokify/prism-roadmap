package canvas

// JTBDCanvas implements the Jobs-to-be-Done methodology (Clayton Christensen, Tony Ulwick).
// It focuses on understanding the jobs customers are trying to accomplish and the outcomes
// they expect from completing those jobs.
type JTBDCanvas struct {
	Metadata Metadata `json:"metadata"`

	// Core job definition
	MainJob     *JobStatement  `json:"mainJob"`               // The primary job to be done
	RelatedJobs []JobStatement `json:"relatedJobs,omitempty"` // Related jobs in the job hierarchy

	// Job map - stages of getting the job done
	JobMap []JobStage `json:"jobMap,omitempty"`

	// Outcome expectations
	DesiredOutcomes   []OutcomeExpectation `json:"desiredOutcomes,omitempty"`   // What success looks like
	UndesiredOutcomes []OutcomeExpectation `json:"undesiredOutcomes,omitempty"` // What to avoid

	// Circumstances and context
	Circumstances []JobCircumstance `json:"circumstances,omitempty"` // When/where the job arises

	// Current solutions and alternatives
	HiringSolutions    []HiringSolution `json:"hiringSolutions,omitempty"`    // Current solutions "hired"
	FiringSolutions    []string         `json:"firingSolutions,omitempty"`    // Solutions being "fired"
	CompetingSolutions []string         `json:"competingSolutions,omitempty"` // Alternative approaches

	// Forces analysis (switching behavior)
	PushForces []Force `json:"pushForces,omitempty"` // Forces pushing away from current solution
	PullForces []Force `json:"pullForces,omitempty"` // Forces pulling toward new solution
	Anxieties  []Force `json:"anxieties,omitempty"`  // Concerns about switching
	Habits     []Force `json:"habits,omitempty"`     // Comfort with current solution

	// Opportunity scoring (ODI - Outcome-Driven Innovation)
	OpportunityScores []OpportunityScore `json:"opportunityScores,omitempty"`

	// Progress and validation
	Interviews      []JTBDInterview `json:"interviews,omitempty"`
	ValidationNotes string          `json:"validationNotes,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// JobStatement represents a job-to-be-done in standard JTBD format.
// Format: [verb] + [object of verb] + [contextual clarifier]
type JobStatement struct {
	ID           string  `json:"id"`
	Statement    string  `json:"statement"`              // The job statement
	Type         JobType `json:"type"`                   // functional, emotional, social, consumption
	Importance   string  `json:"importance,omitempty"`   // high, medium, low
	Satisfaction string  `json:"satisfaction,omitempty"` // current satisfaction level
	Frequency    string  `json:"frequency,omitempty"`    // how often the job arises
	Context      string  `json:"context,omitempty"`      // situational context

	// Job hierarchy
	ParentJobID string   `json:"parentJobId,omitempty"` // For sub-jobs
	ChildJobIDs []string `json:"childJobIds,omitempty"` // For parent jobs
}

// JobType represents the type of job.
type JobType string

// Job type constants.
const (
	JobTypeFunctional  JobType = "functional"  // Core task to accomplish
	JobTypeEmotional   JobType = "emotional"   // How they want to feel
	JobTypeSocial      JobType = "social"      // How they want to be perceived
	JobTypeConsumption JobType = "consumption" // Using/consuming the solution
)

// JobStage represents a stage in the job map (Ulwick's Universal Job Map).
type JobStage struct {
	ID          string   `json:"id"`
	Stage       string   `json:"stage"` // define, locate, prepare, confirm, execute, monitor, modify, conclude
	Name        string   `json:"name"`  // Custom stage name
	Description string   `json:"description,omitempty"`
	Steps       []string `json:"steps,omitempty"`      // Steps within this stage
	PainPoints  []string `json:"painPoints,omitempty"` // Problems at this stage
	Outcomes    []string `json:"outcomes,omitempty"`   // Desired outcomes for this stage
}

// OutcomeExpectation represents a desired (or undesired) outcome.
// Format: [direction] + [unit of measure] + [object of measure] + [contextual clarifier]
type OutcomeExpectation struct {
	ID           string  `json:"id"`
	Statement    string  `json:"statement"`              // The outcome statement
	Direction    string  `json:"direction"`              // minimize, maximize, increase, decrease
	Metric       string  `json:"metric,omitempty"`       // Unit of measure
	Importance   int     `json:"importance,omitempty"`   // 1-10 scale
	Satisfaction int     `json:"satisfaction,omitempty"` // 1-10 scale (current)
	Opportunity  float64 `json:"opportunity,omitempty"`  // Calculated opportunity score
	JobStageRef  string  `json:"jobStageRef,omitempty"`  // Which stage this relates to
}

// JobCircumstance represents the context in which a job arises.
type JobCircumstance struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`           // When/where/why the job arises
	Trigger     string   `json:"trigger,omitempty"`     // What triggers the job
	Frequency   string   `json:"frequency,omitempty"`   // How often
	Urgency     string   `json:"urgency,omitempty"`     // How urgent when it arises
	Constraints []string `json:"constraints,omitempty"` // Limitations in this context
}

// HiringSolution represents a solution currently "hired" to do the job.
type HiringSolution struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Type              string   `json:"type,omitempty"`              // product, service, workaround, DIY
	WhyHired          string   `json:"whyHired,omitempty"`          // Why this solution was chosen
	HowUsed           string   `json:"howUsed,omitempty"`           // How it's used to do the job
	Limitations       []string `json:"limitations,omitempty"`       // Where it falls short
	Workarounds       []string `json:"workarounds,omitempty"`       // Hacks to make it work
	SatisfactionLevel int      `json:"satisfactionLevel,omitempty"` // 1-10
}

// Force represents a force in the forces diagram (push, pull, anxiety, habit).
type Force struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Strength    string `json:"strength,omitempty"` // strong, medium, weak
	Category    string `json:"category,omitempty"` // For grouping
	Quote       string `json:"quote,omitempty"`    // Customer quote
}

// OpportunityScore represents ODI (Outcome-Driven Innovation) scoring.
// Formula: Opportunity = Importance + max(Importance - Satisfaction, 0)
type OpportunityScore struct {
	OutcomeRef   string  `json:"outcomeRef"`        // Reference to outcome
	Importance   float64 `json:"importance"`        // 1-10 scale
	Satisfaction float64 `json:"satisfaction"`      // 1-10 scale
	Opportunity  float64 `json:"opportunity"`       // Calculated score
	Segment      string  `json:"segment,omitempty"` // Customer segment
}

// JTBDInterview captures a Jobs-to-be-Done interview (switch interview).
type JTBDInterview struct {
	ID              string `json:"id"`
	Date            string `json:"date,omitempty"`
	ParticipantType string `json:"participantType"`
	SwitchContext   string `json:"switchContext,omitempty"` // What they switched from/to

	// Timeline
	FirstThought string `json:"firstThought,omitempty"` // When first thought about switching
	EventOne     string `json:"eventOne,omitempty"`     // First event that triggered consideration
	EventTwo     string `json:"eventTwo,omitempty"`     // Second event
	ActiveSearch string `json:"activeSearch,omitempty"` // When actively searched
	DecisionMade string `json:"decisionMade,omitempty"` // When decided
	PurchaseDate string `json:"purchaseDate,omitempty"` // When purchased/switched

	// Forces discovered
	PushForces []string `json:"pushForces,omitempty"`
	PullForces []string `json:"pullForces,omitempty"`
	Anxieties  []string `json:"anxieties,omitempty"`
	Habits     []string `json:"habits,omitempty"`

	// Key quotes and insights
	KeyQuotes []string `json:"keyQuotes,omitempty"`
	Insights  []string `json:"insights,omitempty"`
}

// Version constant for JTBD canvas.
const VersionJTBD1 Version = "jtbd/1.0"

// NewJTBDCanvas creates a new JTBDCanvas with defaults.
func NewJTBDCanvas(id, title string) *JTBDCanvas {
	return &JTBDCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionJTBD1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *JTBDCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// CalculateOpportunityScore calculates the ODI opportunity score.
// Formula: Opportunity = Importance + max(Importance - Satisfaction, 0)
func CalculateOpportunityScore(importance, satisfaction float64) float64 {
	gap := importance - satisfaction
	if gap < 0 {
		gap = 0
	}
	return importance + gap
}

// TopOpportunities returns outcomes sorted by opportunity score (highest first).
func (c *JTBDCanvas) TopOpportunities(limit int) []OutcomeExpectation {
	// Create a copy and sort by opportunity
	outcomes := make([]OutcomeExpectation, len(c.DesiredOutcomes))
	copy(outcomes, c.DesiredOutcomes)

	// Simple bubble sort (good enough for small lists)
	for i := 0; i < len(outcomes)-1; i++ {
		for j := i + 1; j < len(outcomes); j++ {
			if outcomes[j].Opportunity > outcomes[i].Opportunity {
				outcomes[i], outcomes[j] = outcomes[j], outcomes[i]
			}
		}
	}

	if limit > 0 && limit < len(outcomes) {
		return outcomes[:limit]
	}
	return outcomes
}

// OverservedOutcomes returns outcomes where satisfaction exceeds importance.
func (c *JTBDCanvas) OverservedOutcomes() []OutcomeExpectation {
	var overserved []OutcomeExpectation
	for _, o := range c.DesiredOutcomes {
		if o.Satisfaction > o.Importance {
			overserved = append(overserved, o)
		}
	}
	return overserved
}

// UnderservedOutcomes returns outcomes where importance exceeds satisfaction.
func (c *JTBDCanvas) UnderservedOutcomes() []OutcomeExpectation {
	var underserved []OutcomeExpectation
	for _, o := range c.DesiredOutcomes {
		if o.Importance > o.Satisfaction {
			underserved = append(underserved, o)
		}
	}
	return underserved
}

// GetJobStage returns a job stage by ID.
func (c *JTBDCanvas) GetJobStage(id string) *JobStage {
	for i := range c.JobMap {
		if c.JobMap[i].ID == id {
			return &c.JobMap[i]
		}
	}
	return nil
}

// FunctionalJobs returns only functional job statements.
func (c *JTBDCanvas) FunctionalJobs() []JobStatement {
	var jobs []JobStatement
	for _, j := range c.RelatedJobs {
		if j.Type == JobTypeFunctional {
			jobs = append(jobs, j)
		}
	}
	if c.MainJob != nil && c.MainJob.Type == JobTypeFunctional {
		jobs = append([]JobStatement{*c.MainJob}, jobs...)
	}
	return jobs
}

// EmotionalJobs returns only emotional job statements.
func (c *JTBDCanvas) EmotionalJobs() []JobStatement {
	var jobs []JobStatement
	for _, j := range c.RelatedJobs {
		if j.Type == JobTypeEmotional {
			jobs = append(jobs, j)
		}
	}
	if c.MainJob != nil && c.MainJob.Type == JobTypeEmotional {
		jobs = append([]JobStatement{*c.MainJob}, jobs...)
	}
	return jobs
}

// SocialJobs returns only social job statements.
func (c *JTBDCanvas) SocialJobs() []JobStatement {
	var jobs []JobStatement
	for _, j := range c.RelatedJobs {
		if j.Type == JobTypeSocial {
			jobs = append(jobs, j)
		}
	}
	if c.MainJob != nil && c.MainJob.Type == JobTypeSocial {
		jobs = append([]JobStatement{*c.MainJob}, jobs...)
	}
	return jobs
}
