// Package render provides interfaces and utilities for rendering OKR documents
// to various output formats including Marp slides.
package render

import "github.com/grokify/prism-execution/goals/okr"

// Renderer defines the interface for output format renderers.
type Renderer interface {
	// Format returns the output format name (e.g., "marp").
	Format() string
	// FileExtension returns the file extension for this format (e.g., ".md").
	FileExtension() string
	// Render converts an OKR document to the output format.
	Render(doc *okr.OKRDocument, opts *Options) ([]byte, error)
}

// Options contains rendering options common to all renderers.
type Options struct {
	// Theme name (renderer-specific, e.g., "default", "corporate", "minimal")
	Theme string

	// IncludeRisks includes risks slides
	IncludeRisks bool

	// IncludeStatus shows status indicators and colors
	IncludeStatus bool

	// ShowScoreGrades shows letter grades for scores (A, B, C, etc.)
	ShowScoreGrades bool

	// ShowProgressBars shows visual progress bars
	ShowProgressBars bool

	// MaxObjectives limits objectives shown (0 = all)
	MaxObjectives int

	// MaxKeyResults limits key results per objective (0 = all)
	MaxKeyResults int

	// Custom CSS (for Marp/HTML renderers)
	CustomCSS string

	// Additional metadata (renderer-specific)
	Metadata map[string]string
}

// DefaultOptions returns sensible default rendering options.
func DefaultOptions() *Options {
	return &Options{
		Theme:            "default",
		IncludeRisks:     true,
		IncludeStatus:    true,
		ShowScoreGrades:  true,
		ShowProgressBars: true,
		MaxObjectives:    0,
		MaxKeyResults:    0,
		Metadata:         make(map[string]string),
	}
}

// ExecutiveOptions returns options for executive-focused slides (fewer details).
func ExecutiveOptions() *Options {
	return &Options{
		Theme:            "corporate",
		IncludeRisks:     true,
		IncludeStatus:    true,
		ShowScoreGrades:  true,
		ShowProgressBars: true,
		MaxObjectives:    5,
		MaxKeyResults:    3,
		Metadata:         make(map[string]string),
	}
}
