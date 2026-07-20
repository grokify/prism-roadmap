// Package goals provides cascading functionality for goal hierarchies.
package goals

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/goals/okr"
)

// CascadeOptions configures how goals cascade from parent to child.
type CascadeOptions struct {
	// ChildTeam is the name of the team receiving cascaded goals.
	ChildTeam string

	// ChildOwner is the owner/lead for child goals.
	ChildOwner string

	// ChildPeriod is the period for child goals (e.g., "Q2 2026").
	ChildPeriod string

	// NameTemplate is a Go text/template for generating child objective names.
	// Available fields: .ParentTitle, .ChildTeam, .Index
	// Default: "{{.ParentTitle}} - {{.ChildTeam}}"
	NameTemplate string

	// FilterObjectiveIDs limits cascading to specific parent objectives.
	// If empty, all objectives are cascaded.
	FilterObjectiveIDs []string

	// InheritRisks copies risks from parent objectives.
	InheritRisks bool

	// InheritTags copies tags from parent objectives.
	InheritTags bool

	// CreateKeyResultsFromParent creates child key results based on parent KRs.
	CreateKeyResultsFromParent bool
}

// CascadeResult contains the cascaded goals and alignment information.
type CascadeResult struct {
	// ChildGoals contains the generated child goals.
	ChildGoals *Goals

	// Alignment maps child goal IDs to parent goal IDs.
	Alignment map[string]string

	// ParentID is the ID of the parent document.
	ParentID string
}

// CascadeOKR generates child OKRs from a parent OKR document.
func CascadeOKR(parent *okr.OKRDocument, opts CascadeOptions) (*CascadeResult, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent OKR document is nil")
	}

	// Filter objectives if specified
	objectives := parent.Objectives
	if len(opts.FilterObjectiveIDs) > 0 {
		filterSet := make(map[string]bool)
		for _, id := range opts.FilterObjectiveIDs {
			filterSet[id] = true
		}
		filtered := make([]okr.Objective, 0)
		for _, obj := range objectives {
			if filterSet[obj.ID] {
				filtered = append(filtered, obj)
			}
		}
		objectives = filtered
	}

	if len(objectives) == 0 {
		return nil, fmt.Errorf("no objectives to cascade")
	}

	// Generate child objectives
	childObjectives := make([]okr.Objective, 0, len(objectives))
	alignment := make(map[string]string)

	for i, parentObj := range objectives {
		childID := generateChildID(parentObj.ID, opts.ChildTeam, i)

		// Generate child title
		childTitle := generateChildTitle(parentObj.Title, opts.ChildTeam, opts.NameTemplate)

		childObj := okr.Objective{
			ID:          childID,
			Title:       childTitle,
			Description: fmt.Sprintf("Team contribution to: %s", parentObj.Title),
			Rationale:   parentObj.Rationale,
			Category:    parentObj.Category,
			Owner:       opts.ChildOwner,
			Timeframe:   opts.ChildPeriod,
			Status:      okr.StatusDraft,
			ParentID:    parentObj.ID,
			AlignedWith: []string{parentObj.ID},
		}

		// Inherit tags
		if opts.InheritTags && len(parentObj.Tags) > 0 {
			childObj.Tags = append([]string{}, parentObj.Tags...)
		}

		// Inherit risks
		if opts.InheritRisks && len(parentObj.Risks) > 0 {
			childObj.Risks = make([]okr.Risk, len(parentObj.Risks))
			for j, risk := range parentObj.Risks {
				childObj.Risks[j] = okr.Risk{
					ID:          fmt.Sprintf("%s-risk-%d", childID, j+1),
					Title:       risk.Title,
					Description: risk.Description,
					Impact:      risk.Impact,
					Likelihood:  risk.Likelihood,
					Status:      "Identified",
				}
			}
		}

		// Create key results
		if opts.CreateKeyResultsFromParent && len(parentObj.KeyResults) > 0 {
			childObj.KeyResults = make([]okr.KeyResult, 0, len(parentObj.KeyResults))
			for k, parentKR := range parentObj.KeyResults {
				childKR := okr.KeyResult{
					ID:          fmt.Sprintf("%s-kr-%d", childID, k+1),
					Title:       fmt.Sprintf("Contribute to: %s", parentKR.Title),
					Description: parentKR.Description,
					Metric:      parentKR.Metric,
					Unit:        parentKR.Unit,
					Confidence:  okr.ConfidenceMedium,
					Status:      "Not Started",
				}
				// Copy phase targets
				if len(parentKR.PhaseTargets) > 0 {
					childKR.PhaseTargets = make([]okr.PhaseTarget, len(parentKR.PhaseTargets))
					for p, pt := range parentKR.PhaseTargets {
						childKR.PhaseTargets[p] = okr.PhaseTarget{
							PhaseID: pt.PhaseID,
							Status:  "not_started",
						}
					}
				}
				childObj.KeyResults = append(childObj.KeyResults, childKR)
			}
		} else {
			// Create placeholder key results
			childObj.KeyResults = []okr.KeyResult{
				{
					ID:         fmt.Sprintf("%s-kr-1", childID),
					Title:      "Define team-specific key result",
					Confidence: okr.ConfidenceLow,
					Status:     "Not Started",
				},
			}
		}

		childObjectives = append(childObjectives, childObj)
		alignment[childID] = parentObj.ID
	}

	// Build parent ID
	parentID := ""
	if parent.Metadata != nil && parent.Metadata.ID != "" {
		parentID = parent.Metadata.ID
	}

	// Note: child OKR document structure is built by the caller (CLI)
	// This function returns the Goals wrapper for flexibility

	// Convert to Goals wrapper
	okrSet := okr.FromObjectives(childObjectives)
	childGoals := NewOKR(okrSet)

	return &CascadeResult{
		ChildGoals: childGoals,
		Alignment:  alignment,
		ParentID:   parentID,
	}, nil
}

// generateChildID creates a unique ID for a child objective.
func generateChildID(parentID, childTeam string, index int) string {
	teamSlug := strings.ToLower(strings.ReplaceAll(childTeam, " ", "-"))
	if parentID == "" {
		return fmt.Sprintf("%s-obj-%d", teamSlug, index+1)
	}
	return fmt.Sprintf("%s-%s-%d", parentID, teamSlug, index+1)
}

// generateChildTitle creates a title for a child objective.
func generateChildTitle(parentTitle, childTeam, template string) string {
	if template == "" {
		return fmt.Sprintf("%s - %s", parentTitle, childTeam)
	}
	// Simple template replacement (could use text/template for more complex cases)
	result := strings.ReplaceAll(template, "{{.ParentTitle}}", parentTitle)
	result = strings.ReplaceAll(result, "{{.ChildTeam}}", childTeam)
	return result
}

// AlignmentScore calculates how well child goals align with parent goals.
// Returns a score from 0.0 (no alignment) to 1.0 (full alignment).
type AlignmentScore struct {
	// Overall alignment score (0.0 - 1.0).
	Score float64

	// CoveredParentGoals is the count of parent goals with child support.
	CoveredParentGoals int

	// TotalParentGoals is the total count of parent goals.
	TotalParentGoals int

	// OrphanedChildGoals is the count of child goals without parent alignment.
	OrphanedChildGoals int

	// TotalChildGoals is the total count of child goals.
	TotalChildGoals int

	// Issues lists alignment problems found.
	Issues []string
}

// CalculateAlignment calculates the alignment between parent and child goals.
func CalculateAlignment(parent, child *Goals) (*AlignmentScore, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent goals is nil")
	}
	if child == nil {
		return nil, fmt.Errorf("child goals is nil")
	}

	parentItems := parent.GoalItems()
	childItems := child.GoalItems()

	if len(parentItems) == 0 {
		return &AlignmentScore{
			Score:  0,
			Issues: []string{"Parent has no goals"},
		}, nil
	}

	// For OKR, check AlignedWith fields
	if child.IsOKR() && child.OKR != nil {
		return calculateOKRAlignment(parentItems, child.OKR)
	}

	// Default: simple ID-based alignment check
	parentIDs := make(map[string]bool)
	for _, p := range parentItems {
		if p.ID != "" {
			parentIDs[p.ID] = true
		}
	}

	covered := make(map[string]bool)
	orphaned := len(childItems) // Simplified: all orphaned without OKR-specific check

	score := &AlignmentScore{
		CoveredParentGoals: len(covered),
		TotalParentGoals:   len(parentItems),
		OrphanedChildGoals: orphaned,
		TotalChildGoals:    len(childItems),
	}

	if score.TotalParentGoals > 0 {
		score.Score = float64(score.CoveredParentGoals) / float64(score.TotalParentGoals)
	}

	return score, nil
}

// calculateOKRAlignment calculates alignment for OKR goals.
func calculateOKRAlignment(parentItems []GoalItem, childOKR *okr.OKRSet) (*AlignmentScore, error) {
	parentIDs := make(map[string]bool)
	for _, p := range parentItems {
		if p.ID != "" {
			parentIDs[p.ID] = true
		}
	}

	covered := make(map[string]bool)
	orphaned := 0
	var issues []string

	for _, okrItem := range childOKR.OKRs {
		obj := okrItem.Objective
		hasAlignment := false

		// Check AlignedWith
		for _, alignedID := range obj.AlignedWith {
			if parentIDs[alignedID] {
				covered[alignedID] = true
				hasAlignment = true
			} else if alignedID != "" {
				issues = append(issues, fmt.Sprintf("Objective %q references unknown parent %q", obj.Title, alignedID))
			}
		}

		// Check ParentID
		if obj.ParentID != "" {
			if parentIDs[obj.ParentID] {
				covered[obj.ParentID] = true
				hasAlignment = true
			} else {
				issues = append(issues, fmt.Sprintf("Objective %q references unknown parent %q", obj.Title, obj.ParentID))
			}
		}

		if !hasAlignment && len(obj.AlignedWith) == 0 && obj.ParentID == "" {
			orphaned++
			issues = append(issues, fmt.Sprintf("Objective %q has no parent alignment", obj.Title))
		}
	}

	// Check for uncovered parent goals
	for _, p := range parentItems {
		if p.ID != "" && !covered[p.ID] {
			issues = append(issues, fmt.Sprintf("Parent objective %q has no child support", p.Title))
		}
	}

	score := &AlignmentScore{
		CoveredParentGoals: len(covered),
		TotalParentGoals:   len(parentItems),
		OrphanedChildGoals: orphaned,
		TotalChildGoals:    len(childOKR.OKRs),
		Issues:             issues,
	}

	if score.TotalParentGoals > 0 {
		score.Score = float64(score.CoveredParentGoals) / float64(score.TotalParentGoals)
	}

	return score, nil
}
