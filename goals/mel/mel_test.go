package mel

import (
	"path/filepath"
	"reflect"
	"testing"
)

// sample returns a small plan that meets the strict options.
func sample() *Plan {
	return &Plan{
		Metadata: &Metadata{ID: "plan", Name: "Plan", Status: StatusActive},
		Purpose:  "Know whether onboarding works.",
		Results: []Result{
			{ID: "r-impact", Statement: "Owners save time", Level: LevelImpact},
			{ID: "r-long", Statement: "Accounts stay", Level: LevelLongTerm, ContributesTo: []string{"r-impact"}},
			{ID: "r-short", Statement: "Accounts connect a bank", Level: LevelShortTerm, ContributesTo: []string{"r-long"}},
		},
		Indicators: []Indicator{
			ind("i-hours", "r-impact"), ind("i-active", "r-long"), ind("i-connect", "r-short"),
		},
		Assumptions: []Assumption{{ID: "a-1", Statement: "Banks stay reliable", ResultIDs: []string{"r-short"}, Monitoring: "Status pages"}},
		Evaluations: []Evaluation{{
			ID: "e-1", Title: "Final", Type: EvalFinal, Questions: []string{"Did it work?"},
			Evaluator: EvaluatorExternal, StartDate: "2027-10-01", EndDate: "2027-12-31", Status: EvaluationPlanned,
			RelatedIDs: []string{"r-long", "q-1"},
		}},
		Learning: []LearningQuestion{{
			ID: "q-1", Question: "Why do owners stay?", Methods: []string{"Interviews"}, DueDate: "2027-06-30",
			Status: QuestionOpen, RelatedIDs: []string{"r-long", "a-1"},
		}},
		Reviews: []Review{{ID: "rv-1", Name: "Quarterly review", Frequency: "Quarterly", QuestionIDs: []string{"q-1"}}},
		Reports: []Report{{ID: "rp-1", Audience: "Board", Frequency: "Quarterly", IndicatorIDs: []string{"i-active"}}},
	}
}

func ind(id, result string) Indicator {
	return Indicator{
		ID: id, Name: id, ResultID: result, Type: TypeQuantitative, Direction: DirectionIncrease,
		Baseline: &Baseline{Value: "1"}, Targets: []Target{{Value: "2", Date: "2027-12-31"}},
		DataSource: "Analytics", CollectionMethod: "Query", Frequency: "Monthly", Owner: "Ana",
	}
}

func TestHelpers(t *testing.T) {
	p := sample()
	if r := p.Result("r-long"); r == nil || r.Statement != "Accounts stay" {
		t.Errorf("Result(r-long) = %v", r)
	}
	if p.Result("nope") != nil || p.Indicator("nope") != nil {
		t.Error("unknown IDs should give nil")
	}
	if got := p.IndicatorsFor("r-short"); len(got) != 1 || got[0].ID != "i-connect" {
		t.Errorf("IndicatorsFor(r-short) = %v", got)
	}
	if got := p.Contributors("r-long"); len(got) != 1 || got[0].ID != "r-short" {
		t.Errorf("Contributors(r-long) = %v", got)
	}
	if got := p.TopResults(); len(got) != 1 || got[0].ID != "r-impact" {
		t.Errorf("TopResults() = %v", got)
	}
}

func TestLevels(t *testing.T) {
	if LevelRank(LevelOutput) != 0 || LevelRank(LevelImpact) != 4 || LevelRank("nope") != -1 {
		t.Error("LevelRank is wrong")
	}
	if LevelLabel(LevelMediumTerm) != "Medium-term outcome" || LevelLabel("other") != "other" {
		t.Error("LevelLabel is wrong")
	}
}

func TestNew(t *testing.T) {
	p := New("id", "Name")
	if p.Metadata.ID != "id" || p.Metadata.Name != "Name" || p.Metadata.Status != StatusDraft {
		t.Errorf("New() = %+v", p.Metadata)
	}
}

func TestWriteReadRoundTrip(t *testing.T) {
	p := sample()
	path := filepath.Join(t.TempDir(), "plan.mel.json")
	if err := p.WriteFile(path); err != nil {
		t.Fatal(err)
	}
	got, err := ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, p) {
		t.Errorf("round trip changed the plan:\n got %+v\nwant %+v", got, p)
	}
	if _, err := Parse([]byte("{")); err == nil {
		t.Error("Parse accepted invalid JSON")
	}
}
