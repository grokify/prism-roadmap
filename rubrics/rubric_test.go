package rubrics

import (
	"strings"
	"testing"
)

// TestLoadAllRubrics guards every embedded rubric: each must parse into the
// structured-evaluation RubricSet with a name and at least one category.
func TestLoadAllRubrics(t *testing.T) {
	names := Names()
	if len(names) == 0 {
		t.Fatal("no embedded rubrics found")
	}
	for _, name := range names {
		r, err := LoadRubric(name)
		if err != nil {
			t.Errorf("LoadRubric(%q): %v", name, err)
			continue
		}
		if r.Name == "" {
			t.Errorf("rubric %q: empty name", name)
		}
		if len(r.Categories) == 0 {
			t.Errorf("rubric %q: no categories", name)
		}
	}
}

// TestV2MOMRubricWeights verifies the weight arithmetic of the v2mom rubric
// family: category weights sum to 100 and each category's criterion weights sum
// to the category weight, so weighted scoring is well-defined.
func TestV2MOMRubricWeights(t *testing.T) {
	for _, name := range Names() {
		if !strings.HasPrefix(name, "v2mom-") {
			continue
		}
		r, err := LoadRubric(name)
		if err != nil {
			t.Fatalf("LoadRubric(%q): %v", name, err)
		}
		var catSum float64
		for _, c := range r.Categories {
			catSum += c.Weight
			var critSum float64
			for _, cr := range c.Criteria {
				critSum += cr.Weight
			}
			if critSum != c.Weight {
				t.Errorf("rubric %q category %q: criterion weights sum to %v, want %v", name, c.Name, critSum, c.Weight)
			}
		}
		if catSum != 100 {
			t.Errorf("rubric %q: category weights sum to %v, want 100", name, catSum)
		}
	}
}

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
