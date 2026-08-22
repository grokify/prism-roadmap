package assessment

import (
	"testing"
	"time"

	"github.com/plexusone/structured-evaluation/rubric"
)

func TestHasEvidence(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)

	t.Run("no evidence anywhere", func(t *testing.T) {
		a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
		if a.HasEvidence() {
			t.Error("expected HasEvidence() = false for an empty assessment")
		}
	})

	t.Run("moscow answer evidence", func(t *testing.T) {
		a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
		a.MoSCoWAnswers = []ThresholdAnswer{{LevelID: "must", Satisfied: true, EvidenceIDs: []string{"EV-1"}}}
		if !a.HasEvidence() {
			t.Error("expected HasEvidence() = true from MoSCoWAnswers")
		}
	})

	t.Run("reach evidence", func(t *testing.T) {
		a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
		a.RICE = &RICEAssessment{Reach: Reach{Fraction: 0.5, EvidenceIDs: []string{"EV-1"}}}
		if !a.HasEvidence() {
			t.Error("expected HasEvidence() = true from Reach")
		}
	})

	t.Run("impact answer evidence", func(t *testing.T) {
		a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
		a.RICE = &RICEAssessment{ImpactAnswers: []ThresholdAnswer{{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-1"}}}}
		if !a.HasEvidence() {
			t.Error("expected HasEvidence() = true from ImpactAnswers")
		}
	})

	t.Run("dimension answer evidence", func(t *testing.T) {
		a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
		a.Dimensions = []DimensionAssignment{{DimensionID: "kano", Answers: []DimensionAnswer{
			{OptionID: "must_be", Answer: true, EvidenceIDs: []string{"EV-1"}},
		}}}
		if !a.HasEvidence() {
			t.Error("expected HasEvidence() = true from Dimensions[].Answers")
		}
	})

	t.Run("okr contribution evidence", func(t *testing.T) {
		a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
		a.Contributions = []OKRContribution{{ObjectiveID: "OBJ-1", Strength: ContributionHigh, EvidenceIDs: []string{"EV-1"}}}
		if !a.HasEvidence() {
			t.Error("expected HasEvidence() = true from Contributions")
		}
	})
}

func TestHasRubricAnswers(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)

	empty := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	if empty.HasRubricAnswers() {
		t.Error("expected HasRubricAnswers() = false for an empty assessment")
	}

	withMoSCoW := NewOpportunityAssessment("OA-2", OpportunityRef{SpecID: "OPP-2"}, "Title", now)
	withMoSCoW.MoSCoWAnswers = []ThresholdAnswer{{LevelID: "must", Satisfied: true}}
	if !withMoSCoW.HasRubricAnswers() {
		t.Error("expected HasRubricAnswers() = true from MoSCoWAnswers")
	}

	withRICE := NewOpportunityAssessment("OA-3", OpportunityRef{SpecID: "OPP-3"}, "Title", now)
	withRICE.RICE = &RICEAssessment{ConfidenceAnswers: []ThresholdAnswer{{LevelID: "high", Satisfied: true}}}
	if !withRICE.HasRubricAnswers() {
		t.Error("expected HasRubricAnswers() = true from RICE.ConfidenceAnswers")
	}

	withDimension := NewOpportunityAssessment("OA-4", OpportunityRef{SpecID: "OPP-4"}, "Title", now)
	withDimension.Dimensions = []DimensionAssignment{{DimensionID: "kano", Answers: []DimensionAnswer{{OptionID: "x", Answer: true}}}}
	if !withDimension.HasRubricAnswers() {
		t.Error("expected HasRubricAnswers() = true from Dimensions[].Answers")
	}
}

func TestNewOpportunityReportSectionOrder(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)

	report := NewOpportunityReport(now, *a, nil)

	wantSectionIDs := []string{"recommendation", "definition", "prioritization", "portfolio-context", "strategic-capability", "risks-decision"}
	if len(report.Sections) != len(wantSectionIDs) {
		t.Fatalf("Sections = %+v, want %d entries", report.Sections, len(wantSectionIDs))
	}
	for i, want := range wantSectionIDs {
		if report.Sections[i].ID != want {
			t.Errorf("Sections[%d].ID = %q, want %q", i, report.Sections[i].ID, want)
		}
	}

	wantAppendixIDs := []string{"evidence-log", "rubric-trace", "provenance", "related"}
	if len(report.Appendices) != len(wantAppendixIDs) {
		t.Fatalf("Appendices = %+v, want %d entries", report.Appendices, len(wantAppendixIDs))
	}
	for i, want := range wantAppendixIDs {
		if report.Appendices[i].ID != want {
			t.Errorf("Appendices[%d].ID = %q, want %q", i, report.Appendices[i].ID, want)
		}
	}
}

func TestNewOpportunityReportPresenceForMinimalAssessment(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)

	report := NewOpportunityReport(now, *a, nil)

	alwaysOn := map[string]bool{"recommendation": true, "definition": true, "risks-decision": true}
	for _, s := range report.Sections {
		want := alwaysOn[s.ID]
		if s.Present != want {
			t.Errorf("section %q Present = %v, want %v for a minimal assessment", s.ID, s.Present, want)
		}
	}
	for _, a := range report.Appendices {
		if a.ID == "related" {
			continue // always false from this constructor by design
		}
		if a.Present {
			t.Errorf("appendix %q Present = true, want false for a minimal assessment", a.ID)
		}
	}
}

func TestNewOpportunityReportPresenceForFullAssessment(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	a.MoSCoWAnswers = []ThresholdAnswer{{LevelID: "must", Satisfied: true, EvidenceIDs: []string{"EV-1"}}}
	a.RICE = &RICEAssessment{Reach: Reach{Fraction: 0.5, EvidenceIDs: []string{"EV-2"}}}
	a.Dimensions = []DimensionAssignment{NewDimensionAssignment(KanoDimension(), nil)}
	a.Contributions = []OKRContribution{{ObjectiveID: "OBJ-1", Strength: ContributionHigh}}
	a.Capabilities = []CapabilityReference{{CapabilityID: "auth", Relation: CapabilityEnables}}
	a.Judge = rubric.NewJudgeMetadata("claude-sonnet-5")

	report := NewOpportunityReport(now, *a, nil)

	for _, s := range report.Sections {
		if !s.Present {
			t.Errorf("section %q Present = false, want true for a fully-populated assessment", s.ID)
		}
	}
	for _, appendix := range report.Appendices {
		if appendix.ID == "related" {
			continue
		}
		if !appendix.Present {
			t.Errorf("appendix %q Present = false, want true for a fully-populated assessment", appendix.ID)
		}
	}
}

func TestPresentSectionsAndAppendices(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	report := NewOpportunityReport(now, *a, nil)

	present := report.PresentSections()
	for _, s := range present {
		if !s.Present {
			t.Errorf("PresentSections() included a non-present section: %+v", s)
		}
	}
	if len(present) >= len(report.Sections) {
		t.Errorf("expected PresentSections() to filter out at least the not-yet-populated sections, got %d of %d", len(present), len(report.Sections))
	}

	presentAppendices := report.PresentAppendices()
	for _, appendix := range presentAppendices {
		if !appendix.Present {
			t.Errorf("PresentAppendices() included a non-present appendix: %+v", appendix)
		}
	}
}

func TestOpportunityReportValidate(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	valid := NewOpportunityReport(now, *a, nil)
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if err := (OpportunityReport{}).Validate(); err == nil {
		t.Error("expected error for zero-value report")
	}

	missingAssessmentID := valid
	missingAssessmentID.Assessment.ID = ""
	if err := missingAssessmentID.Validate(); err == nil {
		t.Error("expected error propagated from Assessment.Validate()")
	}
}
