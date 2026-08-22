package assessment

// KanoDimension is the built-in Kano portfolio dimension: which kind of
// customer satisfaction/value an opportunity creates (Must-be, Performance,
// Attractive, Indifferent, Reverse). This definition exists for
// display/reference (its Options carry no per-option Questions) — Kano
// classification itself is resolved by ResolveKano, not
// DimensionDefinition.ResolveCategory, because Kano's categories are
// derived by pattern-matching across shared cross-cutting characteristics
// (ideation doc), not independent per-option criteria the way Market
// Investment Horizon's are (see MarketInvestmentHorizonDimension).
//
// Kano is portfolio-descriptive only — it never enters Opportunity Rank
// (RankingPolicy.Rank uses only MoSCoW + RICE); it explains the resulting
// portfolio's composition and can inform a human tie-break.
func KanoDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "kano", Name: "Kano", Version: "1.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{ID: "must_be", Label: "Must-be"},
			{ID: "performance", Label: "Performance"},
			{ID: "attractive", Label: "Attractive"},
			{ID: "indifferent", Label: "Indifferent"},
			{ID: "reverse", Label: "Reverse"},
		},
	}
}

// KanoCharacteristic is a judge's evidence-backed answer to one of Kano's
// eight cross-cutting characterization questions.
type KanoCharacteristic struct {
	Answer      bool     `json:"answer"`
	Rationale   string   `json:"rationale,omitempty"`
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// isYes reports whether this characteristic was answered true AND
// evidence-backed (the same evidence discipline used throughout this
// package: an unsupported true claim does not count).
func (c KanoCharacteristic) isYes() bool {
	return c.Answer && len(c.EvidenceIDs) > 0
}

// KanoAnswers captures a judge's answers to Kano's eight characterization
// questions (ideation doc):
//
//   - Expected: would customers reasonably consider this a basic
//     expectation for a platform of this type?
//   - AbsenceDissatisfaction: would absence or poor implementation
//     materially reduce customer satisfaction?
//   - PresenceSatisfaction: would providing/improving this increase
//     customer satisfaction?
//   - MoreIsBetter: would progressively better performance/capability
//     produce progressively greater satisfaction?
//   - Unexpected: would customers generally NOT expect this as baseline?
//   - Delight: could this materially delight customers or differentiate
//     the product?
//   - Indifferent: would most affected customers care little whether this
//     exists or improves?
//   - ReversePreference: would a meaningful portion of target customers
//     prefer this not exist or not be enabled?
//
// A negative claim (Answer: false) needs no evidence to count as "no" —
// declining something doesn't require proof, mirroring Reach's
// zero-fraction case elsewhere in this package. Only "yes" claims require
// evidence.
type KanoAnswers struct {
	Expected               KanoCharacteristic `json:"expected"`
	AbsenceDissatisfaction KanoCharacteristic `json:"absenceDissatisfaction"`
	PresenceSatisfaction   KanoCharacteristic `json:"presenceSatisfaction"`
	MoreIsBetter           KanoCharacteristic `json:"moreIsBetter"`
	Unexpected             KanoCharacteristic `json:"unexpected"`
	Delight                KanoCharacteristic `json:"delight"`
	Indifferent            KanoCharacteristic `json:"indifferent"`
	ReversePreference      KanoCharacteristic `json:"reversePreference"`
}

// ResolveKano classifies an opportunity's Kano category from its
// characterization answers, following the documented decision patterns
// (ideation doc):
//
//	Must-be:      Expected=Y, AbsenceDissatisfaction=Y, PresenceSatisfaction=N
//	Performance:  AbsenceDissatisfaction=Y, PresenceSatisfaction=Y, MoreIsBetter=Y
//	Attractive:   Unexpected=Y, PresenceSatisfaction=Y, Delight=Y
//	Indifferent:  AbsenceDissatisfaction=N, PresenceSatisfaction=N, Indifferent=Y
//	Reverse:      PresenceSatisfaction=N, ReversePreference=Y
//
// Returns Ambiguous if the answer pattern matches more than one category
// (a sign the answers are internally inconsistent and need review) and an
// unresolved (zero-value) selection if too little was answered to match
// any pattern.
func ResolveKano(a KanoAnswers) CategorySelection {
	var matched []string

	if a.Expected.isYes() && a.AbsenceDissatisfaction.isYes() && !a.PresenceSatisfaction.isYes() {
		matched = append(matched, "must_be")
	}
	if a.AbsenceDissatisfaction.isYes() && a.PresenceSatisfaction.isYes() && a.MoreIsBetter.isYes() {
		matched = append(matched, "performance")
	}
	if a.Unexpected.isYes() && a.PresenceSatisfaction.isYes() && a.Delight.isYes() {
		matched = append(matched, "attractive")
	}
	if !a.AbsenceDissatisfaction.isYes() && !a.PresenceSatisfaction.isYes() && a.Indifferent.isYes() {
		matched = append(matched, "indifferent")
	}
	if !a.PresenceSatisfaction.isYes() && a.ReversePreference.isYes() {
		matched = append(matched, "reverse")
	}

	switch len(matched) {
	case 0:
		return CategorySelection{}
	case 1:
		return CategorySelection{OptionID: matched[0], Resolved: true}
	default:
		return CategorySelection{Ambiguous: true, AmbiguousOptionIDs: matched}
	}
}
