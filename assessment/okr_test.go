package assessment

import (
	"testing"

	"github.com/grokify/prism-roadmap/goals/okr"
)

func testOKRSet() *okr.OKRSet {
	return &okr.OKRSet{
		OKRs: []okr.OKR{
			{
				Objective: okr.Objective{
					ID:    "OBJ-3",
					Title: "Enterprise readiness",
					KeyResults: []okr.KeyResult{
						{ID: "KR-3.1", Title: "Enterprise RBAC adoption"},
						{ID: "KR-3.2", Title: "Implementation time"},
					},
				},
			},
		},
	}
}

func TestContributionStrengthIsValid(t *testing.T) {
	tests := []struct {
		s    ContributionStrength
		want bool
	}{
		{ContributionHigh, true},
		{ContributionMedium, true},
		{ContributionLow, true},
		{"", false},
		{"critical", false},
	}
	for _, tt := range tests {
		if got := tt.s.IsValid(); got != tt.want {
			t.Errorf("ContributionStrength(%q).IsValid() = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestOKRContributionValidate(t *testing.T) {
	valid := OKRContribution{ObjectiveID: "OBJ-3", Strength: ContributionHigh}
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := (OKRContribution{Strength: ContributionHigh}).Validate(); err == nil {
		t.Error("expected error for missing objectiveId")
	}
	if err := (OKRContribution{ObjectiveID: "OBJ-3"}).Validate(); err == nil {
		t.Error("expected error for missing/invalid strength")
	}
}

func TestOKRContributionValidateAgainst(t *testing.T) {
	set := testOKRSet()

	tests := []struct {
		name    string
		c       OKRContribution
		wantErr bool
	}{
		{"objective only, exists", OKRContribution{ObjectiveID: "OBJ-3", Strength: ContributionHigh}, false},
		{"objective + valid KR", OKRContribution{ObjectiveID: "OBJ-3", KeyResultID: "KR-3.1", Strength: ContributionMedium}, false},
		{"unknown objective", OKRContribution{ObjectiveID: "OBJ-999", Strength: ContributionHigh}, true},
		{"known objective, unknown KR", OKRContribution{ObjectiveID: "OBJ-3", KeyResultID: "KR-999", Strength: ContributionHigh}, true},
		{"invalid shape short-circuits", OKRContribution{ObjectiveID: "OBJ-3"}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.ValidateAgainst(set)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateAgainst() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOKRContributionValidateAgainstNilSet(t *testing.T) {
	c := OKRContribution{ObjectiveID: "OBJ-3", Strength: ContributionHigh}
	if err := c.ValidateAgainst(nil); err == nil {
		t.Error("expected error for nil OKRSet")
	}
}
