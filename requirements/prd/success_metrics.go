package prd

// SuccessMetrics contains structured success metrics organized by type.
// These metrics help define and measure the success of the product or feature.
type SuccessMetrics struct {
	// NorthStar contains the primary metrics that define success.
	// These are the most important metrics that the team focuses on.
	NorthStar []Metric `json:"northStar"`

	// Supporting contains metrics that support the north star metrics.
	// These help explain or contribute to north star metric performance.
	Supporting []Metric `json:"supporting,omitempty"`

	// Guardrail contains metrics that should not degrade.
	// These ensure that improvements to north star metrics don't come at
	// the expense of other important aspects of the product.
	Guardrail []Metric `json:"guardrail,omitempty"`
}

// Metric represents a measurable success metric.
type Metric struct {
	// ID is the unique identifier for this metric.
	ID string `json:"id"`

	// Name is the name of this metric.
	Name string `json:"name"`

	// Description provides details about what this metric measures.
	Description string `json:"description,omitempty"`

	// Baseline is the current or starting value of this metric.
	Baseline string `json:"baseline,omitempty"`

	// Target is the goal value for this metric.
	Target string `json:"target"`

	// MeasurementMethod describes how this metric is measured.
	MeasurementMethod string `json:"measurementMethod,omitempty"`
}
