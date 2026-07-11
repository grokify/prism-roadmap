package journey

// Initiative represents a planned effort that enables capability transitions.
// Initiatives connect capability evolution to actual execution by teams.
type Initiative struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Description       string              `json:"description,omitempty"`
	Type              InitiativeType      `json:"type,omitempty"`
	Status            InitiativeStatus    `json:"status,omitempty"`
	Periods           []string            `json:"periods,omitempty"` // Period IDs when active
	StartPeriod       string              `json:"startPeriod,omitempty"`
	EndPeriod         string              `json:"endPeriod,omitempty"`
	OwnerTeam         string              `json:"ownerTeam,omitempty"`         // Primary responsible team
	ContributingTeams []string            `json:"contributingTeams,omitempty"` // Supporting teams
	Advances          []CapabilityAdvance `json:"advances,omitempty"`          // Capability improvements
	Outputs           []string            `json:"outputs,omitempty"`           // Deliverables
	ExpectedOutcomes  []string            `json:"expectedOutcomes,omitempty"`  // Outcome journey IDs
	RequiredCapacity  *Capacity           `json:"requiredCapacity,omitempty"`
	Confidence        float64             `json:"confidence,omitempty"` // 0.0-1.0
	Priority          string              `json:"priority,omitempty"`   // P0, P1, P2, P3
	Tags              []string            `json:"tags,omitempty"`
	Links             []Link              `json:"links,omitempty"` // External references
	Metadata          map[string]any      `json:"metadata,omitempty"`
}

// InitiativeType classifies initiatives.
type InitiativeType string

const (
	InitiativeTypePlatform   InitiativeType = "platform"
	InitiativeTypeProduct    InitiativeType = "product"
	InitiativeTypeProcess    InitiativeType = "process"
	InitiativeTypeEnablement InitiativeType = "enablement"
	InitiativeTypeCompliance InitiativeType = "compliance"
	InitiativeTypeDebt       InitiativeType = "technical_debt"
	InitiativeTypeExperiment InitiativeType = "experiment"
)

// InitiativeStatus tracks initiative progress.
type InitiativeStatus string

const (
	InitiativeStatusDraft      InitiativeStatus = "draft"
	InitiativeStatusProposed   InitiativeStatus = "proposed"
	InitiativeStatusApproved   InitiativeStatus = "approved"
	InitiativeStatusInProgress InitiativeStatus = "in_progress"
	InitiativeStatusCompleted  InitiativeStatus = "completed"
	InitiativeStatusOnHold     InitiativeStatus = "on_hold"
	InitiativeStatusCancelled  InitiativeStatus = "cancelled"
)

// CapabilityAdvance describes how an initiative improves a capability.
type CapabilityAdvance struct {
	CapabilityID   string `json:"capabilityId"`
	CapabilityName string `json:"capabilityName,omitempty"`
	From           string `json:"from"` // Starting maturity level
	To             string `json:"to"`   // Target maturity level
	Description    string `json:"description,omitempty"`
}

// Link represents an external reference (e.g., to a project tracker, doc).
type Link struct {
	Type  string `json:"type,omitempty"` // jira, confluence, github, notion
	URL   string `json:"url"`
	Title string `json:"title,omitempty"`
}

// Capacity represents resource requirements.
type Capacity struct {
	StoryPoints *int     `json:"storyPoints,omitempty"`
	FTEs        *float64 `json:"ftes,omitempty"`       // Full-time equivalents
	FTEMonths   *float64 `json:"fteMonths,omitempty"`  // FTE-months of effort
	CustomUnit  string   `json:"customUnit,omitempty"` // For other units
	CustomValue *float64 `json:"customValue,omitempty"`
	Notes       string   `json:"notes,omitempty"`
}

// TotalFTEMonths returns FTE-months, calculating from FTEs and periods if needed.
func (c *Capacity) TotalFTEMonths(periodCount int, monthsPerPeriod float64) float64 {
	if c == nil {
		return 0
	}
	if c.FTEMonths != nil {
		return *c.FTEMonths
	}
	if c.FTEs != nil && periodCount > 0 {
		return *c.FTEs * float64(periodCount) * monthsPerPeriod
	}
	return 0
}

// InitiativeSummary provides a high-level view of an initiative.
type InitiativeSummary struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Status          string   `json:"status"`
	OwnerTeam       string   `json:"ownerTeam"`
	Periods         []string `json:"periods"`
	CapabilityCount int      `json:"capabilityCount"`
	OutcomeCount    int      `json:"outcomeCount"`
	MaturityGain    int      `json:"maturityGain"` // Total maturity levels gained
}

// Summarize creates a summary of the initiative.
func (i *Initiative) Summarize() InitiativeSummary {
	maturityGain := 0
	for _, adv := range i.Advances {
		fromLevel := parseMaturityLevel(adv.From)
		toLevel := parseMaturityLevel(adv.To)
		maturityGain += toLevel - fromLevel
	}

	return InitiativeSummary{
		ID:              i.ID,
		Name:            i.Name,
		Status:          string(i.Status),
		OwnerTeam:       i.OwnerTeam,
		Periods:         i.Periods,
		CapabilityCount: len(i.Advances),
		OutcomeCount:    len(i.ExpectedOutcomes),
		MaturityGain:    maturityGain,
	}
}
