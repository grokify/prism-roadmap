// Package cli provides the exported Cobra command tree for the PRISM roadmap CLI.
package cli

import (
	"github.com/spf13/cobra"
)

// Set by GoReleaser ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// RootCmd is the root command for PRISM roadmap operations.
// It can be imported and added as a subcommand to other CLI tools.
var RootCmd = &cobra.Command{
	Use:   "roadmap",
	Short: "PRISM Roadmap - Goals, requirements, and planning documents",
	Long: `PRISM Roadmap provides tools for working with roadmaps and planning documents.

Document types:
  Goals (define what the roadmap achieves):
    - OKR (Objectives and Key Results)
    - V2MOM (Vision, Values, Methods, Obstacles, Measures)

  Requirements (execute roadmap items):
    - PRD (Product Requirements Document)
    - MRD (Market Requirements Document)
    - TRD (Technical Requirements Document)

It can convert document JSON files to markdown with Pandoc-compatible YAML
frontmatter, generate Marp presentations, and validate files against their
respective schemas.`,
	Version: version,
}

func init() {
	RootCmd.SetVersionTemplate("splan version {{.Version}} (commit: " + commit + ", built: " + date + ")\n")

	// Add top-level commands
	RootCmd.AddCommand(requirementsCmd)
	RootCmd.AddCommand(goalsCmd)
	RootCmd.AddCommand(journeyCmd)
	RootCmd.AddCommand(schemaCmd)
	RootCmd.AddCommand(mergeCmd)
	RootCmd.AddCommand(rmiCmd)

	// Add requirements subcommands
	requirementsCmd.AddCommand(prdCmd)
	requirementsCmd.AddCommand(mrdCmd)
	requirementsCmd.AddCommand(trdCmd)

	// Add goals subcommands
	goalsCmd.AddCommand(v2momCmd)
	goalsCmd.AddCommand(okrCmd)
	goalsCmd.AddCommand(goalsCascadeCmd)
	goalsCmd.AddCommand(goalsAlignCmd)
	goalsCmd.AddCommand(goalsValidateCmd)
	goalsCmd.AddCommand(goalsInitCmd)

	// Add journey subcommands
	journeyCmd.AddCommand(journeyValidateCmd)
	journeyCmd.AddCommand(journeyGenerateCmd)
	journeyCmd.AddCommand(journeyCheckCmd)
	journeyCmd.AddCommand(journeyInitCmd)
	journeyGenerateCmd.AddCommand(journeyGenerateMarkdownCmd)
}
