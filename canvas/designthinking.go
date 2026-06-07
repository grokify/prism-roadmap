package canvas

// DesignThinkingCanvas implements the Stanford d.school Design Thinking methodology.
// It covers the five phases: Empathize, Define, Ideate, Prototype, Test.
type DesignThinkingCanvas struct {
	Metadata Metadata `json:"metadata"`

	// Phase 1: Empathize
	EmpathyMaps  []EmpathyMap    `json:"empathyMaps,omitempty"`
	Interviews   []DTInterview   `json:"interviews,omitempty"`
	Observations []DTObservation `json:"observations,omitempty"`

	// Phase 2: Define
	ProblemStatement string       `json:"problemStatement,omitempty"` // Point of View (POV) statement
	UserNeeds        []DTUserNeed `json:"userNeeds,omitempty"`
	Insights         []DTInsight  `json:"insights,omitempty"`
	HowMightWe       []string     `json:"howMightWe,omitempty"` // HMW questions

	// Phase 3: Ideate
	IdeationSessions []IdeationSession `json:"ideationSessions,omitempty"`
	Ideas            []DTIdea          `json:"ideas,omitempty"`
	SelectedIdeas    []string          `json:"selectedIdeas,omitempty"` // IDs of ideas selected for prototyping

	// Phase 4: Prototype
	Prototypes []DTPrototype `json:"prototypes,omitempty"`

	// Phase 5: Test
	Tests          []DTTest `json:"tests,omitempty"`
	IterationCount int      `json:"iterationCount,omitempty"` // Number of iteration cycles

	// Current state
	CurrentPhase DTPhase `json:"currentPhase,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// DTPhase represents a Design Thinking phase.
type DTPhase string

// Design Thinking phase constants.
const (
	DTPhaseEmpathize DTPhase = "empathize"
	DTPhaseDefine    DTPhase = "define"
	DTPhaseIdeate    DTPhase = "ideate"
	DTPhasePrototype DTPhase = "prototype"
	DTPhaseTest      DTPhase = "test"
)

// EmpathyMap captures user empathy research using the four quadrant model.
type EmpathyMap struct {
	ID          string `json:"id"`
	PersonaName string `json:"personaName"`    // Who are we empathizing with?
	Goal        string `json:"goal,omitempty"` // What do they need to accomplish?

	// Four quadrants
	Says   []string `json:"says,omitempty"`   // What they say (quotes)
	Thinks []string `json:"thinks,omitempty"` // What they might be thinking
	Does   []string `json:"does,omitempty"`   // Actions and behaviors
	Feels  []string `json:"feels,omitempty"`  // Emotional state

	// Extended quadrants (from newer models)
	Sees  []string `json:"sees,omitempty"`  // Environment, what they see around them
	Hears []string `json:"hears,omitempty"` // What others say to them

	// Pains and Gains
	Pains []string `json:"pains,omitempty"` // Frustrations, obstacles
	Gains []string `json:"gains,omitempty"` // Wants, needs, goals

	// Source
	Source     string `json:"source,omitempty"`     // Interview, observation, survey
	Date       string `json:"date,omitempty"`       // When empathy research was conducted
	PersonaRef string `json:"personaRef,omitempty"` // Reference to PRD persona
}

// DTInterview captures a user interview.
type DTInterview struct {
	ID              string   `json:"id"`
	ParticipantType string   `json:"participantType"` // Type of user
	Date            string   `json:"date,omitempty"`
	Duration        string   `json:"duration,omitempty"`
	KeyQuotes       []string `json:"keyQuotes,omitempty"`
	Observations    []string `json:"observations,omitempty"`
	Insights        []string `json:"insights,omitempty"`
	EmpathyMapRef   string   `json:"empathyMapRef,omitempty"`
}

// DTObservation captures a field observation.
type DTObservation struct {
	ID          string   `json:"id"`
	Context     string   `json:"context"`     // Where and when observed
	Description string   `json:"description"` // What was observed
	Behaviors   []string `json:"behaviors,omitempty"`
	Surprises   []string `json:"surprises,omitempty"` // Unexpected findings
	Questions   []string `json:"questions,omitempty"` // Questions raised
}

// DTUserNeed represents a user need identified during the Define phase.
type DTUserNeed struct {
	ID           string   `json:"id"`
	User         string   `json:"user"`               // Who has this need?
	Need         string   `json:"need"`               // What do they need?
	Insight      string   `json:"insight"`            // Because (insight)
	Priority     string   `json:"priority,omitempty"` // high, medium, low
	Validated    bool     `json:"validated,omitempty"`
	EvidenceRefs []string `json:"evidenceRefs,omitempty"` // References to supporting empathy data
}

// DTInsight represents an insight derived from empathy research.
type DTInsight struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Type        string   `json:"type,omitempty"`     // behavioral, emotional, contextual
	Strength    string   `json:"strength,omitempty"` // strong, medium, weak
	Sources     []string `json:"sources,omitempty"`  // Interview/observation IDs
	QuoteRefs   []string `json:"quoteRefs,omitempty"`
}

// IdeationSession represents a brainstorming or ideation session.
type IdeationSession struct {
	ID           string   `json:"id"`
	Date         string   `json:"date,omitempty"`
	Method       string   `json:"method"`                // brainstorm, brainwrite, SCAMPER, mindmap, etc.
	HMWQuestion  string   `json:"hmwQuestion,omitempty"` // How Might We question being addressed
	Participants []string `json:"participants,omitempty"`
	Duration     string   `json:"duration,omitempty"`
	IdeaCount    int      `json:"ideaCount,omitempty"`
	IdeaRefs     []string `json:"ideaRefs,omitempty"` // IDs of ideas generated
}

// DTIdea represents an idea generated during ideation.
type DTIdea struct {
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	Category     string   `json:"category,omitempty"`     // feature, service, experience, process
	SessionRef   string   `json:"sessionRef,omitempty"`   // Ideation session where generated
	Votes        int      `json:"votes,omitempty"`        // Dot voting count
	Feasibility  string   `json:"feasibility,omitempty"`  // high, medium, low
	Impact       string   `json:"impact,omitempty"`       // high, medium, low
	Selected     bool     `json:"selected,omitempty"`     // Selected for prototyping
	CombinedWith []string `json:"combinedWith,omitempty"` // IDs of ideas this was combined with
}

// DTPrototype represents a prototype created during the Prototype phase.
type DTPrototype struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`               // paper, digital, physical, storyboard, role-play, wizard-of-oz
	Fidelity    string   `json:"fidelity,omitempty"` // low, medium, high
	Description string   `json:"description,omitempty"`
	IdeaRefs    []string `json:"ideaRefs,omitempty"` // Ideas being prototyped
	Iteration   int      `json:"iteration,omitempty"`
	Status      string   `json:"status,omitempty"`   // planned, in-progress, complete
	ImageRef    string   `json:"imageRef,omitempty"` // Reference to prototype image/file
	TestRefs    []string `json:"testRefs,omitempty"` // Tests conducted with this prototype
}

// DTTest represents a test session during the Test phase.
type DTTest struct {
	ID               string   `json:"id"`
	PrototypeRef     string   `json:"prototypeRef"` // Prototype being tested
	Date             string   `json:"date,omitempty"`
	Participants     int      `json:"participants,omitempty"`
	Method           string   `json:"method,omitempty"`    // usability, A/B, interview, observation
	Questions        []string `json:"questions,omitempty"` // What we wanted to learn
	Findings         []string `json:"findings,omitempty"`  // What we learned
	PositiveFeedback []string `json:"positiveFeedback,omitempty"`
	NegativeFeedback []string `json:"negativeFeedback,omitempty"`
	Suggestions      []string `json:"suggestions,omitempty"`
	NextIteration    string   `json:"nextIteration,omitempty"`  // What to change for next iteration
	ShouldContinue   bool     `json:"shouldContinue,omitempty"` // Continue with this direction?
}

// Version constant for Design Thinking canvas.
const VersionDesignThinking1 Version = "designthinking/1.0"

// NewDesignThinkingCanvas creates a new DesignThinkingCanvas with defaults.
func NewDesignThinkingCanvas(id, title string) *DesignThinkingCanvas {
	return &DesignThinkingCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionDesignThinking1,
		},
		CurrentPhase: DTPhaseEmpathize,
	}
}

// GetPRDReference returns the PRD reference.
func (c *DesignThinkingCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// TotalIdeas returns the count of all ideas.
func (c *DesignThinkingCanvas) TotalIdeas() int {
	return len(c.Ideas)
}

// SelectedIdeaCount returns the count of selected ideas.
func (c *DesignThinkingCanvas) SelectedIdeaCount() int {
	count := 0
	for _, idea := range c.Ideas {
		if idea.Selected {
			count++
		}
	}
	return count
}

// PrototypeCount returns the count of prototypes.
func (c *DesignThinkingCanvas) PrototypeCount() int {
	return len(c.Prototypes)
}

// CompletedTests returns tests that have been completed.
func (c *DesignThinkingCanvas) CompletedTests() []DTTest {
	var completed []DTTest
	for _, test := range c.Tests {
		if len(test.Findings) > 0 {
			completed = append(completed, test)
		}
	}
	return completed
}

// GetEmpathyMap returns an empathy map by ID.
func (c *DesignThinkingCanvas) GetEmpathyMap(id string) *EmpathyMap {
	for i := range c.EmpathyMaps {
		if c.EmpathyMaps[i].ID == id {
			return &c.EmpathyMaps[i]
		}
	}
	return nil
}

// GetIdea returns an idea by ID.
func (c *DesignThinkingCanvas) GetIdea(id string) *DTIdea {
	for i := range c.Ideas {
		if c.Ideas[i].ID == id {
			return &c.Ideas[i]
		}
	}
	return nil
}

// HighImpactIdeas returns ideas marked as high impact.
func (c *DesignThinkingCanvas) HighImpactIdeas() []DTIdea {
	var high []DTIdea
	for _, idea := range c.Ideas {
		if idea.Impact == "high" {
			high = append(high, idea)
		}
	}
	return high
}
