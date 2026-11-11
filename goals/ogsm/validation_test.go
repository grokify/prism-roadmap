package ogsm

import "testing"

func TestValidateSampleIsClean(t *testing.T) {
	if errs := sample().Validate(StrictValidationOptions()); len(errs) > 0 {
		t.Errorf("sample should validate cleanly under strict options: %v", errs)
	}
}

func TestValidate(t *testing.T) {
	strict := StrictValidationOptions()
	tests := []struct {
		name    string
		edit    func(o *OGSM)
		opts    *ValidationOptions
		path    string
		isError bool
	}{
		{"no objective", func(o *OGSM) { o.Objective.Statement = " " }, nil, "objective.statement", true},
		{"bad objective date", func(o *OGSM) { o.Objective.DueDate = "2029" }, nil, "objective.dueDate", true},
		{"no objective date warns", func(o *OGSM) { o.Objective.DueDate = "" }, nil, "objective.dueDate", false},
		{"no objective date strict", func(o *OGSM) { o.Objective.DueDate = "" }, strict, "objective.dueDate", true},
		{"no goals", func(o *OGSM) { o.Goals = nil; o.Strategies[0].GoalIDs = nil; o.Strategies[1].GoalIDs = nil }, nil, "goals", true},
		{"no strategies", func(o *OGSM) { o.Strategies = nil; o.Measures = nil; o.Actions = nil }, nil, "strategies", true},
		{"goal without id", func(o *OGSM) { o.Goals[0].ID = "" }, nil, "goals[0].id", true},
		{"duplicate id across parts", func(o *OGSM) { o.Actions[1].ID = "g1"; o.Actions[1].Dependencies = nil }, nil, "actions[1].id", true},
		{"goal without name", func(o *OGSM) { o.Goals[0].Name = "" }, nil, "goals[0].name", true},
		{"goal without target warns", func(o *OGSM) { o.Goals[0].Target = "" }, nil, "goals[0].target", false},
		{"goal without target strict", func(o *OGSM) { o.Goals[0].Target = "" }, strict, "goals[0].target", true},
		{"goal without owner strict", func(o *OGSM) { o.Goals[0].Owner = "" }, strict, "goals[0].owner", true},
		{"bad goal direction", func(o *OGSM) { o.Goals[0].Direction = "up" }, nil, "goals[0].direction", true},
		{"bad goal status", func(o *OGSM) { o.Goals[0].Status = "Green" }, nil, "goals[0].status", true},
		{"goal progress above one", func(o *OGSM) { o.Goals[0].Progress = 2 }, nil, "goals[0].progress", true},
		{"goal no strategy serves strict", func(o *OGSM) { o.Strategies[1].GoalIDs = []string{"g1"} }, strict, "goals[1]", true},
		{"strategy without name", func(o *OGSM) { o.Strategies[0].Name = "" }, nil, "strategies[0].name", true},
		{"bad priority", func(o *OGSM) { o.Strategies[0].Priority = "high" }, nil, "strategies[0].priority", true},
		{"strategy serves unknown goal", func(o *OGSM) { o.Strategies[0].GoalIDs = []string{"nope"} }, nil, "strategies[0].goalIds[0]", true},
		{"strategy serves no goal warns", func(o *OGSM) { o.Strategies[0].GoalIDs = nil }, nil, "strategies[0].goalIds", false},
		{"strategy without measure strict", func(o *OGSM) { o.Measures = o.Measures[1:] }, strict, "strategies[0]", true},
		{"strategy without action warns", func(o *OGSM) { o.Actions = o.Actions[:1] }, nil, "strategies[1]", false},
		{"measure without strategy", func(o *OGSM) { o.Measures[0].StrategyID = "" }, nil, "measures[0].strategyId", true},
		{"measure with unknown strategy", func(o *OGSM) { o.Measures[0].StrategyID = "nope" }, nil, "measures[0].strategyId", true},
		{"measure without target strict", func(o *OGSM) { o.Measures[0].Target = "" }, strict, "measures[0].target", true},
		{"bad measure due date", func(o *OGSM) { o.Measures[0].DueDate = "soon" }, nil, "measures[0].dueDate", true},
		{"action without name", func(o *OGSM) { o.Actions[0].Name = "" }, nil, "actions[0].name", true},
		{"action without strategy", func(o *OGSM) { o.Actions[0].StrategyID = "" }, nil, "actions[0].strategyId", true},
		{"action without owner strict", func(o *OGSM) { o.Actions[0].Owner = "" }, strict, "actions[0].owner", true},
		{"bad action status", func(o *OGSM) { o.Actions[0].Status = "Doing" }, nil, "actions[0].status", true},
		{"action ends before it starts", func(o *OGSM) { o.Actions[0].StartDate, o.Actions[0].EndDate = "2027-03-01", "2027-02-01" }, nil, "actions[0].endDate", true},
		{"unknown dependency", func(o *OGSM) { o.Actions[1].Dependencies = []string{"zz"} }, nil, "actions[1].dependencies[0]", true},
		{"self dependency", func(o *OGSM) { o.Actions[1].Dependencies = []string{"a2"} }, nil, "actions[1].dependencies[0]", true},
		{"bad document status", func(o *OGSM) { o.Metadata.Status = "Live" }, nil, "metadata.status", true},
		{"too many goals", func(o *OGSM) {
			for _, id := range []string{"g3", "g4", "g5", "g6"} {
				o.Goals = append(o.Goals, Goal{ID: id, Name: id, Owner: "x", Target: "1"})
				o.Strategies[0].GoalIDs = append(o.Strategies[0].GoalIDs, id)
			}
		}, strict, "goals", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := sample()
			tt.edit(o)
			errs := o.Validate(tt.opts)
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
	errs := []ValidationError{{Path: "a", Message: "bad", IsError: true}, {Path: "b", Message: "meh"}}
	if errs[0].Error() != "a: bad" || len(Errors(errs)) != 1 || len(Warnings(errs)) != 1 || IsValid(errs) || !IsValid(errs[1:]) {
		t.Error("ValidationError helpers are wrong")
	}
}
