package goals

import (
	"testing"

	"github.com/grokify/prism-roadmap/goals/bsc"
)

func bscGoals() *Goals {
	return NewBSC(&bsc.BSC{
		Vision:       "Be the app owners keep",
		Perspectives: []bsc.Perspective{{ID: "customer", Name: "Customer"}},
		Objectives:   []bsc.Objective{{ID: "o1", Name: "Keep owners", PerspectiveID: "customer", Owner: "Ana", Priority: "P0"}},
		Measures: []bsc.Measure{{ID: "m1", Name: "90-day retention", ObjectiveID: "o1", Type: bsc.TypeLag,
			Baseline: "52", Target: "65", Current: "56", Unit: "%", Progress: 0.3}},
	})
}

func TestBSCGoals(t *testing.T) {
	g := bscGoals()
	if !g.IsBSC() || g.IsOGSM() || g.Framework != FrameworkBSC {
		t.Errorf("framework = %q", g.Framework)
	}
	var nilGoals *Goals
	if nilGoals.IsBSC() {
		t.Error("nil.IsBSC() should be false")
	}
	items := g.GoalItems()
	if len(items) != 1 || items[0].ID != "o1" || items[0].Owner != "Ana" || items[0].Priority != "P0" ||
		len(items[0].Tags) != 1 || items[0].Tags[0] != "Customer" {
		t.Errorf("GoalItems() = %+v, want the objective tagged with its perspective", items)
	}
	results := g.ResultItems()
	if len(results) != 1 || results[0].ID != "m1" || results[0].GoalID != "o1" || results[0].Target != "65" || results[0].Score != 0.3 {
		t.Errorf("ResultItems() = %+v", results)
	}
	if g.GoalLabel() != "Objectives" || g.ResultLabel() != "Measures" {
		t.Errorf("labels = %q, %q", g.GoalLabel(), g.ResultLabel())
	}
	empty := &Goals{Framework: FrameworkBSC}
	if empty.GoalItems() != nil || empty.ResultItems() != nil {
		t.Error("a Balanced Scorecard framework without data should have no items")
	}
}
