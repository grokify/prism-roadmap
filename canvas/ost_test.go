package canvas

import (
	"encoding/json"
	"testing"
)

func TestOpportunitySolutionTree(t *testing.T) {
	ost := NewOpportunitySolutionTree("ost-1", "User Activation OST")

	ost.Outcome = OSTOutcome{
		ID:          "outcome-1",
		Description: "Increase user activation to 60%",
		Metric:      "% users completing onboarding",
		Target:      "60%",
		Opportunities: []OSTOpportunity{
			{
				ID:          "opp-1",
				Description: "Users don't understand value prop",
				Priority:    1,
				Solutions: []OSTSolution{
					{
						ID:          "sol-1",
						Description: "Interactive tutorial",
						Status:      "validated",
						Experiments: []OSTExperiment{
							{
								ID:         "exp-1",
								Hypothesis: "Tutorial increases completion by 20%",
								Method:     "A/B test",
								Status:     "completed",
								Result:     "success",
							},
						},
					},
					{
						ID:          "sol-2",
						Description: "Personalized onboarding",
						Status:      "testing",
						Experiments: []OSTExperiment{
							{
								ID:         "exp-2",
								Hypothesis: "Personalization increases engagement",
								Method:     "prototype",
								Status:     "running",
							},
						},
					},
				},
			},
			{
				ID:          "opp-2",
				Description: "Setup is too complex",
				Priority:    2,
			},
		},
	}

	// Test metadata
	if ost.Metadata.ID != "ost-1" {
		t.Errorf("ID = %v, want %v", ost.Metadata.ID, "ost-1")
	}
	if ost.Metadata.Version != VersionOST1 {
		t.Errorf("Version = %v, want %v", ost.Metadata.Version, VersionOST1)
	}

	// Test AllOpportunities
	opps := ost.AllOpportunities()
	if len(opps) != 2 {
		t.Errorf("AllOpportunities() length = %v, want 2", len(opps))
	}

	// Test AllSolutions
	sols := ost.AllSolutions()
	if len(sols) != 2 {
		t.Errorf("AllSolutions() length = %v, want 2", len(sols))
	}

	// Test AllExperiments
	exps := ost.AllExperiments()
	if len(exps) != 2 {
		t.Errorf("AllExperiments() length = %v, want 2", len(exps))
	}

	// Test PrioritizedOpportunities
	prioritized := ost.PrioritizedOpportunities()
	if len(prioritized) != 2 {
		t.Errorf("PrioritizedOpportunities() length = %v, want 2", len(prioritized))
	}
	if prioritized[0].Priority != 1 {
		t.Errorf("First priority = %v, want 1", prioritized[0].Priority)
	}
	if prioritized[1].Priority != 2 {
		t.Errorf("Second priority = %v, want 2", prioritized[1].Priority)
	}

	// Test ValidatedSolutions
	validated := ost.ValidatedSolutions()
	if len(validated) != 1 {
		t.Errorf("ValidatedSolutions() length = %v, want 1", len(validated))
	}
	if validated[0].ID != "sol-1" {
		t.Errorf("Validated solution ID = %v, want sol-1", validated[0].ID)
	}

	// Test ExperimentsByStatus
	running := ost.RunningExperiments()
	if len(running) != 1 {
		t.Errorf("RunningExperiments() length = %v, want 1", len(running))
	}

	completed := ost.CompletedExperiments()
	if len(completed) != 1 {
		t.Errorf("CompletedExperiments() length = %v, want 1", len(completed))
	}
}

func TestOSTJSON(t *testing.T) {
	original := &OpportunitySolutionTree{
		Metadata: Metadata{
			ID:      "ost-json-test",
			Title:   "JSON Test OST",
			Version: VersionOST1,
		},
		Outcome: OSTOutcome{
			ID:          "o1",
			Description: "Test outcome",
			Opportunities: []OSTOpportunity{
				{
					ID:          "opp1",
					Description: "Test opportunity",
					Solutions: []OSTSolution{
						{
							ID:          "sol1",
							Description: "Test solution",
							Experiments: []OSTExperiment{
								{
									ID:         "exp1",
									Hypothesis: "Test hypothesis",
									Method:     "A/B test",
									Status:     "planned",
								},
							},
						},
					},
				},
			},
		},
	}

	// Marshal
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal() error: %v", err)
	}

	// Unmarshal
	var parsed OpportunitySolutionTree
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("Unmarshal() error: %v", err)
	}

	// Validate structure
	if parsed.Metadata.ID != "ost-json-test" {
		t.Errorf("Metadata.ID = %v, want %v", parsed.Metadata.ID, "ost-json-test")
	}
	if len(parsed.Outcome.Opportunities) != 1 {
		t.Errorf("Opportunities length = %v, want 1", len(parsed.Outcome.Opportunities))
	}
	if len(parsed.Outcome.Opportunities[0].Solutions) != 1 {
		t.Errorf("Solutions length = %v, want 1", len(parsed.Outcome.Opportunities[0].Solutions))
	}
	if len(parsed.Outcome.Opportunities[0].Solutions[0].Experiments) != 1 {
		t.Errorf("Experiments length = %v, want 1", len(parsed.Outcome.Opportunities[0].Solutions[0].Experiments))
	}
}

func TestOSTPRDReference(t *testing.T) {
	ost := NewOpportunitySolutionTree("ost-prd", "PRD Linked OST")
	ost.PRDRef = &PRDReference{
		PRDID:      "prd-123",
		FeatureIDs: []string{"feat-1"},
	}
	ost.Outcome = OSTOutcome{
		ID:     "o1",
		OKRRef: "okr-quarterly-goals",
		Opportunities: []OSTOpportunity{
			{
				ID: "opp1",
				Solutions: []OSTSolution{
					{
						ID:             "sol1",
						RequirementRef: "req-123",
					},
				},
			},
		},
	}

	ref := ost.GetPRDReference()
	if ref == nil {
		t.Fatal("GetPRDReference() returned nil")
	}
	if ref.PRDID != "prd-123" {
		t.Errorf("PRDID = %v, want %v", ref.PRDID, "prd-123")
	}

	if ost.Outcome.OKRRef != "okr-quarterly-goals" {
		t.Errorf("OKRRef = %v, want %v", ost.Outcome.OKRRef, "okr-quarterly-goals")
	}

	sols := ost.AllSolutions()
	if sols[0].RequirementRef != "req-123" {
		t.Errorf("RequirementRef = %v, want %v", sols[0].RequirementRef, "req-123")
	}
}
