package goals

import (
	"testing"

	"github.com/grokify/prism-roadmap/goals/ogsm"
)

func ogsmGoals() *Goals {
	return NewOGSM(&ogsm.OGSM{
		Objective:  ogsm.Objective{Statement: "Be the app owners keep"},
		Goals:      []ogsm.Goal{{ID: "g1", Name: "Retention", Baseline: "52", Target: "65", Unit: "%"}},
		Strategies: []ogsm.Strategy{{ID: "s1", Name: "Onboard", Owner: "Ana", Priority: "P0", GoalIDs: []string{"g1"}}},
		Measures:   []ogsm.Measure{{ID: "m1", Name: "Connect rate", StrategyID: "s1", Current: "51", Progress: 0.2}},
	})
}

func TestOGSMGoals(t *testing.T) {
	g := ogsmGoals()
	if !g.IsOGSM() || g.IsVSEM() || g.Framework != FrameworkOGSM {
		t.Errorf("framework = %q", g.Framework)
	}
	var nilGoals *Goals
	if nilGoals.IsOGSM() {
		t.Error("nil.IsOGSM() should be false")
	}
	items := g.GoalItems()
	if len(items) != 1 || items[0].ID != "s1" || items[0].Owner != "Ana" || items[0].Priority != "P0" {
		t.Errorf("GoalItems() = %+v", items)
	}
	results := g.ResultItems()
	if len(results) != 2 {
		t.Fatalf("want a goal and a measure, got %+v", results)
	}
	if results[0].ID != "g1" || results[0].GoalID != "" || results[0].Target != "65" {
		t.Errorf("goal result = %+v, want no strategy", results[0])
	}
	if results[1].ID != "m1" || results[1].GoalID != "s1" || results[1].Score != 0.2 {
		t.Errorf("measure result = %+v", results[1])
	}
	if g.GoalLabel() != "Strategies" || g.ResultLabel() != "Goals and Measures" {
		t.Errorf("labels = %q, %q", g.GoalLabel(), g.ResultLabel())
	}
	if (&Goals{Framework: FrameworkOGSM}).ResultItems() != nil {
		t.Error("an OGSM framework without data should have no results")
	}
}
