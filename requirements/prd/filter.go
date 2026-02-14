package prd

// tagMatchFunc is a function that determines if entity tags match filter tags.
type tagMatchFunc func(entityTags, filterTags []string) bool

// FilterByTags returns a new Document containing only entities that have at least
// one of the specified tags (OR logic). If no tags are provided, returns a copy
// of the original document. Metadata tags are always preserved.
func (d Document) FilterByTags(tags ...string) Document {
	return d.filterByTagsInternal(tags, hasAnyTag)
}

// FilterByTagsAll returns a new Document containing only entities that have ALL
// of the specified tags (AND logic). If no tags are provided, returns a copy
// of the original document.
func (d Document) FilterByTagsAll(tags ...string) Document {
	return d.filterByTagsInternal(tags, hasAllTags)
}

// filterByTagsInternal is the common implementation for tag filtering.
func (d Document) filterByTagsInternal(tags []string, matchFunc tagMatchFunc) Document {
	if len(tags) == 0 {
		return d
	}

	filtered := Document{
		Metadata:         d.Metadata,
		ExecutiveSummary: d.ExecutiveSummary,
		Objectives:       Objectives{},
	}

	// Filter Personas
	filtered.Personas = filterSliceByTagsFunc(d.Personas, tags, func(p Persona) []string { return p.Tags }, matchFunc)

	// Filter UserStories
	filtered.UserStories = filterSliceByTagsFunc(d.UserStories, tags, func(us UserStory) []string { return us.Tags }, matchFunc)

	// Filter Requirements
	filtered.Requirements = Requirements{
		Functional:    filterSliceByTagsFunc(d.Requirements.Functional, tags, func(fr FunctionalRequirement) []string { return fr.Tags }, matchFunc),
		NonFunctional: filterSliceByTagsFunc(d.Requirements.NonFunctional, tags, func(nfr NonFunctionalRequirement) []string { return nfr.Tags }, matchFunc),
		Compliance:    filterSliceByTagsFunc(d.Requirements.Compliance, tags, func(cr ComplianceRequirement) []string { return cr.Tags }, matchFunc),
	}

	// Filter Roadmap
	filteredPhases := make([]Phase, 0, len(d.Roadmap.Phases))
	for _, phase := range d.Roadmap.Phases {
		if matchFunc(phase.Tags, tags) {
			filteredPhases = append(filteredPhases, phase)
		} else {
			// Check if any deliverables match
			filteredDeliverables := filterSliceByTagsFunc(phase.Deliverables, tags, func(del Deliverable) []string { return del.Tags }, matchFunc)
			if len(filteredDeliverables) > 0 {
				phaseCopy := phase
				phaseCopy.Deliverables = filteredDeliverables
				filteredPhases = append(filteredPhases, phaseCopy)
			}
		}
	}
	filtered.Roadmap = Roadmap{Phases: filteredPhases}

	// Filter Objectives.OKRs
	filteredOKRs := filterOKRsByTagsFunc(d.Objectives.OKRs, tags, matchFunc)
	if len(filteredOKRs) > 0 {
		filtered.Objectives.OKRs = filteredOKRs
	}

	// Filter optional sections (no tags on these, keep all)
	filtered.Assumptions = d.Assumptions
	filtered.OutOfScope = d.OutOfScope
	filtered.TechArchitecture = d.TechArchitecture
	filtered.UXRequirements = d.UXRequirements
	filtered.Glossary = d.Glossary
	filtered.CustomSections = d.CustomSections

	// Filter Risks
	filtered.Risks = filterSliceByTagsFunc(d.Risks, tags, func(r Risk) []string { return r.Tags }, matchFunc)

	// Extended sections (no tags, keep all)
	filtered.Problem = d.Problem
	filtered.Market = d.Market
	filtered.Solution = d.Solution
	filtered.Decisions = d.Decisions
	filtered.Reviews = d.Reviews
	filtered.RevisionHistory = d.RevisionHistory
	filtered.Goals = d.Goals
	filtered.CurrentState = d.CurrentState
	filtered.SecurityModel = d.SecurityModel
	filtered.Appendices = d.Appendices
	filtered.OpenItems = d.OpenItems

	return filtered
}

// hasAnyTag returns true if entityTags contains at least one of filterTags.
func hasAnyTag(entityTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(filterTags))
	for _, t := range filterTags {
		tagSet[t] = struct{}{}
	}
	for _, et := range entityTags {
		if _, ok := tagSet[et]; ok {
			return true
		}
	}
	return false
}

// hasAllTags returns true if entityTags contains ALL of filterTags.
func hasAllTags(entityTags, filterTags []string) bool {
	tagSet := make(map[string]struct{}, len(entityTags))
	for _, t := range entityTags {
		tagSet[t] = struct{}{}
	}
	for _, ft := range filterTags {
		if _, ok := tagSet[ft]; !ok {
			return false
		}
	}
	return true
}

// filterSliceByTagsFunc filters a slice of items using the provided match function.
func filterSliceByTagsFunc[T any](items []T, tags []string, getTags func(T) []string, matchFunc tagMatchFunc) []T {
	if len(items) == 0 {
		return nil
	}
	result := make([]T, 0, len(items))
	for _, item := range items {
		if matchFunc(getTags(item), tags) {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// filterOKRsByTagsFunc filters OKRs using the provided match function.
func filterOKRsByTagsFunc(okrs []OKR, tags []string, matchFunc tagMatchFunc) []OKR {
	if len(okrs) == 0 {
		return nil
	}
	result := make([]OKR, 0, len(okrs))
	for _, okr := range okrs {
		// Check if Objective has matching tags
		if matchFunc(okr.Objective.Tags, tags) {
			result = append(result, okr)
			continue
		}
		// Check if any KeyResult has matching tags
		filteredKRs := filterSliceByTagsFunc(okr.KeyResults, tags, func(kr KeyResult) []string { return kr.Tags }, matchFunc)
		if len(filteredKRs) > 0 {
			okrCopy := OKR{
				Objective:  okr.Objective,
				KeyResults: filteredKRs,
			}
			result = append(result, okrCopy)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return result
}
