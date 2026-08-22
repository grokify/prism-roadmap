package assessment

import (
	"testing"
	"time"

	"github.com/grokify/prism-roadmap/prioritization"
)

// assessmentWithEffort builds a minimal assessment with a given Person-Days
// effort (via a RICE assessment whose gate always passes) for aggregation
// tests. Dimensions/Contributions/Capabilities are set by the caller.
func assessmentWithEffort(id string, personDays float64) *OpportunityAssessment {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment(id, OpportunityRef{SpecID: "OPP-" + id}, "Title "+id, now)
	a.RICE = &RICEAssessment{
		Effort: EffortEstimate{Expected: personDays},
	}
	return a
}

func TestOpportunityAssessmentPersonDays(t *testing.T) {
	withRICE := assessmentWithEffort("A", 12)
	if got := withRICE.PersonDays(); got != 12 {
		t.Errorf("PersonDays() = %v, want 12", got)
	}

	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	withoutRICE := NewOpportunityAssessment("B", OpportunityRef{SpecID: "OPP-B"}, "Title", now)
	if got := withoutRICE.PersonDays(); got != 0 {
		t.Errorf("PersonDays() = %v, want 0 when no RICE assessment recorded", got)
	}
}

func TestOpportunityAssessmentDimensionAssignment(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("A", OpportunityRef{SpecID: "OPP-A"}, "Title", now)
	a.Dimensions = []DimensionAssignment{
		NewDimensionAssignment(KanoDimension(), nil),
	}
	if got := a.DimensionAssignment("kano"); got == nil {
		t.Fatal("expected kano assignment to be found")
	}
	if got := a.DimensionAssignment("market-investment-horizon"); got != nil {
		t.Errorf("expected nil for unassigned dimension, got %+v", got)
	}
}

func TestSelectedOptionIDsCategory(t *testing.T) {
	resolved := DimensionAssignment{Category: &CategorySelection{OptionID: "must_be", Resolved: true}}
	if got := resolved.SelectedOptionIDs(); len(got) != 1 || got[0] != "must_be" {
		t.Errorf("SelectedOptionIDs() = %v, want [must_be]", got)
	}

	unresolved := DimensionAssignment{Category: &CategorySelection{}}
	if got := unresolved.SelectedOptionIDs(); len(got) != 0 {
		t.Errorf("SelectedOptionIDs() = %v, want empty for unresolved category", got)
	}
}

func TestSelectedOptionIDsTags(t *testing.T) {
	tagged := DimensionAssignment{Tags: []string{"ai", "growth"}}
	got := tagged.SelectedOptionIDs()
	if len(got) != 2 || got[0] != "ai" || got[1] != "growth" {
		t.Errorf("SelectedOptionIDs() = %v, want [ai growth]", got)
	}
}

func TestComputeDimensionDistributionCategory(t *testing.T) {
	a1 := assessmentWithEffort("A1", 10)
	a1.Dimensions = []DimensionAssignment{{DimensionID: "kano", DimensionVersion: "1.0", Category: &CategorySelection{OptionID: "must_be", Resolved: true}}}

	a2 := assessmentWithEffort("A2", 30)
	a2.Dimensions = []DimensionAssignment{{DimensionID: "kano", DimensionVersion: "1.0", Category: &CategorySelection{OptionID: "performance", Resolved: true}}}

	a3 := assessmentWithEffort("A3", 10) // no kano assignment — unclassified

	dist := ComputeDimensionDistribution("kano", []*OpportunityAssessment{a1, a2, a3})

	if dist.DimensionID != "kano" || dist.DimensionVersion != "1.0" {
		t.Errorf("identity = %+v", dist)
	}
	if dist.UnclassifiedPersonDays != 10 {
		t.Errorf("UnclassifiedPersonDays = %v, want 10", dist.UnclassifiedPersonDays)
	}
	if len(dist.Buckets) != 2 {
		t.Fatalf("expected 2 buckets, got %d: %+v", len(dist.Buckets), dist.Buckets)
	}

	byOption := make(map[string]DistributionBucket)
	for _, b := range dist.Buckets {
		byOption[b.OptionID] = b
	}
	// total = 10 + 30 + 10 = 50
	if b := byOption["must_be"]; b.PersonDays != 10 || b.Fraction != 0.2 || b.OpportunityCount != 1 {
		t.Errorf("must_be bucket = %+v, want personDays=10 fraction=0.2 count=1", b)
	}
	if b := byOption["performance"]; b.PersonDays != 30 || b.Fraction != 0.6 || b.OpportunityCount != 1 {
		t.Errorf("performance bucket = %+v, want personDays=30 fraction=0.6 count=1", b)
	}
	// Category fractions + unclassified fraction should sum to 1.0.
	sumFraction := byOption["must_be"].Fraction + byOption["performance"].Fraction + (dist.UnclassifiedPersonDays / 50)
	if sumFraction != 1.0 {
		t.Errorf("fractions + unclassified = %v, want 1.0", sumFraction)
	}
}

func TestComputeDimensionDistributionTagsCanExceedOne(t *testing.T) {
	a1 := assessmentWithEffort("A1", 10)
	a1.Dimensions = []DimensionAssignment{{DimensionID: "strategic-themes", Tags: []string{"ai", "growth"}}}

	a2 := assessmentWithEffort("A2", 10)
	a2.Dimensions = []DimensionAssignment{{DimensionID: "strategic-themes", Tags: []string{"ai"}}}

	dist := ComputeDimensionDistribution("strategic-themes", []*OpportunityAssessment{a1, a2})

	byOption := make(map[string]DistributionBucket)
	for _, b := range dist.Buckets {
		byOption[b.OptionID] = b
	}
	// total = 20; "ai" appears in both (20 PD, fraction 1.0), "growth" in one (10 PD, fraction 0.5)
	if b := byOption["ai"]; b.PersonDays != 20 || b.Fraction != 1.0 {
		t.Errorf("ai bucket = %+v, want personDays=20 fraction=1.0", b)
	}
	if b := byOption["growth"]; b.PersonDays != 10 || b.Fraction != 0.5 {
		t.Errorf("growth bucket = %+v, want personDays=10 fraction=0.5", b)
	}
	sum := byOption["ai"].Fraction + byOption["growth"].Fraction
	if sum <= 1.0 {
		t.Errorf("expected tag fractions to sum past 1.0, got %v", sum)
	}
}

func TestComputeDimensionDistributionZeroTotal(t *testing.T) {
	dist := ComputeDimensionDistribution("kano", nil)
	if len(dist.Buckets) != 0 || dist.UnclassifiedPersonDays != 0 {
		t.Errorf("expected empty distribution for no assessments, got %+v", dist)
	}
}

func TestComputeCapabilityOverlay(t *testing.T) {
	a1 := assessmentWithEffort("A1", 10)
	a1.Capabilities = []CapabilityReference{{CapabilityID: "authorization", Relation: CapabilityEnables}}

	a2 := assessmentWithEffort("A2", 20)
	a2.Capabilities = []CapabilityReference{
		{CapabilityID: "authorization", Relation: CapabilityImproves},
		{CapabilityID: "audit", Relation: CapabilityEnables},
	}

	overlay := ComputeCapabilityOverlay([]*OpportunityAssessment{a1, a2})
	if len(overlay) != 2 {
		t.Fatalf("expected 2 capabilities, got %d: %+v", len(overlay), overlay)
	}

	byID := make(map[string]CapabilityInvestment)
	for _, inv := range overlay {
		byID[inv.CapabilityID] = inv
	}
	// total = 30
	if inv := byID["authorization"]; inv.PersonDays != 30 || len(inv.OpportunityIDs) != 2 {
		t.Errorf("authorization = %+v, want personDays=30, 2 opportunities", inv)
	}
	if inv := byID["audit"]; inv.PersonDays != 20 || inv.Fraction != 20.0/30.0 {
		t.Errorf("audit = %+v", inv)
	}
}

func TestComputeObjectiveInvestment(t *testing.T) {
	a1 := assessmentWithEffort("A1", 10)
	a1.Contributions = []OKRContribution{{ObjectiveID: "OBJ-1", Strength: ContributionHigh}}

	a2 := assessmentWithEffort("A2", 10)
	a2.Contributions = []OKRContribution{{ObjectiveID: "OBJ-1", Strength: ContributionMedium}}

	overlay := ComputeObjectiveInvestment([]*OpportunityAssessment{a1, a2})
	if len(overlay) != 1 || overlay[0].ObjectiveID != "OBJ-1" || overlay[0].PersonDays != 20 {
		t.Errorf("ComputeObjectiveInvestment() = %+v", overlay)
	}
}

func TestNewReportDatasetDerivesOverrideLog(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	policy := DefaultRankingPolicy()
	ranking := []OpportunityRank{
		{
			RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1", MoSCoW: prioritization.MoSCoWMustHave},
			FinalRank:         1,
			Override:          &RankOverride{AssessmentID: "OA-1", FinalRank: 1, Rationale: "r", ApprovedBy: "a"},
		},
		{
			RankedOpportunity: RankedOpportunity{AssessmentID: "OA-2", MoSCoW: prioritization.MoSCoWShouldHave},
			FinalRank:         2,
		},
	}

	dataset := NewReportDataset(now, policy, ranking)
	if len(dataset.OverrideLog) != 1 || dataset.OverrideLog[0].AssessmentID != "OA-1" {
		t.Errorf("OverrideLog = %+v, want exactly the OA-1 override", dataset.OverrideLog)
	}
	if dataset.RankingPolicyID != policy.ID || dataset.RankingPolicyVersion != policy.Version {
		t.Errorf("policy identity not carried through: %+v", dataset)
	}
}

func TestReportDatasetValidate(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	valid := NewReportDataset(now, DefaultRankingPolicy(), nil)
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := (ReportDataset{}).Validate(); err == nil {
		t.Error("expected error for zero-value dataset")
	}
}

func TestComputeDeltasRankMovesAddedRemoved(t *testing.T) {
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	previous := ReportDataset{
		GeneratedAt: t1,
		Ranking: []OpportunityRank{
			{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1"}, FinalRank: 1},
			{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-2"}, FinalRank: 2},
			{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-3"}, FinalRank: 3}, // removed this cycle
		},
	}
	current := ReportDataset{
		GeneratedAt: t2,
		Ranking: []OpportunityRank{
			{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-2"}, FinalRank: 1}, // moved 2 -> 1
			{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1"}, FinalRank: 2}, // moved 1 -> 2
			{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-4"}, FinalRank: 3}, // added
		},
	}

	deltas := ComputeDeltas(previous, current)
	if deltas.PreviousGeneratedAt != t1 {
		t.Errorf("PreviousGeneratedAt = %v, want %v", deltas.PreviousGeneratedAt, t1)
	}
	if len(deltas.RankMoves) != 2 {
		t.Fatalf("expected 2 rank moves, got %+v", deltas.RankMoves)
	}
	if len(deltas.Added) != 1 || deltas.Added[0] != "OA-4" {
		t.Errorf("Added = %v, want [OA-4]", deltas.Added)
	}
	if len(deltas.Removed) != 1 || deltas.Removed[0] != "OA-3" {
		t.Errorf("Removed = %v, want [OA-3]", deltas.Removed)
	}
}

func TestComputeDeltasExcludesExcludedOpportunities(t *testing.T) {
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	previous := ReportDataset{GeneratedAt: t1, Ranking: []OpportunityRank{
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1", Excluded: ExclusionWont}, FinalRank: 0},
	}}
	current := ReportDataset{GeneratedAt: t2, Ranking: []OpportunityRank{
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1", Excluded: ExclusionWont}, FinalRank: 0},
	}}
	deltas := ComputeDeltas(previous, current)
	if len(deltas.RankMoves) != 0 || len(deltas.Added) != 0 || len(deltas.Removed) != 0 {
		t.Errorf("expected no deltas for an excluded opportunity present both cycles, got %+v", deltas)
	}
}

func TestComputeDeltasDistributionShifts(t *testing.T) {
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)

	previous := ReportDataset{
		GeneratedAt: t1,
		Distributions: []DimensionDistribution{
			{DimensionID: "market-investment-horizon", Buckets: []DistributionBucket{
				{OptionID: "ktlo", Fraction: 0.71},
				{OptionID: "tam_expansion", Fraction: 0.05},
			}},
		},
	}
	current := ReportDataset{
		GeneratedAt: t2,
		Distributions: []DimensionDistribution{
			{DimensionID: "market-investment-horizon", Buckets: []DistributionBucket{
				{OptionID: "ktlo", Fraction: 0.63},
				{OptionID: "tam_expansion", Fraction: 0.05}, // unchanged — should not appear
			}},
		},
	}

	deltas := ComputeDeltas(previous, current)
	if len(deltas.DistributionShifts) != 1 {
		t.Fatalf("expected 1 shift (unchanged tam_expansion excluded), got %+v", deltas.DistributionShifts)
	}
	shift := deltas.DistributionShifts[0]
	if shift.OptionID != "ktlo" || shift.PreviousFraction != 0.71 || shift.CurrentFraction != 0.63 {
		t.Errorf("shift = %+v", shift)
	}
}

func TestComputeDeltasNewDistributionOption(t *testing.T) {
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	previous := ReportDataset{GeneratedAt: t1, Distributions: []DimensionDistribution{
		{DimensionID: "kano", Buckets: []DistributionBucket{{OptionID: "must_be", Fraction: 1.0}}},
	}}
	current := ReportDataset{GeneratedAt: t2, Distributions: []DimensionDistribution{
		{DimensionID: "kano", Buckets: []DistributionBucket{
			{OptionID: "must_be", Fraction: 0.8},
			{OptionID: "attractive", Fraction: 0.2}, // new this cycle, previous fraction implied 0
		}},
	}}
	deltas := ComputeDeltas(previous, current)
	var attractiveShift *DistributionShift
	for i := range deltas.DistributionShifts {
		if deltas.DistributionShifts[i].OptionID == "attractive" {
			attractiveShift = &deltas.DistributionShifts[i]
		}
	}
	if attractiveShift == nil || attractiveShift.PreviousFraction != 0 || attractiveShift.CurrentFraction != 0.2 {
		t.Errorf("attractive shift = %+v, want previous=0 current=0.2", attractiveShift)
	}
}
