package rubrics

import "testing"

func TestLoadRichRubric(t *testing.T) {
	r, err := LoadRubric("opportunity-spec")
	if err != nil {
		t.Fatalf("LoadRubric: %v", err)
	}
	if r.Name == "" {
		t.Error("rubric name is empty")
	}
	if len(r.Categories) == 0 {
		t.Fatal("rubric has no categories")
	}

	// The rich format carries weighted criteria with indicators.
	var sawCriteria, sawIndicators bool
	for _, c := range r.Categories {
		if c.Weight == 0 {
			t.Errorf("category %q has zero weight", c.Name)
		}
		for _, cr := range c.Criteria {
			sawCriteria = true
			if cr.Pass.Description == "" {
				t.Errorf("criterion %q has no pass description", cr.Name)
			}
			if len(cr.Pass.Indicators) > 0 {
				sawIndicators = true
			}
		}
	}
	if !sawCriteria {
		t.Error("expected weighted criteria in the rich rubric")
	}
	if !sawIndicators {
		t.Error("expected pass-level indicators in the rich rubric")
	}
}

func TestRichRubricScoreThresholds(t *testing.T) {
	r, err := LoadRubric("opportunity-spec")
	if err != nil {
		t.Fatalf("LoadRubric: %v", err)
	}
	if r.PassCriteria.ScoreThresholds == nil {
		t.Fatal("score thresholds not parsed")
	}
	if r.PassCriteria.ScoreThresholds.Pass != 80 || r.PassCriteria.ScoreThresholds.Partial != 60 {
		t.Errorf("thresholds = %+v, want pass 80 partial 60", r.PassCriteria.ScoreThresholds)
	}
}
