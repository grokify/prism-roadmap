package assessment

import (
	"testing"
	"time"

	"github.com/ProductBuildersHQ/compass-rice/catalog"
	"github.com/ProductBuildersHQ/compass-rice/rice"
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
