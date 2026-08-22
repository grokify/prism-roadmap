package assessment

import (
	"testing"
	"time"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/prioritization"
	"github.com/plexusone/structured-evaluation/rubric"
)

func TestRefFromSpec(t *testing.T) {
	spec := canvas.NewOpportunitySpec("OPP-001", "Unified Authorization Platform")
	ref := RefFromSpec(spec)
	if ref.SpecID != "OPP-001" {
		t.Errorf("SpecID = %q, want OPP-001", ref.SpecID)
	}
	if ref.RMIID != "" {
		t.Errorf("RMIID = %q, want empty", ref.RMIID)
	}
}

func TestNewOpportunityAssessment(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	ref := OpportunityRef{SpecID: "OPP-001", RMIID: "RMI-MYREPO-007"}

	a := NewOpportunityAssessment("OA-001", ref, "Unified Authorization Platform", now)

	if a.ID != "OA-001" {
		t.Errorf("ID = %q, want OA-001", a.ID)
	}
	if a.Opportunity != ref {
		t.Errorf("Opportunity = %+v, want %+v", a.Opportunity, ref)
	}
	if a.Cycle.Number != 1 {
		t.Errorf("Cycle.Number = %d, want 1", a.Cycle.Number)
	}
	if !a.Cycle.AssessedAt.Equal(now) {
		t.Errorf("Cycle.AssessedAt = %v, want %v", a.Cycle.AssessedAt, now)
	}
	if a.Cycle.SupersedesID != "" {
		t.Errorf("Cycle.SupersedesID = %q, want empty for first cycle", a.Cycle.SupersedesID)
	}
	if !a.Cycle.Current {
		t.Error("Cycle.Current = false, want true for a new assessment")
	}
	if a.Judge != nil {
		t.Error("Judge should be nil until a rubric has been run")
	}
}

func TestNextCycle(t *testing.T) {
	t1 := time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC)
	t2 := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	ref := OpportunityRef{SpecID: "OPP-001"}

	first := NewOpportunityAssessment("OA-001", ref, "Unified Authorization Platform", t1)
	first.Judge = rubric.NewJudgeMetadata("claude-sonnet-5").WithRubric("platform-rice", "v1.0")

	second := first.NextCycle("OA-002", t2)

	if second.Cycle.Number != 2 {
		t.Errorf("Cycle.Number = %d, want 2", second.Cycle.Number)
	}
	if second.Cycle.SupersedesID != "OA-001" {
		t.Errorf("Cycle.SupersedesID = %q, want OA-001", second.Cycle.SupersedesID)
	}
	if !second.Cycle.Current {
		t.Error("new cycle should be current")
	}
	if second.Opportunity != ref {
		t.Errorf("Opportunity should carry forward: got %+v, want %+v", second.Opportunity, ref)
	}
	if second.Title != first.Title {
		t.Errorf("Title should carry forward: got %q, want %q", second.Title, first.Title)
	}
	if second.Judge != nil {
		t.Error("Judge should NOT carry forward — each cycle gets its own judge run")
	}

	first.MarkSuperseded()
	if first.Cycle.Current {
		t.Error("MarkSuperseded should flip Current to false")
	}
	if !second.Cycle.Current {
		t.Error("MarkSuperseded on the previous assessment should not affect the new one")
	}
}

func TestOpportunityAssessmentValidate(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	valid := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)

	tests := []struct {
		name    string
		mutate  func(*OpportunityAssessment)
		wantErr bool
	}{
		{"valid", func(a *OpportunityAssessment) {}, false},
		{"missing id", func(a *OpportunityAssessment) { a.ID = "" }, true},
		{"missing specId", func(a *OpportunityAssessment) { a.Opportunity.SpecID = "" }, true},
		{"missing title", func(a *OpportunityAssessment) { a.Title = "" }, true},
		{"zero cycle number", func(a *OpportunityAssessment) { a.Cycle.Number = 0 }, true},
		{"zero assessedAt", func(a *OpportunityAssessment) { a.Cycle.AssessedAt = time.Time{} }, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := *valid
			tt.mutate(&a)
			err := a.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOpportunityAssessmentMoSCoWDefaultsToWont(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)
	if got := a.MoSCoW(); got != prioritization.MoSCoWWontHave {
		t.Errorf("MoSCoW() = %q, want wontHave when no answers recorded", got)
	}
}

func TestOpportunityAssessmentMoSCoWResolvesFromAnswers(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)
	a.MoSCoWAnswers = []ThresholdAnswer{
		{LevelID: "must", Satisfied: true, CriterionMet: "KTLO", EvidenceIDs: []string{"EV-1"}},
	}
	if got := a.MoSCoW(); got != prioritization.MoSCoWMustHave {
		t.Errorf("MoSCoW() = %q, want mustHave", got)
	}
}

func TestOpportunityAssessmentToRankInputNoRICE(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)
	a.MoSCoWAnswers = []ThresholdAnswer{
		{LevelID: "must", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
	}

	in := a.ToRankInput()
	if in.AssessmentID != "OA-001" || in.Title != "Title" {
		t.Errorf("ToRankInput() identity = %+v", in)
	}
	if in.MoSCoW != prioritization.MoSCoWMustHave {
		t.Errorf("MoSCoW = %q, want mustHave", in.MoSCoW)
	}
	if in.RICE.Computable {
		t.Error("expected RICE.Computable = false when a.RICE is nil")
	}
}

func TestOpportunityAssessmentToRankInputWithRICE(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)
	a.MoSCoWAnswers = []ThresholdAnswer{
		{LevelID: "should", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
	}
	a.RICE = &RICEAssessment{
		Reach: Reach{Fraction: 0.5, EvidenceIDs: []string{"EV-2"}},
		ImpactAnswers: []ThresholdAnswer{
			{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-3"}},
		},
		ConfidenceAnswers: []ThresholdAnswer{
			{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-4"}},
		},
		Effort: EffortEstimate{
			Expected: 10,
			Gate: EstimabilityGate{
				ScopeDefined: true, ImplementationIdentified: true, DependenciesIdentified: true,
				TestingIdentified: true, DeploymentIdentified: true,
			},
		},
	}

	in := a.ToRankInput()
	if !in.RICE.Computable {
		t.Fatalf("expected RICE.Computable = true, got reason: %s", in.RICE.Reason)
	}
	// (0.5 * 2.0 * 1.0) / 10 = 0.10
	if in.RICE.Score != 0.10 {
		t.Errorf("RICE.Score = %v, want 0.10", in.RICE.Score)
	}
	if in.MoSCoW != prioritization.MoSCoWShouldHave {
		t.Errorf("MoSCoW = %q, want shouldHave", in.MoSCoW)
	}
}

func TestOpportunityAssessmentDimensionsContributionsCapabilities(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)

	a.Dimensions = []DimensionAssignment{
		NewDimensionAssignment(KanoDimension(), nil),
	}
	a.Contributions = []OKRContribution{
		{ObjectiveID: "OBJ-3", Strength: ContributionHigh},
	}
	a.Capabilities = []CapabilityReference{
		{CapabilityID: "authorization", Relation: CapabilityEnables},
	}

	if len(a.Dimensions) != 1 || a.Dimensions[0].DimensionID != "kano" {
		t.Errorf("Dimensions = %+v", a.Dimensions)
	}
	if len(a.Contributions) != 1 || a.Contributions[0].ObjectiveID != "OBJ-3" {
		t.Errorf("Contributions = %+v", a.Contributions)
	}
	if len(a.Capabilities) != 1 || a.Capabilities[0].CapabilityID != "authorization" {
		t.Errorf("Capabilities = %+v", a.Capabilities)
	}
}

func TestNextCycleDoesNotCarryForwardPrioritizationOrLinks(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	first := NewOpportunityAssessment("OA-001", OpportunityRef{SpecID: "OPP-001"}, "Title", now)
	first.MoSCoWAnswers = []ThresholdAnswer{{LevelID: "must", Satisfied: true, EvidenceIDs: []string{"EV-1"}}}
	first.RICE = &RICEAssessment{Reach: Reach{Fraction: 0.5, EvidenceIDs: []string{"EV-2"}}}
	first.Dimensions = []DimensionAssignment{NewDimensionAssignment(KanoDimension(), nil)}
	first.Contributions = []OKRContribution{{ObjectiveID: "OBJ-1", Strength: ContributionLow}}
	first.Capabilities = []CapabilityReference{{CapabilityID: "x", Relation: CapabilityEnables}}

	second := first.NextCycle("OA-002", now.AddDate(0, 3, 0))

	if second.MoSCoWAnswers != nil || second.RICE != nil || second.Dimensions != nil ||
		second.Contributions != nil || second.Capabilities != nil {
		t.Errorf("NextCycle() carried forward stale judgments: %+v", second)
	}
}

func TestEvidenceReferencesEmpty(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	if got := a.EvidenceReferences(); len(got) != 0 {
		t.Errorf("EvidenceReferences() = %v, want empty", got)
	}
}

func TestEvidenceReferencesAllSources(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-42", OpportunityRef{SpecID: "OPP-1"}, "Title", now)

	a.MoSCoWAnswers = []ThresholdAnswer{
		{LevelID: "must", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
	}
	a.RICE = &RICEAssessment{
		Reach: Reach{Fraction: 0.5, EvidenceIDs: []string{"EV-2"}},
		ImpactAnswers: []ThresholdAnswer{
			{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-3"}},
		},
		ConfidenceAnswers: []ThresholdAnswer{
			{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-4"}},
		},
	}
	a.Dimensions = []DimensionAssignment{
		{DimensionID: "kano", Answers: []DimensionAnswer{
			{OptionID: "must_be", QuestionID: "expected", Answer: true, EvidenceIDs: []string{"EV-5"}},
		}},
	}
	a.Contributions = []OKRContribution{
		{ObjectiveID: "OBJ-1", KeyResultID: "KR-1.1", Strength: ContributionHigh, EvidenceIDs: []string{"EV-6"}},
	}

	refs := a.EvidenceReferences()
	if len(refs) != 6 {
		t.Fatalf("EvidenceReferences() = %+v, want 6 entries", refs)
	}

	byEvidence := make(map[string]EvidenceRef)
	for _, r := range refs {
		if r.AssessmentID != "OA-42" {
			t.Errorf("ref %+v has wrong AssessmentID", r)
		}
		byEvidence[r.EvidenceID] = r
	}

	wantQuestionIDs := map[string]string{
		"EV-1": "moscow.must",
		"EV-2": "rice.reach",
		"EV-3": "rice.impact.high",
		"EV-4": "rice.confidence.medium",
		"EV-5": "dimension.kano.expected",
		"EV-6": "okr.OBJ-1.KR-1.1",
	}
	for evID, wantQID := range wantQuestionIDs {
		ref, ok := byEvidence[evID]
		if !ok {
			t.Errorf("missing ref for %s", evID)
			continue
		}
		if ref.QuestionID != wantQID {
			t.Errorf("ref for %s QuestionID = %q, want %q", evID, ref.QuestionID, wantQID)
		}
	}
}

func TestEvidenceReferencesOKRWithoutKeyResult(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	a.Contributions = []OKRContribution{
		{ObjectiveID: "OBJ-1", Strength: ContributionLow, EvidenceIDs: []string{"EV-1"}},
	}
	refs := a.EvidenceReferences()
	if len(refs) != 1 || refs[0].QuestionID != "okr.OBJ-1" {
		t.Errorf("EvidenceReferences() = %+v, want QuestionID okr.OBJ-1 (no KR suffix)", refs)
	}
}

func TestEvidenceReferencesMultipleEvidenceIDsPerAnswer(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	a := NewOpportunityAssessment("OA-1", OpportunityRef{SpecID: "OPP-1"}, "Title", now)
	a.MoSCoWAnswers = []ThresholdAnswer{
		{LevelID: "must", Satisfied: true, EvidenceIDs: []string{"EV-1", "EV-2"}},
	}
	refs := a.EvidenceReferences()
	if len(refs) != 2 {
		t.Fatalf("EvidenceReferences() = %+v, want 2 entries (one per evidence ID)", refs)
	}
	for _, r := range refs {
		if r.QuestionID != "moscow.must" {
			t.Errorf("ref %+v QuestionID mismatch", r)
		}
	}
}
