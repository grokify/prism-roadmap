package assessment

import "testing"

func TestMarketInvestmentHorizonDimensionValidates(t *testing.T) {
	def := MarketInvestmentHorizonDimension()
	if err := def.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
	for _, id := range []string{"ktlo", "sam_som", "tam_expansion"} {
		if def.OptionByID(id) == nil {
			t.Errorf("missing option %q", id)
		}
	}
}

func TestMarketInvestmentHorizonResolveViaGenericCategory(t *testing.T) {
	def := MarketInvestmentHorizonDimension()
	answers := []DimensionAnswer{
		{OptionID: "ktlo", QuestionID: "sustains-existing-business", Answer: false},
		{OptionID: "tam_expansion", QuestionID: "enables-new-market", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Resolved || sel.OptionID != "tam_expansion" {
		t.Errorf("ResolveCategory() = %+v, want resolved=tam_expansion", sel)
	}
}

func TestMarketInvestmentHorizonAmbiguousWhenMultipleCriteriaMatch(t *testing.T) {
	def := MarketInvestmentHorizonDimension()
	// FedRAMP-style initiative: plausibly reads as both KTLO (compliance)
	// and TAM Expansion (new addressable market) — this should surface as
	// ambiguous, not be silently resolved to one.
	answers := []DimensionAnswer{
		{OptionID: "ktlo", QuestionID: "sustains-existing-business", Answer: true, EvidenceIDs: []string{"EV-1"}},
		{OptionID: "tam_expansion", QuestionID: "enables-new-market", Answer: true, EvidenceIDs: []string{"EV-2"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Ambiguous {
		t.Errorf("ResolveCategory() = %+v, want ambiguous", sel)
	}
}

func TestMarketInvestmentHorizonAssignment(t *testing.T) {
	def := MarketInvestmentHorizonDimension()
	answers := []DimensionAnswer{
		{OptionID: "sam_som", QuestionID: "improves-current-market-position", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	a := NewDimensionAssignment(def, answers)
	if a.DimensionID != "market-investment-horizon" || a.Category == nil || a.Category.OptionID != "sam_som" {
		t.Errorf("NewDimensionAssignment() = %+v", a)
	}
}
