package canvas

import (
	"testing"
)

func TestNewJTBDCanvas(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-001", "Customer Onboarding Jobs")
	if canvas.Metadata.ID != "jtbd-001" {
		t.Errorf("expected ID jtbd-001, got %s", canvas.Metadata.ID)
	}
	if canvas.Metadata.Title != "Customer Onboarding Jobs" {
		t.Errorf("expected title 'Customer Onboarding Jobs', got %s", canvas.Metadata.Title)
	}
	if canvas.Metadata.Version != VersionJTBD1 {
		t.Errorf("expected version %s, got %s", VersionJTBD1, canvas.Metadata.Version)
	}
}

func TestCalculateOpportunityScore(t *testing.T) {
	tests := []struct {
		name         string
		importance   float64
		satisfaction float64
		expected     float64
	}{
		{"underserved", 9, 3, 15}, // 9 + max(9-3, 0) = 9 + 6 = 15
		{"appropriately served", 7, 7, 7},  // 7 + max(7-7, 0) = 7 + 0 = 7
		{"overserved", 5, 8, 5},    // 5 + max(5-8, 0) = 5 + 0 = 5
		{"high importance low satisfaction", 10, 1, 19}, // 10 + 9 = 19
		{"low importance high satisfaction", 2, 9, 2},   // 2 + 0 = 2
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateOpportunityScore(tt.importance, tt.satisfaction)
			if score != tt.expected {
				t.Errorf("expected score %v, got %v", tt.expected, score)
			}
		})
	}
}

func TestJTBDCanvasJobFiltering(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-002", "Job Filtering Test")
	canvas.MainJob = &JobStatement{
		ID:        "main-1",
		Statement: "Get work done efficiently",
		Type:      JobTypeFunctional,
	}
	canvas.RelatedJobs = []JobStatement{
		{ID: "r1", Statement: "Feel confident in my work", Type: JobTypeEmotional},
		{ID: "r2", Statement: "Be seen as productive", Type: JobTypeSocial},
		{ID: "r3", Statement: "Complete tasks quickly", Type: JobTypeFunctional},
		{ID: "r4", Statement: "Feel in control", Type: JobTypeEmotional},
	}

	// Test FunctionalJobs
	functional := canvas.FunctionalJobs()
	if len(functional) != 2 {
		t.Errorf("expected 2 functional jobs, got %d", len(functional))
	}

	// Test EmotionalJobs
	emotional := canvas.EmotionalJobs()
	if len(emotional) != 2 {
		t.Errorf("expected 2 emotional jobs, got %d", len(emotional))
	}

	// Test SocialJobs
	social := canvas.SocialJobs()
	if len(social) != 1 {
		t.Errorf("expected 1 social job, got %d", len(social))
	}
}

func TestJTBDCanvasOpportunityAnalysis(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-003", "Opportunity Analysis")
	canvas.DesiredOutcomes = []OutcomeExpectation{
		{ID: "o1", Statement: "Minimize time to complete", Importance: 9, Satisfaction: 3, Opportunity: 15},
		{ID: "o2", Statement: "Maximize accuracy", Importance: 8, Satisfaction: 8, Opportunity: 8},
		{ID: "o3", Statement: "Minimize errors", Importance: 7, Satisfaction: 9, Opportunity: 7},
		{ID: "o4", Statement: "Minimize rework", Importance: 10, Satisfaction: 2, Opportunity: 18},
	}

	// Test TopOpportunities
	top := canvas.TopOpportunities(2)
	if len(top) != 2 {
		t.Errorf("expected 2 top opportunities, got %d", len(top))
	}
	if top[0].ID != "o4" {
		t.Errorf("expected top opportunity to be o4 (score 18), got %s", top[0].ID)
	}
	if top[1].ID != "o1" {
		t.Errorf("expected second opportunity to be o1 (score 15), got %s", top[1].ID)
	}

	// Test TopOpportunities with no limit
	allTop := canvas.TopOpportunities(0)
	if len(allTop) != 4 {
		t.Errorf("expected 4 opportunities with no limit, got %d", len(allTop))
	}

	// Test UnderservedOutcomes
	underserved := canvas.UnderservedOutcomes()
	if len(underserved) != 2 {
		t.Errorf("expected 2 underserved outcomes, got %d", len(underserved))
	}

	// Test OverservedOutcomes
	overserved := canvas.OverservedOutcomes()
	if len(overserved) != 1 {
		t.Errorf("expected 1 overserved outcome, got %d", len(overserved))
	}
	if overserved[0].ID != "o3" {
		t.Errorf("expected overserved outcome to be o3, got %s", overserved[0].ID)
	}
}

func TestJTBDCanvasGetJobStage(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-004", "Job Map Test")
	canvas.JobMap = []JobStage{
		{ID: "s1", Stage: "define", Name: "Define the job"},
		{ID: "s2", Stage: "locate", Name: "Locate resources"},
		{ID: "s3", Stage: "execute", Name: "Execute the job"},
	}

	// Test finding existing stage
	stage := canvas.GetJobStage("s2")
	if stage == nil {
		t.Fatal("expected to find stage s2")
	}
	if stage.Name != "Locate resources" {
		t.Errorf("expected stage name 'Locate resources', got %s", stage.Name)
	}

	// Test not finding non-existent stage
	notFound := canvas.GetJobStage("s99")
	if notFound != nil {
		t.Error("expected GetJobStage to return nil for non-existent stage")
	}
}

func TestJTBDCanvasGetPRDReference(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-005", "PRD Ref Test")

	// Test nil PRDRef
	if canvas.GetPRDReference() != nil {
		t.Error("expected nil PRDReference when not set")
	}

	// Test with PRDRef
	canvas.PRDRef = &PRDReference{
		PRDID:      "prd-001",
		FeatureIDs: []string{"f1", "f2"},
	}
	ref := canvas.GetPRDReference()
	if ref == nil {
		t.Fatal("expected non-nil PRDReference")
	}
	if ref.PRDID != "prd-001" {
		t.Errorf("expected PRDID prd-001, got %s", ref.PRDID)
	}
	if len(ref.FeatureIDs) != 2 {
		t.Errorf("expected 2 feature IDs, got %d", len(ref.FeatureIDs))
	}
}

func TestJTBDCanvasForces(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-006", "Forces Test")
	canvas.PushForces = []Force{
		{ID: "push1", Description: "Current solution is slow", Strength: "strong"},
		{ID: "push2", Description: "Costs too much", Strength: "medium"},
	}
	canvas.PullForces = []Force{
		{ID: "pull1", Description: "New solution is faster", Strength: "strong"},
	}
	canvas.Anxieties = []Force{
		{ID: "anx1", Description: "Learning curve", Strength: "medium"},
	}
	canvas.Habits = []Force{
		{ID: "hab1", Description: "Comfortable with current workflow", Strength: "strong"},
	}

	if len(canvas.PushForces) != 2 {
		t.Errorf("expected 2 push forces, got %d", len(canvas.PushForces))
	}
	if len(canvas.PullForces) != 1 {
		t.Errorf("expected 1 pull force, got %d", len(canvas.PullForces))
	}
	if len(canvas.Anxieties) != 1 {
		t.Errorf("expected 1 anxiety, got %d", len(canvas.Anxieties))
	}
	if len(canvas.Habits) != 1 {
		t.Errorf("expected 1 habit, got %d", len(canvas.Habits))
	}
}

func TestJTBDCanvasHiringSolutions(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-007", "Hiring Solutions Test")
	canvas.HiringSolutions = []HiringSolution{
		{
			ID:                "sol1",
			Name:              "Spreadsheet",
			Type:              "workaround",
			WhyHired:          "Flexible and familiar",
			Limitations:       []string{"Not collaborative", "Error prone"},
			SatisfactionLevel: 5,
		},
		{
			ID:                "sol2",
			Name:              "Legacy software",
			Type:              "product",
			SatisfactionLevel: 3,
		},
	}
	canvas.FiringSolutions = []string{"Manual tracking"}
	canvas.CompetingSolutions = []string{"Competitor A", "Competitor B"}

	if len(canvas.HiringSolutions) != 2 {
		t.Errorf("expected 2 hiring solutions, got %d", len(canvas.HiringSolutions))
	}
	if len(canvas.FiringSolutions) != 1 {
		t.Errorf("expected 1 firing solution, got %d", len(canvas.FiringSolutions))
	}
	if len(canvas.CompetingSolutions) != 2 {
		t.Errorf("expected 2 competing solutions, got %d", len(canvas.CompetingSolutions))
	}
	if canvas.HiringSolutions[0].SatisfactionLevel != 5 {
		t.Errorf("expected satisfaction level 5, got %d", canvas.HiringSolutions[0].SatisfactionLevel)
	}
}

func TestJTBDCanvasInterviews(t *testing.T) {
	canvas := NewJTBDCanvas("jtbd-008", "Interview Test")
	canvas.Interviews = []JTBDInterview{
		{
			ID:              "int1",
			Date:            "2024-01-15",
			ParticipantType: "power user",
			SwitchContext:   "Switched from legacy system",
			FirstThought:    "2023-06-01",
			DecisionMade:    "2023-12-15",
			PurchaseDate:    "2024-01-01",
			PushForces:      []string{"Old system was slow"},
			PullForces:      []string{"New system promised speed"},
			KeyQuotes:       []string{"I finally feel productive"},
		},
	}

	if len(canvas.Interviews) != 1 {
		t.Errorf("expected 1 interview, got %d", len(canvas.Interviews))
	}
	if len(canvas.Interviews[0].PushForces) != 1 {
		t.Errorf("expected 1 push force in interview, got %d", len(canvas.Interviews[0].PushForces))
	}
	if len(canvas.Interviews[0].KeyQuotes) != 1 {
		t.Errorf("expected 1 key quote, got %d", len(canvas.Interviews[0].KeyQuotes))
	}
}

func TestJobTypeConstants(t *testing.T) {
	// Verify job type constants are correct strings
	if JobTypeFunctional != "functional" {
		t.Errorf("expected JobTypeFunctional to be 'functional', got %s", JobTypeFunctional)
	}
	if JobTypeEmotional != "emotional" {
		t.Errorf("expected JobTypeEmotional to be 'emotional', got %s", JobTypeEmotional)
	}
	if JobTypeSocial != "social" {
		t.Errorf("expected JobTypeSocial to be 'social', got %s", JobTypeSocial)
	}
	if JobTypeConsumption != "consumption" {
		t.Errorf("expected JobTypeConsumption to be 'consumption', got %s", JobTypeConsumption)
	}
}
