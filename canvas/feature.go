package canvas

// FeatureCanvas follows Nikita Efimov's structure for feature definition.
// Licensed under CC BY-SA 4.0.
// See: https://medium.com/@nikita_ixotec/featurecanvas-3a05e10bac90
type FeatureCanvas struct {
	Metadata      Metadata `json:"metadata"`
	IdeaStatement string   `json:"ideaStatement"` // Top banner - concise feature idea

	// Problem Area (left side)
	Situations   []Situation  `json:"situations"`   // When/where does problem occur?
	Problems     []Problem    `json:"problems"`     // What problems exist?
	Value        []ValueItem  `json:"value"`        // What value will this provide?
	Capabilities []Capability `json:"capabilities"` // What must the solution do?

	// Solution Area (right side)
	Restrictions []string `json:"restrictions"` // What can't we do? (business/legal)
	Limitations  []string `json:"limitations"`  // What are technical limitations?

	// PRD integration
	PRDRef *PRDReference `json:"prdRef,omitempty"`
}

// Situation describes when and where a problem occurs.
type Situation struct {
	ID          string `json:"id"`
	Description string `json:"description"`         // When/where does problem occur?
	Actor       string `json:"actor,omitempty"`     // Who experiences this?
	Context     string `json:"context,omitempty"`   // What context triggers it?
	Trigger     string `json:"trigger,omitempty"`   // What action/event triggers it?
	Frequency   string `json:"frequency,omitempty"` // How often does this happen?
}

// FeatureProblem describes a problem in the Feature Canvas context.
// Note: Reuses Problem from opportunity.go but adds feature-specific fields.
type FeatureProblem struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	AffectedUser string   `json:"affectedUser,omitempty"` // Who is affected?
	Impact       string   `json:"impact,omitempty"`       // What's the impact?
	SituationRef string   `json:"situationRef,omitempty"` // Links to Situation.ID
	Evidence     []string `json:"evidence,omitempty"`     // User feedback, data, etc.
}

// ValueItem describes the value a feature provides.
type ValueItem struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Beneficiary string `json:"beneficiary,omitempty"` // Who benefits? (user, business, both)
	Type        string `json:"type,omitempty"`        // functional, emotional, social
	Metric      string `json:"metric,omitempty"`      // How to measure?
	TargetValue string `json:"targetValue,omitempty"` // What's the target?
}

// Capability describes what the solution must do.
type Capability struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`            // What must the solution do?
	Priority     string   `json:"priority,omitempty"`     // MoSCoW: must, should, could, wont
	Rationale    string   `json:"rationale,omitempty"`    // Why is this needed?
	Dependencies []string `json:"dependencies,omitempty"` // Other capability IDs
}

// NewFeatureCanvas creates a new FeatureCanvas with defaults.
func NewFeatureCanvas(id, title string) *FeatureCanvas {
	return &FeatureCanvas{
		Metadata: Metadata{
			ID:      id,
			Title:   title,
			Version: VersionFeature1,
		},
	}
}

// GetPRDReference returns the PRD reference.
func (c *FeatureCanvas) GetPRDReference() *PRDReference {
	return c.PRDRef
}

// MustHaveCapabilities returns capabilities with "must" priority.
func (c *FeatureCanvas) MustHaveCapabilities() []Capability {
	var must []Capability
	for _, cap := range c.Capabilities {
		if cap.Priority == "must" {
			must = append(must, cap)
		}
	}
	return must
}

// ShouldHaveCapabilities returns capabilities with "should" priority.
func (c *FeatureCanvas) ShouldHaveCapabilities() []Capability {
	var should []Capability
	for _, cap := range c.Capabilities {
		if cap.Priority == "should" {
			should = append(should, cap)
		}
	}
	return should
}

// HasRestrictions returns true if there are restrictions defined.
func (c *FeatureCanvas) HasRestrictions() bool {
	return len(c.Restrictions) > 0
}

// HasLimitations returns true if there are limitations defined.
func (c *FeatureCanvas) HasLimitations() bool {
	return len(c.Limitations) > 0
}
