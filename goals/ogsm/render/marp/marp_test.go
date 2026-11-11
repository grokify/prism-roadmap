package marp

import (
	"bytes"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/ogsm"
)

func TestRender(t *testing.T) {
	o, err := ogsm.ReadFile("../../../../examples/product.ogsm.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(o, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{
		"marp: true",
		`header: "Company | 2027-2029"`,
		`footer: "OGSM | Bookkeeping App OGSM 2027-2029"`,
		"## The Plan on One Page",
		"## Objective and Goals",
		"**By:** 2029-12-31",
		"## Strategy 1: Get every new account to a reconciled month in its first week",
		"| Ship guided bank import | Onboarding squad | 2027-01-04 to 2027-03-26 | In Progress |",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q", want)
		}
	}
	// Six slides (title, one page, objective and goals, three strategies)
	// need five breaks, plus the line that closes the front matter.
	if got := strings.Count(s, "\n---\n"); got != 6 {
		t.Errorf("want 6 separator lines, got %d", got)
	}
	again, err := New().Render(o, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(out, again) {
		t.Error("rendering twice gave different output")
	}
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
	if New().Format() != "marp" || New().FileExtension() != ".md" {
		t.Error("Format or FileExtension is wrong")
	}
}
