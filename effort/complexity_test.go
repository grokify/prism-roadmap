package effort

import "testing"

func TestComplexityFactorsScore(t *testing.T) {
	tests := []struct {
		name       string
		complexity ComplexityFactors
		want       float64
	}{
		{
			name:       "no complexity",
			complexity: ComplexityFactors{},
			want:       0,
		},
		{
			name: "new architecture only",
			complexity: ComplexityFactors{
				NewArchitecture: true,
			},
			want: 2.0,
		},
		{
			name: "all flags",
			complexity: ComplexityFactors{
				NewArchitecture: true, // 2.0
				NewDesignUX:     true, // 1.0
				NewBillingSKU:   true, // 1.5
				NewIntegration:  true, // 1.0
				NewCompliance:   true, // 1.0
				NewDataModel:    true, // 1.0
			},
			want: 7.5,
		},
		{
			name: "with dependencies",
			complexity: ComplexityFactors{
				Dependencies: []Dependency{
					{TeamID: "team-1"},
					{TeamID: "team-2"},
				},
			},
			want: 1.0, // 2 deps × 0.5
		},
		{
			name: "mixed",
			complexity: ComplexityFactors{
				NewArchitecture: true, // 2.0
				NewDesignUX:     true, // 1.0
				Dependencies: []Dependency{
					{TeamID: "team-1"},
				}, // 0.5
			},
			want: 3.5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.complexity.ComplexityScore()
			if got != tt.want {
				t.Errorf("ComplexityScore() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestComplexityFactorsHasComplexity(t *testing.T) {
	empty := ComplexityFactors{}
	if empty.HasComplexity() {
		t.Error("Empty ComplexityFactors HasComplexity() = true, want false")
	}

	withArch := ComplexityFactors{NewArchitecture: true}
	if !withArch.HasComplexity() {
		t.Error("ComplexityFactors with NewArchitecture HasComplexity() = false, want true")
	}

	withDeps := ComplexityFactors{
		Dependencies: []Dependency{{TeamID: "team-1"}},
	}
	if !withDeps.HasComplexity() {
		t.Error("ComplexityFactors with dependencies HasComplexity() = false, want true")
	}
}

func TestComplexityFactorsBlockingDependencies(t *testing.T) {
	c := ComplexityFactors{
		Dependencies: []Dependency{
			{TeamID: "team-1", Type: DependencyTypeBlocking},
			{TeamID: "team-2", Type: DependencyTypeInformational},
			{TeamID: "team-3", Type: DependencyTypeBlocking},
		},
	}

	blocking := c.BlockingDependencies()
	if len(blocking) != 2 {
		t.Errorf("BlockingDependencies() returned %d, want 2", len(blocking))
	}
}

func TestComplexityFactorsTotalDependencyDays(t *testing.T) {
	c := ComplexityFactors{
		Dependencies: []Dependency{
			{TeamID: "team-1", EstimatedDays: 5},
			{TeamID: "team-2", EstimatedDays: 10},
			{TeamID: "team-3", EstimatedDays: 3},
		},
	}

	total := c.TotalDependencyDays()
	if total != 18 {
		t.Errorf("TotalDependencyDays() = %d, want 18", total)
	}
}

func TestDependencyIsBlocking(t *testing.T) {
	tests := []struct {
		name string
		dep  Dependency
		want bool
	}{
		{
			name: "blocking pending",
			dep:  Dependency{Type: DependencyTypeBlocking, Status: DependencyStatusPending},
			want: true,
		},
		{
			name: "blocking resolved",
			dep:  Dependency{Type: DependencyTypeBlocking, Status: DependencyStatusResolved},
			want: false,
		},
		{
			name: "informational pending",
			dep:  Dependency{Type: DependencyTypeInformational, Status: DependencyStatusPending},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dep.IsBlocking()
			if got != tt.want {
				t.Errorf("IsBlocking() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewBlockingDependency(t *testing.T) {
	dep := NewBlockingDependency("team-1", "Platform Team", 5)

	if dep.TeamID != "team-1" {
		t.Errorf("TeamID = %s, want team-1", dep.TeamID)
	}
	if dep.TeamName != "Platform Team" {
		t.Errorf("TeamName = %s, want Platform Team", dep.TeamName)
	}
	if dep.Type != DependencyTypeBlocking {
		t.Errorf("Type = %s, want blocking", dep.Type)
	}
	if dep.EstimatedDays != 5 {
		t.Errorf("EstimatedDays = %d, want 5", dep.EstimatedDays)
	}
	if dep.Status != DependencyStatusPending {
		t.Errorf("Status = %s, want pending", dep.Status)
	}
}
