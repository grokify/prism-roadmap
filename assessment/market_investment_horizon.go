package assessment

// MarketInvestmentHorizonDimension is the built-in Market Investment
// Horizon (MIH) portfolio dimension: where relative to the current market
// an investment creates or protects value, as a formal four-rung ordinal
// ladder — KTLO (protect the existing business) < SOM (capture demand
// obtainable now within customers/segments we already reach) < SAM (expand
// into serviceable segments our model can address but do not yet reach) <
// TAM Expansion (enter markets not effectively addressable today). Adapted
// from an industry investment-classification practice combining KTLO with
// TAM/SAM/SOM market concepts (ideation doc); this is prism-roadmap's own
// framework, not a published external standard.
//
// The four options are defined as INCREMENTAL rings (the ring, not the
// nested disk), so they stay mutually exclusive: an initiative is
// classified by the furthest-out horizon it opens. SOM ("win now within
// reach we already have") and SAM ("extend into reach we do not yet have")
// are deliberately worded to not overlap.
//
// A 3-way rollup — KTLO / SAM+SOM / TAM — is the canonical presentation
// grain (see Rollup helpers); the 4-way split is the storage grain. The
// rollup is also the natural fallback when SOM vs. SAM is genuinely
// ambiguous.
//
// Unlike Kano, MIH's categories each have their own independent judge
// criterion, so this dimension resolves through the generic
// DimensionDefinition.ResolveCategory rather than a bespoke resolver — no
// separate ResolveMIH function is needed. If more than one criterion is
// satisfied (e.g. an initiative plausibly reads as both KTLO and TAM
// Expansion), ResolveCategory reports that as Ambiguous rather than
// silently picking one; that ambiguity is itself useful portfolio
// information, not just a rubric defect.
//
// Like Kano, MIH is portfolio-descriptive only — it never enters
// Opportunity Rank.
// MIH option IDs, in ordinal order: KTLO < SOM < SAM < TAM Expansion.
const (
	MIHOptionKTLO = "ktlo"
	MIHOptionSOM  = "som"
	MIHOptionSAM  = "sam"
	MIHOptionTAM  = "tam_expansion"

	// MIHRollupSAMSOM is the 3-way rollup key SOM and SAM collapse to — the
	// same "sam_som" id the legacy combined bucket used, so pre-split data
	// and rolled-up 4-way views share one key.
	MIHRollupSAMSOM = "sam_som"
)

// MIHRollup maps a 4-way MIH option ID to its 3-way rollup key
// (KTLO / SAM+SOM / TAM Expansion). SOM and SAM — and the legacy combined
// "sam_som" — collapse to MIHRollupSAMSOM; KTLO and TAM pass through.
func MIHRollup(optionID string) string {
	switch optionID {
	case MIHOptionSOM, MIHOptionSAM, MIHRollupSAMSOM:
		return MIHRollupSAMSOM
	default:
		return optionID
	}
}

func MarketInvestmentHorizonDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "market-investment-horizon", Name: "Market Investment Horizon", Version: "2.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{
				ID: "ktlo", Label: "KTLO",
				Questions: []DimensionQuestion{
					{
						ID:       "sustains-existing-business",
						Question: "Is this initiative primarily required to sustain, secure, support, comply, or operate the existing business?",
					},
				},
			},
			{
				ID: "som", Label: "SOM",
				Questions: []DimensionQuestion{
					{
						ID:       "captures-obtainable-demand-now",
						Question: "Is this initiative primarily intended to capture demand we can realistically win now, within customers and segments we already reach and sell to?",
					},
				},
			},
			{
				ID: "sam", Label: "SAM",
				Questions: []DimensionQuestion{
					{
						ID:       "expands-serviceable-reach",
						Question: "Is this initiative primarily intended to expand into serviceable segments our business model can address but that we do not yet effectively reach or win?",
					},
				},
			},
			{
				ID: "tam_expansion", Label: "TAM Expansion",
				Questions: []DimensionQuestion{
					{
						ID:       "enables-new-market",
						Question: "Does this initiative enable the product to address a meaningful market, segment, geography, regulatory environment, or use case that it cannot effectively address today?",
					},
				},
			},
		},
	}
}
