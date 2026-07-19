// Package rubrics exposes the canonical canvas and document evaluation rubrics
// as an embedded filesystem so downstream tools (for example VisionSpec) can
// consume them directly instead of duplicating copies.
package rubrics

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed *.rubric.yaml
var rubricFS embed.FS

const rubricSuffix = ".rubric.yaml"

// FS returns the embedded rubric filesystem. Callers with a filesystem-based
// loader (such as VisionSpec's rubrics.NewEmbedFSLoader) can wrap this directly.
func FS() embed.FS { return rubricFS }

// Get returns the rubric content for a canvas or document by name, with or
// without the ".rubric.yaml" suffix (for example "bmc" or "bmc.rubric.yaml").
func Get(name string) (string, error) {
	filename := name
	if !strings.HasSuffix(filename, rubricSuffix) {
		filename += rubricSuffix
	}
	data, err := rubricFS.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("rubric %q: %w", name, err)
	}
	return string(data), nil
}

// Names returns the available rubric names without the ".rubric.yaml" suffix, sorted.
func Names() []string {
	entries, err := fs.ReadDir(rubricFS, ".")
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), rubricSuffix) {
			names = append(names, strings.TrimSuffix(e.Name(), rubricSuffix))
		}
	}
	sort.Strings(names)
	return names
}
