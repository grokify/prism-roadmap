package rubrics

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// Rubric is the rich, LLM-as-a-Judge evaluation rubric for a canvas or document.
// It is the canonical rubric format for the ecosystem: weighted categories, each
// with weighted criteria, each scored at pass / partial / fail with a description
// and concrete indicators. Downstream tools (for example VisionSpec) import this
// type rather than reimplementing it.
type Rubric struct {
	Name        string     `json:"name" yaml:"name"`
	Version     string     `json:"version,omitempty" yaml:"version,omitempty"`
	Description string     `json:"description,omitempty" yaml:"description,omitempty"`
	Thresholds  Thresholds `json:"thresholds,omitempty" yaml:"thresholds,omitempty"`
	Categories  []Category `json:"categories" yaml:"categories"`
}

// Thresholds are the overall scoring cutoffs (0-100). A score at or above Pass
// passes; at or above Partial is partial; below Partial fails.
type Thresholds struct {
	Pass    int `json:"pass" yaml:"pass"`
	Partial int `json:"partial" yaml:"partial"`
}

// Category is a weighted evaluation dimension containing one or more criteria.
type Category struct {
	ID          string      `json:"id,omitempty" yaml:"id,omitempty"`
	Name        string      `json:"name" yaml:"name"`
	Weight      float64     `json:"weight" yaml:"weight"`
	Description string      `json:"description,omitempty" yaml:"description,omitempty"`
	Criteria    []Criterion `json:"criteria" yaml:"criteria"`
}

// Criterion is a weighted, individually scored check within a category.
type Criterion struct {
	ID      string  `json:"id,omitempty" yaml:"id,omitempty"`
	Name    string  `json:"name" yaml:"name"`
	Weight  float64 `json:"weight" yaml:"weight"`
	Pass    Level   `json:"pass" yaml:"pass"`
	Partial Level   `json:"partial" yaml:"partial"`
	Fail    Level   `json:"fail" yaml:"fail"`
}

// Level is a scoring band for a criterion: what it means and the concrete
// indicators an evaluator looks for.
type Level struct {
	Description string   `json:"description" yaml:"description"`
	Indicators  []string `json:"indicators,omitempty" yaml:"indicators,omitempty"`
}

// Parse parses rich rubric YAML.
func Parse(data []byte) (*Rubric, error) {
	var r Rubric
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parse rubric: %w", err)
	}
	return &r, nil
}

// LoadRubric loads and parses the rich rubric for a canvas or document by name
// (with or without the ".rubric.yaml" suffix).
func LoadRubric(name string) (*Rubric, error) {
	content, err := Get(name)
	if err != nil {
		return nil, err
	}
	r, err := Parse([]byte(content))
	if err != nil {
		return nil, fmt.Errorf("rubric %q: %w", name, err)
	}
	return r, nil
}
