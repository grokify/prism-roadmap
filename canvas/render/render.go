// Package render provides interfaces and utilities for rendering canvas types
// to various output formats including D2, Mermaid, and Lit components.
package render

import (
	"errors"
	"sync"

	"github.com/grokify/prism-roadmap/canvas"
)

// Format identifies the output format.
type Format string

// Output format constants.
const (
	FormatD2       Format = "d2"       // D2 diagram language text
	FormatD2SVG    Format = "d2svg"    // D2 rendered to SVG (via d2 CLI)
	FormatSVG      Format = "svg"      // Native SVG output
	FormatMermaid  Format = "mermaid"  // Mermaid diagram text
	FormatLit      Format = "lit"      // Lit component JSON data
	FormatMarkdown Format = "markdown" // Markdown table layout
)

// FileExtensions maps formats to file extensions.
var FileExtensions = map[Format]string{
	FormatD2:       ".d2",
	FormatD2SVG:    ".svg",
	FormatSVG:      ".svg",
	FormatMermaid:  ".mmd",
	FormatLit:      ".json",
	FormatMarkdown: ".md",
}

// Renderer converts a canvas to a specific output format.
type Renderer interface {
	// Format returns the output format this renderer produces.
	Format() Format

	// FileExtension returns the file extension for this format.
	FileExtension() string

	// Supports returns true if this renderer supports the given canvas type.
	Supports(canvasType canvas.CanvasType) bool

	// Render converts a canvas to the output format.
	Render(c *canvas.Canvas, opts *Options) ([]byte, error)
}

// Registry holds available renderers and dispatches to them.
type Registry struct {
	mu        sync.RWMutex
	renderers map[Format]Renderer
}

// NewRegistry creates a new renderer registry.
func NewRegistry() *Registry {
	return &Registry{
		renderers: make(map[Format]Renderer),
	}
}

// Register adds a renderer to the registry.
func (r *Registry) Register(renderer Renderer) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.renderers[renderer.Format()] = renderer
}

// Get returns a renderer for the given format.
func (r *Registry) Get(format Format) (Renderer, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	renderer, ok := r.renderers[format]
	return renderer, ok
}

// Formats returns all registered formats.
func (r *Registry) Formats() []Format {
	r.mu.RLock()
	defer r.mu.RUnlock()
	formats := make([]Format, 0, len(r.renderers))
	for f := range r.renderers {
		formats = append(formats, f)
	}
	return formats
}

// Render renders a canvas to the specified format.
func (r *Registry) Render(c *canvas.Canvas, format Format, opts *Options) ([]byte, error) {
	if c == nil {
		return nil, errors.New("canvas is nil")
	}

	renderer, ok := r.Get(format)
	if !ok {
		return nil, errors.New("no renderer for format: " + string(format))
	}

	if !renderer.Supports(c.Type) {
		return nil, errors.New("renderer does not support canvas type: " + string(c.Type))
	}

	if opts == nil {
		opts = DefaultOptions()
	}

	return renderer.Render(c, opts)
}

// SupportedFormats returns formats that support the given canvas type.
func (r *Registry) SupportedFormats(canvasType canvas.CanvasType) []Format {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var supported []Format
	for format, renderer := range r.renderers {
		if renderer.Supports(canvasType) {
			supported = append(supported, format)
		}
	}
	return supported
}

// DefaultRegistry is the global renderer registry.
var DefaultRegistry = NewRegistry()

// Register registers a renderer with the default registry.
func Register(renderer Renderer) {
	DefaultRegistry.Register(renderer)
}

// Render renders a canvas using the default registry.
func Render(c *canvas.Canvas, format Format, opts *Options) ([]byte, error) {
	return DefaultRegistry.Render(c, format, opts)
}
