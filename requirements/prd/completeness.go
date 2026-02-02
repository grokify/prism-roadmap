package prd

import (
	"fmt"
	"strings"
)

// CompletenessReport contains the results of a PRD completeness check.
type CompletenessReport struct {
	OverallScore     float64          `json:"overallScore"`     // 0-100
	Grade            string           `json:"grade"`            // A, B, C, D, F
	Sections         []SectionScore   `json:"sections"`         // Per-section scores
	Recommendations  []Recommendation `json:"recommendations"`  // Improvement suggestions
	Summary          string           `json:"summary"`          // Human-readable summary
	RequiredComplete int              `json:"requiredComplete"` // Count of complete required sections
	RequiredTotal    int              `json:"requiredTotal"`    // Total required sections
	OptionalComplete int              `json:"optionalComplete"` // Count of complete optional sections
	OptionalTotal    int              `json:"optionalTotal"`    // Total optional sections
}

// SectionScore represents the completeness score for a document section.
type SectionScore struct {
	Name        string   `json:"name"`
	Score       float64  `json:"score"`     // 0-100
	MaxPoints   float64  `json:"maxPoints"` // Weight in overall score
	Required    bool     `json:"required"`
	Status      string   `json:"status"` // complete, partial, missing
	Issues      []string `json:"issues,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// Recommendation provides specific guidance on improving the PRD.
type Recommendation struct {
	Section  string            `json:"section"`
	Priority RecommendPriority `json:"priority"` // critical, high, medium, low
	Message  string            `json:"message"`
	Guidance string            `json:"guidance,omitempty"`
}

// RecommendPriority represents recommendation priority levels.
type RecommendPriority string

const (
	RecommendCritical RecommendPriority = "critical"
	RecommendHigh     RecommendPriority = "high"
	RecommendMedium   RecommendPriority = "medium"
	RecommendLow      RecommendPriority = "low"
)

// CheckCompleteness analyzes the PRD and returns a completeness report.
func (d *Document) CheckCompleteness() CompletenessReport {
	report := CompletenessReport{
		RequiredTotal: 7,
		OptionalTotal: 8,
	}

	// Check each section
	report.Sections = append(report.Sections, d.checkMetadata())
	report.Sections = append(report.Sections, d.checkExecutiveSummary())
	report.Sections = append(report.Sections, d.checkObjectives())
	report.Sections = append(report.Sections, d.checkPersonas())
	report.Sections = append(report.Sections, d.checkUserStories())
	report.Sections = append(report.Sections, d.checkRequirements())
	report.Sections = append(report.Sections, d.checkRoadmap())

	// Optional sections
	report.Sections = append(report.Sections, d.checkAssumptions())
	report.Sections = append(report.Sections, d.checkInScope())
	report.Sections = append(report.Sections, d.checkOutOfScope())
	report.Sections = append(report.Sections, d.checkTechnicalArchitecture())
	report.Sections = append(report.Sections, d.checkUXRequirements())
	report.Sections = append(report.Sections, d.checkRisks())
	report.Sections = append(report.Sections, d.checkGlossary())
	report.Sections = append(report.Sections, d.checkRelatedDocuments())

	// Calculate overall score
	var totalPoints, earnedPoints float64
	for _, section := range report.Sections {
		totalPoints += section.MaxPoints
		earnedPoints += (section.Score / 100) * section.MaxPoints

		if section.Status == "complete" {
			if section.Required {
				report.RequiredComplete++
			} else {
				report.OptionalComplete++
			}
		}

		// Collect recommendations from issues
		for _, issue := range section.Issues {
			priority := RecommendMedium
			if section.Required && section.Score < 50 {
				priority = RecommendCritical
			} else if section.Required {
				priority = RecommendHigh
			}
			report.Recommendations = append(report.Recommendations, Recommendation{
				Section:  section.Name,
				Priority: priority,
				Message:  issue,
			})
		}

		// Add suggestions as low priority recommendations
		for _, suggestion := range section.Suggestions {
			report.Recommendations = append(report.Recommendations, Recommendation{
				Section:  section.Name,
				Priority: RecommendLow,
				Message:  suggestion,
			})
		}
	}

	report.OverallScore = (earnedPoints / totalPoints) * 100
	report.Grade = scoreToGrade(report.OverallScore)
	report.Summary = generateSummary(report)

	return report
}

func scoreToGrade(score float64) string {
	switch {
	case score >= 90:
		return "A"
	case score >= 80:
		return "B"
	case score >= 70:
		return "C"
	case score >= 60:
		return "D"
	default:
		return "F"
	}
}

func generateSummary(report CompletenessReport) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("PRD Completeness: %.1f%% (Grade: %s)\n", report.OverallScore, report.Grade))
	sb.WriteString(fmt.Sprintf("Required sections: %d/%d complete\n", report.RequiredComplete, report.RequiredTotal))
	sb.WriteString(fmt.Sprintf("Optional sections: %d/%d complete\n", report.OptionalComplete, report.OptionalTotal))

	criticalCount := 0
	for _, r := range report.Recommendations {
		if r.Priority == RecommendCritical {
			criticalCount++
		}
	}
	if criticalCount > 0 {
		sb.WriteString(fmt.Sprintf("\nCritical issues requiring attention: %d", criticalCount))
	}

	return sb.String()
}

func (d *Document) checkMetadata() SectionScore {
	score := SectionScore{
		Name:      "Metadata",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 5.0

	// Check required fields
	if d.Metadata.ID != "" {
		points++
	} else {
		score.Issues = append(score.Issues, "Missing document ID")
	}

	if d.Metadata.Title != "" {
		points++
	} else {
		score.Issues = append(score.Issues, "Missing document title")
	}

	if d.Metadata.Version != "" {
		points++
	} else {
		score.Issues = append(score.Issues, "Missing version number")
	}

	if d.Metadata.Status != "" {
		points++
	} else {
		score.Issues = append(score.Issues, "Missing document status")
	}

	if len(d.Metadata.Authors) > 0 {
		points++
	} else {
		score.Issues = append(score.Issues, "No authors specified")
	}

	// Bonus for optional fields
	if len(d.Metadata.Reviewers) > 0 {
		score.Suggestions = append(score.Suggestions, "")
	} else {
		score.Suggestions = append(score.Suggestions, "Consider adding reviewers for accountability")
	}

	if len(d.Metadata.Tags) > 0 {
		score.Suggestions = append(score.Suggestions, "")
	} else {
		score.Suggestions = append(score.Suggestions, "Consider adding tags for discoverability")
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	// Clean up empty suggestions
	cleanSuggestions := []string{}
	for _, s := range score.Suggestions {
		if s != "" {
			cleanSuggestions = append(cleanSuggestions, s)
		}
	}
	score.Suggestions = cleanSuggestions

	return score
}

func (d *Document) checkExecutiveSummary() SectionScore {
	score := SectionScore{
		Name:      "Executive Summary",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 5.0

	// Problem statement (weighted more heavily)
	if d.ExecutiveSummary.ProblemStatement != "" {
		if len(d.ExecutiveSummary.ProblemStatement) >= 100 {
			points += 1.5
		} else {
			points += 0.75
			score.Issues = append(score.Issues, "Problem statement is brief; consider expanding with more context")
		}
	} else {
		score.Issues = append(score.Issues, "Missing problem statement - this is critical for stakeholder alignment")
	}

	// Proposed solution
	if d.ExecutiveSummary.ProposedSolution != "" {
		if len(d.ExecutiveSummary.ProposedSolution) >= 100 {
			points += 1.5
		} else {
			points += 0.75
			score.Issues = append(score.Issues, "Proposed solution is brief; consider adding more detail")
		}
	} else {
		score.Issues = append(score.Issues, "Missing proposed solution")
	}

	// Expected outcomes
	if len(d.ExecutiveSummary.ExpectedOutcomes) >= 3 {
		points += 1
	} else if len(d.ExecutiveSummary.ExpectedOutcomes) > 0 {
		points += 0.5
		score.Issues = append(score.Issues, "Consider adding more expected outcomes (recommend 3+)")
	} else {
		score.Issues = append(score.Issues, "Missing expected outcomes")
	}

	// Optional but valuable
	if d.ExecutiveSummary.ValueProposition != "" {
		points += 0.5
	} else {
		score.Suggestions = append(score.Suggestions, "Consider adding a value proposition")
	}

	if d.ExecutiveSummary.TargetAudience != "" {
		points += 0.5
	} else {
		score.Suggestions = append(score.Suggestions, "Consider specifying the target audience")
	}

	score.Score = (points / maxPoints) * 100
	if score.Score > 100 {
		score.Score = 100
	}
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkObjectives() SectionScore {
	score := SectionScore{
		Name:      "Objectives",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 6.0

	// Check OKRs
	numOKRs := len(d.Objectives.OKRs)
	if numOKRs >= 2 {
		points += 2
	} else if numOKRs > 0 {
		points += 1
		score.Issues = append(score.Issues, "Consider adding more OKRs (recommend 2+)")
	} else {
		score.Issues = append(score.Issues, "Missing OKRs - define objectives and key results")
	}

	// Count total key results
	totalKeyResults := 0
	for _, okr := range d.Objectives.OKRs {
		totalKeyResults += len(okr.KeyResults)
	}

	// Check key results
	if totalKeyResults >= 3 {
		points += 2
	} else if totalKeyResults > 0 {
		points += 1
		score.Issues = append(score.Issues, "Consider adding more key results (recommend 3+)")
	} else if numOKRs > 0 {
		score.Issues = append(score.Issues, "OKRs missing key results - how will you measure success?")
	}

	// Check key result quality
	for _, okr := range d.Objectives.OKRs {
		for _, kr := range okr.KeyResults {
			if kr.Target == "" {
				score.Issues = append(score.Issues, fmt.Sprintf("Key result '%s' missing target value", kr.Description))
				break
			}
		}
	}

	// Check for phase targets (roadmap alignment)
	hasPhaseTargets := false
	for _, okr := range d.Objectives.OKRs {
		for _, kr := range okr.KeyResults {
			if len(kr.PhaseTargets) > 0 {
				hasPhaseTargets = true
				break
			}
		}
		if hasPhaseTargets {
			break
		}
	}
	if hasPhaseTargets {
		points += 2
	} else if numOKRs > 0 {
		score.Issues = append(score.Issues, "Consider adding phase targets to key results for roadmap alignment")
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkPersonas() SectionScore {
	score := SectionScore{
		Name:      "Personas",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 6.0

	personaCount := len(d.Personas)

	if personaCount >= 3 {
		points += 2
	} else if personaCount >= 2 {
		points += 1.5
	} else if personaCount == 1 {
		points += 1
		score.Issues = append(score.Issues, "Only one persona defined; consider adding more for broader coverage")
	} else {
		score.Issues = append(score.Issues, "No personas defined - this is critical for user-centered design")
	}

	// Check persona quality
	hasPrimary := false
	for _, persona := range d.Personas {
		if persona.IsPrimary {
			hasPrimary = true
		}

		personaIssues := checkPersonaQuality(persona)
		for _, issue := range personaIssues {
			score.Issues = append(score.Issues, fmt.Sprintf("Persona '%s': %s", persona.Name, issue))
		}
	}

	// Award points for complete personas
	completePersonas := 0
	for _, persona := range d.Personas {
		if isPersonaComplete(persona) {
			completePersonas++
		}
	}
	points += float64(completePersonas) * 1.5
	if points > maxPoints {
		points = maxPoints
	}

	if !hasPrimary && personaCount > 0 {
		score.Suggestions = append(score.Suggestions, "Consider marking one persona as primary")
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func checkPersonaQuality(p Persona) []string {
	var issues []string

	if len(p.Goals) == 0 {
		issues = append(issues, "missing goals")
	}
	if len(p.PainPoints) == 0 {
		issues = append(issues, "missing pain points")
	}
	if p.Description == "" {
		issues = append(issues, "missing description")
	}

	return issues
}

func isPersonaComplete(p Persona) bool {
	return p.ID != "" &&
		p.Name != "" &&
		p.Role != "" &&
		p.Description != "" &&
		len(p.Goals) > 0 &&
		len(p.PainPoints) > 0
}

func (d *Document) checkUserStories() SectionScore {
	score := SectionScore{
		Name:      "User Stories",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 6.0

	storyCount := len(d.UserStories)

	if storyCount >= 10 {
		points += 2
	} else if storyCount >= 5 {
		points += 1.5
	} else if storyCount > 0 {
		points += 1
		score.Issues = append(score.Issues, "Limited user stories; consider adding more for comprehensive coverage")
	} else {
		score.Issues = append(score.Issues, "No user stories defined")
	}

	// Check story quality
	storiesWithAC := 0
	storiesWithPersona := 0
	storiesWithPhase := 0

	personaIDs := make(map[string]bool)
	for _, p := range d.Personas {
		personaIDs[p.ID] = true
	}

	phaseIDs := make(map[string]bool)
	for _, p := range d.Roadmap.Phases {
		phaseIDs[p.ID] = true
	}

	for _, story := range d.UserStories {
		if len(story.AcceptanceCriteria) > 0 {
			storiesWithAC++
		}
		if personaIDs[story.PersonaID] {
			storiesWithPersona++
		}
		if phaseIDs[story.PhaseID] {
			storiesWithPhase++
		}
	}

	// Points for acceptance criteria
	if storyCount > 0 {
		acRatio := float64(storiesWithAC) / float64(storyCount)
		if acRatio >= 0.9 {
			points += 2
		} else if acRatio >= 0.5 {
			points += 1
			score.Issues = append(score.Issues, fmt.Sprintf("%.0f%% of stories have acceptance criteria; aim for 90%%+", acRatio*100))
		} else {
			score.Issues = append(score.Issues, fmt.Sprintf("Only %.0f%% of stories have acceptance criteria", acRatio*100))
		}

		// Points for persona linkage
		personaRatio := float64(storiesWithPersona) / float64(storyCount)
		if personaRatio >= 0.9 {
			points += 1
		} else if personaRatio >= 0.5 {
			points += 0.5
			score.Issues = append(score.Issues, "Some user stories not linked to valid personas")
		} else {
			score.Issues = append(score.Issues, "Most user stories not linked to valid personas")
		}

		// Points for phase linkage
		phaseRatio := float64(storiesWithPhase) / float64(storyCount)
		if phaseRatio >= 0.9 {
			points += 1
		} else if phaseRatio >= 0.5 {
			points += 0.5
			score.Issues = append(score.Issues, "Some user stories not linked to roadmap phases")
		}
	}

	score.Score = (points / maxPoints) * 100
	if score.Score > 100 {
		score.Score = 100
	}
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkRequirements() SectionScore {
	score := SectionScore{
		Name:      "Requirements",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 6.0

	frCount := len(d.Requirements.Functional)
	nfrCount := len(d.Requirements.NonFunctional)

	// Functional requirements
	if frCount >= 10 {
		points += 2
	} else if frCount >= 5 {
		points += 1.5
	} else if frCount > 0 {
		points += 1
		score.Issues = append(score.Issues, "Limited functional requirements; consider adding more detail")
	} else {
		score.Issues = append(score.Issues, "No functional requirements defined")
	}

	// Non-functional requirements
	if nfrCount >= 5 {
		points += 2
	} else if nfrCount >= 3 {
		points += 1.5
	} else if nfrCount > 0 {
		points += 1
		score.Issues = append(score.Issues, "Limited non-functional requirements; consider performance, security, scalability")
	} else {
		score.Issues = append(score.Issues, "No non-functional requirements defined")
	}

	// Check NFR categories coverage
	categories := make(map[NFRCategory]bool)
	for _, nfr := range d.Requirements.NonFunctional {
		categories[nfr.Category] = true
	}

	essentialCategories := []NFRCategory{NFRPerformance, NFRSecurity, NFRReliability}
	missingEssential := []string{}
	for _, cat := range essentialCategories {
		if !categories[cat] {
			missingEssential = append(missingEssential, string(cat))
		}
	}

	if len(missingEssential) == 0 {
		points += 2
	} else if len(missingEssential) < len(essentialCategories) {
		points += 1
		score.Issues = append(score.Issues, fmt.Sprintf("Missing NFR categories: %s", strings.Join(missingEssential, ", ")))
	} else {
		score.Issues = append(score.Issues, "Missing essential NFR categories: performance, security, reliability")
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkRoadmap() SectionScore {
	score := SectionScore{
		Name:      "Roadmap",
		MaxPoints: 10,
		Required:  true,
	}

	points := 0.0
	maxPoints := 6.0

	phaseCount := len(d.Roadmap.Phases)

	if phaseCount >= 3 {
		points += 2
	} else if phaseCount >= 2 {
		points += 1.5
	} else if phaseCount == 1 {
		points += 1
		score.Issues = append(score.Issues, "Only one phase defined; consider breaking into milestones")
	} else {
		score.Issues = append(score.Issues, "No roadmap phases defined")
	}

	// Check phase quality
	phasesWithDeliverables := 0
	phasesWithSuccessCriteria := 0
	phasesWithGoals := 0

	for _, phase := range d.Roadmap.Phases {
		if len(phase.Deliverables) > 0 {
			phasesWithDeliverables++
		}
		if len(phase.SuccessCriteria) > 0 {
			phasesWithSuccessCriteria++
		}
		if len(phase.Goals) > 0 {
			phasesWithGoals++
		}
	}

	if phaseCount > 0 {
		// Deliverables
		if phasesWithDeliverables == phaseCount {
			points += 1.5
		} else if phasesWithDeliverables > 0 {
			points += 0.75
			score.Issues = append(score.Issues, "Some phases missing deliverables")
		} else {
			score.Issues = append(score.Issues, "Phases missing deliverables")
		}

		// Success criteria
		if phasesWithSuccessCriteria == phaseCount {
			points += 1.5
		} else if phasesWithSuccessCriteria > 0 {
			points += 0.75
			score.Issues = append(score.Issues, "Some phases missing success criteria")
		} else {
			score.Issues = append(score.Issues, "Phases missing success criteria - how will you know when done?")
		}

		// Goals
		if phasesWithGoals == phaseCount {
			points += 1
		} else if phasesWithGoals > 0 {
			points += 0.5
			score.Issues = append(score.Issues, "Some phases missing goals")
		} else {
			score.Issues = append(score.Issues, "Phases missing goals")
		}
	}

	score.Score = (points / maxPoints) * 100
	if score.Score > 100 {
		score.Score = 100
	}
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkAssumptions() SectionScore {
	score := SectionScore{
		Name:      "Assumptions & Constraints",
		MaxPoints: 5,
		Required:  false,
	}

	if d.Assumptions == nil {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider documenting assumptions and constraints")
		return score
	}

	points := 0.0
	maxPoints := 4.0

	if len(d.Assumptions.Assumptions) >= 3 {
		points += 2
	} else if len(d.Assumptions.Assumptions) > 0 {
		points += 1
		score.Suggestions = append(score.Suggestions, "Consider documenting more assumptions")
	}

	if len(d.Assumptions.Constraints) >= 2 {
		points += 1.5
	} else if len(d.Assumptions.Constraints) > 0 {
		points += 0.75
	}

	if len(d.Assumptions.Dependencies) > 0 {
		points += 0.5
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkInScope() SectionScore {
	score := SectionScore{
		Name:      "In Scope",
		MaxPoints: 5,
		Required:  false,
	}

	if len(d.InScope) == 0 {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider documenting what's explicitly in scope for clarity")
		return score
	}

	if len(d.InScope) >= 5 {
		score.Score = 100
	} else if len(d.InScope) >= 3 {
		score.Score = 80
	} else {
		score.Score = 60
	}
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkOutOfScope() SectionScore {
	score := SectionScore{
		Name:      "Out of Scope",
		MaxPoints: 5,
		Required:  false,
	}

	if len(d.OutOfScope) == 0 {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider documenting what's explicitly out of scope to prevent scope creep")
		return score
	}

	if len(d.OutOfScope) >= 5 {
		score.Score = 100
	} else if len(d.OutOfScope) >= 3 {
		score.Score = 80
	} else {
		score.Score = 60
	}
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkRelatedDocuments() SectionScore {
	score := SectionScore{
		Name:      "Related Documents",
		MaxPoints: 5,
		Required:  false,
	}

	if len(d.RelatedDocuments) == 0 {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider linking related PRDs, TRDs, or design docs")
		return score
	}

	points := 0.0
	maxPoints := 4.0

	// Points for having related documents
	if len(d.RelatedDocuments) >= 3 {
		points += 2
	} else {
		points += 1
	}

	// Check document quality - count complete entries
	completeCount := 0
	for _, doc := range d.RelatedDocuments {
		if doc.ID != "" && doc.Title != "" && doc.Type != "" && doc.Relationship != "" {
			completeCount++
		}
	}

	if completeCount == len(d.RelatedDocuments) {
		points += 2
	} else if completeCount > 0 {
		points += 1
		score.Issues = append(score.Issues, "Some related documents missing required fields (id, title, type, relationship)")
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkTechnicalArchitecture() SectionScore {
	score := SectionScore{
		Name:      "Technical Architecture",
		MaxPoints: 5,
		Required:  false,
	}

	if d.TechArchitecture == nil {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider adding technical architecture overview for engineering context")
		return score
	}

	points := 0.0
	maxPoints := 7.5

	if d.TechArchitecture.Overview != "" {
		points += 1
	}

	if len(d.TechArchitecture.IntegrationPoints) > 0 {
		points += 1
	}

	if d.TechArchitecture.SystemDiagram != "" {
		points += 1
	}

	if d.TechArchitecture.SecurityDesign != "" {
		points += 0.5
	}

	if d.TechArchitecture.ScalabilityDesign != "" {
		points += 0.5
	}

	// New fields for microservices architecture
	if len(d.TechArchitecture.Services) > 0 {
		points += 1
	}

	if len(d.TechArchitecture.APIs) > 0 {
		points += 1
	}

	if len(d.TechArchitecture.StorageArchitecture) > 0 {
		points += 0.5
	}

	if d.TechArchitecture.GitOps != nil && d.TechArchitecture.GitOps.Enabled {
		points += 0.5
	}

	if d.TechArchitecture.Orchestration != nil {
		points += 0.5
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkUXRequirements() SectionScore {
	score := SectionScore{
		Name:      "UX Requirements",
		MaxPoints: 5,
		Required:  false,
	}

	if d.UXRequirements == nil {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider adding UX requirements for user-facing products")
		return score
	}

	points := 0.0
	maxPoints := 4.0

	if len(d.UXRequirements.DesignPrinciples) > 0 {
		points += 1
	}

	if len(d.UXRequirements.Wireframes) > 0 {
		points += 1.5
	}

	if len(d.UXRequirements.InteractionFlows) > 0 {
		points += 1
	}

	if d.UXRequirements.Accessibility.Standard != "" {
		points += 0.5
	}

	score.Score = (points / maxPoints) * 100
	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkRisks() SectionScore {
	score := SectionScore{
		Name:      "Risks",
		MaxPoints: 5,
		Required:  false,
	}

	riskCount := len(d.Risks)

	if riskCount == 0 {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider documenting project risks and mitigations")
		return score
	}

	// Check risk quality
	risksWithMitigation := 0
	for _, risk := range d.Risks {
		if risk.Mitigation != "" {
			risksWithMitigation++
		}
	}

	if riskCount >= 5 && risksWithMitigation == riskCount {
		score.Score = 100
	} else if riskCount >= 3 && float64(risksWithMitigation)/float64(riskCount) >= 0.8 {
		score.Score = 80
	} else if riskCount >= 2 {
		score.Score = 60
		if risksWithMitigation < riskCount {
			score.Issues = append(score.Issues, "Some risks missing mitigation strategies")
		}
	} else {
		score.Score = 40
	}

	score.Status = getStatus(score.Score)

	return score
}

func (d *Document) checkGlossary() SectionScore {
	score := SectionScore{
		Name:      "Glossary",
		MaxPoints: 5,
		Required:  false,
	}

	termCount := len(d.Glossary)

	if termCount == 0 {
		score.Score = 0
		score.Status = "missing"
		score.Suggestions = append(score.Suggestions, "Consider adding a glossary for domain-specific terms")
		return score
	}

	if termCount >= 10 {
		score.Score = 100
	} else if termCount >= 5 {
		score.Score = 80
	} else {
		score.Score = 60
	}

	score.Status = getStatus(score.Score)

	return score
}

func getStatus(score float64) string {
	switch {
	case score >= 80:
		return "complete"
	case score >= 40:
		return "partial"
	default:
		return "missing"
	}
}

// FormatReport returns a human-readable string representation of the report.
func (r *CompletenessReport) FormatReport() string {
	var sb strings.Builder

	// Header
	sb.WriteString("=" + strings.Repeat("=", 60) + "\n")
	sb.WriteString("PRD COMPLETENESS REPORT\n")
	sb.WriteString("=" + strings.Repeat("=", 60) + "\n\n")

	// Overall score
	sb.WriteString(fmt.Sprintf("Overall Score: %.1f%% (Grade: %s)\n", r.OverallScore, r.Grade))
	sb.WriteString(fmt.Sprintf("Required Sections: %d/%d complete\n", r.RequiredComplete, r.RequiredTotal))
	sb.WriteString(fmt.Sprintf("Optional Sections: %d/%d complete\n\n", r.OptionalComplete, r.OptionalTotal))

	// Section breakdown
	sb.WriteString("-" + strings.Repeat("-", 60) + "\n")
	sb.WriteString("SECTION BREAKDOWN\n")
	sb.WriteString("-" + strings.Repeat("-", 60) + "\n\n")

	// Required sections first
	sb.WriteString("Required Sections:\n")
	for _, section := range r.Sections {
		if section.Required {
			status := getStatusIcon(section.Status)
			sb.WriteString(fmt.Sprintf("  %s %-25s %5.1f%% (%s)\n",
				status, section.Name, section.Score, section.Status))
		}
	}

	sb.WriteString("\nOptional Sections:\n")
	for _, section := range r.Sections {
		if !section.Required {
			status := getStatusIcon(section.Status)
			sb.WriteString(fmt.Sprintf("  %s %-25s %5.1f%% (%s)\n",
				status, section.Name, section.Score, section.Status))
		}
	}

	// Recommendations
	if len(r.Recommendations) > 0 {
		sb.WriteString("\n" + "-" + strings.Repeat("-", 60) + "\n")
		sb.WriteString("RECOMMENDATIONS\n")
		sb.WriteString("-" + strings.Repeat("-", 60) + "\n\n")

		// Group by priority
		critical := filterByPriority(r.Recommendations, RecommendCritical)
		high := filterByPriority(r.Recommendations, RecommendHigh)
		medium := filterByPriority(r.Recommendations, RecommendMedium)

		if len(critical) > 0 {
			sb.WriteString("CRITICAL (must fix):\n")
			for _, rec := range critical {
				sb.WriteString(fmt.Sprintf("  [!] %s: %s\n", rec.Section, rec.Message))
			}
			sb.WriteString("\n")
		}

		if len(high) > 0 {
			sb.WriteString("HIGH (should fix):\n")
			for _, rec := range high {
				sb.WriteString(fmt.Sprintf("  [*] %s: %s\n", rec.Section, rec.Message))
			}
			sb.WriteString("\n")
		}

		if len(medium) > 0 {
			sb.WriteString("MEDIUM (consider):\n")
			for _, rec := range medium {
				sb.WriteString(fmt.Sprintf("  [-] %s: %s\n", rec.Section, rec.Message))
			}
			sb.WriteString("\n")
		}
	}

	sb.WriteString("=" + strings.Repeat("=", 60) + "\n")

	return sb.String()
}

func getStatusIcon(status string) string {
	switch status {
	case "complete":
		return "[+]"
	case "partial":
		return "[~]"
	default:
		return "[ ]"
	}
}

func filterByPriority(recs []Recommendation, priority RecommendPriority) []Recommendation {
	var filtered []Recommendation
	for _, rec := range recs {
		if rec.Priority == priority {
			filtered = append(filtered, rec)
		}
	}
	return filtered
}
