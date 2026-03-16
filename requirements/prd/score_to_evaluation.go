package prd

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/plexusone/structured-evaluation/evaluation"
)

// ScoreToEvaluationReport converts deterministic scoring results to an EvaluationReport.
// This allows the existing deterministic scoring to output in the standardized format
// that can be combined with LLM-based evaluations.
func ScoreToEvaluationReport(doc *Document, filename string) *evaluation.EvaluationReport {
	// Get the deterministic scoring result
	result := Score(doc)

	report := evaluation.NewEvaluationReport("prd", filepath.Base(filename))

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
	report.Metadata.GeneratedBy = "srequirements (deterministic)"

	// Convert category scores
	for _, cs := range result.CategoryScores {
		evalCategory := evaluation.CategoryScore{
			Category:      cs.Category,
			Score:         cs.Score,
			MaxScore:      cs.MaxScore,
			Weight:        cs.Weight,
			Justification: cs.Justification,
			Evidence:      cs.Evidence,
		}
		evalCategory.ComputeStatus()
		report.Categories = append(report.Categories, evalCategory)
	}

	// Convert revision triggers to findings
	for _, rt := range result.RevisionTriggers {
		severity := severityFromString(rt.Severity)
		finding := evaluation.Finding{
			ID:             rt.IssueID,
			Category:       rt.Category,
			Severity:       severity,
			Title:          rt.Description,
			Description:    rt.Description,
			Recommendation: generateFixRecommendationForCategory(rt.Category),
			Owner:          rt.RecommendedOwner,
			Effort:         estimateEffortForCategory(rt.Category),
		}
		report.Findings = append(report.Findings, finding)
	}

	// Set weighted score
	report.WeightedScore = result.WeightedScore

	// Finalize the report (computes decision, next steps, summary)
	rerunCommand := fmt.Sprintf("srequirements prd score %s", filename)
	report.Finalize(rerunCommand)

	return report
}

func severityFromString(s string) evaluation.Severity {
	switch s {
	case "blocker":
		return evaluation.SeverityCritical
	case "major":
		return evaluation.SeverityHigh
	case "minor":
		return evaluation.SeverityMedium
	default:
		return evaluation.SeverityLow
	}
}

func generateFixRecommendationForCategory(category string) string {
	recommendations := map[string]string{
		"problem_definition":    "Add detailed problem statement with evidence and root cause analysis",
		"user_understanding":    "Define at least 3 personas with pain points and behaviors; add user stories",
		"market_awareness":      "Add competitive analysis with 3-5 alternatives and differentiation points",
		"solution_fit":          "Document solution options with pros/cons and selection rationale",
		"scope_discipline":      "Define clear objectives and out-of-scope items",
		"requirements_quality":  "Add functional requirements with acceptance criteria and NFRs",
		"ux_coverage":           "Add UX requirements including wireframes, flows, and accessibility",
		"technical_feasibility": "Document technical architecture with integration points and tech stack",
		"metrics_quality":       "Define success metrics with targets, baselines, and measurement methods",
		"risk_management":       "Identify risks with mitigations; document assumptions and constraints",
	}
	if rec, ok := recommendations[category]; ok {
		return rec
	}
	return "Review and improve this category"
}

func estimateEffortForCategory(category string) string {
	highEffort := map[string]bool{
		"technical_feasibility": true,
		"requirements_quality":  true,
		"market_awareness":      true,
	}
	if highEffort[category] {
		return "high"
	}
	return "medium"
}
