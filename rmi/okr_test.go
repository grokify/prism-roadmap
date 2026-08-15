package rmi

import (
	"testing"

	"github.com/grokify/prism-roadmap/goals/okr"
	"github.com/grokify/prism-roadmap/prioritization"
)

// scoredItem builds a roadmap item with a RICE score linked to the given
// objective IDs.
func scoredItem(id string, reach int, impact prioritization.ImpactLevel, effort float64, objRefs ...string) RoadmapItem {
	item := NewRoadmapItem(id, id, prioritization.MoSCoWMustHave)
	item.WithRICE(prioritization.NewRICEScore(id, reach, impact, prioritization.ConfidenceHigh, effort))
	for _, ref := range objRefs {
		item.AddObjectiveRef(ref)
	}
	return *item
}

func TestRoadmapItem_SupportsObjective(t *testing.T) {
	item := NewRoadmapItem("r1", "R1", prioritization.MoSCoWMustHave)
	item.AddObjectiveRef("obj-1")

	if !item.SupportsObjective("obj-1") {
		t.Error("SupportsObjective(obj-1) = false, want true")
	}
	if item.SupportsObjective("obj-2") {
		t.Error("SupportsObjective(obj-2) = true, want false")
	}
}

func TestRICEByObjective(t *testing.T) {
	items := []RoadmapItem{
		scoredItem("r1", 1000, prioritization.ImpactHigh, 2, "obj-1"),    // score 1000
		scoredItem("r2", 500, prioritization.ImpactMedium, 1, "obj-1"),   // score 500
		scoredItem("r3", 1000, prioritization.ImpactMassive, 1, "obj-2"), // score 3000
	}
	objectives := []okr.Objective{
		{ID: "obj-1", Title: "Grow activation"},
		{ID: "obj-2", Title: "Reduce churn"},
		{ID: "obj-3", Title: "Unbacked objective"},
	}

	rollups := RICEByObjective(items, objectives)
	if len(rollups) != 3 {
		t.Fatalf("RICEByObjective returned %d rollups, want 3", len(rollups))
	}

	// Ordered by TotalRICE desc: obj-2 (3000) > obj-1 (1500) > obj-3 (0).
	if rollups[0].ObjectiveID != "obj-2" {
		t.Errorf("rollups[0] = %s, want obj-2", rollups[0].ObjectiveID)
	}
	if rollups[0].TotalRICE != 3000 {
		t.Errorf("obj-2 TotalRICE = %v, want 3000", rollups[0].TotalRICE)
	}

	obj1 := rollups[1]
	if obj1.ObjectiveID != "obj-1" {
		t.Fatalf("rollups[1] = %s, want obj-1", obj1.ObjectiveID)
	}
	if obj1.TotalRICE != 1500 {
		t.Errorf("obj-1 TotalRICE = %v, want 1500", obj1.TotalRICE)
	}
	if obj1.ScoredCount != 2 {
		t.Errorf("obj-1 ScoredCount = %d, want 2", obj1.ScoredCount)
	}
	if obj1.MeanRICE != 750 {
		t.Errorf("obj-1 MeanRICE = %v, want 750", obj1.MeanRICE)
	}
	if obj1.TopRICE != 1000 {
		t.Errorf("obj-1 TopRICE = %v, want 1000", obj1.TopRICE)
	}

	// obj-3 has no linked items: included with zero aggregates.
	last := rollups[2]
	if last.ObjectiveID != "obj-3" || last.TotalRICE != 0 || len(last.Items) != 0 {
		t.Errorf("obj-3 rollup = %+v, want empty aggregates", last)
	}
}

func TestUnlinkedScoredItems(t *testing.T) {
	items := []RoadmapItem{
		scoredItem("linked", 1000, prioritization.ImpactHigh, 2, "obj-1"),
		scoredItem("orphan", 500, prioritization.ImpactMedium, 1), // no objective ref
	}
	// An item with no RICE score should not count as an orphan.
	noScore := NewRoadmapItem("noscore", "noscore", prioritization.MoSCoWMustHave)
	items = append(items, *noScore)

	orphans := UnlinkedScoredItems(items)
	if len(orphans) != 1 {
		t.Fatalf("UnlinkedScoredItems returned %d, want 1", len(orphans))
	}
	if orphans[0].ID != "orphan" {
		t.Errorf("orphan = %s, want orphan", orphans[0].ID)
	}
}

func TestService_RICEByObjective(t *testing.T) {
	svc := NewService()
	svc.Set().Add(scoredItem("r1", 1000, prioritization.ImpactHigh, 2, "obj-1"))

	doc := &okr.OKRDocument{Objectives: []okr.Objective{{ID: "obj-1", Title: "Grow"}}}
	rollups := svc.RICEByObjective(doc)
	if len(rollups) != 1 || rollups[0].TotalRICE != 1000 {
		t.Errorf("Service.RICEByObjective = %+v, want one rollup with TotalRICE 1000", rollups)
	}

	if svc.RICEByObjective(nil) != nil {
		t.Error("Service.RICEByObjective(nil) should return nil")
	}
}
