package assessment

import (
	"fmt"
	"time"
)

// NarrativeSlot is one LLM-fillable prose block within a report, always
// anchored to specific computed facts or cited evidence — "every material
// generated assertion should be traceable either to deterministic
// computation or cited evidence" (ideation doc). A narrative pass fills
// Text once; the anchors travel with it so a reviewer, or a later
// regeneration, can verify the claim wasn't invented.
type NarrativeSlot struct {
	ID string `json:"id"`

	// Text is empty until an LLM narrative pass fills it in — the slot
	// itself, and its anchors, are produced deterministically regardless.
	Text string `json:"text,omitempty"`

	// DerivedFrom names the computed fact(s) this narrative interprets
	// (e.g. "rank.finalRank", "rice.score") — free-form path-like
	// references into the report's data, not a strict schema; renderers
	// document their own path conventions.
	DerivedFrom []string `json:"derivedFrom,omitempty"`

	// EvidenceIDs cite specific Evidence records the narrative draws on,
	// beyond what DerivedFrom implies.
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// ReportSection is one section (or appendix) of a report's deterministic
// table of contents. Present reflects data availability, computed when the
// report is assembled (e.g. no Dimensions recorded -> the portfolio-context
// section is Present=false and a renderer omits it) — the section ORDER
// and ID set never changes between reports; only which sections actually
// render does.
type ReportSection struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	Present bool   `json:"present"`

	Narrative *NarrativeSlot `json:"narrative,omitempty"`
}

// OpportunityReport is the per-opportunity 6-pager's renderer-ready input:
// one assessment (and its rank, if computed) plus the deterministic section
// list with per-section presence and narrative slots. Assembled by
// omniroadmap (RMI-OMNIROADMAP-007) from one OpportunityAssessment and its
// OpportunityRank; the actual markdown/text render is a pure function of
// this struct — the LLM narrative pass fills NarrativeSlot.Text values, it
// does not decide section presence or invent facts (prism-roadmap PRD:
// "the report is always a pure function of the IR").
//
// Unlike ReportDataset (the portfolio-wide compile output,
// RMI-PRISMROADMAP-010), an OpportunityReport is scoped to a single
// opportunity — it is the "audit this one ranking end-to-end" artifact
// (prism-roadmap PRD success criteria: "a contested ranking can be
// defended end-to-end: rank → framework classifications → rubric answers →
// evidence records → source links").
type OpportunityReport struct {
	GeneratedAt time.Time `json:"generatedAt"`

	Assessment OpportunityAssessment `json:"assessment"`

	// Rank is nil if this opportunity has not yet been ranked (e.g. it was
	// excluded, or no compile has run since this assessment cycle).
	Rank *OpportunityRank `json:"rank,omitempty"`

	// Sections is the fixed 6-pager body, in order: Recommendation Summary,
	// Opportunity Definition, Prioritization Rationale, Portfolio Context,
	// Strategic & Capability Alignment, Risks/Assumptions & Decision.
	Sections []ReportSection `json:"sections"`

	// Appendices is the fixed appendix list, in order: Full Evidence Log,
	// Rubric Answer Trace, Assessment Provenance, Related Opportunities.
	Appendices []ReportSection `json:"appendices"`
}

// NewOpportunityReport assembles an OpportunityReport for one assessment
// (and its rank, if already computed), with section/appendix Present flags
// derived from what data is actually recorded. A renderer with additional
// context this package doesn't hold — e.g. the full canvas.OpportunitySpec
// content for risks/assumptions, or cross-opportunity relationships for the
// Related Opportunities appendix — may refine Present further; the values
// here are the deterministic default given only the assessment/rank.
func NewOpportunityReport(generatedAt time.Time, assessment OpportunityAssessment, rank *OpportunityRank) OpportunityReport {
	hasPrioritization := assessment.RICE != nil || len(assessment.MoSCoWAnswers) > 0
	hasPortfolioContext := len(assessment.Dimensions) > 0
	hasStrategicCapability := len(assessment.Contributions) > 0 || len(assessment.Capabilities) > 0

	return OpportunityReport{
		GeneratedAt: generatedAt,
		Assessment:  assessment,
		Rank:        rank,
		Sections: []ReportSection{
			{ID: "recommendation", Title: "Recommendation Summary", Present: true},
			{ID: "definition", Title: "Opportunity Definition", Present: true},
			{ID: "prioritization", Title: "Prioritization Rationale", Present: hasPrioritization},
			{ID: "portfolio-context", Title: "Portfolio Context", Present: hasPortfolioContext},
			{ID: "strategic-capability", Title: "Strategic & Capability Alignment", Present: hasStrategicCapability},
			{ID: "risks-decision", Title: "Risks, Assumptions & Decision", Present: true},
		},
		Appendices: []ReportSection{
			{ID: "evidence-log", Title: "Full Evidence Log", Present: assessment.HasEvidence()},
			{ID: "rubric-trace", Title: "Rubric Answer Trace", Present: assessment.HasRubricAnswers()},
			{ID: "provenance", Title: "Assessment Provenance", Present: assessment.Judge != nil},
			// Related Opportunities requires cross-referencing the full
			// assessment corpus (shared capabilities/objectives), which
			// this single-opportunity constructor doesn't have — a
			// renderer with corpus access sets Present after computing it.
			{ID: "related", Title: "Related Opportunities", Present: false},
		},
	}
}

// Validate returns an error if required fields are missing.
func (r OpportunityReport) Validate() error {
	if r.GeneratedAt.IsZero() {
		return fmt.Errorf("generatedAt is required")
	}
	if err := r.Assessment.Validate(); err != nil {
		return fmt.Errorf("assessment: %w", err)
	}
	return nil
}

// PresentSections returns only the sections with Present == true, in
// order — the actual render list for a renderer.
func (r OpportunityReport) PresentSections() []ReportSection {
	return presentOnly(r.Sections)
}

// PresentAppendices returns only the appendices with Present == true, in
// order.
func (r OpportunityReport) PresentAppendices() []ReportSection {
	return presentOnly(r.Appendices)
}

func presentOnly(sections []ReportSection) []ReportSection {
	var out []ReportSection
	for _, s := range sections {
		if s.Present {
			out = append(out, s)
		}
	}
	return out
}
