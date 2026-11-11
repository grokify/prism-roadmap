package goals

import (
	"testing"

	"github.com/grokify/prism-roadmap/goals/vsem"
)

func vsemGoals() *Goals {
	return NewVSEM(&vsem.VSEM{
		Vision: "Owners close their books in an hour a month.",
		Strategies: []vsem.Strategy{
			{ID: "s1", Name: "Onboard", Description: "First week", Owner: "Ana", Priority: "P0"},
			{ID: "s2", Name: "Retain", Owner: "Ben"},
		},
		Metrics: []vsem.Metric{
			{ID: "m0", Name: "Active at 90 days", Baseline: "52", Target: "60", Unit: "%"},
			{ID: "m1", Name: "Connect rate", StrategyID: "s1", Current: "51", Progress: 0.24, Status: "On Track"},
		},
	})
}

func TestNewVSEM(t *testing.T) {
	g := vsemGoals()
	if !g.IsVSEM() || g.IsOKR() || g.IsV2MOM() || g.Framework != FrameworkVSEM {
		t.Errorf("NewVSEM gave framework %q", g.Framework)
	}
	var nilGoals *Goals
	if nilGoals.IsVSEM() {
		t.Error("nil.IsVSEM() should be false")
	}
	if (&Goals{Framework: FrameworkVSEM}).GoalItems() != nil {
		t.Error("a VSEM framework without data should have no goal items")
	}
}

func TestGoalItemsVSEM(t *testing.T) {
	items := vsemGoals().GoalItems()
	if len(items) != 2 {
		t.Fatalf("want 2 goal items, got %d", len(items))
	}
	want := GoalItem{ID: "s1", Title: "Onboard", Description: "First week", Owner: "Ana", Priority: "P0"}
	if items[0].ID != want.ID || items[0].Title != want.Title || items[0].Description != want.Description ||
		items[0].Owner != want.Owner || items[0].Priority != want.Priority {
		t.Errorf("GoalItems()[0] = %+v, want %+v", items[0], want)
	}
}

func TestResultItemsVSEM(t *testing.T) {
	items := vsemGoals().ResultItems()
	if len(items) != 2 {
		t.Fatalf("want 2 result items, got %d", len(items))
	}
	if items[0].GoalID != "" || items[0].Baseline != "52" || items[0].Target != "60" || items[0].Unit != "%" {
		t.Errorf("vision metric = %+v, want no goal and its baseline, target, and unit", items[0])
	}
	if items[1].GoalID != "s1" || items[1].Current != "51" || items[1].Score != 0.24 || items[1].Status != "On Track" {
		t.Errorf("strategy metric = %+v", items[1])
	}
}

func TestLabelsVSEM(t *testing.T) {
	g := vsemGoals()
	if g.GoalLabel() != "Strategies" || g.ResultLabel() != "Metrics" {
		t.Errorf("labels = %q, %q", g.GoalLabel(), g.ResultLabel())
	}
	if g.GoalCount() != 2 || g.ResultCount() != 2 {
		t.Errorf("counts = %d, %d", g.GoalCount(), g.ResultCount())
	}
}
