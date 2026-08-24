package assessment

import (
	"testing"
	"time"

	"github.com/ProductBuildersHQ/compass-rice/rice"
)

func TestParseProfileAssignmentStatus(t *testing.T) {
	for _, s := range []string{"proposed", "Confirmed", " REJECTED "} {
		if _, err := ParseProfileAssignmentStatus(s); err != nil {
			t.Errorf("ParseProfileAssignmentStatus(%q) error = %v, want nil", s, err)
		}
	}
	if _, err := ParseProfileAssignmentStatus("bogus"); err == nil {
		t.Error("ParseProfileAssignmentStatus(\"bogus\") = nil error, want error")
	}
}

func TestProposeProfileAssignment(t *testing.T) {
	p := ProposeProfileAssignment("OS-001", "customer/b2b/v1", "primarily a retention play", "judge-session-42")
	if p.Status != ProfileAssignmentProposed {
		t.Errorf("Status = %q, want %q", p.Status, ProfileAssignmentProposed)
	}
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestProfileAssignmentConfirm(t *testing.T) {
	p := ProposeProfileAssignment("OS-001", "customer/b2b/v1", "primarily a retention play", "judge-session-42")
	confirmed := p.Confirm("pm@example.com", time.Now())
	if confirmed.Status != ProfileAssignmentConfirmed {
		t.Errorf("Status = %q, want %q", confirmed.Status, ProfileAssignmentConfirmed)
	}
	if confirmed.ConfirmedBy != "pm@example.com" {
		t.Errorf("ConfirmedBy = %q, want pm@example.com", confirmed.ConfirmedBy)
	}
	if err := confirmed.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
	// original is untouched
	if p.Status != ProfileAssignmentProposed {
		t.Errorf("original Status = %q after Confirm(), want unchanged %q", p.Status, ProfileAssignmentProposed)
	}
}

func TestProfileAssignmentReject(t *testing.T) {
	p := ProposeProfileAssignment("OS-001", "customer/b2b/v1", "primarily a retention play", "judge-session-42")
	rejected := p.Reject("pm@example.com", time.Now(), "actually a platform play, not customer")
	if rejected.Status != ProfileAssignmentRejected {
		t.Errorf("Status = %q, want %q", rejected.Status, ProfileAssignmentRejected)
	}
	if err := rejected.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestProfileAssignmentValidateRequiresConfirmerForConfirmedStatus(t *testing.T) {
	p := ProfileAssignment{
		SpecID:     "OS-001",
		ProfileID:  "customer/b2b/v1",
		Rationale:  "r",
		ProposedBy: "judge",
		Status:     ProfileAssignmentConfirmed,
	}
	if err := p.Validate(); err == nil {
		t.Error("Validate() with confirmed status but no ConfirmedBy/ConfirmedAt = nil error, want error")
	}
}

func TestProfileAssignmentValidateInvalidProfileID(t *testing.T) {
	p := ProposeProfileAssignment("OS-001", "bogus/v1", "r", "judge")
	if err := p.Validate(); err == nil {
		t.Error("Validate() with invalid profileId = nil error, want error")
	}
}

func TestProfileAssignmentValidateInvalidSecondary(t *testing.T) {
	p := ProposeProfileAssignment("OS-001", "customer/b2b/v1", "r", "judge")
	p.Secondary = []rice.Profile{"not-a-real-profile"}
	if err := p.Validate(); err == nil {
		t.Error("Validate() with invalid secondary profile = nil error, want error")
	}
}

func TestProfileAssignmentValidateValidSecondary(t *testing.T) {
	p := ProposeProfileAssignment("OS-001", "platform/internal/v1", "leverage across teams", "judge")
	p.Secondary = []rice.Profile{rice.ProfileCustomer, rice.ProfileRisk}
	if err := p.Validate(); err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
}

func TestProfileAssignmentValidateMissingRequiredFields(t *testing.T) {
	cases := []ProfileAssignment{
		{ProfileID: "customer/b2b/v1", Rationale: "r", ProposedBy: "j", Status: ProfileAssignmentProposed},        // no SpecID
		{SpecID: "OS-001", ProfileID: "customer/b2b/v1", ProposedBy: "j", Status: ProfileAssignmentProposed},      // no Rationale
		{SpecID: "OS-001", ProfileID: "customer/b2b/v1", Rationale: "r", Status: ProfileAssignmentProposed},       // no ProposedBy
		{SpecID: "OS-001", ProfileID: "customer/b2b/v1", Rationale: "r", ProposedBy: "j", Status: "not-a-status"}, // bad status
	}
	for i, c := range cases {
		if err := c.Validate(); err == nil {
			t.Errorf("case %d: Validate() = nil error, want error", i)
		}
	}
}
