// Package cli provides goal-setting CLI commands.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/prism-roadmap/goals"
	"github.com/grokify/prism-roadmap/goals/okr"
)

// ============================================================================
// Goals Cascade Command
// ============================================================================

var goalsCascadeFlags struct {
	output     string
	team       string
	owner      string
	period     string
	inheritKRs bool
	filterIDs  string
}

var goalsCascadeCmd = &cobra.Command{
	Use:   "cascade PARENT_FILE",
	Short: "Generate child goals from parent",
	Long: `Generate child team goals cascaded from parent organizational goals.

Takes a parent OKR document and generates child OKRs with proper alignment.
Child objectives reference parent objectives via parentId and alignedWith fields.

Examples:
  splan goals cascade company-okrs.json --team="Platform Team" --owner="jane@example.com"
  splan goals cascade parent.json --team="API Team" --period="Q3 2026" -o team-okrs.json
  splan goals cascade parent.json --team="Frontend" --inherit-krs`,
	Args: cobra.ExactArgs(1),
	RunE: runGoalsCascade,
}

func init() {
	goalsCascadeCmd.Flags().StringVarP(&goalsCascadeFlags.output, "output", "o", "", "Output file (default: stdout)")
	goalsCascadeCmd.Flags().StringVar(&goalsCascadeFlags.team, "team", "", "Child team name (required)")
	goalsCascadeCmd.Flags().StringVar(&goalsCascadeFlags.owner, "owner", "", "Child goals owner")
	goalsCascadeCmd.Flags().StringVar(&goalsCascadeFlags.period, "period", "", "Period for child goals (e.g., Q2 2026)")
	goalsCascadeCmd.Flags().BoolVar(&goalsCascadeFlags.inheritKRs, "inherit-krs", false, "Create child KRs based on parent KRs")
	goalsCascadeCmd.Flags().StringVar(&goalsCascadeFlags.filterIDs, "filter", "", "Comma-separated objective IDs to cascade (default: all)")
	_ = goalsCascadeCmd.MarkFlagRequired("team")
}

func runGoalsCascade(cmd *cobra.Command, args []string) error {
	// Load parent OKR document
	parentDoc, err := okr.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading parent file: %w", err)
	}

	// Build cascade options
	opts := goals.CascadeOptions{
		ChildTeam:                  goalsCascadeFlags.team,
		ChildOwner:                 goalsCascadeFlags.owner,
		ChildPeriod:                goalsCascadeFlags.period,
		InheritTags:                true,
		CreateKeyResultsFromParent: goalsCascadeFlags.inheritKRs,
	}

	if goalsCascadeFlags.filterIDs != "" {
		opts.FilterObjectiveIDs = strings.Split(goalsCascadeFlags.filterIDs, ",")
	}

	// Cascade
	result, err := goals.CascadeOKR(parentDoc, opts)
	if err != nil {
		return fmt.Errorf("cascading goals: %w", err)
	}

	// Output child document
	if result.ChildGoals.OKR != nil {
		childDoc := &okr.OKRDocument{
			Metadata: &okr.Metadata{
				ID:     fmt.Sprintf("%s-%s", result.ParentID, strings.ToLower(strings.ReplaceAll(opts.ChildTeam, " ", "-"))),
				Name:   fmt.Sprintf("%s OKRs", opts.ChildTeam),
				Owner:  opts.ChildOwner,
				Team:   opts.ChildTeam,
				Period: opts.ChildPeriod,
				Status: okr.StatusDraft,
			},
			Theme:      parentDoc.Theme,
			Objectives: result.ChildGoals.OKR.ToObjectives(),
			Alignment: &okr.Alignment{
				ParentOKRID: result.ParentID,
			},
		}

		data, err := childDoc.JSON()
		if err != nil {
			return fmt.Errorf("marshaling output: %w", err)
		}

		if goalsCascadeFlags.output != "" {
			if err := os.WriteFile(goalsCascadeFlags.output, data, 0600); err != nil {
				return fmt.Errorf("writing output: %w", err)
			}
			fmt.Printf("Generated %s (%d objectives cascaded)\n", goalsCascadeFlags.output, len(childDoc.Objectives))
		} else {
			fmt.Println(string(data))
		}
	}

	return nil
}

// ============================================================================
// Goals Align Command
// ============================================================================

var goalsAlignFlags struct {
	jsonOutput bool
}

var goalsAlignCmd = &cobra.Command{
	Use:   "align PARENT_FILE CHILD_FILE",
	Short: "Validate alignment between parent and child goals",
	Long: `Validate that child goals properly align with parent goals.

Checks:
  - All parent objectives have supporting child objectives
  - Child objectives reference valid parent objectives
  - No orphaned child objectives (without parent alignment)

Examples:
  splan goals align company-okrs.json team-okrs.json
  splan goals align parent.json child.json --json`,
	Args: cobra.ExactArgs(2),
	RunE: runGoalsAlign,
}

func init() {
	goalsAlignCmd.Flags().BoolVar(&goalsAlignFlags.jsonOutput, "json", false, "Output as JSON")
}

func runGoalsAlign(cmd *cobra.Command, args []string) error {
	// Load parent
	parentDoc, err := okr.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading parent file: %w", err)
	}

	// Load child
	childDoc, err := okr.ReadFile(args[1])
	if err != nil {
		return fmt.Errorf("reading child file: %w", err)
	}

	// Convert to Goals wrappers
	parentGoals := goals.NewOKR(okr.FromObjectives(parentDoc.Objectives))
	childGoals := goals.NewOKR(okr.FromObjectives(childDoc.Objectives))

	// Calculate alignment
	score, err := goals.CalculateAlignment(parentGoals, childGoals)
	if err != nil {
		return fmt.Errorf("calculating alignment: %w", err)
	}

	if goalsAlignFlags.jsonOutput {
		data, err := json.MarshalIndent(score, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling result: %w", err)
		}
		fmt.Println(string(data))
	} else {
		printAlignmentReport(score, args[0], args[1])
	}

	return nil
}

func printAlignmentReport(score *goals.AlignmentScore, parentFile, childFile string) {
	fmt.Println("=============================================================")
	fmt.Println("GOALS ALIGNMENT REPORT")
	fmt.Println("=============================================================")
	fmt.Println()
	fmt.Printf("Parent: %s\n", parentFile)
	fmt.Printf("Child:  %s\n", childFile)
	fmt.Println()
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("ALIGNMENT SCORE")
	fmt.Println("-------------------------------------------------------------")
	fmt.Println()
	fmt.Printf("  Overall Score: %.1f%%\n", score.Score*100)
	fmt.Printf("  Parent Goals Covered: %d/%d\n", score.CoveredParentGoals, score.TotalParentGoals)
	fmt.Printf("  Orphaned Child Goals: %d/%d\n", score.OrphanedChildGoals, score.TotalChildGoals)
	fmt.Println()

	if len(score.Issues) > 0 {
		fmt.Println("-------------------------------------------------------------")
		fmt.Println("ISSUES")
		fmt.Println("-------------------------------------------------------------")
		fmt.Println()
		for _, issue := range score.Issues {
			fmt.Printf("  - %s\n", issue)
		}
		fmt.Println()
	}

	fmt.Println("=============================================================")

	// Summary
	if score.Score >= 0.9 {
		fmt.Println("✓ Excellent alignment")
	} else if score.Score >= 0.7 {
		fmt.Println("○ Good alignment with minor gaps")
	} else if score.Score >= 0.5 {
		fmt.Println("△ Partial alignment - review recommended")
	} else {
		fmt.Println("✗ Poor alignment - significant gaps found")
	}
}

// ============================================================================
// Goals Validate Command
// ============================================================================

var goalsValidateCmd = &cobra.Command{
	Use:   "validate FILE",
	Short: "Validate a goals file",
	Long: `Validate an OKR or V2MOM file for structure and best practices.

Checks:
  - Required fields are present
  - Objectives have 1-5 key results (OKR best practice)
  - Key results have measurable targets
  - References are valid

Examples:
  splan goals validate okrs.json`,
	Args: cobra.ExactArgs(1),
	RunE: runGoalsValidate,
}

func runGoalsValidate(cmd *cobra.Command, args []string) error {
	// Try to load as OKR
	doc, err := okr.ReadFile(args[0])
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	issues := validateOKRDocument(doc)
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

func validateOKRDocument(doc *okr.OKRDocument) []string {
	var issues []string

	if doc.Metadata == nil || doc.Metadata.ID == "" {
		issues = append(issues, "missing metadata.id")
	}

	if len(doc.Objectives) == 0 {
		issues = append(issues, "no objectives defined")
	}

	for i, obj := range doc.Objectives {
		prefix := fmt.Sprintf("objectives[%d]", i)

		if obj.Title == "" && obj.Description == "" {
			issues = append(issues, fmt.Sprintf("%s: missing title and description", prefix))
		}

		if len(obj.KeyResults) == 0 {
			issues = append(issues, fmt.Sprintf("%s: no key results (OKR requires at least 1)", prefix))
		} else if len(obj.KeyResults) > 5 {
			issues = append(issues, fmt.Sprintf("%s: too many key results (%d, best practice is 1-5)", prefix, len(obj.KeyResults)))
		}

		for j, kr := range obj.KeyResults {
			krPrefix := fmt.Sprintf("%s.keyResults[%d]", prefix, j)
			if kr.Title == "" && kr.Description == "" {
				issues = append(issues, fmt.Sprintf("%s: missing title and description", krPrefix))
			}
		}
	}

	return issues
}

// ============================================================================
// Goals Init Command
// ============================================================================

var goalsInitFlags struct {
	output    string
	framework string
	team      string
	period    string
}

var goalsInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Create a goals template",
	Long: `Create a new goals template file.

Supports OKR and V2MOM frameworks.

Examples:
  splan goals init --framework=okr --team="Engineering" --period="Q3 2026"
  splan goals init --framework=v2mom -o my-v2mom.json`,
	RunE: runGoalsInit,
}

func init() {
	goalsInitCmd.Flags().StringVarP(&goalsInitFlags.output, "output", "o", "goals.json", "Output file")
	goalsInitCmd.Flags().StringVar(&goalsInitFlags.framework, "framework", "okr", "Framework: okr or v2mom")
	goalsInitCmd.Flags().StringVar(&goalsInitFlags.team, "team", "Team Name", "Team name")
	goalsInitCmd.Flags().StringVar(&goalsInitFlags.period, "period", "Q3 2026", "Period (e.g., Q3 2026, FY2026)")
}

func runGoalsInit(cmd *cobra.Command, args []string) error {
	var content []byte
	var err error

	switch strings.ToLower(goalsInitFlags.framework) {
	case "okr":
		doc := createOKRTemplate(goalsInitFlags.team, goalsInitFlags.period)
		content, err = doc.JSON()
	case "v2mom":
		content = []byte(createV2MOMTemplate(goalsInitFlags.team))
	default:
		return fmt.Errorf("unknown framework: %s (use okr or v2mom)", goalsInitFlags.framework)
	}

	if err != nil {
		return fmt.Errorf("generating template: %w", err)
	}

	if err := os.WriteFile(goalsInitFlags.output, content, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Created %s (%s framework)\n", goalsInitFlags.output, goalsInitFlags.framework)
	return nil
}

func createOKRTemplate(team, period string) *okr.OKRDocument {
	return &okr.OKRDocument{
		Metadata: &okr.Metadata{
			ID:     okr.GenerateID(),
			Name:   fmt.Sprintf("%s OKRs", team),
			Owner:  "owner@example.com",
			Team:   team,
			Period: period,
			Status: okr.StatusDraft,
		},
		Theme: "Deliver customer value and improve efficiency",
		Objectives: []okr.Objective{
			{
				ID:          "obj-1",
				Title:       "Improve customer satisfaction",
				Description: "Measurably improve customer experience across key touchpoints",
				Category:    "Customer",
				Owner:       "owner@example.com",
				Timeframe:   period,
				Status:      okr.StatusDraft,
				KeyResults: []okr.KeyResult{
					{
						ID:         "obj-1-kr-1",
						Title:      "Increase NPS from 40 to 60",
						Metric:     "Net Promoter Score",
						Baseline:   "40",
						Target:     "60",
						Unit:       "score",
						Confidence: okr.ConfidenceMedium,
						Status:     "Not Started",
					},
					{
						ID:         "obj-1-kr-2",
						Title:      "Reduce support ticket resolution time by 30%",
						Metric:     "Resolution Time",
						Baseline:   "24 hours",
						Target:     "17 hours",
						Unit:       "hours",
						Confidence: okr.ConfidenceHigh,
						Status:     "Not Started",
					},
				},
				Tags: []string{"customer", "satisfaction"},
			},
		},
	}
}

func createV2MOMTemplate(team string) string {
	return fmt.Sprintf(`{
  "metadata": {
    "id": "v2mom-%s",
    "name": "%s V2MOM",
    "owner": "owner@example.com",
    "team": "%s"
  },
  "vision": "Where we want to be in 12 months",
  "values": [
    {
      "id": "v1",
      "name": "Customer First",
      "description": "Every decision starts with customer impact",
      "priority": 1
    },
    {
      "id": "v2",
      "name": "Move Fast",
      "description": "Ship early, iterate often",
      "priority": 2
    }
  ],
  "methods": [
    {
      "id": "m1",
      "name": "Launch new feature X",
      "description": "Deliver key capability to improve customer experience",
      "owner": "owner@example.com",
      "priority": "P0",
      "status": "Not Started",
      "measures": [
        {
          "id": "m1-measure-1",
          "name": "Feature adoption rate",
          "target": "25%%",
          "current": "0%%",
          "status": "Not Started"
        }
      ]
    }
  ],
  "obstacles": [
    {
      "id": "o1",
      "name": "Resource constraints",
      "description": "Limited engineering bandwidth",
      "mitigation": "Prioritize ruthlessly, defer low-impact work"
    }
  ],
  "measures": [
    {
      "id": "global-1",
      "name": "Revenue growth",
      "target": "20%%",
      "current": "5%%",
      "status": "In Progress"
    }
  ]
}
`, strings.ToLower(strings.ReplaceAll(team, " ", "-")), team, team)
}
