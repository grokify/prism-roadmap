// Package effort provides effort estimation and complexity tracking types.
package effort

// ComplexityFactors captures implementation complexity signals for a roadmap item.
// These factors help identify items that may require additional planning or resources.
type ComplexityFactors struct {
	// Boolean complexity indicators
	NewArchitecture bool `json:"new_architecture"` // Requires new architectural patterns
	NewDesignUX     bool `json:"new_design_ux"`    // Requires new design/UX work
	NewBillingSKU   bool `json:"new_billing_sku"`  // Requires new billing/pricing SKU

	// External dependencies
	Dependencies []Dependency `json:"dependencies,omitempty"`

	// Additional factors (extensible)
	NewIntegration bool `json:"new_integration,omitempty"` // Requires new third-party integration
	NewCompliance  bool `json:"new_compliance,omitempty"`  // Requires compliance/security review
	NewDataModel   bool `json:"new_data_model,omitempty"`  // Requires data model changes

	// Notes
	Notes string `json:"notes,omitempty"`
}

// ComplexityScore calculates a numeric complexity penalty.
// Higher score = more complex.
// Formula: (new_arch × 2) + (new_ux × 1) + (new_sku × 1.5) + (dep_count × 0.5) + extensions
func (c *ComplexityFactors) ComplexityScore() float64 {
	score := 0.0
	if c.NewArchitecture {
		score += 2.0
	}
	if c.NewDesignUX {
		score += 1.0
	}
	if c.NewBillingSKU {
		score += 1.5
	}
	if c.NewIntegration {
		score += 1.0
	}
	if c.NewCompliance {
		score += 1.0
	}
	if c.NewDataModel {
		score += 1.0
	}
	score += float64(len(c.Dependencies)) * 0.5
	return score
}

// HasComplexity returns true if any complexity factors are set.
func (c *ComplexityFactors) HasComplexity() bool {
	return c.NewArchitecture ||
		c.NewDesignUX ||
		c.NewBillingSKU ||
		c.NewIntegration ||
		c.NewCompliance ||
		c.NewDataModel ||
		len(c.Dependencies) > 0
}

// BlockingDependencies returns dependencies that are blocking.
func (c *ComplexityFactors) BlockingDependencies() []Dependency {
	var result []Dependency
	for _, d := range c.Dependencies {
		if d.Type == DependencyTypeBlocking {
			result = append(result, d)
		}
	}
	return result
}

// TotalDependencyDays returns the total estimated days for all dependencies.
func (c *ComplexityFactors) TotalDependencyDays() int {
	total := 0
	for _, d := range c.Dependencies {
		total += d.EstimatedDays
	}
	return total
}

// Dependency represents a cross-team or external dependency.
type Dependency struct {
	TeamID        string           `json:"team_id"`                  // Unique team identifier
	TeamName      string           `json:"team_name,omitempty"`      // Human-readable team name
	Type          DependencyType   `json:"type"`                     // blocking, informational
	Description   string           `json:"description,omitempty"`    // What is needed
	EstimatedDays int              `json:"estimated_days,omitempty"` // Estimated resolution time
	Status        DependencyStatus `json:"status,omitempty"`         // pending, confirmed, resolved
	Owner         string           `json:"owner,omitempty"`          // Point of contact
	Notes         string           `json:"notes,omitempty"`
}

// DependencyType indicates how the dependency affects the item.
type DependencyType string

const (
	DependencyTypeBlocking      DependencyType = "blocking"      // Cannot proceed without resolution
	DependencyTypeInformational DependencyType = "informational" // FYI, no blocking impact
)

// DependencyStatus indicates the current state of the dependency.
type DependencyStatus string

const (
	DependencyStatusPending    DependencyStatus = "pending"     // Not yet acknowledged
	DependencyStatusConfirmed  DependencyStatus = "confirmed"   // Acknowledged by dependent team
	DependencyStatusInProgress DependencyStatus = "in_progress" // Being worked on
	DependencyStatusResolved   DependencyStatus = "resolved"    // Completed
	DependencyStatusBlocked    DependencyStatus = "blocked"     // Blocked on something else
)

// IsResolved returns true if the dependency is resolved.
func (d *Dependency) IsResolved() bool {
	return d.Status == DependencyStatusResolved
}

// IsBlocking returns true if the dependency is blocking and not resolved.
func (d *Dependency) IsBlocking() bool {
	return d.Type == DependencyTypeBlocking && d.Status != DependencyStatusResolved
}

// NewDependency creates a new dependency.
func NewDependency(teamID string, depType DependencyType) *Dependency {
	return &Dependency{
		TeamID: teamID,
		Type:   depType,
		Status: DependencyStatusPending,
	}
}

// NewBlockingDependency creates a new blocking dependency.
func NewBlockingDependency(teamID, teamName string, estimatedDays int) *Dependency {
	return &Dependency{
		TeamID:        teamID,
		TeamName:      teamName,
		Type:          DependencyTypeBlocking,
		EstimatedDays: estimatedDays,
		Status:        DependencyStatusPending,
	}
}
