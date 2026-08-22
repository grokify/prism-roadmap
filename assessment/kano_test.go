package assessment

import "testing"

func yes() KanoCharacteristic {
	return KanoCharacteristic{Answer: true, EvidenceIDs: []string{"EV-1"}}
}

func no() KanoCharacteristic {
	return KanoCharacteristic{Answer: false}
}

func TestResolveKanoMustBe(t *testing.T) {
	// Reliable authentication: expected, absence causes dissatisfaction,
	// presence alone doesn't especially satisfy (ideation doc example).
	a := KanoAnswers{
		Expected:               yes(),
		AbsenceDissatisfaction: yes(),
		PresenceSatisfaction:   no(),
	}
	sel := ResolveKano(a)
	if !sel.Resolved || sel.OptionID != "must_be" {
		t.Errorf("ResolveKano() = %+v, want resolved=must_be", sel)
	}
}

func TestResolveKanoPerformance(t *testing.T) {
	// Audit-search performance: absence dissatisfies, presence satisfies,
	// and more is progressively better (ideation doc example).
	a := KanoAnswers{
		AbsenceDissatisfaction: yes(),
		PresenceSatisfaction:   yes(),
		MoreIsBetter:           yes(),
	}
	sel := ResolveKano(a)
	if !sel.Resolved || sel.OptionID != "performance" {
		t.Errorf("ResolveKano() = %+v, want resolved=performance", sel)
	}
}

func TestResolveKanoAttractive(t *testing.T) {
	a := KanoAnswers{
		Unexpected:           yes(),
		PresenceSatisfaction: yes(),
		Delight:              yes(),
	}
	sel := ResolveKano(a)
	if !sel.Resolved || sel.OptionID != "attractive" {
		t.Errorf("ResolveKano() = %+v, want resolved=attractive", sel)
	}
}

func TestResolveKanoIndifferent(t *testing.T) {
	a := KanoAnswers{
		AbsenceDissatisfaction: no(),
		PresenceSatisfaction:   no(),
		Indifferent:            yes(),
	}
	sel := ResolveKano(a)
	if !sel.Resolved || sel.OptionID != "indifferent" {
		t.Errorf("ResolveKano() = %+v, want resolved=indifferent", sel)
	}
}

func TestResolveKanoReverse(t *testing.T) {
	a := KanoAnswers{
		PresenceSatisfaction: no(),
		ReversePreference:    yes(),
	}
	sel := ResolveKano(a)
	if !sel.Resolved || sel.OptionID != "reverse" {
		t.Errorf("ResolveKano() = %+v, want resolved=reverse", sel)
	}
}

func TestResolveKanoUnresolved(t *testing.T) {
	sel := ResolveKano(KanoAnswers{})
	if sel.Resolved || sel.Ambiguous {
		t.Errorf("ResolveKano(zero value) = %+v, want unresolved", sel)
	}
}

func TestResolveKanoUnsupportedAnswerDoesNotCount(t *testing.T) {
	a := KanoAnswers{
		Expected:               KanoCharacteristic{Answer: true}, // no evidence — must not count
		AbsenceDissatisfaction: yes(),
		PresenceSatisfaction:   no(),
	}
	sel := ResolveKano(a)
	if sel.Resolved {
		t.Errorf("ResolveKano() = %+v, want unresolved (Expected was unsupported)", sel)
	}
}

func TestResolveKanoAmbiguous(t *testing.T) {
	// A contrived contradictory answer set matching both Indifferent and
	// Reverse simultaneously.
	a := KanoAnswers{
		AbsenceDissatisfaction: no(),
		PresenceSatisfaction:   no(),
		Indifferent:            yes(),
		ReversePreference:      yes(),
	}
	sel := ResolveKano(a)
	if !sel.Ambiguous || sel.Resolved {
		t.Errorf("ResolveKano() = %+v, want ambiguous", sel)
	}
	if len(sel.AmbiguousOptionIDs) != 2 {
		t.Errorf("AmbiguousOptionIDs = %v, want 2 entries", sel.AmbiguousOptionIDs)
	}
}

func TestKanoDimensionOptionsMatchResolverIDs(t *testing.T) {
	def := KanoDimension()
	for _, id := range []string{"must_be", "performance", "attractive", "indifferent", "reverse"} {
		if def.OptionByID(id) == nil {
			t.Errorf("KanoDimension() missing option %q referenced by ResolveKano", id)
		}
	}
}

func TestKanoDimensionValidates(t *testing.T) {
	if err := KanoDimension().Validate(); err != nil {
		t.Errorf("KanoDimension().Validate() = %v", err)
	}
}
