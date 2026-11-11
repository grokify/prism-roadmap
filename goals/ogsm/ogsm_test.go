package ogsm

import (
	"path/filepath"
	"reflect"
	"testing"
)

// sample returns a small OGSM that meets the strict options.
func sample() *OGSM {
	return &OGSM{
		Metadata:  &Metadata{ID: "plan", Name: "Plan", Status: StatusActive},
		Objective: Objective{Statement: "Be the app owners keep", DueDate: "2029-12-31"},
		Goals: []Goal{
			{ID: "g1", Name: "Retention", Owner: "Ana", Target: "65"},
			{ID: "g2", Name: "Paying accounts", Owner: "Ben", Target: "100000"},
		},
		Strategies: []Strategy{
			{ID: "s1", Name: "Onboard", GoalIDs: []string{"g1"}},
			{ID: "s2", Name: "Partners", GoalIDs: []string{"g2"}},
		},
		Measures: []Measure{
			{ID: "m1", Name: "Connect rate", StrategyID: "s1", Owner: "Ana", Target: "75"},
			{ID: "m2", Name: "Partner accounts", StrategyID: "s2", Owner: "Ben", Target: "15000"},
		},
		Actions: []Action{
			{ID: "a1", Name: "Guided import", StrategyID: "s1", Owner: "Squad"},
			{ID: "a2", Name: "Partner program", StrategyID: "s2", Owner: "Partnerships", Dependencies: []string{"a1"}},
		},
	}
}

func TestHelpers(t *testing.T) {
	o := sample()
	if g := o.Goal("g2"); g == nil || g.Name != "Paying accounts" {
		t.Errorf("Goal(g2) = %v", g)
	}
	if o.Goal("nope") != nil || o.Strategy("nope") != nil {
		t.Error("unknown IDs should give nil")
	}
	if s := o.Strategy("s1"); s == nil || s.Name != "Onboard" {
		t.Errorf("Strategy(s1) = %v", s)
	}
	if m := o.MeasuresFor("s2"); len(m) != 1 || m[0].ID != "m2" {
		t.Errorf("MeasuresFor(s2) = %v", m)
	}
	if a := o.ActionsFor("s1"); len(a) != 1 || a[0].ID != "a1" {
		t.Errorf("ActionsFor(s1) = %v", a)
	}
	if s := o.StrategiesFor("g1"); len(s) != 1 || s[0].ID != "s1" {
		t.Errorf("StrategiesFor(g1) = %v", s)
	}
}

func TestNew(t *testing.T) {
	o := New("id", "Name", "Objective")
	if o.Metadata.ID != "id" || o.Metadata.Status != StatusDraft || o.Objective.Statement != "Objective" {
		t.Errorf("New() = %+v %+v", o, o.Metadata)
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	o := sample()
	path := filepath.Join(t.TempDir(), "plan.ogsm.json")
	if err := o.WriteFile(path); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, o) {
		t.Errorf("round trip changed the OGSM:\n got %+v\nwant %+v", got, o)
	}
	if _, err := Parse([]byte("{")); err == nil {
		t.Error("Parse accepted invalid JSON")
	}
}
