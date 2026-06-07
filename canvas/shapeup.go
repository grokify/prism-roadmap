package canvas

// Shape Up is Basecamp's product development methodology created by Ryan Singer.
// It emphasizes fixed time, variable scope, and betting on pitches.
//
// References:
//   - https://basecamp.com/shapeup
//   - Ryan Singer: "Shape Up: Stop Running in Circles and Ship Work that Matters"
//
// Core concepts:
//   - Appetite: Fixed time budget (2 or 6 weeks), not estimates
//   - Shaping: Define the problem and solution at the right level of abstraction
//   - Betting: Leadership bets on shaped pitches, not backlogs
//   - Building: Small teams with autonomy during the cycle
//   - Cool-down: 2-week break between cycles for bugs, exploration

// ShapeUpPitch is the core artifact in Shape Up - a shaped problem and solution
// ready for betting. Pitches are written at the right level of abstraction:
// concrete enough to evaluate, abstract enough to leave room for builders.
type ShapeUpPitch struct {
	Metadata Metadata `json:"metadata"`

	// Problem definition
	Problem     SUProblem      `json:"problem"`
	Appetite    SUAppetite     `json:"appetite"`
	Solution    SUSolution     `json:"solution"`
	RabbitHoles []SURabbitHole `json:"rabbitHoles,omitempty"` // Things to avoid
	NoGos       []string       `json:"noGos,omitempty"`       // Out of scope

	// Betting context
	BettingStatus string `json:"bettingStatus,omitempty"` // pitched, bet, declined, deferred
	CycleRef      string `json:"cycleRef,omitempty"`      // Which cycle this was bet for

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// SUProblem describes the raw idea or problem to be solved.
type SUProblem struct {
	// The raw idea or request that started this
	RawIdea string `json:"rawIdea"`

	// The shaped problem statement
	Statement string `json:"statement"`

	// Why this matters now
	WhyNow string `json:"whyNow,omitempty"`

	// Who has this problem
	Audience string `json:"audience,omitempty"`

	// Evidence this is a real problem
	Evidence []string `json:"evidence,omitempty"`

	// Current workarounds
	CurrentWorkarounds []string `json:"currentWorkarounds,omitempty"`
}

// SUAppetite defines the time budget for the work.
// In Shape Up, appetite is a constraint, not an estimate.
type SUAppetite struct {
	// Duration in weeks: typically 2 (small batch) or 6 (big batch)
	Weeks int `json:"weeks"`

	// Size classification: "small-batch" (2 weeks) or "big-batch" (6 weeks)
	Size string `json:"size"` // small-batch, big-batch

	// Why this appetite is appropriate
	Rationale string `json:"rationale,omitempty"`

	// What would we cut if we only had less time?
	TradeoffIfSmaller string `json:"tradeoffIfSmaller,omitempty"`

	// What could we add if we had more time? (usually not relevant)
	ExpansionIfLarger string `json:"expansionIfLarger,omitempty"`
}

// SUSolution describes the shaped solution at the right level of abstraction.
type SUSolution struct {
	// High-level approach - what are we building?
	Approach string `json:"approach"`

	// Breadboarding: abstract UI flows showing affordances and connections
	// Not mockups - just boxes and arrows showing the flow
	Breadboards []SUBreadboard `json:"breadboards,omitempty"`

	// Fat marker sketches: rough visual concepts
	// Intentionally low-fidelity to leave room for builders
	FatMarkerSketches []SUSketch `json:"fatMarkerSketches,omitempty"`

	// Key elements the solution must include
	MustInclude []string `json:"mustInclude,omitempty"`

	// Elements that are nice to have but can be cut
	NiceToHave []string `json:"niceToHave,omitempty"`

	// Technical approach if relevant
	TechnicalApproach string `json:"technicalApproach,omitempty"`
}

// SUBreadboard is an abstract UI flow diagram showing places, affordances,
// and connection lines. It's like a flowchart for UI concepts.
type SUBreadboard struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Places      []string `json:"places"`      // Screens or states
	Affordances []string `json:"affordances"` // Actions users can take
	Flow        string   `json:"flow"`        // Description of the flow
}

// SUSketch is a fat marker sketch - intentionally rough to leave room
// for designers and builders to fill in details.
type SUSketch struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ImageRef    string `json:"imageRef,omitempty"` // Reference to sketch image
	Notes       string `json:"notes,omitempty"`
}

// SURabbitHole is a known complexity or risk that should be avoided.
// Identifying rabbit holes upfront prevents scope creep during building.
type SURabbitHole struct {
	ID           string `json:"id"`
	Description  string `json:"description"`
	WhyDangerous string `json:"whyDangerous,omitempty"` // Why this could blow up the project
	Avoidance    string `json:"avoidance,omitempty"`    // How to avoid or handle it
}

// ShapeUpBet represents a decision to bet on a pitch for a cycle.
// The betting table is where leadership decides what to work on.
type ShapeUpBet struct {
	Metadata Metadata `json:"metadata"`

	// Reference to the pitch being bet on
	PitchRef   string `json:"pitchRef"`
	PitchTitle string `json:"pitchTitle,omitempty"`

	// Cycle information
	Cycle SUCycle `json:"cycle"`

	// Team assignment
	Team SUTeam `json:"team"`

	// Betting decision
	Decision    string `json:"decision"`              // bet, declined, deferred
	Rationale   string `json:"rationale,omitempty"`   // Why this decision
	DecidedBy   string `json:"decidedBy,omitempty"`   // Who made the call
	DecidedDate string `json:"decidedDate,omitempty"` // When decided

	// If deferred, when to reconsider
	DeferredUntil string `json:"deferredUntil,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// SUCycle represents a development cycle (typically 6 weeks + 2 week cool-down).
type SUCycle struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`     // e.g., "Q1 2024 Cycle 1"
	StartDate string `json:"startDate"`          // ISO date
	EndDate   string `json:"endDate"`            // ISO date
	Weeks     int    `json:"weeks"`              // Usually 6
	CoolDown  int    `json:"coolDown,omitempty"` // Usually 2 weeks
}

// SUTeam is the small team assigned to build a bet.
// Shape Up uses small, autonomous teams (typically 1 designer + 1-2 programmers).
type SUTeam struct {
	Designer    string   `json:"designer,omitempty"`
	Programmers []string `json:"programmers,omitempty"`
	Lead        string   `json:"lead,omitempty"` // Who's accountable
}

// ShapeUpScope represents scopes and tasks during the building phase.
// Scopes are discovered during building, not upfront.
type ShapeUpScope struct {
	Metadata Metadata `json:"metadata"`

	// Reference to the bet this scope belongs to
	BetRef string `json:"betRef"`

	// The hill chart tracks progress across scopes
	Scopes []SUScope `json:"scopes"`

	// Overall status
	Status      string `json:"status,omitempty"` // in_progress, completed, circuit_breaker
	LastUpdated string `json:"lastUpdated,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// SUScope is a meaningful slice of work discovered during building.
// Scopes are "hills" that teams climb - uncertain at first, then downhill.
type SUScope struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`

	// Hill chart position: 0-100 where 50 is the top of the hill
	// 0-50: uphill (figuring things out, uncertainty)
	// 50-100: downhill (executing known work)
	HillPosition int `json:"hillPosition"`

	// Tasks within the scope (discovered during building)
	Tasks []SUTask `json:"tasks,omitempty"`

	// Status
	Status string `json:"status,omitempty"` // not_started, uphill, downhill, done

	// Notes and learnings
	Notes string `json:"notes,omitempty"`
}

// SUTask is a task within a scope.
type SUTask struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
	Assignee    string `json:"assignee,omitempty"`
}

// NewShapeUpPitch creates a new pitch with defaults.
func NewShapeUpPitch(id, title string) *ShapeUpPitch {
	return &ShapeUpPitch{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionShapeUp1,
		},
	}
}

// NewShapeUpBet creates a new bet with defaults.
func NewShapeUpBet(id, title string) *ShapeUpBet {
	return &ShapeUpBet{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionShapeUp1,
		},
	}
}

// NewShapeUpScope creates a new scope tracker with defaults.
func NewShapeUpScope(id, title string) *ShapeUpScope {
	return &ShapeUpScope{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionShapeUp1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (p *ShapeUpPitch) GetPRDReference() *PRDReference {
	return p.PRDRef
}

// GetPRDReference returns the PRD reference.
func (b *ShapeUpBet) GetPRDReference() *PRDReference {
	return b.PRDRef
}

// GetPRDReference returns the PRD reference.
func (s *ShapeUpScope) GetPRDReference() *PRDReference {
	return s.PRDRef
}

// IsSmallBatch returns true if this is a small batch (2-week) pitch.
func (p *ShapeUpPitch) IsSmallBatch() bool {
	return p.Appetite.Size == "small-batch" || p.Appetite.Weeks <= 2
}

// IsBigBatch returns true if this is a big batch (6-week) pitch.
func (p *ShapeUpPitch) IsBigBatch() bool {
	return p.Appetite.Size == "big-batch" || p.Appetite.Weeks >= 6
}

// IsBet returns true if this pitch was bet on.
func (p *ShapeUpPitch) IsBet() bool {
	return p.BettingStatus == "bet"
}

// IsDeclined returns true if this pitch was declined.
func (p *ShapeUpPitch) IsDeclined() bool {
	return p.BettingStatus == "declined"
}

// HasRabbitHoles returns true if rabbit holes have been identified.
func (p *ShapeUpPitch) HasRabbitHoles() bool {
	return len(p.RabbitHoles) > 0
}

// OverallProgress returns the average hill position across all scopes.
func (s *ShapeUpScope) OverallProgress() int {
	if len(s.Scopes) == 0 {
		return 0
	}
	total := 0
	for _, scope := range s.Scopes {
		total += scope.HillPosition
	}
	return total / len(s.Scopes)
}

// UphillScopes returns scopes still in the uncertain/uphill phase.
func (s *ShapeUpScope) UphillScopes() []SUScope {
	var uphill []SUScope
	for _, scope := range s.Scopes {
		if scope.HillPosition < 50 {
			uphill = append(uphill, scope)
		}
	}
	return uphill
}

// DownhillScopes returns scopes in the execution/downhill phase.
func (s *ShapeUpScope) DownhillScopes() []SUScope {
	var downhill []SUScope
	for _, scope := range s.Scopes {
		if scope.HillPosition >= 50 {
			downhill = append(downhill, scope)
		}
	}
	return downhill
}

// DoneScopes returns completed scopes.
func (s *ShapeUpScope) DoneScopes() []SUScope {
	var done []SUScope
	for _, scope := range s.Scopes {
		if scope.HillPosition == 100 || scope.Status == "done" {
			done = append(done, scope)
		}
	}
	return done
}
