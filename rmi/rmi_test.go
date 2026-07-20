package rmi

import (
	"testing"

	"github.com/grokify/prism-roadmap/effort"
	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/grokify/prism-roadmap/signal"
)

func TestRMIStatusIsActive(t *testing.T) {
	tests := []struct {
		status RMIStatus
		want   bool
	}{
		{RMIStatusPlanned, true},
		{RMIStatusInProgress, true},
		{RMIStatusBlocked, true},
		{RMIStatusCompleted, false},
		{RMIStatusCancelled, false},
		{RMIStatusDeferred, false},
	}

	for _, tt := range tests {
		got := tt.status.IsActive()
		if got != tt.want {
			t.Errorf("RMIStatus(%s).IsActive() = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestRMIStatusIsClosed(t *testing.T) {
	tests := []struct {
		status RMIStatus
		want   bool
	}{
		{RMIStatusCompleted, true},
		{RMIStatusCancelled, true},
		{RMIStatusPlanned, false},
		{RMIStatusInProgress, false},
		{RMIStatusBlocked, false},
		{RMIStatusDeferred, false},
	}

	for _, tt := range tests {
		got := tt.status.IsClosed()
		if got != tt.want {
			t.Errorf("RMIStatus(%s).IsClosed() = %v, want %v", tt.status, got, tt.want)
		}
	}
}

func TestRoadmapItemPriorityScore(t *testing.T) {
	// Simple test with MoSCoW only
	item := NewRoadmapItem("rmi-1", "Test Item", prioritization.MoSCoWMustHave)
	score := item.PriorityScore()
	if score != 4.0 { // MustHave weight = 4
		t.Errorf("PriorityScore() = %f, want 4.0", score)
	}

	// With RICE score
	rice := prioritization.NewRICEScore("rmi-1", 1000, prioritization.ImpactHigh, prioritization.ConfidenceHigh, 2)
	item.WithRICE(rice)
	scoreWithRICE := item.PriorityScore()
	if scoreWithRICE <= 4.0 {
		t.Errorf("PriorityScore with RICE = %f, should be > 4.0", scoreWithRICE)
	}

	// With market signal
	ms := signal.NewMarketSignalFromIdea(100, 5, 500000, "aha")
	item.WithMarketSignal(ms)
	scoreWithMS := item.PriorityScore()
	if scoreWithMS <= scoreWithRICE {
		t.Errorf("PriorityScore with market signal = %f, should be > %f", scoreWithMS, scoreWithRICE)
	}

	// With complexity (should reduce score)
	complexity := &effort.ComplexityFactors{
		NewArchitecture: true,
		NewDesignUX:     true,
	}
	item.WithComplexity(complexity)
	scoreWithComplexity := item.PriorityScore()
	if scoreWithComplexity >= scoreWithMS {
		t.Errorf("PriorityScore with complexity = %f, should be < %f", scoreWithComplexity, scoreWithMS)
	}
}

func TestRoadmapItemEffectiveEffortDays(t *testing.T) {
	item := NewRoadmapItem("rmi-1", "Test", prioritization.MoSCoWMustHave)

	// No effort set
	if days := item.EffectiveEffortDays(); days != 0 {
		t.Errorf("EffectiveEffortDays() = %d, want 0", days)
	}

	// With effort
	item.WithEffort(effort.NewEffortEstimate(15, effort.ConfidenceMedium))
	if days := item.EffectiveEffortDays(); days != 15 {
		t.Errorf("EffectiveEffortDays() = %d, want 15", days)
	}
}

func TestRoadmapItemHasBlockingDependencies(t *testing.T) {
	item := NewRoadmapItem("rmi-1", "Test", prioritization.MoSCoWMustHave)

	// No complexity
	if item.HasBlockingDependencies() {
		t.Error("HasBlockingDependencies() = true, want false (no complexity)")
	}

	// With non-blocking dependency
	item.WithComplexity(&effort.ComplexityFactors{
		Dependencies: []effort.Dependency{
			{TeamID: "team-1", Type: effort.DependencyTypeInformational},
		},
	})
	if item.HasBlockingDependencies() {
		t.Error("HasBlockingDependencies() = true, want false (informational)")
	}

	// With blocking dependency
	item.WithComplexity(&effort.ComplexityFactors{
		Dependencies: []effort.Dependency{
			{TeamID: "team-1", Type: effort.DependencyTypeBlocking, Status: effort.DependencyStatusPending},
		},
	})
	if !item.HasBlockingDependencies() {
		t.Error("HasBlockingDependencies() = false, want true")
	}

	// With resolved blocking dependency
	item.WithComplexity(&effort.ComplexityFactors{
		Dependencies: []effort.Dependency{
			{TeamID: "team-1", Type: effort.DependencyTypeBlocking, Status: effort.DependencyStatusResolved},
		},
	})
	if item.HasBlockingDependencies() {
		t.Error("HasBlockingDependencies() = true, want false (resolved)")
	}
}

func TestRoadmapItemIsActionable(t *testing.T) {
	// Actionable item
	item := NewRoadmapItem("rmi-1", "Test", prioritization.MoSCoWMustHave)
	if !item.IsActionable() {
		t.Error("IsActionable() = false, want true")
	}

	// Won't have is not actionable
	wontHave := NewRoadmapItem("rmi-2", "Test", prioritization.MoSCoWWontHave)
	if wontHave.IsActionable() {
		t.Error("Won't have item IsActionable() = true, want false")
	}

	// Completed is not actionable
	completed := NewRoadmapItem("rmi-3", "Test", prioritization.MoSCoWMustHave)
	completed.Status = RMIStatusCompleted
	if completed.IsActionable() {
		t.Error("Completed item IsActionable() = true, want false")
	}

	// Item with blocking deps is not actionable
	blocked := NewRoadmapItem("rmi-4", "Test", prioritization.MoSCoWMustHave)
	blocked.WithComplexity(&effort.ComplexityFactors{
		Dependencies: []effort.Dependency{
			{TeamID: "team-1", Type: effort.DependencyTypeBlocking, Status: effort.DependencyStatusPending},
		},
	})
	if blocked.IsActionable() {
		t.Error("Blocked item IsActionable() = true, want false")
	}
}

func TestRoadmapItemValidate(t *testing.T) {
	tests := []struct {
		name    string
		item    RoadmapItem
		wantErr bool
	}{
		{
			name: "valid",
			item: RoadmapItem{
				ID:     "rmi-1",
				Name:   "Test Item",
				MoSCoW: prioritization.MoSCoWMustHave,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			item: RoadmapItem{
				Name:   "Test Item",
				MoSCoW: prioritization.MoSCoWMustHave,
			},
			wantErr: true,
		},
		{
			name: "missing name",
			item: RoadmapItem{
				ID:     "rmi-1",
				MoSCoW: prioritization.MoSCoWMustHave,
			},
			wantErr: true,
		},
		{
			name: "missing MoSCoW",
			item: RoadmapItem{
				ID:   "rmi-1",
				Name: "Test Item",
			},
			wantErr: true,
		},
		{
			name: "invalid MoSCoW",
			item: RoadmapItem{
				ID:     "rmi-1",
				Name:   "Test Item",
				MoSCoW: "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRoadmapItem(t *testing.T) {
	item := NewRoadmapItem("rmi-1", "SSO Integration", prioritization.MoSCoWMustHave)

	if item.ID != "rmi-1" {
		t.Errorf("ID = %s, want rmi-1", item.ID)
	}
	if item.Name != "SSO Integration" {
		t.Errorf("Name = %s, want SSO Integration", item.Name)
	}
	if item.MoSCoW != prioritization.MoSCoWMustHave {
		t.Errorf("MoSCoW = %s, want must_have", item.MoSCoW)
	}
	if item.Status != RMIStatusPlanned {
		t.Errorf("Status = %s, want planned", item.Status)
	}
	if item.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
}

func TestRoadmapItemAddRefs(t *testing.T) {
	item := NewRoadmapItem("rmi-1", "Test", prioritization.MoSCoWMustHave)

	item.AddCapabilityRef("cap-1")
	item.AddCapabilityRef("cap-2")
	if len(item.CapabilityRefs) != 2 {
		t.Errorf("CapabilityRefs length = %d, want 2", len(item.CapabilityRefs))
	}

	item.AddGoalRef("goal-1")
	if len(item.GoalRefs) != 1 {
		t.Errorf("GoalRefs length = %d, want 1", len(item.GoalRefs))
	}

	item.AddIdeaRef("idea-1")
	item.AddIdeaRef("idea-2")
	item.AddIdeaRef("idea-3")
	if len(item.IdeaRefs) != 3 {
		t.Errorf("IdeaRefs length = %d, want 3", len(item.IdeaRefs))
	}
}

func TestValidRMIStatuses(t *testing.T) {
	statuses := ValidRMIStatuses()
	if len(statuses) != 6 {
		t.Errorf("ValidRMIStatuses() returned %d statuses, want 6", len(statuses))
	}
}
