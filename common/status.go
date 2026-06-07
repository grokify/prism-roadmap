package common

import (
	pf "github.com/grokify/priority-frameworks"
	core "github.com/grokify/prism-core"
)

// Status represents the document lifecycle status.
// Used across PRD, MRD, and TRD documents.
type Status string

// Status constants imported from prism-core.
const (
	StatusDraft      Status = Status(core.StatusDraft)
	StatusInReview   Status = Status(core.StatusInReview)
	StatusApproved   Status = Status(core.StatusApproved)
	StatusDeprecated Status = Status(core.StatusDeprecated)
)

// Priority represents priority levels.
// Now backed by priority-frameworks Severity framework.
type Priority string

// Priority constants using priority-frameworks.
const (
	PriorityCritical Priority = "critical"
	PriorityHigh     Priority = "high"
	PriorityMedium   Priority = "medium"
	PriorityLow      Priority = "low"
)

// ValidPriority checks if a priority value is valid.
func ValidPriority(priority Priority) bool {
	if priority == "" {
		return true
	}
	f := pf.Severity()
	return f.IndexOf(string(priority)) >= 0
}

// PriorityWeight returns a numeric weight for sorting priorities.
// Higher weight = higher priority.
func PriorityWeight(priority Priority) int {
	f := pf.Severity()
	idx := f.IndexOf(string(priority))
	if idx < 0 {
		return 0
	}
	return len(f.Levels) - idx
}

// PriorityFramework returns the Severity priority framework.
func PriorityFramework() *pf.Framework {
	return pf.Severity()
}
