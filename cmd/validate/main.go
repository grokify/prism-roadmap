// Package main provides a CLI for cross-repo validation of entity references.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/prism-roadmap/integration"
	"github.com/grokify/prism-roadmap/rmi"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "validate",
	Short: "Cross-repo validation for prism-roadmap",
	Long:  `Validates entity references between prism-roadmap and external systems like productcontext.`,
}

var refsCmd = &cobra.Command{
	Use:   "refs <roadmap-file>",
	Short: "Validate entity references in a roadmap file",
	Long:  `Checks that all idea, goal, and capability references in RMIs point to valid entities.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runRefsValidation,
}

var ideasCmd = &cobra.Command{
	Use:   "ideas <roadmap-file> --ideas-json <ideas-file>",
	Short: "Validate idea references against an ideas export",
	Long:  `Checks that all idea references in RMIs exist in the provided ideas JSON file.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runIdeasValidation,
}

var summaryCmd = &cobra.Command{
	Use:   "summary <roadmap-file>",
	Short: "Show validation summary for a roadmap",
	Long:  `Displays a summary of entity references and potential issues.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runSummary,
}

var (
	ideasFile  string
	outputJSON bool
)

func init() {
	rootCmd.AddCommand(refsCmd)
	rootCmd.AddCommand(ideasCmd)
	rootCmd.AddCommand(summaryCmd)

	ideasCmd.Flags().StringVar(&ideasFile, "ideas-json", "", "Path to ideas JSON export file")
	_ = ideasCmd.MarkFlagRequired("ideas-json")

	rootCmd.PersistentFlags().BoolVar(&outputJSON, "json", false, "Output results as JSON")
}

// ValidationResult holds the validation results.
type ValidationResult struct {
	Valid          bool                       `json:"valid"`
	TotalItems     int                        `json:"totalItems"`
	ItemsWithRefs  int                        `json:"itemsWithRefs"`
	IdeaRefs       *integration.RefValidation `json:"ideaRefs,omitempty"`
	GoalRefs       *integration.RefValidation `json:"goalRefs,omitempty"`
	CapabilityRefs *integration.RefValidation `json:"capabilityRefs,omitempty"`
	OrphanedIdeas  []string                   `json:"orphanedIdeas,omitempty"`
	Errors         []string                   `json:"errors,omitempty"`
}

func runRefsValidation(cmd *cobra.Command, args []string) error {
	roadmapFile := args[0]

	svc, err := rmi.NewServiceFromFile(roadmapFile)
	if err != nil {
		return fmt.Errorf("loading roadmap: %w", err)
	}

	result := &ValidationResult{
		Valid:      true,
		TotalItems: len(svc.List(rmi.ListFilter{})),
	}

	// Collect all refs
	var ideaRefs, goalRefs, capabilityRefs []string
	for _, item := range svc.List(rmi.ListFilter{}) {
		if item.MarketSignal != nil && len(item.MarketSignal.IdeaIDs) > 0 {
			ideaRefs = append(ideaRefs, item.MarketSignal.IdeaIDs...)
			result.ItemsWithRefs++
		}
	}

	// Create validator (without pre-registered IDs, this just shows what refs exist)
	validator := integration.NewRefValidator()

	// Show what refs are in use
	if len(ideaRefs) > 0 {
		result.IdeaRefs = &integration.RefValidation{
			Valid:   false,
			Missing: uniqueStrings(ideaRefs), // All are "missing" without a source of truth
		}
		result.Valid = false
	}
	if len(goalRefs) > 0 {
		refResult := validator.ValidateGoalRefs(goalRefs)
		result.GoalRefs = &refResult
		if !refResult.Valid {
			result.Valid = false
		}
	}
	if len(capabilityRefs) > 0 {
		refResult := validator.ValidateCapabilityRefs(capabilityRefs)
		result.CapabilityRefs = &refResult
		if !refResult.Valid {
			result.Valid = false
		}
	}

	return outputResult(result)
}

func runIdeasValidation(cmd *cobra.Command, args []string) error {
	roadmapFile := args[0]

	svc, err := rmi.NewServiceFromFile(roadmapFile)
	if err != nil {
		return fmt.Errorf("loading roadmap: %w", err)
	}

	// Load ideas from JSON file
	ideasData, err := os.ReadFile(ideasFile)
	if err != nil {
		return fmt.Errorf("reading ideas file: %w", err)
	}

	var ideas []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(ideasData, &ideas); err != nil {
		return fmt.Errorf("parsing ideas JSON: %w", err)
	}

	// Create validator with idea IDs
	validator := integration.NewRefValidator()
	ideaIDs := make(map[string]bool)
	for _, idea := range ideas {
		validator.RegisterIdea(idea.ID)
		ideaIDs[idea.ID] = true
	}

	result := &ValidationResult{
		Valid:      true,
		TotalItems: len(svc.List(rmi.ListFilter{})),
	}

	// Collect idea refs from RMIs
	var ideaRefs []string
	referencedIdeaIDs := make(map[string]bool)
	for _, item := range svc.List(rmi.ListFilter{}) {
		if item.MarketSignal != nil && len(item.MarketSignal.IdeaIDs) > 0 {
			ideaRefs = append(ideaRefs, item.MarketSignal.IdeaIDs...)
			for _, id := range item.MarketSignal.IdeaIDs {
				referencedIdeaIDs[id] = true
			}
			result.ItemsWithRefs++
		}
	}

	// Validate refs
	if len(ideaRefs) > 0 {
		refResult := validator.ValidateIdeaRefs(ideaRefs)
		result.IdeaRefs = &refResult
		if !refResult.Valid {
			result.Valid = false
		}
	}

	// Find orphaned ideas (in productcontext but not referenced by any RMI)
	for _, idea := range ideas {
		if !referencedIdeaIDs[idea.ID] {
			result.OrphanedIdeas = append(result.OrphanedIdeas, idea.ID)
		}
	}

	return outputResult(result)
}

func runSummary(cmd *cobra.Command, args []string) error {
	roadmapFile := args[0]

	svc, err := rmi.NewServiceFromFile(roadmapFile)
	if err != nil {
		return fmt.Errorf("loading roadmap: %w", err)
	}

	summary := svc.Summary()
	items := svc.List(rmi.ListFilter{})

	// Count ref statistics
	itemsWithMarketSignal := 0
	totalIdeaRefs := 0
	uniqueIdeaRefs := make(map[string]bool)

	for _, item := range items {
		if item.MarketSignal != nil {
			if item.MarketSignal.Score > 0 {
				itemsWithMarketSignal++
			}
			for _, id := range item.MarketSignal.IdeaIDs {
				totalIdeaRefs++
				uniqueIdeaRefs[id] = true
			}
		}
	}

	if outputJSON {
		data := map[string]any{
			"totalItems":            summary.TotalItems,
			"itemsWithMarketSignal": itemsWithMarketSignal,
			"totalIdeaRefs":         totalIdeaRefs,
			"uniqueIdeaRefs":        len(uniqueIdeaRefs),
			"moscowCounts":          summary.MoSCoWCounts,
			"statusCounts":          summary.StatusCounts,
		}
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(data)
	}

	fmt.Printf("Roadmap Summary: %s\n", roadmapFile)
	fmt.Printf("================\n\n")
	fmt.Printf("Total Items:              %d\n", summary.TotalItems)
	fmt.Printf("Items with Market Signal: %d\n", itemsWithMarketSignal)
	fmt.Printf("Total Idea References:    %d\n", totalIdeaRefs)
	fmt.Printf("Unique Ideas Referenced:  %d\n", len(uniqueIdeaRefs))
	fmt.Printf("\nMoSCoW Distribution:\n")
	for k, v := range summary.MoSCoWCounts {
		fmt.Printf("  %s: %d\n", k, v)
	}
	fmt.Printf("\nStatus Distribution:\n")
	for k, v := range summary.StatusCounts {
		fmt.Printf("  %s: %d\n", k, v)
	}

	return nil
}

func outputResult(result *ValidationResult) error {
	if outputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(result)
	}

	if result.Valid {
		fmt.Println("✓ Validation passed")
	} else {
		fmt.Println("✗ Validation failed")
	}

	fmt.Printf("\nTotal RMIs: %d\n", result.TotalItems)
	fmt.Printf("RMIs with refs: %d\n", result.ItemsWithRefs)

	if result.IdeaRefs != nil {
		fmt.Printf("\nIdea References:\n")
		if result.IdeaRefs.Valid {
			fmt.Println("  ✓ All valid")
		} else {
			fmt.Printf("  ✗ Missing: %s\n", strings.Join(result.IdeaRefs.Missing, ", "))
		}
	}

	if len(result.OrphanedIdeas) > 0 {
		fmt.Printf("\nOrphaned Ideas (not referenced by any RMI): %d\n", len(result.OrphanedIdeas))
		for _, id := range result.OrphanedIdeas {
			fmt.Printf("  - %s\n", id)
		}
	}

	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, e := range result.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	if !result.Valid {
		os.Exit(1)
	}
	return nil
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
