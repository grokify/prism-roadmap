package assessment

import "sort"

// EvidenceRef is a citation of an Evidence record by an assessment rubric
// answer — e.g. "assessment OA-018, question moscow.must.ktlo cites EV-042".
// The persistence layer (omniroadmap) is the source of these refs (they are
// discovered by walking assessment rubric answers); EvidenceIndex is a
// query helper over a set of them, not itself a store.
type EvidenceRef struct {
	EvidenceID   string `json:"evidenceId"`
	AssessmentID string `json:"assessmentId"`
	QuestionID   string `json:"questionId,omitempty"`
}

// EvidenceIndex answers "which assessments/questions cite this evidence" in
// memory, given a flat list of refs. Built fresh from the current ref set —
// it does not persist anything itself (prism-roadmap PRD FR5).
type EvidenceIndex struct {
	byEvidence map[string][]EvidenceRef
}

// NewEvidenceIndex builds an index over the given refs.
func NewEvidenceIndex(refs []EvidenceRef) *EvidenceIndex {
	idx := &EvidenceIndex{byEvidence: make(map[string][]EvidenceRef, len(refs))}
	for _, r := range refs {
		idx.byEvidence[r.EvidenceID] = append(idx.byEvidence[r.EvidenceID], r)
	}
	return idx
}

// CitedBy returns every ref citing the given evidence ID, ordered by
// assessment ID then question ID for stable output.
func (idx *EvidenceIndex) CitedBy(evidenceID string) []EvidenceRef {
	refs := append([]EvidenceRef(nil), idx.byEvidence[evidenceID]...)
	sort.Slice(refs, func(i, j int) bool {
		if refs[i].AssessmentID != refs[j].AssessmentID {
			return refs[i].AssessmentID < refs[j].AssessmentID
		}
		return refs[i].QuestionID < refs[j].QuestionID
	})
	return refs
}

// CitationCount returns how many refs cite the given evidence ID.
func (idx *EvidenceIndex) CitationCount(evidenceID string) int {
	return len(idx.byEvidence[evidenceID])
}

// AssessmentIDs returns the distinct assessment IDs citing the given
// evidence ID, sorted.
func (idx *EvidenceIndex) AssessmentIDs(evidenceID string) []string {
	seen := make(map[string]bool)
	var ids []string
	for _, r := range idx.byEvidence[evidenceID] {
		if !seen[r.AssessmentID] {
			seen[r.AssessmentID] = true
			ids = append(ids, r.AssessmentID)
		}
	}
	sort.Strings(ids)
	return ids
}
