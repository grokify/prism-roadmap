package markdown

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/ogsm"
	"github.com/grokify/prism-roadmap/goals/ogsm/render"
)

var update = flag.Bool("update", false, "rewrite the golden files")

func TestRenderGolden(t *testing.T) {
	o, err := ogsm.ReadFile("../../../../examples/product.ogsm.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err := New().Render(o, nil)
	if err != nil {
		t.Fatal(err)
	}
	golden := filepath.Join("testdata", "product.md")
	if *update {
		if err := os.MkdirAll("testdata", 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(golden, got, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	want, err := os.ReadFile(golden)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Errorf("output differs from %s; run go test -update to accept\n%s", golden, got)
	}
	again, err := New().Render(o, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, again) {
		t.Error("rendering twice gave different output")
	}
}

func TestOnePageOptionsLeaveOutActionsAndDescriptions(t *testing.T) {
	o, err := ogsm.ReadFile("../../../../examples/product.ogsm.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(o, render.OnePageOptions())
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, unwanted := range []string{"**Actions**", "Growth has come from sign-ups"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("one-page output should not contain %q", unwanted)
		}
	}
}

// The one-page table has a row per goal or strategy, whichever is more, with
// the objective in the first row only.
func TestOnePage(t *testing.T) {
	o := &ogsm.OGSM{
		Objective:  ogsm.Objective{Statement: "Win | keep"},
		Goals:      []ogsm.Goal{{ID: "g", Name: "Retention", Target: "65", Unit: "%", DueDate: "2029-12-31"}},
		Strategies: []ogsm.Strategy{{ID: "s1", Name: "Onboard"}, {ID: "s2", Name: "Partners"}},
		Measures:   []ogsm.Measure{{ID: "m", Name: "Connect", StrategyID: "s1", Target: "75", Unit: "%"}},
	}
	got := OnePage(o)
	want := "| Objective | Goals | Strategies | Measures |\n" +
		"|-----------|-------|------------|----------|\n" +
		"| Win \\| keep | Retention: 65% by 2029-12-31 | Onboard | Connect: 75% |\n" +
		"|  |  | Partners |  |\n"
	if got != want {
		t.Errorf("OnePage =\n%s\nwant\n%s", got, want)
	}
}

func TestRenderNil(t *testing.T) {
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
	if New().Format() != "markdown" || New().FileExtension() != ".md" || Title(&ogsm.OGSM{}) != "OGSM" {
		t.Error("Format, FileExtension, or Title is wrong")
	}
	if Timing("", "2027-01-01") != "by 2027-01-01" || Value("6", "hours") != "6 hours" {
		t.Error("Timing or Value is wrong")
	}
}
