package assessment

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ProductBuildersHQ/compass-rice/rice"
	"github.com/plexusone/structured-evaluation/claims"
	"github.com/plexusone/structured-evaluation/rubric"
)

// CompassAssessment is one cycle's COMPASS-RICE judgment: the raw evidence a
// judge (an LLM or a human) produced for one compass-rice profile, the
// deterministic Normalized result that evidence normalizes to, the judge's
// reasoning and verified claims, and whether a human has cleared it for
// scoring. compass-rice's own architecture ("compass-rice defines the
// methodology; the consumer executes it") makes EvidenceJSON the source of
// truth — Normalized is always re-derivable from it via the profile's
// Normalizer, never hand-edited independently.
type CompassAssessment struct {
	// ProfileID pins the profile, evidence model, and normalization
	// algorithm version this assessment was judged against, e.g.
	// "customer/b2b/v1". Must match Normalized.ProfileID.
	ProfileID rice.ProfileID `json:"profileId"`

	// EvidenceJSON is the profile-typed evidence document the judge
	// produced — the source of truth, re-normalizable at any time via the
	// profile's Normalizer (see the compass-rice catalog package).
	EvidenceJSON json.RawMessage `json:"evidenceJson"`

	// Normalized is the deterministic evidence -> Reach/Impact/Confidence/
	// Effort result compass-rice's Normalizer produced from EvidenceJSON.
	Normalized rice.Normalized `json:"normalized"`

	// Categories are the judge's reasoning per rubric category ("reach",
	// "impact"), with citations — compass-rice judge.Output's own
	// provenance trail, carried through unchanged.
	Categories []rubric.CategoryResult `json:"categories,omitempty"`

	// Claims backs the countable fields in the evidence; compass-rice's
	// provenance.Confidence derives Normalized.Confidence from this list,
	// never from a judge's self-reported confidence.
	Claims []*claims.Claim `json:"claims,omitempty"`

	// NeedsHumanReview is captured at ingest from compass-rice
	// judge.Output.NeedsHumanReview() — a category fell below its
	// confidence threshold, or carries a reason code requiring human
	// review. ResolveCompassRICE treats this assessment as uncomputable
	// until HumanReview is set, even though Normalized already holds a
	// valid score.
	NeedsHumanReview bool `json:"needsHumanReview,omitempty"`

	// HumanReview is set when a PM accepts an assessment that was flagged
	// NeedsHumanReview. Acceptance is recorded on a new assessment cycle
	// (never a mutation of this one), matching OpportunityAssessment's
	// existing cycle-immutability discipline.
	HumanReview *CompassHumanReview `json:"humanReview,omitempty"`
}

// CompassHumanReview records a PM's acceptance of a COMPASS-RICE assessment
// that was flagged for human review before it could be trusted for scoring.
type CompassHumanReview struct {
	ReviewedBy string    `json:"reviewedBy"`
	ReviewedAt time.Time `json:"reviewedAt"`
	Note       string    `json:"note,omitempty"`
}

// Validate returns an error if c is not internally consistent: ProfileID
// parses, Normalized passes its own validation, and Normalized.ProfileID
// matches ProfileID — a stored assessment can never silently carry a score
// for a different profile than the one it claims.
func (c CompassAssessment) Validate() error {
	if _, _, _, err := c.ProfileID.Parse(); err != nil {
		return fmt.Errorf("compass: profileId: %w", err)
	}
	if err := c.Normalized.Validate(); err != nil {
		return fmt.Errorf("compass: normalized: %w", err)
	}
	if c.Normalized.ProfileID != c.ProfileID {
		return fmt.Errorf("compass: normalized.profileId %q does not match profileId %q", c.Normalized.ProfileID, c.ProfileID)
	}
	if c.HumanReview != nil {
		if c.HumanReview.ReviewedBy == "" {
			return fmt.Errorf("compass: humanReview.reviewedBy is required")
		}
		if c.HumanReview.ReviewedAt.IsZero() {
			return fmt.Errorf("compass: humanReview.reviewedAt is required")
		}
	}
	return nil
}

// ResolveCompassRICE is the single gate deciding whether a CompassAssessment
// is trusted for ranking. It never falls back to the legacy ladder RICE
// scale — compass-rice's 0-100 Reach band and the legacy 0..1 Reach fraction
// are different regimes, and mixing them in one ranked list would silently
// distort relative ordering (an opportunity scored on one scale would not
// really outrank one scored on the other). An opportunity with no COMPASS
// assessment, or one still awaiting human review, is uncomputable — not
// silently scored on a different methodology.
func ResolveCompassRICE(c *CompassAssessment) RICEScoreResult {
	if c == nil {
		return RICEScoreResult{Reason: "no COMPASS assessment recorded"}
	}
	if c.NeedsHumanReview && c.HumanReview == nil {
		return RICEScoreResult{Reason: "COMPASS assessment awaiting human review", ProfileID: string(c.ProfileID)}
	}
	if err := c.Validate(); err != nil {
		return RICEScoreResult{Reason: err.Error(), ProfileID: string(c.ProfileID)}
	}
	score, err := c.Normalized.Score()
	if err != nil {
		return RICEScoreResult{Reason: err.Error(), ProfileID: string(c.ProfileID)}
	}
	return RICEScoreResult{
		Score:      score,
		Computable: true,
		ProfileID:  string(c.ProfileID),
	}
}
