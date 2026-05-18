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
	PRDSchemaID = "https://github.com/grokify/prism-execution/schema/prd.schema.json"

	// MRDSchemaID is the canonical ID for the MRD schema (placeholder).
	MRDSchemaID = "https://github.com/grokify/prism-execution/schema/mrd.schema.json"

	// TRDSchemaID is the canonical ID for the TRD schema (placeholder).
	TRDSchemaID = "https://github.com/grokify/prism-execution/schema/trd.schema.json"

	// OKRSchemaID is the canonical ID for the OKR schema.
	OKRSchemaID = "https://github.com/grokify/prism-execution/schema/okr.schema.json"

	// V2MOMSchemaID is the canonical ID for the V2MOM schema.
	V2MOMSchemaID = "https://github.com/grokify/prism-execution/schema/v2mom.schema.json"
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
