package assessment

// MarketInvestmentHorizonDimension is the built-in Market Investment
// Horizon (MIH) portfolio dimension: where relative to the current market
// an investment creates or protects value — KTLO (protect the existing
// business), SAM+SOM (grow within markets already served), or TAM
// Expansion (enter markets not effectively addressable today). Adapted
// from an industry investment-classification practice combining KTLO with
// TAM/SAM/SOM market concepts (ideation doc); this is prism-roadmap's own
// framework, not a published external standard.
//
// Unlike Kano, MIH's three categories each have their own independent
// judge criterion, so this dimension resolves through the generic
// DimensionDefinition.ResolveCategory rather than a bespoke resolver — no
// separate ResolveMIH function is needed. If more than one criterion is
// satisfied (e.g. an initiative plausibly reads as both KTLO and TAM
// Expansion), ResolveCategory reports that as Ambiguous rather than
// silently picking one; that ambiguity is itself useful portfolio
// information, not just a rubric defect.
//
// Like Kano, MIH is portfolio-descriptive only — it never enters
// Opportunity Rank.
func MarketInvestmentHorizonDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "market-investment-horizon", Name: "Market Investment Horizon", Version: "1.0", Kind: DimensionKindCategory,
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
				ID: "sam_som", Label: "SAM + SOM",
				Questions: []DimensionQuestion{
					{
						ID:       "improves-current-market-position",
						Question: "Is this initiative primarily intended to improve adoption, retention, expansion, competitiveness, or customer value within markets we already address?",
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
