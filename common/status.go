package common

import core "github.com/grokify/prism-core"

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
type Priority string

// Priority constants imported from prism-core.
const (
	PriorityCritical Priority = Priority(core.PriorityCritical)
	PriorityHigh     Priority = Priority(core.PriorityHigh)
	PriorityMedium   Priority = Priority(core.PriorityMedium)
	PriorityLow      Priority = Priority(core.PriorityLow)
)

// ValidPriority checks if a priority value is valid.
func ValidPriority(priority Priority) bool {
	return core.ValidPriority(string(priority))
}

// PriorityWeight returns a numeric weight for sorting priorities.
func PriorityWeight(priority Priority) int {
	return core.PriorityWeight(string(priority))
}
