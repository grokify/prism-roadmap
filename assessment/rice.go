package assessment

import (
	"fmt"
	"strings"

	"github.com/grokify/prism-roadmap/prioritization"
)

// ImpactLadder is the RICE Impact classification ladder: an initiative's
// consequence for an affected customer, judged independently of Reach so
// Reach's breadth never inflates Impact's magnitude (ideation doc: "avoid
// putting Reach into Impact"). Level IDs match prioritization.ImpactLevel's
// string values so ResolveImpact can parse them directly — this ladder
// supplies the evidence-backed classification; prioritization.ImpactLevel
// still owns the level→multiplier mapping (3.0/2.0/1.0/0.5/0.25), so the
// numbers are defined in exactly one place.
func ImpactLadder() *Ladder {
	return &Ladder{
		ID:   "rice-impact",
		Name: "RICE Impact",
		Levels: []ThresholdLevel{
			{
				ID: "massive", Label: "Massive",
				Criteria: []string{
					"Without this, an affected customer cannot adopt, cannot continue using, or cannot meet a mandatory requirement",
					"Removes a credible critical/severe security risk",
					"Removes a credible catastrophic/critical availability or data-loss risk",
					"Enables a fundamentally new platform capability that was previously impossible",
				},
			},
			{
				ID: "high", Label: "High",
				Criteria: []string{
					"Removes a major limitation in an existing customer workflow/capability",
					"Produces a substantial measurable reliability/performance improvement",
					"Substantially reduces a significant security/compliance risk",
					"Enables an important new enterprise use case",
					"Eliminates substantial recurring customer operational effort",
				},
			},
			{
				ID: "medium", Label: "Medium",
				Criteria: []string{
					"Produces a clearly identifiable, meaningful, measurable customer/platform outcome, though the existing capability remains usable without it",
				},
			},
			{
				ID: "low", Label: "Low",
				Criteria: []string{
					"Produces an observable but minor improvement without materially changing customer capability or outcomes",
				},
			},
			{
				ID: "minimal", Label: "Minimal",
				Criteria: []string{
					"Produces some identifiable benefit, but the effect is negligible or primarily convenience/cleanup",
				},
			},
		},
	}
}

// ResolveImpact evaluates ImpactLadder() against the given answers and
// returns the resulting prioritization.ImpactLevel. ok is false if no level
// was satisfied with supporting evidence.
func ResolveImpact(answers []ThresholdAnswer) (level prioritization.ImpactLevel, ok bool) {
	lvl, _, evaluated := ImpactLadder().Evaluate(answers)
	if !evaluated {
		return "", false
	}
	parsed, err := prioritization.ParseImpactLevel(lvl.ID)
	if err != nil {
		return "", false
	}
	return parsed, true
}

// ConfidenceLadder is the RICE Confidence classification ladder: the
// strength/completeness of evidence supporting an assessment's Reach and
// Impact claims — an evidence-quality judgment, not the judge's own
// subjective certainty (kept separate: an LLM can be highly confident it
// correctly applied the rubric while concluding evidence is Low — that is
// two different fields, not one). Level IDs match
// prioritization.ConfidenceLevel's string values.
func ConfidenceLadder() *Ladder {
	return &Ladder{
		ID:   "rice-confidence",
		Name: "RICE Confidence",
		Levels: []ThresholdLevel{
			{
				ID: "high", Label: "High", AllCriteria: true,
				Criteria: []string{
					"Reach and Impact claims are supported by direct, current, verifiable evidence",
					"No material unsupported assumptions remain",
					"No contradictory evidence",
				},
			},
			{
				ID: "medium", Label: "Medium",
				Criteria: []string{
					"The most important Reach and Impact claims are supported by credible evidence, but one or more material assumptions remain, or evidence is incomplete/indirect",
				},
			},
			{
				ID: "low", Label: "Low",
				Criteria: []string{
					"One or more material Reach or Impact claims rely substantially on estimates, inference, weak/stale evidence, or unsupported assumptions",
				},
			},
		},
	}
}

// ResolveConfidence evaluates ConfidenceLadder() against the given answers
// and returns the resulting prioritization.ConfidenceLevel. ok is false if
// no level was satisfied with supporting evidence — callers should also
// check ConfidenceInsufficientEvidence first (see RICEAssessment), which is
// a distinct review state rather than an unresolved ladder.
func ResolveConfidence(answers []ThresholdAnswer) (level prioritization.ConfidenceLevel, ok bool) {
	lvl, _, evaluated := ConfidenceLadder().Evaluate(answers)
	if !evaluated {
		return "", false
	}
	parsed, err := prioritization.ParseConfidenceLevel(lvl.ID)
	if err != nil {
		return "", false
	}
	return parsed, true
}

// Reach is the RICE Reach input: the fraction of the relevant customer
// population materially affected, directly or indirectly, within the
// scoring horizon. Unlike Impact/Confidence, Reach is not a ladder — it is
// a continuous percentage backed by evidence (deployment inventory,
// customer/account data), per prism-roadmap PRD FR2.
type Reach struct {
	// Fraction is the reach percentage expressed as 0..1 (0.72 = 72%).
	Fraction float64 `json:"fraction"`

	// Population describes the denominator (e.g. "active paying customer
	// accounts"), so the percentage is reproducible.
	Population string `json:"population,omitempty"`

	Rationale string `json:"rationale,omitempty"`

	// EvidenceIDs cite the supporting Evidence records. Required whenever
	// Fraction > 0 — a claimed nonzero reach with no evidence is exactly
	// the "PM assertion" the ideation doc's evidence ladder ranks weakest.
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// Validate returns an error if Reach is out of range or an unsupported
// nonzero claim.
func (r Reach) Validate() error {
	if r.Fraction < 0 || r.Fraction > 1 {
		return fmt.Errorf("fraction must be within 0..1, got %v", r.Fraction)
	}
	if r.Fraction > 0 && len(r.EvidenceIDs) == 0 {
		return fmt.Errorf("a nonzero reach fraction requires at least one evidence citation")
	}
	return nil
}

// EffortComponent is one decomposed unit of work contributing to an Effort
// estimate (e.g. "Application engineering", 2.0 person-days).
type EffortComponent struct {
	Name       string  `json:"name"`
	PersonDays float64 `json:"personDays"`
}

// EstimabilityGate records whether an initiative's plan had enough detail
// to produce a defensible Effort estimate before it is trusted for RICE. An
// LLM that "happily turns a three-sentence PLAN into 17 PD" gives false
// precision (ideation doc) — the gate exists so under-specified initiatives
// come back with a validation error instead of a fabricated number.
type EstimabilityGate struct {
	ScopeDefined             bool `json:"scopeDefined"`
	ImplementationIdentified bool `json:"implementationIdentified"`
	DependenciesIdentified   bool `json:"dependenciesIdentified"`
	TestingIdentified        bool `json:"testingIdentified"`
	DeploymentIdentified     bool `json:"deploymentIdentified"`
}

// Passed reports whether the plan had enough detail for a defensible
// estimate. All five checks must hold.
func (g EstimabilityGate) Passed() bool {
	return g.ScopeDefined && g.ImplementationIdentified && g.DependenciesIdentified &&
		g.TestingIdentified && g.DeploymentIdentified
}

// MissingChecks lists which gate checks failed, for a judge or reviewer to
// address before an estimate can be trusted. Derived from the booleans
// rather than stored, so it can never drift out of sync with them.
func (g EstimabilityGate) MissingChecks() []string {
	var missing []string
	if !g.ScopeDefined {
		missing = append(missing, "scope not sufficiently defined")
	}
	if !g.ImplementationIdentified {
		missing = append(missing, "major implementation work not identified")
	}
	if !g.DependenciesIdentified {
		missing = append(missing, "dependencies not identified")
	}
	if !g.TestingIdentified {
		missing = append(missing, "testing/acceptance work not identified")
	}
	if !g.DeploymentIdentified {
		missing = append(missing, "deployment/rollout not identified")
	}
	return missing
}

// EffortEstimate is the RICE Effort input: total labor investment (not
// calendar duration — four engineers for one month equals one engineer for
// four months) to deliver an initiative to its defined Definition of Done.
// Effort is Person-Days, not person-months, matching this org's existing
// estimation unit and giving finer resolution.
type EffortEstimate struct {
	Components []EffortComponent `json:"components,omitempty"`

	// Low, Expected, High give an uncertainty range; RICE calculations use
	// Expected.
	Low      float64 `json:"low,omitempty"`
	Expected float64 `json:"expected"`
	High     float64 `json:"high,omitempty"`

	Gate EstimabilityGate `json:"gate"`
}

// TotalPersonDays sums the decomposed components.
func (e EffortEstimate) TotalPersonDays() float64 {
	var total float64
	for _, c := range e.Components {
		total += c.PersonDays
	}
	return total
}

// Validate returns an error if the estimability gate failed or Expected is
// not usable.
func (e EffortEstimate) Validate() error {
	if !e.Gate.Passed() {
		return fmt.Errorf("estimability gate not passed: %s", strings.Join(e.Gate.MissingChecks(), "; "))
	}
	if e.Expected <= 0 {
		return fmt.Errorf("expected person-days must be > 0")
	}
	return nil
}

// RICEAssessment is the evidence-backed RICE input for one opportunity: the
// judge/decomposition output (ladder answers, Reach, EffortEstimate) that
// ComputeRICE resolves into a score. The LLM/judge only ever produces this
// struct's fields — never a score directly (prism-roadmap PRD: "the LLM
// never emits a number").
type RICEAssessment struct {
	Reach             Reach             `json:"reach"`
	ImpactAnswers     []ThresholdAnswer `json:"impactAnswers"`
	ConfidenceAnswers []ThresholdAnswer `json:"confidenceAnswers,omitempty"`

	// ConfidenceInsufficientEvidence marks Confidence as a review state
	// rather than resolved via ConfidenceAnswers. This is NOT the same as
	// RICE Confidence = 0 — treating it as zero would make an
	// under-evidenced but potentially important initiative score as if it
	// had no merit and silently vanish from the roadmap (ideation doc).
	ConfidenceInsufficientEvidence bool `json:"confidenceInsufficientEvidence,omitempty"`

	Effort EffortEstimate `json:"effort"`
}

// RICEScoreResult is the outcome of ComputeRICE: either a numeric score, or
// a reason it could not be computed. Computable is false whenever Reason is
// set — callers should treat an uncomputable result as "needs more
// evidence/detail before ranking," not as a score of zero.
type RICEScoreResult struct {
	Score      float64 `json:"score,omitempty"`
	Computable bool    `json:"computable"`
	Reason     string  `json:"reason,omitempty"`

	Impact     prioritization.ImpactLevel     `json:"impact,omitempty"`
	Confidence prioritization.ConfidenceLevel `json:"confidence,omitempty"`

	// ProfileID is set when this score was produced by ResolveCompassRICE
	// rather than ComputeRICE — e.g. "customer/b2b/v1" — so a ranked list
	// can show which scores are cross-profile-comparable COMPASS-RICE
	// scores versus the legacy single-scale RICE. Empty for ComputeRICE
	// results.
	ProfileID string `json:"profileId,omitempty"`
}

// ComputeRICE deterministically resolves a RICEAssessment's ladder answers
// and effort estimate into a RICE score: (Reach × Impact × Confidence) /
// Effort, using prioritization.ImpactLevel/ConfidenceLevel's existing
// multiplier tables so the numeric mapping is defined in exactly one place.
// This is the only place the formula is evaluated — the judge supplies
// inputs, this function supplies the number (prism-roadmap PRD FR2/FR3).
func ComputeRICE(a RICEAssessment) RICEScoreResult {
	if err := a.Reach.Validate(); err != nil {
		return RICEScoreResult{Reason: "reach: " + err.Error()}
	}

	impact, ok := ResolveImpact(a.ImpactAnswers)
	if !ok {
		return RICEScoreResult{Reason: "no Impact threshold satisfied with supporting evidence"}
	}

	if a.ConfidenceInsufficientEvidence {
		return RICEScoreResult{
			Reason: "confidence marked insufficient evidence — needs human review before ranking",
			Impact: impact,
		}
	}

	confidence, ok := ResolveConfidence(a.ConfidenceAnswers)
	if !ok {
		return RICEScoreResult{
			Reason: "no Confidence threshold satisfied with supporting evidence",
			Impact: impact,
		}
	}

	if err := a.Effort.Validate(); err != nil {
		return RICEScoreResult{
			Reason:     "effort: " + err.Error(),
			Impact:     impact,
			Confidence: confidence,
		}
	}

	score := (a.Reach.Fraction * impact.Multiplier() * confidence.Multiplier()) / a.Effort.Expected
	return RICEScoreResult{
		Score:      score,
		Computable: true,
		Impact:     impact,
		Confidence: confidence,
	}
}
