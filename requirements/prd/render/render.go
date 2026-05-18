// Package render provides interfaces and utilities for rendering PRD documents
// to various output formats including Marp slides.
package render

import "github.com/grokify/prism-execution/requirements/prd"

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Format returns the output format name (e.g., "marp").
	Format() string
	// FileExtension returns the file extension for this format (e.g., ".md").
	FileExtension() string
	// Render converts a PRD document to the output format.
	Render(doc *prd.Document, opts *Options) ([]byte, error)
}

// Options contains rendering options common to all renderers.
type Options struct {
	// Theme name (renderer-specific, e.g., "default", "corporate", "minimal")
	Theme string

	// IncludeGoals includes goals alignment section in slides
	IncludeGoals bool

	// IncludeRoadmap includes roadmap/timeline slides
	IncludeRoadmap bool

	// IncludeRisks includes risks slide
	IncludeRisks bool

	// IncludeRequirements includes detailed requirements slides
	IncludeRequirements bool

	// MaxPersonas limits the number of personas shown (0 = all)
	MaxPersonas int

	// MaxRequirements limits requirements per slide (0 = all)
	MaxRequirements int

	// Custom CSS (for Marp/HTML renderers)
	CustomCSS string

	// Additional metadata (renderer-specific)
	Metadata map[string]string
}

// DefaultOptions returns sensible default rendering options.
func DefaultOptions() *Options {
	return &Options{
		Theme:               "default",
		IncludeGoals:        true,
		IncludeRoadmap:      true,
		IncludeRisks:        true,
		IncludeRequirements: true,
		MaxPersonas:         5,
		MaxRequirements:     10,
		Metadata:            make(map[string]string),
	}
}

// ExecutiveOptions returns options for executive-focused slides (fewer details).
func ExecutiveOptions() *Options {
	return &Options{
		Theme:               "corporate",
		IncludeGoals:        true,
		IncludeRoadmap:      true,
		IncludeRisks:        true,
		IncludeRequirements: false, // Skip detailed requirements
		MaxPersonas:         3,
		MaxRequirements:     0,
		Metadata:            make(map[string]string),
	}
}
