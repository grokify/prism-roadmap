// Package schema provides embedded JSON Schema files for structured requirements documents.
// These schemas can be used for validation and documentation purposes.
package schema

import (
	_ "embed"
)

// PRD Schema

//go:embed prd.schema.json
var PRDSchemaJSON []byte

// PRDSchema returns the PRD JSON Schema as a string.
func PRDSchema() string {
	return string(PRDSchemaJSON)
}

// PRDSchemaBytes returns the PRD JSON Schema as a byte slice.
func PRDSchemaBytes() []byte {
	return PRDSchemaJSON
}

// SchemaID constants for referencing schemas.
const (
	// PRDSchemaID is the canonical ID for the PRD schema.
	PRDSchemaID = "https://github.com/grokify/prism-roadmap/schema/prd.schema.json"

	// MRDSchemaID is the canonical ID for the MRD schema (placeholder).
	MRDSchemaID = "https://github.com/grokify/prism-roadmap/schema/mrd.schema.json"

	// TRDSchemaID is the canonical ID for the TRD schema (placeholder).
	TRDSchemaID = "https://github.com/grokify/prism-roadmap/schema/trd.schema.json"

	// OKRSchemaID is the canonical ID for the OKR schema.
	OKRSchemaID = "https://github.com/grokify/prism-roadmap/schema/okr.schema.json"

	// V2MOMSchemaID is the canonical ID for the V2MOM schema.
	V2MOMSchemaID = "https://github.com/grokify/prism-roadmap/schema/v2mom.schema.json"

	// Canvas schema IDs

	// CanvasSchemaID is the canonical ID for the Canvas wrapper schema.
	CanvasSchemaID = "https://github.com/grokify/prism-roadmap/schema/canvas.schema.json"

	// BMCSchemaID is the canonical ID for the Business Model Canvas schema.
	BMCSchemaID = "https://github.com/grokify/prism-roadmap/schema/bmc.schema.json"

	// OpportunitySchemaID is the canonical ID for the Opportunity Canvas schema.
	OpportunitySchemaID = "https://github.com/grokify/prism-roadmap/schema/opportunity.schema.json"

	// FeatureSchemaID is the canonical ID for the Feature Canvas schema.
	FeatureSchemaID = "https://github.com/grokify/prism-roadmap/schema/feature.schema.json"

	// LeanUXSchemaID is the canonical ID for the Lean UX Canvas schema.
	LeanUXSchemaID = "https://github.com/grokify/prism-roadmap/schema/leanux.schema.json"

	// OSTSchemaID is the canonical ID for the Opportunity Solution Tree schema.
	OSTSchemaID = "https://github.com/grokify/prism-roadmap/schema/ost.schema.json"

	// Shape Up canvas schema IDs

	// ShapeUpPitchSchemaID is the canonical ID for the Shape Up Pitch schema.
	ShapeUpPitchSchemaID = "https://github.com/grokify/prism-roadmap/schema/shapeup-pitch.schema.json"

	// ShapeUpBetSchemaID is the canonical ID for the Shape Up Bet schema.
	ShapeUpBetSchemaID = "https://github.com/grokify/prism-roadmap/schema/shapeup-bet.schema.json"

	// ShapeUpScopeSchemaID is the canonical ID for the Shape Up Scope schema.
	ShapeUpScopeSchemaID = "https://github.com/grokify/prism-roadmap/schema/shapeup-scope.schema.json"

	// Continuous Discovery canvas schema IDs

	// DiscoverySnapshotSchemaID is the canonical ID for the Discovery Snapshot schema.
	DiscoverySnapshotSchemaID = "https://github.com/grokify/prism-roadmap/schema/discovery-snapshot.schema.json"

	// AssumptionMapSchemaID is the canonical ID for the Assumption Map schema.
	AssumptionMapSchemaID = "https://github.com/grokify/prism-roadmap/schema/assumption-map.schema.json"

	// ExperienceMapSchemaID is the canonical ID for the Experience Map schema.
	ExperienceMapSchemaID = "https://github.com/grokify/prism-roadmap/schema/experience-map.schema.json"

	// Lean Startup schema ID

	// LeanStartupSchemaID is the canonical ID for the Lean Startup schema.
	LeanStartupSchemaID = "https://github.com/grokify/prism-roadmap/schema/leanstartup.schema.json"

	// Design Thinking schema ID

	// DesignThinkingSchemaID is the canonical ID for the Design Thinking schema.
	DesignThinkingSchemaID = "https://github.com/grokify/prism-roadmap/schema/designthinking.schema.json"

	// Jobs-to-be-Done schema ID

	// JTBDSchemaID is the canonical ID for the JTBD schema.
	JTBDSchemaID = "https://github.com/grokify/prism-roadmap/schema/jtbd.schema.json"

	// Journey Roadmap schema ID

	// JourneyRoadmapSchemaID is the canonical ID for the Journey Roadmap schema.
	JourneyRoadmapSchemaID = "https://github.com/grokify/prism-roadmap/schema/journey-roadmap.schema.json"
)

// OKR Schema

//go:embed okr.schema.json
var OKRSchemaJSON []byte

// OKRSchema returns the OKR JSON Schema as a string.
func OKRSchema() string {
	return string(OKRSchemaJSON)
}

// OKRSchemaBytes returns the OKR JSON Schema as a byte slice.
func OKRSchemaBytes() []byte {
	return OKRSchemaJSON
}

// V2MOM Schema

//go:embed v2mom.schema.json
var V2MOMSchemaJSON []byte

// V2MOMSchema returns the V2MOM JSON Schema as a string.
func V2MOMSchema() string {
	return string(V2MOMSchemaJSON)
}

// V2MOMSchemaBytes returns the V2MOM JSON Schema as a byte slice.
func V2MOMSchemaBytes() []byte {
	return V2MOMSchemaJSON
}

// Canvas Schemas
// Note: Schema JSON files are generated from Go types.
// Run `go generate ./schema/...` to regenerate.

// TODO: Add canvas schema embeds when generated:
//
//	//go:embed canvas.schema.json
//	var CanvasSchemaJSON []byte
//
//	//go:embed bmc.schema.json
//	var BMCSchemaJSON []byte
//
//	//go:embed opportunity.schema.json
//	var OpportunitySchemaJSON []byte
//
//	//go:embed feature.schema.json
//	var FeatureSchemaJSON []byte
//
//	//go:embed leanux.schema.json
//	var LeanUXSchemaJSON []byte
//
//	//go:embed ost.schema.json
//	var OSTSchemaJSON []byte

// TODO: Add MRD and TRD schemas when created.
// When mrd.schema.json is added:
//
//	//go:embed mrd.schema.json
//	var MRDSchemaJSON []byte
//
//	func MRDSchema() string { return string(MRDSchemaJSON) }
//
// When trd.schema.json is added:
//
//	//go:embed trd.schema.json
//	var TRDSchemaJSON []byte
//
//	func TRDSchema() string { return string(TRDSchemaJSON) }

// Journey Roadmap Schema

//go:embed journey-roadmap.schema.json
var JourneyRoadmapSchemaJSON []byte

// JourneyRoadmapSchema returns the Journey Roadmap JSON Schema as a string.
func JourneyRoadmapSchema() string {
	return string(JourneyRoadmapSchemaJSON)
}

// JourneyRoadmapSchemaBytes returns the Journey Roadmap JSON Schema as a byte slice.
func JourneyRoadmapSchemaBytes() []byte {
	return JourneyRoadmapSchemaJSON
}
