package journey

// Dependency represents a dependency between two entities.
// The "from" entity is blocked by or depends on the "to" entity.
// This follows the convention that the blocked entity declares/generates the dependency.
type Dependency struct {
	ID                 string           `json:"id,omitempty"`
	From               EntityRef        `json:"from"` // The blocked/dependent entity
	To                 EntityRef        `json:"to"`   // The blocking/required entity
	Type               DependencyType   `json:"type"`
	Description        string           `json:"description,omitempty"`
	Status             DependencyStatus `json:"status,omitempty"`
	Risk               DependencyRisk   `json:"risk,omitempty"`
	ExpectedResolution string           `json:"expectedResolution,omitempty"` // Period ID when expected resolved
	Notes              string           `json:"notes,omitempty"`
}

// DependencyType classifies the nature of the dependency.
type DependencyType string

const (
	// DependencyRequires indicates a technical or logical requirement.
	// "Multi-Agent requires Agent Runtime" - can't build without it.
	DependencyRequires DependencyType = "requires"

	// DependencyBlockedBy indicates sequencing/scheduling constraint.
	// "SSO rollout blocked by Identity Platform" - must wait for completion.
	DependencyBlockedBy DependencyType = "blocked_by"

	// DependencyResource indicates team capacity dependency.
	// "Feature X blocked by Platform Team capacity" - competing for resources.
	DependencyResource DependencyType = "resource"

	// DependencyExternal indicates external dependency (vendor, compliance, etc).
	// "Integration blocked by Vendor API availability"
	DependencyExternal DependencyType = "external"

	// DependencyInforms indicates a soft/informational dependency.
	// "Decision X informs approach for Feature Y" - not blocking, but influential.
	DependencyInforms DependencyType = "informs"

	// DependencyContributes indicates contribution to a shared goal.
	// "Team A's work contributes to Platform Maturity"
	DependencyContributes DependencyType = "contributes"
)

// DependencyStatus tracks the resolution state.
type DependencyStatus string

const (
	DependencyStatusPending  DependencyStatus = "pending"  // Not yet resolved
	DependencyStatusResolved DependencyStatus = "resolved" // Dependency satisfied
	DependencyStatusBlocked  DependencyStatus = "blocked"  // Actively blocking progress
	DependencyStatusAtRisk   DependencyStatus = "at_risk"  // May cause delays
	DependencyStatusWaived   DependencyStatus = "waived"   // Accepted without resolution
)

// DependencyRisk indicates the risk level of the dependency.
type DependencyRisk string

const (
	DependencyRiskLow      DependencyRisk = "low"
	DependencyRiskMedium   DependencyRisk = "medium"
	DependencyRiskHigh     DependencyRisk = "high"
	DependencyRiskCritical DependencyRisk = "critical"
)

// EntityRef references an entity in the roadmap.
type EntityRef struct {
	Type EntityType `json:"type"`
	ID   string     `json:"id"`
	Name string     `json:"name,omitempty"` // For display purposes
}

// EntityType identifies what kind of entity is referenced.
type EntityType string

const (
	EntityTypeCapability EntityType = "capability"
	EntityTypeInitiative EntityType = "initiative"
	EntityTypeTeam       EntityType = "team"
	EntityTypeMilestone  EntityType = "milestone"
	EntityTypeOutcome    EntityType = "outcome"
	EntityTypeExternal   EntityType = "external"
	EntityTypeDecision   EntityType = "decision"
)

// DependencyGraph represents all dependencies for analysis.
type DependencyGraph struct {
	Dependencies []Dependency `json:"dependencies"`
}

// GetBlockers returns all dependencies where the given entity is blocked.
func (g *DependencyGraph) GetBlockers(entityType EntityType, entityID string) []Dependency {
	var result []Dependency
	for _, dep := range g.Dependencies {
		if dep.From.Type == entityType && dep.From.ID == entityID {
			result = append(result, dep)
		}
	}
	return result
}

// GetBlocking returns all dependencies where the given entity is blocking others.
func (g *DependencyGraph) GetBlocking(entityType EntityType, entityID string) []Dependency {
	var result []Dependency
	for _, dep := range g.Dependencies {
		if dep.To.Type == entityType && dep.To.ID == entityID {
			result = append(result, dep)
		}
	}
	return result
}

// GetCriticalPath returns dependencies that are blocking and high/critical risk.
func (g *DependencyGraph) GetCriticalPath() []Dependency {
	var result []Dependency
	for _, dep := range g.Dependencies {
		if dep.Status == DependencyStatusBlocked &&
			(dep.Risk == DependencyRiskHigh || dep.Risk == DependencyRiskCritical) {
			result = append(result, dep)
		}
	}
	return result
}

// GetResourceDependencies returns all resource/capacity dependencies.
func (g *DependencyGraph) GetResourceDependencies() []Dependency {
	var result []Dependency
	for _, dep := range g.Dependencies {
		if dep.Type == DependencyResource {
			result = append(result, dep)
		}
	}
	return result
}

// GetTeamDependencies returns all dependencies involving a specific team.
func (g *DependencyGraph) GetTeamDependencies(teamID string) []Dependency {
	var result []Dependency
	for _, dep := range g.Dependencies {
		if (dep.From.Type == EntityTypeTeam && dep.From.ID == teamID) ||
			(dep.To.Type == EntityTypeTeam && dep.To.ID == teamID) {
			result = append(result, dep)
		}
	}
	return result
}
