// Package render provides interfaces and utilities for rendering V2MOM documents
// to various output formats.
package render

import (
	"github.com/grokify/prism-execution/goals/v2mom"
)

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Render converts a V2MOM to the target format.
	Render(v *v2mom.V2MOM, opts *Options) ([]byte, error)
	// Format returns the output format name (e.g., "marp", "confluence").
	Format() string
	// FileExtension returns the file extension for this format (e.g., ".md", ".html").
	FileExtension() string
}

// Options contains rendering options common to all renderers.
type Options struct {
	// Theme name (renderer-specific, e.g., "default", "corporate", "minimal")
	Theme string

	// Terminology mode: "v2mom", "okr", or "hybrid"
	Terminology string

	// Structure handling: override structure detection ("flat", "nested", "hybrid")
	Structure string

	// Include project/roadmap slides
	IncludeProjects bool

	// Include status indicators (colors, badges)
	IncludeStatus bool

	// FlattenMeasures combines global + nested measures into single view
	FlattenMeasures bool

	// Custom CSS (for Marp/HTML renderers)
	CustomCSS string

	// Additional metadata (renderer-specific)
	Metadata map[string]string
}

// DefaultOptions returns sensible default rendering options.
func DefaultOptions() *Options {
	return &Options{
		Theme:           "default",
		Terminology:     v2mom.TerminologyV2MOM,
		IncludeProjects: true,
		IncludeStatus:   true,
		FlattenMeasures: false,
		Metadata:        make(map[string]string),
	}
}

// GetTerminology returns the terminology to use, preferring options over V2MOM metadata.
func (o *Options) GetTerminology(v *v2mom.V2MOM) string {
	if o.Terminology != "" {
		return o.Terminology
	}
	return v.GetTerminology()
}

// GetStructure returns the structure to use, preferring options over V2MOM inference.
func (o *Options) GetStructure(v *v2mom.V2MOM) string {
	if o.Structure != "" {
		return o.Structure
	}
	return v.GetStructure()
}
