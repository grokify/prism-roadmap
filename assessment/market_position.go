package assessment

// MarketPositionDimension is the built-in Market Position (BCG Growth-Share
// Matrix) portfolio dimension: where a PRODUCT LINE sits on market growth ×
// relative share — Star, Cash Cow, Question Mark, Dog — plus Enabler for a
// horizontal platform whose value is leverage across lines rather than its own
// market share. Adapted from the Boston Consulting Group growth-share matrix;
// this is prism-roadmap's own framework, not a published external standard.
//
// GRANULARITY: unlike Kano and MIH (which classify an opportunity/initiative),
// BCG classifies a PRODUCT LINE. The DimensionDefinition here is the reusable
// rubric; the assignment is made once per product line (see
// omniroadmap-core provider.ProductLine / prismstudio ADR-0005) and inherited by
// that line's initiatives. It is portfolio-descriptive only — like MIH it never
// enters Opportunity Rank; it contextualizes it (invest Stars & Question Marks,
// harvest Cash Cows, divest Dogs, sustain Enablers for leverage).
//
// Unlike MIH, BCG is NOT ordinal — the four quadrants are a 2×2 categorical
// split (growth × share), so there is no rollup. Each quadrant's question folds
// both axes into one condition; if more than one matches, ResolveCategory
// reports Ambiguous rather than silently picking one.
func MarketPositionDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "market-position", Name: "Market Position (BCG Growth-Share)", Version: "1.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{
				ID: "star", Label: "Star",
				Questions: []DimensionQuestion{
					{
						ID:       "high-growth-high-share",
						Question: "Is this a high-growth market in which the product line holds a leading or strengthening relative share?",
					},
				},
			},
			{
				ID: "cash_cow", Label: "Cash Cow",
				Questions: []DimensionQuestion{
					{
						ID:       "low-growth-high-share",
						Question: "Is this a low-growth or mature market in which the product line holds a leading relative share (a profit/cash generator)?",
					},
				},
			},
			{
				ID: "question_mark", Label: "Question Mark",
				Questions: []DimensionQuestion{
					{
						ID:       "high-growth-low-share",
						Question: "Is this a high-growth market in which the product line holds a low or unproven relative share (a selective bet)?",
					},
				},
			},
			{
				ID: "dog", Label: "Dog",
				Questions: []DimensionQuestion{
					{
						ID:       "low-growth-low-share",
						Question: "Is this a low-growth market in which the product line holds a low relative share (limited strategic value)?",
					},
				},
			},
			{
				ID: "enabler", Label: "Enabler",
				Questions: []DimensionQuestion{
					{
						ID:       "horizontal-leverage",
						Question: "Is this a horizontal platform/infrastructure line whose value is leverage across other product lines rather than its own market share?",
					},
				},
			},
		},
	}
}

// BCG quadrant option IDs.
const (
	BCGStar         = "star"
	BCGCashCow      = "cash_cow"
	BCGQuestionMark = "question_mark"
	BCGDog          = "dog"
	BCGEnabler      = "enabler"
)
