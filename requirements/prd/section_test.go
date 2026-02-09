package prd

import (
	"strings"
	"testing"
)

func TestGetSectionOrder(t *testing.T) {
	tests := []struct {
		name     string
		prdType  PRDType
		expected []SectionID
	}{
		{
			name:     "default type",
			prdType:  PRDTypeDefault,
			expected: DefaultSectionOrder,
		},
		{
			name:     "strategy type",
			prdType:  PRDTypeStrategy,
			expected: StrategySectionOrder,
		},
		{
			name:     "feature type",
			prdType:  PRDTypeFeature,
			expected: FeatureSectionOrder,
		},
		{
			name:     "technical type",
			prdType:  PRDTypeTechnical,
			expected: TechnicalSectionOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetSectionOrder(tt.prdType)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d sections, got %d", len(tt.expected), len(result))
			}
			for i, id := range result {
				if id != tt.expected[i] {
					t.Errorf("section %d: expected %s, got %s", i, tt.expected[i], id)
				}
			}
		})
	}
}

func TestStrategySectionOrderHasContextFirst(t *testing.T) {
	// Strategy PRDs should have context sections early
	order := StrategySectionOrder

	// Find positions of key sections
	positions := make(map[SectionID]int)
	for i, id := range order {
		positions[id] = i
	}

	// CurrentState should come before Objectives
	if positions[SectionCurrentState] > positions[SectionObjectives] {
		t.Error("CurrentState should come before Objectives in strategy order")
	}

	// Problem should come before Objectives
	if positions[SectionProblem] > positions[SectionObjectives] {
		t.Error("Problem should come before Objectives in strategy order")
	}

	// Market should come before UserStories
	if positions[SectionMarket] > positions[SectionUserStories] {
		t.Error("Market should come before UserStories in strategy order")
	}

	// Solution should come before FunctionalReqs
	if positions[SectionSolution] > positions[SectionFunctionalReqs] {
		t.Error("Solution should come before FunctionalReqs in strategy order")
	}
}

func TestFeatureSectionOrderHasUserNeedsFirst(t *testing.T) {
	// Feature PRDs should have user-focused sections early
	order := FeatureSectionOrder

	// Find positions of key sections
	positions := make(map[SectionID]int)
	for i, id := range order {
		positions[id] = i
	}

	// Problem should be early (within first 5)
	if positions[SectionProblem] > 5 {
		t.Error("Problem should be in first 5 sections for feature order")
	}

	// Personas should come before TechArchitecture
	if positions[SectionPersonas] > positions[SectionTechArchitecture] {
		t.Error("Personas should come before TechArchitecture in feature order")
	}

	// UserStories should come before Solution
	if positions[SectionUserStories] > positions[SectionSolution] {
		t.Error("UserStories should come before Solution in feature order")
	}
}

func TestTechnicalSectionOrderHasArchitectureEarly(t *testing.T) {
	// Technical PRDs should have architecture sections early
	order := TechnicalSectionOrder

	// Find positions of key sections
	positions := make(map[SectionID]int)
	for i, id := range order {
		positions[id] = i
	}

	// TechArchitecture should be in first 10
	if positions[SectionTechArchitecture] > 10 {
		t.Error("TechArchitecture should be in first 10 sections for technical order")
	}

	// SecurityModel should be in first 10
	if positions[SectionSecurityModel] > 10 {
		t.Error("SecurityModel should be in first 10 sections for technical order")
	}

	// TechArchitecture should come before Personas
	if positions[SectionTechArchitecture] > positions[SectionPersonas] {
		t.Error("TechArchitecture should come before Personas in technical order")
	}
}

func TestCompleteSectionOrder(t *testing.T) {
	partial := []SectionID{
		SectionExecutiveSummary,
		SectionProblem,
		SectionSolution,
	}

	template := []SectionID{
		SectionExecutiveSummary,
		SectionObjectives,
		SectionPersonas,
		SectionProblem,
		SectionSolution,
		SectionRoadmap,
	}

	result := CompleteSectionOrder(partial, template)

	// Should start with partial order
	if result[0] != SectionExecutiveSummary {
		t.Error("first section should be executiveSummary")
	}
	if result[1] != SectionProblem {
		t.Error("second section should be problem")
	}
	if result[2] != SectionSolution {
		t.Error("third section should be solution")
	}

	// Should have all sections from template
	if len(result) != 6 {
		t.Errorf("expected 6 sections, got %d", len(result))
	}

	// Should include missing sections
	found := make(map[SectionID]bool)
	for _, id := range result {
		found[id] = true
	}
	if !found[SectionObjectives] {
		t.Error("should include objectives from template")
	}
	if !found[SectionPersonas] {
		t.Error("should include personas from template")
	}
	if !found[SectionRoadmap] {
		t.Error("should include roadmap from template")
	}
}

func TestValidateSectionOrder(t *testing.T) {
	tests := []struct {
		name    string
		order   []string
		invalid []string
	}{
		{
			name:    "all valid",
			order:   []string{"executiveSummary", "problem", "solution"},
			invalid: nil,
		},
		{
			name:    "one invalid",
			order:   []string{"executiveSummary", "invalidSection", "solution"},
			invalid: []string{"invalidSection"},
		},
		{
			name:    "multiple invalid",
			order:   []string{"foo", "executiveSummary", "bar"},
			invalid: []string{"foo", "bar"},
		},
		{
			name:    "empty order",
			order:   []string{},
			invalid: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateSectionOrder(tt.order)
			if len(result) != len(tt.invalid) {
				t.Errorf("expected %d invalid, got %d", len(tt.invalid), len(result))
			}
			for i, id := range result {
				if i < len(tt.invalid) && id != tt.invalid[i] {
					t.Errorf("invalid[%d]: expected %s, got %s", i, tt.invalid[i], id)
				}
			}
		})
	}
}

func TestAllSectionIDs(t *testing.T) {
	ids := AllSectionIDs()

	// Check that all expected sections are present
	expectedSections := []SectionID{
		SectionExecutiveSummary,
		SectionCurrentState,
		SectionProblem,
		SectionMarket,
		SectionSolution,
		SectionSuccessMetrics,
		SectionObjectives,
		SectionNonGoals,
		SectionPersonas,
		SectionUserStories,
		SectionFunctionalReqs,
		SectionNonFunctionalReqs,
		SectionRoadmap,
		SectionTechArchitecture,
		SectionSecurityModel,
		SectionDecisions,
		SectionRisks,
		SectionAssumptions,
		SectionInScope,
		SectionOutOfScope,
		SectionOpenItems,
		SectionRelatedDocuments,
		SectionAppendices,
		SectionGlossary,
		SectionReviews,
		SectionRevisionHistory,
		SectionCustom,
	}

	for _, id := range expectedSections {
		if !ids[id] {
			t.Errorf("expected section %s to be in AllSectionIDs", id)
		}
	}
}

func TestSectionDisplayNames(t *testing.T) {
	// All sections should have display names
	for id := range AllSectionIDs() {
		if id == SectionCustom {
			continue // Custom sections have dynamic titles
		}
		name := SectionDisplayNames[id]
		if name == "" {
			t.Errorf("section %s should have a display name", id)
		}
	}
}

func TestSectionAnchors(t *testing.T) {
	// All sections should have anchors
	for id := range AllSectionIDs() {
		if id == SectionCustom {
			continue // Custom sections have dynamic anchors
		}
		anchor := SectionAnchors[id]
		if anchor == "" {
			t.Errorf("section %s should have an anchor", id)
		}
		// Anchors should be lowercase with hyphens
		if strings.ToLower(anchor) != anchor {
			t.Errorf("anchor %s should be lowercase", anchor)
		}
		if strings.Contains(anchor, " ") {
			t.Errorf("anchor %s should not contain spaces", anchor)
		}
	}
}

func TestDocumentHasSection(t *testing.T) {
	doc := &Document{}

	// Always-present sections
	if !doc.HasSection(SectionExecutiveSummary) {
		t.Error("ExecutiveSummary should always be present")
	}
	if !doc.HasSection(SectionObjectives) {
		t.Error("Objectives should always be present")
	}
	if !doc.HasSection(SectionRoadmap) {
		t.Error("Roadmap should always be present")
	}

	// Optional sections should be absent on empty document
	if doc.HasSection(SectionProblem) {
		t.Error("Problem should not be present on empty document")
	}
	if doc.HasSection(SectionMarket) {
		t.Error("Market should not be present on empty document")
	}
	if doc.HasSection(SectionTechArchitecture) {
		t.Error("TechArchitecture should not be present on empty document")
	}

	// Add optional content and check again
	doc.Problem = &ProblemDefinition{Statement: "test"}
	if !doc.HasSection(SectionProblem) {
		t.Error("Problem should be present after setting")
	}

	doc.Personas = []Persona{{ID: "p1", Name: "Test"}}
	if !doc.HasSection(SectionPersonas) {
		t.Error("Personas should be present after adding")
	}
}

func TestDocumentGetSectionOrder(t *testing.T) {
	t.Run("default order", func(t *testing.T) {
		doc := &Document{}
		order := doc.GetSectionOrder()
		if len(order) != len(DefaultSectionOrder) {
			t.Errorf("expected %d sections, got %d", len(DefaultSectionOrder), len(order))
		}
	})

	t.Run("strategy type", func(t *testing.T) {
		doc := &Document{}
		doc.Metadata.PRDType = PRDTypeStrategy
		order := doc.GetSectionOrder()
		if len(order) != len(StrategySectionOrder) {
			t.Errorf("expected %d sections, got %d", len(StrategySectionOrder), len(order))
		}
		// First section should be ExecutiveSummary
		if order[0] != SectionExecutiveSummary {
			t.Errorf("first section should be executiveSummary, got %s", order[0])
		}
		// Second should be CurrentState for strategy
		if order[1] != SectionCurrentState {
			t.Errorf("second section should be currentState for strategy, got %s", order[1])
		}
	})

	t.Run("custom order", func(t *testing.T) {
		doc := &Document{}
		doc.Metadata.SectionOrder = []string{
			"executiveSummary",
			"problem",
			"solution",
			"objectives",
		}
		order := doc.GetSectionOrder()

		// Custom sections should come first
		if order[0] != SectionExecutiveSummary {
			t.Errorf("first section should be executiveSummary, got %s", order[0])
		}
		if order[1] != SectionProblem {
			t.Errorf("second section should be problem, got %s", order[1])
		}
		if order[2] != SectionSolution {
			t.Errorf("third section should be solution, got %s", order[2])
		}
		if order[3] != SectionObjectives {
			t.Errorf("fourth section should be objectives, got %s", order[3])
		}

		// Should include remaining sections from default order
		if len(order) != len(DefaultSectionOrder) {
			t.Errorf("should include all sections, got %d", len(order))
		}
	})
}

func TestDocumentGetActiveSections(t *testing.T) {
	doc := &Document{}

	// Empty document should only have always-present sections
	active := doc.GetActiveSections()
	expectedAlways := []SectionID{
		SectionExecutiveSummary,
		SectionObjectives,
		SectionRoadmap,
	}

	// Check that always-present sections are included
	found := make(map[SectionID]bool)
	for _, id := range active {
		found[id] = true
	}
	for _, id := range expectedAlways {
		if !found[id] {
			t.Errorf("expected always-present section %s", id)
		}
	}

	// Add some optional content
	doc.Problem = &ProblemDefinition{Statement: "test"}
	doc.Personas = []Persona{{ID: "p1", Name: "Test"}}
	doc.UserStories = []UserStory{{ID: "us1"}}

	active = doc.GetActiveSections()
	found = make(map[SectionID]bool)
	for _, id := range active {
		found[id] = true
	}

	if !found[SectionProblem] {
		t.Error("Problem should be active after setting")
	}
	if !found[SectionPersonas] {
		t.Error("Personas should be active after adding")
	}
	if !found[SectionUserStories] {
		t.Error("UserStories should be active after adding")
	}
}

func TestToMarkdownWithDifferentPRDTypes(t *testing.T) {
	// Create a document with some optional sections
	doc := &Document{
		Metadata: Metadata{
			ID:    "TEST-001",
			Title: "Test PRD",
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem",
			ProposedSolution: "Test solution",
		},
		Problem: &ProblemDefinition{
			Statement: "Detailed problem statement",
		},
		CurrentState: &CurrentState{
			Overview: "Current state overview",
		},
	}

	t.Run("default order has problem late", func(t *testing.T) {
		doc.Metadata.PRDType = PRDTypeDefault
		md := doc.ToMarkdown(MarkdownOptions{IncludeFrontmatter: false})

		// In default order, Executive Summary comes before Problem Definition
		execPos := strings.Index(md, "## Executive Summary")
		probPos := strings.Index(md, "## Problem Definition")
		currPos := strings.Index(md, "## Current State")

		if execPos == -1 || probPos == -1 || currPos == -1 {
			t.Error("expected all sections to be present")
		}

		// In default order, CurrentState comes before Problem
		if currPos > probPos {
			t.Error("in default order, CurrentState should come before Problem")
		}
	})

	t.Run("strategy order has context early", func(t *testing.T) {
		doc.Metadata.PRDType = PRDTypeStrategy
		md := doc.ToMarkdown(MarkdownOptions{IncludeFrontmatter: false})

		// In strategy order, CurrentState and Problem should come right after Executive Summary
		execPos := strings.Index(md, "## Executive Summary")
		currPos := strings.Index(md, "## Current State")
		probPos := strings.Index(md, "## Problem Definition")
		objPos := strings.Index(md, "## Objectives and Goals")

		if execPos == -1 || currPos == -1 || probPos == -1 || objPos == -1 {
			t.Error("expected all sections to be present")
		}

		// Executive Summary should be first
		if execPos > currPos || execPos > probPos {
			t.Error("Executive Summary should come first")
		}

		// CurrentState and Problem should come before Objectives
		if currPos > objPos {
			t.Error("in strategy order, CurrentState should come before Objectives")
		}
		if probPos > objPos {
			t.Error("in strategy order, Problem should come before Objectives")
		}
	})
}

func TestListSections(t *testing.T) {
	sections := ListSections()

	if len(sections) != len(DefaultSectionOrder) {
		t.Errorf("expected %d sections, got %d", len(DefaultSectionOrder), len(sections))
	}

	// Check that all sections have IDs and display names
	for _, s := range sections {
		if s.ID == "" {
			t.Error("section should have ID")
		}
		if s.ID != SectionCustom && s.DisplayName == "" {
			t.Errorf("section %s should have display name", s.ID)
		}
	}
}
