package prd

import (
	"strings"
	"testing"
	"time"
)

func TestGenerateProblem(t *testing.T) {
	doc := &Document{
		Problem: &ProblemDefinition{
			ID:         "PROB-001",
			Statement:  "Users struggle to find relevant documents",
			UserImpact: "Productivity loss of 2 hours per week",
			Confidence: 0.85,
			Evidence: []Evidence{
				{
					Type:       EvidenceSurvey,
					Source:     "Q4 User Survey",
					Summary:    "73% of users report difficulty",
					SampleSize: 500,
					Strength:   StrengthHigh,
					Date:       "2024-01",
				},
			},
			RootCauses: []string{
				"Poor search algorithm",
				"Lack of metadata",
			},
			AffectedSegments: []string{
				"Enterprise users",
				"Remote workers",
			},
			SecondaryProblems: []ProblemDefinition{
				{
					Statement:  "Document versioning is confusing",
					UserImpact: "Users work with outdated versions",
				},
			},
		},
	}

	result := doc.generateProblem()

	checks := []string{
		"## Problem Definition",
		"### Problem Statement",
		"Users struggle to find relevant documents",
		"### User Impact",
		"Productivity loss of 2 hours per week",
		"**Confidence:** 85%",
		"### Evidence",
		"| Type | Source |",
		"survey",
		"Q4 User Survey",
		"73% of users report difficulty",
		"500",
		"high",
		"### Root Causes",
		"Poor search algorithm",
		"### Affected Segments",
		"Enterprise users",
		"### Secondary Problems",
		"Document versioning is confusing",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateProblem() missing expected content: %q", check)
		}
	}
}

func TestGenerateMarket(t *testing.T) {
	doc := &Document{
		Market: &MarketDefinition{
			Alternatives: []Alternative{
				{
					ID:           "ALT-001",
					Name:         "Competitor A",
					Type:         AlternativeCompetitor,
					Description:  "Market leader",
					Strengths:    []string{"Brand recognition", "Large user base"},
					Weaknesses:   []string{"Expensive", "Complex"},
					WhyNotChosen: "Too expensive for target market",
				},
			},
			Differentiation: []string{
				"AI-powered search",
				"Better mobile experience",
			},
			MarketRisks: []string{
				"Market saturation",
				"Regulatory changes",
			},
		},
	}

	result := doc.generateMarket()

	checks := []string{
		"## Market Analysis",
		"### Alternatives",
		"| ID | Name | Type |",
		"ALT-001",
		"Competitor A",
		"competitor",
		"Too expensive for target market",
		"#### Competitor A",
		"**Strengths:**",
		"Brand recognition",
		"**Weaknesses:**",
		"Expensive",
		"### Differentiation",
		"AI-powered search",
		"### Market Risks",
		"Market saturation",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateMarket() missing expected content: %q", check)
		}
	}
}

func TestGenerateSolution(t *testing.T) {
	doc := &Document{
		Solution: &SolutionDefinition{
			SolutionOptions: []SolutionOption{
				{
					ID:              "SOL-001",
					Name:            "Build In-House",
					Description:     "Custom solution",
					EstimatedEffort: "6 months",
					Benefits:        []string{"Full control", "Custom features"},
					Tradeoffs:       []string{"Higher cost", "Longer timeline"},
					Risks:           []string{"Resource constraints"},
				},
				{
					ID:              "SOL-002",
					Name:            "Buy Off-the-Shelf",
					Description:     "Commercial solution",
					EstimatedEffort: "2 months",
				},
			},
			SelectedSolutionID: "SOL-001",
			SolutionRationale:  "Custom solution provides better long-term value",
			Confidence:         0.75,
		},
	}

	// Test without text icons
	opts := MarkdownOptions{UseTextIcons: false}
	result := doc.generateSolution(opts)

	checks := []string{
		"## Solution",
		"### Solution Options",
		"| ID | Name | Description | Effort | Selected |",
		"SOL-001",
		"Build In-House",
		"Custom solution",
		"6 months",
		"✅ Selected",
		"#### Build In-House",
		"**Benefits:**",
		"Full control",
		"**Tradeoffs:**",
		"Higher cost",
		"**Risks:**",
		"Resource constraints",
		"### Solution Rationale",
		"Custom solution provides better long-term value",
		"**Confidence:** 75%",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateSolution() missing expected content: %q", check)
		}
	}

	// Test with text icons
	opts.UseTextIcons = true
	result = doc.generateSolution(opts)

	if !strings.Contains(result, "[*] Selected") {
		t.Error("generateSolution() with UseTextIcons should use [*] Selected")
	}
}

func TestGenerateDecisions(t *testing.T) {
	date := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
	doc := &Document{
		Decisions: &DecisionsDefinition{
			Records: []DecisionRecord{
				{
					ID:                     "DEC-001",
					Decision:               "Use PostgreSQL for primary database",
					Rationale:              "Better JSON support and reliability",
					Status:                 DecisionAccepted,
					Date:                   date,
					MadeBy:                 "Architecture Team",
					AlternativesConsidered: []string{"MySQL", "MongoDB"},
				},
			},
		},
	}

	result := doc.generateDecisions()

	checks := []string{
		"## Decisions",
		"| ID | Decision | Rationale | Status | Date | Made By |",
		"DEC-001",
		"Use PostgreSQL for primary database",
		"Better JSON support and reliability",
		"accepted",
		"2024-01-15",
		"Architecture Team",
		"**DEC-001 - Alternatives Considered:**",
		"MySQL",
		"MongoDB",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateDecisions() missing expected content: %q", check)
		}
	}
}

func TestGenerateReviews(t *testing.T) {
	doc := &Document{
		Reviews: &ReviewsDefinition{
			ReviewBoardSummary: "PRD meets quality standards with minor issues",
			Decision:           ReviewApprove,
			QualityScores: &QualityScores{
				ProblemDefinition:    8.5,
				UserUnderstanding:    9.0,
				MarketAwareness:      7.5,
				SolutionFit:          8.0,
				ScopeDiscipline:      8.5,
				RequirementsQuality:  8.0,
				UXCoverage:           7.0,
				TechnicalFeasibility: 8.5,
				MetricsQuality:       7.5,
				RiskManagement:       8.0,
				OverallScore:         8.1,
			},
			Blockers: []Blocker{
				{
					ID:          "BLK-001",
					Category:    "Security",
					Description: "Missing threat model",
				},
			},
			RevisionTriggers: []RevisionTrigger{
				{
					IssueID:          "REV-001",
					Category:         "UX",
					Severity:         "minor",
					Description:      "Need more wireframes",
					RecommendedOwner: "Design Team",
				},
			},
		},
	}

	// Test without text icons
	opts := MarkdownOptions{UseTextIcons: false}
	result := doc.generateReviews(opts)

	checks := []string{
		"## Reviews",
		"### Review Board Summary",
		"PRD meets quality standards with minor issues",
		"**Decision:** ✅ Approved",
		"### Quality Scores",
		"| Dimension | Score |",
		"| Problem Definition | 8.5 |",
		"| **Overall Score** | **8.1** |",
		"### Blockers",
		"**BLK-001** (Security): Missing threat model",
		"### Revision Triggers",
		"| Issue ID | Category | Severity | Description | Recommended Owner |",
		"REV-001",
		"UX",
		"minor",
		"Need more wireframes",
		"Design Team",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateReviews() missing expected content: %q", check)
		}
	}

	// Test with text icons
	opts.UseTextIcons = true
	result = doc.generateReviews(opts)

	if !strings.Contains(result, "[APPROVED]") {
		t.Error("generateReviews() with UseTextIcons should use [APPROVED]")
	}
}

func TestGenerateRevisionHistory(t *testing.T) {
	date := time.Date(2024, 1, 10, 0, 0, 0, 0, time.UTC)
	doc := &Document{
		RevisionHistory: []RevisionRecord{
			{
				Version: "1.0",
				Date:    date,
				Author:  "John Doe",
				Trigger: TriggerInitial,
				Changes: []string{"Initial draft", "Added personas"},
			},
		},
	}

	result := doc.generateRevisionHistory()

	checks := []string{
		"## Revision History",
		"| Version | Date | Author | Trigger | Changes |",
		"1.0",
		"2024-01-10",
		"John Doe",
		"initial",
		"Initial draft; Added personas",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateRevisionHistory() missing expected content: %q", check)
		}
	}
}

func TestGenerateNonGoals(t *testing.T) {
	doc := &Document{
		NonGoals: []NonGoal{
			{
				ID:          "NG-001",
				Title:       "Mobile App",
				Description: "Native mobile application",
				Rationale:   "Focus on web first",
				FuturePhase: "Phase 2",
			},
			{
				ID:    "NG-002",
				Title: "Multi-language support",
			},
		},
	}

	result := doc.generateNonGoals()

	checks := []string{
		"## Non-Goals",
		"| ID | Title | Description | Rationale | Future Phase |",
		"NG-001",
		"Mobile App",
		"Native mobile application",
		"Focus on web first",
		"Phase 2",
		"NG-002",
		"Multi-language support",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateNonGoals() missing expected content: %q", check)
		}
	}
}

func TestGenerateSuccessMetrics(t *testing.T) {
	doc := &Document{
		SuccessMetrics: &SuccessMetrics{
			NorthStar: []Metric{
				{
					ID:                "NS-001",
					Name:              "User Retention",
					Description:       "Monthly active users returning",
					Baseline:          "65%",
					Target:            "80%",
					MeasurementMethod: "Analytics dashboard",
				},
			},
			Supporting: []Metric{
				{
					ID:       "SUP-001",
					Name:     "NPS Score",
					Baseline: "40",
					Target:   "60",
				},
			},
			Guardrail: []Metric{
				{
					ID:       "GR-001",
					Name:     "Page Load Time",
					Baseline: "2s",
					Target:   "<3s",
				},
			},
		},
	}

	result := doc.generateSuccessMetrics()

	checks := []string{
		"## Success Metrics",
		"### North Star Metrics",
		"*Primary metrics that define success.*",
		"| ID | Name | Description | Baseline | Target | Measurement Method |",
		"NS-001",
		"User Retention",
		"Monthly active users returning",
		"65%",
		"80%",
		"Analytics dashboard",
		"### Supporting Metrics",
		"*Metrics that support the north star metrics.*",
		"SUP-001",
		"NPS Score",
		"### Guardrail Metrics",
		"*Metrics that should not degrade.*",
		"GR-001",
		"Page Load Time",
	}

	for _, check := range checks {
		if !strings.Contains(result, check) {
			t.Errorf("generateSuccessMetrics() missing expected content: %q", check)
		}
	}
}

func TestToMarkdownIncludesExtendedSections(t *testing.T) {
	doc := &Document{
		Metadata: Metadata{
			ID:      "PRD-001",
			Title:   "Test PRD",
			Version: "1.0",
			Status:  StatusDraft,
		},
		ExecutiveSummary: ExecutiveSummary{
			ProblemStatement: "Test problem",
			ProposedSolution: "Test solution",
		},
		Problem: &ProblemDefinition{
			Statement: "Detailed problem",
		},
		Market: &MarketDefinition{
			Differentiation: []string{"Key differentiator"},
		},
		Solution: &SolutionDefinition{
			SolutionRationale: "Why this solution",
		},
		Decisions: &DecisionsDefinition{
			Records: []DecisionRecord{
				{ID: "DEC-001", Decision: "Test decision"},
			},
		},
		Reviews: &ReviewsDefinition{
			Decision: ReviewApprove,
		},
		RevisionHistory: []RevisionRecord{
			{Version: "1.0", Changes: []string{"Initial"}},
		},
		NonGoals: []NonGoal{
			{ID: "NG-001", Title: "Not doing this"},
		},
		SuccessMetrics: &SuccessMetrics{
			NorthStar: []Metric{
				{ID: "M-001", Name: "Primary Metric", Target: "100%"},
			},
		},
	}

	opts := DefaultMarkdownOptions()
	result := doc.ToMarkdown(opts)

	// Check TOC includes new sections
	tocChecks := []string{
		"[Problem Definition](#problem-definition)",
		"[Market Analysis](#market-analysis)",
		"[Solution](#solution)",
		"[Decisions](#decisions)",
		"[Reviews](#reviews)",
		"[Revision History](#revision-history)",
		"[Non-Goals](#non-goals)",
		"[Success Metrics](#success-metrics)",
	}

	for _, check := range tocChecks {
		if !strings.Contains(result, check) {
			t.Errorf("ToMarkdown() TOC missing: %q", check)
		}
	}

	// Check section headers are present
	sectionChecks := []string{
		"## Problem Definition",
		"## Market Analysis",
		"## Solution",
		"## Decisions",
		"## Reviews",
		"## Revision History",
		"## Non-Goals",
		"## Success Metrics",
	}

	for _, check := range sectionChecks {
		if !strings.Contains(result, check) {
			t.Errorf("ToMarkdown() missing section: %q", check)
		}
	}
}
