// Package cli provides the exported Cobra command tree for the PRISM roadmap CLI.
package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/grokify/prism-roadmap/rmi"
)

// ============================================================================
// RMI Parent Command
// ============================================================================

var rmiCmd = &cobra.Command{
	Use:   "rmi",
	Short: "Work with Roadmap Items (RMI)",
	Long: `Commands for managing Roadmap Items with full prioritization support.

RMI combines:
  - MoSCoW prioritization (strategic priority)
  - RICE scoring (quantitative prioritization)
  - Market signals (customer demand)
  - Effort estimation (implementation cost)
  - Complexity factors (risk/dependency tracking)
  - Cross-module references (capabilities, goals, ideas)`,
}

// ============================================================================
// RMI Init Command
// ============================================================================

var rmiInitFlags struct {
	description string
	quarter     string
	output      string
}

var rmiInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new roadmap items file",
	Long: `Create a new RMI JSON file with empty items list.

Example:
  splan rmi init -o roadmap.json
  splan rmi init --description "Q3 2026 Roadmap" --quarter "Q3 2026"`,
	RunE: runRMIInit,
}

// ============================================================================
// RMI List Command
// ============================================================================

var rmiListFlags struct {
	file     string
	moscow   string
	status   string
	quarter  string
	tag      string
	owner    string
	limit    int
	sortBy   string
	sortDesc bool
	json     bool
}

var rmiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List roadmap items",
	Long: `List roadmap items from a JSON file with optional filtering.

Sort options:
  priority      - Sort by composite priority score (default)
  rice          - Sort by RICE score
  market_signal - Sort by market signal score
  created       - Sort by creation date

Examples:
  splan rmi list -f roadmap.json
  splan rmi list -f roadmap.json --moscow must_have
  splan rmi list -f roadmap.json --quarter "Q3 2026" --sort-by priority
  splan rmi list -f roadmap.json --status in_progress --json`,
	RunE: runRMIList,
}

// ============================================================================
// RMI Get Command
// ============================================================================

var rmiGetFlags struct {
	file string
	json bool
}

var rmiGetCmd = &cobra.Command{
	Use:   "get ID",
	Short: "Get a roadmap item by ID",
	Long: `Retrieve detailed information about a specific roadmap item.

Examples:
  splan rmi get rmi-1 -f roadmap.json
  splan rmi get rmi-1 -f roadmap.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runRMIGet,
}

// ============================================================================
// RMI Create Command
// ============================================================================

var rmiCreateFlags struct {
	file        string
	id          string
	name        string
	description string
	moscow      string
	quarter     string
	owner       string
	tags        []string
}

var rmiCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new roadmap item",
	Long: `Create a new roadmap item and add it to a JSON file.

MoSCoW priorities: must_have, should_have, could_have, wont_have

Examples:
  splan rmi create -f roadmap.json --id rmi-1 --name "SSO Integration" --moscow must_have
  splan rmi create -f roadmap.json --id rmi-2 --name "Dark Mode" --moscow should_have --quarter "Q3 2026"`,
	RunE: runRMICreate,
}

// ============================================================================
// RMI Update Command
// ============================================================================

var rmiUpdateFlags struct {
	file        string
	name        string
	description string
	moscow      string
	status      string
	quarter     string
	owner       string
	progress    int
	tags        []string
}

var rmiUpdateCmd = &cobra.Command{
	Use:   "update ID",
	Short: "Update a roadmap item",
	Long: `Update fields of an existing roadmap item.

Status options: planned, in_progress, completed, blocked, cancelled, deferred

Examples:
  splan rmi update rmi-1 -f roadmap.json --status in_progress
  splan rmi update rmi-1 -f roadmap.json --moscow must_have --progress 50`,
	Args: cobra.ExactArgs(1),
	RunE: runRMIUpdate,
}

// ============================================================================
// RMI Delete Command
// ============================================================================

var rmiDeleteFlags struct {
	file string
}

var rmiDeleteCmd = &cobra.Command{
	Use:   "delete ID",
	Short: "Delete a roadmap item",
	Long: `Remove a roadmap item from the JSON file.

Examples:
  splan rmi delete rmi-1 -f roadmap.json`,
	Args: cobra.ExactArgs(1),
	RunE: runRMIDelete,
}

// ============================================================================
// RMI Summary Command
// ============================================================================

var rmiSummaryFlags struct {
	file string
	json bool
}

var rmiSummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show roadmap summary statistics",
	Long: `Display aggregated statistics for all roadmap items.

Examples:
  splan rmi summary -f roadmap.json
  splan rmi summary -f roadmap.json --json`,
	RunE: runRMISummary,
}

// ============================================================================
// RMI Top Command
// ============================================================================

var rmiTopFlags struct {
	file   string
	sortBy string
	limit  int
	json   bool
}

var rmiTopCmd = &cobra.Command{
	Use:   "top",
	Short: "Show top priority items",
	Long: `Display top roadmap items by various scoring methods.

Sort options:
  priority      - Composite priority score (default)
  rice          - RICE score
  market_signal - Market signal score

Examples:
  splan rmi top -f roadmap.json
  splan rmi top -f roadmap.json --sort-by rice --limit 5`,
	RunE: runRMITop,
}

// ============================================================================
// RMI Validate Command
// ============================================================================

var rmiValidateCmd = &cobra.Command{
	Use:   "validate FILE",
	Short: "Validate a roadmap items JSON file",
	Long: `Validate a roadmap items JSON file for structural correctness.

Checks:
  - All items have required fields (ID, Name)
  - No duplicate IDs
  - Valid MoSCoW priorities when set (MoSCoW is optional)
  - Valid RICE/Effort configurations if present

Examples:
  splan rmi validate roadmap.json`,
	Args: cobra.ExactArgs(1),
	RunE: runRMIValidate,
}

// ============================================================================
// Init functions
// ============================================================================

func init() {
	// RMI init flags
	rmiInitCmd.Flags().StringVarP(&rmiInitFlags.output, "output", "o", "roadmap.json", "Output file path")
	rmiInitCmd.Flags().StringVar(&rmiInitFlags.description, "description", "", "Roadmap description")
	rmiInitCmd.Flags().StringVar(&rmiInitFlags.quarter, "quarter", "", "Target quarter (e.g., Q3 2026)")

	// RMI list flags
	rmiListCmd.Flags().StringVarP(&rmiListFlags.file, "file", "f", "roadmap.json", "Input JSON file")
	rmiListCmd.Flags().StringVar(&rmiListFlags.moscow, "moscow", "", "Filter by MoSCoW priority")
	rmiListCmd.Flags().StringVar(&rmiListFlags.status, "status", "", "Filter by status")
	rmiListCmd.Flags().StringVar(&rmiListFlags.quarter, "quarter", "", "Filter by quarter")
	rmiListCmd.Flags().StringVar(&rmiListFlags.tag, "tag", "", "Filter by tag")
	rmiListCmd.Flags().StringVar(&rmiListFlags.owner, "owner", "", "Filter by owner")
	rmiListCmd.Flags().IntVarP(&rmiListFlags.limit, "limit", "n", 0, "Limit number of results")
	rmiListCmd.Flags().StringVar(&rmiListFlags.sortBy, "sort-by", "priority", "Sort by (priority, rice, market_signal, created)")
	rmiListCmd.Flags().BoolVar(&rmiListFlags.sortDesc, "desc", false, "Sort descending (reverse order)")
	rmiListCmd.Flags().BoolVar(&rmiListFlags.json, "json", false, "Output as JSON")

	// RMI get flags
	rmiGetCmd.Flags().StringVarP(&rmiGetFlags.file, "file", "f", "roadmap.json", "Input JSON file")
	rmiGetCmd.Flags().BoolVar(&rmiGetFlags.json, "json", false, "Output as JSON")

	// RMI create flags
	rmiCreateCmd.Flags().StringVarP(&rmiCreateFlags.file, "file", "f", "roadmap.json", "JSON file to update")
	rmiCreateCmd.Flags().StringVar(&rmiCreateFlags.id, "id", "", "Item ID (required)")
	rmiCreateCmd.Flags().StringVar(&rmiCreateFlags.name, "name", "", "Item name (required)")
	rmiCreateCmd.Flags().StringVar(&rmiCreateFlags.description, "description", "", "Item description")
	rmiCreateCmd.Flags().StringVar(&rmiCreateFlags.moscow, "moscow", "should_have", "MoSCoW priority")
	rmiCreateCmd.Flags().StringVar(&rmiCreateFlags.quarter, "quarter", "", "Target quarter")
	rmiCreateCmd.Flags().StringVar(&rmiCreateFlags.owner, "owner", "", "Owner")
	rmiCreateCmd.Flags().StringSliceVar(&rmiCreateFlags.tags, "tags", nil, "Tags (comma-separated)")
	_ = rmiCreateCmd.MarkFlagRequired("id")
	_ = rmiCreateCmd.MarkFlagRequired("name")

	// RMI update flags
	rmiUpdateCmd.Flags().StringVarP(&rmiUpdateFlags.file, "file", "f", "roadmap.json", "JSON file to update")
	rmiUpdateCmd.Flags().StringVar(&rmiUpdateFlags.name, "name", "", "New name")
	rmiUpdateCmd.Flags().StringVar(&rmiUpdateFlags.description, "description", "", "New description")
	rmiUpdateCmd.Flags().StringVar(&rmiUpdateFlags.moscow, "moscow", "", "New MoSCoW priority")
	rmiUpdateCmd.Flags().StringVar(&rmiUpdateFlags.status, "status", "", "New status")
	rmiUpdateCmd.Flags().StringVar(&rmiUpdateFlags.quarter, "quarter", "", "New quarter")
	rmiUpdateCmd.Flags().StringVar(&rmiUpdateFlags.owner, "owner", "", "New owner")
	rmiUpdateCmd.Flags().IntVar(&rmiUpdateFlags.progress, "progress", -1, "Progress percentage (0-100)")
	rmiUpdateCmd.Flags().StringSliceVar(&rmiUpdateFlags.tags, "tags", nil, "Tags (comma-separated)")

	// RMI delete flags
	rmiDeleteCmd.Flags().StringVarP(&rmiDeleteFlags.file, "file", "f", "roadmap.json", "JSON file to update")

	// RMI summary flags
	rmiSummaryCmd.Flags().StringVarP(&rmiSummaryFlags.file, "file", "f", "roadmap.json", "Input JSON file")
	rmiSummaryCmd.Flags().BoolVar(&rmiSummaryFlags.json, "json", false, "Output as JSON")

	// RMI top flags
	rmiTopCmd.Flags().StringVarP(&rmiTopFlags.file, "file", "f", "roadmap.json", "Input JSON file")
	rmiTopCmd.Flags().StringVar(&rmiTopFlags.sortBy, "sort-by", "priority", "Sort by (priority, rice, market_signal)")
	rmiTopCmd.Flags().IntVarP(&rmiTopFlags.limit, "limit", "n", 10, "Number of items to show")
	rmiTopCmd.Flags().BoolVar(&rmiTopFlags.json, "json", false, "Output as JSON")

	// Add subcommands
	rmiCmd.AddCommand(rmiInitCmd)
	rmiCmd.AddCommand(rmiListCmd)
	rmiCmd.AddCommand(rmiGetCmd)
	rmiCmd.AddCommand(rmiCreateCmd)
	rmiCmd.AddCommand(rmiUpdateCmd)
	rmiCmd.AddCommand(rmiDeleteCmd)
	rmiCmd.AddCommand(rmiSummaryCmd)
	rmiCmd.AddCommand(rmiTopCmd)
	rmiCmd.AddCommand(rmiValidateCmd)
}

// ============================================================================
// Command Implementations
// ============================================================================

func runRMIInit(cmd *cobra.Command, args []string) error {
	// Check if file already exists
	if _, err := os.Stat(rmiInitFlags.output); err == nil {
		return fmt.Errorf("file already exists: %s (use -o to specify a different output path)", rmiInitFlags.output)
	}

	set := rmi.NewRoadmapItemSet()
	set.Description = rmiInitFlags.description
	set.Quarter = rmiInitFlags.quarter

	if err := set.WriteFile(rmiInitFlags.output); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	fmt.Printf("Created: %s\n", rmiInitFlags.output)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Add items: splan rmi create -f " + rmiInitFlags.output + " --id rmi-1 --name \"My Item\" --moscow must_have")
	fmt.Println("  2. List items: splan rmi list -f " + rmiInitFlags.output)
	fmt.Println("  3. View summary: splan rmi summary -f " + rmiInitFlags.output)

	return nil
}

func runRMIList(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiListFlags.file)
	if err != nil {
		return err
	}

	filter := rmi.ListFilter{
		MoSCoW:   prioritization.MoSCoWPriority(rmiListFlags.moscow),
		Status:   rmi.RMIStatus(rmiListFlags.status),
		Quarter:  rmiListFlags.quarter,
		Tag:      rmiListFlags.tag,
		Owner:    rmiListFlags.owner,
		Limit:    rmiListFlags.limit,
		SortBy:   rmiListFlags.sortBy,
		SortDesc: rmiListFlags.sortDesc,
	}

	items := svc.List(filter)

	if rmiListFlags.json {
		output, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(output))
		return nil
	}

	if len(items) == 0 {
		fmt.Println("No items found.")
		return nil
	}

	fmt.Printf("Found %d item(s):\n\n", len(items))
	for _, item := range items {
		statusIcon := "[ ]"
		switch item.Status {
		case rmi.RMIStatusInProgress:
			statusIcon = "[~]"
		case rmi.RMIStatusCompleted:
			statusIcon = "[x]"
		case rmi.RMIStatusBlocked:
			statusIcon = "[!]"
		case rmi.RMIStatusCancelled:
			statusIcon = "[-]"
		case rmi.RMIStatusDeferred:
			statusIcon = "[>]"
		}

		fmt.Printf("%s %s (%s) - %s\n", statusIcon, item.ID, item.MoSCoW, item.Name)
		if item.Quarter != "" {
			fmt.Printf("    Quarter: %s\n", item.Quarter)
		}
		if item.Owner != "" {
			fmt.Printf("    Owner: %s\n", item.Owner)
		}
		if len(item.Tags) > 0 {
			fmt.Printf("    Tags: %s\n", strings.Join(item.Tags, ", "))
		}
		fmt.Printf("    Priority Score: %.2f\n", item.PriorityScore())
	}

	return nil
}

func runRMIGet(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiGetFlags.file)
	if err != nil {
		return err
	}

	item, err := svc.Get(args[0])
	if err != nil {
		return err
	}

	if rmiGetFlags.json {
		output, err := json.MarshalIndent(item, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(output))
		return nil
	}

	fmt.Printf("ID:          %s\n", item.ID)
	fmt.Printf("Name:        %s\n", item.Name)
	if item.Description != "" {
		fmt.Printf("Description: %s\n", item.Description)
	}
	fmt.Printf("MoSCoW:      %s\n", item.MoSCoW)
	fmt.Printf("Status:      %s\n", item.Status)
	if item.Quarter != "" {
		fmt.Printf("Quarter:     %s\n", item.Quarter)
	}
	if item.Owner != "" {
		fmt.Printf("Owner:       %s\n", item.Owner)
	}
	if item.Progress != nil {
		fmt.Printf("Progress:    %d%%\n", *item.Progress)
	}
	if len(item.Tags) > 0 {
		fmt.Printf("Tags:        %s\n", strings.Join(item.Tags, ", "))
	}
	fmt.Printf("Priority:    %.2f\n", item.PriorityScore())
	if item.RICE != nil {
		fmt.Printf("RICE Score:  %.2f\n", item.RICE.Score)
	}
	if item.MarketSignal != nil {
		fmt.Printf("Market Sig:  %.2f\n", item.MarketSignal.Score)
	}
	if item.Effort != nil {
		fmt.Printf("Effort:      %d person-days\n", item.EffectiveEffortDays())
	}
	fmt.Printf("Created:     %s\n", item.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated:     %s\n", item.UpdatedAt.Format("2006-01-02 15:04:05"))

	return nil
}

func runRMICreate(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiCreateFlags.file)
	if err != nil {
		// File doesn't exist, create new service
		svc = rmi.NewService()
	}

	// MoSCoW is optional — the flag defaults to should_have, but an explicit
	// --moscow="" creates an unprioritized item. Validate only when set.
	moscow := prioritization.MoSCoWPriority(rmiCreateFlags.moscow)
	if moscow != prioritization.MoSCoWUnspecified && !prioritization.IsValidMoSCoWPriority(moscow) {
		return fmt.Errorf("invalid moscow priority: %s (expected: must_have, should_have, could_have, wont_have)", rmiCreateFlags.moscow)
	}

	input := rmi.CreateInput{
		ID:          rmiCreateFlags.id,
		Name:        rmiCreateFlags.name,
		Description: rmiCreateFlags.description,
		MoSCoW:      moscow,
		Quarter:     rmiCreateFlags.quarter,
		Owner:       rmiCreateFlags.owner,
		Tags:        rmiCreateFlags.tags,
	}

	item, err := svc.Create(input)
	if err != nil {
		return err
	}

	if err := svc.SaveAs(rmiCreateFlags.file); err != nil {
		return err
	}

	fmt.Printf("Created item: %s\n", item.ID)
	fmt.Printf("  Name:   %s\n", item.Name)
	fmt.Printf("  MoSCoW: %s\n", item.MoSCoW)

	return nil
}

func runRMIUpdate(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiUpdateFlags.file)
	if err != nil {
		return err
	}

	input := rmi.UpdateInput{}

	if rmiUpdateFlags.name != "" {
		input.Name = &rmiUpdateFlags.name
	}
	if rmiUpdateFlags.description != "" {
		input.Description = &rmiUpdateFlags.description
	}
	if rmiUpdateFlags.moscow != "" {
		moscow := prioritization.MoSCoWPriority(rmiUpdateFlags.moscow)
		if !prioritization.IsValidMoSCoWPriority(moscow) {
			return fmt.Errorf("invalid moscow priority: %s", rmiUpdateFlags.moscow)
		}
		input.MoSCoW = &moscow
	}
	if rmiUpdateFlags.status != "" {
		status := rmi.RMIStatus(rmiUpdateFlags.status)
		input.Status = &status
	}
	if rmiUpdateFlags.quarter != "" {
		input.Quarter = &rmiUpdateFlags.quarter
	}
	if rmiUpdateFlags.owner != "" {
		input.Owner = &rmiUpdateFlags.owner
	}
	if rmiUpdateFlags.progress >= 0 && rmiUpdateFlags.progress <= 100 {
		input.Progress = &rmiUpdateFlags.progress
	}
	if rmiUpdateFlags.tags != nil {
		input.Tags = rmiUpdateFlags.tags
	}

	item, updated, err := svc.Update(args[0], input)
	if err != nil {
		return err
	}

	if !updated {
		fmt.Println("No changes made.")
		return nil
	}

	if err := svc.Save(); err != nil {
		return err
	}

	fmt.Printf("Updated item: %s\n", item.ID)

	return nil
}

func runRMIDelete(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiDeleteFlags.file)
	if err != nil {
		return err
	}

	if err := svc.Delete(args[0]); err != nil {
		return err
	}

	if err := svc.Save(); err != nil {
		return err
	}

	fmt.Printf("Deleted item: %s\n", args[0])

	return nil
}

func runRMISummary(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiSummaryFlags.file)
	if err != nil {
		return err
	}

	summary := svc.Summary()

	if rmiSummaryFlags.json {
		output, err := json.MarshalIndent(summary, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(output))
		return nil
	}

	fmt.Printf("Roadmap Summary: %s\n\n", rmiSummaryFlags.file)

	fmt.Printf("Total Items:        %d\n", summary.TotalItems)
	fmt.Printf("Total Effort:       %d person-days\n", summary.TotalEffortDays)
	fmt.Printf("Avg Priority Score: %.2f\n\n", summary.AvgPriorityScore)

	fmt.Println("By MoSCoW Priority:")
	for _, moscow := range prioritization.AllMoSCoWPriorities() {
		count := summary.MoSCoWCounts[moscow]
		if count > 0 {
			fmt.Printf("  %s: %d\n", moscow, count)
		}
	}

	fmt.Println("\nBy Status:")
	for _, status := range rmi.ValidRMIStatuses() {
		count := summary.StatusCounts[status]
		if count > 0 {
			fmt.Printf("  %s: %d\n", status, count)
		}
	}

	if len(summary.QuarterCounts) > 0 {
		fmt.Println("\nBy Quarter:")
		for quarter, count := range summary.QuarterCounts {
			fmt.Printf("  %s: %d\n", quarter, count)
		}
	}

	fmt.Println()
	fmt.Printf("Actionable: %d\n", summary.ActionableCount)
	fmt.Printf("Blocked:    %d\n", summary.BlockedCount)

	return nil
}

func runRMITop(cmd *cobra.Command, args []string) error {
	svc, err := rmi.NewServiceFromFile(rmiTopFlags.file)
	if err != nil {
		return err
	}

	var items []rmi.RoadmapItem
	switch rmiTopFlags.sortBy {
	case "rice":
		items = svc.TopByRICE(rmiTopFlags.limit)
	case "market_signal":
		items = svc.TopByMarketSignal(rmiTopFlags.limit)
	default:
		items = svc.TopByPriority(rmiTopFlags.limit)
	}

	if rmiTopFlags.json {
		output, err := json.MarshalIndent(items, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling JSON: %w", err)
		}
		fmt.Println(string(output))
		return nil
	}

	if len(items) == 0 {
		fmt.Println("No items found.")
		return nil
	}

	fmt.Printf("Top %d items by %s:\n\n", len(items), rmiTopFlags.sortBy)

	for i, item := range items {
		var score float64
		switch rmiTopFlags.sortBy {
		case "rice":
			if item.RICE != nil {
				score = item.RICE.Score
			}
		case "market_signal":
			if item.MarketSignal != nil {
				score = item.MarketSignal.Score
			}
		default:
			score = item.PriorityScore()
		}
		fmt.Printf("%2d. %s (%s) - Score: %.2f\n", i+1, item.Name, item.MoSCoW, score)
		fmt.Printf("    ID: %s\n", item.ID)
		if item.Quarter != "" {
			fmt.Printf("    Quarter: %s\n", item.Quarter)
		}
	}

	return nil
}

func runRMIValidate(cmd *cobra.Command, args []string) error {
	filepath := args[0]

	// Check file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filepath)
	}

	set, err := rmi.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading file: %w", err)
	}

	if err := set.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("Valid RMI file: %s\n", filepath)
	fmt.Printf("  Items: %d\n", len(set.Items))
	if set.Description != "" {
		fmt.Printf("  Description: %s\n", set.Description)
	}
	if set.Quarter != "" {
		fmt.Printf("  Quarter: %s\n", set.Quarter)
	}

	// Summary of MoSCoW distribution
	summary := set.MoSCoWSummary()
	fmt.Println("  MoSCoW Distribution:")
	for _, moscow := range prioritization.AllMoSCoWPriorities() {
		count := summary[moscow]
		if count > 0 {
			fmt.Printf("    %s: %d\n", moscow, count)
		}
	}

	return nil
}
