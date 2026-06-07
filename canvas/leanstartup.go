package canvas

// LeanStartupCanvas implements Eric Ries' Lean Startup methodology.
// It tracks the Build-Measure-Learn feedback loop, MVP experiments,
// and pivot/persevere decisions.
type LeanStartupCanvas struct {
	Metadata Metadata `json:"metadata"`

	// Vision and strategy
	Vision             string `json:"vision"`                       // Long-term vision
	Strategy           string `json:"strategy,omitempty"`           // Current strategy to achieve vision
	ProductMarketFit   string `json:"productMarketFit,omitempty"`   // Current PMF status
	TargetCustomer     string `json:"targetCustomer"`               // Who are we building for?
	ProblemHypothesis  string `json:"problemHypothesis"`            // What problem are we solving?
	SolutionHypothesis string `json:"solutionHypothesis,omitempty"` // How are we solving it?

	// Core hypotheses to validate
	ValueHypothesis  *LSValueHypothesis  `json:"valueHypothesis,omitempty"`  // Does the product deliver value?
	GrowthHypothesis *LSGrowthHypothesis `json:"growthHypothesis,omitempty"` // Can we grow sustainably?

	// Validation activities
	MVPs        []MVP          `json:"mvps,omitempty"`        // MVP iterations
	Experiments []LSExperiment `json:"experiments,omitempty"` // Build-Measure-Learn experiments

	// Direction changes
	Pivots       []Pivot `json:"pivots,omitempty"`       // Pivot history
	CurrentPivot string  `json:"currentPivot,omitempty"` // ID of current pivot being evaluated

	// Innovation accounting
	Metrics       []LSMetric `json:"metrics,omitempty"`       // Key metrics being tracked
	LearningGoals []string   `json:"learningGoals,omitempty"` // Current learning objectives

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// LSValueHypothesis tracks the value hypothesis validation.
// The value hypothesis tests whether a product delivers value to customers.
type LSValueHypothesis struct {
	ID          string `json:"id"`
	Statement   string `json:"statement"`             // The value we believe we provide
	Evidence    string `json:"evidence,omitempty"`    // Evidence collected
	Validated   *bool  `json:"validated,omitempty"`   // True if validated, false if invalidated, nil if untested
	Confidence  string `json:"confidence,omitempty"`  // high, medium, low
	NextStep    string `json:"nextStep,omitempty"`    // What to do next
	CustomerRef string `json:"customerRef,omitempty"` // Reference to customer segment
}

// LSGrowthHypothesis tracks the growth hypothesis validation.
// The growth hypothesis tests how new customers discover the product.
type LSGrowthHypothesis struct {
	ID          string      `json:"id"`
	GrowthModel GrowthModel `json:"growthModel"`          // sticky, viral, paid
	Statement   string      `json:"statement"`            // The growth mechanism we believe will work
	Evidence    string      `json:"evidence,omitempty"`   // Evidence collected
	Validated   *bool       `json:"validated,omitempty"`  // True if validated, false if invalidated, nil if untested
	Confidence  string      `json:"confidence,omitempty"` // high, medium, low
	NextStep    string      `json:"nextStep,omitempty"`   // What to do next
}

// GrowthModel represents the type of growth engine.
type GrowthModel string

// Growth model constants from Lean Startup.
const (
	GrowthModelSticky GrowthModel = "sticky" // Retention-based growth
	GrowthModelViral  GrowthModel = "viral"  // Word-of-mouth/sharing growth
	GrowthModelPaid   GrowthModel = "paid"   // Paid acquisition growth
)

// MVP represents a Minimum Viable Product iteration.
type MVP struct {
	ID             string   `json:"id"`
	Name           string   `json:"name"`
	Type           MVPType  `json:"type"`                     // Type of MVP
	Description    string   `json:"description"`              // What the MVP includes
	Goal           string   `json:"goal"`                     // What we're trying to learn
	Audience       string   `json:"audience,omitempty"`       // Target audience for this MVP
	Timeline       string   `json:"timeline,omitempty"`       // Expected duration
	Status         string   `json:"status"`                   // planned, building, measuring, learning, complete
	Results        string   `json:"results,omitempty"`        // What we measured
	Learnings      string   `json:"learnings,omitempty"`      // What we learned
	Decision       string   `json:"decision,omitempty"`       // pivot, persevere, or iterate
	NextMVP        string   `json:"nextMVP,omitempty"`        // ID of next MVP if iterating
	Iteration      int      `json:"iteration,omitempty"`      // Iteration number
	HypothesisRefs []string `json:"hypothesisRefs,omitempty"` // Hypotheses being tested
}

// MVPType represents the type of MVP.
type MVPType string

// MVP type constants.
const (
	MVPTypeLandingPage    MVPType = "landing-page"    // Simple landing page to gauge interest
	MVPTypeSmokeTest      MVPType = "smoke-test"      // Fake door test
	MVPTypeVideo          MVPType = "video"           // Explainer video MVP
	MVPTypeConcierge      MVPType = "concierge"       // Manual service delivery
	MVPTypeWizardOfOz     MVPType = "wizard-of-oz"    // Automated-looking manual process
	MVPTypePiecemeal      MVPType = "piecemeal"       // Cobbled together from existing services
	MVPTypeSingleFeature  MVPType = "single-feature"  // One core feature
	MVPTypePrototype      MVPType = "prototype"       // Interactive prototype
	MVPTypeCrowdfunding   MVPType = "crowdfunding"    // Pre-sales validation
	MVPTypeEmail          MVPType = "email"           // Email-based product delivery
	MVPTypeSpreadsheet    MVPType = "spreadsheet"     // Spreadsheet-powered product
	MVPTypePaperPrototype MVPType = "paper-prototype" // Paper/sketch prototype
)

// LSExperiment represents a Build-Measure-Learn experiment.
type LSExperiment struct {
	ID string `json:"id"`

	// Build phase
	BuildDescription string `json:"buildDescription"`      // What we're building
	BuildCost        string `json:"buildCost,omitempty"`   // Resources required
	BuildTime        string `json:"buildTime,omitempty"`   // Time to build
	BuildStatus      string `json:"buildStatus,omitempty"` // planned, in-progress, complete
	BuildOutput      string `json:"buildOutput,omitempty"` // What was built

	// Measure phase
	MeasureMethod    string   `json:"measureMethod"`             // How we're measuring
	MeasureMetrics   []string `json:"measureMetrics,omitempty"`  // Specific metrics
	MeasureBaseline  string   `json:"measureBaseline,omitempty"` // Current baseline
	MeasureTarget    string   `json:"measureTarget,omitempty"`   // Target to validate
	MeasureActual    string   `json:"measureActual,omitempty"`   // Actual measurement
	MeasureStatus    string   `json:"measureStatus,omitempty"`   // planned, in-progress, complete
	ActionableMetric bool     `json:"actionableMetric"`          // True if metric is actionable (not vanity)

	// Learn phase
	LearnHypothesis string `json:"learnHypothesis"`          // Hypothesis being tested
	LearnInsight    string `json:"learnInsight,omitempty"`   // What we learned
	LearnValidated  *bool  `json:"learnValidated,omitempty"` // Was hypothesis validated?
	LearnDecision   string `json:"learnDecision,omitempty"`  // pivot, persevere, or continue

	// Overall
	Status    ExperimentStatus `json:"status"`
	MVPRef    string           `json:"mvpRef,omitempty"`    // Reference to MVP if applicable
	StartDate string           `json:"startDate,omitempty"` // When experiment started
	EndDate   string           `json:"endDate,omitempty"`   // When experiment ended
	Owner     string           `json:"owner,omitempty"`     // Who owns this experiment
}

// Pivot represents a strategic direction change.
type Pivot struct {
	ID             string    `json:"id"`
	Type           PivotType `json:"type"`                     // Type of pivot
	Date           string    `json:"date,omitempty"`           // When the pivot occurred
	FromState      string    `json:"fromState"`                // What we pivoted from
	ToState        string    `json:"toState"`                  // What we pivoted to
	Reason         string    `json:"reason"`                   // Why we pivoted
	Evidence       string    `json:"evidence,omitempty"`       // Data supporting the pivot
	ExperimentRefs []string  `json:"experimentRefs,omitempty"` // Experiments that informed the pivot
	Outcome        string    `json:"outcome,omitempty"`        // Result of the pivot
	Status         string    `json:"status,omitempty"`         // proposed, approved, executed, evaluated
}

// PivotType represents the type of pivot from Lean Startup.
type PivotType string

// Pivot type constants from Eric Ries' Lean Startup.
const (
	PivotTypeZoomIn          PivotType = "zoom-in"          // Single feature becomes the product
	PivotTypeZoomOut         PivotType = "zoom-out"         // Product becomes a feature
	PivotTypeCustomerSegment PivotType = "customer-segment" // Different customer than originally planned
	PivotTypeCustomerNeed    PivotType = "customer-need"    // Different problem for same customer
	PivotTypePlatform        PivotType = "platform"         // Application to platform (or vice versa)
	PivotTypeBusinessArch    PivotType = "business-arch"    // High margin/low volume to low margin/high volume (or vice versa)
	PivotTypeValueCapture    PivotType = "value-capture"    // Change in monetization
	PivotTypeEngineOfGrowth  PivotType = "engine-of-growth" // Change growth model (viral, sticky, paid)
	PivotTypeChannel         PivotType = "channel"          // Change distribution channel
	PivotTypeTechnology      PivotType = "technology"       // Change underlying technology
)

// LSMetric represents a metric tracked for innovation accounting.
type LSMetric struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`                 // actionable, vanity, leading, lagging
	Definition string `json:"definition"`           // How it's calculated
	Current    string `json:"current,omitempty"`    // Current value
	Target     string `json:"target,omitempty"`     // Target value
	Baseline   string `json:"baseline,omitempty"`   // Starting value
	Trend      string `json:"trend,omitempty"`      // improving, declining, stable
	Cohort     string `json:"cohort,omitempty"`     // Cohort this metric applies to
	UpdateFreq string `json:"updateFreq,omitempty"` // How often updated (daily, weekly, etc.)
}

// Version constant for Lean Startup canvas.
const VersionLeanStartup1 Version = "leanstartup/1.0"

// NewLeanStartupCanvas creates a new LeanStartupCanvas with defaults.
func NewLeanStartupCanvas(id, title string) *LeanStartupCanvas {
	return &LeanStartupCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionLeanStartup1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *LeanStartupCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// CompletedMVPs returns MVPs that have been completed.
func (c *LeanStartupCanvas) CompletedMVPs() []MVP {
	var completed []MVP
	for _, mvp := range c.MVPs {
		if mvp.Status == "complete" {
			completed = append(completed, mvp)
		}
	}
	return completed
}

// ActiveMVPs returns MVPs that are currently in progress.
func (c *LeanStartupCanvas) ActiveMVPs() []MVP {
	var active []MVP
	for _, mvp := range c.MVPs {
		if mvp.Status == "building" || mvp.Status == "measuring" || mvp.Status == "learning" {
			active = append(active, mvp)
		}
	}
	return active
}

// PivotHistory returns all executed pivots in chronological order.
func (c *LeanStartupCanvas) PivotHistory() []Pivot {
	var history []Pivot
	for _, pivot := range c.Pivots {
		if pivot.Status == "executed" || pivot.Status == "evaluated" {
			history = append(history, pivot)
		}
	}
	return history
}

// ActionableMetrics returns only metrics that are actionable (not vanity).
func (c *LeanStartupCanvas) ActionableMetrics() []LSMetric {
	var actionable []LSMetric
	for _, m := range c.Metrics {
		if m.Type == "actionable" || m.Type == "leading" {
			actionable = append(actionable, m)
		}
	}
	return actionable
}

// IsValueHypothesisValidated returns true if the value hypothesis has been validated.
func (c *LeanStartupCanvas) IsValueHypothesisValidated() bool {
	return c.ValueHypothesis != nil && c.ValueHypothesis.Validated != nil && *c.ValueHypothesis.Validated
}

// IsGrowthHypothesisValidated returns true if the growth hypothesis has been validated.
func (c *LeanStartupCanvas) IsGrowthHypothesisValidated() bool {
	return c.GrowthHypothesis != nil && c.GrowthHypothesis.Validated != nil && *c.GrowthHypothesis.Validated
}

// HasProductMarketFit returns true if both value and growth hypotheses are validated.
func (c *LeanStartupCanvas) HasProductMarketFit() bool {
	return c.IsValueHypothesisValidated() && c.IsGrowthHypothesisValidated()
}
