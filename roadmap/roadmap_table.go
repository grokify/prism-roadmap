package roadmap

import (
	"fmt"
	"slices"
	"sort"
	"strings"
)

// TableOptions configures roadmap table generation.
type TableOptions struct {
	// IncludeStatus adds status indicators to deliverables
	IncludeStatus bool
	// IncludeEmptySwimlanes shows rows even if no deliverables of that type exist
	IncludeEmptySwimlanes bool
	// SwimlaneOrder specifies the order of swimlanes (nil = alphabetical)
	SwimlaneOrder []DeliverableType
	// MaxTitleLen truncates deliverable titles (0 = no limit)
	MaxTitleLen int
	// IncludeOKRs adds Objectives and Key Results swimlanes derived from PhaseTargets.
	// This is primarily used by PRD documents that have OKR integration.
	IncludeOKRs bool
	// UseTextIcons uses ASCII text instead of emoji for status icons.
	// Enable this for Pandoc/LaTeX PDF generation compatibility.
	UseTextIcons bool
}

// DefaultTableOptions returns sensible defaults for roadmap table generation.
func DefaultTableOptions() TableOptions {
	return TableOptions{
		IncludeStatus:         true,
		IncludeEmptySwimlanes: false,
		SwimlaneOrder: []DeliverableType{
			DeliverableFeature,
			DeliverableIntegration,
			DeliverableInfrastructure,
			DeliverableDocumentation,
			DeliverableMilestone,
			DeliverableRollout,
		},
		MaxTitleLen: 0, // No truncation by default
	}
}

// ToSwimlaneTable generates a markdown table with phases as columns and
// deliverable types as swimlane rows.
//
// Example output:
//
//	| Swimlane       | **Phase 1**<br>Foundation | **Phase 2**<br>Core Features |
//	|----------------|---------------------------|------------------------------|
//	| Features       | • Auth<br>• Search        | • Dashboard                  |
//	| Infrastructure | • CI/CD                   | • Monitoring                 |
func (r *Roadmap) ToSwimlaneTable(opts TableOptions) string {
	if len(r.Phases) == 0 {
		return ""
	}

	// Collect all unique swimlanes (deliverable types) across all phases
	swimlaneSet := make(map[DeliverableType]bool)
	for _, phase := range r.Phases {
		for _, del := range phase.Deliverables {
			swimlaneSet[del.Type] = true
		}
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

	if len(swimlanes) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header row: | Swimlane | **Phase 1**<br>Description | **Phase 2**<br>Description | ...
	sb.WriteString("| Swimlane |")
	for i, phase := range r.Phases {
		// Format: **Phase N**<br>Phase Name/Description
		header := fmt.Sprintf(" **Phase %d**<br>%s |", i+1, phase.Name)
		sb.WriteString(header)
	}
	sb.WriteString("\n")

	// Separator row
	sb.WriteString("|----------|")
	for range r.Phases {
		sb.WriteString("----------|")
	}
	sb.WriteString("\n")

	// Data rows: one per swimlane
	for _, swimlane := range swimlanes {
		sb.WriteString(fmt.Sprintf("| **%s** |", SwimlaneLabel(swimlane)))

		for _, phase := range r.Phases {
			// Collect deliverables of this type in this phase
			var items []string
			for _, del := range phase.Deliverables {
				if del.Type == swimlane {
					item := del.Title
					if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
						item = item[:opts.MaxTitleLen-3] + "..."
					}
					if opts.IncludeStatus && del.Status != "" {
						item = fmt.Sprintf("%s %s", StatusIconWithOptions(del.Status, opts.UseTextIcons), item)
					}
					// Add bullet point prefix
					items = append(items, "• "+item)
				}
			}
			// Join with <br> for line breaks within cell
			cell := strings.Join(items, "<br>")
			sb.WriteString(fmt.Sprintf(" %s |", cell))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// ToPhaseTable generates a traditional phase-based table showing
// each phase with its deliverables listed.
//
// Example output:
//
//	| Phase   | Status      | Deliverables                          |
//	|---------|-------------|---------------------------------------|
//	| Phase 1 | In Progress | • ✅ Auth<br>• 🔄 Search<br>• ⏳ Docs |
//	| Phase 2 | Planned     | • Dashboard<br>• Monitoring           |
func (r *Roadmap) ToPhaseTable(opts TableOptions) string {
	if len(r.Phases) == 0 {
		return ""
	}

	var sb strings.Builder

	// Header
	sb.WriteString("| Phase | Status | Deliverables |\n")
	sb.WriteString("|-------|--------|---------------|\n")

	for i, phase := range r.Phases {
		var items []string
		for _, del := range phase.Deliverables {
			item := del.Title
			if opts.MaxTitleLen > 0 && len(item) > opts.MaxTitleLen {
				item = item[:opts.MaxTitleLen-3] + "..."
			}
			if opts.IncludeStatus && del.Status != "" {
				item = fmt.Sprintf("%s %s", StatusIconWithOptions(del.Status, opts.UseTextIcons), item)
			}
			// Add bullet point prefix
			items = append(items, "• "+item)
		}

		status := string(phase.Status)
		if status == "" {
			status = "planned"
		}

		// Join with <br> for line breaks within cell
		sb.WriteString(fmt.Sprintf("| **Phase %d**<br>%s | %s | %s |\n",
			i+1, phase.Name, status, strings.Join(items, "<br>")))
	}

	return sb.String()
}

// SwimlaneLabel converts a DeliverableType to a human-readable label.
func SwimlaneLabel(dt DeliverableType) string {
	switch dt {
	case DeliverableFeature:
		return "Features"
	case DeliverableDocumentation:
		return "Documentation"
	case DeliverableInfrastructure:
		return "Infrastructure"
	case DeliverableIntegration:
		return "Integrations"
	case DeliverableMilestone:
		return "Milestones"
	case DeliverableRollout:
		return "Rollout"
	default:
		// Capitalize the first letter
		s := string(dt)
		if len(s) > 0 {
			return strings.ToUpper(s[:1]) + s[1:]
		}
		return s
	}
}

// StatusIcon returns an emoji/icon for the deliverable status.
func StatusIcon(status DeliverableStatus) string {
	return StatusIconWithOptions(status, false)
}

// StatusIconWithOptions returns an icon for the deliverable status.
// If useText is true, returns ASCII text instead of emoji for PDF compatibility.
func StatusIconWithOptions(status DeliverableStatus, useText bool) string {
	if useText {
		switch status {
		case DeliverableCompleted:
			return "[DONE]"
		case DeliverableInProgress:
			return "[WIP]"
		case DeliverableBlocked:
			return "[BLOCKED]"
		case DeliverableNotStarted:
			return "[TODO]"
		default:
			return ""
		}
	}
	switch status {
	case DeliverableCompleted:
		return "✅"
	case DeliverableInProgress:
		return "🔄"
	case DeliverableBlocked:
		return "🚫"
	case DeliverableNotStarted:
		return "⏳"
	default:
		return ""
	}
}

// PhaseTargetStatusIcon returns an emoji/icon for the phase target status.
func PhaseTargetStatusIcon(status string) string {
	return PhaseTargetStatusIconWithOptions(status, false)
}

// PhaseTargetStatusIconWithOptions returns an icon for the phase target status.
// If useText is true, returns ASCII text instead of emoji for PDF compatibility.
func PhaseTargetStatusIconWithOptions(status string, useText bool) string {
	if useText {
		switch status {
		case "achieved":
			return "[DONE]"
		case "in_progress":
			return "[WIP]"
		case "missed":
			return "[MISSED]"
		case "not_started":
			return "[TODO]"
		default:
			return ""
		}
	}
	switch status {
	case "achieved":
		return "✅"
	case "in_progress":
		return "🔄"
	case "missed":
		return "❌"
	case "not_started":
		return "⏳"
	default:
		return ""
	}
}

// StatusLegend returns a markdown table explaining the status icons.
func StatusLegend() string {
	return StatusLegendWithOptions(false)
}

// StatusLegendWithOptions returns a markdown table explaining the status icons.
// If useText is true, shows ASCII text icons instead of emoji.
func StatusLegendWithOptions(useText bool) string {
	if useText {
		return `| Icon | Status |
|------|--------|
| [DONE] | Completed / Achieved |
| [WIP] | In Progress |
| [TODO] | Not Started |
| [BLOCKED] | Blocked |
| [MISSED] | Missed |
`
	}
	return `| Icon | Status |
|------|--------|
| ✅ | Completed / Achieved |
| 🔄 | In Progress |
| ⏳ | Not Started |
| 🚫 | Blocked |
| ❌ | Missed |
`
}
