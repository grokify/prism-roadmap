package integration

import (
	"testing"
)

func TestIdeaToMarketSignal(t *testing.T) {
	idea := IdeaData{
		ID:            "idea-1",
		Source:        "aha",
		Votes:         100,
		CustomerCount: 5,
		ARRImpact:     500000, // $5,000 in cents
	}

	signal := IdeaToMarketSignal(idea)

	if signal.TotalVotes != 100 {
		t.Errorf("TotalVotes = %d, want 100", signal.TotalVotes)
	}
	if signal.CustomerCount != 5 {
		t.Errorf("CustomerCount = %d, want 5", signal.CustomerCount)
	}
	if signal.TotalARR != 500000 {
		t.Errorf("TotalARR = %d, want 500000", signal.TotalARR)
	}
	if signal.IdeaCount != 1 {
		t.Errorf("IdeaCount = %d, want 1", signal.IdeaCount)
	}
	if len(signal.Sources) != 1 || signal.Sources[0] != "aha" {
		t.Errorf("Sources = %v, want [aha]", signal.Sources)
	}

	// Verify score calculation
	// Formula: (votes × 0.1) + (customer_count × 1.0) + (arr / 10_000_000)
	// = (100 × 0.1) + (5 × 1.0) + (500000 / 10_000_000)
	// = 10 + 5 + 0.05 = 15.05
	expectedScore := 15.05
	if signal.Score < 15.04 || signal.Score > 15.06 {
		t.Errorf("Score = %f, want ~%f", signal.Score, expectedScore)
	}
}

func TestIdeasToMarketSignal(t *testing.T) {
	ideas := []IdeaData{
		{ID: "idea-1", Source: "aha", Votes: 50, CustomerCount: 2, ARRImpact: 200000},
		{ID: "idea-2", Source: "aha", Votes: 30, CustomerCount: 3, ARRImpact: 300000},
		{ID: "idea-3", Source: "productboard", Votes: 20, CustomerCount: 1, ARRImpact: 100000},
	}

	signal := IdeasToMarketSignal(ideas)

	if signal.TotalVotes != 100 {
		t.Errorf("TotalVotes = %d, want 100", signal.TotalVotes)
	}
	if signal.CustomerCount != 6 {
		t.Errorf("CustomerCount = %d, want 6", signal.CustomerCount)
	}
	if signal.TotalARR != 600000 {
		t.Errorf("TotalARR = %d, want 600000", signal.TotalARR)
	}
	if signal.IdeaCount != 3 {
		t.Errorf("IdeaCount = %d, want 3", signal.IdeaCount)
	}

	// Should have deduplicated sources
	if len(signal.Sources) != 2 {
		t.Errorf("len(Sources) = %d, want 2", len(signal.Sources))
	}
}

func TestIdeasToMarketSignal_EmptyList(t *testing.T) {
	signal := IdeasToMarketSignal([]IdeaData{})

	if signal.TotalVotes != 0 {
		t.Errorf("TotalVotes = %d, want 0", signal.TotalVotes)
	}
	if signal.IdeaCount != 0 {
		t.Errorf("IdeaCount = %d, want 0", signal.IdeaCount)
	}
	if signal.Score != 0 {
		t.Errorf("Score = %f, want 0", signal.Score)
	}
}

func TestConversionBatch(t *testing.T) {
	batch := NewConversionBatch()

	batch.Add("idea-1", "aha", "rmi-1", "alice")
	batch.Add("idea-2", "aha", "rmi-1", "alice") // same RMI
	batch.Add("idea-3", "productboard", "rmi-2", "bob")

	if len(batch.Conversions) != 3 {
		t.Errorf("len(Conversions) = %d, want 3", len(batch.Conversions))
	}

	// Test GetByIdeaID
	c := batch.GetByIdeaID("idea-1")
	if c == nil {
		t.Fatal("GetByIdeaID(idea-1) = nil, want non-nil")
	}
	if c.RMIID != "rmi-1" {
		t.Errorf("RMIID = %s, want rmi-1", c.RMIID)
	}

	// Test GetByIdeaID not found
	c = batch.GetByIdeaID("nonexistent")
	if c != nil {
		t.Error("GetByIdeaID(nonexistent) should be nil")
	}

	// Test GetByRMIID
	conversions := batch.GetByRMIID("rmi-1")
	if len(conversions) != 2 {
		t.Errorf("GetByRMIID(rmi-1) len = %d, want 2", len(conversions))
	}

	// Test GetByRMIID not found
	conversions = batch.GetByRMIID("nonexistent")
	if len(conversions) != 0 {
		t.Error("GetByRMIID(nonexistent) should return empty slice")
	}
}

func TestRefValidator(t *testing.T) {
	v := NewRefValidator()

	// Register valid IDs
	v.RegisterIdea("idea-1")
	v.RegisterIdea("idea-2")
	v.RegisterGoal("goal-1")
	v.RegisterCapability("cap-1")

	// Test ValidateIdeaRefs - all valid
	result := v.ValidateIdeaRefs([]string{"idea-1", "idea-2"})
	if !result.Valid {
		t.Errorf("ValidateIdeaRefs([idea-1, idea-2]) Valid = false, want true")
	}

	// Test ValidateIdeaRefs - some missing
	result = v.ValidateIdeaRefs([]string{"idea-1", "idea-3", "idea-4"})
	if result.Valid {
		t.Errorf("ValidateIdeaRefs([idea-1, idea-3, idea-4]) Valid = true, want false")
	}
	if len(result.Missing) != 2 {
		t.Errorf("len(Missing) = %d, want 2", len(result.Missing))
	}

	// Test ValidateGoalRefs
	result = v.ValidateGoalRefs([]string{"goal-1"})
	if !result.Valid {
		t.Error("ValidateGoalRefs([goal-1]) Valid = false, want true")
	}

	result = v.ValidateGoalRefs([]string{"goal-2"})
	if result.Valid {
		t.Error("ValidateGoalRefs([goal-2]) Valid = true, want false")
	}

	// Test ValidateCapabilityRefs
	result = v.ValidateCapabilityRefs([]string{"cap-1"})
	if !result.Valid {
		t.Error("ValidateCapabilityRefs([cap-1]) Valid = false, want true")
	}

	result = v.ValidateCapabilityRefs([]string{"cap-2"})
	if result.Valid {
		t.Error("ValidateCapabilityRefs([cap-2]) Valid = true, want false")
	}
}

func TestRefValidator_EmptyRefs(t *testing.T) {
	v := NewRefValidator()

	// Empty refs should be valid
	result := v.ValidateIdeaRefs([]string{})
	if !result.Valid {
		t.Error("ValidateIdeaRefs([]) Valid = false, want true")
	}

	result = v.ValidateGoalRefs(nil)
	if !result.Valid {
		t.Error("ValidateGoalRefs(nil) Valid = false, want true")
	}
}

func TestMarketSignalScoreConsistency(t *testing.T) {
	// Test that the score is consistent whether calculated from
	// a single idea or multiple ideas with the same total

	singleIdea := IdeaData{
		ID:            "idea-1",
		Source:        "aha",
		Votes:         100,
		CustomerCount: 10,
		ARRImpact:     1000000, // $10,000
	}

	multipleIdeas := []IdeaData{
		{ID: "idea-1", Source: "aha", Votes: 30, CustomerCount: 3, ARRImpact: 300000},
		{ID: "idea-2", Source: "aha", Votes: 40, CustomerCount: 4, ARRImpact: 400000},
		{ID: "idea-3", Source: "aha", Votes: 30, CustomerCount: 3, ARRImpact: 300000},
	}

	singleSignal := IdeaToMarketSignal(singleIdea)
	aggregatedSignal := IdeasToMarketSignal(multipleIdeas)

	// Total values should be equal
	if singleSignal.TotalVotes != aggregatedSignal.TotalVotes {
		t.Errorf("TotalVotes: single=%d, aggregated=%d", singleSignal.TotalVotes, aggregatedSignal.TotalVotes)
	}
	if singleSignal.CustomerCount != aggregatedSignal.CustomerCount {
		t.Errorf("CustomerCount: single=%d, aggregated=%d", singleSignal.CustomerCount, aggregatedSignal.CustomerCount)
	}
	if singleSignal.TotalARR != aggregatedSignal.TotalARR {
		t.Errorf("TotalARR: single=%d, aggregated=%d", singleSignal.TotalARR, aggregatedSignal.TotalARR)
	}

	// Scores should be equal
	if singleSignal.Score != aggregatedSignal.Score {
		t.Errorf("Score: single=%f, aggregated=%f", singleSignal.Score, aggregatedSignal.Score)
	}

	// But IdeaCount should differ
	if singleSignal.IdeaCount == aggregatedSignal.IdeaCount {
		t.Errorf("IdeaCount should differ: single=%d, aggregated=%d", singleSignal.IdeaCount, aggregatedSignal.IdeaCount)
	}
}

func TestConversionRoundTrip(t *testing.T) {
	// Simulate a full conversion workflow:
	// 1. Create idea
	// 2. Track conversion
	// 3. Validate RMI references back to idea

	batch := NewConversionBatch()
	validator := NewRefValidator()

	// Register ideas
	ideas := []IdeaData{
		{ID: "idea-1", Source: "aha", Votes: 50, CustomerCount: 2, ARRImpact: 200000},
		{ID: "idea-2", Source: "productboard", Votes: 30, CustomerCount: 1, ARRImpact: 100000},
	}

	for _, idea := range ideas {
		validator.RegisterIdea(idea.ID)
	}

	// Convert ideas to RMI
	rmiID := "rmi-1"
	for _, idea := range ideas {
		batch.Add(idea.ID, idea.Source, rmiID, "system")
	}

	// Validate that RMI can reference its source ideas
	rmiIdeaRefs := []string{}
	for _, c := range batch.GetByRMIID(rmiID) {
		rmiIdeaRefs = append(rmiIdeaRefs, c.IdeaID)
	}

	result := validator.ValidateIdeaRefs(rmiIdeaRefs)
	if !result.Valid {
		t.Errorf("RMI idea refs validation failed: missing=%v", result.Missing)
	}

	// Create market signal from the converted ideas
	signal := IdeasToMarketSignal(ideas)
	if signal.IdeaCount != 2 {
		t.Errorf("MarketSignal IdeaCount = %d, want 2", signal.IdeaCount)
	}
}
