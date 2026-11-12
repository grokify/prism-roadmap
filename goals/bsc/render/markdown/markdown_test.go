package markdown

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/bsc"
	"github.com/grokify/prism-roadmap/goals/bsc/render"
)

var update = flag.Bool("update", false, "rewrite the golden files")

func TestRenderGolden(t *testing.T) {
	b, err := bsc.ReadFile("../../../../examples/product.bsc.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err := New().Render(b, nil)
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
	again, err := New().Render(b, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, again) {
		t.Error("rendering twice gave different output")
	}
}

func TestCompactOptionsLeaveOutInitiativesAndDescriptions(t *testing.T) {
	b, err := bsc.ReadFile("../../../../examples/product.bsc.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(b, render.CompactOptions())
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, unwanted := range []string{"## Initiatives", "Accounts that reach a reconciled month early stay."} {
		if strings.Contains(s, unwanted) {
			t.Errorf("compact output should not contain %q", unwanted)
		}
	}
	if !strings.Contains(s, "## Strategy Map") {
		t.Error("compact output should keep the strategy map")
	}
}

// The strategy map groups objectives by perspective, top first, and draws an
// arrow for each link, with positional node IDs and escaped labels.
func TestStrategyMap(t *testing.T) {
	b := &bsc.BSC{
		Perspectives: []bsc.Perspective{{ID: "c", Name: "Customer"}, {ID: "p", Name: "Process"}, {ID: "l", Name: "Learning"}},
		Objectives: []bsc.Objective{
			{ID: "keep", Name: `Keep "owners"`, PerspectiveID: "c"},
			{ID: "on-board", Name: "Onboard", PerspectiveID: "p", ContributesTo: []string{"keep", "missing"}},
		},
	}
	want := "```mermaid\nflowchart BT\n" +
		"  subgraph p1[\"Customer\"]\n    o1[\"Keep #quot;owners#quot;\"]\n  end\n" +
		"  subgraph p2[\"Process\"]\n    o2[\"Onboard\"]\n  end\n" +
		"  o2 --> o1\n```\n"
	if got := StrategyMap(b); got != want {
		t.Errorf("StrategyMap =\n%s\nwant\n%s", got, want)
	}
}

func TestHelpers(t *testing.T) {
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
	if New().Format() != "markdown" || New().FileExtension() != ".md" || Title(&bsc.BSC{}) != "Balanced Scorecard" {
		t.Error("Format, FileExtension, or Title is wrong")
	}
	if Timing("", "2027-01-01") != "by 2027-01-01" || Timing("2027-01-01", "") != "from 2027-01-01" || Value("6", "hours") != "6 hours" {
		t.Error("Timing or Value is wrong")
	}
	m := bsc.Measure{Target: "65", Unit: "%", Thresholds: &bsc.Thresholds{Green: "63", Red: "55"}}
	if got := Target(m); got != "65% (green at 63%, red at 55%)" {
		t.Errorf("Target = %q", got)
	}
	if got := Purpose(&bsc.BSC{Vision: " Be  kept ", Strategy: "Retain"}); len(got) != 2 || got[0].Text != "Be kept" {
		t.Errorf("Purpose = %v", got)
	}
}
