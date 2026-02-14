package prd

import (
	"fmt"
	"strings"
)

// CategoryWeight defines the weight for each scoring category.
type CategoryWeight struct {
	Category string
	Weight   float64
}

// DefaultWeights returns the standard category weights.
// Weights sum to 1.0.
func DefaultWeights() []CategoryWeight {
	return []CategoryWeight{
		{Category: "problem_definition", Weight: 0.20},
		{Category: "solution_fit", Weight: 0.15},
		{Category: "user_understanding", Weight: 0.10},
		{Category: "market_awareness", Weight: 0.10},
		{Category: "scope_discipline", Weight: 0.10},
		{Category: "requirements_quality", Weight: 0.10},
		{Category: "metrics_quality", Weight: 0.10},
		{Category: "ux_coverage", Weight: 0.05},
		{Category: "technical_feasibility", Weight: 0.05},
		{Category: "risk_management", Weight: 0.05},
	}
}

// Thresholds for scoring decisions.
const (
	ThresholdApprove     = 8.0
	ThresholdRevise      = 6.5
	ThresholdHumanReview = 6.5
	ThresholdBlocker     = 3.0
)

// CategoryScore represents a score for a single category.
type CategoryScore struct {
	Category       string  `json:"category"`
	Weight         float64 `json:"weight"`
	Score          float64 `json:"score"`
	MaxScore       float64 `json:"maxScore"`
	Justification  string  `json:"justification"`
	Evidence       string  `json:"evidence"`
	BelowThreshold bool    `json:"belowThreshold"`
}

// ScoringResult contains the complete scoring output.
type ScoringResult struct {
	CategoryScores   []CategoryScore   `json:"categoryScores"`
	WeightedScore    float64           `json:"weightedScore"`
	Decision         string            `json:"decision"`
	Blockers         []string          `json:"blockers"`
	RevisionTriggers []RevisionTrigger `json:"revisionTriggers"`
	Summary          string            `json:"summary"`
}

// Score evaluates a Document and returns scoring results.
func Score(doc *Document) *ScoringResult {
	result := &ScoringResult{
		CategoryScores: make([]CategoryScore, 0),
	}

	weights := DefaultWeights()
	var totalWeightedScore float64
	var totalWeight float64
	revisionID := 1

	for _, w := range weights {
		score := scoreCategory(doc, w.Category)
		score.Weight = w.Weight
		score.MaxScore = 10.0

		if score.Score <= ThresholdBlocker {
			score.BelowThreshold = true
			result.Blockers = append(result.Blockers, fmt.Sprintf("%s: %s", w.Category, score.Justification))
		}

		// Generate revision triggers for low scores
		if score.Score < 7.0 {
			severity := "minor"
			if score.Score < 5.0 {
				severity = "major"
			}
			if score.Score <= ThresholdBlocker {
				severity = "blocker"
			}

			result.RevisionTriggers = append(result.RevisionTriggers, RevisionTrigger{
				IssueID:          fmt.Sprintf("REV-%d", revisionID),
				Category:         w.Category,
				Severity:         severity,
				Description:      score.Justification,
				RecommendedOwner: categoryToOwner(w.Category),
			})
			revisionID++
		}

		result.CategoryScores = append(result.CategoryScores, score)
		totalWeightedScore += score.Score * w.Weight
		totalWeight += w.Weight
	}

	if totalWeight > 0 {
		result.WeightedScore = totalWeightedScore / totalWeight
	}

	// Determine decision
	if len(result.Blockers) > 0 {
		result.Decision = "reject"
	} else if result.WeightedScore >= ThresholdApprove {
		result.Decision = "approve"
	} else if result.WeightedScore >= ThresholdRevise {
		result.Decision = "revise"
	} else {
		result.Decision = "human_review"
	}

	result.Summary = generateScoringSummary(result)

	return result
}

func scoreCategory(doc *Document, category string) CategoryScore {
	switch category {
	case "problem_definition":
		return scoreProblemDefinition(doc)
	case "user_understanding":
		return scoreUserUnderstanding(doc)
	case "market_awareness":
		return scoreMarketAwareness(doc)
	case "solution_fit":
		return scoreSolutionFit(doc)
	case "scope_discipline":
		return scoreScopeDiscipline(doc)
	case "requirements_quality":
		return scoreRequirementsQuality(doc)
	case "ux_coverage":
		return scoreUXCoverage(doc)
	case "technical_feasibility":
		return scoreTechnicalFeasibility(doc)
	case "metrics_quality":
		return scoreMetricsQuality(doc)
	case "risk_management":
		return scoreRiskManagement(doc)
	default:
		return CategoryScore{Category: category, Score: 5.0, Justification: "Unknown category"}
	}
}

func scoreProblemDefinition(doc *Document) CategoryScore {
	score := CategoryScore{Category: "problem_definition"}
	var points float64
	var evidence []string

	// Check problem statement from executive summary or Problem field
	hasStatement := false
	if doc.ExecutiveSummary.ProblemStatement != "" {
		points += 3.0
		evidence = append(evidence, "Problem statement present in executive summary")
		hasStatement = true
	}

	// Check detailed problem definition if available
	if doc.Problem != nil {
		if doc.Problem.Statement != "" && !hasStatement {
			points += 3.0
			evidence = append(evidence, "Detailed problem statement present")
		}

		if doc.Problem.UserImpact != "" {
			points += 2.0
			evidence = append(evidence, "User impact documented")
		}

		if len(doc.Problem.Evidence) > 0 {
			points += 2.0
			evidence = append(evidence, fmt.Sprintf("%d evidence sources", len(doc.Problem.Evidence)))

			// Bonus for high-strength evidence
			for _, e := range doc.Problem.Evidence {
				if e.Strength == StrengthHigh {
					points += 0.5
					break
				}
			}
		}

		if doc.Problem.Confidence >= 0.7 {
			points += 1.0
			evidence = append(evidence, fmt.Sprintf("Confidence: %.0f%%", doc.Problem.Confidence*100))
		}

		if len(doc.Problem.RootCauses) > 0 {
			points += 1.0
			evidence = append(evidence, fmt.Sprintf("%d root causes identified", len(doc.Problem.RootCauses)))
		}
	} else {
		// Fallback: check executive summary completeness
		if len(doc.ExecutiveSummary.ExpectedOutcomes) > 0 {
			points += 1.0
			evidence = append(evidence, "Expected outcomes defined")
		}
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("problem_definition", score.Score)

	return score
}

func scoreUserUnderstanding(doc *Document) CategoryScore {
	score := CategoryScore{Category: "user_understanding"}
	var points float64
	var evidence []string

	// Check personas
	if len(doc.Personas) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d personas defined", len(doc.Personas)))

		// Check persona quality
		hasPainPoints := false
		hasBehaviors := false
		hasPrimary := false
		for _, persona := range doc.Personas {
			if len(persona.PainPoints) > 0 {
				hasPainPoints = true
			}
			if len(persona.Behaviors) > 0 {
				hasBehaviors = true
			}
			if persona.IsPrimary {
				hasPrimary = true
			}
		}
		if hasPainPoints {
			points += 2.0
			evidence = append(evidence, "Pain points documented")
		}
		if hasBehaviors {
			points += 1.0
			evidence = append(evidence, "Behaviors documented")
		}
		if hasPrimary {
			points += 1.0
			evidence = append(evidence, "Primary persona identified")
		}

		// Bonus for multiple personas
		if len(doc.Personas) >= 3 {
			points += 1.0
			evidence = append(evidence, "Multiple personas (3+)")
		}
	}

	// Check user stories
	if len(doc.UserStories) > 0 {
		points += 2.0
		evidence = append(evidence, fmt.Sprintf("%d user stories", len(doc.UserStories)))
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("user_understanding", score.Score)

	return score
}

func scoreMarketAwareness(doc *Document) CategoryScore {
	score := CategoryScore{Category: "market_awareness"}
	var points float64
	var evidence []string

	if doc.Market == nil {
		score.Score = 0
		score.Justification = "No market information defined"
		return score
	}

	// Check alternatives
	if len(doc.Market.Alternatives) > 0 {
		points += 4.0
		evidence = append(evidence, fmt.Sprintf("%d alternatives analyzed", len(doc.Market.Alternatives)))

		// Check for different types
		hasCompetitor := false
		hasWorkaround := false
		for _, alt := range doc.Market.Alternatives {
			if alt.Type == AlternativeCompetitor {
				hasCompetitor = true
			}
			if alt.Type == AlternativeWorkaround || alt.Type == AlternativeDoNothing {
				hasWorkaround = true
			}
		}
		if hasCompetitor && hasWorkaround {
			points += 2.0
			evidence = append(evidence, "Both competitors and alternatives covered")
		}
	}

	// Check differentiation
	if len(doc.Market.Differentiation) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d differentiation points", len(doc.Market.Differentiation)))
	}

	// Check market risks
	if len(doc.Market.MarketRisks) > 0 {
		points += 1.0
		evidence = append(evidence, "Market risks identified")
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("market_awareness", score.Score)

	return score
}

func scoreSolutionFit(doc *Document) CategoryScore {
	score := CategoryScore{Category: "solution_fit"}
	var points float64
	var evidence []string

	// Check executive summary solution
	if doc.ExecutiveSummary.ProposedSolution != "" {
		points += 2.0
		evidence = append(evidence, "Proposed solution in executive summary")
	}

	if doc.Solution == nil {
		score.Score = minFloat(points, 10.0)
		score.Evidence = strings.Join(evidence, "; ")
		score.Justification = generateJustification("solution_fit", score.Score)
		return score
	}

	// Check solution options
	if len(doc.Solution.SolutionOptions) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d solution options", len(doc.Solution.SolutionOptions)))

		if len(doc.Solution.SolutionOptions) >= 2 {
			points += 1.0
			evidence = append(evidence, "Multiple options considered")
		}
	}

	// Check selected solution
	if doc.Solution.SelectedSolutionID != "" {
		points += 2.0
		evidence = append(evidence, "Solution selected")
	}

	// Check rationale
	if doc.Solution.SolutionRationale != "" {
		points += 2.0
		evidence = append(evidence, "Selection rationale provided")
	}

	// Check problems addressed
	for _, opt := range doc.Solution.SolutionOptions {
		if len(opt.ProblemsAddressed) > 0 {
			points += 1.0
			evidence = append(evidence, "Problem mapping present")
			break
		}
	}

	// Check confidence
	if doc.Solution.Confidence >= 0.7 {
		points += 1.0
		evidence = append(evidence, fmt.Sprintf("Confidence: %.0f%%", doc.Solution.Confidence*100))
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("solution_fit", score.Score)

	return score
}

func scoreScopeDiscipline(doc *Document) CategoryScore {
	score := CategoryScore{Category: "scope_discipline"}
	var points float64
	var evidence []string

	// Check objectives (goals) from OKRs
	totalOKRs := len(doc.Objectives.OKRs)
	if totalOKRs > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d OKRs defined", totalOKRs))
	}

	// Check out of scope (non-goals)
	if len(doc.OutOfScope) > 0 {
		points += 4.0
		evidence = append(evidence, fmt.Sprintf("%d out-of-scope items defined", len(doc.OutOfScope)))
	}

	// Check key results (success criteria) from OKRs
	totalKeyResults := 0
	for _, okr := range doc.Objectives.OKRs {
		totalKeyResults += len(okr.KeyResults)
	}
	if totalKeyResults > 0 {
		points += 2.0
		evidence = append(evidence, fmt.Sprintf("%d key results defined", totalKeyResults))
	}

	// Check solution tradeoffs
	if doc.Solution != nil {
		for _, opt := range doc.Solution.SolutionOptions {
			if len(opt.Tradeoffs) > 0 {
				points += 1.0
				evidence = append(evidence, "Tradeoffs documented")
				break
			}
		}
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("scope_discipline", score.Score)

	return score
}

func scoreRequirementsQuality(doc *Document) CategoryScore {
	score := CategoryScore{Category: "requirements_quality"}
	var points float64
	var evidence []string

	// Check functional requirements
	if len(doc.Requirements.Functional) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d functional requirements", len(doc.Requirements.Functional)))

		// Check acceptance criteria
		hasAC := false
		hasTraceability := false
		hasPriority := false
		for _, req := range doc.Requirements.Functional {
			if len(req.AcceptanceCriteria) > 0 {
				hasAC = true
			}
			if len(req.UserStoryIDs) > 0 {
				hasTraceability = true
			}
			if req.Priority != "" {
				hasPriority = true
			}
		}
		if hasAC {
			points += 2.0
			evidence = append(evidence, "Acceptance criteria present")
		}
		if hasTraceability {
			points += 1.0
			evidence = append(evidence, "Traceability to user stories")
		}
		if hasPriority {
			points += 1.0
			evidence = append(evidence, "Priorities assigned")
		}
	}

	// Check NFRs
	if len(doc.Requirements.NonFunctional) > 0 {
		points += 2.0
		evidence = append(evidence, fmt.Sprintf("%d NFRs", len(doc.Requirements.NonFunctional)))

		// Check for essential NFR categories
		categories := make(map[NFRCategory]bool)
		for _, nfr := range doc.Requirements.NonFunctional {
			categories[nfr.Category] = true
		}
		essentialCount := 0
		if categories[NFRPerformance] {
			essentialCount++
		}
		if categories[NFRSecurity] {
			essentialCount++
		}
		if categories[NFRReliability] {
			essentialCount++
		}
		if essentialCount >= 2 {
			points += 1.0
			evidence = append(evidence, "Essential NFR categories covered")
		}
	}

	// Check compliance requirements
	if len(doc.Requirements.Compliance) > 0 {
		points += 1.5
		evidence = append(evidence, fmt.Sprintf("%d compliance requirements", len(doc.Requirements.Compliance)))

		// Check for compliance coverage
		complianceCategories := make(map[ComplianceCategory]bool)
		hasEvidenceReqs := false
		for _, cr := range doc.Requirements.Compliance {
			complianceCategories[cr.Category] = true
			if len(cr.EvidenceRequirements) > 0 {
				hasEvidenceReqs = true
			}
		}

		// Bonus for multiple compliance categories
		if len(complianceCategories) >= 2 {
			points += 0.5
			evidence = append(evidence, "Multiple compliance categories covered")
		}

		// Bonus for evidence requirements
		if hasEvidenceReqs {
			points += 0.5
			evidence = append(evidence, "Evidence requirements documented")
		}
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("requirements_quality", score.Score)

	return score
}

func scoreUXCoverage(doc *Document) CategoryScore {
	score := CategoryScore{Category: "ux_coverage"}
	var points float64
	var evidence []string

	if doc.UXRequirements == nil {
		score.Score = 0
		score.Justification = "No UX requirements defined"
		return score
	}

	// Check design principles
	if len(doc.UXRequirements.DesignPrinciples) > 0 {
		points += 2.0
		evidence = append(evidence, "Design principles defined")
	}

	// Check wireframes
	if len(doc.UXRequirements.Wireframes) > 0 {
		points += 2.0
		evidence = append(evidence, fmt.Sprintf("%d wireframes", len(doc.UXRequirements.Wireframes)))
	}

	// Check interaction flows
	if len(doc.UXRequirements.InteractionFlows) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d interaction flows", len(doc.UXRequirements.InteractionFlows)))
	}

	// Check accessibility
	if doc.UXRequirements.Accessibility.Standard != "" {
		points += 2.0
		evidence = append(evidence, "Accessibility requirements defined")
	}

	// Check brand guidelines
	if doc.UXRequirements.BrandGuidelines != "" {
		points += 1.0
		evidence = append(evidence, "Brand guidelines referenced")
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("ux_coverage", score.Score)

	return score
}

func scoreTechnicalFeasibility(doc *Document) CategoryScore {
	score := CategoryScore{Category: "technical_feasibility"}
	var points float64
	var evidence []string

	if doc.TechArchitecture == nil {
		score.Score = 0
		score.Justification = "No technical architecture defined"
		return score
	}

	// Check overview
	if doc.TechArchitecture.Overview != "" {
		points += 2.0
		evidence = append(evidence, "Architecture overview present")
	}

	// Check system diagram
	if doc.TechArchitecture.SystemDiagram != "" {
		points += 2.0
		evidence = append(evidence, "System diagram provided")
	}

	// Check integration points
	if len(doc.TechArchitecture.IntegrationPoints) > 0 {
		points += 2.0
		evidence = append(evidence, fmt.Sprintf("%d integration points", len(doc.TechArchitecture.IntegrationPoints)))
	}

	// Check technology stack
	if hasTechnologyStack(doc.TechArchitecture.TechnologyStack) {
		points += 2.0
		evidence = append(evidence, "Technology stack defined")
	}

	// Check security design
	if doc.TechArchitecture.SecurityDesign != "" {
		points += 1.0
		evidence = append(evidence, "Security design addressed")
	}

	// Check scalability design
	if doc.TechArchitecture.ScalabilityDesign != "" {
		points += 1.0
		evidence = append(evidence, "Scalability design addressed")
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("technical_feasibility", score.Score)

	return score
}

func scoreMetricsQuality(doc *Document) CategoryScore {
	score := CategoryScore{Category: "metrics_quality"}
	var points float64
	var evidence []string

	// Count key results from OKRs
	totalKeyResults := 0
	for _, okr := range doc.Objectives.OKRs {
		totalKeyResults += len(okr.KeyResults)
	}

	if totalKeyResults > 0 {
		points += 4.0
		evidence = append(evidence, fmt.Sprintf("%d key results defined", totalKeyResults))

		// Check key result quality
		hasTargets := false
		hasBaseline := false
		hasMeasurement := false
		for _, okr := range doc.Objectives.OKRs {
			for _, kr := range okr.KeyResults {
				if kr.Target != "" {
					hasTargets = true
				}
				if kr.Baseline != "" {
					hasBaseline = true
				}
				if kr.MeasurementMethod != "" {
					hasMeasurement = true
				}
			}
		}
		if hasTargets {
			points += 2.0
			evidence = append(evidence, "Targets defined for key results")
		}
		if hasBaseline {
			points += 2.0
			evidence = append(evidence, "Baselines documented")
		}
		if hasMeasurement {
			points += 2.0
			evidence = append(evidence, "Measurement methods specified")
		}
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("metrics_quality", score.Score)

	return score
}

func scoreRiskManagement(doc *Document) CategoryScore {
	score := CategoryScore{Category: "risk_management"}
	var points float64
	var evidence []string

	// Check assumptions
	if doc.Assumptions != nil && len(doc.Assumptions.Assumptions) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d assumptions documented", len(doc.Assumptions.Assumptions)))

		// Check for validated assumptions
		hasValidated := false
		for _, a := range doc.Assumptions.Assumptions {
			if a.Validated {
				hasValidated = true
				break
			}
		}
		if hasValidated {
			points += 1.0
			evidence = append(evidence, "Some assumptions validated")
		}
	}

	// Check risks
	if len(doc.Risks) > 0 {
		points += 3.0
		evidence = append(evidence, fmt.Sprintf("%d risks identified", len(doc.Risks)))

		// Check for mitigations
		hasMitigation := false
		for _, r := range doc.Risks {
			if r.Mitigation != "" {
				hasMitigation = true
				break
			}
		}
		if hasMitigation {
			points += 1.0
			evidence = append(evidence, "Mitigations documented")
		}
	}

	// Check constraints
	if doc.Assumptions != nil && len(doc.Assumptions.Constraints) > 0 {
		points += 2.0
		evidence = append(evidence, fmt.Sprintf("%d constraints documented", len(doc.Assumptions.Constraints)))
	}

	score.Score = minFloat(points, 10.0)
	score.Evidence = strings.Join(evidence, "; ")
	score.Justification = generateJustification("risk_management", score.Score)

	return score
}

func generateJustification(category string, score float64) string {
	categoryNames := map[string]string{
		"problem_definition":    "Problem Definition",
		"user_understanding":    "User Understanding",
		"market_awareness":      "Market Awareness",
		"solution_fit":          "Solution Fit",
		"scope_discipline":      "Scope Discipline",
		"requirements_quality":  "Requirements Quality",
		"ux_coverage":           "UX Coverage",
		"technical_feasibility": "Technical Feasibility",
		"metrics_quality":       "Metrics Quality",
		"risk_management":       "Risk Management",
	}

	name := categoryNames[category]
	if name == "" {
		name = category
	}

	switch {
	case score >= 8:
		return fmt.Sprintf("%s is strong and well-documented", name)
	case score >= 6:
		return fmt.Sprintf("%s is adequate but could be improved", name)
	case score >= 4:
		return fmt.Sprintf("%s has significant gaps that should be addressed", name)
	case score >= 2:
		return fmt.Sprintf("%s is weak and requires substantial work", name)
	default:
		return fmt.Sprintf("%s is missing or fundamentally incomplete", name)
	}
}

func categoryToOwner(category string) string {
	owners := map[string]string{
		"problem_definition":    "problem-discovery",
		"user_understanding":    "user-research",
		"market_awareness":      "market-intel",
		"solution_fit":          "solution-ideation",
		"scope_discipline":      "prd-lead",
		"requirements_quality":  "requirements",
		"ux_coverage":           "ux-journey",
		"technical_feasibility": "tech-feasibility",
		"metrics_quality":       "metrics-success",
		"risk_management":       "risk-compliance",
	}
	if owner, ok := owners[category]; ok {
		return owner
	}
	return "prd-lead"
}

func generateScoringSummary(result *ScoringResult) string {
	var parts []string

	parts = append(parts, fmt.Sprintf("Overall score: %.1f/10", result.WeightedScore))

	if len(result.Blockers) > 0 {
		parts = append(parts, fmt.Sprintf("%d blocking issues found", len(result.Blockers)))
	}

	strong := 0
	weak := 0
	for _, cs := range result.CategoryScores {
		if cs.Score >= 8 {
			strong++
		} else if cs.Score < 6 {
			weak++
		}
	}

	if strong > 0 {
		parts = append(parts, fmt.Sprintf("%d categories are strong", strong))
	}
	if weak > 0 {
		parts = append(parts, fmt.Sprintf("%d categories need improvement", weak))
	}

	switch result.Decision {
	case "approve":
		parts = append(parts, "PRD is ready for approval")
	case "revise":
		parts = append(parts, "PRD needs targeted revisions before approval")
	case "reject":
		parts = append(parts, "PRD has blocking issues that must be resolved")
	case "human_review":
		parts = append(parts, "PRD requires human review due to low overall score")
	}

	return strings.Join(parts, ". ") + "."
}

func minFloat(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// hasTechnologyStack checks if any technology stack components are defined.
func hasTechnologyStack(ts TechnologyStack) bool {
	return len(ts.Frontend) > 0 ||
		len(ts.Backend) > 0 ||
		len(ts.Database) > 0 ||
		len(ts.Infrastructure) > 0 ||
		len(ts.DevOps) > 0 ||
		len(ts.Monitoring) > 0
}
