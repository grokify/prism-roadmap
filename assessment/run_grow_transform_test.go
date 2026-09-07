package assessment

import "testing"

func TestRunGrowTransformDimensionValidates(t *testing.T) {
	def := RunGrowTransformDimension()
	if err := def.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
	for _, id := range []string{"run", "grow", "transform"} {
		if def.OptionByID(id) == nil {
			t.Errorf("missing option %q", id)
		}
	}
}

func TestRunGrowTransformResolveViaGenericCategory(t *testing.T) {
	def := RunGrowTransformDimension()
	answers := []DimensionAnswer{
		{OptionID: "run", QuestionID: "operates-existing-business", Answer: false},
		{OptionID: "transform", QuestionID: "enables-new-business", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Resolved || sel.OptionID != RGTOptionTransform {
		t.Errorf("ResolveCategory() = %+v, want resolved=transform", sel)
	}
}

func TestRunGrowTransformAmbiguousWhenMultipleCriteriaMatch(t *testing.T) {
	def := RunGrowTransformDimension()
	// Modernization that both keeps existing services running and enables
	// a new market — surfaces as ambiguous, not silently resolved.
	answers := []DimensionAnswer{
		{OptionID: "run", QuestionID: "operates-existing-business", Answer: true, EvidenceIDs: []string{"EV-1"}},
		{OptionID: "transform", QuestionID: "enables-new-business", Answer: true, EvidenceIDs: []string{"EV-2"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Ambiguous {
		t.Errorf("ResolveCategory() = %+v, want ambiguous", sel)
	}
}
