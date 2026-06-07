// Package export provides interfaces for exporting canvases to external systems.
package export

import (
	"context"
	"fmt"
	"sync"

	"github.com/grokify/prism-roadmap/canvas"
)

// CanvasProvider exports canvases to external systems like Aha.io, Miro, etc.
type CanvasProvider interface {
	// Name returns the provider identifier (e.g., "aha", "miro").
	Name() string

	// SupportedTypes returns which canvas types this provider supports.
	SupportedTypes() []canvas.CanvasType

	// CreateCanvas creates a new canvas in the external system.
	// Returns the external system's ID for the created canvas.
	CreateCanvas(ctx context.Context, c *canvas.Canvas) (string, error)

	// UpdateCanvas updates an existing canvas in the external system.
	UpdateCanvas(ctx context.Context, externalID string, c *canvas.Canvas) error

	// GetCanvas retrieves a canvas from the external system.
	// The returned canvas may have limited data depending on what the external system stores.
	GetCanvas(ctx context.Context, externalID string) (*canvas.Canvas, error)

	// DeleteCanvas removes a canvas from the external system.
	DeleteCanvas(ctx context.Context, externalID string) error
}

// Registry holds registered canvas providers.
type Registry struct {
	mu        sync.RWMutex
	providers map[string]CanvasProvider
}

// NewRegistry creates a new provider registry.
func NewRegistry() *Registry {
	return &Registry{
		providers: make(map[string]CanvasProvider),
	}
}

// Register adds a provider to the registry.
func (r *Registry) Register(p CanvasProvider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[p.Name()] = p
}

// Get returns a provider by name.
func (r *Registry) Get(name string) (CanvasProvider, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	p, ok := r.providers[name]
	return p, ok
}

// List returns all registered provider names.
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// Export creates a canvas using the named provider.
func (r *Registry) Export(ctx context.Context, providerName string, c *canvas.Canvas) (string, error) {
	p, ok := r.Get(providerName)
	if !ok {
		return "", fmt.Errorf("provider %q not registered", providerName)
	}
	return p.CreateCanvas(ctx, c)
}

// Sync updates a canvas using the named provider.
func (r *Registry) Sync(ctx context.Context, providerName string, externalID string, c *canvas.Canvas) error {
	p, ok := r.Get(providerName)
	if !ok {
		return fmt.Errorf("provider %q not registered", providerName)
	}
	return p.UpdateCanvas(ctx, externalID, c)
}
