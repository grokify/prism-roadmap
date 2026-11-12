// Package render defines the renderer interface and options for MEL plans.
// The markdown and pandoc subpackages implement it.
package render

import "github.com/grokify/prism-roadmap/goals/mel"

// Renderer renders a MEL plan to an output format.
type Renderer interface {
	// Format returns the output format name, e.g. "markdown" or "pandoc".
	Format() string
	// FileExtension returns the file extension for this format, e.g. ".md".
	FileExtension() string
	// Render converts a MEL plan to the output format.
	Render(p *mel.Plan, opts *Options) ([]byte, error)
}

// Options are rendering options common to all renderers.
type Options struct {
	// IncludeReferenceSheets adds a full reference sheet for each indicator
	// after the indicator summary.
	IncludeReferenceSheets bool
	// IncludeDescriptions shows result descriptions and the rationale of each
	// learning question.
	IncludeDescriptions bool
	// OmitHeader leaves out the document title and the metadata line under it,
	// for a caller that embeds the output under a heading of its own. Markdown only.
	OmitHeader bool
	// HeadingOffset pushes every heading down this many levels (capped at level 6),
	// for embedding the output inside a larger document. Markdown only.
	HeadingOffset int
}

// DefaultOptions returns the full plan: reference sheets and descriptions.
func DefaultOptions() *Options {
	return &Options{IncludeReferenceSheets: true, IncludeDescriptions: true}
}

// SummaryOptions returns a short plan: the indicator summary without
// reference sheets, and no descriptions.
func SummaryOptions() *Options {
	return &Options{}
}
