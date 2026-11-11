package vsem

import (
	"strings"
	"testing"
)

func TestValidateSampleIsClean(t *testing.T) {
	if errs := sample().Validate(StrictValidationOptions()); len(errs) > 0 {
		t.Errorf("sample should validate cleanly under strict options: %v", errs)
	}
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		edit    func(v *VSEM)
		opts    *ValidationOptions
		path    string
		isError bool
	}{
		{"no vision or mission", func(v *VSEM) { v.Vision = "" }, nil, "vision", true},
		{"no strategies", func(v *VSEM) { v.Strategies = nil; v.Execution = nil; v.Metrics = v.Metrics[:1] }, nil, "strategies", true},
		{"strategy without a name", func(v *VSEM) { v.Strategies[0].Name = " " }, nil, "strategies[0].name", true},
		{"strategy without an id", func(v *VSEM) { v.Strategies[1].ID = "" }, nil, "strategies[1].id", true},
		{"duplicate id across sections", func(v *VSEM) { v.Metrics[0].ID = "e1" }, nil, "metrics[0].id", true},
		{"bad priority", func(v *VSEM) { v.Strategies[0].Priority = "urgent" }, nil, "strategies[0].priority", true},
		{"resource without a name", func(v *VSEM) { v.Strategies[0].Resources[0].Name = "" }, nil, "strategies[0].resources[0].name", true},
		{"no resources warns", func(v *VSEM) { v.Strategies[0].Resources = nil }, nil, "strategies[0].resources", false},
		{"no resources strict", func(v *VSEM) { v.Strategies[0].Resources = nil }, StrictValidationOptions(), "strategies[0].resources", true},
		{"strategy without execution warns", func(v *VSEM) { v.Execution = v.Execution[:1]; v.Execution[0].Dependencies = nil }, nil, "strategies[1]", false},
		{"execution without strategy", func(v *VSEM) { v.Execution[0].StrategyID = "" }, nil, "execution[0].strategyId", true},
		{"execution with unknown strategy", func(v *VSEM) { v.Execution[0].StrategyID = "nope" }, nil, "execution[0].strategyId", true},
		{"bad kind", func(v *VSEM) { v.Execution[0].Kind = "epic" }, nil, "execution[0].kind", true},
		{"bad execution status", func(v *VSEM) { v.Execution[0].Status = "Doing" }, nil, "execution[0].status", true},
		{"execution without owner warns", func(v *VSEM) { v.Execution[0].Owner = "" }, nil, "execution[0].owner", false},
		{"execution without owner strict", func(v *VSEM) { v.Execution[0].Owner = "" }, StrictValidationOptions(), "execution[0].owner", true},
		{"bad start date", func(v *VSEM) { v.Execution[0].StartDate = "2027-13-01" }, nil, "execution[0].startDate", true},
		{"end before start", func(v *VSEM) { v.Execution[0].StartDate = "2027-03-01"; v.Execution[0].EndDate = "2027-02-01" }, nil, "execution[0].endDate", true},
		{"unknown dependency", func(v *VSEM) { v.Execution[2].Dependencies = []string{"zz"} }, nil, "execution[2].dependencies[0]", true},
		{"self dependency", func(v *VSEM) { v.Execution[2].Dependencies = []string{"e3"} }, nil, "execution[2].dependencies[0]", true},
		{"metric with unknown strategy", func(v *VSEM) { v.Metrics[1].StrategyID = "nope" }, nil, "metrics[1].strategyId", true},
		{"metric without owner warns", func(v *VSEM) { v.Metrics[1].Owner = "" }, nil, "metrics[1].owner", false},
		{"metric without owner strict", func(v *VSEM) { v.Metrics[1].Owner = "" }, StrictValidationOptions(), "metrics[1].owner", true},
		{"metric without target warns", func(v *VSEM) { v.Metrics[1].Target = "" }, nil, "metrics[1].target", false},
		{"bad direction", func(v *VSEM) { v.Metrics[1].Direction = "up" }, nil, "metrics[1].direction", true},
		{"bad metric status", func(v *VSEM) { v.Metrics[1].Status = "Green" }, nil, "metrics[1].status", true},
		{"progress above one", func(v *VSEM) { v.Metrics[1].Progress = 1.5 }, nil, "metrics[1].progress", true},
		{"bad due date", func(v *VSEM) { v.Metrics[1].DueDate = "soon" }, nil, "metrics[1].dueDate", true},
		{"bad document status", func(v *VSEM) { v.Metadata.Status = "Live" }, nil, "metadata.status", true},
		{"too many strategies", func(v *VSEM) {
			for i := range 4 {
				id := string(rune('a' + i))
				v.Strategies = append(v.Strategies, Strategy{ID: id, Name: id})
			}
		}, StrictValidationOptions(), "strategies", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := sample()
			tt.edit(v)
			errs := v.Validate(tt.opts)
			for _, e := range errs {
				if e.Path == tt.path && e.IsError == tt.isError {
					return
				}
			}
			kind := "warning"
			if tt.isError {
				kind = "error"
			}
			t.Errorf("want an %s at %s, got %v", kind, tt.path, errs)
		})
	}
}

func TestValidationErrorAndFilters(t *testing.T) {
	errs := []ValidationError{
		{Path: "a", Message: "bad", IsError: true},
		{Path: "b", Message: "meh"},
	}
	if got := errs[0].Error(); got != "a: bad" {
		t.Errorf("Error() = %q", got)
	}
	if len(Errors(errs)) != 1 || len(Warnings(errs)) != 1 {
		t.Errorf("Errors/Warnings split wrong: %v / %v", Errors(errs), Warnings(errs))
	}
	if IsValid(errs) {
		t.Error("IsValid should be false with an error")
	}
	if !IsValid(errs[1:]) {
		t.Error("IsValid should be true with only warnings")
	}
}

func TestMissionAloneIsEnough(t *testing.T) {
	v := sample()
	v.Vision = ""
	v.Mission = "Take the busywork out of bookkeeping."
	for _, e := range Errors(v.Validate(nil)) {
		if strings.HasPrefix(e.Path, "vision") {
			t.Errorf("a mission without a vision should be accepted: %v", e)
		}
	}
}
