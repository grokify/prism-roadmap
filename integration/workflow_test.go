package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/grokify/prism-roadmap/rmi"
	"github.com/grokify/prism-roadmap/signal"
)

// TestFullWorkflow tests the complete workflow from ideas to prioritized roadmap items.
func TestFullWorkflow(t *testing.T) {
	// Step 1: Start with ideas from multiple sources
	ideas := []IdeaData{
		{
			ID:            "idea-1",
			Source:        "aha",
			Votes:         100,
			CustomerCount: 10,
			ARRImpact:     500000, // $5,000
			CustomerIDs:   []string{"cust-1", "cust-2"},
		},
		{
			ID:            "idea-2",
			Source:        "aha",
			Votes:         50,
			CustomerCount: 5,
			ARRImpact:     200000, // $2,000
			CustomerIDs:   []string{"cust-1"},
		},
		{
			ID:            "idea-3",
			Source:        "productboard",
			Votes:         75,
			CustomerCount: 8,
			ARRImpact:     350000, // $3,500
			CustomerIDs:   []string{"cust-3"},
		},
	}

	// Step 2: Aggregate ideas into market signals
	combinedSignal := IdeasToMarketSignal(ideas)

	if combinedSignal.TotalVotes != 225 {
		t.Errorf("TotalVotes = %d, want 225", combinedSignal.TotalVotes)
	}
	if combinedSignal.CustomerCount != 23 {
		t.Errorf("CustomerCount = %d, want 23", combinedSignal.CustomerCount)
	}
	if combinedSignal.TotalARR != 1050000 {
		t.Errorf("TotalARR = %d, want 1050000", combinedSignal.TotalARR)
	}

	// Step 3: Create RMI service and track conversions
	svc := rmi.NewService()
	batch := NewConversionBatch()

	// Step 4: Create roadmap item from aggregated ideas
	rmiInput := rmi.CreateInput{
		ID:          "rmi-1",
		Name:        "API Rate Limiting",
		Description: "Implement rate limiting based on customer feedback",
		MoSCoW:      prioritization.MoSCoWMustHave,
		Quarter:     "Q3 2025",
		Owner:       "platform-team",
		Tags:        []string{"api", "performance"},
	}

	item, err := svc.Create(rmiInput)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Step 5: Attach market signal to RMI
	item.MarketSignal = combinedSignal

	// Track conversions
	for _, idea := range ideas {
		batch.Add(idea.ID, idea.Source, item.ID, "workflow-test")
	}

	// Step 6: Validate conversions
	validator := NewRefValidator()
	for _, idea := range ideas {
		validator.RegisterIdea(idea.ID)
	}

	// Get idea refs for the RMI
	rmiConversions := batch.GetByRMIID(item.ID)
	ideaRefs := make([]string, 0, len(rmiConversions))
	for _, c := range rmiConversions {
		ideaRefs = append(ideaRefs, c.IdeaID)
	}

	refResult := validator.ValidateIdeaRefs(ideaRefs)
	if !refResult.Valid {
		t.Errorf("Idea refs validation failed: missing=%v", refResult.Missing)
	}

	// Step 7: Verify RMI has correct market signal
	if item.MarketSignal.Score <= 0 {
		t.Error("MarketSignal.Score should be positive")
	}
	if len(batch.Conversions) != 3 {
		t.Errorf("len(Conversions) = %d, want 3", len(batch.Conversions))
	}
}

// TestPrioritizationWorkflow tests the prioritization workflow.
func TestPrioritizationWorkflow(t *testing.T) {
	svc := rmi.NewService()

	// Create items with different priorities
	items := []rmi.CreateInput{
		{
			ID:     "rmi-must",
			Name:   "Critical Security Fix",
			MoSCoW: prioritization.MoSCoWMustHave,
		},
		{
			ID:     "rmi-should",
			Name:   "Performance Improvement",
			MoSCoW: prioritization.MoSCoWShouldHave,
		},
		{
			ID:     "rmi-could",
			Name:   "Nice-to-have Feature",
			MoSCoW: prioritization.MoSCoWCouldHave,
		},
	}

	for _, input := range items {
		if _, err := svc.Create(input); err != nil {
			t.Fatalf("Create(%s) error = %v", input.ID, err)
		}
	}

	// Attach market signals
	mustItem, _ := svc.Get("rmi-must")
	mustItem.MarketSignal = signal.NewMarketSignalFromIdea(200, 20, 1000000, "aha")

	shouldItem, _ := svc.Get("rmi-should")
	shouldItem.MarketSignal = signal.NewMarketSignalFromIdea(100, 10, 500000, "aha")

	couldItem, _ := svc.Get("rmi-could")
	couldItem.MarketSignal = signal.NewMarketSignalFromIdea(50, 5, 100000, "aha")

	// Get top by market signal
	topItems := svc.TopByMarketSignal(3)

	if len(topItems) != 3 {
		t.Fatalf("len(TopByMarketSignal) = %d, want 3", len(topItems))
	}

	// Verify ordering by market signal score
	if topItems[0].ID != "rmi-must" {
		t.Errorf("topItems[0].ID = %s, want rmi-must", topItems[0].ID)
	}
}

// TestRoadmapPersistence tests saving and loading roadmap data.
func TestRoadmapPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "workflow-test-*")
	if err != nil {
		t.Fatalf("creating temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "roadmap.json")

	// Create and save
	svc := rmi.NewService()
	_, err = svc.Create(rmi.CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWMustHave,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := svc.SaveAs(filePath); err != nil {
		t.Fatalf("SaveAs() error = %v", err)
	}

	// Load and verify
	svc2, err := rmi.NewServiceFromFile(filePath)
	if err != nil {
		t.Fatalf("NewServiceFromFile() error = %v", err)
	}

	item, err := svc2.Get("rmi-1")
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}

	if item.Name != "Test Item" {
		t.Errorf("Name = %s, want Test Item", item.Name)
	}
}

// TestStatusTransitions tests RMI status transitions.
func TestStatusTransitions(t *testing.T) {
	svc := rmi.NewService()

	_, err := svc.Create(rmi.CreateInput{
		ID:     "rmi-1",
		Name:   "Test Item",
		MoSCoW: prioritization.MoSCoWMustHave,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	// Verify initial status
	item, _ := svc.Get("rmi-1")
	if item.Status != rmi.RMIStatusPlanned {
		t.Errorf("Initial Status = %s, want %s", item.Status, rmi.RMIStatusPlanned)
	}

	// Transition to in_progress
	inProgress := rmi.RMIStatusInProgress
	item, _, _ = svc.Update("rmi-1", rmi.UpdateInput{Status: &inProgress})
	if item.Status != rmi.RMIStatusInProgress {
		t.Errorf("Status = %s, want %s", item.Status, rmi.RMIStatusInProgress)
	}

	// Transition to completed
	completed := rmi.RMIStatusCompleted
	item, _, _ = svc.Update("rmi-1", rmi.UpdateInput{Status: &completed})
	if item.Status != rmi.RMIStatusCompleted {
		t.Errorf("Status = %s, want %s", item.Status, rmi.RMIStatusCompleted)
	}
}

// TestSummaryAggregation tests the summary aggregation across items.
func TestSummaryAggregation(t *testing.T) {
	svc := rmi.NewService()

	// Create items with different MoSCoW priorities
	moscowInputs := []rmi.CreateInput{
		{ID: "must-1", Name: "Must 1", MoSCoW: prioritization.MoSCoWMustHave},
		{ID: "must-2", Name: "Must 2", MoSCoW: prioritization.MoSCoWMustHave},
		{ID: "should-1", Name: "Should 1", MoSCoW: prioritization.MoSCoWShouldHave},
		{ID: "could-1", Name: "Could 1", MoSCoW: prioritization.MoSCoWCouldHave},
		{ID: "wont-1", Name: "Wont 1", MoSCoW: prioritization.MoSCoWWontHave},
	}

	for _, input := range moscowInputs {
		if _, err := svc.Create(input); err != nil {
			t.Fatalf("Create(%s) error = %v", input.ID, err)
		}
	}

	summary := svc.Summary()

	if summary.TotalItems != 5 {
		t.Errorf("TotalItems = %d, want 5", summary.TotalItems)
	}
	if summary.MoSCoWCounts[prioritization.MoSCoWMustHave] != 2 {
		t.Errorf("MoSCoWCounts[must_have] = %d, want 2", summary.MoSCoWCounts[prioritization.MoSCoWMustHave])
	}
	if summary.MoSCoWCounts[prioritization.MoSCoWShouldHave] != 1 {
		t.Errorf("MoSCoWCounts[should_have] = %d, want 1", summary.MoSCoWCounts[prioritization.MoSCoWShouldHave])
	}
}
