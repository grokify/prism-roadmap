package common

import core "github.com/grokify/prism-core"

// RiskProbability represents risk probability levels.
type RiskProbability string

// Risk probability constants imported from prism-core.
const (
	RiskProbabilityLow    RiskProbability = RiskProbability(core.RiskProbabilityLow)
	RiskProbabilityMedium RiskProbability = RiskProbability(core.RiskProbabilityMedium)
	RiskProbabilityHigh   RiskProbability = RiskProbability(core.RiskProbabilityHigh)
)

// RiskImpact represents risk impact levels.
type RiskImpact string

// Risk impact constants imported from prism-core.
const (
	RiskImpactLow      RiskImpact = RiskImpact(core.RiskImpactLow)
	RiskImpactMedium   RiskImpact = RiskImpact(core.RiskImpactMedium)
	RiskImpactHigh     RiskImpact = RiskImpact(core.RiskImpactHigh)
	RiskImpactCritical RiskImpact = RiskImpact(core.RiskImpactCritical)
)

// RiskStatus represents risk status.
type RiskStatus string

// Risk status constants imported from prism-core.
const (
	RiskStatusOpen      RiskStatus = RiskStatus(core.RiskStatusOpen)
	RiskStatusMitigated RiskStatus = RiskStatus(core.RiskStatusMitigated)
	RiskStatusAccepted  RiskStatus = RiskStatus(core.RiskStatusAccepted)
	RiskStatusClosed    RiskStatus = RiskStatus(core.RiskStatusClosed)
)

// Risk represents a project risk.
// Used across PRD, MRD, and TRD documents.
// Extends core.Risk with additional document-specific fields.
type Risk struct {
	ID          string          `json:"id"`
	Description string          `json:"description"`
	Probability RiskProbability `json:"probability"`
	Impact      RiskImpact      `json:"impact"`
	Mitigation  string          `json:"mitigation"`
	Owner       string          `json:"owner,omitempty"`
	Status      RiskStatus      `json:"status,omitempty"`
	Category    string          `json:"category,omitempty"` // Market, Competitive, Technical, etc.
	DueDate     string          `json:"dueDate,omitempty"`
	Tags        []string        `json:"tags,omitempty"`
	Notes       string          `json:"notes,omitempty"`

	// AppendixRefs references appendices with additional details for this risk.
	AppendixRefs []string `json:"appendixRefs,omitempty"`
}

// ValidRiskProbability checks if a risk probability value is valid.
func ValidRiskProbability(probability RiskProbability) bool {
	return core.ValidRiskProbability(string(probability))
}

// ValidRiskImpact checks if a risk impact value is valid.
func ValidRiskImpact(impact RiskImpact) bool {
	return core.ValidRiskImpact(string(impact))
}

// ValidRiskStatus checks if a risk status value is valid.
func ValidRiskStatus(status RiskStatus) bool {
	return core.ValidRiskStatus(string(status))
}

// RiskScore calculates a numeric score based on probability and impact.
func RiskScore(probability RiskProbability, impact RiskImpact) int {
	return core.RiskScore(string(probability), string(impact))
}

// RiskSeverity returns a severity level based on risk score.
func RiskSeverity(score int) string {
	return core.RiskSeverity(score)
}
