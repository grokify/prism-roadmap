package schema

import (
	"encoding/json"
	"testing"
)

func TestPRDSchemaEmbedded(t *testing.T) {
	if len(PRDSchemaJSON) == 0 {
		t.Fatal("PRDSchemaJSON is empty")
	}

	var schema map[string]any
	if err := json.Unmarshal(PRDSchemaJSON, &schema); err != nil {
		t.Fatalf("PRDSchemaJSON is not valid JSON: %v", err)
	}

	// Schema uses $ref to point to Document definition
	expectedKeys := []string{"$schema", "$id", "$ref", "$defs"}
	for _, key := range expectedKeys {
		if _, ok := schema[key]; !ok {
			t.Errorf("PRDSchemaJSON missing expected key: %s", key)
		}
	}
}

func TestPRDSchemaExtendedSections(t *testing.T) {
	var schema map[string]any
	if err := json.Unmarshal(PRDSchemaJSON, &schema); err != nil {
		t.Fatalf("PRDSchemaJSON is not valid JSON: %v", err)
	}

	// Get $defs.Document.properties since schema uses $ref
	defs, ok := schema["$defs"].(map[string]any)
	if !ok {
		t.Fatal("$defs is not an object")
	}

	document, ok := defs["Document"].(map[string]any)
	if !ok {
		t.Fatal("$defs.Document is not an object")
	}

	props, ok := document["properties"].(map[string]any)
	if !ok {
		t.Fatal("$defs.Document.properties is not an object")
	}

	// Check that extended sections are reflected from Go types (camelCase)
	extendedSections := []string{"problem", "market", "solution", "decisions", "reviews", "revisionHistory", "goals"}
	for _, section := range extendedSections {
		if _, ok := props[section]; !ok {
			t.Errorf("PRDSchemaJSON properties missing extended section: %s", section)
		}
	}
}

func TestPRDSchemaFunction(t *testing.T) {
	schema := PRDSchema()
	if schema == "" {
		t.Fatal("PRDSchema() returned empty string")
	}
	if schema != string(PRDSchemaJSON) {
		t.Error("PRDSchema() does not match PRDSchemaJSON")
	}
}

func TestPRDSchemaBytesFunction(t *testing.T) {
	schemaBytes := PRDSchemaBytes()
	if len(schemaBytes) == 0 {
		t.Fatal("PRDSchemaBytes() returned empty slice")
	}
	if string(schemaBytes) != string(PRDSchemaJSON) {
		t.Error("PRDSchemaBytes() does not match PRDSchemaJSON")
	}
}

func TestAssessmentSchemasEmbedded(t *testing.T) {
	tests := []struct {
		name string
		data []byte
	}{
		{"OpportunityAssessmentSchemaJSON", OpportunityAssessmentSchemaJSON},
		{"EvidenceSchemaJSON", EvidenceSchemaJSON},
		{"DimensionDefinitionSchemaJSON", DimensionDefinitionSchemaJSON},
		{"OpportunityRankSchemaJSON", OpportunityRankSchemaJSON},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.data) == 0 {
				t.Fatalf("%s is empty", tt.name)
			}
			var schema map[string]any
			if err := json.Unmarshal(tt.data, &schema); err != nil {
				t.Fatalf("%s is not valid JSON: %v", tt.name, err)
			}
			for _, key := range []string{"$schema", "$id", "$ref", "$defs"} {
				if _, ok := schema[key]; !ok {
					t.Errorf("%s missing expected key: %s", tt.name, key)
				}
			}
		})
	}
}

func TestOpportunityAssessmentSchemaFunctions(t *testing.T) {
	if OpportunityAssessmentSchema() != string(OpportunityAssessmentSchemaJSON) {
		t.Error("OpportunityAssessmentSchema() does not match OpportunityAssessmentSchemaJSON")
	}
	if string(OpportunityAssessmentSchemaBytes()) != string(OpportunityAssessmentSchemaJSON) {
		t.Error("OpportunityAssessmentSchemaBytes() does not match OpportunityAssessmentSchemaJSON")
	}
}
