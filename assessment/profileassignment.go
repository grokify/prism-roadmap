package assessment

import (
	"fmt"
	"strings"
	"time"

	"github.com/ProductBuildersHQ/compass-rice/rice"
)

// ProfileAssignmentStatus tracks a ProfileAssignment through the two-phase
// selection workflow: an LLM judge proposes a profile, a PM confirms (or
// rejects) it. Only a confirmed assignment's ProfileID may be trusted by a
// consumer to gate scoring — compass-rice PRD D5's "exactly one profile
// generates the canonical RICE score" is a PM decision, not something an
// LLM's proposal alone can authorize.
type ProfileAssignmentStatus string

const (
	ProfileAssignmentProposed  ProfileAssignmentStatus = "proposed"
	ProfileAssignmentConfirmed ProfileAssignmentStatus = "confirmed"
	ProfileAssignmentRejected  ProfileAssignmentStatus = "rejected"
)

// IsValid reports whether s is one of the three recognized statuses.
func (s ProfileAssignmentStatus) IsValid() bool {
	switch s {
	case ProfileAssignmentProposed, ProfileAssignmentConfirmed, ProfileAssignmentRejected:
		return true
	default:
		return false
	}
}

// ParseProfileAssignmentStatus parses a string into a ProfileAssignmentStatus
// (case-insensitive). Returns an error if the value is not recognized.
func ParseProfileAssignmentStatus(s string) (ProfileAssignmentStatus, error) {
	status := ProfileAssignmentStatus(strings.ToLower(strings.TrimSpace(s)))
	if !status.IsValid() {
		return "", fmt.Errorf("invalid profile assignment status: %q", s)
	}
	return status, nil
}

// ProfileAssignment records an opportunity's primary COMPASS-RICE investment
// thesis — the profile whose Normalizer produces its canonical score
// (compass-rice PRD D5: "never run competing profiles to produce competing
// scores"). It is spec-scoped and survives assessment cycles, like
// RankOverride: a re-assessment does not require re-selecting the profile,
// and a profile correction is recorded by moving Status back to proposed
// with a new Rationale, not by mutating history.
type ProfileAssignment struct {
	// SpecID is the canvas.OpportunitySpec.Metadata.ID this assignment
	// applies to.
	SpecID string `json:"specId"`

	// ProfileID is the primary investment thesis — the profile whose
	// Normalizer produces this opportunity's canonical COMPASS-RICE score.
	// Only trustworthy for scoring once Status is confirmed (see
	// omniroadmap's compile-time gating, which checks this against
	// CompassAssessment.ProfileID).
	ProfileID rice.ProfileID `json:"profileId"`

	// Secondary lists other profiles with real value for this opportunity
	// ("this is primarily Platform, but also has Customer and Risk value")
	// — context only, per compass-rice's integration guide: never run
	// through their own Normalizer to produce competing scores.
	Secondary []rice.Profile `json:"secondary,omitempty"`

	// Rationale explains why ProfileID is the primary thesis.
	Rationale string `json:"rationale"`

	// ProposedBy identifies who or what proposed this assignment (a judge
	// identity, e.g. a model name/session, or a human).
	ProposedBy string `json:"proposedBy"`

	Status ProfileAssignmentStatus `json:"status"`

	// ConfirmedBy/ConfirmedAt are set when a PM confirms (or rejects) the
	// proposal. Required whenever Status is confirmed or rejected.
	ConfirmedBy string    `json:"confirmedBy,omitempty"`
	ConfirmedAt time.Time `json:"confirmedAt,omitempty"`

	// EvidenceIDs optionally cite the evidence motivating the profile
	// choice, auditable like any other assignment input.
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// ProposeProfileAssignment creates a new proposed ProfileAssignment — the
// judge phase of the two-phase workflow. A PM must separately call Confirm
// or Reject before ProfileID is trustworthy for scoring.
func ProposeProfileAssignment(specID string, profileID rice.ProfileID, rationale, proposedBy string) ProfileAssignment {
	return ProfileAssignment{
		SpecID:     specID,
		ProfileID:  profileID,
		Rationale:  rationale,
		ProposedBy: proposedBy,
		Status:     ProfileAssignmentProposed,
	}
}

// Confirm returns a copy of p with Status set to confirmed and the PM's
// identity/timestamp recorded — the human phase of the two-phase workflow.
func (p ProfileAssignment) Confirm(confirmedBy string, confirmedAt time.Time) ProfileAssignment {
	confirmed := p
	confirmed.Status = ProfileAssignmentConfirmed
	confirmed.ConfirmedBy = confirmedBy
	confirmed.ConfirmedAt = confirmedAt
	return confirmed
}

// Reject returns a copy of p with Status set to rejected and the PM's
// identity/timestamp recorded, along with a Rationale explaining why —
// e.g. the judge proposed the wrong profile and a new proposal is needed.
func (p ProfileAssignment) Reject(confirmedBy string, confirmedAt time.Time, rationale string) ProfileAssignment {
	rejected := p
	rejected.Status = ProfileAssignmentRejected
	rejected.ConfirmedBy = confirmedBy
	rejected.ConfirmedAt = confirmedAt
	rejected.Rationale = rationale
	return rejected
}

// Validate returns an error if p is not well-formed for its Status.
func (p ProfileAssignment) Validate() error {
	if p.SpecID == "" {
		return fmt.Errorf("specId is required")
	}
	if _, _, _, err := p.ProfileID.Parse(); err != nil {
		return fmt.Errorf("profileId: %w", err)
	}
	for _, s := range p.Secondary {
		if !s.Valid() {
			return fmt.Errorf("secondary: invalid profile %q", s)
		}
	}
	if p.Rationale == "" {
		return fmt.Errorf("rationale is required")
	}
	if p.ProposedBy == "" {
		return fmt.Errorf("proposedBy is required")
	}
	if !p.Status.IsValid() {
		return fmt.Errorf("invalid status %q", p.Status)
	}
	if p.Status == ProfileAssignmentConfirmed || p.Status == ProfileAssignmentRejected {
		if p.ConfirmedBy == "" {
			return fmt.Errorf("confirmedBy is required for status %q", p.Status)
		}
		if p.ConfirmedAt.IsZero() {
			return fmt.Errorf("confirmedAt is required for status %q", p.Status)
		}
	}
	return nil
}
