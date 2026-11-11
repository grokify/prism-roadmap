package markdown

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/vsem"
	"github.com/grokify/prism-roadmap/goals/vsem/render"
)

var update = flag.Bool("update", false, "rewrite the golden files")

func TestRenderGolden(t *testing.T) {
	v, err := vsem.ReadFile("../../../../examples/product.vsem.json")
	if err != nil {
		t.Fatal(err)
	}
	got, err := New().Render(v, nil)
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
	again, err := New().Render(v, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, again) {
		t.Error("rendering twice gave different output")
	}
}

func TestRenderOptions(t *testing.T) {
	v := &vsem.VSEM{
		Mission:    "Take the busywork out of bookkeeping.",
		Strategies: []vsem.Strategy{{ID: "s", Name: "Onboard", Description: "Hidden detail", Resources: []vsem.Resource{{Name: "Engineering", Allocation: "2 squads"}}}},
		Execution:  []vsem.Execution{{ID: "e", Name: "Import | export", StrategyID: "s"}},
	}
	out, err := New().Render(v, &render.Options{})
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{"# VSEM", "**Mission:** Take the busywork", `Import \| export`} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q:\n%s", want, s)
		}
	}
	for _, unwanted := range []string{"Hidden detail", "Engineering", "Metrics for the Vision"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("output should not contain %q:\n%s", unwanted, s)
		}
	}
}

func TestRenderNil(t *testing.T) {
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
}

func TestTimingAndValue(t *testing.T) {
	for _, tt := range []struct{ start, end, want string }{
		{"2027-01-01", "2027-02-01", "2027-01-01 to 2027-02-01"},
		{"2027-01-01", "", "from 2027-01-01"},
		{"", "2027-02-01", "by 2027-02-01"},
		{"", "", ""},
	} {
		if got := Timing(tt.start, tt.end); got != tt.want {
			t.Errorf("Timing(%q, %q) = %q, want %q", tt.start, tt.end, got, tt.want)
		}
	}
	for _, tt := range []struct{ v, unit, want string }{
		{"45", "%", "45%"},
		{"2,000", "accounts", "2,000 accounts"},
		{"7", "", "7"},
		{"", "%", ""},
	} {
		if got := Value(tt.v, tt.unit); got != tt.want {
			t.Errorf("Value(%q, %q) = %q, want %q", tt.v, tt.unit, got, tt.want)
		}
	}
}

func TestFormat(t *testing.T) {
	r := New()
	if r.Format() != "markdown" || r.FileExtension() != ".md" {
		t.Errorf("Format() = %q, FileExtension() = %q", r.Format(), r.FileExtension())
	}
}
