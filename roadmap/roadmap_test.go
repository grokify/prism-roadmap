package roadmap

import (
	"strings"
	"testing"
)

func TestToSwimlaneTable(t *testing.T) {
	r := createTestRoadmap()
	opts := DefaultTableOptions()

	table := r.ToSwimlaneTable(opts)

	// Check header contains phase names
	if !strings.Contains(table, "**Phase 1**") {
		t.Error("Missing Phase 1 header")
	}
	if !strings.Contains(table, "**Phase 2**") {
		t.Error("Missing Phase 2 header")
	}

	// Check swimlane labels
	if !strings.Contains(table, "**Features**") {
		t.Error("Missing Features swimlane")
	}
	if !strings.Contains(table, "**Infrastructure**") {
		t.Error("Missing Infrastructure swimlane")
	}

	// Check deliverable titles
	if !strings.Contains(table, "Auth") {
		t.Error("Missing Auth deliverable")
	}
	if !strings.Contains(table, "Dashboard") {
		t.Error("Missing Dashboard deliverable")
	}
}

func TestToSwimlaneTableEmpty(t *testing.T) {
	r := &Roadmap{}
	opts := DefaultTableOptions()

	table := r.ToSwimlaneTable(opts)

	if table != "" {
		t.Errorf("Expected empty table for empty roadmap, got: %s", table)
	}
}

func TestToSwimlaneTableWithStatus(t *testing.T) {
	r := createTestRoadmap()
	opts := DefaultTableOptions()
	opts.IncludeStatus = true

	table := r.ToSwimlaneTable(opts)

	// Check status icons
	if !strings.Contains(table, "✅") {
		t.Error("Missing completed status icon")
	}
	if !strings.Contains(table, "🔄") {
		t.Error("Missing in-progress status icon")
	}
}

func TestToPhaseTable(t *testing.T) {
	r := createTestRoadmap()
	opts := DefaultTableOptions()

	table := r.ToPhaseTable(opts)

	// Check header
	if !strings.Contains(table, "| Phase | Status | Deliverables |") {
		t.Error("Missing table header")
	}

	// Check phase names
	if !strings.Contains(table, "Foundation") {
		t.Error("Missing Foundation phase")
	}
	if !strings.Contains(table, "Core Features") {
		t.Error("Missing Core Features phase")
	}

	// Check status
	if !strings.Contains(table, "in_progress") {
		t.Error("Missing in_progress status")
	}
}

func TestToPhaseTableEmpty(t *testing.T) {
	r := &Roadmap{}
	opts := DefaultTableOptions()

	table := r.ToPhaseTable(opts)

	if table != "" {
		t.Errorf("Expected empty table for empty roadmap, got: %s", table)
	}
}

func TestSwimlaneLabel(t *testing.T) {
	tests := []struct {
		input    DeliverableType
		expected string
	}{
		{DeliverableFeature, "Features"},
		{DeliverableDocumentation, "Documentation"},
		{DeliverableInfrastructure, "Infrastructure"},
		{DeliverableIntegration, "Integrations"},
		{DeliverableMilestone, "Milestones"},
		{DeliverableRollout, "Rollout"},
		{DeliverableType("custom"), "Custom"},
		{DeliverableType(""), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			result := SwimlaneLabel(tt.input)
			if result != tt.expected {
				t.Errorf("SwimlaneLabel(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStatusIcon(t *testing.T) {
	tests := []struct {
		input    DeliverableStatus
		expected string
	}{
		{DeliverableCompleted, "✅"},
		{DeliverableInProgress, "🔄"},
		{DeliverableBlocked, "🚫"},
		{DeliverableNotStarted, "⏳"},
		{DeliverableStatus(""), ""},
		{DeliverableStatus("unknown"), ""},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			result := StatusIcon(tt.input)
			if result != tt.expected {
				t.Errorf("StatusIcon(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestPhaseTargetStatusIcon(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"achieved", "✅"},
		{"in_progress", "🔄"},
		{"missed", "❌"},
		{"not_started", "⏳"},
		{"", ""},
		{"unknown", ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := PhaseTargetStatusIcon(tt.input)
			if result != tt.expected {
				t.Errorf("PhaseTargetStatusIcon(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestStatusLegend(t *testing.T) {
	legend := StatusLegend()

	// Check it contains expected icons
	if !strings.Contains(legend, "✅") {
		t.Error("Legend missing completed icon")
	}
	if !strings.Contains(legend, "🔄") {
		t.Error("Legend missing in-progress icon")
	}
	if !strings.Contains(legend, "⏳") {
		t.Error("Legend missing not-started icon")
	}
	if !strings.Contains(legend, "🚫") {
		t.Error("Legend missing blocked icon")
	}
	if !strings.Contains(legend, "❌") {
		t.Error("Legend missing missed icon")
	}
}

func TestMaxTitleLen(t *testing.T) {
	r := &Roadmap{
		Phases: []Phase{
			{
				ID:   "p1",
				Name: "Test Phase",
				Deliverables: []Deliverable{
					{
						ID:    "d1",
						Title: "This is a very long deliverable title that should be truncated",
						Type:  DeliverableFeature,
					},
				},
			},
		},
	}

	opts := TableOptions{
		IncludeStatus: false,
		MaxTitleLen:   20,
	}

	table := r.ToSwimlaneTable(opts)

	// Should be truncated with "..."
	if !strings.Contains(table, "...") {
		t.Error("Expected truncated title with ellipsis")
	}
	// Should not contain the full title
	if strings.Contains(table, "that should be truncated") {
		t.Error("Title was not truncated")
	}
}

func TestRolloutStatusDeploymentPercent(t *testing.T) {
	tests := []struct {
		name     string
		rollout  *RolloutStatus
		expected float64
	}{
		{"nil rollout", nil, 0},
		{"zero total", &RolloutStatus{TotalCustomers: 0, DeployedCustomers: 5}, 0},
		{"50 percent", &RolloutStatus{TotalCustomers: 100, DeployedCustomers: 50}, 50},
		{"100 percent", &RolloutStatus{TotalCustomers: 50, DeployedCustomers: 50}, 100},
		{"partial", &RolloutStatus{TotalCustomers: 80, DeployedCustomers: 45}, 56.25},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rollout.DeploymentPercent()
			if result != tt.expected {
				t.Errorf("DeploymentPercent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRolloutStatusAdoptionPercent(t *testing.T) {
	tests := []struct {
		name     string
		rollout  *RolloutStatus
		expected float64
	}{
		{"nil rollout", nil, 0},
		{"zero total", &RolloutStatus{TotalCustomers: 0, AdoptedCustomers: 5}, 0},
		{"30 percent", &RolloutStatus{TotalCustomers: 100, AdoptedCustomers: 30}, 30},
		{"full adoption", &RolloutStatus{TotalCustomers: 50, AdoptedCustomers: 50}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rollout.AdoptionPercent()
			if result != tt.expected {
				t.Errorf("AdoptionPercent() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRolloutStatusAdoptionOfDeployed(t *testing.T) {
	tests := []struct {
		name     string
		rollout  *RolloutStatus
		expected float64
	}{
		{"nil rollout", nil, 0},
		{"zero deployed", &RolloutStatus{DeployedCustomers: 0, AdoptedCustomers: 5}, 0},
		{"50 percent of deployed", &RolloutStatus{DeployedCustomers: 80, AdoptedCustomers: 40}, 50},
		{"full adoption of deployed", &RolloutStatus{DeployedCustomers: 50, AdoptedCustomers: 50}, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rollout.AdoptionOfDeployed()
			if result != tt.expected {
				t.Errorf("AdoptionOfDeployed() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestRolloutStatusIsFullyDeployed(t *testing.T) {
	tests := []struct {
		name     string
		rollout  *RolloutStatus
		expected bool
	}{
		{"nil rollout", nil, false},
		{"not fully deployed", &RolloutStatus{TotalCustomers: 100, DeployedCustomers: 50}, false},
		{"fully deployed", &RolloutStatus{TotalCustomers: 50, DeployedCustomers: 50}, true},
		{"over deployed", &RolloutStatus{TotalCustomers: 50, DeployedCustomers: 55}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.rollout.IsFullyDeployed()
			if result != tt.expected {
				t.Errorf("IsFullyDeployed() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestDeliverableWithRollout(t *testing.T) {
	d := Deliverable{
		ID:     "feature-1",
		Title:  "New Dashboard",
		Type:   DeliverableFeature,
		Status: DeliverableInProgress,
		Rollout: &RolloutStatus{
			TotalCustomers:    100,
			DeployedCustomers: 75,
			AdoptedCustomers:  45,
			Status:            RolloutStageRollingOut,
			Waves: []RolloutWave{
				{ID: "wave-1", Name: "Beta", TargetCustomers: 10, DeployedCustomers: 10, Status: "completed"},
				{ID: "wave-2", Name: "GA Wave 1", TargetCustomers: 40, DeployedCustomers: 40, Status: "completed"},
				{ID: "wave-3", Name: "GA Wave 2", TargetCustomers: 50, DeployedCustomers: 25, Status: "in_progress"},
			},
		},
	}

	if d.Rollout.DeploymentPercent() != 75 {
		t.Errorf("Expected 75%% deployed, got %v%%", d.Rollout.DeploymentPercent())
	}

	if d.Rollout.AdoptionPercent() != 45 {
		t.Errorf("Expected 45%% adoption, got %v%%", d.Rollout.AdoptionPercent())
	}

	if d.Rollout.AdoptionOfDeployed() != 60 {
		t.Errorf("Expected 60%% adoption of deployed, got %v%%", d.Rollout.AdoptionOfDeployed())
	}

	if len(d.Rollout.Waves) != 3 {
		t.Errorf("Expected 3 waves, got %d", len(d.Rollout.Waves))
	}
}

func createTestRoadmap() *Roadmap {
	return &Roadmap{
		Phases: []Phase{
			{
				ID:     "phase-1",
				Name:   "Foundation",
				Type:   PhaseTypeGeneric,
				Status: PhaseStatusInProgress,
				Goals:  []string{"Establish base infrastructure"},
				Deliverables: []Deliverable{
					{
						ID:     "d1",
						Title:  "Auth",
						Type:   DeliverableFeature,
						Status: DeliverableCompleted,
					},
					{
						ID:     "d2",
						Title:  "CI/CD",
						Type:   DeliverableInfrastructure,
						Status: DeliverableInProgress,
					},
				},
				SuccessCriteria: []string{"All tests passing"},
			},
			{
				ID:     "phase-2",
				Name:   "Core Features",
				Type:   PhaseTypeGeneric,
				Status: PhaseStatusPlanned,
				Goals:  []string{"Build core functionality"},
				Deliverables: []Deliverable{
					{
						ID:     "d3",
						Title:  "Dashboard",
						Type:   DeliverableFeature,
						Status: DeliverableNotStarted,
					},
					{
						ID:     "d4",
						Title:  "Monitoring",
						Type:   DeliverableInfrastructure,
						Status: DeliverableNotStarted,
					},
				},
				SuccessCriteria: []string{"Feature complete"},
			},
		},
	}
}
