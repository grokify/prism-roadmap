package assessment

import (
	"testing"

	"github.com/grokify/prism-roadmap/prioritization"
)

func TestResolveMoSCoWPriority(t *testing.T) {
	tests := []struct {
		name    string
		answers []ThresholdAnswer
		want    prioritization.MoSCoWPriority
	}{
		{
			name: "KTLO criterion satisfies Must",
			answers: []ThresholdAnswer{
				{LevelID: "must", Satisfied: true, CriterionMet: "KTLO", EvidenceIDs: []string{"EV-1"}},
			},
			want: prioritization.MoSCoWMustHave,
		},
		{
			name: "should satisfied when must is not",
			answers: []ThresholdAnswer{
				{LevelID: "must", Satisfied: false},
				{LevelID: "should", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
			},
			want: prioritization.MoSCoWShouldHave,
		},
		{
			name: "could satisfied when must and should are not",
			answers: []ThresholdAnswer{
				{LevelID: "must", Satisfied: false},
				{LevelID: "should", Satisfied: false},
				{LevelID: "could", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
			},
			want: prioritization.MoSCoWCouldHave,
		},
		{
			name:    "no criteria satisfied floors at wont",
			answers: nil,
			want:    prioritization.MoSCoWWontHave,
		},
		{
			name: "satisfied without evidence floors at wont",
			answers: []ThresholdAnswer{
				{LevelID: "must", Satisfied: true}, // no evidence
			},
			want: prioritization.MoSCoWWontHave,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ResolveMoSCoWPriority(tt.answers)
			if got != tt.want {
				t.Errorf("ResolveMoSCoWPriority() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestMoSCoWLadderLevelIDsParseable(t *testing.T) {
	for _, lvl := range MoSCoWLadder().Levels {
		if _, err := prioritization.ParseMoSCoWPriority(lvl.ID); err != nil {
			t.Errorf("level ID %q does not parse via prioritization.ParseMoSCoWPriority: %v", lvl.ID, err)
		}
	}
}
