package assessment

import (
	"testing"
	"time"

	"github.com/plexusone/structured-evaluation/claims"
)

func TestNewEvidenceDefaults(t *testing.T) {
	e := NewEvidence("EV-001", "72% of customers run in AWS")
	if e.ID != "EV-001" {
		t.Errorf("ID = %q, want EV-001", e.ID)
	}
	if e.Claim.Text != "72% of customers run in AWS" {
		t.Errorf("Claim.Text = %q", e.Claim.Text)
	}
	if e.Claim.Category != claims.ClaimTechnicalFinding {
		t.Errorf("Claim.Category = %q, want default ClaimTechnicalFinding", e.Claim.Category)
	}
	if e.Claim.Verdict != claims.VerdictUnverified {
		t.Errorf("Claim.Verdict = %q, want VerdictUnverified before a source is attached", e.Claim.Verdict)
	}
}

func TestEvidenceWithSourceComputesVerdict(t *testing.T) {
	tests := []struct {
		name        string
		reliability claims.ReliabilityTier
		wantVerdict claims.Verdict
	}{
		{"authoritative reliable", claims.ReliabilityAuthoritative, claims.VerdictVerified},
		{"high reliability", claims.ReliabilityHigh, claims.VerdictVerified},
		{"medium requires review", claims.ReliabilityMedium, claims.VerdictNeedsReview},
		{"low is rejected", claims.ReliabilityLow, claims.VerdictRejected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEvidence("EV-002", "claim text").
				WithSource("https://example.com/doc", EvidenceSystemGoogleDocs, claims.ExternalCommunity, tt.reliability)
			if e.Claim.Verdict != tt.wantVerdict {
				t.Errorf("Verdict = %q, want %q", e.Claim.Verdict, tt.wantVerdict)
			}
			if e.SourceURI() != "https://example.com/doc" {
				t.Errorf("SourceURI() = %q", e.SourceURI())
			}
		})
	}
}

func TestEvidenceFluentChain(t *testing.T) {
	captured := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)
	e := NewEvidence("EV-003", "Customer X contract requires FedRAMP by Q3 FY27").
		WithCategory(claims.ClaimRiskAssessment).
		WithSource("https://docs.google.com/document/d/abc", EvidenceSystemGoogleDocs, claims.ExternalCommunity, claims.ReliabilityHigh).
		WithExcerpt("Section 4.2: Provider shall achieve FedRAMP Moderate...").
		WithCapturedAt(captured).
		WithCapturedBy("pm:jwang").
		WithSensitivity(SensitivityRestricted)

	if e.Claim.Category != claims.ClaimRiskAssessment {
		t.Errorf("Category = %q", e.Claim.Category)
	}
	if e.Excerpt() != "Section 4.2: Provider shall achieve FedRAMP Moderate..." {
		t.Errorf("Excerpt() = %q", e.Excerpt())
	}
	if e.CapturedAtTime() == nil || !e.CapturedAtTime().Equal(captured) {
		t.Errorf("CapturedAtTime() = %v, want %v", e.CapturedAtTime(), captured)
	}
	if e.CapturedBy != "pm:jwang" {
		t.Errorf("CapturedBy = %q", e.CapturedBy)
	}
	if e.Sensitivity != SensitivityRestricted {
		t.Errorf("Sensitivity = %q", e.Sensitivity)
	}
}

func TestWithExcerptAndCapturedAtNoOpWithoutSource(t *testing.T) {
	e := NewEvidence("EV-004", "unsourced claim").
		WithExcerpt("should not be set").
		WithCapturedAt(time.Now())

	if e.Excerpt() != "" {
		t.Errorf("Excerpt() = %q, want empty (no source attached)", e.Excerpt())
	}
	if e.CapturedAtTime() != nil {
		t.Errorf("CapturedAtTime() = %v, want nil (no source attached)", e.CapturedAtTime())
	}
}

func TestSensitivityRenderable(t *testing.T) {
	tests := []struct {
		s    Sensitivity
		want bool
	}{
		{SensitivityPublic, true},
		{SensitivityInternal, true},
		{SensitivityRestricted, false},
		{"", false}, // zero value defaults to not-renderable
	}
	for _, tt := range tests {
		if got := tt.s.Renderable(); got != tt.want {
			t.Errorf("Sensitivity(%q).Renderable() = %v, want %v", tt.s, got, tt.want)
		}
	}
}

func TestRenderableExcerpt(t *testing.T) {
	renderable := NewEvidence("EV-005", "claim").
		WithSource("https://example.com", EvidenceSystemWiki, claims.ExternalCommunity, claims.ReliabilityHigh).
		WithExcerpt("public excerpt").
		WithSensitivity(SensitivityPublic)

	if excerpt, ok := renderable.RenderableExcerpt(); !ok || excerpt != "public excerpt" {
		t.Errorf("RenderableExcerpt() = (%q, %v), want (\"public excerpt\", true)", excerpt, ok)
	}

	restricted := NewEvidence("EV-006", "claim").
		WithSource("https://example.com", EvidenceSystemContract, claims.ExternalCommunity, claims.ReliabilityHigh).
		WithExcerpt("restricted excerpt").
		WithSensitivity(SensitivityRestricted)

	if excerpt, ok := restricted.RenderableExcerpt(); ok {
		t.Errorf("RenderableExcerpt() = (%q, %v), want renderable=false for restricted evidence", excerpt, ok)
	}

	noExcerpt := NewEvidence("EV-007", "claim").WithSensitivity(SensitivityPublic)
	if excerpt, ok := noExcerpt.RenderableExcerpt(); ok || excerpt != "" {
		t.Errorf("RenderableExcerpt() = (%q, %v), want (\"\", false) when no excerpt captured", excerpt, ok)
	}
}

func TestDefaultValidityWindow(t *testing.T) {
	if w := DefaultValidityWindow(EvidenceSystemContract); w != 0 {
		t.Errorf("contract window = %v, want 0 (never expires)", w)
	}
	if w := DefaultValidityWindow(EvidenceSystemAnalytics); w != 90*24*time.Hour {
		t.Errorf("analytics window = %v, want 90 days", w)
	}
	if w := DefaultValidityWindow(EvidenceSystemGitHub); w != 180*24*time.Hour {
		t.Errorf("github window = %v, want 180 days", w)
	}
}

func TestIsStale(t *testing.T) {
	now := time.Date(2026, 8, 19, 0, 0, 0, 0, time.UTC)

	t.Run("zero window never expires", func(t *testing.T) {
		e := NewEvidence("EV-010", "claim").
			WithSource("uri", EvidenceSystemContract, claims.ExternalCommunity, claims.ReliabilityHigh).
			WithCapturedAt(now.AddDate(-5, 0, 0))
		if e.IsStale(now, 0) {
			t.Error("expected not stale with zero window")
		}
	})

	t.Run("no capture time is always stale", func(t *testing.T) {
		e := NewEvidence("EV-011", "claim").
			WithSource("uri", EvidenceSystemAnalytics, claims.ExternalCommunity, claims.ReliabilityHigh)
		if !e.IsStale(now, 90*24*time.Hour) {
			t.Error("expected stale when no capture time recorded")
		}
	})

	t.Run("within window is fresh", func(t *testing.T) {
		e := NewEvidence("EV-012", "claim").
			WithSource("uri", EvidenceSystemAnalytics, claims.ExternalCommunity, claims.ReliabilityHigh).
			WithCapturedAt(now.AddDate(0, 0, -30))
		if e.IsStale(now, 90*24*time.Hour) {
			t.Error("expected fresh within window")
		}
	})

	t.Run("beyond window is stale", func(t *testing.T) {
		e := NewEvidence("EV-013", "claim").
			WithSource("uri", EvidenceSystemAnalytics, claims.ExternalCommunity, claims.ReliabilityHigh).
			WithCapturedAt(now.AddDate(0, 0, -120))
		if !e.IsStale(now, 90*24*time.Hour) {
			t.Error("expected stale beyond window")
		}
	})
}

func TestIsVerified(t *testing.T) {
	verified := NewEvidence("EV-020", "claim").
		WithSource("uri", EvidenceSystemGitHub, claims.ExternalFramework, claims.ReliabilityAuthoritative)
	if !verified.IsVerified() {
		t.Error("expected verified evidence to report IsVerified() = true")
	}

	unverified := NewEvidence("EV-021", "claim")
	if unverified.IsVerified() {
		t.Error("expected unsourced evidence to report IsVerified() = false")
	}
}

func TestEvidenceValidate(t *testing.T) {
	if err := (&Evidence{}).Validate(); err == nil {
		t.Error("expected error for missing ID")
	}
	if err := (&Evidence{ID: "EV-030"}).Validate(); err == nil {
		t.Error("expected error for missing claim text")
	}
	valid := NewEvidence("EV-031", "claim text")
	if err := valid.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
