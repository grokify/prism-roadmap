package prioritization

import (
	"testing"
)

func TestClassifyKano(t *testing.T) {
	tests := []struct {
		name           string
		functional     KanoResponse
		dysfunctional  KanoResponse
		wantCategory   KanoCategory
	}{
		// Must-Be examples
		{
			name:          "Must-Be: Expect + Dislike",
			functional:    KanoExpect,
			dysfunctional: KanoDislike,
			wantCategory:  KanoMustBe,
		},
		{
			name:          "Must-Be: Neutral + Dislike",
			functional:    KanoNeutral,
			dysfunctional: KanoDislike,
			wantCategory:  KanoMustBe,
		},
		// Performance examples
		{
			name:          "Performance: Like + Dislike",
			functional:    KanoLike,
			dysfunctional: KanoDislike,
			wantCategory:  KanoPerformance,
		},
		// Attractive examples
		{
			name:          "Attractive: Like + Expect",
			functional:    KanoLike,
			dysfunctional: KanoExpect,
			wantCategory:  KanoAttractive,
		},
		{
			name:          "Attractive: Like + Neutral",
			functional:    KanoLike,
			dysfunctional: KanoNeutral,
			wantCategory:  KanoAttractive,
		},
		// Indifferent examples
		{
			name:          "Indifferent: Expect + Expect",
			functional:    KanoExpect,
			dysfunctional: KanoExpect,
			wantCategory:  KanoIndifferent,
		},
		{
			name:          "Indifferent: Neutral + Neutral",
			functional:    KanoNeutral,
			dysfunctional: KanoNeutral,
			wantCategory:  KanoIndifferent,
		},
		// Reverse examples
		{
			name:          "Reverse: Dislike + Like",
			functional:    KanoDislike,
			dysfunctional: KanoLike,
			wantCategory:  KanoReverse,
		},
		{
			name:          "Reverse: Expect + Like",
			functional:    KanoExpect,
			dysfunctional: KanoLike,
			wantCategory:  KanoReverse,
		},
		// Questionable examples
		{
			name:          "Questionable: Like + Like",
			functional:    KanoLike,
			dysfunctional: KanoLike,
			wantCategory:  KanoQuestionable,
		},
		{
			name:          "Questionable: Dislike + Dislike",
			functional:    KanoDislike,
			dysfunctional: KanoDislike,
			wantCategory:  KanoQuestionable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ClassifyKano(tt.functional, tt.dysfunctional)
			if got != tt.wantCategory {
				t.Errorf("ClassifyKano(%v, %v) = %v, want %v",
					tt.functional, tt.dysfunctional, got, tt.wantCategory)
			}
		})
	}
}

func TestKanoFeature_Classify(t *testing.T) {
	feature := KanoFeature{
		FeatureID:             "feature-1",
		FeatureName:           "Dark Mode",
		FunctionalResponse:    KanoLike,
		DysfunctionalResponse: KanoNeutral,
	}

	category := feature.Classify()

	if category != KanoAttractive {
		t.Errorf("Classify() = %v, want %v", category, KanoAttractive)
	}

	if feature.Category != KanoAttractive {
		t.Errorf("feature.Category = %v, want %v", feature.Category, KanoAttractive)
	}
}

func TestKanoAnalysis_GetByCategory(t *testing.T) {
	analysis := NewKanoAnalysis()

	analysis.Add(KanoFeature{
		FeatureID:             "must-1",
		FunctionalResponse:    KanoExpect,
		DysfunctionalResponse: KanoDislike,
	})
	analysis.Add(KanoFeature{
		FeatureID:             "delight-1",
		FunctionalResponse:    KanoLike,
		DysfunctionalResponse: KanoNeutral,
	})
	analysis.Add(KanoFeature{
		FeatureID:             "must-2",
		FunctionalResponse:    KanoNeutral,
		DysfunctionalResponse: KanoDislike,
	})

	mustHaves := analysis.MustHaves()
	if len(mustHaves) != 2 {
		t.Errorf("MustHaves() returned %d items, want 2", len(mustHaves))
	}

	delighters := analysis.Delighters()
	if len(delighters) != 1 {
		t.Errorf("Delighters() returned %d items, want 1", len(delighters))
	}
}

func TestKanoCategory_Priority(t *testing.T) {
	// Must-Be should have highest priority
	if KanoMustBe.Priority() <= KanoPerformance.Priority() {
		t.Error("Must-Be should have higher priority than Performance")
	}

	// Performance should have higher priority than Attractive
	if KanoPerformance.Priority() <= KanoAttractive.Priority() {
		t.Error("Performance should have higher priority than Attractive")
	}

	// Reverse should have zero priority
	if KanoReverse.Priority() != 0 {
		t.Error("Reverse should have zero priority")
	}
}

func TestKanoFeature_CalculateCoefficients(t *testing.T) {
	feature := KanoFeature{
		FeatureID:        "test",
		AttractiveCount:  20,
		PerformanceCount: 30,
		MustBeCount:      40,
		IndifferentCount: 10,
	}

	feature.CalculateCoefficients()

	// Satisfaction = (A+O)/(A+O+M+I) = (20+30)/(20+30+40+10) = 50/100 = 0.5
	if feature.SatisfactionCoefficient != 0.5 {
		t.Errorf("SatisfactionCoefficient = %v, want 0.5", feature.SatisfactionCoefficient)
	}

	// Dissatisfaction = -1*(O+M)/(A+O+M+I) = -1*(30+40)/100 = -0.7
	if feature.DissatisfactionCoefficient != -0.7 {
		t.Errorf("DissatisfactionCoefficient = %v, want -0.7", feature.DissatisfactionCoefficient)
	}
}
