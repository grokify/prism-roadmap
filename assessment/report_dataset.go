package assessment

import (
	"fmt"
	"sort"
	"time"
)

// ReportDataset is the computed-facts input to both report renderers (the
// opportunity 6-pager and the portfolio review — RMI-PRISMROADMAP-011/012).
// It is produced by omniroadmap's compile step (RMI-OMNIROADMAP-004) from
// the corpus of OpportunityAssessments, using the deterministic functions in
// this file; every report render is a pure function of one ReportDataset —
// no fact is computed by a renderer, and no number is invented by an LLM
// narrative pass (prism-roadmap PRD: "the report is always a pure function
// of the IR").
type ReportDataset struct {
	// GeneratedAt is supplied by the caller, not computed internally — this
	// stack avoids hidden time.Now()-style nondeterminism so the same
	// inputs always reproduce the same dataset.
	GeneratedAt time.Time `json:"generatedAt"`

	// RankingPolicyID/Version identify which ranking policy produced
	// Ranking, so a later policy change doesn't silently reinterpret a
	// stored dataset.
	RankingPolicyID      string `json:"rankingPolicyId"`
	RankingPolicyVersion string `json:"rankingPolicyVersion"`

	// Ranking is the full calculated-and-governed ranking
	// (RankingPolicy.Rank + ApplyOverrides output) across every assessed
	// opportunity in this compile.
	Ranking []OpportunityRank `json:"ranking"`

	// Distributions are % Person-Day breakdowns per portfolio dimension
	// (Kano, Market Investment Horizon, custom) — descriptive portfolio
	// composition, never a ranking input.
	Distributions []DimensionDistribution `json:"distributions,omitempty"`

	// CapabilityOverlay is % Person-Day investment per referenced
	// capability — the Engineering-facing "what does this roadmap do for
	// us" view (ideation doc).
	CapabilityOverlay []CapabilityInvestment `json:"capabilityOverlay,omitempty"`

	// ObjectiveInvestment is % Person-Day investment per OKR objective —
	// answers "does the roadmap reflect the strategy we say we're
	// pursuing."
	ObjectiveInvestment []ObjectiveInvestment `json:"objectiveInvestment,omitempty"`

	// Deltas compares this dataset against the previous review cycle, if
	// any. Nil for a first review.
	Deltas *ReportDeltas `json:"deltas,omitempty"`

	// OverrideLog lists every governance override applied in this compile
	// (see OverridesFromRanking), so a portfolio review can show
	// calculated-vs-final and why.
	OverrideLog []RankOverride `json:"overrideLog,omitempty"`
}

// NewReportDataset assembles a ReportDataset from a calculated-and-governed
// ranking, deriving OverrideLog from it. Callers add Distributions,
// CapabilityOverlay, ObjectiveInvestment, and Deltas via the Compute*
// functions below.
func NewReportDataset(generatedAt time.Time, policy RankingPolicy, ranking []OpportunityRank) ReportDataset {
	return ReportDataset{
		GeneratedAt:          generatedAt,
		RankingPolicyID:      policy.ID,
		RankingPolicyVersion: policy.Version,
		Ranking:              ranking,
		OverrideLog:          OverridesFromRanking(ranking),
	}
}

// Validate returns an error if required fields are missing.
func (d ReportDataset) Validate() error {
	if d.GeneratedAt.IsZero() {
		return fmt.Errorf("generatedAt is required")
	}
	if d.RankingPolicyID == "" {
		return fmt.Errorf("rankingPolicyId is required")
	}
	if d.RankingPolicyVersion == "" {
		return fmt.Errorf("rankingPolicyVersion is required")
	}
	return nil
}

// OverridesFromRanking extracts every applied override from a ranking list.
func OverridesFromRanking(ranking []OpportunityRank) []RankOverride {
	var out []RankOverride
	for _, r := range ranking {
		if r.Override != nil {
			out = append(out, *r.Override)
		}
	}
	return out
}

// TotalPersonDays sums PersonDays() across the given assessments.
func TotalPersonDays(assessments []*OpportunityAssessment) float64 {
	var total float64
	for _, a := range assessments {
		total += a.PersonDays()
	}
	return total
}

// DistributionBucket is one option's share of a dimension's distribution.
type DistributionBucket struct {
	OptionID         string  `json:"optionId"`
	PersonDays       float64 `json:"personDays"`
	OpportunityCount int     `json:"opportunityCount"`

	// Fraction is PersonDays / the total Person-Days across every
	// assessment passed to ComputeDimensionDistribution (not just
	// classified ones) — for a Category dimension, bucket fractions plus
	// UnclassifiedPersonDays/total sum to 1.0; for a Tags dimension,
	// fractions can sum past 1.0 since one opportunity may occupy several
	// buckets (prism-roadmap PRD FR4).
	Fraction float64 `json:"fraction"`
}

// DimensionDistribution is one portfolio dimension's composition across a
// set of assessments, weighted by Person-Days (ideation doc: "% Effort is
// probably more important than % projects" — ten 2-PD items shouldn't look
// more significant than one 100-PD item).
type DimensionDistribution struct {
	DimensionID      string `json:"dimensionId"`
	DimensionVersion string `json:"dimensionVersion,omitempty"`

	Buckets []DistributionBucket `json:"buckets,omitempty"`

	// UnclassifiedPersonDays is Person-Days from assessments with no
	// assignment recorded for this dimension — surfaced explicitly rather
	// than silently excluded from the percentages (no silent caps).
	UnclassifiedPersonDays float64 `json:"unclassifiedPersonDays,omitempty"`
}

// ComputeDimensionDistribution aggregates Person-Days by option for one
// dimension across the given assessments. Assessments with no RICE
// assessment recorded contribute 0 Person-Days (visible in the totals, not
// silently skipped).
func ComputeDimensionDistribution(dimensionID string, assessments []*OpportunityAssessment) DimensionDistribution {
	dist := DimensionDistribution{DimensionID: dimensionID}
	bucketByOption := make(map[string]*DistributionBucket)
	total := TotalPersonDays(assessments)

	for _, a := range assessments {
		pd := a.PersonDays()

		assignment := a.DimensionAssignment(dimensionID)
		if assignment == nil {
			dist.UnclassifiedPersonDays += pd
			continue
		}
		if dist.DimensionVersion == "" {
			dist.DimensionVersion = assignment.DimensionVersion
		}

		optionIDs := assignment.SelectedOptionIDs()
		if len(optionIDs) == 0 {
			dist.UnclassifiedPersonDays += pd
			continue
		}
		for _, optID := range optionIDs {
			b, ok := bucketByOption[optID]
			if !ok {
				b = &DistributionBucket{OptionID: optID}
				bucketByOption[optID] = b
			}
			b.PersonDays += pd
			b.OpportunityCount++
		}
	}

	for _, b := range bucketByOption {
		if total > 0 {
			b.Fraction = b.PersonDays / total
		}
		dist.Buckets = append(dist.Buckets, *b)
	}
	sort.Slice(dist.Buckets, func(i, j int) bool { return dist.Buckets[i].OptionID < dist.Buckets[j].OptionID })
	return dist
}

// CapabilityInvestment is one capability's total Person-Day investment
// across the assessments that reference it (any relation).
type CapabilityInvestment struct {
	CapabilityID   string   `json:"capabilityId"`
	PersonDays     float64  `json:"personDays"`
	Fraction       float64  `json:"fraction"`
	OpportunityIDs []string `json:"opportunityIds,omitempty"`
}

// ObjectiveInvestment is one OKR objective's total Person-Day investment
// across the assessments contributing to it.
type ObjectiveInvestment struct {
	ObjectiveID    string   `json:"objectiveId"`
	PersonDays     float64  `json:"personDays"`
	Fraction       float64  `json:"fraction"`
	OpportunityIDs []string `json:"opportunityIds,omitempty"`
}

// idInvestment is the shared aggregation shape behind ComputeCapabilityOverlay
// and ComputeObjectiveInvestment — both are "sum Person-Days by an ID this
// assessment references," differing only in which IDs and which exported
// struct the result maps onto.
type idInvestment struct {
	ID             string
	PersonDays     float64
	Fraction       float64
	OpportunityIDs []string
}

// aggregateByID sums Person-Days across assessments, grouped by the IDs
// idsFor extracts from each one (e.g. capability IDs, objective IDs).
func aggregateByID(assessments []*OpportunityAssessment, idsFor func(*OpportunityAssessment) []string) []idInvestment {
	total := TotalPersonDays(assessments)
	byID := make(map[string]*idInvestment)

	for _, a := range assessments {
		pd := a.PersonDays()
		for _, id := range idsFor(a) {
			inv, ok := byID[id]
			if !ok {
				inv = &idInvestment{ID: id}
				byID[id] = inv
			}
			inv.PersonDays += pd
			inv.OpportunityIDs = append(inv.OpportunityIDs, a.ID)
		}
	}

	var out []idInvestment
	for _, inv := range byID {
		sort.Strings(inv.OpportunityIDs)
		if total > 0 {
			inv.Fraction = inv.PersonDays / total
		}
		out = append(out, *inv)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out
}

// ComputeCapabilityOverlay aggregates Person-Days by referenced capability
// across the given assessments.
func ComputeCapabilityOverlay(assessments []*OpportunityAssessment) []CapabilityInvestment {
	aggs := aggregateByID(assessments, func(a *OpportunityAssessment) []string {
		ids := make([]string, len(a.Capabilities))
		for i, c := range a.Capabilities {
			ids[i] = c.CapabilityID
		}
		return ids
	})
	out := make([]CapabilityInvestment, len(aggs))
	for i, agg := range aggs {
		out[i] = CapabilityInvestment{
			CapabilityID:   agg.ID,
			PersonDays:     agg.PersonDays,
			Fraction:       agg.Fraction,
			OpportunityIDs: agg.OpportunityIDs,
		}
	}
	return out
}

// ComputeObjectiveInvestment aggregates Person-Days by OKR objective across
// the given assessments' Contributions.
func ComputeObjectiveInvestment(assessments []*OpportunityAssessment) []ObjectiveInvestment {
	aggs := aggregateByID(assessments, func(a *OpportunityAssessment) []string {
		ids := make([]string, len(a.Contributions))
		for i, c := range a.Contributions {
			ids[i] = c.ObjectiveID
		}
		return ids
	})
	out := make([]ObjectiveInvestment, len(aggs))
	for i, agg := range aggs {
		out[i] = ObjectiveInvestment{
			ObjectiveID:    agg.ID,
			PersonDays:     agg.PersonDays,
			Fraction:       agg.Fraction,
			OpportunityIDs: agg.OpportunityIDs,
		}
	}
	return out
}

// RankMove is one opportunity's rank change between review cycles.
type RankMove struct {
	AssessmentID string `json:"assessmentId"`
	PreviousRank int    `json:"previousRank"`
	CurrentRank  int    `json:"currentRank"`
}

// DistributionShift is one dimension option's Person-Day-fraction change
// between review cycles.
type DistributionShift struct {
	DimensionID      string  `json:"dimensionId"`
	OptionID         string  `json:"optionId"`
	PreviousFraction float64 `json:"previousFraction"`
	CurrentFraction  float64 `json:"currentFraction"`
}

// ReportDeltas summarizes what changed since the previous review cycle's
// dataset — a recurring review's most useful section (ideation doc:
// "instead of every review starting from zero... What changed, why did it
// change, and does everyone agree with the consequences?").
type ReportDeltas struct {
	PreviousGeneratedAt time.Time `json:"previousGeneratedAt"`

	RankMoves []RankMove `json:"rankMoves,omitempty"`
	Added     []string   `json:"added,omitempty"`   // AssessmentIDs new this cycle
	Removed   []string   `json:"removed,omitempty"` // AssessmentIDs no longer present

	DistributionShifts []DistributionShift `json:"distributionShifts,omitempty"`
}

// ComputeDeltas compares current against previous — an earlier
// ReportDataset for the same portfolio — and returns what changed.
// Excluded opportunities (see RankedOpportunity.Excluded) are not compared,
// since they carry no meaningful rank to move.
func ComputeDeltas(previous, current ReportDataset) ReportDeltas {
	deltas := ReportDeltas{PreviousGeneratedAt: previous.GeneratedAt}

	prevRank := make(map[string]int)
	for _, r := range previous.Ranking {
		if r.Excluded == "" {
			prevRank[r.AssessmentID] = r.FinalRank
		}
	}
	currRank := make(map[string]int)
	for _, r := range current.Ranking {
		if r.Excluded == "" {
			currRank[r.AssessmentID] = r.FinalRank
		}
	}

	for id, cur := range currRank {
		if prev, ok := prevRank[id]; ok {
			if prev != cur {
				deltas.RankMoves = append(deltas.RankMoves, RankMove{AssessmentID: id, PreviousRank: prev, CurrentRank: cur})
			}
		} else {
			deltas.Added = append(deltas.Added, id)
		}
	}
	for id := range prevRank {
		if _, ok := currRank[id]; !ok {
			deltas.Removed = append(deltas.Removed, id)
		}
	}
	sort.Slice(deltas.RankMoves, func(i, j int) bool { return deltas.RankMoves[i].AssessmentID < deltas.RankMoves[j].AssessmentID })
	sort.Strings(deltas.Added)
	sort.Strings(deltas.Removed)

	prevDist := make(map[string]DimensionDistribution, len(previous.Distributions))
	for _, d := range previous.Distributions {
		prevDist[d.DimensionID] = d
	}
	for _, curD := range current.Distributions {
		prevD, ok := prevDist[curD.DimensionID]
		if !ok {
			continue
		}
		prevFractions := make(map[string]float64, len(prevD.Buckets))
		for _, b := range prevD.Buckets {
			prevFractions[b.OptionID] = b.Fraction
		}
		for _, b := range curD.Buckets {
			prevFraction := prevFractions[b.OptionID] // 0 if not previously present
			if prevFraction != b.Fraction {
				deltas.DistributionShifts = append(deltas.DistributionShifts, DistributionShift{
					DimensionID:      curD.DimensionID,
					OptionID:         b.OptionID,
					PreviousFraction: prevFraction,
					CurrentFraction:  b.Fraction,
				})
			}
		}
	}
	sort.Slice(deltas.DistributionShifts, func(i, j int) bool {
		if deltas.DistributionShifts[i].DimensionID != deltas.DistributionShifts[j].DimensionID {
			return deltas.DistributionShifts[i].DimensionID < deltas.DistributionShifts[j].DimensionID
		}
		return deltas.DistributionShifts[i].OptionID < deltas.DistributionShifts[j].OptionID
	})

	return deltas
}
