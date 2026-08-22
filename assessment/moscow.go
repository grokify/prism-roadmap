package assessment

import "github.com/grokify/prism-roadmap/prioritization"

// MoSCoWLadder is the MoSCoW priority-tier classification ladder: an
// objective, evidence-backed alternative to informal MoSCoW, where "everyone
// argues their item is a Must" (ideation doc). Must means a binding
// obligation, not merely "very important." Level IDs match values accepted
// by prioritization.ParseMoSCoWPriority.
//
// Won't/Not Now is intentionally absent as a ladder rung: it means
// "insufficient current priority," not a criterion to test for, so it is
// the ladder's floor (see ResolveMoSCoWPriority) rather than a rung with its
// own criteria.
func MoSCoWLadder() *Ladder {
	return &Ladder{
		ID:   "moscow",
		Name: "MoSCoW",
		Levels: []ThresholdLevel{
			{
				ID: "must", Label: "Must",
				Criteria: []string{
					"Required to keep the service operational, supported, secure, or maintainable (KTLO)",
					"Required to satisfy a binding regulatory, security, privacy, or compliance obligation",
					"Required by an existing customer/vendor contractual commitment",
					"Not doing this leaves an unacceptable availability, security, data-loss, or operational risk",
					"Required because a critical technology or dependency is approaching end-of-life/end-of-support",
				},
			},
			{
				ID: "should", Label: "Should",
				Criteria: []string{
					"Would provide meaningful value to a substantial portion of target customers",
					"A material capability commonly offered by relevant competitors",
					"Relevant analyst frameworks/reports identify this capability as important or expected",
					"Consistent demand from multiple customers/prospects",
					"Materially advances an approved product/company objective",
				},
			},
			{
				ID: "could", Label: "Could",
				Criteria: []string{
					"Demand is limited to a relatively small customer/prospect segment",
					"Provides value without materially affecting competitive position",
					"Unlikely to materially affect analyst evaluations or category expectations",
					"Customers can achieve their goals reasonably well without it",
				},
			},
		},
	}
}

// ResolveMoSCoWPriority evaluates MoSCoWLadder() against the given answers
// and returns the resulting prioritization.MoSCoWPriority. When no level is
// satisfied with evidence, returns MoSCoWWontHave — MoSCoW's "Won't/Not Now"
// tier is the ladder's floor, not a criterion the judge tests for.
func ResolveMoSCoWPriority(answers []ThresholdAnswer) prioritization.MoSCoWPriority {
	lvl, _, ok := MoSCoWLadder().Evaluate(answers)
	if !ok {
		return prioritization.MoSCoWWontHave
	}
	priority, err := prioritization.ParseMoSCoWPriority(lvl.ID)
	if err != nil {
		return prioritization.MoSCoWWontHave
	}
	return priority
}
