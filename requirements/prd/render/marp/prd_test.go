package marp

import (
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/requirements/prd"
	"github.com/grokify/prism-roadmap/requirements/prd/render"
)

func TestPRDRenderer_Format(t *testing.T) {
	r := NewPRDRenderer()
	if got := r.Format(); got != "marp" {
		t.Errorf("Format() = %v, want %v", got, "marp")
	}
}

func TestPRDRenderer_FileExtension(t *testing.T) {
	r := NewPRDRenderer()
	if got := r.FileExtension(); got != ".md" {
		t.Errorf("FileExtension() = %v, want %v", got, ".md")
	}
}

func TestPRDRenderer_Render(t *testing.T) {
	doc := createTestPRD()
	r := NewPRDRenderer()

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
	if !strings.Contains(content, "# Test PRD") {
		t.Error("Missing title in slides")
	}

	// Check problem slide
	if !strings.Contains(content, "The Problem") {
		t.Error("Missing problem slide")
	}
	if !strings.Contains(content, "Users need better testing") {
		t.Error("Missing problem statement")
	}

	// Check solution slide
	if !strings.Contains(content, "The Solution") {
		t.Error("Missing solution slide")
	}
	if !strings.Contains(content, "Implement a testing framework") {
		t.Error("Missing proposed solution")
	}

	// Check objectives slide
	if !strings.Contains(content, "Objectives") {
		t.Error("Missing objectives slide")
	}

	// Check metrics slide
	if !strings.Contains(content, "Key Results") {
		t.Error("Missing metrics slide")
	}

	// Check summary slide
	if !strings.Contains(content, "Summary") {
		t.Error("Missing summary slide")
	}
}

func TestPRDRenderer_RenderWithOptions(t *testing.T) {
	doc := createTestPRD()
	r := NewPRDRenderer()

	opts := &render.Options{
		Theme:               "corporate",
		IncludeGoals:        false,
		IncludeRoadmap:      false,
		IncludeRisks:        false,
		IncludeRequirements: false,
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

	// Check that requirements are not included (option disabled)
	if strings.Contains(content, "Key Requirements") {
		t.Error("Requirements slide should be excluded when IncludeRequirements=false")
	}
}

func TestPRDRenderer_RenderWithPersonas(t *testing.T) {
	doc := createTestPRD()
	doc.Personas = []prd.Persona{
		{
			ID:        "PER-1",
			Name:      "Test Developer",
			Role:      "Software Engineer",
			IsPrimary: true,
			Goals:     []string{"Write reliable code", "Fast feedback"},
			PainPoints: []string{
				"Tests take too long",
				"Hard to debug failures",
			},
		},
	}

	r := NewPRDRenderer()
	output, err := r.Render(doc, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check personas slide
	if !strings.Contains(content, "Target Personas") {
		t.Error("Missing personas slide")
	}
	if !strings.Contains(content, "Test Developer") {
		t.Error("Missing persona name")
	}
	if !strings.Contains(content, "Software Engineer") {
		t.Error("Missing persona role")
	}
}

func TestPRDRenderer_RenderWithRisks(t *testing.T) {
	doc := createTestPRD()
	doc.Risks = []prd.Risk{
		{
			ID:          "R-1",
			Description: "Technical complexity",
			Impact:      "high",
			Mitigation:  "Incremental delivery",
		},
	}

	r := NewPRDRenderer()
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
	if !strings.Contains(content, "Technical complexity") {
		t.Error("Missing risk description")
	}
}

func TestPRDRenderer_Themes(t *testing.T) {
	doc := createTestPRD()
	r := NewPRDRenderer()

	themes := []struct {
		name  string
		color string
	}{
		{"default", "#5a67d8"}, // Uses structureddocs theme
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

func createTestPRD() *prd.Document {
	return &prd.Document{
		Metadata: prd.Metadata{
			ID:      "PRD-TEST-001",
			Title:   "Test PRD",
			Version: "1.0.0",
			Status:  prd.StatusDraft,
			Authors: []prd.Person{
				{Name: "Test Author", Role: "PM"},
			},
		},
		ExecutiveSummary: prd.ExecutiveSummary{
			ProblemStatement: "Users need better testing",
			ProposedSolution: "Implement a testing framework",
			ExpectedOutcomes: []string{"Faster feedback", "Better quality"},
			TargetAudience:   "Software developers",
			ValueProposition: "Save time with automated testing",
		},
		Objectives: prd.Objectives{
			OKRs: []prd.OKR{
				{
					Objective: prd.Objective{ID: "O-1", Description: "Increase developer productivity"},
					KeyResults: []prd.KeyResult{
						{ID: "KR-1", Description: "Test Coverage", Target: ">80%", Baseline: "50%"},
					},
				},
				{
					Objective: prd.Objective{ID: "O-2", Description: "Reduce test execution time"},
					KeyResults: []prd.KeyResult{
						{ID: "KR-2", Description: "Execution time", Target: "<5min"},
					},
				},
			},
		},
		Requirements: prd.Requirements{
			Functional: []prd.FunctionalRequirement{
				{
					ID:       "FR-1",
					Title:    "Test Runner",
					Priority: prd.MoSCoWMust,
				},
			},
			NonFunctional: []prd.NonFunctionalRequirement{
				{
					ID:          "NFR-1",
					Category:    prd.NFRPerformance,
					Description: "Tests complete in <5 seconds",
				},
			},
		},
		Roadmap: prd.Roadmap{
			Phases: []prd.Phase{
				{
					ID:    "P1",
					Name:  "MVP",
					Goals: []string{"Basic test runner"},
				},
			},
		},
	}
}
