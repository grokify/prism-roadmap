package assessment

import (
	"reflect"
	"testing"
)

func TestEvidenceIndexCitedBy(t *testing.T) {
	refs := []EvidenceRef{
		{EvidenceID: "EV-042", AssessmentID: "OA-002", QuestionID: "rice.impact.high"},
		{EvidenceID: "EV-042", AssessmentID: "OA-018", QuestionID: "moscow.must.ktlo"},
		{EvidenceID: "EV-042", AssessmentID: "OA-018", QuestionID: "moscow.must.compliance"},
		{EvidenceID: "EV-099", AssessmentID: "OA-018", QuestionID: "rice.reach"},
	}
	idx := NewEvidenceIndex(refs)

	got := idx.CitedBy("EV-042")
	want := []EvidenceRef{
		{EvidenceID: "EV-042", AssessmentID: "OA-002", QuestionID: "rice.impact.high"},
		{EvidenceID: "EV-042", AssessmentID: "OA-018", QuestionID: "moscow.must.compliance"},
		{EvidenceID: "EV-042", AssessmentID: "OA-018", QuestionID: "moscow.must.ktlo"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("CitedBy(EV-042) = %+v, want %+v", got, want)
	}

	if got := idx.CitationCount("EV-042"); got != 3 {
		t.Errorf("CitationCount(EV-042) = %d, want 3", got)
	}
	if got := idx.CitationCount("EV-999"); got != 0 {
		t.Errorf("CitationCount(EV-999) = %d, want 0", got)
	}
}

func TestEvidenceIndexAssessmentIDs(t *testing.T) {
	refs := []EvidenceRef{
		{EvidenceID: "EV-042", AssessmentID: "OA-018", QuestionID: "q1"},
		{EvidenceID: "EV-042", AssessmentID: "OA-018", QuestionID: "q2"},
		{EvidenceID: "EV-042", AssessmentID: "OA-002", QuestionID: "q1"},
	}
	idx := NewEvidenceIndex(refs)

	got := idx.AssessmentIDs("EV-042")
	want := []string{"OA-002", "OA-018"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("AssessmentIDs(EV-042) = %v, want %v (deduped, sorted)", got, want)
	}

	if got := idx.AssessmentIDs("EV-nonexistent"); got != nil {
		t.Errorf("AssessmentIDs(nonexistent) = %v, want nil", got)
	}
}
