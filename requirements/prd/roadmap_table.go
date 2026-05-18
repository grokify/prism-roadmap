package prd

import (
	"fmt"
	"slices"
	"sort"
	"strings"

	"github.com/grokify/prism-execution/roadmap"
)

// RoadmapTableOptions configures roadmap table generation.
// This is an alias for backward compatibility.
type RoadmapTableOptions = roadmap.TableOptions

// DefaultRoadmapTableOptions returns sensible defaults for roadmap table generation.
func DefaultRoadmapTableOptions() RoadmapTableOptions {
	return roadmap.DefaultTableOptions()
}

// StatusLegend returns a markdown table explaining the status icons.
func StatusLegend() string {
	return roadmap.StatusLegend()
}

// StatusLegendWithOptions returns a markdown table explaining the status icons.
// If useText is true, shows ASCII text icons instead of emoji.
func StatusLegendWithOptions(useText bool) string {
	return roadmap.StatusLegendWithOptions(useText)
}

// ToSwimlaneTableWithOKRs generates a markdown table with phases as columns,
// deliverable types as swimlane rows, and OKR swimlanes auto-derived from PhaseTargets.
// Deprecated: Use ToSwimlaneTableWithGoals for framework-agnostic goal support.
func (d *Document) ToSwimlaneTableWithOKRs(opts RoadmapTableOptions) string {
	return d.ToSwimlaneTableWithGoals(opts)
}

// ToSwimlaneTableWithGoals generates a markdown table with phases as columns,
// deliverable types as swimlane rows, and goal swimlanes auto-derived from PhaseTargets.
// This method supports both OKR and V2MOM frameworks via the Goals abstraction.
//
// Labels are dynamic based on the framework:
//   - OKR: "Objectives" and "Key Results"
//   - V2MOM: "Methods" and "Measures"
//
// Example output:
//
//	| Swimlane       | **Phase 1**<br>Q1 2026 | **Phase 2**<br>Q2 2026 |
//	|----------------|------------------------|------------------------|
//	| Features       | • Auth<br>• Dashboard  | • Reporting            |
//	| Objectives     | • O1: Market leader    |                        |
//	| Key Results    | • KR1.1: Share → 15%   | • KR1.1: Share → 20%   |
func (d *Document) ToSwimlaneTableWithGoals(opts RoadmapTableOptions) string {
	if len(d.Roadmap.Phases) == 0 {
		return ""
	}

	// Get goals using the framework-agnostic method
	productGoals := d.GetProductGoals()

	// Build phase ID to index map
	phaseIDToIndex := make(map[string]int)
	for i, phase := range d.Roadmap.Phases {
		phaseIDToIndex[phase.ID] = i
	}

	// Collect all unique swimlanes (deliverable types) across all phases
	swimlaneSet := make(map[DeliverableType]bool)
	for _, phase := range d.Roadmap.Phases {
		for _, del := range phase.Deliverables {
			swimlaneSet[del.Type] = true
		}
	}

	// Check if we have goals with PhaseTargets
	hasGoalsWithPhaseTargets := false
	if productGoals != nil {
		resultsByPhase := productGoals.ResultItemsByPhase()
		hasGoalsWithPhaseTargets = len(resultsByPhase) > 0
	}

	// Determine swimlane order
	var swimlanes []DeliverableType
	if len(opts.SwimlaneOrder) > 0 {
		for _, st := range opts.SwimlaneOrder {
			if swimlaneSet[st] || opts.IncludeEmptySwimlanes {
				swimlanes = append(swimlanes, st)
			}
		}
		// Add any swimlanes not in the specified order
		for st := range swimlaneSet {
			if !slices.Contains(opts.SwimlaneOrder, st) {
				swimlanes = append(swimlanes, st)
			}
		}
	} else {
		for st := range swimlaneSet {
			swimlanes = append(swimlanes, st)
		}
		sort.Slice(swimlanes, func(i, j int) bool {
			return string(swimlanes[i]) < string(swimlanes[j])
		})
	}

	var sb strings.Builder

	// Header row: | Swimlane | **Phase 1**<br>Description | **Phase 2**<br>Description | ...
	sb.WriteString("| Swimlane |")
	for i, phase := range d.Roadmap.Phases {
		header := fmt.Sprintf(" **Phase %d**<br>%s |", i+1, phase.Name)
		sb.WriteString(header)
	}
	sb.WriteString("\n")

	// Separator row
	sb.WriteString("|----------|")
	for range d.Roadmap.Phases {
		sb.WriteString("----------|")
	}
	sb.WriteString("\n")

	// Data rows: one per deliverable swimlane
	for _, swimlane := range swimlanes {
		sb.WriteString(fmt.Sprintf("| **%s** |", roadmap.SwimlaneLabel(swimlane)))

		for _, phase := range d.Roadmap.Phases {
			var items []string
			for _, del := range phase.Deliverables {
				if del.Type == swimlane {
					item := del.Title
					if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
						item = item[:opts.MaxTitleLen-3] + "..."
					}
					if opts.IncludeStatus && del.Status != "" {
						item = fmt.Sprintf("%s %s", roadmap.StatusIconWithOptions(del.Status, opts.UseTextIcons), item)
					}
					items = append(items, "• "+item)
				}
			}
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	// Add goal swimlanes if we have goals with PhaseTargets
	if hasGoalsWithPhaseTargets && productGoals != nil {
		// Get dynamic labels based on framework
		goalLabel := productGoals.GoalLabel()     // "Objectives" or "Methods"
		resultLabel := productGoals.ResultLabel() // "Key Results" or "Measures"

		// Goals swimlane: show goal in every phase where any of its results has a target
		sb.WriteString(fmt.Sprintf("| **%s** |", goalLabel))
		goalsByPhase := make(map[int][]string) // phase index -> goal titles
		goalItems := productGoals.GoalItems()
		resultItems := productGoals.ResultItems()

		// Build a map of goalID -> goal title
		goalTitles := make(map[string]string)
		for i, g := range goalItems {
			goalTitles[g.ID] = fmt.Sprintf("G%d: %s", i+1, g.Title)
		}

		// Find phases for each goal based on its results
		goalPhases := make(map[string]map[int]bool) // goalID -> set of phase indices
		for _, r := range resultItems {
			if r.PhaseID != "" {
				if idx, ok := phaseIDToIndex[r.PhaseID]; ok {
					if goalPhases[r.GoalID] == nil {
						goalPhases[r.GoalID] = make(map[int]bool)
					}
					goalPhases[r.GoalID][idx] = true
				}
			}
		}

		// Add goal labels to phases
		for goalID, phases := range goalPhases {
			if title, ok := goalTitles[goalID]; ok {
				for phaseIdx := range phases {
					goalsByPhase[phaseIdx] = append(goalsByPhase[phaseIdx], "• "+title)
				}
			}
		}

		for i := range d.Roadmap.Phases {
			items := goalsByPhase[i]
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")

		// Results swimlane: show result targets for each phase
		sb.WriteString(fmt.Sprintf("| **%s** |", resultLabel))
		resultsByPhase := productGoals.ResultItemsByPhase()
		for _, phase := range d.Roadmap.Phases {
			var items []string
			if results, ok := resultsByPhase[phase.ID]; ok {
				for _, r := range results {
					// Format: R1: Title → Target
					label := fmt.Sprintf("%s → %s", r.Title, r.PhaseTarget)
					if opts.IncludeStatus && r.Status != "" {
						label = fmt.Sprintf("%s %s", roadmap.PhaseTargetStatusIconWithOptions(r.Status, opts.UseTextIcons), label)
					}
					items = append(items, "• "+label)
				}
			}
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
