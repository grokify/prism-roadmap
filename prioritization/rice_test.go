package prioritization

import (
	"testing"
)

func TestRICEScore_Calculate(t *testing.T) {
	tests := []struct {
		name      string
		score     RICEScore
		wantScore float64
		wantClose bool // Allow for floating point comparison
	}{
		{
			name: "Standard calculation",
			score: RICEScore{
				FeatureID:  "feat-1",
				Reach:      1000,
				Impact:     ImpactHigh,     // 2x
				Confidence: ConfidenceHigh, // 100%
				Effort:     2,
			},
			wantScore: 1000, // (1000 * 2 * 1.0) / 2 = 1000
		},
		{
			name: "Low confidence reduces score",
			score: RICEScore{
				FeatureID:  "feat-2",
				Reach:      1000,
				Impact:     ImpactHigh,    // 2x
				Confidence: ConfidenceLow, // 50%
				Effort:     2,
			},
			wantScore: 500, // (1000 * 2 * 0.5) / 2 = 500
		},
		{
			name: "Massive impact multiplier",
			score: RICEScore{
				FeatureID:  "feat-3",
				Reach:      500,
				Impact:     ImpactMassive, // 3x
				Confidence: ConfidenceHigh,
				Effort:     1.5,
			},
			wantScore: 1000, // (500 * 3 * 1.0) / 1.5 = 1000
		},
		{
			name: "Zero effort returns zero",
			score: RICEScore{
				FeatureID:  "feat-4",
				Reach:      1000,
				Impact:     ImpactHigh,
				Confidence: ConfidenceHigh,
				Effort:     0,
			},
			wantScore: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.score.Calculate()
			if got != tt.wantScore {
				t.Errorf("Calculate() = %v, want %v", got, tt.wantScore)
			}
		})
	}
}

func TestRICEScoreSet_TopN(t *testing.T) {
	set := NewRICEScoreSet()

	set.Add(*NewRICEScore("low", 100, ImpactLow, ConfidenceLow, 5))       // Low score
	set.Add(*NewRICEScore("high", 1000, ImpactHigh, ConfidenceHigh, 1))   // High score
	set.Add(*NewRICEScore("mid", 500, ImpactMedium, ConfidenceMedium, 2)) // Mid score

	top2 := set.TopN(2)

	if len(top2) != 2 {
		t.Errorf("TopN(2) returned %d items, want 2", len(top2))
	}

	if top2[0].FeatureID != "high" {
		t.Errorf("First feature should be 'high', got %s", top2[0].FeatureID)
	}
}

func TestImpactLevel_Multiplier(t *testing.T) {
	tests := []struct {
		level ImpactLevel
		want  float64
	}{
		{ImpactMassive, 3.0},
		{ImpactHigh, 2.0},
		{ImpactMedium, 1.0},
		{ImpactLow, 0.5},
		{ImpactMinimal, 0.25},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			if got := tt.level.Multiplier(); got != tt.want {
				t.Errorf("Multiplier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseImpactLevel(t *testing.T) {
	tests := []struct {
		in      string
		want    ImpactLevel
		wantErr bool
	}{
		{"massive", ImpactMassive, false},
		{"HIGH", ImpactHigh, false}, // case-insensitive
		{"  medium  ", ImpactMedium, false},
		{"bogus", "", true},
		{"", "", true},
	}
	for _, tt := range tests {
		got, err := ParseImpactLevel(tt.in)
		if tt.wantErr {
			if err == nil {
				t.Errorf("ParseImpactLevel(%q) expected error, got nil", tt.in)
			}
			continue
		}
		if err != nil {
			t.Errorf("ParseImpactLevel(%q) unexpected error: %v", tt.in, err)
		}
		if got != tt.want {
			t.Errorf("ParseImpactLevel(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseConfidenceLevel(t *testing.T) {
	if got, err := ParseConfidenceLevel("High"); err != nil || got != ConfidenceHigh {
		t.Errorf("ParseConfidenceLevel(High) = %q, %v; want high, nil", got, err)
	}
	if _, err := ParseConfidenceLevel("certain"); err == nil {
		t.Error("ParseConfidenceLevel(certain) expected error, got nil")
	}
}

func TestRICEScore_ValidateInvalidLevels(t *testing.T) {
	// A non-empty but unrecognized impact/confidence must be rejected rather
	// than silently multiplied by the default fallback.
	bad := RICEScore{FeatureID: "f", Reach: 10, Impact: "gigantic", Confidence: ConfidenceHigh, Effort: 1}
	if err := bad.Validate(); err == nil {
		t.Error("Validate() with invalid impact = nil, want error")
	}

	bad = RICEScore{FeatureID: "f", Reach: 10, Impact: ImpactHigh, Confidence: "certain", Effort: 1}
	if err := bad.Validate(); err == nil {
		t.Error("Validate() with invalid confidence = nil, want error")
	}

	// A fully valid score still passes.
	good := RICEScore{FeatureID: "f", Reach: 10, Impact: ImpactHigh, Confidence: ConfidenceHigh, Effort: 1}
	if err := good.Validate(); err != nil {
		t.Errorf("Validate() with valid levels = %v, want nil", err)
	}

	// Empty impact/confidence remains allowed (partial scoring in progress).
	partial := RICEScore{FeatureID: "f", Reach: 10, Effort: 1}
	if err := partial.Validate(); err != nil {
		t.Errorf("Validate() with empty levels = %v, want nil", err)
	}
}

func TestConfidenceLevel_Multiplier(t *testing.T) {
	tests := []struct {
		level ConfidenceLevel
		want  float64
	}{
		{ConfidenceHigh, 1.0},
		{ConfidenceMedium, 0.8},
		{ConfidenceLow, 0.5},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			if got := tt.level.Multiplier(); got != tt.want {
				t.Errorf("Multiplier() = %v, want %v", got, tt.want)
			}
		})
	}
}
