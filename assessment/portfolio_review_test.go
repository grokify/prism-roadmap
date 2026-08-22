package assessment

import (
	"testing"
	"time"
)

func minimalDataset(t time.Time) ReportDataset {
	return NewReportDataset(t, DefaultRankingPolicy(), nil)
}

func TestNewPortfolioReviewAgendaOrder(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	review := NewPortfolioReview(minimalDataset(now))

	wantAgendaIDs := []string{
		"executive-summary", "decision-requested", "prioritized-roadmap",
		"changes-since-previous", "methodology", "portfolio-composition",
		"capability-stack", "strategic-alignment", "key-opportunities",
		"governance-overrides", "recommendations", "decisions-open-questions",
	}
	if len(review.Agenda) != len(wantAgendaIDs) {
		t.Fatalf("Agenda = %+v, want %d entries", review.Agenda, len(wantAgendaIDs))
	}
	for i, want := range wantAgendaIDs {
		if review.Agenda[i].ID != want {
			t.Errorf("Agenda[%d].ID = %q, want %q", i, review.Agenda[i].ID, want)
		}
	}

	wantAppendixIDs := []string{
		"complete-ranked-roadmap", "distribution-detail", "investment-detail",
		"override-log", "methodology-detail",
	}
	if len(review.Appendices) != len(wantAppendixIDs) {
		t.Fatalf("Appendices = %+v, want %d entries", review.Appendices, len(wantAppendixIDs))
	}
	for i, want := range wantAppendixIDs {
		if review.Appendices[i].ID != want {
			t.Errorf("Appendices[%d].ID = %q, want %q", i, review.Appendices[i].ID, want)
		}
	}
}

func TestNewPortfolioReviewPresenceMinimalDataset(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	review := NewPortfolioReview(minimalDataset(now))

	alwaysOn := map[string]bool{
		"executive-summary": true, "decision-requested": true, "methodology": true,
		"recommendations": true, "decisions-open-questions": true,
	}
	for _, s := range review.Agenda {
		want := alwaysOn[s.ID]
		if s.Present != want {
			t.Errorf("agenda %q Present = %v, want %v for a minimal dataset", s.ID, s.Present, want)
		}
	}

	alwaysOnAppendix := map[string]bool{"methodology-detail": true}
	for _, a := range review.Appendices {
		want := alwaysOnAppendix[a.ID]
		if a.Present != want {
			t.Errorf("appendix %q Present = %v, want %v for a minimal dataset", a.ID, a.Present, want)
		}
	}
}

func TestNewPortfolioReviewPresenceFullDataset(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	dataset := NewReportDataset(now, DefaultRankingPolicy(), []OpportunityRank{
		{
			RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1"},
			FinalRank:         1,
			Override:          &RankOverride{AssessmentID: "OA-1", FinalRank: 1, Rationale: "r", ApprovedBy: "a"},
		},
	})
	dataset.Distributions = []DimensionDistribution{{DimensionID: "kano"}}
	dataset.CapabilityOverlay = []CapabilityInvestment{{CapabilityID: "auth"}}
	dataset.ObjectiveInvestment = []ObjectiveInvestment{{ObjectiveID: "OBJ-1"}}
	dataset.Deltas = &ReportDeltas{PreviousGeneratedAt: now.AddDate(0, -3, 0)}

	review := NewPortfolioReview(dataset)

	for _, s := range review.Agenda {
		if !s.Present {
			t.Errorf("agenda %q Present = false, want true for a fully-populated dataset", s.ID)
		}
	}
	for _, a := range review.Appendices {
		if !a.Present {
			t.Errorf("appendix %q Present = false, want true for a fully-populated dataset", a.ID)
		}
	}
}

func TestPortfolioReviewPresentAgendaAndAppendices(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	review := NewPortfolioReview(minimalDataset(now))

	present := review.PresentAgenda()
	for _, s := range present {
		if !s.Present {
			t.Errorf("PresentAgenda() included a non-present section: %+v", s)
		}
	}
	if len(present) >= len(review.Agenda) {
		t.Errorf("expected PresentAgenda() to filter, got %d of %d", len(present), len(review.Agenda))
	}

	presentAppendices := review.PresentAppendices()
	for _, a := range presentAppendices {
		if !a.Present {
			t.Errorf("PresentAppendices() included a non-present appendix: %+v", a)
		}
	}
}

func TestPortfolioReviewValidate(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	valid := NewPortfolioReview(minimalDataset(now))
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := (PortfolioReview{}).Validate(); err == nil {
		t.Error("expected error for a zero-value review (invalid dataset)")
	}
}

func TestPresentationProjectionMatchesPresentAgenda(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	review := NewPortfolioReview(minimalDataset(now))

	slides := review.PresentationProjection()
	present := review.PresentAgenda()

	if len(slides) != len(present) {
		t.Fatalf("PresentationProjection() = %d slides, want %d (matching PresentAgenda())", len(slides), len(present))
	}
	for i, slide := range slides {
		if slide.ID != present[i].ID || slide.SourceSectionID != present[i].ID || slide.Headline != present[i].Title {
			t.Errorf("slide[%d] = %+v, want to mirror agenda section %+v", i, slide, present[i])
		}
	}
}

func TestPresentationProjectionEmptyForNoPresentSections(t *testing.T) {
	// Construct a review with every section explicitly absent to exercise
	// the zero-slide case.
	review := PortfolioReview{
		Dataset: minimalDataset(time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)),
		Agenda:  []ReportSection{{ID: "x", Title: "X", Present: false}},
	}
	if slides := review.PresentationProjection(); len(slides) != 0 {
		t.Errorf("PresentationProjection() = %+v, want empty when no agenda section is present", slides)
	}
}
