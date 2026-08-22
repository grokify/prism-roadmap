package assessment

import (
	"reflect"
	"testing"
)

func testCategoryDefinition() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "test-category", Name: "Test Category", Version: "1.0", Kind: DimensionKindCategory,
		Options: []DimensionOption{
			{ID: "alpha", Label: "Alpha", Questions: []DimensionQuestion{{ID: "q1", Question: "Is it alpha?"}}},
			{ID: "beta", Label: "Beta", Questions: []DimensionQuestion{{ID: "q1", Question: "Is it beta?"}}},
			{ID: "gamma", Label: "Gamma", Questions: []DimensionQuestion{{ID: "q1", Question: "Is it gamma?"}}},
		},
	}
}

func testTagsDefinition() *DimensionDefinition {
	return &DimensionDefinition{
		ID: "test-tags", Name: "Test Tags", Version: "1.0", Kind: DimensionKindTags,
		Options: []DimensionOption{
			{ID: "ai", Label: "AI"},
			{ID: "growth", Label: "Growth"},
			{ID: "excellence", Label: "Excellence"},
		},
	}
}

func TestDimensionDefinitionValidate(t *testing.T) {
	valid := testCategoryDefinition()
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	tests := []struct {
		name   string
		mutate func(*DimensionDefinition)
	}{
		{"missing id", func(d *DimensionDefinition) { d.ID = "" }},
		{"missing name", func(d *DimensionDefinition) { d.Name = "" }},
		{"missing version", func(d *DimensionDefinition) { d.Version = "" }},
		{"invalid kind", func(d *DimensionDefinition) { d.Kind = "bogus" }},
		{"no options", func(d *DimensionDefinition) { d.Options = nil }},
		{"empty option id", func(d *DimensionDefinition) { d.Options[0].ID = "" }},
		{"duplicate option id", func(d *DimensionDefinition) { d.Options[1].ID = d.Options[0].ID }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			def := *valid
			def.Options = append([]DimensionOption{}, valid.Options...)
			tt.mutate(&def)
			if err := def.Validate(); err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestDimensionDefinitionOptionByID(t *testing.T) {
	def := testCategoryDefinition()
	if opt := def.OptionByID("beta"); opt == nil || opt.Label != "Beta" {
		t.Errorf("OptionByID(beta) = %+v", opt)
	}
	if opt := def.OptionByID("nonexistent"); opt != nil {
		t.Errorf("OptionByID(nonexistent) = %+v, want nil", opt)
	}
}

func TestResolveCategoryUnresolved(t *testing.T) {
	def := testCategoryDefinition()
	sel := def.ResolveCategory(nil)
	if sel.Resolved || sel.Ambiguous || sel.OptionID != "" {
		t.Errorf("ResolveCategory(nil) = %+v, want all zero values", sel)
	}
}

func TestResolveCategoryResolved(t *testing.T) {
	def := testCategoryDefinition()
	answers := []DimensionAnswer{
		{OptionID: "alpha", QuestionID: "q1", Answer: false},
		{OptionID: "beta", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Resolved || sel.OptionID != "beta" || sel.Ambiguous {
		t.Errorf("ResolveCategory() = %+v, want resolved=beta", sel)
	}
}

func TestResolveCategoryUnsupportedAnswerDoesNotCount(t *testing.T) {
	def := testCategoryDefinition()
	answers := []DimensionAnswer{
		{OptionID: "alpha", QuestionID: "q1", Answer: true}, // no evidence
		{OptionID: "beta", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	sel := def.ResolveCategory(answers)
	if !sel.Resolved || sel.OptionID != "beta" {
		t.Errorf("ResolveCategory() = %+v, want resolved=beta (alpha unsupported)", sel)
	}
}

func TestResolveCategoryAmbiguous(t *testing.T) {
	def := testCategoryDefinition()
	answers := []DimensionAnswer{
		{OptionID: "alpha", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-1"}},
		{OptionID: "gamma", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-2"}},
	}
	sel := def.ResolveCategory(answers)
	if sel.Resolved || !sel.Ambiguous {
		t.Errorf("ResolveCategory() = %+v, want ambiguous", sel)
	}
	want := []string{"alpha", "gamma"} // definition order
	if !reflect.DeepEqual(sel.AmbiguousOptionIDs, want) {
		t.Errorf("AmbiguousOptionIDs = %v, want %v", sel.AmbiguousOptionIDs, want)
	}
}

func TestResolveTags(t *testing.T) {
	def := testTagsDefinition()
	answers := []DimensionAnswer{
		{OptionID: "ai", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-1"}},
		{OptionID: "growth", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-2"}},
		{OptionID: "excellence", QuestionID: "q1", Answer: false},
	}
	tags := def.ResolveTags(answers)
	want := []string{"ai", "growth"}
	if !reflect.DeepEqual(tags, want) {
		t.Errorf("ResolveTags() = %v, want %v (multiple tags is normal, not ambiguous)", tags, want)
	}
}

func TestNewDimensionAssignmentCategory(t *testing.T) {
	def := testCategoryDefinition()
	answers := []DimensionAnswer{
		{OptionID: "beta", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	a := NewDimensionAssignment(def, answers)
	if a.DimensionID != "test-category" || a.DimensionVersion != "1.0" {
		t.Errorf("assignment identity = %+v", a)
	}
	if a.Category == nil || !a.Category.Resolved || a.Category.OptionID != "beta" {
		t.Errorf("Category = %+v, want resolved=beta", a.Category)
	}
	if a.Tags != nil {
		t.Errorf("Tags = %v, want nil for a category-kind dimension", a.Tags)
	}
}

func TestNewDimensionAssignmentTags(t *testing.T) {
	def := testTagsDefinition()
	answers := []DimensionAnswer{
		{OptionID: "ai", QuestionID: "q1", Answer: true, EvidenceIDs: []string{"EV-1"}},
	}
	a := NewDimensionAssignment(def, answers)
	if a.Category != nil {
		t.Errorf("Category = %+v, want nil for a tags-kind dimension", a.Category)
	}
	want := []string{"ai"}
	if !reflect.DeepEqual(a.Tags, want) {
		t.Errorf("Tags = %v, want %v", a.Tags, want)
	}
}
