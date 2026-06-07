package canvas

import (
	"encoding/json"
	"errors"
)

// CanvasType identifies which canvas framework is in use.
type CanvasType string

// Canvas type constants.
const (
	CanvasTypeBMC                   CanvasType = "bmc"
	CanvasTypeOpportunity           CanvasType = "opportunity"
	CanvasTypeOpportunityAssessment CanvasType = "opportunity-assessment"
	CanvasTypeOpportunitySpec       CanvasType = "opportunity-spec"
	CanvasTypeFeature               CanvasType = "feature"
	CanvasTypeLeanUX                CanvasType = "leanux"
	CanvasTypeOST                   CanvasType = "ost"

	// Shape Up (Basecamp)
	CanvasTypeShapeUpPitch CanvasType = "shapeup-pitch"
	CanvasTypeShapeUpBet   CanvasType = "shapeup-bet"
	CanvasTypeShapeUpScope CanvasType = "shapeup-scope"

	// Continuous Discovery (Teresa Torres)
	CanvasTypeDiscoverySnapshot CanvasType = "discovery-snapshot"
	CanvasTypeAssumptionMap     CanvasType = "assumption-map"
	CanvasTypeExperienceMap     CanvasType = "experience-map"

	// Lean Startup (Eric Ries)
	CanvasTypeLeanStartup CanvasType = "leanstartup"

	// Design Thinking (Stanford d.school)
	CanvasTypeDesignThinking CanvasType = "designthinking"

	// Jobs-to-be-Done (Christensen, Ulwick)
	CanvasTypeJTBD CanvasType = "jtbd"
)

// Canvas is a framework-agnostic container for strategic planning canvases.
// This is a discriminated union - exactly one field is set based on Type.
type Canvas struct {
	Type                  CanvasType               `json:"type"`
	BMC                   *BusinessModelCanvas     `json:"bmc,omitempty"`
	Opportunity           *OpportunityCanvas       `json:"opportunity,omitempty"`
	OpportunityAssessment *OpportunityAssessment   `json:"opportunityAssessment,omitempty"`
	OpportunitySpec       *OpportunitySpec         `json:"opportunitySpec,omitempty"`
	Feature               *FeatureCanvas           `json:"feature,omitempty"`
	LeanUX                *LeanUXCanvas            `json:"leanux,omitempty"`
	OST                   *OpportunitySolutionTree `json:"ost,omitempty"`

	// Shape Up (Basecamp)
	ShapeUpPitch *ShapeUpPitch `json:"shapeupPitch,omitempty"`
	ShapeUpBet   *ShapeUpBet   `json:"shapeupBet,omitempty"`
	ShapeUpScope *ShapeUpScope `json:"shapeupScope,omitempty"`

	// Continuous Discovery (Teresa Torres)
	DiscoverySnapshot *DiscoverySnapshot `json:"discoverySnapshot,omitempty"`
	AssumptionMap     *AssumptionMap     `json:"assumptionMap,omitempty"`
	ExperienceMap     *ExperienceMap     `json:"experienceMap,omitempty"`

	// Lean Startup (Eric Ries)
	LeanStartup *LeanStartupCanvas `json:"leanStartup,omitempty"`

	// Design Thinking (Stanford d.school)
	DesignThinking *DesignThinkingCanvas `json:"designThinking,omitempty"`

	// Jobs-to-be-Done (Christensen, Ulwick)
	JTBD *JTBDCanvas `json:"jtbd,omitempty"`
}

// NewBMC creates a Canvas wrapping a BusinessModelCanvas.
func NewBMC(c *BusinessModelCanvas) *Canvas {
	return &Canvas{
		Type: CanvasTypeBMC,
		BMC:  c,
	}
}

// NewOpportunity creates a Canvas wrapping an OpportunityCanvas.
func NewOpportunity(c *OpportunityCanvas) *Canvas {
	return &Canvas{
		Type:        CanvasTypeOpportunity,
		Opportunity: c,
	}
}

// NewOpportunityAssessment creates a Canvas wrapping an OpportunityAssessment.
func NewOpportunityAssessmentCanvas(c *OpportunityAssessment) *Canvas {
	return &Canvas{
		Type:                  CanvasTypeOpportunityAssessment,
		OpportunityAssessment: c,
	}
}

// NewOpportunitySpecCanvas creates a Canvas wrapping an OpportunitySpec.
func NewOpportunitySpecCanvas(c *OpportunitySpec) *Canvas {
	return &Canvas{
		Type:            CanvasTypeOpportunitySpec,
		OpportunitySpec: c,
	}
}

// NewFeature creates a Canvas wrapping a FeatureCanvas.
func NewFeature(c *FeatureCanvas) *Canvas {
	return &Canvas{
		Type:    CanvasTypeFeature,
		Feature: c,
	}
}

// NewLeanUX creates a Canvas wrapping a LeanUXCanvas.
func NewLeanUX(c *LeanUXCanvas) *Canvas {
	return &Canvas{
		Type:   CanvasTypeLeanUX,
		LeanUX: c,
	}
}

// NewOST creates a Canvas wrapping an OpportunitySolutionTree.
func NewOST(c *OpportunitySolutionTree) *Canvas {
	return &Canvas{
		Type: CanvasTypeOST,
		OST:  c,
	}
}

// NewShapeUpPitchCanvas creates a Canvas wrapping a ShapeUpPitch.
func NewShapeUpPitchCanvas(c *ShapeUpPitch) *Canvas {
	return &Canvas{
		Type:         CanvasTypeShapeUpPitch,
		ShapeUpPitch: c,
	}
}

// NewShapeUpBetCanvas creates a Canvas wrapping a ShapeUpBet.
func NewShapeUpBetCanvas(c *ShapeUpBet) *Canvas {
	return &Canvas{
		Type:       CanvasTypeShapeUpBet,
		ShapeUpBet: c,
	}
}

// NewShapeUpScopeCanvas creates a Canvas wrapping a ShapeUpScope.
func NewShapeUpScopeCanvas(c *ShapeUpScope) *Canvas {
	return &Canvas{
		Type:         CanvasTypeShapeUpScope,
		ShapeUpScope: c,
	}
}

// NewDiscoverySnapshotCanvas creates a Canvas wrapping a DiscoverySnapshot.
func NewDiscoverySnapshotCanvas(c *DiscoverySnapshot) *Canvas {
	return &Canvas{
		Type:              CanvasTypeDiscoverySnapshot,
		DiscoverySnapshot: c,
	}
}

// NewAssumptionMapCanvas creates a Canvas wrapping an AssumptionMap.
func NewAssumptionMapCanvas(c *AssumptionMap) *Canvas {
	return &Canvas{
		Type:          CanvasTypeAssumptionMap,
		AssumptionMap: c,
	}
}

// NewExperienceMapCanvas creates a Canvas wrapping an ExperienceMap.
func NewExperienceMapCanvas(c *ExperienceMap) *Canvas {
	return &Canvas{
		Type:          CanvasTypeExperienceMap,
		ExperienceMap: c,
	}
}

// NewLeanStartupCanvas creates a Canvas wrapping a LeanStartupCanvas.
func NewLeanStartupCanvasWrapper(c *LeanStartupCanvas) *Canvas {
	return &Canvas{
		Type:        CanvasTypeLeanStartup,
		LeanStartup: c,
	}
}

// NewDesignThinkingCanvasWrapper creates a Canvas wrapping a DesignThinkingCanvas.
func NewDesignThinkingCanvasWrapper(c *DesignThinkingCanvas) *Canvas {
	return &Canvas{
		Type:           CanvasTypeDesignThinking,
		DesignThinking: c,
	}
}

// NewJTBDCanvasWrapper creates a Canvas wrapping a JTBDCanvas.
func NewJTBDCanvasWrapper(c *JTBDCanvas) *Canvas {
	return &Canvas{
		Type: CanvasTypeJTBD,
		JTBD: c,
	}
}

// IsBMC returns true if this canvas contains a Business Model Canvas.
func (c *Canvas) IsBMC() bool {
	return c != nil && c.Type == CanvasTypeBMC
}

// IsOpportunity returns true if this canvas contains an Opportunity Canvas.
func (c *Canvas) IsOpportunity() bool {
	return c != nil && c.Type == CanvasTypeOpportunity
}

// IsOpportunityAssessment returns true if this canvas contains an Opportunity Assessment.
func (c *Canvas) IsOpportunityAssessment() bool {
	return c != nil && c.Type == CanvasTypeOpportunityAssessment
}

// IsOpportunitySpec returns true if this canvas contains an OpportunitySpec.
func (c *Canvas) IsOpportunitySpec() bool {
	return c != nil && c.Type == CanvasTypeOpportunitySpec
}

// IsFeature returns true if this canvas contains a Feature Canvas.
func (c *Canvas) IsFeature() bool {
	return c != nil && c.Type == CanvasTypeFeature
}

// IsLeanUX returns true if this canvas contains a Lean UX Canvas.
func (c *Canvas) IsLeanUX() bool {
	return c != nil && c.Type == CanvasTypeLeanUX
}

// IsOST returns true if this canvas contains an Opportunity Solution Tree.
func (c *Canvas) IsOST() bool {
	return c != nil && c.Type == CanvasTypeOST
}

// IsShapeUpPitch returns true if this canvas contains a ShapeUpPitch.
func (c *Canvas) IsShapeUpPitch() bool {
	return c != nil && c.Type == CanvasTypeShapeUpPitch
}

// IsShapeUpBet returns true if this canvas contains a ShapeUpBet.
func (c *Canvas) IsShapeUpBet() bool {
	return c != nil && c.Type == CanvasTypeShapeUpBet
}

// IsShapeUpScope returns true if this canvas contains a ShapeUpScope.
func (c *Canvas) IsShapeUpScope() bool {
	return c != nil && c.Type == CanvasTypeShapeUpScope
}

// IsDiscoverySnapshot returns true if this canvas contains a DiscoverySnapshot.
func (c *Canvas) IsDiscoverySnapshot() bool {
	return c != nil && c.Type == CanvasTypeDiscoverySnapshot
}

// IsAssumptionMap returns true if this canvas contains an AssumptionMap.
func (c *Canvas) IsAssumptionMap() bool {
	return c != nil && c.Type == CanvasTypeAssumptionMap
}

// IsExperienceMap returns true if this canvas contains an ExperienceMap.
func (c *Canvas) IsExperienceMap() bool {
	return c != nil && c.Type == CanvasTypeExperienceMap
}

// IsLeanStartup returns true if this canvas contains a LeanStartupCanvas.
func (c *Canvas) IsLeanStartup() bool {
	return c != nil && c.Type == CanvasTypeLeanStartup
}

// IsDesignThinking returns true if this canvas contains a DesignThinkingCanvas.
func (c *Canvas) IsDesignThinking() bool {
	return c != nil && c.Type == CanvasTypeDesignThinking
}

// IsJTBD returns true if this canvas contains a JTBDCanvas.
func (c *Canvas) IsJTBD() bool {
	return c != nil && c.Type == CanvasTypeJTBD
}

// GetMetadata returns the metadata from the inner canvas, if present.
func (c *Canvas) GetMetadata() *Metadata {
	if c == nil {
		return nil
	}
	switch c.Type {
	case CanvasTypeBMC:
		if c.BMC != nil {
			return &c.BMC.Metadata
		}
	case CanvasTypeOpportunity:
		if c.Opportunity != nil {
			return &c.Opportunity.Metadata
		}
	case CanvasTypeOpportunityAssessment:
		if c.OpportunityAssessment != nil {
			return &c.OpportunityAssessment.Metadata
		}
	case CanvasTypeOpportunitySpec:
		if c.OpportunitySpec != nil {
			return &c.OpportunitySpec.Metadata
		}
	case CanvasTypeFeature:
		if c.Feature != nil {
			return &c.Feature.Metadata
		}
	case CanvasTypeLeanUX:
		if c.LeanUX != nil {
			return &c.LeanUX.Metadata
		}
	case CanvasTypeOST:
		if c.OST != nil {
			return &c.OST.Metadata
		}
	case CanvasTypeShapeUpPitch:
		if c.ShapeUpPitch != nil {
			return &c.ShapeUpPitch.Metadata
		}
	case CanvasTypeShapeUpBet:
		if c.ShapeUpBet != nil {
			return &c.ShapeUpBet.Metadata
		}
	case CanvasTypeShapeUpScope:
		if c.ShapeUpScope != nil {
			return &c.ShapeUpScope.Metadata
		}
	case CanvasTypeDiscoverySnapshot:
		if c.DiscoverySnapshot != nil {
			return &c.DiscoverySnapshot.Metadata
		}
	case CanvasTypeAssumptionMap:
		if c.AssumptionMap != nil {
			return &c.AssumptionMap.Metadata
		}
	case CanvasTypeExperienceMap:
		if c.ExperienceMap != nil {
			return &c.ExperienceMap.Metadata
		}
	case CanvasTypeLeanStartup:
		if c.LeanStartup != nil {
			return &c.LeanStartup.Metadata
		}
	case CanvasTypeDesignThinking:
		if c.DesignThinking != nil {
			return &c.DesignThinking.Metadata
		}
	case CanvasTypeJTBD:
		if c.JTBD != nil {
			return &c.JTBD.Metadata
		}
	}
	return nil
}

// GetInnerCanvas returns the inner canvas as an interface{}.
func (c *Canvas) GetInnerCanvas() any {
	if c == nil {
		return nil
	}
	switch c.Type {
	case CanvasTypeBMC:
		return c.BMC
	case CanvasTypeOpportunity:
		return c.Opportunity
	case CanvasTypeOpportunityAssessment:
		return c.OpportunityAssessment
	case CanvasTypeOpportunitySpec:
		return c.OpportunitySpec
	case CanvasTypeFeature:
		return c.Feature
	case CanvasTypeLeanUX:
		return c.LeanUX
	case CanvasTypeOST:
		return c.OST
	case CanvasTypeShapeUpPitch:
		return c.ShapeUpPitch
	case CanvasTypeShapeUpBet:
		return c.ShapeUpBet
	case CanvasTypeShapeUpScope:
		return c.ShapeUpScope
	case CanvasTypeDiscoverySnapshot:
		return c.DiscoverySnapshot
	case CanvasTypeAssumptionMap:
		return c.AssumptionMap
	case CanvasTypeExperienceMap:
		return c.ExperienceMap
	case CanvasTypeLeanStartup:
		return c.LeanStartup
	case CanvasTypeDesignThinking:
		return c.DesignThinking
	case CanvasTypeJTBD:
		return c.JTBD
	}
	return nil
}

// Validate checks that the canvas has valid structure.
func (c *Canvas) Validate() error {
	if c == nil {
		return errors.New("canvas is nil")
	}
	if c.Type == "" {
		return errors.New("canvas type is required")
	}

	// Check that exactly one inner canvas is set
	count := 0
	if c.BMC != nil {
		count++
	}
	if c.Opportunity != nil {
		count++
	}
	if c.OpportunityAssessment != nil {
		count++
	}
	if c.OpportunitySpec != nil {
		count++
	}
	if c.Feature != nil {
		count++
	}
	if c.LeanUX != nil {
		count++
	}
	if c.OST != nil {
		count++
	}
	if c.ShapeUpPitch != nil {
		count++
	}
	if c.ShapeUpBet != nil {
		count++
	}
	if c.ShapeUpScope != nil {
		count++
	}
	if c.DiscoverySnapshot != nil {
		count++
	}
	if c.AssumptionMap != nil {
		count++
	}
	if c.ExperienceMap != nil {
		count++
	}
	if c.LeanStartup != nil {
		count++
	}
	if c.DesignThinking != nil {
		count++
	}
	if c.JTBD != nil {
		count++
	}

	if count == 0 {
		return errors.New("no inner canvas set")
	}
	if count > 1 {
		return errors.New("multiple inner canvases set; exactly one required")
	}

	// Verify type matches inner canvas
	switch c.Type {
	case CanvasTypeBMC:
		if c.BMC == nil {
			return errors.New("type is 'bmc' but BMC field is nil")
		}
	case CanvasTypeOpportunity:
		if c.Opportunity == nil {
			return errors.New("type is 'opportunity' but Opportunity field is nil")
		}
	case CanvasTypeOpportunityAssessment:
		if c.OpportunityAssessment == nil {
			return errors.New("type is 'opportunity-assessment' but OpportunityAssessment field is nil")
		}
	case CanvasTypeOpportunitySpec:
		if c.OpportunitySpec == nil {
			return errors.New("type is 'opportunity-spec' but OpportunitySpec field is nil")
		}
	case CanvasTypeFeature:
		if c.Feature == nil {
			return errors.New("type is 'feature' but Feature field is nil")
		}
	case CanvasTypeLeanUX:
		if c.LeanUX == nil {
			return errors.New("type is 'leanux' but LeanUX field is nil")
		}
	case CanvasTypeOST:
		if c.OST == nil {
			return errors.New("type is 'ost' but OST field is nil")
		}
	case CanvasTypeShapeUpPitch:
		if c.ShapeUpPitch == nil {
			return errors.New("type is 'shapeup-pitch' but ShapeUpPitch field is nil")
		}
	case CanvasTypeShapeUpBet:
		if c.ShapeUpBet == nil {
			return errors.New("type is 'shapeup-bet' but ShapeUpBet field is nil")
		}
	case CanvasTypeShapeUpScope:
		if c.ShapeUpScope == nil {
			return errors.New("type is 'shapeup-scope' but ShapeUpScope field is nil")
		}
	case CanvasTypeDiscoverySnapshot:
		if c.DiscoverySnapshot == nil {
			return errors.New("type is 'discovery-snapshot' but DiscoverySnapshot field is nil")
		}
	case CanvasTypeAssumptionMap:
		if c.AssumptionMap == nil {
			return errors.New("type is 'assumption-map' but AssumptionMap field is nil")
		}
	case CanvasTypeExperienceMap:
		if c.ExperienceMap == nil {
			return errors.New("type is 'experience-map' but ExperienceMap field is nil")
		}
	case CanvasTypeLeanStartup:
		if c.LeanStartup == nil {
			return errors.New("type is 'leanstartup' but LeanStartup field is nil")
		}
	case CanvasTypeDesignThinking:
		if c.DesignThinking == nil {
			return errors.New("type is 'designthinking' but DesignThinking field is nil")
		}
	case CanvasTypeJTBD:
		if c.JTBD == nil {
			return errors.New("type is 'jtbd' but JTBD field is nil")
		}
	default:
		return errors.New("unknown canvas type: " + string(c.Type))
	}

	return nil
}

// MarshalJSON implements json.Marshaler.
func (c *Canvas) MarshalJSON() ([]byte, error) {
	type alias Canvas
	return json.Marshal((*alias)(c))
}

// UnmarshalJSON implements json.Unmarshaler.
func (c *Canvas) UnmarshalJSON(data []byte) error {
	type alias Canvas
	if err := json.Unmarshal(data, (*alias)(c)); err != nil {
		return err
	}
	return c.Validate()
}
