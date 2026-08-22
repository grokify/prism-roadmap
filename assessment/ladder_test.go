package assessment

import "testing"

func testLadder() *Ladder {
	return &Ladder{
		ID:   "test",
		Name: "Test Ladder",
		Levels: []ThresholdLevel{
			{ID: "high", Label: "High", Criteria: []string{"crit-high-1", "crit-high-2"}},
			{ID: "medium", Label: "Medium", Criteria: []string{"crit-medium-1"}},
			{ID: "low", Label: "Low", Criteria: []string{"crit-low-1"}},
		},
	}
}

func TestLadderEvaluateTopDown(t *testing.T) {
	l := testLadder()
	answers := []ThresholdAnswer{
		{LevelID: "high", Satisfied: false},
		{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-001"}},
		{LevelID: "low", Satisfied: true, EvidenceIDs: []string{"EV-002"}},
	}
	level, answer, ok := l.Evaluate(answers)
	if !ok {
		t.Fatal("expected a level to be satisfied")
	}
	if level.ID != "medium" {
		t.Errorf("level.ID = %q, want medium (highest satisfied)", level.ID)
	}
	if answer.LevelID != "medium" {
		t.Errorf("answer.LevelID = %q, want medium", answer.LevelID)
	}
}

func TestLadderEvaluateNoneSatisfied(t *testing.T) {
	l := testLadder()
	level, answer, ok := l.Evaluate(nil)
	if ok || level != nil || answer != nil {
		t.Errorf("expected ok=false, level=nil, answer=nil for no answers; got ok=%v level=%v answer=%v", ok, level, answer)
	}
}

func TestLadderEvaluateSatisfiedWithoutEvidenceIsUnsupported(t *testing.T) {
	l := testLadder()
	answers := []ThresholdAnswer{
		{LevelID: "high", Satisfied: true}, // no EvidenceIDs — must not count
		{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-001"}},
	}
	level, _, ok := l.Evaluate(answers)
	if !ok {
		t.Fatal("expected medium to be satisfied")
	}
	if level.ID != "medium" {
		t.Errorf("level.ID = %q, want medium (high was unsupported and should be skipped)", level.ID)
	}
}

func TestLadderUnsupportedAnswers(t *testing.T) {
	l := testLadder()
	answers := []ThresholdAnswer{
		{LevelID: "high", Satisfied: true},                                  // unsupported
		{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-1"}}, // supported
		{LevelID: "low", Satisfied: false},                                  // not satisfied, not unsupported
	}
	unsupported := l.UnsupportedAnswers(answers)
	if len(unsupported) != 1 || unsupported[0].LevelID != "high" {
		t.Errorf("UnsupportedAnswers() = %+v, want exactly [high]", unsupported)
	}
}

func TestLadderLevelByID(t *testing.T) {
	l := testLadder()
	if lvl := l.LevelByID("medium"); lvl == nil || lvl.Label != "Medium" {
		t.Errorf("LevelByID(medium) = %+v", lvl)
	}
	if lvl := l.LevelByID("nonexistent"); lvl != nil {
		t.Errorf("LevelByID(nonexistent) = %+v, want nil", lvl)
	}
}
