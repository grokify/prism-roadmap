package canvas

// OpportunityAssessment follows Marty Cagan's SVPG 10-question framework
// for evaluating product opportunities before committing resources.
//
// Reference: https://www.svpg.com/assessing-product-opportunities/
//
// Grid layout (2 columns x 5 rows):
//
//	| 1. Value Proposition     | 2. Target Market         |
//	| 3. Market Size           | 4. Competitive Landscape |
//	| 5. Our Differentiator    | 6. Market Window         |
//	| 7. Go-to-Market Strategy | 8. Metrics/Revenue       |
//	| 9. Solution Requirements | 10. Recommendation       |
type OpportunityAssessment struct {
	Metadata Metadata `json:"metadata"`

	// Box 1: Value Proposition
	// "Exactly what problem will this solve?"
	ValueProposition OAValueProposition `json:"valueProposition"`

	// Box 2: Target Market
	// "For whom do we solve that problem?"
	TargetMarket OATargetMarket `json:"targetMarket"`

	// Box 3: Market Size
	// "How big is the opportunity?"
	MarketSize OAMarketSize `json:"marketSize"`

	// Box 4: Competitive Landscape
	// "What alternatives are out there?"
	CompetitiveLandscape OACompetitiveLandscape `json:"competitiveLandscape"`

	// Box 5: Our Differentiator
	// "Why are we best suited to pursue this?"
	Differentiator OADifferentiator `json:"differentiator"`

	// Box 6: Market Window
	// "Why now?"
	MarketWindow OAMarketWindow `json:"marketWindow"`

	// Box 7: Go-to-Market Strategy
	// "How will we get this product to market?"
	GoToMarket OAGoToMarketStrategy `json:"goToMarket"`

	// Box 8: Metrics/Revenue Strategy
	// "How will we measure success/make money from this product?"
	MetricsRevenue OAMetricsRevenue `json:"metricsRevenue"`

	// Box 9: Solution Requirements
	// "What factors are critical to success?"
	SolutionRequirements OASolutionRequirements `json:"solutionRequirements"`

	// Box 10: Recommendation
	// "Given the above, what's the recommendation?"
	Recommendation OARecommendation `json:"recommendation"`

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// OAValueProposition describes the problem being solved (Box 1).
type OAValueProposition struct {
	ProblemStatement string   `json:"problemStatement"`          // Clear statement of the problem
	CustomerPains    []string `json:"customerPains,omitempty"`   // Specific pain points addressed
	DesiredOutcomes  []string `json:"desiredOutcomes,omitempty"` // What customers want to achieve
	CurrentState     string   `json:"currentState,omitempty"`    // How things work today
	FutureState      string   `json:"futureState,omitempty"`     // How things will work with solution
}

// OATargetMarket describes who has the problem (Box 2).
type OATargetMarket struct {
	PrimarySegment    string           `json:"primarySegment"`              // Main target segment
	SecondarySegments []string         `json:"secondarySegments,omitempty"` // Additional segments
	Personas          []OATargetPerson `json:"personas,omitempty"`          // Specific user personas
	Industries        []string         `json:"industries,omitempty"`        // Target industries
	Geography         []string         `json:"geography,omitempty"`         // Geographic focus
	CompanySize       string           `json:"companySize,omitempty"`       // SMB, Mid-market, Enterprise
}

// OATargetPerson describes a target user persona.
type OATargetPerson struct {
	Name         string   `json:"name"`                   // Persona name/title
	Role         string   `json:"role,omitempty"`         // Job role
	Goals        []string `json:"goals,omitempty"`        // What they're trying to achieve
	Frustrations []string `json:"frustrations,omitempty"` // Current pain points
	PersonaRef   string   `json:"personaRef,omitempty"`   // Link to PRD persona
}

// OAMarketSize describes the opportunity size (Box 3).
type OAMarketSize struct {
	TAM          string `json:"tam,omitempty"`          // Total Addressable Market
	SAM          string `json:"sam,omitempty"`          // Serviceable Addressable Market
	SOM          string `json:"som,omitempty"`          // Serviceable Obtainable Market
	GrowthRate   string `json:"growthRate,omitempty"`   // Market growth rate
	MarketTrends string `json:"marketTrends,omitempty"` // Key trends affecting size
	DataSources  string `json:"dataSources,omitempty"`  // Where estimates come from
	Assumptions  string `json:"assumptions,omitempty"`  // Key assumptions in sizing
}

// OACompetitiveLandscape describes alternatives in the market (Box 4).
type OACompetitiveLandscape struct {
	DirectCompetitors   []OACompetitor `json:"directCompetitors,omitempty"`   // Head-to-head competitors
	IndirectCompetitors []OACompetitor `json:"indirectCompetitors,omitempty"` // Alternative solutions
	Substitutes         []string       `json:"substitutes,omitempty"`         // Non-product alternatives (manual, etc.)
	MarketDynamics      string         `json:"marketDynamics,omitempty"`      // Competitive dynamics
	EntryBarriers       []string       `json:"entryBarriers,omitempty"`       // Barriers to entry
}

// OACompetitor describes a competitive alternative.
type OACompetitor struct {
	Name        string   `json:"name"`
	Strengths   []string `json:"strengths,omitempty"`
	Weaknesses  []string `json:"weaknesses,omitempty"`
	Positioning string   `json:"positioning,omitempty"` // How they position themselves
	MarketShare string   `json:"marketShare,omitempty"` // Estimated market share
}

// OADifferentiator describes why we're best suited (Box 5).
type OADifferentiator struct {
	CoreStrengths      []string `json:"coreStrengths,omitempty"`      // What we do well
	UniqueCapabilities []string `json:"uniqueCapabilities,omitempty"` // What only we can do
	StrategicAssets    []string `json:"strategicAssets,omitempty"`    // Assets we can leverage
	TeamExpertise      string   `json:"teamExpertise,omitempty"`      // Relevant team experience
	TechnologyEdge     string   `json:"technologyEdge,omitempty"`     // Technical advantages
	UnfairAdvantage    string   `json:"unfairAdvantage,omitempty"`    // What's hard to replicate
}

// OAMarketWindow describes timing considerations (Box 6).
type OAMarketWindow struct {
	WhyNow            string   `json:"whyNow"`                      // Primary reason for timing
	MarketTriggers    []string `json:"marketTriggers,omitempty"`    // Events creating opportunity
	TechnologyShifts  []string `json:"technologyShifts,omitempty"`  // Tech changes enabling this
	RegulatoryChanges []string `json:"regulatoryChanges,omitempty"` // Regulatory drivers
	CompetitorMoves   []string `json:"competitorMoves,omitempty"`   // Competitor timing
	WindowDuration    string   `json:"windowDuration,omitempty"`    // How long window is open
	UrgencyLevel      string   `json:"urgencyLevel,omitempty"`      // high, medium, low
}

// OAGoToMarketStrategy describes how to reach customers (Box 7).
type OAGoToMarketStrategy struct {
	Strategy            string   `json:"strategy"`                      // Overall GTM approach
	Channels            []string `json:"channels,omitempty"`            // Distribution channels
	SalesModel          string   `json:"salesModel,omitempty"`          // Self-serve, sales-led, PLG, etc.
	Partnerships        []string `json:"partnerships,omitempty"`        // Key partnerships needed
	LaunchApproach      string   `json:"launchApproach,omitempty"`      // How to launch (beta, GA, etc.)
	CustomerAcquisition string   `json:"customerAcquisition,omitempty"` // How to acquire customers
	EstimatedCAC        string   `json:"estimatedCAC,omitempty"`        // Customer acquisition cost
}

// OAMetricsRevenue describes success measures and business model (Box 8).
type OAMetricsRevenue struct {
	// Success Metrics
	PrimaryMetric     string   `json:"primaryMetric,omitempty"`     // North star metric
	SecondaryMetrics  []string `json:"secondaryMetrics,omitempty"`  // Supporting metrics
	LeadingIndicators []string `json:"leadingIndicators,omitempty"` // Early signals of success

	// Revenue Strategy
	RevenueModel     string `json:"revenueModel,omitempty"`     // Subscription, usage, etc.
	PricingStrategy  string `json:"pricingStrategy,omitempty"`  // How to price
	EstimatedRevenue string `json:"estimatedRevenue,omitempty"` // Revenue potential
	EstimatedMargin  string `json:"estimatedMargin,omitempty"`  // Margin expectations
	TimeToRevenue    string `json:"timeToRevenue,omitempty"`    // When revenue starts
}

// OASolutionRequirements describes critical success factors (Box 9).
type OASolutionRequirements struct {
	MustHaveCapabilities  []string `json:"mustHaveCapabilities,omitempty"`  // Non-negotiable features
	TechnicalRequirements []string `json:"technicalRequirements,omitempty"` // Technical needs
	IntegrationNeeds      []string `json:"integrationNeeds,omitempty"`      // Required integrations
	ComplianceNeeds       []string `json:"complianceNeeds,omitempty"`       // Regulatory/compliance
	ScalabilityNeeds      string   `json:"scalabilityNeeds,omitempty"`      // Scale requirements
	PerformanceNeeds      string   `json:"performanceNeeds,omitempty"`      // Performance requirements
	KeyDependencies       []string `json:"keyDependencies,omitempty"`       // External dependencies
	CriticalRisks         []string `json:"criticalRisks,omitempty"`         // Risks to mitigate
}

// OARecommendation provides the go/no-go decision (Box 10).
type OARecommendation struct {
	Decision       string   `json:"decision"`                 // go, no-go, conditional, defer
	Rationale      string   `json:"rationale"`                // Why this decision
	Confidence     string   `json:"confidence,omitempty"`     // high, medium, low
	KeyAssumptions []string `json:"keyAssumptions,omitempty"` // Assumptions behind decision
	NextSteps      []string `json:"nextSteps,omitempty"`      // What to do next
	Conditions     []string `json:"conditions,omitempty"`     // Conditions for conditional go
	ReviewDate     string   `json:"reviewDate,omitempty"`     // When to revisit if deferred

	// Investment ask
	TimelineEstimate  string `json:"timelineEstimate,omitempty"`  // Time to deliver
	ResourcesRequired string `json:"resourcesRequired,omitempty"` // Team/resources needed
	InvestmentAsk     string `json:"investmentAsk,omitempty"`     // Budget requested
}

// NewOpportunityAssessment creates a new OpportunityAssessment with defaults.
func NewOpportunityAssessment(id, title string) *OpportunityAssessment {
	return &OpportunityAssessment{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionOpportunityAssessment1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *OpportunityAssessment) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// IsGo returns true if the recommendation is to proceed.
func (c *OpportunityAssessment) IsGo() bool {
	return c.Recommendation.Decision == "go"
}

// IsNoGo returns true if the recommendation is not to proceed.
func (c *OpportunityAssessment) IsNoGo() bool {
	return c.Recommendation.Decision == "no-go"
}

// IsConditional returns true if the recommendation is conditional.
func (c *OpportunityAssessment) IsConditional() bool {
	return c.Recommendation.Decision == "conditional"
}

// IsDeferred returns true if the decision is deferred.
func (c *OpportunityAssessment) IsDeferred() bool {
	return c.Recommendation.Decision == "defer"
}

// HasHighConfidence returns true if confidence is high.
func (c *OpportunityAssessment) HasHighConfidence() bool {
	return c.Recommendation.Confidence == "high"
}

// CompetitorCount returns the total number of competitors identified.
func (c *OpportunityAssessment) CompetitorCount() int {
	return len(c.CompetitiveLandscape.DirectCompetitors) +
		len(c.CompetitiveLandscape.IndirectCompetitors)
}
