package schema

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/invopop/jsonschema"

	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/goals/v2mom"
	"github.com/grokify/prism-execution/requirements/prd"
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
