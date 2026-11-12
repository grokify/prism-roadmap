package markdown

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/goals/dmaic"
	"github.com/grokify/prism-roadmap/goals/dmaic/render"
)

var update = flag.Bool("update", false, "rewrite the golden files")

func TestRenderGolden(t *testing.T) {
	doc, err := dmaic.ReadFile("../../../../examples/manufacturing.dmaic.json")
	if err != nil {
		t.Fatal(err)
	}
	for _, tt := range []struct {
		golden string
		opts   *render.Options
	}{
		{"manufacturing.md", nil},
		{"manufacturing-six-sigma.md", render.SixSigmaOptions()},
	} {
		t.Run(tt.golden, func(t *testing.T) {
			got, err := New().Render(doc, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			golden := filepath.Join("testdata", tt.golden)
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
			again, err := New().Render(doc, tt.opts)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, again) {
				t.Error("rendering twice gave different output")
			}
		})
	}
}

func TestRenderPortable(t *testing.T) {
	doc, err := dmaic.ReadFile("../../../../examples/manufacturing.dmaic.json")
	if err != nil {
		t.Fatal(err)
	}
	out, err := New().Render(doc, render.SixSigmaOptions())
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if strings.HasPrefix(s, "---") {
		t.Error("output should not start with front matter")
	}
	for _, unwanted := range []string{`\textcolor`, `\newpage`, `\pagebreak`, "```{=latex}", "$\\"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("output should not contain %q", unwanted)
		}
	}
}

func TestRenderOptions(t *testing.T) {
	doc := &dmaic.DMAICDocument{
		Categories: []dmaic.Category{{
			Name: "Security",
			Metrics: []dmaic.Metric{{
				ID: "vulns", Name: "Open | critical vulnerabilities", Phase: dmaic.PhaseAnalyze, Status: dmaic.StatusYellow,
				RootCauses: []dmaic.RootCause{{Description: "Insufficient testing", Category: "Process", Impact: "High", Validated: true}},
			}},
		}},
		Initiatives: []dmaic.Initiative{{Name: "Add a security scan", MetricIDs: []string{"vulns", "other"}}},
	}
	out, err := New().Render(doc, &render.Options{})
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	for _, want := range []string{"# DMAIC Metrics Report", "## Security", `Open \| critical vulnerabilities`, "### Open | critical vulnerabilities (Yellow)", "Insufficient testing"} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q:\n%s", want, s)
		}
	}
	for _, unwanted := range []string{"## Root Cause Analysis", "## Improvement Initiatives", "## DMAIC Phase Distribution", "## Appendix"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("output should not contain %q:\n%s", unwanted, s)
		}
	}

	out, err = New().Render(doc, &render.Options{IncludeRootCauses: true, IncludeInitiatives: true, GroupByPhase: true})
	if err != nil {
		t.Fatal(err)
	}
	s = string(out)
	for _, want := range []string{"## Root Cause Analysis", "## Improvement Initiatives", "**Linked metrics:** Open | critical vulnerabilities, `other`", "| **A** - Analyze | Analyze root causes | 1 |"} {
		if !strings.Contains(s, want) {
			t.Errorf("output lacks %q:\n%s", want, s)
		}
	}
}

func TestRenderEmbedded(t *testing.T) {
	doc, err := dmaic.ReadFile("../../../../examples/manufacturing.dmaic.json")
	if err != nil {
		t.Fatal(err)
	}
	opts := render.DefaultOptions()
	opts.OmitHeader = true
	opts.HeadingOffset = 1
	out, err := New().Render(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	s := string(out)
	if !strings.HasPrefix(s, "### Overview\n\n") {
		t.Errorf("embedded output should start with %q:\n%s", "### Overview", s)
	}
	for _, unwanted := range []string{"# Manufacturing Quality Metrics", "**Team:**"} {
		if strings.Contains(s, unwanted) {
			t.Errorf("embedded output should not contain %q", unwanted)
		}
	}
	for _, line := range strings.Split(s, "\n") {
		if strings.HasPrefix(line, "# ") || strings.HasPrefix(line, "## ") {
			t.Errorf("heading %q was not pushed down", line)
		}
	}
	if !strings.Contains(s, "\n### Production Quality\n") {
		t.Errorf("embedded output lacks the category heading at level 3:\n%s", s)
	}

	opts.HeadingOffset = 0
	flat, err := New().Render(doc, opts)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(flat), "## Overview\n\n") {
		t.Errorf("OmitHeader alone should start with %q:\n%s", "## Overview", flat)
	}
}

func TestRenderNil(t *testing.T) {
	if _, err := New().Render(nil, nil); err == nil {
		t.Error("Render(nil) should fail")
	}
}

func TestValueAndHealth(t *testing.T) {
	for _, tt := range []struct {
		v    float64
		unit string
		want string
	}{
		{94.2, "%", "94.2%"},
		{3400, "dpmo", "3400"},
		{1.333, "", "1.33"},
	} {
		if got := Value(tt.v, tt.unit); got != tt.want {
			t.Errorf("Value(%v, %q) = %q, want %q", tt.v, tt.unit, got, tt.want)
		}
	}
	if got := Health(0.75); got != "75%" {
		t.Errorf("Health(0.75) = %q", got)
	}
	if StatusText("") != "Unknown" || StatusText(dmaic.StatusRed) != "Red" {
		t.Error("StatusText gave the wrong name")
	}
	if PhaseAbbrev(dmaic.PhaseControl) != "C" || PhaseAbbrev("") != "?" {
		t.Error("PhaseAbbrev gave the wrong letter")
	}
	if TrendText(dmaic.TrendLowerBetter) != "Lower is better" || TrendText("") != "N/A" || TrendText("custom") != "custom" {
		t.Error("TrendText gave the wrong description")
	}
}

func TestFormat(t *testing.T) {
	r := New()
	if r.Format() != "markdown" || r.FileExtension() != ".md" {
		t.Errorf("Format() = %q, FileExtension() = %q", r.Format(), r.FileExtension())
	}
}
