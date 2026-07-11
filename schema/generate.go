package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/goals/okr"
	"github.com/grokify/prism-roadmap/goals/v2mom"
	"github.com/grokify/prism-roadmap/journey"
	"github.com/grokify/prism-roadmap/requirements/prd"
)

// Generator creates JSON Schema files from Go types.
type Generator struct {
	// Reflector is the jsonschema reflector used for generation.
	Reflector *jsonschema.Reflector
}

// NewGenerator creates a new schema generator with default settings.
func NewGenerator() *Generator {
	r := &jsonschema.Reflector{
		DoNotReference:             false,
		ExpandedStruct:             false,
		RequiredFromJSONSchemaTags: true,
	}
	return &Generator{Reflector: r}
}

// GeneratePRDSchema generates JSON Schema for the PRD Document type.
func (g *Generator) GeneratePRDSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&prd.Document{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for prd.Document")
	}

	// Set schema metadata
	schema.ID = jsonschema.ID(PRDSchemaID)
	schema.Title = "Structured PRD"
	schema.Description = "Schema for structured Product Requirements Documents"

	return schema, nil
}

// GeneratePRDSchemaJSON generates JSON Schema for PRD and returns it as JSON bytes.
func (g *Generator) GeneratePRDSchemaJSON() ([]byte, error) {
	schema, err := g.GeneratePRDSchema()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(schema, "", "  ")
}

// WritePRDSchema generates and writes the PRD schema to a file.
func (g *Generator) WritePRDSchema(path string) error {
	data, err := g.GeneratePRDSchemaJSON()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// GenerateAll generates all schema files to the specified directory.
func (g *Generator) GenerateAll(dir string) error {
	// Generate PRD schema
	prdPath := filepath.Join(dir, "prd.schema.json")
	if err := g.WritePRDSchema(prdPath); err != nil {
		return fmt.Errorf("generating PRD schema: %w", err)
	}

	// Generate OKR schema
	okrPath := filepath.Join(dir, "okr.schema.json")
	if err := g.WriteOKRSchema(okrPath); err != nil {
		return fmt.Errorf("generating OKR schema: %w", err)
	}

	// Generate V2MOM schema
	v2momPath := filepath.Join(dir, "v2mom.schema.json")
	if err := g.WriteV2MOMSchema(v2momPath); err != nil {
		return fmt.Errorf("generating V2MOM schema: %w", err)
	}

	// Generate canvas schemas (Shape Up, Continuous Discovery)
	if err := g.GenerateCanvasSchemas(dir); err != nil {
		return fmt.Errorf("generating canvas schemas: %w", err)
	}

	// Generate Journey Roadmap schema
	journeyPath := filepath.Join(dir, "journey-roadmap.schema.json")
	if err := g.WriteJourneyRoadmapSchema(journeyPath); err != nil {
		return fmt.Errorf("generating Journey Roadmap schema: %w", err)
	}

	// TODO: Add MRD and TRD schema generation when types are ready
	// mrdPath := filepath.Join(dir, "mrd.schema.json")
	// trdPath := filepath.Join(dir, "trd.schema.json")

	return nil
}

// GenerateOKRSchema generates JSON Schema for the OKR Document type.
func (g *Generator) GenerateOKRSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&okr.OKRDocument{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for okr.OKRDocument")
	}

	// Set schema metadata
	schema.ID = jsonschema.ID(OKRSchemaID)
	schema.Title = "OKR Document"
	schema.Description = "Schema for OKR (Objectives and Key Results) documents"

	return schema, nil
}

// GenerateOKRSchemaJSON generates JSON Schema for OKR and returns it as JSON bytes.
func (g *Generator) GenerateOKRSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateOKRSchema()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(schema, "", "  ")
}

// WriteOKRSchema generates and writes the OKR schema to a file.
func (g *Generator) WriteOKRSchema(path string) error {
	data, err := g.GenerateOKRSchemaJSON()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// GenerateV2MOMSchema generates JSON Schema for the V2MOM Document type.
func (g *Generator) GenerateV2MOMSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&v2mom.V2MOM{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for v2mom.V2MOM")
	}

	// Set schema metadata
	schema.ID = jsonschema.ID(V2MOMSchemaID)
	schema.Title = "V2MOM Document"
	schema.Description = "Schema for V2MOM (Vision, Values, Methods, Obstacles, Measures) documents"

	return schema, nil
}

// GenerateV2MOMSchemaJSON generates JSON Schema for V2MOM and returns it as JSON bytes.
func (g *Generator) GenerateV2MOMSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateV2MOMSchema()
	if err != nil {
		return nil, err
	}

	return json.MarshalIndent(schema, "", "  ")
}

// WriteV2MOMSchema generates and writes the V2MOM schema to a file.
func (g *Generator) WriteV2MOMSchema(path string) error {
	data, err := g.GenerateV2MOMSchemaJSON()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// Shape Up Canvas Schema Generation

// GenerateShapeUpPitchSchema generates JSON Schema for the ShapeUpPitch type.
func (g *Generator) GenerateShapeUpPitchSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.ShapeUpPitch{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.ShapeUpPitch")
	}

	schema.ID = jsonschema.ID(ShapeUpPitchSchemaID)
	schema.Title = "Shape Up Pitch"
	schema.Description = "Schema for Shape Up Pitch canvas (Basecamp methodology)"

	return schema, nil
}

// GenerateShapeUpPitchSchemaJSON generates JSON Schema for ShapeUpPitch and returns it as JSON bytes.
func (g *Generator) GenerateShapeUpPitchSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateShapeUpPitchSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteShapeUpPitchSchema generates and writes the ShapeUpPitch schema to a file.
func (g *Generator) WriteShapeUpPitchSchema(path string) error {
	return g.writeSchema(path, g.GenerateShapeUpPitchSchemaJSON)
}

// GenerateShapeUpBetSchema generates JSON Schema for the ShapeUpBet type.
func (g *Generator) GenerateShapeUpBetSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.ShapeUpBet{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.ShapeUpBet")
	}

	schema.ID = jsonschema.ID(ShapeUpBetSchemaID)
	schema.Title = "Shape Up Bet"
	schema.Description = "Schema for Shape Up Bet decisions (Basecamp methodology)"

	return schema, nil
}

// GenerateShapeUpBetSchemaJSON generates JSON Schema for ShapeUpBet and returns it as JSON bytes.
func (g *Generator) GenerateShapeUpBetSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateShapeUpBetSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteShapeUpBetSchema generates and writes the ShapeUpBet schema to a file.
func (g *Generator) WriteShapeUpBetSchema(path string) error {
	return g.writeSchema(path, g.GenerateShapeUpBetSchemaJSON)
}

// GenerateShapeUpScopeSchema generates JSON Schema for the ShapeUpScope type.
func (g *Generator) GenerateShapeUpScopeSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.ShapeUpScope{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.ShapeUpScope")
	}

	schema.ID = jsonschema.ID(ShapeUpScopeSchemaID)
	schema.Title = "Shape Up Scope"
	schema.Description = "Schema for Shape Up Scope with hill chart tracking (Basecamp methodology)"

	return schema, nil
}

// GenerateShapeUpScopeSchemaJSON generates JSON Schema for ShapeUpScope and returns it as JSON bytes.
func (g *Generator) GenerateShapeUpScopeSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateShapeUpScopeSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteShapeUpScopeSchema generates and writes the ShapeUpScope schema to a file.
func (g *Generator) WriteShapeUpScopeSchema(path string) error {
	return g.writeSchema(path, g.GenerateShapeUpScopeSchemaJSON)
}

// Continuous Discovery Canvas Schema Generation

// GenerateDiscoverySnapshotSchema generates JSON Schema for the DiscoverySnapshot type.
func (g *Generator) GenerateDiscoverySnapshotSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.DiscoverySnapshot{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.DiscoverySnapshot")
	}

	schema.ID = jsonschema.ID(DiscoverySnapshotSchemaID)
	schema.Title = "Discovery Snapshot"
	schema.Description = "Schema for weekly Continuous Discovery snapshot (Teresa Torres methodology)"

	return schema, nil
}

// GenerateDiscoverySnapshotSchemaJSON generates JSON Schema for DiscoverySnapshot and returns it as JSON bytes.
func (g *Generator) GenerateDiscoverySnapshotSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateDiscoverySnapshotSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteDiscoverySnapshotSchema generates and writes the DiscoverySnapshot schema to a file.
func (g *Generator) WriteDiscoverySnapshotSchema(path string) error {
	return g.writeSchema(path, g.GenerateDiscoverySnapshotSchemaJSON)
}

// GenerateAssumptionMapSchema generates JSON Schema for the AssumptionMap type.
func (g *Generator) GenerateAssumptionMapSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.AssumptionMap{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.AssumptionMap")
	}

	schema.ID = jsonschema.ID(AssumptionMapSchemaID)
	schema.Title = "Assumption Map"
	schema.Description = "Schema for Assumption Map (DVFUE matrix) from Continuous Discovery"

	return schema, nil
}

// GenerateAssumptionMapSchemaJSON generates JSON Schema for AssumptionMap and returns it as JSON bytes.
func (g *Generator) GenerateAssumptionMapSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateAssumptionMapSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteAssumptionMapSchema generates and writes the AssumptionMap schema to a file.
func (g *Generator) WriteAssumptionMapSchema(path string) error {
	return g.writeSchema(path, g.GenerateAssumptionMapSchemaJSON)
}

// GenerateExperienceMapSchema generates JSON Schema for the ExperienceMap type.
func (g *Generator) GenerateExperienceMapSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.ExperienceMap{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.ExperienceMap")
	}

	schema.ID = jsonschema.ID(ExperienceMapSchemaID)
	schema.Title = "Experience Map"
	schema.Description = "Schema for Experience Map (customer journey) from Continuous Discovery"

	return schema, nil
}

// GenerateExperienceMapSchemaJSON generates JSON Schema for ExperienceMap and returns it as JSON bytes.
func (g *Generator) GenerateExperienceMapSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateExperienceMapSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteExperienceMapSchema generates and writes the ExperienceMap schema to a file.
func (g *Generator) WriteExperienceMapSchema(path string) error {
	return g.writeSchema(path, g.GenerateExperienceMapSchemaJSON)
}

// writeSchema is a helper to write schema data to a file.
func (g *Generator) writeSchema(path string, genFunc func() ([]byte, error)) error {
	data, err := genFunc()
	if err != nil {
		return fmt.Errorf("generating schema: %w", err)
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	return nil
}

// GenerateCanvasSchemas generates all canvas schemas to the specified directory.
func (g *Generator) GenerateCanvasSchemas(dir string) error {
	// Shape Up schemas
	if err := g.WriteShapeUpPitchSchema(filepath.Join(dir, "shapeup-pitch.schema.json")); err != nil {
		return fmt.Errorf("generating ShapeUpPitch schema: %w", err)
	}
	if err := g.WriteShapeUpBetSchema(filepath.Join(dir, "shapeup-bet.schema.json")); err != nil {
		return fmt.Errorf("generating ShapeUpBet schema: %w", err)
	}
	if err := g.WriteShapeUpScopeSchema(filepath.Join(dir, "shapeup-scope.schema.json")); err != nil {
		return fmt.Errorf("generating ShapeUpScope schema: %w", err)
	}

	// Continuous Discovery schemas
	if err := g.WriteDiscoverySnapshotSchema(filepath.Join(dir, "discovery-snapshot.schema.json")); err != nil {
		return fmt.Errorf("generating DiscoverySnapshot schema: %w", err)
	}
	if err := g.WriteAssumptionMapSchema(filepath.Join(dir, "assumption-map.schema.json")); err != nil {
		return fmt.Errorf("generating AssumptionMap schema: %w", err)
	}
	if err := g.WriteExperienceMapSchema(filepath.Join(dir, "experience-map.schema.json")); err != nil {
		return fmt.Errorf("generating ExperienceMap schema: %w", err)
	}

	// Lean Startup schema
	if err := g.WriteLeanStartupSchema(filepath.Join(dir, "leanstartup.schema.json")); err != nil {
		return fmt.Errorf("generating LeanStartup schema: %w", err)
	}

	// Design Thinking schema
	if err := g.WriteDesignThinkingSchema(filepath.Join(dir, "designthinking.schema.json")); err != nil {
		return fmt.Errorf("generating DesignThinking schema: %w", err)
	}

	// Jobs-to-be-Done schema
	if err := g.WriteJTBDSchema(filepath.Join(dir, "jtbd.schema.json")); err != nil {
		return fmt.Errorf("generating JTBD schema: %w", err)
	}

	return nil
}

// Lean Startup Schema Generation

// GenerateLeanStartupSchema generates JSON Schema for the LeanStartupCanvas type.
func (g *Generator) GenerateLeanStartupSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.LeanStartupCanvas{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.LeanStartupCanvas")
	}

	schema.ID = jsonschema.ID(LeanStartupSchemaID)
	schema.Title = "Lean Startup Canvas"
	schema.Description = "Schema for Lean Startup canvas (Eric Ries methodology)"

	return schema, nil
}

// GenerateLeanStartupSchemaJSON generates JSON Schema for LeanStartup and returns it as JSON bytes.
func (g *Generator) GenerateLeanStartupSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateLeanStartupSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteLeanStartupSchema generates and writes the LeanStartup schema to a file.
func (g *Generator) WriteLeanStartupSchema(path string) error {
	return g.writeSchema(path, g.GenerateLeanStartupSchemaJSON)
}

// Design Thinking Schema Generation

// GenerateDesignThinkingSchema generates JSON Schema for the DesignThinkingCanvas type.
func (g *Generator) GenerateDesignThinkingSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.DesignThinkingCanvas{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.DesignThinkingCanvas")
	}

	schema.ID = jsonschema.ID(DesignThinkingSchemaID)
	schema.Title = "Design Thinking Canvas"
	schema.Description = "Schema for Design Thinking canvas (Stanford d.school methodology)"

	return schema, nil
}

// GenerateDesignThinkingSchemaJSON generates JSON Schema for DesignThinking and returns it as JSON bytes.
func (g *Generator) GenerateDesignThinkingSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateDesignThinkingSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteDesignThinkingSchema generates and writes the DesignThinking schema to a file.
func (g *Generator) WriteDesignThinkingSchema(path string) error {
	return g.writeSchema(path, g.GenerateDesignThinkingSchemaJSON)
}

// Jobs-to-be-Done Schema Generation

// GenerateJTBDSchema generates JSON Schema for the JTBDCanvas type.
func (g *Generator) GenerateJTBDSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&canvas.JTBDCanvas{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for canvas.JTBDCanvas")
	}

	schema.ID = jsonschema.ID(JTBDSchemaID)
	schema.Title = "Jobs-to-be-Done Canvas"
	schema.Description = "Schema for Jobs-to-be-Done canvas (Christensen, Ulwick methodology)"

	return schema, nil
}

// GenerateJTBDSchemaJSON generates JSON Schema for JTBD and returns it as JSON bytes.
func (g *Generator) GenerateJTBDSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateJTBDSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteJTBDSchema generates and writes the JTBD schema to a file.
func (g *Generator) WriteJTBDSchema(path string) error {
	return g.writeSchema(path, g.GenerateJTBDSchemaJSON)
}

// Journey Roadmap Schema Generation

// GenerateJourneyRoadmapSchema generates JSON Schema for the JourneyRoadmap type.
func (g *Generator) GenerateJourneyRoadmapSchema() (*jsonschema.Schema, error) {
	schema := g.Reflector.Reflect(&journey.JourneyRoadmap{})
	if schema == nil {
		return nil, fmt.Errorf("failed to generate schema for journey.JourneyRoadmap")
	}

	schema.ID = jsonschema.ID(JourneyRoadmapSchemaID)
	schema.Title = "Journey Roadmap"
	schema.Description = "Schema for capability evolution roadmaps that model how capabilities, outcomes, and business value evolve over time through transitions enabled by initiatives"

	return schema, nil
}

// GenerateJourneyRoadmapSchemaJSON generates JSON Schema for JourneyRoadmap and returns it as JSON bytes.
func (g *Generator) GenerateJourneyRoadmapSchemaJSON() ([]byte, error) {
	schema, err := g.GenerateJourneyRoadmapSchema()
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(schema, "", "  ")
}

// WriteJourneyRoadmapSchema generates and writes the JourneyRoadmap schema to a file.
func (g *Generator) WriteJourneyRoadmapSchema(path string) error {
	return g.writeSchema(path, g.GenerateJourneyRoadmapSchemaJSON)
}
