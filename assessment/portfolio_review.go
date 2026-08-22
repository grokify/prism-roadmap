package assessment

import "fmt"

// PortfolioReview is the portfolio-wide review's renderer-ready input: one
// ReportDataset (RMI-PRISMROADMAP-010) plus a deterministic agenda with
// per-section presence and narrative slots — the whole-roadmap counterpart
// to OpportunityReport (RMI-PRISMROADMAP-011). Both reuse ReportSection and
// NarrativeSlot, so document and slide renderers share one rendering
// contract.
//
// The agenda is generic by design: it never names "Kano" or "Market
// Investment Horizon" specifically, because ReportDataset.Distributions
// already holds every portfolio dimension — built-in or custom — uniformly.
// A renderer walks whatever dimensions are actually present rather than
// this contract hardcoding which ones exist (prism-roadmap PRD FR4: a new
// custom dimension needs no schema change).
//
// omniroadmap (RMI-OMNIROADMAP-008) assembles this from a ReportDataset;
// the presentation projection (RMI-OMNIROADMAP-009, see
// PresentationProjection) renders the SAME PortfolioReview into slide form
// rather than document form — one dataset, two renders (ideation doc: "the
// presentation becomes a projection of the same document model").
type PortfolioReview struct {
	Dataset ReportDataset `json:"dataset"`

	// Agenda is the fixed section order: Executive Summary, Decision
	// Requested, Prioritized Roadmap, Changes Since Previous Review,
	// Prioritization Methodology, Portfolio Composition, Capability Stack,
	// Strategic Alignment, Key Opportunities, Governance & Overrides,
	// Recommendations, Decisions/Open Questions.
	Agenda []ReportSection `json:"agenda"`

	// Appendices is the fixed appendix order: Complete Ranked Roadmap,
	// Portfolio Distribution Detail, Capability & Objective Investment
	// Detail, Override Log, Methodology.
	Appendices []ReportSection `json:"appendices"`
}

// NewPortfolioReview assembles a PortfolioReview from a ReportDataset, with
// Present flags derived from which parts of the dataset actually have data.
func NewPortfolioReview(dataset ReportDataset) PortfolioReview {
	hasRanking := len(dataset.Ranking) > 0
	hasDeltas := dataset.Deltas != nil
	hasDistributions := len(dataset.Distributions) > 0
	hasCapabilityOverlay := len(dataset.CapabilityOverlay) > 0
	hasObjectiveInvestment := len(dataset.ObjectiveInvestment) > 0
	hasOverrides := len(dataset.OverrideLog) > 0

	return PortfolioReview{
		Dataset: dataset,
		Agenda: []ReportSection{
			{ID: "executive-summary", Title: "Executive Summary", Present: true},
			{ID: "decision-requested", Title: "Decision Requested", Present: true},
			{ID: "prioritized-roadmap", Title: "Prioritized Roadmap", Present: hasRanking},
			{ID: "changes-since-previous", Title: "Changes Since Previous Review", Present: hasDeltas},
			{ID: "methodology", Title: "Prioritization Methodology", Present: true},
			{ID: "portfolio-composition", Title: "Portfolio Composition", Present: hasDistributions},
			{ID: "capability-stack", Title: "Capability Stack", Present: hasCapabilityOverlay},
			{ID: "strategic-alignment", Title: "Strategic Alignment", Present: hasObjectiveInvestment},
			{ID: "key-opportunities", Title: "Key Opportunities", Present: hasRanking},
			{ID: "governance-overrides", Title: "Governance & Overrides", Present: hasOverrides},
			{ID: "recommendations", Title: "Recommendations", Present: true},
			{ID: "decisions-open-questions", Title: "Decisions / Open Questions", Present: true},
		},
		Appendices: []ReportSection{
			{ID: "complete-ranked-roadmap", Title: "Complete Ranked Roadmap", Present: hasRanking},
			{ID: "distribution-detail", Title: "Portfolio Distribution Detail", Present: hasDistributions},
			{ID: "investment-detail", Title: "Capability & Objective Investment Detail", Present: hasCapabilityOverlay || hasObjectiveInvestment},
			{ID: "override-log", Title: "Override Log", Present: hasOverrides},
			{ID: "methodology-detail", Title: "Methodology", Present: true},
		},
	}
}

// Validate returns an error if the underlying dataset is invalid.
func (r PortfolioReview) Validate() error {
	if err := r.Dataset.Validate(); err != nil {
		return fmt.Errorf("dataset: %w", err)
	}
	return nil
}

// PresentAgenda returns only the agenda sections with Present == true, in
// order — the actual render list for a document renderer.
func (r PortfolioReview) PresentAgenda() []ReportSection {
	return presentOnly(r.Agenda)
}

// PresentAppendices returns only the appendices with Present == true, in
// order.
func (r PortfolioReview) PresentAppendices() []ReportSection {
	return presentOnly(r.Appendices)
}

// PresentationSlide is one slide in the portfolio review's presentation
// projection — the same PortfolioReview content, rendered as slides
// instead of document sections (RMI-OMNIROADMAP-009). SourceSectionID
// keeps the slide traceable back to the agenda section it came from, even
// when a renderer splits one data-heavy section (e.g. portfolio-composition
// with many dimensions) across several slides.
type PresentationSlide struct {
	ID              string         `json:"id"`
	SourceSectionID string         `json:"sourceSectionId"`
	Headline        string         `json:"headline"`
	Narrative       *NarrativeSlot `json:"narrative,omitempty"`
}

// PresentationProjection derives the deterministic default deck — one
// slide per present agenda section, in order — from this review. A
// renderer may split or merge slides further (e.g. one slide per
// dimension within portfolio-composition); this is the default layout the
// contract guarantees, not the only valid one.
func (r PortfolioReview) PresentationProjection() []PresentationSlide {
	present := r.PresentAgenda()
	slides := make([]PresentationSlide, len(present))
	for i, s := range present {
		slides[i] = PresentationSlide{
			ID:              s.ID,
			SourceSectionID: s.ID,
			Headline:        s.Title,
			Narrative:       s.Narrative,
		}
	}
	return slides
}
