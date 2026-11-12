// Package render defines the renderer interface and options for VSEM
// documents. The markdown and marp subpackages implement it.
package render

import "github.com/grokify/prism-roadmap/goals/vsem"

// Renderer renders a VSEM to an output format.
type Renderer interface {
	// Format returns the output format name, e.g. "markdown" or "marp".
	Format() string
	// FileExtension returns the file extension for this format, e.g. ".md".
	FileExtension() string
	// Render converts a VSEM to the output format.
	Render(v *vsem.VSEM, opts *Options) ([]byte, error)
}

// Options are rendering options common to all renderers.
type Options struct {
	// Theme is the renderer-specific theme name, e.g. "default" or "corporate".
	Theme string
	// IncludeDescriptions shows the descriptions of strategies, execution
	// items, and metrics.
	IncludeDescriptions bool
	// IncludeResources shows the resources each strategy commits.
	IncludeResources bool
	// OmitHeader leaves out the document title and the metadata line under it,
	// for a caller that embeds the output under a heading of its own. Markdown only.
	OmitHeader bool
	// HeadingOffset pushes every heading down this many levels (capped at level 6),
	// for embedding the output inside a larger document. Markdown only.
	HeadingOffset int
}

// DefaultOptions returns the full rendering: descriptions and resources.
func DefaultOptions() *Options {
	return &Options{
		Theme:               "default",
		IncludeDescriptions: true,
		IncludeResources:    true,
	}
}

// OnePageOptions returns options for a single-page summary: names, owners,
// timing, and targets, without descriptions.
func OnePageOptions() *Options {
	return &Options{
		Theme:            "default",
		IncludeResources: true,
	}
}
