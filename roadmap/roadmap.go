// Package roadmap provides types for product and project roadmaps.
// Roadmaps can be used standalone or embedded in PRD/MRD/TRD documents.
package roadmap

import "time"

// Roadmap contains the product roadmap with phases.
type Roadmap struct {
	Phases []Phase `json:"phases"`
}

// PhaseType represents the type of roadmap phase.
type PhaseType string

const (
	PhaseTypeGeneric   PhaseType = "generic"   // Phase 1, 2, 3
	PhaseTypeQuarter   PhaseType = "quarter"   // Q1 2026, Q2 2026
	PhaseTypeMonth     PhaseType = "month"     // January 2026
	PhaseTypeSprint    PhaseType = "sprint"    // Sprint 1, Sprint 2
	PhaseTypeMilestone PhaseType = "milestone" // MVP, GA, etc.
)

// Phase represents a roadmap phase.
type Phase struct {
	ID              string        `json:"id"`   // e.g., "phase-1", "q1-2026"
	Name            string        `json:"name"` // e.g., "MVP", "Q1 2026"
	Type            PhaseType     `json:"type"`
	StartDate       *time.Time    `json:"startDate,omitempty"`
	EndDate         *time.Time    `json:"endDate,omitempty"`
	Goals           []string      `json:"goals"`
	Deliverables    []Deliverable `json:"deliverables"`
	SuccessCriteria []string      `json:"successCriteria"`
	Dependencies    []string      `json:"dependencies,omitempty"` // Dependent phase IDs
	Risks           []Risk        `json:"risks,omitempty"`
	Status          PhaseStatus   `json:"status,omitempty"`
	Progress        *int          `json:"progress,omitempty"` // 0-100 percentage
	Tags            []string      `json:"tags,omitempty"`     // For filtering by topic/domain
	Notes           string        `json:"notes,omitempty"`
}

// PhaseStatus represents the current status of a phase.
type PhaseStatus string

const (
	PhaseStatusPlanned    PhaseStatus = "planned"
	PhaseStatusInProgress PhaseStatus = "in_progress"
	PhaseStatusCompleted  PhaseStatus = "completed"
	PhaseStatusDelayed    PhaseStatus = "delayed"
	PhaseStatusCancelled  PhaseStatus = "cancelled"
)

// Deliverable represents a phase deliverable.
type Deliverable struct {
	ID          string            `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Type        DeliverableType   `json:"type"`
	Status      DeliverableStatus `json:"status,omitempty"`
	Tags        []string          `json:"tags,omitempty"`    // For filtering by topic/domain
	Rollout     *RolloutStatus    `json:"rollout,omitempty"` // B2B deployment tracking
}

// RolloutStatus tracks deployment and adoption for B2B SaaS deliverables.
// Deployment = feature is available to customer (rolled out).
// Adoption = customer is actively using the feature.
type RolloutStatus struct {
	TotalCustomers    int           `json:"totalCustomers"`             // Total customers in rollout scope
	DeployedCustomers int           `json:"deployedCustomers"`          // Customers with feature available
	AdoptedCustomers  int           `json:"adoptedCustomers,omitempty"` // Customers actively using feature
	Status            RolloutStage  `json:"status,omitempty"`           // Current rollout stage
	StartDate         string        `json:"startDate,omitempty"`        // Rollout start date (ISO 8601)
	TargetDate        string        `json:"targetDate,omitempty"`       // Target completion date
	Notes             string        `json:"notes,omitempty"`            // Rollout notes
	Waves             []RolloutWave `json:"waves,omitempty"`            // Phased rollout waves
}

// RolloutStage represents the current stage of a rollout.
type RolloutStage string

const (
	RolloutStageNotStarted RolloutStage = "not_started"
	RolloutStageRollingOut RolloutStage = "rolling_out"
	RolloutStageDeployed   RolloutStage = "deployed" // 100% deployed, adoption ongoing
	RolloutStageAdopted    RolloutStage = "adopted"  // Target adoption achieved
	RolloutStagePaused     RolloutStage = "paused"   // Rollout paused
	RolloutStageRolledBack RolloutStage = "rolled_back"
)

// RolloutWave represents a phased rollout wave (e.g., beta, GA waves).
type RolloutWave struct {
	ID                string `json:"id"`
	Name              string `json:"name"`                        // e.g., "Beta", "Wave 1", "GA"
	TargetCustomers   int    `json:"targetCustomers"`             // Customers in this wave
	DeployedCustomers int    `json:"deployedCustomers,omitempty"` // Deployed in this wave
	StartDate         string `json:"startDate,omitempty"`
	EndDate           string `json:"endDate,omitempty"`
	Status            string `json:"status,omitempty"` // planned, in_progress, completed
}

// DeploymentPercent returns the percentage of customers with the feature deployed.
func (r *RolloutStatus) DeploymentPercent() float64 {
	if r == nil || r.TotalCustomers == 0 {
		return 0
	}
	return float64(r.DeployedCustomers) / float64(r.TotalCustomers) * 100
}

// AdoptionPercent returns the percentage of customers actively using the feature.
// This is relative to total customers, not deployed customers.
func (r *RolloutStatus) AdoptionPercent() float64 {
	if r == nil || r.TotalCustomers == 0 {
		return 0
	}
	return float64(r.AdoptedCustomers) / float64(r.TotalCustomers) * 100
}

// AdoptionOfDeployed returns adoption as a percentage of deployed customers.
// This measures adoption among customers who have access to the feature.
func (r *RolloutStatus) AdoptionOfDeployed() float64 {
	if r == nil || r.DeployedCustomers == 0 {
		return 0
	}
	return float64(r.AdoptedCustomers) / float64(r.DeployedCustomers) * 100
}

// IsFullyDeployed returns true if all customers have the feature deployed.
func (r *RolloutStatus) IsFullyDeployed() bool {
	if r == nil {
		return false
	}
	return r.DeployedCustomers >= r.TotalCustomers
}

// DeliverableType represents types of deliverables.
type DeliverableType string

const (
	DeliverableFeature        DeliverableType = "feature"
	DeliverableDocumentation  DeliverableType = "documentation"
	DeliverableInfrastructure DeliverableType = "infrastructure"
	DeliverableIntegration    DeliverableType = "integration"
	DeliverableMilestone      DeliverableType = "milestone"
	DeliverableRollout        DeliverableType = "rollout"
)

// DeliverableStatus represents the status of a deliverable.
type DeliverableStatus string

const (
	DeliverableNotStarted DeliverableStatus = "not_started"
	DeliverableInProgress DeliverableStatus = "in_progress"
	DeliverableCompleted  DeliverableStatus = "completed"
	DeliverableBlocked    DeliverableStatus = "blocked"
)

// Risk represents a risk associated with a roadmap phase.
// This is a simplified risk type for roadmap use; document-level risks
// in PRD/MRD/TRD may have additional fields.
type Risk struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Probability string   `json:"probability"` // Low, Medium, High
	Impact      string   `json:"impact"`      // Low, Medium, High, Critical
	Mitigation  string   `json:"mitigation"`
	Status      string   `json:"status,omitempty"` // Identified, Mitigating, Resolved, Accepted
	Tags        []string `json:"tags,omitempty"`   // For filtering by topic/domain
}
