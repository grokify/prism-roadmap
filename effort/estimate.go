package effort

import "fmt"

// EffortEstimate captures effort estimation for a roadmap item.
type EffortEstimate struct {
	PersonDays int        `json:"person_days"`           // Estimated person-days
	TShirtSize TShirtSize `json:"tshirt_size,omitempty"` // T-shirt sizing
	Confidence Confidence `json:"confidence"`            // Estimate confidence level
	Unit       EffortUnit `json:"unit,omitempty"`        // person-days, person-weeks, etc.

	// Additional context
	Justification string `json:"justification,omitempty"`  // How estimate was derived
	EstimatedBy   string `json:"estimated_by,omitempty"`   // Who made the estimate
	EstimatedDate string `json:"estimated_date,omitempty"` // When estimated (ISO 8601)
	Notes         string `json:"notes,omitempty"`
}

// TShirtSize provides rough effort estimation using t-shirt sizing.
type TShirtSize string

const (
	TShirtXS TShirtSize = "XS" // < 3 days
	TShirtS  TShirtSize = "S"  // 3-5 days
	TShirtM  TShirtSize = "M"  // 1-2 weeks (5-10 days)
	TShirtL  TShirtSize = "L"  // 2-4 weeks (10-20 days)
	TShirtXL TShirtSize = "XL" // > 4 weeks (20+ days)
)

// String returns the string representation.
func (t TShirtSize) String() string {
	return string(t)
}

// PersonDays returns the median person-days for a t-shirt size.
func (t TShirtSize) PersonDays() int {
	switch t {
	case TShirtXS:
		return 2
	case TShirtS:
		return 4
	case TShirtM:
		return 8
	case TShirtL:
		return 15
	case TShirtXL:
		return 30
	default:
		return 0
	}
}

// MinDays returns the minimum days for this t-shirt size.
func (t TShirtSize) MinDays() int {
	switch t {
	case TShirtXS:
		return 1
	case TShirtS:
		return 3
	case TShirtM:
		return 5
	case TShirtL:
		return 10
	case TShirtXL:
		return 20
	default:
		return 0
	}
}

// MaxDays returns the maximum days for this t-shirt size.
func (t TShirtSize) MaxDays() int {
	switch t {
	case TShirtXS:
		return 2
	case TShirtS:
		return 5
	case TShirtM:
		return 10
	case TShirtL:
		return 20
	case TShirtXL:
		return 60 // Capped at ~3 months
	default:
		return 0
	}
}

// Description returns a human-readable description.
func (t TShirtSize) Description() string {
	switch t {
	case TShirtXS:
		return "Extra Small: less than 3 days"
	case TShirtS:
		return "Small: 3-5 days"
	case TShirtM:
		return "Medium: 1-2 weeks"
	case TShirtL:
		return "Large: 2-4 weeks"
	case TShirtXL:
		return "Extra Large: more than 4 weeks"
	default:
		return "Unknown size"
	}
}

// Confidence represents the confidence level of an estimate.
type Confidence string

const (
	ConfidenceLow    Confidence = "low"    // Gut feel, limited information
	ConfidenceMedium Confidence = "medium" // Some data, reasonable assumptions
	ConfidenceHigh   Confidence = "high"   // Strong data, similar past work
)

// String returns the string representation.
func (c Confidence) String() string {
	return string(c)
}

// Multiplier returns a confidence multiplier for risk calculations.
// Lower confidence = higher risk multiplier.
func (c Confidence) Multiplier() float64 {
	switch c {
	case ConfidenceHigh:
		return 1.0
	case ConfidenceMedium:
		return 1.5
	case ConfidenceLow:
		return 2.0
	default:
		return 1.5
	}
}

// Description returns a human-readable description.
func (c Confidence) Description() string {
	switch c {
	case ConfidenceHigh:
		return "High confidence - strong data or similar past work"
	case ConfidenceMedium:
		return "Medium confidence - some data, reasonable assumptions"
	case ConfidenceLow:
		return "Low confidence - limited information, gut feel"
	default:
		return "Unknown confidence level"
	}
}

// EffortUnit represents the unit of effort measurement.
type EffortUnit string

const (
	EffortUnitPersonDays   EffortUnit = "person-days"
	EffortUnitPersonWeeks  EffortUnit = "person-weeks"
	EffortUnitPersonMonths EffortUnit = "person-months"
	EffortUnitStoryPoints  EffortUnit = "story-points"
)

// RiskAdjustedDays returns person-days adjusted for confidence risk.
func (e *EffortEstimate) RiskAdjustedDays() float64 {
	return float64(e.PersonDays) * e.Confidence.Multiplier()
}

// IsComplete returns true if required fields are set.
func (e *EffortEstimate) IsComplete() bool {
	return e.PersonDays > 0 || e.TShirtSize != ""
}

// Validate returns an error if the estimate is invalid.
func (e *EffortEstimate) Validate() error {
	if e.PersonDays < 0 {
		return fmt.Errorf("person_days must be non-negative")
	}
	if e.TShirtSize != "" && !IsValidTShirtSize(e.TShirtSize) {
		return fmt.Errorf("invalid tshirt_size: %s", e.TShirtSize)
	}
	if e.Confidence != "" && !IsValidConfidence(e.Confidence) {
		return fmt.Errorf("invalid confidence: %s", e.Confidence)
	}
	return nil
}

// EffectivePersonDays returns the person-days, using t-shirt size if explicit days not set.
func (e *EffortEstimate) EffectivePersonDays() int {
	if e.PersonDays > 0 {
		return e.PersonDays
	}
	if e.TShirtSize != "" {
		return e.TShirtSize.PersonDays()
	}
	return 0
}

// NewEffortEstimate creates a new effort estimate.
func NewEffortEstimate(personDays int, confidence Confidence) *EffortEstimate {
	return &EffortEstimate{
		PersonDays: personDays,
		Confidence: confidence,
		Unit:       EffortUnitPersonDays,
	}
}

// NewEffortEstimateFromTShirt creates an effort estimate from t-shirt size.
func NewEffortEstimateFromTShirt(size TShirtSize, confidence Confidence) *EffortEstimate {
	return &EffortEstimate{
		PersonDays: size.PersonDays(),
		TShirtSize: size,
		Confidence: confidence,
		Unit:       EffortUnitPersonDays,
	}
}

// ValidTShirtSizes returns all valid t-shirt sizes.
func ValidTShirtSizes() []TShirtSize {
	return []TShirtSize{
		TShirtXS,
		TShirtS,
		TShirtM,
		TShirtL,
		TShirtXL,
	}
}

// IsValidTShirtSize returns true if the size is valid.
func IsValidTShirtSize(s TShirtSize) bool {
	for _, valid := range ValidTShirtSizes() {
		if s == valid {
			return true
		}
	}
	return false
}

// ValidConfidenceLevels returns all valid confidence levels.
func ValidConfidenceLevels() []Confidence {
	return []Confidence{
		ConfidenceHigh,
		ConfidenceMedium,
		ConfidenceLow,
	}
}

// IsValidConfidence returns true if the confidence is valid.
func IsValidConfidence(c Confidence) bool {
	for _, valid := range ValidConfidenceLevels() {
		if c == valid {
			return true
		}
	}
	return false
}

// PersonDaysFromTShirt returns estimated person-days for a t-shirt size.
func PersonDaysFromTShirt(size TShirtSize) int {
	return size.PersonDays()
}

// TShirtFromPersonDays returns the appropriate t-shirt size for given person-days.
func TShirtFromPersonDays(days int) TShirtSize {
	switch {
	case days <= 2:
		return TShirtXS
	case days <= 5:
		return TShirtS
	case days <= 10:
		return TShirtM
	case days <= 20:
		return TShirtL
	default:
		return TShirtXL
	}
}
