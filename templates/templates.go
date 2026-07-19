// Package templates exposes the canonical canvas and document templates as an
// embedded filesystem so downstream tools (for example VisionSpec) can consume
// them directly instead of duplicating copies.
package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed *.md
var templateFS embed.FS

// FS returns the embedded template filesystem. Callers that already have a
// filesystem-based loader (such as VisionSpec's templates.NewEmbedFSLoader) can
// wrap this directly.
func FS() embed.FS { return templateFS }

// Get returns the template content for a canvas or document by name, with or
// without the ".md" suffix (for example "bmc" or "bmc.md").
func Get(name string) (string, error) {
	filename := name
	if !strings.HasSuffix(filename, ".md") {
		filename += ".md"
	}
	data, err := templateFS.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("template %q: %w", name, err)
	}
	return string(data), nil
}

// Names returns the available template names without the ".md" suffix, sorted.
func Names() []string {
	entries, err := fs.ReadDir(templateFS, ".")
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".md") {
			names = append(names, strings.TrimSuffix(e.Name(), ".md"))
		}
	}
	sort.Strings(names)
	return names
}
