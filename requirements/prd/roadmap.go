package prd

import (
	"github.com/grokify/prism-execution/roadmap"
)

// Type aliases for backward compatibility.
// These allow existing code using prd.Roadmap, prd.Phase, etc. to continue working.

// Roadmap contains the product roadmap with phases.
type Roadmap = roadmap.Roadmap

// Phase represents a roadmap phase.
type Phase = roadmap.Phase

// PhaseType represents the type of roadmap phase.
type PhaseType = roadmap.PhaseType

// PhaseStatus represents the current status of a phase.
type PhaseStatus = roadmap.PhaseStatus

// Deliverable represents a phase deliverable.
type Deliverable = roadmap.Deliverable

// DeliverableType represents types of deliverables.
type DeliverableType = roadmap.DeliverableType

// DeliverableStatus represents the status of a deliverable.
type DeliverableStatus = roadmap.DeliverableStatus

// Constants re-exported for backward compatibility.
const (
	// PhaseType constants
	PhaseTypeGeneric   = roadmap.PhaseTypeGeneric
	PhaseTypeQuarter   = roadmap.PhaseTypeQuarter
	PhaseTypeMonth     = roadmap.PhaseTypeMonth
	PhaseTypeSprint    = roadmap.PhaseTypeSprint
	PhaseTypeMilestone = roadmap.PhaseTypeMilestone

	// PhaseStatus constants
	PhaseStatusPlanned    = roadmap.PhaseStatusPlanned
	PhaseStatusInProgress = roadmap.PhaseStatusInProgress
	PhaseStatusCompleted  = roadmap.PhaseStatusCompleted
	PhaseStatusDelayed    = roadmap.PhaseStatusDelayed
	PhaseStatusCancelled  = roadmap.PhaseStatusCancelled

	// DeliverableType constants
	DeliverableFeature        = roadmap.DeliverableFeature
	DeliverableDocumentation  = roadmap.DeliverableDocumentation
	DeliverableInfrastructure = roadmap.DeliverableInfrastructure
	DeliverableIntegration    = roadmap.DeliverableIntegration
	DeliverableMilestone      = roadmap.DeliverableMilestone
	DeliverableRollout        = roadmap.DeliverableRollout

	// DeliverableStatus constants
	DeliverableNotStarted = roadmap.DeliverableNotStarted
	DeliverableInProgress = roadmap.DeliverableInProgress
	DeliverableCompleted  = roadmap.DeliverableCompleted
	DeliverableBlocked    = roadmap.DeliverableBlocked
)
