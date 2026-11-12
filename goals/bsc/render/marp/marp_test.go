package marp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/bsc"
	"github.com/grokify/prism-roadmap/goals/bsc/render"
)

func TestRender(t *testing.T) {
	b, err := bsc.ReadFile("../../../../examples/product.bsc.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{
		"marp: true",
		`header: "Bookkeeping App | FY2027"`,
		`footer: "Balanced Scorecard | Bookkeeping App Balanced Scorecard FY2027"`,
		"## Mission, Vision, and Strategy",
		"## Strategy Map",
		"| **Learning and Growth** | Build the team's skills in bank data integrations<br>Put customer feedback in front of every squad |",
		"## Internal Process",
		"| Make bank sync reliable | Bank syncs that complete without help | lead | 99% | 94% | Behind |",
		"| Ship guided bank import | Get every new account to a reconciled month in its first week; Keep small businesses using the app |",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q", want)
		}
	}
	// Nine slides (title, purpose, strategy map, four perspectives,
	// initiatives) need seven breaks, plus the line that closes the front
	// matter.
	if got := strings.Count(s, "\n---\n"); got != 8 {
		t.Errorf("want 8 separator lines, got %d", got)
	}
	again, err := New().Render(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, again) {
		t.Error("rendering twice gave different output")
	}
	compact, err := New().Render(b, render.CompactOptions())
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(compact), "## Initiatives") {
		t.Error("compact slides should leave out the initiatives")
	}
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
	if New().Format() != "marp" || New().FileExtension() != ".md" {
		t.Error("Format or FileExtension is wrong")
	}
}
