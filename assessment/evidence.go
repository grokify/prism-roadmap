// Package assessment provides the Opportunity Assessment IR: an
// evidence-backed, rubric-driven prioritization record composing
// structured-evaluation's rubric/claims types with prism-roadmap's
// MoSCoW/RICE/Kano/Market-Investment-Horizon domain semantics.
//
// See docs/specs/initiatives/INIT-PRISMROADMAP-001/ for the PRD/TRD driving
// this package's design.
package assessment

import (
	"fmt"
	"time"

	"github.com/plexusone/structured-evaluation/claims"
)

// EvidenceSystem identifies the platform hosting an evidence source, for
// humans resolving the URI and for system-specific staleness windows
// (DefaultValidityWindow). Distinct from claims.ExternalSourceType, which
// categorizes a source's general authority rather than where it lives.
type EvidenceSystem string

const (
	EvidenceSystemWiki       EvidenceSystem = "wiki"
	EvidenceSystemConfluence EvidenceSystem = "confluence"
	EvidenceSystemGoogleDocs EvidenceSystem = "google-docs"
	EvidenceSystemGitHub     EvidenceSystem = "github"
	EvidenceSystemGitLab     EvidenceSystem = "gitlab"
	EvidenceSystemSlack      EvidenceSystem = "slack"
	EvidenceSystemJira       EvidenceSystem = "jira"
	EvidenceSystemAnalytics  EvidenceSystem = "analytics-dashboard"
	EvidenceSystemContract   EvidenceSystem = "contract"
	EvidenceSystemOther      EvidenceSystem = "other"
)

// Sensitivity controls whether an evidence excerpt may be rendered verbatim
// into a generated report (prism-roadmap TRD D4/D6). The zero value is
// deliberately not renderable — an evidence record with no sensitivity set
// is treated conservatively rather than assumed safe to publish.
type Sensitivity string

const (
	SensitivityPublic     Sensitivity = "public"
	SensitivityInternal   Sensitivity = "internal"
	SensitivityRestricted Sensitivity = "restricted"
)

// Renderable reports whether an excerpt at this sensitivity level may be
// inlined into a report body. Restricted evidence still exists in the IR for
// judges/reviewers; a renderer must cite the evidence ID with an access note
// instead of the quoted text.
func (s Sensitivity) Renderable() bool {
	return s == SensitivityPublic || s == SensitivityInternal
}

// Evidence is a first-class, independently referenceable record supporting
// one or more rubric answers across one or more assessments. It wraps
// structured-evaluation's claims.Claim rather than reinventing claim/verdict
// semantics:
//
//   - Claim.Text is the specific assertion this evidence supports.
//   - Claim.Category classifies the assertion (statistical, risk-assessment, etc).
//   - Claim.Verdict + Claim.Validation carry the verification state.
//   - Claim.Validation.External (URL, QuotedText, Reliability, AccessedAt)
//     carries the source pointer, the captured excerpt, the reliability
//     tier, and when it was captured.
//
// Reverse lookup ("which assessments cite this evidence") is intentionally
// NOT stored on Evidence — it is a query the persistence layer (omniroadmap)
// answers by indexing assessment→evidence references, keeping this struct
// the single normative record rather than something every citing assessment
// must keep in sync.
type Evidence struct {
	// ID is the evidence identifier, conventionally EV-NNN, unique within an
	// assessment corpus (an omniroadmap-assigned sequence).
	ID string `json:"id"`

	// Claim is the underlying structured-evaluation claim: text, category,
	// verdict, and validation (including the external source, quoted
	// excerpt, reliability tier, and access timestamp).
	Claim claims.Claim `json:"claim"`

	// System identifies the platform hosting the source.
	System EvidenceSystem `json:"system,omitempty"`

	// Sensitivity gates whether the excerpt may render verbatim in a report.
	Sensitivity Sensitivity `json:"sensitivity,omitempty"`

	// CapturedBy identifies who or what captured this evidence — a judge
	// model identifier, or a person's handle. Distinct from
	// Claim.Validation.External.AccessedAt, which records when.
	CapturedBy string `json:"capturedBy,omitempty"`
}

// NewEvidence creates an Evidence record for the given claim text, defaulting
// to claims.ClaimTechnicalFinding. Use WithCategory to override, and
// WithSource/WithExcerpt/WithCapturedAt/WithCapturedBy/WithSensitivity to
// attach the rest (RMI-PRISMROADMAP-001, prism-roadmap TRD D4).
func NewEvidence(id, claimText string) *Evidence {
	return &Evidence{
		ID:    id,
		Claim: *claims.NewClaim(id, claimText, claims.ClaimTechnicalFinding, claims.Location{}),
	}
}

// WithCategory overrides the claim category (default: ClaimTechnicalFinding).
func (e *Evidence) WithCategory(c claims.ClaimCategory) *Evidence {
	e.Claim.Category = c
	return e
}

// WithSource attaches an external source and (re-)computes the claim's
// verdict from it via claims.DetermineVerdict.
func (e *Evidence) WithSource(uri string, system EvidenceSystem, sourceType claims.ExternalSourceType, reliability claims.ReliabilityTier) *Evidence {
	e.System = system
	e.Claim.Validation = &claims.Validation{
		Type: claims.SourceExternal,
		External: &claims.ExternalValidation{
			URL:         uri,
			SourceType:  sourceType,
			Reliability: reliability,
		},
	}
	e.Claim.Verdict = claims.DetermineVerdict(e.Claim.Validation)
	return e
}

// WithExcerpt sets the captured excerpt — the specific quoted text a judge
// evaluated, kept short (a sentence or two, not a copy of the source) so it
// survives link rot and lets a reviewer sanity-check the citation without
// opening the source. No-op if WithSource hasn't been called yet.
func (e *Evidence) WithExcerpt(text string) *Evidence {
	if ext := e.external(); ext != nil {
		ext.QuotedText = text
	}
	return e
}

// WithCapturedAt records when this evidence was captured, driving staleness
// checks (IsStale). No-op if WithSource hasn't been called yet.
func (e *Evidence) WithCapturedAt(t time.Time) *Evidence {
	if ext := e.external(); ext != nil {
		ext.AccessedAt = &t
	}
	return e
}

// WithCapturedBy records who or what captured this evidence.
func (e *Evidence) WithCapturedBy(who string) *Evidence {
	e.CapturedBy = who
	return e
}

// WithSensitivity sets the rendering-sensitivity level.
func (e *Evidence) WithSensitivity(s Sensitivity) *Evidence {
	e.Sensitivity = s
	return e
}

// external returns the external validation record, or nil if none is set.
func (e *Evidence) external() *claims.ExternalValidation {
	if e.Claim.Validation == nil {
		return nil
	}
	return e.Claim.Validation.External
}

// SourceURI returns the source URI, or "" if none is set.
func (e *Evidence) SourceURI() string {
	if ext := e.external(); ext != nil {
		return ext.URL
	}
	return ""
}

// Excerpt returns the captured excerpt, or "" if none is set.
func (e *Evidence) Excerpt() string {
	if ext := e.external(); ext != nil {
		return ext.QuotedText
	}
	return ""
}

// CapturedAtTime returns when this evidence was captured, or nil if unset.
func (e *Evidence) CapturedAtTime() *time.Time {
	if ext := e.external(); ext != nil {
		return ext.AccessedAt
	}
	return nil
}

// RenderableExcerpt returns the captured excerpt and whether it may be
// rendered verbatim in a report body given this evidence's sensitivity.
// When not renderable, callers should cite the evidence ID with an access
// note instead (prism-roadmap TRD D6).
func (e *Evidence) RenderableExcerpt() (excerpt string, renderable bool) {
	excerpt = e.Excerpt()
	return excerpt, excerpt != "" && e.Sensitivity.Renderable()
}

// IsVerified reports whether the underlying claim's verdict is passing.
func (e *Evidence) IsVerified() bool {
	return e.Claim.Verdict.IsPassing()
}

// DefaultValidityWindow returns how long evidence from a given system is
// considered fresh before an omniroadmap staleness sweep should flag
// assessments citing it (prism-roadmap PRD FR11). A zero duration means the
// evidence does not expire — a signed contract doesn't go stale the way a
// deployment inventory does.
func DefaultValidityWindow(system EvidenceSystem) time.Duration {
	switch system {
	case EvidenceSystemContract:
		return 0
	case EvidenceSystemAnalytics:
		return 90 * 24 * time.Hour // ~1 quarter: deployment/usage data ages fast
	default:
		return 180 * 24 * time.Hour // ~2 quarters: docs, code, tickets
	}
}

// IsStale reports whether this evidence's capture time has exceeded window.
// A zero window means the evidence never expires. Evidence with no recorded
// capture time is always considered stale — a staleness sweep can't vouch
// for freshness it never measured.
func (e *Evidence) IsStale(now time.Time, window time.Duration) bool {
	if window <= 0 {
		return false
	}
	capturedAt := e.CapturedAtTime()
	if capturedAt == nil {
		return true
	}
	return now.Sub(*capturedAt) > window
}

// Validate returns an error if required fields are missing.
func (e *Evidence) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("id is required")
	}
	if e.Claim.Text == "" {
		return fmt.Errorf("claim text is required")
	}
	return nil
}
