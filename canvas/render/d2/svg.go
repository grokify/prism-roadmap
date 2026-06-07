package d2

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

// D2SVGRenderer renders canvas types to SVG via the D2 CLI.
type D2SVGRenderer struct {
	d2Path string // Path to d2 CLI binary
}

// NewD2SVGRenderer creates a new D2-to-SVG renderer.
// Returns nil if d2 CLI is not found.
func NewD2SVGRenderer() *D2SVGRenderer {
	path, err := exec.LookPath("d2")
	if err != nil {
		return nil
	}
	return &D2SVGRenderer{d2Path: path}
}

// NewD2SVGRendererWithPath creates a D2SVGRenderer with a specific d2 path.
func NewD2SVGRendererWithPath(d2Path string) *D2SVGRenderer {
	return &D2SVGRenderer{d2Path: d2Path}
}

// Format returns the output format.
func (r *D2SVGRenderer) Format() render.Format {
	return render.FormatD2SVG
}

// FileExtension returns the file extension for SVG files.
func (r *D2SVGRenderer) FileExtension() string {
	return ".svg"
}

// Supports returns true for all canvas types.
func (r *D2SVGRenderer) Supports(canvasType canvas.CanvasType) bool {
	return true
}

// Available returns true if the d2 CLI is available.
func (r *D2SVGRenderer) Available() bool {
	return r != nil && r.d2Path != ""
}

// Render converts a canvas to SVG via D2.
func (r *D2SVGRenderer) Render(c *canvas.Canvas, opts *render.Options) ([]byte, error) {
	if !r.Available() {
		return nil, errors.New("d2 CLI not available")
	}

	if opts == nil {
		opts = render.DefaultOptions()
	}

	// First, generate D2 text
	d2Renderer := NewD2Renderer()
	d2Text, err := d2Renderer.Render(c, opts)
	if err != nil {
		return nil, fmt.Errorf("generating D2: %w", err)
	}

	// Build d2 command arguments
	args := []string{}

	// Add theme if specified
	if opts.D2Theme > 0 {
		args = append(args, "--theme", strconv.Itoa(opts.D2Theme))
	}

	// Add layout if specified
	if opts.D2Layout != "" {
		args = append(args, "--layout", opts.D2Layout)
	}

	// Read from stdin, write to stdout
	args = append(args, "-", "-")

	// Execute d2
	// #nosec G204 -- d2Path is validated via exec.LookPath, args are from controlled Options
	cmd := exec.Command(r.d2Path, args...)
	cmd.Stdin = bytes.NewReader(d2Text)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("d2 CLI error: %w: %s", err, stderr.String())
	}

	return stdout.Bytes(), nil
}

func init() {
	if renderer := NewD2SVGRenderer(); renderer != nil {
		render.Register(renderer)
	}
}
