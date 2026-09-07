package assessment

// RGT option IDs.
const (
	RGTOptionRun       = "run"
	RGTOptionGrow      = "grow"
	RGTOptionTransform = "transform"
)

// RunGrowTransformDimension is the built-in Run/Grow/Transform (RGT)
// portfolio dimension: Gartner's executive IT investment classification by
// business outcome — Run (continue operating the business as it exists
// today) / Grow (enhance existing capabilities to support growth of the
// current business) / Transform (enable new business models and markets
// the business cannot pursue today). This models Gartner's published
// framework faithfully (see e.g. Gartner IT Key Metrics /
// https://www.gartner.com/imagesrv/ncsc/pdf/ITBudget_Sample_Report.pdf);
// it is an external industry standard, not a prism-roadmap invention.
//
// RGT's known limitation — the reason ProductDevelopmentInvestmentMixDimension
// exists alongside it — is that Run hides the distinction this package
// wants visible: 50% Run could be 40% necessary maintenance + 10% toil or
// 10% maintenance + 40% toil, radically different engineering positions.
// Gartner also notes modernization/automation can still classify as Run
// when it supports existing services, so automation investment is not
// separable inside RGT. Use RGT as the executive roll-up grain and
// Product Development Investment Mix as the finer-resolution grain; the projection
// convention is documented on ProductDevelopmentInvestmentMixDimension and PDIMToRGT.
//
// Like Kano and MIH, RGT is portfolio-descriptive only — it never enters
// Opportunity Rank. Each category has its own independent judge criterion,
// so it resolves through the generic DimensionDefinition.ResolveCategory;
// more than one satisfied criterion reports Ambiguous rather than silently
// picking one.
func RunGrowTransformDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "run-grow-transform", Name: "Run/Grow/Transform (Gartner)", Version: "1.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{
				ID: RGTOptionRun, Label: "Run",
				Questions: []DimensionQuestion{
					{
						ID:       "operates-existing-business",
						Question: "Is this investment primarily required to continue operating the business as it exists today — keeping existing systems, services, and obligations running?",
					},
				},
			},
			{
				ID: RGTOptionGrow, Label: "Grow",
				Questions: []DimensionQuestion{
					{
						ID:       "enhances-for-growth",
						Question: "Is this investment primarily intended to enhance or expand existing capabilities to support growth of the current business?",
					},
				},
			},
			{
				ID: RGTOptionTransform, Label: "Transform",
				Questions: []DimensionQuestion{
					{
						ID:       "enables-new-business",
						Question: "Is this investment primarily intended to enable new business models, markets, or capabilities the business cannot pursue today?",
					},
				},
			},
		},
	}
}
