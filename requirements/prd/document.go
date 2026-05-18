// Package prd provides data types for structured Product Requirements Documents.
package prd

import (
	"time"

	"github.com/grokify/prism-execution/common"
	"github.com/grokify/prism-execution/goals"
	"github.com/grokify/prism-execution/goals/okr"
)

// Person is an alias for common.Person for backwards compatibility.
type Person = common.Person

// Approver is an alias for common.Approver for backwards compatibility.
type Approver = common.Approver

// OKR type aliases from goals/okr for backward compatibility.
// These allow existing PRD code to continue using prd.OKR, prd.Objective, etc.
type (
	// OKR represents an Objective with its Key Results.
	OKR = okr.OKR

	// Objective represents a qualitative, inspirational goal.
	Objective = okr.Objective

	// KeyResult represents a measurable outcome that indicates objective achievement.
	KeyResult = okr.KeyResult

	// PhaseTarget represents a Key Result target for a specific roadmap phase.
	PhaseTarget = okr.PhaseTarget
)

// Common type aliases for backward compatibility.
// These allow existing PRD code to continue using prd.Status, prd.Risk, etc.
type (
	// Status represents the document lifecycle status.
	Status = common.Status

	// Risk represents a project risk.
	Risk = common.Risk

	// RiskProbability represents risk probability levels.
	RiskProbability = common.RiskProbability

	// RiskImpact represents risk impact levels.
	RiskImpact = common.RiskImpact

	// RiskStatus represents risk status.
	RiskStatus = common.RiskStatus

	// Assumption represents a condition assumed to be true.
	Assumption = common.Assumption

	// Constraint represents a limitation on the project.
	Constraint = common.Constraint

	// ConstraintType represents types of constraints.
	ConstraintType = common.ConstraintType

	// GlossaryTerm defines a glossary entry.
	GlossaryTerm = common.GlossaryTerm

	// CustomSection allows project-specific sections.
	CustomSection = common.CustomSection

	// OpenItem represents a pending decision or question.
	OpenItem = common.OpenItem

	// Option represents one possible choice for an open item.
	Option = common.Option

	// OpenItemStatus represents the status of an open item.
	OpenItemStatus = common.OpenItemStatus

	// OpenItemResolution documents how an open item was resolved.
	OpenItemResolution = common.OpenItemResolution

	// EffortLevel represents effort estimates.
	EffortLevel = common.EffortLevel

	// RiskLevel represents risk levels for options.
	RiskLevel = common.RiskLevel

	// DecisionRecord documents a completed decision.
	DecisionRecord = common.DecisionRecord

	// DecisionStatus represents the status of a decision.
	DecisionStatus = common.DecisionStatus

	// Priority represents priority levels.
	Priority = common.Priority

	// NonGoal represents an explicit out-of-scope item.
	NonGoal = common.NonGoal
)

// Goals type aliases from goals package for backward compatibility.
type (
	// Goals is a framework-agnostic container for organizational goals.
	Goals = goals.Goals

	// GoalItem represents a high-level goal (Objective or Method).
	GoalItem = goals.GoalItem

	// ResultItem represents a measurable result (Key Result or Measure).
	ResultItem = goals.ResultItem

	// Framework identifies the goal-setting framework in use.
	Framework = goals.Framework
)

// Framework constants re-exported for convenience.
const (
	FrameworkOKR   = goals.FrameworkOKR
	FrameworkV2MOM = goals.FrameworkV2MOM
)

// Document represents a complete Product Requirements Document.
type Document struct {
	Metadata         Metadata         `json:"metadata"`
	ExecutiveSummary ExecutiveSummary `json:"executiveSummary"`
	Objectives       Objectives       `json:"objectives"`
	Personas         []Persona        `json:"personas"`
	UserStories      []UserStory      `json:"userStories"`
	Requirements     Requirements     `json:"requirements"`
	Roadmap          Roadmap          `json:"roadmap"`

	// ProductGoals contains the product goals using the framework-agnostic Goals wrapper.
	// This supports either OKR or V2MOM frameworks. When set, this takes precedence
	// over the legacy Objectives field for roadmap rendering and other goal-related features.
	ProductGoals *Goals `json:"productGoals,omitempty"`

	// Optional sections
	Assumptions      *AssumptionsConstraints `json:"assumptions,omitempty"`
	InScope          []string                `json:"inScope,omitempty"`
	OutOfScope       []string                `json:"outOfScope,omitempty"`
	TechArchitecture *TechnicalArchitecture  `json:"technicalArchitecture,omitempty"`

	// RelatedDocuments lists related requirements documents.
	RelatedDocuments []RelatedDocument `json:"relatedDocuments,omitempty"`
	UXRequirements   *UXRequirements   `json:"uxRequirements,omitempty"`
	Risks            []Risk            `json:"risks,omitempty"`
	Glossary         []GlossaryTerm    `json:"glossary,omitempty"`

	// Custom sections for project-specific needs
	CustomSections []CustomSection `json:"customSections,omitempty"`

	// Extended sections (from agent-team-prd merge)
	// These provide additional structure for problem definition, market analysis,
	// solution evaluation, decision tracking, and quality reviews.

	// Problem provides detailed problem definition with evidence.
	Problem *ProblemDefinition `json:"problem,omitempty"`

	// Market contains market analysis and competitive landscape.
	Market *MarketDefinition `json:"market,omitempty"`

	// Solution contains solution options and selection rationale.
	Solution *SolutionDefinition `json:"solution,omitempty"`

	// Decisions contains decision records for the PRD.
	Decisions *DecisionsDefinition `json:"decisions,omitempty"`

	// OpenItems contains pending decisions that need resolution.
	OpenItems []OpenItem `json:"openItems,omitempty"`

	// Reviews contains review outcomes and quality assessments.
	Reviews *ReviewsDefinition `json:"reviews,omitempty"`

	// RevisionHistory tracks changes to the PRD over time.
	RevisionHistory []RevisionRecord `json:"revisionHistory,omitempty"`

	// Goals contains alignment with strategic goals (V2MOM, OKR).
	Goals *GoalsAlignment `json:"goals,omitempty"`

	// CurrentState documents the existing state before the proposed solution.
	CurrentState *CurrentState `json:"currentState,omitempty"`

	// SecurityModel documents security architecture and threat model.
	// This section is strongly recommended for all PRDs.
	SecurityModel *SecurityModel `json:"securityModel,omitempty"`

	// Appendices contains supplementary information and domain-specific data.
	Appendices []Appendix `json:"appendices,omitempty"`

	// NonGoals contains structured non-goals with rationale.
	// More structured than OutOfScope for explicit exclusions.
	NonGoals []NonGoal `json:"nonGoals,omitempty"`

	// SuccessMetrics contains structured success metrics organized by type.
	SuccessMetrics *SuccessMetrics `json:"successMetrics,omitempty"`
}

// Status constants re-exported from common for backward compatibility.
const (
	StatusDraft      = common.StatusDraft
	StatusInReview   = common.StatusInReview
	StatusApproved   = common.StatusApproved
	StatusDeprecated = common.StatusDeprecated
)

// Risk constants re-exported from common for backward compatibility.
const (
	RiskProbabilityLow    = common.RiskProbabilityLow
	RiskProbabilityMedium = common.RiskProbabilityMedium
	RiskProbabilityHigh   = common.RiskProbabilityHigh

	RiskImpactLow      = common.RiskImpactLow
	RiskImpactMedium   = common.RiskImpactMedium
	RiskImpactHigh     = common.RiskImpactHigh
	RiskImpactCritical = common.RiskImpactCritical

	RiskStatusOpen      = common.RiskStatusOpen
	RiskStatusMitigated = common.RiskStatusMitigated
	RiskStatusAccepted  = common.RiskStatusAccepted
	RiskStatusClosed    = common.RiskStatusClosed
)

// Constraint type constants re-exported from common for backward compatibility.
const (
	ConstraintTechnical  = common.ConstraintTechnical
	ConstraintBudget     = common.ConstraintBudget
	ConstraintTimeline   = common.ConstraintTimeline
	ConstraintRegulatory = common.ConstraintRegulatory
	ConstraintResource   = common.ConstraintResource
	ConstraintLegal      = common.ConstraintLegal
)

// OpenItem status constants re-exported from common for backward compatibility.
const (
	OpenItemStatusOpen         = common.OpenItemStatusOpen
	OpenItemStatusInDiscussion = common.OpenItemStatusInDiscussion
	OpenItemStatusBlocked      = common.OpenItemStatusBlocked
	OpenItemStatusResolved     = common.OpenItemStatusResolved
	OpenItemStatusDeferred     = common.OpenItemStatusDeferred
)

// Effort level constants re-exported from common for backward compatibility.
const (
	EffortLow    = common.EffortLow
	EffortMedium = common.EffortMedium
	EffortHigh   = common.EffortHigh
)

// Risk level constants re-exported from common for backward compatibility.
const (
	RiskLevelLow    = common.RiskLevelLow
	RiskLevelMedium = common.RiskLevelMedium
	RiskLevelHigh   = common.RiskLevelHigh
)

// Decision status constants re-exported from common for backward compatibility.
const (
	DecisionProposed   = common.DecisionProposed
	DecisionAccepted   = common.DecisionAccepted
	DecisionSuperseded = common.DecisionSuperseded
	DecisionDeprecated = common.DecisionDeprecated
)

// Priority constants re-exported from common for backward compatibility.
const (
	PriorityCritical = common.PriorityCritical
	PriorityHigh     = common.PriorityHigh
	PriorityMedium   = common.PriorityMedium
	PriorityLow      = common.PriorityLow
)

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

	// SemanticVersioning indicates the Version field follows Semantic Versioning (semver.org).
	SemanticVersioning bool `json:"semanticVersioning,omitempty"`

	// PRDType determines the section ordering template for markdown generation.
	// Options: "", "strategy", "feature", "technical"
	// If empty, uses default order for backward compatibility.
	PRDType PRDType `json:"prdType,omitempty"`

	// SectionOrder specifies a custom section order for markdown generation.
	// If provided, overrides the template order from PRDType.
	// Sections not listed are appended in template order.
	// Use section IDs like "executiveSummary", "problem", "market", etc.
	SectionOrder []string `json:"sectionOrder,omitempty"`
}

// ExecutiveSummary provides high-level product overview.
type ExecutiveSummary struct {
	ProblemStatement string   `json:"problemStatement"`
	ProposedSolution string   `json:"proposedSolution"`
	ExpectedOutcomes []string `json:"expectedOutcomes"`
	TargetAudience   string   `json:"targetAudience,omitempty"`
	ValueProposition string   `json:"valueProposition,omitempty"`
}

// Objectives defines business and product goals using OKR structure.
// Deprecated: Use ProductGoals field with goals.Goals wrapper for new PRDs.
type Objectives struct {
	// OKRs contains Objectives and Key Results in nested OKR format.
	OKRs []OKR `json:"okrs"`
}

// GetProductGoals returns the product goals, preferring ProductGoals if set,
// otherwise converting legacy Objectives to Goals format.
// This enables framework-agnostic rendering in roadmaps and other sections.
func (d *Document) GetProductGoals() *Goals {
	if d.ProductGoals != nil {
		return d.ProductGoals
	}
	// Convert legacy Objectives to Goals format
	if len(d.Objectives.OKRs) > 0 {
		okrSet := &okr.OKRSet{OKRs: d.Objectives.OKRs}
		return goals.NewOKR(okrSet)
	}
	return nil
}

// HasProductGoals returns true if the document has product goals defined
// (either in ProductGoals or legacy Objectives).
func (d *Document) HasProductGoals() bool {
	return d.ProductGoals != nil || len(d.Objectives.OKRs) > 0
}

// HasSection returns true if the document has content for the given section.
func (d *Document) HasSection(id SectionID) bool {
	switch id {
	case SectionExecutiveSummary:
		return true // Always present
	case SectionObjectives:
		return true // Always present
	case SectionPersonas:
		return len(d.Personas) > 0
	case SectionUserStories:
		return len(d.UserStories) > 0
	case SectionFunctionalReqs:
		return len(d.Requirements.Functional) > 0
	case SectionNonFunctionalReqs:
		return len(d.Requirements.NonFunctional) > 0
	case SectionComplianceReqs:
		return len(d.Requirements.Compliance) > 0
	case SectionRequirementsByPhase:
		return d.hasRequirementsWithPhases()
	case SectionRoadmap:
		return true // Always present
	case SectionTechArchitecture:
		return d.TechArchitecture != nil
	case SectionAssumptions:
		return d.Assumptions != nil
	case SectionInScope:
		return len(d.InScope) > 0
	case SectionOutOfScope:
		return len(d.OutOfScope) > 0
	case SectionRisks:
		return len(d.Risks) > 0
	case SectionOpenItems:
		return len(d.OpenItems) > 0
	case SectionCurrentState:
		return d.CurrentState != nil
	case SectionSecurityModel:
		return d.SecurityModel != nil
	case SectionAppendices:
		return len(d.Appendices) > 0
	case SectionGlossary:
		return len(d.Glossary) > 0
	case SectionRelatedDocuments:
		return len(d.RelatedDocuments) > 0
	case SectionProblem:
		return d.Problem != nil
	case SectionMarket:
		return d.Market != nil
	case SectionSolution:
		return d.Solution != nil
	case SectionDecisions:
		return d.Decisions != nil && len(d.Decisions.Records) > 0
	case SectionReviews:
		return d.Reviews != nil
	case SectionRevisionHistory:
		return len(d.RevisionHistory) > 0
	case SectionNonGoals:
		return len(d.NonGoals) > 0
	case SectionSuccessMetrics:
		return d.SuccessMetrics != nil
	case SectionCustom:
		return len(d.CustomSections) > 0
	default:
		return false
	}
}

// GetSectionOrder returns the effective section order for this document.
// Priority: 1) Custom SectionOrder, 2) PRDType template, 3) DefaultSectionOrder
func (d *Document) GetSectionOrder() []SectionID {
	// Get base template
	template := GetSectionOrder(d.Metadata.PRDType)

	// If custom order specified, use it with completion
	if len(d.Metadata.SectionOrder) > 0 {
		partial := make([]SectionID, len(d.Metadata.SectionOrder))
		for i, id := range d.Metadata.SectionOrder {
			partial[i] = SectionID(id)
		}
		return CompleteSectionOrder(partial, template)
	}

	return template
}

// GetActiveSections returns only sections that have content, in order.
func (d *Document) GetActiveSections() []SectionID {
	order := d.GetSectionOrder()
	var active []SectionID
	for _, id := range order {
		if d.HasSection(id) {
			active = append(active, id)
		}
	}
	return active
}

// hasRequirementsWithPhases returns true if any requirement has a phase assigned.
func (d *Document) hasRequirementsWithPhases() bool {
	for _, fr := range d.Requirements.Functional {
		if fr.PhaseID != "" {
			return true
		}
	}
	for _, nfr := range d.Requirements.NonFunctional {
		if nfr.PhaseID != "" {
			return true
		}
	}
	for _, cr := range d.Requirements.Compliance {
		if cr.PhaseID != "" {
			return true
		}
	}
	return false
}
