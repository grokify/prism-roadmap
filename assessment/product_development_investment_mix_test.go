package assessment

import "testing"

func TestProductDevelopmentInvestmentMixDimensionValidates(t *testing.T) {
	def := ProductDevelopmentInvestmentMixDimension()
	if err := def.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
	for _, id := range []string{"innovate", "improve", "automate", "maintain", "toil"} {
		if def.OptionByID(id) == nil {
			t.Errorf("missing option %q", id)
		}
	}
}

func TestPDIMRollup(t *testing.T) {
	cases := map[string]string{
		"innovate": "innovate", "improve": "improve", "automate": "automate",
		"maintain": "maintain", "toil": "maintain",
	}
	for in, want := range cases {
		if got := PDIMRollup(in); got != want {
			t.Errorf("PDIMRollup(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestPDIMResolveAndAssignment(t *testing.T) {
	def := ProductDevelopmentInvestmentMixDimension()
	answers := []DimensionAnswer{
		{OptionID: "toil", QuestionID: "repetitive-manual-sustaining-work", Answer: false},
		{OptionID: "automate", QuestionID: "reduces-recurring-effort", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	a := NewDimensionAssignment(def, answers)
	if a.DimensionID != "product-development-investment-mix" || a.Category == nil || a.Category.OptionID != PDIMOptionAutomate {
		t.Errorf("NewDimensionAssignment() = %+v, want resolved=automate", a)
	}
}

func TestPDIMAmbiguousWhenMultipleCriteriaMatch(t *testing.T) {
	def := ProductDevelopmentInvestmentMixDimension()
	// Work that reads as both deliberate maintenance and toil — surfaces
	// as ambiguous for review, not silently resolved.
	answers := []DimensionAnswer{
		{OptionID: "maintain", QuestionID: "preserves-system-health", Answer: true, EvidenceIDs: []string{"EV-1"}},
		{OptionID: "toil", QuestionID: "repetitive-manual-sustaining-work", Answer: true, EvidenceIDs: []string{"EV-2"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Ambiguous {
		t.Errorf("ResolveCategory() = %+v, want ambiguous", sel)
	}
}

func TestPDIMToRGT(t *testing.T) {
	cases := map[string]string{
		"innovate": "transform",
		"improve":  "grow",
		"automate": "run",
		"maintain": "run",
		"toil":     "run",
	}
	for in, want := range cases {
		got, ok := PDIMToRGT(in)
		if !ok || got != want {
			t.Errorf("PDIMToRGT(%q) = %q, %v, want %q, true", in, got, ok, want)
		}
	}
	if got, ok := PDIMToRGT("bogus"); ok || got != "" {
		t.Errorf("PDIMToRGT(bogus) = %q, %v, want \"\", false", got, ok)
	}
}

func TestPDIMToSREWork(t *testing.T) {
	cases := map[string]string{
		"innovate": "engineering",
		"improve":  "engineering",
		"automate": "engineering",
		"maintain": "engineering",
		"toil":     "toil",
	}
	for in, want := range cases {
		got, ok := PDIMToSREWork(in)
		if !ok || got != want {
			t.Errorf("PDIMToSREWork(%q) = %q, %v, want %q, true", in, got, ok, want)
		}
	}
	if got, ok := PDIMToSREWork("bogus"); ok || got != "" {
		t.Errorf("PDIMToSREWork(bogus) = %q, %v, want \"\", false", got, ok)
	}
}

func TestToilReductionValidate(t *testing.T) {
	valid := ToilReduction{ToilSource: "manual certificate rotation", BaselineHoursPerMonth: 80, TargetHoursPerMonth: 10}
	if err := valid.Validate(); err != nil {
		t.Errorf("Validate() = %v, want nil", err)
	}
	invalid := []ToilReduction{
		{BaselineHoursPerMonth: 80},                                           // missing source
		{ToilSource: "x", BaselineHoursPerMonth: -1},                          // negative baseline
		{ToilSource: "x", BaselineHoursPerMonth: 10, TargetHoursPerMonth: -1}, // negative target
		{ToilSource: "x", BaselineHoursPerMonth: 10, TargetHoursPerMonth: 20}, // target > baseline
	}
	for i, tr := range invalid {
		if err := tr.Validate(); err == nil {
			t.Errorf("case %d: Validate() = nil, want error", i)
		}
	}
}

func TestToilReductionHoursReclaimedAndPayback(t *testing.T) {
	tr := ToilReduction{ToilSource: "manual environment provisioning", BaselineHoursPerMonth: 80, TargetHoursPerMonth: 10}
	if got := tr.HoursReclaimedPerMonth(); got != 70 {
		t.Errorf("HoursReclaimedPerMonth() = %v, want 70", got)
	}

	// 20 engineer-days automating 5 engineer-days/month of toil → 4 months.
	payback := ToilReduction{ToilSource: "manual release validation", BaselineHoursPerMonth: 40, TargetHoursPerMonth: 0}
	months, err := payback.PaybackPeriodMonths(160)
	if err != nil {
		t.Fatalf("PaybackPeriodMonths() error: %v", err)
	}
	if months != 4 {
		t.Errorf("PaybackPeriodMonths() = %v, want 4", months)
	}

	noReclaim := ToilReduction{ToilSource: "x", BaselineHoursPerMonth: 10, TargetHoursPerMonth: 10}
	if _, err := noReclaim.PaybackPeriodMonths(100); err == nil {
		t.Error("PaybackPeriodMonths() with zero reclaimed hours: want error, got nil")
	}
	if _, err := payback.PaybackPeriodMonths(-1); err == nil {
		t.Error("PaybackPeriodMonths(-1): want error, got nil")
	}
}
