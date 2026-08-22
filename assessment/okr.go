package assessment

import (
	"fmt"

	"github.com/grokify/prism-roadmap/goals/okr"
)

// ContributionStrength is how strongly an opportunity is expected to move a
// key result — categorical, on the same coarse scale as this package's
// other judgments, not a fabricated percentage.
type ContributionStrength string

const (
	ContributionHigh   ContributionStrength = "high"
	ContributionMedium ContributionStrength = "medium"
	ContributionLow    ContributionStrength = "low"
)

// IsValid reports whether s is a recognized contribution strength.
func (s ContributionStrength) IsValid() bool {
	switch s {
	case ContributionHigh, ContributionMedium, ContributionLow:
		return true
	default:
		return false
	}
}

// OKRContribution links an opportunity to one Objective — and optionally
// one specific KeyResult within it — from the goals/okr package, with an
// evidence-backed strength.
//
// Explicitly NOT a RankingPolicy input: OKR alignment answers "what are we
// trying to achieve," not "which investment is best toward it" — folding it
// into the ranking formula would double-count against RICE/MoSCoW and erase
// the intended hierarchy (ideation doc: "OKRs determine what matters; RICE
// helps determine which projects are the best investments toward what
// matters"). What it does enable: a portfolio rollup of investment by
// objective (e.g. "32% of Person-Days advance Objective 3" — an
// omniroadmap compile-step aggregation over many opportunities'
// contributions, out of scope for this type), and — once opportunities
// ship — comparing predicted KR movement against actual, closing the loop
// from prioritization decision to measured outcome.
type OKRContribution struct {
	ObjectiveID string `json:"objectiveId"`

	// KeyResultID is optional — an opportunity can contribute to an
	// objective's overall intent without being tied to one specific KR.
	KeyResultID string `json:"keyResultId,omitempty"`

	Strength ContributionStrength `json:"strength"`

	Rationale   string   `json:"rationale,omitempty"`
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// Validate returns an error if required fields are missing or invalid.
// This checks shape only; use ValidateAgainst to additionally confirm the
// referenced Objective/KeyResult actually exist in a given OKRSet.
func (c OKRContribution) Validate() error {
	if c.ObjectiveID == "" {
		return fmt.Errorf("objectiveId is required")
	}
	if !c.Strength.IsValid() {
		return fmt.Errorf("invalid strength: %q", c.Strength)
	}
	return nil
}

// ValidateAgainst checks that this contribution's ObjectiveID (and
// KeyResultID, if set) refer to entries that actually exist in set —
// catching a typo'd or stale OKR reference that Validate's shape check
// alone cannot.
func (c OKRContribution) ValidateAgainst(set *okr.OKRSet) error {
	if err := c.Validate(); err != nil {
		return err
	}
	if set == nil {
		return fmt.Errorf("okr set is required")
	}
	for _, obj := range set.ToObjectives() {
		if obj.ID != c.ObjectiveID {
			continue
		}
		if c.KeyResultID == "" {
			return nil
		}
		for _, kr := range obj.KeyResults {
			if kr.ID == c.KeyResultID {
				return nil
			}
		}
		return fmt.Errorf("keyResultId %q not found under objective %q", c.KeyResultID, c.ObjectiveID)
	}
	return fmt.Errorf("objectiveId %q not found", c.ObjectiveID)
}
