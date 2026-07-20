package svg

import (
	"strings"
	"testing"

	"github.com/grokify/prism-roadmap/canvas"
	"github.com/grokify/prism-roadmap/canvas/render"
)

func TestSupportsBMCAndOpportunitySpec(t *testing.T) {
	r := NewSVGRenderer()
	for _, ct := range []canvas.CanvasType{canvas.CanvasTypeBMC, canvas.CanvasTypeOpportunitySpec} {
		if !r.Supports(ct) {
			t.Errorf("Supports(%s) = false, want true", ct)
		}
	}
}

func TestRenderBMCGrid(t *testing.T) {
	c := &canvas.Canvas{
		Type: canvas.CanvasTypeBMC,
		BMC: &canvas.BusinessModelCanvas{
			Metadata:          canvas.Metadata{Title: "Acme BMC"},
			CustomerSegments:  []canvas.CustomerSegment{{Name: "SMB teams"}},
			ValuePropositions: []canvas.ValueProposition{{Description: "Faster onboarding"}},
			RevenueStreams:    []canvas.RevenueStream{{Description: "Subscription"}},
		},
	}

	out, err := NewSVGRenderer().Render(c, render.DefaultOptions())
	if err != nil {
		t.Fatalf("render BMC: %v", err)
	}
	svg := string(out)

	for _, want := range []string{
		"<svg", "</svg>", "Acme BMC",
		"Key Partnerships", "Value Propositions", "Customer Segments",
		"Cost Structure", "Revenue Streams",
		"SMB teams", "Faster onboarding", "Subscription",
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("BMC SVG missing %q", want)
		}
	}
}

func TestRenderOpportunitySpecGrid(t *testing.T) {
	c := &canvas.Canvas{
		Type: canvas.CanvasTypeOpportunitySpec,
		OpportunitySpec: &canvas.OpportunitySpec{
			Metadata:       canvas.Metadata{Title: "Acme Opportunity"},
			UserValue:      canvas.OSUserValue{ValueStatement: "Cut setup time in half"},
			Recommendation: canvas.OSRecommendation{Decision: "Proceed to MVP"},
		},
	}

	out, err := NewSVGRenderer().Render(c, render.DefaultOptions())
	if err != nil {
		t.Fatalf("render OpportunitySpec: %v", err)
	}
	svg := string(out)

	for _, want := range []string{
		"<svg", "</svg>", "Acme Opportunity",
		"1 · Users &amp; Problem", "4 · User Value", "12 · Recommendation",
		"Cut setup time in half", "Proceed to MVP",
	} {
		if !strings.Contains(svg, want) {
			t.Errorf("OpportunitySpec SVG missing %q", want)
		}
	}
}

func TestRenderBMCMissingDataErrors(t *testing.T) {
	c := &canvas.Canvas{Type: canvas.CanvasTypeBMC}
	if _, err := NewSVGRenderer().Render(c, render.DefaultOptions()); err == nil {
		t.Error("expected error rendering BMC with nil data")
	}
}
