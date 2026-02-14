package prd

import (
	"fmt"
	"regexp"
)

// tagPattern matches valid kebab-case tags:
// - Lowercase alphanumeric (a-z, 0-9)
// - Segments separated by single hyphens
// - No leading, trailing, or consecutive hyphens
// Examples: "mvp", "phase-1", "2024-q1", "backend-api"
var tagPattern = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// ValidateTag checks if a tag follows kebab-case conventions.
// Valid tags are lowercase alphanumeric with single hyphens between segments.
func ValidateTag(tag string) error {
	if tag == "" {
		return fmt.Errorf("tag cannot be empty")
	}
	if !tagPattern.MatchString(tag) {
		return fmt.Errorf("invalid tag %q: must be lowercase alphanumeric with hyphens (e.g., 'my-tag', 'phase-1')", tag)
	}
	return nil
}

// ValidateTags checks multiple tags and returns all validation errors.
func ValidateTags(tags []string) []error {
	var errs []error
	for _, tag := range tags {
		if err := ValidateTag(tag); err != nil {
			errs = append(errs, err)
		}
	}
	return errs
}

// ValidationResult contains validation errors and warnings.
type ValidationResult struct {
	Valid    bool                `json:"valid"`
	Errors   []ValidationError   `json:"errors,omitempty"`
	Warnings []ValidationWarning `json:"warnings,omitempty"`
}

// ValidationError represents a validation failure.
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationWarning represents a non-blocking issue.
type ValidationWarning struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Validate checks the Document for structural and content issues.
func Validate(doc *Document) *ValidationResult {
	result := &ValidationResult{Valid: true}

	// Required metadata fields
	if doc.Metadata.ID == "" {
		result.addError("metadata.id", "Document ID is required")
	}

	if doc.Metadata.Title == "" {
		result.addError("metadata.title", "Title is required")
	} else if len(doc.Metadata.Title) < 5 {
		result.addError("metadata.title", "Title must be at least 5 characters")
	}

	if len(doc.Metadata.Authors) == 0 {
		result.addWarning("metadata.authors", "No authors specified")
	}

	if doc.Metadata.Status == "" {
		result.addError("metadata.status", "Status is required")
	}

	// Executive summary
	if doc.ExecutiveSummary.ProblemStatement == "" {
		result.addWarning("executive_summary.problem_statement", "Problem statement is empty")
	}

	if doc.ExecutiveSummary.ProposedSolution == "" {
		result.addWarning("executive_summary.proposed_solution", "Proposed solution is empty")
	}

	// Objectives (OKRs)
	if len(doc.Objectives.OKRs) == 0 {
		result.addWarning("objectives", "No OKRs defined")
	} else {
		// Check that OKRs have key results
		for i, okr := range doc.Objectives.OKRs {
			if len(okr.KeyResults) == 0 {
				result.addWarning(fmt.Sprintf("objectives.okrs[%d]", i), "OKR has no key results defined")
			}
		}
	}

	// Personas
	if len(doc.Personas) == 0 {
		result.addWarning("personas", "No personas defined")
	}

	// User stories
	if len(doc.UserStories) == 0 {
		result.addWarning("user_stories", "No user stories defined")
	}

	// Validate IDs are unique
	result.validateIDs(doc)

	// Validate traceability
	result.validateTraceability(doc)

	// Validate tags
	result.validateTags(doc)

	return result
}

func (r *ValidationResult) addError(field, message string) {
	r.Valid = false
	r.Errors = append(r.Errors, ValidationError{Field: field, Message: message})
}

func (r *ValidationResult) addWarning(field, message string) {
	r.Warnings = append(r.Warnings, ValidationWarning{Field: field, Message: message})
}

// validateIDs checks for duplicate and malformed IDs.
func (r *ValidationResult) validateIDs(doc *Document) {
	ids := make(map[string]string) // id -> location

	checkID := func(id, location string) {
		if id == "" {
			return
		}
		if existingLoc, exists := ids[id]; exists {
			r.addError(location, fmt.Sprintf("Duplicate ID '%s' (also at %s)", id, existingLoc))
		}
		ids[id] = location
	}

	// Check OKR IDs
	for i, okr := range doc.Objectives.OKRs {
		checkID(okr.Objective.ID, fmt.Sprintf("objectives.okrs[%d].objective.id", i))
		for j, kr := range okr.KeyResults {
			checkID(kr.ID, fmt.Sprintf("objectives.okrs[%d].key_results[%d].id", i, j))
		}
	}

	// Check persona IDs
	for i, p := range doc.Personas {
		checkID(p.ID, fmt.Sprintf("personas[%d].id", i))
	}

	// Check user story IDs
	for i, s := range doc.UserStories {
		checkID(s.ID, fmt.Sprintf("user_stories[%d].id", i))
	}

	// Check functional requirement IDs
	for i, req := range doc.Requirements.Functional {
		checkID(req.ID, fmt.Sprintf("requirements.functional[%d].id", i))
	}

	// Check NFR IDs
	for i, nfr := range doc.Requirements.NonFunctional {
		checkID(nfr.ID, fmt.Sprintf("requirements.non_functional[%d].id", i))
	}

	// Check compliance requirement IDs
	for i, cr := range doc.Requirements.Compliance {
		checkID(cr.ID, fmt.Sprintf("requirements.compliance[%d].id", i))
	}

	// Check roadmap phase IDs
	for i, phase := range doc.Roadmap.Phases {
		checkID(phase.ID, fmt.Sprintf("roadmap.phases[%d].id", i))
	}

	// Check problem ID (if present)
	if doc.Problem != nil && doc.Problem.ID != "" {
		checkID(doc.Problem.ID, "problem.id")
	}

	// Check market alternative IDs
	if doc.Market != nil {
		for i, alt := range doc.Market.Alternatives {
			checkID(alt.ID, fmt.Sprintf("market.alternatives[%d].id", i))
		}
	}

	// Check solution option IDs
	if doc.Solution != nil {
		for i, opt := range doc.Solution.SolutionOptions {
			checkID(opt.ID, fmt.Sprintf("solution.solution_options[%d].id", i))
		}
	}

	// Check decision IDs
	if doc.Decisions != nil {
		for i, d := range doc.Decisions.Records {
			checkID(d.ID, fmt.Sprintf("decisions.records[%d].id", i))
		}
	}
}

// validateTraceability checks cross-references between sections.
func (r *ValidationResult) validateTraceability(doc *Document) {
	// Build sets of defined IDs
	definedIDs := make(map[string]bool)

	// Collect all defined IDs from OKRs
	for _, okr := range doc.Objectives.OKRs {
		if okr.Objective.ID != "" {
			definedIDs[okr.Objective.ID] = true
		}
		for _, kr := range okr.KeyResults {
			if kr.ID != "" {
				definedIDs[kr.ID] = true
			}
		}
	}
	for _, p := range doc.Personas {
		if p.ID != "" {
			definedIDs[p.ID] = true
		}
	}
	for _, s := range doc.UserStories {
		if s.ID != "" {
			definedIDs[s.ID] = true
		}
	}
	for _, phase := range doc.Roadmap.Phases {
		if phase.ID != "" {
			definedIDs[phase.ID] = true
		}
	}
	if doc.Problem != nil && doc.Problem.ID != "" {
		definedIDs[doc.Problem.ID] = true
	}
	if doc.Solution != nil {
		for _, opt := range doc.Solution.SolutionOptions {
			if opt.ID != "" {
				definedIDs[opt.ID] = true
			}
		}
	}

	// Check user story references to personas
	for i, story := range doc.UserStories {
		if story.PersonaID != "" && !definedIDs[story.PersonaID] {
			r.addWarning(
				fmt.Sprintf("user_stories[%d].persona_id", i),
				fmt.Sprintf("Reference to undefined persona: %s", story.PersonaID),
			)
		}
		if story.PhaseID != "" && !definedIDs[story.PhaseID] {
			r.addWarning(
				fmt.Sprintf("user_stories[%d].phase_id", i),
				fmt.Sprintf("Reference to undefined phase: %s", story.PhaseID),
			)
		}
	}

	// Check functional requirement references
	for i, req := range doc.Requirements.Functional {
		for _, storyID := range req.UserStoryIDs {
			if storyID != "" && !definedIDs[storyID] {
				r.addWarning(
					fmt.Sprintf("requirements.functional[%d].user_story_ids", i),
					fmt.Sprintf("Reference to undefined user story: %s", storyID),
				)
			}
		}
		if req.PhaseID != "" && !definedIDs[req.PhaseID] {
			r.addWarning(
				fmt.Sprintf("requirements.functional[%d].phase_id", i),
				fmt.Sprintf("Reference to undefined phase: %s", req.PhaseID),
			)
		}
	}

	// Check solution problem references
	if doc.Solution != nil {
		for i, opt := range doc.Solution.SolutionOptions {
			for _, probID := range opt.ProblemsAddressed {
				if probID != "" && !definedIDs[probID] {
					r.addWarning(
						fmt.Sprintf("solution.solution_options[%d].problems_addressed", i),
						fmt.Sprintf("Reference to undefined problem: %s", probID),
					)
				}
			}
		}

		// Check selected solution exists
		if doc.Solution.SelectedSolutionID != "" {
			found := false
			for _, opt := range doc.Solution.SolutionOptions {
				if opt.ID == doc.Solution.SelectedSolutionID {
					found = true
					break
				}
			}
			if !found {
				r.addError(
					"solution.selected_solution_id",
					fmt.Sprintf("Selected solution '%s' not found in solution options", doc.Solution.SelectedSolutionID),
				)
			}
		}
	}
}

// validateTags checks all tags in the document follow kebab-case conventions.
func (r *ValidationResult) validateTags(doc *Document) {
	checkTags := func(tags []string, location string) {
		for _, tag := range tags {
			if err := ValidateTag(tag); err != nil {
				r.addError(location, err.Error())
			}
		}
	}

	// Metadata tags
	checkTags(doc.Metadata.Tags, "metadata.tags")

	// Persona tags
	for i, p := range doc.Personas {
		checkTags(p.Tags, fmt.Sprintf("personas[%d].tags", i))
	}

	// User story tags
	for i, s := range doc.UserStories {
		checkTags(s.Tags, fmt.Sprintf("user_stories[%d].tags", i))
	}

	// Functional requirement tags
	for i, req := range doc.Requirements.Functional {
		checkTags(req.Tags, fmt.Sprintf("requirements.functional[%d].tags", i))
	}

	// Non-functional requirement tags
	for i, nfr := range doc.Requirements.NonFunctional {
		checkTags(nfr.Tags, fmt.Sprintf("requirements.non_functional[%d].tags", i))
	}

	// Compliance requirement tags
	for i, cr := range doc.Requirements.Compliance {
		checkTags(cr.Tags, fmt.Sprintf("requirements.compliance[%d].tags", i))
	}

	// Roadmap phase and deliverable tags
	for i, phase := range doc.Roadmap.Phases {
		checkTags(phase.Tags, fmt.Sprintf("roadmap.phases[%d].tags", i))
		for j, del := range phase.Deliverables {
			checkTags(del.Tags, fmt.Sprintf("roadmap.phases[%d].deliverables[%d].tags", i, j))
		}
	}

	// OKR tags (objectives and key results)
	for i, okr := range doc.Objectives.OKRs {
		checkTags(okr.Objective.Tags, fmt.Sprintf("objectives.okrs[%d].objective.tags", i))
		for j, kr := range okr.KeyResults {
			checkTags(kr.Tags, fmt.Sprintf("objectives.okrs[%d].key_results[%d].tags", i, j))
		}
	}

	// Risk tags
	for i, risk := range doc.Risks {
		checkTags(risk.Tags, fmt.Sprintf("risks[%d].tags", i))
	}
}
