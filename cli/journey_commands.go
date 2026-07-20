// Package cli provides the exported Cobra command tree for the PRISM roadmap CLI.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/grokify/prism-roadmap/journey"
)

// ============================================================================
// Journey Parent Command
// ============================================================================

var journeyCmd = &cobra.Command{
	Use:   "journey",
	Short: "Work with journey roadmaps",
	Long: `Commands for generating and validating journey roadmaps.

Journey roadmaps model capability evolution over time, tracking:
  - Capability maturity progression
  - Outcome journeys with business impact
  - Initiatives driving improvements
  - Dependencies and risks
  - Team capacity and assignments`,
}

// ============================================================================
// Journey Validate Command
// ============================================================================

var journeyValidateCmd = &cobra.Command{
	Use:   "validate FILE",
	Short: "Validate a journey roadmap JSON file",
	Long: `Validate a journey roadmap JSON file against structural rules.

Checks:
  - Required fields (id, name)
  - Period references are valid
  - Capability journey targets reference existing periods
  - Dependency references are valid
  - Team hierarchy is consistent

Examples:
  splan journey validate roadmap.json`,
	Args: cobra.ExactArgs(1),
	RunE: runJourneyValidate,
}

func runJourneyValidate(cmd *cobra.Command, args []string) error {
	roadmap, err := loadJourneyRoadmap(args[0])
	if err != nil {
		return err
	}

	issues := validateJourneyRoadmap(roadmap)
	if len(issues) > 0 {
		fmt.Println("Validation issues found:")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
		return fmt.Errorf("validation failed with %d issues", len(issues))
	}

	fmt.Printf("✓ %s is valid\n", args[0])
	return nil
}

func loadJourneyRoadmap(filename string) (*journey.JourneyRoadmap, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading file %s: %w", filename, err)
	}

	var roadmap journey.JourneyRoadmap
	if err := json.Unmarshal(data, &roadmap); err != nil {
		return nil, fmt.Errorf("parsing JSON from %s: %w", filename, err)
	}

	return &roadmap, nil
}

func validateJourneyRoadmap(r *journey.JourneyRoadmap) []string {
	var issues []string

	// Required fields
	if r.ID == "" {
		issues = append(issues, "missing required field: id")
	}
	if r.Name == "" {
		issues = append(issues, "missing required field: name")
	}

	// Build period lookup
	periodIDs := make(map[string]bool)
	if r.TimeModel != nil {
		for _, p := range r.TimeModel.Periods {
			if p.ID == "" {
				issues = append(issues, "period missing id")
			} else {
				periodIDs[p.ID] = true
			}
		}
	}

	// Validate capability journey targets
	for _, cj := range r.CapabilityJourneys {
		if cj.ID == "" {
			issues = append(issues, "capability journey missing id")
		}
		for _, t := range cj.TargetStates {
			if t.PeriodID != "" && !periodIDs[t.PeriodID] {
				issues = append(issues, fmt.Sprintf("capability %s: target references unknown period %s", cj.ID, t.PeriodID))
			}
		}
	}

	// Validate dependencies
	entityIDs := buildEntityIDSet(r)
	for _, d := range r.Dependencies {
		if d.From.ID != "" && !entityIDs[d.From.ID] {
			issues = append(issues, fmt.Sprintf("dependency references unknown entity: %s", d.From.ID))
		}
		if d.To.ID != "" && !entityIDs[d.To.ID] {
			issues = append(issues, fmt.Sprintf("dependency references unknown entity: %s", d.To.ID))
		}
	}

	// Validate team hierarchy
	teamIDs := make(map[string]bool)
	for _, t := range r.Teams {
		if t.ID == "" {
			issues = append(issues, "team missing id")
		} else {
			teamIDs[t.ID] = true
		}
	}
	for _, t := range r.Teams {
		if t.ParentID != "" && !teamIDs[t.ParentID] {
			issues = append(issues, fmt.Sprintf("team %s: references unknown parent %s", t.ID, t.ParentID))
		}
	}

	return issues
}

func buildEntityIDSet(r *journey.JourneyRoadmap) map[string]bool {
	ids := make(map[string]bool)
	for _, cj := range r.CapabilityJourneys {
		ids[cj.ID] = true
	}
	for _, oj := range r.OutcomeJourneys {
		ids[oj.ID] = true
	}
	for _, i := range r.Initiatives {
		ids[i.ID] = true
	}
	for _, t := range r.Teams {
		ids[t.ID] = true
	}
	return ids
}

// ============================================================================
// Journey Generate Command
// ============================================================================

var journeyGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output from journey roadmap",
	Long:  `Generate various output formats from a journey roadmap JSON file.`,
}

var journeyGenerateMarkdownFlags struct {
	output   string
	sections string
}

var journeyGenerateMarkdownCmd = &cobra.Command{
	Use:   "markdown FILE",
	Short: "Generate markdown document",
	Long: `Generate a markdown document from a journey roadmap JSON file.

Sections (comma-separated):
  all        - All sections (default)
  summary    - Executive summary
  timeline   - Capability timeline table
  storyboard - Period-by-period narrative
  deps       - Dependencies analysis
  teams      - Team assignments
  risks      - Risk register

Examples:
  splan journey generate markdown roadmap.json
  splan journey generate markdown roadmap.json -o roadmap.md
  splan journey generate markdown roadmap.json --sections=summary,timeline`,
	Args: cobra.ExactArgs(1),
	RunE: runJourneyGenerateMarkdown,
}

func init() {
	journeyGenerateMarkdownCmd.Flags().StringVarP(&journeyGenerateMarkdownFlags.output, "output", "o", "", "Output file (default: stdout)")
	journeyGenerateMarkdownCmd.Flags().StringVar(&journeyGenerateMarkdownFlags.sections, "sections", "all", "Sections to include (comma-separated)")
}

func runJourneyGenerateMarkdown(cmd *cobra.Command, args []string) error {
	roadmap, err := loadJourneyRoadmap(args[0])
	if err != nil {
		return err
	}

	sections := strings.Split(journeyGenerateMarkdownFlags.sections, ",")
	sectionSet := make(map[string]bool)
	for _, s := range sections {
		sectionSet[strings.TrimSpace(s)] = true
	}
	includeAll := sectionSet["all"]

	var sb strings.Builder

	// Title and vision
	sb.WriteString(fmt.Sprintf("# %s\n\n", roadmap.Name))
	if roadmap.Vision != "" {
		sb.WriteString(fmt.Sprintf("*%s*\n\n", roadmap.Vision))
	}
	if roadmap.Description != "" {
		sb.WriteString(fmt.Sprintf("%s\n\n", roadmap.Description))
	}

	// Executive summary
	if includeAll || sectionSet["summary"] {
		summary := journey.GenerateExecutiveSummary(roadmap)
		if summary != nil {
			sb.WriteString("## Executive Summary\n\n")
			if summary.CurrentStateOverview != "" {
				sb.WriteString(fmt.Sprintf("%s\n\n", summary.CurrentStateOverview))
			}
			if len(summary.KeyTransformations) > 0 {
				sb.WriteString("### Key Transformations\n\n")
				for _, t := range summary.KeyTransformations {
					sb.WriteString(fmt.Sprintf("- %s\n", t))
				}
				sb.WriteString("\n")
			}
			if len(summary.ExpectedOutcomes) > 0 {
				sb.WriteString("### Expected Outcomes\n\n")
				for _, o := range summary.ExpectedOutcomes {
					sb.WriteString(fmt.Sprintf("- %s\n", o))
				}
				sb.WriteString("\n")
			}
			if summary.Timeline != "" {
				sb.WriteString(fmt.Sprintf("**Timeline:** %s\n\n", summary.Timeline))
			}
		}
	}

	// Timeline table
	if includeAll || sectionSet["timeline"] {
		if len(roadmap.CapabilityJourneys) > 0 && roadmap.TimeModel != nil && len(roadmap.TimeModel.Periods) > 0 {
			sb.WriteString("## Capability Timeline\n\n")
			sb.WriteString(renderTimelineTable(roadmap))
			sb.WriteString("\n")
		}
	}

	// Storyboard
	if includeAll || sectionSet["storyboard"] {
		cards := journey.BuildStoryboard(roadmap)
		if len(cards) > 0 {
			sb.WriteString("## Journey Narrative\n\n")
			for _, card := range cards {
				sb.WriteString(fmt.Sprintf("### %s\n\n", card.PeriodLabel))
				if card.Headline != "" {
					sb.WriteString(fmt.Sprintf("**%s**\n\n", card.Headline))
				}
				if card.UserImpact != "" {
					sb.WriteString(fmt.Sprintf("%s\n\n", card.UserImpact))
				}
				if len(card.MaturityChanges) > 0 {
					sb.WriteString("**Maturity Changes:**\n\n")
					for _, mc := range card.MaturityChanges {
						sb.WriteString(fmt.Sprintf("- %s: %s → %s\n", mc.CapabilityName, mc.From, mc.To))
					}
					sb.WriteString("\n")
				}
			}
		}
	}

	// Dependencies
	if includeAll || sectionSet["deps"] {
		if len(roadmap.Dependencies) > 0 {
			sb.WriteString("## Dependencies\n\n")
			sb.WriteString("| From | To | Type | Status | Risk |\n")
			sb.WriteString("|------|-----|------|--------|------|\n")
			for _, d := range roadmap.Dependencies {
				fromName := d.From.Name
				if fromName == "" {
					fromName = d.From.ID
				}
				toName := d.To.Name
				if toName == "" {
					toName = d.To.ID
				}
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
					fromName, toName, d.Type, d.Status, d.Risk))
			}
			sb.WriteString("\n")
		}
	}

	// Teams
	if includeAll || sectionSet["teams"] {
		if len(roadmap.Teams) > 0 {
			sb.WriteString("## Teams\n\n")
			for _, t := range roadmap.Teams {
				level := ""
				if t.Level != "" {
					level = fmt.Sprintf(" (%s)", t.Level)
				}
				sb.WriteString(fmt.Sprintf("### %s%s\n\n", t.Name, level))
				if t.Description != "" {
					sb.WriteString(fmt.Sprintf("%s\n\n", t.Description))
				}
				if t.LeaderName != "" {
					sb.WriteString(fmt.Sprintf("**Lead:** %s\n\n", t.LeaderName))
				}
			}
		}
	}

	// Risks
	if includeAll || sectionSet["risks"] {
		if len(roadmap.Risks) > 0 {
			sb.WriteString("## Risk Register\n\n")
			sb.WriteString("| ID | Description | Probability | Impact | Status |\n")
			sb.WriteString("|----|-------------|-------------|--------|--------|\n")
			for _, r := range roadmap.Risks {
				sb.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
					r.ID, truncate(r.Description, 50), r.Probability, r.Impact, r.Status))
			}
			sb.WriteString("\n")
		}
	}

	output := sb.String()

	if journeyGenerateMarkdownFlags.output != "" {
		if err := os.WriteFile(journeyGenerateMarkdownFlags.output, []byte(output), 0600); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Printf("Generated %s\n", journeyGenerateMarkdownFlags.output)
	} else {
		fmt.Print(output)
	}

	return nil
}

func renderTimelineTable(r *journey.JourneyRoadmap) string {
	if r.TimeModel == nil || len(r.TimeModel.Periods) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header
	sb.WriteString("| Capability |")
	for _, p := range r.TimeModel.Periods {
		sb.WriteString(fmt.Sprintf(" %s |", p.Label))
	}
	sb.WriteString("\n")

	// Separator
	sb.WriteString("|------------|")
	for range r.TimeModel.Periods {
		sb.WriteString("------|")
	}
	sb.WriteString("\n")

	// Rows
	for _, cj := range r.CapabilityJourneys {
		sb.WriteString(fmt.Sprintf("| %s |", cj.Name))

		// Build target lookup
		targetByPeriod := make(map[string]*journey.TargetState)
		for i := range cj.TargetStates {
			targetByPeriod[cj.TargetStates[i].PeriodID] = &cj.TargetStates[i]
		}

		for _, p := range r.TimeModel.Periods {
			if t, ok := targetByPeriod[p.ID]; ok {
				confidence := ""
				if t.Confidence > 0 {
					confidence = fmt.Sprintf(" (%.0f%%)", t.Confidence*100)
				}
				sb.WriteString(fmt.Sprintf(" %s%s |", t.MaturityLevel, confidence))
			} else {
				sb.WriteString(" - |")
			}
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

// ============================================================================
// Journey Check Command
// ============================================================================

var journeyCheckFlags struct {
	jsonOutput bool
}

var journeyCheckCmd = &cobra.Command{
	Use:   "check FILE",
	Short: "Check journey roadmap completeness",
	Long: `Check a journey roadmap for completeness and quality.

Evaluates:
  - Required sections presence
  - Capability coverage (all periods have targets)
  - Initiative coverage (capabilities have initiatives)
  - Dependency completeness
  - Risk coverage
  - Narrative quality

Examples:
  splan journey check roadmap.json
  splan journey check roadmap.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runJourneyCheck,
}

func init() {
	journeyCheckCmd.Flags().BoolVar(&journeyCheckFlags.jsonOutput, "json", false, "Output as JSON")
}

type JourneyCheckResult struct {
	Filename           string         `json:"filename"`
	Score              float64        `json:"score"`
	Grade              string         `json:"grade"`
	Sections           map[string]int `json:"sections"`
	CapabilityCoverage float64        `json:"capabilityCoverage"`
	InitiativeCoverage float64        `json:"initiativeCoverage"`
	Recommendations    []string       `json:"recommendations"`
}

func runJourneyCheck(cmd *cobra.Command, args []string) error {
	roadmap, err := loadJourneyRoadmap(args[0])
	if err != nil {
		return err
	}

	result := checkJourneyCompleteness(args[0], roadmap)

	if journeyCheckFlags.jsonOutput {
		data, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling result: %w", err)
		}
		fmt.Println(string(data))
	} else {
		printJourneyCheckResult(result)
	}

	return nil
}

func checkJourneyCompleteness(filename string, r *journey.JourneyRoadmap) JourneyCheckResult {
	result := JourneyCheckResult{
		Filename:        filename,
		Sections:        make(map[string]int),
		Recommendations: []string{},
	}

	// Count sections
	result.Sections["capabilityJourneys"] = len(r.CapabilityJourneys)
	result.Sections["outcomeJourneys"] = len(r.OutcomeJourneys)
	result.Sections["initiatives"] = len(r.Initiatives)
	result.Sections["dependencies"] = len(r.Dependencies)
	result.Sections["teams"] = len(r.Teams)
	result.Sections["risks"] = len(r.Risks)
	if r.TimeModel != nil {
		result.Sections["periods"] = len(r.TimeModel.Periods)
	}
	if r.Narrative != nil {
		result.Sections["narrative"] = 1
	}

	// Calculate capability coverage
	if len(r.CapabilityJourneys) > 0 && r.TimeModel != nil && len(r.TimeModel.Periods) > 0 {
		totalTargets := 0
		expectedTargets := len(r.CapabilityJourneys) * len(r.TimeModel.Periods)
		for _, cj := range r.CapabilityJourneys {
			totalTargets += len(cj.TargetStates)
		}
		if expectedTargets > 0 {
			result.CapabilityCoverage = float64(totalTargets) / float64(expectedTargets) * 100
		}
	}

	// Calculate initiative coverage
	if len(r.CapabilityJourneys) > 0 {
		capabilitiesWithInitiatives := make(map[string]bool)
		for _, init := range r.Initiatives {
			for _, adv := range init.Advances {
				capabilitiesWithInitiatives[adv.CapabilityID] = true
			}
		}
		result.InitiativeCoverage = float64(len(capabilitiesWithInitiatives)) / float64(len(r.CapabilityJourneys)) * 100
	}

	// Calculate overall score
	score := 0.0
	maxScore := 0.0

	// Required: name, id, periods, capabilities
	maxScore += 40
	if r.Name != "" {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Add roadmap name")
	}
	if r.ID != "" {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Add roadmap id")
	}
	if r.TimeModel != nil && len(r.TimeModel.Periods) > 0 {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Define time periods")
	}
	if len(r.CapabilityJourneys) > 0 {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Add capability journeys")
	}

	// Optional: vision, initiatives, teams, risks, narrative
	maxScore += 60
	if r.Vision != "" {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Add vision statement")
	}
	if len(r.Initiatives) > 0 {
		score += 15
	} else {
		result.Recommendations = append(result.Recommendations, "Add initiatives driving capability improvements")
	}
	if len(r.Teams) > 0 {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Add team assignments")
	}
	if len(r.Risks) > 0 {
		score += 10
	} else {
		result.Recommendations = append(result.Recommendations, "Add risk register")
	}
	if r.Narrative != nil {
		score += 15
	} else {
		result.Recommendations = append(result.Recommendations, "Add narrative for executive communication")
	}

	result.Score = score / maxScore * 100
	result.Grade = scoreToGrade(result.Score)

	return result
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

func printJourneyCheckResult(r JourneyCheckResult) {
	fmt.Println("=============================================================")
	fmt.Println("JOURNEY ROADMAP COMPLETENESS REPORT")
	fmt.Println("=============================================================")
	fmt.Println()
	fmt.Printf("Overall Score: %.1f%% (Grade: %s)\n", r.Score, r.Grade)
	fmt.Println()
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("SECTIONS")
	fmt.Println("-------------------------------------------------------------")
	fmt.Println()
	for section, count := range r.Sections {
		status := "✓"
		if count == 0 {
			status = "✗"
		}
		fmt.Printf("  [%s] %-25s %d\n", status, section, count)
	}
	fmt.Println()
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("COVERAGE")
	fmt.Println("-------------------------------------------------------------")
	fmt.Println()
	fmt.Printf("  Capability targets: %.1f%%\n", r.CapabilityCoverage)
	fmt.Printf("  Initiative coverage: %.1f%%\n", r.InitiativeCoverage)
	fmt.Println()

	if len(r.Recommendations) > 0 {
		fmt.Println("-------------------------------------------------------------")
		fmt.Println("RECOMMENDATIONS")
		fmt.Println("-------------------------------------------------------------")
		fmt.Println()
		for _, rec := range r.Recommendations {
			fmt.Printf("  - %s\n", rec)
		}
		fmt.Println()
	}

	fmt.Println("=============================================================")
}

// ============================================================================
// Journey Init Command
// ============================================================================

var journeyInitFlags struct {
	output string
}

var journeyInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a journey roadmap template",
	Long: `Create a new journey roadmap JSON template file.

Examples:
  splan journey init
  splan journey init -o my-roadmap.json`,
	RunE: runJourneyInit,
}

func init() {
	journeyInitCmd.Flags().StringVarP(&journeyInitFlags.output, "output", "o", "journey-roadmap.json", "Output file")
}

var journeyTemplate = `{
  "id": "{{ .ID }}",
  "name": "{{ .Name }}",
  "vision": "Where we want to be",
  "description": "Overview of the roadmap",
  "timeModel": {
    "type": "quarterly",
    "fiscalYear": "FY2026",
    "periods": [
      { "id": "now", "label": "Current State", "isCurrent": true },
      { "id": "2026-q1", "label": "Q1 2026", "startDate": "2026-01-01", "endDate": "2026-03-31" },
      { "id": "2026-q2", "label": "Q2 2026", "startDate": "2026-04-01", "endDate": "2026-06-30" },
      { "id": "2026-q3", "label": "Q3 2026", "startDate": "2026-07-01", "endDate": "2026-09-30" },
      { "id": "2026-q4", "label": "Q4 2026", "startDate": "2026-10-01", "endDate": "2026-12-31" }
    ]
  },
  "capabilityJourneys": [
    {
      "id": "cap-1",
      "name": "Example Capability",
      "description": "Description of the capability",
      "currentLevel": 2,
      "targets": [
        { "periodId": "2026-q1", "targetLevel": 3, "confidence": "high" },
        { "periodId": "2026-q2", "targetLevel": 4, "confidence": "medium" }
      ]
    }
  ],
  "initiatives": [
    {
      "id": "init-1",
      "name": "Example Initiative",
      "description": "Work to improve capability",
      "capabilityIds": ["cap-1"],
      "startPeriod": "2026-q1",
      "endPeriod": "2026-q2"
    }
  ],
  "teams": [
    {
      "id": "team-1",
      "name": "Platform Team",
      "level": "squad"
    }
  ],
  "risks": [
    {
      "id": "risk-1",
      "description": "Example risk",
      "probability": "medium",
      "impact": "high",
      "status": "identified"
    }
  ]
}
`

func runJourneyInit(cmd *cobra.Command, args []string) error {
	tmpl, err := template.New("journey").Parse(journeyTemplate)
	if err != nil {
		return fmt.Errorf("parsing template: %w", err)
	}

	data := struct {
		ID   string
		Name string
	}{
		ID:   "roadmap-2026",
		Name: "Capability Roadmap 2026",
	}

	file, err := os.Create(journeyInitFlags.output)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("writing template: %w", err)
	}

	fmt.Printf("Created %s\n", journeyInitFlags.output)
	return nil
}
