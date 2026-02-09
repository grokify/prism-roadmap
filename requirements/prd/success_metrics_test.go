package prd

import (
	"encoding/json"
	"testing"
)

func TestSuccessMetricsJSONRoundTrip(t *testing.T) {
	original := &SuccessMetrics{
		NorthStar: []Metric{
			{
				ID:                "NS-001",
				Name:              "User Retention",
				Description:       "Monthly active users returning",
				Baseline:          "65%",
				Target:            "80%",
				MeasurementMethod: "Analytics dashboard",
			},
		},
		Supporting: []Metric{
			{
				ID:       "SUP-001",
				Name:     "NPS Score",
				Baseline: "40",
				Target:   "60",
			},
		},
		Guardrail: []Metric{
			{
				ID:                "GR-001",
				Name:              "Page Load Time",
				Description:       "95th percentile load time",
				Baseline:          "2s",
				Target:            "<3s",
				MeasurementMethod: "APM monitoring",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal SuccessMetrics: %v", err)
	}

	// Unmarshal back
	var decoded SuccessMetrics
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal SuccessMetrics: %v", err)
	}

	// Verify north star metrics
	if len(decoded.NorthStar) != 1 {
		t.Errorf("Expected 1 north star metric, got %d", len(decoded.NorthStar))
	}
	if decoded.NorthStar[0].ID != "NS-001" {
		t.Errorf("Expected north star ID 'NS-001', got '%s'", decoded.NorthStar[0].ID)
	}
	if decoded.NorthStar[0].Name != "User Retention" {
		t.Errorf("Expected north star name 'User Retention', got '%s'", decoded.NorthStar[0].Name)
	}
	if decoded.NorthStar[0].Target != "80%" {
		t.Errorf("Expected target '80%%', got '%s'", decoded.NorthStar[0].Target)
	}

	// Verify supporting metrics
	if len(decoded.Supporting) != 1 {
		t.Errorf("Expected 1 supporting metric, got %d", len(decoded.Supporting))
	}
	if decoded.Supporting[0].ID != "SUP-001" {
		t.Errorf("Expected supporting ID 'SUP-001', got '%s'", decoded.Supporting[0].ID)
	}

	// Verify guardrail metrics
	if len(decoded.Guardrail) != 1 {
		t.Errorf("Expected 1 guardrail metric, got %d", len(decoded.Guardrail))
	}
	if decoded.Guardrail[0].ID != "GR-001" {
		t.Errorf("Expected guardrail ID 'GR-001', got '%s'", decoded.Guardrail[0].ID)
	}
}

func TestMetricJSONRoundTrip(t *testing.T) {
	original := Metric{
		ID:                "M-001",
		Name:              "Test Metric",
		Description:       "A test metric for validation",
		Baseline:          "10",
		Target:            "100",
		MeasurementMethod: "Manual count",
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal Metric: %v", err)
	}

	// Unmarshal back
	var decoded Metric
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Metric: %v", err)
	}

	// Verify all fields
	if decoded.ID != original.ID {
		t.Errorf("ID mismatch: got '%s', want '%s'", decoded.ID, original.ID)
	}
	if decoded.Name != original.Name {
		t.Errorf("Name mismatch: got '%s', want '%s'", decoded.Name, original.Name)
	}
	if decoded.Description != original.Description {
		t.Errorf("Description mismatch: got '%s', want '%s'", decoded.Description, original.Description)
	}
	if decoded.Baseline != original.Baseline {
		t.Errorf("Baseline mismatch: got '%s', want '%s'", decoded.Baseline, original.Baseline)
	}
	if decoded.Target != original.Target {
		t.Errorf("Target mismatch: got '%s', want '%s'", decoded.Target, original.Target)
	}
	if decoded.MeasurementMethod != original.MeasurementMethod {
		t.Errorf("MeasurementMethod mismatch: got '%s', want '%s'", decoded.MeasurementMethod, original.MeasurementMethod)
	}
}

func TestSuccessMetricsOmitEmpty(t *testing.T) {
	// Test that optional fields are omitted when empty
	minimal := &SuccessMetrics{
		NorthStar: []Metric{
			{
				ID:     "NS-001",
				Name:   "Core Metric",
				Target: "100%",
			},
		},
	}

	data, err := json.Marshal(minimal)
	if err != nil {
		t.Fatalf("Failed to marshal minimal SuccessMetrics: %v", err)
	}

	// Check that omitted fields are not present
	jsonStr := string(data)
	if jsonStr == "" {
		t.Error("JSON output is empty")
	}

	// Unmarshal to map to check structure
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Failed to unmarshal to map: %v", err)
	}

	// Supporting and Guardrail should be omitted
	if _, exists := result["supporting"]; exists {
		t.Error("Expected 'supporting' to be omitted when empty")
	}
	if _, exists := result["guardrail"]; exists {
		t.Error("Expected 'guardrail' to be omitted when empty")
	}

	// NorthStar should be present
	if _, exists := result["northStar"]; !exists {
		t.Error("Expected 'northStar' to be present")
	}
}

func TestDocumentWithSuccessMetrics(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-001",
			Title:   "Test PRD",
			Version: "1.0",
			Status:  StatusDraft,
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem",
			ProposedSolution: "Test solution",
		},
		SuccessMetrics: &SuccessMetrics{
			NorthStar: []Metric{
				{
					ID:     "NS-001",
					Name:   "Revenue",
					Target: "$1M ARR",
				},
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal Document with SuccessMetrics: %v", err)
	}

	// Unmarshal back
	var decoded Document
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Document: %v", err)
	}

	// Verify success metrics are preserved
	if decoded.SuccessMetrics == nil {
		t.Fatal("Expected SuccessMetrics to be non-nil")
	}
	if len(decoded.SuccessMetrics.NorthStar) != 1 {
		t.Errorf("Expected 1 north star metric, got %d", len(decoded.SuccessMetrics.NorthStar))
	}
	if decoded.SuccessMetrics.NorthStar[0].Target != "$1M ARR" {
		t.Errorf("Expected target '$1M ARR', got '%s'", decoded.SuccessMetrics.NorthStar[0].Target)
	}
}

func TestDocumentWithNonGoals(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-001",
			Title:   "Test PRD",
			Version: "1.0",
			Status:  StatusDraft,
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem",
			ProposedSolution: "Test solution",
		},
		NonGoals: []NonGoal{
			{
				ID:          "NG-001",
				Title:       "Mobile App",
				Description: "Native mobile application",
				Rationale:   "Focus on web first",
				FuturePhase: "Phase 2",
			},
			{
				ID:    "NG-002",
				Title: "Multi-language support",
			},
		},
	}

	// Marshal to JSON
	data, err := json.Marshal(doc)
	if err != nil {
		t.Fatalf("Failed to marshal Document with NonGoals: %v", err)
	}

	// Unmarshal back
	var decoded Document
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal Document: %v", err)
	}

	// Verify non-goals are preserved
	if len(decoded.NonGoals) != 2 {
		t.Fatalf("Expected 2 non-goals, got %d", len(decoded.NonGoals))
	}
	if decoded.NonGoals[0].ID != "NG-001" {
		t.Errorf("Expected first non-goal ID 'NG-001', got '%s'", decoded.NonGoals[0].ID)
	}
	if decoded.NonGoals[0].Rationale != "Focus on web first" {
		t.Errorf("Expected rationale 'Focus on web first', got '%s'", decoded.NonGoals[0].Rationale)
	}
	if decoded.NonGoals[1].FuturePhase != "" {
		t.Errorf("Expected empty FuturePhase, got '%s'", decoded.NonGoals[1].FuturePhase)
	}
}
