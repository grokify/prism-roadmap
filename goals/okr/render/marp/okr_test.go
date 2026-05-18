package marp

import (
	"strings"
	"testing"

	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/goals/okr/render"
)

func TestRenderer_Format(t *testing.T) {
	r := New()
	if got := r.Format(); got != "marp" {
		t.Errorf("Format() = %v, want %v", got, "marp")
	}
}

func TestRenderer_FileExtension(t *testing.T) {
	r := New()
	if got := r.FileExtension(); got != ".md" {
		t.Errorf("FileExtension() = %v, want %v", got, ".md")
	}
}

func TestRenderer_Render(t *testing.T) {
	doc := createTestOKR()
	r := New()

	output, err := r.Render(doc, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check Marp front matter
	if !strings.Contains(content, "marp: true") {
		t.Error("Missing marp: true in front matter")
	}
	if !strings.Contains(content, "paginate: true") {
		t.Error("Missing paginate: true in front matter")
	}

	// Check title slide
	if !strings.Contains(content, "Test OKR") {
		t.Error("Missing title in slides")
	}
	if !strings.Contains(content, "Objectives and Key Results") {
		t.Error("Missing OKR subtitle")
	}

	// Check overview slide
	if !strings.Contains(content, "OKR Overview") {
		t.Error("Missing overview slide")
	}

	// Check objectives
	if !strings.Contains(content, "Improve Product Quality") {
		t.Error("Missing objective title")
	}

	// Check key results
	if !strings.Contains(content, "Key Results") {
		t.Error("Missing key results section")
	}
	if !strings.Contains(content, "Increase test coverage") {
		t.Error("Missing key result")
	}

	// Check summary slide
	if !strings.Contains(content, "Summary") {
		t.Error("Missing summary slide")
	}
}

func TestRenderer_RenderWithOptions(t *testing.T) {
	doc := createTestOKR()
	r := New()

	opts := &render.Options{
		Theme:            "corporate",
		IncludeRisks:     false,
		ShowScoreGrades:  true,
		ShowProgressBars: true,
	}

	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check corporate theme colors are used
	if !strings.Contains(content, "#1a365d") {
		t.Error("Missing corporate theme primary color")
	}
}

func TestRenderer_RenderWithRisks(t *testing.T) {
	doc := createTestOKR()
	doc.Risks = []okr.Risk{
		{
			ID:          "R-1",
			Title:       "Resource constraints",
			Description: "Team may be understaffed",
			Impact:      "High",
			Likelihood:  "Medium",
			Mitigation:  "Hire contractors",
		},
	}

	r := New()
	opts := render.DefaultOptions()
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check risks slide
	if !strings.Contains(content, "Risks") {
		t.Error("Missing risks slide")
	}
	if !strings.Contains(content, "Resource constraints") {
		t.Error("Missing risk title")
	}
}

func TestRenderer_Themes(t *testing.T) {
	doc := createTestOKR()
	r := New()

	themes := []struct {
		name  string
		color string
	}{
		{"default", "#5a67d8"},
		{"corporate", "#1a365d"},
		{"minimal", "#2d3748"},
	}

	for _, theme := range themes {
		t.Run(theme.name, func(t *testing.T) {
			opts := &render.Options{Theme: theme.name}
			output, err := r.Render(doc, opts)
			if err != nil {
				t.Fatalf("Render() error = %v", err)
			}

			if !strings.Contains(string(output), theme.color) {
				t.Errorf("Theme %s should contain color %s", theme.name, theme.color)
			}
		})
	}
}

func TestRenderer_ScoreCalculation(t *testing.T) {
	doc := createTestOKR()
	doc.Objectives[0].KeyResults = []okr.KeyResult{
		{ID: "KR-1", Title: "First KR", Score: 0.8},
		{ID: "KR-2", Title: "Second KR", Score: 0.6},
	}

	r := New()
	output, err := r.Render(doc, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Should show 70% progress (average of 0.8 and 0.6)
	if !strings.Contains(content, "70%") {
		t.Error("Missing calculated progress percentage")
	}
}

func TestRenderer_MaxObjectives(t *testing.T) {
	doc := &okr.OKRDocument{
		Metadata: &okr.Metadata{
			Name:  "Large OKR",
			Owner: "Test Owner",
		},
		Objectives: []okr.Objective{
			{ID: "O1", Title: "Objective 1", KeyResults: []okr.KeyResult{{Title: "KR1"}}},
			{ID: "O2", Title: "Objective 2", KeyResults: []okr.KeyResult{{Title: "KR2"}}},
			{ID: "O3", Title: "Objective 3", KeyResults: []okr.KeyResult{{Title: "KR3"}}},
			{ID: "O4", Title: "Objective 4", KeyResults: []okr.KeyResult{{Title: "KR4"}}},
			{ID: "O5", Title: "Objective 5", KeyResults: []okr.KeyResult{{Title: "KR5"}}},
		},
	}

	r := New()
	opts := &render.Options{
		MaxObjectives: 3,
	}

	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Should have O1, O2, O3 but not O4, O5 in detailed slides
	// (Overview still shows all, but detailed slides are limited)
	if strings.Count(content, "<!-- _class: objective -->") > 3 {
		t.Error("Should have at most 3 objective detail slides")
	}
}

func createTestOKR() *okr.OKRDocument {
	return &okr.OKRDocument{
		Metadata: &okr.Metadata{
			ID:     "OKR-TEST-001",
			Name:   "Test OKR",
			Owner:  "Test Owner",
			Team:   "Engineering",
			Period: "2025-Q1",
			Status: okr.StatusActive,
		},
		Theme: "Deliver excellent product quality",
		Objectives: []okr.Objective{
			{
				ID:          "O1",
				Title:       "Improve Product Quality",
				Description: "Make our product more reliable and user-friendly",
				Owner:       "Engineering Lead",
				KeyResults: []okr.KeyResult{
					{
						ID:         "KR1",
						Title:      "Increase test coverage to 80%",
						Target:     "80%",
						Score:      0.7,
						Confidence: okr.ConfidenceHigh,
						Status:     "On Track",
					},
					{
						ID:         "KR2",
						Title:      "Reduce bug count by 50%",
						Target:     "50% reduction",
						Score:      0.5,
						Confidence: okr.ConfidenceMedium,
						Status:     "At Risk",
					},
				},
			},
		},
	}
}
