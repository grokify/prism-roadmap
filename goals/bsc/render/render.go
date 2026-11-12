// Package render defines the renderer interface and options for Balanced
// Scorecards. The markdown and marp subpackages implement it.
package render

import "github.com/grokify/prism-roadmap/goals/bsc"

// Renderer renders a Balanced Scorecard to an output format.
type Renderer interface {
	// Format returns the output format name, e.g. "markdown" or "marp".
	Format() string
	// FileExtension returns the file extension for this format, e.g. ".md".
	FileExtension() string
	// Render converts a Balanced Scorecard to the output format.
	Render(b *bsc.BSC, opts *Options) ([]byte, error)
}

// Options are rendering options common to all renderers.
type Options struct {
	// Theme is the renderer-specific theme name, e.g. "default" or "corporate".
	Theme string
	// IncludeDescriptions shows the descriptions of perspectives, themes,
	// objectives, and initiatives.
	IncludeDescriptions bool
	// IncludeStrategyMap shows the strategy map.
	IncludeStrategyMap bool
	// IncludeInitiatives shows the initiatives.
	IncludeInitiatives bool
}

// DefaultOptions returns the full rendering: descriptions, the strategy map,
// and initiatives.
func DefaultOptions() *Options {
	return &Options{Theme: "default", IncludeDescriptions: true, IncludeStrategyMap: true, IncludeInitiatives: true}
}

// CompactOptions returns the strategy map and the scorecard only: no
// descriptions or initiatives.
func CompactOptions() *Options {
	return &Options{Theme: "default", IncludeStrategyMap: true}
}
