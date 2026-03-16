// Command splan works with structured planning documents (PRD, MRD, TRD, OKR, V2MOM, Roadmap).
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/grokify/structured-plan/goals/okr"
	okrrender "github.com/grokify/structured-plan/goals/okr/render"
	okrmarp "github.com/grokify/structured-plan/goals/okr/render/marp"
	"github.com/grokify/structured-plan/goals/v2mom"
	v2momrender "github.com/grokify/structured-plan/goals/v2mom/render"
	v2mommarp "github.com/grokify/structured-plan/goals/v2mom/render/marp"
	"github.com/grokify/structured-plan/requirements/mrd"
	"github.com/grokify/structured-plan/requirements/prd"
	"github.com/grokify/structured-plan/requirements/prd/render/terminal"
	"github.com/grokify/structured-plan/requirements/trd"
	"github.com/grokify/structured-plan/schema"
	"github.com/plexusone/structured-evaluation/evaluation"
)

// Set by GoReleaser ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "splan",
	Short: "Structured planning document generator and validator",
	Long: `splan is a CLI tool for working with structured planning documents.

Document types:
  Requirements:
    - PRD (Product Requirements Document)
    - MRD (Market Requirements Document)
    - TRD (Technical Requirements Document)

  Goals:
    - OKR (Objectives and Key Results) - coming soon
    - V2MOM (Vision, Values, Methods, Obstacles, Measures)

  Roadmap (coming soon):
    - Standalone roadmaps for portfolio/product planning

It can convert document JSON files to markdown with Pandoc-compatible YAML
frontmatter, generate Marp presentations, and validate files against their
respective schemas.

Example usage:
  splan requirements prd generate myproduct.prd.json
  splan goals v2mom validate my-v2mom.json
  splan goals v2mom generate marp my-v2mom.json -o slides.md
  splan schema generate --type prd`,
	Version: version,
}

func init() {
	rootCmd.SetVersionTemplate("splan version {{.Version}} (commit: " + commit + ", built: " + date + ")\n")
}

// ============================================================================
// Requirements Parent Command
// ============================================================================

var requirementsCmd = &cobra.Command{
	Use:     "requirements",
	Aliases: []string{"req"},
	Short:   "Work with requirements documents (PRD, MRD, TRD)",
	Long: `Commands for generating and validating requirements documents.

Supported document types:
  - prd: Product Requirements Document
  - mrd: Market Requirements Document
  - trd: Technical Requirements Document`,
}

func init() {
	// Add top-level commands
	rootCmd.AddCommand(requirementsCmd)
	rootCmd.AddCommand(goalsCmd)
	rootCmd.AddCommand(schemaCmd)
	rootCmd.AddCommand(mergeCmd)

	// Add requirements subcommands
	requirementsCmd.AddCommand(prdCmd)
	requirementsCmd.AddCommand(mrdCmd)
	requirementsCmd.AddCommand(trdCmd)

	// Add goals subcommands
	goalsCmd.AddCommand(v2momCmd)
	goalsCmd.AddCommand(okrCmd)
}

// ============================================================================
// Merge Command
// ============================================================================

var mergeFlags struct {
	output string
}

var mergeCmd = &cobra.Command{
	Use:   "merge [files...]",
	Short: "Merge multiple JSON files into one",
	Long: `Merge multiple JSON files into one.

The files are merged in the order they are provided. For nested objects,
values are recursively merged. For arrays, values are concatenated.`,
	Example: `  splan merge file1.json file2.json -o merged.json
  splan merge base.prd.json overrides.json -o final.prd.json`,
	Args: cobra.MinimumNArgs(1),
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeFlags.output, "output", "o", "merged.json", "Output file name")
}

func runMerge(cmd *cobra.Command, args []string) error {
	var mergedData map[string]interface{}

	for _, file := range args {
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("reading file %s: %w", file, err)
		}

		var currentData map[string]interface{}
		if err := json.Unmarshal(data, &currentData); err != nil {
			return fmt.Errorf("unmarshaling json from file %s: %w", file, err)
		}

		if mergedData == nil {
			mergedData = currentData
		} else {
			mergedData = deepMerge(mergedData, currentData)
		}
	}

	mergedJSON, err := json.MarshalIndent(mergedData, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling merged data to json: %w", err)
	}

	if err := os.WriteFile(mergeFlags.output, mergedJSON, 0600); err != nil {
		return fmt.Errorf("writing merged json to file %s: %w", mergeFlags.output, err)
	}

	fmt.Printf("Successfully merged %d files into %s\n", len(args), mergeFlags.output)
	return nil
}

func deepMerge(a, b map[string]interface{}) map[string]interface{} {
	for k, v := range b {
		if va, ok := a[k]; ok {
			// Both have this key - check types for merging
			switch vaTyped := va.(type) {
			case map[string]interface{}:
				// Both are maps - recursively merge
				if vMap, ok := v.(map[string]interface{}); ok {
					a[k] = deepMerge(vaTyped, vMap)
					continue
				}
			case []interface{}:
				// Both are arrays - concatenate them
				if vSlice, ok := v.([]interface{}); ok {
					a[k] = append(vaTyped, vSlice...)
					continue
				}
			}
		}
		// Key doesn't exist in a, or types don't match - use b's value
		a[k] = v
	}
	return a
}

// ============================================================================
// Goals Parent Command
// ============================================================================

var goalsCmd = &cobra.Command{
	Use:   "goals",
	Short: "Work with goal frameworks (OKR, V2MOM)",
	Long: `Commands for generating and validating goal framework documents.

Supported frameworks:
  - v2mom: Vision, Values, Methods, Obstacles, Measures
  - okr: Objectives and Key Results (coming soon)`,
}

// ============================================================================
// V2MOM Commands
// ============================================================================

var v2momCmd = &cobra.Command{
	Use:   "v2mom",
	Short: "Work with V2MOM documents",
	Long: `Commands for generating and validating V2MOM (Vision, Values, Methods, Obstacles, Measures) documents.

V2MOM supports:
  - JSON validation against the V2MOM schema
  - Marp markdown slide generation
  - Both traditional flat V2MOM and OKR-aligned nested structures
  - Multiple terminology modes (V2MOM, OKR, hybrid)`,
}

var v2momValidateFlags struct {
	structure string
}

var v2momValidateCmd = &cobra.Command{
	Use:   "validate FILE",
	Short: "Validate a V2MOM JSON file",
	Long: `Validate a V2MOM JSON file against the schema and structural rules.

Structure modes:
  flat    - Traditional V2MOM (measures/obstacles at V2MOM level only)
  nested  - OKR-aligned (measures under Methods, global obstacles allowed)
  hybrid  - Both levels allowed (default)

Examples:
  splan goals v2mom validate my-v2mom.json
  splan goals v2mom validate my-v2mom.json --structure=nested`,
	Args: cobra.ExactArgs(1),
	RunE: runV2MOMValidate,
}

var v2momGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output from V2MOM",
	Long:  `Generate various output formats from a V2MOM JSON file.`,
}

var v2momGenerateMarpFlags struct {
	output      string
	theme       string
	terminology string
}

var v2momGenerateMarpCmd = &cobra.Command{
	Use:   "marp FILE",
	Short: "Generate Marp markdown slides",
	Long: `Generate Marp markdown presentation slides from a V2MOM JSON file.

Themes:
  default   - Clean gradient theme (default)
  corporate - Professional blue theme
  minimal   - Simple grayscale theme

Terminology:
  v2mom  - Use V2MOM terms: Methods, Measures, Obstacles (default)
  okr    - Use OKR terms: Objectives, Key Results, Risks
  hybrid - Show both: Methods (Objectives), Measures (Key Results)

Examples:
  splan goals v2mom generate marp my-v2mom.json
  splan goals v2mom generate marp my-v2mom.json -o slides.md
  splan goals v2mom generate marp my-v2mom.json --theme=corporate --terminology=okr`,
	Args: cobra.ExactArgs(1),
	RunE: runV2MOMMarpGenerate,
}

var v2momInitFlags struct {
	name      string
	output    string
	structure string
}

var v2momInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new V2MOM template",
	Long: `Create a new V2MOM JSON template file with example content.

Structure modes:
  flat   - Traditional V2MOM with measures/obstacles at top level
  nested - OKR-aligned with measures under methods (default)
  hybrid - Both levels with examples

Examples:
  splan goals v2mom init
  splan goals v2mom init --name "FY2026 Product Strategy"
  splan goals v2mom init --name "Engineering Goals" -o engineering-v2mom.json --structure=nested`,
	RunE: runV2MOMInit,
}

func init() {
	// V2MOM validate flags
	v2momValidateCmd.Flags().StringVar(&v2momValidateFlags.structure, "structure", "", "Structure mode to validate against (flat, nested, hybrid)")

	// V2MOM generate marp flags
	v2momGenerateMarpCmd.Flags().StringVarP(&v2momGenerateMarpFlags.output, "output", "o", "", "Output file path (default: stdout)")
	v2momGenerateMarpCmd.Flags().StringVar(&v2momGenerateMarpFlags.theme, "theme", "default", "Slide theme (default, corporate, minimal)")
	v2momGenerateMarpCmd.Flags().StringVar(&v2momGenerateMarpFlags.terminology, "terminology", "", "Display terminology (v2mom, okr, hybrid)")

	// V2MOM init flags
	v2momInitCmd.Flags().StringVar(&v2momInitFlags.name, "name", "My V2MOM", "Name for the V2MOM")
	v2momInitCmd.Flags().StringVarP(&v2momInitFlags.output, "output", "o", "v2mom.json", "Output file path")
	v2momInitCmd.Flags().StringVar(&v2momInitFlags.structure, "structure", "nested", "Structure mode (flat, nested, hybrid)")

	// Add subcommands
	v2momGenerateCmd.AddCommand(v2momGenerateMarpCmd)
	v2momCmd.AddCommand(v2momValidateCmd)
	v2momCmd.AddCommand(v2momGenerateCmd)
	v2momCmd.AddCommand(v2momInitCmd)
}

func runV2MOMValidate(cmd *cobra.Command, args []string) error {
	filepath := args[0]

	// Check file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filepath)
	}

	// Read and parse V2MOM
	v, err := v2mom.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading V2MOM: %w", err)
	}

	// Set up validation options
	opts := v2mom.DefaultValidationOptions()
	if v2momValidateFlags.structure != "" {
		opts.Structure = v2momValidateFlags.structure
	}

	// Validate
	errs := v.Validate(opts)

	// Report results
	errors := v2mom.Errors(errs)
	warnings := v2mom.Warnings(errs)

	if len(warnings) > 0 {
		fmt.Println("Warnings:")
		for _, w := range warnings {
			fmt.Printf("  - %s\n", w)
		}
		fmt.Println()
	}

	if len(errors) > 0 {
		fmt.Println("Errors:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d error(s)", len(errors))
	}

	// Print success info
	fmt.Printf("Valid V2MOM: %s\n", filepath)
	fmt.Printf("  Structure: %s\n", v.GetStructure())
	fmt.Printf("  Methods: %d\n", len(v.Methods))
	fmt.Printf("  Total Measures: %d\n", len(v.AllMeasures()))
	fmt.Printf("  Total Obstacles: %d\n", len(v.AllObstacles()))

	if v.Metadata != nil && v.Metadata.Name != "" {
		fmt.Printf("  Name: %s\n", v.Metadata.Name)
	}

	return nil
}

func runV2MOMMarpGenerate(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Check file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputPath)
	}

	// Read and parse V2MOM
	v, err := v2mom.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading V2MOM: %w", err)
	}

	// Create renderer and options
	renderer := v2mommarp.New()
	opts := v2momrender.DefaultOptions()
	opts.Theme = v2momGenerateMarpFlags.theme
	if v2momGenerateMarpFlags.terminology != "" {
		opts.Terminology = v2momGenerateMarpFlags.terminology
	}

	// Render
	output, err := renderer.Render(v, opts)
	if err != nil {
		return fmt.Errorf("rendering Marp: %w", err)
	}

	// Write output
	if v2momGenerateMarpFlags.output != "" {
		// Ensure output directory exists
		dir := filepath.Dir(v2momGenerateMarpFlags.output)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}
		}

		if err := os.WriteFile(v2momGenerateMarpFlags.output, output, 0600); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Printf("Generated: %s\n", v2momGenerateMarpFlags.output)
	} else {
		// Write to stdout
		fmt.Print(string(output))
	}

	return nil
}

func runV2MOMInit(cmd *cobra.Command, args []string) error {
	// Check if file already exists
	if _, err := os.Stat(v2momInitFlags.output); err == nil {
		return fmt.Errorf("file already exists: %s (use -o to specify a different output path)", v2momInitFlags.output)
	}

	// Create template based on structure
	var template *v2mom.V2MOM

	switch v2momInitFlags.structure {
	case "flat":
		template = createFlatV2MOMTemplate(v2momInitFlags.name)
	case "hybrid":
		template = createHybridV2MOMTemplate(v2momInitFlags.name)
	default: // "nested"
		template = createNestedV2MOMTemplate(v2momInitFlags.name)
	}

	// Write to file
	if err := template.WriteFile(v2momInitFlags.output); err != nil {
		return fmt.Errorf("writing template: %w", err)
	}

	fmt.Printf("Created: %s\n", v2momInitFlags.output)
	fmt.Printf("  Structure: %s\n", v2momInitFlags.structure)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit the file to add your vision, values, methods, and measures")
	fmt.Println("  2. Run 'splan goals v2mom validate " + v2momInitFlags.output + "' to check your V2MOM")
	fmt.Println("  3. Run 'splan goals v2mom generate marp " + v2momInitFlags.output + " -o slides.md' to create slides")

	return nil
}

func createNestedV2MOMTemplate(name string) *v2mom.V2MOM {
	now := time.Now()
	return &v2mom.V2MOM{
		Schema: "../schema/v2mom.schema.json",
		Metadata: &v2mom.Metadata{
			Name:        name,
			Author:      "Your Name",
			Team:        "Your Team",
			FiscalYear:  fmt.Sprintf("FY%d", now.Year()),
			Quarter:     "Q1",
			Version:     "1.0.0",
			Status:      v2mom.StatusDraft,
			Structure:   v2mom.StructureNested,
			Terminology: v2mom.TerminologyV2MOM,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Vision: "Describe your vision here - what do you want to achieve?",
		Values: []v2mom.Value{
			{
				Name:        "Value 1",
				Description: "What's most important to you?",
				Priority:    1,
			},
			{
				Name:        "Value 2",
				Description: "What's the second most important principle?",
				Priority:    2,
			},
		},
		Methods: []v2mom.Method{
			{
				ID:          "method-1",
				Name:        "First Method/Objective",
				Description: "How will you achieve your vision? What's the first major initiative?",
				Priority:    v2mom.PriorityP0,
				Status:      "Not Started",
				Measures: []v2mom.Measure{
					{
						ID:       "m1-kr1",
						Name:     "Key Result 1",
						Target:   "Define your target",
						Status:   "Not Started",
						Progress: 0,
					},
					{
						ID:       "m1-kr2",
						Name:     "Key Result 2",
						Target:   "Define your target",
						Status:   "Not Started",
						Progress: 0,
					},
				},
				Obstacles: []v2mom.Obstacle{
					{
						Name:       "Method-specific obstacle",
						Severity:   "Medium",
						Mitigation: "How will you address this?",
					},
				},
			},
			{
				ID:          "method-2",
				Name:        "Second Method/Objective",
				Description: "What's the second major initiative?",
				Priority:    v2mom.PriorityP1,
				Status:      "Not Started",
				Measures: []v2mom.Measure{
					{
						ID:       "m2-kr1",
						Name:     "Key Result 1",
						Target:   "Define your target",
						Status:   "Not Started",
						Progress: 0,
					},
				},
			},
		},
		Obstacles: []v2mom.Obstacle{
			{
				ID:          "obs-global-1",
				Name:        "Global obstacle",
				Description: "What's preventing success across multiple methods?",
				Severity:    "High",
				Likelihood:  "Medium",
				Mitigation:  "How will you mitigate this risk?",
				Status:      "Identified",
			},
		},
	}
}

func createFlatV2MOMTemplate(name string) *v2mom.V2MOM {
	now := time.Now()
	return &v2mom.V2MOM{
		Schema: "../schema/v2mom.schema.json",
		Metadata: &v2mom.Metadata{
			Name:        name,
			Author:      "Your Name",
			Team:        "Your Team",
			FiscalYear:  fmt.Sprintf("FY%d", now.Year()),
			Quarter:     "Q1",
			Version:     "1.0.0",
			Status:      v2mom.StatusDraft,
			Structure:   v2mom.StructureFlat,
			Terminology: v2mom.TerminologyV2MOM,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Vision: "Describe your vision here - what do you want to achieve?",
		Values: []v2mom.Value{
			{
				Name:        "Value 1",
				Description: "What's most important to you?",
				Priority:    1,
			},
			{
				Name:        "Value 2",
				Description: "What's the second most important principle?",
				Priority:    2,
			},
		},
		Methods: []v2mom.Method{
			{
				Name:        "First Method",
				Description: "How will you achieve your vision?",
				Priority:    v2mom.PriorityP0,
			},
			{
				Name:        "Second Method",
				Description: "What's the second major action?",
				Priority:    v2mom.PriorityP1,
			},
		},
		Obstacles: []v2mom.Obstacle{
			{
				Name:        "Obstacle 1",
				Description: "What's preventing success?",
				Severity:    "High",
				Mitigation:  "How will you address this?",
			},
		},
		Measures: []v2mom.Measure{
			{
				Name:   "Measure 1",
				Target: "Define your target",
				Status: "Not Started",
			},
			{
				Name:   "Measure 2",
				Target: "Define your target",
				Status: "Not Started",
			},
		},
	}
}

func createHybridV2MOMTemplate(name string) *v2mom.V2MOM {
	now := time.Now()
	return &v2mom.V2MOM{
		Schema: "../schema/v2mom.schema.json",
		Metadata: &v2mom.Metadata{
			Name:        name,
			Author:      "Your Name",
			Team:        "Your Team",
			FiscalYear:  fmt.Sprintf("FY%d", now.Year()),
			Quarter:     "Q1",
			Version:     "1.0.0",
			Status:      v2mom.StatusDraft,
			Structure:   v2mom.StructureHybrid,
			Terminology: v2mom.TerminologyV2MOM,
			CreatedAt:   now,
			UpdatedAt:   now,
		},
		Vision: "Describe your vision here - what do you want to achieve?",
		Values: []v2mom.Value{
			{
				Name:        "Value 1",
				Description: "What's most important to you?",
				Priority:    1,
			},
		},
		Methods: []v2mom.Method{
			{
				ID:          "method-1",
				Name:        "Method with nested measures",
				Description: "This method has its own key results",
				Priority:    v2mom.PriorityP0,
				Status:      "Not Started",
				Measures: []v2mom.Measure{
					{
						Name:   "Method-specific KR",
						Target: "Define target",
					},
				},
			},
			{
				ID:          "method-2",
				Name:        "Method without nested measures",
				Description: "This method uses global measures",
				Priority:    v2mom.PriorityP1,
			},
		},
		Obstacles: []v2mom.Obstacle{
			{
				Name:        "Global obstacle",
				Description: "Affects multiple methods",
				Severity:    "High",
				Mitigation:  "Mitigation strategy",
			},
		},
		Measures: []v2mom.Measure{
			{
				Name:        "North star metric",
				Description: "Global measure spanning all methods",
				Target:      "Define target",
			},
		},
	}
}

// ============================================================================
// OKR Commands
// ============================================================================

var okrCmd = &cobra.Command{
	Use:   "okr",
	Short: "Work with OKR documents",
	Long: `Commands for generating and validating OKR (Objectives and Key Results) documents.

OKR is a goal-setting framework that helps teams define objectives and track
measurable key results. Originally developed at Intel and popularized by Google.`,
}

var okrValidateFlags struct {
	strict bool
}

var okrValidateCmd = &cobra.Command{
	Use:   "validate FILE",
	Short: "Validate an OKR JSON file",
	Long: `Validate an OKR JSON file against the schema and structural rules.

Examples:
  splan goals okr validate my-okrs.json
  splan goals okr validate my-okrs.json --strict`,
	Args: cobra.ExactArgs(1),
	RunE: runOKRValidate,
}

var okrGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate output from OKR",
	Long:  `Generate various output formats from an OKR JSON file.`,
}

var okrGenerateMarpFlags struct {
	output string
	theme  string
}

var okrGenerateMarpCmd = &cobra.Command{
	Use:   "marp FILE",
	Short: "Generate Marp markdown slides",
	Long: `Generate Marp markdown presentation slides from an OKR JSON file.

Themes:
  default   - Clean gradient theme (default)
  corporate - Professional blue theme
  minimal   - Simple grayscale theme

Examples:
  splan goals okr generate marp my-okrs.json
  splan goals okr generate marp my-okrs.json -o slides.md
  splan goals okr generate marp my-okrs.json --theme=corporate`,
	Args: cobra.ExactArgs(1),
	RunE: runOKRMarpGenerate,
}

var okrInitFlags struct {
	name   string
	output string
	period string
}

var okrInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new OKR template",
	Long: `Create a new OKR JSON template file with example content.

Examples:
  splan goals okr init
  splan goals okr init --name "Q1 2026 Engineering OKRs"
  splan goals okr init --name "Product Team" -o product-okrs.json --period 2026-Q1`,
	RunE: runOKRInit,
}

func init() {
	// OKR validate flags
	okrValidateCmd.Flags().BoolVar(&okrValidateFlags.strict, "strict", false, "Enable strict validation mode")

	// OKR generate marp flags
	okrGenerateMarpCmd.Flags().StringVarP(&okrGenerateMarpFlags.output, "output", "o", "", "Output file path (default: stdout)")
	okrGenerateMarpCmd.Flags().StringVar(&okrGenerateMarpFlags.theme, "theme", "default", "Slide theme (default, corporate, minimal)")

	// OKR init flags
	okrInitCmd.Flags().StringVar(&okrInitFlags.name, "name", "My OKRs", "Name for the OKR document")
	okrInitCmd.Flags().StringVarP(&okrInitFlags.output, "output", "o", "okrs.json", "Output file path")
	okrInitCmd.Flags().StringVar(&okrInitFlags.period, "period", "", "Period (e.g., 2026-Q1)")

	// Add subcommands
	okrGenerateCmd.AddCommand(okrGenerateMarpCmd)
	okrCmd.AddCommand(okrValidateCmd)
	okrCmd.AddCommand(okrGenerateCmd)
	okrCmd.AddCommand(okrInitCmd)
}

func runOKRValidate(cmd *cobra.Command, args []string) error {
	filepath := args[0]

	// Check file exists
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", filepath)
	}

	// Read and parse OKR
	doc, err := okr.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("reading OKR: %w", err)
	}

	// Set up validation options
	opts := okr.DefaultValidationOptions()
	if okrValidateFlags.strict {
		opts = okr.StrictValidationOptions()
	}

	// Validate
	errs := doc.Validate(opts)

	// Report results
	errors := okr.Errors(errs)
	warnings := okr.Warnings(errs)

	if len(warnings) > 0 {
		fmt.Println("Warnings:")
		for _, w := range warnings {
			fmt.Printf("  - %s\n", w)
		}
		fmt.Println()
	}

	if len(errors) > 0 {
		fmt.Println("Errors:")
		for _, e := range errors {
			fmt.Printf("  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d error(s)", len(errors))
	}

	// Print success info
	fmt.Printf("Valid OKR: %s\n", filepath)
	fmt.Printf("  Objectives: %d\n", len(doc.Objectives))
	fmt.Printf("  Total Key Results: %d\n", len(doc.AllKeyResults()))
	fmt.Printf("  Overall Progress: %.0f%%\n", doc.CalculateOverallProgress()*100)

	if doc.Metadata != nil && doc.Metadata.Name != "" {
		fmt.Printf("  Name: %s\n", doc.Metadata.Name)
	}

	return nil
}

func runOKRMarpGenerate(cmd *cobra.Command, args []string) error {
	inputPath := args[0]

	// Check file exists
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return fmt.Errorf("file not found: %s", inputPath)
	}

	// Read and parse OKR
	doc, err := okr.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("reading OKR: %w", err)
	}

	// Create renderer and options
	renderer := okrmarp.New()
	opts := okrrender.DefaultOptions()
	opts.Theme = okrGenerateMarpFlags.theme

	// Render
	output, err := renderer.Render(doc, opts)
	if err != nil {
		return fmt.Errorf("rendering Marp: %w", err)
	}

	// Write output
	if okrGenerateMarpFlags.output != "" {
		// Ensure output directory exists
		dir := filepath.Dir(okrGenerateMarpFlags.output)
		if dir != "." && dir != "" {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return fmt.Errorf("creating output directory: %w", err)
			}
		}

		if err := os.WriteFile(okrGenerateMarpFlags.output, output, 0600); err != nil {
			return fmt.Errorf("writing output: %w", err)
		}
		fmt.Printf("Generated: %s\n", okrGenerateMarpFlags.output)
	} else {
		// Write to stdout
		fmt.Print(string(output))
	}

	return nil
}

func runOKRInit(cmd *cobra.Command, args []string) error {
	// Check if file already exists
	if _, err := os.Stat(okrInitFlags.output); err == nil {
		return fmt.Errorf("file already exists: %s (use -o to specify a different output path)", okrInitFlags.output)
	}

	// Create template
	now := time.Now()
	period := okrInitFlags.period
	if period == "" {
		quarter := (now.Month()-1)/3 + 1
		period = fmt.Sprintf("%d-Q%d", now.Year(), quarter)
	}

	template := &okr.OKRDocument{
		Schema: "../schema/okr.schema.json",
		Metadata: &okr.Metadata{
			ID:         okr.GenerateID(),
			Name:       okrInitFlags.name,
			Owner:      "Your Name",
			Team:       "Your Team",
			Period:     period,
			PeriodType: "quarter",
			Version:    "1.0.0",
			Status:     okr.StatusDraft,
			CreatedAt:  now,
			UpdatedAt:  now,
		},
		Theme: "Achieve great results this quarter",
		Objectives: []okr.Objective{
			{
				ID:          "obj-1",
				Title:       "First Objective",
				Description: "What do you want to achieve? Make it inspirational.",
				Category:    "Product",
				Owner:       "Your Name",
				Status:      okr.StatusDraft,
				KeyResults: []okr.KeyResult{
					{
						ID:          "kr-1-1",
						Description: "First measurable key result",
						Baseline:    "0",
						Target:      "100",
						Current:     "0",
						Unit:        "%",
						Confidence:  okr.ConfidenceMedium,
					},
					{
						ID:          "kr-1-2",
						Description: "Second measurable key result",
						Baseline:    "0",
						Target:      "50",
						Current:     "0",
						Unit:        "count",
						Confidence:  okr.ConfidenceMedium,
					},
				},
			},
			{
				ID:          "obj-2",
				Title:       "Second Objective",
				Description: "Another goal to achieve this quarter.",
				Category:    "Engineering",
				Owner:       "Your Name",
				Status:      okr.StatusDraft,
				KeyResults: []okr.KeyResult{
					{
						ID:          "kr-2-1",
						Description: "Key result for objective 2",
						Baseline:    "0",
						Target:      "10",
						Current:     "0",
						Unit:        "features",
						Confidence:  okr.ConfidenceHigh,
					},
				},
			},
		},
		Risks: []okr.Risk{
			{
				ID:          "risk-1",
				Description: "What might prevent success?",
				Impact:      "High",
				Likelihood:  "Medium",
				Mitigation:  "How will you address this?",
			},
		},
	}

	// Write to file
	if err := template.WriteFile(okrInitFlags.output); err != nil {
		return fmt.Errorf("writing template: %w", err)
	}

	fmt.Printf("Created: %s\n", okrInitFlags.output)
	fmt.Printf("  Period: %s\n", period)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit the file to add your objectives and key results")
	fmt.Println("  2. Run 'splan goals okr validate " + okrInitFlags.output + "' to check your OKRs")
	fmt.Println("  3. Run 'splan goals okr generate marp " + okrInitFlags.output + " -o slides.md' to create slides")

	return nil
}

// ============================================================================
// Shared Types
// ============================================================================

// Shared generate flags
type generateFlags struct {
	output           string
	margin           string
	mainFont         string
	sansFont         string
	monoFont         string
	fontFamily       string
	noFrontmatter    bool
	noTOC            bool
	noSwimlane       bool
	descLen          int
	swimlaneNoStatus bool
	textIcons        bool
	prdType          string
	sectionOrder     string
}

// ============================================================================
// PRD Commands
// ============================================================================

var prdCmd = &cobra.Command{
	Use:   "prd",
	Short: "Work with Product Requirements Documents",
	Long:  `Commands for generating and validating Product Requirements Documents (PRD).`,
}

var prdGenerateFlags generateFlags

var prdGenerateCmd = &cobra.Command{
	Use:   "generate <input.json>",
	Short: "Convert PRD JSON to markdown",
	Long: `Generate markdown from a Product Requirements Document (PRD).

The output includes YAML frontmatter compatible with Pandoc for PDF generation.
By default, the output file has the same name as the input with a .md extension.

Section ordering:
  --type controls the section ordering template:
    - strategy: Context-first (CurrentState, Problem, Market early)
    - feature: User-needs-first (Problem, Personas, UserStories early)
    - technical: Architecture-focused (TechArchitecture early)

  --order allows custom section ordering (comma-separated IDs).
  Use 'splan requirements prd list-sections' to see available IDs.`,
	Example: `  splan requirements prd generate myproduct.prd.json
  splan requirements prd generate myproduct.json -o output.md
  splan requirements prd generate myproduct.json --no-frontmatter
  splan requirements prd generate myproduct.json --text-icons  # For Pandoc PDF
  splan requirements prd generate myproduct.json --type=strategy
  splan requirements prd generate myproduct.json --order=executiveSummary,problem,solution`,
	Args: cobra.ExactArgs(1),
	RunE: runPRDGenerate,
}

var prdValidateCmd = &cobra.Command{
	Use:     "validate <input.json>",
	Short:   "Validate PRD structure",
	Long:    `Validate a Product Requirements Document by parsing it and checking required fields.`,
	Example: `  splan requirements prd validate myproduct.prd.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runPRDValidate,
}

var prdCheckFlags struct {
	json bool
}

var prdCheckCmd = &cobra.Command{
	Use:   "check <input.json>",
	Short: "Check PRD completeness",
	Long: `Analyze a Product Requirements Document for completeness and quality.

This command evaluates each section of the PRD and provides:
  - Overall completeness score (0-100%) and letter grade
  - Per-section scores for required and optional sections
  - Specific recommendations for improvement

The check examines:
  - Required sections: metadata, executive summary, objectives, personas,
    user stories, requirements, and roadmap
  - Optional sections: assumptions, out of scope, technical architecture,
    UX requirements, risks, and glossary
  - Quality indicators: depth of content, cross-references between sections,
    acceptance criteria coverage, and NFR category coverage`,
	Example: `  splan requirements prd check myproduct.prd.json
  splan requirements prd check myproduct.prd.json --json`,
	Args: cobra.ExactArgs(1),
	RunE: runPRDCheck,
}

var prdScoreFlags struct {
	format string
}

var prdFilterFlags struct {
	output      string
	includeTags []string
	excludeTags []string
	matchAll    bool
}

var prdListSectionsCmd = &cobra.Command{
	Use:   "list-sections",
	Short: "List available section IDs for --order flag",
	Long: `List all available section IDs that can be used with the --order flag.

Each section ID can be used in a comma-separated list with the --order flag
to customize the section ordering in the generated markdown.`,
	Example: `  splan requirements prd list-sections
  splan requirements prd generate myproduct.json --order=executiveSummary,problem,solution`,
	RunE: runPRDListSections,
}

var prdFilterCmd = &cobra.Command{
	Use:   "filter <input.json>",
	Short: "Filter PRD by tags and output filtered JSON",
	Long: `Filter a Product Requirements Document by tags and output the filtered JSON.

By default, uses OR logic (union) - includes entities with ANY of the specified tags.
Use --all to require ALL tags (intersection).

Entities that support tags: personas, user stories, requirements, roadmap phases,
deliverables, OKRs (objectives and key results), and risks.`,
	Example: `  # Include items with tag1 OR tag2 (union)
  splan requirements prd filter input.json --include tag1,tag2

  # Include items with tag1 AND tag2 (intersection)
  splan requirements prd filter input.json --include tag1,tag2 --all

  # Output to specific file
  splan requirements prd filter input.json --include mvp -o filtered.json`,
	Args: cobra.ExactArgs(1),
	RunE: runPRDFilter,
}

var prdScoreCmd = &cobra.Command{
	Use:   "score <input.json>",
	Short: "Score PRD quality with actionable feedback",
	Long: `Score a Product Requirements Document against 10 quality dimensions.

This command provides an actionable workflow:
  (a) Scores - weighted scores across 10 categories
  (b) Gaps - specific missing elements per category
  (c) Fix recommendations - what to improve with effort estimates
  (d) Re-run instructions - command to verify fixes

Output formats:
  - terminal (default): Box-format terminal output with status icons
  - json: Full JSON report for programmatic use
  - markdown: Markdown report for documentation

Quality categories (with weights):
  - Problem Definition (20%)    - Solution Fit (15%)
  - User Understanding (10%)    - Market Awareness (10%)
  - Scope Discipline (10%)      - Requirements Quality (10%)
  - Metrics Quality (10%)       - UX Coverage (5%)
  - Technical Feasibility (5%)  - Risk Management (5%)

Decision thresholds:
  - Approve: >= 8.0
  - Revise:  >= 6.5
  - Reject:  < 3.0 (any blocker)`,
	Example: `  splan requirements prd score myproduct.prd.json
  splan requirements prd score myproduct.prd.json --format=json
  splan requirements prd score myproduct.prd.json --format=markdown`,
	Args: cobra.ExactArgs(1),
	RunE: runPRDScore,
}

func init() {
	// PRD generate flags
	prdGenerateCmd.Flags().StringVarP(&prdGenerateFlags.output, "output", "o", "", "Output markdown file path (default: input with .md extension)")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.margin, "margin", "2cm", "Page margin for Pandoc")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.mainFont, "mainfont", "Helvetica", "Main font family")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.sansFont, "sansfont", "Helvetica", "Sans-serif font family")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.monoFont, "monofont", "Courier New", "Monospace font family")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.fontFamily, "fontfamily", "helvet", "LaTeX font family")
	prdGenerateCmd.Flags().BoolVar(&prdGenerateFlags.noFrontmatter, "no-frontmatter", false, "Disable YAML frontmatter generation")
	prdGenerateCmd.Flags().BoolVar(&prdGenerateFlags.noTOC, "no-toc", false, "Disable Table of Contents generation")
	prdGenerateCmd.Flags().IntVar(&prdGenerateFlags.descLen, "desc-len", prd.DefaultDescriptionMaxLen, "Max length for description fields in tables (0 = no limit)")
	prdGenerateCmd.Flags().BoolVar(&prdGenerateFlags.noSwimlane, "no-swimlane", false, "Disable swimlane table view in roadmap section")
	prdGenerateCmd.Flags().BoolVar(&prdGenerateFlags.swimlaneNoStatus, "swimlane-no-status", false, "Hide status icons in swimlane table")
	prdGenerateCmd.Flags().BoolVar(&prdGenerateFlags.textIcons, "text-icons", false, "Use ASCII text icons instead of emoji for Pandoc/LaTeX PDF compatibility")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.prdType, "type", "", "PRD type for section ordering (strategy, feature, technical)")
	prdGenerateCmd.Flags().StringVar(&prdGenerateFlags.sectionOrder, "order", "", "Custom section order (comma-separated IDs)")

	prdCmd.AddCommand(prdGenerateCmd)
	prdCmd.AddCommand(prdListSectionsCmd)
	prdCmd.AddCommand(prdValidateCmd)
	prdCmd.AddCommand(prdCheckCmd)
	prdCmd.AddCommand(prdScoreCmd)
	prdCmd.AddCommand(prdFilterCmd)

	// PRD check flags
	prdCheckCmd.Flags().BoolVar(&prdCheckFlags.json, "json", false, "Output report as JSON")

	// PRD filter flags
	prdFilterCmd.Flags().StringVarP(&prdFilterFlags.output, "output", "o", "", "Output JSON file path (default: stdout)")
	prdFilterCmd.Flags().StringSliceVarP(&prdFilterFlags.includeTags, "include", "i", nil, "Tags to include (comma-separated)")
	prdFilterCmd.Flags().BoolVarP(&prdFilterFlags.matchAll, "all", "a", false, "Require ALL tags (AND logic) instead of ANY (OR logic)")

	// PRD score flags
	prdScoreCmd.Flags().StringVarP(&prdScoreFlags.format, "format", "f", "terminal", "Output format (terminal, json, markdown)")
}

func runPRDGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Determine output file
	output := prdGenerateFlags.output
	if output == "" {
		output = deriveOutputPath(inputFile)
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Apply PRD type flag (overrides JSON value)
	if prdGenerateFlags.prdType != "" {
		prdType := prd.PRDType(prdGenerateFlags.prdType)
		switch prdType {
		case prd.PRDTypeStrategy, prd.PRDTypeFeature, prd.PRDTypeTechnical:
			doc.Metadata.PRDType = prdType
		default:
			return fmt.Errorf("invalid --type value: %s (expected strategy, feature, or technical)", prdGenerateFlags.prdType)
		}
	}

	// Apply custom section order flag (overrides JSON value)
	if prdGenerateFlags.sectionOrder != "" {
		orderIDs := strings.Split(prdGenerateFlags.sectionOrder, ",")
		for i := range orderIDs {
			orderIDs[i] = strings.TrimSpace(orderIDs[i])
		}
		// Validate section IDs
		invalid := prd.ValidateSectionOrder(orderIDs)
		if len(invalid) > 0 {
			return fmt.Errorf("invalid section IDs in --order: %s (use 'splan requirements prd list-sections' to see valid IDs)", strings.Join(invalid, ", "))
		}
		doc.Metadata.SectionOrder = orderIDs
	}

	// Handle TOC option (default: enabled, disabled with --no-toc)
	includeTOC := !prdGenerateFlags.noTOC
	// Handle swimlane option (default: enabled, disabled with --no-swimlane)
	includeSwimlane := !prdGenerateFlags.noSwimlane

	opts := prd.MarkdownOptions{
		IncludeFrontmatter:   !prdGenerateFlags.noFrontmatter,
		Margin:               prdGenerateFlags.margin,
		MainFont:             prdGenerateFlags.mainFont,
		SansFont:             prdGenerateFlags.sansFont,
		MonoFont:             prdGenerateFlags.monoFont,
		FontFamily:           prdGenerateFlags.fontFamily,
		DescriptionMaxLen:    prdGenerateFlags.descLen,
		IncludeSwimlaneTable: includeSwimlane,
		IncludeTOC:           &includeTOC,
		UseTextIcons:         prdGenerateFlags.textIcons,
	}

	// Configure swimlane table options
	if includeSwimlane {
		tableOpts := prd.DefaultRoadmapTableOptions()
		tableOpts.IncludeStatus = !prdGenerateFlags.swimlaneNoStatus
		tableOpts.UseTextIcons = prdGenerateFlags.textIcons
		opts.RoadmapTableOptions = &tableOpts
	}

	markdown := doc.ToMarkdown(opts)

	if err := os.WriteFile(output, []byte(markdown), 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

func runPRDValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required (at least one author)")
	}
	if doc.ExecutiveSummary.ProblemStatement == "" {
		errors = append(errors, "executive_summary.problem_statement is required")
	}
	if doc.ExecutiveSummary.ProposedSolution == "" {
		errors = append(errors, "executive_summary.proposed_solution is required")
	}
	if len(doc.Personas) == 0 {
		errors = append(errors, "personas is required (at least one persona)")
	}
	if len(doc.UserStories) == 0 {
		errors = append(errors, "user_stories is required (at least one user story)")
	}
	if len(doc.Roadmap.Phases) == 0 {
		errors = append(errors, "roadmap.phases is required (at least one phase)")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Printf("Valid PRD: %s\n", inputFile)
	fmt.Printf("  Title: %s\n", doc.Metadata.Title)
	fmt.Printf("  Version: %s\n", doc.Metadata.Version)
	fmt.Printf("  Personas: %d\n", len(doc.Personas))
	fmt.Printf("  User Stories: %d\n", len(doc.UserStories))
	fmt.Printf("  Functional Requirements: %d\n", len(doc.Requirements.Functional))
	fmt.Printf("  Non-Functional Requirements: %d\n", len(doc.Requirements.NonFunctional))
	fmt.Printf("  Phases: %d\n", len(doc.Roadmap.Phases))

	return nil
}

func runPRDCheck(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	report := doc.CheckCompleteness()

	if prdCheckFlags.json {
		output, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling report: %w", err)
		}
		fmt.Println(string(output))
	} else {
		fmt.Print(report.FormatReport())
	}

	// Return non-zero exit code if PRD has critical issues
	if report.Grade == "F" {
		return fmt.Errorf("PRD completeness check failed (Grade: F)")
	}

	return nil
}

func runPRDScore(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Generate evaluation report from deterministic scoring
	report := prd.ScoreToEvaluationReport(&doc, inputFile)

	switch strings.ToLower(prdScoreFlags.format) {
	case "json":
		output, err := json.MarshalIndent(report, "", "  ")
		if err != nil {
			return fmt.Errorf("marshaling report: %w", err)
		}
		fmt.Println(string(output))

	case "markdown":
		fmt.Print(formatEvaluationReportMarkdown(report))

	case "terminal", "":
		renderer := terminal.New(os.Stdout)
		if err := renderer.Render(report); err != nil {
			return fmt.Errorf("rendering report: %w", err)
		}

	default:
		return fmt.Errorf("unknown format: %s (expected terminal, json, or markdown)", prdScoreFlags.format)
	}

	// Return non-zero exit code if PRD has blocking issues
	if !report.Decision.Passed {
		return fmt.Errorf("PRD evaluation: %s", report.Decision.Rationale)
	}

	return nil
}

func runPRDFilter(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Validate filter tags
	if len(prdFilterFlags.includeTags) > 0 {
		for _, tag := range prdFilterFlags.includeTags {
			if err := prd.ValidateTag(tag); err != nil {
				return fmt.Errorf("invalid filter tag: %w", err)
			}
		}
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc prd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	// Apply filter
	var filtered prd.Document
	if len(prdFilterFlags.includeTags) > 0 {
		if prdFilterFlags.matchAll {
			filtered = doc.FilterByTagsAll(prdFilterFlags.includeTags...)
		} else {
			filtered = doc.FilterByTags(prdFilterFlags.includeTags...)
		}
	} else {
		filtered = doc
	}

	// Marshal output
	output, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling filtered JSON: %w", err)
	}

	// Write output
	if prdFilterFlags.output != "" {
		if err := os.WriteFile(prdFilterFlags.output, output, 0600); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		fmt.Printf("Filtered PRD written to: %s\n", prdFilterFlags.output)
	} else {
		fmt.Println(string(output))
	}

	return nil
}

func runPRDListSections(cmd *cobra.Command, args []string) error {
	sections := prd.ListSections()

	fmt.Println("Available PRD section IDs for --order flag:")
	fmt.Println()
	fmt.Println("PRD Types (--type flag):")
	fmt.Println("  strategy  - Context-first ordering (CurrentState, Problem, Market early)")
	fmt.Println("  feature   - User-needs-first ordering (Problem, Personas, UserStories early)")
	fmt.Println("  technical - Architecture-focused ordering (TechArchitecture early)")
	fmt.Println()
	fmt.Println("Section IDs:")
	fmt.Println("  ID                        Display Name")
	fmt.Println("  --------------------------  --------------------------------")
	for _, s := range sections {
		if s.DisplayName != "" {
			fmt.Printf("  %-26s  %s\n", s.ID, s.DisplayName)
		}
	}
	fmt.Println()
	fmt.Println("Example:")
	fmt.Println("  splan requirements prd generate myproduct.json --order=executiveSummary,problem,solution,objectives")

	return nil
}

func formatEvaluationReportMarkdown(report *evaluation.EvaluationReport) string {
	var b strings.Builder

	b.WriteString("# PRD Evaluation Report\n\n")
	b.WriteString(fmt.Sprintf("**Document**: %s\n", report.Metadata.Document))
	if report.Metadata.DocumentTitle != "" {
		b.WriteString(fmt.Sprintf("**Title**: %s\n", report.Metadata.DocumentTitle))
	}
	b.WriteString(fmt.Sprintf("**Score**: %.1f / 10.0\n", report.WeightedScore))
	b.WriteString(fmt.Sprintf("**Decision**: %s\n\n", strings.ToUpper(string(report.Decision.Status))))

	// Category scores table
	b.WriteString("## Category Scores\n\n")
	b.WriteString("| Category | Score | Weight | Status |\n")
	b.WriteString("|----------|-------|--------|--------|\n")
	for _, cs := range report.Categories {
		status := "✅"
		if cs.Status == evaluation.ScoreStatusWarn {
			status = "⚠️"
		} else if cs.Status == evaluation.ScoreStatusFail {
			status = "❌"
		}
		b.WriteString(fmt.Sprintf("| %s | %.1f | %.0f%% | %s |\n",
			cs.Category, cs.Score, cs.Weight*100, status))
	}
	b.WriteString("\n")

	// Findings by severity
	if len(report.Findings) > 0 {
		b.WriteString("## Findings\n\n")
		for _, sev := range evaluation.AllSeverities() {
			sevFindings := []evaluation.Finding{}
			for _, f := range report.Findings {
				if f.Severity == sev {
					sevFindings = append(sevFindings, f)
				}
			}
			if len(sevFindings) > 0 {
				b.WriteString(fmt.Sprintf("### %s %s\n\n", sev.Icon(), strings.ToUpper(string(sev))))
				for _, f := range sevFindings {
					b.WriteString(fmt.Sprintf("- **[%s]** %s\n", f.Category, f.Title))
					if f.Recommendation != "" {
						b.WriteString(fmt.Sprintf("  - Fix: %s\n", f.Recommendation))
					}
					if f.Owner != "" {
						b.WriteString(fmt.Sprintf("  - Owner: %s\n", f.Owner))
					}
				}
				b.WriteString("\n")
			}
		}
	}

	// Next steps
	b.WriteString("## Next Steps\n\n")
	if len(report.NextSteps.Immediate) > 0 {
		b.WriteString("### Immediate Actions\n\n")
		for _, action := range report.NextSteps.Immediate {
			b.WriteString(fmt.Sprintf("- [ ] 🔴 %s\n", action.Action))
		}
		b.WriteString("\n")
	}

	if len(report.NextSteps.Recommended) > 0 {
		b.WriteString("### Recommended\n\n")
		for _, action := range report.NextSteps.Recommended {
			b.WriteString(fmt.Sprintf("- [ ] %s\n", action.Action))
		}
		b.WriteString("\n")
	}

	b.WriteString("### Re-run Evaluation\n\n")
	b.WriteString(fmt.Sprintf("```bash\n%s\n```\n", report.NextSteps.RerunCommand))

	return b.String()
}

// ============================================================================
// MRD Commands
// ============================================================================

var mrdCmd = &cobra.Command{
	Use:   "mrd",
	Short: "Work with Market Requirements Documents",
	Long:  `Commands for generating and validating Market Requirements Documents (MRD).`,
}

var mrdGenerateFlags generateFlags

var mrdGenerateCmd = &cobra.Command{
	Use:   "generate <input.json>",
	Short: "Convert MRD JSON to markdown",
	Long: `Generate markdown from a Market Requirements Document (MRD).

The output includes YAML frontmatter compatible with Pandoc for PDF generation.
By default, the output file has the same name as the input with a .md extension.`,
	Example: `  splan requirements mrd generate market-analysis.mrd.json
  splan requirements mrd generate market.json -o output.md
  splan requirements mrd generate market.json --no-frontmatter`,
	Args: cobra.ExactArgs(1),
	RunE: runMRDGenerate,
}

var mrdValidateCmd = &cobra.Command{
	Use:     "validate <input.json>",
	Short:   "Validate MRD structure",
	Long:    `Validate a Market Requirements Document by parsing it and checking required fields.`,
	Example: `  splan requirements mrd validate market-analysis.mrd.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runMRDValidate,
}

func init() {
	// MRD generate flags
	mrdGenerateCmd.Flags().StringVarP(&mrdGenerateFlags.output, "output", "o", "", "Output markdown file path (default: input with .md extension)")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.margin, "margin", "2cm", "Page margin for Pandoc")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.mainFont, "mainfont", "Helvetica", "Main font family")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.sansFont, "sansfont", "Helvetica", "Sans-serif font family")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.monoFont, "monofont", "Courier New", "Monospace font family")
	mrdGenerateCmd.Flags().StringVar(&mrdGenerateFlags.fontFamily, "fontfamily", "helvet", "LaTeX font family")
	mrdGenerateCmd.Flags().BoolVar(&mrdGenerateFlags.noFrontmatter, "no-frontmatter", false, "Disable YAML frontmatter generation")

	mrdCmd.AddCommand(mrdGenerateCmd)
	mrdCmd.AddCommand(mrdValidateCmd)
}

func runMRDGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Determine output file
	output := mrdGenerateFlags.output
	if output == "" {
		output = deriveOutputPath(inputFile)
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc mrd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	opts := mrd.MarkdownOptions{
		IncludeFrontmatter: !mrdGenerateFlags.noFrontmatter,
		Margin:             mrdGenerateFlags.margin,
		MainFont:           mrdGenerateFlags.mainFont,
		SansFont:           mrdGenerateFlags.sansFont,
		MonoFont:           mrdGenerateFlags.monoFont,
		FontFamily:         mrdGenerateFlags.fontFamily,
	}
	markdown := doc.ToMarkdown(opts)

	if err := os.WriteFile(output, []byte(markdown), 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

func runMRDValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc mrd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required (at least one author)")
	}
	if doc.ExecutiveSummary.MarketOpportunity == "" {
		errors = append(errors, "executive_summary.market_opportunity is required")
	}
	if doc.ExecutiveSummary.ProposedOffering == "" {
		errors = append(errors, "executive_summary.proposed_offering is required")
	}
	if doc.MarketOverview.TAM.Value == "" {
		errors = append(errors, "market_overview.tam.value is required")
	}
	if len(doc.TargetMarket.PrimarySegments) == 0 {
		errors = append(errors, "target_market.primary_segments is required (at least one segment)")
	}
	if len(doc.CompetitiveLandscape.Competitors) == 0 {
		errors = append(errors, "competitive_landscape.competitors is required (at least one competitor)")
	}
	if len(doc.MarketRequirements) == 0 {
		errors = append(errors, "market_requirements is required (at least one requirement)")
	}
	if doc.Positioning.Statement == "" {
		errors = append(errors, "positioning.statement is required")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Printf("Valid MRD: %s\n", inputFile)
	fmt.Printf("  Title: %s\n", doc.Metadata.Title)
	fmt.Printf("  Version: %s\n", doc.Metadata.Version)
	fmt.Printf("  TAM: %s\n", doc.MarketOverview.TAM.Value)
	fmt.Printf("  Primary Segments: %d\n", len(doc.TargetMarket.PrimarySegments))
	fmt.Printf("  Buyer Personas: %d\n", len(doc.TargetMarket.BuyerPersonas))
	fmt.Printf("  Competitors: %d\n", len(doc.CompetitiveLandscape.Competitors))
	fmt.Printf("  Market Requirements: %d\n", len(doc.MarketRequirements))
	fmt.Printf("  Success Metrics: %d\n", len(doc.SuccessMetrics))

	return nil
}

// ============================================================================
// TRD Commands
// ============================================================================

var trdCmd = &cobra.Command{
	Use:   "trd",
	Short: "Work with Technical Requirements Documents",
	Long:  `Commands for generating and validating Technical Requirements Documents (TRD).`,
}

var trdGenerateFlags generateFlags

var trdGenerateCmd = &cobra.Command{
	Use:   "generate <input.json>",
	Short: "Convert TRD JSON to markdown",
	Long: `Generate markdown from a Technical Requirements Document (TRD).

The output includes YAML frontmatter compatible with Pandoc for PDF generation.
By default, the output file has the same name as the input with a .md extension.`,
	Example: `  splan requirements trd generate architecture.trd.json
  splan requirements trd generate tech-spec.json -o output.md
  splan requirements trd generate tech-spec.json --no-frontmatter`,
	Args: cobra.ExactArgs(1),
	RunE: runTRDGenerate,
}

var trdValidateCmd = &cobra.Command{
	Use:     "validate <input.json>",
	Short:   "Validate TRD structure",
	Long:    `Validate a Technical Requirements Document by parsing it and checking required fields.`,
	Example: `  splan requirements trd validate architecture.trd.json`,
	Args:    cobra.ExactArgs(1),
	RunE:    runTRDValidate,
}

func init() {
	// TRD generate flags
	trdGenerateCmd.Flags().StringVarP(&trdGenerateFlags.output, "output", "o", "", "Output markdown file path (default: input with .md extension)")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.margin, "margin", "2cm", "Page margin for Pandoc")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.mainFont, "mainfont", "Helvetica", "Main font family")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.sansFont, "sansfont", "Helvetica", "Sans-serif font family")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.monoFont, "monofont", "Courier New", "Monospace font family")
	trdGenerateCmd.Flags().StringVar(&trdGenerateFlags.fontFamily, "fontfamily", "helvet", "LaTeX font family")
	trdGenerateCmd.Flags().BoolVar(&trdGenerateFlags.noFrontmatter, "no-frontmatter", false, "Disable YAML frontmatter generation")

	trdCmd.AddCommand(trdGenerateCmd)
	trdCmd.AddCommand(trdValidateCmd)
}

func runTRDGenerate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Determine output file
	output := trdGenerateFlags.output
	if output == "" {
		output = deriveOutputPath(inputFile)
	}

	// Read input file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc trd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	opts := trd.MarkdownOptions{
		IncludeFrontmatter: !trdGenerateFlags.noFrontmatter,
		Margin:             trdGenerateFlags.margin,
		MainFont:           trdGenerateFlags.mainFont,
		SansFont:           trdGenerateFlags.sansFont,
		MonoFont:           trdGenerateFlags.monoFont,
		FontFamily:         trdGenerateFlags.fontFamily,
	}
	markdown := doc.ToMarkdown(opts)

	if err := os.WriteFile(output, []byte(markdown), 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	fmt.Printf("Generated: %s\n", output)
	return nil
}

func runTRDValidate(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("reading input file: %w", err)
	}

	var doc trd.Document
	if err := json.Unmarshal(data, &doc); err != nil {
		return fmt.Errorf("parsing JSON: %w", err)
	}

	var errors []string

	if doc.Metadata.ID == "" {
		errors = append(errors, "metadata.id is required")
	}
	if doc.Metadata.Title == "" {
		errors = append(errors, "metadata.title is required")
	}
	if doc.Metadata.Version == "" {
		errors = append(errors, "metadata.version is required")
	}
	if len(doc.Metadata.Authors) == 0 {
		errors = append(errors, "metadata.authors is required (at least one author)")
	}
	if doc.ExecutiveSummary.Purpose == "" {
		errors = append(errors, "executive_summary.purpose is required")
	}
	if doc.ExecutiveSummary.Scope == "" {
		errors = append(errors, "executive_summary.scope is required")
	}
	if doc.Architecture.Overview == "" {
		errors = append(errors, "architecture.overview is required")
	}
	if len(doc.Architecture.Components) == 0 {
		errors = append(errors, "architecture.components is required (at least one component)")
	}
	if doc.SecurityDesign.Overview == "" {
		errors = append(errors, "security_design.overview is required")
	}
	if len(doc.Performance.Requirements) == 0 {
		errors = append(errors, "performance.requirements is required (at least one requirement)")
	}
	if len(doc.Deployment.Environments) == 0 {
		errors = append(errors, "deployment.environments is required (at least one environment)")
	}

	if len(errors) > 0 {
		fmt.Fprintf(os.Stderr, "Validation failed for %s:\n", inputFile)
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "  - %s\n", e)
		}
		return fmt.Errorf("validation failed with %d errors", len(errors))
	}

	fmt.Printf("Valid TRD: %s\n", inputFile)
	fmt.Printf("  Title: %s\n", doc.Metadata.Title)
	fmt.Printf("  Version: %s\n", doc.Metadata.Version)
	fmt.Printf("  Components: %d\n", len(doc.Architecture.Components))
	fmt.Printf("  APIs: %d\n", len(doc.APISpecifications))
	fmt.Printf("  Performance Requirements: %d\n", len(doc.Performance.Requirements))
	fmt.Printf("  Environments: %d\n", len(doc.Deployment.Environments))
	fmt.Printf("  Integrations: %d\n", len(doc.Integration))

	return nil
}

// ============================================================================
// Schema Commands
// ============================================================================

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate JSON Schema files from Go types",
	Long: `Commands for generating JSON Schema files from Go type definitions.

This implements the Go-first approach where Go structs are the source of truth
and JSON Schema is generated from them.`,
}

var schemaGenerateFlags struct {
	output  string
	docType string
}

var schemaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate JSON Schema from Go types",
	Long: `Generate JSON Schema files from Go type definitions.

By default, generates all schema files (PRD, OKR, V2MOM) to the schema/ directory.
Use --type to generate a specific document type's schema.`,
	Example: `  splan schema generate
  splan schema generate -o ./schema/
  splan schema generate --type prd -o prd.schema.json
  splan schema generate --type okr -o okr.schema.json
  splan schema generate --type v2mom -o v2mom.schema.json`,
	RunE: runSchemaGenerate,
}

func init() {
	schemaGenerateCmd.Flags().StringVarP(&schemaGenerateFlags.output, "output", "o", ".", "Output directory or file path")
	schemaGenerateCmd.Flags().StringVarP(&schemaGenerateFlags.docType, "type", "t", "all", "Document type to generate (prd, okr, v2mom, mrd, trd, or all)")

	schemaCmd.AddCommand(schemaGenerateCmd)
}

func runSchemaGenerate(cmd *cobra.Command, args []string) error {
	gen := schema.NewGenerator()
	output := schemaGenerateFlags.output
	docType := strings.ToLower(schemaGenerateFlags.docType)

	switch docType {
	case "prd":
		// Single PRD schema
		path := output
		if isDir(output) {
			path = filepath.Join(output, "prd.schema.json")
		}
		if err := gen.WritePRDSchema(path); err != nil {
			return fmt.Errorf("generating PRD schema: %w", err)
		}
		fmt.Printf("Generated: %s\n", path)

	case "all":
		// All schemas to directory
		dir := output
		if !isDir(dir) {
			dir = filepath.Dir(output)
		}
		if err := gen.GenerateAll(dir); err != nil {
			return fmt.Errorf("generating schemas: %w", err)
		}
		fmt.Printf("Generated schemas in: %s\n", dir)

	case "okr":
		// Single OKR schema
		path := output
		if isDir(output) {
			path = filepath.Join(output, "okr.schema.json")
		}
		if err := gen.WriteOKRSchema(path); err != nil {
			return fmt.Errorf("generating OKR schema: %w", err)
		}
		fmt.Printf("Generated: %s\n", path)

	case "v2mom":
		// Single V2MOM schema
		path := output
		if isDir(output) {
			path = filepath.Join(output, "v2mom.schema.json")
		}
		if err := gen.WriteV2MOMSchema(path); err != nil {
			return fmt.Errorf("generating V2MOM schema: %w", err)
		}
		fmt.Printf("Generated: %s\n", path)

	case "mrd", "trd":
		return fmt.Errorf("schema generation for %s is not yet implemented", docType)

	default:
		return fmt.Errorf("unknown document type: %s (expected prd, okr, v2mom, mrd, trd, or all)", docType)
	}

	return nil
}

// ============================================================================
// Utility Functions
// ============================================================================

func deriveOutputPath(inputFile string) string {
	ext := filepath.Ext(inputFile)
	if ext == ".json" {
		base := strings.TrimSuffix(inputFile, ext)
		return base + ".md"
	}
	return inputFile + ".md"
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}
