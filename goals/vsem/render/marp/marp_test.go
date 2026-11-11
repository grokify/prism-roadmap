package marp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/vsem"
)

func TestRender(t *testing.T) {
	v, err := vsem.ReadFile("../../../../examples/product.vsem.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(v, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{
		"marp: true",
		`header: "Product | FY2027"`,
		`footer: "VSEM | Bookkeeping App FY2027"`,
		"<!-- _class: title -->",
		"## The Plan on One Page",
		"> Every small business closes its books in an hour a month.",
		"## Strategy 1: Get new accounts to a reconciled month in their first week",
		"**Resources:** Engineering: 2 squads; Design: 1 designer",
		"| Guided bank import | Onboarding squad | 2027-01-04 to 2027-03-26 | In Progress |",
		"| New accounts that connect a bank in their first session | Head of Onboarding | 45% | 70% | 51% | On Track |",
		"## Metrics for the Vision",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q", want)
		}
	}
	if strings.HasSuffix(strings.TrimSpace(s), "---") {
		t.Error("output ends with an empty slide")
	}
	// Seven slides (title, overview, vision, three strategies, vision
	// metrics) need six breaks, plus the line that closes the front matter.
	if got := strings.Count(s, "\n---\n"); got != 7 {
		t.Errorf("want 7 separator lines, got %d", got)
	}
	again, err := New().Render(v, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, again) {
		t.Error("rendering twice gave different output")
	}
}

func TestRenderNil(t *testing.T) {
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
}

func TestFormat(t *testing.T) {
	r := New()
	if r.Format() != "marp" || r.FileExtension() != ".md" {
		t.Errorf("Format() = %q, FileExtension() = %q", r.Format(), r.FileExtension())
	}
}
