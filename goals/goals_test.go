package goals

import (
	"testing"

	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/goals/v2mom"
)

func TestNewOKR(t *testing.T) {
	okrSet := &okr.OKRSet{
		OKRs: []okr.OKR{
			{
				Objective: okr.Objective{
					ID:          "O1",
					Title:       "Increase market share",
					Description: "Become the market leader in our segment",
					Owner:       "Product Team",
					Status:      "Active",
				},
				KeyResults: []okr.KeyResult{
					{
						ID:       "KR1",
						Title:    "Market share growth",
						Metric:   "Market share percentage",
						Baseline: "10%",
						Target:   "20%",
						Current:  "15%",
						Score:    0.5,
						Status:   "On Track",
					},
				},
			},
		},
	}

	g := NewOKR(okrSet)

	if !g.IsOKR() {
		t.Error("Expected IsOKR() to return true")
	}
	if g.IsV2MOM() {
		t.Error("Expected IsV2MOM() to return false")
	}
	if g.Framework != FrameworkOKR {
		t.Errorf("Expected framework %s, got %s", FrameworkOKR, g.Framework)
	}
	if g.OKR == nil {
		t.Error("Expected OKR to be set")
	}
}

func TestNewV2MOM(t *testing.T) {
	v := &v2mom.V2MOM{
		Vision: "Be the best",
		Methods: []v2mom.Method{
			{
				ID:          "M1",
				Name:        "Expand to new markets",
				Description: "Enter APAC region",
				Owner:       "Sales Team",
				Status:      "In Progress",
				Priority:    "P1",
			},
		},
		Measures: []v2mom.Measure{
			{
				ID:       "ME1",
				Name:     "Revenue growth",
				Baseline: "$1M",
				Target:   "$2M",
				Current:  "$1.5M",
				Progress: 0.5,
				Status:   "On Track",
			},
		},
	}

	g := NewV2MOM(v)

	if g.IsOKR() {
		t.Error("Expected IsOKR() to return false")
	}
	if !g.IsV2MOM() {
		t.Error("Expected IsV2MOM() to return true")
	}
	if g.Framework != FrameworkV2MOM {
		t.Errorf("Expected framework %s, got %s", FrameworkV2MOM, g.Framework)
	}
	if g.V2MOM == nil {
		t.Error("Expected V2MOM to be set")
	}
}

func TestGoalItemsOKR(t *testing.T) {
	okrSet := &okr.OKRSet{
		OKRs: []okr.OKR{
			{
				Objective: okr.Objective{
					ID:          "O1",
					Title:       "Increase market share",
					Description: "Become the market leader",
					Owner:       "Product Team",
					Status:      "Active",
					Tags:        []string{"growth", "market"},
				},
			},
			{
				Objective: okr.Objective{
					ID:          "O2",
					Title:       "Improve customer satisfaction",
					Description: "Achieve highest NPS in industry",
					Owner:       "Support Team",
					Status:      "Active",
				},
			},
		},
	}

	g := NewOKR(okrSet)
	items := g.GoalItems()

	if len(items) != 2 {
		t.Errorf("Expected 2 goal items, got %d", len(items))
	}

	if items[0].ID != "O1" {
		t.Errorf("Expected ID O1, got %s", items[0].ID)
	}
	if items[0].Title != "Increase market share" {
		t.Errorf("Expected title 'Increase market share', got %s", items[0].Title)
	}
	if items[0].Owner != "Product Team" {
		t.Errorf("Expected owner 'Product Team', got %s", items[0].Owner)
	}
}

func TestGoalItemsV2MOM(t *testing.T) {
	v := &v2mom.V2MOM{
		Vision: "Be the best",
		Methods: []v2mom.Method{
			{
				ID:          "M1",
				Name:        "Expand to new markets",
				Description: "Enter APAC region",
				Owner:       "Sales Team",
				Status:      "In Progress",
				Priority:    "P1",
			},
			{
				ID:          "M2",
				Name:        "Launch new product line",
				Description: "Release v2 platform",
				Owner:       "Engineering",
				Priority:    "P0",
			},
		},
	}

	g := NewV2MOM(v)
	items := g.GoalItems()

	if len(items) != 2 {
		t.Errorf("Expected 2 goal items, got %d", len(items))
	}

	if items[0].ID != "M1" {
		t.Errorf("Expected ID M1, got %s", items[0].ID)
	}
	if items[0].Title != "Expand to new markets" {
		t.Errorf("Expected title 'Expand to new markets', got %s", items[0].Title)
	}
	if items[0].Priority != "P1" {
		t.Errorf("Expected priority P1, got %s", items[0].Priority)
	}
}

func TestResultItemsOKR(t *testing.T) {
	okrSet := &okr.OKRSet{
		OKRs: []okr.OKR{
			{
				Objective: okr.Objective{
					ID:    "O1",
					Title: "Increase market share",
				},
				KeyResults: []okr.KeyResult{
					{
						ID:       "KR1",
						Title:    "Market share growth",
						Metric:   "Market share percentage",
						Baseline: "10%",
						Target:   "20%",
						Current:  "15%",
						Score:    0.5,
						Status:   "On Track",
					},
					{
						ID:       "KR2",
						Title:    "Revenue increase",
						Metric:   "Annual revenue",
						Baseline: "$1M",
						Target:   "$2M",
						Score:    0.7,
					},
				},
			},
		},
	}

	g := NewOKR(okrSet)
	items := g.ResultItems()

	if len(items) != 2 {
		t.Errorf("Expected 2 result items, got %d", len(items))
	}

	if items[0].ID != "KR1" {
		t.Errorf("Expected ID KR1, got %s", items[0].ID)
	}
	if items[0].GoalID != "O1" {
		t.Errorf("Expected GoalID O1, got %s", items[0].GoalID)
	}
	if items[0].Score != 0.5 {
		t.Errorf("Expected score 0.5, got %f", items[0].Score)
	}
}

func TestResultItemsWithPhaseTargets(t *testing.T) {
	okrSet := &okr.OKRSet{
		OKRs: []okr.OKR{
			{
				Objective: okr.Objective{
					ID:    "O1",
					Title: "Increase market share",
				},
				KeyResults: []okr.KeyResult{
					{
						ID:     "KR1",
						Title:  "Market share growth",
						Target: "20%",
						PhaseTargets: []okr.PhaseTarget{
							{PhaseID: "P1", Target: "15%", Status: "achieved"},
							{PhaseID: "P2", Target: "20%", Status: "in_progress"},
						},
					},
				},
			},
		},
	}

	g := NewOKR(okrSet)
	items := g.ResultItems()

	// Should have 2 items (one per phase target)
	if len(items) != 2 {
		t.Errorf("Expected 2 result items (one per phase), got %d", len(items))
	}

	if items[0].PhaseID != "P1" {
		t.Errorf("Expected PhaseID P1, got %s", items[0].PhaseID)
	}
	if items[0].PhaseTarget != "15%" {
		t.Errorf("Expected PhaseTarget 15%%, got %s", items[0].PhaseTarget)
	}
	if items[0].Status != "achieved" {
		t.Errorf("Expected status 'achieved', got %s", items[0].Status)
	}
}

func TestResultItemsV2MOM(t *testing.T) {
	v := &v2mom.V2MOM{
		Vision: "Be the best",
		Methods: []v2mom.Method{
			{
				ID:   "M1",
				Name: "Expand markets",
				Measures: []v2mom.Measure{
					{
						ID:       "ME1",
						Name:     "New market revenue",
						Baseline: "$0",
						Target:   "$1M",
						Current:  "$500K",
						Progress: 0.5,
					},
				},
			},
		},
		Measures: []v2mom.Measure{
			{
				ID:       "ME2",
				Name:     "Total revenue",
				Baseline: "$5M",
				Target:   "$10M",
				Current:  "$7M",
				Progress: 0.4,
			},
		},
	}

	g := NewV2MOM(v)
	items := g.ResultItems()

	// Should have 2 items (1 global + 1 nested)
	if len(items) != 2 {
		t.Errorf("Expected 2 result items, got %d", len(items))
	}

	// Global measure should be first
	if items[0].ID != "ME2" {
		t.Errorf("Expected ID ME2, got %s", items[0].ID)
	}
	if items[0].GoalID != "" {
		t.Errorf("Expected empty GoalID for global measure, got %s", items[0].GoalID)
	}

	// Nested measure should have GoalID
	if items[1].GoalID != "M1" {
		t.Errorf("Expected GoalID M1 for nested measure, got %s", items[1].GoalID)
	}
}

func TestLabels(t *testing.T) {
	okrGoals := NewOKR(&okr.OKRSet{})
	v2momGoals := NewV2MOM(&v2mom.V2MOM{})

	if okrGoals.GoalLabel() != "Objectives" {
		t.Errorf("Expected OKR GoalLabel 'Objectives', got %s", okrGoals.GoalLabel())
	}
	if okrGoals.ResultLabel() != "Key Results" {
		t.Errorf("Expected OKR ResultLabel 'Key Results', got %s", okrGoals.ResultLabel())
	}

	if v2momGoals.GoalLabel() != "Methods" {
		t.Errorf("Expected V2MOM GoalLabel 'Methods', got %s", v2momGoals.GoalLabel())
	}
	if v2momGoals.ResultLabel() != "Measures" {
		t.Errorf("Expected V2MOM ResultLabel 'Measures', got %s", v2momGoals.ResultLabel())
	}
}

func TestNilGoals(t *testing.T) {
	var g *Goals

	if g.IsOKR() {
		t.Error("Expected nil.IsOKR() to return false")
	}
	if g.IsV2MOM() {
		t.Error("Expected nil.IsV2MOM() to return false")
	}
	if g.GoalLabel() != "Goals" {
		t.Errorf("Expected nil GoalLabel 'Goals', got %s", g.GoalLabel())
	}
	if g.ResultLabel() != "Results" {
		t.Errorf("Expected nil ResultLabel 'Results', got %s", g.ResultLabel())
	}
	if len(g.GoalItems()) != 0 {
		t.Error("Expected nil.GoalItems() to return empty slice")
	}
	if len(g.ResultItems()) != 0 {
		t.Error("Expected nil.ResultItems() to return empty slice")
	}
}

func TestResultItemsByPhase(t *testing.T) {
	okrSet := &okr.OKRSet{
		OKRs: []okr.OKR{
			{
				Objective: okr.Objective{ID: "O1", Title: "Goal 1"},
				KeyResults: []okr.KeyResult{
					{
						ID:    "KR1",
						Title: "Result 1",
						PhaseTargets: []okr.PhaseTarget{
							{PhaseID: "P1", Target: "10"},
							{PhaseID: "P2", Target: "20"},
						},
					},
					{
						ID:    "KR2",
						Title: "Result 2",
						PhaseTargets: []okr.PhaseTarget{
							{PhaseID: "P1", Target: "5"},
						},
					},
				},
			},
		},
	}

	g := NewOKR(okrSet)
	byPhase := g.ResultItemsByPhase()

	if len(byPhase["P1"]) != 2 {
		t.Errorf("Expected 2 items for P1, got %d", len(byPhase["P1"]))
	}
	if len(byPhase["P2"]) != 1 {
		t.Errorf("Expected 1 item for P2, got %d", len(byPhase["P2"]))
	}
}
