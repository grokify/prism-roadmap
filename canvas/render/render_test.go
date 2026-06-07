package render

import (
	"testing"

	"github.com/grokify/prism-roadmap/canvas"
)

// mockRenderer is a simple renderer for testing
type mockRenderer struct {
	format       Format
	extension    string
	supportedAll bool
}

func (m *mockRenderer) Format() Format                     { return m.format }
func (m *mockRenderer) FileExtension() string              { return m.extension }
func (m *mockRenderer) Supports(ct canvas.CanvasType) bool { return m.supportedAll }
func (m *mockRenderer) Render(c *canvas.Canvas, opts *Options) ([]byte, error) {
	return []byte("mock output"), nil
}

func TestRegistry(t *testing.T) {
	registry := NewRegistry()

	// Register mock renderers
	d2Mock := &mockRenderer{format: FormatD2, extension: ".d2", supportedAll: true}
	mermaidMock := &mockRenderer{format: FormatMermaid, extension: ".mmd", supportedAll: true}

	registry.Register(d2Mock)
	registry.Register(mermaidMock)

	// Test Get
	renderer, ok := registry.Get(FormatD2)
	if !ok {
		t.Error("Get(FormatD2) should return true")
	}
	if renderer.Format() != FormatD2 {
		t.Errorf("Format() = %v, want %v", renderer.Format(), FormatD2)
	}

	// Test missing format
	_, ok = registry.Get(FormatLit)
	if ok {
		t.Error("Get(FormatLit) should return false for unregistered format")
	}

	// Test Formats
	formats := registry.Formats()
	if len(formats) != 2 {
		t.Errorf("Formats() length = %v, want 2", len(formats))
	}
}

func TestRegistryRender(t *testing.T) {
	registry := NewRegistry()
	registry.Register(&mockRenderer{format: FormatD2, extension: ".d2", supportedAll: true})

	// Create test canvas
	c := canvas.NewOST(&canvas.OpportunitySolutionTree{
		Metadata: canvas.Metadata{ID: "test", Title: "Test"},
	})

	// Test render
	output, err := registry.Render(c, FormatD2, nil)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if string(output) != "mock output" {
		t.Errorf("Render() = %v, want 'mock output'", string(output))
	}

	// Test nil canvas
	_, err = registry.Render(nil, FormatD2, nil)
	if err == nil {
		t.Error("Render(nil) should return error")
	}

	// Test missing format
	_, err = registry.Render(c, FormatLit, nil)
	if err == nil {
		t.Error("Render() with unregistered format should return error")
	}
}

func TestSupportedFormats(t *testing.T) {
	registry := NewRegistry()

	// Register one that supports all, one that supports none
	registry.Register(&mockRenderer{format: FormatD2, supportedAll: true})
	registry.Register(&mockRenderer{format: FormatMermaid, supportedAll: false})

	supported := registry.SupportedFormats(canvas.CanvasTypeOST)
	if len(supported) != 1 {
		t.Errorf("SupportedFormats() length = %v, want 1", len(supported))
	}
	if supported[0] != FormatD2 {
		t.Errorf("Supported format = %v, want %v", supported[0], FormatD2)
	}
}

func TestFileExtensions(t *testing.T) {
	tests := []struct {
		format Format
		want   string
	}{
		{FormatD2, ".d2"},
		{FormatD2SVG, ".svg"},
		{FormatMermaid, ".mmd"},
		{FormatLit, ".json"},
		{FormatMarkdown, ".md"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			got := FileExtensions[tt.format]
			if got != tt.want {
				t.Errorf("FileExtensions[%v] = %v, want %v", tt.format, got, tt.want)
			}
		})
	}
}

func TestDefaultRegistry(t *testing.T) {
	// DefaultRegistry should be initialized
	if DefaultRegistry == nil {
		t.Fatal("DefaultRegistry is nil")
	}

	// Test global Register function
	mock := &mockRenderer{format: "test-format", supportedAll: true}
	Register(mock)

	renderer, ok := DefaultRegistry.Get("test-format")
	if !ok {
		t.Error("Global Register() should add to DefaultRegistry")
	}
	if renderer.Format() != "test-format" {
		t.Error("Registered renderer has wrong format")
	}
}
