package prd

import (
	"time"

	"github.com/plexusone/structured-evaluation/evaluation"
)

// EvaluationCategory defines metadata for a PRD evaluation category.
type EvaluationCategory struct {
	ID          string  // Category identifier
	Name        string  // Human-readable name
	Description string  // What this category evaluates
	Weight      float64 // Weight in overall score (0.0-1.0)
	Owner       string  // Suggested agent/team owner
}

// StandardCategories returns the standard PRD evaluation categories.
// These match the sections defined in the PRD schema.
func StandardCategories() []EvaluationCategory {
	return []EvaluationCategory{
		{
			ID:          "problem_definition",
			Name:        "Problem Definition",
			Description: "Clarity of problem statement, supporting evidence, root cause analysis, and user impact",
			Weight:      0.20,
			Owner:       "problem-discovery",
		},
		{
			ID:          "solution_fit",
			Name:        "Solution Fit",
			Description: "Quality of solution options, selection rationale, and alignment with problem",
			Weight:      0.15,
			Owner:       "solution-ideation",
		},
		{
			ID:          "user_understanding",
			Name:        "User Understanding",
			Description: "Depth of personas, user stories, pain points, and behavioral insights",
			Weight:      0.10,
			Owner:       "user-research",
		},
		{
			ID:          "market_awareness",
			Name:        "Market Awareness",
			Description: "Competitive analysis, alternatives assessment, and differentiation",
			Weight:      0.10,
			Owner:       "market-intel",
		},
		{
			ID:          "scope_discipline",
			Name:        "Scope Discipline",
			Description: "Clear objectives, out-of-scope items, and success criteria",
			Weight:      0.10,
			Owner:       "prd-lead",
		},
		{
			ID:          "requirements_quality",
			Name:        "Requirements Quality",
			Description: "Functional and non-functional requirements with acceptance criteria",
			Weight:      0.10,
			Owner:       "requirements",
		},
		{
			ID:          "metrics_quality",
			Name:        "Metrics Quality",
			Description: "Success metrics with targets, baselines, and measurement methods",
			Weight:      0.10,
			Owner:       "metrics-success",
		},
		{
			ID:          "ux_coverage",
			Name:        "UX Coverage",
			Description: "Design principles, wireframes, interaction flows, and accessibility",
			Weight:      0.05,
			Owner:       "ux-journey",
		},
		{
			ID:          "technical_feasibility",
			Name:        "Technical Feasibility",
			Description: "Architecture overview, integrations, technology stack, and security design",
			Weight:      0.05,
			Owner:       "tech-feasibility",
		},
		{
			ID:          "risk_management",
			Name:        "Risk Management",
			Description: "Assumptions, constraints, risks, and mitigations",
			Weight:      0.05,
			Owner:       "risk-compliance",
		},
	}
}

// GenerateEvaluationTemplate creates an EvaluationReport template from a PRD document.
// The template includes all standard categories plus custom sections.
// Scores are initialized as partial - they will be filled in by the LLM judge.
func GenerateEvaluationTemplate(doc *Document, filename string) *evaluation.EvaluationReport {
	report := evaluation.NewEvaluationReport("prd", filename)

	// Set metadata from document
	if doc.Metadata.ID != "" {
		report.Metadata.DocumentID = doc.Metadata.ID
	}
	if doc.Metadata.Title != "" {
		report.Metadata.DocumentTitle = doc.Metadata.Title
	}
	if doc.Metadata.Version != "" {
		report.Metadata.DocumentVersion = doc.Metadata.Version
	}
	report.Metadata.GeneratedAt = time.Now().UTC()
	report.Metadata.GeneratedBy = "structured-requirements"

	// Add standard categories (to be scored by LLM judge)
	for _, cat := range StandardCategories() {
		cr := evaluation.NewCategoryResult(
			cat.ID,
			evaluation.ScorePartial, // Placeholder, to be filled by LLM judge
			"",                      // Reasoning to be filled by LLM judge
		)
		report.AddCategoryResult(*cr)
	}

	// Add custom sections as categories
	for _, section := range doc.CustomSections {
		cr := evaluation.NewCategoryResult(
			"custom:"+section.ID,
			evaluation.ScorePartial, // Placeholder, to be filled by LLM judge
			"",                      // Reasoning to be filled by LLM judge
		)
		report.AddCategoryResult(*cr)
	}

	// Set default pass criteria
	report.PassCriteria = evaluation.DefaultPassCriteria()

	return report
}

// GenerateEvaluationTemplateWithWeights creates a template with custom category weights.
// Note: Weights are stored in the EvaluationCategory metadata, not in CategoryResult.
// The new evaluation model uses categorical scores (pass/partial/fail) rather than
// weighted numeric scores.
func GenerateEvaluationTemplateWithWeights(doc *Document, filename string, _ map[string]float64) *evaluation.EvaluationReport {
	// Weights are no longer stored in CategoryResult in the new API.
	// The evaluation model now uses categorical scoring (pass/partial/fail).
	// Return the standard template.
	return GenerateEvaluationTemplate(doc, filename)
}

// CategoryDescriptions returns a map of category IDs to descriptions.
// Useful for providing context to LLM judges.
func CategoryDescriptions() map[string]string {
	descs := make(map[string]string)
	for _, cat := range StandardCategories() {
		descs[cat.ID] = cat.Description
	}
	return descs
}

// CategoryOwners returns a map of category IDs to suggested owners.
// Useful for assigning findings to responsible teams.
func CategoryOwners() map[string]string {
	owners := make(map[string]string)
	for _, cat := range StandardCategories() {
		owners[cat.ID] = cat.Owner
	}
	return owners
}

// GetCategoriesFromDocument extracts the list of categories that should be evaluated
// based on what's present in the document. This includes standard categories and
// any custom sections defined in the PRD.
func GetCategoriesFromDocument(doc *Document) []EvaluationCategory {
	categories := StandardCategories()

	// Add custom sections
	for _, section := range doc.CustomSections {
		categories = append(categories, EvaluationCategory{
			ID:          "custom:" + section.ID,
			Name:        section.Title,
			Description: section.Description,
			Weight:      0.05,
			Owner:       "prd-lead",
		})
	}

	return categories
}
