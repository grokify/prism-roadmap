package prd

import (
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/roadmap"
)

func TestToSwimlaneTable(t *testing.T) {
	roadmap := Roadmap{
		Phases: []Phase{
			{
				ID:   "phase-1",
				Name: "MVP",
				Deliverables: []Deliverable{
					{ID: "d1", Title: "User Auth", Type: DeliverableFeature, Status: DeliverableCompleted},
					{ID: "d2", Title: "CI/CD Pipeline", Type: DeliverableInfrastructure, Status: DeliverableInProgress},
					{ID: "d3", Title: "API Docs", Type: DeliverableDocumentation},
				},
			},
			{
				ID:   "phase-2",
				Name: "Beta",
				Deliverables: []Deliverable{
					{ID: "d4", Title: "Dashboard", Type: DeliverableFeature, Status: DeliverableNotStarted},
					{ID: "d5", Title: "Monitoring", Type: DeliverableInfrastructure},
					{ID: "d6", Title: "Stripe Integration", Type: DeliverableIntegration},
				},
			},
			{
				ID:   "phase-3",
				Name: "GA",
				Deliverables: []Deliverable{
					{ID: "d7", Title: "Analytics", Type: DeliverableFeature},
					{ID: "d8", Title: "GA Release", Type: DeliverableMilestone, Status: DeliverableNotStarted},
				},
			},
		},
	}

	opts := DefaultRoadmapTableOptions()
	table := roadmap.ToSwimlaneTable(opts)

	// Check header - now uses **Phase N**<br>Name format
	if !strings.Contains(table, "| Swimlane |") {
		t.Error("Expected Swimlane header")
	}
	if !strings.Contains(table, "**Phase 1**<br>MVP") {
		t.Error("Expected Phase 1 with MVP description")
	}
	if !strings.Contains(table, "**Phase 2**<br>Beta") {
		t.Error("Expected Phase 2 with Beta description")
	}
	if !strings.Contains(table, "**Phase 3**<br>GA") {
		t.Error("Expected Phase 3 with GA description")
	}

	// Check swimlane rows
	if !strings.Contains(table, "**Features**") {
		t.Error("Expected Features swimlane")
	}
	if !strings.Contains(table, "**Infrastructure**") {
		t.Error("Expected Infrastructure swimlane")
	}
	if !strings.Contains(table, "**Documentation**") {
		t.Error("Expected Documentation swimlane")
	}
	if !strings.Contains(table, "**Integrations**") {
		t.Error("Expected Integrations swimlane")
	}
	if !strings.Contains(table, "**Milestones**") {
		t.Error("Expected Milestones swimlane")
	}

	// Check status icons
	if !strings.Contains(table, "✅") {
		t.Error("Expected completed status icon")
	}
	if !strings.Contains(table, "🔄") {
		t.Error("Expected in-progress status icon")
	}
	if !strings.Contains(table, "⏳") {
		t.Error("Expected not-started status icon")
	}

	// Check deliverable content
	if !strings.Contains(table, "User Auth") {
		t.Error("Expected User Auth deliverable")
	}
	if !strings.Contains(table, "Dashboard") {
		t.Error("Expected Dashboard deliverable")
	}
}

func TestToSwimlaneTableEmpty(t *testing.T) {
	roadmap := Roadmap{}
	table := roadmap.ToSwimlaneTable(DefaultRoadmapTableOptions())
	if table != "" {
		t.Error("Expected empty string for empty roadmap")
	}
}

func TestToSwimlaneTableNoStatus(t *testing.T) {
	roadmap := Roadmap{
		Phases: []Phase{
			{
				ID:   "phase-1",
				Name: "Phase 1",
				Deliverables: []Deliverable{
					{ID: "d1", Title: "Feature A", Type: DeliverableFeature, Status: DeliverableCompleted},
				},
			},
		},
	}

	opts := DefaultRoadmapTableOptions()
	opts.IncludeStatus = false
	table := roadmap.ToSwimlaneTable(opts)

	// Should not contain status icons
	if strings.Contains(table, "✅") || strings.Contains(table, "🔄") {
		t.Error("Expected no status icons when IncludeStatus is false")
	}
	// Should still contain the deliverable
	if !strings.Contains(table, "Feature A") {
		t.Error("Expected Feature A deliverable")
	}
}

func TestToPhaseTable(t *testing.T) {
	roadmap := Roadmap{
		Phases: []Phase{
			{
				ID:     "phase-1",
				Name:   "MVP",
				Status: PhaseStatusInProgress,
				Deliverables: []Deliverable{
					{ID: "d1", Title: "User Auth", Type: DeliverableFeature, Status: DeliverableCompleted},
					{ID: "d2", Title: "API Docs", Type: DeliverableDocumentation},
				},
			},
			{
				ID:     "phase-2",
				Name:   "Beta",
				Status: PhaseStatusPlanned,
				Deliverables: []Deliverable{
					{ID: "d3", Title: "Dashboard", Type: DeliverableFeature},
				},
			},
		},
	}

	opts := DefaultRoadmapTableOptions()
	table := roadmap.ToPhaseTable(opts)

	// Check header
	if !strings.Contains(table, "| Phase | Status | Deliverables |") {
		t.Error("Expected table header")
	}

	// Check phases - now uses **Phase N**<br>Name format
	if !strings.Contains(table, "**Phase 1**<br>MVP") {
		t.Error("Expected Phase 1 with MVP description")
	}
	if !strings.Contains(table, "**Phase 2**<br>Beta") {
		t.Error("Expected Phase 2 with Beta description")
	}

	// Check deliverables use bullet list format
	if !strings.Contains(table, "• ") {
		t.Error("Expected bullet points in deliverables")
	}
	if !strings.Contains(table, "User Auth") {
		t.Error("Expected User Auth")
	}
	if !strings.Contains(table, "API Docs") {
		t.Error("Expected API Docs")
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
		{DeliverableType("custom"), "Custom"},
	}

	for _, tt := range tests {
		got := roadmap.SwimlaneLabel(tt.input)
		if got != tt.expected {
			t.Errorf("SwimlaneLabel(%q) = %q, want %q", tt.input, got, tt.expected)
		}
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
		{DeliverableStatus("unknown"), ""},
	}

	for _, tt := range tests {
		got := roadmap.StatusIcon(tt.input)
		if got != tt.expected {
			t.Errorf("StatusIcon(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestStatusLegend(t *testing.T) {
	legend := StatusLegend()

	// Check it's a valid markdown table
	if !strings.Contains(legend, "| Icon | Status |") {
		t.Error("Legend should have Icon and Status headers")
	}

	// Check all icons are present
	icons := []string{"✅", "🔄", "⏳", "🚫", "❌"}
	for _, icon := range icons {
		if !strings.Contains(legend, icon) {
			t.Errorf("Legend should contain icon %s", icon)
		}
	}

	// Check status labels
	labels := []string{"Completed", "In Progress", "Not Started", "Blocked", "Missed"}
	for _, label := range labels {
		if !strings.Contains(legend, label) {
			t.Errorf("Legend should contain label %q", label)
		}
	}
}
