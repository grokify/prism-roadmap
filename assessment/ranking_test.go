package assessment

import (
	"testing"

	"github.com/grokify/prism-roadmap/prioritization"
)

func computable(score float64) RICEScoreResult {
	return RICEScoreResult{Score: score, Computable: true}
}

func TestRankingPolicyMoSCoWTierBeatsRICE(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-A", MoSCoW: prioritization.MoSCoWCouldHave, RICE: computable(0.90)},
		{AssessmentID: "OA-B", MoSCoW: prioritization.MoSCoWMustHave, RICE: computable(0.10)},
	}
	ranked := policy.Rank(inputs)

	if len(ranked) != 2 {
		t.Fatalf("expected 2 ranked opportunities, got %d", len(ranked))
	}
	// A Could with RICE 0.90 must NOT outrank a Must with RICE 0.10.
	if ranked[0].AssessmentID != "OA-B" || ranked[0].CalculatedRank != 1 {
		t.Errorf("rank 1 = %+v, want OA-B (Must beats Could regardless of RICE)", ranked[0])
	}
	if ranked[1].AssessmentID != "OA-A" || ranked[1].CalculatedRank != 2 {
		t.Errorf("rank 2 = %+v, want OA-A", ranked[1])
	}
}

func TestRankingPolicyRICEOrdersWithinTier(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-LOW", MoSCoW: prioritization.MoSCoWShouldHave, RICE: computable(0.10)},
		{AssessmentID: "OA-HIGH", MoSCoW: prioritization.MoSCoWShouldHave, RICE: computable(0.50)},
	}
	ranked := policy.Rank(inputs)
	if ranked[0].AssessmentID != "OA-HIGH" || ranked[1].AssessmentID != "OA-LOW" {
		t.Errorf("expected OA-HIGH before OA-LOW within the same tier, got %v then %v", ranked[0].AssessmentID, ranked[1].AssessmentID)
	}
}

func TestRankingPolicyExcludesWont(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-MUST", MoSCoW: prioritization.MoSCoWMustHave, RICE: computable(0.10)},
		{AssessmentID: "OA-WONT", MoSCoW: prioritization.MoSCoWWontHave, RICE: computable(0.99)},
		{AssessmentID: "OA-UNSET", MoSCoW: prioritization.MoSCoWUnspecified, RICE: computable(0.99)},
	}
	ranked := policy.Rank(inputs)
	if len(ranked) != 3 {
		t.Fatalf("expected all 3 inputs to appear in output (excluded, not dropped), got %d", len(ranked))
	}

	byID := make(map[string]RankedOpportunity)
	for _, r := range ranked {
		byID[r.AssessmentID] = r
	}
	if byID["OA-WONT"].Excluded != ExclusionWont || byID["OA-WONT"].CalculatedRank != 0 {
		t.Errorf("OA-WONT = %+v, want Excluded=wont, CalculatedRank=0", byID["OA-WONT"])
	}
	if byID["OA-UNSET"].Excluded != ExclusionWont {
		t.Errorf("OA-UNSET = %+v, want Excluded=wont (unspecified MoSCoW is not rankable)", byID["OA-UNSET"])
	}
	if byID["OA-MUST"].Excluded != "" || byID["OA-MUST"].CalculatedRank != 1 {
		t.Errorf("OA-MUST = %+v, want Excluded=\"\", CalculatedRank=1", byID["OA-MUST"])
	}
}

func TestRankingPolicyExcludesUncomputableRICE(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-NOEV", MoSCoW: prioritization.MoSCoWMustHave, RICE: RICEScoreResult{Computable: false, Reason: "no evidence"}},
	}
	ranked := policy.Rank(inputs)
	if len(ranked) != 1 || ranked[0].Excluded != ExclusionRICEUncomputable {
		t.Errorf("ranked = %+v, want 1 item excluded as rice_uncomputable", ranked)
	}
	if ranked[0].CalculatedRank != 0 {
		t.Errorf("CalculatedRank = %d, want 0 for an excluded (not silently zero-scored) item", ranked[0].CalculatedRank)
	}
}

func TestRankingPolicyTieDetection(t *testing.T) {
	policy := DefaultRankingPolicy() // 5% band
	inputs := []RankInput{
		{AssessmentID: "OA-1", MoSCoW: prioritization.MoSCoWShouldHave, RICE: computable(0.310)},
		{AssessmentID: "OA-2", MoSCoW: prioritization.MoSCoWShouldHave, RICE: computable(0.314)}, // within 5% of 0.310
		{AssessmentID: "OA-3", MoSCoW: prioritization.MoSCoWShouldHave, RICE: computable(0.100)}, // not tied with either
	}
	ranked := policy.Rank(inputs)

	byID := make(map[string]RankedOpportunity)
	for _, r := range ranked {
		byID[r.AssessmentID] = r
	}

	if len(byID["OA-2"].TiedWith) != 1 || byID["OA-2"].TiedWith[0] != "OA-1" {
		t.Errorf("OA-2.TiedWith = %v, want [OA-1]", byID["OA-2"].TiedWith)
	}
	if len(byID["OA-1"].TiedWith) != 1 || byID["OA-1"].TiedWith[0] != "OA-2" {
		t.Errorf("OA-1.TiedWith = %v, want [OA-2] (tie relation is symmetric)", byID["OA-1"].TiedWith)
	}
	if len(byID["OA-3"].TiedWith) != 0 {
		t.Errorf("OA-3.TiedWith = %v, want empty (not within band of either)", byID["OA-3"].TiedWith)
	}
}

func TestRankingPolicyTiesDoNotCrossTiers(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-MUST", MoSCoW: prioritization.MoSCoWMustHave, RICE: computable(0.100)},
		{AssessmentID: "OA-SHOULD", MoSCoW: prioritization.MoSCoWShouldHave, RICE: computable(0.101)}, // numerically tied but different tier
	}
	ranked := policy.Rank(inputs)
	for _, r := range ranked {
		if len(r.TiedWith) != 0 {
			t.Errorf("%s.TiedWith = %v, want empty — ties must not cross MoSCoW tiers", r.AssessmentID, r.TiedWith)
		}
	}
}

func TestRankingPolicyValidate(t *testing.T) {
	if err := DefaultRankingPolicy().Validate(); err != nil {
		t.Errorf("unexpected error for default policy: %v", err)
	}
	if err := (RankingPolicy{Version: "1.0", TieBandFraction: 0.05}).Validate(); err == nil {
		t.Error("expected error for missing ID")
	}
	if err := (RankingPolicy{ID: "x", TieBandFraction: -1}).Validate(); err == nil {
		t.Error("expected error for negative tie band")
	}
}

func TestApplyOverrides(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-1", MoSCoW: prioritization.MoSCoWMustHave, RICE: computable(0.50)},
		{AssessmentID: "OA-2", MoSCoW: prioritization.MoSCoWMustHave, RICE: computable(0.40)},
		{AssessmentID: "OA-3", MoSCoW: prioritization.MoSCoWWontHave, RICE: computable(0.99)},
	}
	calculated := policy.Rank(inputs)

	overrides := []RankOverride{
		{AssessmentID: "OA-2", FinalRank: 1, Rationale: "strategic TAM expansion priority", ApprovedBy: "vp-product"},
		{AssessmentID: "OA-1", FinalRank: 2, Rationale: "deprioritized pending budget", ApprovedBy: "vp-product"},
	}
	final := ApplyOverrides(calculated, overrides)

	if final[0].AssessmentID != "OA-2" || final[0].FinalRank != 1 {
		t.Errorf("final[0] = %+v, want OA-2 at rank 1", final[0])
	}
	if final[0].Override == nil || final[0].Override.Rationale != "strategic TAM expansion priority" {
		t.Errorf("final[0].Override = %+v, want the OA-2 override recorded", final[0].Override)
	}
	if final[1].AssessmentID != "OA-1" || final[1].FinalRank != 2 {
		t.Errorf("final[1] = %+v, want OA-1 at rank 2", final[1])
	}

	// Excluded items must survive unchanged and cannot be overridden.
	var wont *OpportunityRank
	for i := range final {
		if final[i].AssessmentID == "OA-3" {
			wont = &final[i]
		}
	}
	if wont == nil {
		t.Fatal("OA-3 (excluded) missing from ApplyOverrides output")
	}
	if wont.Excluded != ExclusionWont || wont.FinalRank != 0 || wont.Override != nil {
		t.Errorf("OA-3 = %+v, want unchanged excluded state with no override applied", *wont)
	}
}

func TestApplyOverridesWithoutOverrideKeepsCalculatedRank(t *testing.T) {
	policy := DefaultRankingPolicy()
	inputs := []RankInput{
		{AssessmentID: "OA-1", MoSCoW: prioritization.MoSCoWMustHave, RICE: computable(0.50)},
	}
	final := ApplyOverrides(policy.Rank(inputs), nil)
	if final[0].FinalRank != final[0].CalculatedRank || final[0].Override != nil {
		t.Errorf("expected FinalRank == CalculatedRank and no Override when none supplied, got %+v", final[0])
	}
}

func TestRankOverrideValidate(t *testing.T) {
	valid := RankOverride{AssessmentID: "OA-1", FinalRank: 1, Rationale: "reason", ApprovedBy: "someone"}
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := (RankOverride{Rationale: "r", ApprovedBy: "a"}).Validate(); err == nil {
		t.Error("expected error for missing AssessmentID")
	}
	if err := (RankOverride{AssessmentID: "OA-1", ApprovedBy: "a"}).Validate(); err == nil {
		t.Error("expected error for missing Rationale")
	}
	if err := (RankOverride{AssessmentID: "OA-1", Rationale: "r"}).Validate(); err == nil {
		t.Error("expected error for missing ApprovedBy")
	}
}

func TestRankCollisions(t *testing.T) {
	ranks := []OpportunityRank{
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1"}, FinalRank: 1},
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-2"}, FinalRank: 1}, // collision
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-3"}, FinalRank: 2},
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-4", Excluded: ExclusionWont}, FinalRank: 0},
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-5", Excluded: ExclusionWont}, FinalRank: 0}, // excluded, must not count as a collision
	}
	collisions := RankCollisions(ranks)
	if len(collisions) != 1 || collisions[0] != 1 {
		t.Errorf("RankCollisions() = %v, want [1]", collisions)
	}
}

func TestRankCollisionsNoneWhenUnique(t *testing.T) {
	ranks := []OpportunityRank{
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-1"}, FinalRank: 1},
		{RankedOpportunity: RankedOpportunity{AssessmentID: "OA-2"}, FinalRank: 2},
	}
	if got := RankCollisions(ranks); len(got) != 0 {
		t.Errorf("RankCollisions() = %v, want empty", got)
	}
}
