package assessment

import "fmt"

// PDIM option IDs, plus the 4-way maintenance rollup key.
const (
	PDIMOptionInnovate = "innovate"
	PDIMOptionImprove  = "improve"
	PDIMOptionAutomate = "automate"
	PDIMOptionMaintain = "maintain"
	PDIMOptionToil     = "toil"

	// PDIMRollupMaintain is the rollup key non-toil maintenance and toil
	// collapse to in the 4-bucket executive view (Innovate / Improve /
	// Automate / Maintain), where Maintenance = non-toil maintenance +
	// toil. It reuses the "maintain" id so 4-bucket views and rolled-up
	// 5-way views share one key.
	PDIMRollupMaintain = PDIMOptionMaintain
)

// PDIMRollup maps a 5-way PDIM option ID to its 4-bucket rollup key
// (Innovate / Improve / Automate / Maintain). Toil collapses into
// Maintain — Maintenance is the combination of non-toil maintenance and
// toil; the other categories pass through.
func PDIMRollup(optionID string) string {
	if optionID == PDIMOptionToil {
		return PDIMRollupMaintain
	}
	return optionID
}

// ProductDevelopmentInvestmentMixDimension is the built-in Product Development Investment Mix
// (PDIM) portfolio dimension: where product-development capacity is being
// invested — Innovate (create new outcomes/capabilities) / Improve (make
// existing outcomes materially better) / Automate (reduce the recurring
// human effort required to produce the same outcome) / Maintain (deliberate
// non-repetitive work preserving existing system health) / Toil
// (repetitive, manual, automatable sustaining work with no enduring value).
// This is the mix of WORK TYPES within the product portfolio, not an
// allocation across products. "Product development" names the joint
// function — product management, engineering, design, docs — so the
// framework stays neutral between the product and engineering teams that
// present it together: product toil (manual status reporting, repetitive
// data pulls, hand-assembled release notes) is Toil exactly as
// operational toil is.
// This is prism-roadmap's own framework (ideation doc), synthesizing two
// external standards it ships alongside: Gartner's Run/Grow/Transform
// (RunGrowTransformDimension) provides the strategic investment lens and
// Google SRE's work classification (SREWorkDimension) provides the
// engineering-efficiency lens.
//
// The key distinction between the categories is whether work changes the
// OUTCOME (Innovate, Improve) or changes the COST of producing the same
// outcome (Automate, and negatively, Toil). Automation is an investment;
// toil is a workload — the framework exists to make their causal loop
// visible and measurable:
//
//	Toil identified → Automation investment → Toil eliminated →
//	Capacity reclaimed → Reinvested in Improve/Innovate
//
// Conceptually Maintenance = non-toil maintenance + toil (not all
// maintenance is toil: a major database migration is deliberate sustaining
// engineering; manually rotating certificates every month is toil). Toil
// is nonetheless promoted to a first-class category at the storage grain
// because exposing and reducing it is an explicit objective — toil is
// consumed capacity to reclaim, not a deliberate portfolio choice, and it
// is where automation investment captures value (see ToilReduction). The
// 5-way split is the canonical presentation while toil reduction is an
// active objective; PDIMRollup provides the cleaner 4-bucket executive view
// (Innovate / Improve / Automate / Maintain) once toil is small.
//
// PDIM is the ASSESSMENT grain for the three frameworks: an opportunity
// judged once at PDIM resolution deterministically projects onto the other
// two views (PDIMToRGT, PDIMToSREWork), guaranteeing the executive roll-up
// and the engineering-health view stay mutually consistent:
//
//	PDIM        → RGT (convention)   → SRE Work (definitional)
//	Innovate     Transform            Engineering
//	Improve      Grow                 Engineering
//	Automate     Run                  Engineering (toil elimination)
//	Maintain     Run                  Engineering
//	Toil         Run                  Toil
//
// The projections are derived-by-default, explicit-override: when an
// assessment also carries a native RunGrowTransformDimension or
// SREWorkDimension assignment, that explicit judgment takes precedence
// over the PDIM-derived one (e.g. automation that unlocks growth capacity
// can be explicitly judged RGT Grow even though the convention says Run).
//
// Like Kano and MIH, PDIM is portfolio-descriptive only — it never enters
// Opportunity Rank. Each category has its own independent judge criterion,
// so it resolves through the generic DimensionDefinition.ResolveCategory;
// more than one satisfied criterion reports Ambiguous rather than silently
// picking one.
func ProductDevelopmentInvestmentMixDimension() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "product-development-investment-mix", Name: "Product Development Investment Mix", Version: "1.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{
				ID: PDIMOptionInnovate, Label: "Innovate",
				Questions: []DimensionQuestion{
					{
						ID:       "creates-new-outcome",
						Question: "Is this work primarily intended to create a new outcome or capability that does not exist today?",
					},
				},
			},
			{
				ID: PDIMOptionImprove, Label: "Improve",
				Questions: []DimensionQuestion{
					{
						ID:       "improves-existing-outcome",
						Question: "Is this work primarily intended to make an existing outcome or capability materially better?",
					},
				},
			},
			{
				ID: PDIMOptionAutomate, Label: "Automate",
				Questions: []DimensionQuestion{
					{
						ID:       "reduces-recurring-effort",
						Question: "Is this work primarily an investment in reducing the recurring human effort required to produce the same outcome — e.g. automating an identified source of toil?",
					},
				},
			},
			{
				ID: PDIMOptionMaintain, Label: "Maintain",
				Questions: []DimensionQuestion{
					{
						ID:       "preserves-system-health",
						Question: "Is this deliberate, non-repetitive work necessary to preserve existing system health and capability — e.g. upgrades, migrations, or patches?",
					},
				},
			},
			{
				ID: PDIMOptionToil, Label: "Toil",
				Questions: []DimensionQuestion{
					{
						ID:       "repetitive-manual-sustaining-work",
						Question: "Is this repetitive, manual, automatable work required to keep the system operating that produces no enduring value?",
					},
				},
			},
		},
	}
}

// PDIMToSREWork deterministically projects a PDIM option onto Google SRE's
// work-classification rollup grain (SREWorkRollupEngineering or
// SREWorkOptionToil). Unlike PDIMToRGT this projection is definitionally
// sound, not merely conventional: PDIM's Toil is Google SRE's toil
// definition verbatim (generalized beyond production services to any
// product-delivery work), and every other PDIM category is deliberate
// engineering-style work with enduring value in SRE terms. SRE's Overhead
// has no PDIM preimage — administrative overhead sits outside PDIM's
// universe, which classifies product-development investment only. Returns
// ok=false for an unknown option ID.
func PDIMToSREWork(optionID string) (sreWorkID string, ok bool) {
	switch optionID {
	case PDIMOptionInnovate, PDIMOptionImprove, PDIMOptionAutomate, PDIMOptionMaintain:
		return SREWorkRollupEngineering, true
	case PDIMOptionToil:
		return SREWorkOptionToil, true
	default:
		return "", false
	}
}

// PDIMToRGT deterministically projects a PDIM option onto its conventional
// Gartner RGT category: Innovate → Transform, Improve → Grow,
// Automate/Maintain/Toil → Run. This is a documented CONVENTION, not a
// definitional equivalence — RGT classifies by business outcome while PDIM
// classifies by what the work does to the outcome vs. its production
// cost, so the axes can genuinely diverge (automation that unlocks growth
// capacity reads as Grow; a new capability sold to the current business
// reads as Grow, not Transform). The convention follows Gartner's own
// guidance that modernization supporting existing services classifies as
// Run. An explicit RunGrowTransformDimension assignment on the same
// assessment takes precedence over this projection; use the projection to
// fill the RGT executive roll-up for opportunities assessed only at the
// PDIM grain. Returns ok=false for an unknown option ID.
func PDIMToRGT(optionID string) (rgtID string, ok bool) {
	switch optionID {
	case PDIMOptionInnovate:
		return RGTOptionTransform, true
	case PDIMOptionImprove:
		return RGTOptionGrow, true
	case PDIMOptionAutomate, PDIMOptionMaintain, PDIMOptionToil:
		return RGTOptionRun, true
	default:
		return "", false
	}
}

// ToilReduction quantifies the value an automation investment captures
// from an identified source of toil — the measurable link in the PDIM loop
// (Toil → Automate → Capacity Reclaimed). Attach it to an Automate-
// classified initiative to state, up front, which toil it targets and how
// much recurring capacity it reclaims; that turns automation into a
// capital-style investment with an ROI and payback period rather than a
// cost line. Follows this package's evidence discipline: cite the evidence
// behind the baseline measurement.
type ToilReduction struct {
	// ToilSource names the recurring manual work being eliminated (e.g.
	// "manual certificate rotation").
	ToilSource string `json:"toilSource"`

	// BaselineHoursPerMonth is the recurring effort the toil consumes
	// today.
	BaselineHoursPerMonth float64 `json:"baselineHoursPerMonth"`

	// TargetHoursPerMonth is the recurring effort expected to remain after
	// the automation lands (0 for full elimination).
	TargetHoursPerMonth float64 `json:"targetHoursPerMonth"`

	// EvidenceIDs cite the evidence behind the baseline measurement.
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// Validate returns an error if the reduction is not usable.
func (t ToilReduction) Validate() error {
	if t.ToilSource == "" {
		return fmt.Errorf("toilSource is required")
	}
	if t.BaselineHoursPerMonth < 0 {
		return fmt.Errorf("baselineHoursPerMonth must be >= 0, got %v", t.BaselineHoursPerMonth)
	}
	if t.TargetHoursPerMonth < 0 {
		return fmt.Errorf("targetHoursPerMonth must be >= 0, got %v", t.TargetHoursPerMonth)
	}
	if t.TargetHoursPerMonth > t.BaselineHoursPerMonth {
		return fmt.Errorf("targetHoursPerMonth (%v) must not exceed baselineHoursPerMonth (%v)", t.TargetHoursPerMonth, t.BaselineHoursPerMonth)
	}
	return nil
}

// HoursReclaimedPerMonth is the recurring capacity the automation frees
// for reinvestment: baseline minus target.
func (t ToilReduction) HoursReclaimedPerMonth() float64 {
	return t.BaselineHoursPerMonth - t.TargetHoursPerMonth
}

// PaybackPeriodMonths is the capital-investment view of an automation
// initiative: how many months of reclaimed capacity repay the one-time
// automation effort (automation investment / recurring toil savings) —
// e.g. 20 engineer-days spent automating 5 engineer-days/month of toil
// pays back in 4 months. Returns an error when the reduction reclaims no
// hours (payback is undefined) or the effort is negative.
func (t ToilReduction) PaybackPeriodMonths(automationEffortHours float64) (float64, error) {
	if automationEffortHours < 0 {
		return 0, fmt.Errorf("automationEffortHours must be >= 0, got %v", automationEffortHours)
	}
	reclaimed := t.HoursReclaimedPerMonth()
	if reclaimed <= 0 {
		return 0, fmt.Errorf("no hours reclaimed (baseline %v, target %v): payback period is undefined", t.BaselineHoursPerMonth, t.TargetHoursPerMonth)
	}
	return automationEffortHours / reclaimed, nil
}
