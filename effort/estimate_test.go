package effort

import "testing"

func TestTShirtSizePersonDays(t *testing.T) {
	tests := []struct {
		size TShirtSize
		want int
	}{
		{TShirtXS, 2},
		{TShirtS, 4},
		{TShirtM, 8},
		{TShirtL, 15},
		{TShirtXL, 30},
		{"invalid", 0},
	}

	for _, tt := range tests {
		got := tt.size.PersonDays()
		if got != tt.want {
			t.Errorf("TShirtSize(%s).PersonDays() = %d, want %d", tt.size, got, tt.want)
		}
	}
}

func TestTShirtSizeMinMaxDays(t *testing.T) {
	tests := []struct {
		size    TShirtSize
		wantMin int
		wantMax int
	}{
		{TShirtXS, 1, 2},
		{TShirtS, 3, 5},
		{TShirtM, 5, 10},
		{TShirtL, 10, 20},
		{TShirtXL, 20, 60},
	}

	for _, tt := range tests {
		gotMin := tt.size.MinDays()
		gotMax := tt.size.MaxDays()
		if gotMin != tt.wantMin {
			t.Errorf("TShirtSize(%s).MinDays() = %d, want %d", tt.size, gotMin, tt.wantMin)
		}
		if gotMax != tt.wantMax {
			t.Errorf("TShirtSize(%s).MaxDays() = %d, want %d", tt.size, gotMax, tt.wantMax)
		}
	}
}

func TestConfidenceMultiplier(t *testing.T) {
	tests := []struct {
		confidence Confidence
		want       float64
	}{
		{ConfidenceHigh, 1.0},
		{ConfidenceMedium, 1.5},
		{ConfidenceLow, 2.0},
		{"invalid", 1.5}, // default
	}

	for _, tt := range tests {
		got := tt.confidence.Multiplier()
		if got != tt.want {
			t.Errorf("Confidence(%s).Multiplier() = %f, want %f", tt.confidence, got, tt.want)
		}
	}
}

func TestEffortEstimateRiskAdjustedDays(t *testing.T) {
	tests := []struct {
		name     string
		estimate EffortEstimate
		want     float64
	}{
		{
			name: "high confidence",
			estimate: EffortEstimate{
				PersonDays: 10,
				Confidence: ConfidenceHigh,
			},
			want: 10.0,
		},
		{
			name: "medium confidence",
			estimate: EffortEstimate{
				PersonDays: 10,
				Confidence: ConfidenceMedium,
			},
			want: 15.0,
		},
		{
			name: "low confidence",
			estimate: EffortEstimate{
				PersonDays: 10,
				Confidence: ConfidenceLow,
			},
			want: 20.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.estimate.RiskAdjustedDays()
			if got != tt.want {
				t.Errorf("RiskAdjustedDays() = %f, want %f", got, tt.want)
			}
		})
	}
}

func TestEffortEstimateEffectivePersonDays(t *testing.T) {
	tests := []struct {
		name     string
		estimate EffortEstimate
		want     int
	}{
		{
			name: "explicit days",
			estimate: EffortEstimate{
				PersonDays: 15,
				TShirtSize: TShirtM, // Should be ignored
			},
			want: 15,
		},
		{
			name: "from t-shirt",
			estimate: EffortEstimate{
				TShirtSize: TShirtM,
			},
			want: 8,
		},
		{
			name:     "neither set",
			estimate: EffortEstimate{},
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.estimate.EffectivePersonDays()
			if got != tt.want {
				t.Errorf("EffectivePersonDays() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestEffortEstimateValidate(t *testing.T) {
	tests := []struct {
		name     string
		estimate EffortEstimate
		wantErr  bool
	}{
		{
			name: "valid",
			estimate: EffortEstimate{
				PersonDays: 10,
				TShirtSize: TShirtM,
				Confidence: ConfidenceMedium,
			},
			wantErr: false,
		},
		{
			name: "negative days",
			estimate: EffortEstimate{
				PersonDays: -5,
			},
			wantErr: true,
		},
		{
			name: "invalid t-shirt",
			estimate: EffortEstimate{
				TShirtSize: "XXL",
			},
			wantErr: true,
		},
		{
			name: "invalid confidence",
			estimate: EffortEstimate{
				Confidence: "very-high",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.estimate.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEffortEstimate(t *testing.T) {
	estimate := NewEffortEstimate(15, ConfidenceMedium)

	if estimate.PersonDays != 15 {
		t.Errorf("PersonDays = %d, want 15", estimate.PersonDays)
	}
	if estimate.Confidence != ConfidenceMedium {
		t.Errorf("Confidence = %s, want medium", estimate.Confidence)
	}
	if estimate.Unit != EffortUnitPersonDays {
		t.Errorf("Unit = %s, want person-days", estimate.Unit)
	}
}

func TestNewEffortEstimateFromTShirt(t *testing.T) {
	estimate := NewEffortEstimateFromTShirt(TShirtL, ConfidenceHigh)

	if estimate.PersonDays != 15 {
		t.Errorf("PersonDays = %d, want 15", estimate.PersonDays)
	}
	if estimate.TShirtSize != TShirtL {
		t.Errorf("TShirtSize = %s, want L", estimate.TShirtSize)
	}
	if estimate.Confidence != ConfidenceHigh {
		t.Errorf("Confidence = %s, want high", estimate.Confidence)
	}
}

func TestTShirtFromPersonDays(t *testing.T) {
	tests := []struct {
		days int
		want TShirtSize
	}{
		{1, TShirtXS},
		{2, TShirtXS},
		{3, TShirtS},
		{5, TShirtS},
		{6, TShirtM},
		{10, TShirtM},
		{11, TShirtL},
		{20, TShirtL},
		{21, TShirtXL},
		{100, TShirtXL},
	}

	for _, tt := range tests {
		got := TShirtFromPersonDays(tt.days)
		if got != tt.want {
			t.Errorf("TShirtFromPersonDays(%d) = %s, want %s", tt.days, got, tt.want)
		}
	}
}

func TestValidTShirtSizes(t *testing.T) {
	sizes := ValidTShirtSizes()
	if len(sizes) != 5 {
		t.Errorf("ValidTShirtSizes() returned %d sizes, want 5", len(sizes))
	}
}

func TestIsValidTShirtSize(t *testing.T) {
	if !IsValidTShirtSize(TShirtM) {
		t.Error("IsValidTShirtSize(M) = false, want true")
	}
	if IsValidTShirtSize("XXL") {
		t.Error("IsValidTShirtSize(XXL) = true, want false")
	}
}

func TestValidConfidenceLevels(t *testing.T) {
	levels := ValidConfidenceLevels()
	if len(levels) != 3 {
		t.Errorf("ValidConfidenceLevels() returned %d levels, want 3", len(levels))
	}
}

func TestIsValidConfidence(t *testing.T) {
	if !IsValidConfidence(ConfidenceHigh) {
		t.Error("IsValidConfidence(high) = false, want true")
	}
	if IsValidConfidence("very-high") {
		t.Error("IsValidConfidence(very-high) = true, want false")
	}
}
