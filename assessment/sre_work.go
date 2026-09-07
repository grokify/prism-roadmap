package assessment

// SRE Work option IDs, plus the 2-way engineering rollup key.
const (
	SREWorkOptionSoftwareEngineering = "software_engineering"
	SREWorkOptionSystemsEngineering  = "systems_engineering"
	SREWorkOptionToil                = "toil"
	SREWorkOptionOverhead            = "overhead"

	// SREWorkRollupEngineering is the rollup key software and systems
	// engineering collapse to — Google's headline split is Engineering vs.
	// Toil vs. Overhead (the "at least 50% engineering" rule), with the
	// software/systems distinction as the finer storage grain.
	SREWorkRollupEngineering = "engineering"
)

// SREWorkRollup maps a 4-way SRE Work option ID to its 3-way rollup key
// (Engineering / Toil / Overhead). Software and systems engineering
// collapse to SREWorkRollupEngineering; toil and overhead pass through.
func SREWorkRollup(optionID string) string {
	switch optionID {
	case SREWorkOptionSoftwareEngineering, SREWorkOptionSystemsEngineering, SREWorkRollupEngineering:
		return SREWorkRollupEngineering
	default:
		return optionID
	}
}

// SREWorkDimension is the built-in SRE Work Classification portfolio
// dimension: Google SRE's categorization of where engineering time goes —
// Software Engineering / Systems Engineering / Toil / Overhead. This
// models the classification from the Google SRE Book's "Eliminating Toil"
// chapter (https://sre.google/sre-book/eliminating-toil/; see also the SRE
// Workbook, https://sre.google/workbook/eliminating-toil/) faithfully; it
// is an external industry standard, not a prism-roadmap invention. The
// name reflects the concept modeled; Google SRE is the source.
//
// Toil is the load-bearing category: work that is manual, repetitive,
// automatable, tactical, devoid of enduring value, and that scales
// linearly with service growth. Google's discipline — measure toil,
// prioritize it, engineer it out, and cap it (SREs spend at least 50% of
// their time on engineering work; the threshold is Google's, not
// universal) — is the industry precedent for treating toil as a place to
// capture value, which ProductDevelopmentInvestmentMixDimension builds on. The
// 3-way rollup (Engineering / Toil / Overhead, see SREWorkRollup) is the
// canonical presentation grain; the software/systems split is the storage
// grain.
//
// Like Kano and MIH, SRE Work is portfolio-descriptive only — it never
// enters Opportunity Rank. Each category has its own independent judge
// criterion, so it resolves through the generic
// DimensionDefinition.ResolveCategory; more than one satisfied criterion
// reports Ambiguous rather than silently picking one.
func SREWorkDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "sre-work", Name: "SRE Work Classification (Google SRE)", Version: "1.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{
				ID: SREWorkOptionSoftwareEngineering, Label: "Software Engineering",
				Questions: []DimensionQuestion{
					{
						ID:       "lasting-code-improvement",
						Question: "Does this work primarily involve writing or changing code (with its design and documentation) to produce a lasting improvement — e.g. automation, tooling, or service features for reliability or performance?",
					},
				},
			},
			{
				ID: SREWorkOptionSystemsEngineering, Label: "Systems Engineering",
				Questions: []DimensionQuestion{
					{
						ID:       "lasting-system-improvement",
						Question: "Does this work primarily involve configuring production systems or consulting on system architecture to produce a lasting improvement — e.g. monitoring setup, load-balancer configuration, or OS tuning?",
					},
				},
			},
			{
				ID: SREWorkOptionToil, Label: "Toil",
				Questions: []DimensionQuestion{
					{
						ID:       "manual-repetitive-automatable",
						Question: "Is this work manual, repetitive, automatable, tactical, and without enduring value — tied to running a production service and scaling roughly linearly with service growth?",
					},
				},
			},
			{
				ID: SREWorkOptionOverhead, Label: "Overhead",
				Questions: []DimensionQuestion{
					{
						ID:       "administrative-not-service-tied",
						Question: "Is this administrative work not tied directly to running a production service — e.g. meetings, hiring, HR paperwork, or training?",
					},
				},
			},
		},
	}
}
