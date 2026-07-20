package rubrics

import (
	"fmt"

	serubric "github.com/plexusone/structured-evaluation/rubric"
	"gopkg.in/yaml.v3"
)

// LoadRubric loads and parses the rich rubric for a canvas or document by name
// (with or without the ".rubric.yaml" suffix) into the canonical
// structured-evaluation RubricSet — the shared rubric spec used across the
// ecosystem. Canvas rubrics are authored against that spec so downstream tools
// (for example VisionSpec) consume them without any translation layer.
func LoadRubric(name string) (*serubric.RubricSet, error) {
	content, err := Get(name)
	if err != nil {
		return nil, err
	}
	var rs serubric.RubricSet
	if err := yaml.Unmarshal([]byte(content), &rs); err != nil {
		return nil, fmt.Errorf("rubric %q: %w", name, err)
	}
	return &rs, nil
}
