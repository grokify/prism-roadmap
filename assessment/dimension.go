package assessment

import "fmt"

// DimensionKind distinguishes mutually-exclusive Category dimensions (0..1
// selection per opportunity, pie-chartable) from multi-select Tag
// dimensions (0..N, bar/coverage — selections can sum to more than 100%).
type DimensionKind string

const (
	DimensionKindCategory DimensionKind = "category"
	DimensionKindTags     DimensionKind = "tags"
)

// DimensionQuestion is one Y/N judge question characterizing a
// DimensionOption (e.g. Kano's "Would customers reasonably consider this a
// basic expectation?"). Each option's questions are evaluated
// independently and cited on their own; a dimension resolves its selection
// from the resulting PATTERN of answers (see ResolveCategory/ResolveTags),
// not from a single ordered "highest wins" ladder — Kano's categories in
// particular are not an ordered hierarchy the way MoSCoW/Impact/Confidence
// are (ideation doc).
type DimensionQuestion struct {
	ID       string `json:"id"`
	Question string `json:"question"`
}

// DimensionOption is one selectable value within a portfolio dimension
// (e.g. Kano's "Must-be", Market Investment Horizon's "KTLO"), with the
// judge questions used to determine whether an opportunity belongs to it.
type DimensionOption struct {
	ID        string              `json:"id"`
	Label     string              `json:"label"`
	Questions []DimensionQuestion `json:"questions,omitempty"`
}

// DimensionAnswer is a judge's Y/N answer to one option's question, with
// evidence — mirrors ThresholdAnswer's evidence discipline: an answer of
// true without an evidence citation does not count toward selecting the
// option (see satisfiedOptions).
type DimensionAnswer struct {
	OptionID    string   `json:"optionId"`
	QuestionID  string   `json:"questionId"`
	Answer      bool     `json:"answer"`
	Rationale   string   `json:"rationale,omitempty"`
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// DimensionDefinition is a versioned portfolio dimension: what an
// opportunity can be classified as (Kind == category) or tagged with (Kind
// == tags). Built-in dimensions (Kano, Market Investment Horizon —
// RMI-PRISMROADMAP-006) and custom organization-defined dimensions (e.g. a
// "2026 Strategic Priority: AI/Growth/Excellence" category) use the same
// shape — a new dimension needs no schema change (prism-roadmap PRD FR4).
//
// Assessments reference a DimensionDefinition by ID+Version
// (DimensionAssignment) rather than copying it, so a definition changing
// later doesn't retroactively alter what a past assignment meant.
type DimensionDefinition struct {
	ID      string        `json:"id"`
	Name    string        `json:"name"`
	Version string        `json:"version"`
	Kind    DimensionKind `json:"kind"`

	Options []DimensionOption `json:"options"`
}

// OptionByID returns an option by ID, or nil if not found.
func (d *DimensionDefinition) OptionByID(id string) *DimensionOption {
	for i := range d.Options {
		if d.Options[i].ID == id {
			return &d.Options[i]
		}
	}
	return nil
}

// Validate returns an error if the definition is not usable.
func (d *DimensionDefinition) Validate() error {
	if d.ID == "" {
		return fmt.Errorf("id is required")
	}
	if d.Name == "" {
		return fmt.Errorf("name is required")
	}
	if d.Version == "" {
		return fmt.Errorf("version is required")
	}
	if d.Kind != DimensionKindCategory && d.Kind != DimensionKindTags {
		return fmt.Errorf("kind must be %q or %q, got %q", DimensionKindCategory, DimensionKindTags, d.Kind)
	}
	if len(d.Options) == 0 {
		return fmt.Errorf("at least one option is required")
	}
	seen := make(map[string]bool, len(d.Options))
	for _, opt := range d.Options {
		if opt.ID == "" {
			return fmt.Errorf("option with empty ID")
		}
		if seen[opt.ID] {
			return fmt.Errorf("duplicate option ID %q", opt.ID)
		}
		seen[opt.ID] = true
	}
	return nil
}

// satisfiedOptions returns the IDs of options with at least one
// evidence-backed true answer, in definition order.
func (d *DimensionDefinition) satisfiedOptions(answers []DimensionAnswer) []string {
	byOption := make(map[string]bool, len(answers))
	for _, a := range answers {
		if a.Answer && len(a.EvidenceIDs) > 0 {
			byOption[a.OptionID] = true
		}
	}
	var ids []string
	for _, opt := range d.Options {
		if byOption[opt.ID] {
			ids = append(ids, opt.ID)
		}
	}
	return ids
}

// CategorySelection is the outcome of resolving a Category-kind
// dimension's answers into a single selected option.
type CategorySelection struct {
	// OptionID is set when Resolved is true.
	OptionID string `json:"optionId,omitempty"`

	// Resolved is true when exactly one option was satisfied with
	// evidence.
	Resolved bool `json:"resolved"`

	// Ambiguous is true when more than one option's questions were
	// satisfied. A Category dimension's options are supposed to be
	// mutually exclusive by design, so this is a signal the rubric or
	// judge output needs review — not something to silently resolve by
	// picking the first match.
	Ambiguous bool `json:"ambiguous,omitempty"`

	// AmbiguousOptionIDs lists every option satisfied when Ambiguous is
	// true.
	AmbiguousOptionIDs []string `json:"ambiguousOptionIds,omitempty"`
}

// ResolveCategory determines which single option (if any) a Category-kind
// dimension's answers satisfy: an option is satisfied when at least one of
// its questions is answered true with evidence. More than one satisfied
// option is reported as Ambiguous rather than resolved arbitrarily.
func (d *DimensionDefinition) ResolveCategory(answers []DimensionAnswer) CategorySelection {
	satisfied := d.satisfiedOptions(answers)
	switch len(satisfied) {
	case 0:
		return CategorySelection{}
	case 1:
		return CategorySelection{OptionID: satisfied[0], Resolved: true}
	default:
		return CategorySelection{Ambiguous: true, AmbiguousOptionIDs: satisfied}
	}
}

// ResolveTags returns every option ID satisfied by the given answers for a
// Tags-kind dimension. Unlike ResolveCategory, more than one satisfied
// option is the normal case, not an ambiguity.
func (d *DimensionDefinition) ResolveTags(answers []DimensionAnswer) []string {
	return d.satisfiedOptions(answers)
}

// DimensionAssignment is one opportunity's classification against a single
// DimensionDefinition: which version was used, the raw judge answers (for
// audit), and the resolved selection.
type DimensionAssignment struct {
	DimensionID      string `json:"dimensionId"`
	DimensionVersion string `json:"dimensionVersion"`

	Answers []DimensionAnswer `json:"answers,omitempty"`

	// Category is set when the source dimension's Kind is category.
	Category *CategorySelection `json:"category,omitempty"`

	// Tags is set when the source dimension's Kind is tags.
	Tags []string `json:"tags,omitempty"`
}

// SelectedOptionIDs returns the option IDs this assignment selects: the
// resolved Category option (if any — empty if unresolved or ambiguous), or
// the full Tags list.
func (d DimensionAssignment) SelectedOptionIDs() []string {
	if d.Category != nil && d.Category.Resolved {
		return []string{d.Category.OptionID}
	}
	return d.Tags
}

// NewDimensionAssignment resolves answers against def and returns the
// assignment appropriate to def.Kind.
func NewDimensionAssignment(def *DimensionDefinition, answers []DimensionAnswer) DimensionAssignment {
	a := DimensionAssignment{
		DimensionID:      def.ID,
		DimensionVersion: def.Version,
		Answers:          answers,
	}
	if def.Kind == DimensionKindTags {
		a.Tags = def.ResolveTags(answers)
		return a
	}
	sel := def.ResolveCategory(answers)
	a.Category = &sel
	return a
}
