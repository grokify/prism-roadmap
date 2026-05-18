package marp

import (
	"strings"
	"testing"

	"github.com/grokify/prism-execution/goals/okr"
	"github.com/grokify/prism-execution/goals/v2mom"
	"github.com/grokify/prism-execution/requirements/prd"
	"github.com/grokify/prism-execution/requirements/prd/render"
)

func TestPRDGoalsRenderer_Format(t *testing.T) {
	r := NewPRDGoalsRenderer()
	if got := r.Format(); got != "marp-prd-goals" {
		t.Errorf("Format() = %v, want %v", got, "marp-prd-goals")
	}
}

func TestPRDGoalsRenderer_FileExtension(t *testing.T) {
	r := NewPRDGoalsRenderer()
	if got := r.FileExtension(); got != ".md" {
		t.Errorf("FileExtension() = %v, want %v", got, ".md")
	}
}

func TestPRDGoalsRenderer_RenderBasic(t *testing.T) {
	doc := createTestPRDWithoutGoals()
	r := NewPRDGoalsRenderer()

	output, err := r.Render(doc, nil)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check Marp front matter
	if !strings.Contains(content, "marp: true") {
		t.Error("Missing marp: true in front matter")
	}

	// Check title slide
	if !strings.Contains(content, "# Test PRD with Goals") {
		t.Error("Missing title in slides")
	}

	// Check agenda slide
	if !strings.Contains(content, "Agenda") {
		t.Error("Missing agenda slide")
	}

	// Check problem slide
	if !strings.Contains(content, "The Problem") {
		t.Error("Missing problem slide")
	}

	// Check solution slide
	if !strings.Contains(content, "The Solution") {
		t.Error("Missing solution slide")
	}
}

func TestPRDGoalsRenderer_RenderWithV2MOM(t *testing.T) {
	doc := createTestPRDWithV2MOM()
	r := NewPRDGoalsRenderer()

	opts := render.DefaultOptions()
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check title indicates goals alignment
	if !strings.Contains(content, "_with Goals Alignment_") {
		t.Error("Missing 'with Goals Alignment' subtitle")
	}

	// Check V2MOM Vision slide
	if !strings.Contains(content, "V2MOM: Vision & Values") {
		t.Error("Missing V2MOM Vision slide")
	}
	if !strings.Contains(content, "Be the best testing platform") {
		t.Error("Missing V2MOM vision statement")
	}

	// Check V2MOM Methods slide
	if !strings.Contains(content, "V2MOM: Methods") {
		t.Error("Missing V2MOM Methods slide")
	}
	if !strings.Contains(content, "Build automated testing") {
		t.Error("Missing V2MOM method")
	}

	// Check agenda includes V2MOM
	if !strings.Contains(content, "V2MOM Alignment") {
		t.Error("Agenda should mention V2MOM Alignment")
	}
}

func TestPRDGoalsRenderer_RenderWithOKR(t *testing.T) {
	doc := createTestPRDWithOKR()
	r := NewPRDGoalsRenderer()

	opts := render.DefaultOptions()
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check OKR Overview slide
	if !strings.Contains(content, "OKR Alignment") {
		t.Error("Missing OKR Alignment slide")
	}

	// Check OKR Key Results slide
	if !strings.Contains(content, "OKR: Key Results") {
		t.Error("Missing OKR Key Results slide")
	}

	// Check objectives are shown
	if !strings.Contains(content, "Improve testing capabilities") {
		t.Error("Missing OKR objective")
	}

	// Check key results are shown
	if !strings.Contains(content, "80% coverage") {
		t.Error("Missing OKR key result target")
	}

	// Check agenda includes OKR
	if !strings.Contains(content, "OKR Alignment") {
		t.Error("Agenda should mention OKR Alignment")
	}
}

func TestPRDGoalsRenderer_RenderWithBothGoals(t *testing.T) {
	doc := createTestPRDWithBothGoals()
	r := NewPRDGoalsRenderer()

	opts := render.DefaultOptions()
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check both V2MOM and OKR are present
	if !strings.Contains(content, "V2MOM: Vision & Values") {
		t.Error("Missing V2MOM slide when both goals present")
	}
	if !strings.Contains(content, "OKR Alignment") {
		t.Error("Missing OKR slide when both goals present")
	}

	// Check summary mentions both
	if !strings.Contains(content, "V2MOM") && !strings.Contains(content, "OKR") {
		t.Error("Summary should mention goals alignment")
	}
}

func TestPRDGoalsRenderer_RenderWithAlignedObjectives(t *testing.T) {
	doc := createTestPRDWithAlignedObjectives()
	r := NewPRDGoalsRenderer()

	opts := render.DefaultOptions()
	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// Check alignment summary slide
	if !strings.Contains(content, "Goals Alignment Summary") {
		t.Error("Missing Goals Alignment Summary slide")
	}
	if !strings.Contains(content, "BO-1") {
		t.Error("Missing aligned PRD objective reference")
	}
}

func TestPRDGoalsRenderer_GoalsDisabled(t *testing.T) {
	doc := createTestPRDWithV2MOM()
	r := NewPRDGoalsRenderer()

	opts := &render.Options{
		IncludeGoals: false,
	}

	output, err := r.Render(doc, opts)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	content := string(output)

	// V2MOM should not appear when goals disabled
	if strings.Contains(content, "V2MOM: Vision") {
		t.Error("V2MOM should not appear when IncludeGoals=false")
	}
}

func createTestPRDWithoutGoals() *prd.Document {
	return &prd.Document{
		Metadata: prd.Metadata{
			ID:      "PRD-GOALS-001",
			Title:   "Test PRD with Goals",
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
		},
		Objectives: prd.Objectives{
			OKRs: []prd.OKR{
				{
					Objective: prd.Objective{ID: "O-1", Description: "Increase developer productivity"},
					KeyResults: []prd.KeyResult{
						{ID: "KR-1", Description: "Test Coverage", Target: ">80%"},
					},
				},
			},
		},
	}
}

func createTestPRDWithV2MOM() *prd.Document {
	doc := createTestPRDWithoutGoals()
	doc.Goals = &prd.GoalsAlignment{
		V2MOM: &v2mom.V2MOM{
			Vision: "Be the best testing platform",
			Values: []v2mom.Value{
				{Name: "Quality", Description: "Deliver high-quality software"},
				{Name: "Speed", Description: "Fast feedback loops"},
			},
			Methods: []v2mom.Method{
				{Name: "Build automated testing", Priority: "P1"},
				{Name: "Integrate CI/CD", Priority: "P2"},
			},
			Obstacles: []v2mom.Obstacle{
				{Name: "Legacy systems", Severity: "High"},
			},
		},
	}
	return doc
}

func createTestPRDWithOKR() *prd.Document {
	doc := createTestPRDWithoutGoals()
	doc.Goals = &prd.GoalsAlignment{
		OKR: &okr.OKRDocument{
			Theme: "Deliver excellent testing capabilities",
			Objectives: []okr.Objective{
				{
					ID:    "O1",
					Title: "Improve testing capabilities",
					KeyResults: []okr.KeyResult{
						{
							ID:     "KR1",
							Title:  "Achieve test coverage",
							Target: "80% coverage",
							Score:  0.6,
						},
						{
							ID:     "KR2",
							Title:  "Reduce test runtime",
							Target: "50% reduction",
							Score:  0.4,
						},
					},
				},
			},
		},
	}
	return doc
}

func createTestPRDWithBothGoals() *prd.Document {
	doc := createTestPRDWithV2MOM()
	doc.Goals.OKR = &okr.OKRDocument{
		Theme: "Deliver excellent testing capabilities",
		Objectives: []okr.Objective{
			{
				ID:    "O1",
				Title: "Improve testing capabilities",
				KeyResults: []okr.KeyResult{
					{ID: "KR1", Title: "Achieve 80% coverage", Score: 0.6},
				},
			},
		},
	}
	return doc
}

func createTestPRDWithAlignedObjectives() *prd.Document {
	doc := createTestPRDWithOKR()
	doc.Goals.AlignedObjectives = map[string]string{
		"BO-1": "O1",
		"SM-1": "KR1",
	}
	return doc
}
