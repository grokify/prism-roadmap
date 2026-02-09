package prd

// SectionID identifies a section for ordering purposes.
type SectionID string

// Section ID constants for all PRD sections.
const (
	SectionExecutiveSummary  SectionID = "executiveSummary"
	SectionCurrentState      SectionID = "currentState"
	SectionProblem           SectionID = "problem"
	SectionMarket            SectionID = "market"
	SectionSolution          SectionID = "solution"
	SectionSuccessMetrics    SectionID = "successMetrics"
	SectionObjectives        SectionID = "objectives"
	SectionNonGoals          SectionID = "nonGoals"
	SectionPersonas          SectionID = "personas"
	SectionUserStories       SectionID = "userStories"
	SectionFunctionalReqs    SectionID = "functionalRequirements"
	SectionNonFunctionalReqs SectionID = "nonFunctionalRequirements"
	SectionRoadmap           SectionID = "roadmap"
	SectionTechArchitecture  SectionID = "technicalArchitecture"
	SectionSecurityModel     SectionID = "securityModel"
	SectionDecisions         SectionID = "decisions"
	SectionRisks             SectionID = "risks"
	SectionAssumptions       SectionID = "assumptions"
	SectionInScope           SectionID = "inScope"
	SectionOutOfScope        SectionID = "outOfScope"
	SectionOpenItems         SectionID = "openItems"
	SectionRelatedDocuments  SectionID = "relatedDocuments"
	SectionAppendices        SectionID = "appendices"
	SectionGlossary          SectionID = "glossary"
	SectionReviews           SectionID = "reviews"
	SectionRevisionHistory   SectionID = "revisionHistory"
	SectionCustom            SectionID = "custom"
)

// PRDType determines the section ordering template for markdown generation.
// This is distinct from DocumentType which categorizes related documents.
type PRDType string

const (
	PRDTypeDefault   PRDType = ""          // Uses DefaultSectionOrder
	PRDTypeStrategy  PRDType = "strategy"  // Context-first ordering
	PRDTypeFeature   PRDType = "feature"   // User-needs-first ordering
	PRDTypeTechnical PRDType = "technical" // Architecture-focused ordering
)

// SectionDisplayNames maps section IDs to display names for TOC and headers.
var SectionDisplayNames = map[SectionID]string{
	SectionExecutiveSummary:  "Executive Summary",
	SectionCurrentState:      "Current State",
	SectionProblem:           "Problem Definition",
	SectionMarket:            "Market Analysis",
	SectionSolution:          "Solution",
	SectionSuccessMetrics:    "Success Metrics",
	SectionObjectives:        "Objectives and Goals",
	SectionNonGoals:          "Non-Goals",
	SectionPersonas:          "Personas",
	SectionUserStories:       "User Stories",
	SectionFunctionalReqs:    "Functional Requirements",
	SectionNonFunctionalReqs: "Non-Functional Requirements",
	SectionRoadmap:           "Roadmap",
	SectionTechArchitecture:  "Technical Architecture",
	SectionSecurityModel:     "Security Model",
	SectionDecisions:         "Decisions",
	SectionRisks:             "Risk Assessment",
	SectionAssumptions:       "Assumptions and Constraints",
	SectionInScope:           "In Scope",
	SectionOutOfScope:        "Out of Scope",
	SectionOpenItems:         "Open Items",
	SectionRelatedDocuments:  "Related Documents",
	SectionAppendices:        "Appendices",
	SectionGlossary:          "Glossary",
	SectionReviews:           "Reviews",
	SectionRevisionHistory:   "Revision History",
}

// SectionAnchors maps section IDs to markdown anchors for TOC links.
var SectionAnchors = map[SectionID]string{
	SectionExecutiveSummary:  "executive-summary",
	SectionCurrentState:      "current-state",
	SectionProblem:           "problem-definition",
	SectionMarket:            "market-analysis",
	SectionSolution:          "solution",
	SectionSuccessMetrics:    "success-metrics",
	SectionObjectives:        "objectives-and-goals",
	SectionNonGoals:          "non-goals",
	SectionPersonas:          "personas",
	SectionUserStories:       "user-stories",
	SectionFunctionalReqs:    "functional-requirements",
	SectionNonFunctionalReqs: "non-functional-requirements",
	SectionRoadmap:           "roadmap",
	SectionTechArchitecture:  "technical-architecture",
	SectionSecurityModel:     "security-model",
	SectionDecisions:         "decisions",
	SectionRisks:             "risk-assessment",
	SectionAssumptions:       "assumptions-and-constraints",
	SectionInScope:           "in-scope",
	SectionOutOfScope:        "out-of-scope",
	SectionOpenItems:         "open-items",
	SectionRelatedDocuments:  "related-documents",
	SectionAppendices:        "appendices",
	SectionGlossary:          "glossary",
	SectionReviews:           "reviews",
	SectionRevisionHistory:   "revision-history",
}

// DefaultSectionOrder is the current v0.8.0 order for backward compatibility.
var DefaultSectionOrder = []SectionID{
	SectionExecutiveSummary,
	SectionObjectives,
	SectionPersonas,
	SectionUserStories,
	SectionFunctionalReqs,
	SectionNonFunctionalReqs,
	SectionRoadmap,
	SectionTechArchitecture,
	SectionAssumptions,
	SectionInScope,
	SectionOutOfScope,
	SectionRisks,
	SectionOpenItems,
	SectionCurrentState,
	SectionSecurityModel,
	SectionAppendices,
	SectionGlossary,
	SectionRelatedDocuments,
	SectionProblem,
	SectionMarket,
	SectionSolution,
	SectionDecisions,
	SectionReviews,
	SectionRevisionHistory,
	SectionNonGoals,
	SectionSuccessMetrics,
	SectionCustom,
}

// StrategySectionOrder places context (CurrentState, Problem, Market) early.
// Best for strategy PRDs, business cases, and product vision documents.
var StrategySectionOrder = []SectionID{
	SectionExecutiveSummary,
	SectionCurrentState,   // Context first
	SectionProblem,        // Then the problem
	SectionMarket,         // Then the landscape
	SectionSolution,       // Then the approach
	SectionSuccessMetrics, // How we measure success
	SectionObjectives,
	SectionNonGoals,
	SectionPersonas,
	SectionUserStories,
	SectionFunctionalReqs,
	SectionNonFunctionalReqs,
	SectionTechArchitecture,
	SectionSecurityModel,
	SectionRoadmap,
	SectionDecisions,
	SectionRisks,
	SectionAssumptions,
	SectionInScope,
	SectionOutOfScope,
	SectionOpenItems,
	SectionRelatedDocuments,
	SectionAppendices,
	SectionGlossary,
	SectionReviews,
	SectionRevisionHistory,
	SectionCustom,
}

// FeatureSectionOrder focuses on user needs before solution details.
// Best for feature PRDs and user-facing product changes.
var FeatureSectionOrder = []SectionID{
	SectionExecutiveSummary,
	SectionProblem,
	SectionSuccessMetrics,
	SectionObjectives,
	SectionPersonas,
	SectionUserStories,
	SectionFunctionalReqs,
	SectionNonFunctionalReqs,
	SectionSolution,
	SectionTechArchitecture,
	SectionSecurityModel,
	SectionRoadmap,
	SectionNonGoals,
	SectionRisks,
	SectionAssumptions,
	SectionInScope,
	SectionOutOfScope,
	SectionOpenItems,
	SectionDecisions,
	SectionCurrentState,
	SectionMarket,
	SectionRelatedDocuments,
	SectionAppendices,
	SectionGlossary,
	SectionReviews,
	SectionRevisionHistory,
	SectionCustom,
}

// TechnicalSectionOrder emphasizes architecture and technical details.
// Best for technical PRDs, infrastructure changes, and platform work.
var TechnicalSectionOrder = []SectionID{
	SectionExecutiveSummary,
	SectionProblem,
	SectionCurrentState,
	SectionSolution,
	SectionTechArchitecture, // Architecture early
	SectionSecurityModel,
	SectionObjectives,
	SectionSuccessMetrics,
	SectionFunctionalReqs,
	SectionNonFunctionalReqs,
	SectionPersonas,
	SectionUserStories,
	SectionRoadmap,
	SectionNonGoals,
	SectionRisks,
	SectionAssumptions,
	SectionInScope,
	SectionOutOfScope,
	SectionOpenItems,
	SectionDecisions,
	SectionMarket,
	SectionRelatedDocuments,
	SectionAppendices,
	SectionGlossary,
	SectionReviews,
	SectionRevisionHistory,
	SectionCustom,
}

// GetSectionOrder returns the section order for a PRD type.
func GetSectionOrder(prdType PRDType) []SectionID {
	switch prdType {
	case PRDTypeStrategy:
		return StrategySectionOrder
	case PRDTypeFeature:
		return FeatureSectionOrder
	case PRDTypeTechnical:
		return TechnicalSectionOrder
	default:
		return DefaultSectionOrder
	}
}

// AllSectionIDs returns a set of all known section IDs.
func AllSectionIDs() map[SectionID]bool {
	return map[SectionID]bool{
		SectionExecutiveSummary:  true,
		SectionCurrentState:      true,
		SectionProblem:           true,
		SectionMarket:            true,
		SectionSolution:          true,
		SectionSuccessMetrics:    true,
		SectionObjectives:        true,
		SectionNonGoals:          true,
		SectionPersonas:          true,
		SectionUserStories:       true,
		SectionFunctionalReqs:    true,
		SectionNonFunctionalReqs: true,
		SectionRoadmap:           true,
		SectionTechArchitecture:  true,
		SectionSecurityModel:     true,
		SectionDecisions:         true,
		SectionRisks:             true,
		SectionAssumptions:       true,
		SectionInScope:           true,
		SectionOutOfScope:        true,
		SectionOpenItems:         true,
		SectionRelatedDocuments:  true,
		SectionAppendices:        true,
		SectionGlossary:          true,
		SectionReviews:           true,
		SectionRevisionHistory:   true,
		SectionCustom:            true,
	}
}

// ValidateSectionOrder checks if all section IDs in the order are valid.
// Returns a list of invalid section IDs.
func ValidateSectionOrder(order []string) []string {
	known := AllSectionIDs()
	var invalid []string
	for _, id := range order {
		if !known[SectionID(id)] {
			invalid = append(invalid, id)
		}
	}
	return invalid
}

// CompleteSectionOrder takes a partial order and appends any missing sections
// from the template order, returning a complete section list.
func CompleteSectionOrder(partial []SectionID, template []SectionID) []SectionID {
	seen := make(map[SectionID]bool)
	for _, id := range partial {
		seen[id] = true
	}

	result := make([]SectionID, len(partial))
	copy(result, partial)

	for _, id := range template {
		if !seen[id] {
			result = append(result, id)
			seen[id] = true
		}
	}
	return result
}

// ListSections returns all section IDs with their display names, sorted by default order.
func ListSections() []struct {
	ID          SectionID
	DisplayName string
} {
	result := make([]struct {
		ID          SectionID
		DisplayName string
	}, len(DefaultSectionOrder))

	for i, id := range DefaultSectionOrder {
		result[i] = struct {
			ID          SectionID
			DisplayName string
		}{
			ID:          id,
			DisplayName: SectionDisplayNames[id],
		}
	}
	return result
}
