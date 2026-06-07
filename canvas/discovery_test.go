package canvas

import (
	"testing"
)

func TestNewDiscoverySnapshot(t *testing.T) {
	snapshot := NewDiscoverySnapshot("ds-001", "Week 3 Discovery")
	if snapshot.Metadata.ID != "ds-001" {
		t.Errorf("expected ID ds-001, got %s", snapshot.Metadata.ID)
	}
	if snapshot.Metadata.Title != "Week 3 Discovery" {
		t.Errorf("expected title 'Week 3 Discovery', got %s", snapshot.Metadata.Title)
	}
	if snapshot.Metadata.Version != VersionDiscovery1 {
		t.Errorf("expected version %s, got %s", VersionDiscovery1, snapshot.Metadata.Version)
	}
}

func TestDiscoverySnapshotInterviews(t *testing.T) {
	snapshot := NewDiscoverySnapshot("ds-002", "Week 4 Discovery")

	// No interviews yet
	if snapshot.HasWeeklyTouchpoint() {
		t.Error("expected HasWeeklyTouchpoint() to be false when empty")
	}
	if snapshot.InterviewCount() != 0 {
		t.Errorf("expected 0 interviews, got %d", snapshot.InterviewCount())
	}

	// Add interviews
	snapshot.Interviews = []CDInterview{
		{
			ID:              "int-1",
			Date:            "2024-01-15",
			ParticipantType: "enterprise user",
			Stories: []CDStory{
				{ID: "s1", Situation: "Onboarding", Behavior: "Struggled to find settings"},
				{ID: "s2", Situation: "Daily use", Behavior: "Created workaround"},
			},
		},
		{
			ID:              "int-2",
			Date:            "2024-01-17",
			ParticipantType: "new user",
			Stories: []CDStory{
				{ID: "s3", Situation: "First login", Behavior: "Clicked help immediately"},
			},
		},
	}

	if !snapshot.HasWeeklyTouchpoint() {
		t.Error("expected HasWeeklyTouchpoint() to be true")
	}
	if snapshot.InterviewCount() != 2 {
		t.Errorf("expected 2 interviews, got %d", snapshot.InterviewCount())
	}
	if snapshot.TotalStories() != 3 {
		t.Errorf("expected 3 total stories, got %d", snapshot.TotalStories())
	}
}

func TestDiscoverySnapshotAssumptionTests(t *testing.T) {
	snapshot := NewDiscoverySnapshot("ds-003", "Week 5 Discovery")
	snapshot.AssumptionTests = []CDAssumptionTest{
		{
			ID:     "at-1",
			Status: "completed",
			Result: "validated",
		},
		{
			ID:     "at-2",
			Status: "running",
		},
		{
			ID:     "at-3",
			Status: "completed",
			Result: "invalidated",
		},
	}

	completed := snapshot.CompletedTests()
	if len(completed) != 2 {
		t.Errorf("expected 2 completed tests, got %d", len(completed))
	}
}

func TestNewAssumptionMap(t *testing.T) {
	am := NewAssumptionMap("am-001", "Checkout Flow Assumptions")
	if am.Metadata.ID != "am-001" {
		t.Errorf("expected ID am-001, got %s", am.Metadata.ID)
	}
	if am.Metadata.Version != VersionDiscovery1 {
		t.Errorf("expected version %s, got %s", VersionDiscovery1, am.Metadata.Version)
	}
}

func TestAssumptionMapAllAssumptions(t *testing.T) {
	am := NewAssumptionMap("am-002", "Feature Assumptions")
	am.Desirability = []CDAssumption{
		{ID: "d1", Description: "Users want this feature"},
	}
	am.Viability = []CDAssumption{
		{ID: "v1", Description: "This will increase revenue"},
		{ID: "v2", Description: "Support costs won't increase"},
	}
	am.Feasibility = []CDAssumption{
		{ID: "f1", Description: "We can build this in 6 weeks"},
	}

	all := am.AllAssumptions()
	if len(all) != 4 {
		t.Errorf("expected 4 total assumptions, got %d", len(all))
	}
}

func TestAssumptionMapUnvalidated(t *testing.T) {
	am := NewAssumptionMap("am-003", "Test Assumptions")
	am.Desirability = []CDAssumption{
		{ID: "d1", Description: "Users want this", Validated: false},
		{ID: "d2", Description: "Users will pay", Validated: true},
	}
	am.Feasibility = []CDAssumption{
		{ID: "f1", Description: "Tech is ready", Validated: false},
	}

	unvalidated := am.UnvalidatedAssumptions()
	if len(unvalidated) != 2 {
		t.Errorf("expected 2 unvalidated assumptions, got %d", len(unvalidated))
	}
}

func TestAssumptionMapHighRisk(t *testing.T) {
	am := NewAssumptionMap("am-004", "Risk Assessment")
	am.Desirability = []CDAssumption{
		{ID: "d1", Description: "High risk", Importance: "high", Confidence: "low", Validated: false},
		{ID: "d2", Description: "Low risk", Importance: "low", Confidence: "high", Validated: false},
		{ID: "d3", Description: "Already validated", Importance: "high", Confidence: "low", Validated: true},
	}
	am.Viability = []CDAssumption{
		{ID: "v1", Description: "Another high risk", Importance: "high", Confidence: "low", Validated: false},
	}

	highRisk := am.HighRiskAssumptions()
	if len(highRisk) != 2 {
		t.Errorf("expected 2 high-risk assumptions, got %d", len(highRisk))
	}

	riskiest := am.RiskiestAssumption()
	if riskiest == nil {
		t.Error("expected RiskiestAssumption() to return non-nil")
	}
	if riskiest.ID != "d1" {
		t.Errorf("expected riskiest assumption ID 'd1', got %s", riskiest.ID)
	}
}

func TestNewExperienceMap(t *testing.T) {
	em := NewExperienceMap("em-001", "Onboarding Journey")
	if em.Metadata.ID != "em-001" {
		t.Errorf("expected ID em-001, got %s", em.Metadata.ID)
	}
	if em.Metadata.Version != VersionDiscovery1 {
		t.Errorf("expected version %s, got %s", VersionDiscovery1, em.Metadata.Version)
	}
}

func TestExperienceMapPhases(t *testing.T) {
	em := NewExperienceMap("em-002", "User Onboarding")
	em.Experience = "First-time user onboarding"
	em.Phases = []EMPhase{
		{
			ID:         "p1",
			Name:       "Awareness",
			Order:      1,
			Feeling:    "curious",
			PainPoints: []string{"Hard to find information"},
		},
		{
			ID:            "p2",
			Name:          "Sign Up",
			Order:         2,
			Feeling:       "frustrated",
			PainPoints:    []string{"Too many fields", "Password requirements confusing"},
			Opportunities: []string{"Simplify form", "Social login"},
		},
		{
			ID:      "p3",
			Name:    "First Use",
			Order:   3,
			Feeling: "overwhelmed",
		},
	}

	if len(em.Phases) != 3 {
		t.Errorf("expected 3 phases, got %d", len(em.Phases))
	}
	if len(em.Phases[1].PainPoints) != 2 {
		t.Errorf("expected 2 pain points in Sign Up phase, got %d", len(em.Phases[1].PainPoints))
	}
}
