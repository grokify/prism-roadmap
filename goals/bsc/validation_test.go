package bsc

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
	strict := StrictValidationOptions()
	tests := []struct {
		name    string
		edit    func(b *BSC)
		opts    *ValidationOptions
		path    string
		isError bool
	}{
		{"no vision warns", func(b *BSC) { b.Vision = "" }, nil, "vision", false},
		{"no vision strict", func(b *BSC) { b.Vision = "" }, strict, "vision", true},
		{"mission is enough", func(b *BSC) { b.Vision, b.Mission = "", "Why" }, strict, "", false},
		{"no perspectives", func(b *BSC) { b.Perspectives = nil }, nil, "perspectives", true},
		{"perspective without name", func(b *BSC) { b.Perspectives[0].Name = "" }, nil, "perspectives[0].name", true},
		{"bad perspective kind", func(b *BSC) { b.Perspectives[0].Kind = "people" }, nil, "perspectives[0].kind", true},
		{"empty perspective warns", func(b *BSC) {
			b.Perspectives = append(b.Perspectives, Perspective{ID: "learning", Name: "Learning"})
		}, nil, "perspectives[2]", false},
		{"empty perspective strict", func(b *BSC) {
			b.Perspectives = append(b.Perspectives, Perspective{ID: "learning", Name: "Learning"})
		}, strict, "perspectives[2]", true},
		{"theme without name", func(b *BSC) { b.Themes[0].Name = "" }, nil, "themes[0].name", true},
		{"unused theme warns", func(b *BSC) { b.Themes = append(b.Themes, Theme{ID: "t2", Name: "Growth"}) }, strict, "themes[1]", false},
		{"no objectives", func(b *BSC) { b.Objectives = nil; b.Measures = nil; b.Initiatives = nil }, nil, "objectives", true},
		{"duplicate id across parts", func(b *BSC) { b.Measures[1].ID = "o1" }, nil, "measures[1].id", true},
		{"objective without perspective", func(b *BSC) { b.Objectives[0].PerspectiveID = "" }, nil, "objectives[0].perspectiveId", true},
		{"objective in unknown perspective", func(b *BSC) { b.Objectives[0].PerspectiveID = "nope" }, nil, "objectives[0].perspectiveId", true},
		{"objective in unknown theme", func(b *BSC) { b.Objectives[0].ThemeIDs = []string{"nope"} }, nil, "objectives[0].themeIds[0]", true},
		{"objective without owner strict", func(b *BSC) { b.Objectives[0].Owner = "" }, strict, "objectives[0].owner", true},
		{"bad priority", func(b *BSC) { b.Objectives[0].Priority = "high" }, nil, "objectives[0].priority", true},
		{"contributes to unknown", func(b *BSC) { b.Objectives[1].ContributesTo = []string{"nope"} }, nil, "objectives[1].contributesTo[0]", true},
		{"contributes to itself", func(b *BSC) { b.Objectives[1].ContributesTo = []string{"o2"} }, nil, "objectives[1].contributesTo[0]", true},
		{"link runs down warns", func(b *BSC) { b.Objectives[0].ContributesTo = []string{"o2"}; b.Objectives[1].ContributesTo = nil }, nil, "objectives[0].contributesTo[0]", false},
		{"cycle", func(b *BSC) { b.Objectives[0].ContributesTo = []string{"o2"} }, nil, "objectives", true},
		{"lower objective contributes to nothing warns", func(b *BSC) { b.Objectives[1].ContributesTo = nil }, nil, "objectives[1].contributesTo", false},
		{"lower objective contributes to nothing strict", func(b *BSC) { b.Objectives[1].ContributesTo = nil }, strict, "objectives[1].contributesTo", true},
		{"nothing contributes to upper objective strict", func(b *BSC) { b.Objectives[1].ContributesTo = nil }, strict, "objectives[0]", true},
		{"objective without measure strict", func(b *BSC) { b.Measures = b.Measures[1:] }, strict, "objectives[0]", true},
		{"objective without initiative warns even strict", func(b *BSC) { b.Initiatives[0].ObjectiveIDs = []string{"o2"} }, strict, "objectives[0]", false},
		{"measure without objective", func(b *BSC) { b.Measures[0].ObjectiveID = "" }, nil, "measures[0].objectiveId", true},
		{"measure with unknown objective", func(b *BSC) { b.Measures[0].ObjectiveID = "nope" }, nil, "measures[0].objectiveId", true},
		{"measure without type warns", func(b *BSC) { b.Measures[0].Type = "" }, nil, "measures[0].type", false},
		{"measure without type strict", func(b *BSC) { b.Measures[0].Type = "" }, strict, "measures[0].type", true},
		{"bad measure type", func(b *BSC) { b.Measures[0].Type = "leading" }, nil, "measures[0].type", true},
		{"only lag measures strict", func(b *BSC) { b.Measures[1].Type = TypeLag }, strict, "measures", true},
		{"measure without target strict", func(b *BSC) { b.Measures[0].Target = "" }, strict, "measures[0].target", true},
		{"measure without owner warns", func(b *BSC) { b.Measures[0].Owner = "" }, nil, "measures[0].owner", false},
		{"bad direction", func(b *BSC) { b.Measures[0].Direction = "up" }, nil, "measures[0].direction", true},
		{"bad measure status", func(b *BSC) { b.Measures[0].Status = "Green" }, nil, "measures[0].status", true},
		{"bad due date", func(b *BSC) { b.Measures[0].DueDate = "soon" }, nil, "measures[0].dueDate", true},
		{"progress above one", func(b *BSC) { b.Measures[0].Progress = 2 }, nil, "measures[0].progress", true},
		{"initiative without name", func(b *BSC) { b.Initiatives[0].Name = "" }, nil, "initiatives[0].name", true},
		{"initiative without objectives", func(b *BSC) { b.Initiatives[0].ObjectiveIDs = nil }, nil, "initiatives[0].objectiveIds", true},
		{"initiative with unknown objective", func(b *BSC) { b.Initiatives[0].ObjectiveIDs = []string{"nope"} }, nil, "initiatives[0].objectiveIds[0]", true},
		{"initiative without owner strict", func(b *BSC) { b.Initiatives[0].Owner = "" }, strict, "initiatives[0].owner", true},
		{"bad initiative status", func(b *BSC) { b.Initiatives[0].Status = "Doing" }, nil, "initiatives[0].status", true},
		{"initiative ends before it starts", func(b *BSC) { b.Initiatives[0].StartDate, b.Initiatives[0].EndDate = "2027-03-01", "2027-02-01" }, nil, "initiatives[0].endDate", true},
		{"unknown dependency", func(b *BSC) { b.Initiatives[1].Dependencies = []string{"zz"} }, nil, "initiatives[1].dependencies[0]", true},
		{"self dependency", func(b *BSC) { b.Initiatives[1].Dependencies = []string{"i2"} }, nil, "initiatives[1].dependencies[0]", true},
		{"bad document status", func(b *BSC) { b.Metadata.Status = "Live" }, nil, "metadata.status", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := sample()
			tt.edit(b)
			errs := b.Validate(tt.opts)
			if tt.path == "" {
				if len(Errors(errs)) > 0 {
					t.Errorf("want no errors, got %v", errs)
				}
				return
			}
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

func TestCycleNamesItsPath(t *testing.T) {
	b := sample()
	b.Objectives[0].ContributesTo = []string{"o2"}
	for _, e := range b.Validate(nil) {
		if e.Path == "objectives" && strings.Contains(e.Message, "o1 -> o2 -> o1") {
			return
		}
	}
	t.Errorf("cycle not reported with its path: %v", b.Validate(nil))
}

func TestMaxMeasuresWarns(t *testing.T) {
	opts := StrictValidationOptions()
	opts.MaxMeasures = 1
	for _, e := range sample().Validate(opts) {
		if e.Path == "measures" && !e.IsError && strings.Contains(e.Message, "more than a team can steer by") {
			return
		}
	}
	t.Error("no warning for too many measures")
}

func TestValidationErrorAndFilters(t *testing.T) {
	errs := []ValidationError{{Path: "a", Message: "bad", IsError: true}, {Path: "b", Message: "meh"}}
	if errs[0].Error() != "a: bad" || len(Errors(errs)) != 1 || len(Warnings(errs)) != 1 || IsValid(errs) || !IsValid(errs[1:]) {
		t.Error("ValidationError helpers are wrong")
	}
}
