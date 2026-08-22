package assessment

// ThresholdLevel is one rung of a Ladder: a named level and the citable
// criteria a judge tests against to determine whether an opportunity meets
// it. Criteria are evaluated as an OR-set by default (any one satisfies the
// level) — set AllCriteria to require all of them (AND).
type ThresholdLevel struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Criteria    []string `json:"criteria"`
	AllCriteria bool     `json:"allCriteria,omitempty"`
}

// ThresholdAnswer is a judge's Y/N answer for one ladder level. Satisfied
// alone is never trusted — Evaluate treats a Satisfied answer with no
// EvidenceIDs as unsupported (prism-roadmap PRD FR2: "a threshold cannot
// return YES unless the judge can cite supporting evidence").
type ThresholdAnswer struct {
	LevelID string `json:"levelId"`

	// Satisfied is the judge's raw Y/N answer for this level.
	Satisfied bool `json:"satisfied"`

	// CriterionMet names which of the level's Criteria was matched.
	CriterionMet string `json:"criterionMet,omitempty"`

	Rationale string `json:"rationale,omitempty"`

	// EvidenceIDs cite the assessment.Evidence records supporting this
	// answer (RMI-PRISMROADMAP-001).
	EvidenceIDs []string `json:"evidenceIds,omitempty"`
}

// Ladder is an ordered set of threshold levels, highest-priority first,
// evaluated top-down: the first level whose answer is both Satisfied and
// evidence-backed wins (prism-roadmap PRD FR2 / ideation design: "evaluate
// top-down... stop at first YES" — chosen over asking all levels
// independently and checking for monotonic nesting, which the design this
// package implements superseded for the same reasons: fewer contradictory
// outputs, lower judge cost).
//
// The LLM/judge only ever produces ThresholdAnswers; Evaluate is the
// deterministic code that turns them into a classification — the judge
// never emits the resulting level or its numeric value directly.
type Ladder struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Levels []ThresholdLevel `json:"levels"`
}

// LevelByID returns a level by ID, or nil if not found.
func (l *Ladder) LevelByID(id string) *ThresholdLevel {
	for i := range l.Levels {
		if l.Levels[i].ID == id {
			return &l.Levels[i]
		}
	}
	return nil
}

// Evaluate scans Levels top-down and returns the first level whose answer
// is Satisfied and cites at least one evidence ID. Returns ok=false if no
// level was satisfied with evidence — callers decide the floor semantics
// (e.g. RICE Impact floors at "no impact"; MoSCoW floors at Won't/Not Now).
func (l *Ladder) Evaluate(answers []ThresholdAnswer) (level *ThresholdLevel, answer *ThresholdAnswer, ok bool) {
	byID := indexAnswers(answers)
	for i := range l.Levels {
		lvl := &l.Levels[i]
		a, exists := byID[lvl.ID]
		if !exists || !a.Satisfied || len(a.EvidenceIDs) == 0 {
			continue
		}
		return lvl, &a, true
	}
	return nil, nil, false
}

// UnsupportedAnswers returns answers marked Satisfied but missing evidence
// — these were NOT counted toward the classification and should be sent
// back for re-judgment rather than silently dropped.
func (l *Ladder) UnsupportedAnswers(answers []ThresholdAnswer) []ThresholdAnswer {
	var out []ThresholdAnswer
	for _, a := range answers {
		if a.Satisfied && len(a.EvidenceIDs) == 0 {
			out = append(out, a)
		}
	}
	return out
}

func indexAnswers(answers []ThresholdAnswer) map[string]ThresholdAnswer {
	byID := make(map[string]ThresholdAnswer, len(answers))
	for _, a := range answers {
		byID[a.LevelID] = a
	}
	return byID
}
