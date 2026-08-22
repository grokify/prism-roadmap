package assessment

import (
	"testing"

	core "github.com/grokify/prism-core"
)

func TestCapabilityRelationConstantsMatchPrismCore(t *testing.T) {
	tests := []struct {
		local string
		want  string
	}{
		{CapabilityEnables, core.CapabilityRelationEnables},
		{CapabilityImproves, core.CapabilityRelationImproves},
		{CapabilityDependsOn, core.CapabilityRelationDependsOn},
	}
	for _, tt := range tests {
		if tt.local != tt.want {
			t.Errorf("local constant %q does not match prism-core value %q", tt.local, tt.want)
		}
		if !core.ValidCapabilityRelation(tt.local) {
			t.Errorf("core.ValidCapabilityRelation(%q) = false, want true", tt.local)
		}
	}
}

func TestValidateCapabilityReference(t *testing.T) {
	valid := CapabilityReference{CapabilityID: "authorization", Relation: CapabilityEnables}
	if err := ValidateCapabilityReference(valid); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if err := ValidateCapabilityReference(CapabilityReference{Relation: CapabilityEnables}); err == nil {
		t.Error("expected error for missing capabilityId")
	}
	if err := ValidateCapabilityReference(CapabilityReference{CapabilityID: "authorization"}); err == nil {
		t.Error("expected error for missing/invalid relation")
	}
}

func TestCapabilityReferenceIsTypeIdenticalToPrismCore(t *testing.T) {
	// Compile-time proof the alias is a true type identity, not a
	// look-alike: a core.CapabilityRef slice must be directly assignable
	// to a []CapabilityReference with no conversion.
	var fromCore []core.CapabilityRef = []core.CapabilityRef{
		{CapabilityID: "authorization", Relation: core.CapabilityRelationEnables},
	}
	var viaAssessment []CapabilityReference = fromCore
	if viaAssessment[0].CapabilityID != "authorization" {
		t.Errorf("unexpected value after direct assignment: %+v", viaAssessment[0])
	}
}
