package canvas

// Continuous Discovery is Teresa Torres's framework for integrating customer
// research into the product development process through weekly touchpoints,
// opportunity solution trees, and assumption testing.
//
// References:
//   - https://www.producttalk.org/
//   - Teresa Torres: "Continuous Discovery Habits"
//
// Core concepts:
//   - Weekly touchpoints: Regular (weekly) customer conversations
//   - Opportunity Solution Trees: Visualizing the path from outcome to solutions
//   - Assumption testing: Small experiments to reduce risk before building
//   - Experience mapping: Understanding customer journeys
//
// Note: OpportunitySolutionTree is defined in ost.go

// DiscoverySnapshot captures a week's worth of discovery activities.
// This is the continuous discovery "heartbeat" - what was learned this week.
type DiscoverySnapshot struct {
	Metadata Metadata `json:"metadata"`

	// Time period
	Week      string `json:"week"`                // e.g., "2024-W03"
	StartDate string `json:"startDate,omitempty"` // ISO date
	EndDate   string `json:"endDate,omitempty"`   // ISO date

	// Interviews conducted
	Interviews []CDInterview `json:"interviews,omitempty"`

	// Opportunities discovered or updated
	OpportunitiesDiscovered []CDOpportunityUpdate `json:"opportunitiesDiscovered,omitempty"`

	// Assumptions tested
	AssumptionTests []CDAssumptionTest `json:"assumptionTests,omitempty"`

	// Key learnings
	KeyLearnings []string `json:"keyLearnings,omitempty"`

	// Decisions made based on discovery
	Decisions []CDDecision `json:"decisions,omitempty"`

	// Link to OST being updated
	OSTRef string `json:"ostRef,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// CDInterview captures a customer interview or touchpoint.
type CDInterview struct {
	ID   string `json:"id"`
	Date string `json:"date"` // ISO date

	// Participant info (anonymized)
	ParticipantType string `json:"participantType"`         // e.g., "enterprise user", "new user"
	ParticipantID   string `json:"participantId,omitempty"` // Anonymous ID for tracking
	Segment         string `json:"segment,omitempty"`       // Market segment
	RecruitSource   string `json:"recruitSource,omitempty"` // How recruited

	// Interview structure
	InterviewType string `json:"interviewType,omitempty"` // story-based, usability, concept
	Duration      string `json:"duration,omitempty"`      // e.g., "30 min"

	// Story-based interview capture
	Stories []CDStory `json:"stories,omitempty"`

	// Opportunities surfaced
	OpportunitiesSurfaced []string `json:"opportunitiesSurfaced,omitempty"` // Opportunity IDs or descriptions

	// Raw notes
	Notes string `json:"notes,omitempty"`

	// Interview quality
	Quality      string `json:"quality,omitempty"` // high, medium, low
	QualityNotes string `json:"qualityNotes,omitempty"`
}

// CDStory captures a specific story from an interview.
// In continuous discovery, we collect stories (past behavior) not opinions.
type CDStory struct {
	ID        string `json:"id"`
	Situation string `json:"situation"` // What was the context?
	Behavior  string `json:"behavior"`  // What did they do?
	Outcome   string `json:"outcome"`   // What happened?

	// Emotions and pain points
	Emotions   []string `json:"emotions,omitempty"`   // How did they feel?
	PainPoints []string `json:"painPoints,omitempty"` // What was frustrating?
	Workaround string   `json:"workaround,omitempty"` // How did they work around it?

	// Opportunities this reveals
	OpportunityRefs []string `json:"opportunityRefs,omitempty"`

	// Verbatim quotes
	Quotes []string `json:"quotes,omitempty"`
}

// CDOpportunityUpdate tracks changes to opportunities from discovery.
type CDOpportunityUpdate struct {
	OpportunityID string `json:"opportunityId"`
	Action        string `json:"action"` // discovered, strengthened, weakened, retired

	// What changed
	Description string `json:"description,omitempty"`

	// Evidence
	EvidenceCount int      `json:"evidenceCount,omitempty"` // How many times seen
	Sources       []string `json:"sources,omitempty"`       // Interview IDs, analytics, etc.

	// Priority change
	PriorityBefore int `json:"priorityBefore,omitempty"`
	PriorityAfter  int `json:"priorityAfter,omitempty"`
}

// CDAssumptionTest represents a small experiment to test an assumption.
type CDAssumptionTest struct {
	ID string `json:"id"`

	// What assumption are we testing?
	Assumption CDAssumption `json:"assumption"`

	// Experiment design
	Method          string `json:"method"`                 // prototype, fake-door, survey, interview, data-analysis
	Hypothesis      string `json:"hypothesis"`             // If X, then Y, because Z
	SuccessCriteria string `json:"successCriteria"`        // How we'll know it worked
	Participants    string `json:"participants,omitempty"` // Who's involved
	SampleSize      int    `json:"sampleSize,omitempty"`   // How many
	Duration        string `json:"duration,omitempty"`     // How long

	// Status and results
	Status    string `json:"status"` // planned, running, completed
	StartDate string `json:"startDate,omitempty"`
	EndDate   string `json:"endDate,omitempty"`

	// Results
	Result   string `json:"result,omitempty"`   // validated, invalidated, inconclusive
	Data     string `json:"data,omitempty"`     // Raw data/observations
	Learning string `json:"learning,omitempty"` // What did we learn?

	// Next steps
	NextStep         string `json:"nextStep,omitempty"`         // What should we do next?
	PivotDescription string `json:"pivotDescription,omitempty"` // If pivoting, what changed
}

// CDAssumption represents something we believe but haven't validated.
// Based on the Assumption Mapping technique in Continuous Discovery.
type CDAssumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`

	// Assumption type from Teresa Torres framework
	Type string `json:"type"` // desirability, viability, feasibility, usability, ethical

	// Risk assessment
	Importance string `json:"importance,omitempty"` // high, medium, low (if wrong, how bad?)
	Confidence string `json:"confidence,omitempty"` // high, medium, low (how sure are we?)

	// For prioritization: test high-importance, low-confidence assumptions first
	Priority int `json:"priority,omitempty"` // Calculated or assigned

	// Evidence
	Evidence        []string `json:"evidence,omitempty"`        // What supports this?
	CounterEvidence []string `json:"counterEvidence,omitempty"` // What contradicts this?

	// Status
	Validated      bool   `json:"validated"`
	ValidationDate string `json:"validationDate,omitempty"`
}

// CDDecision captures a decision made based on discovery learnings.
type CDDecision struct {
	ID          string `json:"id"`
	Date        string `json:"date"` // ISO date
	Description string `json:"description"`
	Rationale   string `json:"rationale,omitempty"` // Why this decision

	// What informed this decision
	InterviewRefs []string `json:"interviewRefs,omitempty"`
	TestRefs      []string `json:"testRefs,omitempty"`

	// Decision type
	Type string `json:"type,omitempty"` // pursue, pivot, persevere, kill

	// Impact
	OSTChanges []string `json:"ostChanges,omitempty"` // Changes to opportunity solution tree
}

// AssumptionMap organizes assumptions by type for systematic testing.
type AssumptionMap struct {
	Metadata Metadata `json:"metadata"`

	// Solution being evaluated
	SolutionRef string `json:"solutionRef,omitempty"` // Link to solution in OST

	// Assumptions by type
	Desirability []CDAssumption `json:"desirability,omitempty"` // Will customers want this?
	Viability    []CDAssumption `json:"viability,omitempty"`    // Will it work for the business?
	Feasibility  []CDAssumption `json:"feasibility,omitempty"`  // Can we build it?
	Usability    []CDAssumption `json:"usability,omitempty"`    // Can users figure it out?
	Ethical      []CDAssumption `json:"ethical,omitempty"`      // Should we build it?

	// Prioritized test order
	TestOrder []string `json:"testOrder,omitempty"` // Assumption IDs in order to test

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// ExperienceMap captures a customer's end-to-end experience.
// Used to identify opportunities across the journey.
type ExperienceMap struct {
	Metadata Metadata `json:"metadata"`

	// What experience are we mapping?
	Experience string `json:"experience"` // e.g., "First-time user onboarding"

	// Persona
	PersonaRef         string `json:"personaRef,omitempty"`
	PersonaDescription string `json:"personaDescription,omitempty"`

	// Journey phases
	Phases []EMPhase `json:"phases"`

	// Overall pain points and opportunities
	TopPainPoints    []string `json:"topPainPoints,omitempty"`
	TopOpportunities []string `json:"topOpportunities,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// EMPhase is a phase in the experience map.
type EMPhase struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Order       int    `json:"order"`

	// What happens in this phase
	Actions []EMAction `json:"actions,omitempty"`

	// Customer state
	Thinking string   `json:"thinking,omitempty"` // What are they thinking?
	Feeling  string   `json:"feeling,omitempty"`  // What are they feeling? (positive/negative)
	Emotions []string `json:"emotions,omitempty"` // Specific emotions

	// Pain points and opportunities
	PainPoints    []string `json:"painPoints,omitempty"`
	Opportunities []string `json:"opportunities,omitempty"`

	// Touchpoints
	Touchpoints []string `json:"touchpoints,omitempty"` // How they interact with us
}

// EMAction is a specific action within a phase.
type EMAction struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Actor       string `json:"actor,omitempty"`   // Who does this
	Channel     string `json:"channel,omitempty"` // Where (web, mobile, support, etc.)
	Pain        string `json:"pain,omitempty"`    // What's hard about this
	Delight     string `json:"delight,omitempty"` // What's good about this
}

// NewDiscoverySnapshot creates a new discovery snapshot.
func NewDiscoverySnapshot(id, title string) *DiscoverySnapshot {
	return &DiscoverySnapshot{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionDiscovery1,
		},
	}
}

// NewAssumptionMap creates a new assumption map.
func NewAssumptionMap(id, title string) *AssumptionMap {
	return &AssumptionMap{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionDiscovery1,
		},
	}
}

// NewExperienceMap creates a new experience map.
func NewExperienceMap(id, title string) *ExperienceMap {
	return &ExperienceMap{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionDiscovery1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (d *DiscoverySnapshot) GetPRDReference() *PRDReference {
	return d.PRDRef
}

// GetPRDReference returns the PRD reference.
func (a *AssumptionMap) GetPRDReference() *PRDReference {
	return a.PRDRef
}

// GetPRDReference returns the PRD reference.
func (e *ExperienceMap) GetPRDReference() *PRDReference {
	return e.PRDRef
}

// InterviewCount returns the number of interviews this week.
func (d *DiscoverySnapshot) InterviewCount() int {
	return len(d.Interviews)
}

// HasWeeklyTouchpoint returns true if at least one interview was conducted.
func (d *DiscoverySnapshot) HasWeeklyTouchpoint() bool {
	return len(d.Interviews) > 0
}

// TotalStories returns the total number of stories collected.
func (d *DiscoverySnapshot) TotalStories() int {
	count := 0
	for _, interview := range d.Interviews {
		count += len(interview.Stories)
	}
	return count
}

// CompletedTests returns assumption tests that are completed.
func (d *DiscoverySnapshot) CompletedTests() []CDAssumptionTest {
	var completed []CDAssumptionTest
	for _, test := range d.AssumptionTests {
		if test.Status == "completed" {
			completed = append(completed, test)
		}
	}
	return completed
}

// AllAssumptions returns all assumptions across all types.
func (a *AssumptionMap) AllAssumptions() []CDAssumption {
	var all []CDAssumption
	all = append(all, a.Desirability...)
	all = append(all, a.Viability...)
	all = append(all, a.Feasibility...)
	all = append(all, a.Usability...)
	all = append(all, a.Ethical...)
	return all
}

// UnvalidatedAssumptions returns assumptions not yet validated.
func (a *AssumptionMap) UnvalidatedAssumptions() []CDAssumption {
	var unvalidated []CDAssumption
	for _, assumption := range a.AllAssumptions() {
		if !assumption.Validated {
			unvalidated = append(unvalidated, assumption)
		}
	}
	return unvalidated
}

// HighRiskAssumptions returns high-importance, low-confidence assumptions.
func (a *AssumptionMap) HighRiskAssumptions() []CDAssumption {
	var highRisk []CDAssumption
	for _, assumption := range a.AllAssumptions() {
		if assumption.Importance == "high" && assumption.Confidence == "low" && !assumption.Validated {
			highRisk = append(highRisk, assumption)
		}
	}
	return highRisk
}

// RiskiestAssumption returns the assumption that should be tested first.
func (a *AssumptionMap) RiskiestAssumption() *CDAssumption {
	highRisk := a.HighRiskAssumptions()
	if len(highRisk) == 0 {
		return nil
	}
	// Return first high-risk assumption (could be enhanced with priority sorting)
	return &highRisk[0]
}
