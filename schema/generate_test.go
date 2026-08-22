package schema

import (
	"encoding/json"
	"os"
	"testing"
)

func TestGeneratePRDSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GeneratePRDSchema()
	if err != nil {
		t.Fatalf("GeneratePRDSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Structured PRD" {
		t.Errorf("expected title 'Structured PRD', got %q", schema.Title)
	}

	if string(schema.ID) != PRDSchemaID {
		t.Errorf("expected ID %q, got %q", PRDSchemaID, schema.ID)
	}
}

func TestGeneratePRDSchemaJSON(t *testing.T) {
	gen := NewGenerator()

	data, err := gen.GeneratePRDSchemaJSON()
	if err != nil {
		t.Fatalf("GeneratePRDSchemaJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("generated JSON is empty")
	}

	// Verify it's valid JSON
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("generated JSON is invalid: %v", err)
	}

	// Check for expected top-level keys (invopop/jsonschema uses $ref pattern)
	expectedKeys := []string{"$schema", "$id", "$ref", "$defs", "title"}
	for _, key := range expectedKeys {
		if _, ok := schema[key]; !ok {
			t.Errorf("generated schema missing key: %s", key)
		}
	}

	// Verify $ref points to Document definition
	ref, ok := schema["$ref"].(string)
	if !ok || ref != "#/$defs/Document" {
		t.Errorf("expected $ref to be '#/$defs/Document', got %v", schema["$ref"])
	}
}

func TestGeneratorReflectsExtendedSections(t *testing.T) {
	gen := NewGenerator()

	data, err := gen.GeneratePRDSchemaJSON()
	if err != nil {
		t.Fatalf("GeneratePRDSchemaJSON failed: %v", err)
	}

	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("generated JSON is invalid: %v", err)
	}

	// Navigate to $defs/Document/properties (invopop/jsonschema structure)
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("$defs is not an object")
	}

	doc, ok := defs["Document"].(map[string]any)
	if !ok {
		t.Fatal("$defs/Document is not an object")
	}

	props, ok := doc["properties"].(map[string]any)
	if !ok {
		t.Fatal("$defs/Document/properties is not an object")
	}

	// Check that extended sections are reflected from Go types
	extendedSections := []string{"problem", "market", "solution", "decisions", "reviews", "revisionHistory", "goals"}
	for _, section := range extendedSections {
		if _, ok := props[section]; !ok {
			t.Errorf("generated schema missing extended section: %s", section)
		}
	}
}

func TestGenerateShapeUpPitchSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateShapeUpPitchSchema()
	if err != nil {
		t.Fatalf("GenerateShapeUpPitchSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Shape Up Pitch" {
		t.Errorf("expected title 'Shape Up Pitch', got %q", schema.Title)
	}

	if string(schema.ID) != ShapeUpPitchSchemaID {
		t.Errorf("expected ID %q, got %q", ShapeUpPitchSchemaID, schema.ID)
	}

	// Verify JSON generation
	data, err := gen.GenerateShapeUpPitchSchemaJSON()
	if err != nil {
		t.Fatalf("GenerateShapeUpPitchSchemaJSON failed: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("generated JSON is empty")
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("generated JSON is invalid: %v", err)
	}
}

func TestGenerateShapeUpBetSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateShapeUpBetSchema()
	if err != nil {
		t.Fatalf("GenerateShapeUpBetSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Shape Up Bet" {
		t.Errorf("expected title 'Shape Up Bet', got %q", schema.Title)
	}
}

func TestGenerateShapeUpScopeSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateShapeUpScopeSchema()
	if err != nil {
		t.Fatalf("GenerateShapeUpScopeSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Shape Up Scope" {
		t.Errorf("expected title 'Shape Up Scope', got %q", schema.Title)
	}
}

func TestGenerateDiscoverySnapshotSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateDiscoverySnapshotSchema()
	if err != nil {
		t.Fatalf("GenerateDiscoverySnapshotSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Discovery Snapshot" {
		t.Errorf("expected title 'Discovery Snapshot', got %q", schema.Title)
	}
}

func TestGenerateAssumptionMapSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateAssumptionMapSchema()
	if err != nil {
		t.Fatalf("GenerateAssumptionMapSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Assumption Map" {
		t.Errorf("expected title 'Assumption Map', got %q", schema.Title)
	}
}

func TestGenerateExperienceMapSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateExperienceMapSchema()
	if err != nil {
		t.Fatalf("GenerateExperienceMapSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Experience Map" {
		t.Errorf("expected title 'Experience Map', got %q", schema.Title)
	}
}

func TestGenerateLeanStartupSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateLeanStartupSchema()
	if err != nil {
		t.Fatalf("GenerateLeanStartupSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Lean Startup Canvas" {
		t.Errorf("expected title 'Lean Startup Canvas', got %q", schema.Title)
	}

	if string(schema.ID) != LeanStartupSchemaID {
		t.Errorf("expected ID %q, got %q", LeanStartupSchemaID, schema.ID)
	}
}

func TestGenerateDesignThinkingSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateDesignThinkingSchema()
	if err != nil {
		t.Fatalf("GenerateDesignThinkingSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Design Thinking Canvas" {
		t.Errorf("expected title 'Design Thinking Canvas', got %q", schema.Title)
	}

	if string(schema.ID) != DesignThinkingSchemaID {
		t.Errorf("expected ID %q, got %q", DesignThinkingSchemaID, schema.ID)
	}
}

func TestGenerateJTBDSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateJTBDSchema()
	if err != nil {
		t.Fatalf("GenerateJTBDSchema failed: %v", err)
	}

	if schema == nil {
		t.Fatal("schema is nil")
	}

	if schema.Title != "Jobs-to-be-Done Canvas" {
		t.Errorf("expected title 'Jobs-to-be-Done Canvas', got %q", schema.Title)
	}

	if string(schema.ID) != JTBDSchemaID {
		t.Errorf("expected ID %q, got %q", JTBDSchemaID, schema.ID)
	}
}

func TestGenerateCanvasSchemas(t *testing.T) {
	gen := NewGenerator()

	// Create temp directory for testing
	tempDir := t.TempDir()

	err := gen.GenerateCanvasSchemas(tempDir)
	if err != nil {
		t.Fatalf("GenerateCanvasSchemas failed: %v", err)
	}

	// Check that all schema files were created
	expectedFiles := []string{
		"shapeup-pitch.schema.json",
		"shapeup-bet.schema.json",
		"shapeup-scope.schema.json",
		"discovery-snapshot.schema.json",
		"assumption-map.schema.json",
		"experience-map.schema.json",
		"leanstartup.schema.json",
		"designthinking.schema.json",
		"jtbd.schema.json",
	}

	for _, fileName := range expectedFiles {
		filePath := tempDir + "/" + fileName
		data, err := readFile(filePath)
		if err != nil {
			t.Errorf("file %s does not exist: %v", fileName, err)
			continue
		}

		// Verify it's valid JSON
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Errorf("file %s is not valid JSON: %v", fileName, err)
		}
	}
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func TestGenerateOpportunityAssessmentSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateOpportunityAssessmentSchema()
	if err != nil {
		t.Fatalf("GenerateOpportunityAssessmentSchema failed: %v", err)
	}
	if schema == nil {
		t.Fatal("schema is nil")
	}
	if schema.Title != "Opportunity Assessment" {
		t.Errorf("expected title 'Opportunity Assessment', got %q", schema.Title)
	}
	if string(schema.ID) != OpportunityAssessmentSchemaID {
		t.Errorf("expected ID %q, got %q", OpportunityAssessmentSchemaID, schema.ID)
	}
}

func TestGenerateOpportunityAssessmentSchemaJSON(t *testing.T) {
	gen := NewGenerator()
	data, err := gen.GenerateOpportunityAssessmentSchemaJSON()
	if err != nil {
		t.Fatalf("GenerateOpportunityAssessmentSchemaJSON failed: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("generated JSON is empty")
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatalf("generated JSON is not valid: %v", err)
	}
}

func TestGenerateEvidenceSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateEvidenceSchema()
	if err != nil {
		t.Fatalf("GenerateEvidenceSchema failed: %v", err)
	}
	if schema == nil {
		t.Fatal("schema is nil")
	}
	if string(schema.ID) != EvidenceSchemaID {
		t.Errorf("expected ID %q, got %q", EvidenceSchemaID, schema.ID)
	}
}

func TestGenerateDimensionDefinitionSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateDimensionDefinitionSchema()
	if err != nil {
		t.Fatalf("GenerateDimensionDefinitionSchema failed: %v", err)
	}
	if schema == nil {
		t.Fatal("schema is nil")
	}
	if string(schema.ID) != DimensionDefinitionSchemaID {
		t.Errorf("expected ID %q, got %q", DimensionDefinitionSchemaID, schema.ID)
	}
}

func TestGenerateOpportunityRankSchema(t *testing.T) {
	gen := NewGenerator()

	schema, err := gen.GenerateOpportunityRankSchema()
	if err != nil {
		t.Fatalf("GenerateOpportunityRankSchema failed: %v", err)
	}
	if schema == nil {
		t.Fatal("schema is nil")
	}
	if string(schema.ID) != OpportunityRankSchemaID {
		t.Errorf("expected ID %q, got %q", OpportunityRankSchemaID, schema.ID)
	}
}

func TestGenerateAssessmentSchemas(t *testing.T) {
	gen := NewGenerator()
	dir := t.TempDir()

	if err := gen.GenerateAssessmentSchemas(dir); err != nil {
		t.Fatalf("GenerateAssessmentSchemas failed: %v", err)
	}

	for _, fileName := range []string{
		"opportunity-assessment.schema.json",
		"evidence.schema.json",
		"dimension-definition.schema.json",
		"opportunity-rank.schema.json",
	} {
		path := dir + "/" + fileName
		data, err := readFile(path)
		if err != nil {
			t.Errorf("file %s does not exist: %v", fileName, err)
			continue
		}
		var schema map[string]any
		if err := json.Unmarshal(data, &schema); err != nil {
			t.Errorf("file %s is not valid JSON: %v", fileName, err)
		}
	}
}
