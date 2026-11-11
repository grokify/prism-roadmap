package vsem

import (
	"path/filepath"
	"reflect"
	"testing"
)

func sample() *VSEM {
	return &VSEM{
		Metadata: &Metadata{ID: "plan", Name: "Plan", Status: StatusActive},
		Vision:   "Owners close their books in an hour a month.",
		Strategies: []Strategy{
			{ID: "s1", Name: "Onboard", Owner: "Ana", Resources: []Resource{{Name: "Engineering", Allocation: "2 squads"}}},
			{ID: "s2", Name: "Retain", Owner: "Ben", Resources: []Resource{{Name: "Design"}}},
		},
		Execution: []Execution{
			{ID: "e1", Name: "Import", StrategyID: "s1", Owner: "Squad A"},
			{ID: "e2", Name: "Reminders", StrategyID: "s2", Owner: "Squad B"},
			{ID: "e3", Name: "Walkthrough", StrategyID: "s1", Owner: "Squad A", Dependencies: []string{"e1"}},
		},
		Metrics: []Metric{
			{ID: "m0", Name: "Active at 90 days", Owner: "Cy", Target: "60"},
			{ID: "m1", Name: "Connect rate", StrategyID: "s1", Owner: "Ana", Target: "70"},
			{ID: "m2", Name: "Weekly reconcile", StrategyID: "s2", Owner: "Ben", Target: "50"},
		},
	}
}

func TestHelpers(t *testing.T) {
	v := sample()
	if got := v.Strategy("s2"); got == nil || got.Name != "Retain" {
		t.Errorf("Strategy(s2) = %v", got)
	}
	if v.Strategy("nope") != nil {
		t.Error("Strategy(nope) should be nil")
	}
	var ids []string
	for _, e := range v.ExecutionFor("s1") {
		ids = append(ids, e.ID)
	}
	if !reflect.DeepEqual(ids, []string{"e1", "e3"}) {
		t.Errorf("ExecutionFor(s1) = %v, want document order [e1 e3]", ids)
	}
	if m := v.MetricsFor("s2"); len(m) != 1 || m[0].ID != "m2" {
		t.Errorf("MetricsFor(s2) = %v", m)
	}
	if m := v.VisionMetrics(); len(m) != 1 || m[0].ID != "m0" {
		t.Errorf("VisionMetrics() = %v", m)
	}
}

func TestStatement(t *testing.T) {
	v := &VSEM{Vision: "V", Mission: "M"}
	if v.Statement() != "V" {
		t.Errorf("Statement() = %q, want the vision", v.Statement())
	}
	v.Vision = ""
	if v.Statement() != "M" {
		t.Errorf("Statement() = %q, want the mission when there is no vision", v.Statement())
	}
}

func TestNew(t *testing.T) {
	v := New("id", "Name", "Vision")
	if v.Metadata.ID != "id" || v.Metadata.Name != "Name" || v.Metadata.Status != StatusDraft || v.Vision != "Vision" {
		t.Errorf("New() = %+v %+v", v, v.Metadata)
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	v := sample()
	path := filepath.Join(t.TempDir(), "plan.vsem.json")
	if err := v.WriteFile(path); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, v) {
		t.Errorf("round trip changed the VSEM:\n got %+v\nwant %+v", got, v)
	}
}

func TestParseRejectsInvalidJSON(t *testing.T) {
	if _, err := Parse([]byte("{")); err == nil {
		t.Error("Parse accepted invalid JSON")
	}
}
