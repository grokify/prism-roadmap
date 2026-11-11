package markdown

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/mel"
	"github.com/grokify/prism-roadmap/goals/mel/render"
)

var update = flag.Bool("update", false, "rewrite the golden files")

func TestRenderGolden(t *testing.T) {
	p, err := mel.ReadFile("../../../../examples/product.mel.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err := New().Render(p, nil)
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
	again, err := New().Render(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, again) {
		t.Error("rendering twice gave different output")
	}
}

func TestSummaryLeavesOutSheetsAndDescriptions(t *testing.T) {
	p, err := mel.ReadFile("../../../../examples/product.mel.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(p, render.SummaryOptions())
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, unwanted := range []string{"## Indicator Reference Sheets", "Why it matters", "Retention past the first quarter"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("summary should not contain %q", unwanted)
		}
	}
	for _, want := range []string{"## Indicators", "## Learning Agenda", "| How it will be answered |"} {
		if !strings.Contains(s, want) {
			t.Errorf("summary lacks %q", want)
		}
	}
}

// A cycle is a validation error, but rendering one must still finish.
func TestRenderSurvivesACycle(t *testing.T) {
	p := &mel.Plan{Results: []mel.Result{
		{ID: "a", Statement: "A", Level: mel.LevelImpact, ContributesTo: []string{"b"}},
		{ID: "b", Statement: "B", Level: mel.LevelOutput, ContributesTo: []string{"a"}},
		{ID: "top", Statement: "Top", Level: mel.LevelImpact},
		{ID: "c", Statement: "C", Level: mel.LevelOutput, ContributesTo: []string{"top", "a"}},
	}}
	if _, err := New().Render(p, nil); err != nil {
		t.Fatal(err)
	}
}

// A result shared by two branches is drawn in full once and referred to after.
func TestSharedResultIsDrawnOnce(t *testing.T) {
	p := &mel.Plan{Results: []mel.Result{
		{ID: "a", Statement: "A", Level: mel.LevelImpact},
		{ID: "b", Statement: "B", Level: mel.LevelImpact},
		{ID: "s", Statement: "Shared", Level: mel.LevelShortTerm, ContributesTo: []string{"a", "b"}},
		{ID: "o", Statement: "Out", Level: mel.LevelOutput, ContributesTo: []string{"s"}},
	}}
	out, err := New().Render(p, nil)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if strings.Count(s, "Out (`o`)") != 1 {
		t.Errorf("the shared branch's child should be drawn once:\n%s", s)
	}
	if !strings.Contains(s, "  - **Short-term outcome:** Shared (`s`), shown above\n") {
		t.Errorf("the second appearance should refer back:\n%s", s)
	}
}

func TestRenderNil(t *testing.T) {
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
}

func TestHelpers(t *testing.T) {
	if Timing("2027-01-01", "2027-02-01") != "2027-01-01 to 2027-02-01" || Timing("2027-01-01", "") != "from 2027-01-01" ||
		Timing("", "2027-02-01") != "by 2027-02-01" || Timing("", "") != "" {
		t.Error("Timing is wrong")
	}
	if Value("45", "%") != "45%" || Value("6", "hours") != "6 hours" || Value("7", "") != "7" || Value("", "%") != "" {
		t.Error("Value is wrong")
	}
	final, ok := FinalTarget(mel.Indicator{Targets: []mel.Target{{Value: "1", Date: "2027-12-31"}, {Value: "2", Date: "2027-06-30"}}})
	if !ok || final.Value != "1" {
		t.Errorf("FinalTarget = %+v", final)
	}
	if _, ok := FinalTarget(mel.Indicator{}); ok {
		t.Error("FinalTarget of no targets should be false")
	}
	if Title(&mel.Plan{}) != "MEL Plan" {
		t.Error("Title default is wrong")
	}
	r := New()
	if r.Format() != "markdown" || r.FileExtension() != ".md" {
		t.Error("Format or FileExtension is wrong")
	}
}
