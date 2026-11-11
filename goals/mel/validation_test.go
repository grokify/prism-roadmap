package mel

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
		edit    func(p *Plan)
		opts    *ValidationOptions
		path    string
		isError bool
	}{
		{"no results", func(p *Plan) {
			p.Results = nil
			p.Indicators = nil
			p.Assumptions = nil
			p.Evaluations[0].RelatedIDs = nil
			p.Learning[0].RelatedIDs = nil
		}, nil, "results", true},
		{"result without statement", func(p *Plan) { p.Results[0].Statement = "" }, nil, "results[0].statement", true},
		{"result without level", func(p *Plan) { p.Results[0].Level = "" }, nil, "results[0].level", true},
		{"bad level", func(p *Plan) { p.Results[0].Level = "goal" }, nil, "results[0].level", true},
		{"missing id", func(p *Plan) { p.Reviews[0].ID = "" }, nil, "reviews[0].id", true},
		{"duplicate id across parts", func(p *Plan) { p.Reports[0].ID = "q-1" }, nil, "reports[0].id", true},
		{"contributes to itself", func(p *Plan) { p.Results[1].ContributesTo = []string{"r-long"} }, nil, "results[1].contributesTo[0]", true},
		{"contributes to unknown", func(p *Plan) { p.Results[1].ContributesTo = []string{"nope"} }, nil, "results[1].contributesTo[0]", true},
		{"contributes to an indicator", func(p *Plan) { p.Results[1].ContributesTo = []string{"i-hours"} }, nil, "results[1].contributesTo[0]", true},
		{"contributes downward warns", func(p *Plan) { p.Results[0].ContributesTo = []string{"r-short"} }, nil, "results[0].contributesTo[0]", false},
		{"cycle", func(p *Plan) { p.Results[0].ContributesTo = []string{"r-short"} }, nil, "results", true},
		{"result without indicator warns", func(p *Plan) { p.Indicators = p.Indicators[1:] }, nil, "results[0]", false},
		{"result without indicator strict", func(p *Plan) { p.Indicators = p.Indicators[1:] }, strict, "results[0]", true},
		{"indicator without result", func(p *Plan) { p.Indicators[0].Measures = "" }, nil, "indicators[0].measures", true},
		{"indicator measures a review", func(p *Plan) { p.Indicators[0].Measures = "rv-1" }, nil, "indicators[0].measures", true},
		{"activity without name", func(p *Plan) { p.Activities[0].Name = "" }, nil, "activities[0].name", true},
		{"activity produces unknown", func(p *Plan) { p.Activities[0].ResultIDs = []string{"nope"} }, nil, "activities[0].resultIds[0]", true},
		{"bad activity status", func(p *Plan) { p.Activities[0].Status = "Doing" }, nil, "activities[0].status", true},
		{"activity ends before it starts", func(p *Plan) { p.Activities[0].EndDate = "2026-12-01" }, nil, "activities[0].endDate", true},
		{"question about an unknown activity", func(p *Plan) { p.Learning[0].RelatedIDs = []string{"act-nope"} }, nil, "learningAgenda[0].relatedIds[0]", true},
		{"bad indicator type", func(p *Plan) { p.Indicators[0].Type = "number" }, nil, "indicators[0].type", true},
		{"bad direction", func(p *Plan) { p.Indicators[0].Direction = "up" }, nil, "indicators[0].direction", true},
		{"baseline without value", func(p *Plan) { p.Indicators[0].Baseline.Value = "" }, nil, "indicators[0].baseline.value", true},
		{"no baseline warns", func(p *Plan) { p.Indicators[0].Baseline = nil }, nil, "indicators[0].baseline", false},
		{"no baseline strict", func(p *Plan) { p.Indicators[0].Baseline = nil }, strict, "indicators[0].baseline", true},
		{"no target strict", func(p *Plan) { p.Indicators[0].Targets = nil }, strict, "indicators[0].targets", true},
		{"bad target date", func(p *Plan) { p.Indicators[0].Targets[0].Date = "soon" }, nil, "indicators[0].targets[0].date", true},
		{"no data source warns", func(p *Plan) { p.Indicators[0].DataSource = "" }, nil, "indicators[0].dataSource", false},
		{"no method strict", func(p *Plan) { p.Indicators[0].CollectionMethod = "" }, strict, "indicators[0].collectionMethod", true},
		{"no frequency strict", func(p *Plan) { p.Indicators[0].Frequency = "" }, strict, "indicators[0].frequency", true},
		{"no owner strict", func(p *Plan) { p.Indicators[0].Owner = "" }, strict, "indicators[0].owner", true},
		{"bad quality date", func(p *Plan) { p.Indicators[0].DataQuality = &DataQuality{LastAssessed: "2027"} }, nil, "indicators[0].dataQuality.lastAssessed", true},
		{"assumption unknown result", func(p *Plan) { p.Assumptions[0].ResultIDs = []string{"nope"} }, nil, "assumptions[0].resultIds[0]", true},
		{"assumption not monitored strict", func(p *Plan) { p.Assumptions[0].Monitoring = "" }, strict, "assumptions[0].monitoring", true},
		{"no evaluations warns", func(p *Plan) { p.Evaluations = nil }, nil, "evaluations", false},
		{"evaluation without questions strict", func(p *Plan) { p.Evaluations[0].Questions = nil }, strict, "evaluations[0].questions", true},
		{"evaluation without dates warns", func(p *Plan) { p.Evaluations[0].StartDate, p.Evaluations[0].EndDate = "", "" }, nil, "evaluations[0]", false},
		{"evaluation without dates strict", func(p *Plan) { p.Evaluations[0].StartDate, p.Evaluations[0].EndDate = "", "" }, strict, "evaluations[0]", true},
		{"evaluation ends before it starts", func(p *Plan) { p.Evaluations[0].EndDate = "2027-01-01" }, nil, "evaluations[0].endDate", true},
		{"bad evaluation type", func(p *Plan) { p.Evaluations[0].Type = "summative" }, nil, "evaluations[0].type", true},
		{"bad evaluator", func(p *Plan) { p.Evaluations[0].Evaluator = "vendor" }, nil, "evaluations[0].evaluator", true},
		{"evaluation related to a review", func(p *Plan) { p.Evaluations[0].RelatedIDs = []string{"rv-1"} }, nil, "evaluations[0].relatedIds[0]", true},
		{"no learning questions strict", func(p *Plan) { p.Learning = nil; p.Reviews[0].QuestionIDs = nil; p.Evaluations[0].RelatedIDs = nil }, strict, "learningAgenda", true},
		{"question without due date strict", func(p *Plan) { p.Learning[0].DueDate = "" }, strict, "learningAgenda[0].dueDate", true},
		{"question with bad due date", func(p *Plan) { p.Learning[0].DueDate = "June" }, nil, "learningAgenda[0].dueDate", true},
		{"question without methods warns", func(p *Plan) { p.Learning[0].Methods = nil }, strict, "learningAgenda[0].methods", false},
		{"answered without answer warns", func(p *Plan) { p.Learning[0].Status = QuestionAnswered }, nil, "learningAgenda[0].answer", false},
		{"bad question status", func(p *Plan) { p.Learning[0].Status = "Done" }, nil, "learningAgenda[0].status", true},
		{"no reviews strict", func(p *Plan) { p.Reviews = nil }, strict, "reviews", true},
		{"review without frequency warns", func(p *Plan) { p.Reviews[0].Frequency = "" }, nil, "reviews[0].frequency", false},
		{"review takes up an indicator", func(p *Plan) { p.Reviews[0].QuestionIDs = []string{"i-hours"} }, nil, "reviews[0].questionIds[0]", true},
		{"report without audience", func(p *Plan) { p.Reports[0].Audience = "" }, nil, "reports[0].audience", true},
		{"report unknown indicator", func(p *Plan) { p.Reports[0].IndicatorIDs = []string{"nope"} }, nil, "reports[0].indicatorIds[0]", true},
		{"bad document status", func(p *Plan) { p.Metadata.Status = "Live" }, nil, "metadata.status", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := sample()
			tt.edit(p)
			errs := p.Validate(tt.opts)
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
	if errs[0].Error() != "a: bad" {
		t.Errorf("Error() = %q", errs[0].Error())
	}
	if len(Errors(errs)) != 1 || len(Warnings(errs)) != 1 || IsValid(errs) || !IsValid(errs[1:]) {
		t.Error("filters are wrong")
	}
}
