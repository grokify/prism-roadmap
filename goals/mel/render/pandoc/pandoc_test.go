package pandoc

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/grokify/prism-roadmap/goals/mel"
)

func TestRender(t *testing.T) {
	p, err := mel.ReadFile("../../../../examples/product.mel.json")
	if err != nil {
		t.Fatal(err)
	}
	p.Metadata.Name = "Onboarding MEL Plan: owners' first year"
	p.Metadata.UpdatedAt = time.Date(2027, 1, 15, 9, 0, 0, 0, time.UTC)
	out, err := New().Render(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{
		"---\ntitle: 'Onboarding MEL Plan: owners'' first year'\n",
		"subtitle: 'Product · FY2027'\n",
		"author: 'Head of Product'\n",
		"date: '2027-01-15'\n",
		"geometry:\n  - margin=2cm\n",
		"toc: true\n",
		"---\n\n## Purpose\n",
		"## Results Framework",
		"## Indicator Reference Sheets",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q", want)
		}
	}
	if strings.Contains(s, "\n# ") {
		t.Error("the body should not repeat the title as a heading")
	}
	again, err := New().Render(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, again) {
		t.Error("rendering twice gave different output")
	}
}

func TestRenderWithoutMetadata(t *testing.T) {
	out, err := New().Render(&mel.Plan{}, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(out), "---\ntitle: 'MEL Plan'\ndocumentclass") {
		t.Errorf("front matter = %.80q", out)
	}
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
	if New().Format() != "pandoc" || New().FileExtension() != ".md" {
		t.Error("Format or FileExtension is wrong")
	}
}
