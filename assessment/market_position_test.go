package assessment

import "testing"

func TestMarketPositionDimensionValidates(t *testing.T) {
	def := MarketPositionDimension()
	if err := def.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
	for _, id := range []string{"star", "cash_cow", "question_mark", "dog", "enabler"} {
		if def.OptionByID(id) == nil {
			t.Errorf("missing option %q", id)
		}
	}
}

func TestMarketPositionResolveViaGenericCategory(t *testing.T) {
	def := MarketPositionDimension()
	answers := []DimensionAnswer{
		{OptionID: "star", QuestionID: "high-growth-high-share", Answer: false},
		{OptionID: "cash_cow", QuestionID: "low-growth-high-share", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Resolved || sel.OptionID != BCGCashCow {
		t.Errorf("ResolveCategory() = %+v, want resolved=cash_cow", sel)
	}
}

func TestMarketPositionAmbiguousWhenMultipleCriteriaMatch(t *testing.T) {
	def := MarketPositionDimension()
	// A platform line that reads as both a Star and an Enabler — surfaces
	// as ambiguous for review, not silently resolved.
	answers := []DimensionAnswer{
		{OptionID: "star", QuestionID: "high-growth-high-share", Answer: true, EvidenceIDs: []string{"EV-1"}},
		{OptionID: "enabler", QuestionID: "horizontal-leverage", Answer: true, EvidenceIDs: []string{"EV-2"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Ambiguous {
		t.Errorf("ResolveCategory() = %+v, want ambiguous", sel)
	}
}
