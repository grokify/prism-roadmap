package journey

// Team represents a team that can own or contribute to initiatives.
// Teams can be organized in a hierarchy (VP → Director → Team).
type Team struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Type        TeamType       `json:"type,omitempty"`
	Level       TeamLevel      `json:"level,omitempty"`    // Org level
	ParentID    string         `json:"parentId,omitempty"` // Parent team for hierarchy
	ChildIDs    []string       `json:"childIds,omitempty"` // Child teams
	LeaderID    string         `json:"leaderId,omitempty"` // Team lead/manager
	LeaderName  string         `json:"leaderName,omitempty"`
	Capacity    *TeamCapacity  `json:"capacity,omitempty"`
	Skills      []string       `json:"skills,omitempty"` // e.g., "backend", "ml", "security"
	CostCenter  string         `json:"costCenter,omitempty"`
	Location    string         `json:"location,omitempty"`
	Tags        []string       `json:"tags,omitempty"`
	Metadata    map[string]any `json:"metadata,omitempty"`
}

// TeamType classifies the team.
type TeamType string

const (
	TeamTypeEngineering TeamType = "engineering"
	TeamTypePlatform    TeamType = "platform"
	TeamTypeProduct     TeamType = "product"
	TeamTypeDesign      TeamType = "design"
	TeamTypeData        TeamType = "data"
	TeamTypeInfra       TeamType = "infrastructure"
	TeamTypeSecurity    TeamType = "security"
	TeamTypeQA          TeamType = "qa"
	TeamTypeDevOps      TeamType = "devops"
	TeamTypeSRE         TeamType = "sre"
	TeamTypeExternal    TeamType = "external" // Contractors, vendors
)

// TeamLevel indicates position in org hierarchy.
type TeamLevel string

const (
	TeamLevelOrganization TeamLevel = "organization" // Top-level org
	TeamLevelDivision     TeamLevel = "division"     // Division/BU
	TeamLevelDepartment   TeamLevel = "department"   // VP-level
	TeamLevelGroup        TeamLevel = "group"        // Director-level
	TeamLevelTeam         TeamLevel = "team"         // Individual team
	TeamLevelSquad        TeamLevel = "squad"        // Sub-team
)

// TeamCapacity represents a team's available capacity.
type TeamCapacity struct {
	// FTEs is the number of full-time equivalent engineers.
	FTEs float64 `json:"ftes,omitempty"`

	// StoryPointsPerSprint is the team's velocity in story points.
	StoryPointsPerSprint int `json:"storyPointsPerSprint,omitempty"`

	// StoryPointsPerQuarter for quarterly planning.
	StoryPointsPerQuarter int `json:"storyPointsPerQuarter,omitempty"`

	// AllocatedPercent is how much of capacity is already committed (0-100).
	AllocatedPercent float64 `json:"allocatedPercent,omitempty"`

	// ReservedPercent is capacity reserved for maintenance, on-call, etc.
	ReservedPercent float64 `json:"reservedPercent,omitempty"`

	// EffectiveDate when this capacity snapshot was taken.
	EffectiveDate string `json:"effectiveDate,omitempty"`

	// Notes for context.
	Notes string `json:"notes,omitempty"`
}

// AvailableCapacity returns the capacity available for new work.
func (tc *TeamCapacity) AvailableCapacity() *TeamCapacity {
	if tc == nil {
		return nil
	}

	availablePercent := 100.0 - tc.AllocatedPercent - tc.ReservedPercent
	if availablePercent < 0 {
		availablePercent = 0
	}

	factor := availablePercent / 100.0

	available := &TeamCapacity{
		FTEs:                  tc.FTEs * factor,
		StoryPointsPerSprint:  int(float64(tc.StoryPointsPerSprint) * factor),
		StoryPointsPerQuarter: int(float64(tc.StoryPointsPerQuarter) * factor),
		EffectiveDate:         tc.EffectiveDate,
	}

	return available
}

// TeamAllocation represents capacity allocated to a specific initiative.
type TeamAllocation struct {
	TeamID       string    `json:"teamId"`
	InitiativeID string    `json:"initiativeId"`
	Periods      []string  `json:"periods,omitempty"`    // Active periods
	Allocation   *Capacity `json:"allocation,omitempty"` // Allocated capacity
	Role         string    `json:"role,omitempty"`       // "owner", "contributor", "reviewer"
	Notes        string    `json:"notes,omitempty"`
}

// TeamHierarchy provides methods to traverse team structure.
type TeamHierarchy struct {
	Teams map[string]*Team `json:"teams"`
}

// NewTeamHierarchy creates a hierarchy from a list of teams.
func NewTeamHierarchy(teams []Team) *TeamHierarchy {
	h := &TeamHierarchy{
		Teams: make(map[string]*Team, len(teams)),
	}
	for i := range teams {
		h.Teams[teams[i].ID] = &teams[i]
	}
	return h
}

// GetParent returns the parent team.
func (h *TeamHierarchy) GetParent(teamID string) *Team {
	team, ok := h.Teams[teamID]
	if !ok || team.ParentID == "" {
		return nil
	}
	return h.Teams[team.ParentID]
}

// GetChildren returns direct child teams.
func (h *TeamHierarchy) GetChildren(teamID string) []*Team {
	var children []*Team
	for _, team := range h.Teams {
		if team.ParentID == teamID {
			children = append(children, team)
		}
	}
	return children
}

// GetDescendants returns all descendant teams (recursive).
func (h *TeamHierarchy) GetDescendants(teamID string) []*Team {
	var descendants []*Team
	children := h.GetChildren(teamID)
	for _, child := range children {
		descendants = append(descendants, child)
		descendants = append(descendants, h.GetDescendants(child.ID)...)
	}
	return descendants
}

// GetAncestors returns all ancestor teams (up to root).
func (h *TeamHierarchy) GetAncestors(teamID string) []*Team {
	var ancestors []*Team
	parent := h.GetParent(teamID)
	for parent != nil {
		ancestors = append(ancestors, parent)
		parent = h.GetParent(parent.ID)
	}
	return ancestors
}

// GetByLevel returns all teams at a specific org level.
func (h *TeamHierarchy) GetByLevel(level TeamLevel) []*Team {
	var result []*Team
	for _, team := range h.Teams {
		if team.Level == level {
			result = append(result, team)
		}
	}
	return result
}

// AggregateCapacity sums capacity for a team and all descendants.
func (h *TeamHierarchy) AggregateCapacity(teamID string) *TeamCapacity {
	team, ok := h.Teams[teamID]
	if !ok {
		return nil
	}

	total := &TeamCapacity{}
	if team.Capacity != nil {
		total.FTEs = team.Capacity.FTEs
		total.StoryPointsPerSprint = team.Capacity.StoryPointsPerSprint
		total.StoryPointsPerQuarter = team.Capacity.StoryPointsPerQuarter
	}

	for _, desc := range h.GetDescendants(teamID) {
		if desc.Capacity != nil {
			total.FTEs += desc.Capacity.FTEs
			total.StoryPointsPerSprint += desc.Capacity.StoryPointsPerSprint
			total.StoryPointsPerQuarter += desc.Capacity.StoryPointsPerQuarter
		}
	}

	return total
}
