// Package mrd provides data types for structured Market Requirements Documents.
package mrd

import (
	"time"

	"github.com/grokify/prism-execution/common"
)

// Person is an alias for common.Person for backwards compatibility.
type Person = common.Person

// Approver is an alias for common.Approver for backwards compatibility.
type Approver = common.Approver

// Status is an alias for common.Status for backwards compatibility.
type Status = common.Status

// GlossaryTerm is an alias for common.GlossaryTerm for backwards compatibility.
type GlossaryTerm = common.GlossaryTerm

// CustomSection is an alias for common.CustomSection for backwards compatibility.
type CustomSection = common.CustomSection

// Status constants re-exported from common for backward compatibility.
const (
	StatusDraft      = common.StatusDraft
	StatusInReview   = common.StatusInReview
	StatusApproved   = common.StatusApproved
	StatusDeprecated = common.StatusDeprecated
)

// Document represents a complete Market Requirements Document.
type Document struct {
	Metadata             Metadata             `json:"metadata"`
	ExecutiveSummary     ExecutiveSummary     `json:"executiveSummary"`
	MarketOverview       MarketOverview       `json:"marketOverview"`
	TargetMarket         TargetMarket         `json:"targetMarket"`
	CompetitiveLandscape CompetitiveLandscape `json:"competitiveLandscape"`
	MarketRequirements   []MarketRequirement  `json:"marketRequirements"`
	Positioning          Positioning          `json:"positioning"`
	GoToMarket           *GoToMarket          `json:"goToMarket,omitempty"`
	SuccessMetrics       []SuccessMetric      `json:"successMetrics"`

	// Optional sections
	Risks          []Risk          `json:"risks,omitempty"`
	Assumptions    []Assumption    `json:"assumptions,omitempty"`
	Glossary       []GlossaryTerm  `json:"glossary,omitempty"`
	CustomSections []CustomSection `json:"customSections,omitempty"`
}

// Note: Status type and constants are defined in common/ and aliased above.

// Metadata contains document metadata.
type Metadata struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Version   string     `json:"version"`
	Status    Status     `json:"status"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	Authors   []Person   `json:"authors"`
	Reviewers []Person   `json:"reviewers,omitempty"`
	Approvers []Approver `json:"approvers,omitempty"`
	Tags      []string   `json:"tags,omitempty"`
}

// ExecutiveSummary provides high-level market overview.
type ExecutiveSummary struct {
	MarketOpportunity string   `json:"marketOpportunity"`
	ProposedOffering  string   `json:"proposedOffering"`
	KeyFindings       []string `json:"keyFindings"`
	Recommendation    string   `json:"recommendation,omitempty"`
}

// MarketOverview contains market size and growth analysis.
type MarketOverview struct {
	TAM         MarketSize `json:"tam"`                   // Total Addressable Market
	SAM         MarketSize `json:"sam"`                   // Serviceable Addressable Market
	SOM         MarketSize `json:"som"`                   // Serviceable Obtainable Market
	GrowthRate  string     `json:"growthRate,omitempty"`  // e.g., "46.3% CAGR"
	MarketStage string     `json:"marketStage,omitempty"` // Emerging, Growth, Mature, Declining
	Trends      []Trend    `json:"trends,omitempty"`
	Drivers     []string   `json:"drivers,omitempty"`  // What's driving growth
	Barriers    []string   `json:"barriers,omitempty"` // Barriers to entry
}

// MarketSize represents a market size measurement.
type MarketSize struct {
	Value  string `json:"value"`            // e.g., "$9.5B"
	Year   int    `json:"year,omitempty"`   // Reference year
	Source string `json:"source,omitempty"` // Citation
	Notes  string `json:"notes,omitempty"`
}

// Trend represents a market trend.
type Trend struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Impact      string `json:"impact,omitempty"`    // High, Medium, Low
	Timeframe   string `json:"timeframe,omitempty"` // Near-term, Mid-term, Long-term
}

// TargetMarket defines the target market segments.
type TargetMarket struct {
	PrimarySegments   []MarketSegment `json:"primarySegments"`
	SecondarySegments []MarketSegment `json:"secondarySegments,omitempty"`
	BuyerPersonas     []BuyerPersona  `json:"buyerPersonas,omitempty"`
	Verticals         []string        `json:"verticals,omitempty"`       // Industry verticals
	GeographicFocus   []string        `json:"geographicFocus,omitempty"` // Regions
	CompanySize       []string        `json:"companySize,omitempty"`     // SMB, Mid-Market, Enterprise
}

// MarketSegment represents a market segment.
type MarketSegment struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Size        string   `json:"size,omitempty"`       // Segment size
	Growth      string   `json:"growth,omitempty"`     // Segment growth rate
	Needs       []string `json:"needs,omitempty"`      // Key needs
	Challenges  []string `json:"challenges,omitempty"` // Key challenges
	Tags        []string `json:"tags,omitempty"`       // For filtering by topic/domain
}

// BuyerPersona represents a market-focused buyer persona.
type BuyerPersona struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Title              string   `json:"title"` // Job title
	Description        string   `json:"description"`
	BuyingRole         string   `json:"buyingRole"` // Decision Maker, Influencer, User, Gatekeeper
	BudgetAuthority    bool     `json:"budgetAuthority"`
	PainPoints         []string `json:"painPoints"`
	Goals              []string `json:"goals"`
	BuyingCriteria     []string `json:"buyingCriteria,omitempty"`
	InformationSources []string `json:"informationSources,omitempty"` // Where they get info
	Tags               []string `json:"tags,omitempty"`               // For filtering by topic/domain
}

// CompetitiveLandscape contains competitive analysis.
type CompetitiveLandscape struct {
	Overview        string       `json:"overview"`
	Competitors     []Competitor `json:"competitors"`
	MarketPosition  string       `json:"marketPosition,omitempty"`  // Our current position
	Differentiators []string     `json:"differentiators,omitempty"` // Key differentiators
	CompetitiveGaps []string     `json:"competitiveGaps,omitempty"` // Gaps to address
}

// Competitor represents a competitor analysis.
type Competitor struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Category    string   `json:"category,omitempty"` // Direct, Indirect, Substitute
	MarketShare string   `json:"marketShare,omitempty"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
	Pricing     string   `json:"pricing,omitempty"`
	Positioning string   `json:"positioning,omitempty"`
	ThreatLevel string   `json:"threatLevel,omitempty"` // High, Medium, Low
	Tags        []string `json:"tags,omitempty"`        // For filtering by topic/domain
}

// MarketRequirement represents a market-level requirement.
type MarketRequirement struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Priority    Priority `json:"priority"`
	Category    string   `json:"category,omitempty"`   // Capability, Integration, Support, etc.
	Source      string   `json:"source,omitempty"`     // Customer feedback, competitor analysis, etc.
	Validation  string   `json:"validation,omitempty"` // How this was validated
	Segments    []string `json:"segments,omitempty"`   // Which segments need this
	Personas    []string `json:"personas,omitempty"`   // Which personas need this
	Tags        []string `json:"tags,omitempty"`       // For filtering by topic/domain
}

// Priority represents requirement priority.
type Priority string

const (
	PriorityMust   Priority = "must"
	PriorityShould Priority = "should"
	PriorityCould  Priority = "could"
	PriorityWont   Priority = "wont"
)

// Positioning defines market positioning strategy.
type Positioning struct {
	Statement       string   `json:"statement"` // Positioning statement
	TargetAudience  string   `json:"targetAudience"`
	Category        string   `json:"category"` // Market category
	KeyBenefits     []string `json:"keyBenefits"`
	Differentiators []string `json:"differentiators"`
	ProofPoints     []string `json:"proofPoints,omitempty"` // Evidence supporting claims
	Tagline         string   `json:"tagline,omitempty"`
}

// GoToMarket contains go-to-market strategy elements.
type GoToMarket struct {
	LaunchStrategy       string           `json:"launchStrategy,omitempty"`
	LaunchTiming         string           `json:"launchTiming,omitempty"`
	PricingStrategy      *PricingStrategy `json:"pricingStrategy,omitempty"`
	DistributionChannels []string         `json:"distributionChannels,omitempty"`
	PartnerStrategy      string           `json:"partnerStrategy,omitempty"`
	MarketingStrategy    string           `json:"marketingStrategy,omitempty"`
	SalesStrategy        string           `json:"salesStrategy,omitempty"`
	Milestones           []Milestone      `json:"milestones,omitempty"`
}

// PricingStrategy defines pricing approach.
type PricingStrategy struct {
	Model       string        `json:"model"` // Subscription, Usage, Perpetual, Freemium
	Tiers       []PricingTier `json:"tiers,omitempty"`
	Positioning string        `json:"positioning,omitempty"` // Premium, Mid-market, Value
	Rationale   string        `json:"rationale,omitempty"`
}

// PricingTier represents a pricing tier.
type PricingTier struct {
	Name        string   `json:"name"`
	Price       string   `json:"price"`
	Billing     string   `json:"billing,omitempty"` // Monthly, Annual
	Features    []string `json:"features,omitempty"`
	TargetBuyer string   `json:"targetBuyer,omitempty"`
}

// Milestone represents a go-to-market milestone.
type Milestone struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	TargetDate  time.Time `json:"targetDate,omitempty"`
	Status      string    `json:"status,omitempty"`
	Tags        []string  `json:"tags,omitempty"` // For filtering by topic/domain
}

// SuccessMetric defines market success metrics.
type SuccessMetric struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Metric            string   `json:"metric"`
	Target            string   `json:"target"`
	Timeframe         string   `json:"timeframe,omitempty"`
	MeasurementMethod string   `json:"measurementMethod,omitempty"`
	Tags              []string `json:"tags,omitempty"` // For filtering by topic/domain
}

// Risk represents a market risk.
type Risk struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Probability string   `json:"probability"` // Low, Medium, High
	Impact      string   `json:"impact"`      // Low, Medium, High, Critical
	Mitigation  string   `json:"mitigation"`
	Category    string   `json:"category,omitempty"` // Market, Competitive, Regulatory, etc.
	Tags        []string `json:"tags,omitempty"`     // For filtering by topic/domain
}

// Assumption represents a market assumption.
type Assumption struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Rationale   string `json:"rationale,omitempty"`
	Validated   bool   `json:"validated,omitempty"`
	Risk        string `json:"risk,omitempty"` // What if wrong
}

// Note: GlossaryTerm and CustomSection types are defined in common/ and aliased above.
