// Package terminal provides terminal rendering for PRD evaluation reports.
package terminal

import (
	"fmt"
	"io"
	"strings"

	"github.com/plexusone/structured-evaluation/evaluation"
)

const boxWidth = 78 // Inner width between border characters

// Renderer renders evaluation reports to terminal with box formatting.
type Renderer struct {
	w io.Writer
}

// New creates a new terminal renderer.
func New(w io.Writer) *Renderer {
	return &Renderer{w: w}
}

// Render outputs the evaluation report in box format.
func (r *Renderer) Render(report *evaluation.EvaluationReport) error {
	var b strings.Builder

	// Header
	b.WriteString(header())
	b.WriteString("\n")
	b.WriteString(centerLine(strings.ToUpper(report.ReviewType) + " EVALUATION"))
	b.WriteString("\n")
	b.WriteString(separator())
	b.WriteString("\n")

	// Document info
	b.WriteString(paddedLine(fmt.Sprintf("Document: %s", truncate(report.Metadata.Document, 60))))
	b.WriteString("\n")
	if report.Metadata.DocumentTitle != "" {
		b.WriteString(paddedLine(fmt.Sprintf("Title:    %s", truncate(report.Metadata.DocumentTitle, 60))))
		b.WriteString("\n")
	}
	b.WriteString(paddedLine(fmt.Sprintf("Score:    %.1f / 10.0", report.WeightedScore)))
	b.WriteString("\n")

	// Decision with finding counts
	counts := report.Decision.FindingCounts
	decisionLine := fmt.Sprintf("Decision: %s", strings.ToUpper(string(report.Decision.Status)))
	if counts.Total > 0 {
		decisionLine += fmt.Sprintf(" (%d Critical, %d High, %d Medium)",
			counts.Critical, counts.High, counts.Medium)
	}
	b.WriteString(paddedLine(decisionLine))
	b.WriteString("\n")

	// Category scores
	b.WriteString(separator())
	b.WriteString("\n")
	b.WriteString(paddedLine("CATEGORY SCORES"))
	b.WriteString("\n")
	b.WriteString(separator())
	b.WriteString("\n")

	for _, cs := range report.Categories {
		line := formatCategoryLine(cs)
		b.WriteString(paddedLine(line))
		b.WriteString("\n")
	}

	// Findings by severity
	if len(report.Findings) > 0 {
		b.WriteString(separator())
		b.WriteString("\n")
		b.WriteString(paddedLine(fmt.Sprintf("FINDINGS (%d Critical, %d High, %d Medium)",
			counts.Critical, counts.High, counts.Medium)))
		b.WriteString("\n")
		b.WriteString(separator())
		b.WriteString("\n")

		// Group by severity
		for _, sev := range evaluation.AllSeverities() {
			for _, f := range report.Findings {
				if f.Severity == sev {
					b.WriteString(paddedLine(fmt.Sprintf("%s %-8s [%s]",
						f.Severity.Icon(), strings.ToUpper(string(f.Severity)), f.Category)))
					b.WriteString("\n")
					b.WriteString(paddedLine(fmt.Sprintf("          %s", truncate(f.Title, 60))))
					b.WriteString("\n")
					if f.Recommendation != "" {
						b.WriteString(paddedLine(fmt.Sprintf("          → %s", truncate(f.Recommendation, 58))))
						b.WriteString("\n")
					}
					b.WriteString(paddedLine(""))
					b.WriteString("\n")
				}
			}
		}
	}

	// Next steps
	if len(report.NextSteps.Immediate) > 0 || report.NextSteps.RerunCommand != "" {
		b.WriteString(separator())
		b.WriteString("\n")
		b.WriteString(paddedLine("NEXT STEPS"))
		b.WriteString("\n")
		b.WriteString(separator())
		b.WriteString("\n")

		for _, action := range report.NextSteps.Immediate {
			prefix := "🔴"
			b.WriteString(paddedLine(fmt.Sprintf("  %s %s", prefix, truncate(action.Action, 65))))
			b.WriteString("\n")
		}

		if report.NextSteps.RerunCommand != "" {
			b.WriteString(paddedLine(""))
			b.WriteString("\n")
			b.WriteString(paddedLine(fmt.Sprintf("Re-run: %s", report.NextSteps.RerunCommand)))
			b.WriteString("\n")
		}
	}

	// Final message
	b.WriteString(separator())
	b.WriteString("\n")
	b.WriteString(centerLine(finalMessage(report)))
	b.WriteString("\n")
	b.WriteString(footer())
	b.WriteString("\n")

	_, err := fmt.Fprint(r.w, b.String())
	return err
}

func formatCategoryLine(cs evaluation.CategoryScore) string {
	name := categoryDisplayName(cs.Category)
	if len(name) > 24 {
		name = name[:21] + "..."
	}

	icon := cs.Status.Icon()
	statusText := string(cs.Status)

	justification := truncate(cs.Justification, 28)

	return fmt.Sprintf("  %-24s %s %-4s %4.1f/%.0f  %s",
		name, icon, statusText, cs.Score, cs.MaxScore, justification)
}

func categoryDisplayName(category string) string {
	names := map[string]string{
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
	if name, ok := names[category]; ok {
		return name
	}
	// Handle custom sections
	if strings.HasPrefix(category, "custom:") {
		return strings.TrimPrefix(category, "custom:")
	}
	return category
}

func finalMessage(report *evaluation.EvaluationReport) string {
	switch report.Decision.Status {
	case evaluation.DecisionPass:
		return fmt.Sprintf("✅ %s PASSED (%.1f/10)", strings.ToUpper(report.ReviewType), report.WeightedScore)
	case evaluation.DecisionConditional:
		return fmt.Sprintf("⚠️ %s CONDITIONAL (%.1f/10)", strings.ToUpper(report.ReviewType), report.WeightedScore)
	case evaluation.DecisionFail:
		return fmt.Sprintf("❌ %s BLOCKED - %d issues to resolve",
			strings.ToUpper(report.ReviewType), report.Decision.FindingCounts.BlockingCount())
	case evaluation.DecisionHumanReview:
		return fmt.Sprintf("👤 %s NEEDS HUMAN REVIEW (%.1f/10)", strings.ToUpper(report.ReviewType), report.WeightedScore)
	default:
		return fmt.Sprintf("📋 %s: %.1f/10", strings.ToUpper(report.ReviewType), report.WeightedScore)
	}
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// Box drawing functions
func header() string {
	return "╔" + strings.Repeat("═", boxWidth) + "╗"
}

func separator() string {
	return "╠" + strings.Repeat("═", boxWidth) + "╣"
}

func footer() string {
	return "╚" + strings.Repeat("═", boxWidth) + "╝"
}

func centerLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen)
	left := padding / 2
	right := padding - left
	return "║" + strings.Repeat(" ", left) + text + strings.Repeat(" ", right) + "║"
}

func paddedLine(text string) string {
	visualLen := visualLength(text)
	padding := max(0, boxWidth-visualLen-1)
	return "║ " + text + strings.Repeat(" ", padding) + "║"
}

// visualLength accounts for emoji double-width display.
func visualLength(s string) int {
	length := 0
	for _, r := range s {
		if r >= 0x1F300 && r <= 0x1FAFF {
			length += 2 // Emoji range
		} else if r >= 0x2600 && r <= 0x27BF {
			length += 2 // Symbol range
		} else {
			length++
		}
	}
	return length
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
