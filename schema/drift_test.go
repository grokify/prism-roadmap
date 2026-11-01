package schema

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestCommittedSchemasMatchGenerator fails when a committed schema differs from
// what the generator writes, as happens after a dependency changes a type the
// schemas reflect. Regenerate with: cd schema && go run cmd/generate/main.go
func TestCommittedSchemasMatchGenerator(t *testing.T) {
	dir := t.TempDir()
	if err := NewGenerator().GenerateAll(dir); err != nil {
		t.Fatal(err)
	}
	generated, err := filepath.Glob(filepath.Join(dir, "*.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	if len(generated) == 0 {
		t.Fatal("the generator wrote no schemas")
	}
	for _, g := range generated {
		name := filepath.Base(g)
		want, err := os.ReadFile(g)
		if err != nil {
			t.Fatal(err)
		}
		got, err := os.ReadFile(name)
		if err != nil {
			t.Errorf("%s: generated but not committed: %v", name, err)
			continue
		}
		if !bytes.Equal(got, want) {
			t.Errorf("%s is out of date; regenerate with: cd schema && go run cmd/generate/main.go", name)
		}
	}
}
