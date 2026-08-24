package assessment

import (
	"fmt"
	"sort"

	"github.com/grokify/prism-roadmap/prioritization"
)

// RankInput is one opportunity's resolved MoSCoW/RICE inputs for ranking —
// the output of ResolveMoSCoWPriority and ComputeRICE for one
// OpportunityAssessment.
type RankInput struct {
	// AssessmentID is the OpportunityAssessment.ID this input was resolved
	// from.
	AssessmentID string `json:"assessmentId"`

	// Title is denormalized for display in a ranked list.
	Title string `json:"title,omitempty"`

	MoSCoW prioritization.MoSCoWPriority `json:"moscow"`
	RICE   RICEScoreResult               `json:"rice"`
}

// ToRankInput resolves this assessment's MoSCoW answers and RICE score into
// a RankInput for RankingPolicy.Rank. When a COMPASS-RICE assessment is
// present (a.Compass != nil), it is resolved via ResolveCompassRICE and
// takes precedence over the legacy ladder RICE — the two use different
// Reach scales (0-100 banded vs. 0..1 fraction) and are never mixed for one
// opportunity. Otherwise falls back to the legacy ComputeRICE(*a.RICE) path
// unchanged. If neither is recorded, RICE.Computable is false with an
// explanatory Reason, matching ComputeRICE's own "never fabricate a score"
// discipline.
func (a *OpportunityAssessment) ToRankInput() RankInput {
	rice := RICEScoreResult{Reason: "no RICE assessment recorded"}
	switch {
	case a.Compass != nil:
		rice = ResolveCompassRICE(a.Compass)
	case a.RICE != nil:
		rice = ComputeRICE(*a.RICE)
	}
	return RankInput{
		AssessmentID: a.ID,
		Title:        a.Title,
		MoSCoW:       a.MoSCoW(),
		RICE:         rice,
	}
}

// ExclusionReason explains why an opportunity did not receive a
// CalculatedRank. An excluded opportunity is never silently dropped from
// RankingPolicy.Rank's return value — it is carried through with the
// reason, per this stack's "no silent caps" convention.
type ExclusionReason string

const (
	// ExclusionWont means MoSCoW resolved to Won't/Not Now (or was never
	// assessed) — outside the current planning horizon by design, not a
	// data gap.
	ExclusionWont ExclusionReason = "wont"

	// ExclusionRICEUncomputable means RICE could not be computed (missing
	// evidence, unresolved Impact/Confidence, insufficient-confidence
	// review state, or a failed effort estimability gate). See
	// RankedOpportunity.RICE.Reason for the specific cause.
	ExclusionRICEUncomputable ExclusionReason = "rice_uncomputable"
)

// RankedOpportunity is one input's place in the calculated ranking, or its
// exclusion reason if it could not be ranked.
type RankedOpportunity struct {
	AssessmentID string `json:"assessmentId"`
	Title        string `json:"title,omitempty"`

	MoSCoW prioritization.MoSCoWPriority `json:"moscow"`
	RICE   RICEScoreResult               `json:"rice"`

	// CalculatedRank is the 1-indexed position in the deterministic
	// portfolio-wide ordering (MoSCoW tier first, RICE score descending
	// within tier). Zero when Excluded is set.
	CalculatedRank int `json:"calculatedRank,omitempty"`

	// TiedWith lists AssessmentIDs in the same MoSCoW tier whose RICE score
	// is within the ranking policy's tie-equivalence band of this one —
	// recorded so a portfolio review can apply a tie-break (e.g. Kano) as a
	// deliberate product decision rather than trusting score precision
	// that isn't really there (prism-roadmap TRD D3: "don't pretend 0.314
	// is meaningfully better than 0.309").
	TiedWith []string `json:"tiedWith,omitempty"`

	// Excluded is non-empty when this opportunity did not receive a
	// CalculatedRank.
	Excluded ExclusionReason `json:"excluded,omitempty"`
}

// RankingPolicy is the deterministic ranking algorithm: MoSCoW tier orders
// first, RICE score orders within a tier. Kano, Market Investment Horizon,
// and custom strategic themes never enter this computation — they describe
// portfolio composition and inform human tie-break decisions, never
// automatic rank (prism-roadmap PRD: "only MoSCoW + RICE determine rank").
type RankingPolicy struct {
	// ID and Version identify this policy for provenance — a
	// recomputation records which policy version produced it, so a later
	// policy change (e.g. a different tie band) doesn't silently
	// reinterpret historical ranks.
	ID      string `json:"id"`
	Version string `json:"version"`

	// TieBandFraction is the RICE-score-relative equivalence band: two
	// scores within this fraction of the larger one are considered tied.
	TieBandFraction float64 `json:"tieBandFraction"`
}

// DefaultRankingPolicy is the reference ranking policy: MoSCoW tier, RICE
// descending within tier, ±5% tie band (prism-roadmap TRD D3).
func DefaultRankingPolicy() RankingPolicy {
	return RankingPolicy{
		ID:              "moscow-rice-v1",
		Version:         "1.0",
		TieBandFraction: 0.05,
	}
}

// Validate returns an error if the policy is not usable.
func (p RankingPolicy) Validate() error {
	if p.ID == "" {
		return fmt.Errorf("id is required")
	}
	if p.Version == "" {
		return fmt.Errorf("version is required")
	}
	if p.TieBandFraction < 0 {
		return fmt.Errorf("tieBandFraction must be >= 0")
	}
	return nil
}

// Rank orders inputs by this policy: MoSCoW tier (Must, Should, Could)
// first, then RICE score descending within tier. Items resolving to
// MoSCoWWontHave/unspecified or with uncomputable RICE are excluded from
// the ordering (see RankedOpportunity.Excluded) rather than silently sorted
// to the bottom with a fabricated low score. The returned slice always
// contains every input: rankable items first in rank order, then excluded
// items in their original input order.
func (p RankingPolicy) Rank(inputs []RankInput) []RankedOpportunity {
	var rankable, excluded []RankedOpportunity

	for _, in := range inputs {
		ro := RankedOpportunity{
			AssessmentID: in.AssessmentID,
			Title:        in.Title,
			MoSCoW:       in.MoSCoW,
			RICE:         in.RICE,
		}
		switch {
		case in.MoSCoW == prioritization.MoSCoWWontHave || in.MoSCoW == prioritization.MoSCoWUnspecified:
			ro.Excluded = ExclusionWont
			excluded = append(excluded, ro)
		case !in.RICE.Computable:
			ro.Excluded = ExclusionRICEUncomputable
			excluded = append(excluded, ro)
		default:
			rankable = append(rankable, ro)
		}
	}

	sort.SliceStable(rankable, func(i, j int) bool {
		wi, wj := rankable[i].MoSCoW.Weight(), rankable[j].MoSCoW.Weight()
		if wi != wj {
			return wi > wj // higher weight (Must=4) ranks first
		}
		return rankable[i].RICE.Score > rankable[j].RICE.Score
	})

	for i := range rankable {
		rankable[i].CalculatedRank = i + 1
	}
	p.annotateTies(rankable)

	return append(rankable, excluded...)
}

// annotateTies fills TiedWith for same-tier items within the policy's tie
// band. rankable must already be sorted by Rank.
func (p RankingPolicy) annotateTies(rankable []RankedOpportunity) {
	for i := range rankable {
		for j := range rankable {
			if i == j || rankable[i].MoSCoW != rankable[j].MoSCoW {
				continue
			}
			if p.withinTieBand(rankable[i].RICE.Score, rankable[j].RICE.Score) {
				rankable[i].TiedWith = append(rankable[i].TiedWith, rankable[j].AssessmentID)
			}
		}
	}
}

// withinTieBand reports whether two RICE scores are within TieBandFraction
// of the larger one.
func (p RankingPolicy) withinTieBand(a, b float64) bool {
	if a == 0 && b == 0 {
		return true
	}
	larger := a
	if b > larger {
		larger = b
	}
	if larger <= 0 {
		return false
	}
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff/larger <= p.TieBandFraction
}

// RankOverride is an explicit governance decision to move an opportunity
// away from its CalculatedRank — recorded as new evidence, never produced
// by quietly reweighting RICE/MoSCoW inputs to manufacture a desired number
// (ideation doc: "don't force leadership to manipulate RICE to manufacture
// the desired answer"). Call Validate before persisting one; ApplyOverrides
// itself does not enforce it, matching this package's convention of keeping
// deterministic transforms separate from field validation.
type RankOverride struct {
	AssessmentID string `json:"assessmentId"`
	FinalRank    int    `json:"finalRank"`
	Rationale    string `json:"rationale"`
	ApprovedBy   string `json:"approvedBy"`

	// EvidenceIDs optionally cite the evidence motivating the override
	// (e.g. a new contractual commitment), so it is auditable like any
	// other ranking input.
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// Validate returns an error if required fields are missing.
func (o RankOverride) Validate() error {
	if o.AssessmentID == "" {
		return fmt.Errorf("assessmentId is required")
	}
	if o.Rationale == "" {
		return fmt.Errorf("rationale is required")
	}
	if o.ApprovedBy == "" {
		return fmt.Errorf("approvedBy is required")
	}
	return nil
}

// OpportunityRank is a ranked opportunity's final, governance-aware
// position: CalculatedRank from RankingPolicy.Rank, and FinalRank equal to
// it unless an explicit RankOverride applies.
type OpportunityRank struct {
	RankedOpportunity
	FinalRank int           `json:"finalRank"`
	Override  *RankOverride `json:"override,omitempty"`
}

// ApplyOverrides produces the final rank list from a calculated ranking
// (RankingPolicy.Rank's output) and a set of overrides keyed by
// AssessmentID. Excluded opportunities are carried through unchanged and
// cannot be overridden — an exclusion means the opportunity lacks a
// computable rank to override; resolve the exclusion (supply the missing
// evidence, pass the estimability gate) and recompute instead.
func ApplyOverrides(ranked []RankedOpportunity, overrides []RankOverride) []OpportunityRank {
	byID := make(map[string]RankOverride, len(overrides))
	for _, o := range overrides {
		byID[o.AssessmentID] = o
	}

	var rankable, excluded []OpportunityRank
	for _, r := range ranked {
		out := OpportunityRank{RankedOpportunity: r, FinalRank: r.CalculatedRank}
		if r.Excluded != "" {
			excluded = append(excluded, out)
			continue
		}
		if o, ok := byID[r.AssessmentID]; ok {
			ov := o
			out.Override = &ov
			out.FinalRank = o.FinalRank
		}
		rankable = append(rankable, out)
	}

	sort.SliceStable(rankable, func(i, j int) bool {
		return rankable[i].FinalRank < rankable[j].FinalRank
	})

	return append(rankable, excluded...)
}

// RankCollisions returns FinalRank values shared by more than one
// non-excluded opportunity, sorted ascending — a sign that overrides need
// reconciling before this ranking is presented (a portfolio review should
// never show two items both at #3).
func RankCollisions(ranks []OpportunityRank) []int {
	counts := make(map[int]int)
	for _, r := range ranks {
		if r.Excluded != "" {
			continue
		}
		counts[r.FinalRank]++
	}
	var collisions []int
	for rank, count := range counts {
		if count > 1 {
			collisions = append(collisions, rank)
		}
	}
	sort.Ints(collisions)
	return collisions
}
