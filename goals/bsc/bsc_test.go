package bsc

import (
	"path/filepath"
	"reflect"
	"testing"
)

// sample returns a small scorecard that meets the strict options with no
// findings: two perspectives, an objective in each linked cause to effect,
// a lag and a lead measure, and an initiative serving both objectives.
func sample() *BSC {
	return &BSC{
		Metadata: &Metadata{ID: "card", Name: "Card", Status: StatusActive},
		Vision:   "Be the app owners keep",
		Perspectives: []Perspective{
			{ID: "customer", Name: "Customer", Kind: KindCustomer},
			{ID: "process", Name: "Internal Process", Kind: KindInternalProcess},
		},
		Themes: []Theme{{ID: "t1", Name: "Retention"}},
		Objectives: []Objective{
			{ID: "o1", Name: "Keep owners", PerspectiveID: "customer", ThemeIDs: []string{"t1"}, Owner: "Ana"},
			{ID: "o2", Name: "Onboard fast", PerspectiveID: "process", ThemeIDs: []string{"t1"}, Owner: "Ben",
				ContributesTo: []string{"o1"}},
		},
		Measures: []Measure{
			{ID: "m1", Name: "90-day retention", ObjectiveID: "o1", Type: TypeLag, Owner: "Ana", Target: "65"},
			{ID: "m2", Name: "First-week reconciliation", ObjectiveID: "o2", Type: TypeLead, Owner: "Ben", Target: "40"},
		},
		Initiatives: []Initiative{
			{ID: "i1", Name: "Guided import", ObjectiveIDs: []string{"o2", "o1"}, Owner: "Squad"},
			{ID: "i2", Name: "Import help", ObjectiveIDs: []string{"o2"}, Owner: "Support", Dependencies: []string{"i1"}},
		},
	}
}

func TestHelpers(t *testing.T) {
	b := sample()
	if p := b.Perspective("process"); p == nil || p.Name != "Internal Process" {
		t.Errorf("Perspective(process) = %v", p)
	}
	if o := b.Objective("o2"); o == nil || o.Name != "Onboard fast" {
		t.Errorf("Objective(o2) = %v", o)
	}
	if th := b.Theme("t1"); th == nil || th.Name != "Retention" {
		t.Errorf("Theme(t1) = %v", th)
	}
	if b.Perspective("nope") != nil || b.Objective("nope") != nil || b.Theme("nope") != nil {
		t.Error("unknown IDs should give nil")
	}
	if got := b.ObjectivesIn("customer"); len(got) != 1 || got[0].ID != "o1" {
		t.Errorf("ObjectivesIn(customer) = %v", got)
	}
	if got := b.ObjectivesInTheme("t1"); len(got) != 2 {
		t.Errorf("ObjectivesInTheme(t1) = %v", got)
	}
	if got := b.ContributorsTo("o1"); len(got) != 1 || got[0].ID != "o2" {
		t.Errorf("ContributorsTo(o1) = %v", got)
	}
	if got := b.MeasuresFor("o2"); len(got) != 1 || got[0].ID != "m2" {
		t.Errorf("MeasuresFor(o2) = %v", got)
	}
	if got := b.InitiativesFor("o2"); len(got) != 2 {
		t.Errorf("InitiativesFor(o2) = %v", got)
	}
	if b.PerspectiveIndex("process") != 1 || b.PerspectiveIndex("nope") != -1 {
		t.Error("PerspectiveIndex is wrong")
	}
}

func TestNew(t *testing.T) {
	b := New("id", "Name")
	if b.Metadata.ID != "id" || b.Metadata.Status != StatusDraft || len(b.Perspectives) != 4 {
		t.Errorf("New() = %+v %+v", b, b.Metadata)
	}
	kinds := []string{KindFinancial, KindCustomer, KindInternalProcess, KindLearningGrowth}
	for i, p := range b.Perspectives {
		if p.Kind != kinds[i] || p.Question == "" {
			t.Errorf("perspective %d = %+v, want kind %s with a question", i, p, kinds[i])
		}
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	b := sample()
	path := filepath.Join(t.TempDir(), "card.bsc.json")
	if err := b.WriteFile(path); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, b) {
		t.Errorf("round trip changed the scorecard:\n got %+v\nwant %+v", got, b)
	}
	if _, err := Parse([]byte("{")); err == nil {
		t.Error("Parse accepted invalid JSON")
	}
}
