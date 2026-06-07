package prioritization

import (
	"fmt"
	"sort"
)

// Kano Model Framework
// Reference: https://www.productplan.com/glossary/kano-model/
// Classifies features by customer satisfaction impact.

// KanoCategory represents the classification of a feature in the Kano Model.
type KanoCategory string

const (
	// KanoMustBe (Basic/Threshold) - Expected features. Their absence causes
	// extreme dissatisfaction, but their presence doesn't increase satisfaction.
	// Example: A car must have brakes.
	KanoMustBe KanoCategory = "must-be"

	// KanoPerformance (One-dimensional) - Linear relationship between
	// fulfillment and satisfaction. More is better.
	// Example: Fuel efficiency - better mileage = more satisfaction.
	KanoPerformance KanoCategory = "performance"

	// KanoAttractive (Delighter/Excitement) - Unexpected features that
	// create high satisfaction when present but don't cause dissatisfaction when absent.
	// Example: Heated steering wheel - nice surprise, not expected.
	KanoAttractive KanoCategory = "attractive"

	// KanoIndifferent - Features that don't significantly affect satisfaction
	// either way. Customers don't care much about these.
	KanoIndifferent KanoCategory = "indifferent"

	// KanoReverse - Features that some customers actively dislike.
	// Their presence causes dissatisfaction.
	// Example: Auto-play videos on websites.
	KanoReverse KanoCategory = "reverse"

	// KanoQuestionable - Contradictory responses that indicate the question
	// wasn't understood or the feature wasn't clear.
	KanoQuestionable KanoCategory = "questionable"
)

// String returns the string representation.
func (k KanoCategory) String() string {
	return string(k)
}

// Description returns a human-readable description of the category.
func (k KanoCategory) Description() string {
	switch k {
	case KanoMustBe:
		return "Basic expectation - absence causes dissatisfaction"
	case KanoPerformance:
		return "Linear satisfaction - more is better"
	case KanoAttractive:
		return "Delighter - unexpected positive surprise"
	case KanoIndifferent:
		return "No significant impact on satisfaction"
	case KanoReverse:
		return "Unwanted - presence causes dissatisfaction"
	case KanoQuestionable:
		return "Unclear response - needs clarification"
	default:
		return "Unknown category"
	}
}

// Priority returns a priority weight for the category (higher = more important to implement).
func (k KanoCategory) Priority() int {
	switch k {
	case KanoMustBe:
		return 5 // Must have - highest priority
	case KanoPerformance:
		return 4 // Important for competitive differentiation
	case KanoAttractive:
		return 3 // Nice to have for delight
	case KanoIndifferent:
		return 1 // Low priority
	case KanoReverse:
		return 0 // Avoid implementing
	case KanoQuestionable:
		return 0 // Needs clarification first
	default:
		return 0
	}
}

// KanoResponse represents a response to a Kano questionnaire.
type KanoResponse string

const (
	KanoLike     KanoResponse = "like"     // I like it
	KanoExpect   KanoResponse = "expect"   // I expect it
	KanoNeutral  KanoResponse = "neutral"  // I am neutral
	KanoTolerate KanoResponse = "tolerate" // I can tolerate it
	KanoDislike  KanoResponse = "dislike"  // I dislike it
)

// KanoFeature represents a feature being evaluated with the Kano Model.
type KanoFeature struct {
	// Feature identification
	FeatureID   string `json:"featureId"`
	FeatureName string `json:"featureName"`
	Description string `json:"description,omitempty"`

	// Kano questionnaire responses
	// Functional question: "If the product HAS this feature, how do you feel?"
	FunctionalResponse KanoResponse `json:"functionalResponse"`
	// Dysfunctional question: "If the product DOES NOT HAVE this feature, how do you feel?"
	DysfunctionalResponse KanoResponse `json:"dysfunctionalResponse"`

	// Classified category (derived from responses)
	Category KanoCategory `json:"category"`

	// Satisfaction coefficients (from aggregate responses)
	// Range: -1 to 1
	SatisfactionCoefficient    float64 `json:"satisfactionCoefficient,omitempty"`    // (A+O)/(A+O+M+I)
	DissatisfactionCoefficient float64 `json:"dissatisfactionCoefficient,omitempty"` // (O+M)/(A+O+M+I)*-1

	// Response counts for aggregate analysis
	MustBeCount       int `json:"mustBeCount,omitempty"`
	PerformanceCount  int `json:"performanceCount,omitempty"`
	AttractiveCount   int `json:"attractiveCount,omitempty"`
	IndifferentCount  int `json:"indifferentCount,omitempty"`
	ReverseCount      int `json:"reverseCount,omitempty"`
	QuestionableCount int `json:"questionableCount,omitempty"`

	// Metadata
	RespondentCount int    `json:"respondentCount,omitempty"`
	Notes           string `json:"notes,omitempty"`
}

// Classify determines the Kano category based on functional/dysfunctional responses.
// Uses the standard Kano evaluation table.
func (f *KanoFeature) Classify() KanoCategory {
	f.Category = ClassifyKano(f.FunctionalResponse, f.DysfunctionalResponse)
	return f.Category
}

// ClassifyKano classifies a feature using the Kano evaluation table.
//
// Kano Evaluation Table:
//
//	                    | Dysfunctional Response
//	Functional Response | Like    | Expect  | Neutral | Tolerate | Dislike
//	--------------------|---------|---------|---------|----------|--------
//	Like                | Q       | A       | A       | A        | O
//	Expect              | R       | I       | I       | I        | M
//	Neutral             | R       | I       | I       | I        | M
//	Tolerate            | R       | I       | I       | I        | M
//	Dislike             | R       | R       | R       | R        | Q
//
// Legend: M=Must-be, O=One-dimensional, A=Attractive, I=Indifferent, R=Reverse, Q=Questionable
func ClassifyKano(functional, dysfunctional KanoResponse) KanoCategory {
	// Evaluation table
	table := map[KanoResponse]map[KanoResponse]KanoCategory{
		KanoLike: {
			KanoLike:     KanoQuestionable,
			KanoExpect:   KanoAttractive,
			KanoNeutral:  KanoAttractive,
			KanoTolerate: KanoAttractive,
			KanoDislike:  KanoPerformance,
		},
		KanoExpect: {
			KanoLike:     KanoReverse,
			KanoExpect:   KanoIndifferent,
			KanoNeutral:  KanoIndifferent,
			KanoTolerate: KanoIndifferent,
			KanoDislike:  KanoMustBe,
		},
		KanoNeutral: {
			KanoLike:     KanoReverse,
			KanoExpect:   KanoIndifferent,
			KanoNeutral:  KanoIndifferent,
			KanoTolerate: KanoIndifferent,
			KanoDislike:  KanoMustBe,
		},
		KanoTolerate: {
			KanoLike:     KanoReverse,
			KanoExpect:   KanoIndifferent,
			KanoNeutral:  KanoIndifferent,
			KanoTolerate: KanoIndifferent,
			KanoDislike:  KanoMustBe,
		},
		KanoDislike: {
			KanoLike:     KanoReverse,
			KanoExpect:   KanoReverse,
			KanoNeutral:  KanoReverse,
			KanoTolerate: KanoReverse,
			KanoDislike:  KanoQuestionable,
		},
	}

	if funcRow, ok := table[functional]; ok {
		if category, ok := funcRow[dysfunctional]; ok {
			return category
		}
	}
	return KanoQuestionable
}

// CalculateCoefficients calculates satisfaction and dissatisfaction coefficients
// from aggregate response counts.
// Satisfaction: (A+O)/(A+O+M+I) - how much satisfaction increases with feature
// Dissatisfaction: -1*(O+M)/(A+O+M+I) - how much satisfaction decreases without feature
func (f *KanoFeature) CalculateCoefficients() {
	total := float64(f.AttractiveCount + f.PerformanceCount + f.MustBeCount + f.IndifferentCount)
	if total == 0 {
		f.SatisfactionCoefficient = 0
		f.DissatisfactionCoefficient = 0
		return
	}

	f.SatisfactionCoefficient = float64(f.AttractiveCount+f.PerformanceCount) / total
	f.DissatisfactionCoefficient = -1 * float64(f.PerformanceCount+f.MustBeCount) / total
}

// IsComplete returns true if the feature has been classified.
func (f *KanoFeature) IsComplete() bool {
	return f.FeatureID != "" && f.Category != ""
}

// KanoAnalysis represents a Kano analysis for multiple features.
type KanoAnalysis struct {
	Features     []KanoFeature `json:"features"`
	Description  string        `json:"description,omitempty"`
	AnalyzedBy   string        `json:"analyzedBy,omitempty"`
	AnalyzedDate string        `json:"analyzedDate,omitempty"`
}

// NewKanoAnalysis creates a new Kano analysis.
func NewKanoAnalysis() *KanoAnalysis {
	return &KanoAnalysis{
		Features: []KanoFeature{},
	}
}

// Add adds a feature to the analysis.
func (a *KanoAnalysis) Add(feature KanoFeature) {
	feature.Classify()
	a.Features = append(a.Features, feature)
}

// ClassifyAll classifies all features.
func (a *KanoAnalysis) ClassifyAll() {
	for i := range a.Features {
		a.Features[i].Classify()
	}
}

// GetByCategory returns features matching the given category.
func (a *KanoAnalysis) GetByCategory(category KanoCategory) []KanoFeature {
	var result []KanoFeature
	for _, f := range a.Features {
		if f.Category == category {
			result = append(result, f)
		}
	}
	return result
}

// MustHaves returns all Must-Be features.
func (a *KanoAnalysis) MustHaves() []KanoFeature {
	return a.GetByCategory(KanoMustBe)
}

// PerformanceFeatures returns all Performance features.
func (a *KanoAnalysis) PerformanceFeatures() []KanoFeature {
	return a.GetByCategory(KanoPerformance)
}

// Delighters returns all Attractive (delighter) features.
func (a *KanoAnalysis) Delighters() []KanoFeature {
	return a.GetByCategory(KanoAttractive)
}

// SortByPriority sorts features by Kano priority (Must-Be first, then Performance, etc.).
func (a *KanoAnalysis) SortByPriority() {
	sort.Slice(a.Features, func(i, j int) bool {
		return a.Features[i].Category.Priority() > a.Features[j].Category.Priority()
	})
}

// SortBySatisfaction sorts features by satisfaction coefficient (highest first).
func (a *KanoAnalysis) SortBySatisfaction() {
	sort.Slice(a.Features, func(i, j int) bool {
		return a.Features[i].SatisfactionCoefficient > a.Features[j].SatisfactionCoefficient
	})
}

// Summary returns a count of features by category.
func (a *KanoAnalysis) Summary() map[KanoCategory]int {
	counts := make(map[KanoCategory]int)
	for _, f := range a.Features {
		counts[f.Category]++
	}
	return counts
}

// Validate returns an error if the analysis is invalid.
func (a *KanoAnalysis) Validate() error {
	for i, f := range a.Features {
		if f.FeatureID == "" {
			return fmt.Errorf("feature %d: featureId is required", i)
		}
	}
	return nil
}

// ValidKanoResponses returns all valid Kano questionnaire responses.
func ValidKanoResponses() []KanoResponse {
	return []KanoResponse{
		KanoLike,
		KanoExpect,
		KanoNeutral,
		KanoTolerate,
		KanoDislike,
	}
}

// ValidKanoCategories returns all valid Kano categories.
func ValidKanoCategories() []KanoCategory {
	return []KanoCategory{
		KanoMustBe,
		KanoPerformance,
		KanoAttractive,
		KanoIndifferent,
		KanoReverse,
		KanoQuestionable,
	}
}
