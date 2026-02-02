package prd

import (
	"encoding/json"
	"testing"
)

func TestCompletenessEmptyDocument(t *testing.T) {
	doc := Document{}
	report := doc.CheckCompleteness()

	// Empty document should have a low score
	if report.OverallScore >= 50 {
		t.Errorf("Empty document score too high: got %.1f%%, want < 50%%", report.OverallScore)
	}

	if report.Grade != "F" {
		t.Errorf("Empty document grade: got %s, want F", report.Grade)
	}

	if report.RequiredComplete != 0 {
		t.Errorf("RequiredComplete: got %d, want 0", report.RequiredComplete)
	}
}

func TestCompletenessMinimalDocument(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:      "prd-001",
			Title:   "Test PRD",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Test Author"}},
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "This is a detailed problem statement that explains the issues we're trying to solve in sufficient detail.",
			ProposedSolution: "This is a detailed proposed solution that explains how we'll address the problem with enough context.",
			ExpectedOutcomes: []string{"Outcome 1", "Outcome 2", "Outcome 3"},
		},
		Objectives: Objectives{
			OKRs: []OKR{
				{
					Objective: Objective{ID: "o-1", Description: "Business objective 1"},
					KeyResults: []KeyResult{
						{ID: "kr-1", Description: "Metric 1", Target: "100"},
						{ID: "kr-2", Description: "Metric 2", Target: "200"},
					},
				},
				{
					Objective: Objective{ID: "o-2", Description: "Product goal 1"},
					KeyResults: []KeyResult{
						{ID: "kr-3", Description: "Metric 3", Target: "300"},
					},
				},
			},
		},
		Personas: []Persona{
			{
				ID:          "p-1",
				Name:        "Test Persona",
				Role:        "Developer",
				Description: "A test persona",
				Goals:       []string{"Goal 1"},
				PainPoints:  []string{"Pain point 1"},
			},
		},
		UserStories: []UserStory{
			{
				ID:        "us-1",
				PersonaID: "p-1",
				Title:     "Test Story",
				AsA:       "developer",
				IWant:     "to test",
				SoThat:    "I can verify",
				Priority:  PriorityHigh,
				PhaseID:   "phase-1",
				AcceptanceCriteria: []AcceptanceCriterion{
					{ID: "ac-1", Description: "Criterion 1"},
				},
			},
		},
		Requirements: Requirements{
			Functional: []FunctionalRequirement{
				{
					ID:          "fr-1",
					Title:       "Functional requirement",
					Description: "Test requirement",
					Priority:    MoSCoWMust,
				},
			},
			NonFunctional: []NonFunctionalRequirement{
				{
					ID:          "nfr-1",
					Category:    NFRPerformance,
					Title:       "Performance requirement",
					Description: "Response time < 200ms",
					Target:      "P95 < 200ms",
					Priority:    MoSCoWMust,
				},
			},
		},
		Roadmap: Roadmap{
			Phases: []Phase{
				{
					ID:              "phase-1",
					Name:            "MVP",
					Type:            PhaseTypeMilestone,
					Goals:           []string{"Launch MVP"},
					Deliverables:    []Deliverable{{ID: "d-1", Title: "Feature 1", Type: DeliverableFeature}},
					SuccessCriteria: []string{"100 users"},
				},
			},
		},
	}

	report := doc.CheckCompleteness()

	// Should have a reasonable score with minimal required fields
	// With just 1 persona, limited stories/requirements, and no optional sections,
	// the score will be in the 50-60% range
	if report.OverallScore < 40 {
		t.Errorf("Minimal document score too low: got %.1f%%, want >= 40%%", report.OverallScore)
	}

	if report.OverallScore > 70 {
		t.Errorf("Minimal document score too high: got %.1f%%, want <= 70%%", report.OverallScore)
	}

	// Grade should be D or F for minimal document (both are acceptable)
	if report.Grade != "D" && report.Grade != "F" {
		t.Errorf("Minimal document grade unexpected: got %s, want D or F", report.Grade)
	}

	// Check that required sections are evaluated
	if len(report.Sections) != 15 {
		t.Errorf("Expected 15 sections (7 required + 8 optional), got %d", len(report.Sections))
	}
}

func TestCompletenessScoreToGrade(t *testing.T) {
	tests := []struct {
		score float64
		grade string
	}{
		{95, "A"},
		{90, "A"},
		{89, "B"},
		{80, "B"},
		{79, "C"},
		{70, "C"},
		{69, "D"},
		{60, "D"},
		{59, "F"},
		{0, "F"},
	}

	for _, tt := range tests {
		got := scoreToGrade(tt.score)
		if got != tt.grade {
			t.Errorf("scoreToGrade(%.1f) = %s, want %s", tt.score, got, tt.grade)
		}
	}
}

func TestCompletenessMetadataSection(t *testing.T) {
	tests := []struct {
		name     string
		metadata Metadata
		minScore float64
		maxScore float64
	}{
		{
			name:     "empty metadata",
			metadata: Metadata{},
			minScore: 0,
			maxScore: 20,
		},
		{
			name: "partial metadata",
			metadata: Metadata{
				ID:    "prd-001",
				Title: "Test",
			},
			minScore: 30,
			maxScore: 60,
		},
		{
			name: "complete metadata",
			metadata: Metadata{
				ID:      "prd-001",
				Title:   "Test PRD",
				Version: "1.0.0",
				Status:  StatusDraft,
				Authors: []Person{{Name: "Author"}},
			},
			minScore: 100,
			maxScore: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc := Document{Metadata: tt.metadata}
			section := doc.checkMetadata()

			if section.Score < tt.minScore || section.Score > tt.maxScore {
				t.Errorf("Metadata score: got %.1f%%, want between %.1f%% and %.1f%%",
					section.Score, tt.minScore, tt.maxScore)
			}

			if section.Required != true {
				t.Error("Metadata should be marked as required")
			}
		})
	}
}

func TestCompletenessPersonaQuality(t *testing.T) {
	tests := []struct {
		name       string
		persona    Persona
		wantIssues int
	}{
		{
			name:       "empty persona",
			persona:    Persona{},
			wantIssues: 3, // missing goals, pain points, description
		},
		{
			name: "complete persona",
			persona: Persona{
				ID:          "p-1",
				Name:        "Test",
				Role:        "Dev",
				Description: "A developer",
				Goals:       []string{"Build things"},
				PainPoints:  []string{"Slow builds"},
			},
			wantIssues: 0,
		},
		{
			name: "missing goals",
			persona: Persona{
				ID:          "p-1",
				Name:        "Test",
				Description: "A developer",
				PainPoints:  []string{"Slow builds"},
			},
			wantIssues: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := checkPersonaQuality(tt.persona)
			if len(issues) != tt.wantIssues {
				t.Errorf("checkPersonaQuality: got %d issues, want %d", len(issues), tt.wantIssues)
			}
		})
	}
}

func TestCompletenessIsPersonaComplete(t *testing.T) {
	complete := Persona{
		ID:          "p-1",
		Name:        "Test",
		Role:        "Dev",
		Description: "A developer",
		Goals:       []string{"Build"},
		PainPoints:  []string{"Bugs"},
	}

	incomplete := Persona{
		ID:   "p-1",
		Name: "Test",
	}

	if !isPersonaComplete(complete) {
		t.Error("Complete persona should return true")
	}

	if isPersonaComplete(incomplete) {
		t.Error("Incomplete persona should return false")
	}
}

func TestCompletenessRequirementsNFRCategories(t *testing.T) {
	doc := Document{
		Requirements: Requirements{
			NonFunctional: []NonFunctionalRequirement{
				{ID: "nfr-1", Category: NFRPerformance},
				{ID: "nfr-2", Category: NFRSecurity},
				{ID: "nfr-3", Category: NFRReliability},
			},
		},
	}

	section := doc.checkRequirements()

	// Should have higher score with all essential categories
	if section.Score < 50 {
		t.Errorf("Requirements with essential NFR categories should score higher: got %.1f%%", section.Score)
	}

	// Check that missing categories are reported
	docMissing := Document{
		Requirements: Requirements{
			NonFunctional: []NonFunctionalRequirement{
				{ID: "nfr-1", Category: NFRUsability},
			},
		},
	}

	sectionMissing := docMissing.checkRequirements()
	hasNFRIssue := false
	for _, issue := range sectionMissing.Issues {
		if len(issue) > 0 {
			hasNFRIssue = true
			break
		}
	}
	if !hasNFRIssue {
		t.Error("Should report missing essential NFR categories")
	}
}

func TestCompletenessOptionalSections(t *testing.T) {
	// Document without optional sections
	docWithout := Document{}
	reportWithout := docWithout.CheckCompleteness()

	// Document with optional sections
	docWith := Document{
		Assumptions: &AssumptionsConstraints{
			Assumptions: []Assumption{
				{ID: "a-1", Description: "Assumption 1"},
				{ID: "a-2", Description: "Assumption 2"},
				{ID: "a-3", Description: "Assumption 3"},
			},
		},
		OutOfScope: []string{"Item 1", "Item 2", "Item 3"},
		Risks: []Risk{
			{ID: "r-1", Description: "Risk 1", Mitigation: "Mitigation 1"},
			{ID: "r-2", Description: "Risk 2", Mitigation: "Mitigation 2"},
		},
		Glossary: []GlossaryTerm{
			{Term: "Term 1", Definition: "Definition 1"},
			{Term: "Term 2", Definition: "Definition 2"},
			{Term: "Term 3", Definition: "Definition 3"},
			{Term: "Term 4", Definition: "Definition 4"},
			{Term: "Term 5", Definition: "Definition 5"},
		},
	}
	reportWith := docWith.CheckCompleteness()

	if reportWith.OptionalComplete <= reportWithout.OptionalComplete {
		t.Error("Document with optional sections should have more complete optional sections")
	}

	if reportWith.OverallScore <= reportWithout.OverallScore {
		t.Error("Document with optional sections should have higher overall score")
	}
}

func TestCompletenessReportJSON(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:      "prd-001",
			Title:   "Test",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Author"}},
		},
	}

	report := doc.CheckCompleteness()

	// Should marshal to JSON without error
	data, err := json.Marshal(report)
	if err != nil {
		t.Errorf("Failed to marshal report to JSON: %v", err)
	}

	// Should unmarshal back
	var unmarshaled CompletenessReport
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal report from JSON: %v", err)
	}

	if unmarshaled.Grade != report.Grade {
		t.Errorf("Grade mismatch after JSON round-trip: got %s, want %s",
			unmarshaled.Grade, report.Grade)
	}
}

func TestCompletenessFormatReport(t *testing.T) {
	doc := Document{
		Metadata: Metadata{
			ID:      "prd-001",
			Title:   "Test PRD",
			Version: "1.0.0",
			Status:  StatusDraft,
			Authors: []Person{{Name: "Author"}},
		},
	}

	report := doc.CheckCompleteness()
	formatted := report.FormatReport()

	// Check that the report contains expected sections
	expectedStrings := []string{
		"PRD COMPLETENESS REPORT",
		"Overall Score:",
		"Grade:",
		"SECTION BREAKDOWN",
		"Required Sections:",
		"Optional Sections:",
		"Metadata",
	}

	for _, expected := range expectedStrings {
		if !containsString(formatted, expected) {
			t.Errorf("Formatted report missing expected string: %s", expected)
		}
	}
}

func TestCompletenessGetStatus(t *testing.T) {
	tests := []struct {
		score  float64
		status string
	}{
		{100, "complete"},
		{80, "complete"},
		{79, "partial"},
		{40, "partial"},
		{39, "missing"},
		{0, "missing"},
	}

	for _, tt := range tests {
		got := getStatus(tt.score)
		if got != tt.status {
			t.Errorf("getStatus(%.1f) = %s, want %s", tt.score, got, tt.status)
		}
	}
}

func TestCompletenessGetStatusIcon(t *testing.T) {
	tests := []struct {
		status string
		icon   string
	}{
		{"complete", "[+]"},
		{"partial", "[~]"},
		{"missing", "[ ]"},
		{"unknown", "[ ]"},
	}

	for _, tt := range tests {
		got := getStatusIcon(tt.status)
		if got != tt.icon {
			t.Errorf("getStatusIcon(%s) = %s, want %s", tt.status, got, tt.icon)
		}
	}
}

func TestCompletenessFilterByPriority(t *testing.T) {
	recs := []Recommendation{
		{Section: "A", Priority: RecommendCritical, Message: "Critical 1"},
		{Section: "B", Priority: RecommendHigh, Message: "High 1"},
		{Section: "C", Priority: RecommendCritical, Message: "Critical 2"},
		{Section: "D", Priority: RecommendMedium, Message: "Medium 1"},
	}

	critical := filterByPriority(recs, RecommendCritical)
	if len(critical) != 2 {
		t.Errorf("filterByPriority(Critical): got %d, want 2", len(critical))
	}

	high := filterByPriority(recs, RecommendHigh)
	if len(high) != 1 {
		t.Errorf("filterByPriority(High): got %d, want 1", len(high))
	}

	low := filterByPriority(recs, RecommendLow)
	if len(low) != 0 {
		t.Errorf("filterByPriority(Low): got %d, want 0", len(low))
	}
}

func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
