// Package render defines the renderer interface and options for OGSM plans.
// The markdown and marp subpackages implement it.
package render

import "github.com/grokify/prism-roadmap/goals/ogsm"

// Renderer renders an OGSM to an output format.
type Renderer interface {
	// Format returns the output format name, e.g. "markdown" or "marp".
	Format() string
	// FileExtension returns the file extension for this format, e.g. ".md".
	FileExtension() string
	// Render converts an OGSM to the output format.
	Render(o *ogsm.OGSM, opts *Options) ([]byte, error)
}

// Options are rendering options common to all renderers.
type Options struct {
	// Theme is the renderer-specific theme name, e.g. "default" or "corporate".
	Theme string
	// IncludeDescriptions shows the descriptions of the objective, goals,
	// strategies, measures, and actions.
	IncludeDescriptions bool
	// IncludeActions shows the actions under each strategy.
	IncludeActions bool
}

// DefaultOptions returns the full rendering: descriptions and actions.
func DefaultOptions() *Options {
	return &Options{Theme: "default", IncludeDescriptions: true, IncludeActions: true}
}

// OnePageOptions returns the four parts only: no descriptions or actions.
func OnePageOptions() *Options {
	return &Options{Theme: "default"}
}
