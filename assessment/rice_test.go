package assessment

import (
	"testing"

	"github.com/grokify/prism-roadmap/prioritization"
)

func TestResolveImpact(t *testing.T) {
	tests := []struct {
		name    string
		answers []ThresholdAnswer
		want    prioritization.ImpactLevel
		wantOk  bool
	}{
		{
			name: "high satisfied",
			answers: []ThresholdAnswer{
				{LevelID: "massive", Satisfied: false},
				{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
			},
			want: prioritization.ImpactHigh, wantOk: true,
		},
		{
			name:    "nothing satisfied",
			answers: nil,
			want:    "", wantOk: false,
		},
		{
			name: "massive without evidence does not count",
			answers: []ThresholdAnswer{
				{LevelID: "massive", Satisfied: true},
				{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
			},
			want: prioritization.ImpactHigh, wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := ResolveImpact(tt.answers)
			if ok != tt.wantOk || got != tt.want {
				t.Errorf("ResolveImpact() = (%q, %v), want (%q, %v)", got, ok, tt.want, tt.wantOk)
			}
		})
	}
}

func TestImpactLadderLevelIDsParseable(t *testing.T) {
	for _, lvl := range ImpactLadder().Levels {
		if _, err := prioritization.ParseImpactLevel(lvl.ID); err != nil {
			t.Errorf("level ID %q does not parse via prioritization.ParseImpactLevel: %v", lvl.ID, err)
		}
	}
}

func TestResolveConfidence(t *testing.T) {
	answers := []ThresholdAnswer{
		{LevelID: "high", Satisfied: false},
		{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-1"}},
	}
	got, ok := ResolveConfidence(answers)
	if !ok || got != prioritization.ConfidenceMedium {
		t.Errorf("ResolveConfidence() = (%q, %v), want (medium, true)", got, ok)
	}
}

func TestConfidenceLadderLevelIDsParseable(t *testing.T) {
	for _, lvl := range ConfidenceLadder().Levels {
		if _, err := prioritization.ParseConfidenceLevel(lvl.ID); err != nil {
			t.Errorf("level ID %q does not parse via prioritization.ParseConfidenceLevel: %v", lvl.ID, err)
		}
	}
}

func TestReachValidate(t *testing.T) {
	tests := []struct {
		name    string
		r       Reach
		wantErr bool
	}{
		{"zero reach needs no evidence", Reach{Fraction: 0}, false},
		{"nonzero reach with evidence", Reach{Fraction: 0.6, EvidenceIDs: []string{"EV-1"}}, false},
		{"nonzero reach without evidence", Reach{Fraction: 0.6}, true},
		{"fraction above 1", Reach{Fraction: 1.5, EvidenceIDs: []string{"EV-1"}}, true},
		{"negative fraction", Reach{Fraction: -0.1}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.r.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEstimabilityGate(t *testing.T) {
	passed := EstimabilityGate{
		ScopeDefined: true, ImplementationIdentified: true, DependenciesIdentified: true,
		TestingIdentified: true, DeploymentIdentified: true,
	}
	if !passed.Passed() {
		t.Error("expected fully-checked gate to pass")
	}
	if got := passed.MissingChecks(); len(got) != 0 {
		t.Errorf("MissingChecks() = %v, want empty", got)
	}

	partial := EstimabilityGate{ScopeDefined: true, ImplementationIdentified: true}
	if partial.Passed() {
		t.Error("expected partially-checked gate to fail")
	}
	missing := partial.MissingChecks()
	if len(missing) != 3 {
		t.Errorf("MissingChecks() = %v, want 3 missing items", missing)
	}
}

func TestEffortEstimateValidate(t *testing.T) {
	gate := EstimabilityGate{
		ScopeDefined: true, ImplementationIdentified: true, DependenciesIdentified: true,
		TestingIdentified: true, DeploymentIdentified: true,
	}

	if err := (EffortEstimate{Expected: 10, Gate: gate}).Validate(); err != nil {
		t.Errorf("unexpected error for valid estimate: %v", err)
	}
	if err := (EffortEstimate{Expected: 10, Gate: EstimabilityGate{}}).Validate(); err == nil {
		t.Error("expected error for unpassed estimability gate")
	}
	if err := (EffortEstimate{Expected: 0, Gate: gate}).Validate(); err == nil {
		t.Error("expected error for zero expected person-days")
	}
}

func TestEffortEstimateTotalPersonDays(t *testing.T) {
	e := EffortEstimate{
		Components: []EffortComponent{
			{Name: "Application", PersonDays: 5},
			{Name: "Testing", PersonDays: 2.5},
		},
	}
	if got := e.TotalPersonDays(); got != 7.5 {
		t.Errorf("TotalPersonDays() = %v, want 7.5", got)
	}
}

func passingGate() EstimabilityGate {
	return EstimabilityGate{
		ScopeDefined: true, ImplementationIdentified: true, DependenciesIdentified: true,
		TestingIdentified: true, DeploymentIdentified: true,
	}
}

func TestComputeRICEHappyPath(t *testing.T) {
	a := RICEAssessment{
		Reach: Reach{Fraction: 0.6, EvidenceIDs: []string{"EV-1"}},
		ImpactAnswers: []ThresholdAnswer{
			{LevelID: "massive", Satisfied: false},
			{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-2"}},
		},
		ConfidenceAnswers: []ThresholdAnswer{
			{LevelID: "high", Satisfied: false},
			{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-3"}},
		},
		Effort: EffortEstimate{Expected: 20, Gate: passingGate()},
	}
	result := ComputeRICE(a)
	if !result.Computable {
		t.Fatalf("expected computable result, got reason: %s", result.Reason)
	}
	// (0.6 * 2.0 * 0.8) / 20 = 0.048
	want := 0.048
	if diff := result.Score - want; diff > 1e-9 || diff < -1e-9 {
		t.Errorf("Score = %v, want %v", result.Score, want)
	}
	if result.Impact != prioritization.ImpactHigh {
		t.Errorf("Impact = %q, want high", result.Impact)
	}
	if result.Confidence != prioritization.ConfidenceMedium {
		t.Errorf("Confidence = %q, want medium", result.Confidence)
	}
}

func TestComputeRICEUnresolvedReach(t *testing.T) {
	a := RICEAssessment{
		Reach: Reach{Fraction: 0.6}, // no evidence
	}
	result := ComputeRICE(a)
	if result.Computable {
		t.Error("expected uncomputable result for unsupported reach")
	}
	if result.Score != 0 {
		t.Errorf("Score = %v, want 0 for uncomputable result (not silently defaulted)", result.Score)
	}
}

func TestComputeRICEUnresolvedImpact(t *testing.T) {
	a := RICEAssessment{
		Reach:         Reach{Fraction: 0.6, EvidenceIDs: []string{"EV-1"}},
		ImpactAnswers: nil, // no impact threshold satisfied
	}
	result := ComputeRICE(a)
	if result.Computable {
		t.Error("expected uncomputable result when no Impact threshold is satisfied")
	}
}

func TestComputeRICEInsufficientConfidenceEvidence(t *testing.T) {
	a := RICEAssessment{
		Reach: Reach{Fraction: 0.6, EvidenceIDs: []string{"EV-1"}},
		ImpactAnswers: []ThresholdAnswer{
			{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-2"}},
		},
		ConfidenceInsufficientEvidence: true,
	}
	result := ComputeRICE(a)
	if result.Computable {
		t.Error("expected uncomputable result when confidence is marked insufficient evidence")
	}
	if result.Impact != prioritization.ImpactHigh {
		t.Errorf("Impact should still be resolved and returned even when Confidence blocks computation, got %q", result.Impact)
	}
}

func TestComputeRICEFailedEstimabilityGate(t *testing.T) {
	a := RICEAssessment{
		Reach: Reach{Fraction: 0.6, EvidenceIDs: []string{"EV-1"}},
		ImpactAnswers: []ThresholdAnswer{
			{LevelID: "high", Satisfied: true, EvidenceIDs: []string{"EV-2"}},
		},
		ConfidenceAnswers: []ThresholdAnswer{
			{LevelID: "medium", Satisfied: true, EvidenceIDs: []string{"EV-3"}},
		},
		Effort: EffortEstimate{Expected: 20, Gate: EstimabilityGate{}}, // unpassed gate
	}
	result := ComputeRICE(a)
	if result.Computable {
		t.Error("expected uncomputable result when the estimability gate fails")
	}
	if result.Impact != prioritization.ImpactHigh || result.Confidence != prioritization.ConfidenceMedium {
		t.Errorf("Impact/Confidence should still resolve even when Effort blocks computation, got impact=%q confidence=%q", result.Impact, result.Confidence)
	}
}
