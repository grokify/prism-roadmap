// Package mel provides types and utilities for monitoring, evaluation, and
// learning (MEL) plans.
//
// A MEL plan says how a team will know whether its work is producing the
// results it intends, and what it will do with what it learns. It has seven
// parts:
//
//   - Results: the results framework, from outputs up to impact, each result
//     naming the higher results it contributes to.
//   - Activities: the work that produces the results, with owners and dates.
//   - Indicators: one reference sheet per indicator, saying what it measures
//     (a result, an activity, or an assumption)
//     and exactly how: definition, unit, disaggregation, baseline, targets,
//     data source, collection method, frequency, the person responsible, data
//     quality checks, and known limitations.
//   - Assumptions: the conditions the results framework depends on, and how
//     each is monitored.
//   - Evaluations: the planned evaluations, each with its questions, type,
//     timing, and method.
//   - Learning agenda: the questions the team most needs answered, how and
//     when each will be answered, and the results or assumptions it bears on.
//   - Reviews and reports: when the team stops to review what it has learned,
//     and what each audience receives and how often.
//
// Monitoring tracks the indicators continuously; evaluation asks, at set
// times, whether and why results happened; learning turns both into decisions.
// Elements refer to each other by ID, and IDs are unique across the plan.
package mel

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Document status constants.
const (
	StatusDraft     = "Draft"
	StatusInReview  = "In Review"
	StatusApproved  = "Approved"
	StatusActive    = "Active"
	StatusCompleted = "Completed"
	StatusArchived  = "Archived"
)

// Result level constants, from the most immediate to the most distant.
const (
	LevelOutput     = "output"      // a direct product of the work
	LevelShortTerm  = "short_term"  // an early change, such as in knowledge or access
	LevelMediumTerm = "medium_term" // a change in behavior or practice
	LevelLongTerm   = "long_term"   // a lasting change in condition
	LevelImpact     = "impact"      // the ultimate change the work contributes to
)

// Indicator type constants.
const (
	TypeQuantitative = "quantitative"
	TypeQualitative  = "qualitative"
)

// Direction constants say which way an indicator should move.
const (
	DirectionIncrease = "increase"
	DirectionDecrease = "decrease"
	DirectionMaintain = "maintain"
)

// Evaluation type constants.
const (
	EvalBaseline      = "baseline"      // establishes the starting point
	EvalProcess       = "process"       // whether the work is delivered as planned
	EvalFormative     = "formative"     // improves the work while it runs
	EvalMidterm       = "midterm"       // takes stock partway through
	EvalFinal         = "final"         // judges results at the end (summative)
	EvalImpact        = "impact"        // attributes change to the work, against a counterfactual
	EvalDevelopmental = "developmental" // supports work that is still being designed
)

// Evaluator constants say who carries out an evaluation.
const (
	EvaluatorInternal = "internal"
	EvaluatorExternal = "external"
	EvaluatorMixed    = "mixed"
)

// Evaluation status constants.
const (
	EvaluationPlanned    = "Planned"
	EvaluationInProgress = "In Progress"
	EvaluationCompleted  = "Completed"
	EvaluationCancelled  = "Cancelled"
)

// Activity status constants.
const (
	ActivityPlanned    = "Planned"
	ActivityInProgress = "In Progress"
	ActivityPaused     = "Paused"
	ActivityCompleted  = "Completed"
	ActivityCancelled  = "Cancelled"
)

// Learning question status constants.
const (
	QuestionOpen     = "Open"
	QuestionAnswered = "Answered"
	QuestionRetired  = "Retired"
)

// Plan is a complete MEL plan.
type Plan struct {
	Schema   string    `json:"$schema,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
	// Purpose says what the plan is for and whom it serves.
	Purpose     string             `json:"purpose,omitempty"`
	Results     []Result           `json:"results"`
	Activities  []Activity         `json:"activities,omitempty"`
	Indicators  []Indicator        `json:"indicators,omitempty"`
	Assumptions []Assumption       `json:"assumptions,omitempty"`
	Evaluations []Evaluation       `json:"evaluations,omitempty"`
	Learning    []LearningQuestion `json:"learningAgenda,omitempty"`
	Reviews     []Review           `json:"reviews,omitempty"`
	Reports     []Report           `json:"reports,omitempty"`
	// Budget is what MEL costs, as written, e.g. "7% of the program budget".
	Budget string `json:"budget,omitempty"`
}

// Metadata contains document metadata.
type Metadata struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name,omitempty"`
	Author  string `json:"author,omitempty"`
	Team    string `json:"team,omitempty"`
	Period  string `json:"period,omitempty"` // e.g. "2026-2028"
	Version string `json:"version,omitempty"`
	Status  string `json:"status,omitempty"`
	// ParentID is the MEL plan this one sits under, e.g. a program plan under
	// an organization's.
	ParentID  string    `json:"parentId,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitzero"`
	UpdatedAt time.Time `json:"updatedAt,omitzero"`
}

// Result is one result in the results framework.
type Result struct {
	ID          string `json:"id"`
	Statement   string `json:"statement"`
	Level       string `json:"level"`
	Description string `json:"description,omitempty"`
	// ContributesTo are the IDs of the higher results this one leads to.
	ContributesTo []string `json:"contributesTo,omitempty"`
}

// Activity is work that produces results.
type Activity struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	// ResultIDs are the results the activity is meant to produce.
	ResultIDs []string `json:"resultIds,omitempty"`
	Owner     string   `json:"owner,omitempty"`
	StartDate string   `json:"startDate,omitempty"` // ISO 8601 date
	EndDate   string   `json:"endDate,omitempty"`   // ISO 8601 date
	Status    string   `json:"status,omitempty"`
}

// Indicator is an indicator reference sheet.
type Indicator struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	// Measures is the ID of the result, activity, or assumption the indicator
	// measures. An activity's indicators track delivery; an assumption's
	// indicators watch whether it still holds.
	Measures   string `json:"measures"`
	Definition string `json:"definition,omitempty"`
	Type       string `json:"type,omitempty"` // quantitative, qualitative
	Unit       string `json:"unit,omitempty"`
	Direction  string `json:"direction,omitempty"`
	// Disaggregation lists how values are broken down, e.g. by sex or region.
	Disaggregation []string  `json:"disaggregation,omitempty"`
	Baseline       *Baseline `json:"baseline,omitempty"`
	Targets        []Target  `json:"targets,omitempty"`
	// DataSource is where the data comes from; CollectionMethod how it is
	// gathered, e.g. "post-session survey" or "export from the CRM".
	DataSource       string `json:"dataSource,omitempty"`
	CollectionMethod string `json:"collectionMethod,omitempty"`
	// Frequency is how often it is collected, e.g. "Monthly".
	Frequency string `json:"frequency,omitempty"`
	// Owner is the person or role responsible for collecting and reporting it.
	Owner       string       `json:"owner,omitempty"`
	DataQuality *DataQuality `json:"dataQuality,omitempty"`
	// Current is the latest value, when the plan is used to track progress.
	Current string `json:"current,omitempty"`
}

// Baseline is an indicator's starting value.
type Baseline struct {
	Value  string `json:"value"`
	Date   string `json:"date,omitempty"` // ISO 8601 date
	Source string `json:"source,omitempty"`
}

// Target is a value an indicator should reach by a date.
type Target struct {
	Value string `json:"value"`
	Date  string `json:"date,omitempty"` // ISO 8601 date
	Label string `json:"label,omitempty"`
}

// DataQuality records how an indicator's data is checked and what it cannot
// show.
type DataQuality struct {
	Checks       []string `json:"checks,omitempty"`
	Limitations  []string `json:"limitations,omitempty"`
	LastAssessed string   `json:"lastAssessed,omitempty"` // ISO 8601 date
}

// Assumption is a condition the results framework depends on, and how it is
// watched.
type Assumption struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
	// ResultIDs are the results that depend on it.
	ResultIDs []string `json:"resultIds,omitempty"`
	// Monitoring says how the team will know whether it still holds.
	Monitoring   string   `json:"monitoring,omitempty"`
	IndicatorIDs []string `json:"indicatorIds,omitempty"`
}

// Evaluation is a planned evaluation.
type Evaluation struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Type      string   `json:"type,omitempty"`
	Questions []string `json:"questions,omitempty"`
	// Method is the design and methods, e.g. "pre-post survey with interviews".
	Method    string `json:"method,omitempty"`
	Evaluator string `json:"evaluator,omitempty"` // internal, external, mixed
	Owner     string `json:"owner,omitempty"`
	StartDate string `json:"startDate,omitempty"` // ISO 8601 date
	EndDate   string `json:"endDate,omitempty"`   // ISO 8601 date
	Status    string `json:"status,omitempty"`
	Budget    string `json:"budget,omitempty"`
	// RelatedIDs are the results, activities, indicators, assumptions, or
	// learning questions the evaluation addresses.
	RelatedIDs []string `json:"relatedIds,omitempty"`
}

// LearningQuestion is one question on the learning agenda.
type LearningQuestion struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	// Rationale says why the answer matters: the decision it informs.
	Rationale string `json:"rationale,omitempty"`
	// Methods say how it will be answered, e.g. "cohort interviews".
	Methods []string `json:"methods,omitempty"`
	Owner   string   `json:"owner,omitempty"`
	DueDate string   `json:"dueDate,omitempty"` // ISO 8601 date
	Status  string   `json:"status,omitempty"`
	Answer  string   `json:"answer,omitempty"`
	// RelatedIDs are the results, activities, indicators, or assumptions it
	// bears on.
	RelatedIDs []string `json:"relatedIds,omitempty"`
}

// Review is a recurring point where the team reviews what it has learned.
type Review struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Purpose      string   `json:"purpose,omitempty"`
	Frequency    string   `json:"frequency,omitempty"`
	Participants []string `json:"participants,omitempty"`
	Owner        string   `json:"owner,omitempty"`
	NextDate     string   `json:"nextDate,omitempty"` // ISO 8601 date
	// QuestionIDs are the learning questions the review takes up.
	QuestionIDs []string `json:"questionIds,omitempty"`
}

// Report is what one audience receives, and how often.
type Report struct {
	ID       string `json:"id"`
	Audience string `json:"audience"`
	// Needs is what the audience needs to know, and the decisions it makes.
	Needs     string `json:"needs,omitempty"`
	Product   string `json:"product,omitempty"` // e.g. "Quarterly dashboard"
	Frequency string `json:"frequency,omitempty"`
	Channel   string `json:"channel,omitempty"`
	Owner     string `json:"owner,omitempty"`
	// IndicatorIDs are the indicators the report shows.
	IndicatorIDs []string `json:"indicatorIds,omitempty"`
}

// New returns a plan with the given ID and name.
func New(id, name string) *Plan {
	return &Plan{Metadata: &Metadata{ID: id, Name: name, Status: StatusDraft}}
}

// ReadFile reads a plan from a JSON file.
func ReadFile(filepath string) (*Plan, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	return Parse(data)
}

// Parse parses plan JSON data.
func Parse(data []byte) (*Plan, error) {
	var p Plan
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("parsing JSON: %w", err)
	}
	return &p, nil
}

// JSON returns the plan as formatted JSON.
func (p *Plan) JSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

// WriteFile writes the plan to a JSON file.
func (p *Plan) WriteFile(filepath string) error {
	data, err := p.JSON()
	if err != nil {
		return fmt.Errorf("marshaling JSON: %w", err)
	}
	if err := os.WriteFile(filepath, data, 0600); err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

// Result returns the result with the given ID, or nil.
func (p *Plan) Result(id string) *Result {
	for i := range p.Results {
		if p.Results[i].ID == id {
			return &p.Results[i]
		}
	}
	return nil
}

// Indicator returns the indicator with the given ID, or nil.
func (p *Plan) Indicator(id string) *Indicator {
	for i := range p.Indicators {
		if p.Indicators[i].ID == id {
			return &p.Indicators[i]
		}
	}
	return nil
}

// Activity returns the activity with the given ID, or nil.
func (p *Plan) Activity(id string) *Activity {
	for i := range p.Activities {
		if p.Activities[i].ID == id {
			return &p.Activities[i]
		}
	}
	return nil
}

// ActivitiesFor returns the activities meant to produce the result, in
// document order.
func (p *Plan) ActivitiesFor(resultID string) []Activity {
	var out []Activity
	for _, a := range p.Activities {
		for _, id := range a.ResultIDs {
			if id == resultID {
				out = append(out, a)
				break
			}
		}
	}
	return out
}

// IndicatorsFor returns the indicators that measure the result, activity, or
// assumption with the given ID, in document order.
func (p *Plan) IndicatorsFor(id string) []Indicator {
	var out []Indicator
	for _, ind := range p.Indicators {
		if ind.Measures == id {
			out = append(out, ind)
		}
	}
	return out
}

// Contributors returns the results that contribute to the result, in document
// order.
func (p *Plan) Contributors(resultID string) []Result {
	var out []Result
	for _, r := range p.Results {
		for _, id := range r.ContributesTo {
			if id == resultID {
				out = append(out, r)
				break
			}
		}
	}
	return out
}

// TopResults returns the results that contribute to no other result, in
// document order: the tops of the results framework.
func (p *Plan) TopResults() []Result {
	var out []Result
	for _, r := range p.Results {
		if len(r.ContributesTo) == 0 {
			out = append(out, r)
		}
	}
	return out
}

// Levels returns the result levels, from the most immediate to the most
// distant.
func Levels() []string {
	return []string{LevelOutput, LevelShortTerm, LevelMediumTerm, LevelLongTerm, LevelImpact}
}

// LevelRank returns a level's position in Levels, or -1 for an unknown level.
func LevelRank(level string) int {
	for i, l := range Levels() {
		if l == level {
			return i
		}
	}
	return -1
}

// LevelLabel returns a level's display label, e.g. "Medium-term outcome".
func LevelLabel(level string) string {
	switch level {
	case LevelOutput:
		return "Output"
	case LevelShortTerm:
		return "Short-term outcome"
	case LevelMediumTerm:
		return "Medium-term outcome"
	case LevelLongTerm:
		return "Long-term outcome"
	case LevelImpact:
		return "Impact"
	}
	return level
}

// DocumentStatuses returns the document statuses.
func DocumentStatuses() []string {
	return []string{StatusDraft, StatusInReview, StatusApproved, StatusActive, StatusCompleted, StatusArchived}
}

// IndicatorTypes returns the indicator types.
func IndicatorTypes() []string { return []string{TypeQuantitative, TypeQualitative} }

// Directions returns the indicator directions.
func Directions() []string { return []string{DirectionIncrease, DirectionDecrease, DirectionMaintain} }

// EvaluationTypes returns the evaluation types.
func EvaluationTypes() []string {
	return []string{EvalBaseline, EvalProcess, EvalFormative, EvalMidterm, EvalFinal, EvalImpact, EvalDevelopmental}
}

// Evaluators returns who can carry out an evaluation.
func Evaluators() []string { return []string{EvaluatorInternal, EvaluatorExternal, EvaluatorMixed} }

// EvaluationStatuses returns the evaluation statuses.
func EvaluationStatuses() []string {
	return []string{EvaluationPlanned, EvaluationInProgress, EvaluationCompleted, EvaluationCancelled}
}

// ActivityStatuses returns the activity statuses.
func ActivityStatuses() []string {
	return []string{ActivityPlanned, ActivityInProgress, ActivityPaused, ActivityCompleted, ActivityCancelled}
}

// QuestionStatuses returns the learning question statuses.
func QuestionStatuses() []string { return []string{QuestionOpen, QuestionAnswered, QuestionRetired} }
