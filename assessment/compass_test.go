package assessment

import (
	"strings"
	"testing"
	"time"

	"github.com/ProductBuildersHQ/compass-rice/catalog"
	"github.com/ProductBuildersHQ/compass-rice/rice"
	"github.com/plexusone/structured-evaluation/claims"
)

const validCustomerB2BEvidence = `{
	"profileId": "customer/b2b/v1",
	"evidence": {
		"eligibleAccounts": 40,
		"affectedAccounts": 12,
		"eligibleArr": 10000000,
		"affectedArr": 3500000,
		"expectedRetentionOrExpansionImprovementPp": 2,
		"verifiedQuantitativeSources": 2,
		"verifiedQualitativeSources": 1,
		"effortPd": 20
	}
}`

func mustNormalize(t *testing.T) rice.Normalized {
	t.Helper()
	n, err := catalog.NormalizeDocument([]byte(validCustomerB2BEvidence))
	if err != nil {
		t.Fatalf("NormalizeDocument() error = %v", err)
	}
	return n
}

func TestCompassAssessmentValidate(t *testing.T) {
	n := mustNormalize(t)
	c := CompassAssessment{ProfileID: n.ProfileID, Normalized: n}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestCompassAssessmentValidateProfileIDMismatch(t *testing.T) {
	n := mustNormalize(t)
	c := CompassAssessment{ProfileID: "operations/v1", Normalized: n}
	if err := c.Validate(); err == nil {
		t.Error("Validate() with mismatched profileId = nil error, want error")
	}
}

func TestCompassAssessmentValidateBadProfileID(t *testing.T) {
	c := CompassAssessment{ProfileID: "not-a-valid-id"}
	if err := c.Validate(); err == nil {
		t.Error("Validate() with invalid profileId = nil error, want error")
	}
}

func TestCompassAssessmentValidateInvalidNormalized(t *testing.T) {
	c := CompassAssessment{ProfileID: "operations/v1", Normalized: rice.Normalized{ProfileID: "operations/v1"}}
	if err := c.Validate(); err == nil {
		t.Error("Validate() with invalid Normalized = nil error, want error")
	}
}

func TestCompassAssessmentValidateHumanReviewRequiresFields(t *testing.T) {
	n := mustNormalize(t)
	c := CompassAssessment{ProfileID: n.ProfileID, Normalized: n, HumanReview: &CompassHumanReview{}}
	if err := c.Validate(); err == nil {
		t.Error("Validate() with empty HumanReview = nil error, want error")
	}
	c.HumanReview = &CompassHumanReview{ReviewedBy: "pm@example.com", ReviewedAt: time.Now()}
	if err := c.Validate(); err != nil {
		t.Errorf("Validate() with complete HumanReview error = %v, want nil", err)
	}
}

func TestOpportunityAssessmentHasEvidenceFromCompassClaims(t *testing.T) {
	n := mustNormalize(t)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OS-001"}, "Test opportunity", time.Now())
	if a.HasEvidence() {
		t.Fatal("HasEvidence() = true before any evidence recorded, want false")
	}
	a.Compass = &CompassAssessment{
		ProfileID:  n.ProfileID,
		Normalized: n,
		Claims:     []*claims.Claim{{ID: "claim-1", Text: "affected ARR is $3.5M"}},
	}
	if !a.HasEvidence() {
		t.Error("HasEvidence() = false with a Compass claim recorded, want true")
	}
}

func TestResolveCompassRICENil(t *testing.T) {
	result := ResolveCompassRICE(nil)
	if result.Computable {
		t.Error("Computable = true for nil CompassAssessment, want false")
	}
	if result.Reason == "" {
		t.Error("Reason is empty for nil CompassAssessment, want a reason")
	}
}

func TestResolveCompassRICENeedsHumanReviewBlocksScoring(t *testing.T) {
	n := mustNormalize(t)
	c := &CompassAssessment{ProfileID: n.ProfileID, Normalized: n, NeedsHumanReview: true}
	result := ResolveCompassRICE(c)
	if result.Computable {
		t.Error("Computable = true for NeedsHumanReview assessment with no HumanReview, want false")
	}
	if !strings.Contains(result.Reason, "human review") {
		t.Errorf("Reason = %q, want it to mention human review", result.Reason)
	}

	c.HumanReview = &CompassHumanReview{ReviewedBy: "pm@example.com", ReviewedAt: time.Now()}
	result = ResolveCompassRICE(c)
	if !result.Computable {
		t.Errorf("Computable = false after HumanReview set, want true (reason: %s)", result.Reason)
	}
}

func TestResolveCompassRICEValid(t *testing.T) {
	n := mustNormalize(t)
	c := &CompassAssessment{ProfileID: n.ProfileID, Normalized: n}
	result := ResolveCompassRICE(c)
	if !result.Computable {
		t.Fatalf("Computable = false, want true (reason: %s)", result.Reason)
	}
	wantScore, err := n.Score()
	if err != nil {
		t.Fatalf("Score() error = %v", err)
	}
	if result.Score != wantScore {
		t.Errorf("Score = %v, want %v", result.Score, wantScore)
	}
	if result.ProfileID != string(n.ProfileID) {
		t.Errorf("ProfileID = %q, want %q", result.ProfileID, n.ProfileID)
	}
}

func TestToRankInputPrefersCompassOverLegacyRICE(t *testing.T) {
	n := mustNormalize(t)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OS-001"}, "Test opportunity", time.Now())
	a.Compass = &CompassAssessment{ProfileID: n.ProfileID, Normalized: n}
	a.RICE = &RICEAssessment{
		Reach:  Reach{Fraction: 0.9, EvidenceIDs: []string{"ev-1"}},
		Effort: EffortEstimate{Expected: 1},
	}

	in := a.ToRankInput()
	if !in.RICE.Computable {
		t.Fatalf("RICE.Computable = false, want true (reason: %s)", in.RICE.Reason)
	}
	if in.RICE.ProfileID != string(n.ProfileID) {
		t.Errorf("RICE.ProfileID = %q, want %q (legacy RICE path taken instead of Compass)", in.RICE.ProfileID, n.ProfileID)
	}
}

func TestToRankInputFallsBackToLegacyRICEWithoutCompass(t *testing.T) {
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OS-001"}, "Test opportunity", time.Now())
	in := a.ToRankInput()
	if in.RICE.Computable {
		t.Error("RICE.Computable = true with no RICE or Compass assessment, want false")
	}
	if in.RICE.Reason != "no RICE assessment recorded" {
		t.Errorf("RICE.Reason = %q, want %q", in.RICE.Reason, "no RICE assessment recorded")
	}
}

func TestOpportunityAssessmentEvidenceReferencesFromCompass(t *testing.T) {
	n := mustNormalize(t)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OS-001"}, "Test opportunity", time.Now())
	a.Compass = &CompassAssessment{
		ProfileID:  n.ProfileID,
		Normalized: n,
		Claims: []*claims.Claim{
			{ID: "claim-1", Text: "affected ARR is $3.5M"},
			{ID: "claim-2", Text: "12 accounts affected"},
			nil,
			{ID: "", Text: "should be skipped, no ID"},
		},
	}
	refs := a.EvidenceReferences()
	if len(refs) != 2 {
		t.Fatalf("EvidenceReferences() len = %d, want 2 (got %+v)", len(refs), refs)
	}
	for _, r := range refs {
		if r.AssessmentID != "OA-001" {
			t.Errorf("ref.AssessmentID = %q, want OA-001", r.AssessmentID)
		}
		if !strings.HasPrefix(r.QuestionID, "compass.customer/b2b/v1") {
			t.Errorf("ref.QuestionID = %q, want prefix compass.customer/b2b/v1", r.QuestionID)
		}
	}
}
