package canvas

import "github.com/grokify/prism-roadmap/prioritization"

// OpportunitySpec is a merged framework combining Jeff Patton's Opportunity Canvas
// (discovery-focused) with Marty Cagan's SVPG Opportunity Assessment (business case-focused).
//
// This 12-box structure provides a comprehensive evaluation for feature-level opportunities,
// bridging user discovery with business validation.
//
// References:
//   - Patton: https://www.jpattonassociates.com/opportunity-canvas/
//   - Cagan: https://www.svpg.com/assessing-product-opportunities/
//
// Grid layout (3 columns x 4 rows):
//
//	| 1. Users & Problem      | 2. Current Solutions | 3. Solution Ideas    |
//	| 4. User Value           | 5. Business Value    | 6. Competitive Edge  |
//	| 7. Market & Timing      | 8. Go-to-Market      | 9. Success Metrics   |
//	| 10. Critical Reqs       | 11. Risks & Assump.  | 12. Recommendation   |
type OpportunitySpec struct {
	Metadata Metadata `json:"metadata"`

	// === Row 1: Discovery ===

	// Box 1: Users & Problem (Patton: Users/Problems merged)
	// "Who has the problem and what problem are we solving?"
	UsersAndProblem OSUsersAndProblem `json:"usersAndProblem"`

	// Box 2: Current Solutions (Patton: Solutions Today + Cagan: Competitive Landscape)
	// "How do people solve this today?"
	CurrentSolutions OSCurrentSolutions `json:"currentSolutions"`

	// Box 3: Solution Ideas (Patton: Solution Ideas)
	// "What are our solution concepts?"
	SolutionIdeas OSSolutionIdeas `json:"solutionIdeas"`

	// === Row 2: Value ===

	// Box 4: User Value (Patton: User Value + Cagan: Value Proposition)
	// "What value does this provide to users?"
	UserValue OSUserValue `json:"userValue"`

	// Box 5: Business Value (Patton: Business Problem/Metrics + Cagan: Market Size)
	// "What value does this provide to the business?"
	BusinessValue OSBusinessValue `json:"businessValue"`

	// Box 6: Competitive Edge (Cagan: Our Differentiator)
	// "Why are we best suited to pursue this?"
	CompetitiveEdge OSCompetitiveEdge `json:"competitiveEdge"`

	// === Row 3: Market ===

	// Box 7: Market & Timing (Cagan: Target Market + Market Window)
	// "Who's the market and why now?"
	MarketAndTiming OSMarketAndTiming `json:"marketAndTiming"`

	// Box 8: Go-to-Market (Patton: Adoption Strategy + Cagan: GTM Strategy)
	// "How will we get this to users?"
	GoToMarket OSGoToMarket `json:"goToMarket"`

	// Box 9: Success Metrics (Patton: User/Business Metrics + Cagan: Metrics/Revenue)
	// "How will we measure success?"
	SuccessMetrics OSSuccessMetrics `json:"successMetrics"`

	// === Row 4: Validation ===

	// Box 10: Critical Requirements (Cagan: Solution Requirements)
	// "What must be true for this to succeed?"
	CriticalRequirements OSCriticalRequirements `json:"criticalRequirements"`

	// Box 11: Risks & Assumptions (Patton: Assumptions + Cagan: Key Assumptions)
	// "What are we betting on and what could go wrong?"
	RisksAndAssumptions OSRisksAndAssumptions `json:"risksAndAssumptions"`

	// Box 12: Recommendation (Cagan: Recommendation + Patton: Budget)
	// "Given all the above, what's the recommendation?"
	Recommendation OSRecommendation `json:"recommendation"`

	// === Prioritization Frameworks ===

	// RICE Score (Intercom)
	// Provides quantitative prioritization based on Reach, Impact, Confidence, Effort.
	// Reference: https://www.intercom.com/blog/rice-simple-prioritization-for-product-managers/
	RICE *prioritization.RICEScore `json:"rice,omitempty"`

	// Kano Analysis (Noriaki Kano)
	// Classifies the feature as Must-Be, Performance, Attractive, Indifferent, or Reverse.
	// Reference: https://www.productplan.com/glossary/kano-model/
	Kano *prioritization.KanoFeature `json:"kano,omitempty"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// OSUsersAndProblem describes who has the problem and what it is (Box 1).
type OSUsersAndProblem struct {
	// Users
	PrimaryUsers      []OSUser `json:"primaryUsers"`                // Primary user types
	SecondaryUsers    []OSUser `json:"secondaryUsers,omitempty"`    // Secondary user types
	AffectedPersonas  []string `json:"affectedPersonas,omitempty"`  // Links to PRD personas
	EstimatedUserBase string   `json:"estimatedUserBase,omitempty"` // How many users affected

	// Problem
	ProblemStatement string      `json:"problemStatement"`        // Clear statement of the problem
	CustomerPains    []string    `json:"customerPains,omitempty"` // Specific pain points
	Problems         []OSProblem `json:"problems,omitempty"`      // Detailed problem breakdown
	CurrentState     string      `json:"currentState,omitempty"`  // How things work today
	DesiredState     string      `json:"desiredState,omitempty"`  // How things should work
	Evidence         string      `json:"evidence,omitempty"`      // Evidence supporting this problem exists
}

// OSUser describes a user type affected by the problem.
type OSUser struct {
	Name         string   `json:"name"`                   // User type/persona name
	Role         string   `json:"role,omitempty"`         // Job role
	Goals        []string `json:"goals,omitempty"`        // What they're trying to achieve
	Frustrations []string `json:"frustrations,omitempty"` // Current pain points
	Context      string   `json:"context,omitempty"`      // In what context do they experience this
	PersonaRef   string   `json:"personaRef,omitempty"`   // Link to PRD persona
}

// OSProblem describes a specific problem instance.
type OSProblem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Frequency   string `json:"frequency,omitempty"` // How often does it occur?
	Severity    string `json:"severity,omitempty"`  // high, medium, low
	Evidence    string `json:"evidence,omitempty"`  // What evidence supports this?
}

// OSCurrentSolutions describes how people solve this today (Box 2).
type OSCurrentSolutions struct {
	DirectCompetitors   []OSCompetitor `json:"directCompetitors,omitempty"`   // Head-to-head competitors
	IndirectCompetitors []OSCompetitor `json:"indirectCompetitors,omitempty"` // Alternative solutions
	Workarounds         []string       `json:"workarounds,omitempty"`         // Manual workarounds
	InternalSolutions   []string       `json:"internalSolutions,omitempty"`   // Internal tools being used
	DoNothing           string         `json:"doNothing,omitempty"`           // Cost of status quo
	MarketDynamics      string         `json:"marketDynamics,omitempty"`      // Competitive dynamics
}

// OSCompetitor describes a competitive alternative.
type OSCompetitor struct {
	Name        string   `json:"name"`
	Type        string   `json:"type,omitempty"` // direct, indirect, substitute
	Strengths   []string `json:"strengths,omitempty"`
	Weaknesses  []string `json:"weaknesses,omitempty"`
	Positioning string   `json:"positioning,omitempty"` // How they position themselves
	MarketShare string   `json:"marketShare,omitempty"` // Estimated market share
}

// OSSolutionIdeas describes potential solution concepts (Box 3).
type OSSolutionIdeas struct {
	Ideas            []OSSolutionIdea `json:"ideas"`                      // Solution concepts
	RecommendedIdea  string           `json:"recommendedIdea,omitempty"`  // Which idea is recommended
	SelectionReason  string           `json:"selectionReason,omitempty"`  // Why this idea
	AlternativesPros []string         `json:"alternativesPros,omitempty"` // Why other ideas might work
	AlternativesCons []string         `json:"alternativesCons,omitempty"` // Why other ideas might not work
}

// OSSolutionIdea describes a potential solution concept.
type OSSolutionIdea struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Approach    string   `json:"approach,omitempty"`    // How would this work?
	Pros        []string `json:"pros,omitempty"`        // Advantages
	Cons        []string `json:"cons,omitempty"`        // Disadvantages
	Feasibility string   `json:"feasibility,omitempty"` // high, medium, low
}

// OSUserValue describes the value to users (Box 4).
type OSUserValue struct {
	ValueStatement  string   `json:"valueStatement"`            // Core value proposition for users
	KeyBenefits     []string `json:"keyBenefits,omitempty"`     // Main benefits
	PainsRelieved   []string `json:"painsRelieved,omitempty"`   // Which pains are addressed
	GainsEnabled    []string `json:"gainsEnabled,omitempty"`    // What new gains are enabled
	JobsToBeDone    []string `json:"jobsToBeDone,omitempty"`    // Jobs this helps accomplish
	UserOutcomes    []string `json:"userOutcomes,omitempty"`    // Expected user outcomes
	AdoptionBarrier string   `json:"adoptionBarrier,omitempty"` // What might prevent adoption
}

// OSBusinessValue describes the value to the business (Box 5).
type OSBusinessValue struct {
	BusinessProblem   string   `json:"businessProblem"`             // Why this matters to the business
	BusinessOutcomes  []string `json:"businessOutcomes,omitempty"`  // Expected business outcomes
	RevenueImpact     string   `json:"revenueImpact,omitempty"`     // Revenue potential
	CostImpact        string   `json:"costImpact,omitempty"`        // Cost savings/implications
	StrategicFit      string   `json:"strategicFit,omitempty"`      // How this fits strategy
	MarketSize        string   `json:"marketSize,omitempty"`        // TAM/SAM/SOM
	GrowthOpportunity string   `json:"growthOpportunity,omitempty"` // Growth potential
}

// OSCompetitiveEdge describes why we're best suited (Box 6).
type OSCompetitiveEdge struct {
	Differentiator     string   `json:"differentiator"`               // What makes us different
	CoreStrengths      []string `json:"coreStrengths,omitempty"`      // What we do well
	UniqueCapabilities []string `json:"uniqueCapabilities,omitempty"` // What only we can do
	StrategicAssets    []string `json:"strategicAssets,omitempty"`    // Assets we can leverage
	TeamExpertise      string   `json:"teamExpertise,omitempty"`      // Relevant team experience
	TechnologyEdge     string   `json:"technologyEdge,omitempty"`     // Technical advantages
	UnfairAdvantage    string   `json:"unfairAdvantage,omitempty"`    // What's hard to replicate
	EntryBarriers      []string `json:"entryBarriers,omitempty"`      // Barriers we create
}

// OSMarketAndTiming describes target market and timing (Box 7).
type OSMarketAndTiming struct {
	// Target Market
	PrimarySegment    string   `json:"primarySegment"`              // Main target segment
	SecondarySegments []string `json:"secondarySegments,omitempty"` // Additional segments
	Industries        []string `json:"industries,omitempty"`        // Target industries
	Geography         []string `json:"geography,omitempty"`         // Geographic focus
	CompanySize       string   `json:"companySize,omitempty"`       // SMB, Mid-market, Enterprise

	// Timing (Why Now?)
	WhyNow            string   `json:"whyNow"`                      // Primary reason for timing
	MarketTriggers    []string `json:"marketTriggers,omitempty"`    // Events creating opportunity
	TechnologyShifts  []string `json:"technologyShifts,omitempty"`  // Tech changes enabling this
	RegulatoryChanges []string `json:"regulatoryChanges,omitempty"` // Regulatory drivers
	CompetitorMoves   []string `json:"competitorMoves,omitempty"`   // Competitor timing
	WindowDuration    string   `json:"windowDuration,omitempty"`    // How long window is open
	UrgencyLevel      string   `json:"urgencyLevel,omitempty"`      // high, medium, low
}

// OSGoToMarket describes how to reach users (Box 8).
type OSGoToMarket struct {
	Strategy            string   `json:"strategy"`                      // Overall GTM approach
	AdoptionPath        string   `json:"adoptionPath,omitempty"`        // How users will find/adopt
	Channels            []string `json:"channels,omitempty"`            // Distribution channels
	SalesModel          string   `json:"salesModel,omitempty"`          // Self-serve, sales-led, PLG
	Partnerships        []string `json:"partnerships,omitempty"`        // Key partnerships needed
	LaunchApproach      string   `json:"launchApproach,omitempty"`      // Beta, GA, etc.
	CustomerAcquisition string   `json:"customerAcquisition,omitempty"` // How to acquire customers
	EstimatedCAC        string   `json:"estimatedCAC,omitempty"`        // Customer acquisition cost
}

// OSSuccessMetrics describes how to measure success (Box 9).
type OSSuccessMetrics struct {
	// User Metrics
	UserMetrics      []string `json:"userMetrics,omitempty"`      // User behavior to track
	AdoptionMetrics  []string `json:"adoptionMetrics,omitempty"`  // Adoption indicators
	SatisfactionGoal string   `json:"satisfactionGoal,omitempty"` // NPS/CSAT target

	// Business Metrics
	BusinessMetrics   []string `json:"businessMetrics,omitempty"`   // Business outcomes to measure
	PrimaryMetric     string   `json:"primaryMetric,omitempty"`     // North star metric
	LeadingIndicators []string `json:"leadingIndicators,omitempty"` // Early signals of success

	// Revenue Metrics
	RevenueModel    string `json:"revenueModel,omitempty"`    // Subscription, usage, etc.
	PricingStrategy string `json:"pricingStrategy,omitempty"` // How to price
	TimeToValue     string `json:"timeToValue,omitempty"`     // When value is realized
}

// OSCriticalRequirements describes what must be true for success (Box 10).
type OSCriticalRequirements struct {
	MustHaveCapabilities  []string `json:"mustHaveCapabilities,omitempty"`  // Non-negotiable features
	TechnicalRequirements []string `json:"technicalRequirements,omitempty"` // Technical needs
	IntegrationNeeds      []string `json:"integrationNeeds,omitempty"`      // Required integrations
	ComplianceNeeds       []string `json:"complianceNeeds,omitempty"`       // Regulatory/compliance
	PerformanceNeeds      string   `json:"performanceNeeds,omitempty"`      // Performance requirements
	ScalabilityNeeds      string   `json:"scalabilityNeeds,omitempty"`      // Scale requirements
	KeyDependencies       []string `json:"keyDependencies,omitempty"`       // External dependencies
	SuccessConditions     []string `json:"successConditions,omitempty"`     // Conditions that must be met
}

// OSRisksAndAssumptions describes risks and assumptions (Box 11).
type OSRisksAndAssumptions struct {
	// Assumptions
	KeyAssumptions     []OSAssumption `json:"keyAssumptions,omitempty"`     // Key assumptions
	RiskiestAssumption string         `json:"riskiestAssumption,omitempty"` // Most critical assumption

	// Risks
	Risks              []OSRisk `json:"risks,omitempty"`              // Identified risks
	HighestRisk        string   `json:"highestRisk,omitempty"`        // Most critical risk
	MitigationStrategy string   `json:"mitigationStrategy,omitempty"` // How to mitigate top risks

	// Validation Plan
	ValidationApproach string   `json:"validationApproach,omitempty"` // How to validate assumptions
	ExperimentsNeeded  []string `json:"experimentsNeeded,omitempty"`  // Experiments to run
}

// OSAssumption represents something assumed true that needs validation.
type OSAssumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"` // user, market, technical, business
	Validated   bool   `json:"validated"`
	Evidence    string `json:"evidence,omitempty"`
	RiskLevel   string `json:"riskLevel,omitempty"` // high, medium, low
}

// OSRisk represents a risk to the opportunity.
type OSRisk struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Category    string `json:"category,omitempty"`    // market, technical, business, regulatory
	Probability string `json:"probability,omitempty"` // high, medium, low
	Impact      string `json:"impact,omitempty"`      // high, medium, low
	Mitigation  string `json:"mitigation,omitempty"`
}

// OSRecommendation provides the go/no-go decision (Box 12).
type OSRecommendation struct {
	// Decision
	Decision   string   `json:"decision"`             // go, no-go, conditional, defer
	Rationale  string   `json:"rationale"`            // Why this decision
	Confidence string   `json:"confidence,omitempty"` // high, medium, low
	Conditions []string `json:"conditions,omitempty"` // Conditions for conditional go
	ReviewDate string   `json:"reviewDate,omitempty"` // When to revisit if deferred

	// Investment Ask (from Patton's Budget)
	TimeEstimate      string `json:"timeEstimate,omitempty"`      // Time to deliver
	TeamSize          string `json:"teamSize,omitempty"`          // Team needed
	ResourcesRequired string `json:"resourcesRequired,omitempty"` // Additional resources
	InvestmentAsk     string `json:"investmentAsk,omitempty"`     // Budget requested
	Constraints       string `json:"constraints,omitempty"`       // Known constraints

	// Next Steps
	NextSteps       []string `json:"nextSteps,omitempty"`       // What to do next
	SuccessCriteria []string `json:"successCriteria,omitempty"` // How we'll know it worked
}

// NewOpportunitySpec creates a new OpportunitySpec with defaults.
func NewOpportunitySpec(id, title string) *OpportunitySpec {
	return &OpportunitySpec{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionOpportunitySpec1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (s *OpportunitySpec) GetPRDReference() *PRDReference {
	return s.PRDRef
}

// IsGo returns true if the recommendation is to proceed.
func (s *OpportunitySpec) IsGo() bool {
	return s.Recommendation.Decision == "go"
}

// IsNoGo returns true if the recommendation is not to proceed.
func (s *OpportunitySpec) IsNoGo() bool {
	return s.Recommendation.Decision == "no-go"
}

// IsConditional returns true if the recommendation is conditional.
func (s *OpportunitySpec) IsConditional() bool {
	return s.Recommendation.Decision == "conditional"
}

// IsDeferred returns true if the decision is deferred.
func (s *OpportunitySpec) IsDeferred() bool {
	return s.Recommendation.Decision == "defer"
}

// HasHighConfidence returns true if confidence is high.
func (s *OpportunitySpec) HasHighConfidence() bool {
	return s.Recommendation.Confidence == "high"
}

// CompetitorCount returns the total number of competitors identified.
func (s *OpportunitySpec) CompetitorCount() int {
	return len(s.CurrentSolutions.DirectCompetitors) +
		len(s.CurrentSolutions.IndirectCompetitors)
}

// UnvalidatedAssumptions returns assumptions that haven't been validated.
func (s *OpportunitySpec) UnvalidatedAssumptions() []OSAssumption {
	var unvalidated []OSAssumption
	for _, a := range s.RisksAndAssumptions.KeyAssumptions {
		if !a.Validated {
			unvalidated = append(unvalidated, a)
		}
	}
	return unvalidated
}

// HighRisks returns risks with high probability or high impact.
func (s *OpportunitySpec) HighRisks() []OSRisk {
	var high []OSRisk
	for _, r := range s.RisksAndAssumptions.Risks {
		if r.Probability == "high" || r.Impact == "high" {
			high = append(high, r)
		}
	}
	return high
}

// === RICE Score Methods ===

// HasRICE returns true if RICE scoring has been completed.
func (s *OpportunitySpec) HasRICE() bool {
	return s.RICE != nil && s.RICE.IsComplete()
}

// GetRICEScore returns the calculated RICE score, or 0 if not set.
func (s *OpportunitySpec) GetRICEScore() float64 {
	if s.RICE == nil {
		return 0
	}
	return s.RICE.Calculate()
}

// SetRICE sets the RICE score with the given values and calculates the score.
func (s *OpportunitySpec) SetRICE(reach int, impact prioritization.ImpactLevel, confidence prioritization.ConfidenceLevel, effort float64) {
	s.RICE = prioritization.NewRICEScore(s.Metadata.ID, reach, impact, confidence, effort)
}

// SetRICEFromRecommendation populates RICE fields from the Recommendation box.
// This extracts Reach from UsersAndProblem.EstimatedUserBase if available.
func (s *OpportunitySpec) SetRICEFromRecommendation() {
	if s.RICE == nil {
		s.RICE = &prioritization.RICEScore{
			FeatureID:   s.Metadata.ID,
			FeatureName: s.Metadata.Title,
		}
	}
	s.RICE.FeatureID = s.Metadata.ID
	s.RICE.FeatureName = s.Metadata.Title
}

// === Kano Model Methods ===

// HasKano returns true if Kano analysis has been completed.
func (s *OpportunitySpec) HasKano() bool {
	return s.Kano != nil && s.Kano.Category != ""
}

// GetKanoCategory returns the Kano category, or empty string if not set.
func (s *OpportunitySpec) GetKanoCategory() prioritization.KanoCategory {
	if s.Kano == nil {
		return ""
	}
	return s.Kano.Category
}

// SetKano sets the Kano analysis with the given responses and classifies.
func (s *OpportunitySpec) SetKano(functional, dysfunctional prioritization.KanoResponse) {
	s.Kano = &prioritization.KanoFeature{
		FeatureID:             s.Metadata.ID,
		FeatureName:           s.Metadata.Title,
		FunctionalResponse:    functional,
		DysfunctionalResponse: dysfunctional,
	}
	s.Kano.Classify()
}

// IsMustHave returns true if this is a Must-Be feature (basic expectation).
func (s *OpportunitySpec) IsMustHave() bool {
	return s.Kano != nil && s.Kano.Category == prioritization.KanoMustBe
}

// IsDelighter returns true if this is an Attractive feature (delighter).
func (s *OpportunitySpec) IsDelighter() bool {
	return s.Kano != nil && s.Kano.Category == prioritization.KanoAttractive
}

// IsPerformance returns true if this is a Performance feature (more is better).
func (s *OpportunitySpec) IsPerformance() bool {
	return s.Kano != nil && s.Kano.Category == prioritization.KanoPerformance
}

// === Combined Prioritization ===

// PrioritizationSummary returns a summary of all prioritization data.
type PrioritizationSummary struct {
	FeatureID     string                        `json:"featureId"`
	FeatureName   string                        `json:"featureName"`
	RICEScore     float64                       `json:"riceScore,omitempty"`
	RICEComplete  bool                          `json:"riceComplete"`
	KanoCategory  prioritization.KanoCategory   `json:"kanoCategory,omitempty"`
	KanoComplete  bool                          `json:"kanoComplete"`
	Recommendation string                       `json:"recommendation,omitempty"`
	Confidence    string                        `json:"confidence,omitempty"`
}

// GetPrioritizationSummary returns a summary of all prioritization data.
func (s *OpportunitySpec) GetPrioritizationSummary() PrioritizationSummary {
	summary := PrioritizationSummary{
		FeatureID:      s.Metadata.ID,
		FeatureName:    s.Metadata.Title,
		RICEComplete:   s.HasRICE(),
		KanoComplete:   s.HasKano(),
		Recommendation: s.Recommendation.Decision,
		Confidence:     s.Recommendation.Confidence,
	}

	if s.HasRICE() {
		summary.RICEScore = s.GetRICEScore()
	}

	if s.HasKano() {
		summary.KanoCategory = s.GetKanoCategory()
	}

	return summary
}
