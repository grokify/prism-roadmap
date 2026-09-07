package assessment

import "testing"

func TestSREWorkDimensionValidates(t *testing.T) {
	def := SREWorkDimension()
	if err := def.Validate(); err != nil {
		t.Errorf("Validate() = %v", err)
	}
	for _, id := range []string{"software_engineering", "systems_engineering", "toil", "overhead"} {
		if def.OptionByID(id) == nil {
			t.Errorf("missing option %q", id)
		}
	}
}

func TestSREWorkRollup(t *testing.T) {
	cases := map[string]string{
		"software_engineering": "engineering",
		"systems_engineering":  "engineering",
		"engineering":          "engineering",
		"toil":                 "toil",
		"overhead":             "overhead",
	}
	for in, want := range cases {
		if got := SREWorkRollup(in); got != want {
			t.Errorf("SREWorkRollup(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSREWorkResolveViaGenericCategory(t *testing.T) {
	def := SREWorkDimension()
	answers := []DimensionAnswer{
		{OptionID: "software_engineering", QuestionID: "lasting-code-improvement", Answer: false},
		{OptionID: "toil", QuestionID: "manual-repetitive-automatable", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Resolved || sel.OptionID != SREWorkOptionToil {
		t.Errorf("ResolveCategory() = %+v, want resolved=toil", sel)
	}
}
