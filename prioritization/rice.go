// Package prioritization provides feature prioritization frameworks
// including RICE scoring, Kano Model, and related prioritization tools.
package prioritization

import (
	"fmt"
	"sort"
	"strings"
)

// RICE Scoring Framework
// Reference: https://www.intercom.com/blog/rice-simple-prioritization-for-product-managers/
// Formula: Score = (Reach × Impact × Confidence) / Effort

// ImpactLevel represents the impact multiplier in RICE scoring.
type ImpactLevel string

const (
	ImpactMassive ImpactLevel = "massive" // 3x - Massive impact
	ImpactHigh    ImpactLevel = "high"    // 2x - High impact
	ImpactMedium  ImpactLevel = "medium"  // 1x - Medium impact
	ImpactLow     ImpactLevel = "low"     // 0.5x - Low impact
	ImpactMinimal ImpactLevel = "minimal" // 0.25x - Minimal impact
)

// ImpactMultiplier returns the numeric multiplier for an impact level.
func (i ImpactLevel) Multiplier() float64 {
	switch i {
	case ImpactMassive:
		return 3.0
	case ImpactHigh:
		return 2.0
	case ImpactMedium:
		return 1.0
	case ImpactLow:
		return 0.5
	case ImpactMinimal:
		return 0.25
	default:
		return 1.0
	}
}

// String returns the string representation.
func (i ImpactLevel) String() string {
	return string(i)
}

// IsValid reports whether i is a recognized impact level.
func (i ImpactLevel) IsValid() bool {
	switch i {
	case ImpactMassive, ImpactHigh, ImpactMedium, ImpactLow, ImpactMinimal:
		return true
	default:
		return false
	}
}

// ParseImpactLevel parses a string into an ImpactLevel (case-insensitive).
// Returns an error if the value is not a recognized level.
func ParseImpactLevel(s string) (ImpactLevel, error) {
	level := ImpactLevel(strings.ToLower(strings.TrimSpace(s)))
	if !level.IsValid() {
		return "", fmt.Errorf("invalid impact level: %q", s)
	}
	return level, nil
}

// ConfidenceLevel represents confidence in RICE estimates.
type ConfidenceLevel string

const (
	ConfidenceHigh   ConfidenceLevel = "high"   // 100% - High confidence, strong data
	ConfidenceMedium ConfidenceLevel = "medium" // 80% - Medium confidence, some data
	ConfidenceLow    ConfidenceLevel = "low"    // 50% - Low confidence, gut feel
)

// Multiplier returns the confidence percentage as a decimal.
func (c ConfidenceLevel) Multiplier() float64 {
	switch c {
	case ConfidenceHigh:
		return 1.0
	case ConfidenceMedium:
		return 0.8
	case ConfidenceLow:
		return 0.5
	default:
		return 0.8
	}
}

// String returns the string representation.
func (c ConfidenceLevel) String() string {
	return string(c)
}

// IsValid reports whether c is a recognized confidence level.
func (c ConfidenceLevel) IsValid() bool {
	switch c {
	case ConfidenceHigh, ConfidenceMedium, ConfidenceLow:
		return true
	default:
		return false
	}
}

// ParseConfidenceLevel parses a string into a ConfidenceLevel (case-insensitive).
// Returns an error if the value is not a recognized level.
func ParseConfidenceLevel(s string) (ConfidenceLevel, error) {
	level := ConfidenceLevel(strings.ToLower(strings.TrimSpace(s)))
	if !level.IsValid() {
		return "", fmt.Errorf("invalid confidence level: %q", s)
	}
	return level, nil
}

// RICEScore represents a RICE prioritization score for a feature.
type RICEScore struct {
	// Feature identification
	FeatureID   string `json:"featureId"`             // Unique feature identifier
	FeatureName string `json:"featureName,omitempty"` // Human-readable name

	// RICE components
	Reach      int             `json:"reach"`      // Number of users/customers affected per time period
	ReachUnit  string          `json:"reachUnit"`  // Time period: "quarter", "month", "year"
	Impact     ImpactLevel     `json:"impact"`     // Impact level: massive, high, medium, low, minimal
	Confidence ConfidenceLevel `json:"confidence"` // Confidence: high, medium, low
	Effort     float64         `json:"effort"`     // Person-months (or person-weeks)
	EffortUnit string          `json:"effortUnit"` // Unit: "person-months", "person-weeks", "story-points"

	// Calculated score
	Score float64 `json:"score"` // Calculated RICE score

	// Supporting context
	ReachJustification      string `json:"reachJustification,omitempty"`      // How reach was estimated
	ImpactJustification     string `json:"impactJustification,omitempty"`     // Why this impact level
	ConfidenceJustification string `json:"confidenceJustification,omitempty"` // What data supports confidence
	EffortJustification     string `json:"effortJustification,omitempty"`     // How effort was estimated

	// Metadata
	ScoredBy   string `json:"scoredBy,omitempty"`   // Who scored this
	ScoredDate string `json:"scoredDate,omitempty"` // When scored
	Notes      string `json:"notes,omitempty"`      // Additional context
}

// Calculate computes the RICE score.
// Formula: Score = (Reach × Impact × Confidence) / Effort
func (r *RICEScore) Calculate() float64 {
	if r.Effort == 0 {
		r.Score = 0
		return 0
	}

	impactMultiplier := r.Impact.Multiplier()
	confidenceMultiplier := r.Confidence.Multiplier()

	r.Score = (float64(r.Reach) * impactMultiplier * confidenceMultiplier) / r.Effort
	return r.Score
}

// IsComplete returns true if all required fields are set.
func (r *RICEScore) IsComplete() bool {
	return r.FeatureID != "" &&
		r.Reach > 0 &&
		r.Impact != "" &&
		r.Confidence != "" &&
		r.Effort > 0
}

// Validate returns an error if the score is invalid.
func (r *RICEScore) Validate() error {
	if r.FeatureID == "" {
		return fmt.Errorf("featureId is required")
	}
	if r.Reach < 0 {
		return fmt.Errorf("reach must be non-negative")
	}
	if r.Effort < 0 {
		return fmt.Errorf("effort must be non-negative")
	}
	// A non-empty but unrecognized impact/confidence silently multiplies by the
	// default fallback in Multiplier(), which hides data-entry errors. Reject it.
	if r.Impact != "" && !r.Impact.IsValid() {
		return fmt.Errorf("invalid impact level: %s", r.Impact)
	}
	if r.Confidence != "" && !r.Confidence.IsValid() {
		return fmt.Errorf("invalid confidence level: %s", r.Confidence)
	}
	return nil
}

// RICEScoreSet represents a collection of RICE scores for ranking.
type RICEScoreSet struct {
	Scores      []RICEScore `json:"scores"`
	Description string      `json:"description,omitempty"`
	ScoredDate  string      `json:"scoredDate,omitempty"`
}

// NewRICEScoreSet creates a new score set.
func NewRICEScoreSet() *RICEScoreSet {
	return &RICEScoreSet{
		Scores: []RICEScore{},
	}
}

// Add adds a score to the set.
func (s *RICEScoreSet) Add(score RICEScore) {
	score.Calculate()
	s.Scores = append(s.Scores, score)
}

// CalculateAll calculates scores for all features.
func (s *RICEScoreSet) CalculateAll() {
	for i := range s.Scores {
		s.Scores[i].Calculate()
	}
}

// SortByScore sorts features by RICE score (highest first).
func (s *RICEScoreSet) SortByScore() {
	sort.Slice(s.Scores, func(i, j int) bool {
		return s.Scores[i].Score > s.Scores[j].Score
	})
}

// TopN returns the top N features by score.
func (s *RICEScoreSet) TopN(n int) []RICEScore {
	s.SortByScore()
	if n > len(s.Scores) {
		n = len(s.Scores)
	}
	return s.Scores[:n]
}

// GetByID returns a score by feature ID.
func (s *RICEScoreSet) GetByID(featureID string) *RICEScore {
	for i := range s.Scores {
		if s.Scores[i].FeatureID == featureID {
			return &s.Scores[i]
		}
	}
	return nil
}

// Rank returns the rank of a feature (1-based, 0 if not found).
func (s *RICEScoreSet) Rank(featureID string) int {
	s.SortByScore()
	for i, score := range s.Scores {
		if score.FeatureID == featureID {
			return i + 1
		}
	}
	return 0
}

// NewRICEScore creates a new RICE score with the given values.
func NewRICEScore(featureID string, reach int, impact ImpactLevel, confidence ConfidenceLevel, effort float64) *RICEScore {
	score := &RICEScore{
		FeatureID:  featureID,
		Reach:      reach,
		ReachUnit:  "quarter",
		Impact:     impact,
		Confidence: confidence,
		Effort:     effort,
		EffortUnit: "person-months",
	}
	score.Calculate()
	return score
}

// ValidImpactLevels returns all valid impact levels.
func ValidImpactLevels() []ImpactLevel {
	return []ImpactLevel{
		ImpactMassive,
		ImpactHigh,
		ImpactMedium,
		ImpactLow,
		ImpactMinimal,
	}
}

// ValidConfidenceLevels returns all valid confidence levels.
func ValidConfidenceLevels() []ConfidenceLevel {
	return []ConfidenceLevel{
		ConfidenceHigh,
		ConfidenceMedium,
		ConfidenceLow,
	}
}
